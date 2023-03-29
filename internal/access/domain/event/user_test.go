package event

import "testing"

func TestEvent(t *testing.T) {
	type testCase struct {
		test     string
		event    Event
		expected string
	}

	testCases := []testCase{
		{
			test:     "UserRegistered event name",
			event:    &UserRegistered{},
			expected: "UserRegistered",
		},
		{
			test:     "UserNameChanged event name",
			event:    &UserNameChanged{},
			expected: "UserNameChanged",
		},
		{
			test:     "UserEmailChanged event name",
			event:    &UserEmailChanged{},
			expected: "UserEmailChanged",
		},
		{
			test:     "UserActivated event name",
			event:    &UserActivated{},
			expected: "UserActivated",
		},
		{
			test:     "UserDeactivated event name",
			event:    &UserDeactivated{},
			expected: "UserDeactivated",
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
