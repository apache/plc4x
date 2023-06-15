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

package spi

import (
	"context"
	"fmt"
	"time"
)

type Expectation interface {
	fmt.Stringer
	GetContext() context.Context
	GetExpiration() time.Time
	GetAcceptsMessage() AcceptsMessage
	GetHandleMessage() HandleMessage
	GetHandleError() HandleError
}

// AcceptsMessage If this function returns true, the message is forwarded to the message handler
type AcceptsMessage func(message Message) bool

// HandleMessage Function for handling the message, returns an error if anything goes wrong
type HandleMessage func(message Message) error

// HandleError Function for handling the message, returns an error if anything goes wrong
type HandleError func(err error) error

// MessageCodec handles sending and retrieving of messages
type MessageCodec interface {
	// Deprecated: use ConnectWithContext
	// Connect connects this codec
	Connect() error
	// ConnectWithContext connects this codec with the supplied context
	ConnectWithContext(ctx context.Context) error
	// Disconnect disconnects this codec
	Disconnect() error
	// IsRunning returns true if the codec (workers are running)
	IsRunning() bool

	// Send is sending a given message
	Send(message Message) error
	// Expect Wait for a given timespan for a message to come in, which returns 'true' for 'acceptMessage'
	// and is then forwarded to the 'handleMessage' function
	Expect(ctx context.Context, acceptsMessage AcceptsMessage, handleMessage HandleMessage, handleError HandleError, ttl time.Duration) error
	// SendRequest A combination that sends a message first and then waits for a response. !!!Important note: the callbacks are blocking calls
	SendRequest(ctx context.Context, message Message, acceptsMessage AcceptsMessage, handleMessage HandleMessage, handleError HandleError, ttl time.Duration) error

	// GetDefaultIncomingMessageChannel gives back the chan where unexpected messages arrive
	GetDefaultIncomingMessageChannel() chan Message
}
