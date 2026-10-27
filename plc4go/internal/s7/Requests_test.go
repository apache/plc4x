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

package s7

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

const testPduSize = 240

// simulatedPlc is a message codec that plays a PLC holding one data block: it answers each read
// with the memory it asks for and applies each write, item by item - and refuses any request or
// response bigger than the PDU, as a CPU does.
type simulatedPlc struct {
	spi.MessageCodec // only SendRequest is used

	t  *testing.T
	db []byte
	// failRequest is the index of a request whose items are all answered with ACCESS_DENIED, or -1.
	failRequest int

	lock     sync.Mutex
	requests []readWriteModel.S7Message
}

func newSimulatedPlc(t *testing.T, size int) *simulatedPlc {
	return &simulatedPlc{t: t, db: make([]byte, size), failRequest: -1}
}

func (p *simulatedPlc) SendRequest(_ context.Context, _ string, message spi.Message, acceptsMessage spi.AcceptsMessage,
	handleMessage spi.HandleMessage, _ spi.HandleError) error {
	request := message.(readWriteModel.TPKTPacket).GetPayload().(readWriteModel.COTPPacketData).GetPayload()
	p.lock.Lock()
	fail := len(p.requests) == p.failRequest
	p.requests = append(p.requests, request)
	p.lock.Unlock()
	assert.LessOrEqual(p.t, request.GetLengthInBytes(context.Background()), uint64(testPduSize), "request bigger than the PDU")

	var response readWriteModel.S7Message
	switch parameter := request.GetParameter().(type) {
	case readWriteModel.S7ParameterReadVarRequest:
		response = p.read(request.GetTpduReference(), parameter, fail)
	case readWriteModel.S7ParameterWriteVarRequest:
		response = p.write(request.GetTpduReference(), parameter, request.GetPayload().(readWriteModel.S7PayloadWriteVarRequest), fail)
	default:
		return errors.Errorf("unexpected request %T", parameter)
	}
	assert.LessOrEqual(p.t, response.GetLengthInBytes(context.Background()), uint64(testPduSize), "response bigger than the PDU")
	packet := readWriteModel.NewTPKTPacket(readWriteModel.NewCOTPPacketData(nil, response, true, uint8(request.GetTpduReference())))
	if !acceptsMessage(packet) {
		return errors.New("the response was not accepted")
	}
	return handleMessage(packet)
}

func (p *simulatedPlc) read(tpdu uint16, parameter readWriteModel.S7ParameterReadVarRequest, fail bool) readWriteModel.S7Message {
	var items []readWriteModel.S7VarPayloadDataItem
	for _, item := range parameter.GetItems() {
		address := item.(readWriteModel.S7VarRequestParameterItemAddress).GetAddress().(readWriteModel.S7AddressAny)
		start := int(address.GetByteAddress())
		switch {
		case fail:
			items = append(items, readWriteModel.NewS7VarPayloadDataItem(readWriteModel.DataTransportErrorCode_ACCESS_DENIED, readWriteModel.DataTransportSize_NULL, nil))
		case address.GetTransportSize() == readWriteModel.TransportSize_BOOL:
			items = append(items, readWriteModel.NewS7VarPayloadDataItem(readWriteModel.DataTransportErrorCode_OK, readWriteModel.DataTransportSize_BIT,
				[]byte{(p.db[start] >> address.GetBitAddress()) & 0x01}))
		default:
			assert.Zero(p.t, address.GetBitAddress(), "only a BOOL read has a bit address")
			end := start + int(payloadBytes(address))
			items = append(items, readWriteModel.NewS7VarPayloadDataItem(readWriteModel.DataTransportErrorCode_OK, readWriteModel.DataTransportSize_BYTE_WORD_DWORD,
				append([]byte(nil), p.db[start:end]...)))
		}
	}
	return readWriteModel.NewS7MessageResponseData(tpdu, readWriteModel.NewS7ParameterReadVarResponse(uint8(len(items))),
		readWriteModel.NewS7PayloadReadVarResponse(items), 0, 0)
}

func (p *simulatedPlc) write(tpdu uint16, parameter readWriteModel.S7ParameterWriteVarRequest, payload readWriteModel.S7PayloadWriteVarRequest, fail bool) readWriteModel.S7Message {
	var statuses []readWriteModel.S7VarPayloadStatusItem
	for i, item := range parameter.GetItems() {
		if fail {
			statuses = append(statuses, readWriteModel.NewS7VarPayloadStatusItem(readWriteModel.DataTransportErrorCode_ACCESS_DENIED))
			continue
		}
		address := item.(readWriteModel.S7VarRequestParameterItemAddress).GetAddress().(readWriteModel.S7AddressAny)
		data := payload.GetItems()[i].GetData()
		start := int(address.GetByteAddress())
		if address.GetTransportSize() == readWriteModel.TransportSize_BOOL {
			mask := byte(1) << address.GetBitAddress()
			if data[0]&0x01 != 0 {
				p.db[start] |= mask
			} else {
				p.db[start] &^= mask
			}
		} else {
			assert.Zero(p.t, address.GetBitAddress(), "only a BOOL write has a bit address")
			copy(p.db[start:], data)
		}
		statuses = append(statuses, readWriteModel.NewS7VarPayloadStatusItem(readWriteModel.DataTransportErrorCode_OK))
	}
	return readWriteModel.NewS7MessageResponseData(tpdu, readWriteModel.NewS7ParameterWriteVarResponse(uint8(len(statuses))),
		readWriteModel.NewS7PayloadWriteVarResponse(statuses), 0, 0)
}

func newTestDriverContext() *DriverContext {
	return &DriverContext{PduSize: testPduSize, ControllerType: readWriteModel.ControllerType_S7_1500}
}

// The tests here stay off spi/testutils: its testing logger turns on colored output process-wide,
// which changes the byte-exact rendering TestS7MessageBytes compares against.
func newTestTransactionManager(t *testing.T) transactions.RequestTransactionManager {
	tm := transactions.NewRequestTransactionManager(1)
	t.Cleanup(func() { assert.NoError(t, tm.Close()) })
	return tm
}

func parseTags(t *testing.T, addresses ...string) (map[string]apiModel.PlcTag, []string) {
	tags := map[string]apiModel.PlcTag{}
	var names []string
	for i, address := range addresses {
		tag, err := NewTagHandler().ParseTag(address)
		require.NoError(t, err, address)
		name := fmt.Sprintf("t%d", i)
		tags[name] = tag
		names = append(names, name)
	}
	return tags, names
}

func readFrom(t *testing.T, plc *simulatedPlc, addresses ...string) apiModel.PlcReadResponse {
	tags, names := parseTags(t, addresses...)
	reader := NewReader(&TpduGenerator{}, plc, newTestTransactionManager(t), newTestDriverContext())
	select {
	case result := <-reader.Read(t.Context(), spiModel.NewDefaultPlcReadRequest(tags, names, reader, nil)):
		require.NoError(t, result.GetErr())
		return result.GetResponse()
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the read never returned")
		return nil
	}
}

func writeTo(t *testing.T, plc *simulatedPlc, values []apiValues.PlcValue, addresses ...string) apiModel.PlcWriteResponse {
	tags, names := parseTags(t, addresses...)
	valueMap := map[string]apiValues.PlcValue{}
	for i, name := range names {
		valueMap[name] = values[i]
	}
	writer := NewWriter(&TpduGenerator{}, plc, newTestTransactionManager(t), newTestDriverContext())
	select {
	case result := <-writer.Write(t.Context(), spiModel.NewDefaultPlcWriteRequest(tags, names, valueMap, writer, nil)):
		require.NoError(t, result.GetErr())
		return result.GetResponse()
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the write never returned")
		return nil
	}
}

// A request that fits into one PDU goes out as one request, as it always did.
func TestReadThatFitsIsOneRequest(t *testing.T) {
	plc := newSimulatedPlc(t, 1024)
	binary.BigEndian.PutUint16(plc.db[20:], 42)

	response := readFrom(t, plc, "%DB1.DBW20:INT", "%DB1:16.3:BOOL")

	assert.Equal(t, int16(42), response.GetValue("t0").GetInt16())
	assert.Len(t, plc.requests, 1)
}

// Every item of a read costs request and response space, so a read of many tags is spread over
// several requests, each the PDU can hold.
func TestReadOfManyTagsIsSplitIntoRequests(t *testing.T) {
	plc := newSimulatedPlc(t, 1024)
	var addresses []string
	for i := range 40 {
		binary.BigEndian.PutUint16(plc.db[i*2:], uint16(i*7))
		addresses = append(addresses, fmt.Sprintf("%%DB1.DBW%d:INT", i*2))
	}

	response := readFrom(t, plc, addresses...)

	for i := range 40 {
		name := fmt.Sprintf("t%d", i)
		assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode(name))
		assert.Equal(t, int16(i*7), response.GetValue(name).GetInt16(), name)
	}
	assert.Greater(t, len(plc.requests), 1)
}

// A tag too big for one PDU is fetched in byte fragments and decoded once they are all in.
func TestReadOfATagBiggerThanOnePdu(t *testing.T) {
	plc := newSimulatedPlc(t, 1024)
	var expected []float32
	for i := range 100 {
		value := float32(i) * 1.5
		expected = append(expected, value)
		binary.BigEndian.PutUint32(plc.db[400+i*4:], math.Float32bits(value))
	}

	response := readFrom(t, plc, "%DB1:400[0..99]:REAL")

	require.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("t0"))
	var values []float32
	for _, item := range response.GetValue("t0").GetList() {
		values = append(values, item.GetFloat32())
	}
	assert.Equal(t, expected, values)
	assert.Greater(t, len(plc.requests), 1)
}

func TestReadOfABoolArrayBiggerThanOnePdu(t *testing.T) {
	plc := newSimulatedPlc(t, 1024)
	var expected []bool
	for i := range 2000 {
		value := i%3 == 0 || i%7 == 0
		expected = append(expected, value)
		if bit := 16*8 + 3 + i; value {
			plc.db[bit/8] |= 1 << (bit % 8)
		}
	}

	response := readFrom(t, plc, "%DB1:16.3[0..1999]:BOOL")

	require.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("t0"))
	var values []bool
	for _, item := range response.GetValue("t0").GetList() {
		values = append(values, item.GetBool())
	}
	assert.Equal(t, expected, values)
}

// A later fragment that succeeds must not cover up an earlier one that failed.
func TestReadOfAFragmentedTagFailsWithAFragment(t *testing.T) {
	plc := newSimulatedPlc(t, 1024)
	plc.failRequest = 0

	response := readFrom(t, plc, "%DB1:0[0..199]:INT")

	assert.Equal(t, apiModel.PlcResponseCode_ACCESS_DENIED, response.GetResponseCode("t0"))
}

// A BOOL array sharing bytes at both ends is written as up to 15 items, more than one PDU holds.
func TestWriteIsSplitIntoRequests(t *testing.T) {
	plc := newSimulatedPlc(t, 64)
	var bools []apiValues.PlcValue
	for range 22 {
		bools = append(bools, spiValues.NewPlcBOOL(true))
	}

	response := writeTo(t, plc, []apiValues.PlcValue{spiValues.NewPlcList(bools)}, "%DB1:16.1[0..21]:BOOL")

	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("t0"))
	assert.Equal(t, []byte{0b11111110, 0xFF, 0b01111111}, plc.db[16:19])
	assert.Greater(t, len(plc.requests), 1)
}

// A write cannot be split without changing what it writes, so an item too big for any PDU fails
// its tag without stopping the others.
func TestWriteOfAnItemBiggerThanOnePdu(t *testing.T) {
	plc := newSimulatedPlc(t, 1024)
	big := make([]byte, 300)

	response := writeTo(t, plc, []apiValues.PlcValue{spiValues.NewPlcRawByteArray(big), spiValues.NewPlcINT(42)},
		"%DB1:0[0..299]:BYTE", "%DB1.DBW400:INT")

	assert.Equal(t, apiModel.PlcResponseCode_INVALID_DATA, response.GetResponseCode("t0"))
	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("t1"))
	assert.Equal(t, uint16(42), binary.BigEndian.Uint16(plc.db[400:]))
}
