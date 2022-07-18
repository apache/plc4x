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

package model

import (
	"encoding/hex"
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func WriteCBusCommand(writeBuffer utils.WriteBuffer, cbusCommand CBusCommand) error {
	return writeToHex("cbusCommand", writeBuffer, cbusCommand)
}

func ReadCBusCommand(readBuffer utils.ReadBuffer, payloadLength uint16, cBusOptions CBusOptions) (CBusCommand, error) {
	rawBytes, err := readBytesFromHex("cbusCommand", readBuffer, payloadLength)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting hex")
	}
	return CBusCommandParse(utils.NewReadBufferByteBased(rawBytes), cBusOptions)
}

func WriteEncodedReply(writeBuffer utils.WriteBuffer, encodedReply EncodedReply) error {
	return writeToHex("encodedReply", writeBuffer, encodedReply)
}

func ReadEncodedReply(readBuffer utils.ReadBuffer, payloadLength uint16, options CBusOptions, requestContext RequestContext) (EncodedReply, error) {
	rawBytes, err := readBytesFromHex("encodedReply", readBuffer, payloadLength)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting hex")
	}
	return EncodedReplyParse(utils.NewReadBufferByteBased(rawBytes), options, requestContext)
}

func WriteCALDataOrSetParameter(writeBuffer utils.WriteBuffer, calDataOrSetParameter CALDataOrSetParameter) error {
	return writeToHex("calDataOrSetParameter", writeBuffer, calDataOrSetParameter)
}

func ReadCALDataOrSetParameter(readBuffer utils.ReadBuffer, payloadLength uint16) (CALDataOrSetParameter, error) {
	rawBytes, err := readBytesFromHex("calDataOrSetParameter", readBuffer, payloadLength)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting hex")
	}
	return CALDataOrSetParameterParse(utils.NewReadBufferByteBased(rawBytes))
}

func readBytesFromHex(logicalName string, readBuffer utils.ReadBuffer, payloadLength uint16) ([]byte, error) {
	if payloadLength == 0 {
		return nil, errors.New("Length is 0")
	}
	hexBytes, err := readBuffer.ReadByteArray(logicalName, int(payloadLength))
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing")
	}
	lastByte := hexBytes[len(hexBytes)-1]
	if (lastByte >= 0x67) && (lastByte <= 0x7A) {
		// We need to reset the alpha
		readBuffer.Reset(readBuffer.GetPos() - 1)
		hexBytes = hexBytes[:len(hexBytes)-1]
	}
	rawBytes := make([]byte, hex.DecodedLen(len(hexBytes)))
	n, err := hex.Decode(rawBytes, hexBytes)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("%d bytes decoded", n)
	return rawBytes, nil
}

func writeToHex(logicalName string, writeBuffer utils.WriteBuffer, serializable utils.Serializable) error {
	wbbb := utils.NewWriteBufferByteBased()
	err := serializable.Serialize(wbbb)
	if err != nil {
		return errors.Wrap(err, "Error serializing")
	}
	bytesToWrite := wbbb.GetBytes()
	hexBytes := make([]byte, hex.EncodedLen(len(bytesToWrite)))
	// usually you use hex.Encode but we want the encoding in uppercase
	//n := hex.Encode(hexBytes, wbbb.GetBytes())
	n := encodeHexUpperCase(hexBytes, bytesToWrite)
	log.Debug().Msgf("%d bytes encoded", n)
	return writeBuffer.WriteByteArray(logicalName, hexBytes)
}

const hextable = "0123456789ABCDEF"

func encodeHexUpperCase(dst, src []byte) int {
	j := 0
	for _, v := range src {
		dst[j] = hextable[v>>4]
		dst[j+1] = hextable[v&0x0f]
		j += 2
	}
	return len(src) * 2
}

func KnowsCALCommandTypeContainer(readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return CALCommandTypeContainerKnows(readUint8)
}

func KnowsLightingCommandTypeContainer(readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return LightingCommandTypeContainerKnows(readUint8)
}

func KnowsSecurityCommandTypeContainer(readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return SecurityCommandTypeContainerKnows(readUint8)
}
