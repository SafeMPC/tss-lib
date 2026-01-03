package signing

import (
	"crypto/ed25519"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/SafeMPC/tss-lib/common"
	"github.com/SafeMPC/tss-lib/eddsa/keygen"
	"github.com/SafeMPC/tss-lib/tss"
)

func TestStandardEd25519Compatibility(t *testing.T) {
	// Use existing test fixtures, set threshold to 0 (allowing single-party signing)
	threshold := 0 // Threshold of 0 means only 1 party is needed for signing
	keys, signPIDs, err := keygen.LoadKeygenTestFixtures(1)
	assert.NoError(t, err, "should load keygen fixtures")

	// Prepare the original message to sign (not pre-hashed)
	originalMessage := []byte("Hello, FROST Ed25519 Standard Compatibility Test!")
	msgBigInt := new(big.Int).SetBytes(originalMessage)

	// Set up single-party context
	pIDs := signPIDs[:1] // Use only the first party
	p2pCtx := tss.NewPeerContext(pIDs)
	parties := make([]*LocalParty, 0, len(pIDs))

	outCh := make(chan tss.Message, len(pIDs))
	endCh := make(chan *common.SignatureData, len(pIDs))

	// Create signing parties
	for i := 0; i < len(pIDs); i++ {
		params := tss.NewParameters(tss.Edwards(), p2pCtx, pIDs[i], len(pIDs), threshold)
		P := NewLocalParty(msgBigInt, params, keys[i], outCh, endCh).(*LocalParty)
		parties = append(parties, P)
	}

	// Start signing protocol
	for _, P := range parties {
		go func(P *LocalParty) {
			if err := P.Start(); err != nil {
				t.Errorf("party %s failed to start: %v", P.PartyID(), err)
			}
		}(P)
	}

	// Process messages until completion
	done := make(chan bool, 1)
	go func() {
		for {
			select {
			case msg := <-outCh:
				// No message forwarding needed in single-party mode
				_ = msg
			case <-endCh:
				done <- true
				return
			}
		}
	}()

	// Wait for signing to complete
	<-done

	// Get signature result
	var signatureData *common.SignatureData
	select {
	case signatureData = <-endCh:
		t.Log("✅ Signature completed successfully")
	default:
		t.Fatal("❌ Signing did not complete")
	}

	// Extract public key from key data
	firstKey := keys[0]
	pubKeyX := firstKey.EDDSAPub.X()
	pubKeyY := firstKey.EDDSAPub.Y()

	// Convert tss-lib public key to standard Ed25519 format
	// Ed25519 public key is 32 bytes of Y coordinate with sign bit in the most significant bit
	pubKeyBytes := make([]byte, 32)
	yBytes := pubKeyY.Bytes()

	// Copy Y coordinate bytes (note byte order conversion)
	for i, b := range yBytes {
		if i < 32 {
			pubKeyBytes[i] = b
		}
	}

	// If the least significant bit of X coordinate is 1, set the sign bit
	if pubKeyX.Bit(0) == 1 {
		pubKeyBytes[31] |= 0x80
	}

	// Verify signature using standard crypto/ed25519.Verify
	valid := ed25519.Verify(ed25519.PublicKey(pubKeyBytes), originalMessage, signatureData.Signature)

	// Assert verification succeeds
	assert.True(t, valid, "tss-lib EdDSA signature should be verifiable with standard Ed25519")

	if valid {
		t.Log("✅ SUCCESS: tss-lib EdDSA signature is compatible with standard Ed25519 verification!")
		t.Logf("Original message: %s", originalMessage)
		t.Logf("Signature length: %d bytes", len(signatureData.Signature))
		t.Logf("Public key (first 8 bytes): %x", pubKeyBytes[:8])
		t.Logf("Message in signature data: %s", signatureData.M)
		t.Logf("Signature R component: %x", signatureData.R)
		t.Logf("Signature S component: %x", signatureData.S)
	} else {
		t.Error("❌ FAILURE: tss-lib EdDSA signature is NOT compatible with standard Ed25519")

		// Debug information
		t.Logf("Signature bytes: %x", signatureData.Signature)
		t.Logf("Public key bytes: %x", pubKeyBytes)
		t.Logf("Message bytes: %x", originalMessage)

		// Try generating signature with standard library for comparison
		stdPubKey, stdPrivKey, _ := ed25519.GenerateKey(nil)
		stdSignature := ed25519.Sign(stdPrivKey, originalMessage)
		stdValid := ed25519.Verify(stdPubKey, originalMessage, stdSignature)
		t.Logf("Standard Ed25519 test (should pass): %v", stdValid)
		t.Logf("Standard signature length: %d bytes", len(stdSignature))
	}
}

func TestEd25519SignatureFormat(t *testing.T) {
	// Test if signature format conforms to Ed25519 standard
	keys, signPIDs, err := keygen.LoadKeygenTestFixturesRandomSet(1, 1)
	assert.NoError(t, err)

	message := []byte("Test Ed25519 signature format")
	msgBigInt := new(big.Int).SetBytes(message)

	pID := signPIDs[0]
	p2pCtx := tss.NewPeerContext([]*tss.PartyID{pID})
	params := tss.NewParameters(tss.Edwards(), p2pCtx, pID, 1, 0)

	out := make(chan tss.Message, 1)
	end := make(chan *common.SignatureData, 1)

	party := NewLocalParty(msgBigInt, params, keys[0], out, end)

	go func() {
		if err := party.Start(); err != nil {
			t.Errorf("party failed to start: %v", err)
		}
	}()

	party.Update(nil)

	var sigData *common.SignatureData
	select {
	case sigData = <-end:
		// Ed25519 signature should be 64 bytes
		assert.Equal(t, 64, len(sigData.Signature), "Ed25519 signature should be 64 bytes")

		// Verify R and S components
		assert.NotEmpty(t, sigData.R, "R component should not be empty")
		assert.NotEmpty(t, sigData.S, "S component should not be empty")

		// Verify message is correctly preserved
		assert.Equal(t, message, sigData.M, "Original message should be preserved in signature data")

		t.Logf("✅ Signature format verified: length=%d, R=%x..., S=%x...",
			len(sigData.Signature), sigData.R[:8], sigData.S[:8])

	default:
		t.Fatal("Signing did not complete")
	}
}
