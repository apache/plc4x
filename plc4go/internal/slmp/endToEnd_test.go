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

// The connection and subscriber tests drive a stub codec and the codec tests call Receive directly,
// so neither of them ever runs the codec's own workers against the connection. These tests do: real
// codec, real transport, bytes in and a plc4x value out.

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/slmp/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
)

// newRunningConnection connects a connection whose codec really runs. SLMP has no handshake, so
// there is nothing to answer before the first request.
func newRunningConnection(t *testing.T) (*Connection, *test.TransportInstance) {
	t.Helper()
	return newRunningConnectionWith(t, DefaultConfiguration())
}

func newRunningConnectionWith(t *testing.T, configuration Configuration) (*Connection, *test.TransportInstance) {
	t.Helper()
	_options := testutils.EnrichOptionsWithOptionsForTesting(t)
	codec, transportInstance := newTestCodec(t)
	tm := transactions.NewRequestTransactionManager(maxConcurrentRequests, _options...)
	t.Cleanup(func() { assert.NoError(t, tm.Close()) })
	connection := NewConnection(codec, configuration, NewTagHandler(), tm,
		map[string][]string{}, _options...)
	require.NoError(t, connection.Connect(testutils.TestContext(t)))
	t.Cleanup(func() { assert.NoError(t, connection.Close()) })
	return connection, transportInstance
}

// waitForSentFrame drains the whole frame the connection has just put on the wire. A frame is
// written in one go, so as soon as any byte shows up all of them are there.
func waitForSentFrame(t *testing.T, transportInstance *test.TransportInstance) []byte {
	t.Helper()
	require.Eventually(t, func() bool {
		return transportInstance.GetNumDrainableBytes() > 0
	}, 20*time.Second, time.Millisecond, "nothing was ever sent")
	return transportInstance.DrainWriteBuffer(transportInstance.GetNumDrainableBytes())
}

// pushFrame serializes a frame and hands it to the transport, the way a device would answer.
func pushFrame(t *testing.T, transportInstance *test.TransportInstance, frame readWriteModel.SlmpMessage) {
	t.Helper()
	theBytes, err := frame.Serialize()
	require.NoError(t, err)
	transportInstance.FillReadBuffer(theBytes)
}

// TestEndToEnd_ReadThroughTheRunningCodec is the whole cycle: request out as the reference frame from
// the SH-080008 Batch Read worked example, the manual's own response bytes back in, and the two
// register values it encodes delivered as a plc4x value.
func TestEndToEnd_ReadThroughTheRunningCodec(t *testing.T) {
	connection, transportInstance := newRunningConnection(t)

	readRequest, err := connection.ReadRequestBuilder().AddTagAddress("registers", "D350:WORD[2]").Build()
	require.NoError(t, err)
	results := readRequest.Execute(testutils.TestContext(t))

	assert.Equal(t, mustDecodeHex(t, batchReadRequestHex), waitForSentFrame(t, transportInstance),
		"the frame on the wire has to be the reference vector byte for byte")
	transportInstance.FillReadBuffer(mustDecodeHex(t, batchReadResponseHex))

	select {
	case result := <-results:
		require.NoError(t, result.GetErr())
		require.NotNil(t, result.GetResponse())
		assert.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("registers"))
		value := result.GetResponse().GetValue("registers")
		require.NotNil(t, value)
		require.True(t, value.IsList())
		require.Len(t, value.GetList(), 2)
		assert.Equal(t, uint16(0x56AB), value.GetList()[0].GetUint16(), "D350")
		assert.Equal(t, uint16(0x170F), value.GetList()[1].GetUint16(), "D351")
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the read never completed")
	}
}

func TestEndToEnd_WriteThroughTheRunningCodec(t *testing.T) {
	connection, transportInstance := newRunningConnection(t)

	writeRequest, err := connection.WriteRequestBuilder().
		AddTagAddress("registers", "D350:WORD[2]", []uint16{0x1234, 0x5678}).Build()
	require.NoError(t, err)
	results := writeRequest.Execute(testutils.TestContext(t))

	assert.Equal(t, mustDecodeHex(t, batchWriteRequestHex), waitForSentFrame(t, transportInstance))
	// A Batch Write success carries the end code and nothing else.
	pushFrame(t, transportInstance, readWriteModel.NewSlmpResponseFrame3E(0x0000, nil))

	select {
	case result := <-results:
		require.NoError(t, result.GetErr())
		require.NotNil(t, result.GetResponse())
		assert.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("registers"))
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the write never completed")
	}
}

// TestEndToEnd_StrayResponseIsDrained is the drain doing its job: a response frame that arrives
// while nothing is waiting for one is discarded instead of piling up in the codec's default incoming
// channel, where it would stay for the life of the connection (and, once 100 of them have piled up,
// turn into a codec warning per further frame).
//
// What the drain does not fix - and no driver can, because a 3E frame carries no correlation id - is
// a late response that arrives while the *next* request is in flight: that one is taken as the next
// request's answer. See the comment on acceptsAnyResponseFrame.
func TestEndToEnd_StrayResponseIsDrained(t *testing.T) {
	connection, transportInstance := newRunningConnection(t)
	incomingMessageChannel := connection.GetMessageCodec().GetDefaultIncomingMessageChannel()

	// Nothing is in flight, so this response can't be matched to anything.
	transportInstance.FillReadBuffer(mustDecodeHex(t, batchReadResponseHex))
	require.Eventually(t, func() bool {
		available, err := transportInstance.GetNumBytesAvailableInBuffer()
		return err == nil && available == 0
	}, 20*time.Second, time.Millisecond, "the codec never picked the stray frame up")

	// Give the frame time to land in the channel if nothing were draining it, then check that it
	// didn't stay there.
	time.Sleep(200 * time.Millisecond)
	assert.Empty(t, incomingMessageChannel, "a frame no request was waiting for has to be drained")
	assert.True(t, connection.IsConnected(), "a stray frame must not take the connection down")
}

// TestEndToEnd_ATimedOutReadDoesNotWedgeTheNextOne is the one thing the one-request-at-a-time design
// could get wrong: the transaction manager runs a single request at a time, so a read whose
// transaction is never released would block every read after it forever. A device that goes quiet
// therefore has to cost exactly one request, not the connection.
func TestEndToEnd_ATimedOutReadDoesNotWedgeTheNextOne(t *testing.T) {
	configuration := DefaultConfiguration()
	configuration.requestTimeout = 250 * time.Millisecond
	connection, transportInstance := newRunningConnectionWith(t, configuration)

	// Nothing answers this one.
	timingOut, err := connection.ReadRequestBuilder().AddTagAddress("hurz", "D350").Build()
	require.NoError(t, err)
	select {
	case result := <-timingOut.Execute(testutils.TestContext(t)):
		require.NoError(t, result.GetErr())
		require.NotNil(t, result.GetResponse())
		assert.Equal(t, apiModel.PlcResponseCode_REQUEST_TIMEOUT, result.GetResponse().GetResponseCode("hurz"))
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the read never completed")
	}
	waitForSentFrame(t, transportInstance)

	// And this one has to get through.
	second, err := connection.ReadRequestBuilder().AddTagAddress("hurz", "D350").Build()
	require.NoError(t, err)
	results := second.Execute(testutils.TestContext(t))
	waitForSentFrame(t, transportInstance)
	pushFrame(t, transportInstance, readWriteModel.NewSlmpResponseFrame3E(0x0000, []byte{0x2A, 0x00}))

	select {
	case result := <-results:
		require.NoError(t, result.GetErr())
		require.NotNil(t, result.GetResponse())
		assert.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("hurz"))
		assert.Equal(t, uint16(42), result.GetResponse().GetValue("hurz").GetUint16())
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the read after a timed-out one never completed")
	}
}
