/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package model

import (
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/pkg/errors"
	"math/big"
)

func ReadPropertyIdentifier(readBuffer utils.ReadBuffer, actualLength uint32) (BACnetPropertyIdentifier, error) {
	bitsToRead := actualLength * 8
	var readUnsignedLong uint32
	var err error
	if bitsToRead <= 8 {
		var readValue uint8
		readValue, err = readBuffer.ReadUint8("propertyIdentifier", 8)
		readUnsignedLong = uint32(readValue)
	} else if bitsToRead <= 16 {
		var readValue uint16
		readValue, err = readBuffer.ReadUint16("propertyIdentifier", 16)
		readUnsignedLong = uint32(readValue)
	} else if bitsToRead <= 32 {
		var readValue uint32
		readValue, err = readBuffer.ReadUint32("propertyIdentifier", 32)
		readUnsignedLong = uint32(readValue)
	} else {
		return 0, errors.Errorf("%d overflows", bitsToRead)
	}
	if err != nil {
		return 0, err
	}

	return BACnetPropertyIdentifier(readUnsignedLong), nil
}

func WritePropertyIdentifier(writeBuffer utils.WriteBuffer, value BACnetPropertyIdentifier) error {
	if value == BACnetPropertyIdentifier_VENDOR_PROPRIETARY_VALUE {
		return nil
	}
	var bitsToWrite uint8
	valueValue := uint64(value)
	if valueValue <= 0xff {
		bitsToWrite = 8
	} else if valueValue <= 0xffff {
		bitsToWrite = 16
	} else if valueValue <= 0xffffffff {
		bitsToWrite = 32
	} else {
		bitsToWrite = 32
	}
	return writeBuffer.WriteUint32("propertyIdentifier", bitsToWrite, uint32(value), utils.WithAdditionalStringRepresentation(value.name()))
}

func WriteProprietaryPropertyIdentifier(writeBuffer utils.WriteBuffer, baCnetPropertyIdentifier BACnetPropertyIdentifier, value uint32) error {
	if baCnetPropertyIdentifier != BACnetPropertyIdentifier_VENDOR_PROPRIETARY_VALUE {
		return nil
	}
	var bitsToWrite uint8
	if value <= 0xff {
		bitsToWrite = 8
	} else if value <= 0xffff {
		bitsToWrite = 16
	} else if value <= 0xffffffff {
		bitsToWrite = 32
	} else {
		bitsToWrite = 32
	}
	return writeBuffer.WriteUint32("proprietaryPropertyIdentifier", bitsToWrite, value, utils.WithAdditionalStringRepresentation(BACnetPropertyIdentifier_VENDOR_PROPRIETARY_VALUE.name()))
}

func ReadProprietaryPropertyIdentifier(readBuffer utils.ReadBuffer, value BACnetPropertyIdentifier, actualLength uint32) (uint32, error) {
	if value != BACnetPropertyIdentifier_VENDOR_PROPRIETARY_VALUE {
		return 0, nil
	}
	// We need to reset our reader to the position we read before
	readBuffer.Reset(readBuffer.GetPos() - uint16(actualLength))
	bitsToRead := (uint8)(actualLength * 8)
	return readBuffer.ReadUint32("proprietaryPropertyIdentifier", bitsToRead)
}

func ReadEventState(readBuffer utils.ReadBuffer, actualLength uint32) (BACnetEventState, error) {
	bitsToRead := actualLength * 8
	var readUnsignedLong uint32
	var err error
	if bitsToRead <= 8 {
		var readValue uint8
		readValue, err = readBuffer.ReadUint8("eventState", 8)
		readUnsignedLong = uint32(readValue)
	} else if bitsToRead <= 16 {
		var readValue uint16
		readValue, err = readBuffer.ReadUint16("eventState", 16)
		readUnsignedLong = uint32(readValue)
	} else if bitsToRead <= 32 {
		var readValue uint32
		readValue, err = readBuffer.ReadUint32("eventState", 32)
		readUnsignedLong = uint32(readValue)
	} else {
		return 0, errors.Errorf("%d overflows", bitsToRead)
	}
	if err != nil {
		return 0, err
	}

	return BACnetEventState(readUnsignedLong), nil
}

func WriteEventState(writeBuffer utils.WriteBuffer, value BACnetEventState) error {
	if value == BACnetEventState_VENDOR_PROPRIETARY_VALUE {
		return nil
	}
	var bitsToWrite uint8
	valueValue := uint64(value)
	if valueValue <= 0xff {
		bitsToWrite = 8
	} else if valueValue <= 0xffff {
		bitsToWrite = 16
	} else if valueValue <= 0xffffffff {
		bitsToWrite = 32
	} else {
		bitsToWrite = 32
	}
	return writeBuffer.WriteUint32("eventState", bitsToWrite, uint32(value), utils.WithAdditionalStringRepresentation(value.name()))
}

func WriteProprietaryEventState(writeBuffer utils.WriteBuffer, baCnetEventState BACnetEventState, value uint32) error {
	if baCnetEventState != BACnetEventState_VENDOR_PROPRIETARY_VALUE {
		return nil
	}
	var bitsToWrite uint8
	if value <= 0xff {
		bitsToWrite = 8
	} else if value <= 0xffff {
		bitsToWrite = 16
	} else if value <= 0xffffffff {
		bitsToWrite = 32
	} else {
		bitsToWrite = 32
	}
	return writeBuffer.WriteUint32("proprietaryEventState", bitsToWrite, value, utils.WithAdditionalStringRepresentation(BACnetEventState_VENDOR_PROPRIETARY_VALUE.name()))
}

func ReadProprietaryEventState(readBuffer utils.ReadBuffer, value BACnetEventState, actualLength uint32) (uint32, error) {
	if value != BACnetEventState_VENDOR_PROPRIETARY_VALUE {
		return 0, nil
	}
	// We need to reset our reader to the position we read before
	readBuffer.Reset(readBuffer.GetPos() - uint16(actualLength))
	bitsToRead := (uint8)(actualLength * 8)
	return readBuffer.ReadUint32("proprietaryEventState", bitsToRead)
}

func ReadEventType(readBuffer utils.ReadBuffer, actualLength uint32) (BACnetEventType, error) {
	bitsToRead := actualLength * 8
	var readUnsignedLong uint32
	var err error
	if bitsToRead <= 8 {
		var readValue uint8
		readValue, err = readBuffer.ReadUint8("eventType", 8)
		readUnsignedLong = uint32(readValue)
	} else if bitsToRead <= 16 {
		var readValue uint16
		readValue, err = readBuffer.ReadUint16("eventType", 16)
		readUnsignedLong = uint32(readValue)
	} else if bitsToRead <= 32 {
		var readValue uint32
		readValue, err = readBuffer.ReadUint32("eventType", 32)
		readUnsignedLong = uint32(readValue)
	} else {
		return 0, errors.Errorf("%d overflows", bitsToRead)
	}
	if err != nil {
		return 0, err
	}

	return BACnetEventType(readUnsignedLong), nil
}

func WriteEventType(writeBuffer utils.WriteBuffer, value BACnetEventType) error {
	if value == BACnetEventType_VENDOR_PROPRIETARY_VALUE {
		return nil
	}
	var bitsToWrite uint8
	valueValue := uint64(value)
	if valueValue <= 0xff {
		bitsToWrite = 8
	} else if valueValue <= 0xffff {
		bitsToWrite = 16
	} else if valueValue <= 0xffffffff {
		bitsToWrite = 32
	} else {
		bitsToWrite = 32
	}
	return writeBuffer.WriteUint32("eventType", bitsToWrite, uint32(value), utils.WithAdditionalStringRepresentation(value.name()))
}

func WriteProprietaryEventType(writeBuffer utils.WriteBuffer, baCnetEventType BACnetEventType, value uint32) error {
	if baCnetEventType != BACnetEventType_VENDOR_PROPRIETARY_VALUE {
		return nil
	}
	var bitsToWrite uint8
	if value <= 0xff {
		bitsToWrite = 8
	} else if value <= 0xffff {
		bitsToWrite = 16
	} else if value <= 0xffffffff {
		bitsToWrite = 32
	} else {
		bitsToWrite = 32
	}
	return writeBuffer.WriteUint32("proprietaryEventType", bitsToWrite, value, utils.WithAdditionalStringRepresentation(BACnetEventType_VENDOR_PROPRIETARY_VALUE.name()))
}

func ReadProprietaryEventType(readBuffer utils.ReadBuffer, value BACnetEventType, actualLength uint32) (uint32, error) {
	if value != BACnetEventType_VENDOR_PROPRIETARY_VALUE {
		return 0, nil
	}
	// We need to reset our reader to the position we read before
	readBuffer.Reset(readBuffer.GetPos() - uint16(actualLength))
	bitsToRead := (uint8)(actualLength * 8)
	return readBuffer.ReadUint32("proprietaryEventType", bitsToRead)
}
func ReadObjectType(readBuffer utils.ReadBuffer) (BACnetObjectType, error) {
	readValue, err := readBuffer.ReadUint16("objectType", 10)
	if err != nil {
		return 0, err
	}
	return BACnetObjectType(readValue), nil
}

func WriteObjectType(writeBuffer utils.WriteBuffer, value BACnetObjectType) error {
	if value == BACnetObjectType_VENDOR_PROPRIETARY_VALUE {
		return nil
	}
	return writeBuffer.WriteUint16("objectType", 10, uint16(value), utils.WithAdditionalStringRepresentation(value.name()))
}

func WriteProprietaryObjectType(writeBuffer utils.WriteBuffer, baCnetObjectType BACnetObjectType, value uint16) error {
	if baCnetObjectType != BACnetObjectType_VENDOR_PROPRIETARY_VALUE {
		return nil
	}
	return writeBuffer.WriteUint16("proprietaryObjectType", 10, value, utils.WithAdditionalStringRepresentation(BACnetObjectType_VENDOR_PROPRIETARY_VALUE.name()))
}

func ReadProprietaryObjectType(readBuffer utils.ReadBuffer, value BACnetObjectType) (uint16, error) {
	if value != BACnetObjectType_VENDOR_PROPRIETARY_VALUE {
		return 0, nil
	}
	// We need to reset our reader to the position we read before
	readBuffer.Reset(readBuffer.GetPos() - 2)
	return readBuffer.ReadUint16("proprietaryObjectType", 10)
}

func IsBACnetConstructedDataClosingTag(readBuffer utils.ReadBuffer, instantTerminate bool, expectedTagNumber byte) bool {
	if instantTerminate {
		return true
	}
	oldPos := readBuffer.GetPos()
	// TODO: add graceful exit if we know already that we are at the end (we might need to add available bytes to reader)
	tagNumber, err := readBuffer.ReadUint8("", 4)
	if err != nil {
		return true
	}
	isContextTag, err := readBuffer.ReadBit("")
	if err != nil {
		return true
	}
	tagValue, err := readBuffer.ReadUint8("", 3)
	if err != nil {
		return true
	}

	foundOurClosingTag := isContextTag && tagNumber == expectedTagNumber && tagValue == 0x7
	readBuffer.Reset(oldPos)
	return foundOurClosingTag
}

func GuessDataType(objectType BACnetObjectType, propertyIdentifier *BACnetContextTagPropertyIdentifier) BACnetDataType {
	// TODO: implement me
	return BACnetDataType_ENUMERATED
}

func ParseVarUint(data []byte) uint32 {
	bigInt := big.NewInt(0)
	return uint32(bigInt.SetBytes(data).Uint64())
}

func WriteVarUint(value uint32) []byte {
	return big.NewInt(int64(value)).Bytes()
}

func CreateBACnetTagHeaderBalanced(isContext bool, id uint8, value uint32) *BACnetTagHeader {
	tagClass := TagClass_APPLICATION_TAGS
	if isContext {
		tagClass = TagClass_CONTEXT_SPECIFIC_TAGS
	}

	var tagNumber uint8
	var extTagNumber *uint8
	if id <= 14 {
		tagNumber = id
	} else {
		tagNumber = 0xF
		extTagNumber = &id
	}

	var lengthValueType uint8
	var extLength *uint8
	var extExtLength *uint16
	var extExtExtLength *uint32
	if value <= 4 {
		lengthValueType = uint8(value)
	} else {
		lengthValueType = 5
		// Depending on the length, we will either write it as an 8 bit, 32 bit, or 64 bit integer
		if value <= 253 {
			_extLength := uint8(value)
			extLength = &_extLength
		} else if value <= 65535 {
			_extLength := uint8(254)
			extLength = &_extLength
			_extExtLength := uint16(value)
			extExtLength = &_extExtLength
		} else {
			_extLength := uint8(255)
			extLength = &_extLength
			extExtExtLength = &value
		}
	}

	return NewBACnetTagHeader(tagNumber, tagClass, lengthValueType, extTagNumber, extLength, extExtLength, extExtExtLength)
}

func CreateBACnetApplicationTagNull() *BACnetApplicationTagNull {
	header := CreateBACnetTagHeaderBalanced(false, uint8(BACnetDataType_NULL), 0)
	return NewBACnetApplicationTagNull(header)
}

func CreateBACnetContextTagNull(tagNumber uint8) *BACnetContextTagNull {
	header := CreateBACnetTagHeaderBalanced(true, tagNumber, 0)
	return NewBACnetContextTagNull(header, tagNumber, true)
}

func CreateBACnetOpeningTag(tagNum uint8) *BACnetOpeningTag {
	var tagNumber uint8
	var extTagNumber *uint8
	if tagNum <= 14 {
		tagNumber = tagNum
	} else {
		tagNumber = 0xF
		extTagNumber = &tagNum
	}
	header := NewBACnetTagHeader(tagNumber, TagClass_CONTEXT_SPECIFIC_TAGS, 0x6, extTagNumber, nil, nil, nil)
	return NewBACnetOpeningTag(header, tagNum, 6)
}

func CreateBACnetClosingTag(tagNum uint8) *BACnetClosingTag {
	var tagNumber uint8
	var extTagNumber *uint8
	if tagNum <= 14 {
		tagNumber = tagNum
	} else {
		tagNumber = 0xF
		extTagNumber = &tagNum
	}
	header := NewBACnetTagHeader(tagNumber, TagClass_CONTEXT_SPECIFIC_TAGS, 0x7, extTagNumber, nil, nil, nil)
	return NewBACnetClosingTag(header, tagNum, 7)
}

func CreateBACnetApplicationTagObjectIdentifier(objectType uint16, instance uint32) *BACnetApplicationTagObjectIdentifier {
	header := NewBACnetTagHeader(uint8(BACnetDataType_BACNET_OBJECT_IDENTIFIER), TagClass_APPLICATION_TAGS, uint8(4), nil, nil, nil, nil)
	objectTypeEnum := BACnetObjectTypeByValue(objectType)
	proprietaryValue := uint16(0)
	if objectType >= 128 || !BACnetObjectTypeKnows(objectType) {
		objectTypeEnum = BACnetObjectType_VENDOR_PROPRIETARY_VALUE
		proprietaryValue = objectType
	}
	payload := NewBACnetTagPayloadObjectIdentifier(objectTypeEnum, proprietaryValue, instance)
	return NewBACnetApplicationTagObjectIdentifier(payload, header)
}

func CreateBACnetContextTagObjectIdentifier(tagNum uint8, objectType uint16, instance uint32) *BACnetContextTagObjectIdentifier {
	header := NewBACnetTagHeader(tagNum, TagClass_CONTEXT_SPECIFIC_TAGS, uint8(4), nil, nil, nil, nil)
	objectTypeEnum := BACnetObjectTypeByValue(objectType)
	proprietaryValue := uint16(0)
	if objectType >= 128 {
		objectTypeEnum = BACnetObjectType_VENDOR_PROPRIETARY_VALUE
		proprietaryValue = objectType
	}
	payload := NewBACnetTagPayloadObjectIdentifier(objectTypeEnum, proprietaryValue, instance)
	return NewBACnetContextTagObjectIdentifier(payload, header, tagNum, true)
}

func CreateBACnetContextTagPropertyIdentifier(tagNum uint8, propertyType uint32) *BACnetContextTagPropertyIdentifier {
	header := NewBACnetTagHeader(tagNum, TagClass_CONTEXT_SPECIFIC_TAGS, uint8(requiredLength(uint(propertyType))), nil, nil, nil, nil)
	propertyTypeEnum := BACnetPropertyIdentifierByValue(propertyType)
	proprietaryValue := uint32(0)
	if !BACnetPropertyIdentifierKnows(propertyType) {
		propertyTypeEnum = BACnetPropertyIdentifier_VENDOR_PROPRIETARY_VALUE
		proprietaryValue = propertyType
	}
	return NewBACnetContextTagPropertyIdentifier(propertyTypeEnum, proprietaryValue, header, tagNum, true, 0)
}

func CreateBACnetApplicationTagUnsignedInteger(value uint) *BACnetApplicationTagUnsignedInteger {
	length, payload := CreateUnsignedPayload(value)
	header := CreateBACnetTagHeaderBalanced(false, uint8(BACnetDataType_UNSIGNED_INTEGER), length)
	return NewBACnetApplicationTagUnsignedInteger(payload, header)
}

func CreateBACnetContextTagUnsignedInteger(tagNumber uint8, value uint) *BACnetContextTagUnsignedInteger {
	length, payload := CreateUnsignedPayload(value)
	header := CreateBACnetTagHeaderBalanced(true, tagNumber, length)
	return NewBACnetContextTagUnsignedInteger(payload, header, tagNumber, true)
}

func CreateUnsignedPayload(value uint) (uint32, *BACnetTagPayloadUnsignedInteger) {
	var length uint32
	var valueUint8 *uint8
	var valueUint16 *uint16
	var valueUint24 *uint32
	var valueUint32 *uint32
	var valueUint40 *uint64
	var valueUint48 *uint64
	var valueUint56 *uint64
	var valueUint64 *uint64
	switch {
	case value < 0x100:
		length = 1
		_valueUint8 := uint8(value)
		valueUint8 = &_valueUint8
	case value < 0x10000:
		length = 2
		_valueUint16 := uint16(value)
		valueUint16 = &_valueUint16
	case value < 0x1000000:
		length = 3
		_valueUint24 := uint32(value)
		valueUint24 = &_valueUint24
		//TODO: support more than 32bit
	default:
		length = 4
		valueUint32_ := uint32(value)
		valueUint32 = &valueUint32_
	}
	payload := NewBACnetTagPayloadUnsignedInteger(valueUint8, valueUint16, valueUint24, valueUint32, valueUint40, valueUint48, valueUint56, valueUint64, length)
	return length, payload
}

func CreateBACnetApplicationTagSignedInteger(value int) *BACnetApplicationTagSignedInteger {
	length, payload := CreateSignedPayload(value)
	header := CreateBACnetTagHeaderBalanced(true, uint8(BACnetDataType_SIGNED_INTEGER), length)
	return NewBACnetApplicationTagSignedInteger(payload, header)
}

func CreateBACnetContextTagSignedInteger(tagNumber uint8, value int) *BACnetContextTagSignedInteger {
	length, payload := CreateSignedPayload(value)
	header := CreateBACnetTagHeaderBalanced(true, tagNumber, length)
	return NewBACnetContextTagSignedInteger(payload, header, tagNumber, true)
}

func CreateSignedPayload(value int) (uint32, *BACnetTagPayloadSignedInteger) {
	var length uint32
	var valueInt8 *int8
	var valueInt16 *int16
	var valueInt24 *int32
	var valueInt32 *int32
	switch {
	case value < 0x100:
		length = 1
		_valueInt8 := int8(value)
		valueInt8 = &_valueInt8
	case value < 0x10000:
		length = 2
		_valueInt16 := int16(value)
		valueInt16 = &_valueInt16
	case value < 0x1000000:
		length = 3
		_valueInt24 := int32(value)
		valueInt24 = &_valueInt24
		//TODO: support more than 32bit
	default:
		length = 4
		_valueInt32 := int32(value)
		valueInt32 = &_valueInt32
	}
	payload := NewBACnetTagPayloadSignedInteger(valueInt8, valueInt16, valueInt24, valueInt32, nil, nil, nil, nil, length)
	return length, payload
}

func CreateBACnetApplicationTagBoolean(value bool) *BACnetApplicationTagBoolean {
	_value := uint32(0)
	if value {
		_value = 1
	}
	header := CreateBACnetTagHeaderBalanced(false, uint8(BACnetDataType_BOOLEAN), _value)
	return NewBACnetApplicationTagBoolean(NewBACnetTagPayloadBoolean(_value), header)
}

func CreateBACnetContextTagBoolean(tagNumber uint8, value bool) *BACnetContextTagBoolean {
	header := CreateBACnetTagHeaderBalanced(true, tagNumber, 1)
	_value := uint8(0)
	if value {
		_value = 1
	}
	return NewBACnetContextTagBoolean(_value, NewBACnetTagPayloadBoolean(uint32(_value)), header, tagNumber, true)
}

func CreateBACnetApplicationTagReal(value float32) *BACnetApplicationTagReal {
	header := CreateBACnetTagHeaderBalanced(false, uint8(BACnetDataType_REAL), 4)
	return NewBACnetApplicationTagReal(NewBACnetTagPayloadReal(value), header)
}

func CreateBACnetContextTagReal(tagNumber uint8, value float32) *BACnetContextTagReal {
	header := CreateBACnetTagHeaderBalanced(true, tagNumber, 4)
	return NewBACnetContextTagReal(NewBACnetTagPayloadReal(value), header, tagNumber, true)
}

func CreateBACnetApplicationTagDouble(value float64) *BACnetApplicationTagDouble {
	header := CreateBACnetTagHeaderBalanced(false, uint8(BACnetDataType_DOUBLE), 8)
	return NewBACnetApplicationTagDouble(NewBACnetTagPayloadDouble(value), header)
}

func CreateBACnetContextTagDouble(tagNumber uint8, value float64) *BACnetContextTagDouble {
	header := CreateBACnetTagHeaderBalanced(true, tagNumber, 8)
	return NewBACnetContextTagDouble(NewBACnetTagPayloadDouble(value), header, tagNumber, true)
}

func CreateBACnetApplicationTagOctetString(value string) *BACnetApplicationTagOctetString {
	header := CreateBACnetTagHeaderBalanced(false, uint8(BACnetDataType_OCTET_STRING), uint32(len(value)+1))
	return NewBACnetApplicationTagOctetString(NewBACnetTagPayloadOctetString(value, uint32(len(value)+1)), header)
}

func CreateBACnetContextTagOctetString(tagNumber uint8, value string) *BACnetContextTagOctetString {
	header := CreateBACnetTagHeaderBalanced(true, tagNumber, uint32(len(value)+1))
	return NewBACnetContextTagOctetString(NewBACnetTagPayloadOctetString(value, uint32(len(value)+1)), header, tagNumber, true)
}

func CreateBACnetApplicationTagCharacterString(baCnetCharacterEncoding BACnetCharacterEncoding, value string) *BACnetApplicationTagCharacterString {
	header := CreateBACnetTagHeaderBalanced(false, uint8(BACnetDataType_CHARACTER_STRING), uint32(len(value)+1))
	return NewBACnetApplicationTagCharacterString(NewBACnetTagPayloadCharacterString(baCnetCharacterEncoding, value, uint32(len(value)+1)), header)
}

func CreateBACnetContextTagCharacterString(tagNumber uint8, baCnetCharacterEncoding BACnetCharacterEncoding, value string) *BACnetContextTagCharacterString {
	header := CreateBACnetTagHeaderBalanced(true, tagNumber, uint32(len(value)+1))
	return NewBACnetContextTagCharacterString(NewBACnetTagPayloadCharacterString(baCnetCharacterEncoding, value, uint32(len(value)+1)), header, tagNumber, true)
}

func CreateBACnetApplicationTagBitString(value []bool) *BACnetApplicationTagBitString {
	numberOfBytesNeeded := (len(value) + 7) / 8
	unusedBits := 8 - (len(value) % 8)
	if unusedBits == 8 {
		unusedBits = 0
	}
	header := CreateBACnetTagHeaderBalanced(false, uint8(BACnetDataType_BIT_STRING), uint32(numberOfBytesNeeded+1))
	return NewBACnetApplicationTagBitString(NewBACnetTagPayloadBitString(uint8(unusedBits), value, make([]bool, unusedBits), uint32(numberOfBytesNeeded+1)), header)
}

func CreateBACnetContextTagBitString(tagNumber uint8, value []bool) *BACnetContextTagBitString {
	numberOfBytesNeeded := (len(value) + 7) / 8
	unusedBits := 8 - (len(value) % 8)
	if unusedBits == 8 {
		unusedBits = 0
	}
	header := CreateBACnetTagHeaderBalanced(true, tagNumber, uint32(numberOfBytesNeeded+1))
	return NewBACnetContextTagBitString(NewBACnetTagPayloadBitString(uint8(unusedBits), value, make([]bool, unusedBits), uint32(numberOfBytesNeeded+1)), header, tagNumber, true)
}

func CreateBACnetApplicationTagEnumerated(value uint32) *BACnetApplicationTagEnumerated {
	length, payload := CreateEnumeratedPayload(value)
	header := CreateBACnetTagHeaderBalanced(false, uint8(BACnetDataType_ENUMERATED), length)
	result := NewBACnetApplicationTagEnumerated(payload, header)
	return CastBACnetApplicationTagEnumerated(result)
}

func CreateBACnetContextTagEnumerated(tagNumber uint8, value uint32) *BACnetContextTagEnumerated {
	length, payload := CreateEnumeratedPayload(value)
	header := CreateBACnetTagHeaderBalanced(true, tagNumber, length)
	result := NewBACnetContextTagEnumerated(payload, header, tagNumber, true)
	return CastBACnetContextTagEnumerated(result)
}

func CreateEnumeratedPayload(value uint32) (uint32, *BACnetTagPayloadEnumerated) {
	length := requiredLength(uint(value))
	data := WriteVarUint(value)
	payload := NewBACnetTagPayloadEnumerated(data, length)
	return length, payload
}

func CreateBACnetApplicationTagDate(year uint16, month, dayOfMonth, dayOfWeek uint8) *BACnetApplicationTagDate {
	header := CreateBACnetTagHeaderBalanced(false, uint8(BACnetDataType_DATE), 4)
	yearMinus1900 := uint8(year - 1900)
	if year == 0xFF {
		yearMinus1900 = 0xFF
	}
	return NewBACnetApplicationTagDate(NewBACnetTagPayloadDate(yearMinus1900, month, dayOfMonth, dayOfWeek), header)
}

func CreateBACnetContextTagDate(tagNumber uint8, year uint16, month, dayOfMonth, dayOfWeek uint8) *BACnetContextTagDate {
	header := CreateBACnetTagHeaderBalanced(true, tagNumber, 4)
	yearMinus1900 := uint8(year - 1900)
	if year == 0xFF {
		yearMinus1900 = 0xFF
	}
	return NewBACnetContextTagDate(NewBACnetTagPayloadDate(yearMinus1900, month, dayOfMonth, dayOfWeek), header, tagNumber, true)
}

func CreateBACnetApplicationTagTime(hour, minute, second, fractional uint8) *BACnetApplicationTagTime {
	header := CreateBACnetTagHeaderBalanced(false, uint8(BACnetDataType_TIME), 4)
	return NewBACnetApplicationTagTime(NewBACnetTagPayloadTime(hour, minute, second, fractional), header)
}

func CreateBACnetContextTagTime(tagNumber uint8, hour, minute, second, fractional uint8) *BACnetContextTagTime {
	header := CreateBACnetTagHeaderBalanced(true, tagNumber, 4)
	return NewBACnetContextTagTime(NewBACnetTagPayloadTime(hour, minute, second, fractional), header, tagNumber, true)
}

func DummyPropertyIdentifier() *BACnetContextTagPropertyIdentifier {
	return NewBACnetContextTagPropertyIdentifier(BACnetPropertyIdentifier_VENDOR_PROPRIETARY_VALUE, 0, nil, 0, true, 0)
}

func requiredLength(value uint) uint32 {
	var length uint32
	switch {
	case value < 0x100:
		length = 1
	case value < 0x10000:
		length = 2
	case value < 0x1000000:
		length = 3
	default:
		length = 4
	}
	return length
}
