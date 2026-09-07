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

package umas

import (
	"math"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/umas/readwrite/model"
)

func TestSession_TransactionIdentifiersStartAtOne(t *testing.T) {
	session := newSession(defaultMaxFrameSize)
	assert.Equal(t, uint16(1), session.nextTransactionIdentifier())
	assert.Equal(t, uint16(2), session.nextTransactionIdentifier())
	assert.Equal(t, uint16(3), session.nextTransactionIdentifier())
}

// The field on the wire is 16 bits wide. Zero is left out so an identifier never collides with an
// uninitialized one, and the counter has to come back around rather than overflow - plc4j's version
// keeps counting past 0xFFFF and only masks the value it returns, which makes its function-key
// tracker (keyed on the unmasked identifier) miss every response after the first wrap.
func TestSession_TransactionIdentifiersWrapWithoutHittingZero(t *testing.T) {
	session := newSession(defaultMaxFrameSize)
	session.transactionIdentifier.Store(math.MaxUint16 - 1)

	assert.Equal(t, uint16(math.MaxUint16), session.nextTransactionIdentifier())
	assert.Equal(t, uint16(1), session.nextTransactionIdentifier(), "the counter has to come back around to one")
	assert.Equal(t, uint16(2), session.nextTransactionIdentifier())
}

func TestSession_TransactionIdentifiersAreUniqueUnderConcurrency(t *testing.T) {
	session := newSession(defaultMaxFrameSize)
	const workers, perWorker = 8, 200

	var mutex sync.Mutex
	seen := map[uint16]bool{}
	var waitGroup sync.WaitGroup
	for range workers {
		waitGroup.Go(func() {
			for range perWorker {
				identifier := session.nextTransactionIdentifier()
				mutex.Lock()
				assert.False(t, seen[identifier], "the identifier %d was handed out twice", identifier)
				seen[identifier] = true
				mutex.Unlock()
			}
		})
	}
	waitGroup.Wait()
	assert.Len(t, seen, workers*perWorker)
	assert.NotContains(t, seen, uint16(0))
}

// The PLC reports its own frame size in the InitComms response, which replaces the configured one.
func TestSession_CommsParametersComeFromThePlc(t *testing.T) {
	session := newSession(defaultMaxFrameSize)
	session.setCommsParameters(260, 0x0304)
	assert.Equal(t, uint16(260), session.getMaxFrameSize())

	// A frame size the echo request couldn't be built from is ignored rather than believed.
	session.setCommsParameters(1, 0x0304)
	assert.Equal(t, uint16(260), session.getMaxFrameSize())
}

func TestSession_ProjectIdentity(t *testing.T) {
	session := newSession(defaultMaxFrameSize)
	assert.Equal(t, uint32(0), session.getHardwareId())
	assert.Equal(t, uint32(0), session.getProjectCrc())
	session.setProjectIdentity(0xDEADBEEF, 0xCAFEBABE)
	assert.Equal(t, uint32(0xDEADBEEF), session.getHardwareId())
	assert.Equal(t, uint32(0xCAFEBABE), session.getProjectCrc())
}

// plc4j declares a pairing key field and never assigns one, so it is zero on every PDU it sends;
// keeping that identical here is what makes the wire bytes match.
func TestSession_ThePairingKeyStaysZero(t *testing.T) {
	assert.Equal(t, uint8(0), newSession(defaultMaxFrameSize).getPairingKey())
}

func TestSession_SymbolsAreLookedUpCaseInsensitively(t *testing.T) {
	session := newSession(defaultMaxFrameSize)
	session.setSymbols([]readWriteModel.UmasUnlocatedVariableReference{
		readWriteModel.NewUmasUnlocatedVariableReference(typeIdReal, 2, 0x100, 0, 0, "G_R32"),
	})

	for _, spelling := range []string{"G_R32", "g_r32", "g_R32"} {
		symbol, ok := session.lookupSymbol(spelling)
		require.True(t, ok, "%s should have been found", spelling)
		assert.Equal(t, "G_R32", symbol.GetValue(), "the original spelling is what the dictionary said")
	}

	_, ok := session.lookupSymbol("g_other")
	assert.False(t, ok)
}

func TestSession_SymbolsAreHandedOutInAStableOrder(t *testing.T) {
	session := newSession(defaultMaxFrameSize)
	session.setSymbols([]readWriteModel.UmasUnlocatedVariableReference{
		readWriteModel.NewUmasUnlocatedVariableReference(typeIdReal, 2, 0x100, 0, 0, "zulu"),
		readWriteModel.NewUmasUnlocatedVariableReference(typeIdReal, 2, 0x104, 0, 0, "alpha"),
		readWriteModel.NewUmasUnlocatedVariableReference(typeIdReal, 2, 0x108, 0, 0, "mike"),
	})

	// Go randomizes map iteration, so a browse would answer in a different order every time.
	for range 5 {
		var names []string
		for _, symbol := range session.symbols() {
			names = append(names, symbol.GetValue())
		}
		assert.Equal(t, []string{"alpha", "mike", "zulu"}, names)
	}
}

// Inside one memory block the symbols are laid out back to back, so the distance to the next
// symbol's offset is the size of this one. The last symbol of a block has no successor and therefore
// no derived size. Ported from plc4j's computeSymbolSizes.
func TestComputeSymbolSizes(t *testing.T) {
	session := newSession(defaultMaxFrameSize)
	session.setSymbols([]readWriteModel.UmasUnlocatedVariableReference{
		readWriteModel.NewUmasUnlocatedVariableReference(typeIdString, 2, 0x100, 0, 0, "first"),
		readWriteModel.NewUmasUnlocatedVariableReference(typeIdString, 2, 0x114, 0, 0, "second"),
		readWriteModel.NewUmasUnlocatedVariableReference(typeIdString, 2, 0x11E, 0, 0, "third"),
		// A different block starts its own run of offsets.
		readWriteModel.NewUmasUnlocatedVariableReference(typeIdString, 3, 0x000, 0, 0, "other"),
		readWriteModel.NewUmasUnlocatedVariableReference(typeIdString, 3, 0x008, 0, 0, "otherToo"),
	})

	assert.Equal(t, map[string]uint16{
		"first":  0x14,
		"second": 0x0A,
		"other":  0x08,
		// "third" and "otherToo" are the last of their blocks, so their size is unknown.
	}, session.symbolSizes)
}

func TestSession_StringBufferSize(t *testing.T) {
	t.Run("the datatype dictionary wins", func(t *testing.T) {
		session := newSession(defaultMaxFrameSize)
		session.setDataTypes([]readWriteModel.UmasDatatypeReference{
			readWriteModel.NewUmasDatatypeReference(20, 0, 0, uint8(typeIdString), "MY_STRING"),
		})
		symbol := readWriteModel.NewUmasUnlocatedVariableReference(customTypeIdBase, 2, 0x100, 0, 0, "g_string")
		session.setSymbols([]readWriteModel.UmasUnlocatedVariableReference{symbol})
		assert.Equal(t, uint16(20), session.stringBufferSize(symbol))
	})

	t.Run("the symbol layout is the fallback", func(t *testing.T) {
		session := newSession(defaultMaxFrameSize)
		symbol := readWriteModel.NewUmasUnlocatedVariableReference(typeIdString, 2, 0x100, 0, 0, "g_string")
		session.setSymbols([]readWriteModel.UmasUnlocatedVariableReference{
			symbol,
			readWriteModel.NewUmasUnlocatedVariableReference(typeIdString, 2, 0x120, 0, 0, "g_next"),
		})
		assert.Equal(t, uint16(0x20), session.stringBufferSize(symbol))
	})

	t.Run("and otherwise the default", func(t *testing.T) {
		session := newSession(defaultMaxFrameSize)
		symbol := readWriteModel.NewUmasUnlocatedVariableReference(typeIdString, 2, 0x100, 0, 0, "g_string")
		session.setSymbols([]readWriteModel.UmasUnlocatedVariableReference{symbol})
		assert.Equal(t, defaultStringBufferSize, session.stringBufferSize(symbol))
	})
}

// The n-th entry of the datatype dictionary describes the type with id customTypeIdBase + n, which
// is what ties a symbol's type id back to the dictionary.
func TestSession_CustomTypeIdsCountFromTheThreshold(t *testing.T) {
	session := newSession(defaultMaxFrameSize)
	session.setDataTypes([]readWriteModel.UmasDatatypeReference{
		readWriteModel.NewUmasDatatypeReference(10, 0, 0, 0, "first"),
		readWriteModel.NewUmasDatatypeReference(20, 0, 0, 0, "second"),
	})
	assert.Equal(t, map[uint16]uint16{
		customTypeIdBase:     10,
		customTypeIdBase + 1: 20,
	}, session.dataTypeSizes)
	assert.Equal(t, uint16(0x1A), customTypeIdBase, "plc4j's CUSTOM_TYPE_THRESHOLD")
}

func TestSession_CustomTypes(t *testing.T) {
	session := newSession(defaultMaxFrameSize)

	t.Run("a primitive is not a custom type", func(t *testing.T) {
		_, known := session.customType(typeIdDint)
		assert.False(t, known)
	})

	t.Run("a struct type carries its members", func(t *testing.T) {
		session.setStructType(customTypeIdBase, "MY_STRUCT", []readWriteModel.UmasUDTDefinition{
			readWriteModel.NewUmasUDTDefinition(typeIdDint, 0, 0, 0, "meta"),
		})
		customType, known := session.customType(customTypeIdBase)
		require.True(t, known)
		assert.Equal(t, "MY_STRUCT", customType.name)
		assert.False(t, customType.isArray)
		require.Len(t, customType.fields, 1)
		assert.Equal(t, "meta", customType.fields[0].GetValue())
	})

	t.Run("an array type carries its element type and bounds", func(t *testing.T) {
		session.setArrayType(customTypeIdBase+1, "MY_ARRAY", typeIdDint,
			[]readWriteModel.UmasArrayDimension{readWriteModel.NewUmasArrayDimension(0, 9)})
		customType, known := session.customType(customTypeIdBase + 1)
		require.True(t, known)
		assert.Equal(t, "MY_ARRAY", customType.name)
		assert.True(t, customType.isArray)
		assert.Equal(t, typeIdDint, customType.elementTypeId)
		require.Len(t, customType.dimensions, 1)
		assert.Equal(t, uint32(9), customType.dimensions[0].GetUpperBound())
		assert.Empty(t, customType.fields, "an array type has no members")
	})
}

func TestSession_String(t *testing.T) {
	session := newSession(defaultMaxFrameSize)
	session.setIdentity("PLC", 1, 2)
	assert.Contains(t, session.String(), "hostname: PLC")
}
