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

func WriteCALReply(writeBuffer utils.WriteBuffer, calReply CALReply) error {
	return writeToHex("calReply", writeBuffer, calReply)
}

func ReadCALReply(readBuffer utils.ReadBuffer, payloadLength uint16, options CBusOptions, requestContext RequestContext) (CALReply, error) {
	rawBytes, err := readBytesFromHex("calReply", readBuffer, payloadLength)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting hex")
	}
	return CALReplyParse(utils.NewReadBufferByteBased(rawBytes), options, requestContext)
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

func WriteMonitoredSAL(writeBuffer utils.WriteBuffer, monitoredSAL MonitoredSAL) error {
	return writeToHex("monitoredSAL", writeBuffer, monitoredSAL)
}

func ReadMonitoredSAL(readBuffer utils.ReadBuffer, payloadLength uint16, options CBusOptions) (MonitoredSAL, error) {
	rawBytes, err := readBytesFromHex("monitoredSAL", readBuffer, payloadLength)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting hex")
	}
	return MonitoredSALParse(utils.NewReadBufferByteBased(rawBytes), options)
}

func WriteStandardFormatStatusReply(writeBuffer utils.WriteBuffer, reply StandardFormatStatusReply) error {
	return writeToHex("reply", writeBuffer, reply)
}

func ReadStandardFormatStatusReply(readBuffer utils.ReadBuffer, payloadLength uint16) (StandardFormatStatusReply, error) {
	rawBytes, err := readBytesFromHex("reply", readBuffer, payloadLength)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting hex")
	}
	return StandardFormatStatusReplyParse(utils.NewReadBufferByteBased(rawBytes))
}

func WriteExtendedFormatStatusReply(writeBuffer utils.WriteBuffer, reply ExtendedFormatStatusReply) error {
	return writeToHex("reply", writeBuffer, reply)
}

func ReadExtendedFormatStatusReply(readBuffer utils.ReadBuffer, payloadLength uint16) (ExtendedFormatStatusReply, error) {
	rawBytes, err := readBytesFromHex("reply", readBuffer, payloadLength)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting hex")
	}
	return ExtendedFormatStatusReplyParse(utils.NewReadBufferByteBased(rawBytes))
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
	hexBytes := make([]byte, hex.EncodedLen(len(wbbb.GetBytes())))
	// usually you use hex.Encode but we want the encoding in uppercase
	//n := hex.Encode(hexBytes, wbbb.GetBytes())
	n := encodeHexUpperCase(hexBytes, wbbb.GetBytes())
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
