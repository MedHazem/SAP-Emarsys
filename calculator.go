package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const (
	workDayStart = 9
	workDayEnd   = 17
	hoursPerDay  = workDayEnd - workDayStart
)

// CalculateDueDate calculates the due date based on business hours (Mon–Fri, 9 AM–5 PM)
func CalculateDueDate(submitTime time.Time, turnaroundHours int) time.Time {
	// Ensure submission is within business hours
	if submitTime.Hour() < workDayStart || submitTime.Hour() >= workDayEnd {
		fmt.Println("Error: Submit time must be within working hours (9 AM - 5 PM).")
		return submitTime
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

	return currentTime
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

// getSubmitTime prompts the user for the submission time and returns it
func getSubmitTime() time.Time {
	layout := "2006-01-02 15:04"
	var submitTime time.Time
	var inputTime string
	var err error
	scanner := bufio.NewScanner(os.Stdin)

	// Loop until a valid date format is entered
	for {
		fmt.Print("Enter submission time (YYYY-MM-DD HH:MM): ")
		scanner.Scan()
		inputTime = scanner.Text()

		// Parse input date-time
		submitTime, err = time.Parse(layout, inputTime)
		if err == nil {
			// Check if the entered day is a weekend
			if submitTime.Weekday() == time.Saturday || submitTime.Weekday() == time.Sunday {
				fmt.Println("Error: The entered date is a weekend. Please enter a working day (Monday to Friday).")
				continue
			}

			// Check if the entered time is within working hours (9 AM - 5 PM)
			if submitTime.Hour() < workDayStart || submitTime.Hour() >= workDayEnd {
				fmt.Println("Error: Submit time must be within working hours (9 AM - 5 PM). Please enter a valid working time.")
				continue
			}
			break
		}
		fmt.Println("Invalid date format. Please use YYYY-MM-DD HH:MM.")
	}
	return submitTime
}

// getTurnaroundHours prompts the user for the turnaround hours and returns it
func getTurnaroundHours() int {
	var turnaroundHours int
	for {
		fmt.Print("Enter turnaround hours (business hours only): ")
		_, err := fmt.Scanln(&turnaroundHours)
		if err == nil && turnaroundHours > 0 {
			break
		}
		fmt.Println("Invalid input. Please enter a positive integer for turnaround hours.")
	}
	return turnaroundHours
}

func main() {
	// Get submission time and turnaround hours from user input
	submitTime := getSubmitTime()
	turnaroundHours := getTurnaroundHours()

	// Calculate the due date
	dueDate := CalculateDueDate(submitTime, turnaroundHours)

	// Print the due date
	fmt.Println("Due Date:", dueDate.Format("2006-01-02 15:04"))
}
