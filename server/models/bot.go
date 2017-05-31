package models

import "github.com/satori/go.uuid"

type Bot struct{
	UserUuid string
	Uuid string
	Name string
	Package string
}

func CreateBotWithNewUuidAndUserUuid(userUuid string) *Bot {
	build := new(Bot)
	build.Uuid = uuid.NewV4().String()
	build.UserUuid = userUuid
	return build;
}
