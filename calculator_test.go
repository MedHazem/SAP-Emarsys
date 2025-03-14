package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandleDueDateRequest(t *testing.T) {
	// Start a test server
	ts := httptest.NewServer(http.HandlerFunc(HandleDueDateRequest))
	defer ts.Close()

	tests := []struct {
		name       string
		request    DueDateRequest
		expectCode int
	}{
		{
			name: "Valid request",
			request: DueDateRequest{
				SubmitTime:      "2025-03-14 10:00",
				TurnaroundHours: 5,
			},
			expectCode: http.StatusOK,
		},
		{
			name: "Submit time on weekend",
			request: DueDateRequest{
				SubmitTime:      "2025-03-16 10:00", // Sunday
				TurnaroundHours: 5,
			},
			expectCode: http.StatusBadRequest,
		},
		{
			name: "Submit time outside business hours",
			request: DueDateRequest{
				SubmitTime:      "2025-03-14 07:00",
				TurnaroundHours: 3,
			},
			expectCode: http.StatusBadRequest,
		},
		{
			name: "Negative turnaround hours",
			request: DueDateRequest{
				SubmitTime:      "2025-03-14 10:00",
				TurnaroundHours: -2,
			},
			expectCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal request
			body, err := json.Marshal(tt.request)
			if err != nil {
				t.Fatalf("Failed to marshal request: %v", err)
			}

			// Send request
			resp, err := http.Post(ts.URL+"/calculate-due-date", "application/json", bytes.NewBuffer(body))
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			// Check response code
			if resp.StatusCode != tt.expectCode {
				t.Errorf("expected status %d, got %d", tt.expectCode, resp.StatusCode)
			}

			// Check response body for successful cases
			if tt.expectCode == http.StatusOK {
				var response DueDateResponse
				if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				// Parse submitTime
				layout := "2006-01-02 15:04"
				submitTime, _ := time.Parse(layout, tt.request.SubmitTime)
				expectedDueDate, _ := CalculateDueDate(submitTime, tt.request.TurnaroundHours)

				if response.DueDate != expectedDueDate.Format(layout) {
					t.Errorf("expected due date %s, got %s", expectedDueDate.Format(layout), response.DueDate)
				}
			}
		})
	}
}
