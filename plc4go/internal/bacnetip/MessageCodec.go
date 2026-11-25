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

// Receive takes one datagram from the transport and parses the BVLC message in it.
//
// BACnet/IP runs over UDP, and a datagram carries exactly one BVLC message (ASHRAE 135 Annex J).
// The udp transport buffers one datagram at a time and reads the next only once the buffer is
// empty, so whatever is buffered is the rest of the current datagram, and nothing more of it will
// ever arrive. Receive therefore consumes it whole, whatever it holds. Waiting for more bytes, as a
// stream codec would, stalls every later read behind a datagram that is shorter than its BVLC
// length says (or shorter than the 4 byte header): the transport never reads the next datagram
// while those bytes sit in the buffer. A malformed datagram is dropped and logged; it does not end
// the connection, as anyone who can reach the socket can send one.
func (m *MessageCodec) Receive(ctx context.Context) (spi.Message, error) {
	transportInstance := m.GetTransportInstance()
	num, err := transportInstance.GetNumBytesAvailableInBuffer()
	if err != nil {
		m.log.Warn().Err(err).Msg("Got error reading")
		return nil, nil
	}
	if num == 0 {
		return nil, nil
	}
	m.log.Debug().Uint32("num", num).Msg("we got num readable bytes")
	datagram, err := transportInstance.Read(ctx, num)
	if err != nil {
		m.log.Debug().Err(err).Msg("Error reading")
		return nil, nil
	}
	if len(datagram) < 4 {
		m.log.Warn().Int("size", len(datagram)).Msg("dropping a datagram shorter than the 4 byte BVLC header")
		return nil, nil
	}
	packetSize := (uint32(datagram[2]) << 8) + uint32(datagram[3])
	if packetSize < 4 || packetSize > num {
		m.log.Warn().
			Uint32("packetSize", packetSize).
			Uint32("size", num).
			Msg("dropping a datagram whose BVLC length does not fit it")
		return nil, nil
	}
	if packetSize < num {
		m.log.Debug().
			Uint32("packetSize", packetSize).
			Uint32("size", num).
			Msg("dropping the bytes after the BVLC message in a datagram")
	}
	ctxForModel := options.GetLoggerContextForModel(ctx, m.log, options.WithPassLoggerToModel(m.passLogToModel))
	bvlcPacket, err := model.BVLCParse[model.BVLC](ctxForModel, datagram[:packetSize])
	if err != nil {
		m.log.Warn().Err(err).Msg("error parsing")
		return nil, nil
	}
	return bvlcPacket, nil
}
