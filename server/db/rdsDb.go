package db

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

func NewRdsDb(addr string) *RdsDb {
	return &RdsDb{pool: &redis.Pool{
		MaxIdle: 3,
		IdleTimeout: 240 * time.Second,
		Dial: func () (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}}
}

func (db RdsDb) Ping(){
	c := db.pool.Get()
	defer c.Close()
	response, err := c.Do("PING")
	if err != nil {
		fmt.Print(err)
	} else {
		fmt.Print(response)
	}
}
func (db RdsDb) GetUserWithApp(app string, appUuid string, generateUser func() *models.User) *models.User {
	c := db.pool.Get()
	defer c.Close()
	appKey := "oauth"+":"+app+":"+appUuid
	userUuid, err := c.Do("GET", appKey)
	if err != nil {
		return nil
	}
	if userUuid != nil {
		userBin, err := c.Do("GET", userUuid)
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
	c := db.pool.Get()
	defer c.Close()
	key := "user:"+uuid
	bini, _ :=c.Do("GET",key)
	usr := &models.User{}
	bin := bini.([]byte)
	json.Unmarshal(bin, usr)
	return usr
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
