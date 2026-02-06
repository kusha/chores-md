package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/kusha/chores-md/internal/parser"
)

func DoneCmd(file string, choreName string, date time.Time, out io.Writer) error {
	result, err := parser.ParseFile(file)
	if err != nil {
		return err
	}

	var found bool
	var matchedName string
	choreNameLower := strings.ToLower(strings.TrimSpace(choreName))

	for _, chore := range result.Chores {
		if strings.ToLower(chore.Name) == choreNameLower {
			found = true
			matchedName = chore.Name
			break
		}
	}

	if !found {
		return fmt.Errorf("chore not found: %q", choreName)
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	dateStr := date.Format("2006-01-02")
	entry := fmt.Sprintf("%s %s\n", dateStr, matchedName)

	if len(content) > 0 && content[len(content)-1] != '\n' {
		entry = "\n" + entry
	}

	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(entry); err != nil {
		return err
	}

	fmt.Fprintf(out, "Done: %q (%s)\n", matchedName, dateStr)
	return nil
}
