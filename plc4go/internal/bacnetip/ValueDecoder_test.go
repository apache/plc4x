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

package bacnetip

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

// taggedStub satisfies the `GetValue() T` shape that taggedEnumToPlcValue
// looks for via reflection. Used to exercise the kinds that don't have a
// real generated counterpart in this codebase (bool, signed, string).
type uintStub struct{ v uint8 }

func (s uintStub) GetValue() uint8 { return s.v }

type signedStub struct{ v int16 }

func (s signedStub) GetValue() int16 { return s.v }

type boolStub struct{ v bool }

func (b boolStub) GetValue() bool { return b.v }

type stringStub struct{ v string }

func (s stringStub) GetValue() string { return s.v }

type noGetValueStub struct{}

func TestTaggedEnumToPlcValue_BACnetBinaryPVTagged(t *testing.T) {
	// The canonical case we hit during integration testing. INACTIVE→0,
	// ACTIVE→1 — both must surface as PlcUDINT, matching how a plain
	// Enumerated application tag is decoded.
	cases := []struct {
		name string
		pv   readWriteModel.BACnetBinaryPV
		want uint32
	}{
		{"INACTIVE", readWriteModel.BACnetBinaryPV_INACTIVE, 0},
		{"ACTIVE", readWriteModel.BACnetBinaryPV_ACTIVE, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Header values don't affect taggedEnumToPlcValue — only GetValue()
			// is consulted via reflection.
			tagged := readWriteModel.NewBACnetBinaryPVTagged(
				readWriteModel.NewBACnetTagHeader(9, 0, 1, nil, nil, nil, nil),
				tc.pv,
			)
			got, ok := taggedEnumToPlcValue(tagged)
			require.True(t, ok, "BACnetBinaryPVTagged should be recognized")
			require.NotNil(t, got)
			assert.Equal(t, tc.want, got.GetUint32(), "PlcUDINT value mismatch")
		})
	}
}

func TestTaggedEnumToPlcValue_NilInput(t *testing.T) {
	got, ok := taggedEnumToPlcValue(nil)
	assert.False(t, ok, "nil input should return ok=false")
	assert.Nil(t, got)
}

func TestTaggedEnumToPlcValue_NoGetValueMethod(t *testing.T) {
	// A struct without GetValue() must not match — it would otherwise
	// silently mis-decode random types as PlcUDINT/PlcSTRING.
	got, ok := taggedEnumToPlcValue(noGetValueStub{})
	assert.False(t, ok)
	assert.Nil(t, got)
}

func TestTaggedEnumToPlcValue_AllKinds(t *testing.T) {
	// Cover each reflection-Kind branch in the switch so future refactors
	// can't accidentally drop one. Uses stub types because the generated
	// model only has uint-returning *Tagged types in our tag-set.
	cases := []struct {
		name string
		in   any
		kind string
		eq   func(t *testing.T, v apiValues.PlcValue)
	}{
		{"uint", uintStub{v: 42}, "udint",
			func(t *testing.T, v apiValues.PlcValue) { assert.Equal(t, uint32(42), v.GetUint32()) }},
		{"int", signedStub{v: -7}, "lint",
			func(t *testing.T, v apiValues.PlcValue) { assert.Equal(t, int64(-7), v.GetInt64()) }},
		{"bool", boolStub{v: true}, "bool",
			func(t *testing.T, v apiValues.PlcValue) { assert.True(t, v.GetBool()) }},
		{"string", stringStub{v: "hello"}, "string",
			func(t *testing.T, v apiValues.PlcValue) { assert.Equal(t, "hello", v.GetString()) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := taggedEnumToPlcValue(tc.in)
			require.True(t, ok, "kind %s should be recognized", tc.kind)
			require.NotNil(t, got)
			tc.eq(t, got)
		})
	}
}

func TestTaggedEnumToPlcValue_RejectsUnsupportedKind(t *testing.T) {
	// A GetValue that returns e.g. a slice or struct isn't a primitive enum
	// — taggedEnumToPlcValue should report ok=false so the caller can fall
	// through to the stringify fallback rather than silently returning nil.
	type sliceStub struct{}
	type sliceStubT = sliceStub
	// Add the method via a wrapper since Go method sets don't allow inline.
	got, ok := taggedEnumToPlcValue(complexReturnStub{})
	assert.False(t, ok, "non-primitive return kind should reject")
	assert.Nil(t, got)
}

type complexReturnStub struct{}

func (complexReturnStub) GetValue() []byte { return []byte{1, 2, 3} }
