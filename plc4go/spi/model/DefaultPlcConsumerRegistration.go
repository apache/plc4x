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

package model

import (
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	"math/rand"
)

var _ apiModel.PlcConsumerRegistration = &DefaultPlcConsumerRegistration{}

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcConsumerRegistration
type DefaultPlcConsumerRegistration struct {
	consumerId    int
	consumer      apiModel.PlcSubscriptionEventConsumer `ignore:"true"` // Function not renderable
	plcSubscriber spi.PlcSubscriber
	handles       []apiModel.PlcSubscriptionHandle
}

func NewDefaultPlcConsumerRegistration(plcSubscriber spi.PlcSubscriber, consumer apiModel.PlcSubscriptionEventConsumer, handles ...apiModel.PlcSubscriptionHandle) apiModel.PlcConsumerRegistration {
	return &DefaultPlcConsumerRegistration{
		// TODO: we need a way to hash the consumer
		consumerId:    rand.Int(),
		consumer:      consumer,
		plcSubscriber: plcSubscriber,
		handles:       handles,
	}
}

func (d *DefaultPlcConsumerRegistration) GetConsumerId() int {
	return d.consumerId
}

func (d *DefaultPlcConsumerRegistration) GetSubscriptionHandles() []apiModel.PlcSubscriptionHandle {
	return d.handles
}

func (d *DefaultPlcConsumerRegistration) Unregister() {
	d.plcSubscriber.Unregister(d)
}
