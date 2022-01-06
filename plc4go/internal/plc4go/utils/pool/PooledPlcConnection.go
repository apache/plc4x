/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package pool

import (
	_default "github.com/apache/plc4x/plc4go/internal/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/pkg/plc4go"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/viney-shih/go-lock"
	"time"
)

type PooledPlcConnectionState int32

const (
	StateInitialized PooledPlcConnectionState = iota
	StateIdle
	StateInUse
	StateInvalid
)

type PooledPlcConnection struct {
	connectionString string
	// Reference to the pool (used for giving back connection)
	connectionPool *PlcConnectionPool

	// The lock for manipulating the pools state.
	lock lock.Mutex
	// The actual connection being pooled.
	activeConnection plc4go.PlcConnection
	// The current state of this connection.
	state PooledPlcConnectionState
	// Queue of waiting clients.
	queue []chan plc4go.PlcConnectionConnectResult
}

func NewPooledPlcConnection(connectionString string, connectionPool *PlcConnectionPool) *PooledPlcConnection {
	return &PooledPlcConnection{
		connectionString: connectionString,
		state:            StateInitialized,
		lock:             lock.NewCASMutex(),
		queue:            make([]chan plc4go.PlcConnectionConnectResult, 0),
		connectionPool:   connectionPool,
	}
}

func (t *PooledPlcConnection) setActiveConnection(activeConnection plc4go.PlcConnection) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.activeConnection = activeConnection
}

func (t *PooledPlcConnection) setState(state PooledPlcConnectionState) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.state = state
}

func (t *PooledPlcConnection) enqueue(ch chan plc4go.PlcConnectionConnectResult) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.queue = append(t.queue, ch)
}

func (t *PooledPlcConnection) Connect() <-chan plc4go.PlcConnectionConnectResult {
	if t.activeConnection == nil {
		panic("No active connection")
	}
	panic("Called 'Connect' on a pooled connection")
}

func (t *PooledPlcConnection) BlockingClose() {
	if t.activeConnection == nil {
		panic("No active connection")
	}

	// Call close and wait for the operation to finish.
	<-t.Close()
}

func (t *PooledPlcConnection) Close() <-chan plc4go.PlcConnectionCloseResult {
	if t.activeConnection == nil {
		panic("No active connection")
	}

	result := make(chan plc4go.PlcConnectionCloseResult)

	go func() {
		// Check if the connection is still alive, if it is, put it back into the pool
		pingResults := t.Ping()
		pingTimeout := time.NewTimer(time.Second * 5)
		select {
		case pingResult := <-pingResults:
			{
				if pingResult.GetErr() != nil {
					// Mark the connection as broken ...
					t.state = StateInvalid
				}
			}
		case <-pingTimeout.C:
			{
				// Mark the connection as broken ...
				t.state = StateInvalid
			}
		}

		// Return the connection to the pool and don't actually close it.
		err := t.connectionPool.returnConnection(t)

		// Finish closing the connection.
		result <- _default.NewDefaultPlcConnectionCloseResult(t, err)
	}()

	return result
}

func (t *PooledPlcConnection) IsConnected() bool {
	if t.activeConnection == nil {
		panic("No active connection")
	}
	return t.activeConnection.IsConnected()
}

func (t *PooledPlcConnection) Ping() <-chan plc4go.PlcConnectionPingResult {
	if t.activeConnection == nil {
		panic("No active connection")
	}
	return t.activeConnection.Ping()
}

func (t *PooledPlcConnection) GetMetadata() model.PlcConnectionMetadata {
	if t.activeConnection == nil {
		panic("No active connection")
	}
	return t.activeConnection.GetMetadata()
}

func (t *PooledPlcConnection) ReadRequestBuilder() model.PlcReadRequestBuilder {
	if t.activeConnection == nil {
		panic("No active connection")
	}
	return t.activeConnection.ReadRequestBuilder()
}

func (t *PooledPlcConnection) WriteRequestBuilder() model.PlcWriteRequestBuilder {
	if t.activeConnection == nil {
		panic("No active connection")
	}
	return t.activeConnection.WriteRequestBuilder()
}

func (t *PooledPlcConnection) SubscriptionRequestBuilder() model.PlcSubscriptionRequestBuilder {
	if t.activeConnection == nil {
		panic("No active connection")
	}
	return t.activeConnection.SubscriptionRequestBuilder()
}

func (t *PooledPlcConnection) UnsubscriptionRequestBuilder() model.PlcUnsubscriptionRequestBuilder {
	if t.activeConnection == nil {
		panic("No active connection")
	}
	return t.activeConnection.UnsubscriptionRequestBuilder()
}

func (t *PooledPlcConnection) BrowseRequestBuilder() model.PlcBrowseRequestBuilder {
	if t.activeConnection == nil {
		panic("No active connection")
	}
	return t.activeConnection.BrowseRequestBuilder()
}
