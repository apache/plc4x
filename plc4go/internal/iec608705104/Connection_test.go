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

package iec608705104

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/iec608705104/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// nextRequest takes the next frame the connection sent through SendRequest, failing rather than
// blocking forever when nothing comes.
func nextRequest(t *testing.T, codec *stubCodec) stubRequest {
	t.Helper()
	select {
	case request := <-codec.requests:
		return request
	case <-time.After(20 * time.Second):
		require.FailNow(t, "no request was sent")
		return stubRequest{}
	}
}

// completeHandshake answers the two handshake round trips the way a station would, and hands back
// the connection once Connect has returned.
func completeHandshake(t *testing.T, connection *Connection, codec *stubCodec) {
	t.Helper()
	connectErrors := make(chan error, 1)
	go func() {
		connectErrors <- connection.Connect(testutils.TestContext(t))
	}()

	testFrame := nextRequest(t, codec)
	activation, ok := testFrame.message.(readWriteModel.APDUUFormatTestFrameActivation)
	require.True(t, ok, "%T is not a test frame activation", testFrame.message)
	assert.Equal(t, commandTestFrameActivation, activation.GetCommand())
	// Only the confirmation which answers this frame will do.
	assert.False(t, testFrame.acceptsMessage(readWriteModel.NewAPDUUFormatStartDataTransferConfirmation(0x0B)))
	testFrameConfirmation := readWriteModel.NewAPDUUFormatTestFrameConfirmation(commandTestFrameConfirmation)
	require.True(t, testFrame.acceptsMessage(testFrameConfirmation))
	require.NoError(t, testFrame.handleMessage(testFrameConfirmation))

	startDataTransfer := nextRequest(t, codec)
	startActivation, ok := startDataTransfer.message.(readWriteModel.APDUUFormatStartDataTransferActivation)
	require.True(t, ok, "%T is not a start-data-transfer activation", startDataTransfer.message)
	assert.Equal(t, commandStartDataTransferActivation, startActivation.GetCommand())
	assert.False(t, startDataTransfer.acceptsMessage(testFrameConfirmation))
	startConfirmation := readWriteModel.NewAPDUUFormatStartDataTransferConfirmation(0x0B)
	require.True(t, startDataTransfer.acceptsMessage(startConfirmation))
	require.NoError(t, startDataTransfer.handleMessage(startConfirmation))

	select {
	case err := <-connectErrors:
		require.NoError(t, err)
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the connect never returned")
	}
}

// newHandshakenConnection is a connection which has completed the handshake and is being fed by the
// stub codec.
func newHandshakenConnection(t *testing.T) (*Connection, *stubCodec) {
	t.Helper()
	_options := testutils.EnrichOptionsWithOptionsForTesting(t)
	codec := newStubCodec()
	connection := NewConnection(DefaultConfiguration(), codec, map[string][]string{}, NewTagHandler(_options...), _options...)
	completeHandshake(t, connection, codec)
	t.Cleanup(func() {
		assert.NoError(t, connection.Close())
	})
	return connection, codec
}

// pushIncoming hands the connection a frame the way the codec's worker would.
func pushIncoming(t *testing.T, codec *stubCodec, message spi.Message) {
	t.Helper()
	select {
	case codec.incoming <- message:
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the incoming channel never drained")
	}
}

// eventuallySent waits for the connection to have sent exactly the given control fields.
func eventuallySent(t *testing.T, codec *stubCodec, want []uint16) {
	t.Helper()
	assert.Eventually(t, func() bool {
		sent := codec.sentCommands()
		if len(sent) != len(want) {
			return false
		}
		for i, command := range want {
			if sent[i] != command {
				return false
			}
		}
		return true
	}, 20*time.Second, 5*time.Millisecond, "sent %v, wanted %v", codec.sentCommands(), want)
}

// The handshake is plc4j's: prove the link with a test frame, then start data transfer. Only once
// STARTDT is confirmed does the station send anything, so only then is the connection usable.
func TestConnection_ConnectHandshake(t *testing.T) {
	_options := testutils.EnrichOptionsWithOptionsForTesting(t)
	codec := newStubCodec()
	connection := NewConnection(DefaultConfiguration(), codec, map[string][]string{}, NewTagHandler(_options...), _options...)

	assert.False(t, connection.IsConnected())
	completeHandshake(t, connection, codec)
	assert.True(t, connection.IsConnected())

	require.NoError(t, connection.Close())
	assert.False(t, connection.IsConnected())
	// STOPDT is the orderly way out and goes out before the transport does. plc4j never sends it.
	assert.Equal(t, []uint16{commandTestFrameActivation, commandStartDataTransferActivation, commandStopDataTransferActivation},
		codec.sentCommands())
}

// A station which never confirms the test frame isn't a station we can work with, and the connect
// has to fail rather than hand out a connection which will never deliver anything.
func TestConnection_ConnectFailsWithoutATestFrameConfirmation(t *testing.T) {
	codec := newStubCodec()
	configuration := DefaultConfiguration()
	configuration.requestTimeout = 250 * time.Millisecond
	connection := NewConnection(configuration, codec, map[string][]string{}, NewTagHandler())

	err := connection.Connect(testutils.TestContext(t))

	assert.Error(t, err)
	assert.False(t, connection.IsConnected())
	assert.Equal(t, []uint16{commandTestFrameActivation}, codec.sentCommands(),
		"the test frame went out even though nothing came back")
}

// A station which confirms the test frame but never starts data transfer fails the connect just the
// same - it would stay silent forever.
func TestConnection_ConnectFailsWithoutAStartDataTransferConfirmation(t *testing.T) {
	codec := newStubCodec()
	configuration := DefaultConfiguration()
	configuration.requestTimeout = 250 * time.Millisecond
	connection := NewConnection(configuration, codec, map[string][]string{}, NewTagHandler())

	connectErrors := make(chan error, 1)
	go func() {
		connectErrors <- connection.Connect(testutils.TestContext(t))
	}()

	testFrame := nextRequest(t, codec)
	require.NoError(t, testFrame.handleMessage(readWriteModel.NewAPDUUFormatTestFrameConfirmation(commandTestFrameConfirmation)))

	select {
	case err := <-connectErrors:
		assert.Error(t, err)
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the connect never returned")
	}
	assert.False(t, connection.IsConnected())
}

// A handshake which errors out fails the connect rather than hanging - every path out of the
// exchange has to complete.
func TestConnection_ConnectFailsOnAHandshakeError(t *testing.T) {
	codec := newStubCodec()
	connection := NewConnection(DefaultConfiguration(), codec, map[string][]string{}, NewTagHandler())

	connectErrors := make(chan error, 1)
	go func() {
		connectErrors <- connection.Connect(testutils.TestContext(t))
	}()

	testFrame := nextRequest(t, codec)
	require.NoError(t, testFrame.handleError(errors.New("the station hung up")))

	select {
	case err := <-connectErrors:
		assert.Error(t, err)
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the connect never returned")
	}
	assert.False(t, connection.IsConnected())
}

// A transport which is gone before the first frame goes out fails the connect immediately.
func TestConnection_ConnectFailsWhenTheFrameCannotBeSent(t *testing.T) {
	codec := newStubCodec()
	codec.failSends()
	connection := NewConnection(DefaultConfiguration(), codec, map[string][]string{}, NewTagHandler())

	err := connection.Connect(testutils.TestContext(t))

	assert.Error(t, err)
	assert.False(t, connection.IsConnected())
}

// A station probes the link whenever its own idle timer fires, and closes the connection if the
// confirmation doesn't come back.
func TestConnection_ConfirmsATestFrame(t *testing.T) {
	connection, codec := newHandshakenConnection(t)
	_ = connection

	pushIncoming(t, codec, readWriteModel.NewAPDUUFormatTestFrameActivation(commandTestFrameActivation))

	eventuallySent(t, codec, []uint16{
		commandTestFrameActivation,
		commandStartDataTransferActivation,
		commandTestFrameConfirmation,
	})
}

// The acknowledgement window: after w unacknowledged I-format frames an S-format acknowledgement goes
// out, carrying the sequence number of the frame we expect next.
//
// plc4j gets both halves of that wrong - it echoes the incoming frame's *receive* sequence number
// (which is the station acknowledging our sends) and it doesn't shift it left by one - so the number
// it sends is neither the right one nor in the right encoding.
func TestConnection_AcknowledgesAfterTheWindow(t *testing.T) {
	connection, codec := newHandshakenConnection(t)

	sentFrames := 0
	for sendSequenceNo := range uint16(defaultAckThreshold) {
		// Every frame also carries the station's own receive sequence number, which has nothing to
		// do with what it has sent us. Making it differ from the send sequence number is what
		// catches a driver echoing the wrong one.
		frame := iFormatFrame(sendSequenceNo, 4242, asduBytes(0x01, 3, 10, informationObjectBytes(13, 0x01)))
		pushIncoming(t, codec, parseApdu(t, frame))
		sentFrames++
	}

	require.Equal(t, defaultAckThreshold, sentFrames)
	eventuallySent(t, codec, []uint16{
		commandTestFrameActivation,
		commandStartDataTransferActivation,
		commandSupervisoryFormat,
	})

	sent := codec.getSent()
	acknowledgement, ok := sent[len(sent)-1].(readWriteModel.APDUSFormat)
	require.True(t, ok, "%T is not an S-format acknowledgement", sent[len(sent)-1])
	// Eight frames numbered 0..7 arrived, so the next one expected is 8.
	assert.Equal(t, uint16(defaultAckThreshold), acknowledgement.GetReceiveSequenceNo()>>1)
	assert.Equal(t, uint16(defaultAckThreshold)<<1, acknowledgement.GetReceiveSequenceNo(),
		"the sequence number sits in the high 15 bits")
	_ = connection
}

// Below the window nothing is acknowledged - that is what the window is for.
func TestConnection_DoesNotAcknowledgeBeforeTheWindow(t *testing.T) {
	_, codec := newHandshakenConnection(t)

	for sendSequenceNo := range uint16(defaultAckThreshold - 1) {
		frame := iFormatFrame(sendSequenceNo, 0, asduBytes(0x01, 3, 10, informationObjectBytes(13, 0x01)))
		pushIncoming(t, codec, parseApdu(t, frame))
	}

	eventuallySent(t, codec, []uint16{commandTestFrameActivation, commandStartDataTransferActivation})
}

// The window is configurable, because a station which insists on a smaller one drops a connection
// which acknowledges too late. plc4j hard-codes 8.
func TestConnection_AcknowledgementWindowIsConfigurable(t *testing.T) {
	_options := testutils.EnrichOptionsWithOptionsForTesting(t)
	codec := newStubCodec()
	configuration := DefaultConfiguration()
	configuration.ackThreshold = 2
	connection := NewConnection(configuration, codec, map[string][]string{}, NewTagHandler(_options...), _options...)
	completeHandshake(t, connection, codec)
	t.Cleanup(func() { assert.NoError(t, connection.Close()) })

	for sendSequenceNo := range uint16(2) {
		frame := iFormatFrame(sendSequenceNo, 0, asduBytes(0x01, 3, 10, informationObjectBytes(13, 0x01)))
		pushIncoming(t, codec, parseApdu(t, frame))
	}

	eventuallySent(t, codec, []uint16{
		commandTestFrameActivation,
		commandStartDataTransferActivation,
		commandSupervisoryFormat,
	})
}

// The sequence number is 15 bits wide and wraps, so the frame after 32767 is frame 0.
func TestConnection_SequenceNumberWraps(t *testing.T) {
	connection, _ := newHandshakenConnection(t)

	next, due := connection.noteReceivedDataFrame(sequenceNumberMask)
	assert.Zero(t, next, "the frame after the last one is frame 0")
	assert.False(t, due)
}

// An S-format frame from the station acknowledges frames we sent. This driver sends none, so there
// is nothing to do but note it - and certainly nothing to answer with.
func TestConnection_IgnoresASupervisoryFrame(t *testing.T) {
	_, codec := newHandshakenConnection(t)

	pushIncoming(t, codec, readWriteModel.NewAPDUSFormat(commandSupervisoryFormat, 26))

	// Give the worker a moment to have done nothing at all.
	assert.Never(t, func() bool {
		return len(codec.sentCommands()) > 2
	}, 200*time.Millisecond, 10*time.Millisecond)
}

// Closing a connection which never completed its handshake must not try to send anything, and must
// not hang.
func TestConnection_CloseWithoutConnect(t *testing.T) {
	codec := newStubCodec()
	connection := NewConnection(DefaultConfiguration(), codec, map[string][]string{}, NewTagHandler())

	require.NoError(t, connection.Close())
	assert.Empty(t, codec.sentCommands())
}

// Closing twice is a no-op the second time round, and in particular does not send STOPDT twice.
func TestConnection_CloseIsIdempotent(t *testing.T) {
	_options := testutils.EnrichOptionsWithOptionsForTesting(t)
	codec := newStubCodec()
	connection := NewConnection(DefaultConfiguration(), codec, map[string][]string{}, NewTagHandler(_options...), _options...)
	completeHandshake(t, connection, codec)

	require.NoError(t, connection.Close())
	require.NoError(t, connection.Close())

	assert.Equal(t, []uint16{commandTestFrameActivation, commandStartDataTransferActivation, commandStopDataTransferActivation},
		codec.sentCommands())
}

func TestConnection_String(t *testing.T) {
	connection := NewConnection(DefaultConfiguration(), newStubCodec(), map[string][]string{}, NewTagHandler())
	assert.Contains(t, connection.String(), "iec608705104.Connection")
}
