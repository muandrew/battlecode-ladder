package data

import "github.com/muandrew/battlecode-ladder/models"

type Db interface {
	GetUserWithApp(app string, appUuid string, generateUser func() *models.User) *models.User
	GetUser(uuid string) *models.User
	GetLatestBot(userUuid string) *models.Bot
	GetLatestCompletedBot(userUuid string) *models.Bot
	SetLatestCompletedBot(model *models.Bot) error
	EnqueueBot(model *models.Bot) error
	AddCompletedMatch(model *models.Match) error
	GetMatches(userUuid string, page int, pageSize int) ([]*models.Match, int)
}
