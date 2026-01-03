// Copyright © 2026 SafeMPC
//
// This file is part of SafeMPC. The full SafeMPC copyright notice, including
// terms governing use, modification, and redistribution, is contained in the
// file LICENSE at the root of the source code distribution tree.

package common

import (
	"math/big"
	"testing"
)

// FuzzSHA512_256 fuzz tests the SHA512_256 function with various byte inputs
func FuzzSHA512_256(f *testing.F) {
	// Add seed corpus
	f.Add([]byte(""))
	f.Add([]byte("test"))
	f.Add([]byte("hello world"))
	f.Add([]byte{0x00, 0x01, 0x02, 0x03})
	f.Add(make([]byte, 100))
	f.Add(make([]byte, 1000))
	f.Add(make([]byte, 10000))

	f.Fuzz(func(t *testing.T, data []byte) {
		// Test with single input
		result := SHA512_256(data)
		if result == nil && len(data) > 0 {
			t.Errorf("SHA512_256 returned nil for non-empty input")
		}
		if len(result) != 32 {
			t.Errorf("SHA512_256 returned wrong length: got %d, want 32", len(result))
		}

		// Test with multiple inputs
		if len(data) > 0 {
			parts := [][]byte{data[:len(data)/2], data[len(data)/2:]}
			result2 := SHA512_256(parts...)
			if result2 == nil {
				t.Errorf("SHA512_256 returned nil for multiple inputs")
			}
			if len(result2) != 32 {
				t.Errorf("SHA512_256 returned wrong length: got %d, want 32", len(result2))
			}
		}
	})
}

// FuzzSHA512_256i fuzz tests the SHA512_256i function with various big.Int inputs
func FuzzSHA512_256i(f *testing.F) {
	// Add seed corpus
	f.Add(int64(0))
	f.Add(int64(1))
	f.Add(int64(123456789))
	f.Add(int64(-1))
	f.Add(int64(1 << 31))
	f.Add(int64(1<<62))

	f.Fuzz(func(t *testing.T, val int64) {
		bi := big.NewInt(val)
		result := SHA512_256i(bi)
		if result == nil && val != 0 {
			t.Errorf("SHA512_256i returned nil for non-zero input")
		}
		if result != nil && result.Sign() < 0 {
			t.Errorf("SHA512_256i returned negative result")
		}
	})
}

// FuzzSHA512_256i_TAGGED fuzz tests the SHA512_256i_TAGGED function
func FuzzSHA512_256i_TAGGED(f *testing.F) {
	// Add seed corpus
	f.Add([]byte("tag"), int64(0))
	f.Add([]byte("tag"), int64(1))
	f.Add([]byte("session"), int64(123456789))
	f.Add([]byte(""), int64(0))
	f.Add(make([]byte, 100), int64(1))

	f.Fuzz(func(t *testing.T, tag []byte, val int64) {
		bi := big.NewInt(val)
		result := SHA512_256i_TAGGED(tag, bi)
		if result == nil && val != 0 {
			t.Errorf("SHA512_256i_TAGGED returned nil for non-zero input")
		}
		if result != nil && result.Sign() < 0 {
			t.Errorf("SHA512_256i_TAGGED returned negative result")
		}
	})
}

