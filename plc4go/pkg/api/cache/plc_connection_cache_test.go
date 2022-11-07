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
	"testing"
	"time"

	"github.com/apache/plc4x/plc4go/internal/simulated"
	"github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/stretchr/testify/assert"
	"github.com/viney-shih/go-lock"
)

var debugTimeout = 1

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
		wantErr     bool
		wantTimeout bool
	}{
		{
			name: "simple",
			fields: fields{
				driverManager: func() plc4go.PlcDriverManager {
					driverManager := plc4go.NewPlcDriverManager()
					driverManager.RegisterDriver(simulated.NewDriver())
					return driverManager
				}(),
			}, args: args{
				connectionString: "simulated://1.2.3.4:42",
			},
			wantErr:     false,
			wantTimeout: false,
		},
		{
			name: "simpleWithTimeout",
			fields: fields{
				driverManager: func() plc4go.PlcDriverManager {
					driverManager := plc4go.NewPlcDriverManager()
					driverManager.RegisterDriver(simulated.NewDriver())
					return driverManager
				}(),
			}, args: args{
				connectionString: "simulated://1.2.3.4:42?connectionDelay=5",
			},
			wantErr:     false,
			wantTimeout: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := NewPlcConnectionCache(tt.fields.driverManager)
			got := cc.GetConnection(tt.args.connectionString)
			select {
			case connectResult := <-got:
				if tt.wantErr && (connectResult.GetErr() == nil) {
					t.Errorf("PlcConnectionCache.GetConnection() = %v, wantErr %v", connectResult.GetErr(), tt.wantErr)
				} else if connectResult.GetErr() != nil {
					t.Errorf("PlcConnectionCache.GetConnection() error = %v, wantErr %v", connectResult.GetErr(), tt.wantErr)
				}
			case <-time.After(10 * time.Second):
				if !tt.wantTimeout {
					t.Errorf("PlcConnectionCache.GetConnection() got timeout")
				}
			}
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
		wantErr     bool
		wantTimeout bool
	}{
		{
			name: "simple",
			fields: fields{
				driverManager: func() plc4go.PlcDriverManager {
					driverManager := plc4go.NewPlcDriverManager()
					driverManager.RegisterDriver(simulated.NewDriver())
					return driverManager
				}(),
			},
			args: args{
				connectionStrings: []string{
					"simulated://1.2.3.4:42",
					"simulated://4.3.2.1:23",
					"simulated://0.8.1.15:7",
				},
			},
			wantErr:     false,
			wantTimeout: false,
		},
		{
			name: "empty close",
			fields: fields{
				driverManager: func() plc4go.PlcDriverManager {
					driverManager := plc4go.NewPlcDriverManager()
					driverManager.RegisterDriver(simulated.NewDriver())
					return driverManager
				}(),
			},
			wantErr:     false,
			wantTimeout: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := NewPlcConnectionCache(tt.fields.driverManager)
			// Connect to all sources first
			for _, connectionString := range tt.args.connectionStrings {
				got := cc.GetConnection(connectionString)
				select {
				case connectResult := <-got:
					if connectResult.GetErr() != nil {
						t.Errorf("PlcConnectionCache.GetConnection() error = %v", connectResult.GetErr())
					} else {
						// Give the connection back.
						connectResult.GetConnection().Close()
					}
				case <-time.After(10 * time.Second):
					if !tt.wantTimeout {
						t.Errorf("PlcConnectionCache.GetConnection() got timeout")
					}
				}
			}
			// Close all connections.
			cacheCloseResults := cc.Close()
			// Wait for all connections to be closed.
			select {
			case cacheCloseResult := <-cacheCloseResults:
				if tt.wantErr && (cacheCloseResult.GetErr() == nil) {
					t.Errorf("PlcConnectionCache.Close() = %v, wantErr %v", cacheCloseResult.GetErr(), tt.wantErr)
				} else if cacheCloseResult.GetErr() != nil {
					t.Errorf("PlcConnectionCache.Close() error = %v, wantErr %v", cacheCloseResult.GetErr(), tt.wantErr)
				}
			case <-time.After(10 * time.Second):
				if !tt.wantTimeout {
					t.Errorf("PlcConnectionCache.Close() got timeout")
				}
			}

		})
	}
}

func readFromPlc(t *testing.T, cache plcConnectionCache, connectionString string, resourceString string) <-chan []spi.TraceEntry {
	ch := make(chan []spi.TraceEntry)

	// Get a connection
	connectionResultChan := cache.GetConnection(connectionString)
	select {
	case connectResult := <-connectionResultChan:
		if connectResult.GetErr() != nil {
			t.Errorf("PlcConnectionCache.GetConnection() error = %v", connectResult.GetErr())
			return nil
		}
		connection := connectResult.GetConnection()
		defer func() {
			closeResults := connection.Close()
			// Wait for the connection to be correctly closed.
			closeResult := <-closeResults
			go func() {
				ch <- (closeResult.(_default.DefaultPlcConnectionCloseResult)).GetTraces()
			}()
		}()

		// Prepare a read request.
		readRequest, err := connection.ReadRequestBuilder().AddFieldQuery("test", resourceString).Build()
		if err != nil {
			t.Errorf("PlcConnectionCache.ReadRequest.Build() error = %v", err)
			return ch
		}

		// Execute the read request.
		execution := readRequest.Execute()
		select {
		case readRequestResult := <-execution:
			err := readRequestResult.GetErr()
			if err != nil {
				t.Errorf("PlcConnectionCache.ReadRequest.Read() error = %v", err)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("PlcConnectionCache.ReadRequest.Read() timeout")
		}
		return ch
	case <-time.After(20 * time.Second):
		t.Errorf("PlcConnectionCache.GetConnection() got timeout")
	}
	return ch
}

func executeAndTestReadFromPlc(t *testing.T, cache plcConnectionCache, connectionString string, resourceString string, expectedTraceEntries []string, expectedNumTotalConnections int) <-chan bool {
	ch := make(chan bool)
	go func() {
		// Read once from the cache.
		tracesChannel := readFromPlc(t, cache, connectionString, resourceString)
		traces := <-tracesChannel

		// In the log we should see one "Successfully connected" entry.
		if len(traces) != len(expectedTraceEntries) {
			t.Errorf("Expected %d 'Successfully connected' entries in the log but got %d", len(expectedTraceEntries), len(traces))
		}
		for i, expectedTraceEntry := range expectedTraceEntries {
			currentTraceEntry := traces[i].Operation + "-" + traces[i].Message
			if expectedTraceEntry != currentTraceEntry {
				t.Errorf("Expected %s as trace entry but got %s", expectedTraceEntry, currentTraceEntry)
			}
		}
		// Now there should be one connection in the cache.
		if len(cache.connections) != expectedNumTotalConnections {
			t.Errorf("Expected %d connections in the cache but got %d", expectedNumTotalConnections, len(cache.connections))
		}
		ch <- true
	}()
	return ch
}

func TestPlcConnectionCache_ReusingAnExistingConnection(t *testing.T) {
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(simulated.NewDriver())
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  time.Second * 5,
		maxWaitTime:   time.Second * 25,
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
	finishedChan := executeAndTestReadFromPlc(t, cache, "simulated://1.2.3.4:42?traceEnabled=true", "RANDOM/test_random:BOOL",
		[]string{
			"connect-started",
			"connect-success",
			"read-started",
			"read-success",
			"ping-started",
			"ping-success",
		}, 1)
	select {
	case _ = <-finishedChan:
	case <-time.After(500 * time.Millisecond * time.Duration(debugTimeout)):
		t.Errorf("Timeout")
	}

	// Request the same connection for a second time.
	finishedChan = executeAndTestReadFromPlc(t, cache, "simulated://1.2.3.4:42?traceEnabled=true", "RANDOM/test_random:BOOL",
		[]string{
			"read-started",
			"read-success",
			"ping-started",
			"ping-success",
		}, 1)
	select {
	case _ = <-finishedChan:
	case <-time.After(500 * time.Millisecond * time.Duration(debugTimeout)):
		t.Errorf("Timeout")
	}

	assert.NotNil(t, cache.GetTracer(), "Tracer should be available")
	traces := cache.GetTracer().GetTraces()
	assert.Equal(t, 5, len(traces), "Unexpected number of trace entries")
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
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(simulated.NewDriver())
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  time.Second * 5,
		maxWaitTime:   time.Second * 25,
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
	firstRun := executeAndTestReadFromPlc(t, cache, "simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true", "RANDOM/test_random:BOOL",
		[]string{
			"connect-started",
			"connect-success",
			"read-started",
			"read-success",
			"ping-started",
			"ping-success",
		}, 1)

	time.Sleep(time.Millisecond * 1)

	// Almost instantly request the same connection for a second time.
	// As the connection takes 100ms, the second connection request will come
	// in while the first is still not finished. So in theory it would have
	// to wait for the first operation to be finished first.
	secondRun := executeAndTestReadFromPlc(t, cache, "simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true", "RANDOM/test_random:BOOL",
		[]string{
			"read-started",
			"read-success",
			"ping-started",
			"ping-success",
		}, 1)
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
	assert.Equal(t, 5, len(traces), "Unexpected number of trace entries")
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
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(simulated.NewDriver())
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  time.Second * 5,
		maxWaitTime:   time.Second * 25,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	cache.EnableTracer()

	// Initially there should be no connection in the cache.
	if len(cache.connections) != 0 {
		t.Errorf("Expected %d connections in the cache but got %d", 0, len(cache.connections))
	}

	connectionResultChan := cache.GetConnection("simulated://1.2.3.4:42?connectionError=hurz&traceEnabled=true")
	select {
	case connectResult := <-connectionResultChan:
		if connectResult.GetErr() == nil {
			t.Error("An error was expected")
			return
		}
		if connectResult.GetErr().Error() != "hurz" {
			t.Errorf("An error '%s' was expected, but got '%s'", "hurz", connectResult.GetErr().Error())
		}
	case <-time.After(20 * time.Second):
		t.Errorf("PlcConnectionCache.GetConnection() got timeout")
	}
}

// In this test, the ping operation used to test the connection before
// putting it back into the cache will return an error, hereby marking
// the connection as invalid
func TestPlcConnectionCache_ReturningConnectionWithPingError(t *testing.T) {
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(simulated.NewDriver())
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  time.Second * 5,
		maxWaitTime:   time.Second * 25,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	cache.EnableTracer()

	// Initially there should be no connection in the cache.
	if len(cache.connections) != 0 {
		t.Errorf("Expected %d connections in the cache but got %d", 0, len(cache.connections))
	}

	connectionResultChan := cache.GetConnection("simulated://1.2.3.4:42?pingError=hurz&traceEnabled=true")
	select {
	case connectResult := <-connectionResultChan:
		if connectResult.GetErr() != nil {
			t.Errorf("PlcConnectionCache.GetConnection() error = %v", connectResult.GetErr())
		}
		connection := connectResult.GetConnection().(*plcConnectionLease)
		if connection != nil {
			connectionCloseResultChan := connection.Close()
			closeResult := <-connectionCloseResultChan
			if closeResult != nil {
				traces := (closeResult.(_default.DefaultPlcConnectionCloseResult)).GetTraces()
				// We expect 4 traces (Connect start & success and Ping start and error.
				if len(traces) != 4 {
					t.Errorf("Expected %d trace entries but got %d", 4, len(traces))
				}
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
	case <-time.After(20 * time.Second):
		t.Errorf("PlcConnectionCache.GetConnection() got timeout")
	}
}

// In this test, we'll make the ping operation take longer than the timeout in the connection cache
// Therefore the error handling should kick in.
func TestPlcConnectionCache_PingTimeout(t *testing.T) {
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(simulated.NewDriver())
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  time.Second * 5,
		maxWaitTime:   time.Second * 25,
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
	firstRun := executeAndTestReadFromPlc(t, cache, "simulated://1.2.3.4:42?pingDelay=10000&traceEnabled=true", "RANDOM/test_random:BOOL",
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

}

// In this test there are multiple requests for the same connection but the first operation fails at returning
// the connection due to a timeout in the ping operation. The second call should get a new connection in this
// case.
func TestPlcConnectionCache_SecondCallGetNewConnectionAfterPingTimeout(t *testing.T) {
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(simulated.NewDriver())
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  time.Second * 5,
		maxWaitTime:   time.Second * 25,
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
	firstRun := executeAndTestReadFromPlc(t, cache, "simulated://1.2.3.4:42?pingDelay=10000&connectionDelay=100&traceEnabled=true", "RANDOM/test_random:BOOL",
		[]string{
			"connect-started",
			"connect-success",
			"read-started",
			"read-success",
			"ping-started",
			"ping-timeout",
		}, 1)

	time.Sleep(time.Millisecond * 1)

	// Almost instantly request the same connection for a second time.
	// As the connection takes 100ms, the second connection request will come
	// in while the first is still not finished. So in theory it would have
	// to wait for the first operation to be finished first.
	secondRun := executeAndTestReadFromPlc(t, cache, "simulated://1.2.3.4:42?pingDelay=10000&connectionDelay=100&traceEnabled=true", "RANDOM/test_random:BOOL",
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
		select {
		case _ = <-secondRun:
		case <-time.After(20 * time.Second * time.Duration(debugTimeout)):
			t.Errorf("Timeout")
		}
		break
	case <-time.After(30 * time.Second * time.Duration(debugTimeout)):
		t.Errorf("Timeout")
	}

	// This should be quite equal to the serial case as the connections are requested serially.
	assert.NotNil(t, cache.GetTracer(), "Tracer should be available")
	traces := cache.GetTracer().GetTraces()
	assert.Equal(t, 5, len(traces), "Unexpected number of trace entries")
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

// In this test the first client requests a connection, but doesn't listen on the response-channel
// This shouldn't block the connection cache.
func TestPlcConnectionCache_FistReadGivesUpBeforeItGetsTheConnectionSoSecondOneTakesOver(t *testing.T) {
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(simulated.NewDriver())
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  time.Second * 5,
		maxWaitTime:   time.Second * 25,
		cacheLock:     lock.NewCASMutex(),
		connections:   make(map[string]*connectionContainer),
		tracer:        nil,
	}
	cache.EnableTracer()

	// Initially there should be no connection in the cache.
	if len(cache.connections) != 0 {
		t.Errorf("Expected %d connections in the cache but got %d", 0, len(cache.connections))
	}

	// Intentionally just ignore the response.
	cache.GetConnection("simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")

	time.Sleep(time.Millisecond * 1)

	// Read once from the cache.
	// NOTE: It doesn't contain the connect-part, as the previous connection handled that.
	firstRun := executeAndTestReadFromPlc(t, cache, "simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true", "RANDOM/test_random:BOOL",
		[]string{
			"read-started",
			"read-success",
			"ping-started",
			"ping-success",
		}, 1)

	select {
	case _ = <-firstRun:
		break
	case <-time.After(30 * time.Second * time.Duration(debugTimeout)):
		t.Errorf("Timeout")
	}
}

func TestPlcConnectionCache_SecondConnectionGivenUpWaiting(t *testing.T) {
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(simulated.NewDriver())
	cache := plcConnectionCache{
		driverManager: driverManager,
		maxLeaseTime:  time.Second * 5,
		maxWaitTime:   time.Second * 25,
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
	firstRun := executeAndTestReadFromPlc(t, cache, "simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true", "RANDOM/test_random:BOOL",
		[]string{
			"connect-started",
			"connect-success",
			"read-started",
			"read-success",
			"ping-started",
			"ping-success",
		}, 1)

	time.Sleep(time.Millisecond * 1)

	// Almost instantly we try to get a new connection but don't listen for the result
	cache.GetConnection("simulated://1.2.3.4:42?connectionDelay=100&traceEnabled=true")

	// Wait for the first operation to finish
	select {
	case _ = <-firstRun:
	case <-time.After(30 * time.Second * time.Duration(debugTimeout)):
		t.Errorf("Timeout")
	}

	// Wait for 1s to have the connection cache timeout (10ms) the lease as nobody's listening.
	time.Sleep(time.Millisecond * 1000)

	// This should be quite equal to the serial case as the connections are requested serially.
	assert.NotNil(t, cache.GetTracer(), "Tracer should be available")
	traces := cache.GetTracer().GetTraces()
	if assert.Equal(t, 5, len(traces), "Unexpected number of trace entries") {
		// First is needs to create a new container for this connection
		assert.Equal(t, "create new cached connection", traces[0].Message, "Unexpected message")
		// Then it gets a lease for the connection
		assert.Equal(t, "lease", traces[1].Message, "Unexpected message")
		// And a second time
		assert.Equal(t, "lease", traces[2].Message, "Unexpected message")
		// Now the delay of 100ms is over, and we should see the first success
		assert.Equal(t, "success", traces[3].Message, "Unexpected message")
		// Now the first operation is finished, and we should see the second give up
		assert.Equal(t, "client given up", traces[4].Message, "Unexpected message")
	} else if len(traces) > 0 {
		var values string
		for _, traceEntry := range traces {
			values = values + traceEntry.Operation + "-" + traceEntry.Message + ", "
		}
		t.Errorf("Got traces: %s", values)
	} else {
		t.Error("No traces")
	}
}

func TestPlcConnectionCache_MaximumWaitTimeReached(t *testing.T) {
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

	// Initially there should be no connection in the cache.
	if len(cache.connections) != 0 {
		t.Errorf("Expected %d connections in the cache but got %d", 0, len(cache.connections))
	}

	// The first and second connection should work fine
	firstConnectionResults := cache.GetConnection("simulated://1.2.3.4:42?connectionDelay=100&pingDelay=4000&traceEnabled=true")

	time.Sleep(time.Millisecond * 1)

	secondConnectionResults := cache.GetConnection("simulated://1.2.3.4:42?connectionDelay=100&pingDelay=4000&traceEnabled=true")

	time.Sleep(time.Millisecond * 1)

	// The third connection should be given up by the cache
	thirdConnectionResults := cache.GetConnection("simulated://1.2.3.4:42?connectionDelay=100&pingDelay=4000&traceEnabled=true")

	// Just make sure the first two connections are returned as soon as they are received
	go func() {
		select {
		case connectionResult := <-firstConnectionResults:
			if assert.NotNil(t, connectionResult) {
				if assert.Nil(t, connectionResult.GetErr()) {
					// Give back the connection.
					connectionResult.GetConnection().Close()
				}
			}
		case <-time.After(5 * time.Second):
			t.Errorf("Timeout")
		}
	}()
	go func() {
		select {
		case connectionResult := <-secondConnectionResults:
			if assert.NotNil(t, connectionResult) {
				if assert.Nil(t, connectionResult.GetErr()) {
					// Give back the connection.
					connectionResult.GetConnection().Close()
				}
			}
		case <-time.After(5 * time.Second):
			t.Errorf("Timeout")
		}
	}()

	// Now wait for the last connection to be timed out by the cache
	select {
	case connectionResult := <-thirdConnectionResults:
		if assert.NotNil(t, connectionResult) {
			assert.Nil(t, connectionResult.GetConnection())
			if assert.NotNil(t, connectionResult.GetErr()) {
				assert.Equal(t, "timeout while waiting for connection", connectionResult.GetErr().Error())
			}
		}
	case <-time.After(15 * time.Second):
		t.Errorf("Timeout")
	}
}
