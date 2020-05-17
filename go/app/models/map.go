package models

import (
	"errors"
	"fmt"
	"path/filepath"

	uuid "github.com/satori/go.uuid"
)

type BcMap struct {
	UUID        string
	Owner       *Competitor
	Competition Competition
	Name        UserString
	Description UserString
}

func CreateBcMap(owner *Competitor, filename string, description string) (*BcMap, error) {
	uFileName, err := NewUserString(filename, BotMaxName, RegexBlacklist(RegexFilterFilename))
	if err != nil {
		return nil, err
	}
	uDesc, err := NewUserString(description, BotMaxDescription, RegexBlacklist(RegexFilterText))
	if err != nil {
		return nil, err
	}
	competition, err := filenameToCompetition(uFileName.GetRawString())
	if err != nil {
		return nil, err
	}
	return &BcMap{
		uuid.NewV4().String(),
		owner,
		competition,
		uFileName,
		uDesc,
	}, nil
}

func filenameToCompetition(filename string) (Competition, error) {
	ext := filepath.Ext(filename)
	fmt.Println(ext)
	switch {
	case ext == ".map17":
		return CompetitionBC17, nil
	default:
		return "", errors.New(fmt.Sprintf("Unknown Extension type: %q", ext))
	}
}
