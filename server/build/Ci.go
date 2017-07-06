package build

import (
	"github.com/jeffail/tunny"
	"github.com/muandrew/battlecode-ladder/models"
	"github.com/muandrew/battlecode-ladder/utils"
	"github.com/muandrew/battlecode-ladder/data"
	"strconv"
	"os"
	"path/filepath"
	"fmt"
	"io"
	"mime/multipart"
)

const (
	folderPermission = 0755
)

type Ci struct {
	db   data.Db
	pool *tunny.WorkPool

	dirBot    string
	dirData   string
	dirMatch  string
	dirUser   string
	dirWorker string
}

func getAndSetupDir(key string, fallback string) (string, error) {
	dir := utils.GetEnv(key)
	if dir == "" {
		dir = fallback
	}
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	os.MkdirAll(dir, folderPermission)
	return dir, nil
}

func NewCi(db data.Db) (*Ci, error) {
	dirData, err := getAndSetupDir("DIR_DATA", "../bcl-data")
	if err != nil {
		return nil, err
	}
	dirBot, err := getAndSetupDir("DIR_BOT", dirData+"/bot")
	if err != nil {
		return nil, err
	}
	dirMatch, err := getAndSetupDir("DIR_MATCH", dirData+"/match")
	if err != nil {
		return nil, err
	}
	dirUser, err := getAndSetupDir("DIR_USER", dirData+"/user")
	if err != nil {
		return nil, err
	}
	dirWorker, err := getAndSetupDir("DIR_WORKER", dirData+"/worker")
	if err != nil {
		return nil, err
	}
	pool, err := tunny.CreateCustomPool(CreateWorkers(dirWorker, 2)).Open()
	if err != nil {
		return nil, err
	}
	return &Ci{
		db,
		pool,
		dirBot,
		dirData,
		dirMatch,
		dirUser,
		dirWorker,
	}, nil
}

func (c Ci) UploadBotSource(file *multipart.FileHeader, bot *models.Bot) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	prefix := c.dirBot + "/" + bot.Uuid
	os.MkdirAll(prefix, folderPermission)
	dst, err := os.Create(prefix + "/source.jar")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}

func (c Ci) SubmitJob(bot *models.Bot) {
	bot.Status.SetQueued()
	c.db.CreateBot(bot)
	c.pool.SendWorkAsync(func(workerId int) {
		bot.Status.SetStart()
		c.db.UpdateBot(bot)
		err := utils.RunShell("sh", []string{"scripts/build-bot.sh", c.dirBot, bot.Uuid})
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
			c.dirBot,
			c.dirMatch,
			c.dirWorker,
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

func (c Ci) GetDirMatches() string {
	return c.dirMatch
}

func SetUpWorkspace(workerDir string, workerId int) {
	utils.FatalRunShell(
		"sh",
		[]string{
			"scripts/setup-worker-match-workspace.sh",
			workerDir,
			strconv.Itoa(workerId),
		},
	)
}
