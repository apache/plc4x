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
	"bytes"
	"context"
	"encoding/binary"
	"math"
	"time"
	"unicode/utf8"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/umas/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

// unknownDataTypeSizeIndex is the dataSizeIndex a read or write of a symbol whose type isn't a
// primitive UmasDataType carries: index 3, four bytes. plc4j hard-codes the same 3 with no
// explanation, and there is no capture to derive a better answer from - a struct or array symbol
// therefore reads its first four bytes.
const unknownDataTypeSizeIndex = uint8(3)

// terminatorByte is the NUL byte a UMAS string ends with.
const terminatorByte = byte(0x00)

// resolveDataType turns the type id of a symbol into a primitive UMAS data type. A custom type of
// the project (id customTypeIdBase and up) is not one, which is what the second return value says.
func resolveDataType(dataTypeId uint16) (readWriteModel.UmasDataType, bool) {
	if dataTypeId > math.MaxUint8 {
		return 0, false
	}
	return readWriteModel.UmasDataTypeByValue(uint8(dataTypeId))
}

// dataSizeIndexFor is the 4 bit size index a read or write reference carries for this type.
func dataSizeIndexFor(dataTypeId uint16) uint8 {
	if dataType, ok := resolveDataType(dataTypeId); ok {
		return dataType.RequestSize()
	}
	return unknownDataTypeSizeIndex
}

// isStringType says whether a symbol of this type is read and written as a byte array rather than as
// a fixed width value. STRING's request size is 17, which doesn't fit the 4 bit dataSizeIndex field,
// so plc4j asks for a byte array of the string's buffer size instead - and so does this driver.
func isStringType(dataTypeId uint16) bool {
	dataType, ok := resolveDataType(dataTypeId)
	return ok && dataType == readWriteModel.UmasDataType_STRING
}

// decodeReadResponse turns the payload of a read response into a plc4x value. Ported from plc4j's
// UmasConnection.parseReadResponse.
//
// The temporal types and STRING are decoded here rather than through the generated DataItem: the
// generated Go DataItem has no result for DATE, TIME_OF_DAY and DATE_AND_TIME at all (its parse
// falls through to "unsupported type") and reads a STRING as a single byte. plc4j hand-decodes the
// same five types for the same reason, so this is a port and not a workaround of the Go model alone.
func decodeReadResponse(ctx context.Context, dataTypeId uint16, block []byte) (apiValues.PlcValue, error) {
	if len(block) == 0 {
		return nil, errors.New("the read response carries no data")
	}
	dataType, ok := resolveDataType(dataTypeId)
	if !ok {
		// A custom type of the project: nothing here knows how to take it apart, so the caller gets
		// the bytes. plc4j answers with a PlcRawByteArray too.
		return spiValues.NewPlcRawByteArray(block), nil
	}

	switch dataType {
	case readWriteModel.UmasDataType_STRING:
		if terminator := bytes.IndexByte(block, terminatorByte); terminator >= 0 {
			return spiValues.NewPlcSTRING(string(block[:terminator])), nil
		}
		return spiValues.NewPlcSTRING(string(block)), nil

	case readWriteModel.UmasDataType_TIME:
		// A uint32 of milliseconds, little endian - not BCD, unlike the other temporal types.
		if len(block) < 4 {
			return nil, errors.Errorf("a TIME needs 4 bytes, got %d", len(block))
		}
		return spiValues.NewPlcTIMEFromMilliseconds(int64(binary.LittleEndian.Uint32(block[0:4]))), nil

	case readWriteModel.UmasDataType_DATE:
		// BCD: day(1) + month(1) + year(2, little endian pair of BCD bytes).
		if len(block) < 4 {
			return nil, errors.Errorf("a DATE needs 4 bytes, got %d", len(block))
		}
		day, err := decodeBcdByte(block[0], "day")
		if err != nil {
			return nil, err
		}
		month, err := decodeBcdByte(block[1], "month")
		if err != nil {
			return nil, err
		}
		year, err := decodeBcd16(block[2], block[3], "year")
		if err != nil {
			return nil, err
		}
		date, err := calendarDate(year, month, day)
		if err != nil {
			return nil, err
		}
		return spiValues.NewPlcDATE(date), nil

	case readWriteModel.UmasDataType_TOD:
		// BCD: centiseconds(1) + seconds(1) + minutes(1) + hours(1).
		if len(block) < 4 {
			return nil, errors.Errorf("a TIME_OF_DAY needs 4 bytes, got %d", len(block))
		}
		centiseconds, err := decodeBcdByte(block[0], "centiseconds")
		if err != nil {
			return nil, err
		}
		seconds, err := decodeBcdByte(block[1], "seconds")
		if err != nil {
			return nil, err
		}
		minutes, err := decodeBcdByte(block[2], "minutes")
		if err != nil {
			return nil, err
		}
		hours, err := decodeBcdByte(block[3], "hours")
		if err != nil {
			return nil, err
		}
		if hours > 23 || minutes > 59 || seconds > 59 || centiseconds > 99 {
			return nil, errors.Errorf("%02d:%02d:%02d.%02d is not a time of day", hours, minutes, seconds, centiseconds)
		}
		// plc4j drops the centiseconds field on the floor here; the mspec names it, so it is kept.
		return spiValues.NewPlcTIME_OF_DAY(time.Date(0, 1, 1, int(hours), int(minutes), int(seconds),
			int(centiseconds)*10*int(time.Millisecond), time.UTC)), nil

	case readWriteModel.UmasDataType_DATE_AND_TIME:
		// BCD: one reserved byte + seconds + minutes + hour + day + month + year(2).
		if len(block) < 8 {
			return nil, errors.Errorf("a DATE_AND_TIME needs 8 bytes, got %d", len(block))
		}
		seconds, err := decodeBcdByte(block[1], "seconds")
		if err != nil {
			return nil, err
		}
		minutes, err := decodeBcdByte(block[2], "minutes")
		if err != nil {
			return nil, err
		}
		hours, err := decodeBcdByte(block[3], "hours")
		if err != nil {
			return nil, err
		}
		day, err := decodeBcdByte(block[4], "day")
		if err != nil {
			return nil, err
		}
		month, err := decodeBcdByte(block[5], "month")
		if err != nil {
			return nil, err
		}
		year, err := decodeBcd16(block[6], block[7], "year")
		if err != nil {
			return nil, err
		}
		if hours > 23 || minutes > 59 || seconds > 59 {
			return nil, errors.Errorf("%02d:%02d:%02d is not a time of day", hours, minutes, seconds)
		}
		date, err := calendarDate(year, month, day)
		if err != nil {
			return nil, err
		}
		return spiValues.NewPlcDATE_AND_TIME(time.Date(date.Year(), date.Month(), date.Day(),
			int(hours), int(minutes), int(seconds), 0, time.UTC)), nil

	case readWriteModel.UmasDataType_UNKNOWN2, readWriteModel.UmasDataType_UNKNOWN3:
		// The dataIo has no case for the UNKNOWN types, in either language, so plc4j fails a read of
		// one outright. The mspec still gives all of them a request size of one byte, and
		// mapToPlcValueType (plc4j's included) calls these two BOOL - so one byte read as a boolean
		// is what the driver already claims a symbol of this type is, rather than a new guess.
		return spiValues.NewPlcBOOL(block[0] != 0), nil

	case readWriteModel.UmasDataType_UNKNOWN11, readWriteModel.UmasDataType_UNKNOWN12,
		readWriteModel.UmasDataType_UNKNOWN13, readWriteModel.UmasDataType_UNKNOWN17,
		readWriteModel.UmasDataType_UNKNOWN18, readWriteModel.UmasDataType_UNKNOWN19,
		readWriteModel.UmasDataType_UNKNOWN20, readWriteModel.UmasDataType_UNKNOWN24:
		// Same as above, for the ones mapToPlcValueType calls BYTE.
		return spiValues.NewPlcBYTE(block[0]), nil

	default:
		// Everything else goes through the generated dataIo. The buffer has to be little endian:
		// none of the DataItem fields pin a byte order, so the buffer is what decides, and UMAS
		// payloads are little endian throughout.
		readBuffer := utils.NewReadBufferByteBased(block, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
		value, err := readWriteModel.DataItemParseWithBuffer(ctx, readBuffer, dataType, 1)
		if err != nil {
			return nil, errors.Wrapf(err, "Error parsing a %s out of the read response", dataType)
		}
		if value == nil {
			return nil, errors.Errorf("the data item parser has no value for a %s", dataType)
		}
		return value, nil
	}
}

// encodeWriteValue turns a plc4x value into the payload bytes of a write request. Ported from
// plc4j's UmasConnection.serializeForType, which spells every type out by hand rather than going
// through the generated dataIo - the temporal types would have no serializer otherwise.
func encodeWriteValue(dataTypeId uint16, value apiValues.PlcValue) ([]byte, error) {
	if value == nil {
		return nil, errors.New("there is no value to write")
	}
	dataType, ok := resolveDataType(dataTypeId)
	if !ok {
		// A custom type of the project. plc4j writes the raw bytes of the value if it has any, and
		// gives up otherwise.
		if value.IsRaw() {
			return value.GetRaw(), nil
		}
		return nil, errors.Errorf("Can't serialize a %s for the unknown UMAS data type %d", value.GetPlcValueType(), dataTypeId)
	}

	switch dataType {
	case readWriteModel.UmasDataType_BOOL, readWriteModel.UmasDataType_EBOOL,
		readWriteModel.UmasDataType_UNKNOWN2, readWriteModel.UmasDataType_UNKNOWN3:
		boolean, ok := asBool(value)
		if !ok {
			return nil, errors.Errorf("a %s is not a %s", value.GetPlcValueType(), dataType)
		}
		if boolean {
			return []byte{1}, nil
		}
		return []byte{0}, nil

	case readWriteModel.UmasDataType_BYTE,
		readWriteModel.UmasDataType_UNKNOWN11, readWriteModel.UmasDataType_UNKNOWN12,
		readWriteModel.UmasDataType_UNKNOWN13, readWriteModel.UmasDataType_UNKNOWN17,
		readWriteModel.UmasDataType_UNKNOWN18, readWriteModel.UmasDataType_UNKNOWN19,
		readWriteModel.UmasDataType_UNKNOWN20, readWriteModel.UmasDataType_UNKNOWN24:
		unsigned, ok := asUnsigned(value, math.MaxUint8)
		if !ok {
			return nil, errors.Errorf("a %s doesn't fit into the single byte a %s takes", value.GetPlcValueType(), dataType)
		}
		return []byte{uint8(unsigned)}, nil

	case readWriteModel.UmasDataType_INT:
		signed, ok := asSigned(value, math.MinInt16, math.MaxInt16)
		if !ok {
			return nil, errors.Errorf("a %s doesn't fit into the two signed bytes an INT takes", value.GetPlcValueType())
		}
		return binary.LittleEndian.AppendUint16(nil, uint16(int16(signed))), nil

	case readWriteModel.UmasDataType_UINT, readWriteModel.UmasDataType_WORD:
		unsigned, ok := asUnsigned(value, math.MaxUint16)
		if !ok {
			return nil, errors.Errorf("a %s doesn't fit into the two unsigned bytes a %s takes", value.GetPlcValueType(), dataType)
		}
		return binary.LittleEndian.AppendUint16(nil, uint16(unsigned)), nil

	case readWriteModel.UmasDataType_DINT:
		signed, ok := asSigned(value, math.MinInt32, math.MaxInt32)
		if !ok {
			return nil, errors.Errorf("a %s doesn't fit into the four signed bytes a DINT takes", value.GetPlcValueType())
		}
		return binary.LittleEndian.AppendUint32(nil, uint32(int32(signed))), nil

	case readWriteModel.UmasDataType_UDINT, readWriteModel.UmasDataType_DWORD:
		unsigned, ok := asUnsigned(value, math.MaxUint32)
		if !ok {
			return nil, errors.Errorf("a %s doesn't fit into the four unsigned bytes a %s takes", value.GetPlcValueType(), dataType)
		}
		return binary.LittleEndian.AppendUint32(nil, uint32(unsigned)), nil

	case readWriteModel.UmasDataType_REAL:
		if !value.IsFloat32() {
			return nil, errors.Errorf("a %s is not a REAL", value.GetPlcValueType())
		}
		return binary.LittleEndian.AppendUint32(nil, math.Float32bits(value.GetFloat32())), nil

	case readWriteModel.UmasDataType_STRING:
		// Asking first: the plain value adapter panics in GetString rather than returning something
		// (spi/values/PlcValueAdapter.go), so a PlcRawByteArray or a PlcNull would take the whole write
		// request down with it instead of failing this one tag.
		if !value.IsString() {
			return nil, errors.Errorf("a %s is not a STRING", value.GetPlcValueType())
		}
		stringValue := value.GetString()
		// The PLC stores one byte per character. plc4j encodes with US_ASCII, which turns every
		// character outside that range into a '?' without saying so; refusing is the honest answer,
		// as there is no capture telling us what code page the device expects for the rest.
		for i := 0; i < len(stringValue); i++ {
			if stringValue[i] >= utf8.RuneSelf {
				return nil, errors.Errorf("Can't write %q as a UMAS STRING: only ASCII characters have a known encoding", stringValue)
			}
			if stringValue[i] == terminatorByte {
				return nil, errors.New("Can't write a UMAS STRING containing a NUL byte, it terminates the string")
			}
		}
		return append([]byte(stringValue), terminatorByte), nil

	case readWriteModel.UmasDataType_TIME:
		// Milliseconds as a uint32, not BCD. A plain number is taken as milliseconds too, which is
		// what plc4j does with its value.getLong().
		if value.IsDuration() {
			milliseconds := value.GetDuration().Milliseconds()
			if milliseconds < 0 || milliseconds > math.MaxUint32 {
				return nil, errors.Errorf("%s doesn't fit into the four bytes of milliseconds a TIME takes", value.GetDuration())
			}
			return binary.LittleEndian.AppendUint32(nil, uint32(milliseconds)), nil
		}
		milliseconds, ok := asUnsigned(value, math.MaxUint32)
		if !ok {
			return nil, errors.Errorf("a %s is neither a TIME nor a count of milliseconds that fits four bytes", value.GetPlcValueType())
		}
		return binary.LittleEndian.AppendUint32(nil, uint32(milliseconds)), nil

	case readWriteModel.UmasDataType_DATE:
		if !value.IsDate() && !value.IsDateTime() && !value.IsTime() {
			return nil, errors.Errorf("a %s is not a DATE", value.GetPlcValueType())
		}
		date := temporalValue(value)
		yearBytes, err := encodeBcdYear(date.Year())
		if err != nil {
			return nil, err
		}
		return []byte{
			encodeBcdByte(uint8(date.Day())),
			encodeBcdByte(uint8(int(date.Month()))),
			yearBytes[0],
			yearBytes[1],
		}, nil

	case readWriteModel.UmasDataType_TOD:
		if !value.IsTime() && !value.IsDateTime() && !value.IsDate() {
			return nil, errors.Errorf("a %s is not a TIME_OF_DAY", value.GetPlcValueType())
		}
		timeOfDay := temporalValue(value)
		centiseconds := (timeOfDay.Nanosecond() / int(10*time.Millisecond)) % 100
		return []byte{
			encodeBcdByte(uint8(centiseconds)),
			encodeBcdByte(uint8(timeOfDay.Second())),
			encodeBcdByte(uint8(timeOfDay.Minute())),
			encodeBcdByte(uint8(timeOfDay.Hour())),
		}, nil

	case readWriteModel.UmasDataType_DATE_AND_TIME:
		if !value.IsDateTime() && !value.IsDate() && !value.IsTime() {
			return nil, errors.Errorf("a %s is not a DATE_AND_TIME", value.GetPlcValueType())
		}
		dateTime := temporalValue(value)
		yearBytes, err := encodeBcdYear(dateTime.Year())
		if err != nil {
			return nil, err
		}
		return []byte{
			// The first byte is reserved and observed as zero in every capture plc4j was built from.
			0x00,
			encodeBcdByte(uint8(dateTime.Second())),
			encodeBcdByte(uint8(dateTime.Minute())),
			encodeBcdByte(uint8(dateTime.Hour())),
			encodeBcdByte(uint8(dateTime.Day())),
			encodeBcdByte(uint8(int(dateTime.Month()))),
			yearBytes[0],
			yearBytes[1],
		}, nil

	default:
		return nil, errors.Errorf("Can't serialize a value for the UMAS data type %s", dataType)
	}
}

// asBool narrows a value to a boolean, saying so when it isn't one.
//
// Asking is what keeps a write honest here: PlcValueAdapter.GetBool answers false rather than
// panicking, so every value which isn't a boolean and isn't a number - a PlcRawByteArray, a PlcNull,
// a PlcSTRING - would otherwise serialize a silent 0x00 to a PLC output and be reported as OK.
//
// A number counts as a boolean, "not zero" being true, because that is what plc4j's serializeForType
// gets out of value.getBoolean() for its numeric values and what the Go numeric values answer to
// GetBool as well - they just don't claim IsBool. The two integer narrowings are reused so the
// DWORD-ish types, whose Is<Integer> predicates all answer false, are covered too.
func asBool(value apiValues.PlcValue) (bool, bool) {
	if value.IsBool() {
		return value.GetBool(), true
	}
	if unsigned, ok := asUnsigned(value, math.MaxUint64); ok {
		return unsigned != 0, true
	}
	if signed, ok := asSigned(value, math.MinInt64, math.MaxInt64); ok {
		return signed != 0, true
	}
	return false, false
}

// asUnsigned narrows a value to an unsigned integer of at most maxValue, saying so when it doesn't
// fit. plc4j's serializeForType calls value.getInteger() and friends, which silently truncate; a
// write that lands a different number in the PLC than the caller asked for is worth refusing.
//
// The two DWORD-ish value types need asking by value type rather than by predicate: PlcDWORD and
// PlcLWORD are built on the plain value adapter instead of the numeric one, so every one of their
// Is<Integer> predicates answers false even though their getters return the number just fine.
func asUnsigned(value apiValues.PlcValue, maxValue uint64) (uint64, bool) {
	switch value.GetPlcValueType() {
	case apiValues.DWORD:
		unsigned := uint64(value.GetUint32())
		return unsigned, unsigned <= maxValue
	case apiValues.LWORD:
		unsigned := value.GetUint64()
		return unsigned, unsigned <= maxValue
	}
	if value.IsUint64() {
		unsigned := value.GetUint64()
		return unsigned, unsigned <= maxValue
	}
	if value.IsInt64() {
		signed := value.GetInt64()
		if signed < 0 {
			return 0, false
		}
		return uint64(signed), uint64(signed) <= maxValue
	}
	return 0, false
}

// asSigned narrows a value to a signed integer inside the given bounds, saying so when it doesn't
// fit. See asUnsigned for why the DWORD-ish types are asked by value type.
func asSigned(value apiValues.PlcValue, minValue int64, maxValue int64) (int64, bool) {
	switch value.GetPlcValueType() {
	case apiValues.DWORD:
		signed := int64(value.GetUint32())
		return signed, signed >= minValue && signed <= maxValue
	case apiValues.LWORD:
		unsigned := value.GetUint64()
		if unsigned > math.MaxInt64 {
			return 0, false
		}
		return int64(unsigned), int64(unsigned) >= minValue && int64(unsigned) <= maxValue
	}
	if value.IsInt64() {
		signed := value.GetInt64()
		return signed, signed >= minValue && signed <= maxValue
	}
	if value.IsUint64() {
		unsigned := value.GetUint64()
		if unsigned > math.MaxInt64 {
			return 0, false
		}
		return int64(unsigned), int64(unsigned) >= minValue && int64(unsigned) <= maxValue
	}
	return 0, false
}

// temporalValue picks whichever of the three time accessors the value actually answers to. The Go
// PlcValue types keep DATE, TIME_OF_DAY and DATE_AND_TIME apart, so a user may hand a write of a
// DATE symbol a PlcDATE_AND_TIME and the other way round; all three carry a time.Time.
func temporalValue(value apiValues.PlcValue) time.Time {
	switch {
	case value.IsDateTime():
		return value.GetDateTime()
	case value.IsDate():
		return value.GetDate()
	default:
		return value.GetTime()
	}
}

// mapToPlcValueType is what a symbol of this type looks like to plc4x. Ported from plc4j's
// UmasConnection.mapToPlcValueType.
func mapToPlcValueType(dataTypeId uint16) apiValues.PlcValueType {
	dataType, ok := resolveDataType(dataTypeId)
	if !ok {
		return apiValues.RAW_BYTE_ARRAY
	}
	switch dataType {
	case readWriteModel.UmasDataType_BOOL, readWriteModel.UmasDataType_EBOOL,
		readWriteModel.UmasDataType_UNKNOWN2, readWriteModel.UmasDataType_UNKNOWN3:
		return apiValues.BOOL
	case readWriteModel.UmasDataType_BYTE,
		readWriteModel.UmasDataType_UNKNOWN11, readWriteModel.UmasDataType_UNKNOWN12,
		readWriteModel.UmasDataType_UNKNOWN13, readWriteModel.UmasDataType_UNKNOWN17,
		readWriteModel.UmasDataType_UNKNOWN18, readWriteModel.UmasDataType_UNKNOWN19,
		readWriteModel.UmasDataType_UNKNOWN20, readWriteModel.UmasDataType_UNKNOWN24:
		return apiValues.BYTE
	case readWriteModel.UmasDataType_INT:
		return apiValues.INT
	case readWriteModel.UmasDataType_UINT:
		return apiValues.UINT
	case readWriteModel.UmasDataType_DINT:
		return apiValues.DINT
	case readWriteModel.UmasDataType_UDINT:
		return apiValues.UDINT
	case readWriteModel.UmasDataType_REAL:
		return apiValues.REAL
	case readWriteModel.UmasDataType_STRING:
		return apiValues.STRING
	case readWriteModel.UmasDataType_TIME:
		return apiValues.TIME
	case readWriteModel.UmasDataType_DATE:
		return apiValues.DATE
	case readWriteModel.UmasDataType_TOD:
		return apiValues.TIME_OF_DAY
	case readWriteModel.UmasDataType_DATE_AND_TIME:
		return apiValues.DATE_AND_TIME
	case readWriteModel.UmasDataType_WORD:
		return apiValues.WORD
	case readWriteModel.UmasDataType_DWORD:
		return apiValues.DWORD
	default:
		return apiValues.RAW_BYTE_ARRAY
	}
}

// decodeBcdByte turns one packed BCD byte into the number it spells: 0x25 is 25.
//
// A nibble above 9 is refused rather than turned into a number the way plc4j's decodeBcdByte does
// (it would read 0x1F as 25). A field which isn't BCD means the payload isn't what we think it is,
// and the caller can report that instead of a plausible looking wrong date.
func decodeBcdByte(bcd byte, what string) (uint8, error) {
	high := (bcd >> 4) & 0x0F
	low := bcd & 0x0F
	if high > 9 || low > 9 {
		return 0, errors.Errorf("0x%02X is not a BCD encoded %s", bcd, what)
	}
	return high*10 + low, nil
}

// decodeBcd16 turns the two packed BCD bytes of a year into the number they spell. The low byte
// carries the last two digits: 0x25 0x20 is 2025.
func decodeBcd16(low byte, high byte, what string) (uint16, error) {
	lowDigits, err := decodeBcdByte(low, what)
	if err != nil {
		return 0, err
	}
	highDigits, err := decodeBcdByte(high, what)
	if err != nil {
		return 0, err
	}
	return uint16(highDigits)*100 + uint16(lowDigits), nil
}

// encodeBcdByte packs a number between 0 and 99 into one BCD byte: 25 becomes 0x25. A caller which
// hands over more than 99 gets the low two digits, which is what plc4j's encodeBcd does; every call
// site here feeds it a calendar field which can't exceed 99.
func encodeBcdByte(value uint8) byte {
	value %= 100
	return byte(((value / 10) << 4) | (value % 10))
}

// encodeBcdYear packs a four digit year into the two BCD bytes a UMAS date carries, low digits
// first.
func encodeBcdYear(year int) ([2]byte, error) {
	if year < 0 || year > 9999 {
		return [2]byte{}, errors.Errorf("the year %d doesn't fit into the four BCD digits a UMAS date carries", year)
	}
	return [2]byte{encodeBcdByte(uint8(year % 100)), encodeBcdByte(uint8(year / 100))}, nil
}

// calendarDate builds a date out of BCD decoded fields, refusing the ones which aren't a date. Go's
// time.Date normalizes out of range fields (month 13 becomes January of the next year), which would
// turn a garbled payload into a perfectly plausible date.
func calendarDate(year uint16, month uint8, day uint8) (time.Time, error) {
	if month < 1 || month > 12 || day < 1 || day > 31 {
		return time.Time{}, errors.Errorf("%04d-%02d-%02d is not a date", year, month, day)
	}
	date := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC)
	if date.Year() != int(year) || date.Month() != time.Month(month) || date.Day() != int(day) {
		return time.Time{}, errors.Errorf("%04d-%02d-%02d is not a date", year, month, day)
	}
	return date, nil
}
