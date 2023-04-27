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

package _default

import (
	"context"
	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
)

// TODO: replace with proper mock
type testConnection struct {
	testMessageCodec testMessageCodec
	connected        bool
}

func (t testConnection) Connect() <-chan plc4go.PlcConnectionConnectResult {
	return t.ConnectWithContext(nil)
}

func (t testConnection) BlockingClose() {
	// NO-OP
}

func (t testConnection) Close() <-chan plc4go.PlcConnectionCloseResult {
	// NO-OP
	return nil
}

func (t testConnection) IsConnected() bool {
	return t.connected
}

func (t testConnection) Ping() <-chan plc4go.PlcConnectionPingResult {
	// NO-OP
	return nil
}

func (t testConnection) GetMetadata() model.PlcConnectionMetadata {
	// NO-OP
	return nil
}

func (t testConnection) ReadRequestBuilder() model.PlcReadRequestBuilder {
	// NO-OP
	return nil
}

func (t testConnection) WriteRequestBuilder() model.PlcWriteRequestBuilder {
	// NO-OP
	return nil
}

func (t testConnection) SubscriptionRequestBuilder() model.PlcSubscriptionRequestBuilder {
	// NO-OP
	return nil
}

func (t testConnection) UnsubscriptionRequestBuilder() model.PlcUnsubscriptionRequestBuilder {
	// NO-OP
	return nil
}

func (t testConnection) BrowseRequestBuilder() model.PlcBrowseRequestBuilder {
	// NO-OP
	return nil
}

func (t testConnection) GetConnection() plc4go.PlcConnection {
	// NO-OP
	return t
}

func (t testConnection) GetMessageCodec() spi.MessageCodec {
	return t.testMessageCodec
}

func (t testConnection) ConnectWithContext(_ context.Context) <-chan plc4go.PlcConnectionConnectResult {
	// NO-OP
	return nil
}
