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
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

//go:generate go tool plc4xGenerator -type=MessageCodec
type MessageCodec struct {
	_default.DefaultCodec
	expectationCounter int32

	passLogToModel bool

	log zerolog.Logger
}

func NewMessageCodec(transportInstance transports.TransportInstance, _options ...options.WithOption) *MessageCodec {
	passLoggerToModel, _ := options.ExtractPassLoggerToModel(_options...)
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	codec := &MessageCodec{
		expectationCounter: 1,
		passLogToModel:     passLoggerToModel,
		log:                customLogger,
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
	tcpAdu := message.(model.ModbusTcpADU)
	// Serialize the request
	theBytes, err := tcpAdu.Serialize()
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
	ti := m.GetTransportInstance()
	if !ti.IsConnected() {
		return nil, errors.New("Transport instance not connected")
	}

	// 1. Fill the buffer
	if err := ti.FillBuffer(ctx, func(pos uint, currentByte byte, reader transports.ExtendedReader) bool {
		m.log.Trace().Uint("pos", pos).Uint8("currentByte", currentByte).Msg("filling")
		numBytesAvailable, err := ti.GetNumBytesAvailableInBuffer()
		if err != nil {
			m.log.Debug().Err(err).Msg("error getting available bytes")
			return false
		}
		return numBytesAvailable < 6
	}); err != nil {
		if err != io.EOF {
			m.log.Debug().Err(err).Msg("error filling buffer")
		}
		// Fall through on errors, we might have enough data...
	}

	// 2. Check buffer status
	numBytesAvail, err := ti.GetNumBytesAvailableInBuffer()
	if err != nil {
		// Yield if we can't check buffer
		if err == io.EOF {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting buffer length")
	}

	// -------------------------------------------------------------------------
	// Discard NIL Packets (Keep-Alives w/padding that leaked into stream)
	// -------------------------------------------------------------------------
	for {
		if numBytesAvail < 6 {
			break
		}
		header, err := ti.PeekReadableBytes(ctx, 6)
		if err != nil {
			m.log.Warn().Err(err).Msg("error peeking header")
			return nil, nil
		}

		// Check for 6 bytes of zeros
		if header[0] == 0 && header[1] == 0 && header[2] == 0 &&
			header[3] == 0 && header[4] == 0 && header[5] == 0 {

			m.log.Debug().Msg("Detected NIL Packet (Keep-Alive). Discarding 6 bytes.")
			if _, err := ti.Read(ctx, 6); err != nil {
				m.log.Warn().Err(err).Msg("Error discarding NIL packet")
				// If we can't read, we are dead.
				return nil, err
			}

			// Refresh num and loop
			numBytesAvail, err = ti.GetNumBytesAvailableInBuffer()
			if err != nil {
				return nil, nil
			}
			continue
		}
		break
	}

	if numBytesAvail < 6 {
		return nil, nil
	}

	// Re-peek the header at the current head
	header, err := ti.PeekReadableBytes(ctx, 6)
	if err != nil {
		m.log.Warn().Err(err).Msg("error peeking header")
		return nil, nil
	}

	// -------------------------------------------------------------------------
	// MBAP SANITY CHECKS
	// -------------------------------------------------------------------------

	// We now peek up to 9 bytes (or whatever is available) to perform the full Consistency Check.
	// IMPORTANT: If we have 8 bytes (Header + Unit + Func), we MUST pass all 8 bytes
	// so that checkPacketConsistency can validate the Function Code and Length limits.
	// If we only pass 6, it will skip the Function check and potentially accept huge lengths.
	var checkBytes []byte
	peekLen := uint32(9)
	if numBytesAvail < 9 {
		peekLen = numBytesAvail
	}
	checkBytes, _ = ti.PeekReadableBytes(ctx, peekLen)

	// Perform Consistency Check
	if !m.checkPacketConsistency(checkBytes) {
		return m.handleDesync(ctx, "Sanity Check Failed", map[string]interface{}{
			"data": utils.Base64Stringer(checkBytes),
		})
	}

	// Length field is big endian encoded WORD
	payloadLength := (uint32(header[4]) << 8) + uint32(header[5])
	packetSize := payloadLength + 6

	// -------------------------------------------------------------------------
	// PARSING
	// -------------------------------------------------------------------------

	// Yield on TCP fragmentation
	if numBytesAvail < packetSize {
		// Wait for more data (standard TCP fragmentation handling)
		var peekedBytes []byte
		if m.log.Debug().Enabled() {
			peekedBytes, _ = ti.PeekReadableBytes(ctx, numBytesAvail)
		}
		m.log.Debug().
			Stringer("dataFragment", utils.Base64Stringer(peekedBytes)).
			Uint32("num", numBytesAvail).
			Uint32("packetSize", packetSize).Msgf("Received fragment. Got: %d Need: %d. Waiting for more data...", numBytesAvail, packetSize)
		return nil, nil
	}

	// Read the entire frame
	frameSlice, err := ti.PeekReadableBytes(ctx, packetSize)
	if err != nil {
		m.log.Warn().Err(err).Msg("Error peeking frame slice")
		return nil, nil
	}

	// Parse the frame
	ctxForModel := options.GetLoggerContextForModel(ctx, m.log, options.WithPassLoggerToModel(m.passLogToModel))
	tcpAdu, err := model.ModbusADUParse[model.ModbusTcpADU](ctxForModel, frameSlice, model.DriverType_MODBUS_TCP, true)
	if err != nil {
		// Parser wasn't happy at packetSize, if there is more available, try parsing all of it
		if numBytesAvail > packetSize {
			extendedSlice, extendedErr := ti.PeekReadableBytes(ctx, numBytesAvail)
			if extendedErr == nil {
				// Try parsing it all
				extendedAdu, extendedParseErr := model.ModbusADUParse[model.ModbusTcpADU](ctxForModel, extendedSlice, model.DriverType_MODBUS_TCP, true)
				if extendedParseErr == nil {
					// Parse succeded...
					// NOTE: MBAP length errors are rare in well-behaved devices, but they happen when
					// firmware miscomputes the length field (or hard-codes to 6), mixes RTU framing
					// with TCP transport, or when intermediate gateways/proxies truncate or merge
					// frames. Duplicate reads and out-of-order slicing in buggy drivers can also leave
					// stale data that makes the reported length disagree with the actual payload.
					// ... looking at you old Moxa and Lantronix "transparent" serial-to-tcp gateways.

					// What is the actual ADU size parsed?
					actualSize := uint32(extendedAdu.GetLengthInBytes(ctxForModel))
					if actualSize > numBytesAvail {
						return nil, nil
					}

					m.log.Info().
						Uint32("reportedSize", packetSize).
						Uint32("actualSize", actualSize).
						Stringer("extendedData", utils.Base64Stringer(extendedSlice)).
						Stringer("consumedData", utils.Base64Stringer(extendedSlice[:actualSize])).
						Msg("MBAP had wrong/hardcoded length. Consumed extended frame with corrected size")
					if _, consumeErr := ti.Read(ctx, actualSize); consumeErr != nil {
						m.log.Debug().Err(consumeErr).Msg("error consuming extended frame")
						return nil, nil
					}
					return extendedAdu, nil
				}
			}
		}

		// Seems unparsable - log and discard
		// NOTE: If checkPacketConsistency PASSED, but Parse FAILED, it means the header
		// structure is valid (Length, Proto, Func), but the Content is bad.
		// We MUST discard this frame to advance the stream.
		// We DO NOT call handleDesync here because we don't want to scan ahead;
		// we just want to eat the bad packet and try the next one.
		m.log.Warn().Err(err).
			Stringer("data", utils.Base64Stringer(frameSlice)).
			Uint32("packetSize", packetSize).
			Msg("Error parsing frame. Discarding invalid frame.")

		// Discard the unparsable frame from the buffer
		if _, discardErr := ti.Read(ctx, packetSize); discardErr != nil {
			m.log.Debug().Err(discardErr).Uint32("packetSize", packetSize).Msg("error discarding unparsable frame")
			return nil, discardErr
		}

		// Yield
		return nil, nil
	}

	// -------------------------------------------------------------------------
	// SUCCESS
	// -------------------------------------------------------------------------
	if _, consumeErr := ti.Read(ctx, packetSize); consumeErr != nil {
		m.log.Debug().Err(consumeErr).Msg("error consuming parsed frame")
		return nil, nil
	}

	return tcpAdu, nil
}

// checkPacketConsistency validates a candidate Modbus TCP packet against the protocol spec.
// It enforces strict relationship rules between the MBAP Header Length and the PDU content (Function Code & Byte Count).
// Returns true if the packet looks valid, false if it is definitely garbage.
// NOTE: When new functions are added to the mspec, this needs to be updated to reflect their structure.
func (m *MessageCodec) checkPacketConsistency(data []byte) bool {
	// 1. Check Protocol ID immediately (Need 4 bytes)
	if len(data) < 4 {
		return false
	}
	// Modbus TCP Protocol ID must be 0
	if data[2] != 0x00 || data[3] != 0x00 {
		return false
	}

	// 2. Check Header Length (Need 6 bytes)
	if len(data) < 6 {
		// Not enough data to check length, but Protocol ID was OK.
		// We shouldn't reject yet if we just don't have the bytes.
		// However, this function is typically called with a slice that *should* contain the header.
		return false
	}

	length := (uint32(data[4]) << 8) + uint32(data[5])
	if length < 2 {
		return false
	}

	// 3. Check Function Code (Need 8 bytes: 6 Header + 1 Unit + 1 Func)
	if len(data) < 8 {
		// We have a valid ProtoID and Length, but not enough data to check content.
		// Return true (benefit of the doubt) because we can't prove it's bad yet.
		return true
	}

	fc := data[7]
	if fc == 0 {
		return false
	}

	// VALIDATION: Check Internal Consistency of Response Structure
	// NOTE: These rules are primarily for Modbus TCP RESPONSES (Client Mode).
	// If acting as a Server (receiving Requests), some rules (like Write Multiple) would be different.

	if fc >= 0x80 {
		// --- EXCEPTION ---
		// Structure: [UnitID] [FuncCode+0x80] [ExceptionCode]
		// Length must be exactly 3.
		return length == 3
	}

	// --- STANDARD FUNCTIONS ---
	switch fc {
	case 0x01, 0x02, 0x03, 0x04, 0x17, 0x14, 0x15:
		// Variable Length Responses
		// Rule: Header Length == ByteCount + 3
		// We need 9 bytes to see the ByteCount at offset 8.
		if len(data) >= 9 {
			byteCount := uint32(data[8])
			if length != byteCount+3 {
				return false
			}
		}
		// Implicit Max Check: Max ByteCount 255 -> Max Length 258.
		if length > 258 {
			return false
		}

	case 0x05, 0x06, 0x0F, 0x10:
		// Fixed Length Responses (Write Single, Write Multi)
		// Structure: [UnitID] [FuncCode] [AddrHi] [AddrLo] [ValHi] [ValLo]
		// Length must be exactly 6.
		if length != 6 {
			// SERVER MODE CAVEAT:
			// If we are a Server receiving a Write Multiple Request (FC 15/16),
			// the length will be > 6. If we strictly return false here, we break Server mode.
			// As a heuristic, if length is > 6 for these codes, we treat it as potentially valid
			// (assuming it's a Request) to be safe, unless it's huge.
			if length > 260 {
				return false
			}
			// Ideally we would enforce Request structure (Len = ByteCount + 7),
			// but that requires more bytes than we might have peeked.
			return true
		}

	case 0x16:
		// Mask Write Register
		// Length must be exactly 8.
		if length != 8 {
			return false
		}

	default:
		// Other/Custom Functions.
		// Max PDU size is 253 -> Length 254. Allow margin.
		if length > 260 {
			return false
		}
	}

	return true
}

// handleDesync handles stream realignment when an invalid header is detected at the head.
// It strictly scans the available buffer for a valid MBAP header using checkPacketConsistency.
// If one is found, it realigns the stream.
// If NO valid header is found in the *entire* buffer, it treats the connection as dead.
func (m *MessageCodec) handleDesync(ctx context.Context, reason string, fields map[string]interface{}) (spi.Message, error) {
	ti := m.GetTransportInstance()

	// Get total available bytes
	numBytesAvail, err := ti.GetNumBytesAvailableInBuffer()
	if err != nil {
		return nil, err // Yield if we can't check buffer
	}

	// Create a logger with context
	fields["reason"] = reason
	fields["bytesAvail"] = numBytesAvail
	opLog := m.log.With().
		Interface("desyncContext", fields).
		Logger()

	// Log on the way out to minimize log spam
	dispMsg := ""
	dispLevel := zerolog.WarnLevel
	defer func() {
		opLog.WithLevel(dispLevel).Msg("Desync detected at stream head: " + dispMsg)
	}()

	// CASE 1: Small Buffer (< 10 bytes).
	// We don't have enough data to scan for a full header+function (need ~8-9 bytes min).
	// We can't definitively say the stream is dead, so we just Discard 1 and Yield.
	if numBytesAvail < 10 {
		dispMsg = "Small buffer - Discard 1"
		dispLevel = zerolog.DebugLevel
		if _, err := ti.Read(ctx, 1); err != nil {
			opLog.Debug().Err(err).Msg("Error reading byte during discard")
			return nil, err // Return error to kill connection if we can't consume
		}
		return nil, nil
	}

	// CASE 2: Scan for Recovery.
	// We peek everything we have.
	allBytes, err := ti.PeekReadableBytes(ctx, numBytesAvail)
	if err != nil {
		return nil, fmt.Errorf("error peeking bytes during desync recovery")
	}

	// Scan Loop: Start at offset 1 (since offset 0 is known bad).
	// We verify candidates up to the end of the buffer.
	// We stop when we don't have enough bytes left to even check the Protocol ID (4 bytes).
	limit := uint32(0)
	if numBytesAvail >= 4 {
		limit = numBytesAvail - 4
	}

	for i := uint32(1); i <= limit; i++ {
		// Use our robust consistency check on the slice starting at i
		if m.checkPacketConsistency(allBytes[i:]) {
			dispMsg = fmt.Sprintf("Found MBAP candidate at +%d", i)
			dispLevel = zerolog.DebugLevel

			// Discard 'i' bytes to align the stream to this candidate
			if _, err := ti.Read(ctx, i); err != nil {
				opLog.Debug().Err(err).Msg("Error discarding garbage during recovery")
				return nil, err // Return error to kill connection
			}

			// Yield. Next Receive() will pick up this valid header.
			return nil, nil
		}
	}

	// CASE 3: Connection is unrecoverable with the tools we have...
	// We scanned the entire available buffer and found NO valid candidates.
	// The connection is sending garbage. Would be nice to destroy it, but DefaultCodec
	// doesn't have any such ability. For now we'll just consume all available bytes
	// and return to make sure we don't spin endlessly.
	dispMsg = fmt.Sprintf("No MBAP candidate found. Discarding all (%d) available bytes", numBytesAvail)
	dispLevel = zerolog.InfoLevel
	if _, err := ti.Read(ctx, numBytesAvail); err != nil {
		opLog.Debug().Err(err).Msg("Error discarding garbage during recovery")
		return nil, err // Return error to kill connection
	}

	return nil, fmt.Errorf("stream desynchronized: %d bytes of garbage with no valid MBAP header found", numBytesAvail)
}
