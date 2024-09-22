package timing

import (
	"testing"
	"time"
)

func TestParseTimeSlot(t *testing.T) {
	tests := []struct {
		input    string
		expected TimeSlot
		hasError bool
	}{
		{
			input: "09:00-10:00",
			expected: TimeSlot{
				Start: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
				End:   time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			hasError: false,
		},
		{
			input:    "invalid",
			expected: TimeSlot{},
			hasError: true,
		},
		{
			input:    "09:00-invalid",
			expected: TimeSlot{},
			hasError: true,
		},
		{
			input:    "invalid-10:00",
			expected: TimeSlot{},
			hasError: true,
		},
	}

	for _, test := range tests {
		result, err := ParseTimeSlot(test.input)
		if (err != nil) != test.hasError {
			t.Errorf("ParseTimeSlot(%s) error = %v, expected error = %v", test.input, err, test.hasError)
		}
		if !test.hasError && (result.Start != test.expected.Start || result.End != test.expected.End) {
			t.Errorf("ParseTimeSlot(%s) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestParseTimeSlots(t *testing.T) {
	tests := []struct {
		input    []string
		expected TimeSlots
		hasError bool
	}{
		{
			input: []string{"09:00-10:00", "11:00-12:00"},
			expected: TimeSlots{
				{
					Start: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
					End:   time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
				},
				{
					Start: time.Date(0, 1, 1, 11, 0, 0, 0, time.UTC),
					End:   time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
				},
			},
			hasError: false,
		},
		{
			input:    []string{"invalid", "11:00-12:00"},
			expected: TimeSlots{},
			hasError: true,
		},
		{
			input:    []string{"09:00-10:00", "invalid"},
			expected: TimeSlots{},
			hasError: true,
		},
		{
			input:    []string{"09:00-10:00", "11:00-invalid"},
			expected: TimeSlots{},
			hasError: true,
		},
	}

	for _, test := range tests {
		result, err := ParseTimeSlots(test.input)
		if (err != nil) != test.hasError {
			t.Errorf("ParseTimeSlots(%v) error = %v, expected error = %v", test.input, err, test.hasError)
		}
		if !test.hasError {
			for i, slot := range result {
				if slot.Start != test.expected[i].Start || slot.End != test.expected[i].End {
					t.Errorf("ParseTimeSlots(%v) = %v, expected %v", test.input, result, test.expected)
					break
				}
			}
		}
	}
}

func TestParseTime(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Time
		hasError bool
	}{
		{
			input:    "09:00",
			expected: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
			hasError: false,
		},
		{
			input:    "invalid",
			expected: time.Time{},
			hasError: true,
		},
		{
			input:    "25:00",
			expected: time.Time{},
			hasError: true,
		},
		{
			input:    "09:60",
			expected: time.Time{},
			hasError: true,
		},
	}

	for _, test := range tests {
		result, err := ParseTime(test.input)
		if (err != nil) != test.hasError {
			t.Errorf("ParseTime(%s) error = %v, expected error = %v", test.input, err, test.hasError)
		}
		if !test.hasError && !result.Equal(test.expected) {
			t.Errorf("ParseTime(%s) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestParseTimes(t *testing.T) {
	tests := []struct {
		input    []string
		expected []time.Time
		hasError bool
	}{
		{
			input: []string{"09:00", "10:00"},
			expected: []time.Time{
				time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
				time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			hasError: false,
		},
		{
			input:    []string{"invalid", "10:00"},
			expected: []time.Time{},
			hasError: true,
		},
		{
			input:    []string{"09:00", "invalid"},
			expected: []time.Time{},
			hasError: true,
		},
		{
			input:    []string{"25:00", "10:00"},
			expected: []time.Time{},
			hasError: true,
		},
	}

	for _, test := range tests {
		result, err := ParseTimes(test.input)
		if (err != nil) != test.hasError {
			t.Errorf("ParseTimes(%v) error = %v, expected error = %v", test.input, err, test.hasError)
		}
		if !test.hasError {
			for i, time := range result {
				if !time.Equal(test.expected[i]) {
					t.Errorf("ParseTimes(%v) = %v, expected %v", test.input, result, test.expected)
					break
				}
			}
		}
	}
}

func TestIsTimeBlocked(t *testing.T) {
	tests := []struct {
		slots    TimeSlots
		time     time.Time
		expected bool
	}{
		{
			slots: TimeSlots{
				{
					Start: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
					End:   time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
				},
			},
			time:     time.Date(0, 1, 1, 9, 30, 0, 0, time.UTC),
			expected: true,
		},
		{
			slots: TimeSlots{
				{
					Start: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
					End:   time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
				},
			},
			time:     time.Date(0, 1, 1, 10, 30, 0, 0, time.UTC),
			expected: false,
		},
		{
			slots: TimeSlots{
				{
					Start: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
					End:   time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
				},
				{
					Start: time.Date(0, 1, 1, 11, 0, 0, 0, time.UTC),
					End:   time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
				},
			},
			time:     time.Date(0, 1, 1, 11, 30, 0, 0, time.UTC),
			expected: true,
		},
		{
			slots: TimeSlots{
				{
					Start: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
					End:   time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
				},
				{
					Start: time.Date(0, 1, 1, 11, 0, 0, 0, time.UTC),
					End:   time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
				},
			},
			time:     time.Date(0, 1, 1, 10, 30, 0, 0, time.UTC),
			expected: false,
		},
	}

	for _, test := range tests {
		result := test.slots.IsTimeBlocked(test.time)
		if result != test.expected {
			t.Errorf("IsTimeBlocked(%v, %v) = %v, expected %v", test.slots, test.time, result, test.expected)
		}
	}
}

func TestIsTimeAllowed(t *testing.T) {
	tests := []struct {
		slots    TimeSlots
		time     time.Time
		expected bool
	}{
		{
			slots: TimeSlots{
				{
					Start: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
					End:   time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
				},
			},
			time:     time.Date(0, 1, 1, 9, 30, 0, 0, time.UTC),
			expected: true,
		},
		{
			slots: TimeSlots{
				{
					Start: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
					End:   time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
				},
			},
			time:     time.Date(0, 1, 1, 10, 30, 0, 0, time.UTC),
			expected: false,
		},
		{
			slots: TimeSlots{
				{
					Start: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
					End:   time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
				},
				{
					Start: time.Date(0, 1, 1, 11, 0, 0, 0, time.UTC),
					End:   time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
				},
			},
			time:     time.Date(0, 1, 1, 11, 30, 0, 0, time.UTC),
			expected: true,
		},
		{
			slots: TimeSlots{
				{
					Start: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
					End:   time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
				},
				{
					Start: time.Date(0, 1, 1, 11, 0, 0, 0, time.UTC),
					End:   time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
				},
			},
			time:     time.Date(0, 1, 1, 10, 30, 0, 0, time.UTC),
			expected: false,
		},
	}

	for _, test := range tests {
		result := test.slots.IsTimeAllowed(test.time)
		if result != test.expected {
			t.Errorf("IsTimeAllowed(%v, %v) = %v, expected %v", test.slots, test.time, result, test.expected)
		}
	}
}

func helperTime(s string) time.Time {
	t, _ := time.Parse("15:04", s)
	return t
}

func TestCreateAllowedTime(t *testing.T) {
	tests := []struct {
		allowed  TimeSlots
		blocked  TimeSlots
		expected bool
	}{
		{
			allowed: TimeSlots{
				{
					Start: helperTime("09:00"),
					End:   helperTime("10:00"),
				},
			},
			blocked:  TimeSlots{},
			expected: true,
		},
		{
			allowed: TimeSlots{
				{
					Start: helperTime("09:00"),
					End:   helperTime("10:00"),
				},
			},
			blocked: TimeSlots{
				{
					Start: helperTime("10:00"),
					End:   helperTime("11:00"),
				},
			},
			expected: true,
		},
		{
			allowed: TimeSlots{
				{
					Start: helperTime("09:00"),
					End:   helperTime("10:00"),
				},
			},
			blocked: TimeSlots{
				{
					Start: helperTime("09:30"),
					End:   helperTime("10:30"),
				},
			},
			expected: true,
		},
	}

	for _, test := range tests {
		result, err := CreateAllowedTime(test.allowed, test.blocked)
		if (err == nil) != test.expected {
			t.Errorf("CreateAllowedTime(%v, %v) error = %v, expected error = %v", test.allowed, test.blocked, err, test.expected)
		}
		if test.expected && !test.allowed.IsTimeAllowed(result) {
			t.Errorf("CreateAllowedTime(%v, %v) = %v, expected time within allowed slots", test.allowed, test.blocked, result)
		}
		if test.expected && test.blocked.IsTimeBlocked(result) {
			t.Errorf("CreateAllowedTime(%v, %v) = %v, expected time not within blocked slots", test.allowed, test.blocked, result)
		}
	}
}
