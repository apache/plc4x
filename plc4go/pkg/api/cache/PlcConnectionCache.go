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
	"runtime/debug"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/config"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/tracer"
)

type PlcConnectionCache interface {
	GetConnection(ctx context.Context, connectionString string) (plc4go.PlcConnection, error)
	Close() error
}

func NewPlcConnectionCache(driverManager plc4go.PlcDriverManager, withConnectionCacheOptions ...WithConnectionCacheOption) PlcConnectionCache {
	var log zerolog.Logger
	if !config.TraceConnectionCache {
		log = zerolog.Nop()
	}
	maxLeaseTime := 5 * time.Second
	cc := &plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  maxLeaseTime,
		maxWaitTime:   maxLeaseTime * 5,
		cacheLock:     &sync.RWMutex{},
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
		log:           log,
		// _options:   _options, // TODO: we might want to migrate the connection cache options to proper options
	}
	for _, option := range withConnectionCacheOptions {
		option(cc)
	}
	return cc
}

type WithConnectionCacheOption func(plcConnectionCache *plcConnectionCache)

func WithMaxLeaseTime(maxLeaseTime time.Duration) WithConnectionCacheOption {
	return func(plcConnectionCache *plcConnectionCache) {
		plcConnectionCache.maxLeaseTime = maxLeaseTime
	}
}

func WithMaxWaitTime(maxWaitTime time.Duration) WithConnectionCacheOption {
	return func(plcConnectionCache *plcConnectionCache) {
		plcConnectionCache.maxWaitTime = maxWaitTime
	}
}

func WithTracer() WithConnectionCacheOption {
	return func(plcConnectionCache *plcConnectionCache) {
		plcConnectionCache.EnableTracer()
	}
}

// Deprecated: use WithCustomLogger
func WithLogger(logger zerolog.Logger) WithConnectionCacheOption {
	return WithCustomLogger(logger)
}

func WithCustomLogger(logger zerolog.Logger) WithConnectionCacheOption {
	return func(plcConnectionCache *plcConnectionCache) {
		plcConnectionCache.log = logger
	}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type plcConnectionCache struct {
	driverManager plc4go.PlcDriverManager

	// Maximum duration a connection can be used per lease.
	// If the connection is used for a longer time, it is forcefully removed from the client.
	maxLeaseTime time.Duration
	maxWaitTime  time.Duration

	cacheLock   *sync.RWMutex
	connections map[string]*connectionContainer
	tracer      tracer.Tracer

	wg sync.WaitGroup // use to track spawned go routines

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

func (c *plcConnectionCache) onConnectionEvent(event connectionEvent) {
	connectionContainerInstance := event.getConnectionContainer()
	if errorEvent, ok := event.(*connectionErrorEvent); ok {
		if c.tracer != nil {
			c.tracer.AddTrace("destroy-connection", errorEvent.getError().Error())
		}
		c.log.Debug().Str("connectionString", connectionContainerInstance.connectionString)
	}
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (c *plcConnectionCache) EnableTracer() {
	c.tracer = tracer.NewTracer(
		"cache",
		append(c._options, options.WithCustomLogger(c.log))...,
	)
}

func (c *plcConnectionCache) GetTracer() tracer.Tracer {
	return c.tracer
}

func (c *plcConnectionCache) GetConnection(ctx context.Context, connectionString string) (plc4go.PlcConnection, error) {
	c.cacheLock.Lock()

	// If a connection for this connection string didn't exist yet, create a new container
	// and make that container connect.
	if _, ok := c.connections[connectionString]; !ok {
		if c.tracer != nil {
			c.tracer.AddTrace("get-connection", "create new cached connection")
		}
		c.log.Debug().Str("connectionString", connectionString).Msg("Create new cached connection")
		// Create a new connection container.
		cc := newConnectionContainer(c.log, c.driverManager, connectionString)
		// Register for connection events (Like connection closed or error).
		cc.addListener(c)
		// Store the new connection container in the cache of connections.
		c.connections[connectionString] = cc
		// Initialize the connection itself.
		c.wg.Go(func() {
			cc.connect(ctx)
		})
	}

	// Get the ConnectionContainer for this connection string.
	connection := c.connections[connectionString]

	// Release the lock again.
	c.cacheLock.Unlock()

	// Try to get a lease on this connection.
	var txId string
	if c.tracer != nil {
		txId = c.tracer.AddTransactionalStartTrace("get-connection", "lease")
	}
	connChan, errChan := connection.lease(ctx)
	maximumWaitTimeout := time.NewTimer(c.maxWaitTime)
	select {
	case <-ctx.Done(): // abort on context cancel
		return nil, ctx.Err()

	case conn := <-connChan: // Wait till we get a lease.
		c.log.Debug().
			Str("connectionString", connectionString).
			Stringer("conn", conn).
			Msg("Successfully got lease to connection")
		if c.tracer != nil {
			c.tracer.AddTransactionalTrace(txId, "get-connection", "success")
		}
		if tie, ok := conn.connection.(spi.TransportInstanceExposer); ok {
			c.log.Trace().Msg("Resetting transport instance")
			tie.GetTransportInstance().Reset()
		}
		return conn, nil

	case err := <-errChan:
		return nil, errors.Wrap(err, "error while trying to get lease on connection")

	case <-ctx.Done():
		return nil, ctx.Err()

	case <-maximumWaitTimeout.C: // Timeout after the maximum waiting time.
		// In this case we need to drain the chan and return it immediate
		c.wg.Go(func() {
			select {
			case <-connChan:
			case <-errChan:
			}
			_ = connection.returnConnection(ctx, StateIdle)
		})
		if c.tracer != nil {
			c.tracer.AddTransactionalTrace(txId, "get-connection", "timeout")
		}
		c.log.Debug().Str("connectionString", connectionString).Msg("Timeout while waiting for connection.")
		return nil, errors.New("timeout while waiting for connection")
	}
}

func (c *plcConnectionCache) Close() error {
	ctx := context.TODO()
	c.log.Debug().Msg("Closing connection cache started.")
	c.log.Trace().Msg("Acquire lock")
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()
	c.log.Trace().Msg("lock acquired")

	if len(c.connections) == 0 {
		c.log.Debug().Msg("Closing connection cache finished.")
		return nil
	}

	var wg sync.WaitGroup
	for _, connectionContainer := range c.connections {
		ccLog := c.log.With().Interface("connectionContainer", connectionContainer).Logger()
		ccLog.Trace().Msg("Closing connection")
		// Mark the connection as being closed to not try to re-establish it.
		connectionContainer.closed = true
		wg.Go(func() {
			defer func() {
				if err := recover(); err != nil {
					c.log.Error().
						Str("stack", string(debug.Stack())).
						Interface("err", err).
						Msg("panic-ed")
				}
			}()
			// Try to get a lease as this way we kow we're not closing the connection
			// while some go func is still using it.
			ccLog.Trace().Msg("getting a lease")
			ctx, cancel := context.WithTimeout(ctx, c.maxWaitTime)
			connChan, errChan := connectionContainer.lease(ctx)
			select {
			// We're just getting the lease as this way we can be sure nobody else is using it.
			// We also really don't care if it worked, or not ... it's just an attempt of being
			// nice.
			case _ = <-connChan:
				ccLog.Debug().Msg("Gracefully closing connection ...")
				// Give back the connection.
				if connectionContainer.connection != nil {
					ccLog.Trace().Msg("closing actual connection")
					if err := connectionContainer.connection.Close(); err != nil {
						ccLog.Debug().Err(err).Msg("Error while closing connection")
					}
				}
			case err := <-errChan:
				ccLog.Debug().Err(err).Msg("Error while trying to get lease on connection, ignoring.")
			// If we're timing out brutally kill the connection.
			case <-ctx.Done():
				ccLog.Debug().Msg("Forcefully closing connection ...")
				// Forcefully close this connection.
				if connectionContainer.connection != nil {
					if err := connectionContainer.connection.Close(); err != nil {
						ccLog.Debug().Err(err).Msg("Error while closing connection")
					}
				}
			}
			cancel()
		})
	}
	wg.Wait()
	c.log.Debug().Msg("Closing connection cache finished.")
	return nil
}

func (c *plcConnectionCache) String() string {
	return fmt.Sprintf("plcConnectionCache{driverManager: %s, maxLeaseTime: %s, maxWaitTime: %s, connections: %s, tracer: %s}", c.driverManager, c.maxLeaseTime, c.maxWaitTime, c.connections, c.tracer)
}
