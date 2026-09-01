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

package iec608705104

import (
	"context"
	"slices"

	"github.com/rs/zerolog"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/iec608705104/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

const (
	// startByte is the fixed first octet of every APCI, the same constant the mspec pins as
	// APDU.startByte.
	startByte = byte(0x68)
	// lengthFieldEnd is how many octets have to be buffered before the length of the frame at the
	// head of the stream can be decided: the start byte plus the length octet itself.
	lengthFieldEnd = uint32(2)
	// apciLengthOverhead is what the length octet does not count: the start byte and itself.
	apciLengthOverhead = uint32(2)
	// minApciLength is the smallest value the length octet may take. Every APDU carries a four octet
	// control field, and the shortest ones (the U-format frames and the S-format acknowledgement)
	// carry nothing else.
	minApciLength = 4
	// maxApciLength is the largest value the length octet can take, which is what a single octet
	// holds. It is also the reason an IEC 60870-5-104 ASDU can never exceed 249 octets.
	maxApciLength = 253
)

var (
	// errNotEnoughData says the frame at the head of the stream isn't complete yet.
	errNotEnoughData = errors.New("not enough data")
	// errOutOfSync says the octets at the head of the stream cannot start an APDU.
	errOutOfSync = errors.New("not the start of an APDU")
)

// MessageCodec frames the IEC 60870-5-104 APCI, which is about as simple as a wire format gets: the
// fixed start octet 0x68, a one octet length which counts everything after itself, and then that
// many octets of control field and ASDU. Ported from plc4j's Iec60870MessageCodec.
//
// Where plc4j throws on a bad start octet - which kills the connection and takes every subscription
// with it - this codec resynchronizes octet by octet the way the modbus RTU and firmata codecs do. A
// controlled station that has been streaming telemetry for a week is not worth dropping over one
// corrupted octet.
type MessageCodec struct {
	_default.DefaultCodec

	passLogToModel bool

	// resyncSkippedBytes counts the octets thrown away since the stream was last in sync. Only
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
	// A no-op (always false) custom handler keeps the receive worker polling the transport even
	// while there is not a single outstanding expectation. IEC 60870-5-104 registers expectations
	// only for the two handshake round trips; everything after them is unsolicited, so without a
	// handler the worker would park right after connecting and the station's telemetry would pile up
	// in the transport buffer with nobody draining it. Returning false leaves all real handling to
	// the default path, which hands the message to the connection's incoming-message worker.
	codec.DefaultCodec = _default.NewDefaultCodec(codec, transportInstance,
		append(slices.Clip(_options), _default.WithCustomMessageHandler(keepReceiveLoopActive))...)
	return codec
}

// keepReceiveLoopActive is a no-op CustomMessageHandler whose only purpose is to make the default
// codec's custom message handling non-nil, so its receive worker doesn't park when the expectations
// drain to zero.
func keepReceiveLoopActive(_ context.Context, _ _default.DefaultCodecRequirements, _ spi.Message) bool {
	return false
}

func (m *MessageCodec) GetCodec() spi.MessageCodec {
	return m
}

func (m *MessageCodec) Send(ctx context.Context, interactionInfo string, message spi.Message) error {
	m.log.Trace().Str("interactionInfo", interactionInfo).Msg("Sending message")
	apdu, ok := message.(readWriteModel.APDU)
	if !ok {
		return errors.Errorf("message is not an APDU, got %T", message)
	}
	// The mspec pins little-endian on the APDU itself, so the model's own Serialize already writes
	// the control field in wire order.
	theBytes, err := apdu.Serialize()
	if err != nil {
		return errors.Wrap(err, "error serializing request")
	}
	if err := m.GetTransportInstance().Write(ctx, theBytes); err != nil {
		return errors.Wrap(err, "error sending request")
	}
	return nil
}

// apduFrameSize is the total length of the APDU starting at buffer[0]. It reports errNotEnoughData
// while the length octet hasn't arrived yet and errOutOfSync when the head cannot start an APDU at
// all - either because the start octet is wrong or because the length octet is impossible.
func apduFrameSize(buffer []byte) (uint32, error) {
	if len(buffer) < int(lengthFieldEnd) {
		if len(buffer) > 0 && buffer[0] != startByte {
			// No point waiting for a length octet behind an octet which cannot start a frame.
			return 0, errOutOfSync
		}
		return 0, errNotEnoughData
	}
	if buffer[0] != startByte {
		return 0, errOutOfSync
	}
	apciLength := uint32(buffer[1])
	if apciLength < minApciLength || apciLength > maxApciLength {
		return 0, errOutOfSync
	}
	return apciLength + apciLengthOverhead, nil
}

// needsMoreBytes says whether the buffer doesn't hold a complete APDU yet. A head which can't start
// one stops the filling: it is the receive loop that resynchronizes, and it can only do that once it
// gets to run.
func (m *MessageCodec) needsMoreBytes(ctx context.Context, ti transports.TransportInstance) bool {
	numBytesAvailable, err := ti.GetNumBytesAvailableInBuffer()
	if err != nil {
		m.log.Debug().Err(err).Msg("error getting available bytes")
		return false
	}
	if numBytesAvailable < lengthFieldEnd {
		return true
	}
	buffer, err := ti.PeekReadableBytes(ctx, lengthFieldEnd)
	if err != nil {
		m.log.Debug().Err(err).Msg("error peeking the length field")
		return false
	}
	size, err := apduFrameSize(buffer)
	if err != nil {
		return errors.Is(err, errNotEnoughData)
	}
	return numBytesAvailable < size
}

// Receive frames one APDU off the stream. A frame is only consumed once it has been parsed, and a
// head which cannot start an APDU is thrown away an octet at a time until the stream lines up again.
func (m *MessageCodec) Receive(ctx context.Context) (spi.Message, error) {
	ti := m.GetTransportInstance()
	if !ti.IsConnected() {
		return nil, errors.New("Transport instance not connected")
	}

	if err := ti.FillBuffer(ctx, func(_ uint, _ byte, _ transports.ExtendedReader) bool {
		return m.needsMoreBytes(ctx, ti)
	}); err != nil {
		if transportError, ok := transports.AsTransportError(err); ok && transportError.Kind() == transports.TransportErrorFatal {
			return nil, err
		}
		// Fall through on non-fatal errors, we might have enough data buffered already.
		m.log.Trace().Err(err).Msg("Error filling buffer, continuing with what's available")
	}

	ctxForModel := options.GetLoggerContextForModel(ctx, m.log, options.WithPassLoggerToModel(m.passLogToModel))

	// Every turn of this loop either returns or consumes at least one octet, so it terminates.
	for {
		numBytesAvailable, err := ti.GetNumBytesAvailableInBuffer()
		if err != nil {
			m.log.Warn().Err(err).Msg("error getting available bytes")
			return nil, nil
		}
		if numBytesAvailable < 1 {
			return nil, nil
		}

		peekSize := min(numBytesAvailable, lengthFieldEnd)
		header, err := ti.PeekReadableBytes(ctx, peekSize)
		if err != nil {
			m.log.Warn().Err(err).Msg("error peeking the length field")
			return nil, nil
		}

		frameSize, err := apduFrameSize(header)
		switch {
		case errors.Is(err, errOutOfSync):
			if skipped, skipErr := m.skipOneByte(ctx, "0x%02X cannot start an APDU", header[0]); !skipped {
				return nil, skipErr
			}
			continue
		case errors.Is(err, errNotEnoughData):
			return nil, nil
		case err != nil:
			return nil, errors.Wrap(err, "error sizing the APDU")
		}
		if numBytesAvailable < frameSize {
			m.log.Debug().
				Uint32("numBytesAvailable", numBytesAvailable).
				Uint32("frameSize", frameSize).
				Msg("Received fragment. Waiting for more data...")
			return nil, nil
		}

		frame, err := ti.PeekReadableBytes(ctx, frameSize)
		if err != nil {
			m.log.Warn().Err(err).Msg("error peeking the frame")
			return nil, nil
		}
		apdu, err := readWriteModel.APDUParse[readWriteModel.APDU](ctxForModel, frame)
		if err != nil {
			m.log.Debug().
				Err(err).
				Stringer("frame", utils.Base64Stringer(frame)).
				Msg("error parsing APDU")
			// The length octet said where the next frame starts, so a frame which doesn't parse is
			// dropped whole rather than resynchronized through - that keeps one malformed ASDU from
			// costing the frames behind it.
			if _, readErr := ti.Read(ctx, frameSize); readErr != nil {
				m.log.Debug().Err(readErr).Msg("error discarding an unparseable frame")
			}
			return nil, nil
		}

		if _, err := ti.Read(ctx, frameSize); err != nil {
			m.log.Debug().Err(err).Msg("error consuming the parsed frame")
			return nil, nil
		}
		m.noteResyncComplete()
		return apdu, nil
	}
}

// skipOneByte throws the octet at the head of the stream away and notes that we are resynchronizing.
// Only the first octet of a run is logged at warn level, so a long stretch of garbage costs one line
// rather than one per octet. It reports whether the octet really went away; when it didn't, nothing
// was consumed and the caller has to yield rather than loop.
func (m *MessageCodec) skipOneByte(ctx context.Context, reason string, args ...any) (bool, error) {
	if m.resyncSkippedBytes == 0 {
		m.log.Warn().Msgf("IEC 60870-5-104 stream out of sync ("+reason+"), resynchronizing octet-wise", args...)
	} else {
		m.log.Trace().Msgf("IEC 60870-5-104 resync skipping an octet ("+reason+")", args...)
	}
	if _, err := m.GetTransportInstance().Read(ctx, 1); err != nil {
		m.log.Debug().Err(err).Msg("error discarding an octet during resync")
		if transportError, ok := transports.AsTransportError(err); ok && transportError.Kind() == transports.TransportErrorFatal {
			return false, err
		}
		return false, nil
	}
	m.resyncSkippedBytes++
	return true, nil
}

func (m *MessageCodec) noteResyncComplete() {
	if m.resyncSkippedBytes > 0 {
		m.log.Warn().Uint64("skippedBytes", m.resyncSkippedBytes).Msg("IEC 60870-5-104 stream resynchronized")
		m.resyncSkippedBytes = 0
	}
}
