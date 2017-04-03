package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/muandrew/battlecode-ladder/oauth"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"github.com/labstack/echo/middleware"
	"time"
	"io/ioutil"
	"fmt"
	"github.com/muandrew/battlecode-ladder/utils"
)

var authProviders map[string]*oauth.OAConfig
var jwtSecret []byte

func main() {
	providers, err := oauth.Init("http://localhost:8080/callback")
	if (err != nil) {
		return
	}
	authProviders = providers

	initSuccess := true
	jwtSecret := []byte(utils.GetRequiredEnv("JWT_SECRET", func() {
		initSuccess = false
	}))
	if (!initSuccess) {
		fmt.Println("Init failed.")
		return
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/login/:app/", getLogin)
	e.GET("/callback/:app/", getCallback)

	r := e.Group("/restricted")
	r.Use(middleware.JWT(jwtSecret))
	r.GET("", restricted)

	e.Logger.Fatal(e.Start(":8080"))
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome " + name + "!")
}

func getLogin(c echo.Context) error {
	app := c.Param("app")
	config := authProviders[app]
	return c.Redirect(http.StatusTemporaryRedirect, config.Config.AuthCodeURL("todo"))
}

func getCallback(c echo.Context) error {
	app := c.Param("app")
	config := authProviders[app]
	token, err := config.Config.Exchange(oauth2.NoContext, c.FormValue("code"))
	if (err == nil) {
		response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
		if (err == nil) {
			defer response.Body.Close()
			contents, _ := ioutil.ReadAll(response.Body)
			//return c.String(http.StatusOK, fmt.Sprintf("Content: %s\n", contents))
			fmt.Println(c.String(http.StatusOK, fmt.Sprintf("Content: %s\n", contents)))
			return user(c, "Test User")
		} else {
			return c.String(http.StatusInternalServerError, "Some error has occured2.")
		}
	} else {
		return c.String(http.StatusInternalServerError, "Some error has occured.")
	}
	//return c.Redirect(http.StatusTemporaryRedirect, "https://www." + app + ".com")
}

func user(c echo.Context, name string) error {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = name
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString(jwtSecret)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}
