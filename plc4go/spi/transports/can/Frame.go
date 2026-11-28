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

// Package can holds what the CAN transports have in common: the fixed 16 byte frame layout they
// hand to a driver.
//
// A CAN transport delivers its frames as a plain byte stream, one fixed-size record after the other,
// so a driver can read them through the ordinary transports.TransportInstance interface:
//
//	offset 0  4 bytes  CAN identifier, big-endian, with the flag bits below in the top three bits
//	offset 4  1 byte   data length code (0..8)
//	offset 5  3 bytes  reserved, zero
//	offset 8  8 bytes  data, zero padded
//
// This is the layout of Linux's struct can_frame with the identifier in network byte order rather
// than host order, which is also what plc4j's CAN transports deliver to its CANopen codec. The
// SocketCAN transport converts the identifier at the socket boundary.
package can

import (
	"encoding/binary"
	"fmt"

	"github.com/apache/plc4x/plc4go/spi/errors"
)

const (
	// FrameSize is the size of one frame record on the stream.
	FrameSize = 16
	// MaxDataLength is the most data a classic CAN frame carries.
	MaxDataLength = 8

	// FlagExtended marks a 29 bit identifier (Linux CAN_EFF_FLAG).
	FlagExtended = uint32(0x80000000)
	// FlagRemote marks a remote transmission request (Linux CAN_RTR_FLAG).
	FlagRemote = uint32(0x40000000)
	// FlagError marks an error frame (Linux CAN_ERR_FLAG).
	FlagError = uint32(0x20000000)

	// StandardIDMask is the 11 bit identifier of a standard frame.
	StandardIDMask = uint32(0x7FF)
	// ExtendedIDMask is the 29 bit identifier of an extended frame.
	ExtendedIDMask = uint32(0x1FFFFFFF)
)

// Frame is one classic CAN frame.
type Frame struct {
	// ID is the identifier without flag bits.
	ID       uint32
	Extended bool
	Remote   bool
	Error    bool
	// Data is at most MaxDataLength bytes. A remote frame carries none.
	Data []byte
	// RequestedLength is the data length code of a remote frame: how many bytes it asks for
	// (0..MaxDataLength). It is used only when Remote is set; a data frame's length is len(Data).
	RequestedLength uint8
}

// Encode lays the frame out as one FrameSize record.
func (f Frame) Encode() ([]byte, error) {
	if len(f.Data) > MaxDataLength {
		return nil, errors.Errorf("a CAN frame carries at most %d data bytes, got %d", MaxDataLength, len(f.Data))
	}
	rawID := f.ID
	if f.Extended {
		if rawID > ExtendedIDMask {
			return nil, errors.Errorf("extended CAN identifier 0x%X is wider than 29 bits", rawID)
		}
		rawID |= FlagExtended
	} else if rawID > StandardIDMask {
		return nil, errors.Errorf("standard CAN identifier 0x%X is wider than 11 bits", rawID)
	}
	dlc := byte(len(f.Data))
	if f.Remote {
		if len(f.Data) != 0 {
			return nil, errors.Errorf("a remote CAN frame carries no data, got %d bytes", len(f.Data))
		}
		if f.RequestedLength > MaxDataLength {
			return nil, errors.Errorf("a remote CAN frame asks for at most %d bytes, got %d", MaxDataLength, f.RequestedLength)
		}
		dlc = f.RequestedLength
		rawID |= FlagRemote
	}
	if f.Error {
		rawID |= FlagError
	}
	record := make([]byte, FrameSize)
	binary.BigEndian.PutUint32(record[0:4], rawID)
	record[4] = dlc
	copy(record[8:], f.Data)
	return record, nil
}

// Decode reads one FrameSize record. A data length code above 8 (CAN FD lengths, or garbage) is an
// error rather than being clipped, so a stream that lost its alignment shows up instead of turning
// into plausible-looking frames.
func Decode(record []byte) (Frame, error) {
	if len(record) < FrameSize {
		return Frame{}, errors.Errorf("a CAN frame record is %d bytes, got %d", FrameSize, len(record))
	}
	rawID := binary.BigEndian.Uint32(record[0:4])
	dlc := int(record[4])
	if dlc > MaxDataLength {
		return Frame{}, errors.Errorf("CAN data length code %d is above %d", dlc, MaxDataLength)
	}
	frame := Frame{
		Extended: rawID&FlagExtended != 0,
		Remote:   rawID&FlagRemote != 0,
		Error:    rawID&FlagError != 0,
	}
	if frame.Remote {
		// The data length code of a remote frame is the length it asks for; the data bytes are padding.
		frame.RequestedLength = uint8(dlc)
		frame.Data = []byte{}
	} else {
		frame.Data = append([]byte(nil), record[8:8+dlc]...)
	}
	if frame.Extended {
		frame.ID = rawID & ExtendedIDMask
	} else {
		frame.ID = rawID & StandardIDMask
	}
	return frame, nil
}

func (f Frame) String() string {
	kind := ""
	switch {
	case f.Error:
		kind = " error"
	case f.Remote:
		kind = fmt.Sprintf(" rtr, dlc: %d", f.RequestedLength)
	}
	if f.Extended {
		return fmt.Sprintf("can{id: 0x%08X%s, data: % X}", f.ID, kind, f.Data)
	}
	return fmt.Sprintf("can{id: 0x%03X%s, data: % X}", f.ID, kind, f.Data)
}
