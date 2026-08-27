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
	"time"

	"github.com/rs/zerolog"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiOptions "github.com/apache/plc4x/plc4go/spi/options"
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
	// Maximum duration a connection may sit idle in the cache before it is
	// discarded and re-established on the next lease (0 = keep forever).
	// Remotes routinely reap idle connections without a FIN/RST reaching us
	// (half-open TCP), which no liveness flag can detect - the death only
	// surfaces on the first write. Refusing to hand out connections that
	// exceeded this age avoids that failure mode and frees connection slots
	// on connection-limited remotes between bursts.
	maxIdleTime time.Duration
	// idleSince records when the connection last became idle.
	idleSince time.Time

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
	c.log.Debug().Str("connectionString", spiOptions.RedactConnectionString(c.connectionString)).Msg("Connecting new cached connection ...")
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
		c.log.Debug().Str("connectionString", spiOptions.RedactConnectionString(c.connectionString)).
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
		// Mark the container as invalid so the next lease request knows it has to
		// re-attempt the connection. Previously the state was left untouched
		// (StateInitialized on the initial connect), which parked the container in a
		// state where lease() queued requests forever and no reconnect was ever
		// attempted - permanently breaking the connection string until restart.
		c.state = StateInvalid
		return
	}

	c.log.Debug().Str("connectionString", spiOptions.RedactConnectionString(c.connectionString)).Msg("Successfully connected new cached connection.")
	// Inject the real connection into the container.
	if connection, ok := connection.(tracedPlcConnection); !ok {
		panic("Return connection doesn't implement the cache.tracedPlcConnection interface")
	} else {
		c.connection = connection
	}
	c.tracerEnabled = c.connection.IsTraceEnabled()
	// Mark the connection as idle for now.
	c.state = StateIdle
	c.idleSince = time.Now()
	// If there is a request in the queue, hand out the connection to that.
	if queueHead := c.nextWaiter(); queueHead != nil {
		c.log.Trace().Int("waitingClientsLen", len(c.queue)).Msg("notifies waiting clients of connection")
		// Mark the connection as being used.
		c.state = StateInUse
		c.leaseCounter++
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

// nextWaiter pops the next queued lease request whose context is still live.
// Requests whose caller has already given up waiting are failed explicitly (their
// error channel is buffered) instead of being handed a lease: sending a lease to
// an abandoned request would park the connection in StateInUse forever, as nobody
// is left to return it. Must be called with c.lock held.
func (c *connectionContainer) nextWaiter() *connectionRequest {
	for len(c.queue) > 0 {
		head := c.queue[0]
		c.queue = c.queue[1:]
		if head.ctx.Err() != nil {
			c.log.Debug().Str("connectionString", spiOptions.RedactConnectionString(c.connectionString)).
				Msg("Skipping lease request with cancelled context")
			select {
			case head.errChan <- head.ctx.Err():
			default:
			}
			continue
		}
		return &head
	}
	return nil
}

// startReconnect discards any stale connection and re-runs connect asynchronously.
// State is set to StateInitialized so concurrent lease requests queue up instead of
// spawning additional connect attempts. Must be called with c.lock held.
func (c *connectionContainer) startReconnect(ctx context.Context) {
	stale := c.connection
	c.connection = nil
	c.state = StateInitialized
	go func() {
		if stale != nil {
			// Close the stale connection so its message-codec workers don't leak.
			if err := stale.Close(); err != nil {
				c.log.Debug().Err(err).
					Str("connectionString", spiOptions.RedactConnectionString(c.connectionString)).
					Msg("Error closing stale connection before reconnect")
			}
		}
		c.connect(ctx)
	}()
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
		// Verify the cached connection is still alive before handing it out. The
		// remote may have dropped it while it sat idle in the cache (idle timeouts
		// and connection-limited gateways do this routinely); without this check
		// the client receives a dead connection and fails on first use.
		// A connection that exceeded the configured max idle time is treated the
		// same way: remotes that silently reap idle connections leave a half-open
		// socket that IsConnected() cannot detect, so past that age the connection
		// is not to be trusted and gets replaced proactively.
		if c.connection == nil || !c.connection.IsConnected() || c.idleExpired() {
			if c.closed {
				// The cache is shutting down - don't reconnect, just report.
				errorChan <- errors.New("connection container is closed")
				break
			}
			if c.idleExpired() {
				c.log.Debug().Str("connectionString", spiOptions.RedactConnectionString(c.connectionString)).
					Dur("maxIdleTime", c.maxIdleTime).
					Time("idleSince", c.idleSince).
					Msg("Cached idle connection exceeded max idle time - reconnecting.")
			} else {
				c.log.Debug().Str("connectionString", spiOptions.RedactConnectionString(c.connectionString)).
					Msg("Cached idle connection is no longer alive - reconnecting.")
			}
			c.queue = append(c.queue, connectionRequest{ctx: ctx, connChan: connectionChan, errChan: errorChan})
			c.startReconnect(ctx)
			break
		}
		c.leaseCounter++
		connection := newPlcConnectionLease(c, c.leaseCounter, c.connection)
		c.state = StateInUse
		// In this case we don'c need to check for blocks
		// as the getConnection function of the connection cache
		// is definitely eagerly waiting for input.
		c.log.Debug().Str("connectionString", spiOptions.RedactConnectionString(c.connectionString)).
			Msg("Got lease instantly as connection was idle.")
		connectionChan <- connection
	case StateInUse, StateInitialized:
		// If the connection is currently busy or not finished initializing,
		// add the new channel to the queue for this connection.
		c.queue = append(c.queue, connectionRequest{ctx: ctx, connChan: connectionChan, errChan: errorChan})
		c.log.Debug().Str("connectionString", spiOptions.RedactConnectionString(c.connectionString)).
			Int("waiting-queue-size", len(c.queue)).
			Msg("Added lease-request to queue.")
	case StateInvalid:
		// A previous connect or reconnect failed. Previously the request was simply
		// ignored here (nothing was ever sent on either channel), so every caller
		// blocked until timeout and the container stayed broken forever. Instead,
		// queue the request and attempt to re-establish the connection.
		if c.closed {
			errorChan <- errors.New("connection container is closed")
			break
		}
		c.log.Debug().Str("connectionString", spiOptions.RedactConnectionString(c.connectionString)).
			Msg("Connection is invalid - attempting reconnect.")
		c.queue = append(c.queue, connectionRequest{ctx: ctx, connChan: connectionChan, errChan: errorChan})
		c.startReconnect(ctx)
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
			Str("connectionString", spiOptions.RedactConnectionString(c.connectionString)).
			Stringer("newState", newState).
			Msg("Client returned a connection, reconnecting.")
		// Close the stale connection before reconnecting. c.connect() overwrites
		// c.connection with a freshly-established one, so without this the previous
		// connection's message-codec workers (ReceiveWork/ExpireWork) keep running,
		// leaking a pair of goroutines on every invalidate->reconnect cycle. Against
		// an endpoint that accepts TCP but is unresponsive at the protocol layer this
		// recurs every poll and grows without bound.
		c.lock.Lock()
		stale := c.connection
		c.connection = nil
		c.lock.Unlock()
		if stale != nil {
			if err := stale.Close(); err != nil {
				c.log.Debug().Err(err).
					Str("connectionString", spiOptions.RedactConnectionString(c.connectionString)).
					Msg("Error closing stale connection before reconnect")
			}
		}
		c.connect(ctx)

		c.lock.Lock()
		defer c.lock.Unlock()
		if c.connection == nil {
			c.state = StateInvalid
			return errors.New("Can't return a broken connection")
		}
		// connect() already dealt with the queue: it either handed the fresh
		// connection to the next waiter (StateInUse) or marked it idle. Handing out
		// another lease here would lease the same connection twice.
		return nil
	default:
		c.log.Debug().Str("connectionString", spiOptions.RedactConnectionString(c.connectionString)).Msg("Client returned valid connection.")
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.connection == nil {
		c.state = StateInvalid
		return errors.New("Can't return a broken connection")
	}

	// Check how many others are waiting for this connection.
	if next := c.nextWaiter(); next != nil {
		c.log.Trace().Int("waitingClientsLen", len(c.queue)).Msg("notifies waiting clients of connection return")
		// There are waiting clients, give the connection to the next client in the line.
		c.state = StateInUse
		c.leaseCounter++
		connection := newPlcConnectionLease(c, c.leaseCounter, c.connection)

		// In this case we don'c need to check for blocks
		// as the getConnection function of the connection cache
		// is definitely eagerly waiting for input.
		next.connChan <- connection
		c.log.Debug().Str("connectionString", spiOptions.RedactConnectionString(c.connectionString)).
			Int("waiting-queue-size", len(c.queue)).
			Msg("Returned connection to the next client waiting.")
	} else {
		// Otherwise, just mark the connection as idle.
		c.log.Debug().Str("connectionString", spiOptions.RedactConnectionString(c.connectionString)).
			Msg("Connection set to 'idle'.")
		c.state = StateIdle
		c.idleSince = time.Now()
	}
	return nil
}

// subscriptionAware is an optional capability of cached connections: drivers
// whose connections carry stateful subscriptions (e.g. BACnet COV, with its
// server-side registrations and client-side refresh timers) implement it so
// the cache can tell "parked but working" apart from "abandoned".
type subscriptionAware interface {
	// ActiveSubscriptionCount reports the number of currently active
	// subscription handles on the connection.
	ActiveSubscriptionCount() int
}

// idleExpired reports whether the connection outstayed the configured max idle
// time. Must be called with c.lock held.
//
// A connection with active subscription handles never expires: its
// subscription state lives on the connection, so reaping it would silently
// destroy server-side registrations and refresh timers, cutting off passive
// updates until the client happens to re-subscribe. Such a connection is not
// "idle" in any meaningful sense even when no lease touched it for a while.
func (c *connectionContainer) idleExpired() bool {
	if c.maxIdleTime <= 0 || time.Since(c.idleSince) <= c.maxIdleTime {
		return false
	}
	if sa, ok := c.connection.(subscriptionAware); ok && sa.ActiveSubscriptionCount() > 0 {
		return false
	}
	return true
}

func (c *connectionContainer) String() string {
	return fmt.Sprintf("connectionContainer{%s:%s, leaseCounter: %d, closed: %t, state: %s}", c.connectionString, c.connection, c.leaseCounter, c.closed, c.state)
}
