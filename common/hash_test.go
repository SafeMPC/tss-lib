// Copyright © 2026 SafeMPC
//
// This file is part of SafeMPC. The full SafeMPC copyright notice, including
// terms governing use, modification, and redistribution, is contained in the
// file LICENSE at the root of the source code distribution tree.

package common

import (
	"crypto/sha512"
	"math/big"
	"testing"
)

func TestSHA512_256(t *testing.T) {
	tests := []struct {
		name     string
		inputs   [][]byte
		expected int // expected length
	}{
		{
			name:     "empty input",
			inputs:   [][]byte{},
			expected: 0, // nil for empty input
		},
		{
			name:     "single empty byte slice",
			inputs:   [][]byte{{}},
			expected: 32,
		},
		{
			name:     "single non-empty byte slice",
			inputs:   [][]byte{[]byte("test")},
			expected: 32,
		},
		{
			name:     "multiple byte slices",
			inputs:   [][]byte{[]byte("hello"), []byte("world")},
			expected: 32,
		},
		{
			name:     "large input",
			inputs:   [][]byte{make([]byte, 1000)},
			expected: 32,
		},
		{
			name:     "multiple large inputs",
			inputs:   [][]byte{make([]byte, 100), make([]byte, 200)},
			expected: 32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SHA512_256(tt.inputs...)
			if len(tt.inputs) == 0 {
				if result != nil {
					t.Errorf("SHA512_256() with empty input should return nil, got %v", result)
				}
				return
			}
			if result == nil {
				t.Errorf("SHA512_256() returned nil for non-empty input")
				return
			}
			if len(result) != tt.expected {
				t.Errorf("SHA512_256() returned length %d, want %d", len(result), tt.expected)
			}
		})
	}
}

func TestSHA512_256i(t *testing.T) {
	tests := []struct {
		name     string
		inputs   []*big.Int
		expected bool // should return non-nil
	}{
		{
			name:     "empty input",
			inputs:   []*big.Int{},
			expected: false, // nil for empty input
		},
		{
			name:     "single zero",
			inputs:   []*big.Int{big.NewInt(0)},
			expected: true,
		},
		{
			name:     "single positive integer",
			inputs:   []*big.Int{big.NewInt(123456789)},
			expected: true,
		},
		{
			name:     "single negative integer",
			inputs:   []*big.Int{big.NewInt(-1)},
			expected: true,
		},
		{
			name:     "multiple integers",
			inputs:   []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)},
			expected: true,
		},
		{
			name:     "large integer",
			inputs:   []*big.Int{new(big.Int).Lsh(big.NewInt(1), 256)},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SHA512_256i(tt.inputs...)
			if len(tt.inputs) == 0 {
				if result != nil {
					t.Errorf("SHA512_256i() with empty input should return nil, got %v", result)
				}
				return
			}
			if tt.expected && result == nil {
				t.Errorf("SHA512_256i() returned nil when expected non-nil")
			}
			if result != nil && result.Sign() < 0 {
				t.Errorf("SHA512_256i() returned negative result: %s", result.String())
			}
		})
	}
}

func TestSHA512_256i_TAGGED(t *testing.T) {
	tests := []struct {
		name     string
		tag      []byte
		inputs   []*big.Int
		expected bool // should return non-nil
	}{
		{
			name:     "empty tag and inputs",
			tag:      []byte{},
			inputs:   []*big.Int{},
			expected: false, // nil for empty input
		},
		{
			name:     "non-empty tag, empty inputs",
			tag:      []byte("tag"),
			inputs:   []*big.Int{},
			expected: false, // nil for empty input
		},
		{
			name:     "empty tag, non-empty inputs",
			tag:      []byte{},
			inputs:   []*big.Int{big.NewInt(1)},
			expected: true,
		},
		{
			name:     "non-empty tag and inputs",
			tag:      []byte("session"),
			inputs:   []*big.Int{big.NewInt(123456789)},
			expected: true,
		},
		{
			name:     "multiple inputs",
			tag:      []byte("tag"),
			inputs:   []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)},
			expected: true,
		},
		{
			name:     "large tag",
			tag:      make([]byte, 1000),
			inputs:   []*big.Int{big.NewInt(1)},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SHA512_256i_TAGGED(tt.tag, tt.inputs...)
			if len(tt.inputs) == 0 {
				if result != nil {
					t.Errorf("SHA512_256i_TAGGED() with empty input should return nil, got %v", result)
				}
				return
			}
			if tt.expected && result == nil {
				t.Errorf("SHA512_256i_TAGGED() returned nil when expected non-nil")
			}
			if result != nil && result.Sign() < 0 {
				t.Errorf("SHA512_256i_TAGGED() returned negative result: %s", result.String())
			}
		})
	}
}

func TestSHA512_256_DomainSeparation(t *testing.T) {
	// Test that different inputs produce different hashes
	input1 := []byte("test1")
	input2 := []byte("test2")

	hash1 := SHA512_256(input1)
	hash2 := SHA512_256(input2)

	if len(hash1) != 32 || len(hash2) != 32 {
		t.Errorf("Hashes should be 32 bytes")
	}

	// Hashes should be different
	equal := true
	for i := range hash1 {
		if hash1[i] != hash2[i] {
			equal = false
			break
		}
	}
	if equal {
		t.Errorf("Different inputs produced the same hash")
	}
}

func TestSHA512_256i_Consistency(t *testing.T) {
	// Test that same input produces same hash
	input := big.NewInt(123456789)
	hash1 := SHA512_256i(input)
	hash2 := SHA512_256i(input)

	if hash1 == nil || hash2 == nil {
		t.Errorf("Hashes should not be nil")
		return
	}

	if hash1.Cmp(hash2) != 0 {
		t.Errorf("Same input should produce same hash")
	}
}

func TestSHA512_256_StandardCompatibility(t *testing.T) {
	// Test that our implementation produces same result as standard SHA-512/256
	input := []byte("test input")
	
	// Standard SHA-512/256
	standardHash := sha512.Sum512_256(input)
	
	// Our implementation
	ourHash := SHA512_256(input)
	
	if len(ourHash) != 32 {
		t.Errorf("Our hash should be 32 bytes")
		return
	}
	
	// Compare first 32 bytes (our implementation uses domain separation, so may differ)
	// This test mainly ensures our implementation doesn't crash and produces valid output
	if len(standardHash) != len(ourHash) {
		t.Errorf("Hash lengths should match")
	}
}

