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

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

func testWriteRequest(t *testing.T, writer *Writer) apiModel.PlcWriteRequest {
	t.Helper()
	tag := NewTag(HoldingRegister, 1, 1, readWriteModel.ModbusDataType_UINT)
	return spiModel.NewDefaultPlcWriteRequest(
		map[string]apiModel.PlcTag{"tag": tag},
		[]string{"tag"},
		map[string]apiValues.PlcValue{"tag": spiValues.NewPlcUINT(42)},
		writer,
		nil,
	)
}

// A failed send must deliver an error result to the caller instead of only
// logging it — a caller without its own deadline otherwise waits forever on
// a channel that never receives anything.
func TestWriter_failedSendDeliversErrorResult(t *testing.T) {
	codec := newCaptureCodec(errors.New("send failed: broken pipe"))
	writer := NewWriter(1, codec)

	results := writer.Write(testutils.TestContext(t), testWriteRequest(t, writer))

	select {
	case result := <-results:
		if result.GetErr() == nil {
			t.Fatal("expected an error result for the failed send")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("no result delivered for a failed send")
	}
}

// When the caller has abandoned the result channel, a duplicate error-handler
// invocation (e.g. expectation timeout racing a handled message, or the
// disconnect fan-out) must not block forever on the full single-slot buffer.
func TestWriter_duplicateErrorHandlerMustNotBlock(t *testing.T) {
	codec := newCaptureCodec(nil)
	writer := NewWriter(1, codec)

	_ = writer.Write(testutils.TestContext(t), testWriteRequest(t, writer))

	var handlers capturedHandlers
	select {
	case handlers = <-codec.handlers:
	case <-time.After(time.Second):
		t.Fatal("SendRequest was never invoked")
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = handlers.handleError(errors.New("timeout"))
		_ = handlers.handleError(errors.New("disconnected"))
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("duplicate error handler blocked on the abandoned result channel")
	}
}
