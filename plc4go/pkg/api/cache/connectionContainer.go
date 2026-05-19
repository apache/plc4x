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
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
)

type connectionContainer struct {
	lock             *sync.RWMutex
	connectionString string
	driverManager    plc4go.PlcDriverManager
	tracerEnabled    bool
	connection       tracedPlcConnection
	leaseCounter     uint32
	closed           bool
	// The current state of this connection.
	state cachedPlcConnectionState
	// Queue of waiting clients.
	queue []connectionRequest
	// Listeners for connection events.
	listeners []connectionListener

	log zerolog.Logger
}

type connectionRequest struct {
	ctx      context.Context
	connChan chan *plcConnectionLease
	errChan  chan error
}

func newConnectionContainer(log zerolog.Logger, driverManager plc4go.PlcDriverManager, connectionString string) *connectionContainer {
	return &connectionContainer{
		driverManager:    driverManager,
		connectionString: connectionString,
		lock:             &sync.RWMutex{},
		leaseCounter:     0,
		closed:           false,
		state:            StateInitialized,

		log: log,
	}
}

func (c *connectionContainer) connect(ctx context.Context) {
	c.log.Debug().Str("connectionString", c.connectionString).Msg("Connecting new cached connection ...")
	// Initialize the new connection.
	connection, err := c.driverManager.GetConnection(ctx, c.connectionString)

	// Allow us to finish this function and return the lock quickly
	// Wait for the connection to be established.

	// Get the lock.
	c.lock.Lock()
	defer c.lock.Unlock()

	// If the connection was successful, pass the active connection into the container.
	// If something went wrong, we have to remove the connection from the cache and return the error.
	if err != nil {
		c.log.Debug().Str("connectionString", c.connectionString).
			Err(err).
			Msg("Error connecting new cached connection.")
		// Tell the connection cache that the connection is no longer available.
		if c.listeners != nil {
			event := &connectionErrorEvent{
				conn: c,
				err:  err,
			}
			for _, listener := range c.listeners {
				listener.onConnectionEvent(event)
			}
		}

		// Send a failure to all waiting clients.
		if len(c.queue) > 0 {
			c.log.Trace().Msg("notifies waiting clients of error")
			for _, waitingClient := range c.queue {
				select {
				case waitingClient.errChan <- err:
					c.log.Trace().Msg("sent error to waiting client")
				case <-waitingClient.ctx.Done():
					c.log.Trace().Msg("waiting client timed out")
				case <-ctx.Done():
					c.log.Trace().Msg("context timed out")
				}
			}
			c.queue = nil
		} else {
			c.log.Trace().Msg("no waiting clients")
		}
		return
	}

	c.log.Debug().Str("connectionString", c.connectionString).Msg("Successfully connected new cached connection.")
	// Inject the real connection into the container.
	if connection, ok := connection.(tracedPlcConnection); !ok {
		panic("Return connection doesn't implement the cache.tracedPlcConnection interface")
	} else {
		c.connection = connection
	}
	c.tracerEnabled = c.connection.IsTraceEnabled()
	// Mark the connection as idle for now.
	c.state = StateIdle
	// If there is a request in the queue, hand out the connection to that.
	if waitingClientsLen := len(c.queue); waitingClientsLen > 0 {
		c.log.Trace().Int("waitingClientsLen", waitingClientsLen).Msg("notifies waiting clients of connection")
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
		queueHead.connChan <- connection
	} else {
		c.log.Trace().Msg("no waiting clients")
	}
}

func (c *connectionContainer) addListener(listener connectionListener) {
	// Get the lock.
	c.lock.Lock()
	defer c.lock.Unlock()
	// Add the listener to the queue
	c.listeners = append(c.listeners, listener)
}

func (c *connectionContainer) lease(ctx context.Context) (chan *plcConnectionLease, chan error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	connectionChan := make(chan *plcConnectionLease, 1)
	errorChan := make(chan error, 1)
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
		connectionChan <- connection
	case StateInUse, StateInitialized:
		// If the connection is currently busy or not finished initializing,
		// add the new channel to the queue for this connection.
		c.queue = append(c.queue, connectionRequest{ctx: ctx, connChan: connectionChan, errChan: errorChan})
		c.log.Debug().Str("connectionString", c.connectionString).
			Int("waiting-queue-size", len(c.queue)).
			Msg("Added lease-request to queue.")
	case StateInvalid:
		c.log.Debug().Str("connectionString", c.connectionString).Msg("No lease because invalid")
	}
	return connectionChan, errorChan
}

func (c *connectionContainer) returnConnection(ctx context.Context, newState cachedPlcConnectionState) error {
	// Intentionally not locking anything, as there are two cases, where the connection is returned:
	// 1) The connection failed to get established (No connection has a lock anyway)
	// 2) The connection is returned, then the one returning it already has a lock on it.
	// If the connection is marked as "invalid", destroy it and remove it from the cache.
	switch newState {
	case StateInitialized, StateInvalid:
		// TODO: Perhaps do a maximum number of retries and then call failConnection()
		c.log.Debug().
			Str("connectionString", c.connectionString).
			Stringer("newState", newState).
			Msg("Client returned a connection, reconnecting.")
		c.connect(ctx)
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
	if waitingClientsLen := len(c.queue); waitingClientsLen > 0 {
		c.log.Trace().Int("waitingClientsLen", waitingClientsLen).Msg("notifies waiting clients of connection return")
		// There are waiting clients, give the connection to the next client in the line.
		next := c.queue[0]
		c.queue = c.queue[1:]
		c.leaseCounter++
		connection := newPlcConnectionLease(c, c.leaseCounter, c.connection)

		// In this case we don'c need to check for blocks
		// as the getConnection function of the connection cache
		// is definitely eagerly waiting for input.
		next.connChan <- connection
		c.log.Debug().Str("connectionString", c.connectionString).
			Int("waiting-queue-size", len(c.queue)).
			Msg("Returned connection to the next client waiting.")
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
