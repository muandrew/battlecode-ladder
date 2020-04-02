package models

import (
	"errors"
	"github.com/satori/go.uuid"
)

const (
	WinnerNone    = -1
	WinnerNeutral = -2
)

type Match struct {
	Uuid        string
	Bots        []*Bot
	MapUuid     string
	Winner      int
	Status      *BuildStatus
	Competition Competition
}

func CreateMatch(bots []*Bot, bcMap *BcMap) (*Match, error) {
	length := len(bots)
	if length < 2 {
		return nil, errors.New("Can't play with just one bot")
	}
	competition := bots[0].Competition
	for _, bot := range bots {
		if bot == nil {
			return nil, errors.New("Nil bot received.");
		}
		if competition != bot.Competition {
			return nil, errors.New("Bots from different competitions can't play with each other.")
		}
	}
	mapUuid := ""
	if bcMap != nil {
		mapUuid = bcMap.Uuid
	}
	return &Match{
		uuid.NewV4().String(),
		bots,
		mapUuid,
		WinnerNone,
		NewBuildStatus(),
		competition,
	}, nil
}
