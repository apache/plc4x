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

package bacnetip

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// TestNativeBacnetRead_SurvivesMalformedDatagrams pins, over a real UDP socket, that one malformed
// datagram neither stalls the reads after it nor ends the connection. The udp transport only reads
// the next datagram once the buffered bytes are consumed, so a codec that waits for the rest of a
// short frame never sees another reply.
func TestNativeBacnetRead_SurvivesMalformedDatagrams(t *testing.T) {
	valid, err := makeWhoIsBVLC(t).Serialize()
	require.NoError(t, err)
	for _, tc := range []struct {
		name     string
		datagram []byte
	}{
		{name: "shorter than the header", datagram: []byte{0x81, 0x0a}},
		{name: "shorter than its BVLC length", datagram: []byte{0x81, 0x0a, 0x00, 0x64, 0x01}},
		{name: "BVLC length below the header", datagram: []byte{0x81, 0x0a, 0x00, 0x02, 0x01}},
		{name: "trailing bytes", datagram: append(append([]byte(nil), valid...), 0xde, 0xad)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			log := testutils.ProduceTestingLogger(t)
			device := startFakeBacnetDevice(t, log)
			t.Cleanup(device.stop)
			connection := connectForeignDeviceTest(t, log, device, map[string][]string{
				"ApduTimeoutMs": {"1000"},
			})
			t.Cleanup(func() { assert.NoError(t, connection.Close()) })

			code, _, ok := readPresentValue(t, connection)
			require.True(t, ok)
			require.Equal(t, apiModel.PlcResponseCode_OK, code)

			// The device's own socket sends it, so even the socket connected to the device takes it.
			device.mu.Lock()
			connectionAddress := device.lastSrc
			device.mu.Unlock()
			_, err := device.conn.WriteToUDP(tc.datagram, connectionAddress)
			require.NoError(t, err)

			code, value, ok := readPresentValue(t, connection)
			require.True(t, ok, "the read after a malformed datagram must complete")
			assert.Equal(t, apiModel.PlcResponseCode_OK, code)
			assert.Equal(t, float32(23.5), value)
			assert.True(t, connection.IsConnected(), "a malformed datagram must not end the connection")
		})
	}
}
