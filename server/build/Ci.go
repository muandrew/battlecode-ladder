package build

import (
	"errors"
	"fmt"
	"github.com/jeffail/tunny"
	"github.com/muandrew/battlecode-ladder/data"
	"github.com/muandrew/battlecode-ladder/models"
	"github.com/muandrew/battlecode-ladder/utils"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"bufio"
)

const (
	folderPermission     = 0755
	forbiddenCharacters  = "~$"
	errorIllegalArgument = utils.Error("Illegal Argument(s)")
)

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
		return "", errors.New(fmt.Sprintf("Only relative and absolute pathing are allowed,"+
			" and if you are using (%s) in your directory structure, consider better names.",
			forbiddenCharacters))
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

func (c *Ci) UploadBotSource(file *multipart.FileHeader, bot *models.Bot) error {
	return c.upload(file, c.dirBot+"/"+bot.Uuid, "source.jar")
}

func (c *Ci) UploadMap(file *multipart.FileHeader, bcMap *models.BcMap) error {
	if file == nil || bcMap == nil {
		return errorIllegalArgument
	}
	err := c.upload(file, c.dirMap+"/"+bcMap.Uuid, file.Filename)
	if err != nil {
		return err
	}
	return c.db.CreateBcMap(bcMap)
}

func (c *Ci) upload(file *multipart.FileHeader, destDir string, destFile string) error {
	if file == nil || destDir == "" || destFile == "" {
		return errorIllegalArgument
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	os.MkdirAll(destDir, folderPermission)
	dst, err := os.Create(destDir + "/" + destFile)
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

func (c *Ci) SubmitJob(bot *models.Bot) {
	bot.Status.SetQueued()
	c.db.CreateBot(bot)
	c.pool.SendWorkAsync(func(workerId int) {
		bot.Status.SetStart()
		c.db.UpdateBot(bot)
		err := utils.RunShell("bash", []string{"scripts/build-bot.sh", c.dirBot, bot.Uuid})
		if err != nil {
			bot.Status.SetFailure()
		} else {
			bot.Status.SetSuccess()
		}
		c.db.UpdateBot(bot)
	}, nil)
}

func (c *Ci) RunMatch(bots []*models.Bot, bcMap *models.BcMap) error {
	if len(bots) != 2 {
		return errors.New("Currently only support 1v1.");
	}
	match, err := models.CreateMatch(bots, bcMap)
	if err != nil {
		return err
	}
	return c.RunMatchWithModel(match, bcMap)
}

func (c *Ci) RunMatchWithModel(match *models.Match, bcMap *models.BcMap) error {
	bot1 := match.Bots[0];
	bot2 := match.Bots[1];
	mapDir := ""
	mapName := ""
	if bcMap != nil {
		basename := bcMap.Name.GetRawString()
		mapDir = c.dirMap + "/" + bcMap.Uuid
		mapName = strings.TrimSuffix(basename, filepath.Ext(basename))
	}
	match.Status.SetQueued()
	c.db.CreateMatch(match)
	c.pool.SendWorkAsync(func(workerId int) {
		match.Status.SetStart()
		c.db.UpdateMatch(match)
		winner := models.WinnerNone
		err := utils.RunShellWithScan(
			"bash",
			[]string{
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
				mapDir,
				mapName,
			},
			func(scanner *bufio.Scanner) {
				for scanner.Scan() {
					line := scanner.Text()
					if strings.Contains(line, "wins (round") {
						index := strings.IndexRune(line, '(')
						if index != -1 {
							switch line[index+1] {
							case 'A':
								winner = 0
							case 'B':
								winner = 1
							default:
								winner = models.WinnerNone
							}
						}
					}
					fmt.Printf("%s\n", line)
				}
			},
			utils.BasicScanFunc,
		)
		if err != nil {
			match.Status.SetFailure()
		} else {
			match.Status.SetSuccess()
			match.Winner = winner
		}
		c.db.UpdateMatch(match)
	}, nil)
	return nil
}

func (c *Ci) RunGame(
	owner *models.Competitor,
	name string,
	description string,
	bots []*models.Bot,
	bcMap *models.BcMap) error {

	if bots == nil {
		return errors.New("Bots should not be empty.")
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
		err = c.RunMatchWithModel(match, nil)
		if err != nil {
			return err
		}
	}
	return nil;
}

func (c *Ci) Close() {
	c.pool.Close()
}

func (c *Ci) GetDirMatches() string {
	return c.dirMatch
}

func SetUpWorkspace(workerDir string, workerId int) {
	utils.FatalRunShell(
		"bash",
		[]string{
			"scripts/setup-worker-match-workspace.sh",
			workerDir,
			strconv.Itoa(workerId),
		},
	)
}
