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

package umas

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/umas/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

// readTags starts a read of the given addresses and hands back the channel the result arrives on. It
// deliberately does not wait: the test has to play the PLC first.
func readTags(t *testing.T, connection *Connection, addresses map[string]string) <-chan apiModel.PlcReadRequestResult {
	t.Helper()
	builder := connection.ReadRequestBuilder()
	for tagName, address := range addresses {
		builder = builder.AddTagAddress(tagName, address)
	}
	request, err := builder.Build()
	require.NoError(t, err)
	return request.Execute(testutils.TestContext(t))
}

// awaitReadResponse insists that the read completes. A read which never delivers a result would hang
// its caller forever, which is the single worst failure mode a driver can have.
func awaitReadResponse(t *testing.T, results <-chan apiModel.PlcReadRequestResult) apiModel.PlcReadResponse {
	t.Helper()
	select {
	case result := <-results:
		require.NotNil(t, result)
		require.NoError(t, result.GetErr())
		require.NotNil(t, result.GetResponse())
		return result.GetResponse()
	case <-time.After(5 * time.Second):
		t.Fatal("the read never delivered a result")
		return nil
	}
}

// A scalar read: the reference says which block and offset, and the size index says how wide the
// value is. The 32 bit symbol offset is split with the low 8 bits in the offset field and everything
// above them in baseOffset - plc4j's split, reproduced.
func TestReader_ReadsAScalar(t *testing.T) {
	connection, codec := newConnectedConnection(t)

	results := readTags(t, connection, map[string]string{"real": "g_r32"})

	request := codec.nextRequest(t)
	readRequest, ok := request.item(t).(readWriteModel.UmasPDUReadVariableRequest)
	require.True(t, ok, "%T is not a read-variable request", request.item(t))
	assert.Equal(t, uint8(0x22), readRequest.GetUmasFunctionKey())
	assert.Equal(t, uint32(testProjectCrc), readRequest.GetCrc(), "a read has to carry the project CRC")
	assert.Equal(t, uint8(1), readRequest.GetVariableCount())
	require.Len(t, readRequest.GetVariables(), 1)

	reference := readRequest.GetVariables()[0]
	assert.Equal(t, scalarReadFlag, reference.GetIsArray())
	assert.Equal(t, uint8(3), reference.GetDataSizeIndex(), "a REAL is four bytes, which is size index 3")
	assert.Equal(t, uint16(2), reference.GetBlock())
	// The symbol sits at 0x00001234.
	assert.Equal(t, uint16(0x12), reference.GetBaseOffset())
	assert.Equal(t, uint8(0x34), reference.GetOffset())
	assert.Nil(t, reference.GetArrayLength(), "a scalar read has no array length")

	// 1.0f, little endian.
	codec.answerWith(t, request, readWriteModel.NewUmasPDUReadVariableResponse(0,
		[]byte{0x00, 0x00, 0x80, 0x3F}))

	response := awaitReadResponse(t, results)
	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("real"))
	assert.Equal(t, spiValues.NewPlcREAL(1.0), response.GetValue("real"))
}

// A read of a symbol whose offset needs more than 24 bits can't be expressed: the baseOffset field is
// 16 bits and the offset field 8. Truncating it would read some other symbol's memory.
func TestReader_RefusesAnOffsetTheReferenceCannotCarry(t *testing.T) {
	connection, _ := newConnectedConnection(t)

	response := awaitReadResponse(t, readTags(t, connection, map[string]string{"far": "g_toofar"}))
	assert.Equal(t, apiModel.PlcResponseCode_INVALID_ADDRESS, response.GetResponseCode("far"))
	assert.Nil(t, response.GetValue("far"))
}

// A STRING is read as a byte array, because STRING's request size of 17 doesn't fit the 4 bit size
// index field. The length comes from the symbol layout, or from the default when the symbol is the
// last of its memory block and has no successor to measure against.
func TestReader_ReadsAStringAsAByteArray(t *testing.T) {
	connection, codec := newConnectedConnection(t)

	results := readTags(t, connection, map[string]string{"text": "g_string"})

	request := codec.nextRequest(t)
	readRequest := request.item(t).(readWriteModel.UmasPDUReadVariableRequest)
	reference := readRequest.GetVariables()[0]
	assert.Equal(t, arrayReadFlag, reference.GetIsArray())
	assert.Equal(t, stringReadSizeIndex, reference.GetDataSizeIndex(), "single bytes")
	require.NotNil(t, reference.GetArrayLength())
	// g_string is the last symbol of block 2, so its size can't be derived and the default applies.
	assert.Equal(t, defaultStringBufferSize, *reference.GetArrayLength())

	codec.answerWith(t, request, readWriteModel.NewUmasPDUReadVariableResponse(0,
		append([]byte("hello"), 0x00, 0x00, 0x00)))

	response := awaitReadResponse(t, results)
	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("text"))
	assert.Equal(t, spiValues.NewPlcSTRING("hello"), response.GetValue("text"))
}

// A symbol whose size the layout does give away asks for exactly that many bytes.
func TestReader_ReadsAStringSizedFromTheSymbolLayout(t *testing.T) {
	connection, codec := newTestConnection(t, DefaultConfiguration())
	connectResult := connect(t, connection)
	fixture := defaultFixture()
	fixture.symbols = []symbolRecord{
		{name: "g_string", dataType: typeIdString, block: 2, offset: 0x00000000},
		// The next symbol of the same block is 0x14 bytes further on, so that is g_string's size.
		{name: "g_after", dataType: typeIdDint, block: 2, offset: 0x00000014},
	}
	runHandshake(t, codec, fixture)
	requireConnected(t, connectResult)

	results := readTags(t, connection, map[string]string{"text": "g_string"})
	request := codec.nextRequest(t)
	reference := request.item(t).(readWriteModel.UmasPDUReadVariableRequest).GetVariables()[0]
	require.NotNil(t, reference.GetArrayLength())
	assert.Equal(t, uint16(0x14), *reference.GetArrayLength())

	codec.answerWith(t, request, readWriteModel.NewUmasPDUReadVariableResponse(0, []byte{'h', 'i', 0x00}))
	response := awaitReadResponse(t, results)
	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("text"))
}

// A symbol of one of the project's own types is opaque to the driver, so the caller gets the bytes.
// The four bytes come from plc4j's hard-coded fallback size index.
func TestReader_ReadsACustomTypeAsRawBytes(t *testing.T) {
	connection, codec := newConnectedConnection(t)

	results := readTags(t, connection, map[string]string{"plant": "g_plant"})

	request := codec.nextRequest(t)
	reference := request.item(t).(readWriteModel.UmasPDUReadVariableRequest).GetVariables()[0]
	assert.Equal(t, scalarReadFlag, reference.GetIsArray())
	assert.Equal(t, unknownDataTypeSizeIndex, reference.GetDataSizeIndex())

	codec.answerWith(t, request, readWriteModel.NewUmasPDUReadVariableResponse(0, []byte{0x01, 0x02, 0x03, 0x04}))

	response := awaitReadResponse(t, results)
	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("plant"))
	assert.Equal(t, spiValues.NewPlcRawByteArray([]byte{0x01, 0x02, 0x03, 0x04}), response.GetValue("plant"))
}

func TestReader_UnknownSymbolIsNotFound(t *testing.T) {
	connection, _ := newConnectedConnection(t)

	response := awaitReadResponse(t, readTags(t, connection, map[string]string{"nope": "g_nothing"}))
	assert.Equal(t, apiModel.PlcResponseCode_NOT_FOUND, response.GetResponseCode("nope"))
	assert.Nil(t, response.GetValue("nope"))
}

func TestReader_RefusedReadIsARemoteError(t *testing.T) {
	connection, codec := newConnectedConnection(t)

	results := readTags(t, connection, map[string]string{"real": "g_r32"})
	codec.answerWith(t, codec.nextRequest(t), readWriteModel.NewUmasPDUErrorResponse(0, []byte{0x01}))

	response := awaitReadResponse(t, results)
	assert.Equal(t, apiModel.PlcResponseCode_REMOTE_ERROR, response.GetResponseCode("real"))
}

func TestReader_TransportErrorIsARemoteError(t *testing.T) {
	connection, codec := newConnectedConnection(t)

	results := readTags(t, connection, map[string]string{"real": "g_r32"})
	codec.failRequest(t, codec.nextRequest(t), errors.New("the PLC hung up"))

	response := awaitReadResponse(t, results)
	assert.Equal(t, apiModel.PlcResponseCode_REMOTE_ERROR, response.GetResponseCode("real"))
}

// A payload the type can't be decoded from is reported as bad data rather than passed off as a value.
func TestReader_UndecodableResponseIsInvalidData(t *testing.T) {
	connection, codec := newConnectedConnection(t)

	results := readTags(t, connection, map[string]string{"real": "g_r32"})
	// Two bytes where a REAL needs four.
	codec.answerWith(t, codec.nextRequest(t), readWriteModel.NewUmasPDUReadVariableResponse(0, []byte{0x01, 0x02}))

	response := awaitReadResponse(t, results)
	assert.Equal(t, apiModel.PlcResponseCode_INVALID_DATA, response.GetResponseCode("real"))
	assert.Nil(t, response.GetValue("real"))
}

// A read whose request the PLC never answers has to time out with a result, not hang.
func TestReader_TimesOutWithAResult(t *testing.T) {
	configuration := DefaultConfiguration()
	configuration.requestTimeout = 200 * time.Millisecond
	connection, codec := newTestConnection(t, configuration)
	connectResult := connect(t, connection)
	runHandshake(t, codec, defaultFixture())
	requireConnected(t, connectResult)

	results := readTags(t, connection, map[string]string{"real": "g_r32"})
	// The request goes out and is never answered.
	codec.nextRequest(t)

	response := awaitReadResponse(t, results)
	assert.Equal(t, apiModel.PlcResponseCode_REQUEST_TIMEOUT, response.GetResponseCode("real"))
}

// Several tags in one request are read one after the other, because the PLC answers one request at a
// time and a variable read response is one undelimited block of bytes.
func TestReader_ReadsSeveralTagsOneAfterTheOther(t *testing.T) {
	connection, codec := newConnectedConnection(t)

	results := readTags(t, connection, map[string]string{"real": "g_r32", "int": "g_b16"})

	// Two requests, one per tag, each answered with a payload of the right width for its type.
	for i := 0; i < 2; i++ {
		request := codec.nextRequest(t)
		readRequest := request.item(t).(readWriteModel.UmasPDUReadVariableRequest)
		require.Len(t, readRequest.GetVariables(), 1, "one variable per request")
		switch readRequest.GetVariables()[0].GetDataSizeIndex() {
		case 3: // the REAL
			codec.answerWith(t, request, readWriteModel.NewUmasPDUReadVariableResponse(0,
				[]byte{0x00, 0x00, 0x80, 0x3F}))
		case 2: // the INT
			codec.answerWith(t, request, readWriteModel.NewUmasPDUReadVariableResponse(0,
				[]byte{0x2A, 0x00}))
		default:
			t.Fatalf("unexpected size index %d", readRequest.GetVariables()[0].GetDataSizeIndex())
		}
	}

	response := awaitReadResponse(t, results)
	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("real"))
	assert.Equal(t, spiValues.NewPlcREAL(1.0), response.GetValue("real"))
	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("int"))
	assert.Equal(t, spiValues.NewPlcINT(42), response.GetValue("int"))
}

// A read on a connection whose dictionary download failed answers NOT_FOUND for every tag rather
// than hanging or guessing.
func TestReader_WithoutADictionaryNothingResolves(t *testing.T) {
	connection, codec := newTestConnection(t, DefaultConfiguration())
	connectResult := connect(t, connection)
	runHandshakeWithoutDictionary(t, codec)
	codec.answerWith(t, codec.nextRequest(t), readWriteModel.NewUmasPDUErrorResponse(0, []byte{0x01}))
	requireConnected(t, connectResult)

	response := awaitReadResponse(t, readTags(t, connection, map[string]string{"real": "g_r32"}))
	assert.Equal(t, apiModel.PlcResponseCode_NOT_FOUND, response.GetResponseCode("real"))
}
