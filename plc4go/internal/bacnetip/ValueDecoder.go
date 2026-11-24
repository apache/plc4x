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
	"fmt"
	"reflect"
	"strings"
	"time"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

// appTagToPlcValue maps a BACnet ApplicationTag (the wire-level primitive
// representation) onto a plc4go PlcValue. A nil tag yields PlcNULL.
//
// The mapping mirrors how plc4j's BACnet driver renders these tags:
//   - Boolean        → PlcBOOL
//   - UnsignedInteger→ PlcULINT
//   - SignedInteger  → PlcLINT
//   - Real           → PlcREAL
//   - Double         → PlcLREAL
//   - CharacterString→ PlcSTRING
//   - OctetString    → PlcRawByteArray
//   - Enumerated     → PlcUDINT (the raw numeric; the property identifier on
//     the request side determines the enum schema)
//   - BitString      → PlcRawByteArray (packed MSB-first)
//   - Date           → PlcDATE       (year-1900 + month + day-of-month), or a
//     PlcStruct of the raw octets when the date is not one calendar day
//     (see bacnetDateToPlcValue)
//   - Time           → PlcTIME_OF_DAY, or a PlcStruct of the raw octets when
//     the time is not one time of day (see bacnetTimeToPlcValue)
//   - ObjectId       → PlcSTRING formatted as "<type>,<instance>"
//   - Null           → PlcNULL
func appTagToPlcValue(tag model.BACnetApplicationTag) apiValues.PlcValue {
	if tag == nil {
		return spiValues.NewPlcNULL()
	}
	switch t := tag.(type) {
	case model.BACnetApplicationTagNull:
		return spiValues.NewPlcNULL()
	case model.BACnetApplicationTagBoolean:
		return spiValues.NewPlcBOOL(t.GetActualValue())
	case model.BACnetApplicationTagUnsignedInteger:
		return spiValues.NewPlcULINT(t.GetActualValue())
	case model.BACnetApplicationTagSignedInteger:
		return spiValues.NewPlcLINT(int64(t.GetPayload().GetActualValue()))
	case model.BACnetApplicationTagReal:
		return spiValues.NewPlcREAL(t.GetActualValue())
	case model.BACnetApplicationTagDouble:
		return spiValues.NewPlcLREAL(t.GetActualValue())
	case model.BACnetApplicationTagCharacterString:
		return spiValues.NewPlcSTRING(t.GetValue())
	case model.BACnetApplicationTagOctetString:
		return spiValues.NewPlcRawByteArray(t.GetPayload().GetOctets())
	case model.BACnetApplicationTagEnumerated:
		return spiValues.NewPlcUDINT(t.GetActualValue())
	case model.BACnetApplicationTagBitString:
		return spiValues.NewPlcRawByteArray(bitsToBytes(t.GetPayload().GetData()))
	case model.BACnetApplicationTagDate:
		return bacnetDateToPlcValue(t.GetPayload())
	case model.BACnetApplicationTagTime:
		return bacnetTimeToPlcValue(t.GetPayload())
	case model.BACnetApplicationTagObjectIdentifier:
		return spiValues.NewPlcSTRING(fmt.Sprintf("%s,%d", t.GetObjectType(), t.GetInstanceNumber()))
	default:
		// Fallback: surface the type name + stringified value so the caller has
		// something to log. New ApplicationTag variants added to the spec should
		// produce a build-time miss and be added to the switch above.
		return spiValues.NewPlcSTRING(fmt.Sprintf("%T:%v", tag, tag))
	}
}

// constructedDataToPlcValue extracts a PlcValue from a BACnetConstructedData.
//
// BACnet has ~650 ConstructedData subtypes generated from the spec, each with a
// property-specific accessor. We use two strategies:
//
//  1. If the data is BACnetConstructedDataUnspecified (the parser's fallback for
//     unmapped (object, property) pairs), walk GetData() — each element wraps an
//     ApplicationTag we can decode via appTagToPlcValue. A single element yields
//     a scalar; multiple elements yield a PlcList.
//
//  2. Otherwise, look for a GetActualValue() method on the typed subtype via
//     reflection. The vast majority of typed ConstructedData expose this and the
//     return value is either a BACnetApplicationTag* (directly decodable) or a
//     more complex composite (DateTime, ObjectReference, ...) for which we
//     surface a string representation until per-type mappings are added.
func constructedDataToPlcValue(data model.BACnetConstructedData) apiValues.PlcValue {
	if data == nil {
		return spiValues.NewPlcNULL()
	}
	if u, ok := data.(model.BACnetConstructedDataUnspecified); ok {
		return elementsToPlcValue(u.GetData())
	}
	// Array-valued properties (OBJECT_LIST, PRIORITY_ARRAY, STATE_TEXT, ...) have
	// no single GetActualValue; decode them to a PlcList (or, for an array-index-0
	// read, the element count) before falling through to the scalar path.
	if pv, ok := arrayConstructedDataToPlcValue(data); ok {
		return pv
	}
	if v, ok := callGetActualValue(data); ok {
		if tag, ok := v.(model.BACnetApplicationTag); ok {
			return appTagToPlcValue(tag)
		}
		// Tagged bit strings (PROTOCOL_SERVICES_SUPPORTED, STATUS_FLAGS,
		// LIMIT_ENABLE, OBJECT_TYPES_SUPPORTED, ...) carry a bit-string payload
		// rather than a single scalar; surface the packed bytes (MSB-first), the
		// same representation used for a BitString application tag.
		if pv, ok := taggedBitStringToPlcValue(v); ok {
			return pv
		}
		if pv, ok := taggedEnumToPlcValue(v); ok {
			return pv
		}
		// Composite types: stringify until we add typed mappings (DateTime,
		// ObjectReference, BACnetTimeStamp, ...).
		return spiValues.NewPlcSTRING(fmt.Sprintf("%v", v))
	}
	return spiValues.NewPlcSTRING(fmt.Sprintf("%T:%v", data, data))
}

// arrayConstructedDataToPlcValue handles array-valued constructed data such as
// OBJECT_LIST. BACnet array properties expose a NumberOfDataElements accessor
// (set only when the client read array index 0, where the device returns just
// the element count) plus a per-property element slice. A whole-array read
// yields a PlcList of the decoded elements; an array-index-0 read yields the
// count as an unsigned PlcValue. Returns ok=false for non-array constructed
// data (no NumberOfDataElements accessor), so the caller falls through to the
// scalar GetActualValue path.
func arrayConstructedDataToPlcValue(data model.BACnetConstructedData) (apiValues.PlcValue, bool) {
	rv := reflect.ValueOf(data)
	if !rv.IsValid() {
		return nil, false
	}
	countMethod := rv.MethodByName("GetNumberOfDataElements")
	if !countMethod.IsValid() || countMethod.Type().NumIn() != 0 || countMethod.Type().NumOut() != 1 {
		return nil, false // not an array-valued constructed data
	}
	// Array-index-0 read: only the element count is present.
	if out := countMethod.Call(nil)[0]; out.IsValid() && !out.IsNil() {
		if tag, ok := out.Interface().(model.BACnetApplicationTag); ok {
			return appTagToPlcValue(tag), true
		}
	}
	// Whole-array read: find the element slice accessor and decode each entry.
	for i := 0; i < rv.NumMethod(); i++ {
		name := rv.Type().Method(i).Name
		if !strings.HasPrefix(name, "Get") || name == "GetNumberOfDataElements" {
			continue
		}
		method := rv.Method(i)
		mt := method.Type()
		if mt.NumIn() != 0 || mt.NumOut() != 1 || mt.Out(0).Kind() != reflect.Slice {
			continue
		}
		out := method.Call(nil)[0]
		elems := make([]apiValues.PlcValue, 0, out.Len())
		for j := 0; j < out.Len(); j++ {
			switch el := out.Index(j).Interface().(type) {
			case model.BACnetApplicationTag:
				elems = append(elems, appTagToPlcValue(el))
			case model.BACnetConstructedData:
				elems = append(elems, constructedDataToPlcValue(el))
			default:
				elems = append(elems, spiValues.NewPlcSTRING(fmt.Sprintf("%v", el)))
			}
		}
		return spiValues.NewPlcList(elems), true
	}
	return nil, false
}

// taggedBitStringToPlcValue handles BACnet *Tagged bit-string wrappers
// (BACnetServicesSupportedTagged, BACnetStatusFlagsTagged, BACnetLimitEnableTagged,
// BACnetObjectTypesSupportedTagged, ...) returned by GetActualValue on typed
// ConstructedData subtypes. These expose a GetPayload() returning a
// BACnetTagPayloadBitString whose GetData() is the unpacked bit slice. We pack
// it MSB-first into bytes and surface it as a PlcRawByteArray — identical to how
// a BitString application tag is decoded — so callers can walk the bits.
func taggedBitStringToPlcValue(v any) (apiValues.PlcValue, bool) {
	bs, ok := v.(interface {
		GetPayload() model.BACnetTagPayloadBitString
	})
	if !ok {
		return nil, false
	}
	payload := bs.GetPayload()
	if payload == nil {
		return nil, false
	}
	return spiValues.NewPlcRawByteArray(bitsToBytes(payload.GetData())), true
}

// taggedEnumToPlcValue handles BACnet's *Tagged enum wrappers (BACnetBinaryPV,
// BACnetReliability, BACnetEventState, ...) returned by GetActualValue on
// typed ConstructedData subtypes. These don't satisfy BACnetApplicationTag
// but expose a GetValue() method whose return is the unwrapped enum integer.
// Surfacing them as PlcUDINT keeps the API consistent with how Enumerated
// application tags are decoded (the property identifier on the request side
// determines the enum schema).
func taggedEnumToPlcValue(v any) (apiValues.PlcValue, bool) {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return nil, false
	}
	method := rv.MethodByName("GetValue")
	if !method.IsValid() || method.Type().NumIn() != 0 || method.Type().NumOut() != 1 {
		return nil, false
	}
	out := method.Call(nil)[0]
	switch out.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return spiValues.NewPlcUDINT(uint32(out.Uint())), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return spiValues.NewPlcLINT(out.Int()), true
	case reflect.Bool:
		return spiValues.NewPlcBOOL(out.Bool()), true
	case reflect.String:
		return spiValues.NewPlcSTRING(out.String()), true
	}
	return nil, false
}

// elementsToPlcValue collapses a list of BACnetConstructedDataElement to a
// PlcValue. A single-element list yields a scalar; multi-element yields PlcList.
func elementsToPlcValue(elements []model.BACnetConstructedDataElement) apiValues.PlcValue {
	switch len(elements) {
	case 0:
		return spiValues.NewPlcNULL()
	case 1:
		return appTagToPlcValue(elements[0].GetApplicationTag())
	}
	items := make([]apiValues.PlcValue, len(elements))
	for i, el := range elements {
		items[i] = appTagToPlcValue(el.GetApplicationTag())
	}
	return spiValues.NewPlcList(items)
}

// callGetActualValue invokes a GetActualValue() method on data via reflection
// if one is defined, returning the single result. Used to walk the typed
// ConstructedData subtypes without writing a 650-arm type switch.
func callGetActualValue(data model.BACnetConstructedData) (any, bool) {
	rv := reflect.ValueOf(data)
	if !rv.IsValid() {
		return nil, false
	}
	method := rv.MethodByName("GetActualValue")
	if !method.IsValid() {
		return nil, false
	}
	if method.Type().NumIn() != 0 || method.Type().NumOut() != 1 {
		return nil, false
	}
	out := method.Call(nil)
	if len(out) != 1 {
		return nil, false
	}
	return out[0].Interface(), true
}

// bitsToBytes packs a []bool BACnet bit-string into bytes MSB-first. Bit 0 of
// the input goes to bit 7 of byte 0 (i.e. the wire order BACnet uses on the
// payload side).
func bitsToBytes(bits []bool) []byte {
	if len(bits) == 0 {
		return []byte{}
	}
	out := make([]byte, (len(bits)+7)/8)
	for i, b := range bits {
		if b {
			out[i/8] |= 1 << (7 - uint(i%8))
		}
	}
	return out
}

// bacnetDateToPlcValue maps a BACnet Date (ASHRAE 135 clause 20.2.12) to a PlcDATE
// when its octets name one calendar day. A Date can also be a pattern: 255 is the
// wildcard for the year, month and day, month 13/14 means every odd/even month and
// day 32/33/34 the last/every odd/every even day of the month. time.Date would
// normalise those (and impossible days such as 2026-02-30) into a different, real
// date, so they are returned as a PlcStruct of the four raw octets instead, named
// after the wire fields: yearMinus1900, month, dayOfMonth and dayOfWeek. The day of
// week is not checked against a concrete date; devices commonly send 255 there.
func bacnetDateToPlcValue(payload model.BACnetTagPayloadDate) apiValues.PlcValue {
	yearMinus1900 := payload.GetYearMinus1900()
	month := payload.GetMonth()
	day := payload.GetDayOfMonth()
	if yearMinus1900 != 0xFF && month >= 1 && month <= 12 && day >= 1 && day <= 31 {
		date := time.Date(int(payload.GetYear()), time.Month(month), int(day), 0, 0, 0, 0, time.UTC)
		if date.Month() == time.Month(month) && date.Day() == int(day) {
			return spiValues.NewPlcDATE(date)
		}
	}
	return spiValues.NewPlcStructOrdered(map[string]apiValues.PlcValue{
		"yearMinus1900": spiValues.NewPlcUSINT(yearMinus1900),
		"month":         spiValues.NewPlcUSINT(month),
		"dayOfMonth":    spiValues.NewPlcUSINT(day),
		"dayOfWeek":     spiValues.NewPlcUSINT(payload.GetDayOfWeek()),
	}, []string{"yearMinus1900", "month", "dayOfMonth", "dayOfWeek"})
}

// bacnetTimeToPlcValue maps a BACnet Time (clause 20.2.13) to a PlcTIME_OF_DAY when
// its octets name one time of day. 255 in any octet is a wildcard, which a duration
// would turn into a wrong time (an hour of 255); such a time, and any octet out of
// range, is returned as a PlcStruct of the raw octets: hour, minute, second and
// fractional (hundredths of a second).
func bacnetTimeToPlcValue(payload model.BACnetTagPayloadTime) apiValues.PlcValue {
	hour, minute, second, fractional := payload.GetHour(), payload.GetMinute(), payload.GetSecond(), payload.GetFractional()
	if hour < 24 && minute < 60 && second < 60 && fractional < 100 {
		return spiValues.NewPlcTIME_OF_DAYFromMillisecondsSinceMidnight(
			((uint32(hour)*60+uint32(minute))*60+uint32(second))*1000 + uint32(fractional)*10)
	}
	return spiValues.NewPlcStructOrdered(map[string]apiValues.PlcValue{
		"hour":       spiValues.NewPlcUSINT(hour),
		"minute":     spiValues.NewPlcUSINT(minute),
		"second":     spiValues.NewPlcUSINT(second),
		"fractional": spiValues.NewPlcUSINT(fractional),
	}, []string{"hour", "minute", "second", "fractional"})
}
