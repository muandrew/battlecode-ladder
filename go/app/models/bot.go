package models

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

const (
	//BotMaxName max size for bot name
	BotMaxName = 60
	//BotMaxDescription max size for bot description
	BotMaxDescription = 140
	//BotMaxNote max size for note on the bot
	BotMaxNote = 140
	//BotMaxPackage max size for bot's package name
	BotMaxPackage = 60
)

//Bot represents a particular build
type Bot struct {
	UUID            string
	Owner           *Competitor
	Package         UserString
	Note            UserString
	Status          *BuildStatus
	Competition     Competition
	CompetitionMeta string
}

//CreateBot creates a bot model
func CreateBot(
	owner *Competitor,
	pkg string, note string,
	competition Competition,
	competitionMeta string,
) (*Bot, error) {
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
		competitionMeta,
	}
	return bot, nil
}
