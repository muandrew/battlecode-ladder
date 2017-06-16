package models

const (
	CompetitorTypeUser = "user"
	CompetitorTypeTeam = "team"
)

type Competitor struct {
	Type string
	Uuid string
}

func NewCompetitor(competitorType string, uuid string) *Competitor {
	return &Competitor{competitorType, uuid}
}
