/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package modbus

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

const (
	// modbusRtuMinSize is the smallest an RTU frame can be: address (1) + function code (1) +
	// CRC (2).
	modbusRtuMinSize = uint32(4)
	// rtuExceptionResponseSize is address + function code + exception code + CRC.
	rtuExceptionResponseSize = 5
	// rtuWriteEchoResponseSize is address + function code + 2 address bytes + 2 value/quantity
	// bytes + CRC.
	rtuWriteEchoResponseSize = 8
	// rtuSizingHeaderSize is address + function code + the third byte, which carries the byte
	// count of a read response. That is everything the size table below looks at.
	rtuSizingHeaderSize = uint32(3)
)

// MessageCodecRtu is the codec of the serial flavor of modbus. An RTU frame is nothing but a
// station address, the PDU and a CRC - there is no length field and no start delimiter, so a frame
// can only be recognized by deriving its length from the function code and then checking the CRC.
//
// Ported from plc4j's ModbusRtuMessageCodec, which in turn mirrors the structure of the Modbus TCP
// codec next door: peek a candidate, hand it to the generated parser (which validates the CRC
// through the checksum field), and only consume it once the parser accepted it. Anything else
// advances the stream by a single byte and tries again - byte-wise resynchronization, which is the
// only recovery a framing without delimiters allows.
//
// Limitation, same as plc4j's: the size table is response-shaped, so this codec frames the
// responses a master receives. It cannot frame a stream of raw RTU requests.
//
//go:generate go tool plc4xGenerator -type=MessageCodecRtu
type MessageCodecRtu struct {
	_default.DefaultCodec

	passLogToModel bool

	// resyncSkippedBytes counts the bytes thrown away since the stream was last in sync. It is
	// only touched from Receive, which the default codec runs on a single worker.
	resyncSkippedBytes uint64

	log zerolog.Logger
}

var (
	_ spi.TransportInstanceExposer = (*MessageCodecRtu)(nil)
)

func NewMessageCodecRtu(transportInstance transports.TransportInstance, _options ...options.WithOption) *MessageCodecRtu {
	passLoggerToModel, _ := options.ExtractPassLoggerToModel(_options...)
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	codec := &MessageCodecRtu{
		passLogToModel: passLoggerToModel,
		log:            customLogger,
	}
	codec.DefaultCodec = _default.NewDefaultCodec(codec, transportInstance, _options...)
	return codec
}

func (m *MessageCodecRtu) GetCodec() spi.MessageCodec {
	return m
}

func (m *MessageCodecRtu) classifyTransportError(err error) transports.TransportErrorKind {
	if err == nil {
		return transports.TransportErrorUnknown
	}
	if transport := m.GetTransportInstance(); transport != nil {
		return transport.ClassifyError(err)
	}
	return transports.TransportErrorUnknown
}

func (m *MessageCodecRtu) isFatalTransportError(err error) bool {
	if err == nil || isEOF(err) {
		return false
	}
	return m.classifyTransportError(err) == transports.TransportErrorFatal
}

func (m *MessageCodecRtu) wrapFatalTransportError(err error, msg string) error {
	if err == nil {
		return nil
	}
	return transports.NewTransportError(transports.TransportErrorFatal, errors.Wrap(err, msg))
}

func (m *MessageCodecRtu) Send(ctx context.Context, interactionInfo string, message spi.Message) error {
	m.log.Trace().Str("interactionInfo", interactionInfo).Msg("Sending message")
	// Cast the message to the correct type of struct
	rtuAdu, ok := message.(model.ModbusRtuADU)
	if !ok {
		return errors.Errorf("message is not a ModbusRtuADU, got %T", message)
	}
	// Serialize the request. The CRC is a checksum field of the generated model, so serializing
	// is all it takes to have a correct one on the wire.
	theBytes, err := rtuAdu.Serialize()
	if err != nil {
		return errors.Wrap(err, "error serializing request")
	}

	// Send it to the PLC
	if err := m.GetTransportInstance().Write(ctx, theBytes); err != nil {
		return errors.Wrap(err, "error sending request")
	}
	return nil
}

// rtuResponseSize computes the size of the response frame that starts with the given sizing
// header. RTU has no length field, so the size follows from the function code:
//   - an exception (function code >= 0x80) is 5 bytes,
//   - a read response (0x01-0x04 and the file-record operations 0x14, 0x15, 0x17, whose responses
//     carry a byte count covering the whole item array just like the plain reads) is 5 bytes plus
//     the byte count in the third header byte,
//   - a write echo (0x05, 0x06, 0x0F, 0x10) is a fixed 8 bytes.
//
// Everything else returns -1, which the caller treats as "the stream is out of sync here" rather
// than as "wait for more data".
func rtuResponseSize(header []byte) int {
	if len(header) < int(rtuSizingHeaderSize) {
		return -1
	}
	functionCode := header[1]
	if functionCode >= 0x80 {
		return rtuExceptionResponseSize
	}
	switch functionCode {
	case 0x01, 0x02, 0x03, 0x04, 0x14, 0x15, 0x17:
		return int(modbusRtuMinSize) + 1 + int(header[2])
	case 0x05, 0x06, 0x0F, 0x10:
		return rtuWriteEchoResponseSize
	default:
		return -1
	}
}

// Receive peeks a candidate frame, only consumes it once the generated parser has validated its
// CRC, and otherwise resynchronizes byte-wise.
//
// While resynchronizing - a previous candidate at this stream position already failed - a header
// that looks valid but demands more bytes than are available is treated as another failed
// candidate rather than as a reason to wait: a byte-shifted candidate can easily look like a
// legitimate function code with a large byte count, and waiting for bytes that were never going to
// arrive would stall the codec forever. Waiting is only correct at a clean stream position, where
// an undersized buffer really does mean the next frame hasn't fully arrived yet. The trade-off is
// that a genuine partial frame sitting directly behind garbage is skipped rather than awaited;
// that response is lost and recovered by the request timeout, which beats stalling every frame
// queued behind a garbage candidate that will never complete. The mirror image holds at a clean
// position: garbage that happens to look like a large read response is waited out until enough
// bytes have arrived to reject it, because there is nothing to tell it apart from a frame that is
// still on its way.
func (m *MessageCodecRtu) Receive(ctx context.Context) (spi.Message, error) {
	ti := m.GetTransportInstance()
	if !ti.IsConnected() {
		return nil, errors.New("Transport instance not connected")
	}

	// 1. Fill the buffer until at least the smallest frame could be in there
	if err := ti.FillBuffer(ctx, func(pos uint, currentByte byte, reader transports.ExtendedReader) bool {
		m.log.Trace().Uint("pos", pos).Uint8("currentByte", currentByte).Msg("filling")
		numBytesAvailable, err := ti.GetNumBytesAvailableInBuffer()
		if err != nil {
			m.log.Debug().Err(err).Msg("error getting available bytes")
			return false
		}
		return numBytesAvailable < modbusRtuMinSize
	}); err != nil {
		if m.isFatalTransportError(err) {
			m.log.Debug().Err(err).Msg("error filling buffer")
			return nil, m.wrapFatalTransportError(err, "error filling buffer")
		}
		// Fall through on non-fatal errors, we might have enough data...
	}

	ctxForModel := options.GetLoggerContextForModel(ctx, m.log, options.WithPassLoggerToModel(m.passLogToModel))

	// 2. Walk the buffer, skipping a byte at a time whenever the head isn't the start of a frame.
	// Every turn of this loop either returns or consumes at least one byte, so it terminates.
	for {
		numBytesAvail, err := ti.GetNumBytesAvailableInBuffer()
		if err != nil && numBytesAvail < modbusRtuMinSize {
			if isEOF(err) {
				m.log.Debug().Msg("transport buffer exhausted while checking availability")
				return nil, nil
			}
			m.log.Warn().Err(err).Msg("error getting buffer length")
			return nil, m.wrapFatalTransportError(err, "error getting buffer length")
		}
		if numBytesAvail < modbusRtuMinSize {
			return nil, nil
		}

		header, err := ti.PeekReadableBytes(ctx, rtuSizingHeaderSize)
		if err != nil {
			m.log.Warn().Err(err).Msg("error peeking header")
			if m.isFatalTransportError(err) {
				return nil, m.wrapFatalTransportError(err, "error peeking header")
			}
			return nil, nil
		}

		expectedSize := rtuResponseSize(header)
		if expectedSize < 0 {
			if skipped, err := m.skipOneByte(ctx, "unknown function code 0x%x", header[1]); !skipped {
				return nil, err
			}
			continue
		}
		if numBytesAvail < uint32(expectedSize) {
			if m.resyncSkippedBytes > 0 {
				if skipped, err := m.skipOneByte(ctx, "candidate frame during resync needs %d bytes, only %d available", expectedSize, numBytesAvail); !skipped {
					return nil, err
				}
				continue
			}
			// Never consume a partial frame - wait for the rest to arrive.
			m.log.Debug().
				Uint32("num", numBytesAvail).
				Int("packetSize", expectedSize).
				Msg("Received fragment. Waiting for more data...")
			return nil, nil
		}

		frame, err := ti.PeekReadableBytes(ctx, uint32(expectedSize))
		if err != nil {
			m.log.Warn().Err(err).Msg("error peeking frame slice")
			if m.isFatalTransportError(err) {
				return nil, m.wrapFatalTransportError(err, "error peeking frame slice")
			}
			return nil, nil
		}

		// The parser validates the CRC on the way through, so a frame that comes out of here is a
		// frame that was really on the wire - which is what makes peek-then-consume safe.
		rtuAdu, err := model.ModbusADUParse[model.ModbusRtuADU](ctxForModel, frame, model.DriverType_MODBUS_RTU, true)
		if err != nil {
			if skipped, skipErr := m.skipOneByte(ctx, "frame failed validation: %s", err); !skipped {
				return nil, skipErr
			}
			continue
		}

		if _, err := ti.Read(ctx, uint32(expectedSize)); err != nil {
			m.log.Debug().Err(err).Msg("error consuming parsed frame")
			return nil, nil
		}
		m.noteResyncComplete()
		return rtuAdu, nil
	}
}

// skipOneByte throws the byte at the head of the stream away and notes that we are resynchronizing.
// Only the first byte of a run is logged at warn level, so a long stretch of garbage costs one line
// rather than one per byte. It reports whether the byte really went away; when it didn't, nothing
// was consumed and the caller has to yield rather than loop (the error is fatal-wrapped or nil,
// following the same rule as the rest of the codec).
func (m *MessageCodecRtu) skipOneByte(ctx context.Context, reason string, args ...any) (bool, error) {
	if m.resyncSkippedBytes == 0 {
		m.log.Warn().Msgf("Modbus RTU stream out of sync ("+reason+"), resynchronizing byte-wise", args...)
	} else {
		m.log.Trace().Msgf("Modbus RTU resync skipping a byte ("+reason+")", args...)
	}
	if _, err := m.GetTransportInstance().Read(ctx, 1); err != nil {
		m.log.Debug().Err(err).Msg("error discarding byte during resync")
		if m.isFatalTransportError(err) {
			return false, m.wrapFatalTransportError(err, "error discarding byte during resync")
		}
		return false, nil
	}
	m.resyncSkippedBytes++
	return true, nil
}

func (m *MessageCodecRtu) noteResyncComplete() {
	if m.resyncSkippedBytes > 0 {
		m.log.Warn().Uint64("skippedBytes", m.resyncSkippedBytes).Msg("Modbus RTU stream resynchronized")
		m.resyncSkippedBytes = 0
	}
}
