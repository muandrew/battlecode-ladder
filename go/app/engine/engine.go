package engine

import (
	"github.com/muandrew/battlecode-legacy-go/models"
)

//Engine abstracts the different types of competitions
type Engine interface {
	Competition() models.Competition
	ActivateAssets()
	BattleBotSetup(
		workerID int,
		workspaceDir string,
		match *models.Match,
	) error
	BattleBotPostProcessing(
		matchPath string,
		match *models.Match,
	) error
	BuildBotSetup(
		workerID int,
		workspaceDir string,
		botUUID string,
	) error
}
