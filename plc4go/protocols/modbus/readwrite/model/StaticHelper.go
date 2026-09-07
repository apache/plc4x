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

package model

import (
	"context"
	"encoding/binary"
	"math/bits"

	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// AsciiLrcCheck computes the longitudinal redundancy check an ASCII frame carries behind its PDU.
// Ported from plc4j StaticHelper.asciiLrcCheck: sum the address and the serialized PDU into an
// 8 bit accumulator (carries discarded) and negate it, which is the two's complement the Modbus
// specification asks for. The hex encoding of the result happens in the codec, not here.
func AsciiLrcCheck(ctx context.Context, address uint8, pdu ModbusPDU) func() (uint8, error) {
	return func() (uint8, error) {
		data, err := checksummedBytes(ctx, address, pdu)
		if err != nil {
			return 0, errors.Wrap(err, "error assembling the bytes to checksum")
		}
		var lrc uint8
		for _, b := range data {
			lrc += b
		}
		return -lrc, nil
	}
}

// checksummedBytes is what both the RTU CRC and the ASCII LRC are computed over: the station
// address followed by the serialized PDU (plc4j StaticHelper.rtuCrcCheck/asciiLrcCheck).
func checksummedBytes(ctx context.Context, address uint8, pdu ModbusPDU) ([]byte, error) {
	if pdu == nil {
		return nil, errors.New("pdu must not be nil")
	}
	wb := utils.NewWriteBufferByteBased(
		utils.WithInitialSizeForByteBasedBuffer(int(pdu.GetLengthInBytes(ctx))+1),
		utils.WithByteOrderForByteBasedBuffer(binary.BigEndian),
	)
	if err := wb.WriteUint8("address", 8, address); err != nil {
		return nil, errors.Wrap(err, "error serializing address")
	}
	if err := pdu.SerializeWithWriteBuffer(ctx, wb); err != nil {
		return nil, errors.Wrap(err, "error serializing pdu")
	}
	return wb.GetBytes(), nil
}

// rtuCrcTable is the CRC-16/MODBUS table (reflected polynomial 0xA001, init 0xFFFF). plc4j keeps
// it as the two byte-wide auchCRCHi/auchCRCLo tables from PI_MBUS_300.pdf page 121; those are just
// the low and high halves of these entries, so one table of shorts computes the same thing.
var rtuCrcTable = func() (table [256]uint16) {
	for i := range table {
		crc := uint16(i)
		for range 8 {
			if crc&1 != 0 {
				crc = (crc >> 1) ^ 0xA001
			} else {
				crc >>= 1
			}
		}
		table[i] = crc
	}
	return
}()

// RtuCrcCheck computes the CRC an RTU frame carries behind its PDU. Modbus transmits the CRC low
// byte first, but the mspec models the field as a big-endian uint16, so the value handed to the
// field is the byte-swapped CRC - which is exactly what plc4j's rtuCrcCheck returns from its
// `(uchCRCHi << 8) | uchCRCLo` (its "Hi" register holds the byte that goes on the wire first).
func RtuCrcCheck(ctx context.Context, address uint8, pdu ModbusPDU) func() (uint16, error) {
	return func() (uint16, error) {
		data, err := checksummedBytes(ctx, address, pdu)
		if err != nil {
			return 0, errors.Wrap(err, "error assembling the bytes to checksum")
		}
		crc := uint16(0xFFFF)
		for _, b := range data {
			crc = (crc >> 8) ^ rtuCrcTable[byte(crc)^b]
		}
		return bits.ReverseBytes16(crc), nil
	}
}
