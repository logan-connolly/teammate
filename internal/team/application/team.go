package application

import (
	"github.com/logan-connolly/teammate/internal/team/application/services"
)

var memoryConfigs = []services.RosterConfiguration{
	services.WithMemoryPlayerRepository(),
	services.WithMemoryTeamRepository(),
}

// TeamApplication holds all services related to team management.
type TeamApplication struct {
	rosterService services.RosterService
}

// NewTeamApplication intitializes the team application.
func NewTeamApplication() (*TeamApplication, error) {
	rs, err := services.NewRosterService(memoryConfigs...)
	if err != nil {
		return &TeamApplication{}, services.ErrInvalidRosterConfig
	}

	return &TeamApplication{
		rosterService: *rs,
	}, nil
}
