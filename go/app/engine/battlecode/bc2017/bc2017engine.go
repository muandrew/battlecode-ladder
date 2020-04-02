package bc2017

import (
	"github.com/muandrew/battlecode-legacy-go/models"
)

//Engine runs battlecode 2017
type Engine struct {
}

//Competition see parent.
func (e *Engine) Competition() models.Competition {
	return models.CompetitionBC17
}

//BattleBot see parent
func (e *Engine) BattleBot() {}

//BuildBot see parent
func (e *Engine) BuildBot() {}
