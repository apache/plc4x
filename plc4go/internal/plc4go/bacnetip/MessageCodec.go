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

package bacnetip

import (
	"github.com/apache/plc4x/plc4go/internal/plc4go/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type MessageCodec struct {
	_default.DefaultCodec
}

func NewMessageCodec(transportInstance transports.TransportInstance) *MessageCodec {
	codec := &MessageCodec{}
	codec.DefaultCodec = _default.NewDefaultCodec(codec, transportInstance)
	return codec
}

func (m *MessageCodec) GetCodec() spi.MessageCodec {
	return m
}

func (m *MessageCodec) Send(message interface{}) error {
	log.Trace().Msg("Sending message")
	// Cast the message to the correct type of struct
	bvlcPacket := model.CastBVLC(message)
	// Serialize the request
	wb := utils.NewWriteBufferByteBased()
	err := bvlcPacket.Serialize(wb)
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

func (m *MessageCodec) Receive() (interface{}, error) {
	log.Trace().Msg("receiving")
	// We need at least 6 bytes in order to know how big the packet is in total
	if num, err := m.GetTransportInstance().GetNumReadableBytes(); (err == nil) && (num >= 4) {
		log.Debug().Msgf("we got %d readable bytes", num)
		data, err := m.GetTransportInstance().PeekReadableBytes(4)
		if err != nil {
			log.Warn().Err(err).Msg("error peeking")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		//Second byte for the size and then add the header size 24
		packetSize := uint32((uint16(data[3]) << 8) + uint16(data[2]))
		if num < packetSize {
			log.Debug().Msgf("Not enough bytes. Got: %d Need: %d\n", num, packetSize)
			return nil, nil
		}
		data, err = m.GetTransportInstance().Read(packetSize)
		if err != nil {
			log.Debug().Err(err).Msg("Error reading")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		rb := utils.NewReadBufferByteBased(data)
		bvlcPacket, err := model.BVLCParse(rb)
		if err != nil {
			log.Warn().Err(err).Msg("error parsing")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		return bvlcPacket, nil
	} else if err != nil {
		log.Warn().Err(err).Msg("Got error reading")
		return nil, nil
	}
	// TODO: maybe we return here a not enough error error
	return nil, nil
}
