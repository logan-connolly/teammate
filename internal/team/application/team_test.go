package application

import (
	"testing"

	"github.com/logan-connolly/teammate/internal/team/application/services"
	"github.com/matryer/is"
)

func withInvalidMemoryConfig() services.RosterConfiguration {
	return func(s *services.RosterService) error {
		return services.ErrInvalidRosterConfig
	}
}

func TestTeamApplication(t *testing.T) {
	t.Run("Init team app", func(t *testing.T) {
		is := is.New(t)
		_, err := NewTeamApplication()
		is.NoErr(err)
	})

	t.Run("Init failure due to bad roster service config", func(t *testing.T) {
		is := is.New(t)
		originalMemoryConfigs := memoryConfigs
		memoryConfigs = []services.RosterConfiguration{withInvalidMemoryConfig()}

		_, err := NewTeamApplication()

		is.Equal(err, services.ErrInvalidRosterConfig)

		// clean up memoryConfigs
		memoryConfigs = originalMemoryConfigs
	})
}
