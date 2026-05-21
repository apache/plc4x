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

package bacnetip

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

//go:generate go tool plc4xGenerator -type=MessageCodec
type MessageCodec struct {
	_default.DefaultCodec

	passLogToModel bool
	log            zerolog.Logger
}

var (
	_ spi.TransportInstanceExposer = (*MessageCodec)(nil)
)

func NewMessageCodec(transportInstance transports.TransportInstance, _options ...options.WithOption) *MessageCodec {
	codec := &MessageCodec{}
	// Register a no-op (always-false) custom handler so the codec's receive
	// loop keeps polling even when there are zero outstanding expectations.
	// The default loop skips reading the transport when both `expectations`
	// is empty AND `customMessageHandling` is nil — that would mean unsolicited
	// COV notifications arrive in the kernel buffer but nobody drains them.
	// Returning false here lets the default expectation-matching path run
	// first, then falls through to defaultIncomingMessageChannel for the
	// Connection's COV-notification poller.
	codec.DefaultCodec = _default.NewDefaultCodec(codec, transportInstance, append(_options, _default.WithCustomMessageHandler(keepReceiveLoopActive))...)
	return codec
}

// keepReceiveLoopActive is a no-op CustomMessageHandler. Its only purpose is
// to set m.customMessageHandling to non-nil so the codec's receive worker
// doesn't park when expectations drain to zero. Returning false defers all
// real handling to HandleMessages → defaultIncomingMessageChannel.
func keepReceiveLoopActive(_ context.Context, _ _default.DefaultCodecRequirements, _ spi.Message) bool {
	return false
}

func (m *MessageCodec) GetCodec() spi.MessageCodec {
	return m
}

func (m *MessageCodec) Send(ctx context.Context, interactionInfo string, message spi.Message) error {
	m.log.Trace().Str("interactionInfo", interactionInfo).Msg("Sending message")
	// Cast the message to the correct type of struct
	bvlcPacket := message.(model.BVLC)
	// Serialize the request
	theBytes, err := bvlcPacket.Serialize()
	if err != nil {
		return errors.Wrap(err, "error serializing request")
	}

	// Send it to the PLC
	err = m.GetTransportInstance().Write(ctx, theBytes)
	if err != nil {
		return errors.Wrap(err, "error sending request")
	}
	return nil
}

func (m *MessageCodec) Receive(ctx context.Context) (spi.Message, error) {
	// We need at least 6 bytes in order to know how big the packet is in total
	if num, err := m.GetTransportInstance().GetNumBytesAvailableInBuffer(); (err == nil) && (num >= 4) {
		m.log.Debug().Uint32("num", num).Msg("we got num readable bytes")
		data, err := m.GetTransportInstance().PeekReadableBytes(ctx, 4)
		if err != nil {
			m.log.Warn().Err(err).Msg("error peeking")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		packetSize := uint32((uint16(data[2]) << 8) + uint16(data[3]))
		if num < packetSize {
			m.log.Debug().
				Uint32("num", num).
				Uint32("packetSize", packetSize).
				Msg("Not enough bytes. Got: num Need: packetSize")
			return nil, nil
		}
		data, err = m.GetTransportInstance().Read(ctx, packetSize)
		if err != nil {
			m.log.Debug().Err(err).Msg("Error reading")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		ctxForModel := options.GetLoggerContextForModel(ctx, m.log, options.WithPassLoggerToModel(m.passLogToModel))
		bvlcPacket, err := model.BVLCParse[model.BVLC](ctxForModel, data)
		if err != nil {
			m.log.Warn().Err(err).Msg("error parsing")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		return bvlcPacket, nil
	} else if err != nil {
		m.log.Warn().Err(err).Msg("Got error reading")
		return nil, nil
	}
	// TODO: maybe we return here a not enough error error
	return nil, nil
}
