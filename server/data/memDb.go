package data

import (
	"github.com/muandrew/battlecode-ladder/models"
)

type MemDb struct {
	apps  map[string]string
	users map[string]*models.User
}

func NewMemDb() *MemDb {
	return &MemDb{make(map[string]string), make(map[string]*models.User)}
}

func (db MemDb) GetUserWithApp(app string, appUuid string, generateUser func() *models.User) *models.User {
	appKey := app + ":" + appUuid
	uuid := db.apps[appKey]
	if uuid == "" {
		user := generateUser()
		db.apps[appKey] = user.Uuid
		db.users[user.Uuid] = user
		return user
	} else {
		user := db.GetUser(uuid)
		return user
	}
}

func (db MemDb) GetUser(uuid string) *models.User {
	return db.users[uuid]
}

func (db MemDb) GetAllUsers() []*models.User {
	allUsers := make([]*models.User, len(db.users))

	i := 0
	for _, v := range db.users {
		allUsers[i] = v
		i++
	}
	return allUsers
}
