//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package plc4go

import "plc4x.apache.org/plc4go-modbus-driver/0.8.0/pkg/plc4go/model"

type PlcConnection interface {
	// Initiate the connection to the PLC
	Connect() <-chan error
	// Close the connection to the PLC (gracefully)
	Close() <-chan error
	// Checks if the connection is currently still connected
	IsConnected() bool

	// Get some metadata regarding the current connection
	GetMetadata() (model.PlcConnectionMetadata, error)

	// Executes a no-op operation to check if the current connection is still able to communicate
	Ping() error

	// Create a builder for assembling read-requests
	ReadRequestBuilder() (model.PlcReadRequestBuilder, error)
	// Create a builder for assembling write-requests
	WriteRequestBuilder() (model.PlcWriteRequestBuilder, error)
	// Create a builder for assembling subscription-requests
	SubscriptionRequestBuilder() (model.PlcSubscriptionRequestBuilder, error)
	// Create a builder for assembling unsubscription-requests
	UnsubscriptionRequestBuilder() (model.PlcUnsubscriptionRequestBuilder, error)
}
