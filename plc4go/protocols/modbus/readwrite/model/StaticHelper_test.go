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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The expected values are the CRCs of the frames in the Modbus RTU parser-serializer testsuite,
// which is the same set plc4j checks its rtuCrcCheck against.
func TestRtuCrcCheck(t *testing.T) {
	tests := []struct {
		name    string
		address uint8
		pdu     ModbusPDU
		want    uint16
	}{
		{
			name:    "read holding registers request",
			address: 1,
			pdu:     NewModbusPDUReadHoldingRegistersRequest(0, 10),
			want:    0xC5CD,
		},
		{
			name:    "read holding registers response",
			address: 1,
			pdu:     NewModbusPDUReadHoldingRegistersResponse(make([]byte, 20)),
			want:    0xA367,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			crc, err := RtuCrcCheck(context.Background(), testCase.address, testCase.pdu)()
			require.NoError(t, err)
			assert.Equal(t, testCase.want, crc)
		})
	}
}

// The CRC covers the address as well as the PDU, so a frame that is only addressed differently has
// a different checksum - that is what keeps a resynchronizing reader from accepting a neighbour's
// frame with the leading byte shifted off.
func TestRtuCrcCheckCoversTheAddress(t *testing.T) {
	pdu := NewModbusPDUReadHoldingRegistersRequest(0, 10)

	first, err := RtuCrcCheck(context.Background(), 1, pdu)()
	require.NoError(t, err)
	second, err := RtuCrcCheck(context.Background(), 2, pdu)()
	require.NoError(t, err)

	assert.NotEqual(t, first, second)
}

func TestRtuCrcCheckRejectsAMissingPdu(t *testing.T) {
	_, err := RtuCrcCheck(context.Background(), 1, nil)()
	assert.Error(t, err)
}
