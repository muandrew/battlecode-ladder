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
	"github.com/muandrew/battlecode-ladder/db"
	"github.com/muandrew/battlecode-ladder/build"
)

type Template struct {
	templates *template.Template
}

func NewInstance() *Template {
	return &Template{
		templates: template.Must(template.ParseGlob("lazy/views/*.html")),
	}
}

func (t *Template) Init(e *echo.Echo, auth echo.MiddlewareFunc, data db.Db, c *build.Ci) {
	e.Renderer = t
	g := e.Group("/lazy")
	g.GET("/", getHello)
	g.GET("/login/",getLogin)
	r := g.Group("/loggedin")
	r.Use(auth)
	r.GET("/", wrapGetLoggedIn(data))
	r.POST("/upload/", wrapPostUpload(data, c))
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

func wrapGetLoggedIn(data db.Db) func(context echo.Context) error {
	return func(c echo.Context) error{
		model := map[string]interface{}{
			"name": auth.GetName(c),
			"latest_complete_build": data.GetLatestCompletedBot(auth.GetUuid(c)),
		}
		return c.Render(http.StatusOK, "loggedin", model)
	}
}

func wrapPostUpload(data db.Db, ci *build.Ci) func(context echo.Context) error {
	return func (c echo.Context) error {
		      userUuid := auth.GetUuid(c)
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}
		botName := c.FormValue("name")
		botPackage := c.FormValue("package")

		bot := models.CreateBotWithNewUuidAndUserUuid(userUuid)
		bot.Name = botName
		bot.Package = botPackage

		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Destination
		prefix := "bl-data/user/"+userUuid+"/bot/"+bot.Uuid
		//todo not be lazy
		os.MkdirAll(prefix, 0777)
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
		data.EnqueueBot(bot)
		ci.SubmitJob(bot)

		return c.Render(http.StatusOK, "uploaded", auth.GetName(c))
	}
}

