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

	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// A caller whose context expired abandons the result channel without reading.
// The send-failure result then fills the single-slot buffer; when the still
// registered expectation later times out, its error handler must not block
// forever on the full channel — those blocked handlers pile up in the codec's
// WaitGroup and wedge Disconnect indefinitely.
func TestReader_lateTimeoutAfterFailedSendMustNotBlock(t *testing.T) {
	codec := newCaptureCodec(errors.New("send failed: broken pipe"))
	reader := NewReader(1, codec)
	tag := NewTag(HoldingRegister, 1, 1, readWriteModel.ModbusDataType_UINT)
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
		t.Fatal("late timeout handler blocked on the abandoned result channel")
	}
}
