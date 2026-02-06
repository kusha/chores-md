// Package model defines the core data types for the chores CLI.
package model

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// Chore represents a recurring household task defined in the markdown file.
type Chore struct {
	Name            string // The chore name from ## header
	FrequencyDays   int    // Frequency converted to days
	FrequencyRaw    string // Original frequency token (e.g., "2w") for display
	DurationMinutes int    // Duration in minutes (optional)
	DurationRaw     string // Original duration token (e.g., "1h30m") for display
	Description     string // Optional description text after the header
	Line            int    // Line number in file for error reporting
}

// Completion represents a single completion entry (date + chore name).
type Completion struct {
	Date      time.Time // The date the chore was completed
	ChoreName string    // The chore name as written in the completion entry
	Line      int       // Line number in file for error reporting
}

// frequencyRegex matches frequency patterns like "1d", "2w", "1m", "3y"
var frequencyRegex = regexp.MustCompile(`^(\d+)([dwmy])$`)

// durationRegex matches duration patterns like "30m", "2h", "1h30m"
var durationRegex = regexp.MustCompile(`^(\d+)h?(?:(\d+)m)?$|^(\d+)m$`)

// ParseFrequency parses a frequency string like "2w" and returns the number of days,
// the original raw string, and any error encountered.
//
// Supported units:
//   - d: days (N * 1)
//   - w: weeks (N * 7)
//   - m: months (N * 30)
//   - y: years (N * 365)
//
// Returns an error for invalid formats or zero/negative values.
func ParseFrequency(s string) (days int, raw string, err error) {
	matches := frequencyRegex.FindStringSubmatch(s)
	if matches == nil {
		return 0, "", fmt.Errorf("invalid frequency format: %q (expected format like 1d, 2w, 1m, 1y)", s)
	}

	n, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, "", fmt.Errorf("invalid frequency number: %q", matches[1])
	}

	if n <= 0 {
		return 0, "", fmt.Errorf("frequency must be positive, got: %d", n)
	}

	unit := matches[2]
	var multiplier int
	switch unit {
	case "d":
		multiplier = 1
	case "w":
		multiplier = 7
	case "m":
		multiplier = 30
	case "y":
		multiplier = 365
	default:
		return 0, "", fmt.Errorf("unknown frequency unit: %q", unit)
	}

	return n * multiplier, s, nil
}

// ParseDuration parses a duration string like "30m", "2h", or "1h30m" and returns the total minutes,
// the original raw string, and any error encountered.
//
// Supported formats:
//   - Xm: minutes only (e.g., "30m", "90m")
//   - Xh: hours only (e.g., "2h")
//   - XhYm: hours and minutes (e.g., "1h30m")
//
// Returns an error for invalid formats or zero/negative values.
func ParseDuration(s string) (minutes int, raw string, err error) {
	if s == "" {
		return 0, "", fmt.Errorf("invalid duration format: %q (expected format like 30m, 2h, 1h30m)", s)
	}

	// Try to match pattern like "30m" or "2h" or "1h30m"
	var hours, mins int

	// Try XhYm format first (e.g., "1h30m")
	if matched := regexp.MustCompile(`^(\d+)h(\d+)m$`).FindStringSubmatch(s); matched != nil {
		h, _ := strconv.Atoi(matched[1])
		m, _ := strconv.Atoi(matched[2])
		hours = h
		mins = m
	} else if matched := regexp.MustCompile(`^(\d+)h$`).FindStringSubmatch(s); matched != nil {
		// Try Xh format (e.g., "2h")
		h, _ := strconv.Atoi(matched[1])
		hours = h
		mins = 0
	} else if matched := regexp.MustCompile(`^(\d+)m$`).FindStringSubmatch(s); matched != nil {
		// Try Xm format (e.g., "30m")
		m, _ := strconv.Atoi(matched[1])
		hours = 0
		mins = m
	} else {
		return 0, "", fmt.Errorf("invalid duration format: %q (expected format like 30m, 2h, 1h30m)", s)
	}

	totalMinutes := hours*60 + mins

	if totalMinutes <= 0 {
		return 0, "", fmt.Errorf("duration must be positive, got: %d minutes", totalMinutes)
	}

	return totalMinutes, s, nil
}

// FormatDuration formats a duration in minutes as a human-readable string.
// Returns "Xh Ym" for combined hours and minutes, "Xh" for hours only, "Xm" for minutes only.
func FormatDuration(minutes int) string {
	if minutes <= 0 {
		return "0m"
	}

	hours := minutes / 60
	mins := minutes % 60

	if hours > 0 && mins > 0 {
		return fmt.Sprintf("%dh %dm", hours, mins)
	} else if hours > 0 {
		return fmt.Sprintf("%dh", hours)
	} else {
		return fmt.Sprintf("%dm", mins)
	}
}
