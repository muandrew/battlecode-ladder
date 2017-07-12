package data

import "github.com/muandrew/battlecode-ladder/models"

type Db interface {
	GetUserWithApp(app string, appUuid string, generateUser func() *models.User) *models.User
	GetUser(uuid string) *models.User
	GetBot(uuid string) *models.Bot
	CreateBot(model *models.Bot) error
	UpdateBot(model *models.Bot) error
	GetBots(userUuid string, page int, pageSize int) ([]*models.Bot, int)
	CreateMatch(model *models.Match) error
	UpdateMatch(model *models.Match) error
	GetMatches(userUuid string, page int, pageSize int) ([]*models.Match, int)
	CreateBcMap(model *models.BcMap) error
	UpdateBcMap(model *models.BcMap) error
	GetBcMaps(userUuid string, page int, pageSize int) ([]*models.BcMap, int)
}
