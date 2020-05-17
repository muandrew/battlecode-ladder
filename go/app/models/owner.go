package models

const (
	//CompetitorTypeUser its a user.
	CompetitorTypeUser = CompetitorType("user")
	//CompetitorTypeTeam its a team.
	CompetitorTypeTeam = CompetitorType("team")
)

//CompetitorType What type of entity is the competitor
type CompetitorType string

func (ct CompetitorType) String() string {
	return string(ct)
}

//Competitor Represents a player
type Competitor struct {
	Type CompetitorType
	UUID string
}

//NewCompetitor Creates a new competitor
func NewCompetitor(competitorType CompetitorType, uuid string) *Competitor {
	return &Competitor{competitorType, uuid}
}

//Equals equality check
func (c *Competitor) Equals(c2 *Competitor) bool {
	return c.Type == c2.Type && c.UUID == c2.UUID
}

//AsValue creates a new instance
func (c *Competitor) AsValue() Competitor {
	return Competitor{
		c.Type,
		c.UUID,
	}
}
