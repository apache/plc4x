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

package ads

import (
	"context"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	internalModel "github.com/apache/plc4x/plc4go/spi/model"
)

func (m *Connection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	return internalModel.NewDefaultPlcSubscriptionRequestBuilder(m.GetPlcTagHandler(), m.GetPlcValueHandler(), m)
}

func (m *Connection) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	// TODO: Implement this ...
	return nil
}

func (m *Connection) Subscribe(ctx context.Context, subscriptionRequest apiModel.PlcSubscriptionRequest) <-chan apiModel.PlcSubscriptionRequestResult {
	// Unfortunately the subscription API is the only part that doesn't support multi-item requests.
	return nil
}

func (m *Connection) Unsubscribe(ctx context.Context, unsubscriptionRequest apiModel.PlcUnsubscriptionRequest) <-chan apiModel.PlcUnsubscriptionRequestResult {
	return nil
}

func (m *Connection) Register(consumer apiModel.PlcSubscriptionEventConsumer, handles []apiModel.PlcSubscriptionHandle) apiModel.PlcConsumerRegistration {
	return nil
}

func (m *Connection) Unregister(registration apiModel.PlcConsumerRegistration) {
}
