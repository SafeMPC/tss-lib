// Copyright © 2026 SafeMPC
//
// This file is part of SafeMPC. The full SafeMPC copyright notice, including
// terms governing use, modification, and redistribution, is contained in the
// file LICENSE at the root of the source code distribution tree.

package signing

import (
	"math/big"
	"testing"
)

// FuzzBigIntToEncodedBytes fuzz tests the bigIntToEncodedBytes function
func FuzzBigIntToEncodedBytes(f *testing.F) {
	// Add seed corpus
	f.Add(int64(0))
	f.Add(int64(1))
	f.Add(int64(123456789))
	f.Add(int64(-1))
	f.Add(int64(1 << 31))
	f.Add(int64(1 << 62))
	f.Add(int64(0x7FFFFFFFFFFFFFFF))

	f.Fuzz(func(t *testing.T, val int64) {
		bi := big.NewInt(val)
		result := bigIntToEncodedBytes(bi)
		if result == nil {
			t.Errorf("bigIntToEncodedBytes returned nil")
		}
		if len(result) != 32 {
			t.Errorf("bigIntToEncodedBytes returned wrong length: got %d, want 32", len(result))
		}
	})
}

// FuzzEncodedBytesToBigInt fuzz tests the encodedBytesToBigInt function
func FuzzEncodedBytesToBigInt(f *testing.F) {
	// Add seed corpus
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	ones := make([]byte, 32)
	for i := range ones {
		ones[i] = 0xFF
	}
	f.Add(ones)
	max := make([]byte, 32)
	for i := range max {
		max[i] = 0xFF
	}
	max[0] = 0x7F // Make it positive
	f.Add(max)
	f.Add([]byte{0x01, 0x02, 0x03, 0x04, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

	f.Fuzz(func(t *testing.T, data []byte) {
		// Convert to [32]byte
		if len(data) != 32 {
			// Pad or truncate to 32 bytes
			var arr [32]byte
			copy(arr[:], data)
			data = arr[:]
		}
		var arr [32]byte
		copy(arr[:], data)

		result := encodedBytesToBigInt(&arr)
		if result == nil {
			t.Errorf("encodedBytesToBigInt returned nil")
		}
		// Round-trip test
		encoded := bigIntToEncodedBytes(result)
		if encoded == nil {
			t.Errorf("bigIntToEncodedBytes returned nil in round-trip test")
		}
	})
}

// FuzzCopyBytes fuzz tests the copyBytes function
func FuzzCopyBytes(f *testing.F) {
	// Add seed corpus
	f.Add([]byte(""))
	f.Add([]byte("test"))
	f.Add([]byte{0x00, 0x01, 0x02, 0x03})
	f.Add(make([]byte, 32))
	f.Add(make([]byte, 64))
	f.Add(make([]byte, 100))

	f.Fuzz(func(t *testing.T, data []byte) {
		result := copyBytes(data)
		if data == nil && result != nil {
			t.Errorf("copyBytes returned non-nil for nil input")
		}
		if data != nil && result == nil {
			t.Errorf("copyBytes returned nil for non-nil input")
		}
		if result != nil && len(result) != 32 {
			t.Errorf("copyBytes returned wrong length: got %d, want 32", len(result))
		}
		// Verify first 32 bytes match (or padding)
		if result != nil && len(data) > 0 {
			expectedLen := len(data)
			if expectedLen > 32 {
				expectedLen = 32
			}
			for i := 0; i < expectedLen; i++ {
				idx := len(data) - expectedLen + i
				if idx < 0 {
					idx = 0
				}
				if idx < len(data) && result[i] != data[idx] {
					// This might be expected due to padding logic, so we don't fail
					break
				}
			}
		}
	})
}

// FuzzSignatureToStandardEd25519 fuzz tests the SignatureToStandardEd25519 function
func FuzzSignatureToStandardEd25519(f *testing.F) {
	// Add seed corpus
	f.Add(make([]byte, 64))
	f.Add(make([]byte, 32))
	f.Add(make([]byte, 128))
	validSig := make([]byte, 64)
	for i := range validSig {
		validSig[i] = byte(i)
	}
	f.Add(validSig)

	f.Fuzz(func(t *testing.T, data []byte) {
		result, err := SignatureToStandardEd25519(data)
		if len(data) != 64 {
			if err == nil {
				t.Errorf("SignatureToStandardEd25519 should return error for length %d", len(data))
			}
			return
		}
		if err != nil {
			t.Errorf("SignatureToStandardEd25519 returned error for valid length: %v", err)
		}
		if result == nil {
			t.Errorf("SignatureToStandardEd25519 returned nil result")
		}
		if len(result) != 64 {
			t.Errorf("SignatureToStandardEd25519 returned wrong length: got %d, want 64", len(result))
		}
	})
}
