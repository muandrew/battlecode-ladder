package models

import (
	"regexp"
	"errors"
	"strings"
)

type UserString string
type UserStringFilter func(interface{}) (string, error)

const(
	RegexFilterPackage="[^a-zA-Z_.]"
	RegexFilterText="[<>]"
)

func RegexBlacklist(regexString string)  func (s string) error {
	return func(s string) error {
		match, err := regexp.MatchString(regexString, s)
		if err != nil { return err}
		if match {
			return errors.New("Illegal character detected.")
		} else {
			return nil
		}
	}
}

func NewUserString(userInputString string, limit int, filters ...func(string) error) (UserString, error) {
	if len(userInputString) > limit {
		return "", errors.New("Above the limit")
	} else {
		userInputString = strings.Replace(userInputString, "<", "&lt", -1)
		userInputString = strings.Replace(userInputString, ">", "&gt", -1)
		for _, filter := range filters {
			err := filter(userInputString)
			if err != nil {
				return "", err
			}
		}
		return UserString(userInputString), nil
	}
}

func (u UserString) GetRawString() string {
	return string(u)
}

func (u UserString) GetPackageFormat() string {
	return string(u)
}

func (u UserString) GetDisplayString() string {
	// do any post processing
	return string(u)
}
