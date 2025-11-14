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
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/viney-shih/go-lock"

	"github.com/apache/plc4x/plc4go/internal/simulated"
	"github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/config"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/tracer"
)

var debugTimeout = 100

func TestPlcConnectionCache_GetConnection(t *testing.T) {
	type fields struct {
		driverManager plc4go.PlcDriverManager
	}
	type args struct {
		connectionString string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		setup       func(t *testing.T, fields *fields, args *args)
		wantErr     bool
		wantTimeout bool
	}{
		{
			name: "simple",
			args: args{
				connectionString: "simulated://1.2.3.4:42",
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				logger := testutils.ProduceTestingLogger(t)
				driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
				t.Cleanup(func() {
					assert.NoError(t, driverManager.Close())
				})
				driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
				fields.driverManager = driverManager
			},
			wantErr:     false,
			wantTimeout: false,
		},
		{
			name: "simpleWithTimeout",
			args: args{
				connectionString: "simulated://1.2.3.4:42?connectionDelay=5",
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				logger := testutils.ProduceTestingLogger(t)
				driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
				t.Cleanup(func() {
					assert.NoError(t, driverManager.Close())
				})
				driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
				fields.driverManager = driverManager
			},
			wantErr:     false,
			wantTimeout: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields, &tt.args)
			}
			cc := NewPlcConnectionCache(tt.fields.driverManager, WithCustomLogger(testutils.ProduceTestingLogger(t)))
			conn, err := cc.GetConnection(t.Context(), tt.args.connectionString)
			if tt.wantErr && (err == nil) {
				t.Errorf("PlcConnectionCache.GetConnection() = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil {
				t.Errorf("PlcConnectionCache.GetConnection() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(conn)
		})
	}
}

func TestPlcConnectionCache_Close(t *testing.T) {
	type fields struct {
		driverManager plc4go.PlcDriverManager
	}
	type args struct {
		connectionStrings []string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		setup       func(t *testing.T, fields *fields, args *args)
		wantErr     bool
		wantTimeout bool
	}{
		{
			name: "simple",
			args: args{
				connectionStrings: []string{
					"simulated://1.2.3.4:42",
					"simulated://4.3.2.1:23",
					"simulated://0.8.1.15:7",
				},
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				logger := testutils.ProduceTestingLogger(t)

				driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
				t.Cleanup(func() {
					assert.NoError(t, driverManager.Close())
				})
				driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
				fields.driverManager = driverManager
			},
			wantErr:     false,
			wantTimeout: false,
		},
		{
			name: "empty close",
			setup: func(t *testing.T, fields *fields, args *args) {
				logger := testutils.ProduceTestingLogger(t)

				driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
				t.Cleanup(func() {
					assert.NoError(t, driverManager.Close())
				})
				driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
				fields.driverManager = driverManager
			},
			wantErr:     false,
			wantTimeout: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields, &tt.args)
			}
			cc := NewPlcConnectionCache(tt.fields.driverManager)
			// Connect to all sources first
			for _, connectionString := range tt.args.connectionStrings {
				conn, err := cc.GetConnection(t.Context(), connectionString)
				if err != nil {
					t.Errorf("PlcConnectionCache.GetConnection() error = %v", err)
				} else {
					// Give the connection back.
					if err := conn.Close(); err != nil {
						t.Log(err)
					}
				}
			}
			// Close all connections.
			assert.NoError(t, cc.Close())
		})
	}
}

func readFromPlc(t *testing.T, ctx context.Context, c *plcConnectionCache, preConnectJob func(), connectionString string, resourceString string) <-chan []tracer.TraceEntry {
	t.Helper()
	t.Log("readFromPlc")
	t.Log("Creating tracer channel")
	tracerChan := make(chan []tracer.TraceEntry, 1)

	if preConnectJob != nil {
		t.Log("Executing preConnectJob")
		preConnectJob()
	}
	t.Log("Getting connection from cache")
	// Get a connection
	connection, err := c.GetConnection(ctx, connectionString)
	if err != nil {
		t.Errorf("PlcConnectionCache.GetConnection() error = %v", err)
		return nil
	}
	t.Log("Got connection from cache")
	defer func() {
		t.Log("Closing connection")
		if err := connection.Close(); err != nil {
			t.Log("Error closing connection", err)
		}
		leasedConnection := connection.(*plcConnectionLease)
		tracerChan <- leasedConnection.GetLastTraces()
		t.Log("Closed connection")
	}()

	// Prepare a read request.
	t.Log("Preparing read request")
	readRequest, err := connection.ReadRequestBuilder().AddTagAddress("test", resourceString).Build()
	if err != nil {
		t.Errorf("PlcConnectionCache.ReadRequest.Build() error = %v", err)
		return tracerChan
	}

	// Execute the read request.
	t.Log("Executing read request")
	execution := readRequest.Execute(ctx)
	select {
	case readRequestResult := <-execution:
		err := readRequestResult.GetErr()
		if err != nil {
			t.Errorf("PlcConnectionCache.ReadRequest.Read() error = %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Errorf("PlcConnectionCache.ReadRequest.Read() timeout")
	}
	return tracerChan
}

func executeAndTestReadFromPlc(t *testing.T, ctx context.Context, c *plcConnectionCache, preConnectJob func(), connectionString string, resourceString string, expectedTraceEntries []string, expectedNumTotalConnections int) <-chan struct{} {
	t.Helper()
	ch := make(chan struct{}, 1)
	c.wg.Go(func() {
		t.Log("Starting goroutine")
		// Read once from the c.
		t.Log("Reading from the cache")
		var traces []tracer.TraceEntry
		select {
		case traces = <-readFromPlc(t, ctx, c, preConnectJob, connectionString, resourceString):
		case <-ctx.Done():
			t.Log("Context done", t.Context().Err())
			ch <- struct{}{}
			return
		case <-t.Context().Done():
			t.Log("Context done", t.Context().Err())
			ch <- struct{}{}
			return
		}

		t.Log("Finished reading from the cache")

		t.Log("Checking trace entries")
		t.Log("actual, expected")
		for i := 0; i < max(len(traces), len(expectedTraceEntries)); i++ {
			actual, expected := "\t\t\t\t\t", "\t\t\t\t\t"
			if i < len(traces)-1 {
				actual = fmt.Sprintf("%s-%s", traces[i].Operation, traces[i].Message)
			}
			if i < len(expectedTraceEntries)-1 {
				expected = fmt.Sprintf("%v", expectedTraceEntries[i])
			}
			t.Logf("%s\t\t%s", actual, expected)
			assert.Equal(t, actual, expected)
		}
		t.Log("Trace entries are as expected")
		t.Log("Checking number of connections in the cache")
		// Now there should be one connection in the c.
		if len(c.connections) != expectedNumTotalConnections {
			t.Errorf("Expected %d connections in the c but got %d", expectedNumTotalConnections, len(c.connections))
		}
		t.Log("Number of connections in the cache is as expected")
		ch <- struct{}{}
	})
	return ch
}

func TestPlcConnectionCache_ReusingAnExistingConnection(t *testing.T) {
	logger := testutils.ProduceTestingLogger(t)
	driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
	cache := &plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  5 * time.Second,
		maxWaitTime:   25 * time.Second,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	cache.EnableTracer()

	// Initially there should be no connection in the cache.
	if len(cache.connections) != 0 {
		t.Errorf("Expected %d connections in the cache but got %d", 0, len(cache.connections))
	}

	// Read once from the cache.
	finishedChan := executeAndTestReadFromPlc(
		t,
		t.Context(),
		cache,
		nil,
		"simulated://1.2.3.4:42?traceEnabled=true",
		"RANDOM/test_random:BOOL",
		[]string{
			"connect-started",
			"connect-success",
			"read-started",
			"read-success",
			"ping-started",
			"ping-success",
		},
		1,
	)
	select {
	case _ = <-finishedChan:
	case <-time.After(500 * time.Millisecond * time.Duration(debugTimeout)):
		t.Errorf("Timeout")
	}

	// Request the same connection for a second time.
	finishedChan = executeAndTestReadFromPlc(
		t,
		t.Context(),
		cache,
		nil,
		"simulated://1.2.3.4:42?traceEnabled=true",
		"RANDOM/test_random:BOOL",
		[]string{
			"read-started",
			"read-success",
			"ping-started",
			"ping-success",
		},
		1,
	)
	select {
	case _ = <-finishedChan:
	case <-time.After(500 * time.Millisecond * time.Duration(debugTimeout)):
		t.Errorf("Timeout")
	}

	assert.NotNil(t, cache.GetTracer(), "Tracer should be available")
	traces := cache.GetTracer().GetTraces()
	require.Equal(t, 5, len(traces), "Unexpected number of trace entries")
	// First is needs to create a new container for this connection
	assert.Equal(t, "create new cached connection", traces[0].Message, "Unexpected message")
	// Then it gets a lease for the connection
	assert.Equal(t, "lease", traces[1].Message, "Unexpected message")
	assert.Equal(t, "success", traces[2].Message, "Unexpected message")
	// And a second time
	assert.Equal(t, "lease", traces[3].Message, "Unexpected message")
	assert.Equal(t, "success", traces[4].Message, "Unexpected message")
}

func TestPlcConnectionCache_MultipleConcurrentConnectionRequests(t *testing.T) {
	logger := testutils.ProduceTestingLogger(t)
	driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
	cache := &plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  5 * time.Second,
		maxWaitTime:   25 * time.Second,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	cache.EnableTracer()

	// Initially there should be no connection in the cache.
	if len(cache.connections) != 0 {
		t.Errorf("Expected %d connections in the cache but got %d", 0, len(cache.connections))
	}

	floodGate := lock.NewCASMutex() // floodgate is use because we want both get connection to get executed in short order
	floodGate.Lock()                // We use a cas mutex write lock to lock the floodgate

	// Read once from the cache.
	firstRun := executeAndTestReadFromPlc(
		t,
		t.Context(),
		cache,
		func() {
			floodGate.RLock()
			defer floodGate.RUnlock()
		},
		"simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true",
		"RANDOM/test_random:BOOL",
		[]string{
			"connect-started",
			"connect-success",
			"read-started",
			"read-success",
			"ping-started",
			"ping-success",
		},
		1,
	)

	// Almost instantly request the same connection for a second time.
	// As the connection takes 100ms, the second connection request will come
	// in while the first is still not finished. So in theory it would have
	// to wait for the first operation to be finished first.
	secondRun := executeAndTestReadFromPlc(
		t,
		t.Context(),
		cache,
		func() {
			floodGate.RLock()
			defer floodGate.RUnlock()
			time.Sleep(10 * time.Millisecond)
		},
		"simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true",
		"RANDOM/test_random:BOOL",
		[]string{
			"read-started",
			"read-success",
			"ping-started",
			"ping-success",
		},
		1,
	)
	floodGate.Unlock()
	select {
	case _ = <-firstRun:
		select {
		case _ = <-secondRun:
		case <-time.After(500 * time.Millisecond * time.Duration(debugTimeout)):
			t.Errorf("Timeout")
		}
		break
	case <-time.After(1 * time.Second * time.Duration(debugTimeout)):
		t.Errorf("Timeout")
	}

	// This should be quite equal to the serial case as the connections are requested serially.
	assert.NotNil(t, cache.GetTracer(), "Tracer should be available")
	traces := cache.GetTracer().GetTraces()
	require.Equal(t, 5, len(traces), "Unexpected number of trace entries")
	// First is needs to create a new container for this connection
	assert.Equal(t, "create new cached connection", traces[0].Message, "Unexpected message")
	// Then it gets a lease for the connection
	assert.Equal(t, "lease", traces[1].Message, "Unexpected message")
	// And a second time
	assert.Equal(t, "lease", traces[2].Message, "Unexpected message")
	// Now the delay of 100ms is over, and we should see the first success
	assert.Equal(t, "success", traces[3].Message, "Unexpected message")
	// Now the first operation is finished, and we should see the second success
	assert.Equal(t, "success", traces[4].Message, "Unexpected message")
}

func TestPlcConnectionCache_ConnectWithError(t *testing.T) {
	logger := testutils.ProduceTestingLogger(t)
	driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  5 * time.Second,
		maxWaitTime:   25 * time.Second,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	cache.EnableTracer()

	// Initially there should be no connection in the cache.
	if len(cache.connections) != 0 {
		t.Errorf("Expected %d connections in the cache but got %d", 0, len(cache.connections))
	}

	conn, err := cache.GetConnection(t.Context(), "simulated://1.2.3.4:42?connectionError=hurz&traceEnabled=true")
	if err == nil {
		t.Error("An error was expected")
		return
	}
	if !strings.Contains(err.Error(), "hurz") {
		t.Errorf("An error '%s' was expected, but got '%s'", "hurz", err)
	}
	t.Log(conn)
}

// In this test, the ping operation used to test the connection before
// putting it back into the cache will return an error, hereby marking
// the connection as invalid
func TestPlcConnectionCache_ReturningConnectionWithPingError(t *testing.T) {
	logger := testutils.ProduceTestingLogger(t)
	driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  5 * time.Second,
		maxWaitTime:   25 * time.Second,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	cache.EnableTracer()

	// Initially there should be no connection in the cache.
	if len(cache.connections) != 0 {
		t.Errorf("Expected %d connections in the cache but got %d", 0, len(cache.connections))
	}

	// In the connection string, we tell the driver to return an error with
	// the given message on executing a ping operation.
	conn, err := cache.GetConnection(t.Context(), "simulated://1.2.3.4:42?pingError=hurz&traceEnabled=true")
	require.NoError(t, err, "Error getting connection from cache: %s", err)
	connection := conn.(*plcConnectionLease)
	if err := connection.Close(); err != nil {
		traces := connection.GetLastTraces()
		// We expect 4 traces (Connect start & success and Ping start and error.
		require.Len(t, traces, 4, "Expected %d trace entries but got %d", 4, len(traces))
		if traces[0].Operation+"-"+traces[0].Message != "connect-started" {
			t.Errorf("Expected '%s' as first trace message, but got '%s'", "connect-started", traces[0])
		}
		if traces[1].Operation+"-"+traces[1].Message != "connect-success" {
			t.Errorf("Expected '%s' as second trace message, but got '%s'", "connect-success", traces[1])
		}
		if traces[2].Operation+"-"+traces[2].Message != "ping-started" {
			t.Errorf("Expected '%s' as third trace message, but got '%s'", "ping-started", traces[2])
		}
		if traces[3].Operation+"-"+traces[3].Message != "ping-error: hurz" {
			t.Errorf("Expected '%s' as fourth trace message, but got '%s'", "ping-error: hurz", traces[3])
		}
	}
}

// In this test, we'll make the ping operation take longer than the timeout in the connection cache
// Therefore the error handling should kick in.
func TestPlcConnectionCache_PingTimeout(t *testing.T) {
	logger := testutils.ProduceTestingLogger(t)
	driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
	cache := &plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  5 * time.Second,
		maxWaitTime:   25 * time.Second,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	cache.EnableTracer()

	// Initially there should be no connection in the cache.
	if len(cache.connections) != 0 {
		t.Errorf("Expected %d connections in the cache but got %d", 0, len(cache.connections))
	}

	ctx, cancelFunc := context.WithTimeout(t.Context(), 100*time.Millisecond)
	t.Cleanup(cancelFunc)

	// Read once from the cache.
	firstRun := executeAndTestReadFromPlc(
		t,
		ctx,
		cache,
		nil,
		"simulated://1.2.3.4:42?pingDelay=10000&traceEnabled=true",
		"RANDOM/test_random:BOOL",
		[]string{
			"connect-started",
			"connect-success",
			"read-started",
			"read-success",
			"ping-started",
			"ping-timeout",
		}, 1)

	select {
	case _ = <-firstRun:
		break
	case <-time.After(20 * time.Second * time.Duration(debugTimeout)):
		t.Errorf("Timeout")
	}
	t.Log("done")
}

// In this test there are multiple requests for the same connection but the first operation fails at returning
// the connection due to a timeout in the ping operation. The second call should get a new connection in this
// case.
func TestPlcConnectionCache_SecondCallGetNewConnectionAfterPingTimeout(t *testing.T) {
	logger := testutils.ProduceTestingLogger(t)
	driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
	cache := &plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  5 * time.Second,
		maxWaitTime:   25 * time.Second,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	cache.EnableTracer()

	// Initially there should be no connection in the cache.
	if len(cache.connections) != 0 {
		t.Errorf("Expected %d connections in the cache but got %d", 0, len(cache.connections))
	}

	ctx, cancelFunc := context.WithTimeout(t.Context(), 10*time.Second)
	t.Cleanup(cancelFunc)

	// Read once from the cache.
	firstRun := executeAndTestReadFromPlc(
		t,
		ctx,
		cache,
		nil,
		"simulated://1.2.3.4:42?pingDelay=10000&connectionDelay=100&traceEnabled=true",
		"RANDOM/test_random:BOOL",
		[]string{
			"connect-started",
			"connect-success",
			"read-started",
			"read-success",
			"ping-started",
			"ping-timeout",
		}, 1)

	time.Sleep(time.Millisecond * 1)

	ctx, cancelFunc = context.WithTimeout(t.Context(), 10*time.Second)
	t.Cleanup(cancelFunc)

	// Almost instantly request the same connection for a second time.
	// As the connection takes 100ms, the second connection request will come
	// in while the first is still not finished. So in theory it would have
	// to wait for the first operation to be finished first.
	secondRun := executeAndTestReadFromPlc(
		t,
		ctx,
		cache,
		nil,
		"simulated://1.2.3.4:42?pingDelay=10000&connectionDelay=100&traceEnabled=true",
		"RANDOM/test_random:BOOL",
		[]string{
			"connect-started",
			"connect-success",
			"read-started",
			"read-success",
			"ping-started",
			"ping-timeout",
		},
		1,
	)
	select {
	case _ = <-firstRun:
		select {
		case _ = <-secondRun:
		case <-time.After(20 * time.Second * time.Duration(debugTimeout)):
			t.Errorf("Timeout")
		}
	case <-time.After(30 * time.Second * time.Duration(debugTimeout)):
		t.Errorf("Timeout")
	}

	// This should be quite equal to the serial case as the connections are requested serially.
	require.NotNil(t, cache.GetTracer(), "Tracer should be available")
	traces := cache.GetTracer().GetTraces()
	require.Equal(t, 5, len(traces), "Unexpected number of trace entries")
	// First is needs to create a new container for this connection
	assert.Equal(t, "create new cached connection", traces[0].Message, "Unexpected message")
	// Then it gets a lease for the connection
	assert.Equal(t, "lease", traces[1].Message, "Unexpected message")
	// And a second time
	assert.Equal(t, "lease", traces[2].Message, "Unexpected message")
	// Now the delay of 100ms is over, and we should see the first success
	assert.Equal(t, "success", traces[3].Message, "Unexpected message")
	// Now the first operation is finished, and we should see the second success
	assert.Equal(t, "success", traces[4].Message, "Unexpected message")
}

func TestPlcConnectionCache_MaximumWaitTimeReached(t *testing.T) {
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
	cache.EnableTracer()

	// Initially there should be no connection in the cache.
	if len(cache.connections) != 0 {
		t.Errorf("Expected %d connections in the cache but got %d", 0, len(cache.connections))
	}

	// The first and second connection should work fine
	firstConn, err := cache.GetConnection(t.Context(), "simulated://1.2.3.4:42?connectionDelay=100&pingDelay=4000&traceEnabled=true")
	require.NoError(t, err)

	time.Sleep(1 * time.Millisecond)

	// Just make sure the first two connections are returned as soon as they are received
	_ = firstConn.Close()

	// Second one blocks
	secondConn, err := cache.GetConnection(t.Context(), "simulated://1.2.3.4:42?connectionDelay=100&pingDelay=4000&traceEnabled=true")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = secondConn.Close()
	})

	time.Sleep(1 * time.Millisecond)

	// The third connection should be given up by the cache
	thrirdConn, err := cache.GetConnection(t.Context(), "simulated://1.2.3.4:42?connectionDelay=100&pingDelay=4000&traceEnabled=true")
	require.Error(t, err)
	require.Nil(t, thrirdConn)
}
