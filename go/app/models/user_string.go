package models

import (
	"errors"
	"regexp"
	"strings"
)

//UserString represents user input, assume this to be dangerous
type UserString string

//UserStringFilter a function used to filter UserStrings
type UserStringFilter func(interface{}) (string, error)

const (
	//RegexFilterFilename a regex based filter, meant to be used for filenames
	RegexFilterFilename = "[^a-zA-Z0-9_.]"
	//RegexFilterPackage a regex based filter, meant to be used for java packages
	RegexFilterPackage = "[^a-zA-Z_.]"
	//RegexFilterText a regex based filter, meant to be used for text
	RegexFilterText = "[<>]"
)

//RegexBlacklist returns error if an illegal character is detected, nil otherwise
func RegexBlacklist(regexString string) func(s string) error {
	return func(s string) error {
		match, err := regexp.MatchString(regexString, s)
		if err != nil {
			return err
		}
		if match {
			return errors.New("Illegal character detected")
		}
		return nil
	}
}

//NewUserString creates a new instance from raw user input
func NewUserString(userInputString string, limit int, filters ...func(string) error) (UserString, error) {
	if len(userInputString) > limit {
		return "", errors.New("Above the limit")
	}

	for _, filter := range filters {
		err := filter(userInputString)
		if err != nil {
			return "", err
		}
	}
	userInputString = strings.Replace(userInputString, "<", "&lt", -1)
	userInputString = strings.Replace(userInputString, ">", "&gt", -1)
	return UserString(userInputString), nil
}

//GetRawString returns the string
//note: right now there is some preprocessing that api may need to be
//changed so this really returns raw string.
func (u UserString) GetRawString() string {
	return string(u)
}

//GetDisplayString returns the string as if it needs to be used for
//display
func (u UserString) GetDisplayString() string {
	// do any post processing
	return string(u)
}
