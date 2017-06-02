package build

import (
	"github.com/jeffail/tunny"
	"github.com/muandrew/battlecode-ladder/models"
	"github.com/muandrew/battlecode-ladder/utils"
	"github.com/muandrew/battlecode-ladder/data"
	"strconv"
	"github.com/satori/go.uuid"
	"time"
)

type Ci struct {
	db    data.Db
	pool *tunny.WorkPool
}

func NewCi(db data.Db) *Ci {
	pool,_ := tunny.CreateCustomPool(CreateWorkers(2)).Open()
	return &Ci{
		db:db,
		pool: pool,
	}
}

func (c Ci) SubmitJob(bot *models.Bot) {
	c.db.EnqueueBot(bot)
	c.pool.SendWorkAsync(func (workerId int) {
		utils.RunShell("sh", []string{"scripts/build-bot.sh", bot.UserUuid, bot.Uuid})
		lb := c.db.GetLatestBot(bot.UserUuid)
		if lb.Uuid == bot.Uuid {
			c.db.SetLatestCompletedBot(bot)
		}
	},nil)
}

func (c Ci) RunMatch(bot1 *models.Bot, bot2 *models.Bot) {
	c.pool.SendWorkAsync(func (workerId int) {
		match := models.CreateMatch(uuid.NewV4().String())
		match.Bots = []*models.Bot{bot1, bot2}
		match.CompletedTimestamp = time.Now().Unix()
		utils.RunShell("sh", []string{
			"scripts/run-match.sh",
			strconv.Itoa(workerId),
			match.Uuid,
			bot1.UserUuid,
			bot1.Uuid,
			bot1.Package,
			bot2.UserUuid,
			bot2.Uuid,
			bot2.Package,
		})
		c.db.AddCompletedMatch(match)
	}, nil)
}

func (c Ci) Close() {
	c.pool.Close()
}
