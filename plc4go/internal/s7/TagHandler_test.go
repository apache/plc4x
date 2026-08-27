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

package s7

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
)

func TestTagHandlerParseTag(t *testing.T) {
	handler := NewTagHandler()

	t.Run("S5TIME address", func(t *testing.T) {
		tag, err := handler.ParseTag("%M0:S5TIME")
		require.NoError(t, err)
		s7Tag := tag.(PlcTag)
		assert.Equal(t, readWriteModel.TransportSize_S5TIME, s7Tag.GetDataType())
		assert.Equal(t, readWriteModel.MemoryArea_FLAGS_MARKERS, s7Tag.GetMemoryArea())
	})
	t.Run("S5TIME in data block", func(t *testing.T) {
		tag, err := handler.ParseTag("%DB10:4:S5TIME")
		require.NoError(t, err)
		s7Tag := tag.(PlcTag)
		assert.Equal(t, readWriteModel.TransportSize_S5TIME, s7Tag.GetDataType())
		assert.Equal(t, uint16(10), s7Tag.GetBlockNumber())
	})
	t.Run("fixed length string", func(t *testing.T) {
		tag, err := handler.ParseTag("%DB69:68:STRING(20)")
		require.NoError(t, err)
		stringTag := tag.(PlcStringTag)
		assert.Equal(t, uint16(20), stringTag.stringLength)
	})
	t.Run("var length string short form", func(t *testing.T) {
		tag, err := handler.ParseTag("%DB69:68:STRING")
		require.NoError(t, err)
		stringTag := tag.(PlcStringTag)
		assert.Equal(t, uint16(254), stringTag.stringLength)
		assert.Equal(t, readWriteModel.TransportSize_STRING, stringTag.GetDataType())
	})
	t.Run("var length string long form with array", func(t *testing.T) {
		tag, err := handler.ParseTag("%DB69.DBX68[0..2]:WSTRING")
		require.NoError(t, err)
		stringTag := tag.(PlcStringTag)
		assert.Equal(t, uint16(254), stringTag.stringLength)
		assert.Equal(t, uint16(3), stringTag.GetNumElements())
		assert.Equal(t, readWriteModel.TransportSize_WSTRING, stringTag.GetDataType())
	})
	t.Run("alarm subscription tag", func(t *testing.T) {
		tag, err := handler.ParseTag("ALM")
		require.NoError(t, err)
		alarmTag := tag.(*AlarmTag)
		assert.Equal(t, AlarmTagPush, alarmTag.GetKind())
	})
	t.Run("alarm query tags", func(t *testing.T) {
		tag, err := handler.ParseTag("QUERY:ALARM_S")
		require.NoError(t, err)
		alarmTag := tag.(*AlarmTag)
		assert.Equal(t, AlarmTagQuery, alarmTag.GetKind())
		assert.Equal(t, readWriteModel.QueryType_ALARM_S, alarmTag.GetQueryType())

		tag, err = handler.ParseTag("query:alarm_8")
		require.NoError(t, err)
		alarmTag = tag.(*AlarmTag)
		assert.Equal(t, readWriteModel.QueryType_ALARM_8, alarmTag.GetQueryType())
	})
	t.Run("the current forms parse", func(t *testing.T) {
		for _, address := range []string{"%Q0.0:BOOL", "%M100[0..9]:INT", "%DB1.DBX0.0:BOOL", "%DB1:0.0:BOOL", "%I0:BYTE"} {
			_, err := handler.ParseTag(address)
			assert.NoError(t, err, address)
		}
	})
	t.Run("oversized byte offset is rejected", func(t *testing.T) {
		_, err := handler.ParseTag("%DB1:70000:INT")
		assert.Error(t, err)
	})
	t.Run("bogus string type suffix is rejected", func(t *testing.T) {
		_, err := handler.ParseTag("%DB1:0:STRINGX")
		assert.Error(t, err)
	})
}

// An S7 address names a byte offset, so a selection that starts past the declared base is
// resolved into the address itself: "%DB1.DBW20[4..7]" is the same read as "%DB1.DBW28[0..3]",
// four words further in. The written indices are gone by the time the tag exists, which is why
// GetArrayInfo reports a shape rather than the indices - see Tag.GetArrayInfo.
func TestTagHandlerConsumesTheSelectionOffset(t *testing.T) {
	handler := NewTagHandler()

	shifted, err := handler.ParseTag("%DB1.DBW20[4..7]:INT")
	require.NoError(t, err)
	equivalent, err := handler.ParseTag("%DB1.DBW28[0..3]:INT")
	require.NoError(t, err)
	assert.Equal(t, equivalent, shifted)

	// A declared base is what the offset is measured from, so [4..7;4] starts at the base and
	// shifts nothing.
	fromDeclaredBase, err := handler.ParseTag("%DB1.DBW20[4..7;4]:INT")
	require.NoError(t, err)
	unshifted, err := handler.ParseTag("%DB1.DBW20[0..3]:INT")
	require.NoError(t, err)
	assert.Equal(t, unshifted, fromDeclaredBase)

	// The shift is in elements, so it scales with the type's size: a DINT is four bytes.
	dwords, err := handler.ParseTag("%DB1.DBD20[2..3]:DINT")
	require.NoError(t, err)
	assert.Equal(t, uint16(28), dwords.(PlcTag).GetByteOffset())
	assert.Equal(t, uint16(2), dwords.(PlcTag).GetNumElements())
}

// A bare index selects one element, which is a scalar - not a count of that many.
func TestTagHandlerReadsABareIndexAsOneElement(t *testing.T) {
	handler := NewTagHandler()

	tag, err := handler.ParseTag("%DB1.DBW20[4]:INT")
	require.NoError(t, err)
	assert.Equal(t, uint16(1), tag.(PlcTag).GetNumElements())
	assert.Equal(t, uint16(28), tag.(PlcTag).GetByteOffset())
	assert.Empty(t, tag.GetArrayInfo(), "one element is a scalar")
}

// Addresses written with the count after the type must fail, naming what to write instead.
func TestTagHandlerRejectsTheOldCountSuffix(t *testing.T) {
	handler := NewTagHandler()

	for _, address := range []string{"%M100:INT[10]", "%DB1.DBW20:INT[4]", "%DB69.DBX68:WSTRING[3]", "%DB1:0:STRING(40)[3]"} {
		t.Run(address, func(t *testing.T) {
			_, err := handler.ParseTag(address)
			require.Error(t, err, address)
			assert.Contains(t, err.Error(), "invalid address", address)
		})
	}
}

// An S7 read is one contiguous byte range, so nothing deeper than a single dimension fits, and a
// selection may not span more than the addressable area.
func TestTagHandlerRejectsWhatS7CannotRead(t *testing.T) {
	handler := NewTagHandler()

	_, err := handler.ParseTag("%DB1.DBW20[0..1][2..3]:INT")
	assert.Error(t, err, "S7 reads one dimension")

	_, err = handler.ParseTag("%DB1.DBW0[0..40000]:INT")
	assert.Error(t, err, "40001 INTs span more than the addressable 65536 bytes")
}

// A rendered address must parse back to the same tag. It did not: the tag rendered as
// "0:INT[8]", which named neither the memory area nor the offset it read, and parsed as nothing.
func TestTagHandler_AddressStringRoundTrips(t *testing.T) {
	handler := NewTagHandler()

	for _, address := range []string{
		"%M100[0..9]:INT",
		"%DB1.DB20[0..3]:INT",
		"%DB1.DB20:INT",
		"%Q0.0:BOOL",
		"%DB1.DB0.0:BOOL",
		"%DB1.DB0[0..2]:STRING(20)",
		"%DB69.DB68[0..2]:WSTRING(254)",
	} {
		t.Run(address, func(t *testing.T) {
			tag, err := handler.ParseTag(address)
			require.NoError(t, err)
			assert.Equal(t, address, tag.GetAddressString())

			reparsed, err := handler.ParseTag(tag.GetAddressString())
			require.NoError(t, err, "the rendered address must parse")
			assert.Equal(t, tag, reparsed)
		})
	}

	// The optional transfer size code is not part of the canonical form - it only repeats what
	// the type already says - so an address carrying one renders without it and still re-parses
	// to the same tag.
	withSizeCode, err := handler.ParseTag("%DB69.DBX68[0..2]:WSTRING(254)")
	require.NoError(t, err)
	assert.Equal(t, "%DB69.DB68[0..2]:WSTRING(254)", withSizeCode.GetAddressString())
	reparsed, err := handler.ParseTag(withSizeCode.GetAddressString())
	require.NoError(t, err)
	assert.Equal(t, withSizeCode, reparsed)

	// The data block is part of the address: a string tag used to be built with a hard-coded
	// block number of 0, so this address read DB0 rather than DB69.
	stringTag, err := handler.ParseTag("%DB69.DBX68:STRING(10)")
	require.NoError(t, err)
	assert.Equal(t, uint16(69), stringTag.(PlcTag).GetBlockNumber())
}

// A fixed-length string is read from the data block its address names.
//
// The long-form branch built the tag with a hard-coded block number of zero, so
// "%DB69.DBX68:STRING(10)" read DB0 and reported the result as though it had come from DB69 -
// wrong data, with nothing to suggest anything had gone wrong. Every other branch of ParseTag
// already parsed the block number; only this one, and the equivalent for WSTRING, did not.
//
// plc4j parses it (S7StringFixedLengthTag.of), so the two bindings disagreed about the same
// address.
func TestTagHandlerParsesTheBlockNumberOfAFixedLengthString(t *testing.T) {
	handler := NewTagHandler()

	for _, c := range []struct {
		address     string
		blockNumber uint16
		byteOffset  uint16
	}{
		{"%DB69.DBX68:STRING(10)", 69, 68},
		{"%DB69.DBX68:WSTRING(10)", 69, 68},
		{"%DB1.DBX0:STRING(20)", 1, 0},
		// The short form was never affected; it is here so a regression in either is caught.
		{"%DB69:68:STRING(10)", 69, 68},
	} {
		t.Run(c.address, func(t *testing.T) {
			tag, err := handler.ParseTag(c.address)
			require.NoError(t, err)
			assert.Equal(t, c.blockNumber, tag.(PlcTag).GetBlockNumber(), "the data block the address names")
			assert.Equal(t, c.byteOffset, tag.(PlcTag).GetByteOffset())
		})
	}
}
