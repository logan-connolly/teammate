package event

import "testing"

func TestPlayerEvent(t *testing.T) {
	type testCase struct {
		test     string
		event    Event
		expected string
	}

	testCases := []testCase{
		{
			test:     "PlayerRegistered event name",
			event:    &PlayerRegistered{},
			expected: "PlayerRegistered",
		},
		{
			test:     "PlayerActivated event name",
			event:    &PlayerActivated{},
			expected: "PlayerActivated",
		},
		{
			test:     "PlayerDeactivated event name",
			event:    &PlayerDeactivated{},
			expected: "PlayerDeactivated",
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
