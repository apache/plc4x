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
	"github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/config"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/tracer"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/viney-shih/go-lock"
	"time"
)

type PlcConnectionCache interface {
	GetConnection(connectionString string) <-chan plc4go.PlcConnectionConnectResult
	Close() <-chan PlcConnectionCacheCloseResult
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
		cacheLock:     lock.NewCASMutex(),
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

func WithMaxLeaseTime(duration time.Duration) WithConnectionCacheOption {
	return func(plcConnectionCache *plcConnectionCache) {
		plcConnectionCache.maxLeaseTime = duration
	}
}

func WithMaxWaitTime(duration time.Duration) WithConnectionCacheOption {
	return func(plcConnectionCache *plcConnectionCache) {
		plcConnectionCache.maxLeaseTime = duration
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

	cacheLock   lock.RWMutex
	connections map[string]*connectionContainer
	tracer      tracer.Tracer

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

func (c *plcConnectionCache) onConnectionEvent(event connectionEvent) {
	connectionContainerInstance := event.getConnectionContainer()
	if errorEvent, ok := event.(connectionErrorEvent); ok {
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

func (c *plcConnectionCache) GetConnection(connectionString string) <-chan plc4go.PlcConnectionConnectResult {
	ch := make(chan plc4go.PlcConnectionConnectResult)

	go func() {
		c.cacheLock.Lock()

		// If a connection for this connection string didn'c exist yet, create a new container
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
			go func(cc2 *connectionContainer) {
				cc2.connect()
			}(cc)
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
		leaseChan := connection.lease()
		maximumWaitTimeout := time.NewTimer(c.maxWaitTime)
		defer utils.CleanupTimer(maximumWaitTimeout)
		select {
		// Wait till we get a lease.
		case connectionResponse := <-leaseChan:
			c.log.Debug().Str("connectionString", connectionString).Msg("Successfully got lease to connection")
			responseTimeout := time.NewTimer(10 * time.Millisecond)
			defer utils.CleanupTimer(responseTimeout)
			select {
			case ch <- connectionResponse:
				if c.tracer != nil {
					c.tracer.AddTransactionalTrace(txId, "get-connection", "success")
				}
			case <-responseTimeout.C:
				// Log a message, that the client has given up
				if c.tracer != nil {
					c.tracer.AddTransactionalTrace(txId, "get-connection", "client given up")
				}
				close(ch)
				c.log.Debug().Str("connectionString", connectionString).Msg("Client not available returning connection to cache.")
				// Return the connection to give another connection the chance to use it.
				if connectionResponse.GetConnection() != nil {
					connectionResponse.GetConnection().Close()
				}
			}

		// Timeout after the maximum waiting time.
		case <-maximumWaitTimeout.C:
			// In this case we need to drain the chan and return it immediate
			go func() {
				<-leaseChan
				_ = connection.returnConnection(StateIdle)
			}()
			if c.tracer != nil {
				c.tracer.AddTransactionalTrace(txId, "get-connection", "timeout")
			}
			c.log.Debug().Str("connectionString", connectionString).Msg("Timeout while waiting for connection.")
			ch <- _default.NewDefaultPlcConnectionCloseResult(nil, errors.New("timeout while waiting for connection"))
		}
	}()

	return ch
}

func (c *plcConnectionCache) Close() <-chan PlcConnectionCacheCloseResult {
	c.log.Debug().Msg("Closing connection cache started.")
	ch := make(chan PlcConnectionCacheCloseResult)

	go func() {
		c.cacheLock.Lock()
		defer c.cacheLock.Unlock()

		if len(c.connections) == 0 {
			responseDeliveryTimeout := time.NewTimer(10 * time.Millisecond)
			defer utils.CleanupTimer(responseDeliveryTimeout)
			select {
			case ch <- newDefaultPlcConnectionCacheCloseResult(c, nil):
			case <-responseDeliveryTimeout.C:
			}
			c.log.Debug().Msg("Closing connection cache finished.")
			return
		}

		for _, cc := range c.connections {
			// Mark the connection as being closed to not try to re-establish it.
			cc.closed = true
			// Try to get a lease as this way we kow we're not closing the connection
			// while some go func is still using it.
			go func(container *connectionContainer) {
				leaseResults := container.lease()
				closeTimeout := time.NewTimer(c.maxWaitTime)
				defer utils.CleanupTimer(closeTimeout)
				select {
				// We're just getting the lease as this way we can be sure nobody else is using it.
				// We also really don'c care if it worked, or not ... it's just an attempt of being
				// nice.
				case _ = <-leaseResults:
					c.log.Debug().Str("connectionString", container.connectionString).Msg("Gracefully closing connection ...")
					// Give back the connection.
					if container.connection != nil {
						container.connection.Close()
					}
				// If we're timing out brutally kill the connection.
				case <-closeTimeout.C:
					c.log.Debug().Str("connectionString", container.connectionString).Msg("Forcefully closing connection ...")
					// Forcefully close this connection.
					if container.connection != nil {
						container.connection.Close()
					}
				}

				responseDeliveryTimeout := time.NewTimer(10 * time.Millisecond)
				defer utils.CleanupTimer(responseDeliveryTimeout)
				select {
				case ch <- newDefaultPlcConnectionCacheCloseResult(c, nil):
				case <-responseDeliveryTimeout.C:
				}
				c.log.Debug().Msg("Closing connection cache finished.")
			}(cc)
		}
	}()

	return ch
}

func (c *plcConnectionCache) String() string {
	return fmt.Sprintf("plcConnectionCache{driverManager: %s, maxLeaseTime: %s, maxWaitTime: %s, connections: %s, tracer: %s}", c.driverManager, c.maxLeaseTime, c.maxWaitTime, c.connections, c.tracer)
}
