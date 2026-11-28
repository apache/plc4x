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

// Hand-written helpers called by the generated Profinet code (STATIC_CALL in
// the mspec). Ported from plc4j's profinet StaticHelper; differences are noted
// where Java's behaviour was a bug.

import (
	"bytes"
	"context"
	"encoding/binary"

	"github.com/apache/plc4x/plc4go/spi/utils"
)

// ArrayLength is the length of a byte array.
func ArrayLength(_ context.Context, arr []byte) int {
	return len(arr)
}

// onesComplementSum adds 16 bit words with end-around carry.
func onesComplementSum(sum uint32, words ...uint32) uint32 {
	for _, w := range words {
		sum += w
		for sum > 0xFFFF {
			sum = (sum & 0xFFFF) + 1
		}
	}
	return sum
}

func addressWords(address IpAddress) (uint32, uint32) {
	data := address.GetData()
	if len(data) < 4 {
		return 0, 0
	}
	return uint32(data[0])<<8 | uint32(data[1]), uint32(data[2])<<8 | uint32(data[3])
}

// CalculateIPv4Checksum is the header checksum of the IPv4 header the Profinet
// UDP frames carry (version 4, header length 5, don't fragment, protocol UDP).
func CalculateIPv4Checksum(_ context.Context, totalLength uint16, identification uint16, timeToLive uint8, sourceAddress IpAddress, destinationAddress IpAddress) uint16 {
	src0, src1 := addressWords(sourceAddress)
	dst0, dst1 := addressWords(destinationAddress)
	sum := onesComplementSum(0,
		0x4500,                     // version and header length
		uint32(totalLength),        //
		uint32(identification),     //
		0x4000,                     // flags (don't fragment) and fragment offset
		uint32(timeToLive)<<8|0x11, // time to live and protocol (UDP)
		0x0000,                     // the checksum itself, zero while computing
		src0, src1, dst0, dst1,     //
	)
	return ^uint16(sum)
}

// CalculateUdpChecksum is the UDP checksum over the pseudo header, the UDP
// header and the serialized DCE/RPC payload (padded to an even length).
// Like Java, it returns 0 (no checksum) if the payload cannot be serialized.
func CalculateUdpChecksum(ctx context.Context, sourceAddress IpAddress, destinationAddress IpAddress, sourcePort uint16, destPort uint16, packetLength uint16, payload DceRpc_Packet) uint16 {
	src0, src1 := addressWords(sourceAddress)
	dst0, dst1 := addressWords(destinationAddress)
	sum := onesComplementSum(0,
		src0, src1, dst0, dst1,
		0x0011, // protocol (UDP)
		uint32(packetLength),
		uint32(sourcePort),
		uint32(destPort),
		uint32(packetLength),
		0x0000, // the checksum itself, zero while computing
	)
	buffer := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := payload.SerializeWithWriteBuffer(ctx, buffer); err != nil {
		return 0x0000
	}
	data := buffer.GetBytes()
	if len(data)%2 != 0 {
		data = append(data, 0x00)
	}
	for i := 0; i+1 < len(data); i += 2 {
		sum = onesComplementSum(sum, uint32(data[i])<<8|uint32(data[i+1]))
	}
	return ^uint16(sum)
}

// GetFrameId maps a frame id value to its range (IEC 61158-6-10 frame id table).
// Java compares RT class 1 against decimal 8000; the range starts at 0x8000.
func GetFrameId(_ context.Context, frameIdValue uint16) PnDcp_FrameId {
	switch v := frameIdValue; {
	case v == 0x0020:
		return PnDcp_FrameId_PTCP_RTSyncPDUWithFollowUp
	case v == 0x0080:
		return PnDcp_FrameId_PTCP_RTSyncPDU
	case v >= 0x0100 && v <= 0x0FFF:
		return PnDcp_FrameId_RT_CLASS_3
	case v >= 0x8000 && v <= 0xBFFF:
		return PnDcp_FrameId_RT_CLASS_1
	case v >= 0xC000 && v <= 0xFBFF:
		return PnDcp_FrameId_RT_CLASS_UDP
	case v == 0xFC01:
		return PnDcp_FrameId_Alarm_High
	case v == 0xFE01:
		return PnDcp_FrameId_Alarm_Low
	case v == 0xFEFC:
		return PnDcp_FrameId_DCP_Hello_ReqPDU
	case v == 0xFEFD:
		return PnDcp_FrameId_DCP_GetSet_PDU
	case v == 0xFEFE:
		return PnDcp_FrameId_DCP_Identify_ReqPDU
	case v == 0xFEFF:
		return PnDcp_FrameId_DCP_Identify_ResPDU
	case v == 0xFF00:
		return PnDcp_FrameId_PTCP_AnnouncePDU
	case v == 0xFF20:
		return PnDcp_FrameId_PTCP_FollowUpPDU
	case v == 0xFF40:
		return PnDcp_FrameId_PTCP_DelayReqPDU
	case v == 0xFF41:
		return PnDcp_FrameId_PTCP_DelayResPDUWithFollowUp
	case v == 0xFF42:
		return PnDcp_FrameId_PTCP_DelayFuResPDUWithFollowUp
	case v == 0xFF43:
		return PnDcp_FrameId_PTCP_DelayResPDUWithoutFollowUp
	case v >= 0xFF80 && v <= 0xFF8F:
		return PnDcp_FrameId_FragmentationFrameId
	}
	return PnDcp_FrameId_RESERVED
}

// IsSysexEnd ends the LLDP TLV list at the End of LLDPDU TLV (second byte 0),
// or at the end of the buffer (peeking past it reads 0).
func IsSysexEnd(_ context.Context, io utils.ReadBuffer) func([]LldpUnit) bool {
	return func([]LldpUnit) bool {
		rb, ok := io.(utils.ReadBufferByteBased)
		if !ok {
			return false
		}
		return rb.PeekByte(1) == 0x00
	}
}

// ParseSysexString parses one LLDP TLV.
func ParseSysexString(_ context.Context, io utils.ReadBuffer) func(context.Context) (LldpUnit, error) {
	return func(ctx context.Context) (LldpUnit, error) {
		return LldpUnitParseWithBuffer[LldpUnit](ctx, io)
	}
}

// SerializeSysexString writes one LLDP TLV.
func SerializeSysexString(ctx context.Context, io utils.WriteBuffer, unit LldpUnit) error {
	return unit.SerializeWithWriteBuffer(ctx, io)
}

// LengthSysexString is the length of the TLV list in bits: the generated code
// adds it to lengthInBits. (Java returns bytes there.)
func LengthSysexString(ctx context.Context, data []LldpUnit) uint64 {
	var lengthInBits uint64
	for _, unit := range data {
		lengthInBits += unit.GetLengthInBits(ctx)
	}
	return lengthInBits
}

// WriteDataUnit writes the cyclic data unit.
func WriteDataUnit(ctx context.Context, writeBuffer utils.WriteBuffer, dataUnit PnIo_CyclicServiceDataUnit) error {
	return dataUnit.SerializeWithWriteBuffer(ctx, writeBuffer)
}

// cyclicTrailerBytes is the trailer of a cyclic frame after the data unit:
// cycle counter (2), data status (1) and transfer status (1).
const cyclicTrailerBytes = 4

// ReadDataUnit reads the cyclic data unit: every byte up to the 4 byte trailer.
func ReadDataUnit(ctx context.Context, readBuffer utils.ReadBuffer) (PnIo_CyclicServiceDataUnit, error) {
	rb, ok := readBuffer.(utils.ReadBufferByteBased)
	if !ok {
		return PnIo_CyclicServiceDataUnitParseWithBuffer(ctx, readBuffer, 0)
	}
	remaining := int64(rb.GetTotalBytes()) - int64(rb.GetPos())
	length := remaining - cyclicTrailerBytes
	if length < 0 {
		length = 0
	}
	return PnIo_CyclicServiceDataUnitParseWithBuffer(ctx, readBuffer, int16(length))
}

var nullUuid = make([]byte, 16)

// IsNullUuid reports whether a UUID is all zeros.
func IsNullUuid(_ context.Context, uuid Uuid) bool {
	return bytes.Equal(uuid.GetData(), nullUuid)
}
