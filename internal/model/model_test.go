package model

import (
	"testing"
)

func TestParseFrequency(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantDays int
		wantRaw  string
		wantErr  bool
	}{
		{"1d", "1d", 1, "1d", false},
		{"2d", "2d", 2, "2d", false},
		{"1w", "1w", 7, "1w", false},
		{"2w", "2w", 14, "2w", false},
		{"1m", "1m", 30, "1m", false},
		{"3m", "3m", 90, "3m", false},
		{"1y", "1y", 365, "1y", false},
		{"0d", "0d", 0, "", true},
		{"invalid", "invalid", 0, "", true},
		{"2x", "2x", 0, "", true},
		{"empty", "", 0, "", true},
		{"negative", "-1d", 0, "", true},
		{"no_unit", "10", 0, "", true},
		{"decimal", "1.5w", 0, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			days, raw, err := ParseFrequency(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseFrequency(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseFrequency(%q) unexpected error: %v", tt.input, err)
				return
			}

			if days != tt.wantDays {
				t.Errorf("ParseFrequency(%q) days = %d, want %d", tt.input, days, tt.wantDays)
			}

			if raw != tt.wantRaw {
				t.Errorf("ParseFrequency(%q) raw = %q, want %q", tt.input, raw, tt.wantRaw)
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantMinutes int
		wantRaw     string
		wantErr     bool
	}{
		{"30m", "30m", 30, "30m", false},
		{"2h", "2h", 120, "2h", false},
		{"1h30m", "1h30m", 90, "1h30m", false},
		{"90m", "90m", 90, "90m", false},
		{"1h", "1h", 60, "1h", false},
		{"3h", "3h", 180, "3h", false},
		{"45m", "45m", 45, "45m", false},
		{"2h15m", "2h15m", 135, "2h15m", false},
		{"0m", "0m", 0, "", true},
		{"0h", "0h", 0, "", true},
		{"0h0m", "0h0m", 0, "", true},
		{"invalid", "abc", 0, "", true},
		{"empty", "", 0, "", true},
		{"negative", "-30m", 0, "", true},
		{"no_unit", "30", 0, "", true},
		{"space_separated", "1h 30m", 0, "", true},
		{"m_only_invalid", "m", 0, "", true},
		{"h_only_invalid", "h", 0, "", true},
		{"decimal_minutes", "30.5m", 0, "", true},
		{"mixed_invalid", "1h30", 0, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			minutes, raw, err := ParseDuration(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseDuration(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseDuration(%q) unexpected error: %v", tt.input, err)
				return
			}

			if minutes != tt.wantMinutes {
				t.Errorf("ParseDuration(%q) minutes = %d, want %d", tt.input, minutes, tt.wantMinutes)
			}

			if raw != tt.wantRaw {
				t.Errorf("ParseDuration(%q) raw = %q, want %q", tt.input, raw, tt.wantRaw)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		wantStr string
	}{
		{"30m", 30, "30m"},
		{"60m_to_1h", 60, "1h"},
		{"90m_to_1h30m", 90, "1h 30m"},
		{"120m_to_2h", 120, "2h"},
		{"45m", 45, "45m"},
		{"75m_to_1h15m", 75, "1h 15m"},
		{"180m_to_3h", 180, "3h"},
		{"1m", 1, "1m"},
		{"zero", 0, "0m"},
		{"negative", -30, "0m"},
		{"large_value", 1440, "24h"},
		{"large_with_minutes", 1485, "24h 45m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatDuration(tt.input)
			if got != tt.wantStr {
				t.Errorf("FormatDuration(%d) = %q, want %q", tt.input, got, tt.wantStr)
			}
		})
	}
}
