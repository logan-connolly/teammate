package event

import (
	"testing"

	"github.com/matryer/is"
)

func TestTeamEvent(t *testing.T) {
	testCases := []struct {
		test     string
		event    Event
		expected string
	}{
		{"TeamCreated event name", &TeamCreated{}, "TeamCreated"},
		{"TeamActivated event name", &TeamActivated{}, "TeamActivated"},
		{"TeamDeactivated event name", &TeamDeactivated{}, "TeamDeactivated"},
		{"PlayerAssignedToTeam event name", &PlayerAssignedToTeam{}, "PlayerAssignedToTeam"},
		{"PlayerUnassignedFromTeam event name", &PlayerUnassignedFromTeam{}, "PlayerUnassignedFromTeam"},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			is.Equal(tc.event.eventName(), tc.expected)
		})
	}
}
