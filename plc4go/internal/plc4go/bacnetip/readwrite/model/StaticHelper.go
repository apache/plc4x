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
	// TODO: check if it's in the known range and if not return (==VENDOR_PROPRIETARY_VALUE)
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
	if baCnetPropertyIdentifier != 0 && baCnetPropertyIdentifier != BACnetPropertyIdentifier_VENDOR_PROPRIETARY_VALUE {
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
	if value != 0 && value != BACnetPropertyIdentifier_VENDOR_PROPRIETARY_VALUE {
		return 0, nil
	}
	// We need to reset our reader to the position we read before
	readBuffer.SetPos(readBuffer.GetPos() - uint16(actualLength))
	bitsToRead := (uint8)(actualLength * 8)
	return readBuffer.ReadUint32("proprietaryPropertyIdentifier", bitsToRead)
}

func OpeningClosingTerminate(readBuffer utils.ReadBuffer, openingTag *BACnetContextTag) bool {
	if openingTag == nil {
		// If we don't have an opening tag at all we can terminate here
		return true
	}
	oldPos := readBuffer.GetPos()
	aByte, _ := readBuffer.ReadByte("")
	readBuffer.SetPos(oldPos)
	return aByte == 0x3F
}

func ParseTags(readBuffer utils.ReadBuffer) *BACnetTag {
	tag, err := BACnetTagParse(readBuffer)
	if err != nil {
		panic(err)
	}
	return tag
}

func WriteTags(writeBuffer utils.WriteBuffer, value *BACnetTag) error {
	return value.Serialize(writeBuffer)
}

func TagsLength(tags []*BACnetTag) uint16 {
	var length uint16
	for _, tag := range tags {
		length += tag.LengthInBytes()
	}
	return length
}
