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

package umas

import (
	"context"
	"encoding/binary"

	"github.com/rs/zerolog"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/umas/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

const (
	// mbapHeaderSize is the part of the Modbus/TCP header a frame's total length can be read from:
	// transaction identifier (2), protocol identifier (2) and the length field (2). The unit
	// identifier the header ends with is already counted by that length field.
	mbapHeaderSize = uint32(6)
	// mbapLengthOverhead is what has to be added to the length field to get the size of the whole
	// frame: the field counts the unit identifier plus the PDU, so everything but the six bytes it
	// is preceded by.
	mbapLengthOverhead = uint32(6)
	// mbapMinimumLength is the smallest value the length field may carry: the unit identifier plus
	// at least a function code.
	mbapMinimumLength = uint32(2)
	// modbusProtocolIdentifier is the constant the mspec pins on the protocol identifier field.
	modbusProtocolIdentifier = uint16(0x0000)
)

// MessageCodec frames the UMAS wire format, which is UMAS tunneled inside Modbus/TCP: a 7 byte MBAP
// header (transaction identifier, protocol identifier, length, unit identifier) followed by a PDU
// whose function code is 0x5A. The length field at offset 4 is big endian and counts the unit
// identifier plus the PDU, so a frame is length + 6 bytes long - the same arithmetic plc4j's
// UmasMessageCodec does.
//
// The byte orders are mixed and the model pins them per field: the MBAP header is big endian per
// the Modbus specification, everything inside the UMAS PDU is little endian. The codec therefore
// never picks an order of its own, it only reads the length field by hand.
//
// A UMAS response carries the generic function key 0xFE and needs the *request's* function key to be
// parsed into the right response type. The codec is the one place which sees both directions, so it
// records the key of every outgoing request in its functionKeyTracker and looks it up again by the
// transaction identifier at the head of the incoming frame.
//
//go:generate go tool plc4xGenerator -type=MessageCodec
type MessageCodec struct {
	_default.DefaultCodec

	functionKeys *functionKeyTracker `ignore:"true"`

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
		functionKeys:   newFunctionKeyTracker(),
		passLogToModel: passLoggerToModel,
		log:            customLogger,
	}
	codec.DefaultCodec = _default.NewDefaultCodec(codec, transportInstance, _options...)
	return codec
}

func (m *MessageCodec) GetCodec() spi.MessageCodec {
	return m
}

// Send serializes one ADU and remembers its UMAS function key on the way out, so Receive can hand
// the parser what it needs to tell the response types apart.
func (m *MessageCodec) Send(ctx context.Context, interactionInfo string, message spi.Message) error {
	m.log.Trace().Str("interactionInfo", interactionInfo).Msg("Sending message")
	adu, ok := message.(readWriteModel.ModbusTcpADU)
	if !ok {
		return errors.Errorf("message is not a ModbusTcpADU, got %T", message)
	}
	// Tracking before the write, not after: the response may be back before Write returns.
	if umasPdu, isUmas := adu.GetPdu().(readWriteModel.UmasPDU); isUmas {
		m.functionKeys.trackRequest(adu.GetTransactionIdentifier(), umasPdu.GetItem().GetUmasFunctionKey())
	}
	theBytes, err := adu.Serialize()
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
	if numBytesAvailable < mbapHeaderSize {
		return nil, nil
	}
	header, err := transportInstance.PeekReadableBytes(ctx, mbapHeaderSize)
	if err != nil {
		m.log.Warn().Err(err).Msg("error peeking the MBAP header")
		return nil, nil
	}
	if !headerLooksLikeMbap(header) {
		// Nothing in this stream can resynchronize us other than throwing bytes away: the frame
		// format has no start marker. Dropping a single byte at a time lets the next Receive find a
		// header which starts one byte further in.
		m.log.Warn().
			Stringer("header", utils.Base64Stringer(header)).
			Msg("The stream head is not a Modbus/TCP header, discarding one byte")
		if _, discardErr := transportInstance.Read(ctx, 1); discardErr != nil {
			m.log.Debug().Err(discardErr).Msg("error discarding a byte")
		}
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
	// The transaction identifier is the first field of the header, big endian, and it is what ties
	// the response to the function key of its request.
	requestFunctionKey := m.functionKeys.consumeFunctionKey(binary.BigEndian.Uint16(data[0:2]))
	ctxForModel := options.GetLoggerContextForModel(ctx, m.log, options.WithPassLoggerToModel(m.passLogToModel))
	adu, err := readWriteModel.ModbusTcpADUParse(ctxForModel, data, requestFunctionKey)
	if err != nil {
		// The frame was consumed, so a malformed packet costs us exactly that packet and the stream
		// stays in sync for the next one.
		m.log.Warn().Err(err).
			Uint8("requestFunctionKey", requestFunctionKey).
			Stringer("data", utils.Base64Stringer(data)).
			Msg("error parsing")
		return nil, nil
	}
	return adu, nil
}

// needsMoreBytes says whether the buffer doesn't hold a complete frame yet: either the length field
// isn't readable, or it says the frame is longer than what has arrived so far.
func (m *MessageCodec) needsMoreBytes(ctx context.Context, transportInstance transports.TransportInstance) bool {
	numBytesAvailable, err := transportInstance.GetNumBytesAvailableInBuffer()
	if err != nil {
		m.log.Debug().Err(err).Msg("error getting available bytes")
		return false
	}
	if numBytesAvailable < mbapHeaderSize {
		return true
	}
	header, err := transportInstance.PeekReadableBytes(ctx, mbapHeaderSize)
	if err != nil {
		m.log.Debug().Err(err).Msg("error peeking the length field")
		return false
	}
	if !headerLooksLikeMbap(header) {
		// Let Receive deal with the garbage rather than waiting for a frame which will never come.
		return false
	}
	return numBytesAvailable < packetSizeFromHeader(header)
}

// headerLooksLikeMbap is the cheapest check that tells a Modbus/TCP header from garbage: the
// protocol identifier is a constant zero and the length field has to cover at least the unit
// identifier and a function code. Without it a desynchronized stream makes the codec wait for a
// frame of up to 65541 bytes which nobody is going to send.
func headerLooksLikeMbap(header []byte) bool {
	if len(header) < int(mbapHeaderSize) {
		return false
	}
	if binary.BigEndian.Uint16(header[2:4]) != modbusProtocolIdentifier {
		return false
	}
	return uint32(binary.BigEndian.Uint16(header[4:6])) >= mbapMinimumLength
}

// packetSizeFromHeader is the total length of the frame whose first bytes are passed in. The
// addition happens in uint32 on purpose: a wire length of 0xFFFF would wrap a uint16 total back to
// almost nothing, making Read consume less than a frame and leaving the receive worker spinning on
// bytes it can never resynchronize from.
func packetSizeFromHeader(header []byte) uint32 {
	return uint32(binary.BigEndian.Uint16(header[4:6])) + mbapLengthOverhead
}
