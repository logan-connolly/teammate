package event

import (
	"testing"

	"github.com/matryer/is"
)

func TestPlayerEvent(t *testing.T) {
	type testCase struct {
		test     string
		event    Event
		expected string
	}

	testCases := []testCase{
		{
			test:     "PlayerCreated event name",
			event:    &PlayerCreated{},
			expected: "PlayerCreated",
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
		{
			test:     "TeamAssignedToPlayer event name",
			event:    &TeamAssignedToPlayer{},
			expected: "TeamAssignedToPlayer",
		},
		{
			test:     "TeamUnassignedFromPlayer event name",
			event:    &TeamUnassignedFromPlayer{},
			expected: "TeamUnassignedFromPlayer",
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
