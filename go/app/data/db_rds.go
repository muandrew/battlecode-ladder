package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/muandrew/battlecode-legacy-go/models"
)

const (
	//AddSet redis command to set.
	AddSet   = "SET"
	addLpush = "LPUSH"
)

//RdsDb and implementation of Db with Redis
//In most cases using Redis is a bad idea as your main
//datastore. Probably also in this case.
type RdsDb struct {
	pool *redis.Pool
}

//NewRdsDb sets up a new redis database
func NewRdsDb(addr string) (*RdsDb, error) {
	rdb := &RdsDb{pool: &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}}
	err := rdb.Ping()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}

//Ping check to see if the connection is working
func (db *RdsDb) Ping() error {
	c := db.pool.Get()
	defer c.Close()
	response, err := c.Do("PING")
	if err != nil {
		return err
	}
	fmt.Printf("redis ping: %s\n", response)
	return nil
}

//GetUserWithApp get a user from the specified app.
func (db *RdsDb) GetUserWithApp(app string, appUUID string, generateUser func() *models.User) *models.User {
	c := db.pool.Get()
	defer c.Close()
	appKey := "oauth" + ":" + app + ":" + appUUID
	userUUID, _ := redis.String(c.Do("GET", appKey))
	if userUUID != "" {
		userBin, err := c.Do("GET", "user:"+userUUID)
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
		key := "user:" + user.UUID
		userBin, _ := json.Marshal(user)
		c.Do("SET", key, userBin)
		c.Do("SET", appKey, user.UUID)
		return user
	}
	return nil
}

//GetUser gets the user model
func (db *RdsDb) GetUser(uuid string) *models.User {
	model := &models.User{}
	db.getModelForKey(model, "user:"+uuid)
	return model
}

//GetBot gets teh bot model
func (db *RdsDb) GetBot(uuid string) *models.Bot {
	model := &models.Bot{}
	err := db.getModelForKey(model, getBotKeyWithUUID(uuid))
	if err != nil {
		return nil
	}
	return model
}

//CreateBot creates a bot entry
func (db *RdsDb) CreateBot(model *models.Bot) error {
	c := db.pool.Get()
	defer c.Close()

	err := SendModel(c, AddSet, getBotKey(model), model)
	if err != nil {
		return err
	}
	err = c.Send(addLpush, getPrefix(model.Owner)+":bot-list", model.UUID)
	if err != nil {
		return err
	}
	_, err = flushAndReceive(c)
	return err
}

//UpdateBot updaates a bot entry
func (db *RdsDb) UpdateBot(model *models.Bot) error {
	return db.setModelForKey(model, getBotKey(model))
}

//GetBots gets a list of bots
func (db *RdsDb) GetBots(userUUID string, page int, pageSize int) ([]*models.Bot, int) {
	c := db.pool.Get()
	defer c.Close()
	length, _ := redis.Int(c.Do("LLEN", "user:"+userUUID+":bot-list"))
	start := page * pageSize
	end := start + pageSize - 1
	botUUIDs, err := redis.Strings(c.Do("LRANGE", "user:"+userUUID+":bot-list", start, end))
	if err != nil {
		return nil, 0
	}
	bots := make([]*models.Bot, len(botUUIDs))

	for i, botUUID := range botUUIDs {
		bot := &models.Bot{}
		err = GetModel(c, getBotKeyWithUUID(botUUID), bot)
		if err != nil {
			return nil, 0
		}
		bots[i] = bot
	}
	return bots, length
}

//GetPublicBots gets a list of public bots
func (db *RdsDb) GetPublicBots(page int, pageSize int) ([]*models.Bot, int) {
	c := db.pool.Get()
	defer c.Close()

	botUUIDs, err := redis.Strings(c.Do("ZREVRANGE", "public:bot-list", 0, -1))
	if err != nil {
		return nil, 0
	}

	length := len(botUUIDs)
	bots := make([]*models.Bot, length)
	for i, botUUID := range botUUIDs {
		bot := &models.Bot{}
		err = GetModel(c, getBotKeyWithUUID(botUUID), bot)
		if err != nil {
			return nil, 0
		}
		bots[i] = bot
	}
	return bots, length
}

//SetPublicBot set a bot as public
func (db *RdsDb) SetPublicBot(userUUID string, botUUID string) (*models.Bot, error) {
	c := db.pool.Get()
	defer c.Close()

	currentBotUUID, _ := redis.String(c.Do("GET", "user:"+userUUID+":public-bot"))

	bot := &models.Bot{}
	err := GetModel(c, getBotKeyWithUUID(botUUID), bot)
	if err != nil {
		return nil, err
	}

	if bot.Owner.UUID != userUUID {
		return nil, errors.New("you can only set your own bot")
	}
	if bot.Status.Status != models.BuildStatusSuccess {
		return nil, errors.New("you should only set successful bots")
	}

	// user removing public bot
	if bot == nil {
		if currentBotUUID != "" {
			err := c.Send("DEL", "user:"+userUUID+":public-bot")
			if err != nil {
				return nil, err
			}
			err = c.Send("ZREM", "public:bot-list", currentBotUUID)
			if err != nil {
				return nil, err
			}
		}
	} else {
		err = c.Send("SET", "user:"+userUUID+":public-bot", bot.UUID)
		if err != nil {
			return nil, err
		}
		if currentBotUUID != "" {
			err = c.Send("ZREM", "public:bot-list", currentBotUUID)
			if err != nil {
				return nil, err
			}
		}
		err = c.Send("ZADD", "public:bot-list", time.Now().Unix(), bot.UUID)
		if err != nil {
			return nil, err
		}
	}
	err = c.Flush()
	if err != nil {
		return nil, err
	}
	return bot, nil
}

//CreateMatch creates a match entry
func (db *RdsDb) CreateMatch(model *models.Match) error {
	c := db.pool.Get()
	defer c.Close()

	err := SendModel(c, AddSet, getMatchKey(model), CreateMatch(model))
	if err != nil {
		return err
	}
	done := make(map[string]bool)
	for _, bot := range model.Bots {
		ownerPrefix := getPrefix(bot.Owner)
		if !done[ownerPrefix] {
			err := c.Send(addLpush, ownerPrefix+":match-list", model.UUID)
			if err != nil {
				return err
			}
			done[ownerPrefix] = true
		}
	}
	_, err = flushAndReceive(c)
	return err
}

//UpdateMatch updates a match entry
func (db *RdsDb) UpdateMatch(model *models.Match) error {
	return db.setModelForKey(CreateMatch(model), getMatchKey(model))
}

//GetMatch gets a match model
func (db *RdsDb) GetMatch(matchUUID string) (*Match, error) {
	model := &Match{}
	err := db.getModelForKey(model, getMatchKeyWithUUID(matchUUID))
	if err != nil {
		return nil, err
	}
	return model, nil
}

//GetDataMatches gets a page of data Match models, they are an intermediate format.
func (db *RdsDb) GetDataMatches(userUUID string, page int, pageSize int) (*Page, error) {
	c := db.pool.Get()
	defer c.Close()
	length, _ := redis.Int(c.Do("LLEN", "user:"+userUUID+":match-list"))
	start := page * pageSize
	end := start + pageSize - 1
	matchUUIDs, err := redis.Strings(c.Do("LRANGE", "user:"+userUUID+":match-list", start, end))
	if err != nil {
		return nil, err
	}
	matches := make([]interface{}, len(matchUUIDs))

	for i, matchUUID := range matchUUIDs {
		rdsMatch := &Match{}
		err = GetModel(c, getMatchKeyWithUUID(matchUUID), rdsMatch)
		if err != nil {
			return nil, err
		}
		matches[i] = rdsMatch
	}
	return &Page{
		matches,
		length,
	}, err
}

//GetMatches gets a page of matches
func (db *RdsDb) GetMatches(userUUID string, page int, pageSize int) ([]*models.Match, int) {
	c := db.pool.Get()
	defer c.Close()
	length, _ := redis.Int(c.Do("LLEN", "user:"+userUUID+":match-list"))
	start := page * pageSize
	end := start + pageSize - 1
	matchUUIDs, err := redis.Strings(c.Do("LRANGE", "user:"+userUUID+":match-list", start, end))
	if err != nil {
		return nil, 0
	}
	matches := make([]*models.Match, len(matchUUIDs))

	for i, matchUUID := range matchUUIDs {
		rdsMatch := &Match{}
		err = GetModel(c, getMatchKeyWithUUID(matchUUID), rdsMatch)
		if err != nil {
			return nil, 0
		}

		bots := make([]*models.Bot, len(rdsMatch.BotUUIDs))
		for j, botUUID := range rdsMatch.BotUUIDs {
			bot := &models.Bot{}
			err = GetModel(c, getBotKeyWithUUID(botUUID), bot)
			if err != nil {
				return nil, 0
			}
			bots[j] = bot
		}
		//noinspection GoStructInitializationWithoutFieldNames
		match := &models.Match{
			rdsMatch.UUID,
			bots,
			rdsMatch.MapUUID,
			rdsMatch.Winner,
			rdsMatch.Status,
			rdsMatch.Competition,
		}
		matches[i] = match
	}
	return matches, length
}

//CreateBcMap creates a new entry
func (db *RdsDb) CreateBcMap(model *models.BcMap) error {
	c := db.pool.Get()
	defer c.Close()

	err := SendModel(c, AddSet, getBcMapKey(model), model)
	if err != nil {
		return err
	}
	err = c.Send(addLpush, getPrefix(model.Owner)+":map-list", model.UUID)
	if err != nil {
		return err
	}
	_, err = flushAndReceive(c)
	return err
}

//UpdateBcMap updates an entry of BcMap
func (db *RdsDb) UpdateBcMap(model *models.BcMap) error {
	return db.setModelForKey(model, getBcMapKey(model))
}

//GetBcMap retrieves an entry of BcMap
func (db *RdsDb) GetBcMap(uuid string) *models.BcMap {
	model := &models.BcMap{}
	err := db.getModelForKey(model, getBcMapWithUUID(uuid))
	if err != nil {
		return nil
	}
	return model
}

//GetBcMaps retrieves a page of BcMap
func (db *RdsDb) GetBcMaps(userUUID string, page int, pageSize int) ([]*models.BcMap, int) {
	c := db.pool.Get()
	defer c.Close()
	length, _ := redis.Int(c.Do("LLEN", "user:"+userUUID+":map-list"))
	start := page * pageSize
	end := start + pageSize - 1
	bcMapUUIDs, err := redis.Strings(c.Do("LRANGE", "user:"+userUUID+":map-list", start, end))
	if err != nil {
		return nil, 0
	}
	bcMaps := make([]*models.BcMap, len(bcMapUUIDs))

	for i, bcMapUUID := range bcMapUUIDs {
		bcMap := &models.BcMap{}
		err = GetModel(c, getBcMapWithUUID(bcMapUUID), bcMap)
		if err != nil {
			return nil, 0
		}
		bcMaps[i] = bcMap
	}
	return bcMaps, length
}

/*
 utility
*/

//GetModel gets a model from Redis
func GetModel(c redis.Conn, key string, model interface{}) error {
	bin, err := c.Do("GET", key)
	if err != nil {
		return err
	}
	if bin != nil {
		return json.Unmarshal(bin.([]byte), model)
	}
	return fmt.Errorf("Couldn't find model for key: %q", key)
}

//SendModel sends a model to redis
func SendModel(c redis.Conn, action string, key string, model interface{}) error {
	bin, err := json.Marshal(model)
	if err != nil {
		return err
	}
	return c.Send(action, key, bin)
}

func flushAndReceive(c redis.Conn) (interface{}, error) {
	err := c.Flush()
	if err != nil {
		return nil, err
	}
	return c.Receive()
}

func (db *RdsDb) getModelForKey(model interface{}, key string) error {
	c := db.pool.Get()
	defer c.Close()
	return GetModel(c, key, model)
}

func (db *RdsDb) setModelForKey(model interface{}, key string) error {
	c := db.pool.Get()
	defer c.Close()
	err := SendModel(c, AddSet, key, model)
	if err != nil {
		return err
	}
	_, err = flushAndReceive(c)
	return err
}

func getPrefix(c *models.Competitor) string {
	return c.Type.String() + ":" + c.UUID
}

func getMatchKey(m *models.Match) string {
	return getMatchKeyWithUUID(m.UUID)
}

func getMatchKeyWithUUID(key string) string {
	return "match:" + key
}

func getBotKey(b *models.Bot) string {
	return getBotKeyWithUUID(b.UUID)
}

func getBotKeyWithUUID(uuid string) string {
	return "bot:" + uuid
}

func getBcMapKey(m *models.BcMap) string {
	return getBcMapWithUUID(m.UUID)
}

func getBcMapWithUUID(uuid string) string {
	return "map:" + uuid
}

//Scan scan for a pattern in Redis
func (db *RdsDb) Scan(pattern string, run func(redis.Conn, string)) error {
	c := db.pool.Get()
	defer c.Close()
	fmt.Printf("Scanning for %q\n", pattern)
	index := 0
	for true {
		reply, err := redis.Values(c.Do("scan", index, "match", pattern))
		if err != nil {
			return err
		}
		index, err = redis.Int(reply[0], nil)
		if err != nil {
			return err
		}
		keys, err := redis.Strings(reply[1], nil)
		if err != nil {
			return err
		}
		for _, key := range keys {
			run(c, key)
		}
		if index == 0 {
			fmt.Printf("done\n")
			break
		} else {
			fmt.Printf("idx %d complete.\n", index)
		}
	}
	return nil
}

//end utility

//deprecate
func (db *RdsDb) pushModelForKey(model interface{}, key string) error {
	c := db.pool.Get()
	defer c.Close()
	err := SendModel(c, addLpush, key, model)
	if err != nil {
		return err
	}
	_, err = flushAndReceive(c)
	return err
}

//end deprecate
