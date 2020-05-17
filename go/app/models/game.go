package models

import uuid "github.com/satori/go.uuid"

const (
	GameTypeRoundRobin = "roundRobin"
)

type Game struct {
	UUID        string
	Owner       *Competitor
	Competition Competition
	Type        string
	Name        UserString
	Description UserString
	Status      *BuildStatus
	Bots        []*Bot
	Matches     []*Match
}

type GameRoundRobin struct {
	*Game
}

func CreateGameRoundRobin(
	owner *Competitor,
	competition Competition,
	name string,
	description string,
	bots []*Bot,
	bcMap *BcMap) (*GameRoundRobin, error) {

	n, err := NewUserString(name, BotMaxName)
	if err != nil {
		return nil, err
	}
	d, err := NewUserString(description, BotMaxDescription)
	if err != nil {
		return nil, err
	}
	numBots := len(bots)
	matches := make([]*Match, numBots*numBots-numBots)
	var idx = 0
	for i, a := range bots {
		for j, b := range bots {
			if i == j {
				continue
			} else {
				match, err := CreateMatch([]*Bot{a, b}, bcMap)
				if err != nil {
					return nil, err
				}
				matches[idx] = match
				idx++
			}
		}
	}
	return &GameRoundRobin{
		&Game{
			uuid.NewV4().String(),
			owner,
			competition,
			GameTypeRoundRobin,
			n,
			d,
			NewBuildStatus(),
			bots,
			matches,
		},
	}, nil
}
