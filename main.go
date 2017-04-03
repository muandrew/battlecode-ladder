package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/muandrew/battlecode-ladder/oauth"
)

var authProviders map[string]*oauth.OAConfig

func main() {
	providers, err := oauth.Init("http://localhost:8080/callback")
	if (err != nil) {
		return
	}
	authProviders = providers

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/login/:app/", getLogin)
	e.GET("/callback/:app/", getAuthorize)
	e.Logger.Fatal(e.Start(":8080"))
}

func getLogin(c echo.Context) error {
	app := c.Param("app")
	config := authProviders[app]
	return c.Redirect(http.StatusTemporaryRedirect, config.Config.AuthCodeURL("todo"))
}

func getAuthorize(c echo.Context) error {
	app := c.Param("app")
	return c.Redirect(http.StatusTemporaryRedirect, "https://www." + app + ".com")
}