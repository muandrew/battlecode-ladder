package models

type UserInfo struct {
	Email         string `json:"email"`
	FamilyName    string `json:"family_name"`
	Gender        string `json:"gender"`
	GivenName     string `json:"given_name"`
	ID            string `json:"id"`
	Link          string `json:"link"`
	Locale        string `json:"locale"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}
