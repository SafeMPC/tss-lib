// Copyright © 2026 SafeMPC
//
// This file is part of SafeMPC. The full SafeMPC copyright notice, including
// terms governing use, modification, and redistribution, is contained in the
// file LICENSE at the root of the source code distribution tree.

package tss

type (
	PeerContext struct {
		partyIDs SortedPartyIDs
	}
)

func NewPeerContext(parties SortedPartyIDs) *PeerContext {
	return &PeerContext{partyIDs: parties}
}

func (p2pCtx *PeerContext) IDs() SortedPartyIDs {
	return p2pCtx.partyIDs
}

func (p2pCtx *PeerContext) SetIDs(ids SortedPartyIDs) {
	p2pCtx.partyIDs = ids
}
