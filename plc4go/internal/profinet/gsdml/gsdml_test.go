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

package gsdml

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
)

// resources are the real vendor GSDML files plc4j's profinet-ng tests use.
const resources = "../../../../plc4j/drivers/profinet-ng/src/test/resources"

func load(t *testing.T, name string) *Index {
	t.Helper()
	path := filepath.Join(resources, name)
	if _, err := os.Stat(path); err != nil {
		t.Skipf("%s is not checked out", path)
	}
	start := time.Now()
	profile, err := ParseFile(path)
	require.NoError(t, err)
	t.Logf("parsed %s in %s", name, time.Since(start))
	require.True(t, profile.IsDeviceProfile())
	return NewIndex(profile)
}

func TestAdvantech(t *testing.T) {
	x := load(t, "GSDML-V2.2-Advantech-ADAM6100-20111216.xml")
	vendor, device, err := x.Profile.Identity()
	require.NoError(t, err)
	assert.Equal(t, uint16(0x01C6), vendor)
	assert.Equal(t, uint16(0x6100), device)

	daps := x.DeviceAccessPointsByIdent(0x61000000)
	require.Len(t, daps, 1)
	assert.Equal(t, "ID_DAP1", daps[0].ID)

	resolved, ok := x.Resolve(0x61510000, 0x61511000)
	require.True(t, ok)
	assert.Equal(t, "IDM_ADAM6151PN_1", resolved.ModuleID())
	assert.Equal(t, "ADAM-6151PN", x.Text(resolved.Module.ModuleInfo.Name.TextID))
	layout, err := resolved.Submodule.IOData.Layout()
	require.NoError(t, err)
	assert.Equal(t, Layout{InputLength: 2, OutputLength: 0, IOPSLength: 1, IOCSLength: 1}, layout)
	valueType, count, err := resolved.Submodule.IOData.Inputs[0].DataItems[0].ValueType()
	require.NoError(t, err)
	assert.Equal(t, apiValues.BOOL, valueType)
	assert.Equal(t, 8, count)

	port, ok := x.Resolve(0x61000000, 0x00000002)
	require.True(t, ok, "a DAP's port submodule resolves")
	assert.Equal(t, "ID_DAP1", port.ModuleID())
	layout, err = port.Submodule.IOData.Layout()
	require.NoError(t, err)
	assert.Zero(t, layout.InputLength+layout.OutputLength)

	_, ok = x.Resolve(0x61510000, 0x99999999)
	assert.False(t, ok)
	assert.Equal(t, "NO_SUCH_TEXT", x.Text("NO_SUCH_TEXT"))
}

func TestSiemensDAPs(t *testing.T) {
	x := load(t, "GSDML-V2.43-SIEMENS-SIMOCODEproVPN-20220817.xml")
	assert.Len(t, x.Profile.ProfileBody.ApplicationProcess.DeviceAccessPoints, 6)
	assert.Len(t, x.DeviceAccessPointsByIdent(0x10), 3, "three DAPs share ident 0x10")
	dap, ok := x.DeviceAccessPoint("DAP_3UF7_GP2_V3")
	require.True(t, ok)
	assert.Equal(t, "0x0000310", dap.ModuleIdentNumber)
}

// WAGO's 750 file is 9 MB with its I/O submodules in a SubmoduleList that modules only reference,
// the form plc4j does not resolve.
func TestWagoReferencedSubmodules(t *testing.T) {
	x := load(t, "GSDML-V2.33-wago-series750_753-20230115.xml")
	resolved, ok := x.Resolve(0x00008200, 0x44010001)
	require.True(t, ok)
	assert.True(t, strings.HasPrefix(resolved.Submodule.ID, "IDT_SUBM_06"), resolved.Submodule.ID)
	assert.NotNil(t, resolved.Module)
	layout, err := resolved.Submodule.IOData.Layout()
	require.NoError(t, err)
	assert.Equal(t, 2, layout.InputLength)
	item := resolved.Submodule.IOData.Inputs[0].DataItems[0]
	assert.Len(t, item.BitDataItems, 16)
	valueType, count, err := item.ValueType()
	require.NoError(t, err)
	assert.Equal(t, apiValues.WORD, valueType)
	assert.Equal(t, 1, count)

	var record *RecordDataItem
	for i := range resolved.Submodule.Records {
		if resolved.Submodule.Records[i].Index == 16384 {
			record = &resolved.Submodule.Records[i]
		}
	}
	require.NotNil(t, record)
	data, err := record.DefaultData()
	require.NoError(t, err)
	require.Len(t, data, 24)
	// The Const bytes, except byte 7: its Const is 0x81, but the diagnosis parameter at byte 7
	// bit 7 defaults to 0 and defaults go on top.
	assert.Equal(t, []byte{0x00, 0x01, 0x00, 0x00, 0x01, 0x01, 0x00, 0x01, 0x03, 0x03, 0x8A, 0x03}, data[:12])
}

// Every submodule of every file has a computable layout and every parameter record default data:
// what the connect step will need.
func TestEveryFileComputes(t *testing.T) {
	for _, name := range []string{
		"GSDML-V2.2-Advantech-ADAM6100-20111216.xml",
		"GSDML-V2.3-SIEMENS-SIMOCODEproVPN-20201104.xml",
		"GSDML-V2.33-wago-series750_753-20230115.xml",
		"GSDML-V2.43-SIEMENS-SIMOCODEproVPN-20220817.xml",
		"gsdml.xml",
	} {
		t.Run(name, func(t *testing.T) {
			x := load(t, name)
			process := &x.Profile.ProfileBody.ApplicationProcess
			var submodules []*SubmoduleItem
			for i := range process.DeviceAccessPoints {
				for j := range process.DeviceAccessPoints[i].VirtualSubmodules {
					submodules = append(submodules, &process.DeviceAccessPoints[i].VirtualSubmodules[j])
				}
			}
			for i := range process.Modules {
				submodules = append(submodules, x.SubmodulesOf(&process.Modules[i])...)
			}
			records := 0
			for _, submodule := range submodules {
				_, err := submodule.IOData.Layout()
				assert.NoError(t, err, submodule.ID)
				for _, set := range append(append([]IODataSet{}, submodule.IOData.Inputs...), submodule.IOData.Outputs...) {
					for _, item := range set.DataItems {
						_, _, err := item.ValueType()
						assert.NoError(t, err, submodule.ID)
					}
				}
				for _, record := range submodule.Records {
					_, err := record.DefaultData()
					assert.NoError(t, err, "%s record %d", submodule.ID, record.Index)
					records++
				}
			}
			t.Logf("%d submodules, %d records", len(submodules), records)
			assert.NotEmpty(t, submodules)
		})
	}
}

func TestFindInDirectory(t *testing.T) {
	if _, err := os.Stat(resources); err != nil {
		t.Skip("plc4j resources are not checked out")
	}
	profile, path, err := FindInDirectory(resources, 0x01C6, 0x6100)
	require.NoError(t, err)
	assert.Equal(t, "GSDML-V2.2-Advantech-ADAM6100-20111216.xml", filepath.Base(path))
	assert.Equal(t, "0x6100", profile.ProfileBody.DeviceIdentity.DeviceID)

	_, path, err = FindInDirectory(resources, 0x002A, 0x0904)
	require.NoError(t, err)
	assert.Equal(t, "GSDML-V2.43-SIEMENS-SIMOCODEproVPN-20220817.xml", filepath.Base(path), "the newer of two versions")

	_, _, err = FindInDirectory(resources, 0xFFFF, 0xFFFF)
	assert.Error(t, err)
	_, _, err = FindInDirectory(filepath.Join(resources, "missing"), 1, 1)
	assert.Error(t, err)
}

func TestDefaultData(t *testing.T) {
	no := false
	record := RecordDataItem{
		Index:  1,
		Length: 12,
		Consts: []RecordConst{{Data: "0xAA,0xBB", ByteOffset: 0}, {Data: "1,2", ByteOffset: 10}},
		Refs: []RecordRef{
			{DataType: "Bit", ByteOffset: 2, BitOffset: 7, DefaultValue: "1"},
			{DataType: "BitArea", ByteOffset: 2, BitOffset: 1, BitLength: 3, DefaultValue: "5"},
			{DataType: "BitArea", ByteOffset: 3, BitOffset: 6, BitLength: 4, DefaultValue: "0xF", Changeable: &no},
			{DataType: "Unsigned16", ByteOffset: 5, DefaultValue: "0x1234"},
			{DataType: "Integer8", ByteOffset: 7, DefaultValue: "-2"},
			{DataType: "Unsigned8", ByteOffset: 8},
			{DataType: "VisibleString", ByteOffset: 9, Length: 1, DefaultValue: "Z"},
		},
	}
	data, err := record.DefaultData()
	require.NoError(t, err)
	// byte 2: bit 7 set, bits 1..3 = 0b101 -> 1000_1010; bytes 3..4: bits 6..9 of a big-endian word
	assert.Equal(t, []byte{0xAA, 0xBB, 0x8A, 0x03, 0xC0, 0x12, 0x34, 0xFE, 0x00, 'Z', 0x01, 0x02}, data)

	for _, bad := range []RecordRef{
		{DataType: "Unsigned8", ByteOffset: 0, DefaultValue: "256"},
		{DataType: "Integer8", ByteOffset: 0, DefaultValue: "-129"},
		{DataType: "Bit", ByteOffset: 0, DefaultValue: "2"},
		{DataType: "Unsigned16", ByteOffset: 3, DefaultValue: "1"},
		{DataType: "Mystery", ByteOffset: 0},
	} {
		_, err := RecordDataItem{Length: 4, Refs: []RecordRef{bad}}.DefaultData()
		assert.Error(t, err, "%+v", bad)
	}
	_, err = RecordDataItem{Length: 2, Consts: []RecordConst{{Data: "1,2,3"}}}.DefaultData()
	assert.Error(t, err, "Const past the end")
}

func TestDataItemSizes(t *testing.T) {
	for _, c := range []struct {
		item DataItem
		want int
	}{
		{DataItem{DataType: "Unsigned8"}, 1},
		{DataItem{DataType: "Float32+Status8"}, 5},
		{DataItem{DataType: "OctetString", Length: 3}, 3},
		{DataItem{DataType: "TimeOfDay with date indication"}, 6},
		{DataItem{DataType: "F_MessageTrailer5Byte"}, 5},
	} {
		size, err := c.item.Size()
		require.NoError(t, err)
		assert.Equal(t, c.want, size, c.item.DataType)
	}
	_, err := DataItem{DataType: "VisibleString"}.Size()
	assert.Error(t, err, "a string needs a Length")
	valueType, count, err := DataItem{DataType: "Date"}.ValueType()
	require.NoError(t, err)
	assert.Equal(t, apiValues.BYTE, valueType, "plc4j can't map Date; here it is raw bytes")
	assert.Equal(t, 7, count)
}
