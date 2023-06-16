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
		log:           log,
		driverManager: driverManager,
		maxLeaseTime:  maxLeaseTime,
		maxWaitTime:   maxLeaseTime * 5,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
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

	log zerolog.Logger
}

func (t *plcConnectionCache) onConnectionEvent(event connectionEvent) {
	connectionContainerInstance := event.getConnectionContainer()
	if errorEvent, ok := event.(connectionErrorEvent); ok {
		if t.tracer != nil {
			t.tracer.AddTrace("destroy-connection", errorEvent.getError().Error())
		}
		t.log.Debug().Str("connectionString", connectionContainerInstance.connectionString)
	}
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (t *plcConnectionCache) EnableTracer() {
	t.tracer = tracer.NewTracer("cache", options.WithCustomLogger(t.log))
}

func (t *plcConnectionCache) GetTracer() tracer.Tracer {
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
			t.log.Debug().Str("connectionString", connectionString).Msg("Create new cached connection")
			// Create a new connection container.
			cc := newConnectionContainer(t.log, t.driverManager, connectionString)
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
		maximumWaitTimeout := time.NewTimer(t.maxWaitTime)
		defer utils.CleanupTimer(maximumWaitTimeout)
		select {
		// Wait till we get a lease.
		case connectionResponse := <-leaseChan:
			t.log.Debug().Str("connectionString", connectionString).Msg("Successfully got lease to connection")
			responseTimeout := time.NewTimer(10 * time.Millisecond)
			defer utils.CleanupTimer(responseTimeout)
			select {
			case ch <- connectionResponse:
				if t.tracer != nil {
					t.tracer.AddTransactionalTrace(txId, "get-connection", "success")
				}
			case <-responseTimeout.C:
				// Log a message, that the client has given up
				if t.tracer != nil {
					t.tracer.AddTransactionalTrace(txId, "get-connection", "client given up")
				}
				close(ch)
				t.log.Debug().Str("connectionString", connectionString).Msg("Client not available returning connection to cache.")
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
			if t.tracer != nil {
				t.tracer.AddTransactionalTrace(txId, "get-connection", "timeout")
			}
			t.log.Debug().Str("connectionString", connectionString).Msg("Timeout while waiting for connection.")
			ch <- _default.NewDefaultPlcConnectionCloseResult(nil, errors.New("timeout while waiting for connection"))
		}
	}()

	return ch
}

func (t *plcConnectionCache) Close() <-chan PlcConnectionCacheCloseResult {
	t.log.Debug().Msg("Closing connection cache started.")
	ch := make(chan PlcConnectionCacheCloseResult)

	go func() {
		t.cacheLock.Lock()
		defer t.cacheLock.Unlock()

		if len(t.connections) == 0 {
			responseDeliveryTimeout := time.NewTimer(10 * time.Millisecond)
			defer utils.CleanupTimer(responseDeliveryTimeout)
			select {
			case ch <- newDefaultPlcConnectionCacheCloseResult(t, nil):
			case <-responseDeliveryTimeout.C:
			}
			t.log.Debug().Msg("Closing connection cache finished.")
			return
		}

		for _, cc := range t.connections {
			// Mark the connection as being closed to not try to re-establish it.
			cc.closed = true
			// Try to get a lease as this way we kow we're not closing the connection
			// while some go func is still using it.
			go func(container *connectionContainer) {
				leaseResults := container.lease()
				closeTimeout := time.NewTimer(t.maxWaitTime)
				defer utils.CleanupTimer(closeTimeout)
				select {
				// We're just getting the lease as this way we can be sure nobody else is using it.
				// We also really don't care if it worked, or not ... it's just an attempt of being
				// nice.
				case _ = <-leaseResults:
					t.log.Debug().Str("connectionString", container.connectionString).Msg("Gracefully closing connection ...")
					// Give back the connection.
					if container.connection != nil {
						container.connection.Close()
					}
				// If we're timing out brutally kill the connection.
				case <-closeTimeout.C:
					t.log.Debug().Str("connectionString", container.connectionString).Msg("Forcefully closing connection ...")
					// Forcefully close this connection.
					if container.connection != nil {
						container.connection.Close()
					}
				}

				responseDeliveryTimeout := time.NewTimer(10 * time.Millisecond)
				defer utils.CleanupTimer(responseDeliveryTimeout)
				select {
				case ch <- newDefaultPlcConnectionCacheCloseResult(t, nil):
				case <-responseDeliveryTimeout.C:
				}
				t.log.Debug().Msg("Closing connection cache finished.")
			}(cc)
		}
	}()

	return ch
}

func (t *plcConnectionCache) String() string {
	return fmt.Sprintf("plcConnectionCache{driverManager: %s, maxLeaseTime: %s, maxWaitTime: %s, connections: %s, tracer: %s}", t.driverManager, t.maxLeaseTime, t.maxWaitTime, t.connections, t.tracer)
}
