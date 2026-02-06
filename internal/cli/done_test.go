package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDoneCmd(t *testing.T) {
	baseContent := `## Kitchen Clean
> 1w

## Bathroom Clean
> 2w
`

	t.Run("appends_entry", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "chores.md")
		if err := os.WriteFile(testFile, []byte(baseContent), 0644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		date := time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC)
		var buf bytes.Buffer

		if err := DoneCmd(testFile, "Kitchen Clean", date, &buf); err != nil {
			t.Fatalf("DoneCmd error: %v", err)
		}

		content, _ := os.ReadFile(testFile)
		if !strings.Contains(string(content), "2026-02-10 Kitchen Clean") {
			t.Errorf("file should contain completion entry, got:\n%s", string(content))
		}

		if !strings.Contains(buf.String(), `Done: "Kitchen Clean"`) {
			t.Errorf("output should confirm completion, got: %s", buf.String())
		}
	})

	t.Run("unknown_chore", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "chores.md")
		if err := os.WriteFile(testFile, []byte(baseContent), 0644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		date := time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC)
		var buf bytes.Buffer

		err := DoneCmd(testFile, "Nonexistent Chore", date, &buf)
		if err == nil {
			t.Fatal("expected error for unknown chore")
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("error should mention 'not found', got: %v", err)
		}
	})

	t.Run("case_insensitive", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "chores.md")
		if err := os.WriteFile(testFile, []byte(baseContent), 0644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		date := time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC)
		var buf bytes.Buffer

		if err := DoneCmd(testFile, "kitchen clean", date, &buf); err != nil {
			t.Fatalf("DoneCmd error (case insensitive): %v", err)
		}

		content, _ := os.ReadFile(testFile)
		if !strings.Contains(string(content), "Kitchen Clean") {
			t.Errorf("should use canonical name from definition, got:\n%s", string(content))
		}
	})

	t.Run("custom_date", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "chores.md")
		if err := os.WriteFile(testFile, []byte(baseContent), 0644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		date := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
		var buf bytes.Buffer

		if err := DoneCmd(testFile, "Kitchen Clean", date, &buf); err != nil {
			t.Fatalf("DoneCmd error: %v", err)
		}

		content, _ := os.ReadFile(testFile)
		if !strings.Contains(string(content), "2026-01-15") {
			t.Errorf("should use provided date, got:\n%s", string(content))
		}
	})

	t.Run("no_trailing_newline", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "chores.md")
		contentNoNewline := strings.TrimSuffix(baseContent, "\n")
		if err := os.WriteFile(testFile, []byte(contentNoNewline), 0644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		date := time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC)
		var buf bytes.Buffer

		if err := DoneCmd(testFile, "Kitchen Clean", date, &buf); err != nil {
			t.Fatalf("DoneCmd error: %v", err)
		}

		content, _ := os.ReadFile(testFile)
		lines := strings.Split(string(content), "\n")
		lastNonEmpty := ""
		for i := len(lines) - 1; i >= 0; i-- {
			if strings.TrimSpace(lines[i]) != "" {
				lastNonEmpty = lines[i]
				break
			}
		}
		if !strings.Contains(lastNonEmpty, "2026-02-10 Kitchen Clean") {
			t.Errorf("entry should be on its own line, last line: %q", lastNonEmpty)
		}
	})
}
