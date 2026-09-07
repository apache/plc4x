/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package fields

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── convBool ───────────────────────────────────────────────────────────────

func TestConvBool(t *testing.T) {
	cases := []struct {
		name string
		in   any
		want bool
	}{
		{"nil", nil, false},
		{"true", true, true},
		{"false", false, false},
		{"int non-zero", int(5), true},
		{"int zero", int(0), false},
		{"int8 negative", int8(-1), true},
		{"uint non-zero", uint(1), true},
		{"uintptr non-zero", uintptr(7), true},
		{"float non-zero", 1.5, true},
		{"float zero", 0.0, false},
		{"complex non-zero", complex(1, 0), true},
		{"complex zero", complex(0, 0), false},
		{"string true", "true", true},
		{"string false", "false", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := convBool(tc.in)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestConvBool_UnsupportedType(t *testing.T) {
	_, err := convBool([]int{1, 2})
	assert.Error(t, err)
}

func TestConvBool_BadString(t *testing.T) {
	_, err := convBool("not-a-bool")
	assert.Error(t, err)
}

// ── convInt64 / signed narrow ──────────────────────────────────────────────

func TestConvInt64(t *testing.T) {
	cases := []struct {
		name string
		in   any
		want int64
	}{
		{"nil", nil, 0},
		{"true", true, 1},
		{"false", false, 0},
		{"int", int(-42), -42},
		{"int8", int8(-1), -1},
		{"int64 max", int64(1<<62 + 1), 1<<62 + 1},
		{"uint", uint(7), 7},
		{"uintptr", uintptr(8), 8},
		{"float trunc", 3.9, 3},
		{"complex real", complex(2.5, 99), 2},
		{"string", "-123", -123},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := convInt64(tc.in)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestConvInt64_UnsupportedType(t *testing.T) {
	_, err := convInt64(struct{}{})
	assert.Error(t, err)
}

func TestSignedNarrowWrappers(t *testing.T) {
	// Use values that trivially fit and verify the narrowing happens (cast).
	i, err := convInt(int64(42))
	require.NoError(t, err)
	assert.Equal(t, 42, i)

	i8, err := convInt8(int(127))
	require.NoError(t, err)
	assert.Equal(t, int8(127), i8)

	i16, err := convInt16(int(-32768))
	require.NoError(t, err)
	assert.Equal(t, int16(-32768), i16)

	i32, err := convInt32(int64(-1))
	require.NoError(t, err)
	assert.Equal(t, int32(-1), i32)
}

func TestSignedNarrowWrappers_Truncate(t *testing.T) {
	// int8 holds values in [-128, 127]; 300 wraps via two's complement narrowing.
	got, err := convInt8(int64(300))
	require.NoError(t, err)
	assert.Equal(t, int8(300-256), got)
}

// ── convUint64 / unsigned narrow ───────────────────────────────────────────

func TestConvUint64(t *testing.T) {
	cases := []struct {
		name string
		in   any
		want uint64
	}{
		{"nil", nil, 0},
		{"true", true, 1},
		{"false", false, 0},
		{"int positive", int(42), 42},
		{"uint8", uint8(255), 255},
		{"uintptr", uintptr(9), 9},
		{"float trunc", 3.9, 3},
		{"complex real", complex(7.0, 1.0), 7},
		{"string", "1234", 1234},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := convUint64(tc.in)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestConvUint64_NegativeWraps(t *testing.T) {
	// int(-1) wraps to uint64(max) via two's complement, matching go-conv behavior.
	got, err := convUint64(int(-1))
	require.NoError(t, err)
	assert.Equal(t, uint64(0xFFFFFFFFFFFFFFFF), got)
}

func TestUnsignedNarrowWrappers(t *testing.T) {
	u, err := convUint(uint64(7))
	require.NoError(t, err)
	assert.Equal(t, uint(7), u)

	u8, err := convUint8(uint64(255))
	require.NoError(t, err)
	assert.Equal(t, uint8(255), u8)

	u16, err := convUint16(uint64(0xFFFF))
	require.NoError(t, err)
	assert.Equal(t, uint16(0xFFFF), u16)

	u32, err := convUint32(uint64(0xDEADBEEF))
	require.NoError(t, err)
	assert.Equal(t, uint32(0xDEADBEEF), u32)

	up, err := convUintptr(uint64(42))
	require.NoError(t, err)
	assert.Equal(t, uintptr(42), up)
}

// ── convFloat64 / float narrow ─────────────────────────────────────────────

func TestConvFloat64(t *testing.T) {
	cases := []struct {
		name string
		in   any
		want float64
	}{
		{"nil", nil, 0},
		{"true", true, 1},
		{"int", int(-3), -3},
		{"uint", uint(7), 7},
		{"float", 1.5, 1.5},
		{"complex real", complex(2.25, 99), 2.25},
		{"string", "3.14", 3.14},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := convFloat64(tc.in)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestConvFloat32(t *testing.T) {
	got, err := convFloat32(float64(1.5))
	require.NoError(t, err)
	assert.Equal(t, float32(1.5), got)
}

// ── convComplex128 / convComplex64 ─────────────────────────────────────────

func TestConvComplex128(t *testing.T) {
	cases := []struct {
		name string
		in   any
		want complex128
	}{
		{"nil", nil, 0},
		{"true", true, 1},
		{"int", int(3), complex(3, 0)},
		{"uint", uint(4), complex(4, 0)},
		{"float", 2.5, complex(2.5, 0)},
		{"complex passthrough", complex(1, 2), complex(1, 2)},
		{"string", "(1+2i)", complex(1, 2)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := convComplex128(tc.in)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestConvComplex64(t *testing.T) {
	got, err := convComplex64(complex(1.5, 2.0))
	require.NoError(t, err)
	assert.Equal(t, complex64(complex(1.5, 2.0)), got)
}
