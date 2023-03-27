package team

import "testing"

func TestEvent(t *testing.T) {
	type testCase struct {
		test     string
		event    Event
		expected string
	}

	testCases := []testCase{
		{
			test:     "TeamRegistered event name",
			event:    &TeamRegistered{},
			expected: "TeamRegistered",
		},
		{
			test:     "TeamActivated event name",
			event:    &TeamActivated{},
			expected: "TeamActivated",
		},
		{
			test:     "TeamDeactivated event name",
			event:    &TeamDeactivated{},
			expected: "TeamDeactivated",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			got := tc.event.eventName()
			if got != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, got)
			}
		})
	}
}
