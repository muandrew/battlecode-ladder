package models

//Competition The game engine
type Competition string

const (
	//CompetitionBC17 Battlecode 2017
	CompetitionBC17 = Competition("bc17")
	//CompetitionICPC2011Q ICPC 2011 Queue
	CompetitionICPC2011Q = Competition("icpc2011q")
)

//AsString retruns a string representation of Compeititon
func (c Competition) AsString() string {
	return string(c)
}
