package models

import (
	"errors"
	"github.com/satori/go.uuid"
)

type Match struct {
	Uuid        string
	Bots        []*Bot
	Status      *BuildStatus
	Competition string
}

func CreateMatch(bots []*Bot) (*Match, error) {
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
	return &Match{
		uuid.NewV4().String(),
		bots,
		NewBuildStatus(),
		competition,
	}, nil
}
