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

	if err := ti.FillBuffer(ctx, func(pos uint, currentByte byte, reader transports.ExtendedReader) bool {
		m.log.Trace().Uint("pos", pos).Uint8("currentByte", currentByte).Msg("filling")
		numBytesAvailable, err := ti.GetNumBytesAvailableInBuffer()
		if err != nil {
			m.log.Debug().Err(err).Msg("error getting available bytes")
			return false
		}
		m.log.Trace().Uint32("numBytesAvailable", numBytesAvailable).Msg("check available bytes < 6")
		return numBytesAvailable < 6
	}); err != nil {
		m.log.Debug().Err(err).Msg("error filling buffer")
	}

	// We need at least 6 bytes in order to know how big the packet is in total
	num, err := ti.GetNumBytesAvailableInBuffer()
	if err == nil {
		if num < 6 {
			// Don't have enough data yet...
			return nil, nil
		}
		m.log.Trace().Uint32("num", num).Msgf("we got %d readable bytes", num)

		// Peek the MBAP header to extract the txcnid, protocol, and length
		header, err := ti.PeekReadableBytes(ctx, 6)
		if err != nil {
			m.log.Warn().Err(err).Msg("error peeking header")
			return nil, nil
		}

		// Modbus TCP Protocol ID (bytes 2 and 3) must be 0x0000.
		// If it is not, we are desynchronized. Discard 1 byte and retry to find alignment.
		if header[2] != 0x00 || header[3] != 0x00 {
			m.log.Warn().
				Hex("proto", header[2:4]).
				Msg("Invalid Protocol ID (expected 0). Stream desynchronized. Discarding 1 byte to realign.")

			// Burn 1 byte to shift the window and try again next cycle
			ti.Read(ctx, 1)
			return nil, nil
		}

		// Interpret the length field (big endian) to determine the full packet size
		packetSize := (uint32(header[4]) << 8) + uint32(header[5]) + 6
		if num < packetSize {
			var peekedBytes []byte
			if m.log.Debug().Enabled() {
				peekedBytes, _ = ti.PeekReadableBytes(ctx, num)
			}
			m.log.Debug().
				Stringer("currentData", utils.Base64Stringer(peekedBytes)).
				Uint32("num", num).
				Uint32("packetSize", packetSize).Msgf("Not enough bytes. Got: %d Need: %d. Waiting for more data...", num, packetSize)
			return nil, nil
		}

		// Peek read the entire frame
		frameSlice, err := ti.PeekReadableBytes(ctx, packetSize)
		if err != nil {
			m.log.Warn().Err(err).Uint32("packetSize", packetSize).Msg("error peeking packet")
			return nil, nil
		}

		// Parse the frame
		ctxForModel := options.GetLoggerContextForModel(ctx, m.log, options.WithPassLoggerToModel(m.passLogToModel))
		tcpAdu, err := model.ModbusADUParse[model.ModbusTcpADU](ctxForModel, frameSlice, model.DriverType_MODBUS_TCP, true)
		if err != nil {
			if errors.Is(err, io.EOF) {
				// Incomplete frame, keep waiting for the remaining bytes.
				m.log.Debug().Uint32("packetSize", packetSize).Stringer("data", utils.Base64Stringer(frameSlice)).Msg("partial packet detected, awaiting more data")
				return nil, nil
			}

			// Did we perhaps only read part of a bigger frame?
			if num > packetSize {
				// Try reading everything
				extendedSlice, extendedErr := ti.PeekReadableBytes(ctx, num)
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

						// What is the actual ADU size parsed?
						actualSize := uint32(extendedAdu.GetLengthInBytes(ctxForModel))
						if actualSize > num {
							// We still need more data
							m.log.Trace().
								Uint32("available", num).
								Uint32("required", actualSize).
								Msg("extended parsed frame requires more bytes, awaiting more data")
							return nil, nil
						}

						// Consume the actual size in bytes from the buffer (all peeks to here)
						if _, consumeErr := ti.Read(ctx, actualSize); consumeErr != nil {
							m.log.Debug().Err(consumeErr).Uint32("actualSize", actualSize).Msg("error consuming parsed frame")
							return nil, nil
						}

						// Success
						m.log.Warn().
							Uint32("reportedSize", packetSize).
							Uint32("actualSize", actualSize).
							Stringer("extendedData", utils.Base64Stringer(extendedSlice)).
							Stringer("consumedData", utils.Base64Stringer(extendedSlice[:actualSize])).
							Msg("consumed extended frame with corrected size")

						// Check for trailing CRC garbage
						m.trailingCRCGarbageCheck(ctx, ti, extendedSlice[:actualSize])

						return extendedAdu, nil
					}
					if errors.Is(extendedParseErr, io.EOF) {
						m.log.Trace().Uint32("buffered", num).Msg("extended parse incomplete, awaiting more data")
						return nil, nil
					}
				}
			}

			// Seems unparsable - log and discard
			m.log.Warn().Err(err).
				Stringer("data", utils.Base64Stringer(frameSlice)).
				Uint32("packetSize", packetSize).
				Msg("error parsing frame, discarding")

			// Discard the unparsable frame from the buffer
			if _, discardErr := ti.Read(ctx, packetSize); discardErr != nil {
				m.log.Debug().Err(discardErr).Uint32("packetSize", packetSize).Msg("error discarding unparsable frame")
			}

			// Check if the unparsable junk also had a CRC tail
			m.trailingCRCGarbageCheck(ctx, ti, frameSlice)

			return nil, nil
		}

		// First parse attempt successful, consume the bytes from the buffer
		if _, consumeErr := ti.Read(ctx, packetSize); consumeErr != nil {
			m.log.Debug().Err(consumeErr).Uint32("packetSize", packetSize).Msg("error consuming parsed frame")
			return nil, nil
		}

		// Check for trailing CRC garbage
		m.trailingCRCGarbageCheck(ctx, ti, frameSlice)

		return tcpAdu, nil
	}

	m.log.Warn().Err(err).Msg("Got error reading")
	return nil, nil
}

// trailingCRCGarbageCheck checks for trailing Modbus RTU CRC bytes that may have leaked
// into the TCP stream, and discards them to maintain synchronization.
func (m *MessageCodec) trailingCRCGarbageCheck(ctx context.Context, ti transports.TransportInstance, consumedBytes []byte) {
	// We only check if we have consumed at least 7 bytes (Header 6 + UnitID 1)
	// and if there are at least 2 bytes waiting in the buffer.
	if len(consumedBytes) > 6 {
		bytesAvailable, err := ti.GetNumBytesAvailableInBuffer()
		if err == nil && bytesAvailable >= 2 {
			// Peek the next 2 bytes
			nextTwoBytes, peekErr := ti.PeekReadableBytes(ctx, 2)
			if peekErr == nil {
				// Calculate what the CRC *would* be for the RTU portion of the frame we just read.
				// RTU Frame = [UnitID] + [PDU]
				// This corresponds to consumedBytes[6:] (Skipping 6-byte MBAP header)
				rtuPayload := consumedBytes[6:]
				expectedCRC := m.calculateModbusCRC(rtuPayload)

				// Compare: If the next 2 bytes in the buffer match the calculated CRC, it's garbage.
				if nextTwoBytes[0] == expectedCRC[0] && nextTwoBytes[1] == expectedCRC[1] {
					m.log.Warn().
						Stringer("data", utils.Base64Stringer(consumedBytes)).
						Hex("crc", nextTwoBytes).
						Msg("Detected leaked Modbus RTU CRC at end of TCP frame. Discarding to maintain sync.")

					// Discard the 2 CRC bytes
					if _, discardErr := ti.Read(ctx, 2); discardErr != nil {
						m.log.Warn().Err(discardErr).Msg("Error discarding CRC bytes")
					}
				}
			}
		}
	}
}

func (m *MessageCodec) calculateModbusCRC(data []byte) []byte {
	crc := uint16(0xFFFF)
	for _, b := range data {
		crc ^= uint16(b)
		for i := 0; i < 8; i++ {
			if (crc & 0x0001) != 0 {
				crc = (crc >> 1) ^ 0xA001
			} else {
				crc >>= 1
			}
		}
	}
	// Return as Little Endian (Low Byte, High Byte) as per Modbus RTU
	return []byte{uint8(crc & 0xFF), uint8(crc >> 8)}
}
