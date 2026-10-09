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

package ads

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	driverModel "github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
)

// symbolEntry encodes one ADS symbol table entry with the given 16 bit flags (TcAdsDef.h
// ADSSYMBOLFLAG_*), no comment and optionally a raw attribute block.
func symbolEntry(name string, flags uint16, attributes []byte) []byte {
	body := binary.LittleEndian.AppendUint32(nil, 0x4040) // group
	body = binary.LittleEndian.AppendUint32(body, 0)      // offset
	body = binary.LittleEndian.AppendUint32(body, 1)      // size
	body = binary.LittleEndian.AppendUint32(body, 33)     // dataType: ADST_BIT
	body = binary.LittleEndian.AppendUint16(body, flags)
	body = binary.LittleEndian.AppendUint16(body, 0) // second flag word
	body = binary.LittleEndian.AppendUint16(body, uint16(len(name)))
	body = binary.LittleEndian.AppendUint16(body, uint16(len("BOOL")))
	body = binary.LittleEndian.AppendUint16(body, 0) // commentLength
	body = append(body, name...)
	body = append(body, 0)
	body = append(body, "BOOL"...)
	body = append(body, 0, 0) // type terminator, empty comment terminator
	body = append(body, attributes...)
	return append(binary.LittleEndian.AppendUint32(nil, uint32(4+len(body))), body...)
}

// Bits 11..8 of the symbol flags are a 4 bit task/context id, not flags. Context 2 (0x0200) used to
// be read as the attribute flag, so the parser ran past the entry into the next one (GH-2773).
func TestSymbolTableEntry_ContextIdIsNotMistakenForAttributes(t *testing.T) {
	for context := uint16(0); context < 16; context++ {
		encoded := symbolEntry("MAIN.a", context<<8, nil)
		entry, err := driverModel.AdsSymbolTableEntryParse(t.Context(), encoded)
		require.NoError(t, err, "context %d", context)
		assert.Equal(t, uint8(context), entry.GetContextMask())
		assert.Nil(t, entry.GetAttributes())
		assert.False(t, entry.GetFlagStatic())
		assert.False(t, entry.GetFlagInitOnReset())
		assert.False(t, entry.GetFlagExtendedFlags())
		assert.Equal(t, len(encoded), int(entry.GetLengthInBytes(t.Context())), "context %d", context)
	}
}

func TestSymbolTableEntry_UpperFlagBitsMatchTcAdsDef(t *testing.T) {
	static, err := driverModel.AdsSymbolTableEntryParse(t.Context(), symbolEntry("MAIN.x", 0x2000, nil))
	require.NoError(t, err)
	assert.True(t, static.GetFlagStatic())
	assert.False(t, static.GetFlagInitOnReset())
	assert.False(t, static.GetFlagExtendedFlags())

	initOnReset, err := driverModel.AdsSymbolTableEntryParse(t.Context(), symbolEntry("MAIN.x", 0x4000, nil))
	require.NoError(t, err)
	assert.False(t, initOnReset.GetFlagStatic())
	assert.True(t, initOnReset.GetFlagInitOnReset())
	assert.False(t, initOnReset.GetFlagExtendedFlags())

	extended, err := driverModel.AdsSymbolTableEntryParse(t.Context(), symbolEntry("MAIN.x", 0x8000, nil))
	require.NoError(t, err)
	assert.False(t, extended.GetFlagStatic())
	assert.False(t, extended.GetFlagInitOnReset())
	assert.True(t, extended.GetFlagExtendedFlags())
}

func TestSymbolTableEntry_AttributeFlagSelectsAttributes(t *testing.T) {
	attributes := binary.LittleEndian.AppendUint16(nil, 1) // numAttributes
	attributes = append(attributes, byte(len("TcDisplayName")), 1)
	attributes = append(attributes, "TcDisplayName"...)
	attributes = append(attributes, 0, 'x', 0)
	entry, err := driverModel.AdsSymbolTableEntryParse(t.Context(), symbolEntry("MAIN.x", 0x1000|0x0300, attributes))
	require.NoError(t, err)
	assert.Equal(t, uint8(3), entry.GetContextMask())
	require.NotNil(t, entry.GetAttributes())
	require.Len(t, entry.GetAttributes().GetAttributes(), 1)
	assert.Equal(t, "TcDisplayName", entry.GetAttributes().GetAttributes()[0].GetName())
}
