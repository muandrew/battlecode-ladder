package models

import (
	"github.com/satori/go.uuid"
)

type Match struct {
	Uuid   string
	Bots   []*Bot
	Status *BuildStatus
}

func CreateMatch(bots []*Bot) *Match {
	return &Match{
		Uuid:   uuid.NewV4().String(),
		Bots:   bots,
		Status: NewBuildStatus(),
	}
}
