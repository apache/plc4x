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

package slmp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/slmp/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

// newTestConnection builds a connection on a stub codec, with a transaction manager of its own so
// nothing leaks between tests.
func newTestConnection(t *testing.T, configuration Configuration) (*Connection, *stubCodec) {
	t.Helper()
	_options := testutils.EnrichOptionsWithOptionsForTesting(t)
	codec := newStubCodec()
	tm := transactions.NewRequestTransactionManager(maxConcurrentRequests, _options...)
	t.Cleanup(func() {
		assert.NoError(t, tm.Close())
	})
	connection := NewConnection(codec, configuration, NewTagHandler(), tm, map[string][]string{}, _options...)
	return connection, codec
}

// newConnectedTestConnection is newTestConnection plus Connect, which for SLMP is just opening the
// transport - there is no handshake to answer.
func newConnectedTestConnection(t *testing.T, configuration Configuration) (*Connection, *stubCodec) {
	t.Helper()
	connection, codec := newTestConnection(t, configuration)
	require.NoError(t, connection.Connect(testutils.TestContext(t)))
	t.Cleanup(func() {
		assert.NoError(t, connection.Close())
	})
	return connection, codec
}

func TestConnection_ConnectHasNoHandshake(t *testing.T) {
	connection, codec := newTestConnection(t, DefaultConfiguration())
	t.Cleanup(func() {
		assert.NoError(t, connection.Close())
	})

	require.NoError(t, connection.Connect(testutils.TestContext(t)))
	assert.True(t, connection.IsConnected())
	// plc4j's SlmpConnection.onConnect only starts the receive loop: SLMP has no session, so nothing
	// goes on the wire until the first read or write.
	assert.Empty(t, codec.getSent(), "connecting must not put anything on the wire")
}

// TestConnection_CloseWithoutConnect covers the path a failed GetConnection leaves behind: Close has
// to be safe on a connection that was never connected.
func TestConnection_CloseWithoutConnect(t *testing.T) {
	connection, _ := newTestConnection(t, DefaultConfiguration())
	assert.NoError(t, connection.Close())
}

func TestConnection_CloseTwice(t *testing.T) {
	connection, _ := newTestConnection(t, DefaultConfiguration())
	require.NoError(t, connection.Connect(testutils.TestContext(t)))
	assert.NoError(t, connection.Close())
	assert.NoError(t, connection.Close())
}

// readResponseFrame builds a 3E read response carrying the passed end code and payload.
func readResponseFrame(endCode uint16, responseData []byte) readWriteModel.SlmpResponseFrame3E {
	return readWriteModel.NewSlmpResponseFrame3E(endCode, responseData)
}

// readTag runs one read through the connection and answers the request the reader sends with the
// passed end code and payload.
func readTag(t *testing.T, address string, endCode uint16, responseData []byte) (apiModel.PlcReadResponse, readWriteModel.SlmpRequestFrame3E) {
	t.Helper()
	connection, codec := newConnectedTestConnection(t, DefaultConfiguration())

	readRequest, err := connection.ReadRequestBuilder().AddTagAddress("hurz", address).Build()
	require.NoError(t, err)

	resultChan := readRequest.Execute(testutils.TestContext(t))
	request := codec.nextRequest(t)
	frame, ok := request.message.(readWriteModel.SlmpRequestFrame3E)
	require.True(t, ok, "%T is not a 3E request frame", request.message)
	codec.answer(t, request, readResponseFrame(endCode, responseData))

	select {
	case result := <-resultChan:
		require.NoError(t, result.GetErr())
		require.NotNil(t, result.GetResponse())
		return result.GetResponse(), frame
	case <-time.After(5 * time.Second):
		t.Fatal("the read didn't finish")
		return nil, nil
	}
}

// TestConnection_ReadBuildsABatchReadFrame is the request side of the cycle: what the reader puts
// on the wire for a given tag.
func TestConnection_ReadBuildsABatchReadFrame(t *testing.T) {
	// D350 as two words is the SH-080008 Batch Read worked example.
	_, frame := readTag(t, "D350[0..1]:WORD", 0x0000, []byte{0xAB, 0x56, 0x0F, 0x17})

	assert.Equal(t, commandBatchRead, frame.GetCommand())
	assert.Equal(t, subCommandWordUnits, frame.GetSubCommand(), "this version only speaks word units")
	assert.Equal(t, defaultMonitoringTimer, frame.GetMonitoringTimer())
	readRequest, ok := frame.GetRequestData().(readWriteModel.SlmpReadRequest)
	require.True(t, ok, "%T is not a batch read request", frame.GetRequestData())
	assert.Equal(t, uint32(350), readRequest.GetHeadDeviceNumber())
	assert.Equal(t, readWriteModel.SlmpDeviceCode_D, readRequest.GetDeviceCode())
	assert.Equal(t, uint16(2), readRequest.GetNumberOfPoints())
}

// TestConnection_ReadAsksForWordsNotElements is the one arithmetic mistake that would send a frame
// half the size it needs: a REAL is two words, so a REAL[4] tag has to ask for eight points.
func TestConnection_ReadAsksForWordsNotElements(t *testing.T) {
	_, frame := readTag(t, "R200[0..3]:REAL", 0x0000, make([]byte, 16))
	readRequest, ok := frame.GetRequestData().(readWriteModel.SlmpReadRequest)
	require.True(t, ok)
	assert.Equal(t, uint16(8), readRequest.GetNumberOfPoints())
	assert.Equal(t, readWriteModel.SlmpDeviceCode_R, readRequest.GetDeviceCode())
	assert.Equal(t, uint32(200), readRequest.GetHeadDeviceNumber())
}

// TestConnection_ReadUsesTheConfiguredMonitoringTimer pins the one configuration value that ends up
// on the wire.
func TestConnection_ReadUsesTheConfiguredMonitoringTimer(t *testing.T) {
	configuration := DefaultConfiguration()
	configuration.monitoringTimer = 0x00FA
	connection, codec := newConnectedTestConnection(t, configuration)

	readRequest, err := connection.ReadRequestBuilder().AddTagAddress("hurz", "D350").Build()
	require.NoError(t, err)
	resultChan := readRequest.Execute(testutils.TestContext(t))

	request := codec.nextRequest(t)
	frame, ok := request.message.(readWriteModel.SlmpRequestFrame3E)
	require.True(t, ok)
	assert.Equal(t, uint16(0x00FA), frame.GetMonitoringTimer())
	codec.answer(t, request, readResponseFrame(0x0000, []byte{0x01, 0x00}))
	<-resultChan
}

func TestConnection_ReadDecodesTheResponse(t *testing.T) {
	tests := []struct {
		name         string
		address      string
		responseData []byte
		assert       func(t *testing.T, response apiModel.PlcReadResponse)
	}{
		{
			name: "a single WORD", address: "D350", responseData: []byte{0xAB, 0x56},
			assert: func(t *testing.T, response apiModel.PlcReadResponse) {
				assert.Equal(t, uint16(0x56AB), response.GetValue("hurz").GetUint16())
			},
		},
		{
			name: "two WORDs", address: "D350[0..1]:WORD", responseData: []byte{0xAB, 0x56, 0x0F, 0x17},
			assert: func(t *testing.T, response apiModel.PlcReadResponse) {
				value := response.GetValue("hurz")
				require.True(t, value.IsList())
				require.Len(t, value.GetList(), 2)
				assert.Equal(t, uint16(0x56AB), value.GetList()[0].GetUint16())
				assert.Equal(t, uint16(0x170F), value.GetList()[1].GetUint16())
			},
		},
		{
			name: "a negative INT", address: "D350:INT", responseData: []byte{0xFF, 0xFF},
			assert: func(t *testing.T, response apiModel.PlcReadResponse) {
				assert.Equal(t, int16(-1), response.GetValue("hurz").GetInt16())
			},
		},
		{
			name: "a REAL over two words", address: "D350:REAL", responseData: []byte{0x00, 0x00, 0x80, 0x3F},
			assert: func(t *testing.T, response apiModel.PlcReadResponse) {
				assert.InDelta(t, float32(1.0), response.GetValue("hurz").GetFloat32(), 0.0001)
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			response, _ := readTag(t, testCase.address, 0x0000, testCase.responseData)
			require.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("hurz"))
			testCase.assert(t, response)
		})
	}
}

// TestConnection_ReadMapsAnErrorEndCode is the error mapping plc4j's SlmpResponseMapper does: every
// non-zero end code is a REMOTE_ERROR and no value comes back, because on an abnormal completion the
// payload is error information rather than device data.
func TestConnection_ReadMapsAnErrorEndCode(t *testing.T) {
	for _, endCode := range []uint16{0x0055, 0xC059, 0xC050, 0xFFFF} {
		t.Run(formatEndCode(endCode), func(t *testing.T) {
			// The payload is the error-information block, which must not be decoded as a value.
			response, _ := readTag(t, "D350", endCode, []byte{0x00, 0xFF, 0xFF, 0x03, 0x00, 0x5E, 0x01, 0x00, 0xA8, 0x01})
			assert.Equal(t, apiModel.PlcResponseCode_REMOTE_ERROR, response.GetResponseCode("hurz"))
			assert.Nil(t, response.GetValue("hurz"))
		})
	}
}

// TestConnection_ReadMapsAShortPayload is the other half of the mapping: a normal completion whose
// payload can't cover the tag is INVALID_DATA rather than a value made up from what did arrive.
func TestConnection_ReadMapsAShortPayload(t *testing.T) {
	tests := []struct {
		name         string
		address      string
		responseData []byte
	}{
		{name: "nothing at all", address: "D350", responseData: nil},
		{name: "half a word", address: "D350", responseData: []byte{0x01}},
		{name: "one word for a REAL", address: "D350:REAL", responseData: []byte{0x01, 0x02}},
		{name: "one word short of a list", address: "D350[0..1]:WORD", responseData: []byte{0x01, 0x02}},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			response, _ := readTag(t, testCase.address, 0x0000, testCase.responseData)
			assert.Equal(t, apiModel.PlcResponseCode_INVALID_DATA, response.GetResponseCode("hurz"))
			assert.Nil(t, response.GetValue("hurz"))
		})
	}
}

// TestConnection_ReadTimesOut is the "always complete the result channel" case for the read path: a
// device that accepts the connection and then says nothing must not hang the caller.
func TestConnection_ReadTimesOut(t *testing.T) {
	configuration := DefaultConfiguration()
	configuration.requestTimeout = 200 * time.Millisecond
	connection, codec := newConnectedTestConnection(t, configuration)

	readRequest, err := connection.ReadRequestBuilder().AddTagAddress("hurz", "D350").Build()
	require.NoError(t, err)

	start := time.Now()
	select {
	case result := <-readRequest.Execute(testutils.TestContext(t)):
		require.NoError(t, result.GetErr())
		require.NotNil(t, result.GetResponse())
		// plc4j reports this as REMOTE_ERROR; plc4go has a code that says what actually happened.
		assert.Equal(t, apiModel.PlcResponseCode_REQUEST_TIMEOUT, result.GetResponse().GetResponseCode("hurz"))
	case <-time.After(5 * time.Second):
		t.Fatal("the read didn't finish")
	}
	assert.Less(t, time.Since(start), 5*time.Second)
	// The request did go out, it just never got answered.
	assert.Len(t, codec.getSent(), 1)
}

// TestConnection_ReadReportsASendFailure covers the other never-answered path: the frame can't even
// be handed to the transport.
func TestConnection_ReadReportsASendFailure(t *testing.T) {
	connection, codec := newConnectedTestConnection(t, DefaultConfiguration())
	codec.failSends()

	readRequest, err := connection.ReadRequestBuilder().AddTagAddress("hurz", "D350").Build()
	require.NoError(t, err)

	select {
	case result := <-readRequest.Execute(testutils.TestContext(t)):
		require.NoError(t, result.GetErr())
		require.NotNil(t, result.GetResponse())
		assert.Equal(t, apiModel.PlcResponseCode_INTERNAL_ERROR, result.GetResponse().GetResponseCode("hurz"))
	case <-time.After(5 * time.Second):
		t.Fatal("the read didn't finish")
	}
}

// TestConnection_ReadIsolatesPerTagFailures is plc4j's partial-failure isolation: a tag that failed
// carries its own code and the rest of the request still comes back. It also pins that a multi-tag
// request goes out as one frame per tag - there is no request optimizer on either side.
func TestConnection_ReadIsolatesPerTagFailures(t *testing.T) {
	connection, codec := newConnectedTestConnection(t, DefaultConfiguration())

	readRequest, err := connection.ReadRequestBuilder().
		AddTagAddress("good", "D350").
		AddTagAddress("bad", "D351").
		AddTagAddress("alsoGood", "D352").
		Build()
	require.NoError(t, err)
	resultChan := readRequest.Execute(testutils.TestContext(t))

	// One frame per tag, in the order the tags were added.
	first := codec.nextRequest(t)
	codec.answer(t, first, readResponseFrame(0x0000, []byte{0x01, 0x00}))
	second := codec.nextRequest(t)
	codec.answer(t, second, readResponseFrame(0xC059, nil))
	third := codec.nextRequest(t)
	codec.answer(t, third, readResponseFrame(0x0000, []byte{0x03, 0x00}))

	select {
	case result := <-resultChan:
		require.NoError(t, result.GetErr())
		response := result.GetResponse()
		require.NotNil(t, response)
		assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("good"))
		assert.Equal(t, uint16(1), response.GetValue("good").GetUint16())
		assert.Equal(t, apiModel.PlcResponseCode_REMOTE_ERROR, response.GetResponseCode("bad"))
		assert.Nil(t, response.GetValue("bad"))
		assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("alsoGood"))
		assert.Equal(t, uint16(3), response.GetValue("alsoGood").GetUint16())
	case <-time.After(5 * time.Second):
		t.Fatal("the read didn't finish")
	}
	assert.Len(t, codec.getSent(), 3, "one Batch Read frame per tag - there is no optimizer")
}

// writeTag runs one write through the connection and answers the request the writer sends.
func writeTag(t *testing.T, address string, value any, endCode uint16, responseData []byte) (apiModel.PlcWriteResponse, readWriteModel.SlmpRequestFrame3E) {
	t.Helper()
	connection, codec := newConnectedTestConnection(t, DefaultConfiguration())

	writeRequest, err := connection.WriteRequestBuilder().AddTagAddress("hurz", address, value).Build()
	require.NoError(t, err)

	resultChan := writeRequest.Execute(testutils.TestContext(t))
	request := codec.nextRequest(t)
	frame, ok := request.message.(readWriteModel.SlmpRequestFrame3E)
	require.True(t, ok, "%T is not a 3E request frame", request.message)
	codec.answer(t, request, readResponseFrame(endCode, responseData))

	select {
	case result := <-resultChan:
		require.NoError(t, result.GetErr())
		require.NotNil(t, result.GetResponse())
		return result.GetResponse(), frame
	case <-time.After(5 * time.Second):
		t.Fatal("the write didn't finish")
		return nil, nil
	}
}

func TestConnection_WriteBuildsABatchWriteFrame(t *testing.T) {
	response, frame := writeTag(t, "D350[0..1]:WORD", []uint16{0x1234, 0x5678}, 0x0000, nil)
	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("hurz"))

	assert.Equal(t, commandBatchWrite, frame.GetCommand())
	assert.Equal(t, subCommandWordUnits, frame.GetSubCommand())
	writeRequest, ok := frame.GetRequestData().(readWriteModel.SlmpWriteRequest)
	require.True(t, ok, "%T is not a batch write request", frame.GetRequestData())
	assert.Equal(t, uint32(350), writeRequest.GetHeadDeviceNumber())
	assert.Equal(t, readWriteModel.SlmpDeviceCode_D, writeRequest.GetDeviceCode())
	assert.Equal(t, uint16(2), writeRequest.GetNumberOfPoints())
	// Little-endian words, and exactly numberOfPoints * 2 bytes of them: the frame's announced point
	// count and its payload length can't be allowed to disagree.
	assert.Equal(t, []byte{0x34, 0x12, 0x78, 0x56}, writeRequest.GetWriteData())
	assert.Len(t, writeRequest.GetWriteData(), int(writeRequest.GetNumberOfPoints())*2)
}

func TestConnection_WriteEncodesTheValue(t *testing.T) {
	tests := []struct {
		name    string
		address string
		value   any
		want    []byte
	}{
		{name: "a WORD", address: "D350", value: uint16(0x56AB), want: []byte{0xAB, 0x56}},
		{name: "a negative INT", address: "D350:INT", value: int16(-1), want: []byte{0xFF, 0xFF}},
		{name: "a DINT over two words", address: "D350:DINT", value: int32(-2), want: []byte{0xFE, 0xFF, 0xFF, 0xFF}},
		{name: "a REAL over two words", address: "D350:REAL", value: float32(1.0), want: []byte{0x00, 0x00, 0x80, 0x3F}},
		{
			name: "a list of REALs", address: "D350[0..1]:REAL", value: []float32{1.0, 2.0},
			want: []byte{0x00, 0x00, 0x80, 0x3F, 0x00, 0x00, 0x00, 0x40},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			response, frame := writeTag(t, testCase.address, testCase.value, 0x0000, nil)
			assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("hurz"))
			writeRequest, ok := frame.GetRequestData().(readWriteModel.SlmpWriteRequest)
			require.True(t, ok)
			assert.Equal(t, testCase.want, writeRequest.GetWriteData())
		})
	}
}

func TestConnection_WriteMapsAnErrorEndCode(t *testing.T) {
	response, _ := writeTag(t, "D350", uint16(1), 0xC059, []byte{0x00, 0xFF})
	assert.Equal(t, apiModel.PlcResponseCode_REMOTE_ERROR, response.GetResponseCode("hurz"))
}

// TestConnection_WriteRefusesASuccessFrameWithAPayload is plc4j's mis-attribution guard: a Batch
// Write success carries no data (SH-080008), so a payload here is the signature of a response that
// belongs to some other request - a 3E frame has no correlation id to tell them apart. Reporting OK
// would claim a write happened on no evidence.
func TestConnection_WriteRefusesASuccessFrameWithAPayload(t *testing.T) {
	response, _ := writeTag(t, "D350", uint16(1), 0x0000, []byte{0xAB, 0x56})
	assert.Equal(t, apiModel.PlcResponseCode_REMOTE_ERROR, response.GetResponseCode("hurz"))
}

// TestConnection_WriteBuilderRefusesAnUncoercibleValue is the first of the two gates a bad value
// runs into: the value handler refuses to coerce it into the tag's type, so no request is built at
// all and nothing reaches the device.
func TestConnection_WriteBuilderRefusesAnUncoercibleValue(t *testing.T) {
	tests := []struct {
		name    string
		address string
		value   any
	}{
		{name: "a negative value for a WORD", address: "D350:WORD", value: int32(-1)},
		{name: "a string for a REAL", address: "D350:REAL", value: "hurz"},
		{name: "too few values for an array tag", address: "D350[0..2]:WORD", value: []uint16{1, 2}},
		{name: "too many values for an array tag", address: "D350[0..1]:WORD", value: []uint16{1, 2, 3}},
		{name: "a scalar for an array tag", address: "D350[0..1]:WORD", value: uint16(1)},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			connection, codec := newConnectedTestConnection(t, DefaultConfiguration())
			writeRequest, err := connection.WriteRequestBuilder().
				AddTagAddress("hurz", testCase.address, testCase.value).Build()
			assert.Error(t, err)
			assert.Nil(t, writeRequest)
			assert.Empty(t, codec.getSent(), "a value that can't be coerced must not reach the device")
		})
	}
}

// TestConnection_WriteRefusesAnUnencodableValue is the second gate, which the builder can't reach:
// a request assembled by hand (as an interceptor or a test harness may do) whose value the tag
// can't carry. The device has to be left alone and the tag has to report INVALID_DATA rather than
// a frame whose payload doesn't match its announced point count.
func TestConnection_WriteRefusesAnUnencodableValue(t *testing.T) {
	tests := []struct {
		name    string
		address string
		value   apiValues.PlcValue
	}{
		{name: "a negative value for a WORD", address: "D350:WORD", value: spiValues.NewPlcDINT(-1)},
		{
			name: "too few values for an array tag", address: "D350[0..2]:WORD",
			value: spiValues.NewPlcList([]apiValues.PlcValue{spiValues.NewPlcWORD(1), spiValues.NewPlcWORD(2)}),
		},
		{name: "a scalar for an array tag", address: "D350[0..1]:WORD", value: spiValues.NewPlcWORD(1)},
		{name: "no value at all", address: "D350:WORD", value: nil},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			connection, codec := newConnectedTestConnection(t, DefaultConfiguration())
			tag, err := connection.GetPlcTagHandler().ParseTag(testCase.address)
			require.NoError(t, err)
			writer := NewWriter(codec, connection.tm, connection.configuration,
				testutils.EnrichOptionsWithOptionsForTesting(t)...)
			writeRequest := spiModel.NewDefaultPlcWriteRequest(
				map[string]apiModel.PlcTag{"hurz": tag},
				[]string{"hurz"},
				map[string]apiValues.PlcValue{"hurz": testCase.value},
				writer, nil)

			select {
			case result := <-writeRequest.Execute(testutils.TestContext(t)):
				require.NoError(t, result.GetErr())
				require.NotNil(t, result.GetResponse())
				assert.Equal(t, apiModel.PlcResponseCode_INVALID_DATA, result.GetResponse().GetResponseCode("hurz"))
			case <-time.After(5 * time.Second):
				t.Fatal("the write didn't finish")
			}
			assert.Empty(t, codec.getSent(), "a value that can't be encoded must not reach the device")
		})
	}
}

// TestConnection_ReadRejectsAForeignTag covers the same guard on the read path: a tag from another
// driver has to come back as INVALID_ADDRESS rather than panicking somewhere downstream.
func TestConnection_ReadRejectsAForeignTag(t *testing.T) {
	connection, codec := newConnectedTestConnection(t, DefaultConfiguration())
	reader := NewReader(codec, connection.tm, connection.configuration,
		testutils.EnrichOptionsWithOptionsForTesting(t)...)
	readRequest := spiModel.NewDefaultPlcReadRequest(
		map[string]apiModel.PlcTag{"hurz": spiModel.NewDefaultPlcSubscriptionTag(apiModel.SubscriptionCyclic, nil, 0)},
		[]string{"hurz"},
		reader, nil)

	select {
	case result := <-readRequest.Execute(testutils.TestContext(t)):
		require.NoError(t, result.GetErr())
		require.NotNil(t, result.GetResponse())
		assert.Equal(t, apiModel.PlcResponseCode_INVALID_ADDRESS, result.GetResponse().GetResponseCode("hurz"))
	case <-time.After(5 * time.Second):
		t.Fatal("the read didn't finish")
	}
	assert.Empty(t, codec.getSent())
}

// TestConnection_WriteTimesOut is the "always complete the result channel" case for the write path.
func TestConnection_WriteTimesOut(t *testing.T) {
	configuration := DefaultConfiguration()
	configuration.requestTimeout = 200 * time.Millisecond
	connection, codec := newConnectedTestConnection(t, configuration)

	writeRequest, err := connection.WriteRequestBuilder().AddTagAddress("hurz", "D350", uint16(1)).Build()
	require.NoError(t, err)

	start := time.Now()
	select {
	case result := <-writeRequest.Execute(testutils.TestContext(t)):
		require.NoError(t, result.GetErr())
		require.NotNil(t, result.GetResponse())
		// A timed-out write may still have been applied by the device; the driver cannot tell.
		assert.Equal(t, apiModel.PlcResponseCode_REQUEST_TIMEOUT, result.GetResponse().GetResponseCode("hurz"))
	case <-time.After(5 * time.Second):
		t.Fatal("the write didn't finish")
	}
	assert.Less(t, time.Since(start), 5*time.Second)
	assert.Len(t, codec.getSent(), 1)
}

func TestConnection_WriteReportsASendFailure(t *testing.T) {
	connection, codec := newConnectedTestConnection(t, DefaultConfiguration())
	codec.failSends()

	writeRequest, err := connection.WriteRequestBuilder().AddTagAddress("hurz", "D350", uint16(1)).Build()
	require.NoError(t, err)

	select {
	case result := <-writeRequest.Execute(testutils.TestContext(t)):
		require.NoError(t, result.GetErr())
		require.NotNil(t, result.GetResponse())
		assert.Equal(t, apiModel.PlcResponseCode_INTERNAL_ERROR, result.GetResponse().GetResponseCode("hurz"))
	case <-time.After(5 * time.Second):
		t.Fatal("the write didn't finish")
	}
}

func TestConnection_WriteIsolatesPerTagFailures(t *testing.T) {
	connection, codec := newConnectedTestConnection(t, DefaultConfiguration())

	writeRequest, err := connection.WriteRequestBuilder().
		AddTagAddress("good", "D350", uint16(1)).
		AddTagAddress("bad", "D351", uint16(2)).
		Build()
	require.NoError(t, err)
	resultChan := writeRequest.Execute(testutils.TestContext(t))

	codec.answer(t, codec.nextRequest(t), readResponseFrame(0x0000, nil))
	codec.answer(t, codec.nextRequest(t), readResponseFrame(0xC059, nil))

	select {
	case result := <-resultChan:
		require.NoError(t, result.GetErr())
		response := result.GetResponse()
		require.NotNil(t, response)
		assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("good"))
		assert.Equal(t, apiModel.PlcResponseCode_REMOTE_ERROR, response.GetResponseCode("bad"))
	case <-time.After(5 * time.Second):
		t.Fatal("the write didn't finish")
	}
	assert.Len(t, codec.getSent(), 2, "one Batch Write frame per tag")
}

// TestConnection_AcceptsAnyResponseFrame documents the only correlation a 3E frame allows. The
// request transaction manager's concurrency of one is what makes it sound; see the comment on
// acceptsAnyResponseFrame for the caveat it inherits.
func TestConnection_AcceptsAnyResponseFrame(t *testing.T) {
	assert.True(t, acceptsAnyResponseFrame(readResponseFrame(0x0000, nil)))
	// A request frame echoed back is not an answer.
	assert.False(t, acceptsAnyResponseFrame(readWriteModel.NewSlmpRequestFrame3E(0, commandBatchRead, 0,
		readWriteModel.NewSlmpReadRequest(350, readWriteModel.SlmpDeviceCode_D, 1))))
}

func TestConnection_String(t *testing.T) {
	connection, _ := newTestConnection(t, DefaultConfiguration())
	assert.Contains(t, connection.String(), "slmp.Connection{")
}
