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

// Runs the driver against CANopen nodes written independently of plc4x, on a SocketCAN
// interface: CANopenNode (CANOPEN_TEST_CANOPENNODE, its node id) and the python canopen library
// (CANOPEN_TEST_PYTHON). tools/canopen-peer runs both on a virtual CAN bus in Docker.

import (
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/transports/socketcan"
)

// peer is the node id in the named environment variable; the test is skipped without it.
func peer(t *testing.T, variable string) (iface string, node uint64) {
	iface = socketCANInterface(t)
	text := os.Getenv(variable)
	if text == "" {
		t.Skipf("set %s to the node id of the peer", variable)
	}
	node, err := strconv.ParseUint(text, 0, 7)
	require.NoError(t, err)
	return iface, node
}

func connectToPeer(t *testing.T, iface string) plc4go.PlcConnection {
	return connectOn(t, socketcan.NewTransport(), iface, "node-id=20")
}

func sdo(node uint64, object string) string {
	return "SDO:" + strconv.FormatUint(node, 10) + ":" + object
}

// heartbeats subscribes to a node's heartbeat and returns the states it reports.
func heartbeats(t *testing.T, connection plc4go.PlcConnection, node uint64) func() []uint8 {
	var mutex sync.Mutex
	var states []uint8
	request, err := connection.SubscriptionRequestBuilder().
		AddEventTagAddress("heartbeat", "HEARTBEAT:"+strconv.FormatUint(node, 10)).
		AddPreRegisteredConsumer("heartbeat", func(event apiModel.PlcSubscriptionEvent) {
			mutex.Lock()
			defer mutex.Unlock()
			states = append(states, event.GetValue("heartbeat").GetValue("state").GetUint8())
		}).Build()
	require.NoError(t, err)
	result := <-request.Execute(t.Context())
	require.NoError(t, result.GetErr())
	require.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("heartbeat"))
	return func() []uint8 {
		mutex.Lock()
		defer mutex.Unlock()
		return append([]uint8(nil), states...)
	}
}

// lastStateIs waits for the node's heartbeats to report state.
func lastStateIs(t *testing.T, states func() []uint8, state uint8, message string) {
	require.Eventually(t, func() bool {
		s := states()
		return len(s) > 0 && s[len(s)-1] == state
	}, 3*time.Second, 20*time.Millisecond, message)
}

// nmtAndHeartbeat sets the node's heartbeat by SDO, then moves it through the NMT states and
// checks its heartbeats follow: stopped (4), pre-operational (0x7F), operational (5).
func nmtAndHeartbeat(t *testing.T, connection plc4go.PlcConnection, node uint64) {
	states := heartbeats(t, connection, node)
	require.Equal(t, apiModel.PlcResponseCode_OK, write(t, connection, sdo(node, "0x1017/0:UNSIGNED16"), 50))
	response := read(t, connection, map[string]string{"heartbeat": sdo(node, "0x1017/0:UNSIGNED16")})
	require.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("heartbeat"))
	assert.EqualValues(t, 50, response.GetValue("heartbeat").GetUint16())
	nmt := "NMT:" + strconv.FormatUint(node, 10)
	require.Equal(t, apiModel.PlcResponseCode_OK, write(t, connection, nmt, 0x02))
	lastStateIs(t, states, 0x04, "stopped after NMT stop")
	require.Equal(t, apiModel.PlcResponseCode_OK, write(t, connection, nmt, 0x80))
	lastStateIs(t, states, 0x7F, "pre-operational after NMT enter pre-operational")
	require.Equal(t, apiModel.PlcResponseCode_OK, write(t, connection, nmt, 0x01))
	lastStateIs(t, states, 0x05, "operational after NMT start")
	// Heartbeats keep coming at about 50 ms.
	before := len(states())
	time.Sleep(500 * time.Millisecond)
	assert.GreaterOrEqual(t, len(states())-before, 5, "heartbeats every 50 ms")
}

// tpdo subscribes to the node's TPDO 1 and returns the values it carries.
func tpdo(t *testing.T, connection plc4go.PlcConnection, node uint64, dataType string) func() []int64 {
	var mutex sync.Mutex
	var values []int64
	request, err := connection.SubscriptionRequestBuilder().
		AddEventTagAddress("pdo", "TRANSMIT_PDO_1:"+strconv.FormatUint(node, 10)+":"+dataType).
		AddPreRegisteredConsumer("pdo", func(event apiModel.PlcSubscriptionEvent) {
			if event.GetResponseCode("pdo") != apiModel.PlcResponseCode_OK {
				return
			}
			mutex.Lock()
			defer mutex.Unlock()
			values = append(values, event.GetValue("pdo").GetInt64())
		}).Build()
	require.NoError(t, err)
	result := <-request.Execute(t.Context())
	require.NoError(t, result.GetErr())
	return func() []int64 {
		mutex.Lock()
		defer mutex.Unlock()
		return append([]int64(nil), values...)
	}
}

func TestPeerCANopenNodeSDO(t *testing.T) {
	iface, node := peer(t, "CANOPEN_TEST_CANOPENNODE")
	connection := connectToPeer(t, iface)

	// Defaults of the DS301 example profile.
	response := read(t, connection, map[string]string{
		"deviceType": sdo(node, "0x1000/0:UNSIGNED32"),
		"sync":       sdo(node, "0x1005/0:UNSIGNED32"),
		"emcy":       sdo(node, "0x1014/0:UNSIGNED32"),
		"sdoServer":  sdo(node, "0x1200/1:UNSIGNED32"),
		"identity":   sdo(node, "0x1018/0:UNSIGNED8"),
		"missing":    sdo(node, "0x6000/1:UNSIGNED16"),
		"noSubIndex": sdo(node, "0x1018/9:UNSIGNED32"),
	})
	for _, name := range []string{"deviceType", "sync", "emcy", "sdoServer", "identity"} {
		require.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode(name), name)
	}
	assert.EqualValues(t, 0, response.GetValue("deviceType").GetUint32())
	assert.EqualValues(t, 0x80, response.GetValue("sync").GetUint32())
	assert.EqualValues(t, 0x80+node, response.GetValue("emcy").GetUint32())
	assert.EqualValues(t, 0x600+node, response.GetValue("sdoServer").GetUint32())
	assert.EqualValues(t, 4, response.GetValue("identity").GetUint8())
	assert.Equal(t, apiModel.PlcResponseCode_NOT_FOUND, response.GetResponseCode("missing"))
	assert.Equal(t, apiModel.PlcResponseCode_NOT_FOUND, response.GetResponseCode("noSubIndex"))

	assert.Equal(t, apiModel.PlcResponseCode_ACCESS_DENIED, write(t, connection, sdo(node, "0x1000/0:UNSIGNED32"), 1), "read only")
	require.Equal(t, apiModel.PlcResponseCode_OK, write(t, connection, sdo(node, "0x1006/0:UNSIGNED32"), 123456))
	response = read(t, connection, map[string]string{"cycle": sdo(node, "0x1006/0:UNSIGNED32")})
	assert.EqualValues(t, 123456, response.GetValue("cycle").GetUint32())
	require.Equal(t, apiModel.PlcResponseCode_OK, write(t, connection, sdo(node, "0x1006/0:UNSIGNED32"), 0))
}

func TestPeerCANopenNodeNMTAndHeartbeat(t *testing.T) {
	iface, node := peer(t, "CANOPEN_TEST_CANOPENNODE")
	nmtAndHeartbeat(t, connectToPeer(t, iface), node)
}

// The driver maps the error register into TPDO 1 by SDO, as a master configures a node, and
// receives the PDOs the node then sends.
func TestPeerCANopenNodeTPDO(t *testing.T) {
	iface, node := peer(t, "CANOPEN_TEST_CANOPENNODE")
	connection := connectToPeer(t, iface)
	nmt := "NMT:" + strconv.FormatUint(node, 10)
	require.Equal(t, apiModel.PlcResponseCode_OK, write(t, connection, nmt, 0x80), "PDOs are configured in pre-operational")
	cobID := uint32(0x180 + node)
	for _, step := range []struct {
		object string
		value  uint32
	}{
		{"0x1800/1:UNSIGNED32", 0x80000000 | cobID}, // disable
		{"0x1A00/0:UNSIGNED8", 0},
		{"0x1A00/1:UNSIGNED32", 0x10010008}, // 0x1001/0, 8 bits
		{"0x1A00/0:UNSIGNED8", 1},
		{"0x1800/2:UNSIGNED8", 254},
		{"0x1800/5:UNSIGNED16", 50},    // event timer, ms
		{"0x1800/1:UNSIGNED32", cobID}, // enable
	} {
		require.Equal(t, apiModel.PlcResponseCode_OK, write(t, connection, sdo(node, step.object), step.value), step.object)
	}
	values := tpdo(t, connection, node, "UNSIGNED8")
	require.Equal(t, apiModel.PlcResponseCode_OK, write(t, connection, nmt, 0x01))
	require.Eventually(t, func() bool { return len(values()) >= 5 }, 3*time.Second, 20*time.Millisecond, "TPDO 1 every 50 ms")
	for _, value := range values() {
		assert.EqualValues(t, 0, value, "the error register is clear")
	}
	require.Equal(t, apiModel.PlcResponseCode_OK, write(t, connection, nmt, 0x02))
	time.Sleep(100 * time.Millisecond)
	count := len(values())
	time.Sleep(300 * time.Millisecond)
	assert.Equal(t, count, len(values()), "no PDOs when stopped")
	// Stopped nodes do not answer SDO; back in pre-operational the PDO is switched off again.
	assert.Equal(t, apiModel.PlcResponseCode_REQUEST_TIMEOUT, read(t, connection, map[string]string{"v": sdo(node, "0x1000/0:UNSIGNED32")}).GetResponseCode("v"))
	require.Equal(t, apiModel.PlcResponseCode_OK, write(t, connection, nmt, 0x80))
	require.Equal(t, apiModel.PlcResponseCode_OK, write(t, connection, sdo(node, "0x1800/1:UNSIGNED32"), 0x80000000|cobID))
}

func TestPeerPythonSDO(t *testing.T) {
	iface, node := peer(t, "CANOPEN_TEST_PYTHON")
	connection := connectToPeer(t, iface)
	response := read(t, connection, map[string]string{
		"deviceType": sdo(node, "0x1000/0:UNSIGNED32"),
		"name":       sdo(node, "0x1008/0:VISIBLE_STRING"),
		"vendor":     sdo(node, "0x1018/1:UNSIGNED32"),
		"number":     sdo(node, "0x2001/0:INTEGER32"),
		"real":       sdo(node, "0x2002/0:REAL32"),
		"missing":    sdo(node, "0x6000/1:UNSIGNED16"),
	})
	for _, name := range []string{"deviceType", "name", "vendor", "number", "real"} {
		require.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode(name), name)
	}
	assert.EqualValues(t, 0x00020192, response.GetValue("deviceType").GetUint32())
	assert.Equal(t, "plc4x canopen python peer", response.GetValue("name").GetString(), "segmented upload")
	assert.EqualValues(t, 0xAB42, response.GetValue("vendor").GetUint32())
	assert.EqualValues(t, -7, response.GetValue("number").GetInt32())
	assert.EqualValues(t, 1.5, response.GetValue("real").GetFloat32())
	assert.Equal(t, apiModel.PlcResponseCode_NOT_FOUND, response.GetResponseCode("missing"))

	assert.Equal(t, apiModel.PlcResponseCode_ACCESS_DENIED, write(t, connection, sdo(node, "0x1008/0:VISIBLE_STRING"), "x"), "read only")
	text := "written by plc4go in several segments"
	require.Equal(t, apiModel.PlcResponseCode_OK, write(t, connection, sdo(node, "0x2000/0:VISIBLE_STRING"), text), "segmented download")
	require.Equal(t, apiModel.PlcResponseCode_OK, write(t, connection, sdo(node, "0x2002/0:REAL32"), float32(-0.25)))
	response = read(t, connection, map[string]string{
		"text": sdo(node, "0x2000/0:VISIBLE_STRING"),
		"real": sdo(node, "0x2002/0:REAL32"),
	})
	assert.Equal(t, text, response.GetValue("text").GetString())
	assert.EqualValues(t, -0.25, response.GetValue("real").GetFloat32())
}

func TestPeerPythonNMTAndHeartbeat(t *testing.T) {
	iface, node := peer(t, "CANOPEN_TEST_PYTHON")
	nmtAndHeartbeat(t, connectToPeer(t, iface), node)
}

// The python node sends 0x2001 as TPDO 1 while operational.
func TestPeerPythonTPDO(t *testing.T) {
	iface, node := peer(t, "CANOPEN_TEST_PYTHON")
	connection := connectToPeer(t, iface)
	nmt := "NMT:" + strconv.FormatUint(node, 10)
	values := tpdo(t, connection, node, "INTEGER32")
	require.Equal(t, apiModel.PlcResponseCode_OK, write(t, connection, sdo(node, "0x2001/0:INTEGER32"), -123456))
	require.Equal(t, apiModel.PlcResponseCode_OK, write(t, connection, nmt, 0x01))
	require.Eventually(t, func() bool {
		v := values()
		return len(v) >= 3 && v[len(v)-1] == -123456
	}, 3*time.Second, 20*time.Millisecond, "TPDO 1 carries 0x2001")
	require.Equal(t, apiModel.PlcResponseCode_OK, write(t, connection, nmt, 0x02))
}
