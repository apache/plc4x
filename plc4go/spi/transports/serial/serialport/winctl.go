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

// Windows comm-control values kept portable (no x/sys/windows import) so the
// mapping is unit-testable on every development platform. Values per
// Microsoft documentation for EscapeCommFunction (x/sys/windows defines the
// same SETRTS/CLRRTS/SETDTR/CLRDTR values, but only under GOOS=windows) and
// GetCommModemStatus (whose MS_*_ON masks x/sys/windows does not define).
const (
	winEscapeSetRTS = 3 // SETRTS
	winEscapeClrRTS = 4 // CLRRTS
	winEscapeSetDTR = 5 // SETDTR
	winEscapeClrDTR = 6 // CLRDTR

	winModemCTSOn  = 0x0010 // MS_CTS_ON
	winModemDSROn  = 0x0020 // MS_DSR_ON
	winModemRingOn = 0x0040 // MS_RING_ON
	winModemRLSDOn = 0x0080 // MS_RLSD_ON (carrier detect)
)

// winEscapeFor returns the EscapeCommFunction code for asserting/clearing
// the DTR (dtr=true) or RTS (dtr=false) line.
func winEscapeFor(dtr, assert bool) uint32 {
	switch {
	case dtr && assert:
		return winEscapeSetDTR
	case dtr:
		return winEscapeClrDTR
	case assert:
		return winEscapeSetRTS
	default:
		return winEscapeClrRTS
	}
}

// modemStatusFromWinBits translates a GetCommModemStatus bit mask.
func modemStatusFromWinBits(bits uint32) ModemStatus {
	return ModemStatus{
		CTS: bits&winModemCTSOn != 0,
		DSR: bits&winModemDSROn != 0,
		RI:  bits&winModemRingOn != 0,
		DCD: bits&winModemRLSDOn != 0,
	}
}
