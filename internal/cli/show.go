package cli

import (
	"fmt"
	"io"
	"time"

	"github.com/kusha/chores-md/internal/model"
	"github.com/kusha/chores-md/internal/parser"
	"github.com/kusha/chores-md/internal/schedule"
)

func ShowCmd(file string, now time.Time, out io.Writer) error {
	result, err := parser.ParseFile(file)
	if err != nil {
		return err
	}

	statuses := schedule.Calculate(result.Chores, result.Completions, now)
	schedule.SortByUrgency(statuses)

	var overdue, dueToday, upcoming, clear []schedule.ChoreStatus
	for _, cs := range statuses {
		switch cs.Status {
		case schedule.StatusOverdue:
			overdue = append(overdue, cs)
		case schedule.StatusDueToday:
			dueToday = append(dueToday, cs)
		case schedule.StatusUpcoming:
			upcoming = append(upcoming, cs)
		case schedule.StatusClear:
			clear = append(clear, cs)
		}
	}

	if len(overdue) > 0 {
		fmt.Fprintln(out, "OVERDUE")
		var totalMinutes int
		for _, cs := range overdue {
			durationStr := ""
			if cs.Chore.DurationMinutes > 0 {
				durationStr = fmt.Sprintf("(~%s) ", model.FormatDuration(cs.Chore.DurationMinutes))
				totalMinutes += cs.Chore.DurationMinutes
			}
			if cs.DaysOverdue == schedule.NeverDoneSentinel {
				fmt.Fprintf(out, "  %s %s(never done)\n", cs.Chore.Name, durationStr)
				fmt.Fprintln(out, "    Last: never")
			} else {
				fmt.Fprintf(out, "  %s %s(%d days overdue)\n", cs.Chore.Name, durationStr, cs.DaysOverdue)
				fmt.Fprintf(out, "    Last: %s\n", cs.LastDone.Format("2006-01-02"))
			}
		}
		if totalMinutes > 0 {
			fmt.Fprintf(out, "  Total: %s\n", model.FormatDuration(totalMinutes))
		}
		fmt.Fprintln(out)
	}

	if len(dueToday) > 0 {
		fmt.Fprintln(out, "DUE TODAY")
		var totalMinutes int
		for _, cs := range dueToday {
			durationStr := ""
			if cs.Chore.DurationMinutes > 0 {
				durationStr = fmt.Sprintf("(~%s) ", model.FormatDuration(cs.Chore.DurationMinutes))
				totalMinutes += cs.Chore.DurationMinutes
			}
			fmt.Fprintf(out, "  %s %s\n", cs.Chore.Name, durationStr)
			fmt.Fprintf(out, "    Last: %s\n", cs.LastDone.Format("2006-01-02"))
		}
		if totalMinutes > 0 {
			fmt.Fprintf(out, "  Total: %s\n", model.FormatDuration(totalMinutes))
		}
		fmt.Fprintln(out)
	}

	if len(upcoming) > 0 {
		fmt.Fprintln(out, "UPCOMING (7 days)")
		var totalMinutes int
		for _, cs := range upcoming {
			durationStr := ""
			if cs.Chore.DurationMinutes > 0 {
				durationStr = fmt.Sprintf("(~%s) ", model.FormatDuration(cs.Chore.DurationMinutes))
				totalMinutes += cs.Chore.DurationMinutes
			}
			fmt.Fprintf(out, "  %s %s(due in %d day", cs.Chore.Name, durationStr, cs.DaysUntil)
			if cs.DaysUntil != 1 {
				fmt.Fprint(out, "s")
			}
			fmt.Fprintln(out, ")")
			fmt.Fprintf(out, "    Last: %s\n", cs.LastDone.Format("2006-01-02"))
		}
		if totalMinutes > 0 {
			fmt.Fprintf(out, "  Total: %s\n", model.FormatDuration(totalMinutes))
		}
		fmt.Fprintln(out)
	}

	if len(clear) > 0 {
		fmt.Fprintln(out, "ALL CLEAR")
		var totalMinutes int
		for _, cs := range clear {
			durationStr := ""
			if cs.Chore.DurationMinutes > 0 {
				durationStr = fmt.Sprintf("(~%s) ", model.FormatDuration(cs.Chore.DurationMinutes))
				totalMinutes += cs.Chore.DurationMinutes
			}
			fmt.Fprintf(out, "  %s %s(due in %d days)\n", cs.Chore.Name, durationStr, cs.DaysUntil)
			fmt.Fprintf(out, "    Last: %s\n", cs.LastDone.Format("2006-01-02"))
		}
		if totalMinutes > 0 {
			fmt.Fprintf(out, "  Total: %s\n", model.FormatDuration(totalMinutes))
		}
		fmt.Fprintln(out)
	}

	return nil
}
