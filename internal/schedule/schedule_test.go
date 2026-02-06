package schedule

import (
	"testing"
	"time"

	"github.com/user/chores/internal/model"
)

func date(y, m, d int) time.Time {
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
}

func TestDaysBetween(t *testing.T) {
	tests := []struct {
		from, to time.Time
		want     int
	}{
		{date(2026, 2, 1), date(2026, 2, 3), 2},
		{date(2026, 2, 3), date(2026, 2, 1), -2},
		{date(2026, 2, 1), date(2026, 2, 1), 0},
		{date(2026, 1, 31), date(2026, 2, 1), 1},
	}

	for _, tt := range tests {
		got := DaysBetween(tt.from, tt.to)
		if got != tt.want {
			t.Errorf("DaysBetween(%v, %v) = %d, want %d", tt.from, tt.to, got, tt.want)
		}
	}
}

func TestCalculate(t *testing.T) {
	now := date(2026, 2, 10)

	t.Run("overdue", func(t *testing.T) {
		chores := []model.Chore{{Name: "Test", FrequencyDays: 7}}
		completions := []model.Completion{{ChoreName: "Test", Date: date(2026, 1, 31)}}

		results := Calculate(chores, completions, now)
		if len(results) != 1 {
			t.Fatalf("got %d results, want 1", len(results))
		}

		cs := results[0]
		if cs.Status != StatusOverdue {
			t.Errorf("status = %v, want StatusOverdue", cs.Status)
		}
		if cs.DaysOverdue != 3 {
			t.Errorf("DaysOverdue = %d, want 3 (10 days since, 7 day freq)", cs.DaysOverdue)
		}
	})

	t.Run("due_today", func(t *testing.T) {
		chores := []model.Chore{{Name: "Test", FrequencyDays: 7}}
		completions := []model.Completion{{ChoreName: "Test", Date: date(2026, 2, 3)}}

		results := Calculate(chores, completions, now)
		cs := results[0]
		if cs.Status != StatusDueToday {
			t.Errorf("status = %v, want StatusDueToday", cs.Status)
		}
	})

	t.Run("upcoming", func(t *testing.T) {
		chores := []model.Chore{{Name: "Test", FrequencyDays: 7}}
		completions := []model.Completion{{ChoreName: "Test", Date: date(2026, 2, 5)}}

		results := Calculate(chores, completions, now)
		cs := results[0]
		if cs.Status != StatusUpcoming {
			t.Errorf("status = %v, want StatusUpcoming", cs.Status)
		}
		if cs.DaysUntil != 2 {
			t.Errorf("DaysUntil = %d, want 2 (5 days since, freq 7, due in 2)", cs.DaysUntil)
		}
	})

	t.Run("clear", func(t *testing.T) {
		chores := []model.Chore{{Name: "Test", FrequencyDays: 14}}
		completions := []model.Completion{{ChoreName: "Test", Date: date(2026, 2, 9)}}

		results := Calculate(chores, completions, now)
		cs := results[0]
		if cs.Status != StatusClear {
			t.Errorf("status = %v, want StatusClear", cs.Status)
		}
		if cs.DaysUntil != 13 {
			t.Errorf("DaysUntil = %d, want 13", cs.DaysUntil)
		}
	})

	t.Run("never_done", func(t *testing.T) {
		chores := []model.Chore{{Name: "Test", FrequencyDays: 7}}
		var completions []model.Completion

		results := Calculate(chores, completions, now)
		cs := results[0]
		if cs.Status != StatusOverdue {
			t.Errorf("status = %v, want StatusOverdue", cs.Status)
		}
		if cs.LastDone != nil {
			t.Errorf("LastDone should be nil for never-done chore")
		}
		if cs.DaysOverdue != NeverDoneSentinel {
			t.Errorf("DaysOverdue = %d, want %d (sentinel)", cs.DaysOverdue, NeverDoneSentinel)
		}
	})

	t.Run("case_insensitive_match", func(t *testing.T) {
		chores := []model.Chore{{Name: "Kitchen Clean", FrequencyDays: 7}}
		completions := []model.Completion{{ChoreName: "kitchen clean", Date: date(2026, 2, 9)}}

		results := Calculate(chores, completions, now)
		cs := results[0]
		if cs.LastDone == nil {
			t.Error("should match completion case-insensitively")
		}
	})
}

func TestSortByUrgency(t *testing.T) {
	t.Run("equal_urgency_alphabetical", func(t *testing.T) {
		statuses := []ChoreStatus{
			{Chore: model.Chore{Name: "Zebra"}, Status: StatusOverdue, DaysOverdue: 3},
			{Chore: model.Chore{Name: "Alpha"}, Status: StatusOverdue, DaysOverdue: 3},
		}

		SortByUrgency(statuses)

		if statuses[0].Chore.Name != "Alpha" {
			t.Errorf("expected Alpha first (alphabetical tie-breaker), got %s", statuses[0].Chore.Name)
		}
	})

	t.Run("never_done_at_end_of_overdue", func(t *testing.T) {
		statuses := []ChoreStatus{
			{Chore: model.Chore{Name: "Never Done"}, Status: StatusOverdue, DaysOverdue: NeverDoneSentinel},
			{Chore: model.Chore{Name: "3 Days Over"}, Status: StatusOverdue, DaysOverdue: 3},
		}

		SortByUrgency(statuses)

		if statuses[0].Chore.Name != "3 Days Over" {
			t.Errorf("expected dated overdue first, got %s", statuses[0].Chore.Name)
		}
		if statuses[1].Chore.Name != "Never Done" {
			t.Errorf("expected never-done at end, got %s", statuses[1].Chore.Name)
		}
	})

	t.Run("status_order", func(t *testing.T) {
		statuses := []ChoreStatus{
			{Chore: model.Chore{Name: "Clear"}, Status: StatusClear, DaysUntil: 10},
			{Chore: model.Chore{Name: "Upcoming"}, Status: StatusUpcoming, DaysUntil: 2},
			{Chore: model.Chore{Name: "DueToday"}, Status: StatusDueToday},
			{Chore: model.Chore{Name: "Overdue"}, Status: StatusOverdue, DaysOverdue: 1},
		}

		SortByUrgency(statuses)

		expected := []string{"Overdue", "DueToday", "Upcoming", "Clear"}
		for i, name := range expected {
			if statuses[i].Chore.Name != name {
				t.Errorf("position %d: got %s, want %s", i, statuses[i].Chore.Name, name)
			}
		}
	})
}
