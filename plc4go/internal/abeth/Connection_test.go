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

package abeth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/abeth/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

const testSessionHandle = uint32(818)

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
	driverContext, err := NewDriverContext(configuration)
	require.NoError(t, err)
	connection := NewConnection(codec, configuration, driverContext, NewTagHandler(), tm, map[string][]string{}, _options...)
	return connection, codec
}

// connectHandshake answers the connection request the connect handshake sends.
func connectHandshake(t *testing.T, codec *stubCodec) {
	t.Helper()
	request := codec.nextRequest(t)
	connectionRequest, ok := request.message.(readWriteModel.CIPEncapsulationConnectionRequest)
	require.True(t, ok, "%T is not a connection request", request.message)
	assert.Equal(t, uint32(0), connectionRequest.GetSessionHandle(), "the handshake can't know a handle yet")
	assert.Equal(t, connectionRequestSenderContext, connectionRequest.GetSenderContext())

	codec.answer(t, request, readWriteModel.NewCIPEncapsulationConnectionResponse(
		testSessionHandle, 0, connectionRequestSenderContext, 0))
}

func TestConnection_ConnectRunsTheHandshake(t *testing.T) {
	connection, codec := newTestConnection(t, DefaultConfiguration())

	connectResult := make(chan error, 1)
	go func() {
		connectResult <- connection.Connect(testutils.TestContext(t))
	}()

	connectHandshake(t, codec)

	select {
	case err := <-connectResult:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("connect didn't finish")
	}
	assert.True(t, connection.IsConnected())
	assert.Equal(t, testSessionHandle, connection.session.getSessionHandle(),
		"the session handle from the response has to be remembered")
}

func TestConnection_ConnectFailsWhenTheHandshakeCannotBeSent(t *testing.T) {
	connection, codec := newTestConnection(t, DefaultConfiguration())
	codec.failSends()

	err := connection.Connect(testutils.TestContext(t))
	assert.Error(t, err)
	assert.Empty(t, codec.getSent())
}

func TestConnection_ConnectFailsWhenTheHandshakeIsRefused(t *testing.T) {
	connection, codec := newTestConnection(t, DefaultConfiguration())

	connectResult := make(chan error, 1)
	go func() {
		connectResult <- connection.Connect(testutils.TestContext(t))
	}()

	request := codec.nextRequest(t)
	require.NoError(t, request.handleError(assert.AnError))

	select {
	case err := <-connectResult:
		assert.Error(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("connect didn't finish")
	}
	assert.Equal(t, uint32(0), connection.session.getSessionHandle())
}

// TestConnection_ConnectTimesOut pins the request timeout down to a real bound: a PLC that accepts
// the TCP connection and then says nothing must not hang the caller forever.
func TestConnection_ConnectTimesOut(t *testing.T) {
	configuration := DefaultConfiguration()
	configuration.requestTimeout = 200 * time.Millisecond
	connection, codec := newTestConnection(t, configuration)

	start := time.Now()
	err := connection.Connect(testutils.TestContext(t))
	assert.Error(t, err)
	assert.Less(t, time.Since(start), 5*time.Second)
	// The request did go out, it just never got answered.
	assert.Len(t, codec.getSent(), 1)
}

// readTag runs one read through the connection, answering the request the reader sends with the
// passed payload and DF1 status.
func readTag(t *testing.T, address string, status uint8, data []uint8) (apiModel.PlcResponseCode, apiModel.PlcReadResponse) {
	t.Helper()
	connection, codec := newTestConnection(t, DefaultConfiguration())

	connectResult := make(chan error, 1)
	go func() {
		connectResult <- connection.Connect(testutils.TestContext(t))
	}()
	connectHandshake(t, codec)
	require.NoError(t, <-connectResult)

	readRequest, err := connection.ReadRequestBuilder().AddTagAddress("hurz", address).Build()
	require.NoError(t, err)
	resultChannel := readRequest.Execute(testutils.TestContext(t))

	request := codec.nextRequest(t)
	readPacket, ok := request.message.(readWriteModel.CIPEncapsulationReadRequest)
	require.True(t, ok, "%T is not a read request", request.message)
	assert.Equal(t, testSessionHandle, readPacket.GetSessionHandle(), "a read has to carry the session handle")
	df1Request, ok := readPacket.GetRequest().(readWriteModel.DF1CommandRequestMessage)
	require.True(t, ok, "%T is not a DF1 command request", readPacket.GetRequest())
	assert.Equal(t, uint16(1), df1Request.GetTransactionCounter(), "the first read of a connection is transaction 1")
	assert.Equal(t, df1SourceAddress, df1Request.GetSourceAddress())

	// A response carrying a different transaction counter must not be taken as the answer.
	assert.False(t, request.acceptsMessage(readWriteModel.NewCIPEncapsulationReadResponse(
		testSessionHandle, 0, emptySenderContext, 0,
		readWriteModel.NewDF1CommandResponseMessageProtectedTypedLogicalRead(
			df1SourceAddress, 0, 0, df1Request.GetTransactionCounter()+1, data))),
		"a response for another transaction must not be accepted")

	codec.answer(t, request, readWriteModel.NewCIPEncapsulationReadResponse(
		testSessionHandle, 0, emptySenderContext, 0,
		readWriteModel.NewDF1CommandResponseMessageProtectedTypedLogicalRead(
			df1SourceAddress, 0, status, df1Request.GetTransactionCounter(), data)))

	select {
	case result := <-resultChannel:
		require.NoError(t, result.GetErr())
		response := result.GetResponse()
		require.NotNil(t, response)
		return response.GetResponseCode("hurz"), response
	case <-time.After(5 * time.Second):
		t.Fatal("the read didn't finish")
		return 0, nil
	}
}

func TestConnection_ReadDecodesTheFileTypes(t *testing.T) {
	tests := []struct {
		name      string
		address   string
		status    uint8
		data      []uint8
		wantCode  apiModel.PlcResponseCode
		wantValue apiValues.PlcValue
	}{
		{
			name:      "a one byte integer is the value itself",
			address:   "N7:3:INTEGER[1]",
			data:      []uint8{42},
			wantCode:  apiModel.PlcResponseCode_OK,
			wantValue: spiValues.NewPlcINT(42),
		},
		{
			name:     "a wider integer comes back as the list of its bytes",
			address:  "N7:3:INTEGER[2]",
			data:     []uint8{42, 1},
			wantCode: apiModel.PlcResponseCode_OK,
			// plc4j hands the byte list to the value handler, which builds a list of INTs from it.
			wantValue: spiValues.NewPlcList([]apiValues.PlcValue{spiValues.NewPlcINT(42), spiValues.NewPlcINT(1)}),
		},
		{
			name:      "a word is two little-endian bytes",
			address:   "N7:3:WORD",
			data:      []uint8{0x2A, 0x01},
			wantCode:  apiModel.PlcResponseCode_OK,
			wantValue: spiValues.NewPlcWORD(298),
		},
		{
			name:      "a word keeps the bits of a value with the top bit set",
			address:   "N7:3:WORD",
			data:      []uint8{0xFB, 0xFF},
			wantCode:  apiModel.PlcResponseCode_OK,
			wantValue: spiValues.NewPlcWORD(65531),
		},
		{
			name:      "a dword is four little-endian bytes",
			address:   "N7:3:DWORD",
			data:      []uint8{0x01, 0x02, 0x03, 0x04},
			wantCode:  apiModel.PlcResponseCode_OK,
			wantValue: spiValues.NewPlcDWORD(67305985),
		},
		{
			name:      "a singlebit picks its bit out of the low byte",
			address:   "N7:3/3:SINGLEBIT",
			data:      []uint8{0x08, 0x00},
			wantCode:  apiModel.PlcResponseCode_OK,
			wantValue: spiValues.NewPlcBOOL(true),
		},
		{
			name:      "a singlebit picks its bit out of the high byte",
			address:   "N7:3/11:SINGLEBIT",
			data:      []uint8{0x00, 0x08},
			wantCode:  apiModel.PlcResponseCode_OK,
			wantValue: spiValues.NewPlcBOOL(true),
		},
		{
			name:      "a singlebit that isn't set is false",
			address:   "N7:3/11:SINGLEBIT",
			data:      []uint8{0x08, 0x00},
			wantCode:  apiModel.PlcResponseCode_OK,
			wantValue: spiValues.NewPlcBOOL(false),
		},
		{
			name:     "a non-zero DF1 status is reported rather than decoded",
			address:  "N7:3:WORD",
			status:   0xF0,
			data:     []uint8{0x2A, 0x01},
			wantCode: apiModel.PlcResponseCode_NOT_FOUND,
		},
		{
			name:     "a word without enough payload is invalid data",
			address:  "N7:3:WORD",
			data:     []uint8{0x2A},
			wantCode: apiModel.PlcResponseCode_INVALID_DATA,
		},
		{
			name:     "a dword without enough payload is invalid data",
			address:  "N7:3:DWORD",
			data:     []uint8{0x01, 0x02, 0x03},
			wantCode: apiModel.PlcResponseCode_INVALID_DATA,
		},
		{
			name:     "an integer without any payload is invalid data",
			address:  "N7:3:INTEGER[1]",
			data:     nil,
			wantCode: apiModel.PlcResponseCode_INVALID_DATA,
		},
		{
			name:     "a file type with no decoder is reported as unsupported",
			address:  "N7:3:STATUS[2]",
			data:     []uint8{0x01, 0x02},
			wantCode: apiModel.PlcResponseCode_UNSUPPORTED,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			code, response := readTag(t, testCase.address, testCase.status, testCase.data)
			assert.Equal(t, testCase.wantCode, code)
			if testCase.wantValue == nil {
				assert.Nil(t, response.GetValue("hurz"), "a failed read must not carry a value")
				return
			}
			value := response.GetValue("hurz")
			require.NotNil(t, value)
			assert.Equal(t, testCase.wantValue, value)
		})
	}
}

// TestConnection_ReadsTagsSequentially pins down that a multi-tag request turns into one request per
// tag, each with its own transaction counter: ab-eth answers one request at a time and the DF1
// command can only ever address a single element.
func TestConnection_ReadsTagsSequentially(t *testing.T) {
	connection, codec := newTestConnection(t, DefaultConfiguration())
	connectResult := make(chan error, 1)
	go func() {
		connectResult <- connection.Connect(testutils.TestContext(t))
	}()
	connectHandshake(t, codec)
	require.NoError(t, <-connectResult)

	readRequest, err := connection.ReadRequestBuilder().
		AddTagAddress("first", "N7:1:WORD").
		AddTagAddress("second", "N7:2:WORD").
		Build()
	require.NoError(t, err)
	resultChannel := readRequest.Execute(testutils.TestContext(t))

	for i, expectedElement := range []uint8{1, 2} {
		request := codec.nextRequest(t)
		readPacket, ok := request.message.(readWriteModel.CIPEncapsulationReadRequest)
		require.True(t, ok, "%T is not a read request", request.message)
		df1Request, ok := readPacket.GetRequest().(readWriteModel.DF1CommandRequestMessage)
		require.True(t, ok)
		assert.Equal(t, uint16(i+1), df1Request.GetTransactionCounter(), "every read gets a fresh counter")
		logicalRead, ok := df1Request.GetCommand().(readWriteModel.DF1RequestProtectedTypedLogicalRead)
		require.True(t, ok)
		assert.Equal(t, expectedElement, logicalRead.GetElementNumber())
		assert.Equal(t, uint8(2), logicalRead.GetByteSize(), "a WORD tag always asks for two bytes")
		assert.Equal(t, uint8(0x89), logicalRead.GetFileType(), "a WORD is read out of an integer file")

		codec.answer(t, request, readWriteModel.NewCIPEncapsulationReadResponse(
			testSessionHandle, 0, emptySenderContext, 0,
			readWriteModel.NewDF1CommandResponseMessageProtectedTypedLogicalRead(
				df1SourceAddress, 0, 0, df1Request.GetTransactionCounter(), []uint8{expectedElement, 0})))
	}

	select {
	case result := <-resultChannel:
		require.NoError(t, result.GetErr())
		response := result.GetResponse()
		require.NotNil(t, response)
		assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("first"))
		assert.Equal(t, spiValues.NewPlcWORD(1), response.GetValue("first"))
		assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("second"))
		assert.Equal(t, spiValues.NewPlcWORD(2), response.GetValue("second"))
	case <-time.After(5 * time.Second):
		t.Fatal("the read didn't finish")
	}
}

// TestConnection_ReadTimesOut checks that a read of a tag the PLC never answers reports a timeout for
// that tag rather than hanging or failing the whole request.
func TestConnection_ReadTimesOut(t *testing.T) {
	configuration := DefaultConfiguration()
	configuration.requestTimeout = 200 * time.Millisecond
	connection, codec := newTestConnection(t, configuration)
	connectResult := make(chan error, 1)
	go func() {
		connectResult <- connection.Connect(testutils.TestContext(t))
	}()
	connectHandshake(t, codec)
	require.NoError(t, <-connectResult)

	readRequest, err := connection.ReadRequestBuilder().AddTagAddress("hurz", "N7:3:WORD").Build()
	require.NoError(t, err)
	resultChannel := readRequest.Execute(testutils.TestContext(t))
	// Take the request off the stub's queue without answering it.
	codec.nextRequest(t)

	select {
	case result := <-resultChannel:
		require.NoError(t, result.GetErr())
		response := result.GetResponse()
		require.NotNil(t, response)
		assert.Equal(t, apiModel.PlcResponseCode_REQUEST_TIMEOUT, response.GetResponseCode("hurz"))
		assert.Nil(t, response.GetValue("hurz"))
	case <-time.After(5 * time.Second):
		t.Fatal("the read didn't finish")
	}
}

func TestSession_TransactionCounterCycles(t *testing.T) {
	s := newSession()
	assert.Equal(t, uint16(1), s.nextTransactionCounter(), "the first counter is 1, never 0")
	assert.Equal(t, uint16(2), s.nextTransactionCounter())

	// Just before the wrap.
	s.transactionCounter.Store(maxTransactionCounter - 1)
	assert.Equal(t, uint16(maxTransactionCounter), s.nextTransactionCounter())
	assert.Equal(t, uint16(1), s.nextTransactionCounter(), "the counter restarts at 1, skipping 0")
}

func TestConnection_CloseWithoutConnect(t *testing.T) {
	connection, _ := newTestConnection(t, DefaultConfiguration())
	// Closing a connection that was never connected has to be harmless: the driver closes on every
	// failed connect attempt.
	assert.NoError(t, connection.Close())
}

// strayMessage is a packet no request is waiting for: what the codec hands to its default incoming
// message channel and what the drain is there to throw away.
func strayMessage() spi.Message {
	return readWriteModel.NewCIPEncapsulationReadResponse(
		testSessionHandle, 0, emptySenderContext, 0,
		readWriteModel.NewDF1CommandResponseMessageProtectedTypedLogicalRead(
			df1SourceAddress, DefaultConfiguration().station, 0, 1, []uint8{0}))
}

// requireNothingDrains insists that nothing is reading the codec's default incoming message channel
// any more: a message handed to it stays there. Both the close and the failed-connect path stop the
// drain synchronously, so a single check is enough - there is no goroutine left that could still
// take the message afterwards.
func requireNothingDrains(t *testing.T, codec *stubCodec) {
	t.Helper()
	incoming := codec.GetDefaultIncomingMessageChannel()
	require.Empty(t, incoming, "the drain should have emptied the channel while it was running")
	incoming <- strayMessage()
	assert.Len(t, incoming, 1, "something is still draining the channel")
}

// TestConnection_StrayMessageDrainStopsOnClose pins the drain's lifetime to the connection's: it has
// to be gone once Close returned, rather than sitting on the codec's channel forever.
func TestConnection_StrayMessageDrainStopsOnClose(t *testing.T) {
	connection, codec := newTestConnection(t, DefaultConfiguration())

	connectResult := make(chan error, 1)
	go func() {
		connectResult <- connection.Connect(testutils.TestContext(t))
	}()
	connectHandshake(t, codec)
	require.NoError(t, <-connectResult)

	// While the connection is up the drain takes everything that lands on the channel.
	codec.GetDefaultIncomingMessageChannel() <- strayMessage()
	require.Eventually(t, func() bool {
		return len(codec.GetDefaultIncomingMessageChannel()) == 0
	}, 5*time.Second, time.Millisecond, "the stray message was never drained")

	require.NoError(t, connection.Close())
	requireNothingDrains(t, codec)
}

// TestConnection_StrayMessageDrainStopsWhenTheHandshakeFails covers the path with no Close in it: the
// driver hands the caller an error instead of a connection, so nobody will ever close this one and
// the drain has to have stopped itself.
func TestConnection_StrayMessageDrainStopsWhenTheHandshakeFails(t *testing.T) {
	connection, codec := newTestConnection(t, DefaultConfiguration())

	connectResult := make(chan error, 1)
	go func() {
		connectResult <- connection.Connect(testutils.TestContext(t))
	}()
	request := codec.nextRequest(t)
	require.NoError(t, request.handleError(assert.AnError))
	require.Error(t, <-connectResult)

	requireNothingDrains(t, codec)
}

// TestConnection_StrayMessageDrainIsNotDoubledByASecondConnect makes sure a reconnect replaces the
// drain instead of stacking a second one on top of it, and that closing still stops all of them.
func TestConnection_StrayMessageDrainIsNotDoubledByASecondConnect(t *testing.T) {
	connection, codec := newTestConnection(t, DefaultConfiguration())

	for i := 0; i < 3; i++ {
		connectResult := make(chan error, 1)
		go func() {
			connectResult <- connection.Connect(testutils.TestContext(t))
		}()
		connectHandshake(t, codec)
		require.NoError(t, <-connectResult)
	}

	require.NoError(t, connection.Close())
	requireNothingDrains(t, codec)
}

// TestConnection_StrayMessageDrainKeepsUpWithMoreThanTheChannelHolds is the finding itself: without a
// drain the channel filled up after its last slot was taken and the codec logged a warning for every
// further stray packet for the rest of the connection's life. With one, far more packets than the
// channel can hold go through it without it ever filling up.
func TestConnection_StrayMessageDrainKeepsUpWithMoreThanTheChannelHolds(t *testing.T) {
	connection, codec := newTestConnection(t, DefaultConfiguration())

	connectResult := make(chan error, 1)
	go func() {
		connectResult <- connection.Connect(testutils.TestContext(t))
	}()
	connectHandshake(t, codec)
	require.NoError(t, <-connectResult)

	incoming := codec.GetDefaultIncomingMessageChannel()
	for i := 0; i < 5*cap(incoming); i++ {
		select {
		case incoming <- strayMessage():
		case <-time.After(5 * time.Second):
			require.FailNow(t, "the channel filled up, so the stray messages aren't being drained")
		}
	}
	require.Eventually(t, func() bool {
		return len(incoming) == 0
	}, 5*time.Second, time.Millisecond, "the stray messages were never drained")

	require.NoError(t, connection.Close())
}
