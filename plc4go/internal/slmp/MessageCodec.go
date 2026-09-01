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

package slmp

import (
	"context"
	"encoding/binary"

	"github.com/rs/zerolog"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/slmp/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

const (
	// frame3EHeaderSize is the fixed part of a 3E frame that precedes the length field's own
	// contribution: the 2-byte subheader, the 5-byte access route and the 2-byte length field
	// itself. Everything the length field counts follows it, so a frame is always
	// frame3EHeaderSize + length bytes long.
	frame3EHeaderSize = uint32(9)
	// frame3ELengthFieldOffset is where the 2-byte data length sits. It is the same offset in a
	// request (requestDataLength) and in a response (responseDataLength), because the subheader and
	// the access route have the same width in both.
	frame3ELengthFieldOffset = uint32(7)
)

// MessageCodec frames the SLMP 3E binary wire format over a TCP byte stream. Ported from plc4j's
// SlmpMessageCodec.
//
// The length field is a *little-endian* uint16 at offset 7 - SLMP is little-endian throughout,
// unlike modbus, whose MBAP length at a similar spot is big-endian. Reading it big-endian would
// turn a 12-byte response into a 3072-byte one and stall the receive loop forever, so this is the
// one thing in the codec worth being loud about.
//
//go:generate go tool plc4xGenerator -type=MessageCodec
type MessageCodec struct {
	_default.DefaultCodec

	passLogToModel bool

	log zerolog.Logger
}

var (
	_ spi.TransportInstanceExposer = (*MessageCodec)(nil)
)

func NewMessageCodec(transportInstance transports.TransportInstance, _options ...options.WithOption) *MessageCodec {
	passLoggerToModel, _ := options.ExtractPassLoggerToModel(_options...)
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	codec := &MessageCodec{
		passLogToModel: passLoggerToModel,
		log:            customLogger,
	}
	codec.DefaultCodec = _default.NewDefaultCodec(codec, transportInstance, _options...)
	return codec
}

func (m *MessageCodec) GetCodec() spi.MessageCodec {
	return m
}

func (m *MessageCodec) Send(ctx context.Context, interactionInfo string, message spi.Message) error {
	m.log.Trace().Str("interactionInfo", interactionInfo).Msg("Sending message")
	frame, ok := message.(readWriteModel.SlmpMessage)
	if !ok {
		return errors.Errorf("message is not an SlmpMessage, got %T", message)
	}
	// The mspec pins the byte order on the message itself, so the model decides the encoding and a
	// plain buffer is enough here.
	theBytes, err := frame.Serialize()
	if err != nil {
		return errors.Wrap(err, "error serializing request")
	}
	if err := m.GetTransportInstance().Write(ctx, theBytes); err != nil {
		return errors.Wrap(err, "error sending request")
	}
	return nil
}

func (m *MessageCodec) Receive(ctx context.Context) (spi.Message, error) {
	transportInstance := m.GetTransportInstance()
	// Pull data from the transport until a whole frame is buffered. Some transports (e.g. the test
	// transport) only surface queued data through fills, so stopping the fill as soon as the length
	// field is readable would leave the rest of the frame sitting in the transport forever.
	if err := transportInstance.FillBuffer(ctx, func(_ uint, _ byte, _ transports.ExtendedReader) bool {
		return m.needsMoreBytes(ctx, transportInstance)
	}); err != nil {
		if transportError, ok := transports.AsTransportError(err); ok && transportError.Kind() == transports.TransportErrorFatal {
			return nil, err
		}
		// Fall through on non-fatal errors, we might have enough data buffered already.
		m.log.Trace().Err(err).Msg("Error filling buffer, continuing with what's available")
	}

	numBytesAvailable, err := transportInstance.GetNumBytesAvailableInBuffer()
	if err != nil {
		m.log.Warn().Err(err).Msg("Got error reading")
		return nil, nil
	}
	if numBytesAvailable < frame3EHeaderSize {
		return nil, nil
	}
	header, err := transportInstance.PeekReadableBytes(ctx, frame3EHeaderSize)
	if err != nil {
		m.log.Warn().Err(err).Msg("error peeking")
		return nil, nil
	}
	frameSize := frameSizeFromHeader(header)
	if numBytesAvailable < frameSize {
		m.log.Debug().
			Uint32("numBytesAvailable", numBytesAvailable).
			Uint32("frameSize", frameSize).
			Msg("Not enough bytes yet")
		return nil, nil
	}
	data, err := transportInstance.Read(ctx, frameSize)
	if err != nil {
		m.log.Debug().Err(err).Msg("Error reading")
		return nil, nil
	}
	frame, err := readWriteModel.SlmpMessageParse[readWriteModel.SlmpMessage](ctx, data)
	if err != nil {
		// The frame was consumed, so a malformed frame costs us exactly that frame and the stream
		// stays in sync for the next one.
		m.log.Warn().Err(err).Msg("error parsing")
		return nil, nil
	}
	return frame, nil
}

// needsMoreBytes says whether the buffer doesn't hold a complete frame yet: either the header isn't
// readable, or the length field says the frame is longer than what has arrived so far.
func (m *MessageCodec) needsMoreBytes(ctx context.Context, transportInstance transports.TransportInstance) bool {
	numBytesAvailable, err := transportInstance.GetNumBytesAvailableInBuffer()
	if err != nil {
		m.log.Debug().Err(err).Msg("error getting available bytes")
		return false
	}
	if numBytesAvailable < frame3EHeaderSize {
		return true
	}
	header, err := transportInstance.PeekReadableBytes(ctx, frame3EHeaderSize)
	if err != nil {
		m.log.Debug().Err(err).Msg("error peeking the length field")
		return false
	}
	return numBytesAvailable < frameSizeFromHeader(header)
}

// frameSizeFromHeader is the total length of the frame whose first frame3EHeaderSize bytes are
// passed in. The addition happens in uint32 on purpose: a wire length of 0xFFFF would wrap a uint16
// total back to almost nothing, making Read consume less than a frame and leaving the receive worker
// spinning on bytes it can never resynchronize from.
func frameSizeFromHeader(header []byte) uint32 {
	return uint32(binary.LittleEndian.Uint16(header[frame3ELengthFieldOffset:])) + frame3EHeaderSize
}
