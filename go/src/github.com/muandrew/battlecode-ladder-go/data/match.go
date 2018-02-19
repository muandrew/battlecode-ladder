package data

import (
	"github.com/muandrew/battlecode-ladder-go/models"
)

type Match struct {
	Uuid        string
	BotUuids    []string
	MapUuid     string
	Winner      int
	Status      *models.BuildStatus
	Competition models.Competition
}

type Matches struct {
	Matches      []*models.Match
	TotalMatches int
}

func CreateMatch(match *models.Match) *Match {
	uuids := make([]string, len(match.Bots))
	for i, bot := range match.Bots {
		uuids[i] = bot.Uuid
	}
	return &Match{
		match.Uuid,
		uuids,
		match.MapUuid,
		match.Winner,
		match.Status,
		match.Competition,
	}
}
