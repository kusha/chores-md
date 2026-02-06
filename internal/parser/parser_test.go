package parser

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	t.Run("two_chores", func(t *testing.T) {
		content := `## Kitchen
> 1w

## Bathroom
> 2w
`
		result, err := Parse(content)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Chores) != 2 {
			t.Errorf("got %d chores, want 2", len(result.Chores))
		}
		if result.Chores[0].Name != "Kitchen" {
			t.Errorf("first chore name = %q, want %q", result.Chores[0].Name, "Kitchen")
		}
		if result.Chores[1].Name != "Bathroom" {
			t.Errorf("second chore name = %q, want %q", result.Chores[1].Name, "Bathroom")
		}
	})

	t.Run("three_completions", func(t *testing.T) {
		content := `## Kitchen
> 1w

2026-02-03 Kitchen
2026-02-02 Kitchen
2026-02-01 Kitchen
`
		result, err := Parse(content)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Completions) != 3 {
			t.Errorf("got %d completions, want 3", len(result.Completions))
		}
	})

	t.Run("chore_without_frequency", func(t *testing.T) {
		content := `## Kitchen

Some description but no frequency.
`
		_, err := Parse(content)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "no frequency") {
			t.Errorf("error should mention 'no frequency', got: %v", err)
		}
		if !strings.Contains(err.Error(), "line") {
			t.Errorf("error should mention line number, got: %v", err)
		}
	})

	t.Run("case_insensitive_completion", func(t *testing.T) {
		content := `## Kitchen Clean
> 1w

2026-02-03 kitchen clean
`
		result, err := Parse(content)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Completions) != 1 {
			t.Errorf("got %d completions, want 1", len(result.Completions))
		}
		if result.Completions[0].ChoreName != "kitchen clean" {
			t.Errorf("completion name = %q, want %q", result.Completions[0].ChoreName, "kitchen clean")
		}
	})

	t.Run("empty_file", func(t *testing.T) {
		result, err := Parse("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Chores) != 0 {
			t.Errorf("got %d chores, want 0", len(result.Chores))
		}
		if len(result.Completions) != 0 {
			t.Errorf("got %d completions, want 0", len(result.Completions))
		}
	})

	t.Run("completion_with_comment", func(t *testing.T) {
		content := `## Kitchen
> 1w

2026-02-03 Kitchen # after dinner
`
		result, err := Parse(content)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Completions) != 1 {
			t.Fatalf("got %d completions, want 1", len(result.Completions))
		}
		if result.Completions[0].ChoreName != "Kitchen" {
			t.Errorf("completion name = %q, want %q (comment should be stripped)", result.Completions[0].ChoreName, "Kitchen")
		}
	})

	t.Run("duplicate_chore", func(t *testing.T) {
		content := `## Kitchen
> 1w

## Kitchen
> 2w
`
		result, err := Parse(content)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Chores) != 1 {
			t.Errorf("got %d chores, want 1 (first wins)", len(result.Chores))
		}
		if result.Chores[0].FrequencyDays != 7 {
			t.Errorf("frequency = %d, want 7 (first definition)", result.Chores[0].FrequencyDays)
		}
		if len(result.Warnings) == 0 {
			t.Error("expected warning for duplicate chore")
		}
	})

	t.Run("undefined_chore_completion", func(t *testing.T) {
		content := `## Kitchen
> 1w

2026-02-03 NonExistent Chore
`
		result, err := Parse(content)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Completions) != 1 {
			t.Errorf("got %d completions, want 1 (undefined completions still parsed)", len(result.Completions))
		}
	})

	t.Run("whitespace_variations", func(t *testing.T) {
		tests := []string{
			"## Kitchen\n>2w\n",
			"## Kitchen\n> 2w\n",
			"## Kitchen\n>  2w\n",
		}
		for _, content := range tests {
			result, err := Parse(content)
			if err != nil {
				t.Errorf("Parse(%q) error: %v", content, err)
				continue
			}
			if result.Chores[0].FrequencyDays != 14 {
				t.Errorf("Parse(%q) frequency = %d, want 14", content, result.Chores[0].FrequencyDays)
			}
		}
	})

	t.Run("frequency_raw_preserved", func(t *testing.T) {
		content := `## Kitchen
> 2w
`
		result, err := Parse(content)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Chores[0].FrequencyRaw != "2w" {
			t.Errorf("FrequencyRaw = %q, want %q", result.Chores[0].FrequencyRaw, "2w")
		}
	})

	t.Run("completion_before_definition", func(t *testing.T) {
		content := `# Log
2026-02-03 Kitchen

# Chores
## Kitchen
> 1w
`
		result, err := Parse(content)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Completions) != 1 {
			t.Errorf("got %d completions, want 1", len(result.Completions))
		}
		if len(result.Chores) != 1 {
			t.Errorf("got %d chores, want 1", len(result.Chores))
		}
	})
}

func TestParseFile(t *testing.T) {
	t.Run("valid_file", func(t *testing.T) {
		result, err := ParseFile("testdata/valid.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Chores) != 2 {
			t.Errorf("got %d chores, want 2", len(result.Chores))
		}
		if len(result.Completions) != 3 {
			t.Errorf("got %d completions, want 3", len(result.Completions))
		}
	})

	t.Run("empty_file", func(t *testing.T) {
		result, err := ParseFile("testdata/empty.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Chores) != 0 || len(result.Completions) != 0 {
			t.Errorf("expected empty result for empty file")
		}
	})

	t.Run("no_frequency_error", func(t *testing.T) {
		_, err := ParseFile("testdata/no_frequency.md")
		if err == nil {
			t.Fatal("expected error for chore without frequency")
		}
	})

	t.Run("duplicates_warning", func(t *testing.T) {
		result, err := ParseFile("testdata/duplicates.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Warnings) == 0 {
			t.Error("expected warning for duplicate chore")
		}
	})

	t.Run("comments_stripped", func(t *testing.T) {
		result, err := ParseFile("testdata/comments.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for _, c := range result.Completions {
			if strings.Contains(c.ChoreName, "#") {
				t.Errorf("completion name %q should not contain comment", c.ChoreName)
			}
		}
	})

	t.Run("nonexistent_file", func(t *testing.T) {
		_, err := ParseFile("testdata/does_not_exist.md")
		if err == nil {
			t.Fatal("expected error for nonexistent file")
		}
	})
}
