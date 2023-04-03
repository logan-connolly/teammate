package event

import (
	"testing"

	"github.com/matryer/is"
)

func TestEvent(t *testing.T) {
	testCases := []struct {
		test     string
		event    Event
		expected string
	}{
		{"UserRegistered event name", &UserRegistered{}, "UserRegistered"},
		{"UserNameChanged event name", &UserNameChanged{}, "UserNameChanged"},
		{"UserEmailChanged event name", &UserEmailChanged{}, "UserEmailChanged"},
		{"UserActivated event name", &UserActivated{}, "UserActivated"},
		{"UserDeactivated event name", &UserDeactivated{}, "UserDeactivated"},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			is.Equal(tc.event.eventName(), tc.expected)
		})
	}
}
