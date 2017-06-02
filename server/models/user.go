package models

import "github.com/satori/go.uuid"

type SetupNewUser func(*User) *User

type User struct{
	Uuid string
	Name string
}

func CreateUserWithNewUuid() *User {
	user := new(User)
	user.Uuid = uuid.NewV4().String()
	return user;
}
