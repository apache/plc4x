//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package knxnetip

import (
	"github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type MessageCodec struct {
	_default.DefaultCodec
	sequenceCounter    int32
	messageInterceptor func(message interface{})
}

func NewMessageCodec(transportInstance transports.TransportInstance, messageInterceptor func(message interface{})) *MessageCodec {
	codec := &MessageCodec{
		messageInterceptor: messageInterceptor,
	}
	codec.DefaultCodec = _default.NewDefaultCodec(
		codec,
		transportInstance,
		_default.WithCustomMessageHandler(CustomMessageHandling),
	)
	return codec
}

func (m *MessageCodec) GetCodec() spi.MessageCodec {
	return m
}

func (m *MessageCodec) Send(message interface{}) error {
	log.Trace().Msg("Sending message")
	// Cast the message to the correct type of struct
	knxMessage := model.CastKnxNetIpMessage(message)
	// Serialize the request
	wb := utils.NewWriteBufferByteBased()
	err := knxMessage.Serialize(wb)
	if err != nil {
		return errors.Wrap(err, "error serializing request")
	}

	// Send it to the PLC
	err = m.GetTransportInstance().Write(wb.GetBytes())
	if err != nil {
		return errors.Wrap(err, "error sending request ")
	}
	return nil
}

func (m *MessageCodec) Receive() (interface{}, error) {
	log.Trace().Msg("receiving")
	// We need at least 6 bytes in order to know how big the packet is in total
	if num, err := m.GetTransportInstance().GetNumReadableBytes(); (err == nil) && (num >= 6) {
		log.Debug().Msgf("we got %d readable bytes", num)
		data, err := m.GetTransportInstance().PeekReadableBytes(6)
		if err != nil {
			log.Warn().Err(err).Msg("error peeking")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		// Get the size of the entire packet
		packetSize := (uint32(data[4]) << 8) + uint32(data[5])
		if num < packetSize {
			log.Trace().Msgf("Not enough bytes. Got: %d Need: %d\n", num, packetSize)
			return nil, nil
		}
		data, err = m.GetTransportInstance().Read(packetSize)
		if err != nil {
			log.Warn().Err(err).Msg("error reading")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		rb := utils.NewReadBufferByteBased(data)
		knxMessage, err := model.KnxNetIpMessageParse(rb)
		if err != nil {
			log.Warn().Err(err).Msg("error parsing message")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		return knxMessage, nil
	} else if err != nil {
		log.Warn().Err(err).Msg("Got error reading")
		return nil, nil
	}
	return nil, nil
}

func CustomMessageHandling(codec *_default.DefaultCodecRequirements, message interface{}) bool {
	// If this message is a simple KNXNet/IP UDP Ack, ignore it for now
	tunnelingResponse := model.CastTunnelingResponse(message)
	if tunnelingResponse != nil {
		return true
	}

	// If this is an incoming tunneling request, automatically send a tunneling ACK back to the gateway
	tunnelingRequest := model.CastTunnelingRequest(message)
	if tunnelingRequest != nil {
		response := model.NewTunnelingResponse(
			model.NewTunnelingResponseDataBlock(
				tunnelingRequest.TunnelingRequestDataBlock.CommunicationChannelId,
				tunnelingRequest.TunnelingRequestDataBlock.SequenceCounter,
				model.Status_NO_ERROR),
		)
		err := (*codec).Send(response)
		if err != nil {
			log.Warn().Err(err).Msg("got an error sending ACK from transport")
		}
	}

	localCodec := (*codec).(*MessageCodec)
	// Handle the packet itself
	// Give a message interceptor a chance to intercept
	if (*localCodec).messageInterceptor != nil {
		(*localCodec).messageInterceptor(message)
	}
	return false
}
