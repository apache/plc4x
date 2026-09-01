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
	"bytes"
	"context"
	"encoding/hex"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

const (
	// asciiFrameStart is the colon every Modbus ASCII frame opens with. It is the one byte of the
	// framing that can't appear inside a frame, which is what makes resynchronization cheap here.
	asciiFrameStart = byte(':')

	// modbusAsciiMinSize is the smallest a frame can be: ':' + address (2 hex chars) + function
	// code (2) + LRC (2) + CR + LF.
	modbusAsciiMinSize = uint32(9)

	// modbusAsciiMinPayloadSize is the smallest the binary form behind the hex can be: address,
	// function code and LRC.
	modbusAsciiMinPayloadSize = 3

	// modbusAsciiMaxSize bounds a frame: the binary form is an address, a PDU of at most 253 bytes
	// and an LRC, each byte spelled as two hex characters, plus the colon and the CR/LF. Nothing
	// longer can be a frame, so a colon that hasn't been terminated within this many bytes is
	// garbage rather than a frame that is still arriving.
	modbusAsciiMaxSize = uint32(1 + 2*(1+253+1) + 2)
)

// asciiHexChars is the alphabet frames are written with. Modbus ASCII is specified with upper case
// hex; readers have to accept both cases, which encoding/hex does.
const asciiHexChars = "0123456789ABCDEF"

// MessageCodecAscii is the codec of the ASCII flavor of modbus. A frame is a colon, the binary ADU
// (address, PDU and LRC) spelled out as hex characters, and a CR/LF terminator - so unlike RTU
// there is a delimiter to look for, and unlike TCP there is no length field.
//
// Ported from plc4j's ModbusAsciiMessageCodec: frame on the CR/LF terminator, strip the hex layer,
// and hand the binary ADU to the generated parser, which validates the LRC through the checksum
// field. Anything that doesn't survive that is not a frame, and the codec resynchronizes onto the
// next colon - which is cheaper than the byte-wise walk RTU has to do, because ASCII has a start
// delimiter that cannot appear inside a frame.
//
// Limitation, same as plc4j's: the parser is asked for responses, so this codec frames what a
// master receives. It cannot frame a stream of raw ASCII requests.
//
//go:generate go tool plc4xGenerator -type=MessageCodecAscii
type MessageCodecAscii struct {
	_default.DefaultCodec

	passLogToModel bool

	// resyncSkippedBytes counts the bytes thrown away since the stream was last in sync. It is
	// only touched from Receive, which the default codec runs on a single worker.
	resyncSkippedBytes uint64

	log zerolog.Logger
}

var (
	_ spi.TransportInstanceExposer = (*MessageCodecAscii)(nil)
)

func NewMessageCodecAscii(transportInstance transports.TransportInstance, _options ...options.WithOption) *MessageCodecAscii {
	passLoggerToModel, _ := options.ExtractPassLoggerToModel(_options...)
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	codec := &MessageCodecAscii{
		passLogToModel: passLoggerToModel,
		log:            customLogger,
	}
	codec.DefaultCodec = _default.NewDefaultCodec(codec, transportInstance, _options...)
	return codec
}

func (m *MessageCodecAscii) GetCodec() spi.MessageCodec {
	return m
}

func (m *MessageCodecAscii) classifyTransportError(err error) transports.TransportErrorKind {
	if err == nil {
		return transports.TransportErrorUnknown
	}
	if transport := m.GetTransportInstance(); transport != nil {
		return transport.ClassifyError(err)
	}
	return transports.TransportErrorUnknown
}

func (m *MessageCodecAscii) isFatalTransportError(err error) bool {
	if err == nil || isEOF(err) {
		return false
	}
	return m.classifyTransportError(err) == transports.TransportErrorFatal
}

func (m *MessageCodecAscii) wrapFatalTransportError(err error, msg string) error {
	if err == nil {
		return nil
	}
	return transports.NewTransportError(transports.TransportErrorFatal, errors.Wrap(err, msg))
}

// encodeAsciiFrame turns the binary form of an ADU into the frame that goes on the wire.
func encodeAsciiFrame(binaryBytes []byte) []byte {
	frame := make([]byte, 0, 1+len(binaryBytes)*2+2)
	frame = append(frame, asciiFrameStart)
	for _, b := range binaryBytes {
		frame = append(frame, asciiHexChars[(b>>4)&0x0F], asciiHexChars[b&0x0F])
	}
	return append(frame, '\r', '\n')
}

func (m *MessageCodecAscii) Send(ctx context.Context, interactionInfo string, message spi.Message) error {
	m.log.Trace().Str("interactionInfo", interactionInfo).Msg("Sending message")
	// Cast the message to the correct type of struct
	asciiAdu, ok := message.(model.ModbusAsciiADU)
	if !ok {
		return errors.Errorf("message is not a ModbusAsciiADU, got %T", message)
	}
	// Serialize the request into its binary form. The LRC is a checksum field of the generated
	// model, so serializing is all it takes to have a correct one; the hex layer goes on top.
	theBytes, err := asciiAdu.Serialize()
	if err != nil {
		return errors.Wrap(err, "error serializing request")
	}

	// Send it to the PLC
	if err := m.GetTransportInstance().Write(ctx, encodeAsciiFrame(theBytes)); err != nil {
		return errors.Wrap(err, "error sending request")
	}
	return nil
}

// asciiFrameEnd reports the offset of the CR of the CR/LF that terminates the frame starting at the
// head of data, or -1 when the terminator hasn't arrived yet. The colon itself is skipped, so a
// frame is never terminated before it has a payload.
func asciiFrameEnd(data []byte) int {
	if len(data) < 2 {
		return -1
	}
	index := bytes.Index(data[1:], []byte{'\r', '\n'})
	if index < 0 {
		return -1
	}
	return index + 1
}

// decodeAsciiPayload turns the hex characters between the colon and the terminator back into the
// binary ADU. Everything it rejects is a malformed frame, which the caller resynchronizes past.
func decodeAsciiPayload(payload []byte) ([]byte, error) {
	if len(payload) < 2*modbusAsciiMinPayloadSize {
		return nil, errors.Errorf("frame carries %d hex characters, needs at least %d", len(payload), 2*modbusAsciiMinPayloadSize)
	}
	if len(payload)%2 != 0 {
		return nil, errors.Errorf("frame carries an odd number of hex characters (%d)", len(payload))
	}
	binaryData := make([]byte, len(payload)/2)
	if _, err := hex.Decode(binaryData, payload); err != nil {
		return nil, errors.Wrap(err, "frame is not hex encoded")
	}
	return binaryData, nil
}

// Receive frames on the colon and the CR/LF terminator, only consumes a frame once the generated
// parser has validated its LRC, and otherwise resynchronizes onto the next colon.
//
// A colon that hasn't been terminated yet is waited for rather than skipped - that is what a frame
// still trickling in over a serial line looks like - but only up to the longest frame the protocol
// allows. Past that the colon can't be the start of a frame, so it is thrown away instead of
// stalling the codec forever on a terminator that was never going to arrive.
func (m *MessageCodecAscii) Receive(ctx context.Context) (spi.Message, error) {
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
		return numBytesAvailable < modbusAsciiMinSize
	}); err != nil {
		if m.isFatalTransportError(err) {
			m.log.Debug().Err(err).Msg("error filling buffer")
			return nil, m.wrapFatalTransportError(err, "error filling buffer")
		}
		// Fall through on non-fatal errors, we might have enough data...
	}

	ctxForModel := options.GetLoggerContextForModel(ctx, m.log, options.WithPassLoggerToModel(m.passLogToModel))

	// 2. Walk the buffer, skipping ahead whenever the head isn't the start of a frame. Every turn
	// of this loop either returns or consumes at least one byte, so it terminates.
	for {
		numBytesAvail, err := ti.GetNumBytesAvailableInBuffer()
		if err != nil && numBytesAvail < modbusAsciiMinSize {
			if isEOF(err) {
				m.log.Debug().Msg("transport buffer exhausted while checking availability")
				return nil, nil
			}
			m.log.Warn().Err(err).Msg("error getting buffer length")
			return nil, m.wrapFatalTransportError(err, "error getting buffer length")
		}
		if numBytesAvail < modbusAsciiMinSize {
			return nil, nil
		}

		data, err := ti.PeekReadableBytes(ctx, numBytesAvail)
		if err != nil {
			m.log.Warn().Err(err).Msg("error peeking buffer")
			if m.isFatalTransportError(err) {
				return nil, m.wrapFatalTransportError(err, "error peeking buffer")
			}
			return nil, nil
		}

		if data[0] != asciiFrameStart {
			// Nothing before the next colon can be part of a frame, so the whole run goes at
			// once rather than a byte at a time. If there is no colon in what has arrived, all
			// of it is garbage.
			skip := uint32(len(data))
			if next := bytes.IndexByte(data, asciiFrameStart); next > 0 {
				skip = uint32(next)
			}
			if skipped, err := m.skipBytes(ctx, skip, "0x%x where a frame should start", data[0]); !skipped {
				return nil, err
			}
			continue
		}

		end := asciiFrameEnd(data)
		if end < 0 {
			if numBytesAvail >= modbusAsciiMaxSize {
				if skipped, err := m.skipBytes(ctx, 1, "no frame terminator within %d bytes", numBytesAvail); !skipped {
					return nil, err
				}
				continue
			}
			// Never consume a partial frame - wait for the terminator to arrive.
			m.log.Debug().
				Uint32("num", numBytesAvail).
				Msg("Received fragment. Waiting for more data...")
			return nil, nil
		}
		frameSize := uint32(end + 2)

		binaryData, err := decodeAsciiPayload(data[1:end])
		if err != nil {
			if skipped, skipErr := m.skipBytes(ctx, 1, "%s", err); !skipped {
				return nil, skipErr
			}
			continue
		}

		// The parser validates the LRC on the way through, so a frame that comes out of here is a
		// frame that was really on the wire - which is what makes peek-then-consume safe.
		asciiAdu, err := model.ModbusADUParse[model.ModbusAsciiADU](ctxForModel, binaryData, model.DriverType_MODBUS_ASCII, true)
		if err != nil {
			if skipped, skipErr := m.skipBytes(ctx, 1, "frame failed validation: %s", err); !skipped {
				return nil, skipErr
			}
			continue
		}

		if _, err := ti.Read(ctx, frameSize); err != nil {
			m.log.Debug().Err(err).Msg("error consuming parsed frame")
			return nil, nil
		}
		m.noteResyncComplete()
		return asciiAdu, nil
	}
}

// skipBytes throws count bytes away from the head of the stream and notes that we are
// resynchronizing. Only the first skip of a run is logged at warn level, so a long stretch of
// garbage costs one line rather than one per skip. It reports whether the bytes really went away;
// when they didn't, nothing was consumed and the caller has to yield rather than loop (the error is
// fatal-wrapped or nil, following the same rule as the rest of the codec).
func (m *MessageCodecAscii) skipBytes(ctx context.Context, count uint32, reason string, args ...any) (bool, error) {
	if count == 0 {
		return false, nil
	}
	if m.resyncSkippedBytes == 0 {
		m.log.Warn().Msgf("Modbus ASCII stream out of sync ("+reason+"), resynchronizing on the next frame start", args...)
	} else {
		m.log.Trace().Msgf("Modbus ASCII resync skipping %d byte(s) ("+reason+")", append([]any{count}, args...)...)
	}
	if _, err := m.GetTransportInstance().Read(ctx, count); err != nil {
		m.log.Debug().Err(err).Msg("error discarding bytes during resync")
		if m.isFatalTransportError(err) {
			return false, m.wrapFatalTransportError(err, "error discarding bytes during resync")
		}
		return false, nil
	}
	m.resyncSkippedBytes += uint64(count)
	return true, nil
}

func (m *MessageCodecAscii) noteResyncComplete() {
	if m.resyncSkippedBytes > 0 {
		m.log.Warn().Uint64("skippedBytes", m.resyncSkippedBytes).Msg("Modbus ASCII stream resynchronized")
		m.resyncSkippedBytes = 0
	}
}
