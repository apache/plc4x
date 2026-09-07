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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Expected values quoted from Microsoft's winbase.h documentation:
//
//	Parity:   NOPARITY=0, ODDPARITY=1, EVENPARITY=2
//	StopBits: ONESTOPBIT=0, ONE5STOPBITS=1, TWOSTOPBITS=2
//	Flags bits: fBinary=1<<0, fParity=1<<1, fOutxCtsFlow=1<<2,
//	            fDtrControl(2 bits)@4 (DTR_CONTROL_ENABLE=1),
//	            fRtsControl(2 bits)@12 (RTS_CONTROL_ENABLE=1, RTS_CONTROL_HANDSHAKE=2)
func TestMakeDCBSettings(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want dcbSettings
	}{
		{
			name: "9600 8N1 no flow control",
			cfg:  Config{BaudRate: 9600, DataBits: 8, StopBits: StopBitsOne, Parity: ParityNone},
			want: dcbSettings{
				BaudRate: 9600,
				ByteSize: 8,
				Parity:   0,                   // NOPARITY
				StopBits: 0,                   // ONESTOPBIT
				Flags:    0x1 | 0x10 | 0x1000, // fBinary | DTR_CONTROL_ENABLE | RTS_CONTROL_ENABLE
				// XonChar/XoffChar/XonLim/XoffLim are always populated, even
				// with XON/XOFF flow control disabled — see makeDCBSettings.
				XonChar: 0x11, XoffChar: 0x13,
				XonLim: 2048, XoffLim: 512,
			},
		},
		{
			name: "19200 7E2",
			cfg:  Config{BaudRate: 19200, DataBits: 7, StopBits: StopBitsTwo, Parity: ParityEven},
			want: dcbSettings{
				BaudRate: 19200,
				ByteSize: 7,
				Parity:   2,                         // EVENPARITY
				StopBits: 2,                         // TWOSTOPBITS
				Flags:    0x1 | 0x2 | 0x10 | 0x1000, // ... | fParity | ...
				XonChar:  0x11, XoffChar: 0x13,
				XonLim: 2048, XoffLim: 512,
			},
		},
		{
			name: "odd parity",
			cfg:  Config{BaudRate: 9600, DataBits: 8, StopBits: StopBitsOne, Parity: ParityOdd},
			want: dcbSettings{
				BaudRate: 9600,
				ByteSize: 8,
				Parity:   1, // ODDPARITY
				StopBits: 0,
				Flags:    0x1 | 0x2 | 0x10 | 0x1000,
				XonChar:  0x11, XoffChar: 0x13,
				XonLim: 2048, XoffLim: 512,
			},
		},
		{
			name: "1.5 stop bits with 5 data bits",
			cfg:  Config{BaudRate: 9600, DataBits: 5, StopBits: StopBitsOnePointFive, Parity: ParityNone},
			want: dcbSettings{
				BaudRate: 9600,
				ByteSize: 5,
				Parity:   0,
				StopBits: 1, // ONE5STOPBITS
				Flags:    0x1 | 0x10 | 0x1000,
				XonChar:  0x11, XoffChar: 0x13,
				XonLim: 2048, XoffLim: 512,
			},
		},
		{
			name: "RTS/CTS flow control switches RTS to handshake and enables CTS sensitivity",
			cfg:  Config{BaudRate: 115200, DataBits: 8, StopBits: StopBitsOne, Parity: ParityNone, RTSCTSFlowControl: true},
			want: dcbSettings{
				BaudRate: 115200,
				ByteSize: 8,
				Parity:   0,
				StopBits: 0,
				Flags:    0x1 | 0x4 | 0x10 | 0x2000, // fBinary | fOutxCtsFlow | DTR_CONTROL_ENABLE | RTS_CONTROL_HANDSHAKE
				XonChar:  0x11, XoffChar: 0x13,
				XonLim: 2048, XoffLim: 512,
			},
		},
		{
			name: "mark parity",
			cfg:  Config{BaudRate: 9600, DataBits: 8, StopBits: StopBitsOne, Parity: ParityMark},
			want: dcbSettings{
				BaudRate: 9600, ByteSize: 8,
				Parity:   3, // MARKPARITY
				StopBits: 0,
				Flags:    0x1 | 0x2 | 0x10 | 0x1000,
				XonChar:  0x11, XoffChar: 0x13,
				XonLim: 2048, XoffLim: 512,
			},
		},
		{
			name: "space parity",
			cfg:  Config{BaudRate: 9600, DataBits: 8, StopBits: StopBitsOne, Parity: ParitySpace},
			want: dcbSettings{
				BaudRate: 9600, ByteSize: 8,
				Parity:   4, // SPACEPARITY
				StopBits: 0,
				Flags:    0x1 | 0x2 | 0x10 | 0x1000,
				XonChar:  0x11, XoffChar: 0x13,
				XonLim: 2048, XoffLim: 512,
			},
		},
		{
			name: "xon/xoff flow control sets fOutX and fInX with DC1/DC3 chars",
			cfg:  Config{BaudRate: 9600, DataBits: 8, StopBits: StopBitsOne, Parity: ParityNone, XONXOFFFlowControl: true},
			want: dcbSettings{
				BaudRate: 9600, ByteSize: 8, Parity: 0, StopBits: 0,
				Flags:   0x1 | 0x10 | 0x100 | 0x200 | 0x1000, // fBinary | DTR_CONTROL_ENABLE | fOutX | fInX | RTS_CONTROL_ENABLE
				XonChar: 0x11, XoffChar: 0x13,
				XonLim: 2048, XoffLim: 512,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, makeDCBSettings(tt.cfg))
		})
	}
}
