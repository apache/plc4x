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
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
)

// plcValueToApplicationTag maps a PlcValue to a BACnet ApplicationTag, the
// inverse of appTagToPlcValue. The caller may pass a hint via the property
// identifier so we can prefer enumerated/character-string encoding when the
// PlcValue type alone is ambiguous (e.g. uint32 → Enumerated vs UnsignedInteger
// for properties like StatusFlags / Units).
//
// Returns an error if the PlcValue type isn't representable in BACnet.
func plcValueToApplicationTag(v apiValues.PlcValue, hint encodingHint) (model.BACnetApplicationTag, error) {
	if v == nil {
		return model.CreateBACnetApplicationTagNull(), nil
	}
	switch v.GetPlcValueType() {
	case apiValues.NULL:
		return model.CreateBACnetApplicationTagNull(), nil
	case apiValues.BOOL:
		return model.CreateBACnetApplicationTagBoolean(v.GetBool()), nil
	case apiValues.BYTE, apiValues.USINT, apiValues.UINT, apiValues.UDINT, apiValues.ULINT, apiValues.WORD, apiValues.DWORD, apiValues.LWORD:
		if hint == hintEnumerated {
			return model.CreateBACnetApplicationTagEnumerated(uint32(v.GetUint64())), nil
		}
		return model.CreateBACnetApplicationTagUnsignedInteger(uint(v.GetUint64())), nil
	case apiValues.SINT, apiValues.INT, apiValues.DINT, apiValues.LINT:
		return model.CreateBACnetApplicationTagSignedInteger(int(v.GetInt64())), nil
	case apiValues.REAL:
		return model.CreateBACnetApplicationTagReal(v.GetFloat32()), nil
	case apiValues.LREAL:
		return model.CreateBACnetApplicationTagDouble(v.GetFloat64()), nil
	case apiValues.STRING, apiValues.CHAR, apiValues.WSTRING, apiValues.WCHAR:
		return model.CreateBACnetApplicationTagCharacterString(model.BACnetCharacterEncoding_ISO_10646, v.GetString()), nil
	case apiValues.RAW_BYTE_ARRAY:
		return model.CreateBACnetApplicationTagOctetString(v.GetRaw()), nil
	default:
		return nil, errors.Errorf("PlcValue type %v cannot be encoded as a BACnet ApplicationTag", v.GetPlcValueType())
	}
}

// encodingHint nudges plcValueToApplicationTag when the source PlcValue type is
// ambiguous. Today the only hint is enumerated; future additions could include
// bit-string-from-uint or date-from-string.
type encodingHint uint8

const (
	hintNone encodingHint = iota
	hintEnumerated
)

// hintForProperty returns the encoding hint appropriate to a BACnet property
// identifier. PRESENT_VALUE on Binary*/Multistate* objects (and many discrete
// status fields) carries Enumerated; everything else defaults to none.
func hintForProperty(_ uint32) encodingHint {
	// Today we don't second-guess the caller. WriteProperty callers should
	// either pass a PlcUDINT and accept the default UnsignedInteger encoding,
	// or use a PlcREAL/PlcBOOL/PlcSTRING that maps unambiguously. Per-property
	// dispatch will be added once the integration tests expose a real device
	// that rejects an UnsignedInteger where Enumerated is required.
	return hintNone
}
