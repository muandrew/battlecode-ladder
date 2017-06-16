package rds

import (
	"github.com/muandrew/battlecode-ladder/models"
)

type Match struct {
	Uuid     string
	BotUuids []string
	Status   *models.BuildStatus
}

func CreateMatch(match *models.Match) *Match {
	uuids := make([]string, len(match.Bots))
	for i, bot := range match.Bots {
		uuids[i] = bot.Uuid
	}
	return &Match{
		Uuid:     match.Uuid,
		BotUuids: uuids,
		Status:   match.Status,
	}
}
