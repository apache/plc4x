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

func GuessDataType(objectType BACnetObjectType) BACnetDataType {
	// TODO: implement me
	return BACnetDataType_BACNET_PROPERTY_IDENTIFIER
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

func CreateBACnetApplicationTagObjectIdentifier(objectType uint16, instance uint32) *BACnetApplicationTagObjectIdentifier {
	header := NewBACnetTagHeader(0xC, TagClass_APPLICATION_TAGS, 4, nil, nil, nil, nil)
	objectTypeEnum := BACnetObjectTypeByValue(objectType)
	if objectType >= 128 || !BACnetObjectTypeKnows(objectType) {
		objectTypeEnum = BACnetObjectType_VENDOR_PROPRIETARY_VALUE
	}
	payload := NewBACnetTagPayloadObjectIdentifier(objectTypeEnum, objectType, instance)
	result := NewBACnetApplicationTagObjectIdentifier(payload, header)
	return CastBACnetApplicationTagObjectIdentifier(result)
}

func CreateBACnetContextTagObjectIdentifier(tagNum uint8, objectType uint16, instance uint32) *BACnetContextTagObjectIdentifier {
	header := NewBACnetTagHeader(tagNum, TagClass_CONTEXT_SPECIFIC_TAGS, 4, nil, nil, nil, nil)
	objectTypeEnum := BACnetObjectTypeByValue(objectType)
	if objectType >= 128 {
		objectTypeEnum = BACnetObjectType_VENDOR_PROPRIETARY_VALUE
	}
	payload := NewBACnetTagPayloadObjectIdentifier(objectTypeEnum, objectType, instance)
	result := NewBACnetContextTagObjectIdentifier(payload, header, tagNum, true)
	return CastBACnetContextTagObjectIdentifier(result)
}

func CreateBACnetContextTagPropertyIdentifier(tagNum uint8, propertyType uint32) *BACnetContextTagPropertyIdentifier {
	header := NewBACnetTagHeader(tagNum, TagClass_CONTEXT_SPECIFIC_TAGS, 4, nil, nil, nil, nil)
	propertyTypeEnum := BACnetPropertyIdentifierByValue(propertyType)
	if !BACnetPropertyIdentifierKnows(propertyType) {
		propertyTypeEnum = BACnetPropertyIdentifier_VENDOR_PROPRIETARY_VALUE
	}
	result := NewBACnetContextTagPropertyIdentifier(propertyTypeEnum, propertyType, header, tagNum, true, 0)
	return CastBACnetContextTagPropertyIdentifier(result)
}

func CreateBACnetApplicationTagEnumerated(value uint32) *BACnetApplicationTagEnumerated {
	length, payload := CreateEnumeratedPayload(value)
	header := CreateBACnetTagHeaderBalanced(false, 0x9, length)
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
	data := WriteVarUint(value)
	payload := NewBACnetTagPayloadEnumerated(data, length)
	return length, payload
}

func CreateBACnetApplicationTagUnsignedInteger(value uint32) *BACnetApplicationTagUnsignedInteger {
	length, payload := CreateUnsignedPayload(value)
	header := CreateBACnetTagHeaderBalanced(false, 0x2, length)
	result := NewBACnetApplicationTagUnsignedInteger(payload, header)
	return CastBACnetApplicationTagUnsignedInteger(result)
}

func CreateBACnetContextTagUnsignedInteger(tagNumber uint8, value uint32) *BACnetContextTagUnsignedInteger {
	length, payload := CreateUnsignedPayload(value)
	header := CreateBACnetTagHeaderBalanced(true, tagNumber, length)
	result := NewBACnetContextTagUnsignedInteger(payload, header, tagNumber, true)
	return CastBACnetContextTagUnsignedInteger(result)
}

func CreateUnsignedPayload(value uint32) (uint32, *BACnetTagPayloadUnsignedInteger) {
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
	default:
		length = 4
		valueUint32 = &value
	}
	payload := NewBACnetTagPayloadUnsignedInteger(valueUint8, valueUint16, valueUint24, valueUint32, valueUint40, valueUint48, valueUint56, valueUint64, length)
	return length, payload
}

func CreateBACnetApplicationTagSignedInteger(value int32) *BACnetApplicationTagSignedInteger {
	length, payload := CreateSignedPayload(value)
	header := CreateBACnetTagHeaderBalanced(true, 0x3, length)
	result := NewBACnetApplicationTagSignedInteger(payload, header)
	return CastBACnetApplicationTagSignedInteger(result)
}

func CreateBACnetContextTagSignedInteger(tagNumber uint8, value int32) *BACnetContextTagSignedInteger {
	length, payload := CreateSignedPayload(value)
	header := CreateBACnetTagHeaderBalanced(true, tagNumber, length)
	result := NewBACnetContextTagSignedInteger(payload, header, tagNumber, true)
	return CastBACnetContextTagSignedInteger(result)
}

func CreateSignedPayload(value int32) (uint32, *BACnetTagPayloadSignedInteger) {
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
	default:
		length = 4
		valueInt32 = &value
	}
	payload := NewBACnetTagPayloadSignedInteger(valueInt8, valueInt16, valueInt24, valueInt32, nil, nil, nil, nil, length)
	return length, payload
}
