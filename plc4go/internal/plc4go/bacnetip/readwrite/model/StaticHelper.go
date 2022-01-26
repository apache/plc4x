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
