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
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcSubscriptionHandle
type DefaultPlcSubscriptionHandle struct {
	handleToRegister model.PlcSubscriptionHandle `ignore:"true"`
	plcSubscriber    spi.PlcSubscriber
}

// NewDefaultPlcSubscriptionHandle can be used when the DefaultPlcSubscriptionHandle is sufficient
func NewDefaultPlcSubscriptionHandle(plcSubscriber spi.PlcSubscriber) *DefaultPlcSubscriptionHandle {
	handle := &DefaultPlcSubscriptionHandle{
		plcSubscriber: plcSubscriber,
	}
	handle.handleToRegister = handle
	return handle
}

// NewDefaultPlcSubscriptionHandleWithHandleToRegister should be used when an extension of DefaultPlcSubscriptionHandle is used
func NewDefaultPlcSubscriptionHandleWithHandleToRegister(plcSubscriber spi.PlcSubscriber, handleToRegister model.PlcSubscriptionHandle) *DefaultPlcSubscriptionHandle {
	return &DefaultPlcSubscriptionHandle{
		handleToRegister: handleToRegister,
		plcSubscriber:    plcSubscriber,
	}
}

// Register registers at the spi.PlcSubscriber
func (d *DefaultPlcSubscriptionHandle) Register(consumer model.PlcSubscriptionEventConsumer) model.PlcConsumerRegistration {
	return d.plcSubscriber.Register(consumer, []model.PlcSubscriptionHandle{d.handleToRegister})
}
