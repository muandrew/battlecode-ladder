package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/muandrew/battlecode-legacy-go/auth"
	"github.com/muandrew/battlecode-legacy-go/build"
	"github.com/muandrew/battlecode-legacy-go/data"
	"github.com/muandrew/battlecode-legacy-go/graphql"
	"github.com/muandrew/battlecode-legacy-go/lazy"
	"github.com/muandrew/battlecode-legacy-go/oauth"
	"github.com/muandrew/battlecode-legacy-go/utils"
)

func main() {
	utils.InitMainEnv()

	initSuccess := true
	onFail := func() {
		initSuccess = false
	}
	jwtSecret := []byte(utils.GetRequiredEnv("JWT_SECRET", onFail))
	db, err := data.NewRdsDb(utils.GetRequiredEnv("REDIS_ADDRESS", onFail))
	if err != nil {
		log.Fatalf("Failed to init redis: %s", err)
	}
	rootAddress := utils.GetRequiredEnv("ROOT_ADDRESS", onFail)
	port := utils.GetRequiredEnv("PORT", onFail)
	if !initSuccess {
		log.Fatalf("Init failed.")
	}
	authentication := auth.NewAuth(db, jwtSecret)

	e := echo.New()
	_, err = oauth.Init(e, rootAddress, "/", authentication)
	if err != nil {
		log.Fatalf("Failed to init oauth: %s", err)
	}

	ci, err := build.NewCi(db)
	if err != nil {
		log.Fatalf("Failed to init Ci: %s", err)
	}
	defer ci.Close()

	t := lazy.NewInstance()
	t.Init(e, authentication, db, ci)
	if utils.IsDev() {
		err = graphql.Init(db, e)
		if err != nil {
			log.Fatalf("Failed to init GraphQL: %s", err)
		}
	}
	e.Static("/bc17", "static/viewer/bc17/res")
	e.Static("/viewer/bc17", "static/viewer/bc17")
	e.Static("/replay", ci.GetDirMatches())
	e.GET("*", getRedirected)
	e.Logger.Fatal(e.Start(":" + port))
}

func getRedirected(c echo.Context) error {
	return c.Redirect(http.StatusTemporaryRedirect, "/lazy/")
}
