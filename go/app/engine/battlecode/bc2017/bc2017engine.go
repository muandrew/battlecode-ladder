package bc2017

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/markbates/pkger"
	"github.com/muandrew/battlecode-legacy-go/models"
	"github.com/muandrew/battlecode-legacy-go/utils"
)

//Engine runs battlecode 2017
type Engine struct {
}

//Competition see parent.
func (eng *Engine) Competition() models.Competition {
	return models.CompetitionBC17
}

//ActivateAssets see parent.
func (eng *Engine) ActivateAssets() {
	pkger.Include("/engine/battlecode/bc2017/assets")
}

//BattleBotSetup see parent
func (eng *Engine) BattleBotSetup(
	workerID int,
	workspaceDir string,
	match *models.Match,
) error {
	err := utils.CopyFromPkgr(
		"/engine/battlecode/bc2017/assets/bot-builder",
		filepath.Join(workspaceDir, "workspace"),
	)
	if err != nil {
		return err
	}
	fileToSource, err := os.Create(filepath.Join(workspaceDir, "source.sh"))
	if err != nil {
		return err
	}
	defer fileToSource.Close()
	for idx, bot := range match.Bots {
		fileToSource.WriteString(fmt.Sprintf(
			"export BOT_%d_NAME=%s\n",
			idx,
			bot.Package.GetRawString(),
		))
	}
	fileToSource.WriteString(fmt.Sprintf(
		"export WORKER_ID=%d\n",
		workerID,
	))

	err = utils.CopyFromPkgr(
		"/engine/battlecode/bc2017/assets/runmatch/run.sh",
		filepath.Join(workspaceDir, "run.sh"),
	)
	return err
}

//BattleBotPostProcessing see parent
func (eng *Engine) BattleBotPostProcessing(
	matchPath string,
	match *models.Match,
) error {
	file, err := os.Open(filepath.Join(matchPath, "result", "log.txt"))
	if err != nil {
		return nil
	}
	winner := models.WinnerNone
	defer file.Close()
	scanner := bufio.NewScanner(file)
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
	}
	match.Winner = winner
	return nil
}

//BuildBotSetup see parent
func (eng *Engine) BuildBotSetup(
	workerID int,
	workspaceDir string,
	botUUID string,
) error {
	err := utils.CopyFromPkgr(
		"/engine/battlecode/bc2017/assets/bot-builder",
		filepath.Join(workspaceDir, "workspace"),
	)
	if err != nil {
		return err
	}
	err = utils.CopyFromPkgr(
		"/engine/battlecode/bc2017/assets/buildbot/run.sh",
		filepath.Join(workspaceDir, "run.sh"),
	)
	return err
}
