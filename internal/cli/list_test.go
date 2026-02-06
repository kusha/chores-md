package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListCmd(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "chores.md")

	content := `## Zebra Task
> 1w

## Alpha Task
> 2d

## Beta Task
> 1m

2026-02-03 Alpha Task
2026-02-01 Zebra Task
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	var buf bytes.Buffer
	if err := ListCmd(testFile, &buf); err != nil {
		t.Fatalf("ListCmd error: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 3 {
		t.Errorf("got %d lines, want 3", len(lines))
	}

	if !strings.HasPrefix(lines[0], "Alpha Task") {
		t.Errorf("first line should be Alpha Task (alphabetical), got: %s", lines[0])
	}
	if !strings.HasPrefix(lines[1], "Beta Task") {
		t.Errorf("second line should be Beta Task, got: %s", lines[1])
	}
	if !strings.HasPrefix(lines[2], "Zebra Task") {
		t.Errorf("third line should be Zebra Task, got: %s", lines[2])
	}

	if !strings.Contains(lines[0], "every 2d") {
		t.Errorf("Alpha Task should show 'every 2d', got: %s", lines[0])
	}
	if !strings.Contains(lines[0], "Last: 2026-02-03") {
		t.Errorf("Alpha Task should show last completion, got: %s", lines[0])
	}

	if !strings.Contains(lines[1], "Last: never") {
		t.Errorf("Beta Task should show 'never', got: %s", lines[1])
	}
}
