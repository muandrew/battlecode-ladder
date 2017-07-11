package models

const (
	CompetitorTypeUser = CompetitorType("user")
	CompetitorTypeTeam = CompetitorType("team")
)

type CompetitorType string

func (ct CompetitorType) String() string {
	return string(ct)
}

type Competitor struct {
	Type CompetitorType
	Uuid string
}

func NewCompetitor(competitorType CompetitorType, uuid string) *Competitor {
	return &Competitor{competitorType, uuid}
}

func (c *Competitor) Equals(c2 *Competitor) bool {
	return c.Type == c2.Type && c.Uuid == c2.Uuid
}

func (c *Competitor) AsValue() Competitor {
	return Competitor{
		c.Type,
		c.Uuid,
	}
}
