package build

import (
	"github.com/jeffail/tunny"
	"github.com/muandrew/battlecode-ladder/models"
	"github.com/muandrew/battlecode-ladder/utils"
	"github.com/muandrew/battlecode-ladder/data"
	"strconv"
)

type Ci struct {
	db   data.Db
	pool *tunny.WorkPool
}

func NewCi(db data.Db) *Ci {
	pool, _ := tunny.CreateCustomPool(CreateWorkers(2)).Open()
	return &Ci{
		db:   db,
		pool: pool,
	}
}

func (c Ci) SubmitJob(bot *models.Bot) {
	bot.Status.SetQueued()
	c.db.CreateBot(bot)
	c.pool.SendWorkAsync(func(workerId int) {
		bot.Status.SetStart()
		c.db.UpdateBot(bot)
		err := utils.RunShell("sh", []string{"scripts/build-bot.sh", bot.Uuid})
		if err != nil {
			bot.Status.SetFailure()
		} else {
			bot.Status.SetSuccess()
		}
		c.db.UpdateBot(bot)
	}, nil)
}

func (c Ci) RunMatch(bot1 *models.Bot, bot2 *models.Bot) {
	match := models.CreateMatch([]*models.Bot{bot1, bot2})
	match.Status.SetQueued()
	c.db.CreateMatch(match)
	c.pool.SendWorkAsync(func(workerId int) {
		match.Status.SetStart()
		c.db.UpdateMatch(match)
		err := utils.RunShell("sh", []string{
			"scripts/run-match.sh",
			strconv.Itoa(workerId),
			match.Uuid,
			bot1.Uuid,
			bot1.Package.GetPackageFormat(),
			bot2.Uuid,
			bot2.Package.GetPackageFormat(),
		})
		if err != nil {
			match.Status.SetFailure()
		} else {
			match.Status.SetSuccess()
		}
		match.Status.SetSuccess()
		c.db.UpdateMatch(match)
	}, nil)
}

func (c Ci) Close() {
	c.pool.Close()
}
