package build

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jeffail/tunny"
	"github.com/labstack/gommon/log"
	"github.com/muandrew/battlecode-legacy-go/data"
	"github.com/muandrew/battlecode-legacy-go/engine"
	"github.com/muandrew/battlecode-legacy-go/models"
	"github.com/muandrew/battlecode-legacy-go/utils"
)

const (
	forbiddenCharacters  = "~$"
	errorIllegalArgument = utils.Error("Illegal Argument(s)")
)

//Ci represents the build system
type Ci struct {
	db   data.Db
	pool *tunny.WorkPool

	dirBot    string
	dirData   string
	dirMap    string
	dirMatch  string
	dirUser   string
	dirWorker string
}

func getAndSetupDir(key string, fallback string) (string, error) {
	dir := utils.GetEnv(key)
	if dir == "" {
		dir = fallback
	} else if strings.ContainsAny(dir, forbiddenCharacters) {
		return "", fmt.Errorf("Only relative and absolute pathing are allowed,"+
			" and if you are using (%s) in your directory structure, consider better names.",
			forbiddenCharacters)
	}
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	os.MkdirAll(dir, utils.FileModeStandardFolder)
	return dir, nil
}

//NewCi creates a new instance of Ci
func NewCi(db data.Db) (*Ci, error) {
	dirData, err := getAndSetupDir("DIR_DATA", "../bcl-data")
	if err != nil {
		return nil, err
	}
	dirBot, err := getAndSetupDir("DIR_BOT", dirData+"/bot")
	if err != nil {
		return nil, err
	}
	dirMap, err := getAndSetupDir("DIR_MAP", dirData+"/map")
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
		dirMap,
		dirMatch,
		dirUser,
		dirWorker,
	}, nil
}

//UploadBotSource uploads a bots source.
func (c *Ci) UploadBotSource(file *multipart.FileHeader, bot *models.Bot) error {
	return c.upload(file, c.botSourcePath(bot.Uuid))
}

//UploadMap uploads a map
func (c *Ci) UploadMap(file *multipart.FileHeader, bcMap *models.BcMap) error {
	if file == nil || bcMap == nil {
		return errorIllegalArgument
	}
	err := c.upload(file, filepath.Join(c.dirMap, bcMap.Uuid, file.Filename))
	if err != nil {
		return err
	}
	return c.db.CreateBcMap(bcMap)
}

func (c *Ci) upload(file *multipart.FileHeader, dest string) error {
	if file == nil || dest == "" {
		return errorIllegalArgument
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	os.MkdirAll(filepath.Dir(dest), utils.FileModeStandardFolder)
	dst, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}

//BuildBot builds a bot
func (c *Ci) BuildBot(eng engine.Engine, bot *models.Bot) {
	bot.Status.SetQueued()
	c.db.CreateBot(bot)
	c.pool.SendWorkAsync(func(workerId int) {
		bot.Status.SetStart()
		c.db.UpdateBot(bot)
		var err error

		// prep the workspace
		workspaceDir := c.workspaceDir(workerId)
		dirReset(workspaceDir)

		// copy the soruces over
		if err == nil {
			err = utils.CopyPlain(
				c.botSourcePath(bot.Uuid),
				filepath.Join(workspaceDir, "source.zip"),
			)
		}
		if err == nil {
			// let the engine do prep work
			err = eng.BuildBotSetup(
				workerId,
				workspaceDir,
				bot.Uuid,
			)
		}
		if err == nil {
			err = c.runRun(workspaceDir)
		}
		if err == nil {
			err = utils.CopyPlain(
				filepath.Join(workspaceDir, "result.zip"),
				c.botResultPath(bot.Uuid),
			)
		}
		// updating model
		if err != nil {
			log.Errorf("ERR: %s", err.Error())
			bot.Status.SetFailure()
		} else {
			bot.Status.SetSuccess()
		}
		c.db.UpdateBot(bot)
	}, nil)
}

//RunMatch runs a single match
func (c *Ci) RunMatch(e engine.Engine, bots []*models.Bot, bcMap *models.BcMap) error {
	if len(bots) != 2 {
		return errors.New("Currently only support 1v1")
	}
	match, err := models.CreateMatch(bots, bcMap)
	if err != nil {
		return err
	}
	return c.RunMatchWithModel(e, match, bcMap)
}

//RunMatchWithModel runs a single match
func (c *Ci) RunMatchWithModel(
	e engine.Engine,
	match *models.Match,
	bcMap *models.BcMap,
) error {
	match.Status.SetQueued()
	c.db.CreateMatch(match)
	c.pool.SendWorkAsync(func(workerId int) {
		match.Status.SetStart()
		c.db.UpdateMatch(match)
		var err error

		// prep the workspace
		workspaceDir := c.workspaceDir(workerId)
		dirReset(workspaceDir)

		//there prob needs to be more specialization with map copy
		if err == nil && bcMap != nil {
			mapFileName := bcMap.Name.GetRawString()
			mapWorkspaceDir := filepath.Join(workspaceDir, "map")
			err = os.MkdirAll(mapWorkspaceDir, utils.FileModeStandardFolder)
			err = utils.CopyPlain(
				filepath.Join(c.mapPath(bcMap.Uuid), mapFileName),
				filepath.Join(mapWorkspaceDir, mapFileName),
			)
		}
		if err == nil {
			// copy over the results
			for idx, bot := range match.Bots {
				err = utils.CopyPlain(
					c.botResultPath(bot.Uuid),
					filepath.Join(workspaceDir, fmt.Sprintf("bot%d.zip", idx)),
				)
				if err != nil {
					break
				}
			}
		}
		if err == nil {
			// allow each engine to run its own setup.
			err = e.BattleBotSetup(
				workerId,
				workspaceDir,
				match,
			)
		}
		if err == nil {
			err = c.runRun(workspaceDir)
		}
		matchPath := c.matchPath(match.Uuid)
		if err == nil {
			err = os.MkdirAll(matchPath, utils.FileModeStandardFolder)
		}
		if err == nil {
			err = utils.CopyPlain(
				filepath.Join(workspaceDir, "result.zip"),
				filepath.Join(matchPath, "result.zip"),
			)
		}
		if err == nil {
			err = utils.Unzip(matchPath, "result.zip", "result")
		}
		if err == nil {
			err = e.BattleBotPostProcessing(
				matchPath,
				match,
			)
		}
		// updating model
		if err != nil {
			log.Errorf("ERR: %s", err.Error())
			match.Status.SetFailure()
		} else {
			match.Status.SetSuccess()
		}
		c.db.UpdateMatch(match)
	}, nil)
	return nil
}

func (c *Ci) runRun(workspaceDir string) error {
	cmd := exec.Command("bash", "run.sh")
	cmd.Dir = workspaceDir
	// work on removing as many variables as possible
	//cmd.Env = []string{}
	return cmd.Run()
}

//RunGame execute a series of matches
func (c *Ci) RunGame(
	eng engine.Engine,
	owner *models.Competitor,
	name string,
	description string,
	bots []*models.Bot,
	bcMap *models.BcMap) error {

	if bots == nil {
		return errors.New("Bots should not be empty")
	}
	game, err := models.CreateGameRoundRobin(
		owner,
		models.CompetitionBC17,
		name,
		description,
		bots,
		bcMap,
	)
	if err != nil {
		return err
	}
	for _, match := range game.Matches {
		err = c.RunMatchWithModel(eng, match, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

//Close call to cleanup all resources
func (c *Ci) Close() {
	c.pool.Close()
}

//GetDirMatches returns the directory where match results are
func (c *Ci) GetDirMatches() string {
	return c.dirMatch
}

func dirReset(dir string) {
	os.RemoveAll(dir)
	os.MkdirAll(dir, utils.FileModeStandardFolder)
}

func (c *Ci) botSourcePath(botUUID string) string {
	return filepath.Join(c.dirBot, botUUID, "source.zip")
}

func (c *Ci) botResultPath(botUUID string) string {
	return filepath.Join(c.dirBot, botUUID, "result.zip")
}

func (c *Ci) mapPath(mapUUID string) string {
	return filepath.Join(c.dirMap, mapUUID)
}

func (c *Ci) matchPath(matchUUID string) string {
	return filepath.Join(c.dirMatch, matchUUID)
}

func (c *Ci) workspaceDir(workerID int) string {
	return filepath.Join(c.dirWorker, strconv.Itoa(workerID))
}
