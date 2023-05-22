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

package simulated

import (
	"context"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/pkg/errors"
)

type Subscriber struct {
	device  *Device
	options map[string][]string
	tracer  *spi.Tracer
}

func NewSubscriber(device *Device, options map[string][]string, tracer *spi.Tracer) *Subscriber {
	return &Subscriber{
		device:  device,
		options: options,
		tracer:  tracer,
	}
}

func (r Subscriber) Subscribe(_ context.Context, subscriptionRequest apiModel.PlcSubscriptionRequest) <-chan apiModel.PlcSubscriptionRequestResult {
	// TODO: handle context
	result := make(chan apiModel.PlcSubscriptionRequestResult)
	result <- spiModel.NewDefaultPlcSubscriptionRequestResult(subscriptionRequest, nil, errors.New("Not Implemented"))
	// TODO: implement me
	return result
}

func (r Subscriber) Unsubscribe(_ context.Context, unsubscriptionRequest apiModel.PlcUnsubscriptionRequest) <-chan apiModel.PlcUnsubscriptionRequestResult {
	// TODO: handle context
	result := make(chan apiModel.PlcUnsubscriptionRequestResult, 1)
	result <- spiModel.NewDefaultPlcUnsubscriptionRequestResult(unsubscriptionRequest, nil, errors.New("Not Implemented"))
	// TODO: implement me
	return result
}

func (r Subscriber) Register(consumer apiModel.PlcSubscriptionEventConsumer, handles []apiModel.PlcSubscriptionHandle) apiModel.PlcConsumerRegistration {
	// TODO: implement me
	return nil
}

func (r Subscriber) Unregister(registration apiModel.PlcConsumerRegistration) {
	// TODO: implement me
}
