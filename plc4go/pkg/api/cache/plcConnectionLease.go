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
	"github.com/apache/plc4x/plc4go/spi/tracer"
	"time"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	_default "github.com/apache/plc4x/plc4go/spi/default"
)

type plcConnectionLease struct {
	// Reference back to the container, so we can give the connection back.
	connectionContainer *connectionContainer
	// Counter for the number of times this connection has been used before.
	leaseId uint32
	// The actual connection being cached.
	connection tracedPlcConnection
}

func newPlcConnectionLease(connectionContainer *connectionContainer, leaseId uint32, connection tracedPlcConnection) *plcConnectionLease {
	p := &plcConnectionLease{
		connectionContainer: connectionContainer,
		leaseId:             leaseId,
		connection:          connection,
	}
	if connection.IsTraceEnabled() {
		connection.GetTracer().SetConnectionId(p.GetConnectionId())
	}
	return p
}

func (t *plcConnectionLease) IsTraceEnabled() bool {
	if t.connection == nil {
		panic("Called 'IsTraceEnabled' on a closed cached connection")
	}
	return t.connection.IsTraceEnabled()
}

func (t *plcConnectionLease) GetTracer() tracer.Tracer {
	if t.connection == nil {
		panic("Called 'GetTracer' on a closed cached connection")
	}
	return t.connection.GetTracer()
}

func (t *plcConnectionLease) GetConnectionId() string {
	if t.connection == nil {
		panic("Called 'GetConnectionId' on a closed cached connection")
	}
	return fmt.Sprintf("%s-%d", t.connection.GetConnectionId(), t.leaseId)
}

func (t *plcConnectionLease) Connect() <-chan plc4go.PlcConnectionConnectResult {
	panic("Called 'Connect' on a cached connection")
}

func (t *plcConnectionLease) ConnectWithContext(_ context.Context) <-chan plc4go.PlcConnectionConnectResult {
	panic("Called 'Connect' on a cached connection")
}

func (t *plcConnectionLease) BlockingClose() {
	if t.connection == nil {
		panic("Called 'BlockingClose' on a closed cached connection")
	}
	// Call close and wait for the operation to finish.
	<-t.Close()
}

func (t *plcConnectionLease) Close() <-chan plc4go.PlcConnectionCloseResult {
	if t.connection == nil {
		panic("Called 'Close' on a closed cached connection")
	}

	result := make(chan plc4go.PlcConnectionCloseResult, 1)

	go func() {
		// Check if the connection is still alive, if it is, put it back into the cache
		pingResults := t.Ping()
		pingTimeout := time.NewTimer(5 * time.Second)
		newState := StateIdle
		select {
		case pingResult := <-pingResults:
			{
				if pingResult.GetErr() != nil {
					newState = StateInvalid
				}
			}
		case <-pingTimeout.C:
			{
				// Add some trace information
				if t.connection.IsTraceEnabled() {
					t.connection.GetTracer().AddTrace("ping", "timeout")
				}
				// Mark the connection as broken ...
				newState = StateInvalid
			}
		}

		// Extract the trace entries from the connection.
		var traces []tracer.TraceEntry
		if t.IsTraceEnabled() {
			_tracer := t.GetTracer()
			// Save all traces.
			traces = _tracer.GetTraces()
			// Clear the log.
			_tracer.ResetTraces()
			// Reset the connection id back to the one without the lease-id.
			_tracer.SetConnectionId(t.connection.GetConnectionId())
		}

		// Return the connection to the connection container and don't actually close it.
		err := t.connectionContainer.returnConnection(newState)

		// Detach the connection from this lease, so it can no longer be used by the client.
		t.connection = nil

		// Finish closing the connection.
		result <- _default.NewDefaultPlcConnectionCloseResultWithTraces(t, err, traces)
	}()

	return result
}

func (t *plcConnectionLease) IsConnected() bool {
	if t.connection == nil {
		return false
	}
	return t.connection.IsConnected()
}

func (t *plcConnectionLease) Ping() <-chan plc4go.PlcConnectionPingResult {
	if t.connection == nil {
		panic("Called 'Ping' on a closed cached connection")
	}
	return t.connection.Ping()
}

func (t *plcConnectionLease) GetMetadata() apiModel.PlcConnectionMetadata {
	if t.connection == nil {
		panic("Called 'GetMetadata' on a closed cached connection")
	}
	return t.connection.GetMetadata()
}

func (t *plcConnectionLease) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	if t.connection == nil {
		panic("Called 'ReadRequestBuilder' on a closed cached connection")
	}
	return t.connection.ReadRequestBuilder()
}

func (t *plcConnectionLease) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	if t.connection == nil {
		panic("Called 'WriteRequestBuilder' on a closed cached connection")
	}
	return t.connection.WriteRequestBuilder()
}

func (t *plcConnectionLease) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	if t.connection == nil {
		panic("Called 'SubscriptionRequestBuilder' on a closed cached connection")
	}
	return t.connection.SubscriptionRequestBuilder()
}

func (t *plcConnectionLease) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	if t.connection == nil {
		panic("Called 'UnsubscriptionRequestBuilder' on a closed cached connection")
	}
	return t.connection.UnsubscriptionRequestBuilder()
}

func (t *plcConnectionLease) BrowseRequestBuilder() apiModel.PlcBrowseRequestBuilder {
	if t.connection == nil {
		panic("Called 'BrowseRequestBuilder' on a closed cached connection")
	}
	return t.connection.BrowseRequestBuilder()
}

func (t *plcConnectionLease) String() string {
	return fmt.Sprintf("plcConnectionLease{connectionContainer: %s, leaseId: %d, connection: %s}", t.connectionContainer, t.leaseId, t.connection)
}
