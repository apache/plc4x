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

package eip

import (
	"context"
	"encoding/binary"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

//go:generate go tool plc4xGenerator -type=MessageCodec
type MessageCodec struct {
	_default.DefaultCodec

	none      bool // TODO: just a empty field to satisfy generator (needs fixing because in this case here we have the delegate)
	byteOrder binary.ByteOrder

	log zerolog.Logger
}

var (
	_ spi.TransportInstanceExposer = (*MessageCodec)(nil)
)

func NewMessageCodec(transportInstance transports.TransportInstance, byteOrder binary.ByteOrder, _options ...options.WithOption) *MessageCodec {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	if byteOrder == nil {
		byteOrder = binary.LittleEndian
	}
	codec := &MessageCodec{
		byteOrder: byteOrder,
		log:       customLogger,
	}
	codec.DefaultCodec = _default.NewDefaultCodec(codec, transportInstance, _options...)
	return codec
}

func (m *MessageCodec) GetCodec() spi.MessageCodec {
	return m
}

func (m *MessageCodec) Send(ctx context.Context, interactionInfo string, message spi.Message) error {
	m.log.Trace().Str("interactionInfo", interactionInfo).Msg("Sending message")
	// Cast the message to the correct type of struct
	eipPacket := message.(model.EipPacket)
	// Serialize the request
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(m.byteOrder))
	err := eipPacket.SerializeWithWriteBuffer(ctx, wb)
	if err != nil {
		return errors.Wrap(err, "error serializing request")
	}

	// Send it to the PLC
	err = m.GetTransportInstance().Write(ctx, wb.GetBytes())
	if err != nil {
		return errors.Wrap(err, "error sending request")
	}
	return nil
}

func (m *MessageCodec) Receive(ctx context.Context) (spi.Message, error) {
	// We need at least 6 bytes in order to know how big the packet is in total
	transportInstance := m.GetTransportInstance()
	// Pull data from the transport until at least the 4-byte header is buffered.
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
	if num, err := transportInstance.GetNumBytesAvailableInBuffer(); (err == nil) && (num >= 4) {
		m.log.Debug().Uint32("num", num).Msg("we got num readable bytes")
		data, err := transportInstance.PeekReadableBytes(ctx, 4)
		if err != nil {
			m.log.Warn().Err(err).Msg("error peeking")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		packetSize := packetSizeFromHeader(data, m.byteOrder)
		if num < packetSize {
			m.log.Debug().Uint32("num", num).Uint32("packetSize", packetSize).Msg("Not enough bytes. Got: num Need: packetSize")
			return nil, nil
		}
		data, err = transportInstance.Read(ctx, packetSize)
		if err != nil {
			m.log.Debug().Err(err).Msg("Error reading")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		rb := utils.NewReadBufferByteBased(data, utils.WithByteOrderForReadBufferByteBased(m.byteOrder))
		eipPacket, err := model.EipPacketParseWithBuffer[model.EipPacket](ctx, rb, true)
		if err != nil {
			m.log.Warn().Err(err).Msg("error parsing")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		return eipPacket, nil
	} else if err != nil {
		m.log.Warn().Err(err).Msg("Got error reading")
		return nil, nil
	}
	// TODO: maybe we return here a not enough error error
	return nil, nil
}

// packetSizeFromHeader computes the full packet size from the first 4 header
// bytes. The addition must happen in uint32: doing it in uint16 lets a wire
// length like 0xFFE8 wrap the total to 0, making Read consume nothing and the
// receive worker spin forever on the same unconsumed bytes.
func packetSizeFromHeader(header []byte, byteOrder binary.ByteOrder) uint32 {
	return uint32(byteOrder.Uint16(header[2:4])) + 24
}
