package oauth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"strings"
	"errors"
	"github.com/muandrew/battlecode-ladder/utils"
	"github.com/labstack/echo"
	"net/http"
	"golang.org/x/net/context"
	gmodels "github.com/muandrew/battlecode-ladder/google/models"
	"github.com/muandrew/battlecode-ladder/models"
	"github.com/muandrew/battlecode-ladder/auth"
)

const (
	oauthID = "_ID"
	oauthSecret = "_SECRET"
)

type OAMap map[string]*OAConfig

type OAConfig struct {
	App    string
	Config *oauth2.Config
}

func NewOAConfig(redirectPath string, app string, scopes []string, endpoint oauth2.Endpoint, fail func()) *OAConfig {
	return &OAConfig{App:app, Config:  &oauth2.Config{
		RedirectURL:  redirectPath + app + "/",
		ClientID:     getKeyForApp(oauthID, app, fail),
		ClientSecret: getKeyForApp(oauthSecret, app, fail),
		Scopes:       scopes,
		Endpoint:     endpoint,
	}}
}

func Init(e *echo.Echo, address string, prefix string, auth *auth.Auth) (OAMap, error) {

	redirectPath := address + prefix + "callback/"

	success := true
	fail := func() {
		success = false
	}

	auths := []*OAConfig{
		NewOAConfig(
			redirectPath,
			"google",
			[]string{
				"https://www.googleapis.com/auth/userinfo.profile",
				"https://www.googleapis.com/auth/userinfo.email",
			},
			google.Endpoint,
			fail,
		),
	}

	m := make(OAMap)
	for _, item := range auths {
		m[item.App] = item
	}
	if success {
		e.GET(prefix + "login/:app/", getGetLogin(m))
		e.GET(prefix + "callback/:app/", getGetCallback(m, auth))
		return m, nil
	} else {
		return nil, errors.New("Initialization Failed.")
	}
}

func getGetLogin(oamap OAMap) func(echo.Context) error {
	return func (c echo.Context) error {
		app := c.Param("app")
		config := oamap[app]
		return c.Redirect(http.StatusTemporaryRedirect, config.Config.AuthCodeURL("todo"))
	}
}

func getGetCallback(oamap OAMap, auth *auth.Auth) func(echo.Context) error {
	return func (c echo.Context) error {
		app := c.Param("app")
		config := oamap[app]
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

		auth.GetUserWithApp(c, app, info.ID, func(user *models.User) *models.User {
			user.Name = info.Name
			return user
		})
		return c.JSON(http.StatusOK, info)
	}
}

func getKeyForApp(key string, app string, fail func()) string {
	return utils.GetRequiredEnv("OAUTH_" + strings.ToUpper(app) + key, fail)
}
