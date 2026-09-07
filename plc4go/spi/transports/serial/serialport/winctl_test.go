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

// Values per Microsoft documentation: EscapeCommFunction dwFunc codes
// (SETRTS=3, CLRRTS=4, SETDTR=5, CLRDTR=6) and GetCommModemStatus masks
// (MS_CTS_ON=0x10, MS_DSR_ON=0x20, MS_RING_ON=0x40, MS_RLSD_ON=0x80).
func TestWinEscapeFor(t *testing.T) {
	assert.EqualValues(t, 5, winEscapeFor(true, true), "SETDTR")
	assert.EqualValues(t, 6, winEscapeFor(true, false), "CLRDTR")
	assert.EqualValues(t, 3, winEscapeFor(false, true), "SETRTS")
	assert.EqualValues(t, 4, winEscapeFor(false, false), "CLRRTS")
}

func TestModemStatusFromWinBits(t *testing.T) {
	assert.Equal(t, ModemStatus{}, modemStatusFromWinBits(0))
	assert.Equal(t,
		ModemStatus{CTS: true, DSR: true, RI: true, DCD: true},
		modemStatusFromWinBits(0x10|0x20|0x40|0x80))
	assert.Equal(t, ModemStatus{CTS: true}, modemStatusFromWinBits(0x10))
	assert.Equal(t, ModemStatus{DSR: true}, modemStatusFromWinBits(0x20))
	assert.Equal(t, ModemStatus{RI: true}, modemStatusFromWinBits(0x40))
	assert.Equal(t, ModemStatus{DCD: true}, modemStatusFromWinBits(0x80))
}
