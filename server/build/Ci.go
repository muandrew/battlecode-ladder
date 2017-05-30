package build

import (
	"github.com/jeffail/tunny"
	"github.com/muandrew/battlecode-ladder/models"
	"github.com/muandrew/battlecode-ladder/utils"
	"github.com/muandrew/battlecode-ladder/data"
)

type Ci struct {
	db    data.Db
	pool *tunny.WorkPool
}

func NewCi(db data.Db) *Ci {
	pool,_ := tunny.CreateCustomPool(CreatePool(2)).Open()
	return &Ci{
		db:db,
		pool: pool,
	}
}

func (c Ci) SubmitJob(bot *models.Bot) {
	c.pool.SendWork(func (workerId int) {
		utils.RunShell("sh", []string{"scripts/build-bot.sh", bot.UserUuid, bot.Uuid})
		lb := c.db.GetLatestBot(bot.UserUuid)
		if lb.Uuid == bot.Uuid {
			c.db.SetLatestCompletedBot(bot)
		}
	})
}

func (c Ci) RunMatch(bot1 *models.Bot, bot2 *models.Bot) {
	c.pool.SendWork(func (workerId int) {

	})
}

func (c Ci) Close() {
	c.pool.Close()
}
