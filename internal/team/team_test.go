package team

import (
	"errors"
	"testing"

	"github.com/logan-connolly/teammate/internal/team/services"
)

func TestMemoryApplication(t *testing.T) {
	t.Run("Init team app", func(t *testing.T) {
		_, err := NewMemoryApplication()

		if !errors.Is(err, nil) {
			t.Errorf("Wanted %v, but got %v", nil, err)
		}
	})

	t.Run("Init failure due to bad roster service config", func(t *testing.T) {
		WithInvalidMemoryConfig := func() services.RosterConfiguration {
			return func(s *services.RosterService) error {
				return services.ErrInvalidRosterConfig
			}
		}
		originalMemoryConfigs := memoryConfigs
		memoryConfigs = []services.RosterConfiguration{WithInvalidMemoryConfig()}

		_, err := NewMemoryApplication()

		if !errors.Is(err, services.ErrInvalidRosterConfig) {
			t.Errorf("Wanted %v, but got %v", services.ErrInvalidRosterConfig, err)
		}

		// clean up memoryConfigs
		memoryConfigs = originalMemoryConfigs
	})
}
