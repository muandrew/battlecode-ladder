package rds

import (
	"github.com/muandrew/battlecode-ladder/models"
)

type Match struct {
	Uuid        string
	BotUuids    []string
	Status      *models.BuildStatus
	Competition string
}

func CreateMatch(match *models.Match) *Match {
	uuids := make([]string, len(match.Bots))
	for i, bot := range match.Bots {
		uuids[i] = bot.Uuid
	}
	return &Match{
		match.Uuid,
		uuids,
		match.Status,
		match.Competition,
	}
}
