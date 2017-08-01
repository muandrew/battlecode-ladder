package rds

import (
	"github.com/muandrew/battlecode-ladder/models"
)

type Match struct {
	Uuid        string
	BotUuids    []string
	MapUuid     string
	Winner      int
	Status      *models.BuildStatus
	Competition models.Competition
}

func CreateMatch(match *models.Match) *Match {
	uuids := make([]string, len(match.Bots))
	for i, bot := range match.Bots {
		uuids[i] = bot.Uuid
	}
	return &Match{
		match.Uuid,
		uuids,
		match.Uuid,
		match.Winner,
		match.Status,
		match.Competition,
	}
}
