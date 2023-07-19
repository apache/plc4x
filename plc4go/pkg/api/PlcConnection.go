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

package plc4go

import (
	"context"
	"fmt"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
)

type PlcConnection interface {
	fmt.Stringer
	// Connect Initiate the connection to the PLC
	Connect() <-chan PlcConnectionConnectResult
	// ConnectWithContext connects connection codec with the supplied context
	ConnectWithContext(ctx context.Context) <-chan PlcConnectionConnectResult
	// BlockingClose Blocking variant of Close (for usage in "defer" statements)
	BlockingClose()
	// Close the connection to the PLC (gracefully)
	Close() <-chan PlcConnectionCloseResult
	// IsConnected Checks if the connection is currently still connected
	IsConnected() bool

	// Ping Executes a no-op operation to check if the current connection is still able to communicate
	Ping() <-chan PlcConnectionPingResult

	// GetMetadata Get some metadata regarding the current connection
	GetMetadata() model.PlcConnectionMetadata

	// ReadRequestBuilder Create a builder for assembling read-requests
	ReadRequestBuilder() model.PlcReadRequestBuilder
	// WriteRequestBuilder Create a builder for assembling write-requests
	WriteRequestBuilder() model.PlcWriteRequestBuilder
	// SubscriptionRequestBuilder Create a builder for assembling subscription-requests
	SubscriptionRequestBuilder() model.PlcSubscriptionRequestBuilder
	// UnsubscriptionRequestBuilder Create a builder for assembling unsubscription-requests
	UnsubscriptionRequestBuilder() model.PlcUnsubscriptionRequestBuilder
	// BrowseRequestBuilder Create a builder for assembling browser-requests
	BrowseRequestBuilder() model.PlcBrowseRequestBuilder
}
