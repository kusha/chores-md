package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestShowCmd(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "chores.md")

	now := time.Date(2026, 2, 10, 12, 0, 0, 0, time.UTC)

	content := `## Overdue Task
> 1w

## Due Today Task
> 1w

## Upcoming Task
> 1w

## Clear Task
> 2w

## Never Done Task
> 1d

2026-01-31 Overdue Task
2026-02-03 Due Today Task
2026-02-05 Upcoming Task
2026-02-09 Clear Task
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	var buf bytes.Buffer
	if err := ShowCmd(testFile, now, &buf); err != nil {
		t.Fatalf("ShowCmd error: %v", err)
	}

	output := buf.String()

	t.Run("contains_overdue_section", func(t *testing.T) {
		if !strings.Contains(output, "OVERDUE") {
			t.Error("missing OVERDUE section")
		}
		if !strings.Contains(output, "3 days overdue") {
			t.Errorf("missing overdue count, got:\n%s", output)
		}
	})

	t.Run("contains_due_today_section", func(t *testing.T) {
		if !strings.Contains(output, "DUE TODAY") {
			t.Error("missing DUE TODAY section")
		}
	})

	t.Run("contains_upcoming_section", func(t *testing.T) {
		if !strings.Contains(output, "UPCOMING (7 days)") {
			t.Error("missing UPCOMING section")
		}
		if !strings.Contains(output, "due in 2 days") {
			t.Errorf("missing upcoming days, got:\n%s", output)
		}
	})

	t.Run("contains_clear_section", func(t *testing.T) {
		if !strings.Contains(output, "ALL CLEAR") {
			t.Error("missing ALL CLEAR section")
		}
	})

	t.Run("never_done_format", func(t *testing.T) {
		if !strings.Contains(output, "(never done)") {
			t.Errorf("never-done chore should show '(never done)', got:\n%s", output)
		}
		if !strings.Contains(output, "Last: never") {
			t.Errorf("never-done chore should show 'Last: never', got:\n%s", output)
		}
	})
}

func TestShowCmd_equal_urgency_sort(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "chores.md")

	now := time.Date(2026, 2, 10, 12, 0, 0, 0, time.UTC)

	content := `## Zebra Task
> 1w

## Alpha Task
> 1w

2026-01-31 Zebra Task
2026-01-31 Alpha Task
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	var buf bytes.Buffer
	if err := ShowCmd(testFile, now, &buf); err != nil {
		t.Fatalf("ShowCmd error: %v", err)
	}

	output := buf.String()
	alphaIdx := strings.Index(output, "Alpha Task")
	zebraIdx := strings.Index(output, "Zebra Task")

	if alphaIdx == -1 || zebraIdx == -1 {
		t.Fatalf("missing tasks in output:\n%s", output)
	}

	if alphaIdx > zebraIdx {
		t.Errorf("Alpha Task should appear before Zebra Task (alphabetical tie-breaker)")
	}
}
