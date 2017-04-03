package oauth

import (
	"golang.org/x/oauth2"
	"os"
	"golang.org/x/oauth2/google"
	"fmt"
	"strings"
	"errors"
)

const (
	oauthID = "_ID"
	oauthSecret = "_SECRET"
)

type OAConfig struct {
	App    string
	Config *oauth2.Config
}

func NewOAConfig(redirectPath string, app string, scopes []string, endpoint oauth2.Endpoint, fail func()) *OAConfig {
	return &OAConfig{App:app, Config:  &oauth2.Config{
		RedirectURL:  redirectPath + "/" + app + "/",
		ClientID:     getKeyForApp(oauthID, app, fail),
		ClientSecret: getKeyForApp(oauthSecret, app, fail),
		Scopes:       scopes,
		Endpoint:     endpoint,
	}}
}

func Init(redirectPath string) (map[string]*OAConfig, error) {
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

	m := make(map[string]*OAConfig)
	for _, item := range auths {
		m[item.App] = item
	}
	if (success) {
		return m, nil
	} else {
		return nil, errors.New("Initialization Failed.")
	}

}

func getKeyForApp(key string, app string, fail func()) string {
	return getRequiredEnv("OAUTH_" + strings.ToUpper(app) + key, fail)
}

func getRequiredEnv(key string, fail func()) string {
	value := os.Getenv(key)
	if value == "" {
		fmt.Printf("Your env variable %s was not configured.\n", key)
		fail()
	}
	return value
}