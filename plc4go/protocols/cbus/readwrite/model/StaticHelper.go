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
	"context"
	"encoding/hex"

	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/pkg/errors"
)

func ReadAndValidateChecksum(ctx context.Context, readBuffer utils.ReadBuffer, message spi.Message, srchk bool) (Checksum, error) {
	if !srchk {
		return nil, nil
	}
	hexBytes, err := readBytesFromHex(ctx, "chksum", readBuffer, false)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to calculate checksum")
	}
	checksum := hexBytes[0]
	actualChecksum, err := getChecksum(message)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to calculate checksum")
	}
	if checksum != actualChecksum {
		return nil, errors.Errorf("Expected checksum %#x doesn't match actual checksum %#x", checksum, actualChecksum)
	}
	return NewChecksum(checksum), nil
}

func CalculateChecksum(ctx context.Context, writeBuffer utils.WriteBuffer, message spi.Message, srchk bool) error {
	if !srchk {
		// Nothing to do when srchck is disabled
		return nil
	}
	checksum, err := getChecksum(message)
	if err != nil {
		return errors.Wrap(err, "Unable to calculate checksum")
	}
	return writeToHex(ctx, "chksum", writeBuffer, []byte{checksum})
}

func getChecksum(message spi.Message) (byte, error) {
	checksum := byte(0x0)
	theBytes, err := message.Serialize()
	if err != nil {
		return 0, errors.Wrap(err, "Error serializing")
	}
	for _, aByte := range theBytes {
		checksum += aByte
	}
	checksum = ^checksum
	checksum++
	return checksum, nil
}

func WriteCBusCommand(ctx context.Context, writeBuffer utils.WriteBuffer, cbusCommand CBusCommand) error {
	return writeSerializableToHex(ctx, "cbusCommand", writeBuffer, cbusCommand)
}

func ReadCBusCommand(ctx context.Context, readBuffer utils.ReadBuffer, cBusOptions CBusOptions, srchk bool) (CBusCommand, error) {
	rawBytes, err := readBytesFromHex(ctx, "cbusCommand", readBuffer, srchk)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting hex")
	}
	return CBusCommandParse(ctx, rawBytes, cBusOptions)
}

func WriteEncodedReply(ctx context.Context, writeBuffer utils.WriteBuffer, encodedReply EncodedReply) error {
	return writeSerializableToHex(ctx, "encodedReply", writeBuffer, encodedReply)
}

func ReadEncodedReply(ctx context.Context, readBuffer utils.ReadBuffer, options CBusOptions, requestContext RequestContext, srchk bool) (EncodedReply, error) {
	rawBytes, err := readBytesFromHex(ctx, "encodedReply", readBuffer, srchk)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting hex")
	}
	return EncodedReplyParse(ctx, rawBytes, options, requestContext)
}

func WriteCALData(ctx context.Context, writeBuffer utils.WriteBuffer, calData CALData) error {
	return writeSerializableToHex(ctx, "calData", writeBuffer, calData)
}

func ReadCALData(ctx context.Context, readBuffer utils.ReadBuffer) (CALData, error) {
	rawBytes, err := readBytesFromHex(ctx, "calData", readBuffer, false)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting hex")
	}
	return CALDataParse(context.TODO(), rawBytes, nil)
}

func readBytesFromHex(ctx context.Context, logicalName string, readBuffer utils.ReadBuffer, srchk bool) ([]byte, error) {
	payloadLength := findHexEnd(ctx, readBuffer)
	if payloadLength == 0 {
		return nil, utils.ParseAssertError{Message: "Length is 0"}
	}
	hexBytes, err := readBuffer.ReadByteArray(logicalName, payloadLength)
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
	Plc4xModelLog.Trace().Msgf("%d bytes decoded", n)
	return rawBytes, nil
}

func findHexEnd(ctx context.Context, readBuffer utils.ReadBuffer) int {
	// TODO: find out if there is a smarter way to find the end...
	oldPos := readBuffer.GetPos()
	payloadLength := 0
	for readBuffer.HasMore(8) {
		hexByte, _ := readBuffer.ReadByte("")
		isHex := hexByte >= 'A' && hexByte <= 'F' || hexByte >= 'a' && hexByte <= 'f'
		isNumber := hexByte >= '0' && hexByte <= '9'
		if !isHex && !isNumber {
			break
		}
		payloadLength++
	}
	readBuffer.Reset(oldPos)
	return payloadLength
}

func writeSerializableToHex(ctx context.Context, logicalName string, writeBuffer utils.WriteBuffer, serializable utils.Serializable) error {
	theBytes, err := serializable.Serialize()
	if err != nil {
		return errors.Wrap(err, "Error serializing")
	}
	return writeToHex(ctx, logicalName, writeBuffer, theBytes)
}

func writeToHex(ctx context.Context, logicalName string, writeBuffer utils.WriteBuffer, bytesToWrite []byte) error {
	hexBytes := make([]byte, hex.EncodedLen(len(bytesToWrite)))
	// usually you use hex.Encode but we want the encoding in uppercase
	//n := hex.Encode(hexBytes, wbbb.GetBytes())
	n := encodeHexUpperCase(hexBytes, bytesToWrite)
	Plc4xModelLog.Trace().Msgf("%d bytes encoded", n)
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

func KnowsCALCommandTypeContainer(ctx context.Context, readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return CALCommandTypeContainerKnows(readUint8)
}

func KnowsLightingCommandTypeContainer(ctx context.Context, readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return LightingCommandTypeContainerKnows(readUint8)
}

func KnowsSecurityCommandTypeContainer(ctx context.Context, readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return SecurityCommandTypeContainerKnows(readUint8)
}

func KnowsMeteringCommandTypeContainer(ctx context.Context, readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return MeteringCommandTypeContainerKnows(readUint8)
}

func KnowsTriggerControlCommandTypeContainer(ctx context.Context, readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return TriggerControlCommandTypeContainerKnows(readUint8)
}

func KnowsEnableControlCommandTypeContainer(ctx context.Context, readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return EnableControlCommandTypeContainerKnows(readUint8)
}

func KnowsTemperatureBroadcastCommandTypeContainer(ctx context.Context, readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return TemperatureBroadcastCommandTypeContainerKnows(readUint8)
}

func KnowsAccessControlCommandTypeContainer(ctx context.Context, readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return AccessControlCommandTypeContainerKnows(readUint8)
}

func KnowsMediaTransportControlCommandTypeContainer(ctx context.Context, readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return MediaTransportControlCommandTypeContainerKnows(readUint8)
}

func KnowsClockAndTimekeepingCommandTypeContainer(ctx context.Context, readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return ClockAndTimekeepingCommandTypeContainerKnows(readUint8)
}

func KnowsTelephonyCommandTypeContainer(ctx context.Context, readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return TelephonyCommandTypeContainerKnows(readUint8)
}

func KnowsAirConditioningCommandTypeContainer(ctx context.Context, readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return AirConditioningCommandTypeContainerKnows(readUint8)
}

func KnowsMeasurementCommandTypeContainer(ctx context.Context, readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return MeasurementCommandTypeContainerKnows(readUint8)
}

func KnowsErrorReportingCommandTypeContainer(ctx context.Context, readBuffer utils.ReadBuffer) bool {
	oldPos := readBuffer.GetPos()
	defer readBuffer.Reset(oldPos)
	readUint8, err := readBuffer.ReadUint8("", 8)
	if err != nil {
		return false
	}
	return ErrorReportingCommandTypeContainerKnows(readUint8)
}
