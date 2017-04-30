package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/muandrew/battlecode-ladder/oauth"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/middleware"
	"fmt"
	"github.com/muandrew/battlecode-ladder/utils"
	"github.com/joho/godotenv"
	"log"
	"github.com/muandrew/battlecode-ladder/db"
	"github.com/muandrew/battlecode-ladder/auth"
)

var data db.Db

func main() {
	err := godotenv.Load("secrets.sh")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	data = db.NewMemDb()

	initSuccess := true
	onFail := func() {
		initSuccess = false
	}
	jwtSecret := []byte(utils.GetRequiredEnv("JWT_SECRET", onFail))
	rootAddress := utils.GetRequiredEnv("ROOT_ADDRESS", onFail)
	port := utils.GetRequiredEnv("PORT", onFail)
	if !initSuccess {
		fmt.Println("Init failed.")
		return
	}
	authentication := auth.NewAuth(data, jwtSecret)

	e := echo.New()
	_, err = oauth.Init(e, rootAddress, "/", authentication)
	if err != nil {
		return
	}

	e.GET("/inspect/", getInspect)
	r := e.Group("/restricted")
	r.Use(middleware.JWTWithConfig(authentication.JwtConfig))
	r.GET("/", restricted)

	e.Static("/", "../client")
	e.Logger.Fatal(e.Start(":"+port))
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome " + name + "!")
}

func getInspect(c echo.Context) error {
	return c.JSON(http.StatusOK, data.GetAllUsers())
}
