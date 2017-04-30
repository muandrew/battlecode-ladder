package models

import "github.com/satori/go.uuid"

type SetupNewUser func(*User) *User

type User struct{
	Uuid string
	Name string
}

var UserDummy = &User{
	Uuid:"0",
	Name:"Dummy User",
}

func CreateUserWithNewUuid() *User {
	user := new(User)
	user.Uuid = uuid.NewV4().String()
	return user;
}
