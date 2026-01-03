// Copyright © 2026 SafeMPC
//
// This file is part of SafeMPC. The full SafeMPC copyright notice, including
// terms governing use, modification, and redistribution, is contained in the
// file LICENSE at the root of the source code distribution tree.

package commitments_test

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/SafeMPC/tss-lib/crypto/commitments"
)

func TestCreateVerify(t *testing.T) {
	one := big.NewInt(1)
	zero := big.NewInt(0)

	commitment := NewHashCommitment(rand.Reader, zero, one)
	pass := commitment.Verify()

	assert.True(t, pass, "must pass")
}

func TestDeCommit(t *testing.T) {
	one := big.NewInt(1)
	zero := big.NewInt(0)

	commitment := NewHashCommitment(rand.Reader, zero, one)
	pass, secrets := commitment.DeCommit()

	assert.True(t, pass, "must pass")

	assert.NotZero(t, len(secrets), "len(secrets) must be non-zero")
}
