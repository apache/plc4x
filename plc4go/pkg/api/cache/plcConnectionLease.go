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
	"sync/atomic"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/tracer"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type plcConnectionLease struct {
	// Reference back to the container, so we can give the connection back.
	connectionContainer *connectionContainer
	// Counter for the number of times this connection has been used before.
	leaseId uint32
	// The actual connection being cached.
	connection tracedPlcConnection
	// the last traces of this connection
	lastTraces []tracer.TraceEntry
	// invalidated indicates the lease was explicitly marked unusable by the caller.
	invalidated atomic.Bool
}

var errConnectionInvalidated = errors.New("connection has been invalidated")

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
	t.ensureConnection("IsTraceEnabled")
	return t.connection.IsTraceEnabled()
}

func (t *plcConnectionLease) GetTracer() tracer.Tracer {
	t.ensureConnection("GetTracer")
	return t.connection.GetTracer()
}

func (t *plcConnectionLease) GetConnectionId() string {
	t.ensureConnection("GetConnectionId")
	return fmt.Sprintf("%s-%d", t.connection.GetConnectionId(), t.leaseId)
}

func (t *plcConnectionLease) Connect(_ context.Context) error {
	panic("Called 'Connect' on a cached connection")
}

func (t *plcConnectionLease) Close() error {
	ctx := context.TODO()
	ctx, cancelFunc := utils.WithNamedTimeout(ctx, "connection close timeout", 5*time.Second)
	defer cancelFunc()

	if t.connection == nil {
		panic("Called 'Close' on a closed cached connection")
	}

	// Check if the connection is still alive, if it is, put it back into the cache
	newState := StateIdle
	if t.isInvalidated() {
		newState = StateInvalid
	} else if err := t.Ping(ctx); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			// Add some trace information
			if t.connection.IsTraceEnabled() {
				t.connection.GetTracer().AddTrace("ping", "timeout")
			}
		}
		newState = StateInvalid
	}

	// Extract the trace entries from the connection unless it was invalidated.
	if !t.isInvalidated() && t.IsTraceEnabled() {
		_tracer := t.GetTracer()
		// Save all traces.
		t.lastTraces = _tracer.GetTraces()
		// Clear the log.
		_tracer.ResetTraces()
		// Reset the connection id back to the one without the lease-id.
		_tracer.SetConnectionId(t.connection.GetConnectionId())
	}

	// Return the connection to the connection container and don't actually close it.
	err := t.connectionContainer.returnConnection(ctx, newState)

	// Detach the connection from this lease, so it can no longer be used by the client.
	t.connection = nil

	return err
}

func (t *plcConnectionLease) GetLastTraces() []tracer.TraceEntry {
	return t.lastTraces
}

func (t *plcConnectionLease) IsConnected() bool {
	if t.connection == nil {
		return false
	}
	if t.isInvalidated() {
		return false
	}
	return t.connection.IsConnected()
}

func (t *plcConnectionLease) Ping(ctx context.Context) error {
	if t.connection == nil {
		panic("Called 'Ping' on a closed cached connection")
	}
	if t.isInvalidated() {
		return errConnectionInvalidated
	}
	return t.connection.Ping(ctx)
}

func (t *plcConnectionLease) GetMetadata() apiModel.PlcConnectionMetadata {
	t.ensureConnection("GetMetadata")
	return t.connection.GetMetadata()
}

func (t *plcConnectionLease) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	t.ensureConnection("ReadRequestBuilder")
	return t.connection.ReadRequestBuilder()
}

func (t *plcConnectionLease) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	t.ensureConnection("WriteRequestBuilder")
	return t.connection.WriteRequestBuilder()
}

func (t *plcConnectionLease) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	t.ensureConnection("SubscriptionRequestBuilder")
	return t.connection.SubscriptionRequestBuilder()
}

func (t *plcConnectionLease) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	t.ensureConnection("UnsubscriptionRequestBuilder")
	return t.connection.UnsubscriptionRequestBuilder()
}

func (t *plcConnectionLease) BrowseRequestBuilder() apiModel.PlcBrowseRequestBuilder {
	t.ensureConnection("BrowseRequestBuilder")
	return t.connection.BrowseRequestBuilder()
}

func (t *plcConnectionLease) Invalidate() {
	if t.connection == nil {
		t.invalidated.Store(true)
		return
	}
	if t.invalidated.Swap(true) {
		return
	}
	t.connection.Invalidate()
}

func (t *plcConnectionLease) ensureConnection(method string) {
	if t.connection == nil {
		panic(fmt.Sprintf("Called '%s' on a closed cached connection", method))
	}
	if t.isInvalidated() {
		panic(fmt.Sprintf("Called '%s' on an invalidated cached connection", method))
	}
}

func (t *plcConnectionLease) isInvalidated() bool {
	return t.invalidated.Load()
}

func (t *plcConnectionLease) String() string {
	return fmt.Sprintf("plcConnectionLease{connectionContainer: %s, leaseId: %d, connection: %s}", t.connectionContainer, t.leaseId, t.connection)
}
