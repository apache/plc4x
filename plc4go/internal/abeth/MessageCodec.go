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

package abeth

import (
	"context"
	"encoding/binary"

	"github.com/rs/zerolog"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/abeth/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

const (
	// cipHeaderSize is the fixed part every CIP encapsulation packet starts with: command type,
	// packet length, session handle, status, 8 bytes of sender context, options and 4 reserved
	// bytes. The mspec spells it out as 'lengthInBytes - 28' on the implicit length field.
	cipHeaderSize = uint32(28)
	// cipLengthFieldEnd is how many bytes have to be buffered before the length of the frame at the
	// head of the stream can be decided: the command type plus the length field itself.
	cipLengthFieldEnd = uint32(4)
)

// MessageCodec frames the Allen-Bradley ETH wire format, which is CIP encapsulation over plain TCP:
// a 28 byte header whose length field at offset 2 counts only what follows the header, so a frame
// is always packetLen + 28 bytes long. Ported from plc4j's AbEthMessageCodec (and the
// ByteLengthEstimator the legacy driver used before it).
//
// Unlike EtherNet/IP, ab-eth is big-endian throughout - the mspec pins the byte order on every
// single field, so the model decides the encoding and the codec only has to agree on the length
// field it peeks at by hand.
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
	packet, ok := message.(readWriteModel.CIPEncapsulationPacket)
	if !ok {
		return errors.Errorf("message is not a CIPEncapsulationPacket, got %T", message)
	}
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := packet.SerializeWithWriteBuffer(ctx, wb); err != nil {
		return errors.Wrap(err, "error serializing request")
	}
	if err := m.GetTransportInstance().Write(ctx, wb.GetBytes()); err != nil {
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
	if numBytesAvailable < cipLengthFieldEnd {
		return nil, nil
	}
	header, err := transportInstance.PeekReadableBytes(ctx, cipLengthFieldEnd)
	if err != nil {
		m.log.Warn().Err(err).Msg("error peeking")
		return nil, nil
	}
	packetSize := packetSizeFromHeader(header)
	if numBytesAvailable < packetSize {
		m.log.Debug().
			Uint32("numBytesAvailable", numBytesAvailable).
			Uint32("packetSize", packetSize).
			Msg("Not enough bytes yet")
		return nil, nil
	}
	data, err := transportInstance.Read(ctx, packetSize)
	if err != nil {
		m.log.Debug().Err(err).Msg("Error reading")
		return nil, nil
	}
	rb := utils.NewReadBufferByteBased(data, utils.WithByteOrderForReadBufferByteBased(binary.BigEndian))
	packet, err := readWriteModel.CIPEncapsulationPacketParseWithBuffer[readWriteModel.CIPEncapsulationPacket](ctx, rb)
	if err != nil {
		// The frame was consumed, so a malformed packet costs us exactly that packet and the stream
		// stays in sync for the next one.
		m.log.Warn().Err(err).Msg("error parsing")
		return nil, nil
	}
	return packet, nil
}

// needsMoreBytes says whether the buffer doesn't hold a complete frame yet: either the length field
// isn't readable, or it says the frame is longer than what has arrived so far.
func (m *MessageCodec) needsMoreBytes(ctx context.Context, transportInstance transports.TransportInstance) bool {
	numBytesAvailable, err := transportInstance.GetNumBytesAvailableInBuffer()
	if err != nil {
		m.log.Debug().Err(err).Msg("error getting available bytes")
		return false
	}
	if numBytesAvailable < cipLengthFieldEnd {
		return true
	}
	header, err := transportInstance.PeekReadableBytes(ctx, cipLengthFieldEnd)
	if err != nil {
		m.log.Debug().Err(err).Msg("error peeking the length field")
		return false
	}
	return numBytesAvailable < packetSizeFromHeader(header)
}

// packetSizeFromHeader is the total length of the frame whose first bytes are passed in. The
// addition has to happen in uint32: a wire length of 0xFFFF would wrap a uint16 total back to
// almost nothing, making Read consume less than a frame and leaving the receive worker spinning on
// bytes it can never resynchronize from.
func packetSizeFromHeader(header []byte) uint32 {
	return uint32(binary.BigEndian.Uint16(header[2:4])) + cipHeaderSize
}
