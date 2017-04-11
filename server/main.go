package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/muandrew/battlecode-ladder/oauth"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/middleware"
	"time"
	"fmt"
	"github.com/muandrew/battlecode-ladder/utils"
	"github.com/joho/godotenv"
	"log"
	"golang.org/x/net/context"
	gmodels "github.com/muandrew/battlecode-ladder/google/models"
	"github.com/muandrew/battlecode-ladder/db"
	"github.com/muandrew/battlecode-ladder/models"
)

var authProviders map[string]*oauth.OAConfig
var jwtSecret []byte
var data db.Db

const jwtCookieName = "xbclauth"

func main() {
	err := godotenv.Load("secrets.sh")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	providers, err := oauth.Init("http://localhost:8080/callback")
	if err != nil {
		return
	}
	authProviders = providers

	data = db.NewMemDb()

	initSuccess := true
	jwtSecret = []byte(utils.GetRequiredEnv("JWT_SECRET", func() {
		initSuccess = false
	}))
	if !initSuccess {
		fmt.Println("Init failed.")
		return
	}

	e := echo.New()
	e.GET("/cookie/", getCookie)
	e.GET("/inspect/", getInspect)
	e.GET("/login/:app/", getLogin)
	e.GET("/callback/:app/", getCallback)

	r := e.Group("/restricted")
	config := middleware.DefaultJWTConfig
	config.SigningKey = jwtSecret
	config.TokenLookup = "cookie:" + jwtCookieName
	r.Use(middleware.JWTWithConfig(config))
	r.GET("/", restricted)

	e.Static("/", "../client")
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
	token, err := config.Config.Exchange(context.TODO(), c.FormValue("code"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Some error has occured.")
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Some error has occured2.")
	}
	info := new(gmodels.UserInfo)
	utils.ReadBody(response, info)
	user := data.GetUserWithApp(app, info.ID, func() *models.User {
		user := models.CreateUserWithNewUuid()
		user.Name = info.Name
		return user
	})
	setJwtInCookie(c, user)
	return c.JSON(http.StatusOK, info)
}

func setJwtInCookie(c echo.Context, user *models.User) string {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["uuid"] = user.Uuid
	claims["name"] = user.Name
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString(jwtSecret)
	if err != nil {
		return ""
	}

	cookie := new(http.Cookie)
	cookie.Name = jwtCookieName
	cookie.Value = t
	cookie.Expires = time.Now().Add(time.Hour * 24)
	cookie.Path = "/"
	c.SetCookie(cookie)
	return t
}

func getCookie(c echo.Context) error {
	t := setJwtInCookie(c, models.UserDummy)
	return c.String(http.StatusOK, "set: " + t)
}

func getInspect(c echo.Context) error {
	return c.JSON(http.StatusOK, data.GetAllUsers())
}
