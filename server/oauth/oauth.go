package oauth

import (
	"errors"
	"github.com/labstack/echo"
	"github.com/muandrew/battlecode-ladder/auth"
	gmodels "github.com/muandrew/battlecode-ladder/google/models"
	"github.com/muandrew/battlecode-ladder/models"
	"github.com/muandrew/battlecode-ladder/utils"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"strings"
)

const (
	oauthID     = "_ID"
	oauthSecret = "_SECRET"
)

type OAMap map[string]*OAConfig

type getUserI func(c echo.Context, authp *auth.Auth, accessToken string) (*models.User, error)

type OAConfig struct {
	App     string
	Config  *oauth2.Config
	GetUser getUserI
}

func NewOAConfig(redirectPath string, app string, getUser getUserI, scopes []string, endpoint oauth2.Endpoint, fail func()) *OAConfig {
	return &OAConfig{
		App:     app,
		GetUser: getUser,
		Config: &oauth2.Config{
			RedirectURL:  redirectPath + app + "/",
			ClientID:     getKeyForApp(oauthID, app, fail),
			ClientSecret: getKeyForApp(oauthSecret, app, fail),
			Scopes:       scopes,
			Endpoint:     endpoint,
		}}
}

func Init(e *echo.Echo, address string, prefix string, authp *auth.Auth) (OAMap, error) {

	redirectPath := address + prefix + "callback/"

	success := true
	fail := func() {
		success = false
	}

	auths := []*OAConfig{
		NewOAConfig(
			redirectPath,
			"google",
			func(c echo.Context, authp *auth.Auth, accessToken string) (*models.User, error) {
				response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
				if err != nil {
					return nil, err
				}
				info := new(gmodels.UserInfo)
				utils.ReadBody(response, info)
				if err != nil {
					return nil, err
				}
				usr := authp.GetUserWithApp(c, "google", info.ID, func() *models.User {
					mUser, _ := models.CreateUser(info.Name)
					return mUser
				})
				return usr, nil
			},
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
		e.GET(prefix+"login/:app/", getGetLogin(m))
		e.GET(prefix+"callback/:app/", getGetCallback(m, authp))
		return m, nil
	} else {
		return nil, errors.New("Initialization Failed.")
	}
}

func getGetLogin(oamap OAMap) func(echo.Context) error {
	return func(c echo.Context) error {
		app := c.Param("app")
		config := oamap[app]
		state := uuid.NewV4().String()
		return c.Redirect(http.StatusTemporaryRedirect, config.Config.AuthCodeURL(state))
	}
}

func getGetCallback(oamap OAMap, authp *auth.Auth) func(echo.Context) error {
	return func(c echo.Context) error {
		app := c.Param("app")
		config := oamap[app]
		//todo state:= c.Param("state")
		token, err := config.Config.Exchange(context.TODO(), c.FormValue("code"))
		if err != nil {
			return err
		}
		_, err = config.GetUser(c, authp, token.AccessToken)
		if err != nil {
			return err
		} else {
			//return c.JSON(http.StatusOK, user)
			//todo stop being so lazy
			return c.Redirect(http.StatusTemporaryRedirect, "/lazy/loggedin/")
		}
	}
}

func getKeyForApp(key string, app string, fail func()) string {
	return utils.GetRequiredEnv("OAUTH_"+strings.ToUpper(app)+key, fail)
}
