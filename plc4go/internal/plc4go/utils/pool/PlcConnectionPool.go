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
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/pkg/plc4go"
	"github.com/pkg/errors"
	"github.com/viney-shih/go-lock"
	"time"
)

type PlcConnectionPool struct {
	driverManager plc4go.PlcDriverManager

	// Maximum duration a connection can be used per lease.
	// If the connection is used for a longer time, it is forcefully removed from the client.
	maxLeaseTime time.Duration
	maxWaitTime  time.Duration

	poolLock    lock.RWMutex
	connections map[string]*ConnectionContainer
	tracer      *spi.Tracer
}

func NewPlcConnectionPool(driverManager plc4go.PlcDriverManager) *PlcConnectionPool {
	return NewPlcConnectionPoolWithMaxLeaseTime(driverManager, time.Second*5)
}

func NewPlcConnectionPoolWithMaxLeaseTime(driverManager plc4go.PlcDriverManager, maxLeaseTime time.Duration) *PlcConnectionPool {
	return &PlcConnectionPool{
		driverManager: driverManager,
		maxLeaseTime:  maxLeaseTime,
		maxWaitTime:   maxLeaseTime * 5,
		poolLock:      lock.NewCASMutex(),
		connections:   make(map[string]*ConnectionContainer),
		tracer:        nil,
	}
}

func (t *PlcConnectionPool) EnableTracer() {
	t.tracer = spi.NewTracer("pool")
}

func (t *PlcConnectionPool) GetTracer() *spi.Tracer {
	return t.tracer
}

func (t *PlcConnectionPool) GetConnection(connectionString string) <-chan plc4go.PlcConnectionConnectResult {
	ch := make(chan plc4go.PlcConnectionConnectResult)

	go func() {
		t.poolLock.Lock()

		// If a connection for this connection string didn't exist yet, create a new container
		// and make that container connect.
		if _, ok := t.connections[connectionString]; !ok {
			if t.tracer != nil {
				t.tracer.AddTrace("get-connection", "create new pooled connection")
			}
			// Create a new connection container.
			connectionContainer := &ConnectionContainer{
				driverManager:    t.driverManager,
				connectionString: connectionString,
				lock:             lock.NewCASMutex(),
				leaseCounter:     0,
				state:            StateInitialized,
				queue:            []chan plc4go.PlcConnectionConnectResult{},
			}
			// Register for connection events (Like connection closed or error).
			connectionContainer.addListener(t)
			// Store the new connection container in the pool of connections.
			t.connections[connectionString] = connectionContainer
			// Initialize the connection itself.
			go func() {
				connectionContainer.connect()
			}()
		}

		// Get the ConnectionContainer for this connection string.
		connection := t.connections[connectionString]

		// Release the lock again.
		t.poolLock.Unlock()

		// Try to get a lease on this connection.
		var txId string
		if t.tracer != nil {
			txId = t.tracer.AddTransactionalStartTrace("get-connection", "lease")
		}
		leaseChan := connection.lease()
		select {
		// Wait till we get a lease.
		case connectionResponse := <-leaseChan:
			if t.tracer != nil {
				t.tracer.AddTransactionalTrace(txId, "get-connection", "success")
			}
			select {
			case ch <- connectionResponse:
			case <-time.After(10 * time.Millisecond):
				// Log a message, that the client has given up
				t.tracer.AddTransactionalTrace(txId, "get-connection", "client given up")
				close(ch)
				// Return the connection to give another connection the chance to use it.
				connectionResponse.GetConnection().Close()
			}

		// Timeout after the maximum waiting time.
		case <-time.After(t.maxWaitTime):
			if t.tracer != nil {
				t.tracer.AddTransactionalTrace(txId, "get-connection", "timeout")
			}
			ch <- _default.NewDefaultPlcConnectionCloseResult(nil, errors.New("timeout while waiting for connection"))
		}
	}()

	return ch
}

// onConnectionEvent: Callback called by the connection container to signal connection events
// that have an impact on the pool itself (Like connections being permanently closed).
func (t *PlcConnectionPool) onConnectionEvent(event ConnectionEvent) {
	connectionContainer := event.GetConnectionContainer()
	if errorEvent, ok := event.(ConnectionErrorEvent); ok {
		if t.tracer != nil {
			t.tracer.AddTrace("destroy-connection", errorEvent.GetError().Error())
		}
		delete(t.connections, connectionContainer.connectionString)
	}
}

type ConnectionListener interface {
	onConnectionEvent(event ConnectionEvent)
}

type ConnectionEvent interface {
	GetConnectionContainer() ConnectionContainer
}

type ConnectionErrorEvent struct {
	conn ConnectionContainer
	err  error
}

func (c ConnectionErrorEvent) GetConnectionContainer() ConnectionContainer {
	return c.conn
}

func (c ConnectionErrorEvent) GetError() error {
	return c.err
}

type ConnectionContainer struct {
	lock             lock.RWMutex
	connectionString string
	driverManager    plc4go.PlcDriverManager
	tracerEnabled    bool
	connection       spi.PlcConnection
	leaseCounter     uint32
	// The current state of this connection.
	state PooledPlcConnectionState
	// Queue of waiting clients.
	queue []chan plc4go.PlcConnectionConnectResult
	// Listeners for connection events.
	listeners []ConnectionListener
}

func (t *ConnectionContainer) connect() {
	// Initialize the new connection.
	connectionResultChan := t.driverManager.GetConnection(t.connectionString)

	// Allow us to finish this function and return the lock quickly
	// Wait for the connection to be established.
	// TODO: Add some timeout handling.
	connectionResult := <-connectionResultChan

	// Get the lock.
	t.lock.Lock()
	defer t.lock.Unlock()

	// If the connection was successful, pass the active connection into the container.
	// If something went wrong, we have to remove the connection from the pool and return the error.
	if connectionResult.GetErr() == nil {
		// Inject the real connection into the container.
		t.connection = connectionResult.GetConnection().(spi.PlcConnection)
		t.tracerEnabled = t.connection.IsTraceEnabled()
		// Mark the connection as idle for now.
		t.state = StateIdle
		// If there is a request in the queue, hand out the connection to that.
		if len(t.queue) > 0 {
			// Get the first in the queue.
			queueHead := t.queue[0]
			t.queue = t.queue[1:]
			// Mark the connection as being used.
			t.state = StateInUse
			// Return the lease to the caller.
			connection := NewLeasedPlcConnection(t, t.leaseCounter, t.connection)
			// In this case we don't need to check for blocks
			// as the getConnection function of the connection pool
			// is definitely eagerly waiting for input.
			queueHead <- _default.NewDefaultPlcConnectionConnectResult(connection, nil)
		}
	} else {
		// Tell the connection pool that the connection is no longer available.
		if t.listeners != nil {
			event := ConnectionErrorEvent{
				conn: *t,
				err:  connectionResult.GetErr(),
			}
			for _, listener := range t.listeners {
				listener.onConnectionEvent(event)
			}
		}

		// Send a failure to all waiting clients.
		if len(t.queue) > 0 {
			for _, waitingClient := range t.queue {
				waitingClient <- _default.NewDefaultPlcConnectionConnectResult(nil, connectionResult.GetErr())
			}
		}
	}
}

func (t *ConnectionContainer) addListener(listener ConnectionListener) {
	// Get the lock.
	t.lock.Lock()
	defer t.lock.Unlock()
	// Add the listener to the queue
	t.listeners = append(t.listeners, listener)
}

func (t *ConnectionContainer) lease() <-chan plc4go.PlcConnectionConnectResult {
	t.lock.Lock()
	defer t.lock.Unlock()

	ch := make(chan plc4go.PlcConnectionConnectResult)
	// Check if the connection is available.
	if t.state == StateIdle {
		t.leaseCounter++
		connection := NewLeasedPlcConnection(t, t.leaseCounter, t.connection)
		// In this case we don't need to check for blocks
		// as the getConnection function of the connection pool
		// is definitely eagerly waiting for input.
		go func() {
			ch <- _default.NewDefaultPlcConnectionConnectResult(connection, nil)
		}()
	} else if t.state == StateInUse || t.state == StateInitialized {
		// If the connection is currently busy or not finished initializing,
		// add the new channel to the queue for this connection.
		t.queue = append(t.queue, ch)
	}
	return ch
}

func (t *ConnectionContainer) returnConnection(state PooledPlcConnectionState) error {
	// Intentionally not locking anything, as there are two cases, where the connection is returned:
	// 1) The connection failed to get established (No connection has a lock anyway)
	// 2) The connection is returned, then the one returning it already has a lock on it.
	// If the connection is marked as "invalid", destroy it and remove it from the pool.
	if state == StateInvalid {
		// TODO: Perhaps do a maximum number of retries and then call failConnection()
		t.connect()
	}

	// Check how many others are waiting for this connection.
	if len(t.queue) > 0 {
		// There are waiting clients, give the connection to the next client in the line.
		next := t.queue[0]
		t.queue = t.queue[1:]
		t.leaseCounter++
		connection := NewLeasedPlcConnection(t, t.leaseCounter, t.connection)
		// Send asynchronously as the receiver might have given up waiting,
		// and we don't want anything to block here. 1ms should be enough for
		// the calling process to reach the blocking read.
		go func() {
			// In this case we don't need to check for blocks
			// as the getConnection function of the connection pool
			// is definitely eagerly waiting for input.
			next <- _default.NewDefaultPlcConnectionConnectResult(connection, nil)
		}()
	} else {
		// Otherwise, just mark the connection as idle.
		t.state = StateIdle
	}
	return nil
}
