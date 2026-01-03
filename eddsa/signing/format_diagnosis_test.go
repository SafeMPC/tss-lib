package signing

import (
	"crypto/ed25519"
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEd25519ByteOrder tests the byte order of standard Ed25519
func TestEd25519ByteOrder(t *testing.T) {
	// Generate standard Ed25519 key pair
	stdPubKey, stdPrivKey, err := ed25519.GenerateKey(rand.Reader)
	assert.NoError(t, err)

	message := []byte("Test")
	stdSig := ed25519.Sign(stdPrivKey, message)

	// Verify
	valid := ed25519.Verify(stdPubKey, message, stdSig)
	assert.True(t, valid)

	t.Logf("\n📊 Standard Ed25519 Format:")
	t.Logf("Public key (32 bytes): %x", stdPubKey)
	t.Logf("Signature (64 bytes): %x", stdSig)
	t.Logf("Signature R (first 32 bytes): %x", stdSig[:32])
	t.Logf("Signature S (last 32 bytes): %x", stdSig[32:])

	// Check byte order: create a known value
	testValue := big.NewInt(0x01020304)
	testBytes := testValue.Bytes()
	t.Logf("\n📊 Byte Order Test:")
	t.Logf("big.Int value: 0x%x", testValue)
	t.Logf("big.Int.Bytes() (big-endian): %x", testBytes)

	// If big.Int.Bytes() is big-endian, then:
	// 0x01020304 should be represented as [01 02 03 04]
	// If it's little-endian, it should be [04 03 02 01]

	if len(testBytes) > 0 {
		if testBytes[0] == 0x01 {
			t.Logf("✅ big.Int.Bytes() uses BIG-ENDIAN (most significant byte first)")
		} else if testBytes[len(testBytes)-1] == 0x01 {
			t.Logf("✅ big.Int.Bytes() uses LITTLE-ENDIAN (least significant byte first)")
		}
	}
}

// TestTssLibSignatureByteOrder tests the byte order of tss-lib signatures
func TestTssLibSignatureByteOrder(t *testing.T) {
	// Create a test value
	testR := big.NewInt(0x0102030405060708)
	testS := big.NewInt(0x0807060504030201)

	// tss-lib format (little-endian)
	rLE := bigIntToEncodedBytes(testR)
	sLE := bigIntToEncodedBytes(testS)

	t.Logf("\n📊 tss-lib Format (little-endian):")
	t.Logf("R value: 0x%x", testR)
	t.Logf("R encoded (little-endian): %x", rLE)
	t.Logf("S value: 0x%x", testS)
	t.Logf("S encoded (little-endian): %x", sLE)

	// Convert to big-endian
	rBE := littleEndianToBigEndian(rLE)
	sBE := littleEndianToBigEndian(sLE)

	t.Logf("\n📊 Converted Format (big-endian):")
	t.Logf("R converted (big-endian): %x", rBE)
	t.Logf("S converted (big-endian): %x", sBE)

	// Check the format of big.Int.Bytes()
	rBigIntBytes := testR.Bytes()
	sBigIntBytes := testS.Bytes()

	t.Logf("\n📊 big.Int.Bytes() Format:")
	t.Logf("R big.Int.Bytes(): %x", rBigIntBytes)
	t.Logf("S big.Int.Bytes(): %x", sBigIntBytes)

	// Compare
	t.Logf("\n📊 Comparison:")
	t.Logf("rLE vs rBE: %v", *rLE != *rBE)
	t.Logf("rBE vs rBigIntBytes (padded): need to check")
}

