package access

import (
	"errors"
	"testing"

	"github.com/logan-connolly/teammate/internal/access/services"
)

func TestMemoryApplication(t *testing.T) {
	t.Run("Init access app", func(t *testing.T) {
		_, err := NewMemoryApplication()

		if !errors.Is(err, nil) {
			t.Errorf("Wanted %v, but got %v", nil, err)
		}
	})

	t.Run("Init failure due to bad registration service config", func(t *testing.T) {
		WithInvalidMemoryConfig := func() services.RegistrationConfiguration {
			return func(s *services.RegistrationService) error {
				return services.ErrInvalidRegistrationConfig
			}
		}
		originalMemoryConfigs := memoryConfigs
		memoryConfigs = []services.RegistrationConfiguration{WithInvalidMemoryConfig()}

		_, err := NewMemoryApplication()

		if !errors.Is(err, services.ErrInvalidRegistrationConfig) {
			t.Errorf("Wanted %v, but got %v", services.ErrInvalidRegistrationConfig, err)
		}

		// clean up memoryConfigs
		memoryConfigs = originalMemoryConfigs
	})
}
