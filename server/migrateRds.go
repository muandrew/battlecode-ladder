package main

import (
	"github.com/muandrew/battlecode-ladder/utils"
	"github.com/muandrew/battlecode-ladder/data"
	"github.com/garyburd/redigo/redis"
	"github.com/muandrew/battlecode-ladder/models"
	"github.com/labstack/gommon/log"
)

func main() {
	utils.InitMainEnv()
	db, _ := data.NewRdsDb(utils.GetRequiredEnvFatal("REDIS_ADDRESS"))
	err := db.Scan("match:*", func (c redis.Conn, key string) {
		match := &models.Match{}
		data.GetModel(c, key, match)
		if match.Competition == "" {
			match.Competition = models.BotCompetitionBC17
		}
		data.SendModel(c, data.AddSet, key, match)
		c.Flush()
		c.Receive()
	})
	logFatal(err)
	err = db.Scan("bot:*", func (c redis.Conn, key string) {
		bot := &models.Bot{}
		data.GetModel(c, key, bot)
		if bot.Competition == "" {
			bot.Competition = models.BotCompetitionBC17
		}
		data.SendModel(c, data.AddSet, key, bot)
		c.Flush()
		c.Receive()
	})
	logFatal(err)
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
