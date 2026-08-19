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

package modbus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
)

func parseTag(t *testing.T, address string) modbusTag {
	t.Helper()
	tag, err := NewTagHandler().ParseTag(address)
	require.NoError(t, err)
	parsed, ok := tag.(modbusTag)
	require.True(t, ok, "expected a modbusTag, got %T", tag)
	return parsed
}

func TestTagHandler_ParseTag(t *testing.T) {
	tests := []struct {
		name     string
		address  string
		tagType  TagType
		wire     uint16
		quantity uint16
		datatype readWriteModel.ModbusDataType
	}{
		{"plc4x coil", "coil:1:BOOL", Coil, 0, 1, readWriteModel.ModbusDataType_BOOL},
		{"numeric coil", "0x00001:BOOL", Coil, 0, 1, readWriteModel.ModbusDataType_BOOL},
		{"plc4x discrete input", "discrete-input:5:BOOL", DiscreteInput, 4, 1, readWriteModel.ModbusDataType_BOOL},
		{"numeric discrete input", "1x00005:BOOL", DiscreteInput, 4, 1, readWriteModel.ModbusDataType_BOOL},
		{"plc4x input register", "input-register:3:INT", InputRegister, 2, 1, readWriteModel.ModbusDataType_INT},
		{"numeric input register", "3x00003:INT", InputRegister, 2, 1, readWriteModel.ModbusDataType_INT},
		{"plc4x holding register", "holding-register:1:REAL", HoldingRegister, 0, 1, readWriteModel.ModbusDataType_REAL},
		{"numeric holding register", "4x00001:REAL", HoldingRegister, 0, 1, readWriteModel.ModbusDataType_REAL},
		// The extended register area is addressed starting at zero, so its logical address is the
		// wire address itself (plc4j ModbusTagExtendedRegister.getLogicalAddress).
		{"plc4x extended register", "extended-register:7:DINT", ExtendedRegister, 7, 1, readWriteModel.ModbusDataType_DINT},
		{"numeric extended register", "6x00007:DINT", ExtendedRegister, 7, 1, readWriteModel.ModbusDataType_DINT},
		{"with quantity", "holding-register:1:REAL[2]", HoldingRegister, 0, 2, readWriteModel.ModbusDataType_REAL},
		{"highest address", "holding-register:65535:INT", HoldingRegister, 65534, 1, readWriteModel.ModbusDataType_INT},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tag := parseTag(t, test.address)
			assert.Equal(t, test.tagType, tag.TagType)
			assert.Equal(t, test.wire, tag.Address)
			assert.Equal(t, test.quantity, tag.Quantity)
			assert.Equal(t, test.datatype, tag.Datatype)
		})
	}
}

// An address without an explicit datatype is legal; the bit areas default to BOOL and the register
// areas to INT, as they do in plc4j's per-area tag factories.
func TestTagHandler_ParseTag_defaultsTheDatatype(t *testing.T) {
	tests := []struct {
		address  string
		datatype readWriteModel.ModbusDataType
	}{
		{"coil:1", readWriteModel.ModbusDataType_BOOL},
		{"0x00001", readWriteModel.ModbusDataType_BOOL},
		{"discrete-input:1", readWriteModel.ModbusDataType_BOOL},
		{"1x00001", readWriteModel.ModbusDataType_BOOL},
		{"input-register:1", readWriteModel.ModbusDataType_INT},
		{"3x00001", readWriteModel.ModbusDataType_INT},
		{"holding-register:1", readWriteModel.ModbusDataType_INT},
		{"4x00001", readWriteModel.ModbusDataType_INT},
		{"extended-register:1", readWriteModel.ModbusDataType_INT},
		{"6x00001", readWriteModel.ModbusDataType_INT},
	}
	for _, test := range tests {
		t.Run(test.address, func(t *testing.T) {
			assert.Equal(t, test.datatype, parseTag(t, test.address).Datatype)
		})
	}
	t.Run("quantity without datatype", func(t *testing.T) {
		tag := parseTag(t, "holding-register:1[4]")
		assert.Equal(t, readWriteModel.ModbusDataType_INT, tag.Datatype)
		assert.Equal(t, uint16(4), tag.Quantity)
	})
}

func TestTagHandler_ParseTag_rejectsInvalidAddresses(t *testing.T) {
	tests := []struct {
		name    string
		address string
	}{
		{"unknown datatype", "holding-register:1:NOPE"},
		{"no area", "1:INT"},
		{"garbage", "this is not an address"},
		{"logical address zero", "coil:0:BOOL"},
		{"numeric logical address zero", "4x00000:INT"},
		{"address beyond the address space", "holding-register:65537:INT"},
		{"extended register address zero", "extended-register:0:INT"},
		{"quantity zero", "holding-register:1:INT[0]"},
		// plc4j rejects a range whose last address reaches the end of the address space
		// (ModbusTagHoldingRegister.of checks address + quantity > REGISTER_MAXADDRESS).
		{"range reaching the end of the address space", "holding-register:65536:INT"},
		{"range running past the address space", "holding-register:65535:INT[2]"},
		{"too many coils", "coil:1:BOOL[2001]"},
		{"too many discrete inputs", "discrete-input:1:BOOL[2001]"},
		{"too many holding registers", "holding-register:1:INT[126]"},
		{"too many input registers", "input-register:1:INT[126]"},
		{"too many extended registers", "extended-register:1:INT[126]"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tag, err := NewTagHandler().ParseTag(test.address)
			assert.Error(t, err)
			assert.Nil(t, tag)
		})
	}
}

// The extended register area is addressed starting at zero rather than one, so its logical address
// is the address that goes onto the wire (plc4j ModbusTagExtendedRegister.getLogicalAddress). Every
// other area is shifted by one.
func TestModbusTag_extendedRegisterAddressesAreNotShifted(t *testing.T) {
	assert.Equal(t, uint16(7), parseTag(t, "extended-register:7:DINT").Address)
	assert.Equal(t, "6x00007:DINT[1]", parseTag(t, "extended-register:7:DINT").GetAddressString())

	assert.Equal(t, uint16(6), parseTag(t, "holding-register:7:DINT").Address)
	assert.Equal(t, "4x00007:DINT[1]", parseTag(t, "holding-register:7:DINT").GetAddressString())
}

// A quantity that was spelled out but doesn't fit into the field must be an error. Silently
// falling back to a single element would read or write the wrong amount of data.
func TestNewModbusPlcTagFromStrings_rejectsAnUnparsableQuantity(t *testing.T) {
	tag, err := NewModbusPlcTagFromStrings(HoldingRegister, "1", "99999999999999999999", "", readWriteModel.ModbusDataType_INT, tagConfig{})
	assert.Error(t, err)
	assert.Nil(t, tag)
}

func TestNewModbusPlcTagFromStrings_defaultsQuantityToOne(t *testing.T) {
	tag, err := NewModbusPlcTagFromStrings(HoldingRegister, "1", "", "", readWriteModel.ModbusDataType_INT, tagConfig{})
	require.NoError(t, err)
	assert.Equal(t, uint16(1), tag.(modbusTag).Quantity)
}

// The Go SPI treats the upper bound of an ArrayInfo as exclusive - DefaultArrayInfo.GetSize is
// UpperBound-LowerBound and every consumer subtracts the two to get the element count. Modbus
// follows that convention, which is why it differs from plc4j's inclusive quantity-1.
func TestModbusTag_GetArrayInfoUsesAnExclusiveUpperBound(t *testing.T) {
	arrayInfo := parseTag(t, "holding-register:1:INT[5]").GetArrayInfo()
	require.Len(t, arrayInfo, 1)
	assert.Equal(t, uint32(0), arrayInfo[0].GetLowerBound())
	assert.Equal(t, uint32(5), arrayInfo[0].GetUpperBound())
	assert.Equal(t, uint32(5), arrayInfo[0].GetSize(), "the size must be the number of elements")

	// Same shape as what the other Go drivers produce for five elements.
	var reference apiModel.ArrayInfo = &spiModel.DefaultArrayInfo{LowerBound: 0, UpperBound: 5}
	assert.Equal(t, reference.GetSize(), arrayInfo[0].GetSize())

	// A single element isn't an array at all.
	assert.Empty(t, parseTag(t, "holding-register:1:INT").GetArrayInfo())
}

// STRING and WSTRING carry the length of one string in parentheses; nothing on the wire announces
// it, so without it the driver couldn't know how many bytes belong to the value (plc4j ModbusTag).
func TestTagHandler_ParseTag_stringLength(t *testing.T) {
	for _, test := range []struct {
		address      string
		datatype     readWriteModel.ModbusDataType
		stringLength uint16
		quantity     uint16
	}{
		{"holding-register:1:STRING(20)", readWriteModel.ModbusDataType_STRING, 20, 1},
		{"holding-register:1:WSTRING(20)", readWriteModel.ModbusDataType_WSTRING, 20, 1},
		{"4x00001:STRING(8)", readWriteModel.ModbusDataType_STRING, 8, 1},
		{"holding-register:1:STRING(20)[3]", readWriteModel.ModbusDataType_STRING, 20, 3},
		{"input-register:7:STRING(1)", readWriteModel.ModbusDataType_STRING, 1, 1},
	} {
		t.Run(test.address, func(t *testing.T) {
			tag := parseTag(t, test.address)
			assert.Equal(t, test.datatype, tag.Datatype)
			assert.Equal(t, test.stringLength, tag.StringLength)
			assert.Equal(t, test.quantity, tag.Quantity)
		})
	}
}

// Every data type that isn't a string has a string length of 1, which leaves the size arithmetic
// unchanged.
func TestTagHandler_ParseTag_nonStringsHaveAStringLengthOfOne(t *testing.T) {
	assert.Equal(t, uint16(1), parseTag(t, "holding-register:1:INT").StringLength)
	assert.Equal(t, uint16(1), parseTag(t, "coil:1").StringLength)
	assert.Equal(t, uint16(1), parseTag(t, "holding-register:1:CHAR[4]").StringLength)
}

func TestTagHandler_ParseTag_rejectsABadStringLength(t *testing.T) {
	for _, test := range []struct {
		name    string
		address string
	}{
		{"string without a length", "holding-register:1:STRING"},
		{"wstring without a length", "holding-register:1:WSTRING"},
		{"string with a length of zero", "holding-register:1:STRING(0)"},
		{"length on a non-string", "holding-register:1:INT(20)"},
		{"length on a char", "holding-register:1:CHAR(20)"},
	} {
		t.Run(test.name, func(t *testing.T) {
			tag, err := NewTagHandler().ParseTag(test.address)
			assert.Error(t, err)
			assert.Nil(t, tag)
		})
	}
}

// A tag may name the unit id it is addressed at and the byte order it is encoded in, written in
// curly braces behind the address (plc4j TagConfigParser).
func TestTagHandler_ParseTag_tagConfig(t *testing.T) {
	t.Run("unit id", func(t *testing.T) {
		tag := parseTag(t, "holding-register:1:INT{unit-id: 7}")
		require.NotNil(t, tag.UnitId)
		assert.Equal(t, uint8(7), *tag.UnitId)
		assert.Nil(t, tag.ByteOrder)
	})
	t.Run("byte order", func(t *testing.T) {
		tag := parseTag(t, "holding-register:1:DINT{byte-order: 'LITTLE_ENDIAN_BYTE_SWAP'}")
		require.NotNil(t, tag.ByteOrder)
		assert.Equal(t, LittleEndianByteSwapOrder, *tag.ByteOrder)
		assert.Nil(t, tag.UnitId)
	})
	t.Run("both, on a numeric address", func(t *testing.T) {
		tag := parseTag(t, `4x00001:REAL{unit-id: 3, byte-order: "BIG_ENDIAN_BYTE_SWAP"}`)
		require.NotNil(t, tag.UnitId)
		require.NotNil(t, tag.ByteOrder)
		assert.Equal(t, uint8(3), *tag.UnitId)
		assert.Equal(t, BigEndianByteSwapOrder, *tag.ByteOrder)
	})
	t.Run("together with a quantity and a string length", func(t *testing.T) {
		tag := parseTag(t, "holding-register:1:STRING(4)[2]{unit-id: 9}")
		assert.Equal(t, uint16(4), tag.StringLength)
		assert.Equal(t, uint16(2), tag.Quantity)
		require.NotNil(t, tag.UnitId)
		assert.Equal(t, uint8(9), *tag.UnitId)
	})
	t.Run("an unquoted enum name", func(t *testing.T) {
		tag := parseTag(t, "holding-register:1:DINT{byte-order: LITTLE_ENDIAN}")
		require.NotNil(t, tag.ByteOrder)
		assert.Equal(t, LittleEndianOrder, *tag.ByteOrder)
	})
	// plc4j's TagConfigParser collects every key into a map and lets each tag pick out the ones it
	// knows, so a setting this driver has no use for must not make the whole address unparsable.
	t.Run("an unknown parameter is ignored", func(t *testing.T) {
		tag := parseTag(t, "holding-register:1:INT{no-such-thing: 1, unit-id: 7}")
		require.NotNil(t, tag.UnitId)
		assert.Equal(t, uint8(7), *tag.UnitId)
		assert.Nil(t, tag.ByteOrder)
	})
	t.Run("an empty config leaves the defaults", func(t *testing.T) {
		tag := parseTag(t, "holding-register:1:INT{}")
		assert.Nil(t, tag.UnitId)
		assert.Nil(t, tag.ByteOrder)
	})
	t.Run("no config at all leaves the defaults", func(t *testing.T) {
		tag := parseTag(t, "holding-register:1:INT")
		assert.Nil(t, tag.UnitId)
		assert.Nil(t, tag.ByteOrder)
	})
}

func TestTagHandler_ParseTag_rejectsABadTagConfig(t *testing.T) {
	for _, test := range []struct {
		name    string
		address string
	}{
		{"unknown byte order", "holding-register:1:INT{byte-order: 'MIDDLE_ENDIAN'}"},
		{"unknown unquoted byte order", "holding-register:1:INT{byte-order: MIDDLE_ENDIAN}"},
		{"unit id beyond a byte", "holding-register:1:INT{unit-id: 256}"},
		{"negative unit id", "holding-register:1:INT{unit-id: -1}"},
	} {
		t.Run(test.name, func(t *testing.T) {
			tag, err := NewTagHandler().ParseTag(test.address)
			assert.Error(t, err)
			assert.Nil(t, tag)
		})
	}
}

// A tag that names neither falls back to what the connection was configured with.
func TestModbusTag_resolvesUnitIdAndByteOrderAgainstTheConnectionDefaults(t *testing.T) {
	plain := parseTag(t, "holding-register:1:INT")
	assert.Equal(t, uint8(2), plain.resolveUnitId(2))
	assert.Equal(t, LittleEndianOrder, plain.resolveByteOrder(LittleEndianOrder))

	own := parseTag(t, "holding-register:1:INT{unit-id: 7, byte-order: 'BIG_ENDIAN'}")
	assert.Equal(t, uint8(7), own.resolveUnitId(2))
	assert.Equal(t, BigEndianOrder, own.resolveByteOrder(LittleEndianOrder))
}

// An address has to survive being printed and parsed again, or a tag couldn't be handed on.
func TestModbusTag_GetAddressStringRoundTrips(t *testing.T) {
	for _, address := range []string{
		"holding-register:1:INT",
		"holding-register:1:STRING(20)[3]",
		"holding-register:1:REAL{unit-id: 7}",
		"holding-register:1:DINT{byte-order: 'LITTLE_ENDIAN_BYTE_SWAP'}",
		"holding-register:1:WSTRING(4)[2]{unit-id: 9, byte-order: 'BIG_ENDIAN_BYTE_SWAP'}",
		"extended-register:7:DINT",
		"extended-register:12345:INT[2]",
	} {
		t.Run(address, func(t *testing.T) {
			tag := parseTag(t, address)
			assert.Equal(t, tag, parseTag(t, tag.GetAddressString()))
		})
	}
}

// The number of registers a request has to ask for grows with the declared string length; with a
// literal 1 in its place a STRING(20) would read a single register.
func TestModbusTag_lengthWordsCountsTheStringLength(t *testing.T) {
	for _, test := range []struct {
		address  string
		expected uint16
	}{
		{"holding-register:1:INT", 1},
		{"holding-register:1:BOOL", 1},
		{"holding-register:1:REAL", 2},
		{"holding-register:1:INT[4]", 4},
		{"holding-register:1:STRING(20)", 10},
		{"holding-register:1:STRING(20)[3]", 30},
		{"holding-register:1:WSTRING(20)", 20},
		{"holding-register:1:STRING(3)", 2},
		// Several values narrower than a byte are packed, so three BOOLs share one register
		// instead of taking one each.
		{"holding-register:1:BOOL[3]", 1},
		{"holding-register:1:BOOL[17]", 2},
		{"holding-register:1:SINT[3]", 2},
	} {
		t.Run(test.address, func(t *testing.T) {
			words, err := parseTag(t, test.address).lengthWords()
			require.NoError(t, err)
			assert.Equal(t, test.expected, words)
		})
	}
}

// A payload that doesn't fit into the 16 bit quantity field of a request is an error, not a
// truncated request that would silently read the wrong amount of data.
func TestModbusTag_lengthWordsRejectsAnOversizedPayload(t *testing.T) {
	_, err := parseTag(t, "holding-register:1:STRING(65535)[125]").lengthWords()
	assert.Error(t, err)
}
