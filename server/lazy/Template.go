package lazy

import (
	"io"
	"github.com/labstack/echo"
	"html/template"
	"net/http"
	"github.com/muandrew/battlecode-ladder/auth"
	"os"
	"fmt"
	"github.com/muandrew/battlecode-ladder/models"
	"github.com/muandrew/battlecode-ladder/build"
	"github.com/muandrew/battlecode-ladder/data"
)

type Template struct {
	templates *template.Template
}

func NewInstance() *Template {
	return &Template{
		templates: template.Must(template.ParseGlob("lazy/views/*.html")),
	}
}

func (t *Template) Init(e *echo.Echo, auth echo.MiddlewareFunc, db data.Db, c *build.Ci) {
	e.Renderer = t
	g := e.Group("/lazy")
	g.Static("/static", "lazy/static")
	g.GET("/", getHello)
	g.GET("/login/", getLogin)
	r := g.Group("/loggedin")
	r.Use(auth)
	r.GET("/", wrapGetLoggedIn(db))
	r.POST("/upload/", wrapPostUpload(c))
	r.POST("/challenge/", wrapPostChallenge(db, c))
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func getHello(c echo.Context) error {
	return c.Render(http.StatusOK, "root", nil)
}

func getLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
}

func wrapGetLoggedIn(db data.Db) func(context echo.Context) error {
	return func(c echo.Context) error {
		uuid := auth.GetUuid(c)
		bots, _ := db.GetBots(uuid, 0, 5)
		matches, length := db.GetMatches(uuid, 0, 5)
		model := map[string]interface{}{
			"name":           auth.GetName(c),
			"uuid":           uuid,
			"latest_bots":    bots,
			"latest_matches": matches,
			"length":         length,
		}
		return c.Render(http.StatusOK, "loggedin", model)
	}
}

func wrapPostUpload(ci *build.Ci) func(context echo.Context) error {
	return func(c echo.Context) error {
		uuid := auth.GetUuid(c)
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}
		bot, err := models.CreateBot(
			models.NewCompetitor(models.CompetitorTypeUser, uuid),
			c.FormValue("package"),
			c.FormValue("name"),
			c.FormValue("description"),
		)
		if err != nil {
			return err
		}

		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Destination
		prefix := "bl-data/bot/" + bot.Uuid
		os.MkdirAll(prefix, 0755)
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

		//todo respond to error
		ci.SubmitJob(bot)

		return c.Render(http.StatusOK, "uploaded", auth.GetName(c))
	}
}

func wrapPostChallenge(db data.Db, ci *build.Ci) func(context echo.Context) error {
	return func(c echo.Context) error {
		botUuid := c.FormValue("botUuid")
		oppUuid := c.FormValue("oppUuid")

		ownBot := db.GetBot(botUuid)
		oppBot := db.GetBot(oppUuid)

		if ownBot != nil && oppBot != nil {
			ci.RunMatch(ownBot, oppBot)
			return c.Render(http.StatusOK, "challenged", nil)
		} else {
			return c.Render(http.StatusOK, "challenge_failed", nil)
		}
	}
}
