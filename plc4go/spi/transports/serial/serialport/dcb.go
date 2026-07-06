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

package serialport

// Windows DCB field values, kept in a portable file (no x/sys/windows
// types) so the mapping is unit-testable on every development platform.
// Values per Microsoft's documentation of the DCB structure (winbase.h).
const (
	dcbNoParity    = 0 // NOPARITY
	dcbOddParity   = 1 // ODDPARITY
	dcbEvenParity  = 2 // EVENPARITY
	dcbMarkParity  = 3 // MARKPARITY
	dcbSpaceParity = 4 // SPACEPARITY

	dcbOneStopBit   = 0 // ONESTOPBIT
	dcbOne5StopBits = 1 // ONE5STOPBITS
	dcbTwoStopBits  = 2 // TWOSTOPBITS

	// Bit positions inside DCB.Flags (the C bitfield block starting at fBinary).
	dcbFlagBinary              = 0x00000001 // fBinary: binary mode, required
	dcbFlagParity              = 0x00000002 // fParity: enable parity checking
	dcbFlagOutxCtsFlow         = 0x00000004 // fOutxCtsFlow: output suspended while CTS low
	dcbFlagDtrControlEnable    = 0x00000010 // fDtrControl = DTR_CONTROL_ENABLE (1 << 4)
	dcbFlagRtsControlEnable    = 0x00001000 // fRtsControl = RTS_CONTROL_ENABLE (1 << 12)
	dcbFlagRtsControlHandshake = 0x00002000 // fRtsControl = RTS_CONTROL_HANDSHAKE (2 << 12)
	dcbFlagOutX                = 0x00000100 // fOutX: honor XOFF/XON received from the device
	dcbFlagInX                 = 0x00000200 // fInX: send XOFF/XON when the RX buffer fills/drains
)

// dcbSettings carries the computed DCB field values; port_windows.go copies
// them into a windows.DCB.
type dcbSettings struct {
	BaudRate uint32
	ByteSize uint8
	Parity   uint8
	StopBits uint8
	Flags    uint32
	XonLim   uint16
	XoffLim  uint16
	XonChar  byte
	XoffChar byte
}

// makeDCBSettings translates an already-normalized Config into DCB field
// values. DTR is asserted on open (DTR_CONTROL_ENABLE); RTS is asserted
// unless hardware flow control hands it to the driver.
func makeDCBSettings(cfg Config) dcbSettings {
	s := dcbSettings{
		BaudRate: uint32(cfg.BaudRate),
		ByteSize: uint8(cfg.DataBits),
		Flags:    dcbFlagBinary | dcbFlagDtrControlEnable,
	}
	switch cfg.Parity {
	case ParityOdd:
		s.Parity = dcbOddParity
		s.Flags |= dcbFlagParity
	case ParityEven:
		s.Parity = dcbEvenParity
		s.Flags |= dcbFlagParity
	case ParityMark:
		s.Parity = dcbMarkParity
		s.Flags |= dcbFlagParity
	case ParitySpace:
		s.Parity = dcbSpaceParity
		s.Flags |= dcbFlagParity
	default:
		s.Parity = dcbNoParity
	}
	switch cfg.StopBits {
	case StopBitsOnePointFive:
		s.StopBits = dcbOne5StopBits
	case StopBitsTwo:
		s.StopBits = dcbTwoStopBits
	default:
		s.StopBits = dcbOneStopBit
	}
	if cfg.RTSCTSFlowControl {
		s.Flags |= dcbFlagOutxCtsFlow | dcbFlagRtsControlHandshake
	} else {
		s.Flags |= dcbFlagRtsControlEnable
	}
	if cfg.XONXOFFFlowControl {
		s.Flags |= dcbFlagOutX | dcbFlagInX
	}
	// XonChar/XoffChar/XonLim/XoffLim are set unconditionally, even when
	// XON/XOFF flow control is disabled (fOutX/fInX clear above). Some
	// Windows serial drivers reject SetCommState (error 87, "the parameter
	// is incorrect") when XonChar == XoffChar == 0x00, which is what a
	// zero-value DCB carries — purejavacomm hit this in the field on a
	// plain 8N1 configuration. DC1/DC3 are the conventional XON/XOFF
	// characters and are harmless when the flags are off.
	s.XonChar = 0x11  // DC1
	s.XoffChar = 0x13 // DC3
	// Conventional buffer thresholds for a 4096-byte queue (SetupComm).
	s.XonLim = 2048
	s.XoffLim = 512
	return s
}
