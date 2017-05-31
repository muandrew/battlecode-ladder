package models

type Match struct{
	Uuid string
	Bots []*Bot
}

func CreateMatch(uuid string) *Match {
	return &Match{Uuid:uuid}
}
