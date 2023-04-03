package services

import (
	"testing"

	"github.com/matryer/is"
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
