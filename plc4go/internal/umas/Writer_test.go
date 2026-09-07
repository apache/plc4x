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
)

// writeTag starts a write and hands back the channel the result arrives on, without waiting: the test
// has to play the PLC first.
func writeTag(t *testing.T, connection *Connection, tagName string, address string, value any) <-chan apiModel.PlcWriteRequestResult {
	t.Helper()
	request, err := connection.WriteRequestBuilder().AddTagAddress(tagName, address, value).Build()
	require.NoError(t, err)
	return request.Execute(testutils.TestContext(t))
}

// awaitWriteResponse insists that the write completes. A write which never delivers a result would
// hang its caller forever.
func awaitWriteResponse(t *testing.T, results <-chan apiModel.PlcWriteRequestResult) apiModel.PlcWriteResponse {
	t.Helper()
	select {
	case result := <-results:
		require.NotNil(t, result)
		require.NoError(t, result.GetErr())
		require.NotNil(t, result.GetResponse())
		return result.GetResponse()
	case <-time.After(5 * time.Second):
		t.Fatal("the write never delivered a result")
		return nil
	}
}

// A scalar write. Note the offset split, which is not the one a read uses: the low 16 bits go into
// baseOffset and the high 16 into offset, where a read puts the low 8 into offset. plc4j documents the
// asymmetry as observed, so it is reproduced.
func TestWriter_WritesAScalar(t *testing.T) {
	connection, codec := newConnectedConnection(t)

	results := writeTag(t, connection, "real", "g_r32", float32(1.0))

	request := codec.nextRequest(t)
	writeRequest, ok := request.item(t).(readWriteModel.UmasPDUWriteVariableRequest)
	require.True(t, ok, "%T is not a write-variable request", request.item(t))
	assert.Equal(t, uint8(0x23), writeRequest.GetUmasFunctionKey())
	assert.Equal(t, uint32(testProjectCrc), writeRequest.GetCrc(), "a write has to carry the project CRC")
	assert.Equal(t, uint8(1), writeRequest.GetVariableCount())
	require.Len(t, writeRequest.GetVariables(), 1)

	reference := writeRequest.GetVariables()[0]
	assert.Equal(t, scalarReadFlag, reference.GetIsArray())
	assert.Equal(t, uint8(3), reference.GetDataSizeIndex())
	assert.Equal(t, uint16(2), reference.GetBlock())
	// The symbol sits at 0x00001234: the low 16 bits into baseOffset, the high 16 into offset.
	assert.Equal(t, uint16(0x1234), reference.GetBaseOffset())
	assert.Equal(t, uint16(0x0000), reference.GetOffset())
	assert.Nil(t, reference.GetArrayLength())
	assert.Equal(t, []byte{0x00, 0x00, 0x80, 0x3F}, reference.GetRecordData())

	codec.answerWith(t, request, readWriteModel.NewUmasPDUWriteVariableResponse(0, nil))

	response := awaitWriteResponse(t, results)
	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("real"))
}

// The same symbol offset lands in different fields on a read than on a write, which is the one place
// the asymmetry is visible side by side.
func TestWriter_SplitsTheOffsetDifferentlyThanAReadDoes(t *testing.T) {
	connection, codec := newConnectedConnection(t)

	// g_far sits at 0x00012345.
	readResults := readTags(t, connection, map[string]string{"far": "g_far"})
	readRequestMessage := codec.nextRequest(t)
	readReference := readRequestMessage.item(t).(readWriteModel.UmasPDUReadVariableRequest).GetVariables()[0]
	assert.Equal(t, uint16(0x0123), readReference.GetBaseOffset(), "a read shifts the offset right by eight")
	assert.Equal(t, uint8(0x45), readReference.GetOffset())
	codec.answerWith(t, readRequestMessage, readWriteModel.NewUmasPDUReadVariableResponse(0,
		[]byte{0x01, 0x00, 0x00, 0x00}))
	awaitReadResponse(t, readResults)

	writeResults := writeTag(t, connection, "far", "g_far", int32(1))
	writeRequestMessage := codec.nextRequest(t)
	writeReference := writeRequestMessage.item(t).(readWriteModel.UmasPDUWriteVariableRequest).GetVariables()[0]
	assert.Equal(t, uint16(0x2345), writeReference.GetBaseOffset(), "a write takes the low sixteen bits")
	assert.Equal(t, uint16(0x0001), writeReference.GetOffset(), "and the high sixteen")
	codec.answerWith(t, writeRequestMessage, readWriteModel.NewUmasPDUWriteVariableResponse(0, nil))
	awaitWriteResponse(t, writeResults)
}

// A STRING goes over as a byte array of its own length, with the array length field saying how many.
func TestWriter_WritesAStringAsAByteArray(t *testing.T) {
	connection, codec := newConnectedConnection(t)

	results := writeTag(t, connection, "text", "g_string", "hi")

	request := codec.nextRequest(t)
	reference := request.item(t).(readWriteModel.UmasPDUWriteVariableRequest).GetVariables()[0]
	assert.Equal(t, arrayReadFlag, reference.GetIsArray())
	assert.Equal(t, stringReadSizeIndex, reference.GetDataSizeIndex())
	require.NotNil(t, reference.GetArrayLength())
	// The string plus its NUL terminator.
	assert.Equal(t, uint16(3), *reference.GetArrayLength())
	assert.Equal(t, []byte{'h', 'i', 0x00}, reference.GetRecordData())

	codec.answerWith(t, request, readWriteModel.NewUmasPDUWriteVariableResponse(0, nil))
	assert.Equal(t, apiModel.PlcResponseCode_OK, awaitWriteResponse(t, results).GetResponseCode("text"))
}

// The payload length and the size index have to agree, because the PLC sizes the payload from the
// index rather than from how many bytes arrived.
func TestWriter_PayloadLengthMatchesTheSizeIndex(t *testing.T) {
	connection, codec := newConnectedConnection(t)

	results := writeTag(t, connection, "int", "g_b16", int16(42))
	request := codec.nextRequest(t)
	reference := request.item(t).(readWriteModel.UmasPDUWriteVariableRequest).GetVariables()[0]
	assert.Equal(t, uint8(2), reference.GetDataSizeIndex())
	assert.Len(t, reference.GetRecordData(), 2)
	assert.Equal(t, reference.GetDataSize(), uint16(len(reference.GetRecordData())))

	codec.answerWith(t, request, readWriteModel.NewUmasPDUWriteVariableResponse(0, nil))
	assert.Equal(t, apiModel.PlcResponseCode_OK, awaitWriteResponse(t, results).GetResponseCode("int"))
}

func TestWriter_UnknownSymbolIsNotFound(t *testing.T) {
	connection, _ := newConnectedConnection(t)

	response := awaitWriteResponse(t, writeTag(t, connection, "nope", "g_nothing", int32(1)))
	assert.Equal(t, apiModel.PlcResponseCode_NOT_FOUND, response.GetResponseCode("nope"))
}

// A value the symbol's type can't hold is refused before anything goes on the wire, rather than
// truncated the way plc4j's getInteger() would.
func TestWriter_RefusesAValueTheTypeCannotHold(t *testing.T) {
	connection, _ := newConnectedConnection(t)

	response := awaitWriteResponse(t, writeTag(t, connection, "int", "g_b16", int32(40000)))
	assert.Equal(t, apiModel.PlcResponseCode_INVALID_DATATYPE, response.GetResponseCode("int"))
}

// A value of the wrong shape for the symbol is a response code for that one tag, not a panic taking
// the whole request down. A []byte is the reachable case: a UMAS tag has no value type, so the value
// handler wraps it as a PlcRawByteArray, whose GetString panics in the plain value adapter.
func TestWriter_RefusesAValueOfTheWrongShape(t *testing.T) {
	connection, _ := newConnectedConnection(t)

	// awaitWriteResponse insists on a response rather than an error, which is what a recovered panic
	// in Writer.Write would have delivered instead.
	response := awaitWriteResponse(t, writeTag(t, connection, "text", "g_string", []byte("hi")))
	assert.Equal(t, apiModel.PlcResponseCode_INVALID_DATATYPE, response.GetResponseCode("text"))
}

func TestWriter_RefusedWriteIsARemoteError(t *testing.T) {
	connection, codec := newConnectedConnection(t)

	results := writeTag(t, connection, "real", "g_r32", float32(1.0))
	codec.answerWith(t, codec.nextRequest(t), readWriteModel.NewUmasPDUErrorResponse(0, []byte{0x01}))

	assert.Equal(t, apiModel.PlcResponseCode_REMOTE_ERROR, awaitWriteResponse(t, results).GetResponseCode("real"))
}

func TestWriter_TransportErrorIsARemoteError(t *testing.T) {
	connection, codec := newConnectedConnection(t)

	results := writeTag(t, connection, "real", "g_r32", float32(1.0))
	codec.failRequest(t, codec.nextRequest(t), errors.New("the PLC hung up"))

	assert.Equal(t, apiModel.PlcResponseCode_REMOTE_ERROR, awaitWriteResponse(t, results).GetResponseCode("real"))
}

// A write the PLC never answers has to time out with a result, not hang.
func TestWriter_TimesOutWithAResult(t *testing.T) {
	configuration := DefaultConfiguration()
	configuration.requestTimeout = 200 * time.Millisecond
	connection, codec := newTestConnection(t, configuration)
	connectResult := connect(t, connection)
	runHandshake(t, codec, defaultFixture())
	requireConnected(t, connectResult)

	results := writeTag(t, connection, "real", "g_r32", float32(1.0))
	codec.nextRequest(t)

	assert.Equal(t, apiModel.PlcResponseCode_REQUEST_TIMEOUT, awaitWriteResponse(t, results).GetResponseCode("real"))
}

// An answer of the wrong kind is an internal error, not a silent success.
func TestWriter_UnexpectedResponseIsAnInternalError(t *testing.T) {
	connection, codec := newConnectedConnection(t)

	results := writeTag(t, connection, "real", "g_r32", float32(1.0))
	codec.answerWith(t, codec.nextRequest(t), readWriteModel.NewUmasPDUSuccessResponse(0, nil))

	assert.Equal(t, apiModel.PlcResponseCode_INTERNAL_ERROR, awaitWriteResponse(t, results).GetResponseCode("real"))
}

// Writing a symbol of one of the project's own types only works from raw bytes of exactly the width
// the fallback size index declares - which is as far as plc4j gets with such a symbol too.
func TestWriter_CustomTypeNeedsRawBytesOfTheRightWidth(t *testing.T) {
	connection, codec := newConnectedConnection(t)

	tooShort := awaitWriteResponse(t, writeTag(t, connection, "plant", "g_plant", []byte{0x01, 0x02}))
	assert.Equal(t, apiModel.PlcResponseCode_INVALID_ADDRESS, tooShort.GetResponseCode("plant"),
		"two bytes don't match the four the size index declares")

	results := writeTag(t, connection, "plant", "g_plant", []byte{0x01, 0x02, 0x03, 0x04})
	request := codec.nextRequest(t)
	reference := request.item(t).(readWriteModel.UmasPDUWriteVariableRequest).GetVariables()[0]
	assert.Equal(t, []byte{0x01, 0x02, 0x03, 0x04}, reference.GetRecordData())
	codec.answerWith(t, request, readWriteModel.NewUmasPDUWriteVariableResponse(0, nil))
	assert.Equal(t, apiModel.PlcResponseCode_OK, awaitWriteResponse(t, results).GetResponseCode("plant"))
}

// The value handler must not coerce a value to the tag's type: a tag parsed from an address has no
// type, and coercing to that would turn every value into a PlcNull before the writer sees it.
func TestValueHandler_KeepsTheValueWhenTheTagHasNoType(t *testing.T) {
	handler := NewValueHandler()
	tag, err := NewTagHandler().ParseTag("g_r32")
	require.NoError(t, err)

	value, err := handler.NewPlcValue(tag, float32(1.5))
	require.NoError(t, err)
	assert.True(t, value.IsFloat32())
	assert.Equal(t, float32(1.5), value.GetFloat32())

	value, err = handler.NewPlcValue(tag, "text")
	require.NoError(t, err)
	assert.Equal(t, "text", value.GetString())

	value, err = handler.NewPlcValue(tag, int64(-5))
	require.NoError(t, err)
	assert.Equal(t, int64(-5), value.GetInt64())

	_, err = handler.NewPlcValue(tag, struct{ Nope bool }{})
	assert.Error(t, err, "a value with no plc4x shape has to be refused, not guessed at")
}
