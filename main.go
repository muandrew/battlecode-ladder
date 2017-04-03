package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/muandrew/battlecode-ladder/oauth"
	"io/ioutil"
	"fmt"
	"golang.org/x/oauth2"
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
	config := authProviders[app]
	token, err := config.Config.Exchange(oauth2.NoContext, c.FormValue("code"))
	if (err == nil) {
		response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
		if (err == nil) {
			defer response.Body.Close()
			contents, _ := ioutil.ReadAll(response.Body)
			return c.String(http.StatusOK, fmt.Sprintf("Content: %s\n", contents))
		} else {
			return c.String(http.StatusInternalServerError, "Some error has occured2.")
		}
	} else {
		return c.String(http.StatusInternalServerError, "Some error has occured.")
	}
	//return c.Redirect(http.StatusTemporaryRedirect, "https://www." + app + ".com")
}