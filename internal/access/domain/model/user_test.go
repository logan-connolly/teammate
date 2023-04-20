package model

import (
	"testing"

	"git.sr.ht/~loges/teammate/internal/access/domain/event"
	"git.sr.ht/~loges/teammate/internal/entity"
	"github.com/google/uuid"
	"github.com/matryer/is"
)

var (
	exampleUUID     = uuid.MustParse("f55e93f8-c952-11ed-afa1-0242ac120002")
	exampleName     = "Mike Ditka"
	exampleEmail    = "ditka@teammate.com"
	anotherName     = "Joe Gibbs"
	anotherEmail    = "gibbs@teammate.com"
	userRegistered  = &event.UserRegistered{ID: exampleUUID, Name: exampleName, Email: exampleEmail}
	userDeactivated = &event.UserDeactivated{ID: exampleUUID}
)

func TestUser_NewUser(t *testing.T) {
	testCases := []struct {
		test        string
		person      *entity.Person
		email       string
		expectedErr error
	}{
		{"Empty name validation", &entity.Person{ID: exampleUUID, Name: ""}, exampleEmail, ErrInputIsEmpty},
		{"Empty email validation", &entity.Person{ID: exampleUUID, Name: exampleName}, "", ErrInputIsEmpty},
		{"Valid email and name", &entity.Person{ID: exampleUUID, Name: exampleName}, exampleEmail, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			_, err := NewUser(tc.person, tc.email)
			is.Equal(err, tc.expectedErr)
		})
	}
}

func TestUser_NewUserEvents(t *testing.T) {
	testCases := []struct {
		test                string
		events              []event.Event
		expectedIsActivated bool
	}{
		{"User registered", []event.Event{userRegistered}, true},
		{"User registered and deactivated", []event.Event{userRegistered, userDeactivated}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			u := NewUserFromEvents(tc.events)
			is.Equal(u.GetID(), exampleUUID)
			is.Equal(u.GetName(), exampleName)
			is.Equal(u.GetEmail(), exampleEmail)
			is.Equal(u.IsActivated(), tc.expectedIsActivated)
		})
	}
}

func TestUser_UpdateName(t *testing.T) {
	testCases := []struct {
		test        string
		user        *User
		newName     string
		expectedErr error
	}{
		{
			"Update name",
			NewUserFromEvents([]event.Event{userRegistered}),
			anotherName,
			nil,
		},
		{
			"Update name with same name",
			NewUserFromEvents([]event.Event{userRegistered, userDeactivated}),
			exampleName,
			ErrUserUpdateFailed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			err := tc.user.UpdateName(tc.newName)
			is.Equal(err, tc.expectedErr)
		})
	}
}

func TestUser_UpdateEmail(t *testing.T) {
	testCases := []struct {
		test        string
		user        *User
		newEmail    string
		expectedErr error
	}{
		{
			"Update email",
			NewUserFromEvents([]event.Event{userRegistered}),
			anotherEmail,
			nil,
		},
		{
			"Update email with same email",
			NewUserFromEvents([]event.Event{userRegistered, userDeactivated}),
			exampleEmail,
			ErrUserUpdateFailed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			err := tc.user.UpdateEmail(tc.newEmail)
			is.Equal(err, tc.expectedErr)
		})
	}
}

func TestUser_Activate(t *testing.T) {
	testCases := []struct {
		test        string
		user        *User
		expectedErr error
	}{
		{
			"Activate active user",
			NewUserFromEvents([]event.Event{userRegistered}),
			ErrUserUpdateFailed,
		},
		{
			"Activate deactivated user",
			NewUserFromEvents([]event.Event{userRegistered, userDeactivated}),
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			err := tc.user.Activate()
			is.Equal(err, tc.expectedErr)
			is.True(tc.user.IsActivated())
		})
	}
}

func TestUser_Deactivate(t *testing.T) {
	testCases := []struct {
		test        string
		user        *User
		expectedErr error
	}{
		{
			"Deactivate active user",
			NewUserFromEvents([]event.Event{userRegistered}),
			nil,
		},
		{
			"Deactivate deactivated user",
			NewUserFromEvents([]event.Event{userRegistered, userDeactivated}),
			ErrUserUpdateFailed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			err := tc.user.Deactivate()
			is.Equal(err, tc.expectedErr)
			is.Equal(tc.user.IsActivated(), false)
		})
	}
}

func TestUser_Events(t *testing.T) {
	t.Run("Event log is populated", func(t *testing.T) {
		is := is.New(t)
		u, err := NewUser(&entity.Person{ID: exampleUUID, Name: exampleName}, exampleEmail)
		if err != nil {
			t.Fatalf("Did not expect an error: %v", err)
		}
		u.Deactivate()
		u.Activate()

		is.Equal(len(u.Events()), 3)
	})
}

func TestUser_Version(t *testing.T) {
	t.Run("Version is properly updated", func(t *testing.T) {
		is := is.New(t)
		u := NewUserFromEvents([]event.Event{userRegistered, userDeactivated})
		is.Equal(u.Version(), 2)
	})
}
