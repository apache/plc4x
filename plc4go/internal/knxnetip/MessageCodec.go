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

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

//go:generate go tool plc4xGenerator -type=MessageCodec
type MessageCodec struct {
	_default.DefaultCodec

	sequenceCounter    int32
	messageInterceptor func(message spi.Message)

	passLogToModel bool
	log            zerolog.Logger
}

var (
	_ spi.TransportInstanceExposer = (*MessageCodec)(nil)
)

func NewMessageCodec(transportInstance transports.TransportInstance, messageInterceptor func(message spi.Message), _options ...options.WithOption) *MessageCodec {
	passLoggerToModel, _ := options.ExtractPassLoggerToModel(_options...)
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
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

func (m *MessageCodec) Send(ctx context.Context, interactionInfo string, message spi.Message) error {
	m.log.Trace().Str("interactionInfo", interactionInfo).Msg("Sending message")
	// Cast the message to the correct type of struct
	knxMessage := message.(model.KnxNetIpMessage)
	// Serialize the request
	theBytes, err := knxMessage.Serialize()
	if err != nil {
		return errors.Wrap(err, "error serializing request")
	}

	// Send it to the PLC
	err = m.GetTransportInstance().Write(ctx, theBytes)
	if err != nil {
		return errors.Wrap(err, "error sending request ")
	}
	return nil
}

func (m *MessageCodec) Receive(ctx context.Context) (spi.Message, error) {
	// We need at least 6 bytes in order to know how big the packet is in total
	if num, err := m.GetTransportInstance().GetNumBytesAvailableInBuffer(); (err == nil) && (num >= 6) {
		m.log.Debug().Uint32("num", num).Msg("we got num readable bytes")
		data, err := m.GetTransportInstance().PeekReadableBytes(ctx, 6)
		if err != nil {
			m.log.Warn().Err(err).Msg("error peeking")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		// Get the size of the entire packet
		packetSize := (uint32(data[4]) << 8) + uint32(data[5])
		// A total length below the 6 byte KNXnet/IP header can never become
		// parseable and would make the receive worker spin forever on the same
		// unconsumed bytes, so treat it as a fatal framing error.
		if packetSize < 6 {
			return nil, transports.NewTransportError(transports.TransportErrorFatal,
				errors.Errorf("invalid KNXnet/IP frame length %d (minimum 6)", packetSize))
		}
		if num < packetSize {
			m.log.Trace().Uint32("num", num).Uint32("packetSize", packetSize).Msg("Not enough bytes. Got: num Need: packetSize")
			return nil, nil
		}
		data, err = m.GetTransportInstance().Read(ctx, packetSize)
		if err != nil {
			m.log.Warn().Err(err).Msg("error reading")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		ctxForModel := options.GetLoggerContextForModel(ctx, m.log, options.WithPassLoggerToModel(m.passLogToModel))
		knxMessage, err := model.KnxNetIpMessageParse[model.KnxNetIpMessage](ctxForModel, data)
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
	return func(ctx context.Context, codec _default.DefaultCodecRequirements, message spi.Message) bool {
		// Only the two tunneling frame types get special treatment here. Every
		// other KNXnet/IP frame (SearchResponse, ConnectionResponse,
		// ConnectionStateResponse, DisconnectRequest/Response, ...) has to fall
		// through unhandled, otherwise the handshake- and keepalive-expectations
		// registered on the codec would never see their responses.
		switch typedMessage := message.(type) {
		case model.TunnelingResponse:
			// A TunnelingResponse is the gateway's ACK for a tunneling-request we sent.
			// It must NOT be reported as "handled" here: DefaultCodec.ReceiveWork skips
			// the expectations entirely for a message the custom handler claims, and
			// everything which correlates its ACK (e.g. the group-address write) waits
			// for exactly this frame. An ACK nobody waits for ends up on the default
			// incoming-message channel, where the connection logs it.
			localLog.Trace().Msg("Passing the tunneling ACK on to the expectations")
		case model.TunnelingRequest:
			// If this is an incoming tunneling request, automatically send a tunneling ACK back to the gateway
			response := model.NewTunnelingResponse(
				model.NewTunnelingResponseDataBlock(
					typedMessage.GetTunnelingRequestDataBlock().GetCommunicationChannelId(),
					typedMessage.GetTunnelingRequestDataBlock().GetSequenceCounter(),
					model.Status_NO_ERROR),
			)
			if err := codec.Send(ctx, "tunneling_request", response); err != nil { // TODO: where is a good place to get this timeout from?
				localLog.Warn().Err(err).Msg("got an error sending ACK from transport")
			}
		}

		// Handle the packet itself
		// Give a message interceptor a chance to intercept
		if localCodec, ok := codec.(*MessageCodec); ok && localCodec.messageInterceptor != nil {
			localCodec.messageInterceptor(message)
		}
		return false
	}
}
