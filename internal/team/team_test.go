package team

import (
	"testing"

	"github.com/logan-connolly/teammate/internal/team/services"
	"github.com/matryer/is"
)

func TestMemoryApplication(t *testing.T) {
	t.Run("Init team app", func(t *testing.T) {
		is := is.New(t)
		_, err := NewMemoryApplication()
		is.NoErr(err)
	})

	t.Run("Init failure due to bad roster service config", func(t *testing.T) {
		is := is.New(t)
		WithInvalidMemoryConfig := func() services.RosterConfiguration {
			return func(s *services.RosterService) error {
				return services.ErrInvalidRosterConfig
			}
		}
		originalMemoryConfigs := memoryConfigs
		memoryConfigs = []services.RosterConfiguration{WithInvalidMemoryConfig()}

		_, err := NewMemoryApplication()

		is.Equal(err, services.ErrInvalidRosterConfig)

		// clean up memoryConfigs
		memoryConfigs = originalMemoryConfigs
	})
}
