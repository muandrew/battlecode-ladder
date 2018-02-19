package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/muandrew/battlecode-ladder-go/data"
	"github.com/muandrew/battlecode-ladder-go/models"
	"net/http"
	"time"
)

const jwtCookieName = "xbclauth"

type Auth struct {
	db             data.Db
	jwtSecret      []byte
	AuthMiddleware echo.MiddlewareFunc
}

func NewAuth(db data.Db, jwtSecret []byte) *Auth {
	config := middleware.DefaultJWTConfig
	config.SigningKey = jwtSecret
	config.TokenLookup = "cookie:" + jwtCookieName
	authMiddleware := middleware.JWTWithConfig(config)
	return &Auth{
		db:             db,
		jwtSecret:      jwtSecret,
		AuthMiddleware: authMiddleware,
	}
}

func (auth Auth) GetUserWithApp(c echo.Context, app string, appUuid string, setupUser models.SetupNewUser) *models.User {
	user := auth.db.GetUserWithApp(app, appUuid, setupUser)
	auth.setJwtInCookie(c, user)
	return user
}

func (auth Auth) setJwtInCookie(c echo.Context, user *models.User) string {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["uuid"] = user.Uuid
	claims["name"] = user.Name
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString(auth.jwtSecret)
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

func GetName(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["name"].(string)
}

func GetUuid(c echo.Context) string {
	user, ok := c.Get("user").(*jwt.Token)
	if ok {
		claims := user.Claims.(jwt.MapClaims)
		return claims["uuid"].(string)
	} else {
		return ""
	}
}
