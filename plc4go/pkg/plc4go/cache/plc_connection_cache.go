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

package cache

import (
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/internal/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/pkg/plc4go"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/pkg/errors"
	"github.com/viney-shih/go-lock"
	"time"
)

type PlcConnectionCache interface {
	GetConnection(connectionString string) <-chan plc4go.PlcConnectionConnectResult
	Close() <-chan PlcConnectionCacheCloseResult
}

type plcConnectionCache struct {
	driverManager plc4go.PlcDriverManager

	// Maximum duration a connection can be used per lease.
	// If the connection is used for a longer time, it is forcefully removed from the client.
	maxLeaseTime time.Duration
	maxWaitTime  time.Duration

	cacheLock   lock.RWMutex
	connections map[string]*connectionContainer
	tracer      *spi.Tracer
}

func NewPlcConnectionCache(driverManager plc4go.PlcDriverManager) PlcConnectionCache {
	return NewPlcConnectionCacheWithMaxLeaseTime(driverManager, time.Second*5)
}

func NewPlcConnectionCacheWithMaxLeaseTime(driverManager plc4go.PlcDriverManager, maxLeaseTime time.Duration) PlcConnectionCache {
	return &plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  maxLeaseTime,
		maxWaitTime:   maxLeaseTime * 5,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
}

func (t *plcConnectionCache) EnableTracer() {
	t.tracer = spi.NewTracer("cache")
}

func (t *plcConnectionCache) GetTracer() *spi.Tracer {
	return t.tracer
}

func (t *plcConnectionCache) GetConnection(connectionString string) <-chan plc4go.PlcConnectionConnectResult {
	ch := make(chan plc4go.PlcConnectionConnectResult)

	go func() {
		t.cacheLock.Lock()

		// If a connection for this connection string didn't exist yet, create a new container
		// and make that container connect.
		if _, ok := t.connections[connectionString]; !ok {
			if t.tracer != nil {
				t.tracer.AddTrace("get-connection", "create new cached connection")
			}
			// Create a new connection container.
			cc := &connectionContainer{
				driverManager:    t.driverManager,
				connectionString: connectionString,
				lock:             lock.NewCASMutex(),
				leaseCounter:     0,
				closed:           false,
				state:            StateInitialized,
				queue:            []chan plc4go.PlcConnectionConnectResult{},
			}
			// Register for connection events (Like connection closed or error).
			cc.addListener(t)
			// Store the new connection container in the cache of connections.
			t.connections[connectionString] = cc
			// Initialize the connection itself.
			go func(cc2 *connectionContainer) {
				cc2.connect()
			}(cc)
		}

		// Get the ConnectionContainer for this connection string.
		connection := t.connections[connectionString]

		// Release the lock again.
		t.cacheLock.Unlock()

		// Try to get a lease on this connection.
		var txId string
		if t.tracer != nil {
			txId = t.tracer.AddTransactionalStartTrace("get-connection", "lease")
		}
		leaseChan := connection.lease()
		select {
		// Wait till we get a lease.
		case connectionResponse := <-leaseChan:
			select {
			case ch <- connectionResponse:
				if t.tracer != nil {
					t.tracer.AddTransactionalTrace(txId, "get-connection", "success")
				}
			case <-time.After(10 * time.Millisecond):
				// Log a message, that the client has given up
				if t.tracer != nil {
					t.tracer.AddTransactionalTrace(txId, "get-connection", "client given up")
				}
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

func (t *plcConnectionCache) Close() <-chan PlcConnectionCacheCloseResult {
	ch := make(chan PlcConnectionCacheCloseResult)

	go func() {
		t.cacheLock.Lock()
		defer t.cacheLock.Unlock()

		for _, cc := range t.connections {
			// Mark the connection as being closed to not try to re-establish it.
			cc.closed = true
			// Try to get a lease as this way we kow we're not closing the connection
			// while some go func is still using it.
			go func(container *connectionContainer) {
				leaseResults := container.lease()
				select {
				// We're just getting the lease as this way we can be sure nobody else is using it.
				// We also really don't care if it worked, or not ... it's just an attempt of being
				// nice.
				case _ = <-leaseResults:
					// Give back the connection.
					container.connection.Close()
				// If we're timing out brutally kill the connection.
				case <-time.After(t.maxWaitTime):
					// Forcefully close this connection.
					container.connection.Close()
				}

				select {
				case ch <- newDefaultPlcConnectionCacheCloseResult(t, nil):
				case <-time.After(time.Millisecond * 10):
				}
			}(cc)
		}
	}()

	return ch
}

// onConnectionEvent: Callback called by the connection container to signal connection events
// that have an impact on the cache itself (Like connections being permanently closed).
func (t *plcConnectionCache) onConnectionEvent(event connectionEvent) {
	connectionContainer := event.getConnectionContainer()
	if errorEvent, ok := event.(connectionErrorEvent); ok {
		if t.tracer != nil {
			t.tracer.AddTrace("destroy-connection", errorEvent.getError().Error())
		}
		delete(t.connections, connectionContainer.connectionString)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// connectionContainer
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type connectionContainer struct {
	lock             lock.RWMutex
	connectionString string
	driverManager    plc4go.PlcDriverManager
	tracerEnabled    bool
	connection       spi.PlcConnection
	leaseCounter     uint32
	closed           bool
	// The current state of this connection.
	state cachedPlcConnectionState
	// Queue of waiting clients.
	queue []chan plc4go.PlcConnectionConnectResult
	// Listeners for connection events.
	listeners []connectionListener
}

func (t *connectionContainer) connect() {
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
	// If something went wrong, we have to remove the connection from the cache and return the error.
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
			connection := newPlcConnectionLease(t, t.leaseCounter, t.connection)
			// In this case we don't need to check for blocks
			// as the getConnection function of the connection cache
			// is definitely eagerly waiting for input.
			queueHead <- _default.NewDefaultPlcConnectionConnectResult(connection, nil)
		}
	} else {
		// Tell the connection cache that the connection is no longer available.
		if t.listeners != nil {
			event := connectionErrorEvent{
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

func (t *connectionContainer) addListener(listener connectionListener) {
	// Get the lock.
	t.lock.Lock()
	defer t.lock.Unlock()
	// Add the listener to the queue
	t.listeners = append(t.listeners, listener)
}

func (t *connectionContainer) lease() <-chan plc4go.PlcConnectionConnectResult {
	t.lock.Lock()
	defer t.lock.Unlock()

	ch := make(chan plc4go.PlcConnectionConnectResult)
	// Check if the connection is available.
	if t.state == StateIdle {
		t.leaseCounter++
		connection := newPlcConnectionLease(t, t.leaseCounter, t.connection)
		// In this case we don't need to check for blocks
		// as the getConnection function of the connection cache
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

func (t *connectionContainer) returnConnection(state cachedPlcConnectionState) error {
	// Intentionally not locking anything, as there are two cases, where the connection is returned:
	// 1) The connection failed to get established (No connection has a lock anyway)
	// 2) The connection is returned, then the one returning it already has a lock on it.
	// If the connection is marked as "invalid", destroy it and remove it from the cache.
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
		connection := newPlcConnectionLease(t, t.leaseCounter, t.connection)
		// Send asynchronously as the receiver might have given up waiting,
		// and we don't want anything to block here. 1ms should be enough for
		// the calling process to reach the blocking read.
		go func() {
			// In this case we don't need to check for blocks
			// as the getConnection function of the connection cache
			// is definitely eagerly waiting for input.
			next <- _default.NewDefaultPlcConnectionConnectResult(connection, nil)
		}()
	} else {
		// Otherwise, just mark the connection as idle.
		t.state = StateIdle
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// plcConnectionLease
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type plcConnectionLease struct {
	// Reference back to the container, so we can give the connection back.
	connectionContainer *connectionContainer
	// Counter for the number of times this connection has been used before.
	leaseId uint32
	// The actual connection being cached.
	connection spi.PlcConnection
}

func newPlcConnectionLease(connectionContainer *connectionContainer, leaseId uint32, connection spi.PlcConnection) *plcConnectionLease {
	p := &plcConnectionLease{
		connectionContainer: connectionContainer,
		leaseId:             leaseId,
		connection:          connection,
	}
	if connection.IsTraceEnabled() {
		connection.GetTracer().SetConnectionId(p.GetConnectionId())
	}
	return p
}

func (t *plcConnectionLease) IsTraceEnabled() bool {
	if t.connection == nil {
		panic("Called 'IsTraceEnabled' on a closed cached connection")
	}
	return t.connection.IsTraceEnabled()
}

func (t *plcConnectionLease) GetTracer() *spi.Tracer {
	if t.connection == nil {
		panic("Called 'GetTracer' on a closed cached connection")
	}
	return t.connection.GetTracer()
}

func (t *plcConnectionLease) GetConnectionId() string {
	if t.connection == nil {
		panic("Called 'GetConnectionId' on a closed cached connection")
	}
	return fmt.Sprintf("%s-%d", t.connection.GetConnectionId(), t.leaseId)
}

func (t *plcConnectionLease) Connect() <-chan plc4go.PlcConnectionConnectResult {
	panic("Called 'Connect' on a cached connection")
}

func (t *plcConnectionLease) BlockingClose() {
	if t.connection == nil {
		panic("Called 'BlockingClose' on a closed cached connection")
	}
	// Call close and wait for the operation to finish.
	<-t.Close()
}

func (t *plcConnectionLease) Close() <-chan plc4go.PlcConnectionCloseResult {
	if t.connection == nil {
		panic("Called 'Close' on a closed cached connection")
	}

	result := make(chan plc4go.PlcConnectionCloseResult)

	go func() {
		// Check if the connection is still alive, if it is, put it back into the cache
		pingResults := t.Ping()
		pingTimeout := time.NewTimer(time.Second * 5)
		newState := StateIdle
		select {
		case pingResult := <-pingResults:
			{
				if pingResult.GetErr() != nil {
					newState = StateInvalid
				}
			}
		case <-pingTimeout.C:
			{
				// Add some trace information
				if t.connection.IsTraceEnabled() {
					t.connection.GetTracer().AddTrace("ping", "timeout")
				}
				// Mark the connection as broken ...
				newState = StateInvalid
			}
		}

		// Extract the trace entries from the connection.
		var traces []spi.TraceEntry
		if t.IsTraceEnabled() {
			tracer := t.GetTracer()
			// Save all traces.
			traces = tracer.GetTraces()
			// Clear the log.
			tracer.ResetTraces()
			// Reset the connection id back to the one without the lease-id.
			tracer.SetConnectionId(t.connection.GetConnectionId())
		}

		// Return the connection to the connection container and don't actually close it.
		err := t.connectionContainer.returnConnection(newState)

		// Finish closing the connection.
		result <- _default.NewDefaultPlcConnectionCloseResultWithTraces(t, err, traces)

		// Detach the connection from this lease, so it can no longer be used by the client.
		t.connection = nil
	}()

	return result
}

func (t *plcConnectionLease) IsConnected() bool {
	if t.connection == nil {
		return false
	}
	return t.connection.IsConnected()
}

func (t *plcConnectionLease) Ping() <-chan plc4go.PlcConnectionPingResult {
	if t.connection == nil {
		panic("Called 'Ping' on a closed cached connection")
	}
	return t.connection.Ping()
}

func (t *plcConnectionLease) GetMetadata() model.PlcConnectionMetadata {
	if t.connection == nil {
		panic("Called 'GetMetadata' on a closed cached connection")
	}
	return t.connection.GetMetadata()
}

func (t *plcConnectionLease) ReadRequestBuilder() model.PlcReadRequestBuilder {
	if t.connection == nil {
		panic("Called 'ReadRequestBuilder' on a closed cached connection")
	}
	return t.connection.ReadRequestBuilder()
}

func (t *plcConnectionLease) WriteRequestBuilder() model.PlcWriteRequestBuilder {
	if t.connection == nil {
		panic("Called 'WriteRequestBuilder' on a closed cached connection")
	}
	return t.connection.WriteRequestBuilder()
}

func (t *plcConnectionLease) SubscriptionRequestBuilder() model.PlcSubscriptionRequestBuilder {
	if t.connection == nil {
		panic("Called 'SubscriptionRequestBuilder' on a closed cached connection")
	}
	return t.connection.SubscriptionRequestBuilder()
}

func (t *plcConnectionLease) UnsubscriptionRequestBuilder() model.PlcUnsubscriptionRequestBuilder {
	if t.connection == nil {
		panic("Called 'UnsubscriptionRequestBuilder' on a closed cached connection")
	}
	return t.connection.UnsubscriptionRequestBuilder()
}

func (t *plcConnectionLease) BrowseRequestBuilder() model.PlcBrowseRequestBuilder {
	if t.connection == nil {
		panic("Called 'BrowseRequestBuilder' on a closed cached connection")
	}
	return t.connection.BrowseRequestBuilder()
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// cachedPlcConnectionState
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type cachedPlcConnectionState int32

const (
	StateInitialized cachedPlcConnectionState = iota
	StateIdle
	StateInUse
	StateInvalid
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Events
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type connectionListener interface {
	onConnectionEvent(event connectionEvent)
}

type connectionEvent interface {
	getConnectionContainer() connectionContainer
}

type connectionErrorEvent struct {
	conn connectionContainer
	err  error
}

func (c connectionErrorEvent) getConnectionContainer() connectionContainer {
	return c.conn
}

func (c connectionErrorEvent) getError() error {
	return c.err
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// PlcConnectionCacheCloseResult / plcConnectionCacheCloseResult
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type PlcConnectionCacheCloseResult interface {
	GetConnectionCache() PlcConnectionCache
	GetErr() error
}

type plcConnectionCacheCloseResult struct {
	connectionCache PlcConnectionCache
	err             error
}

func newDefaultPlcConnectionCacheCloseResult(connectionCache PlcConnectionCache, err error) PlcConnectionCacheCloseResult {
	return &plcConnectionCacheCloseResult{
		connectionCache: connectionCache,
		err:             err,
	}
}

func (p plcConnectionCacheCloseResult) GetConnectionCache() PlcConnectionCache {
	return p.connectionCache
}

func (p plcConnectionCacheCloseResult) GetErr() error {
	return p.err
}
