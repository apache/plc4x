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

package modbus

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// One address the driver can't parse costs its own tag and nothing else. Addresses regularly come
// from somewhere other than the caller - a device profile handed down from a server - so one typo
// among hundreds of good points is ordinary; failing the Build would turn it into an entire scan
// that reads nothing. plc4j answers the same request the same way, by keeping the rejected tag in
// the request with INVALID_ADDRESS (DefaultPlcReadRequest.Builder).
func TestReadRequestBuilder_anUnparseableAddressCostsOnlyItsOwnTag(t *testing.T) {
	codec := newCaptureCodec(nil)
	connection := testConnection(DefaultConfiguration(), codec)

	request, err := connection.ReadRequestBuilder().
		AddTagAddress("good1", "holding-register:1:REAL").
		AddTagAddress("bad", "holding-register:nonsense:REAL").
		AddTagAddress("good2", "holding-register:3:REAL").
		Build()
	require.NoError(t, err, "the bad address must not take the request down with it")
	require.NotNil(t, request)

	results := request.Execute(testutils.TestContext(t))

	// The two good tags are adjacent, so they are read as one block of four registers; the bad one
	// contributes nothing to it, because there is no address to ask about.
	select {
	case handlers := <-codec.handlers:
		requestAdu, ok := handlers.message.(readWriteModel.ModbusTcpADU)
		require.True(t, ok, "expected a ModbusTcpADU, got %T", handlers.message)
		pdu := holdingRegisterRequestOf(t, requestAdu)
		assert.Equal(t, uint16(0), pdu.GetStartingAddress())
		assert.Equal(t, uint16(4), pdu.GetQuantity())
		responseAdu := readWriteModel.NewModbusTcpADU(
			requestAdu.GetTransactionIdentifier(), requestAdu.GetUnitIdentifier(),
			readWriteModel.NewModbusPDUReadHoldingRegistersResponse(append(append([]byte{}, pi...), pi...)))
		require.True(t, handlers.acceptsMessage(responseAdu), "the reader didn't accept the response to its own request")
		require.NoError(t, handlers.handleMessage(responseAdu))
	case <-time.After(time.Second):
		t.Fatal("the good tags were never read")
	}

	var result apiModel.PlcReadRequestResult
	select {
	case result = <-results:
	case <-time.After(2 * time.Second):
		t.Fatal("no result was delivered")
	}

	// The rejection is reported as well as being visible per tag, the way a failed block is.
	require.Error(t, result.GetErr())
	assert.Contains(t, result.GetErr().Error(), "bad", "the error has to name the tag that was rejected")

	response := result.GetResponse()
	require.NotNil(t, response)
	assert.Equal(t, apiModel.PlcResponseCode_INVALID_ADDRESS, response.GetResponseCode("bad"))
	value := response.GetValue("bad")
	require.NotNil(t, value)
	assert.Equal(t, apiValues.NULL, value.GetPlcValueType())

	for _, tagName := range []string{"good1", "good2"} {
		assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode(tagName), tagName)
		assert.InDelta(t, 3.1415927, response.GetValue(tagName).GetFloat32(), 0.0000001, tagName)
	}
}

// The rejected tag keeps the address it was given, so that a caller holding only the response can
// still say which address was refused.
func TestReadRequestBuilder_theRejectedTagKeepsItsAddress(t *testing.T) {
	codec := newCaptureCodec(nil)
	connection := testConnection(DefaultConfiguration(), codec)

	request, err := connection.ReadRequestBuilder().
		AddTagAddress("bad", "not-an-address").
		Build()
	require.NoError(t, err)

	assert.Equal(t, []string{"bad"}, request.GetTagNames())
	assert.Equal(t, "not-an-address", request.GetTag("bad").GetAddressString())
	assert.Equal(t, apiValues.NULL, request.GetTag("bad").GetValueType())
}
