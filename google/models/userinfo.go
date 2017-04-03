package models

type UserInfo struct {
	Id            string
	Email         string
	VerifiedEmail bool
	Name          string
	GivenName     string
	FamilyName    string
	Link          string
	Picture       string
	Gender        string
	Locale        string
}
