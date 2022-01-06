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
	connections map[string]*PooledPlcConnection
}

func NewPlcConnectionPool(driverManager plc4go.PlcDriverManager) *PlcConnectionPool {
	return NewPlcConnectionPoolWithMaxLeaseTime(driverManager, time.Second*20)
}

func NewPlcConnectionPoolWithMaxLeaseTime(driverManager plc4go.PlcDriverManager, maxLeaseTime time.Duration) *PlcConnectionPool {
	return &PlcConnectionPool{
		driverManager: driverManager,
		maxLeaseTime:  maxLeaseTime,
		maxWaitTime:   maxLeaseTime * 5,
		poolLock:      lock.NewCASMutex(),
		connections:   make(map[string]*PooledPlcConnection),
	}
}

func (t *PlcConnectionPool) GetConnection(connectionString string) <-chan plc4go.PlcConnectionConnectResult {
	ch := make(chan plc4go.PlcConnectionConnectResult)

	go func() {
		// Check if we've already got this connection
		t.poolLock.Lock()
		defer t.poolLock.Unlock()
		// Try to get a connection for the given url.
		// If this returns ok, this means there is already a connection
		// available, so we'll try to use that. If it fails, we need
		// to create a completely new connection.
		if connection, ok := t.connections[connectionString]; ok {
			// Use an existing connection

			// Check if the connection is available.
			if connection.state == StateIdle {
				// As soon as we have the lock, return the connection.
				pooledConnection := t.connections[connectionString]
				ch <- _default.NewDefaultPlcConnectionConnectResult(pooledConnection, nil)
			} else if connection.state == StateInUse {
				// If the connection is currently busy, add the new channel to the queue for this connection.
				pooledConnection := t.connections[connectionString]
				pooledConnection.enqueue(ch)
			} else {
				ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.New("timeout"))
			}
		} else {
			// Create a new connection.

			pooledPlcConnection := NewPooledPlcConnection(connectionString, t)
			t.connections[connectionString] = pooledPlcConnection

			// Initialize the new connection.
			connectionResultChan := t.driverManager.GetConnection(connectionString)

			// Allow us to finish this function and return the lock quickly
			go func() {
				// Wait for the connection to be established.
				connectionResult := <-connectionResultChan

				// If the connection was successful, pass the active connection into the container.
				// If something went wrong, we have to remove the connection from the pool and return the error.
				if connectionResult.GetErr() == nil {
					// Make sure we have the lock here so all following operations will execute atomically.
					pooledPlcConnection.lock.Lock()
					defer pooledPlcConnection.lock.Unlock()

					// Inject the real connection into the container.
					pooledPlcConnection.activeConnection = connectionResult.GetConnection()
					// Mark the connection as being used.
					pooledPlcConnection.state = StateInUse

					// Return the pooled connection to the client.
					ch <- _default.NewDefaultPlcConnectionCloseResult(pooledPlcConnection, nil)
				} else {
					// Mark the connection as broken.
					pooledPlcConnection.state = StateInvalid

					// Remove the broken connection from the pool.
					t.poolLock.Lock()
					defer t.poolLock.Unlock()
					delete(t.connections, connectionString)

					// Forward the error to the client.
					ch <- _default.NewDefaultPlcConnectionConnectResult(nil, connectionResult.GetErr())
				}
			}()
		}
	}()

	return ch
}

func (t *PlcConnectionPool) returnConnection(pooledConnection *PooledPlcConnection) error {
	// If the connection is marked as "invalid", destroy it and remove it from the pool.
	if pooledConnection.state == StateInvalid {
		// At least try to close the invalid connection.
		pooledConnection.activeConnection.Close()

		// TODO: Either try to reconnect or cancel all waiting connections

		// Remove the connection from the pool.
		t.poolLock.Lock()
		defer t.poolLock.Unlock()
		delete(t.connections, pooledConnection.connectionString)
		return nil
	}

	pooledConnection.lock.Lock()
	defer pooledConnection.lock.Unlock()
	// Check how many others are waiting for this connection.
	if len(pooledConnection.queue) > 0 {
		// There are waiting clients, give the connection to the next client in the line.
		next := pooledConnection.queue[0]
		pooledConnection.queue = pooledConnection.queue[1:]
		next <- _default.NewDefaultPlcConnectionConnectResult(pooledConnection, nil)
	} else {
		// Otherwise, just mark the connection as idle.
		pooledConnection.state = StateIdle
	}

	return nil
}
