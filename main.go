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
)

var authProviders map[string]*oauth.OAConfig
var jwtSecret []byte

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

	e.Static("/", "client")
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

	_, err = http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	//response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Some error has occured2.")
	}

	//defer response.Body.Close()
	//contents, _ := ioutil.ReadAll(response.Body)
	//return c.String(http.StatusOK, fmt.Sprintf("Content: %s\n", contents))
	//fmt.Println(c.String(http.StatusOK, fmt.Sprintf("Content: %s\n", contents)))
	t := setJwtInCookie(c, "Test User")
	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
	//return c.Redirect(http.StatusTemporaryRedirect, "https://www." + app + ".com")
}

func setJwtInCookie(c echo.Context, name string) string {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = name
	claims["admin"] = true
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
	t := setJwtInCookie(c, "John Doe")
	return c.String(http.StatusOK, "set: " + t)
}

func getInspect(c echo.Context) error {
	jwtCookie, _ := c.Cookie(jwtCookieName)
	name := ""
	if jwtCookie != nil {
		name = jwtCookie.Value
	}
	//user := c.Get("user").(*jwt.Token)
	//claims := user.Claims.(jwt.MapClaims)
	//name := claims["name"].(string)
	return c.String(http.StatusOK, "get: " + name)
}
