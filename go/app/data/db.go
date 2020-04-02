package data

import (
	"github.com/muandrew/battlecode-legacy-go/models"
)

type Db interface {
	GetUserWithApp(app string, appUuid string, generateUser func() *models.User) *models.User
	GetUser(uuid string) *models.User
	CreateBot(model *models.Bot) error
	UpdateBot(model *models.Bot) error
	GetBot(uuid string) *models.Bot
	GetBots(userUuid string, page int, pageSize int) ([]*models.Bot, int)
	GetPublicBots(page int, pageSize int) ([]*models.Bot, int)
	SetPublicBot(userUuid string, botUuid string) (*models.Bot, error)
	CreateMatch(model *models.Match) error
	UpdateMatch(model *models.Match) error
	GetMatch(matchUuid string) (*Match, error)
	GetDataMatches(userUuid string, page int, pageSize int) (*Page, error)
	GetMatches(userUuid string, page int, pageSize int) ([]*models.Match, int)
	CreateBcMap(model *models.BcMap) error
	UpdateBcMap(model *models.BcMap) error
	GetBcMap(uuid string) *models.BcMap
	GetBcMaps(userUuid string, page int, pageSize int) ([]*models.BcMap, int)
}
