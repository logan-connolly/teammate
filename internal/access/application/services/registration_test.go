package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/logan-connolly/teammate/internal/access/domain/model"
	"github.com/logan-connolly/teammate/internal/access/domain/repository"
	"github.com/logan-connolly/teammate/internal/entity"
	"github.com/matryer/is"
)

const (
	name       = "Mark"
	email      = "mark@teammate.com"
	otherName  = "Janet"
	otherEmail = "janet@teammate.com"
)

func badRegistrationConfiguration() RegistrationConfiguration {
	return func(s *RegistrationService) error {
		return ErrInvalidRegistrationConfig
	}
}

func TestNewRegistrationService(t *testing.T) {
	testCases := []struct {
		test        string
		testConfig  func() RegistrationConfiguration
		expectedErr error
	}{
		{"With in-memory user repository", WithMemoryUserRepository, nil},
		{"With bad config", badRegistrationConfiguration, ErrInvalidRegistrationConfig},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			_, err := NewRegistrationService(tc.testConfig())
			is.Equal(err, tc.expectedErr)
		})
	}
}

func TestRegistrationService_RegisterUser(t *testing.T) {
	testCases := []struct {
		test        string
		name        string
		email       string
		expectedErr error
	}{
		{"Email already registered", name, email, repository.ErrUserAlreadyExists},
		{"Email missing", name, "", model.ErrInputIsEmpty},
		{"Name missing", "", email, model.ErrInputIsEmpty},
		{"User successfully registered", otherName, otherEmail, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			is := is.New(t)
			s, _ := NewRegistrationService(WithMemoryUserRepository())
			u, _ := model.NewUser(&entity.Person{ID: uuid.New(), Name: name}, email)
			s.users.Add(u)

			err := s.RegisterUser(tc.name, tc.email)

			is.Equal(err, tc.expectedErr)
		})
	}
}
