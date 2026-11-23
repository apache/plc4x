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

// Synthesis of the primitive data types a TwinCAT2 device does not upload.
//
// TwinCAT2 puts only its complex types in the data-type table, so a symbol, a struct member or a
// direct address that names a primitive names a type the table does not contain. TwinCAT3 uploads
// the primitives, so nothing here is reached on a TwinCAT3 device.

package ads

import (
	"strconv"
	"strings"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
)

// primitiveInfo maps a TwinCAT type name to (plc4x IEC value-type name, byte size).
// TwinCAT2 only uploads complex types in the data-type table, so primitive types
// referenced by symbols / struct members must be synthesized on demand.
func primitiveInfo(name string) (string, uint32, bool) {
	switch name {
	case "BOOL":
		return "BOOL", 1, true
	case "BYTE":
		return "BYTE", 1, true
	case "SINT":
		return "SINT", 1, true
	case "USINT":
		return "USINT", 1, true
	case "WORD":
		return "WORD", 2, true
	case "INT", "INT16":
		return "INT", 2, true
	case "UINT", "UINT16":
		return "UINT", 2, true
	case "DWORD":
		return "DWORD", 4, true
	case "DINT", "INT32":
		return "DINT", 4, true
	case "UDINT", "UINT32":
		return "UDINT", 4, true
	case "REAL":
		return "REAL", 4, true
	case "LREAL":
		return "LREAL", 8, true
	case "LWORD":
		return "LWORD", 8, true
	case "LINT":
		return "LINT", 8, true
	case "ULINT":
		return "ULINT", 8, true
	case "TIME":
		return "TIME", 4, true
	case "TOD", "TIME_OF_DAY":
		return "TIME_OF_DAY", 4, true
	case "DATE":
		return "DATE", 4, true
	case "DT", "DATE_AND_TIME":
		return "DATE_AND_TIME", 4, true
	}
	return "", 0, false
}

// SynthesizePrimitiveType builds a minimal AdsDataTypeTableEntry for a primitive
// (or STRING/WSTRING) type name. Returns nil for unknown / complex types.
// Used as a fallback when a type name isn't present in the uploaded data-type
// table (the normal case on TwinCAT2 for primitives).
//
// The IEC name goes in the MAIN name and the name TwinCAT used in the secondary one:
// getPlcValueForAdsDataTypeTableEntry resolves a value type from the main name, and the
// names TwinCAT2 puts in its tables are not all IEC names ("INT16" for INT, "DT" for
// DATE_AND_TIME). STRING(n)/WSTRING(n) keep their declared form in both, because that is
// where the length is read from.
func SynthesizePrimitiveType(name string) readWriteModel.AdsDataTypeTableEntry {
	var mainName string
	var size uint32
	switch {
	case strings.HasPrefix(name, "STRING(") && strings.HasSuffix(name, ")"):
		n, err := strconv.Atoi(name[len("STRING(") : len(name)-1])
		if err != nil {
			return nil
		}
		mainName = name // keep "STRING(n)" so the length is parsed downstream
		size = uint32(n) + 1
	case strings.HasPrefix(name, "WSTRING(") && strings.HasSuffix(name, ")"):
		n, err := strconv.Atoi(name[len("WSTRING(") : len(name)-1])
		if err != nil {
			return nil
		}
		mainName = name
		size = (uint32(n) + 1) * 2
	default:
		var ok bool
		mainName, size, ok = primitiveInfo(name)
		if !ok {
			return nil
		}
	}
	// Only the size and the two names carry anything here; every other field an uploaded entry
	// would have stays at its zero value, which is what the builder leaves them at.
	entry, err := readWriteModel.NewAdsDataTypeTableEntryBuilder().
		WithSize(size).
		WithMainName(mainName).
		WithSecondaryName(name).
		Build()
	if err != nil {
		return nil
	}
	return entry
}
