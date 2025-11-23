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

package modbus

import (
	"context"
	"fmt"

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
		m.log.Debug().Err(err).Msg("error filling buffer")
	}

	// 2. Check buffer status
	numBytesAvail, err := ti.GetNumBytesAvailableInBuffer()
	if err != nil {
		m.log.Warn().Err(err).Msg("Error getting bytes in buffer")
		return nil, nil
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
				return nil, nil
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

	// --- CHECK 1: Protocol ID Validation ---
	if header[2] != 0x00 || header[3] != 0x00 {
		return m.handleDesync(ctx, "Invalid Protocol ID", map[string]interface{}{
			"p1":   header[2],
			"p2":   header[3],
			"data": utils.Base64Stringer(header),
		})
	}

	// Length field is big endian encoded WORD
	payloadLength := (uint32(header[4]) << 8) + uint32(header[5])
	packetSize := payloadLength + 6

	// --- CHECK 2: Minimum Length ---
	if payloadLength < 2 {
		return m.handleDesync(ctx, "Invalid packet length (<2)", map[string]interface{}{
			"len":  payloadLength,
			"data": utils.Base64Stringer(header),
		})
	}

	// --- CHECK 3: Function Code 0 (False Header) ---
	if numBytesAvail >= 8 {
		peekBytes, err := ti.PeekReadableBytes(ctx, 8)
		if err == nil && peekBytes[7] == 0 {
			return m.handleDesync(ctx, "Invalid Function Code (0)", map[string]interface{}{
				"len":  payloadLength,
				"data": utils.Base64Stringer(peekBytes),
			})
		}
	}

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
		m.log.Warn().Err(err).
			Stringer("data", utils.Base64Stringer(frameSlice)).
			Uint32("packetSize", packetSize).
			Msg("Error parsing frame. Discarding invalid frame.")

		// Discard the unparsable frame from the buffer
		if _, discardErr := ti.Read(ctx, packetSize); discardErr != nil {
			m.log.Debug().Err(discardErr).Uint32("packetSize", packetSize).Msg("error discarding unparsable frame")
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

// handleDesync handles stream realignment when an invalid header is detected at the head.
// It strictly scans the available buffer for a valid MBAP header.
// If one is found, it realigns the stream.
// If NO valid header is found in the *entire* buffer, it treats the connection as dead.
func (m *MessageCodec) handleDesync(ctx context.Context, reason string, fields map[string]interface{}) (spi.Message, error) {
	ti := m.GetTransportInstance()

	// Get total available bytes
	numBytesAvail, err := ti.GetNumBytesAvailableInBuffer()
	if err != nil {
		return nil, nil // Yield if we can't check buffer
	}

	// Create a logger with context
	fields["reason"] = reason
	fields["bytesAvail"] = numBytesAvail
	opLog := m.log.With().
		Interface("desyncContext", fields).
		Logger()

	opLog.Warn().Msg("Desync detected at stream head.")

	// CASE 1: Small Buffer (< 10 bytes).
	// We don't have enough data to scan for a full header+function (need ~8-9 bytes min).
	// We can't definitively say the stream is dead, so we just Discard 1 and Yield.
	if numBytesAvail < 10 {
		opLog.Trace().Msg("Small buffer desync. Discarding 1 byte.")
		if _, err := ti.Read(ctx, 1); err != nil {
			opLog.Debug().Err(err).Msg("Error reading byte during discard")
		}
		return nil, nil
	}

	// CASE 2: Scan for Recovery.
	// We peek everything we have.
	allBytes, err := ti.PeekReadableBytes(ctx, numBytesAvail)
	if err != nil {
		return nil, nil
	}

	// Scan Loop: Start at offset 1 (since offset 0 is known bad).
	// We must stop at (num - 8) because we need to inspect the Function Code at (i + 7).
	// Index 'i' is the start of the candidate MBAP Header.
	for i := uint32(1); i <= numBytesAvail-8; i++ {
		// Check Protocol ID (Bytes 2-3 of the candidate header)
		if allBytes[i+2] == 0x00 && allBytes[i+3] == 0x00 {
			// Check Length (Bytes 4-5)
			length := (uint32(allBytes[i+4]) << 8) + uint32(allBytes[i+5])

			// Check Function Code (Byte 7)
			// MBAP(6) + UnitID(1) + FuncCode(1) -> Offset 7
			fc := allBytes[i+7]

			// VALIDATION LOGIC:
			// 1. Length >= 2 (Must have UnitID + FuncCode)
			// 2. FuncCode != 0 (Modbus function 0 doesn't exist)
			// 3. Strict Packet Structure:
			//    - Standard Function (< 0x80): Allow any length.
			//    - Exception Function (>= 0x80): PDU is always 2 bytes (Func + Code).
			//      Modbus TCP Length = UnitID(1) + PDU(2) = 3 bytes.
			if length >= 2 && fc != 0 && (fc < 0x80 || length == 3) {
				opLog.Debug().Uint32("offset", i).Msg("Found MBAP candidate in stream. Discarding garbage prefix.")

				// Discard 'i' bytes to align the stream to this candidate
				if _, err := ti.Read(ctx, i); err != nil {
					opLog.Debug().Err(err).Msg("Error discarding garbage during recovery")
				}

				// Yield. Next Receive() will pick up this valid header.
				return nil, nil
			}
		}
	}

	// CASE 3: Connection is unrecoverable with the tools we have...
	// We scanned the entire available buffer and found NO valid candidates.
	// The connection is sending garbage. Destroy it.
	return nil, fmt.Errorf("stream desynchronized: %d bytes of garbage with no valid MBAP header found", numBytesAvail)
}
