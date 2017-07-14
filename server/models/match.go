package models

import (
	"errors"
	"github.com/satori/go.uuid"
)

type Match struct {
	Uuid        string
	Bots        []*Bot
	MapUuid     string
	Status      *BuildStatus
	Competition Competition
}

func CreateMatch(bots []*Bot, bcMap *BcMap) (*Match, error) {
	length := len(bots)
	if length < 2 {
		return nil, errors.New("Can't play with just one bot")
	}
	competition := bots[0].Competition
	for i := 1; i < length; i++ {
		if competition != bots[i].Competition {
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
		NewBuildStatus(),
		competition,
	}, nil
}
