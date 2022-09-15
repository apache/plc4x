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
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
	"github.com/viney-shih/go-lock"
	"time"
)

type PlcConnectionCache interface {
	GetConnection(connectionString string) <-chan plc4go.PlcConnectionConnectResult
	Close() <-chan PlcConnectionCacheCloseResult
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
	tracer      *spi.Tracer
}

func (t *plcConnectionCache) onConnectionEvent(event connectionEvent) {
	setCacheLog()
	connectionContainerInstance := event.getConnectionContainer()
	if errorEvent, ok := event.(connectionErrorEvent); ok {
		if t.tracer != nil {
			t.tracer.AddTrace("destroy-connection", errorEvent.getError().Error())
		}
		cacheLog.Debug().Str("connectionString", connectionContainerInstance.connectionString)
	}
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (t *plcConnectionCache) EnableTracer() {
	t.tracer = spi.NewTracer("cache")
}

func (t *plcConnectionCache) GetTracer() *spi.Tracer {
	return t.tracer
}

func (t *plcConnectionCache) GetConnection(connectionString string) <-chan plc4go.PlcConnectionConnectResult {
	setCacheLog()
	ch := make(chan plc4go.PlcConnectionConnectResult)

	go func() {
		t.cacheLock.Lock()

		// If a connection for this connection string didn't exist yet, create a new container
		// and make that container connect.
		if _, ok := t.connections[connectionString]; !ok {
			if t.tracer != nil {
				t.tracer.AddTrace("get-connection", "create new cached connection")
			}
			cacheLog.Debug().Str("connectionString", connectionString).Msg("Create new cached connection")
			// Create a new connection container.
			cc := newConnectionContainer(t.driverManager, connectionString)
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
			cacheLog.Debug().Str("connectionString", connectionString).Msg("Successfully got lease to connection")
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
				cacheLog.Debug().Str("connectionString", connectionString).Msg("Client not available returning connection to cache.")
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
			cacheLog.Debug().Str("connectionString", connectionString).Msg("Timeout while waiting for connection.")
			ch <- _default.NewDefaultPlcConnectionCloseResult(nil, errors.New("timeout while waiting for connection"))
		}
	}()

	return ch
}

func (t *plcConnectionCache) Close() <-chan PlcConnectionCacheCloseResult {
	setCacheLog()
	cacheLog.Debug().Msg("Closing connection cache started.")
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
			cacheLog.Debug().Msg("Closing connection cache finished.")
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
					cacheLog.Debug().Str("connectionString", container.connectionString).Msg("Gracefully closing connection ...")
					// Give back the connection.
					if container.connection != nil {
						container.connection.Close()
					}
				// If we're timing out brutally kill the connection.
				case <-closeTimeout.C:
					cacheLog.Debug().Str("connectionString", container.connectionString).Msg("Forcefully closing connection ...")
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
				cacheLog.Debug().Msg("Closing connection cache finished.")
			}(cc)
		}
	}()

	return ch
}

func (t *plcConnectionCache) String() string {
	return fmt.Sprintf("plcConnectionCache{driverManager: %s, maxLeaseTime: %s, maxWaitTime: %s, connections: %s, tracer: %s}", t.driverManager, t.maxLeaseTime, t.maxWaitTime, t.connections, t.tracer)
}
