package memory

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/logan-connolly/teammate/internal/entity"
	"github.com/logan-connolly/teammate/internal/team/domain/event"
	"github.com/logan-connolly/teammate/internal/team/domain/model"
	"github.com/logan-connolly/teammate/internal/team/domain/repository"
)

var (
	examplePlayerUUID = uuid.MustParse("f55e93f8-c952-11ed-afa1-0242ac120002")
	anotherPlayerUUID = uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479")
	examplePlayerName = "Logan"
	anotherPlayerName = "Emily"
)

func TestMemoryPlayerRepository_Get(t *testing.T) {
	type testCase struct {
		test        string
		person      *entity.Person
		expectedErr error
	}

	testCases := []testCase{
		{
			test:        "No player with this person",
			person:      &entity.Person{ID: anotherPlayerUUID, Name: anotherTeamName},
			expectedErr: repository.ErrPlayerNotFound,
		},
		{
			test:        "Player found",
			person:      &entity.Person{ID: examplePlayerUUID, Name: exampleTeamName},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			repo := NewMemoryPlayerRepository()
			repo.players[examplePlayerUUID] = []event.Event{
				&event.PlayerCreated{ID: examplePlayerUUID, Name: exampleTeamName},
			}

			_, err := repo.Get(tc.person)

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
			test:        "Successfully add a player",
			id:          anotherPlayerUUID,
			name:        anotherPlayerName,
			expectedErr: nil,
		},
		{
			test:        "Player already exists error",
			id:          examplePlayerUUID,
			name:        examplePlayerName,
			expectedErr: repository.ErrPlayerAlreadyExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			r := NewMemoryPlayerRepository()
			p := model.NewPlayerFromEvents([]event.Event{
				&event.PlayerCreated{ID: tc.id, Name: tc.name},
			})
			r.players[examplePlayerUUID] = p.Events()

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
			test:        "Update player",
			register:    true,
			deactivate:  true,
			expectedErr: nil,
		},
		{
			test:        "Player has no changes",
			register:    true,
			deactivate:  false,
			expectedErr: repository.ErrPlayerHasNoUpdates,
		},
		{
			test:        "Player not found",
			register:    false,
			deactivate:  true,
			expectedErr: repository.ErrPlayerNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			r := NewMemoryPlayerRepository()
			p := model.NewPlayerFromEvents([]event.Event{
				&event.PlayerCreated{ID: examplePlayerUUID, Name: exampleTeamName},
			})
			if tc.register {
				r.players[examplePlayerUUID] = p.Events()
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
