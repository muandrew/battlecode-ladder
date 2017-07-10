package models

import (
	"github.com/satori/go.uuid"
	"errors"
)

const (
	BotMaxName         = 60
	BotMaxDescription  = 140
	BotMaxNote         = 140
	BotMaxPackage      = 60
	BotCompetitionBC17 = "bc17"
)

type Bot struct {
	Uuid        string
	Owner       *Competitor
	Package     UserString
	Note        UserString
	Status      *BuildStatus
	Competition string
}

func CreateBot(owner *Competitor, pkg string, note string, competition string) (*Bot, error) {
	if pkg == "" {
		return nil, errors.New("We need a package to run your bot.")
	}
	uPkg, err := NewUserString(pkg, BotMaxPackage, RegexBlacklist(RegexFilterPackage))
	if err != nil {
		return nil, err
	}
	uNote, err := NewUserString(note, BotMaxNote, RegexBlacklist(RegexFilterText))
	if err != nil {
		return nil, err
	}
	bot := &Bot{
		uuid.NewV4().String(),
		owner,
		uPkg,
		uNote,
		NewBuildStatus(),
		competition,
	}
	return bot, nil
}
