package event

import (
	"testing"

	"github.com/matryer/is"
)

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
			is := is.New(t)
			got := tc.event.eventName()
			is.Equal(got, tc.expected)
		})
	}
}
