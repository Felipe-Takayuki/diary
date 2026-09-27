package domain_test

import (
	"testing"

	"diary/internal/domain"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      string
		wantError error
	}{
		{
			name:      "canonical DD-MM-YYYY",
			input:     "24-09-2026",
			want:      "24-09-2026",
			wantError: nil,
		},
		{
			name:      "slash DD/MM/YYYY",
			input:     "24/09/2026",
			want:      "24-09-2026",
			wantError: nil,
		},
		{
			name:      "iso YYYY-MM-DD",
			input:     "2026-09-24",
			want:      "24-09-2026",
			wantError: nil,
		},
		{
			name:      "with spaces",
			input:     "  24-09-2026  ",
			want:      "24-09-2026",
			wantError: nil,
		},
		{
			name:      "empty date",
			input:     "",
			want:      "",
			wantError: domain.ErrEmptyDate,
		},
		{
			name:      "only spaces",
			input:     "   ",
			want:      "",
			wantError: domain.ErrEmptyDate,
		},
		{
			name:      "invalid format",
			input:     "not-a-date",
			want:      "",
			wantError: domain.ErrInvalidDate,
		},
		{
			name:      "invalid day",
			input:     "32-01-2026",
			want:      "",
			wantError: domain.ErrInvalidDate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.ParseDate(tt.input)
			if tt.wantError != nil {
				if err != tt.wantError {
					t.Fatalf("expected error %v, got %v", tt.wantError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.String() != tt.want {
				t.Errorf("got %s, want %s", got.String(), tt.want)
			}
		})
	}
}

func TestDate_LegacyISO(t *testing.T) {
	d, err := domain.ParseDate("24-09-2026")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.LegacyISO() != "2026-09-24" {
		t.Errorf("got %s, want 2026-09-24", d.LegacyISO())
	}
}

func TestDate_AddDaysAndToday(t *testing.T) {
	d, err := domain.ParseDate("27-09-2026")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	next := d.AddDays(1)
	if next.String() != "28-09-2026" {
		t.Errorf("got %s, want 28-09-2026", next.String())
	}

	prev := d.AddDays(-1)
	if prev.String() != "26-09-2026" {
		t.Errorf("got %s, want 26-09-2026", prev.String())
	}

	today := domain.Today()
	if today.IsZero() {
		t.Errorf("expected today to not be zero")
	}
}
