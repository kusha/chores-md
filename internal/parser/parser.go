// Package parser handles parsing of chores.md files.
package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/kusha/chores-md/internal/model"
)

type ParseResult struct {
	Chores      []model.Chore
	Completions []model.Completion
	Warnings    []string
}

var (
	headerRegex     = regexp.MustCompile(`^##\s+(.+)$`)
	frequencyRegex  = regexp.MustCompile(`^>\s*(\d+[dwmy])(?:\s+(.+))?\s*$`)
	completionRegex = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})\s+(.+?)(?:\s*#.*)?$`)
)

func Parse(content string) (*ParseResult, error) {
	result := &ParseResult{}
	choreMap := make(map[string]bool)

	lines := strings.Split(content, "\n")

	var currentChore *model.Chore
	var descLines []string

	for i, line := range lines {
		lineNum := i + 1
		line = strings.TrimRight(line, "\r")

		if matches := headerRegex.FindStringSubmatch(line); matches != nil {
			if currentChore != nil {
				currentChore.Description = strings.TrimSpace(strings.Join(descLines, "\n"))
				result.Chores = append(result.Chores, *currentChore)
			}

			choreName := strings.TrimSpace(matches[1])
			nameKey := strings.ToLower(choreName)

			if choreMap[nameKey] {
				result.Warnings = append(result.Warnings, fmt.Sprintf("line %d: duplicate chore %q (first definition wins)", lineNum, choreName))
				currentChore = nil
				descLines = nil
				continue
			}

			choreMap[nameKey] = true
			currentChore = &model.Chore{
				Name: choreName,
				Line: lineNum,
			}
			descLines = nil
			continue
		}

		if currentChore != nil && currentChore.FrequencyDays == 0 {
			if matches := frequencyRegex.FindStringSubmatch(line); matches != nil {
				days, raw, err := model.ParseFrequency(matches[1])
				if err != nil {
					return nil, fmt.Errorf("line %d: %w", lineNum, err)
				}
				currentChore.FrequencyDays = days
				currentChore.FrequencyRaw = raw

				// Check if there's a duration part (group 2)
				if len(matches) > 2 && matches[2] != "" {
					durationStr := strings.TrimSpace(matches[2])
					if durationStr != "" {
						minutes, durationRaw, err := model.ParseDuration(durationStr)
						if err != nil {
							return nil, fmt.Errorf("line %d: %w", lineNum, err)
						}
						currentChore.DurationMinutes = minutes
						currentChore.DurationRaw = durationRaw
					}
				}
				continue
			}
		}

		if matches := completionRegex.FindStringSubmatch(line); matches != nil {
			dateStr := matches[1]
			choreName := strings.TrimSpace(matches[2])

			date, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				result.Warnings = append(result.Warnings, fmt.Sprintf("line %d: invalid date %q, skipping", lineNum, dateStr))
				continue
			}

			date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

			result.Completions = append(result.Completions, model.Completion{
				Date:      date,
				ChoreName: choreName,
				Line:      lineNum,
			})
			continue
		}

		if currentChore != nil && currentChore.FrequencyDays > 0 && strings.TrimSpace(line) != "" {
			descLines = append(descLines, line)
		}
	}

	if currentChore != nil {
		currentChore.Description = strings.TrimSpace(strings.Join(descLines, "\n"))
		result.Chores = append(result.Chores, *currentChore)
	}

	for i := range result.Chores {
		if result.Chores[i].FrequencyDays == 0 {
			return nil, fmt.Errorf("line %d: chore %q has no frequency defined", result.Chores[i].Line, result.Chores[i].Name)
		}
	}

	return result, nil
}

func ParseFile(path string) (*ParseResult, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sb strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sb.WriteString(scanner.Text())
		sb.WriteString("\n")
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return Parse(sb.String())
}
