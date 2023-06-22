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

package ads

import (
	"context"
	"encoding/binary"

	"github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type MessageCodec struct {
	_default.DefaultCodec

	log zerolog.Logger
}

func NewMessageCodec(transportInstance transports.TransportInstance, _options ...options.WithOption) *MessageCodec {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	codec := &MessageCodec{
		log: customLogger,
	}
	codec.DefaultCodec = _default.NewDefaultCodec(
		codec,
		transportInstance,
		append(_options,
			_default.WithCustomMessageHandler(
				// This just prevents the loop from aborting in the start and by returning false,
				// it makes the message go to the default channel, as this means:
				// The handler hasn't handled the message
				func(codec _default.DefaultCodecRequirements, message spi.Message) bool {
					return false
				}),
		)...,
	)
	return codec
}

func (m *MessageCodec) GetCodec() spi.MessageCodec {
	return m
}

func (m *MessageCodec) Send(message spi.Message) error {
	m.log.Trace().Msg("Sending message")
	// Cast the message to the correct type of struct
	tcpPaket := message.(model.AmsTCPPacket)
	// Serialize the request
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	err := tcpPaket.SerializeWithWriteBuffer(context.Background(), wb)
	if err != nil {
		return errors.Wrap(err, "error serializing request")
	}

	// Send it to the PLC
	err = m.GetTransportInstance().Write(wb.GetBytes())
	if err != nil {
		return errors.Wrap(err, "error sending request")
	}
	return nil
}

func (m *MessageCodec) Receive() (spi.Message, error) {
	transportInstance := m.GetTransportInstance()

	if err := transportInstance.FillBuffer(
		func(pos uint, currentByte byte, reader transports.ExtendedReader) bool {
			numBytesAvailable, err := transportInstance.GetNumBytesAvailableInBuffer()
			if err != nil {
				return false
			}
			return numBytesAvailable < 6
		}); err != nil {
		m.log.Warn().Err(err).Msg("error filling buffer")
	}

	// We need at least 6 bytes in order to know how big the packet is in total
	if num, err := transportInstance.GetNumBytesAvailableInBuffer(); (err == nil) && (num >= 6) {
		m.log.Debug().Msgf("we got %d readable bytes", num)
		data, err := transportInstance.PeekReadableBytes(6)
		if err != nil {
			m.log.Warn().Err(err).Msg("error peeking")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		// Get the size of the entire packet little endian plus size of header
		packetSize := (uint32(data[5]) << 24) + (uint32(data[4]) << 16) + (uint32(data[3]) << 8) + (uint32(data[2])) + 6
		if num < packetSize {
			if err := transportInstance.FillBuffer(
				func(pos uint, currentByte byte, reader transports.ExtendedReader) bool {
					numBytesAvailable, err := transportInstance.GetNumBytesAvailableInBuffer()
					if err != nil {
						return false
					}
					return numBytesAvailable < packetSize
				}); err != nil {
				m.log.Warn().Err(err).Msg("error filling buffer")
			}
		}
		data, err = transportInstance.Read(packetSize)
		if err != nil {
			// TODO: Possibly clean up ...
			return nil, nil
		}
		rb := utils.NewReadBufferByteBased(data, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
		tcpPacket, err := model.AmsTCPPacketParseWithBuffer(context.Background(), rb)
		if err != nil {
			m.log.Warn().Err(err).Msg("error parsing")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		return tcpPacket, nil
	} else if err != nil {
		m.log.Warn().Err(err).Msg("Got error reading")
		return nil, nil
	}
	// TODO: maybe we return here a not enough error error
	return nil, nil
}
