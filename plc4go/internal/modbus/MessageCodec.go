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

package modbus

import (
	"github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type MessageCodec struct {
	_default.DefaultCodec
	expectationCounter int32

	log zerolog.Logger
}

func NewMessageCodec(transportInstance transports.TransportInstance, _options ...options.WithOption) *MessageCodec {
	codec := &MessageCodec{
		expectationCounter: 1,
		log:                options.ExtractCustomLogger(_options...),
	}
	codec.DefaultCodec = _default.NewDefaultCodec(codec, transportInstance)
	return codec
}

func (m *MessageCodec) GetCodec() spi.MessageCodec {
	return m
}

func (m *MessageCodec) Send(message spi.Message) error {
	m.log.Trace().Msg("Sending message")
	// Cast the message to the correct type of struct
	tcpAdu := message.(model.ModbusTcpADU)
	// Serialize the request
	theBytes, err := tcpAdu.Serialize()
	if err != nil {
		return errors.Wrap(err, "error serializing request")
	}

	// Send it to the PLC
	err = m.GetTransportInstance().Write(theBytes)
	if err != nil {
		return errors.Wrap(err, "error sending request")
	}
	return nil
}

func (m *MessageCodec) Receive() (spi.Message, error) {
	// We need at least 6 bytes in order to know how big the packet is in total
	if num, err := m.GetTransportInstance().GetNumBytesAvailableInBuffer(); (err == nil) && (num >= 6) {
		m.log.Debug().Msgf("we got %d readable bytes", num)
		data, err := m.GetTransportInstance().PeekReadableBytes(6)
		if err != nil {
			m.log.Warn().Err(err).Msg("error peeking")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		// Get the size of the entire packet
		packetSize := (uint32(data[4]) << 8) + uint32(data[5]) + 6
		if num < packetSize {
			m.log.Debug().Msgf("Not enough bytes. Got: %d Need: %d\n", num, packetSize)
			return nil, nil
		}
		data, err = m.GetTransportInstance().Read(packetSize)
		if err != nil {
			// TODO: Possibly clean up ...
			return nil, nil
		}
		tcpAdu, err := model.ModbusTcpADUParse(data, model.DriverType_MODBUS_TCP, true)
		if err != nil {
			m.log.Warn().Err(err).Msg("error parsing")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		return tcpAdu, nil
	} else if err != nil {
		m.log.Warn().Err(err).Msg("Got error reading")
		return nil, nil
	}
	// TODO: maybe we return here a not enough error error
	return nil, nil
}
