package schedule

import (
	"sort"
	"strings"
	"time"

	"github.com/kusha/chores-md/internal/model"
)

type Status int

const (
	StatusOverdue Status = iota
	StatusDueToday
	StatusUpcoming
	StatusClear
)

type ChoreStatus struct {
	Chore       model.Chore
	Status      Status
	DaysOverdue int
	DaysUntil   int
	LastDone    *time.Time
}

const NeverDoneSentinel = 999999

func DaysBetween(from, to time.Time) int {
	fromUTC := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	toUTC := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.UTC)
	return int(toUTC.Sub(fromUTC).Hours() / 24)
}

func Calculate(chores []model.Chore, completions []model.Completion, now time.Time) []ChoreStatus {
	completionMap := make(map[string]time.Time)
	for _, c := range completions {
		key := strings.ToLower(c.ChoreName)
		if existing, ok := completionMap[key]; !ok || c.Date.After(existing) {
			completionMap[key] = c.Date
		}
	}

	var results []ChoreStatus

	for _, chore := range chores {
		key := strings.ToLower(chore.Name)
		cs := ChoreStatus{Chore: chore}

		lastDone, hasCompletion := completionMap[key]
		if hasCompletion {
			cs.LastDone = &lastDone
		}

		if !hasCompletion {
			cs.Status = StatusOverdue
			cs.DaysOverdue = NeverDoneSentinel
		} else {
			daysSince := DaysBetween(lastDone, now)
			freq := chore.FrequencyDays

			if daysSince > freq {
				cs.Status = StatusOverdue
				cs.DaysOverdue = daysSince - freq
			} else if daysSince == freq {
				cs.Status = StatusDueToday
			} else {
				daysUntil := freq - daysSince
				if daysUntil <= 7 {
					cs.Status = StatusUpcoming
					cs.DaysUntil = daysUntil
				} else {
					cs.Status = StatusClear
					cs.DaysUntil = daysUntil
				}
			}
		}

		results = append(results, cs)
	}

	return results
}

func SortByUrgency(statuses []ChoreStatus) {
	sort.SliceStable(statuses, func(i, j int) bool {
		si, sj := statuses[i], statuses[j]

		if si.Status != sj.Status {
			return si.Status < sj.Status
		}

		switch si.Status {
		case StatusOverdue:
			if si.DaysOverdue == NeverDoneSentinel && sj.DaysOverdue != NeverDoneSentinel {
				return false
			}
			if si.DaysOverdue != NeverDoneSentinel && sj.DaysOverdue == NeverDoneSentinel {
				return true
			}
			if si.DaysOverdue != sj.DaysOverdue {
				return si.DaysOverdue > sj.DaysOverdue
			}
		case StatusUpcoming, StatusClear:
			if si.DaysUntil != sj.DaysUntil {
				return si.DaysUntil < sj.DaysUntil
			}
		}

		return strings.ToLower(si.Chore.Name) < strings.ToLower(sj.Chore.Name)
	})
}
