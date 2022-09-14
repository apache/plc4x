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

package cache

import (
	"fmt"
	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/rs/zerolog/log"
	"github.com/viney-shih/go-lock"
)

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
	log.Debug().Str("connectionString", t.connectionString).Msg("Connecting new cached connection ...")
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
	if err := connectionResult.GetErr(); err != nil {
		log.Debug().Str("connectionString", t.connectionString).
			Err(err).
			Msg("Error connecting new cached connection.")
		// Tell the connection cache that the connection is no longer available.
		if t.listeners != nil {
			event := connectionErrorEvent{
				conn: *t,
				err:  err,
			}
			for _, listener := range t.listeners {
				listener.onConnectionEvent(event)
			}
		}

		// Send a failure to all waiting clients.
		if len(t.queue) > 0 {
			for _, waitingClient := range t.queue {
				waitingClient <- _default.NewDefaultPlcConnectionConnectResult(nil, err)
			}
			t.queue = nil
		}
		return
	}

	log.Debug().Str("connectionString", t.connectionString).Msg("Successfully connected new cached connection.")
	// Inject the real connection into the container.
	if connection, ok := connectionResult.GetConnection().(spi.PlcConnection); !ok {
		panic("Return connection doesn't implement the spi.PlcConnection interface")
	} else {
		t.connection = connection
	}
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
	switch t.state {
	case StateIdle:
		t.leaseCounter++
		connection := newPlcConnectionLease(t, t.leaseCounter, t.connection)
		t.state = StateInUse
		// In this case we don't need to check for blocks
		// as the getConnection function of the connection cache
		// is definitely eagerly waiting for input.
		log.Debug().Str("connectionString", t.connectionString).
			Msg("Got lease instantly as connection was idle.")
		go func() {
			ch <- _default.NewDefaultPlcConnectionConnectResult(connection, nil)
		}()
	case StateInUse, StateInitialized:
		// If the connection is currently busy or not finished initializing,
		// add the new channel to the queue for this connection.
		t.queue = append(t.queue, ch)
		log.Debug().Str("connectionString", t.connectionString).
			Int("waiting-queue-size", len(t.queue)).
			Msg("Added lease-request to queue.")
	case StateInvalid:
		log.Debug().Str("connectionString", t.connectionString).Msg("No lease because invalid")
	}
	return ch
}

func (t *connectionContainer) returnConnection(state cachedPlcConnectionState) error {
	// Intentionally not locking anything, as there are two cases, where the connection is returned:
	// 1) The connection failed to get established (No connection has a lock anyway)
	// 2) The connection is returned, then the one returning it already has a lock on it.
	// If the connection is marked as "invalid", destroy it and remove it from the cache.
	switch state {
	case StateInitialized, StateInvalid:
		// TODO: Perhaps do a maximum number of retries and then call failConnection()
		log.Debug().Str("connectionString", t.connectionString).
			Msgf("Client returned a %s connection, reconnecting.", state)
		t.connect()
	default:
		log.Debug().Str("connectionString", t.connectionString).Msg("Client returned valid connection.")
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
			log.Debug().Str("connectionString", t.connectionString).
				Int("waiting-queue-size", len(t.queue)).
				Msg("Returned connection to the next client waiting.")
		}()
	} else {
		// Otherwise, just mark the connection as idle.
		log.Debug().Str("connectionString", t.connectionString).
			Msg("Connection set to 'idle'.")
		t.state = StateIdle
	}
	return nil
}

func (t *connectionContainer) String() string {
	return fmt.Sprintf("connectionContainer{%s:%s, leaseCounter: %d, closed: %t, state: %s}", t.connectionString, t.connection, t.leaseCounter, t.closed, t.state)
}
