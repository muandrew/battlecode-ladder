package models

import "github.com/muandrew/battlecode-legacy-go/utils"

/*
RAM Resource Access and Management

Currently only the creator has write access.
*/
type RAM struct {
	owner      *Competitor
	readPublic bool
	whitelist  map[Competitor]ramAccessTier
}

type ramAccessTier int

const (
	//RAMAccessNone no access
	RAMAccessNone = ramAccessTier(0)
	//RAMAccessR read
	RAMAccessR = ramAccessTier(1)
	//RAMAccessRW read and write allowed
	RAMAccessRW   = ramAccessTier(3)
	errorAccess   = utils.Error("Insufficient permission.")
	errorArgument = utils.Error("Incorrect Argument.")
)

//CreateRAM creates a new instance of RAM
func CreateRAM(owner *Competitor) *RAM {
	return &RAM{owner: owner}
}

//ReadAllowed returns true if the actor should be allowed to read.
func (r *RAM) ReadAllowed(actor *Competitor) bool {
	if r.readPublic {
		return true
	}
	return r.userBasedAccess(actor, RAMAccessR)
}

//WriteAllowed returns true if the actor should be allowed to write.
func (r *RAM) WriteAllowed(actor *Competitor) bool {
	return r.userBasedAccess(actor, RAMAccessRW)
}

//SetPublic sets the resource to public if the actor has permission
func (r *RAM) SetPublic(actor *Competitor, public bool) error {
	if !r.WriteAllowed(actor) {
		return errorAccess
	}
	r.readPublic = public
	return nil
}

//SetAccess using the actor's permissions, sets this resources's access for the competitor.
func (r *RAM) SetAccess(actor *Competitor, competitor *Competitor, access ramAccessTier) error {
	if !r.WriteAllowed(actor) {
		return errorAccess
	}
	if competitor == nil {
		return errorArgument
	}
	if access == RAMAccessNone {
		delete(r.whitelist, competitor.AsValue())
	} else {
		r.whitelist[competitor.AsValue()] = access
	}
	return nil
}

//TransferOwnership transfers ownership from actor to competitor if allowed
func (r *RAM) TransferOwnership(actor *Competitor, competitor *Competitor) error {
	if !actor.Equals(r.owner) {
		return errorAccess
	}
	if competitor == nil {
		return errorArgument
	}
	r.owner = competitor
	return nil
}

func (r *RAM) userBasedAccess(competitor *Competitor, access ramAccessTier) bool {
	if competitor == nil {
		return false
	} else if competitor.Equals(r.owner) {
		return true
	} else if r.whitelist[competitor.AsValue()] >= access {
		return true
	} else {
		return false
	}
}
