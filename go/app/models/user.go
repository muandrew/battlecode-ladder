package models

import (
	"github.com/satori/go.uuid"
)

type SetupNewUser func() *User

const (
	UserMaxName = 140
)

type User struct {
	Uuid string
	Name UserString
}

func CreateUser(name string) (*User, error) {
	uName, err := NewUserString(name, UserMaxName)
	if err != nil {
		return nil, err
	}
	return &User{
		uuid.NewV4().String(),
		uName,
	}, nil
}
