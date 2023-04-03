package event

import (
	"testing"

	"github.com/matryer/is"
)

func TestPlayerEvent(t *testing.T) {
	testCases := []struct {
		test     string
		event    Event
		expected string
	}{
		{"PlayerCreated event name", &PlayerCreated{}, "PlayerCreated"},
		{"PlayerActivated event name", &PlayerActivated{}, "PlayerActivated"},
		{"PlayerDeactivated event name", &PlayerDeactivated{}, "PlayerDeactivated"},
		{"TeamAssignedToPlayer event name", &TeamAssignedToPlayer{}, "TeamAssignedToPlayer"},
		{"TeamUnassignedFromPlayer event name", &TeamUnassignedFromPlayer{}, "TeamUnassignedFromPlayer"},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			is.Equal(tc.event.eventName(), tc.expected)
		})
	}
}
