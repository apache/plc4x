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

// The connection and subscriber tests drive a stub codec and the codec tests call Receive directly,
// so neither of them ever runs the codec's own workers against the connection. These tests do: real
// codec, real transport, bytes in and a plc4x value out.

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/abeth/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
)

// newRunningConnection connects a connection whose codec really runs, answering the connect
// handshake the way a PLC would: the connection request goes out on the transport and the response
// carrying the session handle is pushed back once it has.
func newRunningConnection(t *testing.T) (*Connection, *test.TransportInstance) {
	t.Helper()
	_options := testutils.EnrichOptionsWithOptionsForTesting(t)
	codec, transportInstance := newTestCodec(t)
	tm := transactions.NewRequestTransactionManager(maxConcurrentRequests, _options...)
	t.Cleanup(func() { assert.NoError(t, tm.Close()) })
	configuration := DefaultConfiguration()
	driverContext, err := NewDriverContext(configuration)
	require.NoError(t, err)
	connection := NewConnection(codec, configuration, driverContext, NewTagHandler(), tm,
		map[string][]string{}, _options...)

	connectErrors := make(chan error, 1)
	go func() {
		connectErrors <- connection.Connect(testutils.TestContext(t))
	}()

	// Answering only once the request is really on the wire is what keeps this deterministic: by
	// then the expectation matching the response is registered.
	sent := waitForSentPacket(t, transportInstance)
	connectionRequest, err := readWriteModel.CIPEncapsulationPacketParse[readWriteModel.CIPEncapsulationConnectionRequest](t.Context(), sent)
	require.NoError(t, err, "the first packet has to be the connection request")
	pushPacket(t, transportInstance, readWriteModel.NewCIPEncapsulationConnectionResponse(
		testSessionHandle, 0, connectionRequest.GetSenderContext(), 0))

	select {
	case err := <-connectErrors:
		require.NoError(t, err)
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the connect never returned")
	}
	require.Equal(t, testSessionHandle, connection.session.getSessionHandle())
	t.Cleanup(func() { assert.NoError(t, connection.Close()) })
	return connection, transportInstance
}

// waitForSentPacket drains the whole packet the connection has just put on the wire. A packet is
// written in one go, so as soon as any byte shows up all of them are there.
func waitForSentPacket(t *testing.T, transportInstance *test.TransportInstance) []byte {
	t.Helper()
	require.Eventually(t, func() bool {
		return transportInstance.GetNumDrainableBytes() > 0
	}, 20*time.Second, time.Millisecond, "nothing was ever sent")
	return transportInstance.DrainWriteBuffer(transportInstance.GetNumDrainableBytes())
}

// pushPacket serializes a packet and hands it to the transport, the way a PLC would answer.
func pushPacket(t *testing.T, transportInstance *test.TransportInstance, packet readWriteModel.CIPEncapsulationPacket) {
	t.Helper()
	theBytes, err := packet.Serialize()
	require.NoError(t, err)
	transportInstance.FillReadBuffer(theBytes)
}

// A read has to travel through the running codec: request out, response in, value delivered.
func TestEndToEnd_ReadThroughTheRunningCodec(t *testing.T) {
	connection, transportInstance := newRunningConnection(t)

	readRequestBuilder := connection.ReadRequestBuilder()
	readRequestBuilder.AddTagAddress("counter", "N7:3:INTEGER[1]")
	readRequest, err := readRequestBuilder.Build()
	require.NoError(t, err)
	results := readRequest.Execute(testutils.TestContext(t))

	// The response has to carry the transaction counter of the request, which is the only thing
	// tying the two together.
	sent := waitForSentPacket(t, transportInstance)
	sentRead, err := readWriteModel.CIPEncapsulationPacketParse[readWriteModel.CIPEncapsulationReadRequest](t.Context(), sent)
	require.NoError(t, err, "%v is not a read request", sent)
	assert.Equal(t, testSessionHandle, sentRead.GetSessionHandle(), "every packet carries the session handle")
	pushPacket(t, transportInstance, readWriteModel.NewCIPEncapsulationReadResponse(
		testSessionHandle, 0, emptySenderContext, 0,
		readWriteModel.NewDF1CommandResponseMessageProtectedTypedLogicalRead(
			df1SourceAddress, DefaultConfiguration().station, 0,
			sentRead.GetRequest().GetTransactionCounter(), []uint8{42})))

	select {
	case result := <-results:
		require.NoError(t, result.GetErr())
		require.NotNil(t, result.GetResponse())
		assert.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("counter"))
		assert.Equal(t, int16(42), result.GetResponse().GetValue("counter").GetInt16())
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the read never completed")
	}
}
