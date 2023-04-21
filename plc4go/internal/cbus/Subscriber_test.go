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

package cbus

import (
	"context"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewSubscriber(t *testing.T) {
	type args struct {
		connection *Connection
	}
	tests := []struct {
		name string
		args args
		want *Subscriber
	}{
		{
			name: "simple",
			want: &Subscriber{
				consumers: map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewSubscriber(tt.args.connection), "NewSubscriber(%v)", tt.args.connection)
		})
	}
}

func TestSubscriber_Subscribe(t *testing.T) {
	type fields struct {
		connection *Connection
		consumers  map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer
	}
	type args struct {
		in0                 context.Context
		subscriptionRequest apiModel.PlcSubscriptionRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantAsserter func(t *testing.T, results <-chan apiModel.PlcSubscriptionRequestResult) bool
	}{
		{
			name: "just subscribe",
			fields: fields{
				connection: NewConnection(nil, Configuration{}, DriverContext{}, nil, nil, nil),
			},
			args: args{
				in0:                 context.Background(),
				subscriptionRequest: spiModel.NewDefaultPlcSubscriptionRequest(nil, []string{"blub"}, map[string]apiModel.PlcTag{"blub": NewMMIMonitorTag(readWriteModel.NewUnitAddress(1), nil, 1)}, nil, nil, nil),
			},
			wantAsserter: func(t *testing.T, results <-chan apiModel.PlcSubscriptionRequestResult) bool {
				timer := time.NewTimer(2 * time.Second)
				defer timer.Stop()
				select {
				case <-timer.C:
					t.Fail()
				case result := <-results:
					assert.Nilf(t, result.GetErr(), "error %v", result.GetErr())
					assert.NotNil(t, result.GetResponse())
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Subscriber{
				connection: tt.fields.connection,
				consumers:  tt.fields.consumers,
			}
			assert.Truef(t, tt.wantAsserter(t, m.Subscribe(tt.args.in0, tt.args.subscriptionRequest)), "Subscribe(%v, %v)", tt.args.in0, tt.args.subscriptionRequest)
		})
	}
}

func TestSubscriber_Unsubscribe(t *testing.T) {
	type fields struct {
		connection *Connection
		consumers  map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer
	}
	type args struct {
		ctx                   context.Context
		unsubscriptionRequest apiModel.PlcUnsubscriptionRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   <-chan apiModel.PlcUnsubscriptionRequestResult
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Subscriber{
				connection: tt.fields.connection,
				consumers:  tt.fields.consumers,
			}
			assert.Equalf(t, tt.want, m.Unsubscribe(tt.args.ctx, tt.args.unsubscriptionRequest), "Unsubscribe(%v, %v)", tt.args.ctx, tt.args.unsubscriptionRequest)
		})
	}
}

func TestSubscriber_handleMonitoredMMI(t *testing.T) {
	type fields struct {
		connection *Connection
		consumers  map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer
	}
	type args struct {
		calReply model.CALReply
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "handle the MMI short",
			args: args{
				calReply: readWriteModel.NewCALReplyShort(1, nil, nil, nil),
			},
		},
		{
			name: "handle the MMI short with consumer",
			fields: fields{
				consumers: map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer{
					func() *spiModel.DefaultPlcConsumerRegistration {
						registration := spiModel.NewDefaultPlcConsumerRegistration(nil, nil, []apiModel.PlcSubscriptionHandle{
							&SubscriptionHandle{},
						}...)
						return registration.(*spiModel.DefaultPlcConsumerRegistration)
					}(): nil,
				},
			},
			args: args{
				calReply: readWriteModel.NewCALReplyShort(1, nil, nil, nil),
			},
		},
		{
			name: "handle the MMI long unit address",
			args: args{
				calReply: readWriteModel.NewCALReplyLong(
					0,
					readWriteModel.NewUnitAddress(0),
					nil,
					nil,
					nil,
					nil,
					0,
					nil,
					nil,
					nil,
				),
			},
		},
		{
			name: "handle the MMI long bridge address",
			args: args{
				calReply: readWriteModel.NewCALReplyLong(
					1,
					readWriteModel.NewUnitAddress(0),
					readWriteModel.NewBridgeAddress(1),
					nil,
					nil,
					readWriteModel.NewReplyNetwork(
						readWriteModel.NewNetworkRoute(
							readWriteModel.NewNetworkProtocolControlInformation(3, 3),
							[]readWriteModel.BridgeAddress{
								readWriteModel.NewBridgeAddress(2),
								readWriteModel.NewBridgeAddress(3),
							},
						),
						readWriteModel.NewUnitAddress(1),
					),
					0,
					nil,
					nil,
					nil,
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Subscriber{
				connection: tt.fields.connection,
				consumers:  tt.fields.consumers,
			}
			assert.Equalf(t, tt.want, m.handleMonitoredMMI(tt.args.calReply), "handleMonitoredMMI(%v)", tt.args.calReply)
		})
	}
}

func TestSubscriber_offerMMI(t *testing.T) {
	type fields struct {
		connection *Connection
		consumers  map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer
	}
	type args struct {
		unitAddressString  string
		calData            model.CALData
		subscriptionHandle *SubscriptionHandle
		consumer           apiModel.PlcSubscriptionEventConsumer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "offer not fitting tag",
			args: args{
				subscriptionHandle: &SubscriptionHandle{},
			},
		},
		{
			name: "valid monitor tag unmapped",
			args: args{
				subscriptionHandle: &SubscriptionHandle{
					tag: &mmiMonitorTag{},
				},
			},
			want: false,
		},
		// TODO: add other cases
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Subscriber{
				connection: tt.fields.connection,
				consumers:  tt.fields.consumers,
			}
			assert.Equalf(t, tt.want, m.offerMMI(tt.args.unitAddressString, tt.args.calData, tt.args.subscriptionHandle, tt.args.consumer), "offerMMI(%v, %v, %v, %v)", tt.args.unitAddressString, tt.args.calData, tt.args.subscriptionHandle, tt.args.consumer)
		})
	}
}

func TestSubscriber_handleMonitoredSAL(t *testing.T) {
	type fields struct {
		connection *Connection
		consumers  map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer
	}
	type args struct {
		sal model.MonitoredSAL
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "no sal, no consumers",
		},
		{
			name: "handle the SAL short with consumer",
			fields: fields{
				consumers: map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer{
					func() *spiModel.DefaultPlcConsumerRegistration {
						registration := spiModel.NewDefaultPlcConsumerRegistration(nil, nil, []apiModel.PlcSubscriptionHandle{
							&SubscriptionHandle{},
						}...)
						return registration.(*spiModel.DefaultPlcConsumerRegistration)
					}(): nil,
				},
			},
			args: args{
				sal: readWriteModel.NewMonitoredSALLongFormSmartMode(0, readWriteModel.NewUnitAddress(0), nil, readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74, nil, nil, nil, 0, nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Subscriber{
				connection: tt.fields.connection,
				consumers:  tt.fields.consumers,
			}
			assert.Equalf(t, tt.want, m.handleMonitoredSAL(tt.args.sal), "handleMonitoredSAL(%v)", tt.args.sal)
		})
	}
}

func TestSubscriber_offerSAL(t *testing.T) {
	type fields struct {
		connection *Connection
		consumers  map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer
	}
	type args struct {
		sal                model.MonitoredSAL
		subscriptionHandle *SubscriptionHandle
		consumer           apiModel.PlcSubscriptionEventConsumer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "offer wong tag",
			args: args{
				subscriptionHandle: &SubscriptionHandle{},
			},
		},
		{
			name: "offer sal tag",
			args: args{
				sal: readWriteModel.NewMonitoredSALLongFormSmartMode(
					0,
					readWriteModel.NewUnitAddress(0),
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					nil,
					nil,
					readWriteModel.NewSALDataLighting(
						readWriteModel.NewLightingDataOn(2, readWriteModel.LightingCommandTypeContainer_LightingCommandOn),
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: &SubscriptionHandle{
					tag: &salMonitorTag{},
				},
				consumer: func(_ apiModel.PlcSubscriptionEvent) {
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Subscriber{
				connection: tt.fields.connection,
				consumers:  tt.fields.consumers,
			}
			assert.Equalf(t, tt.want, m.offerSAL(tt.args.sal, tt.args.subscriptionHandle, tt.args.consumer), "offerSAL(\n%v\n, %v)", tt.args.sal, tt.args.subscriptionHandle)
		})
	}
}

func TestSubscriber_Register(t *testing.T) {
	type fields struct {
		connection *Connection
		consumers  map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer
	}
	type args struct {
		consumer apiModel.PlcSubscriptionEventConsumer
		handles  []apiModel.PlcSubscriptionHandle
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "register something",
			fields: fields{
				consumers: map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Subscriber{
				connection: tt.fields.connection,
				consumers:  tt.fields.consumers,
			}
			assert.NotNilf(t, m.Register(tt.args.consumer, tt.args.handles), "Register(func(), %v)", tt.args.handles)
		})
	}
}

func TestSubscriber_Unregister(t *testing.T) {
	type fields struct {
		connection *Connection
		consumers  map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer
	}
	type args struct {
		registration apiModel.PlcConsumerRegistration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Subscriber{
				connection: tt.fields.connection,
				consumers:  tt.fields.consumers,
			}
			m.Unregister(tt.args.registration)
		})
	}
}
