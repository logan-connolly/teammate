package player

import (
	"testing"

	"github.com/google/uuid"
)

var (
	exampleUUID = uuid.MustParse("f55e93f8-c952-11ed-afa1-0242ac120002")
	exampleName = "Logan"
)

func TestPlayer_NewPlayer(t *testing.T) {
	type testCase struct {
		test        string
		name        string
		expectedErr error
	}
	testCases := []testCase{
		{
			test:        "Empty name validation",
			name:        "",
			expectedErr: ErrInvalidPerson,
		},
		{
			test:        "Valid name",
			name:        exampleName,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			_, err := NewPlayer(tc.name)
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestPlayer_NewPlayerEvents(t *testing.T) {
	type testCase struct {
		test                string
		events              []Event
		expectedIsActivated bool
	}
	testCases := []testCase{
		{
			test: "Player registered",
			events: []Event{
				&PlayerRegistered{ID: exampleUUID, Name: exampleName},
			},
			expectedIsActivated: true,
		},
		{
			test: "Player registered and deactivated",
			events: []Event{
				&PlayerRegistered{ID: exampleUUID, Name: exampleName},
				&PlayerDeactivated{ID: exampleUUID},
			},
			expectedIsActivated: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			p := NewPlayerFromEvents(tc.events)
			if p.GetID() != exampleUUID {
				t.Errorf("Expected %v, got %v", exampleUUID, p.GetID())
			}
			if p.GetName() != exampleName {
				t.Errorf("Expected %v, got %v", exampleName, p.GetName())
			}
			if p.IsActivated() != tc.expectedIsActivated {
				t.Errorf("Expected %v, got %v", tc.expectedIsActivated, p.IsActivated())
			}
		})
	}
}

func TestPlayer_Activate(t *testing.T) {
	type testCase struct {
		test        string
		player      *Player
		expectedErr error
	}
	testCases := []testCase{
		{
			test: "Activate active player",
			player: NewPlayerFromEvents([]Event{
				&PlayerRegistered{ID: exampleUUID, Name: exampleName},
			}),
			expectedErr: ErrPlayerAlreadyActivated,
		},
		{
			test: "Activate deactivated player",
			player: NewPlayerFromEvents([]Event{
				&PlayerRegistered{ID: exampleUUID, Name: exampleName},
				&PlayerDeactivated{ID: exampleUUID},
			}),
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := tc.player.Activate()
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
			if !tc.player.IsActivated() {
				t.Fatal("player should always be activated in these cases.")
			}
		})
	}
}

func TestPlayer_Deactivate(t *testing.T) {
	type testCase struct {
		test        string
		player      *Player
		expectedErr error
	}
	testCases := []testCase{
		{
			test: "Deactivate active player",
			player: NewPlayerFromEvents([]Event{
				&PlayerRegistered{ID: exampleUUID, Name: exampleName},
			}),
			expectedErr: nil,
		},
		{
			test: "Deactivate deactivated player",
			player: NewPlayerFromEvents([]Event{
				&PlayerRegistered{ID: exampleUUID, Name: exampleName},
				&PlayerDeactivated{ID: exampleUUID},
			}),
			expectedErr: ErrPlayerAlreadyDeactivated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := tc.player.Deactivate()
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
			if tc.player.IsActivated() {
				t.Fatal("player should always be deactivated in these cases.")
			}
		})
	}
}

func TestPlayer_Events(t *testing.T) {
	t.Run("Event log is populated", func(t *testing.T) {
		p, err := NewPlayer(exampleName)
		if err != nil {
			t.Fatalf("Did not expect an error: %v", err)
		}
		p.Deactivate()
		p.Activate()

		want := 3
		got := len(p.Events())

		if want != got {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
}

func TestPlayer_Version(t *testing.T) {
	t.Run("Version is properly updated", func(t *testing.T) {
		p := NewPlayerFromEvents([]Event{
			&PlayerRegistered{ID: exampleUUID, Name: exampleName},
			&PlayerDeactivated{ID: exampleUUID},
			&PlayerActivated{ID: exampleUUID},
		})

		want := 3
		got := p.Version()

		if want != got {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
}
