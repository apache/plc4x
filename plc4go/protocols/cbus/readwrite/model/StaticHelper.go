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
	"github.com/apache/plc4x/plc4go/internal/spi"
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func ReadAndValidateChecksum(readBuffer utils.ReadBuffer, message spi.Message, srchk bool) (Checksum, error) {
	if !srchk {
		return nil, nil
	}
	hexBytes, err := readBytesFromHex("chksum", readBuffer, 2, false)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to calculate checksum")
	}
	checksum := hexBytes[0]
	actualChecksum, err := getChecksum(message)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to calculate checksum")
	}
	if checksum != actualChecksum {
		return nil, errors.Errorf("Expected checksum 0x%x doesn't match actual checksum 0x%x", checksum, actualChecksum)
	}
	return NewChecksum(checksum), nil
}

func CalculateChecksum(writeBuffer utils.WriteBuffer, message spi.Message, srchk bool) error {
	if !srchk {
		// Nothing to do when srchck is disabled
		return nil
	}
	checksum, err := getChecksum(message)
	if err != nil {
		return errors.Wrap(err, "Unable to calculate checksum")
	}
	return writeToHex("chksum", writeBuffer, []byte{checksum})
}

func getChecksum(message spi.Message) (byte, error) {
	checksum := byte(0x0)
	checksumWriteBuffer := utils.NewWriteBufferByteBased()
	err := message.Serialize(checksumWriteBuffer)
	if err != nil {
		return 0, errors.Wrap(err, "Error serializing")
	}
	for _, aByte := range checksumWriteBuffer.GetBytes() {
		checksum += aByte
	}
	checksum = ^checksum
	checksum++
	return checksum, nil
}

func WriteCBusCommand(writeBuffer utils.WriteBuffer, cbusCommand CBusCommand) error {
	return writeSerializableToHex("cbusCommand", writeBuffer, cbusCommand)
}

func ReadCBusCommand(readBuffer utils.ReadBuffer, payloadLength uint16, cBusOptions CBusOptions, srchk bool) (CBusCommand, error) {
	rawBytes, err := readBytesFromHex("cbusCommand", readBuffer, payloadLength, srchk)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting hex")
	}
	return CBusCommandParse(utils.NewReadBufferByteBased(rawBytes), cBusOptions)
}

func WriteEncodedReply(writeBuffer utils.WriteBuffer, encodedReply EncodedReply) error {
	return writeSerializableToHex("encodedReply", writeBuffer, encodedReply)
}

func ReadEncodedReply(readBuffer utils.ReadBuffer, payloadLength uint16, options CBusOptions, requestContext RequestContext, srchk bool) (EncodedReply, error) {
	rawBytes, err := readBytesFromHex("encodedReply", readBuffer, payloadLength, srchk)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting hex")
	}
	return EncodedReplyParse(utils.NewReadBufferByteBased(rawBytes), options, requestContext)
}

func WriteCALData(writeBuffer utils.WriteBuffer, calData CALData) error {
	return writeSerializableToHex("calData", writeBuffer, calData)
}

func ReadCALData(readBuffer utils.ReadBuffer, payloadLength uint16) (CALData, error) {
	rawBytes, err := readBytesFromHex("calData", readBuffer, payloadLength, false)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting hex")
	}
	return CALDataParse(utils.NewReadBufferByteBased(rawBytes), nil)
}

func readBytesFromHex(logicalName string, readBuffer utils.ReadBuffer, payloadLength uint16, srchk bool) ([]byte, error) {
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
	if srchk {
		checksum := byte(0x0)
		for _, aByte := range rawBytes {
			checksum += aByte
		}
		if checksum != 0x0 {
			return nil, errors.New("Checksum validation failed")
		}
		// We need to reset the last to hex bytes
		readBuffer.Reset(readBuffer.GetPos() - 2)
		rawBytes = rawBytes[:len(rawBytes)-1]
	}
	log.Debug().Msgf("%d bytes decoded", n)
	return rawBytes, nil
}

func writeSerializableToHex(logicalName string, writeBuffer utils.WriteBuffer, serializable utils.Serializable) error {
	wbbb := utils.NewWriteBufferByteBased()
	err := serializable.Serialize(wbbb)
	if err != nil {
		return errors.Wrap(err, "Error serializing")
	}
	bytesToWrite := wbbb.GetBytes()
	return writeToHex(logicalName, writeBuffer, bytesToWrite)
}

func writeToHex(logicalName string, writeBuffer utils.WriteBuffer, bytesToWrite []byte) error {
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

func KnowsMeteringCommandTypeContainer(readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return MeteringCommandTypeContainerKnows(readUint8)
}

func KnowsTriggerControlCommandTypeContainer(readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return TriggerControlCommandTypeContainerKnows(readUint8)
}

func KnowsEnableControlCommandTypeContainer(readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return EnableControlCommandTypeContainerKnows(readUint8)
}

func KnowsTemperatureBroadcastCommandTypeContainer(readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return TemperatureBroadcastCommandTypeContainerKnows(readUint8)
}

func KnowsAccessControlCommandTypeContainer(readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return AccessControlCommandTypeContainerKnows(readUint8)
}

func KnowsMediaTransportControlCommandTypeContainer(readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return MediaTransportControlCommandTypeContainerKnows(readUint8)
}

func KnowsClockAndTimekeepingCommandTypeContainer(readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return ClockAndTimekeepingCommandTypeContainerKnows(readUint8)
}

func KnowsTelephonyCommandTypeContainer(readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return TelephonyCommandTypeContainerKnows(readUint8)
}

func KnowsAirConditioningCommandTypeContainer(readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return AirConditioningCommandTypeContainerKnows(readUint8)
}

func KnowsMeasurementCommandTypeContainer(readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return MeasurementCommandTypeContainerKnows(readUint8)
}

func KnowsErrorReportingCommandTypeContainer(readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return ErrorReportingCommandTypeContainerKnows(readUint8)
}
