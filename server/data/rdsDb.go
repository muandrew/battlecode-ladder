package data

import (
	"github.com/muandrew/battlecode-ladder/models"
	"github.com/garyburd/redigo/redis"
	"time"
	"fmt"
	"encoding/json"
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
	err := c.Send("SET", "user:"+model.UserUuid+":latest-bot", model)
	if err != nil {
		return err
	}
	err = c.Send("LPUSH", "user:"+model.UserUuid+":bot-list", model)
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

func (db RdsDb) setModelWithKey(model interface{}, key string) error {
	c := db.pool.Get()
	defer c.Close()
	_, err :=c.Do("SET",key, model)
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
