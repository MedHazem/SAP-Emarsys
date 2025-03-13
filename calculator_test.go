package main

import (
	"testing"
	"time"
)

func TestCalculateDueDate(t *testing.T) {
	layout := "2006-01-02 15:04"
	cases := []struct {
		submitTime      string
		turnaroundHours int
		expectedDueDate string
	}{
		{"2025-03-13 10:00", 2, "2025-03-13 12:00"},  // Same day completion
		{"2025-03-13 16:00", 2, "2025-03-14 10:00"},  // Next day due
		{"2025-03-15 10:00", 2, "INVALID"},           // Weekend submission
		{"2025-03-13 15:00", 10, "2025-03-14 17:00"}, // Spans across multiple days
	}

	for _, tc := range cases {
		submitTime, err := time.Parse(layout, tc.submitTime)
		if err != nil {
			t.Fatalf("Failed to parse submitTime: %v", err)
		}

		if submitTime.Weekday() == time.Saturday || submitTime.Weekday() == time.Sunday {
			if tc.expectedDueDate != "INVALID" {
				t.Errorf("Expected INVALID for weekend submission, but got valid output")
			}
			continue
		}

		dueDate := CalculateDueDate(submitTime, tc.turnaroundHours)
		if dueDate.Format(layout) != tc.expectedDueDate {
			t.Errorf("For submitTime %s with %d hours, expected %s but got %s",
				tc.submitTime, tc.turnaroundHours, tc.expectedDueDate, dueDate.Format(layout))
		}
	}
}
