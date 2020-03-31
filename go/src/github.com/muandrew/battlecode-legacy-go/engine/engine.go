package engine

import (
	"github.com/muandrew/battlecode-legacy-go/models"
)

//Engine abstracts the different types of competitions
type Engine interface {
	Competition() models.Competition
	BattleBot()
	BuildBot()
}
