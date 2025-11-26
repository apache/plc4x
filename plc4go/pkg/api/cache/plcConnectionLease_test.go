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
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/viney-shih/go-lock"

	"github.com/apache/plc4x/plc4go/internal/simulated"
	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/config"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/tracer"
)

func TestLeasedPlcConnection_IsTraceEnabled(t *testing.T) {
	logger := testutils.ProduceTestingLogger(t)
	driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  1 * time.Second,
		maxWaitTime:   5 * time.Second,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	t.Cleanup(func() {
		_ = cache.Close()
	})
	cache.EnableTracer()

	// The first and second connection should work fine
	conn, err := cache.GetConnection(t.Context(), "simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	assert.Nil(t, err)
	if assert.NotNil(t, conn) {
		assert.NoError(t, conn.Close())
	}

	// The first and second connection should work fine
	conn, err = cache.GetConnection(t.Context(), "simulated://1.2.3.4:42?connectionDelay=100")
	assert.Nil(t, err)
	assert.NoError(t, conn.Close())
}

func TestLeasedPlcConnection_GetTracer(t *testing.T) {
	logger := testutils.ProduceTestingLogger(t)
	driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  1 * time.Second,
		maxWaitTime:   5 * time.Second,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	t.Cleanup(func() {
		_ = cache.Close()
	})
	cache.EnableTracer()

	// The first and second connection should work fine
	connection, err := cache.GetConnection(t.Context(), "simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	assert.Nil(t, err)
	if assert.NotNil(t, connection) {
		assert.NoError(t, connection.Close())
	}
}

func TestLeasedPlcConnection_GetConnectionId(t *testing.T) {
	logger := testutils.ProduceTestingLogger(t)
	driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  1 * time.Second,
		maxWaitTime:   5 * time.Second,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	t.Cleanup(func() {
		_ = cache.Close()
	})
	cache.EnableTracer()

	// The first and second connection should work fine
	connection, err := cache.GetConnection(t.Context(), "simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	assert.Nil(t, err)
	assert.NoError(t, connection.Close())
}

func TestLeasedPlcConnection_Connect(t *testing.T) {
	logger := testutils.ProduceTestingLogger(t)
	driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  1 * time.Second,
		maxWaitTime:   5 * time.Second,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	t.Cleanup(func() {
		_ = cache.Close()
	})
	cache.EnableTracer()

	// The first and second connection should work fine
	connection, err := cache.GetConnection(t.Context(), "simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	assert.Nil(t, err)
	assert.NotNil(t, connection)
}

func TestLeasedPlcConnection_BlockingClose(t *testing.T) {
	logger := testutils.ProduceTestingLogger(t)
	driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  1 * time.Second,
		maxWaitTime:   5 * time.Second,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	t.Cleanup(func() {
		_ = cache.Close()
	})
	cache.EnableTracer()

	// The first and second connection should work fine
	connection, err := cache.GetConnection(t.Context(), "simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	assert.Nil(t, err)
	if assert.NotNil(t, connection) {
		assert.NoError(t, connection.Close())
	}
}

func TestLeasedPlcConnection_Close(t *testing.T) {
	logger := testutils.ProduceTestingLogger(t)
	driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  1 * time.Second,
		maxWaitTime:   5 * time.Second,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	t.Cleanup(func() {
		_ = cache.Close()
	})
	cache.EnableTracer()

	// The first and second connection should work fine
	connection, err := cache.GetConnection(t.Context(), "simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	assert.Nil(t, err)
	if assert.NotNil(t, connection) {
		assert.NoError(t, connection.Close())
	}
}

func TestLeasedPlcConnection_IsConnected(t *testing.T) {
	logger := testutils.ProduceTestingLogger(t)
	driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  1 * time.Second,
		maxWaitTime:   5 * time.Second,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	t.Cleanup(func() {
		_ = cache.Close()
	})
	cache.EnableTracer()

	// The first and second connection should work fine
	connection, err := cache.GetConnection(t.Context(), "simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	assert.Nil(t, err)
	if assert.NotNil(t, connection) {
		assert.True(t, connection.IsConnected())
		assert.NoError(t, connection.Close())
		assert.False(t, connection.IsConnected())
	}
}

func TestLeasedPlcConnection_Ping(t *testing.T) {
	logger := testutils.ProduceTestingLogger(t)
	driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  1 * time.Second,
		maxWaitTime:   5 * time.Second,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	t.Cleanup(func() {
		_ = cache.Close()
	})
	cache.EnableTracer()

	// The first and second connection should work fine
	connection, err := cache.GetConnection(t.Context(), "simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	assert.Nil(t, err)
	_ = connection.Ping(t.Context())
	assert.NoError(t, connection.Close())
	func() {
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, r, "Called 'Ping' on a closed cached connection")
			} else {
				t.Errorf("The code did not panic")
			}
		}()
		_ = connection.Ping(t.Context())
	}()
}

func TestLeasedPlcConnection_GetMetadata(t *testing.T) {
	logger := testutils.ProduceTestingLogger(t)
	driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  1 * time.Second,
		maxWaitTime:   5 * time.Second,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	t.Cleanup(func() {
		_ = cache.Close()
	})
	cache.EnableTracer()

	// The first and second connection should work fine
	connection, err := cache.GetConnection(t.Context(), "simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	assert.Nil(t, err)
	if assert.NotNil(t, connection) {
		metadata := connection.GetMetadata()
		if assert.NotNil(t, metadata) {
			attributes := metadata.GetConnectionAttributes()
			assert.NotNil(t, attributes)
		}
		assert.NoError(t, connection.Close())
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

func TestLeasedPlcConnection_InvalidateSkipsPingOnClose(t *testing.T) {
	logger := testutils.ProduceTestingLogger(t)
	container := &connectionContainer{
		lock:             lock.NewCASMutex(),
		connectionString: "dummy://invalidate",
		log:              logger,
	}
	dummyConn := &dummyTracedConnection{}
	container.connection = dummyConn
	container.state = StateIdle
	container.driverManager = &dummyDriverManager{
		factory: func() plc4go.PlcConnection { return &dummyTracedConnection{} },
	}
	lease := newPlcConnectionLease(container, 1, dummyConn)

	lease.Invalidate()
	require.True(t, dummyConn.invalidated)

	require.NoError(t, lease.Close())
	require.Zero(t, dummyConn.pingCount, "ping should not be invoked when lease invalidated")
}

func TestLeasedPlcConnection_PingAfterInvalidate(t *testing.T) {
	logger := testutils.ProduceTestingLogger(t)
	container := &connectionContainer{
		lock:             lock.NewCASMutex(),
		connectionString: "dummy://ping",
		log:              logger,
	}
	dummyConn := &dummyTracedConnection{}
	container.connection = dummyConn
	container.state = StateIdle
	container.driverManager = &dummyDriverManager{
		factory: func() plc4go.PlcConnection { return &dummyTracedConnection{} },
	}
	lease := newPlcConnectionLease(container, 1, dummyConn)

	lease.Invalidate()

	err := lease.Ping(context.Background())
	assert.ErrorIs(t, err, errConnectionInvalidated)
}

func TestLeasedPlcConnection_ReadRequestBuilder(t *testing.T) {
	logger := testutils.ProduceTestingLogger(t)
	driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  1 * time.Second,
		maxWaitTime:   5 * time.Second,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	t.Cleanup(func() {
		_ = cache.Close()
	})
	cache.EnableTracer()

	// The first and second connection should work fine
	connection, err := cache.GetConnection(t.Context(), "simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	assert.Nil(t, err)
	if assert.NotNil(t, connection) {
		builder := connection.ReadRequestBuilder()
		assert.NotNil(t, builder)
		assert.NoError(t, connection.Close())
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

func TestLeasedPlcConnection_WriteRequestBuilder(t *testing.T) {
	logger := testutils.ProduceTestingLogger(t)
	driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  1 * time.Second,
		maxWaitTime:   5 * time.Second,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	t.Cleanup(func() {
		_ = cache.Close()
	})
	cache.EnableTracer()

	// The first and second connection should work fine
	connection, err := cache.GetConnection(t.Context(), "simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	assert.Nil(t, err)
	if assert.NotNil(t, connection) {
		builder := connection.WriteRequestBuilder()
		assert.NotNil(t, builder)
		assert.NoError(t, connection.Close())
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

type dummyDriverManager struct {
	factory func() plc4go.PlcConnection
}

func (d *dummyDriverManager) RegisterDriver(plc4go.PlcDriver) {}

func (d *dummyDriverManager) ListDriverNames() []string { return nil }

func (d *dummyDriverManager) GetDriver(string) (plc4go.PlcDriver, error) {
	return nil, errors.New("not implemented")
}

func (d *dummyDriverManager) GetConnection(context.Context, string) (plc4go.PlcConnection, error) {
	if d.factory != nil {
		return d.factory(), nil
	}
	return &dummyTracedConnection{}, nil
}

func (d *dummyDriverManager) Discover(context.Context, func(apiModel.PlcDiscoveryItem), ...plc4go.WithDiscoveryOption) error {
	return nil
}

func (d *dummyDriverManager) Close() error { return nil }

type dummyTracedConnection struct {
	pingCount   int
	invalidated bool
}

func (d *dummyTracedConnection) String() string { return "dummy" }

func (d *dummyTracedConnection) Close() error { return nil }

func (d *dummyTracedConnection) Connect(context.Context) error { return nil }

func (d *dummyTracedConnection) IsConnected() bool { return !d.invalidated }

func (d *dummyTracedConnection) Ping(context.Context) error {
	d.pingCount++
	return nil
}

func (d *dummyTracedConnection) Invalidate() {
	d.invalidated = true
}

func (d *dummyTracedConnection) GetMetadata() apiModel.PlcConnectionMetadata { return nil }

func (d *dummyTracedConnection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder { return nil }

func (d *dummyTracedConnection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder { return nil }

func (d *dummyTracedConnection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	return nil
}

func (d *dummyTracedConnection) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	return nil
}

func (d *dummyTracedConnection) BrowseRequestBuilder() apiModel.PlcBrowseRequestBuilder { return nil }

func (d *dummyTracedConnection) GetConnectionId() string { return "dummy" }

func (d *dummyTracedConnection) IsTraceEnabled() bool { return false }

func (d *dummyTracedConnection) GetTracer() tracer.Tracer { return nil }

func TestLeasedPlcConnection_SubscriptionRequestBuilder(t *testing.T) {
	logger := testutils.ProduceTestingLogger(t)
	driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  1 * time.Second,
		maxWaitTime:   5 * time.Second,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	t.Cleanup(func() {
		_ = cache.Close()
	})
	cache.EnableTracer()

	// The first and second connection should work fine
	connection, err := cache.GetConnection(t.Context(), "simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	assert.Nil(t, err)
	if assert.NotNil(t, connection) {
		builder := connection.SubscriptionRequestBuilder()
		assert.NotNil(t, builder)
		assert.NoError(t, connection.Close())
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

func TestLeasedPlcConnection_UnsubscriptionRequestBuilder(t *testing.T) {
	logger := testutils.ProduceTestingLogger(t)
	driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  1 * time.Second,
		maxWaitTime:   5 * time.Second,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	t.Cleanup(func() {
		_ = cache.Close()
	})
	cache.EnableTracer()

	// The first and second connection should work fine
	connection, err := cache.GetConnection(t.Context(), "simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	assert.Nil(t, err)
	if assert.NotNil(t, connection) {
		func() {
			defer func() {
				if r := recover(); r != nil {
					assert.Equal(t, r, "not provided by simulated connection")
				} else {
					t.Errorf("The code did not panic")
				}
			}()
			connection.UnsubscriptionRequestBuilder()
		}()
		assert.NoError(t, connection.Close())
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

func TestLeasedPlcConnection_BrowseRequestBuilder(t *testing.T) {
	logger := testutils.ProduceTestingLogger(t)
	driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
	// Reduce the max lease time as this way we also reduce the max wait time.
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  1 * time.Second,
		maxWaitTime:   5 * time.Second,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	t.Cleanup(func() {
		_ = cache.Close()
	})
	cache.EnableTracer()

	// The first and second connection should work fine
	connection, err := cache.GetConnection(t.Context(), "simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")
	assert.Nil(t, err)
	if assert.NotNil(t, connection) {
		func() {
			defer func() {
				if r := recover(); r != nil {
					assert.Equal(t, r, "not provided by simulated connection")
				} else {
					t.Errorf("The code did not panic")
				}
			}()
			connection.BrowseRequestBuilder()
		}()
		assert.NoError(t, connection.Close())
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

func Test_plcConnectionLease_String(t1 *testing.T) {
	type fields struct {
		connectionContainer *connectionContainer
		leaseId             uint32
		connection          tracedPlcConnection
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "String it",
			want: "plcConnectionLease{connectionContainer: <nil>, leaseId: 0, connection: %!s(<nil>)}",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &plcConnectionLease{
				connectionContainer: tt.fields.connectionContainer,
				leaseId:             tt.fields.leaseId,
				connection:          tt.fields.connection,
			}
			assert.Equalf(t1, tt.want, t.String(), "String()")
		})
	}
}
