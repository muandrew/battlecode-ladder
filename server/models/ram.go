package models

import "github.com/muandrew/battlecode-ladder/utils"

/**
Resource Access and Management

Currently only the creator has write access.
 */
type RAM struct {
	owner      *Competitor
	readPublic bool
	whitelist  map[Competitor]ramAccessTier
}

type ramAccessTier int

const (
	RamAccessNone = ramAccessTier(0)
	RamAccessR    = ramAccessTier(1)
	RamAccessRW   = ramAccessTier(2)
	errorAccess   = utils.Error("Insufficient permission.")
	errorArgument = utils.Error("Incorrect Argument.")
)

func CreateRAM(owner *Competitor) *RAM {
	return &RAM{owner: owner}
}

func (r *RAM) ReadAllowed(competitor *Competitor) bool {
	if r.readPublic {
		return true
	}
	return r.userBasedAccess(competitor, RamAccessR)
}

func (r *RAM) WriteAllowed(competitor *Competitor) bool {
	return r.userBasedAccess(competitor, RamAccessRW)
}

func (r *RAM) SetPublic(actor *Competitor, public bool) error {
	if !r.WriteAllowed(actor) {
		return errorAccess
	}
	r.readPublic = public
	return nil
}

func (r *RAM) SetAccess(actor *Competitor, competitor *Competitor, access ramAccessTier) error {
	if !r.WriteAllowed(actor) {
		return errorAccess
	}
	if competitor == nil {
		return errorArgument
	}
	if access == RamAccessNone {
		delete(r.whitelist, competitor.AsValue())
	} else {
		r.whitelist[competitor.AsValue()] = access
	}
	return nil
}

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
