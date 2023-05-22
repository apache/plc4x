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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestDefaultPlcConsumerRegistration_GetConsumerId(t *testing.T) {
	type fields struct {
		consumerId    int
		consumer      apiModel.PlcSubscriptionEventConsumer
		plcSubscriber spi.PlcSubscriber
		handles       []apiModel.PlcSubscriptionHandle
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcConsumerRegistration{
				consumerId:    tt.fields.consumerId,
				consumer:      tt.fields.consumer,
				plcSubscriber: tt.fields.plcSubscriber,
				handles:       tt.fields.handles,
			}
			assert.Equalf(t, tt.want, d.GetConsumerId(), "GetConsumerId()")
		})
	}
}

func TestDefaultPlcConsumerRegistration_GetSubscriptionHandles(t *testing.T) {
	type fields struct {
		consumerId    int
		consumer      apiModel.PlcSubscriptionEventConsumer
		plcSubscriber spi.PlcSubscriber
		handles       []apiModel.PlcSubscriptionHandle
	}
	tests := []struct {
		name   string
		fields fields
		want   []apiModel.PlcSubscriptionHandle
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcConsumerRegistration{
				consumerId:    tt.fields.consumerId,
				consumer:      tt.fields.consumer,
				plcSubscriber: tt.fields.plcSubscriber,
				handles:       tt.fields.handles,
			}
			assert.Equalf(t, tt.want, d.GetSubscriptionHandles(), "GetSubscriptionHandles()")
		})
	}
}

func TestDefaultPlcConsumerRegistration_Unregister(t *testing.T) {
	type fields struct {
		consumerId    int
		consumer      apiModel.PlcSubscriptionEventConsumer
		plcSubscriber spi.PlcSubscriber
		handles       []apiModel.PlcSubscriptionHandle
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "unregister it",
			fields: fields{
				plcSubscriber: func() spi.PlcSubscriber {
					subscriber := NewMockPlcSubscriber(t)
					subscriber.EXPECT().Unregister(mock.Anything).Return()
					return subscriber
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcConsumerRegistration{
				consumerId:    tt.fields.consumerId,
				consumer:      tt.fields.consumer,
				plcSubscriber: tt.fields.plcSubscriber,
				handles:       tt.fields.handles,
			}
			d.Unregister()
		})
	}
}

func TestNewDefaultPlcConsumerRegistration(t *testing.T) {
	type args struct {
		plcSubscriber spi.PlcSubscriber
		consumer      apiModel.PlcSubscriptionEventConsumer
		handles       []apiModel.PlcSubscriptionHandle
	}
	tests := []struct {
		name string
		args args
		want apiModel.PlcConsumerRegistration
	}{
		{
			name: "create it",
			want: &DefaultPlcConsumerRegistration{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDefaultPlcConsumerRegistration(tt.args.plcSubscriber, tt.args.consumer, tt.args.handles...).(*DefaultPlcConsumerRegistration)
			tt.want.(*DefaultPlcConsumerRegistration).consumerId = got.consumerId
			assert.Equalf(t, tt.want, got, "NewDefaultPlcConsumerRegistration(%v, func(), %v)", tt.args.plcSubscriber, tt.args.handles)
		})
	}
}
