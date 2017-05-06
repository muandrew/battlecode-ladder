package lazy

import (
	"io"
	"github.com/labstack/echo"
	"html/template"
	"net/http"
	"github.com/muandrew/battlecode-ladder/auth"
	"os"
	"fmt"
)

type Template struct {
	templates *template.Template
}

func NewInstance() *Template {
	return &Template{
		templates: template.Must(template.ParseGlob("lazy/views/*.html")),
	}
}

func (t *Template) Init(e *echo.Echo, auth echo.MiddlewareFunc) {
	e.Renderer = t
	g := e.Group("/lazy")
	g.GET("/", getHello)
	g.GET("/login/",getLogin)
	r := g.Group("/loggedin")
	r.Use(auth)
	r.GET("/", getLoggedIn)
	r.POST("/upload/", postUpload)
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

func getLoggedIn(c echo.Context) error {
	return c.Render(http.StatusOK, "loggedin", auth.GetName(c))
}

func postUpload(c echo.Context) error {
	userUuid := auth.GetUuid(c)
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	prefix := "user/"+userUuid
	//todo not be lazy
	os.MkdirAll(prefix, 0777)
	dst, err := os.Create(prefix + "/test.txt")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.Render(http.StatusOK, "uploaded", auth.GetName(c))
}
