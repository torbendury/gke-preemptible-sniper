// Package timing provides a simple way to handle allowlisted / blocklisted timeslots and check if a time is within allowed timeslots.
// Users can provide a list of timeslots in the format "HH:MM-HH:MM" as well as for allowlisted and blocklisted timeslots.
// The package provides a function to check if a time is within the allowed timeslots.
package timing

import (
	"time"

	"math/rand"
)

// TimeSlot represents a timeslot in the format "HH:MM-HH:MM"
type TimeSlot struct {
	Start time.Time
	End   time.Time
}

// TimeSlots represents a list of timeslots
type TimeSlots []TimeSlot

// IsTimeAllowed checks if a time is within the allowed timeslots
func (ts TimeSlots) IsTimeAllowed(t time.Time) bool {
	// Extract the hours and minutes from the provided time
	hour, minute, _ := t.Clock()

	for _, slot := range ts {
		// Extract the hours and minutes from the start and end times of the slot
		startHour, startMinute, _ := slot.Start.Clock()
		endHour, endMinute, _ := slot.End.Clock()

		// Create new time.Time objects for comparison
		startTime := time.Date(0, 1, 1, startHour, startMinute, 0, 0, time.UTC)
		endTime := time.Date(0, 1, 1, endHour, endMinute, 0, 0, time.UTC)
		currentTime := time.Date(0, 1, 1, hour, minute, 0, 0, time.UTC)

		if currentTime.After(startTime) && currentTime.Before(endTime) {
			return true
		}
	}
	return false
}

// IsTimeBlocked checks if a time is within the blocklisted timeslots
func (ts TimeSlots) IsTimeBlocked(t time.Time) bool {
	hour, minute, _ := t.Clock()

	for _, slot := range ts {
		startHour, startMinute, _ := slot.Start.Clock()
		endHour, endMinute, _ := slot.End.Clock()

		startTime := time.Date(0, 1, 1, startHour, startMinute, 0, 0, time.UTC)
		endTime := time.Date(0, 1, 1, endHour, endMinute, 0, 0, time.UTC)
		currentTime := time.Date(0, 1, 1, hour, minute, 0, 0, time.UTC)

		if currentTime.After(startTime) && currentTime.Before(endTime) {
			return true
		}
	}
	return false
}

// ParseTimeSlot parses a string in the format "HH:MM-HH:MM" to a TimeSlot
func ParseTimeSlot(s string) (TimeSlot, error) {
	var slot TimeSlot
	var err error
	startTime, err := time.Parse("15:04", s[:5])
	if err != nil {
		return slot, err
	}
	endTime, err := time.Parse("15:04", s[6:])
	if err != nil {
		return slot, err
	}
	slot.Start = startTime
	slot.End = endTime
	return slot, nil
}

// ParseTimeSlots parses a list of strings in the format "HH:MM-HH:MM" to a list of TimeSlots
func ParseTimeSlots(slots []string) (TimeSlots, error) {
	var ts TimeSlots
	for _, slot := range slots {
		t, err := ParseTimeSlot(slot)
		if err != nil {
			return ts, err
		}
		ts = append(ts, t)
	}
	return ts, nil
}

// ParseTime parses a string in the format "HH:MM" to a time.Time
func ParseTime(s string) (time.Time, error) {
	return time.Parse("15:04", s)
}

// ParseTimes parses a list of strings in the format "HH:MM" to a list of time.Time
func ParseTimes(times []string) ([]time.Time, error) {
	var ts []time.Time
	for _, t := range times {
		time, err := ParseTime(t)
		if err != nil {
			return ts, err
		}
		ts = append(ts, time)
	}
	return ts, nil
}

// CreateAllowedTime creates an allowed time.Time which is in an allowed TimeSlot and outside a blocklisted TimeSlot
func CreateAllowedTime(allowed TimeSlots, blocked TimeSlots) (time.Time, error) {
	res := time.Now().Add(time.Hour * 3)                                                                    // Start at least 3 hours from now
	res = res.Add(time.Duration(rand.Intn(18)) * time.Hour).Add(time.Duration(rand.Intn(45)) * time.Minute) // Add random hours and minutes
	for {
		if allowed.IsTimeAllowed(res) && !blocked.IsTimeBlocked(res) {
			return res, nil
		}
		res = res.Add(10 * time.Minute) // search in 10 minute intervals
	}
}
