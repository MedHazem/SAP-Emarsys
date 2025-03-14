package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	workDayStart = 9
	workDayEnd   = 17
	hoursPerDay  = workDayEnd - workDayStart
)

// Request structure to receive input
type DueDateRequest struct {
	SubmitTime      string `json:"submitTime"`
	TurnaroundHours int    `json:"turnaroundHours"`
}

// Response structure to send the result
type DueDateResponse struct {
	DueDate string `json:"dueDate"`
}

// CalculateDueDate calculates the due date based on business hours (Mon–Fri, 9 AM–5 PM)
func CalculateDueDate(submitTime time.Time, turnaroundHours int) (time.Time, error) {
	// If submitTime is on the weekend, return an error
	if submitTime.Weekday() == time.Saturday || submitTime.Weekday() == time.Sunday {
		return submitTime, fmt.Errorf("submit time must be on a weekday (Monday–Friday)")
	}

	// Ensure submission is within business hours
	if submitTime.Hour() < workDayStart || submitTime.Hour() >= workDayEnd {
		return submitTime, fmt.Errorf("submit time must be within working hours (9 AM - 5 PM)")
	}

	// Validate turnaroundHours
	if turnaroundHours <= 0 {
		return submitTime, fmt.Errorf("turnaroundHours must be a positive integer")
	}

	remainingHours := turnaroundHours
	currentTime := submitTime

	for remainingHours > 0 {
		hoursLeftToday := workDayEnd - currentTime.Hour()

		if remainingHours <= hoursLeftToday {
			// If turnaround hours fit in the same day
			currentTime = currentTime.Add(time.Duration(remainingHours) * time.Hour)
			break
		} else {
			// Use remaining hours of the day, then move to the next working day
			remainingHours -= hoursLeftToday
			currentTime = nextWorkingDay(currentTime)
			currentTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), workDayStart, 0, 0, 0, currentTime.Location())
		}
	}

	return currentTime, nil
}

// nextWorkingDay moves to the next working day (Monday–Friday only)
func nextWorkingDay(date time.Time) time.Time {
	for {
		date = date.AddDate(0, 0, 1) // Move to next day

		// Skip weekends
		if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
			continue
		}

		break // Return first valid working day
	}

	return date
}

// HandleDueDateRequest handles the HTTP request for calculating the due date
func HandleDueDateRequest(w http.ResponseWriter, r *http.Request) {
	var req DueDateRequest
	decoder := json.NewDecoder(r.Body)

	// Decode the JSON request body
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse the submitTime
	layout := "2006-01-02 15:04"
	submitTime, err := time.Parse(layout, req.SubmitTime)
	if err != nil {
		http.Error(w, "Invalid submitTime format. Use YYYY-MM-DD HH:MM.", http.StatusBadRequest)
		return
	}

	// Calculate the due date
	dueDate, err := CalculateDueDate(submitTime, req.TurnaroundHours)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Prepare the response
	response := DueDateResponse{
		DueDate: dueDate.Format("2006-01-02 15:04"),
	}

	// Set response header and send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Start the HTTP server
func main() {
	http.HandleFunc("/calculate-due-date", HandleDueDateRequest)

	port := ":8080"
	fmt.Println("Server is running on port 8080")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
