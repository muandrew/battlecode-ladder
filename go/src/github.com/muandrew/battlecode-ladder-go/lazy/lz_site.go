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
	"github.com/muandrew/battlecode-ladder-go/auth"
	"github.com/muandrew/battlecode-ladder-go/build"
	"github.com/muandrew/battlecode-ladder-go/data"
	"github.com/muandrew/battlecode-ladder-go/models"
	"github.com/muandrew/battlecode-ladder-go/utils"
)

type LzSite struct {
	templates *template.Template
}

const (
	failedUpload    = "Upload failed :/"
	failedChallenge = "Challenge failed T.T"
	maxBotsInGame   = 4
)

func NewInstance() *LzSite {
	return &LzSite{
		templates: template.Must(template.ParseGlob("lazy/views/*.html")),
	}
}

func (t *LzSite) Init(e *echo.Echo, a *auth.Auth, db data.Db, c *build.Ci) {
	e.Renderer = t
	g := e.Group("/lazy")
	g.Static("/static", "lazy/static")
	g.GET("/", getHello)
	g.GET("/login/", getLogin)
	r := g.Group("/loggedin")
	r.Use(a.AuthMiddleware)
	r.GET("/", wrapGetLoggedIn(db))
	r.POST("/bot/upload/", wrapPostUpload(c))
	r.POST("/bot/public/", wrapPostMakePublic(db))
	r.GET("/bot/public/", wrapGetPublicBots(db))
	r.POST("/map/upload/", wrapPostMapUpload(c))
	r.POST("/challenge/", wrapPostChallenge(db, c))
	r.POST("/challenge-game/", wrapPostChallengeGame(db, c))

	if utils.IsDev() {
		d := g.Group("/dev")
		d.GET("/login/", wrapGetDevLogin(a))
		d.GET("/script/", getDevScript)
		d.POST("/script/", postDevScript)
	}
}

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

func wrapGetLoggedIn(db data.Db) func(context echo.Context) error {
	return func(c echo.Context) error {
		uuid := auth.GetUuid(c)
		bots, _ := db.GetBots(uuid, 0, 5)
		matches, length := db.GetMatches(uuid, 0, 5)
		maps, length := db.GetBcMaps(uuid, 0, 5)
		model := map[string]interface{}{
			"name":           auth.GetName(c),
			"uuid":           uuid,
			"latest_bots":    bots,
			"latest_matches": matches,
			"latest_maps":    maps,
			"length":         length,
		}

		return c.Render(http.StatusOK, "loggedin", model)
	}
}

//noinspection GoUnusedFunction
func debugResponse(c echo.Context, model interface{}) error {
	raw, _ := json.Marshal(model)
	return c.Render(http.StatusOK, "dev_debug", string(raw))
}

func wrapPostUpload(ci *build.Ci) func(context echo.Context) error {
	return func(c echo.Context) error {
		uuid := auth.GetUuid(c)
		file, err := c.FormFile("file")
		if err != nil {
			return renderFailure(c, failedUpload, err)
		}
		bot, err := models.CreateBot(
			models.NewCompetitor(models.CompetitorTypeUser, uuid),
			c.FormValue("package"),
			c.FormValue("note"),
			models.CompetitionBC17,
		)
		if err != nil {
			return renderFailure(c, failedUpload, err)
		}

		err = ci.UploadBotSource(file, bot)
		if err != nil {
			return renderFailure(c, failedUpload, err)
		}
		ci.SubmitJob(bot)
		return c.Render(http.StatusOK, "uploaded", auth.GetName(c))
	}
}

func wrapPostMakePublic(db data.Db) func(context echo.Context) error {
	return func(c echo.Context) error {
		uuid := auth.GetUuid(c)
		botUuid := c.FormValue("botUuid")
		bot, err := db.SetPublicBot(uuid,botUuid)
		if err != nil {
			return renderFailure(c, "failed to set bot as public: ", err)
		}
		return c.Render(http.StatusOK, "public_bot_set", bot)
	}
}

func wrapGetPublicBots(db data.Db) func(ctx echo.Context) error {
	return func(c echo.Context) error {
		bots,_ := db.GetPublicBots(0,10)
		return c.Render(http.StatusOK, "public_bots", bots)
	}
}

func wrapPostMapUpload(ci *build.Ci) func(context echo.Context) error {
	return func(c echo.Context) error {
		uuid := auth.GetUuid(c)
		file, err := c.FormFile("file")
		if err != nil {
			return renderFailure(c, failedUpload, err)
		}

		bcMap, err := models.CreateBcMap(
			models.NewCompetitor(models.CompetitorTypeUser, uuid),
			file.Filename,
			c.FormValue("description"),
		)
		if err != nil {
			return renderFailure(c, failedUpload, err)
		}

		err = ci.UploadMap(file, bcMap)
		if err != nil {
			return renderFailure(c, failedUpload, err)
		}
		return c.Render(http.StatusOK, "uploaded", auth.GetName(c))
	}
}

func wrapPostChallenge(db data.Db, ci *build.Ci) func(context echo.Context) error {
	return func(c echo.Context) error {
		botUuid := c.FormValue("botUuid")
		oppUuid := c.FormValue("oppUuid")
		mapUuid := c.FormValue("mapUuid")

		ownBot := db.GetBot(botUuid)
		oppBot := db.GetBot(oppUuid)
		bcMap := db.GetBcMap(mapUuid)
		if ownBot == nil || oppBot == nil {
			return errors.New("Couldn't find two bots to play.")
		}
		err := ci.RunMatch([]*models.Bot{ownBot, oppBot}, bcMap)

		if err != nil {
			return renderFailure(c, failedChallenge, err)
		} else {
			return c.Render(http.StatusOK, "challenged", nil)
		}
	}
}

func wrapPostChallengeGame(db data.Db, ci *build.Ci) func(context echo.Context) error {
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
			return renderFailure(c, failedChallenge, err)
		} else {
			return c.Render(http.StatusOK, "challenged", nil)
		}
	}
}

func renderFailure(context echo.Context, title string, err error) error {
	model := map[string]interface{}{
		"title": title,
		"error": err,
	}
	return context.Render(http.StatusOK, "failure", model)
}
