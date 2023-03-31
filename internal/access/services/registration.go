package services

import (
	"errors"

	"github.com/logan-connolly/teammate/internal/access/domain/repository"
	"github.com/logan-connolly/teammate/internal/access/infrastructure/memory"
)

var ErrInvalidRegistrationConfig = errors.New("services: invalid registration configuration")

// RegistrationConfiguration is a function that modifies the service.
type RegistrationConfiguration func(s *RegistrationService) error

// RegistrationService is a implementation of the RegistrationService.
type RegistrationService struct {
	users repository.UserRepository
}

// NewRegistrationService accepts configs and returns a new service.
func NewRegistrationService(cfgs ...RegistrationConfiguration) (*RegistrationService, error) {
	s := &RegistrationService{}

	for _, cfg := range cfgs {
		err := cfg(s)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

// WithUserRepository applies a given user repository to the service.
func WithUserRepository(r repository.UserRepository) RegistrationConfiguration {
	return func(s *RegistrationService) error {
		s.users = r
		return nil
	}
}

// WithMemoryUserRepository applies a memory user repository to the service.
func WithMemoryUserRepository() RegistrationConfiguration {
	r := memory.NewMemoryUserRepository()
	return WithUserRepository(r)
}
