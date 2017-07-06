package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/muandrew/battlecode-ladder/oauth"
	"fmt"
	"github.com/muandrew/battlecode-ladder/utils"
	"github.com/joho/godotenv"
	"log"
	"github.com/muandrew/battlecode-ladder/auth"
	"github.com/muandrew/battlecode-ladder/lazy"
	"github.com/muandrew/battlecode-ladder/build"
	"github.com/muandrew/battlecode-ladder/data"
)

var db data.Db

func main() {
	err := godotenv.Load("bcl-env.sh")
	utils.Initialize("BCL_")
	if err != nil {
		log.Fatalf("Error loading .env file; err: %q", err)
	}

	initSuccess := true
	onFail := func() {
		initSuccess = false
	}
	jwtSecret := []byte(utils.GetRequiredEnv("JWT_SECRET", onFail))
	db, err = data.NewRdsDb(utils.GetRequiredEnv("REDIS_ADDRESS", onFail))
	if err != nil {
		log.Fatalf("Failed to init redis: %s", err)
	}
	rootAddress := utils.GetRequiredEnv("ROOT_ADDRESS", onFail)
	port := utils.GetRequiredEnv("PORT", onFail)
	if !initSuccess {
		fmt.Println("Init failed.")
		return
	}
	authentication := auth.NewAuth(db, jwtSecret)

	e := echo.New()
	_, err = oauth.Init(e, rootAddress, "/", authentication)
	if err != nil {
		return
	}
	//todo error handling
	ci, err := build.NewCi(db)
	if err != nil {
		log.Fatalf("Failed to init Ci: %s", err)
	}
	defer ci.Close()

	t := lazy.NewInstance()
	t.Init(e, authentication, db, ci)

	r := e.Group("/restricted")
	r.Use(authentication.AuthMiddleware)
	r.GET("/", restricted)

	e.Static("/bc17", "viewer/bc17")
	e.Static("/viewer", "viewer")
	e.Static("/replay", ci.GetDirMatches())
	e.GET("*", getRedirected)
	e.Logger.Fatal(e.Start(":" + port))
}

func restricted(c echo.Context) error {
	name := auth.GetName(c)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func getRedirected(c echo.Context) error {
	return c.Redirect(http.StatusTemporaryRedirect, "/lazy/")
}
