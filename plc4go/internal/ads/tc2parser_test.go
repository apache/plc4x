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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	driverModel "github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
)

// A TwinCAT2 device uploads only its complex types, and the names it uses for the primitives are
// not all IEC names ("INT16" for INT, "DT" for DATE_AND_TIME). The synthesized entry therefore
// carries the IEC name as its main name, which is what the value resolution reads, and the name
// the device declared as its secondary name.
func TestSynthesizePrimitiveType_NameAndSize(t *testing.T) {
	for _, test := range []struct {
		declared string
		iec      string
		size     uint32
	}{
		{"BOOL", "BOOL", 1},
		{"BYTE", "BYTE", 1},
		{"SINT", "SINT", 1},
		{"USINT", "USINT", 1},
		{"WORD", "WORD", 2},
		{"INT", "INT", 2},
		{"INT16", "INT", 2},
		{"UINT", "UINT", 2},
		{"UINT16", "UINT", 2},
		{"DWORD", "DWORD", 4},
		{"DINT", "DINT", 4},
		{"INT32", "DINT", 4},
		{"UDINT", "UDINT", 4},
		{"UINT32", "UDINT", 4},
		{"REAL", "REAL", 4},
		{"LREAL", "LREAL", 8},
		{"LWORD", "LWORD", 8},
		{"LINT", "LINT", 8},
		{"ULINT", "ULINT", 8},
		{"TIME", "TIME", 4},
		{"TOD", "TIME_OF_DAY", 4},
		{"TIME_OF_DAY", "TIME_OF_DAY", 4},
		{"DATE", "DATE", 4},
		{"DT", "DATE_AND_TIME", 4},
		{"DATE_AND_TIME", "DATE_AND_TIME", 4},
	} {
		t.Run(test.declared, func(t *testing.T) {
			entry := SynthesizePrimitiveType(test.declared)
			require.NotNil(t, entry)
			assert.Equal(t, test.iec, entry.GetMainName(), "main name carries the IEC type")
			assert.Equal(t, test.declared, entry.GetSecondaryName(), "secondary name carries the declared type")
			assert.Equal(t, test.size, entry.GetSize())
		})
	}
}

// STRING(n) and WSTRING(n) keep their declared form in both names, because that is where the
// length is read from downstream. The size covers the terminating null character.
func TestSynthesizePrimitiveType_Strings(t *testing.T) {
	for _, test := range []struct {
		declared string
		size     uint32
	}{
		{"STRING(1)", 2},
		{"STRING(80)", 81},
		{"STRING(255)", 256},
		{"WSTRING(1)", 4},
		{"WSTRING(10)", 22},
	} {
		t.Run(test.declared, func(t *testing.T) {
			entry := SynthesizePrimitiveType(test.declared)
			require.NotNil(t, entry)
			assert.Equal(t, test.declared, entry.GetMainName())
			assert.Equal(t, test.declared, entry.GetSecondaryName())
			assert.Equal(t, test.size, entry.GetSize())
		})
	}
}

// Anything that is not a primitive is left to the uploaded table.
func TestSynthesizePrimitiveType_UnknownTypesAreNotSynthesized(t *testing.T) {
	for _, name := range []string{
		"",
		"ST_Recipe",
		"POINTER TO BYTE",
		"ARRAY [0..9] OF INT",
		"STRING(n)",
		"STRING(",
		"WSTRING()",
	} {
		t.Run(name, func(t *testing.T) {
			assert.Nil(t, SynthesizePrimitiveType(name))
		})
	}
}

func TestDriverContext_lookupDataType(t *testing.T) {
	uploaded := SynthesizePrimitiveType("LREAL")
	require.NotNil(t, uploaded)
	driverContext := &DriverContext{
		// The synthesis is only reached on a device that reported itself as TwinCAT 2.
		isTwinCat2: true,
		dataTypeTable: map[string]driverModel.AdsDataTypeTableEntry{
			"ST_Recipe": uploaded,
			// A TwinCAT3 device uploads the primitives too; the uploaded entry must win.
			"INT16": uploaded,
		},
	}

	t.Run("uploaded entry wins", func(t *testing.T) {
		entry, ok := driverContext.lookupDataType("ST_Recipe")
		require.True(t, ok)
		assert.Same(t, uploaded, entry)

		entry, ok = driverContext.lookupDataType("INT16")
		require.True(t, ok)
		assert.Same(t, uploaded, entry)
	})

	t.Run("primitive missing from the table is synthesized", func(t *testing.T) {
		entry, ok := driverContext.lookupDataType("DT")
		require.True(t, ok)
		require.NotNil(t, entry)
		assert.Equal(t, "DATE_AND_TIME", entry.GetMainName())
	})

	t.Run("unknown type is not resolved", func(t *testing.T) {
		entry, ok := driverContext.lookupDataType("ST_Unknown")
		assert.False(t, ok)
		assert.Nil(t, entry)
	})

	t.Run("nothing is synthesized for a device that is not TwinCAT2", func(t *testing.T) {
		// A TwinCAT3 device uploads its primitives, so a name missing from the table is a
		// name the device does not have, and inventing an entry for it would hide that.
		twinCat3 := &DriverContext{
			dataTypeTable: map[string]driverModel.AdsDataTypeTableEntry{"ST_Recipe": uploaded},
		}

		entry, ok := twinCat3.lookupDataType("DT")
		assert.False(t, ok, "DT is a TwinCAT2 spelling and must not be invented here")
		assert.Nil(t, entry)

		entry, ok = twinCat3.lookupDataType("ST_Recipe")
		require.True(t, ok, "an uploaded entry still resolves")
		assert.Same(t, uploaded, entry)
	})
}
