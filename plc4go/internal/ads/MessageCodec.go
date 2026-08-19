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

package ads

import (
	"context"
	"encoding/binary"
	"io"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

const (
	// amsTCPHeaderSize is the size of the AMS/TCP header (2 reserved bytes + 4 length bytes).
	amsTCPHeaderSize = 6
	// maxAmsTCPFrameSize is the maximum total frame size (header plus payload) this codec
	// accepts. The length field is a raw 32-bit wire value under peer control, so it must
	// be bounded before any allocation is sized by it.
	maxAmsTCPFrameSize = 16 * 1024 * 1024
)

//go:generate go tool plc4xGenerator -type=MessageCodec
type MessageCodec struct {
	_default.DefaultCodec
	none bool // TODO: just a empty field to satisfy generator (needs fixing because in this case here we have the delegate)

	log zerolog.Logger
}

var (
	_ spi.TransportInstanceExposer = (*MessageCodec)(nil)
)

func NewMessageCodec(transportInstance transports.TransportInstance, _options ...options.WithOption) *MessageCodec {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	codec := &MessageCodec{
		log: customLogger,
	}
	codec.DefaultCodec = _default.NewDefaultCodec(
		codec,
		transportInstance,
		append(_options,
			_default.WithCustomMessageHandler(
				// This just prevents the loop from aborting in the start and by returning false,
				// it makes the message go to the default channel, as this means:
				// The handler hasn't handled the message
				func(ctx context.Context, codec _default.DefaultCodecRequirements, message spi.Message) bool {
					return false
				}),
		)...,
	)
	return codec
}

func (m *MessageCodec) GetCodec() spi.MessageCodec {
	return m
}

func (m *MessageCodec) classifyTransportError(err error) transports.TransportErrorKind {
	if err == nil {
		return transports.TransportErrorUnknown
	}
	if transport := m.GetTransportInstance(); transport != nil {
		return transport.ClassifyError(err)
	}
	return transports.TransportErrorUnknown
}

func (m *MessageCodec) isFatalTransportError(err error) bool {
	if err == nil || transports.ErrorIs(err, io.EOF) {
		return false
	}
	return m.classifyTransportError(err) == transports.TransportErrorFatal
}

func (m *MessageCodec) wrapFatalTransportError(err error, msg string) error {
	if err == nil {
		return nil
	}
	// TODO: Any additional context?
	return transports.NewTransportError(transports.TransportErrorFatal, errors.Wrap(err, msg))
}

func (m *MessageCodec) Send(ctx context.Context, interactionInfo string, message spi.Message) error {
	m.log.Trace().Str("interactionInfo", interactionInfo).Msg("Sending message")
	// Cast the message to the correct type of struct
	tcpPaket := message.(model.AmsTCPPacket)
	// Serialize the request
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	err := tcpPaket.SerializeWithWriteBuffer(ctx, wb)
	if err != nil {
		return errors.Wrap(err, "error serializing request")
	}
	serializedBytes := wb.GetBytes()

	// The generated length arithmetic uses 16-bit accumulation, so payloads over
	// ~8KB wrap the length written into the AMS/TCP header. A frame whose declared
	// length doesn't match the bytes actually serialized would be re-framed by the
	// peer as multiple packets (frame smuggling), so refuse to send it.
	if len(serializedBytes) < amsTCPHeaderSize {
		return errors.Errorf("refusing to send AMS/TCP frame: serialized frame of %d bytes is shorter than the %d byte header", len(serializedBytes), amsTCPHeaderSize)
	}
	declaredLength := binary.LittleEndian.Uint32(serializedBytes[2:6])
	actualLength := uint64(len(serializedBytes) - amsTCPHeaderSize)
	if uint64(declaredLength) != actualLength {
		return errors.Errorf("refusing to send AMS/TCP frame: declared length %d does not match actual payload length %d (payload too large for the 16-bit length arithmetic?)", declaredLength, actualLength)
	}

	// Send it to the PLC
	err = m.GetTransportInstance().Write(ctx, serializedBytes)
	if err != nil {
		return errors.Wrap(err, "error sending request")
	}
	return nil
}

func (m *MessageCodec) Receive(ctx context.Context) (spi.Message, error) {
	transportInstance := m.GetTransportInstance()

	if err := transportInstance.FillBuffer(ctx, func(pos uint, currentByte byte, reader transports.ExtendedReader) bool {
		numBytesAvailable, err := transportInstance.GetNumBytesAvailableInBuffer()
		if err != nil {
			m.log.Warn().Err(err).Msg("error getting bytes while filling buffer")
			return false
		}
		return numBytesAvailable < 6
	}); err != nil {
		m.log.Warn().Err(err).Msg("error filling buffer")
		if m.isFatalTransportError(err) {
			return nil, m.wrapFatalTransportError(err, "error filling buffer")
		}
	}

	// We need at least 6 bytes in order to know how big the packet is in total
	if num, err := transportInstance.GetNumBytesAvailableInBuffer(); (err == nil) && (num >= 6) {
		m.log.Debug().Uint32("num", num).Msg("we got num readable bytes")
		data, err := transportInstance.PeekReadableBytes(ctx, 6)
		if err != nil {
			m.log.Warn().Err(err).Msg("error peeking")
			if m.isFatalTransportError(err) {
				return nil, m.wrapFatalTransportError(err, "error peeking header")
			}
			return nil, nil
		}
		// Get the size of the entire packet little endian plus size of header.
		// Compute in uint64 so a wire length close to 0xFFFFFFFF cannot wrap the
		// packet size to zero (which would make this codec spin forever without
		// consuming a single byte).
		amsTCPLength := (uint64(data[5]) << 24) + (uint64(data[4]) << 16) + (uint64(data[3]) << 8) + (uint64(data[2]))
		if amsTCPLength == 0 || amsTCPLength+amsTCPHeaderSize > maxAmsTCPFrameSize {
			// A zero length can never frame a valid AMS packet and an oversize length
			// would force a huge allocation; both are unrecoverable framing errors.
			return nil, transports.NewTransportError(transports.TransportErrorFatal,
				errors.Errorf("invalid AMS/TCP frame length %d (must be > 0 and at most %d)", amsTCPLength, maxAmsTCPFrameSize-amsTCPHeaderSize))
		}
		packetSize := uint32(amsTCPLength) + amsTCPHeaderSize
		if num < packetSize {
			if err := transportInstance.FillBuffer(ctx, func(pos uint, currentByte byte, reader transports.ExtendedReader) bool {
				numBytesAvailable, err := transportInstance.GetNumBytesAvailableInBuffer()
				if err != nil {
					return false
				}
				return numBytesAvailable < packetSize
			}); err != nil {
				m.log.Warn().Err(err).Msg("error filling buffer")
				if m.isFatalTransportError(err) {
					return nil, errors.Wrap(err, "error filling buffer")
				}
			}
		}
		data, err = transportInstance.Read(ctx, packetSize)
		if err != nil {
			m.log.Warn().Err(err).Msg("error reading packet data")
			if m.isFatalTransportError(err) {
				return nil, m.wrapFatalTransportError(err, "error reading packet data")
			}
			return nil, nil
		}
		rb := utils.NewReadBufferByteBased(data, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
		tcpPacket, err := model.AmsTCPPacketParseWithBuffer(ctx, rb)
		if err != nil {
			m.log.Warn().Err(err).Msg("error parsing")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		return tcpPacket, nil
	} else if err != nil {
		m.log.Warn().Err(err).Msg("Got error reading")
		if m.isFatalTransportError(err) {
			return nil, m.wrapFatalTransportError(err, "error getting readable bytes")
		}
		return nil, nil
	}
	// TODO: maybe we return here a not enough error error
	return nil, nil
}
