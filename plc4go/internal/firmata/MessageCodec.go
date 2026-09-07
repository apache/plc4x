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

package firmata

import (
	"context"
	"io"
	"slices"

	"github.com/rs/zerolog"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/firmata/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

const (
	// firmataMinSize is the smallest a firmata message can be: the one-byte system reset.
	firmataMinSize = uint32(1)
	// sysexEnd terminates a sysex message. It is the only way to know how long one is.
	sysexEnd = byte(0xF7)
)

var (
	// errNotEnoughData says the frame at the head of the stream isn't complete yet.
	errNotEnoughData = errors.New("not enough data")
	// errOutOfSync says the byte at the head of the stream cannot start a firmata message.
	errOutOfSync = errors.New("not the start of a firmata message")
)

// MessageCodec frames the firmata (MIDI-style) wire format. Every message starts with a status
// byte whose high nibble identifies the message type and whose low nibble carries the pin or
// channel, so the length of a message follows from its first byte - except for sysex messages,
// which are framed by a trailing 0xF7 and have to be scanned for.
//
// Ported from plc4j's FirmataMessageCodec (and the ByteLengthEstimator before it):
//   - 0xE0 (analog IO) and 0x90 (digital IO) are 3 bytes,
//   - 0xC0 (subscribe analog) and 0xD0 (subscribe digital) are 2 bytes,
//   - 0xF0 (system message) depends on the low nibble: 0x00 is a sysex terminated by 0xF7,
//     0x04/0x05/0x09 are 3 bytes and 0x0F (system reset) is a single byte.
//
// Anything else is garbage. Where plc4j throws, this codec resynchronizes byte-wise the way the
// modbus RTU codec does: a serial line that was joined mid-message otherwise never recovers.
//
//go:generate go tool plc4xGenerator -type=MessageCodec
type MessageCodec struct {
	_default.DefaultCodec

	passLogToModel bool

	// resyncSkippedBytes counts the bytes thrown away since the stream was last in sync. Only
	// Receive touches it, which the default codec runs on a single worker.
	resyncSkippedBytes uint64

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
	// Register a no-op (always-false) custom handler so the receive worker keeps polling the
	// transport even when there is not a single outstanding expectation. The default worker skips
	// receiving while `expectations` is empty AND `customMessageHandling` is nil, and firmata
	// registers exactly one expectation in its whole life (the connect handshake): without a
	// handler the worker would park right after Connect and the board's unsolicited analog-IO and
	// digital-IO messages would sit in the transport buffer with nobody draining them. Returning
	// false leaves all real handling to the default path, which hands the message to the
	// connection's incoming-message worker via defaultIncomingMessageChannel.
	//
	// The option list is clipped before appending so the extra option can never be written into a
	// backing array the caller still shares with somebody else.
	codec.DefaultCodec = _default.NewDefaultCodec(codec, transportInstance,
		append(slices.Clip(_options), _default.WithCustomMessageHandler(keepReceiveLoopActive))...)
	return codec
}

// keepReceiveLoopActive is a no-op CustomMessageHandler. Its only purpose is to make the default
// codec's customMessageHandling non-nil so its receive worker doesn't park when the expectations
// drain to zero. Returning false defers every message to the default handling.
func keepReceiveLoopActive(_ context.Context, _ _default.DefaultCodecRequirements, _ spi.Message) bool {
	return false
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
	if err == nil || isEOF(err) {
		return false
	}
	return m.classifyTransportError(err) == transports.TransportErrorFatal
}

func (m *MessageCodec) wrapFatalTransportError(err error, msg string) error {
	if err == nil {
		return nil
	}
	return transports.NewTransportError(transports.TransportErrorFatal, errors.Wrap(err, msg))
}

func isEOF(err error) bool {
	if err == nil {
		return false
	}
	return transports.ErrorIs(err, io.EOF)
}

func (m *MessageCodec) Send(ctx context.Context, interactionInfo string, message spi.Message) error {
	m.log.Trace().Str("interactionInfo", interactionInfo).Msg("Sending message")
	firmataMessage, ok := message.(readWriteModel.FirmataMessage)
	if !ok {
		return errors.Errorf("message is not a FirmataMessage, got %T", message)
	}
	theBytes, err := firmataMessage.Serialize()
	if err != nil {
		return errors.Wrap(err, "error serializing request")
	}
	if err := m.GetTransportInstance().Write(ctx, theBytes); err != nil {
		return errors.Wrap(err, "error sending request")
	}
	return nil
}

// firmataFrameSize is the total length of the message starting at buffer[0], derived from the
// status byte at the head of the stream. It reports errNotEnoughData while the length cannot be
// decided yet and errOutOfSync when the head cannot start a message at all.
func firmataFrameSize(buffer []byte) (int, error) {
	if len(buffer) < 1 {
		return 0, errNotEnoughData
	}
	first := buffer[0]
	switch first & 0xF0 {
	case 0xE0, 0x90:
		// Analog and digital IO carry a two byte payload.
		return 3, nil
	case 0xC0, 0xD0:
		// Subscribe analog / digital carry a single byte payload.
		return 2, nil
	case 0xF0:
		switch first & 0x0F {
		case 0x00:
			// Sysex: everything up to and including the terminator belongs to the message. A
			// terminator can never be mistaken for payload, as payload bytes always have their
			// most significant bit clear.
			for i := 1; i < len(buffer); i++ {
				if buffer[i] == sysexEnd {
					return i + 1, nil
				}
			}
			return 0, errNotEnoughData
		case 0x04, 0x05, 0x09:
			// Set pin mode, set digital pin value, protocol version.
			return 3, nil
		case 0x0F:
			// System reset.
			return 1, nil
		default:
			return 0, errOutOfSync
		}
	default:
		// The remaining high nibbles (and every byte with its most significant bit clear, which
		// is a payload byte we joined the stream in the middle of) cannot start a message.
		return 0, errOutOfSync
	}
}

// needsMoreBytes says whether the buffer doesn't hold a complete message yet. A head which can't
// start a message stops the filling: it is the receive loop that resynchronizes, and it can only do
// that once it gets to run.
func (m *MessageCodec) needsMoreBytes(ctx context.Context, ti transports.TransportInstance) bool {
	numBytesAvailable, err := ti.GetNumBytesAvailableInBuffer()
	if err != nil {
		m.log.Debug().Err(err).Msg("error getting available bytes")
		return false
	}
	if numBytesAvailable < firmataMinSize {
		return true
	}
	buffer, err := ti.PeekReadableBytes(ctx, numBytesAvailable)
	if err != nil {
		m.log.Debug().Err(err).Msg("error peeking buffer while filling")
		return false
	}
	size, err := firmataFrameSize(buffer)
	if err != nil {
		return errors.Is(err, errNotEnoughData)
	}
	return numBytesAvailable < uint32(size)
}

// Receive frames one message off the stream. A frame is only consumed once it has been parsed, and
// a head which cannot start a message is thrown away a byte at a time until the stream lines up
// again.
func (m *MessageCodec) Receive(ctx context.Context) (spi.Message, error) {
	ti := m.GetTransportInstance()
	if !ti.IsConnected() {
		return nil, errors.New("Transport instance not connected")
	}

	// 1. Fill the buffer until the message at the head of the stream is complete. Sizing while
	// filling is what a format without a length field needs: a sysex is only as long as the
	// distance to its terminator, so how much is missing can only be known by looking at what is
	// already there (plc4j's FirmataMessageCodec peeks for the same reason).
	if err := ti.FillBuffer(ctx, func(pos uint, currentByte byte, reader transports.ExtendedReader) bool {
		m.log.Trace().Uint("pos", pos).Uint8("currentByte", currentByte).Msg("filling")
		return m.needsMoreBytes(ctx, ti)
	}); err != nil {
		if m.isFatalTransportError(err) {
			m.log.Debug().Err(err).Msg("error filling buffer")
			return nil, m.wrapFatalTransportError(err, "error filling buffer")
		}
		// Fall through on non-fatal errors, we might have enough data...
	}

	ctxForModel := options.GetLoggerContextForModel(ctx, m.log, options.WithPassLoggerToModel(m.passLogToModel))

	// 2. Walk the buffer, skipping a byte at a time whenever the head isn't the start of a
	// message. Every turn either returns or consumes at least one byte, so this terminates.
	for {
		numBytesAvail, err := ti.GetNumBytesAvailableInBuffer()
		if err != nil && numBytesAvail < firmataMinSize {
			if isEOF(err) {
				m.log.Debug().Msg("transport buffer exhausted while checking availability")
				return nil, nil
			}
			m.log.Warn().Err(err).Msg("error getting buffer length")
			return nil, m.wrapFatalTransportError(err, "error getting buffer length")
		}
		if numBytesAvail < firmataMinSize {
			return nil, nil
		}

		// Everything currently buffered, because a sysex message can only be sized by scanning
		// for its terminator.
		buffer, err := ti.PeekReadableBytes(ctx, numBytesAvail)
		if err != nil {
			m.log.Warn().Err(err).Msg("error peeking buffer")
			if m.isFatalTransportError(err) {
				return nil, m.wrapFatalTransportError(err, "error peeking buffer")
			}
			return nil, nil
		}

		expectedSize, err := firmataFrameSize(buffer)
		switch {
		case errors.Is(err, errOutOfSync):
			if skipped, skipErr := m.skipOneByte(ctx, "0x%02X cannot start a message", buffer[0]); !skipped {
				return nil, skipErr
			}
			continue
		case errors.Is(err, errNotEnoughData):
			m.log.Debug().
				Uint32("num", numBytesAvail).
				Msg("Received fragment. Waiting for more data...")
			return nil, nil
		case err != nil:
			return nil, errors.Wrap(err, "error sizing the message")
		}
		if numBytesAvail < uint32(expectedSize) {
			// Never consume a partial message - wait for the rest to arrive. Unlike a modbus RTU
			// frame, the size of a firmata message follows from a single byte whose most
			// significant bit tells a status byte from payload, so a byte-shifted candidate is
			// rejected outright rather than by running out of data.
			m.log.Debug().
				Uint32("num", numBytesAvail).
				Int("messageSize", expectedSize).
				Msg("Received fragment. Waiting for more data...")
			return nil, nil
		}

		frame := buffer[:expectedSize]
		// Incoming messages are what the board sends us, which is the 'response' side of the
		// mspec's typeSwitch - the same argument plc4j's FirmataMessageCodec parses with.
		message, err := readWriteModel.FirmataMessageParse[readWriteModel.FirmataMessage](ctxForModel, frame, true)
		if err != nil {
			m.log.Debug().
				Err(err).
				Stringer("frame", utils.Base64Stringer(frame)).
				Msg("error parsing message")
			if skipped, skipErr := m.skipOneByte(ctx, "message failed to parse: %s", err); !skipped {
				return nil, skipErr
			}
			continue
		}

		if _, err := ti.Read(ctx, uint32(expectedSize)); err != nil {
			m.log.Debug().Err(err).Msg("error consuming parsed message")
			return nil, nil
		}
		m.noteResyncComplete()
		return message, nil
	}
}

// skipOneByte throws the byte at the head of the stream away and notes that we are resynchronizing.
// Only the first byte of a run is logged at warn level, so a long stretch of garbage costs one line
// rather than one per byte. It reports whether the byte really went away; when it didn't, nothing
// was consumed and the caller has to yield rather than loop.
func (m *MessageCodec) skipOneByte(ctx context.Context, reason string, args ...any) (bool, error) {
	if m.resyncSkippedBytes == 0 {
		m.log.Warn().Msgf("Firmata stream out of sync ("+reason+"), resynchronizing byte-wise", args...)
	} else {
		m.log.Trace().Msgf("Firmata resync skipping a byte ("+reason+")", args...)
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

func (m *MessageCodec) noteResyncComplete() {
	if m.resyncSkippedBytes > 0 {
		m.log.Warn().Uint64("skippedBytes", m.resyncSkippedBytes).Msg("Firmata stream resynchronized")
		m.resyncSkippedBytes = 0
	}
}
