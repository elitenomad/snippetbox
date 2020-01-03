package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T)  {
	/*
		Create a table of tests
	 */
	tests := []struct {
		name string
		time time.Time
		expected string
	} {
		{
			name: "UTC",
			time:  time.Date(2020, 1, 3, 14, 00, 0, 651387237, time.UTC),
			expected: "03 Jan 2020 at 14:00",
		},
		{
			name: "Empty",
			time:  time.Time{},
			expected: "",
		},
		{
			name: "IST",
			time:  time.Date(2020, 1, 3, 8, 00, 0, 651387237, time.FixedZone("IST", 5.5 * 60 * 60)),
			expected: "03 Jan 2020 at 02:30",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expected := humanDate(test.time)

			if expected != test.expected {
				t.Errorf("Want %q, got %q", test.expected, expected)
			}
		})
	}
}
