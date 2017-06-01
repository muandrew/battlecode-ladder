package models

type Match struct{
	Uuid string
	Bots []*Bot
	CompletedTimestamp int64
}

func CreateMatch(uuid string) *Match {
	return &Match{Uuid:uuid}
}
