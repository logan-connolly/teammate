package team

import (
	"testing"

	"github.com/google/uuid"
)

var (
	exampleUUID = uuid.MustParse("f55e93f8-c952-11ed-afa1-0242ac120002")
	exampleName = "Virginia"
)

func TestTeam_NewTeam(t *testing.T) {
	type testCase struct {
		test        string
		name        string
		expectedErr error
	}
	testCases := []testCase{
		{
			test:        "Empty name validation",
			name:        "",
			expectedErr: ErrInvalidGroup,
		},
		{
			test:        "Valid name",
			name:        exampleName,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			_, err := NewTeam(tc.name)
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestTeam_NewTeamEvents(t *testing.T) {
	type testCase struct {
		test                string
		events              []Event
		expectedIsActivated bool
	}
	testCases := []testCase{
		{
			test: "Team registered",
			events: []Event{
				&TeamRegistered{ID: exampleUUID, Name: exampleName},
			},
			expectedIsActivated: true,
		},
		{
			test: "Team registered and deactivated",
			events: []Event{
				&TeamRegistered{ID: exampleUUID, Name: exampleName},
				&TeamDeactivated{ID: exampleUUID},
			},
			expectedIsActivated: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			team := NewTeamFromEvents(tc.events)
			if team.GetID() != exampleUUID {
				t.Errorf("Expected %v, got %v", exampleUUID, team.GetID())
			}
			if team.GetName() != exampleName {
				t.Errorf("Expected %v, got %v", exampleName, team.GetName())
			}
			if team.IsActivated() != tc.expectedIsActivated {
				t.Errorf("Expected %v, got %v", tc.expectedIsActivated, team.IsActivated())
			}
		})
	}
}

func TestTeam_Activate(t *testing.T) {
	type testCase struct {
		test        string
		team        *Team
		expectedErr error
	}
	testCases := []testCase{
		{
			test: "Activate active team",
			team: NewTeamFromEvents([]Event{
				&TeamRegistered{ID: exampleUUID, Name: exampleName},
			}),
			expectedErr: ErrTeamAlreadyActivated,
		},
		{
			test: "Activate deactivated team",
			team: NewTeamFromEvents([]Event{
				&TeamRegistered{ID: exampleUUID, Name: exampleName},
				&TeamDeactivated{ID: exampleUUID},
			}),
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := tc.team.Activate()
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
			if !tc.team.IsActivated() {
				t.Fatal("team should always be activated in these cases.")
			}
		})
	}
}

func TestTeam_Deactivate(t *testing.T) {
	type testCase struct {
		test        string
		team        *Team
		expectedErr error
	}
	testCases := []testCase{
		{
			test: "Deactivate active team",
			team: NewTeamFromEvents([]Event{
				&TeamRegistered{ID: exampleUUID, Name: exampleName},
			}),
			expectedErr: nil,
		},
		{
			test: "Deactivate deactivated team",
			team: NewTeamFromEvents([]Event{
				&TeamRegistered{ID: exampleUUID, Name: exampleName},
				&TeamDeactivated{ID: exampleUUID},
			}),
			expectedErr: ErrTeamAlreadyDeactivated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := tc.team.Deactivate()
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
			if tc.team.IsActivated() {
				t.Fatal("team should always be deactivated in these cases.")
			}
		})
	}
}

func TestTeam_Events(t *testing.T) {
	t.Run("Event log is populated", func(t *testing.T) {
		team, err := NewTeam(exampleName)
		if err != nil {
			t.Fatalf("Did not expect an error: %v", err)
		}
		team.Deactivate()
		team.Activate()

		want := 3
		got := len(team.Events())

		if want != got {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
}

func TestTeam_Version(t *testing.T) {
	t.Run("Version is properly updated", func(t *testing.T) {
		team := NewTeamFromEvents([]Event{
			&TeamRegistered{ID: exampleUUID, Name: exampleName},
			&TeamDeactivated{ID: exampleUUID},
			&TeamActivated{ID: exampleUUID},
		})

		want := 3
		got := team.Version()

		if want != got {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
}
