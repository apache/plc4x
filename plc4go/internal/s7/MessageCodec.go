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

package s7

import (
	"context"
	"sync"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

//go:generate go tool plc4xGenerator -type=MessageCodec
type MessageCodec struct {
	_default.DefaultCodec

	// unsolicitedUserData carries pushed UserData messages (cyclic data, alarm indications,
	// mode transitions) that no request expectation claims. Closed on Disconnect.
	unsolicitedUserData     chan model.S7MessageUserData `ignore:"true"`
	unsolicitedCloseOnce    sync.Once                    `ignore:"true"`

	passLogToModel bool
	log            zerolog.Logger
}

var (
	_ spi.TransportInstanceExposer = (*MessageCodec)(nil)
)

func NewMessageCodec(transportInstance transports.TransportInstance, _options ...options.WithOption) *MessageCodec {
	passLoggerToModel, _ := options.ExtractPassLoggerToModel(_options...)
	extractCustomLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	codec := &MessageCodec{
		unsolicitedUserData: make(chan model.S7MessageUserData, 100),
		passLogToModel:      passLoggerToModel,
		log:                 extractCustomLogger,
	}
	codec.DefaultCodec = _default.NewDefaultCodec(codec, transportInstance,
		append(_options, _default.WithCustomMessageHandler(extractUnsolicitedUserData(codec)))...)
	return codec
}

// extractUnsolicitedUserData claims pushed UserData messages (cpuFunctionType 0x00) before
// the request expectations see them and forwards them to the subscription dispatch.
func extractUnsolicitedUserData(codec *MessageCodec) _default.CustomMessageHandler {
	return func(ctx context.Context, _ _default.DefaultCodecRequirements, message spi.Message) bool {
		tpktPacket, ok := message.(model.TPKTPacket)
		if !ok {
			return false
		}
		cotpPacketData, ok := tpktPacket.GetPayload().(model.COTPPacketData)
		if !ok {
			return false
		}
		messageUserData, ok := cotpPacketData.GetPayload().(model.S7MessageUserData)
		if !ok {
			return false
		}
		_, functionType, _, ok := userDataPushKey(messageUserData)
		if !ok || functionType != 0x00 {
			return false
		}
		select {
		case codec.unsolicitedUserData <- messageUserData:
		default:
			codec.log.Warn().Msg("Unsolicited user data message discarded, channel full")
		}
		return true
	}
}

// GetUnsolicitedUserData exposes the push channel drained by the connection's dispatcher.
func (m *MessageCodec) GetUnsolicitedUserData() <-chan model.S7MessageUserData {
	return m.unsolicitedUserData
}

func (m *MessageCodec) Disconnect() error {
	err := m.DefaultCodec.Disconnect()
	m.unsolicitedCloseOnce.Do(func() {
		close(m.unsolicitedUserData)
	})
	return err
}

func (m *MessageCodec) GetCodec() spi.MessageCodec {
	return m
}

func (m *MessageCodec) Send(ctx context.Context, interactionInfo string, message spi.Message) error {
	m.log.Trace().Str("interactionInfo", interactionInfo).Msg("Sending message")
	// Cast the message to the correct type of struct
	tpktPacket := message.(model.TPKTPacket)
	// Serialize the request
	theBytes, err := tpktPacket.Serialize()
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
	transportInstance := m.GetTransportInstance()
	if !transportInstance.IsConnected() {
		return nil, errors.New("Transport instance not connected")
	}
	// Pull data from the transport until at least the 4-byte TPKT header is buffered.
	// Some transports (e.g. the test transport) only surface queued data through fills,
	// so checking the buffer without filling first would starve the receive worker.
	if err := transportInstance.FillBuffer(ctx, func(pos uint, currentByte byte, reader transports.ExtendedReader) bool {
		numBytesAvailable, err := transportInstance.GetNumBytesAvailableInBuffer()
		if err != nil {
			return false
		}
		return numBytesAvailable < 4
	}); err != nil {
		if transportError, ok := transports.AsTransportError(err); ok && transportError.Kind() == transports.TransportErrorFatal {
			return nil, err
		}
		// Fall through on non-fatal errors, we might have enough data buffered already.
		m.log.Trace().Err(err).Msg("Error filling buffer, continuing with what's available")
	}
	// We need at least 4 bytes in order to know how big the packet is in total
	if num, err := m.GetTransportInstance().GetNumBytesAvailableInBuffer(); (err == nil) && (num >= 4) {
		m.log.Debug().Uint32("num", num).Msg("we got %d readable bytes")
		data, err := m.GetTransportInstance().PeekReadableBytes(ctx, 4)
		if err != nil {
			m.log.Warn().Err(err).Msg("error peeking")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		// Get the size of the entire packet
		packetSize := (uint32(data[2]) << 8) + uint32(data[3])
		// A TPKT length below the 4 byte header can never become parseable and
		// would make the receive worker spin forever on the same unconsumed
		// bytes, so treat it as a fatal framing error.
		if packetSize < 4 {
			return nil, transports.NewTransportError(transports.TransportErrorFatal,
				errors.Errorf("invalid TPKT packet length %d (minimum 4)", packetSize))
		}
		if num < packetSize {
			m.log.Debug().Uint32("num", num).Uint32("packetSize", packetSize).Msg("Not enough bytes. Got: num Need: packetSize")
			return nil, nil
		}
		data, err = m.GetTransportInstance().Read(ctx, packetSize)
		if err != nil {
			m.log.Debug().Err(err).Msg("Error reading")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		ctxForModel := options.GetLoggerContextForModel(ctx, m.log, options.WithPassLoggerToModel(m.passLogToModel))
		tpktPacket, err := model.TPKTPacketParse(ctxForModel, data)
		if err != nil {
			m.log.Warn().Err(err).Msg("error parsing")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		return tpktPacket, nil
	} else if err != nil {
		m.log.Warn().Err(err).Msg("Got error reading")
		return nil, nil
	}
	// TODO: maybe we return here a not enough error error
	return nil, nil
}
