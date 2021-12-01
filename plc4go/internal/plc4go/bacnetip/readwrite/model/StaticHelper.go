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
		readValue, err = readBuffer.ReadUint8("", 8)
		readUnsignedLong = uint32(readValue)
	} else if bitsToRead <= 16 {
		var readValue uint16
		readValue, err = readBuffer.ReadUint16("", 16)
		readUnsignedLong = uint32(readValue)
	} else if bitsToRead <= 32 {
		var readValue uint32
		readValue, err = readBuffer.ReadUint32("", 32)
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
	return writeBuffer.WriteUint32("", bitsToWrite, uint32(value), utils.WithAdditionalStringRepresentation(value.name()))
}
