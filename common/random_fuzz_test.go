// Copyright © 2026 SafeMPC
//
// This file is part of SafeMPC. The full SafeMPC copyright notice, including
// terms governing use, modification, and redistribution, is contained in the
// file LICENSE at the root of the source code distribution tree.

package common

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"testing"
)

// FuzzGetRandomBytes fuzz tests the GetRandomBytes function
func FuzzGetRandomBytes(f *testing.F) {
	// Add seed corpus
	f.Add(0)
	f.Add(1)
	f.Add(32)
	f.Add(64)
	f.Add(256)
	f.Add(1024)
	f.Add(10000)

	f.Fuzz(func(t *testing.T, length int) {
		// Use a deterministic reader for fuzzing
		reader := bytes.NewReader(make([]byte, length*2))
		result, err := GetRandomBytes(reader, length)
		if length <= 0 {
			if err == nil {
				t.Errorf("GetRandomBytes should return error for length %d", length)
			}
			return
		}
		if err != nil {
			// This is expected if reader doesn't have enough data
			return
		}
		if result == nil {
			t.Errorf("GetRandomBytes returned nil for valid length")
		}
		if len(result) != length {
			t.Errorf("GetRandomBytes returned wrong length: got %d, want %d", len(result), length)
		}
	})
}

// FuzzGetRandomPositiveInt fuzz tests the GetRandomPositiveInt function
func FuzzGetRandomPositiveInt(f *testing.F) {
	// Add seed corpus
	f.Add(int64(1))
	f.Add(int64(100))
	f.Add(int64(1000))
	f.Add(int64(1 << 31))
	f.Add(int64(1<<62))

	f.Fuzz(func(t *testing.T, boundVal int64) {
		if boundVal <= 0 {
			return // Skip invalid bounds
		}
		bound := big.NewInt(boundVal)
		reader := rand.Reader
		result := GetRandomPositiveInt(reader, bound)
		if result == nil {
			t.Errorf("GetRandomPositiveInt returned nil for valid bound")
		}
		if result.Cmp(bound) >= 0 {
			t.Errorf("GetRandomPositiveInt returned value >= bound: %s >= %s", result.String(), bound.String())
		}
		if result.Sign() < 0 {
			t.Errorf("GetRandomPositiveInt returned negative value")
		}
	})
}

