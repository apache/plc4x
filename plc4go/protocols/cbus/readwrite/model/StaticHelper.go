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
	// TODO: maybe we use a writebuffer hex based
	wbbb := utils.NewWriteBufferByteBased()
	err := cbusCommand.Serialize(wbbb)
	if err != nil {
		return errors.Wrap(err, "Error serializing")
	}
	hexBytes := make([]byte, hex.EncodedLen(len(wbbb.GetBytes())))
	n := hex.Encode(hexBytes, wbbb.GetBytes())
	log.Debug().Msgf("%d bytes encoded", n)
	return writeBuffer.WriteByteArray("cbusCommand", hexBytes)
}

func ReadCBusCommand(readBuffer utils.ReadBuffer, payloadLength uint16, srcchk bool) (CBusCommand, error) {
	hexBytes, err := readBuffer.ReadByteArray("cbusCommand", int(payloadLength))
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing")
	}
	rawBytes := make([]byte, hex.DecodedLen(len(hexBytes)))
	n, err := hex.Decode(rawBytes, hexBytes)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("%d bytes decoded", n)
	return CBusCommandParse(utils.NewReadBufferByteBased(rawBytes), srcchk)
}

func WriteCALReply(writeBuffer utils.WriteBuffer, calReply CALReply) error {
	// TODO: maybe we use a writebuffer hex based
	wbbb := utils.NewWriteBufferByteBased()
	err := calReply.Serialize(wbbb)
	if err != nil {
		return errors.Wrap(err, "Error serializing")
	}
	hexBytes := make([]byte, hex.EncodedLen(len(wbbb.GetBytes())))
	n := hex.Encode(hexBytes, wbbb.GetBytes())
	log.Debug().Msgf("%d bytes encoded", n)
	return writeBuffer.WriteByteArray("calReply", hexBytes)
}

func ReadCALReply(readBuffer utils.ReadBuffer, payloadLength uint16) (CALReply, error) {
	hexBytes, err := readBuffer.ReadByteArray("calReply", int(payloadLength))
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing")
	}
	rawBytes := make([]byte, hex.DecodedLen(len(hexBytes)))
	n, err := hex.Decode(rawBytes, hexBytes)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("%d bytes decoded", n)
	return CALReplyParse(utils.NewReadBufferByteBased(rawBytes))
}
