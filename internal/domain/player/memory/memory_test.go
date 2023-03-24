package memory

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/logan-connolly/teammate/internal/domain/player"
)

var (
	exampleUUID = uuid.MustParse("f55e93f8-c952-11ed-afa1-0242ac120002")
	anotherUUID = uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479")
	exampleName = "Logan"
	anotherName = "Emily"
)

func TestMemoryPlayerRepository_Get(t *testing.T) {
	type testCase struct {
		test        string
		id          uuid.UUID
		expectedErr error
	}

	testCases := []testCase{
		{
			test:        "No Player By ID",
			id:          anotherUUID,
			expectedErr: player.ErrPlayerNotFound,
		}, {
			test:        "Player By ID",
			id:          exampleUUID,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			repo := NewMemoryPlayerRepository()
			repo.players[exampleUUID] = []player.Event{
				&player.PlayerRegistered{ID: exampleUUID, Name: exampleName},
			}

			_, err := repo.Get(tc.id)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestMemoryPlayerRepository_Add(t *testing.T) {
	type testCase struct {
		test        string
		id          uuid.UUID
		name        string
		expectedErr error
	}

	testCases := []testCase{
		{
			test:        "Add Player",
			id:          anotherUUID,
			name:        anotherName,
			expectedErr: nil,
		},
		{
			test:        "Player already exists",
			id:          exampleUUID,
			name:        exampleName,
			expectedErr: player.ErrPlayerAlreadyExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			r := NewMemoryPlayerRepository()
			p := player.NewPlayerFromEvents([]player.Event{
				&player.PlayerRegistered{ID: tc.id, Name: tc.name},
			})
			r.players[exampleUUID] = p.Events()

			err := r.Add(p)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestMemoryPlayerRepository_Update(t *testing.T) {
	type testCase struct {
		test        string
		register    bool
		deactivate  bool
		expectedErr error
	}

	testCases := []testCase{
		{
			test:        "Update Player",
			register:    true,
			deactivate:  true,
			expectedErr: nil,
		},
		{
			test:        "Player not updated",
			register:    true,
			deactivate:  false,
			expectedErr: player.ErrPlayerHasNoUpdates,
		},
		{
			test:        "Player not found",
			register:    false,
			deactivate:  true,
			expectedErr: player.ErrPlayerNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			r := NewMemoryPlayerRepository()
			p := player.NewPlayerFromEvents([]player.Event{
				&player.PlayerRegistered{ID: exampleUUID, Name: exampleName},
			})
			if tc.register {
				r.players[exampleUUID] = p.Events()
			}
			if tc.deactivate {
				p.Deactivate()
			}

			err := r.Update(p)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}