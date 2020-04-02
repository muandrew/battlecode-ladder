package lazy

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/muandrew/battlecode-legacy-go/auth"
	"github.com/muandrew/battlecode-legacy-go/build"
	"github.com/muandrew/battlecode-legacy-go/data"
	"github.com/muandrew/battlecode-legacy-go/engine"
	"github.com/muandrew/battlecode-legacy-go/models"
	"github.com/muandrew/battlecode-legacy-go/utils"
)

//LzSite a lazy HTML impl of battlecode
type LzSite struct {
	templates *template.Template
}

const (
	failedUpload    = "Upload failed :/"
	failedChallenge = "Challenge failed T.T"
	maxBotsInGame   = 4
)

//NewInstance creates a new instance
func NewInstance() *LzSite {
	return &LzSite{
		templates: template.Must(template.ParseGlob("lazy/views/*.html")),
	}
}

//Init starts the LzSite
func (t *LzSite) Init(
	e *echo.Echo,
	a *auth.Auth,
	db data.Db,
	c *build.Ci,
	engines []engine.Engine,
) {
	e.Renderer = t
	g := e.Group("/lazy")
	g.Static("/static", "lazy/static")
	g.GET("/", getHello)
	g.GET("/login/", getLogin)
	loggedInGroup := g.Group("/loggedin")
	loggedInGroup.Use(a.AuthMiddleware)
	loggedInGroup.GET("/", wrapLoggedIn(engines))

	for _, engine := range engines {
		engineGroup := loggedInGroup.Group(fmt.Sprintf("/%s", engine.Competition()))
		engineGroup.GET("/", wrapEngineHome(engine, db))
		engineGroup.POST("/bot/upload/", wrapPostUpload(engine, c))
		engineGroup.POST("/bot/public/", wrapPostMakePublic(engine, db))
		engineGroup.GET("/bot/public/", wrapGetPublicBots(engine, db))
		engineGroup.POST("/map/upload/", wrapPostMapUpload(engine, c))
		engineGroup.POST("/challenge/", wrapPostChallenge(engine, db, c))
		engineGroup.POST("/challenge-game/", wrapPostChallengeGame(engine, db, c))
	}

	if utils.IsDev() {
		d := g.Group("/dev")
		d.GET("/login/", wrapGetDevLogin(a))
		d.GET("/script/", getDevScript)
		d.POST("/script/", postDevScript)
	}
}

//Render renders a page
func (t *LzSite) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func getHello(c echo.Context) error {
	return c.Render(http.StatusOK, "root", nil)
}

func getLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
}

func wrapGetDevLogin(a *auth.Auth) func(context echo.Context) error {
	return func(c echo.Context) error {
		a.GetUserWithApp(
			c,
			"dev",
			"#000000",
			func() *models.User {
				user, _ := models.CreateUser("Dev")
				return user
			},
		)
		return c.Redirect(http.StatusTemporaryRedirect, "/lazy/loggedin/")
	}
}

func getDevScript(c echo.Context) error {
	return c.Render(http.StatusOK, "dev_script", nil)
}

func postDevScript(c echo.Context) error {
	script := c.FormValue("script")
	utils.RunShell("bash", []string{"scripts/" + script})
	return c.Render(http.StatusOK, "dev_script", nil)
}

func wrapLoggedIn(engines []engine.Engine) func(context echo.Context) error {
	if len(engines) > 1 {
		engineNames := []string{}
		for _, engine := range engines {
			engineNames = append(engineNames, engine.Competition().AsString())
		}
		return func(c echo.Context) error {
			return c.Render(http.StatusOK, "choose_engine", engineNames)
		}
	} else {
		return func(c echo.Context) error {
			return c.Redirect(
				http.StatusTemporaryRedirect,
				fmt.Sprintf("/lazy/loggedin/%s/", engines[0].Competition()),
			)
		}
	}
}

func wrapEngineHome(engine engine.Engine, db data.Db) func(context echo.Context) error {
	return func(c echo.Context) error {
		uuid := auth.GetUuid(c)
		bots, _ := db.GetBots(uuid, 0, 5)
		matches, length := db.GetMatches(uuid, 0, 5)
		maps, length := db.GetBcMaps(uuid, 0, 5)
		data := map[string]interface{}{
			"name":           auth.GetName(c),
			"uuid":           uuid,
			"competition":    engine.Competition(),
			"latest_bots":    bots,
			"latest_matches": matches,
			"latest_maps":    maps,
			"length":         length,
		}

		return c.Render(http.StatusOK, "loggedin", data)
	}
}

//noinspection GoUnusedFunction
func debugResponse(c echo.Context, model interface{}) error {
	raw, _ := json.Marshal(model)
	return c.Render(http.StatusOK, "dev_debug", string(raw))
}

func wrapPostUpload(engine engine.Engine, ci *build.Ci) func(context echo.Context) error {
	return func(c echo.Context) error {
		uuid := auth.GetUuid(c)
		file, err := c.FormFile("file")
		if err != nil {
			return renderFailure(c, engine, failedUpload, err)
		}
		bot, err := models.CreateBot(
			models.NewCompetitor(models.CompetitorTypeUser, uuid),
			c.FormValue("package"),
			c.FormValue("note"),
			engine.Competition(),
			"",
		)
		if err != nil {
			return renderFailure(c, engine, failedUpload, err)
		}

		err = ci.UploadBotSource(file, bot)
		if err != nil {
			return renderFailure(c, engine, failedUpload, err)
		}
		ci.SubmitJob(bot)
		data := map[string]interface{}{
			"competition": engine.Competition(),
		}
		return c.Render(http.StatusOK, "uploaded", data)
	}
}

func wrapPostMakePublic(engine engine.Engine, db data.Db) func(context echo.Context) error {
	return func(c echo.Context) error {
		uuid := auth.GetUuid(c)
		botUuid := c.FormValue("botUuid")
		bot, err := db.SetPublicBot(uuid, botUuid)
		if err != nil {
			return renderFailure(c, engine, "failed to set bot as public: ", err)
		}
		return c.Render(http.StatusOK, "public_bot_set", bot)
	}
}

func wrapGetPublicBots(engine engine.Engine, db data.Db) func(ctx echo.Context) error {
	return func(c echo.Context) error {
		bots, _ := db.GetPublicBots(0, 10)
		data := map[string]interface{}{
			"bots":        bots,
			"competition": engine.Competition(),
		}
		return c.Render(http.StatusOK, "public_bots", data)
	}
}

func wrapPostMapUpload(engine engine.Engine, ci *build.Ci) func(context echo.Context) error {
	return func(c echo.Context) error {
		uuid := auth.GetUuid(c)
		file, err := c.FormFile("file")
		if err != nil {
			return renderFailure(c, engine, failedUpload, err)
		}

		bcMap, err := models.CreateBcMap(
			models.NewCompetitor(models.CompetitorTypeUser, uuid),
			file.Filename,
			c.FormValue("description"),
		)
		if err != nil {
			return renderFailure(c, engine, failedUpload, err)
		}

		err = ci.UploadMap(file, bcMap)
		if err != nil {
			return renderFailure(c, engine, failedUpload, err)
		}
		data := map[string]interface{}{
			"competition": engine.Competition(),
		}
		return c.Render(http.StatusOK, "uploaded", data)
	}
}

func wrapPostChallenge(engine engine.Engine, db data.Db, ci *build.Ci) func(context echo.Context) error {
	return func(c echo.Context) error {
		botUuid := c.FormValue("botUuid")
		oppUuid := c.FormValue("oppUuid")
		mapUuid := c.FormValue("mapUuid")

		ownBot := db.GetBot(botUuid)
		oppBot := db.GetBot(oppUuid)
		bcMap := db.GetBcMap(mapUuid)
		if ownBot == nil || oppBot == nil {
			return renderFailure(
				c,
				engine,
				failedChallenge,
				errors.New("Couldn't find two bots to play."),
			)
		}
		err := ci.RunMatch([]*models.Bot{ownBot, oppBot}, bcMap)

		if err != nil {
			return renderFailure(c, engine, failedChallenge, err)
		} else {
			data := map[string]interface{}{
				"competition": engine.Competition(),
			}
			return c.Render(http.StatusOK, "challenged", data)
		}
	}
}

func wrapPostChallengeGame(engine engine.Engine, db data.Db, ci *build.Ci) func(context echo.Context) error {
	return func(c echo.Context) error {
		uuid := auth.GetUuid(c)
		name := c.FormValue("name")
		description := c.FormValue("description")
		formBotUuids := c.FormValue("botUuids")
		mapUuid := c.FormValue("mapUuid")

		botUuids := strings.Split(formBotUuids, ",")
		if len(botUuids) > maxBotsInGame {
			return renderFailure(
				c,
				engine,
				failedChallenge,
				errors.New(fmt.Sprintf(
					"Too many fights the server will explode! The current max is %d",
					maxBotsInGame)))
		}
		bots := make([]*models.Bot, len(botUuids), len(botUuids))
		for i, botUuid := range botUuids {
			bot := db.GetBot(botUuid)
			if bot == nil {
				return renderFailure(
					c,
					engine,
					failedChallenge,
					errors.New(fmt.Sprintf("Couldn't find bot %s", botUuid)))
			} else {
				bots[i] = bot
			}
		}

		bcMap := db.GetBcMap(mapUuid)
		err := ci.RunGame(
			models.NewCompetitor(models.CompetitorTypeUser, uuid),
			name,
			description,
			bots,
			bcMap,
		)

		if err != nil {
			return renderFailure(c, engine, failedChallenge, err)
		} else {
			data := map[string]interface{}{
				"competition": engine.Competition(),
			}
			return c.Render(http.StatusOK, "challenged", data)
		}
	}
}

func renderFailure(
	context echo.Context,
	engine engine.Engine,
	title string,
	err error,
) error {
	data := map[string]interface{}{
		"title":       title,
		"error":       err,
		"competition": engine.Competition(),
	}
	return context.Render(http.StatusOK, "failure", data)
}
