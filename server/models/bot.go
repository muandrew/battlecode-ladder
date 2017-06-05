package models

import "github.com/satori/go.uuid"

const (
	BotMaxName = 60
	BotMaxDescription = 140
	BotMaxPackage = 60
)

type Bot struct {
	Uuid string
	Package UserString
	Name UserString
	Description UserString
	//todo deprected
	UserUuid string
}

func CreateBot(pkg string, name string, description string) (*Bot, error) {
	uPkg, err := NewUserString(pkg, BotMaxPackage, RegexBlacklist(RegexFilterPackage))
	if err != nil {return nil, err}
	uName, err := NewUserString(name, BotMaxName, RegexBlacklist(RegexFilterText))
	if err != nil {return nil, err}
	uDesc, err := NewUserString(description, BotMaxDescription, RegexBlacklist(RegexFilterText))
	if err != nil {return nil, err}
	bot := &Bot{
		"",
		uPkg,
		uName,
		uDesc,
		"",
	}
	if bot.Uuid == "" {
		uuid.NewV4().String()
	}
	return bot, nil
}
