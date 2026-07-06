//go:build linux || darwin || freebsd

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
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

// Expected flag values per POSIX (IEEE Std 1003.1) termios: raw,
// non-canonical mode with receiver enabled and modem status ignored.
func TestMakeTermios(t *testing.T) {
	base := uint64(unix.CLOCAL | unix.CREAD)
	tests := []struct {
		name      string
		cfg       Config
		wantCflag uint64
	}{
		{
			name:      "8N1",
			cfg:       Config{BaudRate: 9600, DataBits: 8, StopBits: StopBitsOne, Parity: ParityNone},
			wantCflag: base | unix.CS8,
		},
		{
			name:      "7E1",
			cfg:       Config{BaudRate: 9600, DataBits: 7, StopBits: StopBitsOne, Parity: ParityEven},
			wantCflag: base | unix.CS7 | unix.PARENB,
		},
		{
			name:      "odd parity sets PARENB and PARODD",
			cfg:       Config{BaudRate: 9600, DataBits: 8, StopBits: StopBitsOne, Parity: ParityOdd},
			wantCflag: base | unix.CS8 | unix.PARENB | unix.PARODD,
		},
		{
			name:      "two stop bits set CSTOPB",
			cfg:       Config{BaudRate: 9600, DataBits: 8, StopBits: StopBitsTwo, Parity: ParityNone},
			wantCflag: base | unix.CS8 | unix.CSTOPB,
		},
		{
			name:      "5 data bits",
			cfg:       Config{BaudRate: 9600, DataBits: 5, StopBits: StopBitsOne, Parity: ParityNone},
			wantCflag: base | unix.CS5,
		},
		{
			name:      "6 data bits",
			cfg:       Config{BaudRate: 9600, DataBits: 6, StopBits: StopBitsOne, Parity: ParityNone},
			wantCflag: base | unix.CS6,
		},
		{
			name:      "hardware flow control sets CRTSCTS",
			cfg:       Config{BaudRate: 9600, DataBits: 8, StopBits: StopBitsOne, Parity: ParityNone, RTSCTSFlowControl: true},
			wantCflag: base | unix.CS8 | unix.CRTSCTS,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := makeTermios(tt.cfg)
			require.NoError(t, err)
			assert.Equal(t, tt.wantCflag, uint64(got.Cflag), "Cflag")
			// Raw mode: input/output/local processing fully disabled per POSIX
			// (no ICANON/ECHO/ISIG, no IXON/ICRNL/INPCK, no OPOST).
			assert.Zero(t, got.Iflag, "Iflag must be 0 (raw)")
			assert.Zero(t, got.Oflag, "Oflag must be 0 (raw)")
			assert.Zero(t, got.Lflag, "Lflag must be 0 (raw)")
			// VMIN/VTIME are documented intent only: the fd is non-blocking, so
			// the runtime poller (deadlines) governs read blocking behavior.
			assert.EqualValues(t, 1, got.Cc[unix.VMIN], "VMIN")
			assert.EqualValues(t, 0, got.Cc[unix.VTIME], "VTIME")
		})
	}
}

func TestMakeTermiosRejectsOnePointFiveStopBits(t *testing.T) {
	_, err := makeTermios(Config{BaudRate: 9600, DataBits: 5, StopBits: StopBitsOnePointFive})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "1.5 stop bits")
}

func TestMakeTermiosXONXOFF(t *testing.T) {
	got, err := makeTermios(Config{BaudRate: 9600, DataBits: 8, StopBits: StopBitsOne, Parity: ParityNone, XONXOFFFlowControl: true})
	require.NoError(t, err)
	// POSIX software flow control: IXON (honor received XOFF/XON) and
	// IXOFF (emit XOFF/XON) with the conventional DC1/DC3 characters.
	assert.EqualValues(t, unix.IXON|unix.IXOFF, got.Iflag, "Iflag")
	assert.EqualValues(t, 0x11, got.Cc[unix.VSTART], "VSTART must be DC1")
	assert.EqualValues(t, 0x13, got.Cc[unix.VSTOP], "VSTOP must be DC3")
}
