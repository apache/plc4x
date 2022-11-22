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
	"github.com/apache/plc4x/plc4go/internal/simulated"
	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/stretchr/testify/assert"
	"github.com/viney-shih/go-lock"
	"testing"
	"time"
)

func TestLeasedPlcConnection_IsTraceEnabled(t *testing.T) {
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(simulated.NewDriver())
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  time.Second * 1,
		maxWaitTime:   time.Second * 5,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	cache.EnableTracer()

	// The first and second connection should work fine
	connectionResults := cache.GetConnection("simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	select {
	case connectionResult := <-connectionResults:
		if assert.NotNil(t, connectionResult) {
			assert.Nil(t, connectionResult.GetErr())
			if assert.NotNil(t, connectionResult.GetConnection()) {
				connection := connectionResult.GetConnection().(spi.PlcConnection)
				assert.True(t, connection.IsTraceEnabled())
				connection.BlockingClose()
				func() {
					defer func() {
						if r := recover(); r != nil {
							assert.Equal(t, r, "Called 'IsTraceEnabled' on a closed cached connection")
						} else {
							t.Errorf("The code did not panic")
						}
					}()
					connection.IsTraceEnabled()
				}()
			}
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout")
	}

	// The first and second connection should work fine
	connectionResults = cache.GetConnection("simulated://1.2.3.4:42?connectionDelay=100")
	select {
	case connectionResult := <-connectionResults:
		if assert.NotNil(t, connectionResult) {
			assert.Nil(t, connectionResult.GetErr())
			if assert.NotNil(t, connectionResult.GetConnection()) {
				connection := connectionResult.GetConnection().(spi.PlcConnection)
				assert.False(t, connection.IsTraceEnabled())
				connection.BlockingClose()
				func() {
					defer func() {
						if r := recover(); r != nil {
							assert.Equal(t, r, "Called 'IsTraceEnabled' on a closed cached connection")
						} else {
							t.Errorf("The code did not panic")
						}
					}()
					connection.IsTraceEnabled()
				}()
			}
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout")
	}
}

func TestLeasedPlcConnection_GetTracer(t *testing.T) {
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(simulated.NewDriver())
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  time.Second * 1,
		maxWaitTime:   time.Second * 5,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	cache.EnableTracer()

	// The first and second connection should work fine
	connectionResults := cache.GetConnection("simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	select {
	case connectionResult := <-connectionResults:
		if assert.NotNil(t, connectionResult) {
			assert.Nil(t, connectionResult.GetErr())
			if assert.NotNil(t, connectionResult.GetConnection()) {
				connection := connectionResult.GetConnection().(spi.PlcConnection)
				assert.NotNil(t, connection.GetTracer())
				connection.BlockingClose()
				func() {
					defer func() {
						if r := recover(); r != nil {
							assert.Equal(t, r, "Called 'GetTracer' on a closed cached connection")
						} else {
							t.Errorf("The code did not panic")
						}
					}()
					connection.GetTracer()
				}()
			}
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout")
	}
}

func TestLeasedPlcConnection_GetConnectionId(t *testing.T) {
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(simulated.NewDriver())
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  time.Second * 1,
		maxWaitTime:   time.Second * 5,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	cache.EnableTracer()

	// The first and second connection should work fine
	connectionResults := cache.GetConnection("simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	select {
	case connectionResult := <-connectionResults:
		if assert.NotNil(t, connectionResult) {
			assert.Nil(t, connectionResult.GetErr())
			if assert.NotNil(t, connectionResult.GetConnection()) {
				connection := connectionResult.GetConnection().(spi.PlcConnection)
				assert.Greater(t, len(connection.GetConnectionId()), 0)
				connection.BlockingClose()
				func() {
					defer func() {
						if r := recover(); r != nil {
							assert.Equal(t, r, "Called 'GetConnectionId' on a closed cached connection")
						} else {
							t.Errorf("The code did not panic")
						}
					}()
					connection.GetConnectionId()
				}()
			}
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout")
	}
}

func TestLeasedPlcConnection_Connect(t *testing.T) {
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(simulated.NewDriver())
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  time.Second * 1,
		maxWaitTime:   time.Second * 5,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	cache.EnableTracer()

	// The first and second connection should work fine
	connectionResults := cache.GetConnection("simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	select {
	case connectionResult := <-connectionResults:
		if assert.NotNil(t, connectionResult) {
			assert.Nil(t, connectionResult.GetErr())
			if assert.NotNil(t, connectionResult.GetConnection()) {
				connection := connectionResult.GetConnection().(spi.PlcConnection)
				func() {
					defer func() {
						if r := recover(); r != nil {
							assert.Equal(t, r, "Called 'Connect' on a cached connection")
						} else {
							t.Errorf("The code did not panic")
						}
					}()
					connection.Connect()
				}()
			}
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout")
	}
}

func TestLeasedPlcConnection_BlockingClose(t *testing.T) {
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(simulated.NewDriver())
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  time.Second * 1,
		maxWaitTime:   time.Second * 5,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	cache.EnableTracer()

	// The first and second connection should work fine
	connectionResults := cache.GetConnection("simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	select {
	case connectionResult := <-connectionResults:
		if assert.NotNil(t, connectionResult) {
			assert.Nil(t, connectionResult.GetErr())
			if assert.NotNil(t, connectionResult.GetConnection()) {
				connection := connectionResult.GetConnection()
				connection.BlockingClose()
				func() {
					defer func() {
						if r := recover(); r != nil {
							assert.Equal(t, r, "Called 'BlockingClose' on a closed cached connection")
						} else {
							t.Errorf("The code did not panic")
						}
					}()
					connection.BlockingClose()
				}()
			}
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout")
	}
}

func TestLeasedPlcConnection_Close(t *testing.T) {
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(simulated.NewDriver())
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  time.Second * 1,
		maxWaitTime:   time.Second * 5,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	cache.EnableTracer()

	// The first and second connection should work fine
	connectionResults := cache.GetConnection("simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	select {
	case connectionResult := <-connectionResults:
		if assert.NotNil(t, connectionResult) {
			assert.Nil(t, connectionResult.GetErr())
			if assert.NotNil(t, connectionResult.GetConnection()) {
				connection := connectionResult.GetConnection()
				connection.BlockingClose()
				func() {
					defer func() {
						if r := recover(); r != nil {
							assert.Equal(t, r, "Called 'Close' on a closed cached connection")
						} else {
							t.Errorf("The code did not panic")
						}
					}()
					connection.Close()
				}()
			}
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout")
	}
}

func TestLeasedPlcConnection_IsConnected(t *testing.T) {
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(simulated.NewDriver())
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  time.Second * 1,
		maxWaitTime:   time.Second * 5,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	cache.EnableTracer()

	// The first and second connection should work fine
	connectionResults := cache.GetConnection("simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	select {
	case connectionResult := <-connectionResults:
		if assert.NotNil(t, connectionResult) {
			assert.Nil(t, connectionResult.GetErr())
			if assert.NotNil(t, connectionResult.GetConnection()) {
				connection := connectionResult.GetConnection()
				assert.True(t, connection.IsConnected())
				connection.BlockingClose()
				assert.False(t, connection.IsConnected())
			}
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout")
	}
}

func TestLeasedPlcConnection_Ping(t *testing.T) {
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(simulated.NewDriver())
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  time.Second * 1,
		maxWaitTime:   time.Second * 5,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	cache.EnableTracer()

	// The first and second connection should work fine
	connectionResults := cache.GetConnection("simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	select {
	case connectionResult := <-connectionResults:
		if assert.NotNil(t, connectionResult) {
			assert.Nil(t, connectionResult.GetErr())
			if assert.NotNil(t, connectionResult.GetConnection()) {
				connection := connectionResult.GetConnection()
				connection.Ping()
				connection.BlockingClose()
				func() {
					defer func() {
						if r := recover(); r != nil {
							assert.Equal(t, r, "Called 'Ping' on a closed cached connection")
						} else {
							t.Errorf("The code did not panic")
						}
					}()
					connection.Ping()
				}()
			}
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout")
	}
}

func TestLeasedPlcConnection_GetMetadata(t *testing.T) {
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(simulated.NewDriver())
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  time.Second * 1,
		maxWaitTime:   time.Second * 5,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	cache.EnableTracer()

	// The first and second connection should work fine
	connectionResults := cache.GetConnection("simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	select {
	case connectionResult := <-connectionResults:
		if assert.NotNil(t, connectionResult) {
			assert.Nil(t, connectionResult.GetErr())
			if assert.NotNil(t, connectionResult.GetConnection()) {
				connection := connectionResult.GetConnection()
				metadata := connection.GetMetadata()
				if assert.NotNil(t, metadata) {
					attributes := metadata.GetConnectionAttributes()
					assert.NotNil(t, attributes)
				}
				connection.BlockingClose()
				func() {
					defer func() {
						if r := recover(); r != nil {
							assert.Equal(t, r, "Called 'GetMetadata' on a closed cached connection")
						} else {
							t.Errorf("The code did not panic")
						}
					}()
					connection.GetMetadata()
				}()
			}
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout")
	}
}

func TestLeasedPlcConnection_ReadRequestBuilder(t *testing.T) {
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(simulated.NewDriver())
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  time.Second * 1,
		maxWaitTime:   time.Second * 5,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	cache.EnableTracer()

	// The first and second connection should work fine
	connectionResults := cache.GetConnection("simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	select {
	case connectionResult := <-connectionResults:
		if assert.NotNil(t, connectionResult) {
			assert.Nil(t, connectionResult.GetErr())
			if assert.NotNil(t, connectionResult.GetConnection()) {
				connection := connectionResult.GetConnection()
				builder := connection.ReadRequestBuilder()
				assert.NotNil(t, builder)
				connection.BlockingClose()
				func() {
					defer func() {
						if r := recover(); r != nil {
							assert.Equal(t, r, "Called 'ReadRequestBuilder' on a closed cached connection")
						} else {
							t.Errorf("The code did not panic")
						}
					}()
					connection.ReadRequestBuilder()
				}()
			}
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout")
	}
}

func TestLeasedPlcConnection_WriteRequestBuilder(t *testing.T) {
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(simulated.NewDriver())
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  time.Second * 1,
		maxWaitTime:   time.Second * 5,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	cache.EnableTracer()

	// The first and second connection should work fine
	connectionResults := cache.GetConnection("simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	select {
	case connectionResult := <-connectionResults:
		if assert.NotNil(t, connectionResult) {
			assert.Nil(t, connectionResult.GetErr())
			if assert.NotNil(t, connectionResult.GetConnection()) {
				connection := connectionResult.GetConnection()
				builder := connection.WriteRequestBuilder()
				assert.NotNil(t, builder)
				connection.BlockingClose()
				func() {
					defer func() {
						if r := recover(); r != nil {
							assert.Equal(t, r, "Called 'WriteRequestBuilder' on a closed cached connection")
						} else {
							t.Errorf("The code did not panic")
						}
					}()
					connection.WriteRequestBuilder()
				}()
			}
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout")
	}
}

func TestLeasedPlcConnection_SubscriptionRequestBuilder(t *testing.T) {
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(simulated.NewDriver())
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  time.Second * 1,
		maxWaitTime:   time.Second * 5,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	cache.EnableTracer()

	// The first and second connection should work fine
	connectionResults := cache.GetConnection("simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	select {
	case connectionResult := <-connectionResults:
		if assert.NotNil(t, connectionResult) {
			assert.Nil(t, connectionResult.GetErr())
			if assert.NotNil(t, connectionResult.GetConnection()) {
				connection := connectionResult.GetConnection()
				func() {
					defer func() {
						if r := recover(); r != nil {
							assert.Equal(t, r, "not implemented")
						} else {
							t.Errorf("The code did not panic")
						}
					}()
					connection.SubscriptionRequestBuilder()
				}()
				connection.BlockingClose()
				func() {
					defer func() {
						if r := recover(); r != nil {
							assert.Equal(t, r, "Called 'SubscriptionRequestBuilder' on a closed cached connection")
						} else {
							t.Errorf("The code did not panic")
						}
					}()
					connection.SubscriptionRequestBuilder()
				}()
			}
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout")
	}
}

func TestLeasedPlcConnection_UnsubscriptionRequestBuilder(t *testing.T) {
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(simulated.NewDriver())
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  time.Second * 1,
		maxWaitTime:   time.Second * 5,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	cache.EnableTracer()

	// The first and second connection should work fine
	connectionResults := cache.GetConnection("simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	select {
	case connectionResult := <-connectionResults:
		if assert.NotNil(t, connectionResult) {
			assert.Nil(t, connectionResult.GetErr())
			if assert.NotNil(t, connectionResult.GetConnection()) {
				connection := connectionResult.GetConnection()
				func() {
					defer func() {
						if r := recover(); r != nil {
							assert.Equal(t, r, "not implemented")
						} else {
							t.Errorf("The code did not panic")
						}
					}()
					connection.UnsubscriptionRequestBuilder()
				}()
				connection.BlockingClose()
				func() {
					defer func() {
						if r := recover(); r != nil {
							assert.Equal(t, r, "Called 'UnsubscriptionRequestBuilder' on a closed cached connection")
						} else {
							t.Errorf("The code did not panic")
						}
					}()
					connection.UnsubscriptionRequestBuilder()
				}()
			}
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout")
	}
}

func TestLeasedPlcConnection_BrowseRequestBuilder(t *testing.T) {
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(simulated.NewDriver())
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  time.Second * 1,
		maxWaitTime:   time.Second * 5,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	cache.EnableTracer()

	// The first and second connection should work fine
	connectionResults := cache.GetConnection("simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	select {
	case connectionResult := <-connectionResults:
		if assert.NotNil(t, connectionResult) {
			assert.Nil(t, connectionResult.GetErr())
			if assert.NotNil(t, connectionResult.GetConnection()) {
				connection := connectionResult.GetConnection()
				func() {
					defer func() {
						if r := recover(); r != nil {
							assert.Equal(t, r, "not implemented")
						} else {
							t.Errorf("The code did not panic")
						}
					}()
					connection.BrowseRequestBuilder()
				}()
				connection.BlockingClose()
				func() {
					defer func() {
						if r := recover(); r != nil {
							assert.Equal(t, r, "Called 'BrowseRequestBuilder' on a closed cached connection")
						} else {
							t.Errorf("The code did not panic")
						}
					}()
					connection.BrowseRequestBuilder()
				}()
			}
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout")
	}
}
