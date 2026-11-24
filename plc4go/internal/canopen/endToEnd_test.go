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

// These tests run the whole driver against a simulated CANopen node on the in-process virtual bus.
// The node is written straight from CiA 301's byte layouts, without the generated model, so a test
// passing means the model's bit fields put the right bytes on the wire, not just that the model
// agrees with itself.

import (
	"context"
	"encoding/binary"
	"fmt"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/can"
	"github.com/apache/plc4x/plc4go/spi/transports/virtualcan"
)

type objectAddress struct {
	index    uint16
	subIndex uint8
}

// simulatedNode is an SDO server with an object dictionary, plus a recorder of everything it sees.
type simulatedNode struct {
	t        *testing.T
	node     uint8
	instance transports.TransportInstance

	mutex      sync.Mutex
	dictionary map[objectAddress][]byte
	seen       []can.Frame

	// download state
	downloadAddress objectAddress
	downloadData    []byte
	// upload state
	uploadData []byte
	toggle     bool
	// flipToggle makes the server answer upload segments with the wrong toggle bit.
	flipToggle bool
	// silent makes the server ignore SDO requests.
	silent bool
}

func startNode(t *testing.T, busName string, node uint8) *simulatedNode {
	t.Helper()
	return startNodeOn(t, virtualcan.NewTransportInstance(busName, testutils.EnrichOptionsWithOptionsForTesting(t)...), node)
}

// startNodeOn runs the simulated node on any CAN transport instance.
func startNodeOn(t *testing.T, instance transports.TransportInstance, node uint8) *simulatedNode {
	t.Helper()
	require.NoError(t, instance.Connect(t.Context()))
	n := &simulatedNode{t: t, node: node, instance: instance, dictionary: map[objectAddress][]byte{}}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		for ctx.Err() == nil {
			readCtx, cancelRead := context.WithTimeout(ctx, 50*time.Millisecond)
			record, err := instance.Read(readCtx, can.FrameSize)
			cancelRead()
			if err != nil {
				continue
			}
			frame, err := can.Decode(record)
			if err != nil {
				t.Errorf("node %d got a malformed record: %v", node, err)
				continue
			}
			n.handle(frame)
		}
	}()
	t.Cleanup(func() {
		cancel()
		<-done
		_ = instance.Close()
	})
	return n
}

func (n *simulatedNode) set(index uint16, subIndex uint8, data []byte) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.dictionary[objectAddress{index, subIndex}] = data
}

func (n *simulatedNode) get(index uint16, subIndex uint8) []byte {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	return n.dictionary[objectAddress{index, subIndex}]
}

func (n *simulatedNode) frames() []can.Frame {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	return append([]can.Frame(nil), n.seen...)
}

func (n *simulatedNode) send(id uint32, data []byte) {
	record, err := can.Frame{ID: id, Data: data}.Encode()
	require.NoError(n.t, err)
	require.NoError(n.t, n.instance.Write(context.Background(), record))
}

func (n *simulatedNode) answer(data []byte) {
	full := make([]byte, 8)
	copy(full, data)
	n.send(0x580+uint32(n.node), full)
}

func (n *simulatedNode) abort(address objectAddress, code uint32) {
	data := make([]byte, 8)
	data[0] = 0x80
	binary.LittleEndian.PutUint16(data[1:], address.index)
	data[3] = address.subIndex
	binary.LittleEndian.PutUint32(data[4:], code)
	n.answer(data)
}

func (n *simulatedNode) handle(frame can.Frame) {
	n.mutex.Lock()
	n.seen = append(n.seen, frame)
	silent := n.silent
	n.mutex.Unlock()
	if frame.ID != 0x600+uint32(n.node) || silent {
		return
	}
	require.Len(n.t, frame.Data, 8, "an SDO request is always 8 bytes")
	data := frame.Data
	address := objectAddress{binary.LittleEndian.Uint16(data[1:]), data[3]}
	n.mutex.Lock()
	defer n.mutex.Unlock()
	switch data[0] >> 5 {
	case 2: // initiate upload
		value, ok := n.dictionary[address]
		if !ok {
			n.mutex.Unlock()
			n.abort(address, 0x06020000)
			n.mutex.Lock()
			return
		}
		response := make([]byte, 8)
		binary.LittleEndian.PutUint16(response[1:], address.index)
		response[3] = address.subIndex
		if len(value) <= 4 {
			response[0] = 0x43 | byte(4-len(value))<<2
			copy(response[4:], value)
		} else {
			response[0] = 0x41
			binary.LittleEndian.PutUint32(response[4:], uint32(len(value)))
			n.uploadData = append([]byte(nil), value...)
			n.toggle = false
		}
		n.mutex.Unlock()
		n.answer(response)
		n.mutex.Lock()
	case 3: // upload segment
		toggle := data[0]&0x10 != 0
		require.Equal(n.t, n.toggle, toggle, "client toggle bit")
		chunk := n.uploadData[:min(7, len(n.uploadData))]
		n.uploadData = n.uploadData[len(chunk):]
		last := len(n.uploadData) == 0
		answerToggle := toggle != n.flipToggle
		header := byte(7-len(chunk)) << 1
		if answerToggle {
			header |= 0x10
		}
		if last {
			header |= 0x01
		}
		n.toggle = !n.toggle
		n.mutex.Unlock()
		n.answer(append([]byte{header}, chunk...))
		n.mutex.Lock()
	case 1: // initiate download
		response := make([]byte, 8)
		response[0] = 0x60
		binary.LittleEndian.PutUint16(response[1:], address.index)
		response[3] = address.subIndex
		if data[0]&0x02 != 0 { // expedited
			size := 4
			if data[0]&0x01 != 0 {
				size = 4 - int(data[0]>>2&0x03)
			}
			n.dictionary[address] = append([]byte(nil), data[4:4+size]...)
		} else {
			n.downloadAddress = address
			n.downloadData = nil
			n.toggle = false
		}
		n.mutex.Unlock()
		n.answer(response)
		n.mutex.Lock()
	case 0: // download segment
		toggle := data[0]&0x10 != 0
		require.Equal(n.t, n.toggle, toggle, "client toggle bit")
		unused := int(data[0] >> 1 & 0x07)
		n.downloadData = append(n.downloadData, data[1:8-unused]...)
		if data[0]&0x01 != 0 {
			n.dictionary[n.downloadAddress] = n.downloadData
		}
		n.toggle = !n.toggle
		header := byte(0x20)
		if toggle {
			header |= 0x10
		}
		n.mutex.Unlock()
		n.answer([]byte{header})
		n.mutex.Lock()
	case 4: // abort from the client
	}
}

func connect(t *testing.T, busName string, query string) plc4go.PlcConnection {
	t.Helper()
	return connectOn(t, virtualcan.NewTransport(), busName, query)
}

// connectOn opens a driver connection over the given CAN transport.
func connectOn(t *testing.T, transport transports.Transport, host string, query string) plc4go.PlcConnection {
	t.Helper()
	driver := NewDriver(testutils.EnrichOptionsWithOptionsForTesting(t)...)
	transportUrl := url.URL{Scheme: transport.GetTransportCode(), Host: host}
	options, err := url.ParseQuery(query)
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	connection, err := driver.GetConnection(ctx, transportUrl,
		map[string]transports.Transport{transport.GetTransportCode(): transport}, options)
	require.NoError(t, err)
	t.Cleanup(func() { _ = connection.Close() })
	return connection
}

func read(t *testing.T, connection plc4go.PlcConnection, addresses map[string]string) apiModel.PlcReadResponse {
	t.Helper()
	builder := connection.ReadRequestBuilder()
	for name, address := range addresses {
		builder.AddTagAddress(name, address)
	}
	request, err := builder.Build()
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	result := <-request.Execute(ctx)
	require.NoError(t, result.GetErr())
	return result.GetResponse()
}

func write(t *testing.T, connection plc4go.PlcConnection, address string, value any) apiModel.PlcResponseCode {
	t.Helper()
	request, err := connection.WriteRequestBuilder().AddTagAddress("tag", address, value).Build()
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	result := <-request.Execute(ctx)
	require.NoError(t, result.GetErr())
	return result.GetResponse().GetResponseCode("tag")
}

func TestReadExpeditedAndSegmented(t *testing.T) {
	bus := t.Name()
	node := startNode(t, bus, 5)
	node.set(0x1018, 1, []byte{0x91, 0x01, 0x00, 0x00})
	node.set(0x1001, 0, []byte{0x04})
	node.set(0x1008, 0, []byte("plc4go test node"))
	node.set(0x2000, 0, []byte{0x00, 0x00, 0x48, 0x41}) // 12.5 as REAL32
	connection := connect(t, bus, "node-id=10")

	response := read(t, connection, map[string]string{
		"vendor":  "SDO:5:0x1018/1:UNSIGNED32",
		"error":   "SDO:5:4097/0:UNSIGNED8",
		"name":    "SDO:5:0x1008/0:VISIBLE_STRING",
		"real":    "SDO:5:0x2000/0:REAL32",
		"missing": "SDO:5:0x6000/1:UNSIGNED16",
		"pdo":     "TRANSMIT_PDO_1:5:UNSIGNED16",
	})
	require.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("vendor"))
	assert.EqualValues(t, 0x191, response.GetValue("vendor").GetUint32())
	require.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("error"))
	assert.EqualValues(t, 4, response.GetValue("error").GetUint8())
	require.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("name"))
	assert.Equal(t, "plc4go test node", response.GetValue("name").GetString())
	require.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("real"))
	assert.InDelta(t, 12.5, response.GetValue("real").GetFloat32(), 0)
	assert.Equal(t, apiModel.PlcResponseCode_NOT_FOUND, response.GetResponseCode("missing"))
	assert.Equal(t, apiModel.PlcResponseCode_UNSUPPORTED, response.GetResponseCode("pdo"))
}

func TestWriteExpeditedAndSegmented(t *testing.T) {
	bus := t.Name()
	node := startNode(t, bus, 7)
	connection := connect(t, bus, "")

	require.Equal(t, apiModel.PlcResponseCode_OK, write(t, connection, "SDO:7:0x1017/0:UNSIGNED16", 1000))
	assert.Equal(t, []byte{0xE8, 0x03}, node.get(0x1017, 0))

	require.Equal(t, apiModel.PlcResponseCode_OK, write(t, connection, "SDO:7:0x2001/2:INTEGER32", -2))
	assert.Equal(t, []byte{0xFE, 0xFF, 0xFF, 0xFF}, node.get(0x2001, 2))

	text := "a value longer than one segment"
	require.Equal(t, apiModel.PlcResponseCode_OK, write(t, connection, "SDO:7:0x2002/0:VISIBLE_STRING", text))
	assert.Equal(t, text, string(node.get(0x2002, 0)))

	response := read(t, connection, map[string]string{"back": "SDO:7:0x2002/0:VISIBLE_STRING"})
	require.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("back"))
	assert.Equal(t, text, response.GetValue("back").GetString())
}

func TestSDOTimeoutAndToggleErrorAbort(t *testing.T) {
	bus := t.Name()
	node := startNode(t, bus, 3)
	node.set(0x1008, 0, []byte("segmented value"))
	connection := connect(t, bus, "request-timeout-ms=100")

	node.mutex.Lock()
	node.flipToggle = true
	node.mutex.Unlock()
	response := read(t, connection, map[string]string{"name": "SDO:3:0x1008/0:VISIBLE_STRING"})
	assert.Equal(t, apiModel.PlcResponseCode_INTERNAL_ERROR, response.GetResponseCode("name"))
	assert.Eventually(t, func() bool { return lastAbortCode(node.frames(), 0x603) == 0x05030000 }, time.Second, 10*time.Millisecond,
		"the driver aborts a transfer whose toggle bit doesn't alternate with 0x05030000")

	node.mutex.Lock()
	node.flipToggle, node.silent = false, true
	node.mutex.Unlock()
	response = read(t, connection, map[string]string{"name": "SDO:3:0x1008/0:VISIBLE_STRING"})
	assert.Equal(t, apiModel.PlcResponseCode_REQUEST_TIMEOUT, response.GetResponseCode("name"))
	assert.Eventually(t, func() bool { return lastAbortCode(node.frames(), 0x603) == 0x05040000 }, time.Second, 10*time.Millisecond,
		"the driver aborts a timed out transfer with 0x05040000")
}

func lastAbortCode(frames []can.Frame, id uint32) uint32 {
	for i := len(frames) - 1; i >= 0; i-- {
		if frames[i].ID == id && len(frames[i].Data) == 8 && frames[i].Data[0] == 0x80 {
			return binary.LittleEndian.Uint32(frames[i].Data[4:])
		}
	}
	return 0
}

func TestNMTAndPDOWrites(t *testing.T) {
	bus := t.Name()
	node := startNode(t, bus, 9)
	connection := connect(t, bus, "")

	require.Equal(t, apiModel.PlcResponseCode_OK, write(t, connection, "NMT:9", 1))
	require.Equal(t, apiModel.PlcResponseCode_OK, write(t, connection, "NMT", 0x80))
	require.Equal(t, apiModel.PlcResponseCode_OK, write(t, connection, "RECEIVE_PDO_1:9:UNSIGNED16", 0x1234))
	assert.Equal(t, apiModel.PlcResponseCode_INVALID_DATA, write(t, connection, "NMT:9", 3))

	require.Eventually(t, func() bool { return len(node.frames()) >= 3 }, time.Second, 10*time.Millisecond)
	frames := node.frames()
	assert.Equal(t, can.Frame{ID: 0x000, Data: []byte{0x01, 0x09}}, frames[0])
	assert.Equal(t, can.Frame{ID: 0x000, Data: []byte{0x80, 0x00}}, frames[1])
	assert.Equal(t, can.Frame{ID: 0x209, Data: []byte{0x34, 0x12}}, frames[2])
}

func TestHeartbeatProducer(t *testing.T) {
	bus := t.Name()
	listener := startNode(t, bus, 1)
	connect(t, bus, "node-id=42&heartbeat=true&heartbeat-interval-ms=30")
	require.Eventually(t, func() bool { return len(listener.frames()) >= 3 }, 2*time.Second, 10*time.Millisecond)
	frames := listener.frames()
	assert.Equal(t, can.Frame{ID: 0x72A, Data: []byte{0x00}}, frames[0], "boot-up first")
	assert.Equal(t, can.Frame{ID: 0x72A, Data: []byte{0x05}}, frames[1], "then operational")
}

func TestSubscriptions(t *testing.T) {
	bus := t.Name()
	node := startNode(t, bus, 4)
	connection := connect(t, bus, "")

	request, err := connection.SubscriptionRequestBuilder().
		AddEventTagAddress("pdo", "TRANSMIT_PDO_1:4:INTEGER16").
		AddChangeOfStateTagAddress("heartbeat", "HEARTBEAT").
		AddEventTagAddress("nmt", "NMT:4").
		AddEventTagAddress("otherPdo", "TRANSMIT_PDO_2:4:UNSIGNED8").
		AddCyclicTagAddress("cyclic", "HEARTBEAT:4", time.Second).
		Build()
	require.NoError(t, err)
	result := <-request.Execute(t.Context())
	require.NoError(t, result.GetErr())
	response := result.GetResponse()
	assert.Equal(t, apiModel.PlcResponseCode_UNSUPPORTED, response.GetResponseCode("cyclic"))

	var mutex sync.Mutex
	events := map[string][]string{}
	for _, name := range []string{"pdo", "heartbeat", "nmt", "otherPdo"} {
		handle, err := response.GetSubscriptionHandle(name)
		require.NoError(t, err)
		handle.Register(func(event apiModel.PlcSubscriptionEvent) {
			mutex.Lock()
			defer mutex.Unlock()
			for _, tagName := range event.GetTagNames() {
				events[tagName] = append(events[tagName], valueText(event.GetValue(tagName)))
			}
		})
	}

	node.send(0x184, []byte{0xFE, 0xFF}) // TPDO1 of node 4: -2
	node.send(0x185, []byte{0x01, 0x00}) // TPDO1 of node 5: not subscribed
	node.send(0x704, []byte{0x05})       // heartbeat operational
	node.send(0x704, []byte{0x05})       // unchanged: not delivered to a change-of-state tag
	node.send(0x705, []byte{0x7F})       // another node, wildcard tag
	node.send(0x000, []byte{0x02, 0x04}) // NMT stop node 4
	node.send(0x000, []byte{0x02, 0x05}) // NMT stop node 5: not subscribed
	node.send(0x080, nil)                // SYNC: ignored
	require.Eventually(t, func() bool {
		mutex.Lock()
		defer mutex.Unlock()
		return len(events["pdo"]) == 1 && len(events["heartbeat"]) == 2 && len(events["nmt"]) == 1
	}, 2*time.Second, 10*time.Millisecond)
	time.Sleep(50 * time.Millisecond)
	mutex.Lock()
	defer mutex.Unlock()
	assert.Equal(t, []string{"-2"}, events["pdo"])
	assert.Len(t, events["heartbeat"], 2)
	assert.Len(t, events["nmt"], 1)
	assert.Empty(t, events["otherPdo"])
}

func valueText(value apiValues.PlcValue) string {
	if value.IsStruct() {
		return fmt.Sprintf("state=%d node=%d", value.GetValue("state").GetUint8(), value.GetValue("node").GetUint8())
	}
	return fmt.Sprintf("%d", value.GetInt64())
}
