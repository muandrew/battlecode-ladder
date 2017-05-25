package build

import (
	"github.com/muandrew/battlecode-ladder/db"
	"github.com/jeffail/tunny"
	"github.com/muandrew/battlecode-ladder/models"
	"github.com/muandrew/battlecode-ladder/utils"
)

type Ci struct {
	d db.Db
	pool *tunny.WorkPool
}

func NewCi(d db.Db) *Ci {
	pool,_ := tunny.CreatePoolGeneric(1).Open()
	return &Ci{
		d:d,
		pool: pool,
	}
}

func (c Ci) SubmitJob(bot *models.Bot) {
	c.pool.SendWork(func () {
		utils.RunShell("sh", []string{"./scripts/build-bot.sh", bot.UserUuid, bot.Uuid})
		lb := c.d.GetLatestBot(bot.UserUuid)
		if lb.Uuid == bot.Uuid {
			c.d.SetLatestCompletedBot(bot)
		}
	})
}

func (c Ci) Close() {
	c.pool.Close()
}
