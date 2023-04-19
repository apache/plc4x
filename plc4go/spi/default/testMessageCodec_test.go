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

package _default

import (
	"context"
	"github.com/apache/plc4x/plc4go/spi/transports"

	"github.com/apache/plc4x/plc4go/spi"
	"time"
)

type testMessageCodec struct {
}

func (t testMessageCodec) Connect() error {
	// NO-OP
	return nil
}

func (t testMessageCodec) ConnectWithContext(ctx context.Context) error {
	// NO-OP
	return nil
}

func (t testMessageCodec) Disconnect() error {
	// NO-OP
	return nil
}

func (t testMessageCodec) IsRunning() bool {
	// NO-OP
	return false
}

func (t testMessageCodec) Send(message spi.Message) error {
	// NO-OP
	return nil
}

func (t testMessageCodec) Expect(ctx context.Context, acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError, ttl time.Duration) error {
	// NO-OP
	return nil
}

func (t testMessageCodec) SendRequest(ctx context.Context, message spi.Message, acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError, ttl time.Duration) error {
	// NO-OP
	return nil
}

func (t testMessageCodec) GetDefaultIncomingMessageChannel() chan spi.Message {
	// NO-OP
	return nil
}

func (t testMessageCodec) GetTransportInstance() transports.TransportInstance {
	return testTransportInstance{}
}
