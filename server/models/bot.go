package models

import (
	"github.com/satori/go.uuid"
)

const (
	BotMaxName         = 60
	BotMaxDescription  = 140
	BotMaxPackage      = 60
	BotCompetitionBC17 = "bc17"
)

type Bot struct {
	Uuid        string
	Owner       *Competitor
	Package     UserString
	Name        UserString
	Description UserString
	Status      *BuildStatus
	Competition string
}

func CreateBot(owner *Competitor, pkg string, name string, description string, competition string) (*Bot, error) {
	uPkg, err := NewUserString(pkg, BotMaxPackage, RegexBlacklist(RegexFilterPackage))
	if err != nil {
		return nil, err
	}
	uName, err := NewUserString(name, BotMaxName, RegexBlacklist(RegexFilterText))
	if err != nil {
		return nil, err
	}
	uDesc, err := NewUserString(description, BotMaxDescription, RegexBlacklist(RegexFilterText))
	if err != nil {
		return nil, err
	}
	bot := &Bot{
		uuid.NewV4().String(),
		owner,
		uPkg,
		uName,
		uDesc,
		NewBuildStatus(),
		competition,
	}
	return bot, nil
}
