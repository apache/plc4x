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

package eip

import (
	"context"
	"encoding/binary"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

//go:generate plc4xGenerator -type=MessageCodec
type MessageCodec struct {
	_default.DefaultCodec
	none bool // TODO: just a empty field to satisfy generator (needs fixing because in this case here we have the delegate)

	log zerolog.Logger
}

func NewMessageCodec(transportInstance transports.TransportInstance, _options ...options.WithOption) *MessageCodec {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	codec := &MessageCodec{
		log: customLogger,
	}
	codec.DefaultCodec = _default.NewDefaultCodec(codec, transportInstance, _options...)
	return codec
}

func (m *MessageCodec) GetCodec() spi.MessageCodec {
	return m
}

func (m *MessageCodec) Send(message spi.Message) error {
	m.log.Trace().Msg("Sending message")
	// Cast the message to the correct type of struct
	eipPacket := message.(model.EipPacket)
	// Serialize the request
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	err := eipPacket.SerializeWithWriteBuffer(context.Background(), wb)
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
	// We need at least 6 bytes in order to know how big the packet is in total
	transportInstance := m.GetTransportInstance()
	if num, err := transportInstance.GetNumBytesAvailableInBuffer(); (err == nil) && (num >= 4) {
		m.log.Debug().Uint32("num", num).Msg("we got num readable bytes")
		data, err := transportInstance.PeekReadableBytes(4)
		if err != nil {
			m.log.Warn().Err(err).Msg("error peeking")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		//Second byte for the size and then add the header size 24
		packetSize := uint32(((uint16(data[3]) << 8) + uint16(data[2])) + 24)
		if num < packetSize {
			m.log.Debug().Uint32("num", num).Uint32("packetSize", packetSize).Msg("Not enough bytes. Got: num Need: packetSize")
			return nil, nil
		}
		data, err = transportInstance.Read(packetSize)
		if err != nil {
			m.log.Debug().Err(err).Msg("Error reading")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		rb := utils.NewReadBufferByteBased(data, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
		eipPacket, err := model.EipPacketParseWithBuffer[model.EipPacket](context.Background(), rb, true)
		if err != nil {
			m.log.Warn().Err(err).Msg("error parsing")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		return eipPacket, nil
	} else if err != nil {
		m.log.Warn().Err(err).Msg("Got error reading")
		return nil, nil
	}
	// TODO: maybe we return here a not enough error error
	return nil, nil
}
