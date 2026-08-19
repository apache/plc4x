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

package firmata

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/firmata/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
)

// stubCodec is a minimal spi.MessageCodec which hands the callbacks passed to SendRequest back to
// the test and lets it push unsolicited messages in. The real codec is exercised against the test
// transport in MessageCodec_test.go; here it would only get in the way, as driving a board means
// driving both ends of the same transport at once.
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

func newStubCodec() *stubCodec {
	return &stubCodec{
		requests: make(chan stubRequest, 8),
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

// getSent hands back everything the connection has sent so far.
func (c *stubCodec) getSent() []spi.Message {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return append([]spi.Message(nil), c.sent...)
}

// failSends makes every following send fail.
func (c *stubCodec) failSends() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.sendErr = errors.New("the transport is gone")
}

// allowSends lets sends through again, so a test can retry what failed while the transport was gone.
func (c *stubCodec) allowSends() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.sendErr = nil
}

// sentBytesOf serializes everything the connection has sent, which is what would have been on the
// wire had the codec been a real one.
func (c *stubCodec) sentBytesOf(t *testing.T) []byte {
	t.Helper()
	var theBytes []byte
	for _, message := range c.getSent() {
		firmataMessage, ok := message.(readWriteModel.FirmataMessage)
		require.True(t, ok, "%T is not a firmata message", message)
		serialized, err := firmataMessage.Serialize()
		require.NoError(t, err)
		theBytes = append(theBytes, serialized...)
	}
	return theBytes
}
