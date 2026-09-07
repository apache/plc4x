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
	"context"

	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// captureCodec is a minimal spi.MessageCodec stub that hands the callbacks
// passed to SendRequest back to the test and returns a configurable error.
type captureCodec struct {
	sendRequestErr error
	handlers       chan capturedHandlers
}

type capturedHandlers struct {
	message        spi.Message
	acceptsMessage spi.AcceptsMessage
	handleMessage  spi.HandleMessage
	handleError    spi.HandleError
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
func (c *captureCodec) SendRequest(_ context.Context, _ string, message spi.Message, acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError) error {
	c.handlers <- capturedHandlers{
		message:        message,
		acceptsMessage: acceptsMessage,
		handleMessage:  handleMessage,
		handleError:    handleError,
	}
	return c.sendRequestErr
}
func (c *captureCodec) GetDefaultIncomingMessageChannel() chan spi.Message { return nil }

// notAnAdu is a spi.Message that is not a modbus ADU, used to make sure nothing in the driver
// type-asserts an incoming message without checking.
type notAnAdu struct{}

func (notAnAdu) String() string                                                    { return "notAnAdu" }
func (notAnAdu) Serialize() ([]byte, error)                                        { return nil, nil }
func (notAnAdu) SerializeWithWriteBuffer(context.Context, utils.WriteBuffer) error { return nil }
func (notAnAdu) GetLengthInBytes(context.Context) uint64                           { return 0 }
func (notAnAdu) GetLengthInBits(context.Context) uint64                            { return 0 }
