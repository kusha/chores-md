package cli

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/kusha/chores-md/internal/model"
	"github.com/kusha/chores-md/internal/parser"
)

func ListCmd(file string, out io.Writer) error {
	result, err := parser.ParseFile(file)
	if err != nil {
		return err
	}

	completionMap := make(map[string]string)
	for _, c := range result.Completions {
		key := strings.ToLower(c.ChoreName)
		dateStr := c.Date.Format("2006-01-02")
		if existing, ok := completionMap[key]; !ok || dateStr > existing {
			completionMap[key] = dateStr
		}
	}

	chores := make([]model.Chore, len(result.Chores))
	copy(chores, result.Chores)
	sort.Slice(chores, func(i, j int) bool {
		return strings.ToLower(chores[i].Name) < strings.ToLower(chores[j].Name)
	})

	for _, chore := range chores {
		key := strings.ToLower(chore.Name)
		lastDone := "never"
		if date, ok := completionMap[key]; ok {
			lastDone = date
		}
		durationStr := ""
		if chore.DurationMinutes > 0 {
			durationStr = " ~" + model.FormatDuration(chore.DurationMinutes)
		}
		fmt.Fprintf(out, "%s\tevery %s%s\tLast: %s\n", chore.Name, chore.FrequencyRaw, durationStr, lastDone)
	}

	return nil
}
