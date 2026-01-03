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

func TestModInt(t *testing.T) {
	mod := big.NewInt(17)
	mi := ModInt(mod)

	if mi == nil {
		t.Errorf("ModInt() returned nil")
	}
}

func TestModInt_Add(t *testing.T) {
	mod := big.NewInt(17)
	mi := ModInt(mod)

	tests := []struct {
		name     string
		x        *big.Int
		y        *big.Int
		expected *big.Int
	}{
		{
			name:     "simple addition",
			x:        big.NewInt(5),
			y:        big.NewInt(7),
			expected: big.NewInt(12),
		},
		{
			name:     "addition with wrap",
			x:        big.NewInt(10),
			y:        big.NewInt(10),
			expected: big.NewInt(3), // (10 + 10) mod 17 = 3
		},
		{
			name:     "zero addition",
			x:        big.NewInt(5),
			y:        big.NewInt(0),
			expected: big.NewInt(5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mi.Add(tt.x, tt.y)
			if result.Cmp(tt.expected) != 0 {
				t.Errorf("Add() = %s, want %s", result.String(), tt.expected.String())
			}
		})
	}
}

func TestModInt_Sub(t *testing.T) {
	mod := big.NewInt(17)
	mi := ModInt(mod)

	tests := []struct {
		name     string
		x        *big.Int
		y        *big.Int
		expected *big.Int
	}{
		{
			name:     "simple subtraction",
			x:        big.NewInt(10),
			y:        big.NewInt(5),
			expected: big.NewInt(5),
		},
		{
			name:     "subtraction with wrap",
			x:        big.NewInt(5),
			y:        big.NewInt(10),
			expected: big.NewInt(12), // (5 - 10) mod 17 = 12
		},
		{
			name:     "zero subtraction",
			x:        big.NewInt(5),
			y:        big.NewInt(0),
			expected: big.NewInt(5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mi.Sub(tt.x, tt.y)
			if result.Cmp(tt.expected) != 0 {
				t.Errorf("Sub() = %s, want %s", result.String(), tt.expected.String())
			}
		})
	}
}

func TestModInt_Mul(t *testing.T) {
	mod := big.NewInt(17)
	mi := ModInt(mod)

	tests := []struct {
		name     string
		x        *big.Int
		y        *big.Int
		expected *big.Int
	}{
		{
			name:     "simple multiplication",
			x:        big.NewInt(3),
			y:        big.NewInt(4),
			expected: big.NewInt(12),
		},
		{
			name:     "multiplication with wrap",
			x:        big.NewInt(5),
			y:        big.NewInt(4),
			expected: big.NewInt(3), // (5 * 4) mod 17 = 3
		},
		{
			name:     "zero multiplication",
			x:        big.NewInt(5),
			y:        big.NewInt(0),
			expected: big.NewInt(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mi.Mul(tt.x, tt.y)
			if result.Cmp(tt.expected) != 0 {
				t.Errorf("Mul() = %s, want %s", result.String(), tt.expected.String())
			}
		})
	}
}

func TestModInt_Div(t *testing.T) {
	mod := big.NewInt(17)
	mi := ModInt(mod)

	tests := []struct {
		name     string
		x        *big.Int
		y        *big.Int
		expected *big.Int
	}{
		{
			name:     "simple division",
			x:        big.NewInt(12),
			y:        big.NewInt(3),
			expected: big.NewInt(4),
		},
		{
			name:     "division with wrap",
			x:        big.NewInt(10),
			y:        big.NewInt(5),
			expected: big.NewInt(2),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mi.Div(tt.x, tt.y)
			// Verify: result * y mod mod = x
			verify := mi.Mul(result, tt.y)
			if verify.Cmp(tt.x) != 0 {
				t.Errorf("Div() verification failed: %s * %s mod %s = %s, want %s",
					result.String(), tt.y.String(), mod.String(), verify.String(), tt.x.String())
			}
		})
	}
}

func TestModInt_Exp(t *testing.T) {
	mod := big.NewInt(17)
	mi := ModInt(mod)

	tests := []struct {
		name     string
		x        *big.Int
		y        *big.Int
		expected *big.Int
	}{
		{
			name:     "simple exponentiation",
			x:        big.NewInt(3),
			y:        big.NewInt(2),
			expected: big.NewInt(9),
		},
		{
			name:     "exponentiation with wrap",
			x:        big.NewInt(5),
			y:        big.NewInt(3),
			expected: big.NewInt(6), // 5^3 mod 17 = 125 mod 17 = 6
		},
		{
			name:     "zero exponent",
			x:        big.NewInt(5),
			y:        big.NewInt(0),
			expected: big.NewInt(1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mi.Exp(tt.x, tt.y)
			if result.Cmp(tt.expected) != 0 {
				t.Errorf("Exp() = %s, want %s", result.String(), tt.expected.String())
			}
		})
	}
}

func TestModInt_ModInverse(t *testing.T) {
	mod := big.NewInt(17)
	mi := ModInt(mod)

	tests := []struct {
		name string
		g    *big.Int
	}{
		{
			name: "inverse of 3",
			g:    big.NewInt(3),
		},
		{
			name: "inverse of 5",
			g:    big.NewInt(5),
		},
		{
			name: "inverse of 7",
			g:    big.NewInt(7),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv := mi.ModInverse(tt.g)
			// Verify: g * inv mod mod = 1
			verify := mi.Mul(tt.g, inv)
			one := big.NewInt(1)
			if verify.Cmp(one) != 0 {
				t.Errorf("ModInverse() verification failed: %s * %s mod %s = %s, want 1",
					tt.g.String(), inv.String(), mod.String(), verify.String())
			}
		})
	}
}

func TestIsInInterval(t *testing.T) {
	tests := []struct {
		name     string
		b        *big.Int
		bound    *big.Int
		expected bool
	}{
		{
			name:     "within interval",
			b:        big.NewInt(5),
			bound:    big.NewInt(10),
			expected: true,
		},
		{
			name:     "at lower bound",
			b:        big.NewInt(0),
			bound:    big.NewInt(10),
			expected: true,
		},
		{
			name:     "at upper bound",
			b:        big.NewInt(10),
			bound:    big.NewInt(10),
			expected: false,
		},
		{
			name:     "above bound",
			b:        big.NewInt(15),
			bound:    big.NewInt(10),
			expected: false,
		},
		{
			name:     "negative",
			b:        big.NewInt(-1),
			bound:    big.NewInt(10),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsInInterval(tt.b, tt.bound)
			if result != tt.expected {
				t.Errorf("IsInInterval(%s, %s) = %v, want %v",
					tt.b.String(), tt.bound.String(), result, tt.expected)
			}
		})
	}
}

func TestAppendBigIntToBytesSlice(t *testing.T) {
	tests := []struct {
		name     string
		common   []byte
		appended *big.Int
	}{
		{
			name:     "empty common",
			common:   []byte{},
			appended: big.NewInt(123),
		},
		{
			name:     "non-empty common",
			common:   []byte{1, 2, 3},
			appended: big.NewInt(456),
		},
		{
			name:     "large integer",
			common:   []byte{1, 2, 3},
			appended: new(big.Int).Lsh(big.NewInt(1), 256),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AppendBigIntToBytesSlice(tt.common, tt.appended)
			expectedLen := len(tt.common) + len(tt.appended.Bytes())
			if len(result) != expectedLen {
				t.Errorf("AppendBigIntToBytesSlice() returned length %d, want %d",
					len(result), expectedLen)
			}
			// Verify first part matches
			for i := 0; i < len(tt.common); i++ {
				if result[i] != tt.common[i] {
					t.Errorf("AppendBigIntToBytesSlice() first part mismatch at index %d", i)
				}
			}
		})
	}
}

