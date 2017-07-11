package models

import (
	"github.com/satori/go.uuid"
)

type BcMap struct {
	Uuid        string
	Access      *RAM
	Competition Competition
	Name        UserString
	Description UserString
}

func CreateBcMap(owner *Competitor, competition Competition, name string, description string) (*BcMap, error) {
	uName, err := NewUserString(name, BotMaxName, RegexBlacklist(RegexFilterText))
	if err != nil {
		return nil, err
	}
	uDesc, err := NewUserString(description, BotMaxDescription, RegexBlacklist(RegexFilterText))
	if err != nil {
		return nil, err
	}
	return &BcMap{
		uuid.NewV4().String(),
		CreateRAM(owner),
		competition,
		uName,
		uDesc,
	}, nil
}
