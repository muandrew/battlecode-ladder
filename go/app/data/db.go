package data

import (
	"github.com/muandrew/battlecode-legacy-go/models"
)

//Db represents an abstract contract for long term storage
type Db interface {
	GetUserWithApp(app string, appUUID string, generateUser func() *models.User) *models.User
	GetUser(uuid string) *models.User
	CreateBot(model *models.Bot) error
	UpdateBot(model *models.Bot) error
	GetBot(uuid string) *models.Bot
	GetBots(userUUID string, page int, pageSize int) ([]*models.Bot, int)
	GetPublicBots(page int, pageSize int) ([]*models.Bot, int)
	SetPublicBot(userUUID string, botUUID string) (*models.Bot, error)
	CreateMatch(model *models.Match) error
	UpdateMatch(model *models.Match) error
	GetMatch(matchUUID string) (*Match, error)
	GetDataMatches(userUUID string, page int, pageSize int) (*Page, error)
	GetMatches(userUUID string, page int, pageSize int) ([]*models.Match, int)
	CreateBcMap(model *models.BcMap) error
	UpdateBcMap(model *models.BcMap) error
	GetBcMap(uuid string) *models.BcMap
	GetBcMaps(userUUID string, page int, pageSize int) ([]*models.BcMap, int)
}
