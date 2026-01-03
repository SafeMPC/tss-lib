// Copyright © 2026 SafeMPC
//
// This file is part of SafeMPC. The full SafeMPC copyright notice, including
// terms governing use, modification, and redistribution, is contained in the
// file LICENSE at the root of the source code distribution tree.

package common

import (
	"math/big"
)

// RejectionSample implements the rejection sampling logic for converting a
// SHA512/256 hash to a value between 0-q
func RejectionSample(q *big.Int, eHash *big.Int) *big.Int { // e' = eHash
	e := eHash.Mod(eHash, q)
	return e
}
