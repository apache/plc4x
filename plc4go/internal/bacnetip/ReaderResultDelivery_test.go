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

package bacnetip

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transactions"
)

// captureCodec is a minimal spi.MessageCodec stub that hands the callbacks
// passed to SendRequest back to the test and returns a configurable error.
type captureCodec struct {
	sendRequestErr error
	handlers       chan capturedHandlers
}

type capturedHandlers struct {
	handleMessage spi.HandleMessage
	handleError   spi.HandleError
}

func newCaptureCodec(sendRequestErr error) *captureCodec {
	return &captureCodec{
		sendRequestErr: sendRequestErr,
		handlers:       make(chan capturedHandlers, 1),
	}
}

func (c *captureCodec) Connect(context.Context) error { return nil }
func (c *captureCodec) Disconnect() error             { return nil }
func (c *captureCodec) IsRunning() bool               { return true }
func (c *captureCodec) Send(context.Context, string, spi.Message) error {
	return nil
}
func (c *captureCodec) Expect(context.Context, string, spi.AcceptsMessage, spi.HandleMessage, spi.HandleError) {
}
func (c *captureCodec) SendRequest(_ context.Context, _ string, _ spi.Message, _ spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError) error {
	c.handlers <- capturedHandlers{handleMessage: handleMessage, handleError: handleError}
	return c.sendRequestErr
}
func (c *captureCodec) GetDefaultIncomingMessageChannel() chan spi.Message { return nil }

// A caller whose context expired abandons the result channel without reading.
// The send-failure result then fills the single-slot buffer; a later
// error-handler invocation (expectation timeout, disconnect fan-out) must not
// block forever on the full channel — those blocked handlers pile up in the
// codec's WaitGroup and wedge Disconnect indefinitely.
func TestReader_lateErrorHandlerAfterFailedSendMustNotBlock(t *testing.T) {
	codec := newCaptureCodec(errors.New("send failed: broken pipe"))
	tm := transactions.NewRequestTransactionManager(1)
	reader := NewReader(&InvokeIdGenerator{}, codec, tm, NewDriverContext(createDefaultConfiguration()), nil)

	objType := readWriteModel.BACnetObjectType_ANALOG_INPUT
	propId := readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE
	tag := plcTag{
		ObjectId:   objectId{ObjectIdType: &objType, ObjectIdInstance: 1},
		Properties: []property{{PropertyIdentifier: &propId}},
	}
	request := spiModel.NewDefaultPlcReadRequest(
		map[string]apiModel.PlcTag{"tag": tag}, []string{"tag"}, reader, nil)

	results := reader.Read(testutils.TestContext(t), request)

	var handlers capturedHandlers
	select {
	case handlers = <-codec.handlers:
	case <-time.After(time.Second):
		t.Fatal("SendRequest was never invoked")
	}

	// Wait for the send-failure result to occupy the channel buffer; the
	// abandoned caller never drains it.
	require.Eventually(t, func() bool { return len(results) == 1 },
		time.Second, time.Millisecond)

	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = handlers.handleError(errors.New("timeout"))
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("late error handler blocked on the abandoned result channel")
	}
}
