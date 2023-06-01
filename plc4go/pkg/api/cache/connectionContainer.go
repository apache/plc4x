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
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/viney-shih/go-lock"
)

type connectionContainer struct {
	lock             lock.RWMutex
	connectionString string
	driverManager    plc4go.PlcDriverManager
	tracerEnabled    bool
	connection       tracedPlcConnection
	leaseCounter     uint32
	closed           bool
	// The current state of this connection.
	state cachedPlcConnectionState
	// Queue of waiting clients.
	queue []chan plc4go.PlcConnectionConnectResult
	// Listeners for connection events.
	listeners []connectionListener

	log zerolog.Logger
}

func newConnectionContainer(log zerolog.Logger, driverManager plc4go.PlcDriverManager, connectionString string) *connectionContainer {
	return &connectionContainer{
		driverManager:    driverManager,
		connectionString: connectionString,
		lock:             lock.NewCASMutex(),
		leaseCounter:     0,
		closed:           false,
		state:            StateInitialized,
		queue:            []chan plc4go.PlcConnectionConnectResult{},

		log: log,
	}
}

func (c *connectionContainer) connect() {
	c.log.Debug().Str("connectionString", c.connectionString).Msg("Connecting new cached connection ...")
	// Initialize the new connection.
	connectionResultChan := c.driverManager.GetConnection(c.connectionString)

	// Allow us to finish this function and return the lock quickly
	// Wait for the connection to be established.
	// TODO: Add some timeout handling.
	connectionResult := <-connectionResultChan

	// Get the lock.
	c.lock.Lock()
	defer c.lock.Unlock()

	// If the connection was successful, pass the active connection into the container.
	// If something went wrong, we have to remove the connection from the cache and return the error.
	if err := connectionResult.GetErr(); err != nil {
		c.log.Debug().Str("connectionString", c.connectionString).
			Err(err).
			Msg("Error connecting new cached connection.")
		// Tell the connection cache that the connection is no longer available.
		if c.listeners != nil {
			event := connectionErrorEvent{
				conn: *c,
				err:  err,
			}
			for _, listener := range c.listeners {
				listener.onConnectionEvent(event)
			}
		}

		// Send a failure to all waiting clients.
		if len(c.queue) > 0 {
			for _, waitingClient := range c.queue {
				waitingClient <- _default.NewDefaultPlcConnectionConnectResult(nil, err)
			}
			c.queue = nil
		}
		return
	}

	c.log.Debug().Str("connectionString", c.connectionString).Msg("Successfully connected new cached connection.")
	// Inject the real connection into the container.
	if connection, ok := connectionResult.GetConnection().(tracedPlcConnection); !ok {
		panic("Return connection doesn'c implement the cache.tracedPlcConnection interface")
	} else {
		c.connection = connection
	}
	c.tracerEnabled = c.connection.IsTraceEnabled()
	// Mark the connection as idle for now.
	c.state = StateIdle
	// If there is a request in the queue, hand out the connection to that.
	if len(c.queue) > 0 {
		// Get the first in the queue.
		queueHead := c.queue[0]
		c.queue = c.queue[1:]
		// Mark the connection as being used.
		c.state = StateInUse
		// Return the lease to the caller.
		connection := newPlcConnectionLease(c, c.leaseCounter, c.connection)
		// In this case we don'c need to check for blocks
		// as the getConnection function of the connection cache
		// is definitely eagerly waiting for input.
		queueHead <- _default.NewDefaultPlcConnectionConnectResult(connection, nil)
	}
}

func (c *connectionContainer) addListener(listener connectionListener) {
	// Get the lock.
	c.lock.Lock()
	defer c.lock.Unlock()
	// Add the listener to the queue
	c.listeners = append(c.listeners, listener)
}

func (c *connectionContainer) lease() <-chan plc4go.PlcConnectionConnectResult {
	c.lock.Lock()
	defer c.lock.Unlock()

	ch := make(chan plc4go.PlcConnectionConnectResult)
	// Check if the connection is available.
	switch c.state {
	case StateIdle:
		c.leaseCounter++
		connection := newPlcConnectionLease(c, c.leaseCounter, c.connection)
		c.state = StateInUse
		// In this case we don'c need to check for blocks
		// as the getConnection function of the connection cache
		// is definitely eagerly waiting for input.
		c.log.Debug().Str("connectionString", c.connectionString).
			Msg("Got lease instantly as connection was idle.")
		go func() {
			ch <- _default.NewDefaultPlcConnectionConnectResult(connection, nil)
		}()
	case StateInUse, StateInitialized:
		// If the connection is currently busy or not finished initializing,
		// add the new channel to the queue for this connection.
		c.queue = append(c.queue, ch)
		c.log.Debug().Str("connectionString", c.connectionString).
			Int("waiting-queue-size", len(c.queue)).
			Msg("Added lease-request to queue.")
	case StateInvalid:
		c.log.Debug().Str("connectionString", c.connectionString).Msg("No lease because invalid")
	}
	return ch
}

func (c *connectionContainer) returnConnection(newState cachedPlcConnectionState) error {
	// Intentionally not locking anything, as there are two cases, where the connection is returned:
	// 1) The connection failed to get established (No connection has a lock anyway)
	// 2) The connection is returned, then the one returning it already has a lock on it.
	// If the connection is marked as "invalid", destroy it and remove it from the cache.
	switch newState {
	case StateInitialized, StateInvalid:
		// TODO: Perhaps do a maximum number of retries and then call failConnection()
		c.log.Debug().Str("connectionString", c.connectionString).
			Msgf("Client returned a %s connection, reconnecting.", newState)
		c.connect()
	default:
		c.log.Debug().Str("connectionString", c.connectionString).Msg("Client returned valid connection.")
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.connection == nil {
		c.state = StateInvalid
		return errors.New("Can'c return a broken connection")
	}

	// Check how many others are waiting for this connection.
	if len(c.queue) > 0 {
		// There are waiting clients, give the connection to the next client in the line.
		next := c.queue[0]
		c.queue = c.queue[1:]
		c.leaseCounter++
		connection := newPlcConnectionLease(c, c.leaseCounter, c.connection)
		// Send asynchronously as the receiver might have given up waiting,
		// and we don'c want anything to block here. 1ms should be enough for
		// the calling process to reach the blocking read.
		go func() {
			// In this case we don'c need to check for blocks
			// as the getConnection function of the connection cache
			// is definitely eagerly waiting for input.
			next <- _default.NewDefaultPlcConnectionConnectResult(connection, nil)
			c.log.Debug().Str("connectionString", c.connectionString).
				Int("waiting-queue-size", len(c.queue)).
				Msg("Returned connection to the next client waiting.")
		}()
	} else {
		// Otherwise, just mark the connection as idle.
		c.log.Debug().Str("connectionString", c.connectionString).
			Msg("Connection set to 'idle'.")
		c.state = StateIdle
	}
	return nil
}

func (c *connectionContainer) String() string {
	return fmt.Sprintf("connectionContainer{%s:%s, leaseCounter: %d, closed: %t, state: %s}", c.connectionString, c.connection, c.leaseCounter, c.closed, c.state)
}
