package models

import uuid "github.com/satori/go.uuid"

//SetupNewUser implement this if you know how to setup a users
type SetupNewUser func() *User

const (
	//UserMaxName The max character limit.
	UserMaxName = 140
)

//User the heart and soul of an app.
type User struct {
	UUID string
	Name UserString
}

//CreateUser creates a new user
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
