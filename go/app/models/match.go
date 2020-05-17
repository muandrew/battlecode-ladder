package models

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

const (
	//WinnerNone no one wins, its a tie.
	WinnerNone = -1
	//WinnerNeutral if there is a neutral force like nature, it won.
	WinnerNeutral = -2
)

//Match represents a single simulation
type Match struct {
	UUID        string
	Bots        []*Bot
	MapUUID     string
	Winner      int
	Status      *BuildStatus
	Competition Competition
}

//CreateMatch creates a new instance of a Match object.
func CreateMatch(bots []*Bot, bcMap *BcMap) (*Match, error) {
	length := len(bots)
	if length < 2 {
		return nil, errors.New("Can't play with just one bot")
	}
	competition := bots[0].Competition
	for _, bot := range bots {
		if bot == nil {
			return nil, errors.New("Nil bot received")
		}
		if competition != bot.Competition {
			return nil, errors.New("Bots from different competitions can't play with each other")
		}
	}
	mapUUID := ""
	if bcMap != nil {
		mapUUID = bcMap.UUID
	}
	return &Match{
		uuid.NewV4().String(),
		bots,
		mapUUID,
		WinnerNone,
		NewBuildStatus(),
		competition,
	}, nil
}
