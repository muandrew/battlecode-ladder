package data

import (
	"github.com/muandrew/battlecode-legacy-go/models"
)

//Match how match is stored
type Match struct {
	UUID        string
	BotUUIDs    []string
	MapUUID     string
	Winner      int
	Status      *models.BuildStatus
	Competition models.Competition
}

//Matches multiple matches
type Matches struct {
	Matches      []*models.Match
	TotalMatches int
}

//CreateMatch creates a new instance
func CreateMatch(match *models.Match) *Match {
	uuids := make([]string, len(match.Bots))
	for i, bot := range match.Bots {
		uuids[i] = bot.UUID
	}
	return &Match{
		match.UUID,
		uuids,
		match.MapUUID,
		match.Winner,
		match.Status,
		match.Competition,
	}
}
