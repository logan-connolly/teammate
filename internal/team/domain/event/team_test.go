package event

import "testing"

func TestTeamEvent(t *testing.T) {
	type testCase struct {
		test     string
		event    Event
		expected string
	}

	testCases := []testCase{
		{
			test:     "TeamCreated event name",
			event:    &TeamCreated{},
			expected: "TeamCreated",
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
		{
			test:     "PlayerAssignedToTeam event name",
			event:    &PlayerAssignedToTeam{},
			expected: "PlayerAssignedToTeam",
		},
		{
			test:     "PlayerUnassignedFromTeam event name",
			event:    &PlayerUnassignedFromTeam{},
			expected: "PlayerUnassignedFromTeam",
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
