package data

import (
	"github.com/muandrew/battlecode-ladder/models"
	"github.com/garyburd/redigo/redis"
	"time"
	"fmt"
	"encoding/json"
)

const (
	addSet = "SET"
	addLpush = "LPUSH"
)

type RdsDb struct {
	pool *redis.Pool
}

func NewRdsDb(addr string) (*RdsDb, error) {
	rdb := &RdsDb{pool: &redis.Pool{
		MaxIdle: 3,
		IdleTimeout: 240 * time.Second,
		Dial: func () (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}}
	err := rdb.Ping()
	if err != nil {
		return nil, err
	} else {
		return rdb, nil
	}
}

func (db RdsDb) Ping() error {
	c := db.pool.Get()
	defer c.Close()
	response, err := c.Do("PING")
	if err != nil {
		return err
	} else {
		fmt.Printf("redis ping: %s\n",response)
		return nil
	}
}
func (db RdsDb) GetUserWithApp(app string, appUuid string, generateUser func() *models.User) *models.User {
	c := db.pool.Get()
	defer c.Close()
	appKey := "oauth"+":"+app+":"+appUuid
	userUuid, err := redis.String(c.Do("GET", appKey))
	if err != nil {
		return nil
	}
	if userUuid != "" {
		userBin, err := c.Do("GET", "user:"+userUuid)
		if err != nil {
			return nil
		}
		if userBin != nil {
			bin := userBin.([]byte)
			usr := &models.User{}
			json.Unmarshal([]byte(bin), usr)
			return usr
		}
	} else {
		user := generateUser()
		key := "user:"+user.Uuid
		userBin, _ := json.Marshal(user)
		c.Do("SET", key, userBin)
		c.Do("SET", appKey, user.Uuid)
		return user
	}
	return nil
}

func (db RdsDb) GetUser(uuid string) *models.User {
	model := &models.User{}
	db.getDeserializeModelWithKey(model, "user:"+uuid)
	return model
}

func (db RdsDb) EnqueueBot(model *models.Bot) error {
	c := db.pool.Get()
	defer c.Close()
	err := db.send(c, addSet, "user:"+model.UserUuid+":latest-bot", model)
	if err != nil {
		return err
	}
	err = db.send(c, addLpush, "user:"+model.UserUuid+":bot-list", model)
	if err != nil {
		return err
	}
	err = c.Flush()
	if err != nil {
		return err
	}
	_, err = c.Receive()
	if err != nil {
		return err
	}
	return nil
}

//func (db RdsDb) SetLatestBuild(model *models.Bot) error {
//	return db.setModelWithKey(model, "user:"+model.UserUuid+":latest-bot")
//}

func (db RdsDb) GetLatestBot(userUuid string) *models.Bot {
	model := &models.Bot{}
	db.getDeserializeModelWithKey(model, "user:"+userUuid+":latest-bot")
	return model
}

func (db RdsDb) SetLatestCompletedBot(model *models.Bot) error {
	return db.setModelWithKey(model, "user:"+model.UserUuid+":latest-complete-bot")
}

func (db RdsDb) GetLatestCompletedBot(userUuid string) *models.Bot {
	model := &models.Bot{}
	db.getDeserializeModelWithKey(model, "user:"+userUuid+":latest-complete-bot")
	return model
}

func (db RdsDb) AddCompletedMatch(model *models.Match) error {
	c := db.pool.Get()
	defer c.Close()

	err := c.Send("SET", "match:"+model.Uuid)
	if err != nil {
		return err
	}
	bin, err := json.Marshal(model)
	if err != nil {
		return err
	}
	for _,bot := range model.Bots {
		err := c.Send("LPUSH", "user:"+bot.UserUuid +":match-list", bin)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db RdsDb) GetAllUsers() []*models.User {
	//allUsers := make([]*models.User, len(db.users))
	//
	//i := 0
	//for _,v := range db.users {
	//	allUsers[i] = v
	//	i++
	//}
	//return allUsers
	return nil
}


func (db RdsDb) send(c redis.Conn, action string, key string, model interface{}) error {
	bin, err := json.Marshal(model)
	if err != nil {
		return err
	}
	return c.Send(action, key, bin)
}

func (db RdsDb) setModelWithKey(model interface{}, key string) error {
	c := db.pool.Get()
	defer c.Close()
	bin, err := json.Marshal(model)
	if err != nil {
		return err
	}
	_, err =c.Do("SET",key, bin)
	return err
}

func (db RdsDb) pushModelWithKey(model interface{}, key string) error {
	c := db.pool.Get()
	defer c.Close()
	bin, err := json.Marshal(model)
	if err != nil {
		return err
	}
	_, err = c.Do("LPUSH",key, bin)
	return err
}

func (db RdsDb) getDeserializeModelWithKey(model interface{}, key string) error {
	c := db.pool.Get()
	defer c.Close()
	bin, err :=c.Do("GET",key)
	if err != nil {
		return err
	}
	if bin != nil {
		json.Unmarshal(bin.([]byte), model)
	}
	return nil
}

func (db RdsDb) GetMatches(userUuid string, page int, pageSize int) ([]*models.Match, int) {
	c := db.pool.Get()
	defer c.Close()
	length, _ := redis.Int(c.Do("LLEN","user:"+userUuid +":match-list"))
	start := page * pageSize
	end := start + pageSize - 1
	rawBinArr, _ := c.Do("LRANGE","user:"+userUuid +":match-list", start, end)
	binArr := rawBinArr.([]interface{})
	matches := make([]*models.Match, len(binArr))
	for i, bin := range binArr {
		match := &models.Match{}
		json.Unmarshal(bin.([]byte), match)
		matches[i] = match
	}
	return matches, length
}
