// Copyright © 2026 SafeMPC
//
// This file is part of SafeMPC. The full SafeMPC copyright notice, including
// terms governing use, modification, and redistribution, is contained in the
// file LICENSE at the root of the source code distribution tree.

package tss

type Round interface {
	Params() *Parameters
	Start() *Error
	Update() (bool, *Error)
	RoundNumber() int
	CanAccept(msg ParsedMessage) bool
	CanProceed() bool
	NextRound() Round
	WaitingFor() []*PartyID
	WrapError(err error, culprits ...*PartyID) *Error
}
