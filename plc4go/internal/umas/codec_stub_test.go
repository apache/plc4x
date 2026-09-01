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
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/umas/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
)

// stubCodec is a minimal spi.MessageCodec which hands the callbacks passed to SendRequest back to the
// test. The real codec is exercised against the test transport in MessageCodec_test.go; here it would
// only get in the way, as driving a PLC means driving both ends of the same transport at once.
type stubCodec struct {
	mutex    sync.Mutex
	sent     []spi.Message
	sendErr  error
	requests chan stubRequest
	incoming chan spi.Message
}

type stubRequest struct {
	message        spi.Message
	acceptsMessage spi.AcceptsMessage
	handleMessage  spi.HandleMessage
	handleError    spi.HandleError
}

// adu is the request as what it really is, so a test can look at the UMAS item inside it.
func (r stubRequest) adu(t *testing.T) readWriteModel.ModbusTcpADU {
	t.Helper()
	adu, ok := r.message.(readWriteModel.ModbusTcpADU)
	require.True(t, ok, "%T is not a ModbusTcpADU", r.message)
	return adu
}

// item is the UMAS PDU item the request carries.
func (r stubRequest) item(t *testing.T) readWriteModel.UmasPDUItem {
	t.Helper()
	pdu, ok := r.adu(t).GetPdu().(readWriteModel.UmasPDU)
	require.True(t, ok, "%T is not a UMAS PDU", r.adu(t).GetPdu())
	return pdu.GetItem()
}

func newStubCodec() *stubCodec {
	return &stubCodec{
		// Deep enough to hold the whole connect handshake without a reader.
		requests: make(chan stubRequest, 64),
		incoming: make(chan spi.Message, 8),
	}
}

func (c *stubCodec) Connect(context.Context) error { return nil }

func (c *stubCodec) Disconnect() error { return nil }

func (c *stubCodec) IsRunning() bool { return true }

func (c *stubCodec) Send(_ context.Context, _ string, message spi.Message) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.sendErr != nil {
		return c.sendErr
	}
	c.sent = append(c.sent, message)
	return nil
}

func (c *stubCodec) Expect(context.Context, string, spi.AcceptsMessage, spi.HandleMessage, spi.HandleError) {
}

func (c *stubCodec) SendRequest(_ context.Context, _ string, message spi.Message, acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError) error {
	c.mutex.Lock()
	sendErr := c.sendErr
	if sendErr == nil {
		c.sent = append(c.sent, message)
	}
	c.mutex.Unlock()
	if sendErr != nil {
		return sendErr
	}
	c.requests <- stubRequest{
		message:        message,
		acceptsMessage: acceptsMessage,
		handleMessage:  handleMessage,
		handleError:    handleError,
	}
	return nil
}

func (c *stubCodec) GetDefaultIncomingMessageChannel() chan spi.Message {
	return c.incoming
}

// failSends makes every following send fail.
func (c *stubCodec) failSends() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.sendErr = errors.New("the transport is gone")
}

// nextRequest waits for the connection to have sent a request through SendRequest.
func (c *stubCodec) nextRequest(t *testing.T) stubRequest {
	t.Helper()
	select {
	case request := <-c.requests:
		return request
	case <-time.After(5 * time.Second):
		t.Fatal("no request was sent")
		return stubRequest{}
	}
}

// answerWith wraps a UMAS item in a response ADU carrying the request's own transaction identifier
// and hands it to the pending request, insisting that the request's accepts predicate really
// considers it an answer - which is what makes these tests catch a broken correlation.
func (c *stubCodec) answerWith(t *testing.T, request stubRequest, item readWriteModel.UmasPDUItem) {
	t.Helper()
	response := readWriteModel.NewModbusTcpADU(
		request.adu(t).GetTransactionIdentifier(),
		request.adu(t).GetUnitIdentifier(),
		readWriteModel.NewUmasPDU(item))
	require.True(t, request.acceptsMessage(response), "the request doesn't accept its own answer")
	require.NoError(t, request.handleMessage(response))
}

// answerWithPdu is answerWith for a response which isn't a UMAS PDU at all, e.g. a Modbus exception.
func (c *stubCodec) answerWithPdu(t *testing.T, request stubRequest, pdu readWriteModel.ModbusPDU) {
	t.Helper()
	response := readWriteModel.NewModbusTcpADU(
		request.adu(t).GetTransactionIdentifier(),
		request.adu(t).GetUnitIdentifier(),
		pdu)
	require.True(t, request.acceptsMessage(response), "the request doesn't accept its own answer")
	require.NoError(t, request.handleMessage(response))
}

// failRequest reports a transport error to the pending request.
func (c *stubCodec) failRequest(t *testing.T, request stubRequest, err error) {
	t.Helper()
	require.NoError(t, request.handleError(err))
}
