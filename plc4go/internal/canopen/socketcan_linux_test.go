//go:build linux

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

package canopen

// The same scenarios as endToEnd_test.go, over a real SocketCAN interface. They run only when
// CANOPEN_TEST_INTERFACE names one, typically a virtual CAN interface:
//
//	ip link add dev vcan0 type vcan && ip link set up vcan0
//	CANOPEN_TEST_INTERFACE=vcan0 go test ./internal/canopen/

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports/socketcan"
)

func socketCANInterface(t *testing.T) string {
	name := os.Getenv("CANOPEN_TEST_INTERFACE")
	if name == "" {
		t.Skip("set CANOPEN_TEST_INTERFACE to a (virtual) CAN interface to run the SocketCAN tests")
	}
	return name
}

func TestSocketCANReadWrite(t *testing.T) {
	name := socketCANInterface(t)
	node := startNodeOn(t, socketcan.NewTransportInstance(name, socketcan.NewTransport(), testutils.EnrichOptionsWithOptionsForTesting(t)...), 21)
	node.set(0x1018, 1, []byte{0x91, 0x01, 0x00, 0x00})
	node.set(0x1008, 0, []byte("socketcan segmented"))
	connection := connectOn(t, socketcan.NewTransport(), name, "node-id=20")

	response := read(t, connection, map[string]string{
		"vendor": "SDO:21:0x1018/1:UNSIGNED32",
		"name":   "SDO:21:0x1008/0:VISIBLE_STRING",
	})
	require.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("vendor"))
	assert.EqualValues(t, 0x191, response.GetValue("vendor").GetUint32())
	require.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("name"))
	assert.Equal(t, "socketcan segmented", response.GetValue("name").GetString())

	require.Equal(t, apiModel.PlcResponseCode_OK, write(t, connection, "SDO:21:0x2003/0:VISIBLE_STRING", "written over SocketCAN"))
	assert.Equal(t, "written over SocketCAN", string(node.get(0x2003, 0)))
}
