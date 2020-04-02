package models

import (
	"github.com/satori/go.uuid"
	"path/filepath"
	"errors"
	"fmt"
)

type BcMap struct {
	Uuid        string
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
