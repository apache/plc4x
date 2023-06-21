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

package knxnetip

import (
	"context"

	"github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type MessageCodec struct {
	_default.DefaultCodec
	sequenceCounter    int32
	messageInterceptor func(message spi.Message)

	passLogToModel bool
	log            zerolog.Logger
}

func NewMessageCodec(transportInstance transports.TransportInstance, messageInterceptor func(message spi.Message), _options ...options.WithOption) *MessageCodec {
	passLoggerToModel, _ := options.ExtractPassLoggerToModel(_options...)
	customLogger, _ := options.ExtractCustomLogger(_options...)
	codec := &MessageCodec{
		messageInterceptor: messageInterceptor,
		passLogToModel:     passLoggerToModel,
		log:                customLogger,
	}
	codec.DefaultCodec = _default.NewDefaultCodec(
		codec,
		transportInstance,
		append(_options, _default.WithCustomMessageHandler(CustomMessageHandling(customLogger)))...,
	)
	return codec
}

func (m *MessageCodec) GetCodec() spi.MessageCodec {
	return m
}

func (m *MessageCodec) Send(message spi.Message) error {
	m.log.Trace().Msg("Sending message")
	// Cast the message to the correct type of struct
	knxMessage := message.(model.KnxNetIpMessage)
	// Serialize the request
	theBytes, err := knxMessage.Serialize()
	if err != nil {
		return errors.Wrap(err, "error serializing request")
	}

	// Send it to the PLC
	err = m.GetTransportInstance().Write(theBytes)
	if err != nil {
		return errors.Wrap(err, "error sending request ")
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
		packetSize := (uint32(data[4]) << 8) + uint32(data[5])
		if num < packetSize {
			m.log.Trace().Msgf("Not enough bytes. Got: %d Need: %d\n", num, packetSize)
			return nil, nil
		}
		data, err = m.GetTransportInstance().Read(packetSize)
		if err != nil {
			m.log.Warn().Err(err).Msg("error reading")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		ctxForModel := options.GetLoggerContextForModel(context.TODO(), m.log, options.WithPassLoggerToModel(m.passLogToModel))
		knxMessage, err := model.KnxNetIpMessageParse(ctxForModel, data)
		if err != nil {
			m.log.Warn().Err(err).Msg("error parsing message")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		return knxMessage, nil
	} else if err != nil {
		m.log.Warn().Err(err).Msg("Got error reading")
		return nil, nil
	}
	return nil, nil
}

func CustomMessageHandling(localLog zerolog.Logger) _default.CustomMessageHandler {
	return func(codec _default.DefaultCodecRequirements, message spi.Message) bool {
		// If this message is a simple KNXNet/IP UDP Ack, ignore it for now
		tunnelingResponse := message.(model.TunnelingResponse)
		if tunnelingResponse != nil {
			return true
		}

		// If this is an incoming tunneling request, automatically send a tunneling ACK back to the gateway
		tunnelingRequest := message.(model.TunnelingRequest)
		if tunnelingRequest != nil {
			response := model.NewTunnelingResponse(
				model.NewTunnelingResponseDataBlock(
					tunnelingRequest.GetTunnelingRequestDataBlock().GetCommunicationChannelId(),
					tunnelingRequest.GetTunnelingRequestDataBlock().GetSequenceCounter(),
					model.Status_NO_ERROR),
			)
			err := codec.Send(response)
			if err != nil {
				localLog.Warn().Err(err).Msg("got an error sending ACK from transport")
			}
		}

		localCodec := codec.(*MessageCodec)
		// Handle the packet itself
		// Give a message interceptor a chance to intercept
		if (*localCodec).messageInterceptor != nil {
			(*localCodec).messageInterceptor(message)
		}
		return false
	}
}
