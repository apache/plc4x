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

package profinet

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/internal/profinet/gsdml"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

func TestParseTag(t *testing.T) {
	for address, want := range map[string]Tag{
		"1.1.INPUT.0:BOOL[8]":  {Slot: 1, Subslot: 1, Index: 0, ValueType: apiValues.BOOL, Count: 8},
		"1.1.INPUT[8]:BOOL":    {Slot: 1, Subslot: 1, Index: 0, ValueType: apiValues.BOOL, Count: 8},
		"2.1.OUTPUT.3:REAL":    {Slot: 2, Subslot: 1, Output: true, Index: 3, ValueType: apiValues.REAL, Count: 1},
		"0.32768.INPUT:USINT":  {Slot: 0, Subslot: 32768, ValueType: apiValues.USINT, Count: 1},
		"5.1.OUTPUT.1:BYTE[4]": {Slot: 5, Subslot: 1, Output: true, Index: 1, ValueType: apiValues.BYTE, Count: 4},
	} {
		tag, err := parseTag(address)
		require.NoError(t, err, address)
		assert.Equal(t, want, tag, address)
		again, err := parseTag(tag.GetAddressString())
		require.NoError(t, err)
		assert.Equal(t, want, again, "round trip of %s", address)
	}
	for _, bad := range []string{"", "1.1:BOOL", "1.1.SIDEWAYS:BOOL", "1.1.INPUT:NOPE", "70000.1.INPUT:BOOL", "1.1.INPUT:BOOL[0]"} {
		_, err := parseTag(bad)
		assert.Error(t, err, bad)
	}
}

func TestDecodeEncode(t *testing.T) {
	// An Unsigned16 used as bits: bit 0 is the least significant bit of the big-endian value.
	bits, err := decodeItem(Tag{ValueType: apiValues.BOOL, Count: 16}, []byte{0x80, 0x01})
	require.NoError(t, err)
	list := bits.GetList()
	assert.True(t, list[0].GetBool())
	assert.False(t, list[1].GetBool())
	assert.True(t, list[15].GetBool())

	real, err := decodeItem(Tag{ValueType: apiValues.REAL, Count: 1}, []byte{0x41, 0x48, 0x00, 0x00})
	require.NoError(t, err)
	assert.InDelta(t, 12.5, real.GetFloat32(), 0)

	i16, err := decodeItem(Tag{ValueType: apiValues.INT, Count: 1}, []byte{0xFF, 0xFE})
	require.NoError(t, err)
	assert.EqualValues(t, -2, i16.GetInt16())

	_, err = decodeItem(Tag{ValueType: apiValues.UDINT, Count: 1}, []byte{1, 2})
	assert.Error(t, err, "size mismatch")

	data, mask, err := encodeItem(Tag{ValueType: apiValues.BOOL, Count: 2}, 1, spiValues.NewPlcList([]apiValues.PlcValue{spiValues.NewPlcBOOL(true), spiValues.NewPlcBOOL(false)}))
	require.NoError(t, err)
	assert.Equal(t, []byte{0x01}, data)
	assert.Equal(t, []byte{0x03}, mask, "only the written bits")

	data, mask, err = encodeItem(Tag{ValueType: apiValues.UINT, Count: 1}, 2, spiValues.NewPlcUINT(0x1234))
	require.NoError(t, err)
	assert.Equal(t, []byte{0x12, 0x34}, data)
	assert.Equal(t, []byte{0xFF, 0xFF}, mask)

	_, _, err = encodeItem(Tag{ValueType: apiValues.UINT, Count: 1}, 4, spiValues.NewPlcUINT(1))
	assert.Error(t, err)
}

func TestLocate(t *testing.T) {
	submodule := &gsdml.SubmoduleItem{IOData: gsdml.IOData{
		Inputs:  []gsdml.IODataSet{{DataItems: []gsdml.DataItem{{DataType: "Unsigned8", UseAsBits: true}, {DataType: "Float32"}}}},
		Outputs: []gsdml.IODataSet{{DataItems: []gsdml.DataItem{{DataType: "Unsigned16"}}}},
	}}
	layout, err := submodule.IOData.Layout()
	require.NoError(t, err)
	p, err := planFrom([]realSubmodule{{Slot: 0, Subslot: 1}, {Slot: 1, Subslot: 1}}, func(real realSubmodule) (gsdml.Resolved, gsdml.Layout, error) {
		if real.Slot == 0 {
			return gsdml.Resolved{Submodule: &gsdml.SubmoduleItem{}}, gsdml.Layout{IOPSLength: 1, IOCSLength: 1}, nil
		}
		return gsdml.Resolved{Submodule: submodule}, layout, nil
	})
	require.NoError(t, err)
	// Input frame: DAP IOPS at 0; slot 1 data at 1..5, IOPS 6, IOCS for its outputs 7.
	location, err := p.locate(Tag{Slot: 1, Subslot: 1, Index: 1, ValueType: apiValues.REAL, Count: 1})
	require.NoError(t, err)
	assert.Equal(t, 2, location.offset)
	assert.Equal(t, 4, location.size)
	// Output frame: DAP IOCS 0, slot 1 IOCS 1, then its output data at 2.
	location, err = p.locate(Tag{Slot: 1, Subslot: 1, Output: true, ValueType: apiValues.UINT, Count: 1})
	require.NoError(t, err)
	assert.Equal(t, 2, location.offset)
	_, err = p.locate(Tag{Slot: 1, Subslot: 1, Index: 2})
	assert.Error(t, err)
	_, err = p.locate(Tag{Slot: 9, Subslot: 1})
	assert.Error(t, err)
}
