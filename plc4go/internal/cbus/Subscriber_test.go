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
	"github.com/rs/zerolog/log"
	"testing"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/stretchr/testify/assert"
)

func TestNewSubscriber(t *testing.T) {
	type args struct {
		addSubscriber func(subscriber *Subscriber)
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
				log:       log.Logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewSubscriber(tt.args.addSubscriber), "NewSubscriber(%t)", tt.args.addSubscriber != nil)
		})
	}
}

func TestSubscriber_Subscribe(t *testing.T) {
	type fields struct {
		addSubscriber func(subscriber *Subscriber)
		consumers     map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer
	}
	type args struct {
		in0                 context.Context
		subscriptionRequest apiModel.PlcSubscriptionRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		setup        func(t *testing.T, fields *fields, args *args)
		wantAsserter func(t *testing.T, results <-chan apiModel.PlcSubscriptionRequestResult) bool
	}{
		{
			name: "just subscribe",
			args: args{
				in0:                 context.Background(),
				subscriptionRequest: spiModel.NewDefaultPlcSubscriptionRequest(nil, []string{"blub"}, map[string]apiModel.PlcSubscriptionTag{"blub": NewMMIMonitorTag(readWriteModel.NewUnitAddress(1), nil, 1)}, nil, nil, nil),
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				fields.addSubscriber = func(subscriber *Subscriber) {
					assert.NotNil(t, subscriber)
				}
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
			if tt.setup != nil {
				tt.setup(t, &tt.fields, &tt.args)
			}
			m := &Subscriber{
				addSubscriber: tt.fields.addSubscriber,
				consumers:     tt.fields.consumers,
			}
			assert.Truef(t, tt.wantAsserter(t, m.Subscribe(tt.args.in0, tt.args.subscriptionRequest)), "Subscribe(%v, %v)", tt.args.in0, tt.args.subscriptionRequest)
		})
	}
}

func TestSubscriber_Unsubscribe(t *testing.T) {
	type fields struct {
		addSubscriber func(subscriber *Subscriber)
		consumers     map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer
	}
	type args struct {
		ctx                   context.Context
		unsubscriptionRequest apiModel.PlcUnsubscriptionRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantAsserter func(t *testing.T, results <-chan apiModel.PlcUnsubscriptionRequestResult) bool
	}{
		{
			name: "just do it",
			wantAsserter: func(t *testing.T, results <-chan apiModel.PlcUnsubscriptionRequestResult) bool {
				assert.NotNil(t, results)
				// TODO: add tests once implemented
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Subscriber{
				addSubscriber: tt.fields.addSubscriber,
				consumers:     tt.fields.consumers,
			}
			assert.Truef(t, tt.wantAsserter(t, m.Unsubscribe(tt.args.ctx, tt.args.unsubscriptionRequest)), "Unsubscribe(%v, %v)", tt.args.ctx, tt.args.unsubscriptionRequest)
		})
	}
}

func TestSubscriber_handleMonitoredMMI(t *testing.T) {
	type fields struct {
		addSubscriber func(subscriber *Subscriber)
		consumers     map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer
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
			name: "handle the MMI short with consumerProvider",
			fields: fields{
				consumers: map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer{
					func() *spiModel.DefaultPlcConsumerRegistration {
						registration := spiModel.NewDefaultPlcConsumerRegistration(nil, nil, []apiModel.PlcSubscriptionHandle{
							&SubscriptionHandle{
								DefaultPlcSubscriptionHandle: spiModel.NewDefaultPlcSubscriptionHandle(nil).(*spiModel.DefaultPlcSubscriptionHandle),
							},
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
				addSubscriber: tt.fields.addSubscriber,
				consumers:     tt.fields.consumers,
			}
			assert.Equalf(t, tt.want, m.handleMonitoredMMI(tt.args.calReply), "handleMonitoredMMI(%v)", tt.args.calReply)
		})
	}
}

func TestSubscriber_offerMMI(t *testing.T) {
	type fields struct {
		addSubscriber func(subscriber *Subscriber)
		consumers     map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer
	}
	type args struct {
		unitAddressString  string
		calData            model.CALData
		subscriptionHandle *SubscriptionHandle
		consumerProvider   func(t *testing.T) apiModel.PlcSubscriptionEventConsumer
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
				subscriptionHandle: NewSubscriptionHandle(
					nil,
					"nada",
					nil,
					apiModel.SubscriptionEvent,
					0,
				),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(_ apiModel.PlcSubscriptionEvent) {
						t.Error("should not be called")
					}
				},
			},
		},
		{
			name: "valid monitor tag wrong address",
			args: args{
				unitAddressString: "banana",
				subscriptionHandle: NewSubscriptionHandle(
					nil,
					"tag",
					NewMMIMonitorTag(
						readWriteModel.NewUnitAddress(13),
						nil,
						1,
					),
					apiModel.SubscriptionEvent,
					0,
				),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(_ apiModel.PlcSubscriptionEvent) {
						t.Error("should not be called")
					}
				},
			},
			want: false,
		},
		{
			name: "valid monitor tag unmapped",
			args: args{
				subscriptionHandle: NewSubscriptionHandle(
					nil,
					"nada",
					nil,
					apiModel.SubscriptionEvent,
					0,
				),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(_ apiModel.PlcSubscriptionEvent) {
						t.Error("should not be called")
					}
				},
			},
			want: false,
		},
		{
			name: "valid monitor tag cal unrelated",
			args: args{
				unitAddressString: "u13",
				calData:           readWriteModel.NewCALDataReset(readWriteModel.CALCommandTypeContainer_CALCommandGetStatus, nil, nil),
				subscriptionHandle: NewSubscriptionHandle(
					nil,
					"tag",
					NewMMIMonitorTag(
						readWriteModel.NewUnitAddress(13),
						nil,
						1,
					),
					apiModel.SubscriptionEvent,
					0,
				),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: false,
		},
		{
			name: "valid monitor tag cal status",
			args: args{
				unitAddressString: "u13",
				calData: readWriteModel.NewCALDataStatus(
					readWriteModel.ApplicationIdContainer_LIGHTING_3A,
					0,
					[]readWriteModel.StatusByte{
						readWriteModel.NewStatusByte(readWriteModel.GAVState_DOES_NOT_EXIST, readWriteModel.GAVState_OFF, readWriteModel.GAVState_ON, readWriteModel.GAVState_ERROR),
					},
					readWriteModel.CALCommandTypeContainer_CALCommandIdentify,
					nil,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(
					nil,
					"tag",
					NewMMIMonitorTag(
						readWriteModel.NewUnitAddress(13),
						nil,
						1,
					),
					apiModel.SubscriptionEvent,
					0,
				),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "valid monitor tag cal status wrong application",
			args: args{
				unitAddressString: "u13",
				calData: readWriteModel.NewCALDataStatus(
					readWriteModel.ApplicationIdContainer_LIGHTING_3A,
					0,
					[]readWriteModel.StatusByte{
						readWriteModel.NewStatusByte(readWriteModel.GAVState_DOES_NOT_EXIST, readWriteModel.GAVState_OFF, readWriteModel.GAVState_ON, readWriteModel.GAVState_ERROR),
					},
					readWriteModel.CALCommandTypeContainer_CALCommandIdentify,
					nil,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(
					nil,
					"tag",
					NewMMIMonitorTag(
						readWriteModel.NewUnitAddress(13),
						func() *readWriteModel.ApplicationIdContainer {
							a := readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74
							return &a
						}(),
						1,
					),
					apiModel.SubscriptionEvent,
					0,
				),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: false,
		},
		{
			name: "valid monitor tag cal status extended binary",
			args: args{
				unitAddressString: "u13",
				calData: readWriteModel.NewCALDataStatusExtended(
					readWriteModel.StatusCoding_BINARY_BY_THIS_SERIAL_INTERFACE,
					readWriteModel.ApplicationIdContainer_LIGHTING_3A,
					0,
					[]readWriteModel.StatusByte{
						readWriteModel.NewStatusByte(readWriteModel.GAVState_DOES_NOT_EXIST, readWriteModel.GAVState_OFF, readWriteModel.GAVState_ON, readWriteModel.GAVState_ERROR),
					},
					nil,
					readWriteModel.CALCommandTypeContainer_CALCommandIdentify,
					nil,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(
					nil,
					"tag",
					NewMMIMonitorTag(
						readWriteModel.NewUnitAddress(13),
						nil,
						1,
					),
					apiModel.SubscriptionEvent,
					0,
				),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "valid monitor tag cal status extended level",
			args: args{
				unitAddressString: "u13",
				calData: readWriteModel.NewCALDataStatusExtended(
					readWriteModel.StatusCoding_LEVEL_BY_THIS_SERIAL_INTERFACE,
					readWriteModel.ApplicationIdContainer_LIGHTING_3A,
					0,
					nil,
					[]readWriteModel.LevelInformation{
						readWriteModel.NewLevelInformationAbsent(13),
						readWriteModel.NewLevelInformationCorrupted(1, 2, 3, 4, 5),
						readWriteModel.NewLevelInformationNormal(readWriteModel.LevelInformationNibblePair_Value_0, readWriteModel.LevelInformationNibblePair_Value_2, 13),
					},
					readWriteModel.CALCommandTypeContainer_CALCommandIdentify,
					nil,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(
					nil,
					"tag",
					NewMMIMonitorTag(
						readWriteModel.NewUnitAddress(13),
						nil,
						1,
					),
					apiModel.SubscriptionEvent,
					0,
				),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Subscriber{
				addSubscriber: tt.fields.addSubscriber,
				consumers:     tt.fields.consumers,
			}
			assert.Equalf(t, tt.want, m.offerMMI(tt.args.unitAddressString, tt.args.calData, tt.args.subscriptionHandle, tt.args.consumerProvider(t)), "offerMMI(%v,\n%v\n, \n%v\n, func())", tt.args.unitAddressString, tt.args.calData, tt.args.subscriptionHandle)
		})
	}
}

func TestSubscriber_handleMonitoredSAL(t *testing.T) {
	type fields struct {
		addSubscriber func(subscriber *Subscriber)
		consumers     map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer
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
			name: "handle the SAL short with consumerProvider",
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
				addSubscriber: tt.fields.addSubscriber,
				consumers:     tt.fields.consumers,
			}
			assert.Equalf(t, tt.want, m.handleMonitoredSAL(tt.args.sal), "handleMonitoredSAL(%v)", tt.args.sal)
		})
	}
}

func TestSubscriber_offerSAL(t *testing.T) {
	type fields struct {
		addSubscriber func(subscriber *Subscriber)
		consumers     map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer
	}
	type args struct {
		sal                model.MonitoredSAL
		subscriptionHandle *SubscriptionHandle
		consumerProvider   func(t *testing.T) apiModel.PlcSubscriptionEventConsumer
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
				subscriptionHandle: NewSubscriptionHandle(nil, "", nil, apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						t.Fail()
					}
				},
			},
		},
		{
			name: "offer sal tag short",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataLighting(
						readWriteModel.NewLightingDataOn(2, readWriteModel.LightingCommandTypeContainer_LightingCommandOn),
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag short sal access control",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataAccessControl(
						readWriteModel.NewAccessControlDataAccessPointClosed(readWriteModel.AccessControlCommandTypeContainer_AccessControlCommandAccessPointClosed, 0, 0),
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag short sal air conditioning",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataAirConditioning(
						readWriteModel.NewAirConditioningDataRefresh(0, readWriteModel.AirConditioningCommandTypeContainer_AirConditioningCommandRefresh),
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag short sal audio & video",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataAudioAndVideo(
						readWriteModel.NewLightingDataOff(0, readWriteModel.LightingCommandTypeContainer_LightingCommandOn),
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag short sal clock and timekeeping",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataClockAndTimekeeping(
						readWriteModel.NewClockAndTimekeepingDataRequestRefresh(
							1,
							0,
						),
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag short sal enable control",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataEnableControl(
						readWriteModel.NewEnableControlData(0, 0, 0),
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag short sal error reporting",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataErrorReporting(
						readWriteModel.NewErrorReportingDataGeneric(
							readWriteModel.NewErrorReportingSystemCategory(
								readWriteModel.ErrorReportingSystemCategoryClass_INPUT_UNITS,
								readWriteModel.NewErrorReportingSystemCategoryTypeInputUnits(
									readWriteModel.ErrorReportingSystemCategoryTypeForInputUnits_RESERVED_2,
								),
								readWriteModel.ErrorReportingSystemCategoryVariant_RESERVED_0,
							),
							true,
							true,
							true,
							readWriteModel.ErrorReportingSeverity_ALL_OK,
							1,
							2,
							3,
							4,
						),
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag short sal free usage",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataFreeUsage(nil),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag short heating",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataHeating(
						readWriteModel.NewLightingDataOn(2, readWriteModel.LightingCommandTypeContainer_LightingCommandOn),
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag short actuator",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataHvacActuator(
						readWriteModel.NewLightingDataOn(2, readWriteModel.LightingCommandTypeContainer_LightingCommandOn),
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag short irrigation control",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataIrrigationControl(
						readWriteModel.NewLightingDataOn(2, readWriteModel.LightingCommandTypeContainer_LightingCommandOn),
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag short lighting",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataLighting(
						readWriteModel.NewLightingDataOn(2, readWriteModel.LightingCommandTypeContainer_LightingCommandOn),
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag short measurement",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataMeasurement(
						readWriteModel.NewMeasurementDataChannelMeasurementData(
							0,
							0,
							readWriteModel.MeasurementUnits_ANGLE_DEGREES,
							0,
							0,
							0,
							0,
						),
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag short media transport",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataMediaTransport(
						readWriteModel.NewMediaTransportControlDataFastForward(
							0,
							0,
							0,
						),
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag short metering",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataMetering(
						readWriteModel.NewMeteringDataGasConsumption(
							0,
							0,
							0,
						),
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag short pools spas ponds",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataPoolsSpasPondsFountainsControl(
						readWriteModel.NewLightingDataOn(
							0,
							0,
						),
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag short reserved",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataReserved(
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag short reserved",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataRoomControlSystem(
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag short reserved",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataSecurity(
						readWriteModel.NewSecurityDataAlarmOn(
							0,
							0,
						),
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag short telephony",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataTelephonyStatusAndControl(
						readWriteModel.NewTelephonyDataDivert(
							"1234",
							0,
							0,
						),
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag short temperature broadcast",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataTemperatureBroadcast(
						readWriteModel.NewTemperatureBroadcastData(
							0,
							0,
							0,
						),
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag short testing",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataTesting(
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag short trigger control",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataTriggerControl(
						readWriteModel.NewTriggerControlDataTriggerEvent(
							0,
							0,
							0,
						),
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag short ventilation",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataVentilation(
						readWriteModel.NewLightingDataOn(
							0,
							0,
						),
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag short wrong unit address",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataLighting(
						readWriteModel.NewLightingDataOn(2, readWriteModel.LightingCommandTypeContainer_LightingCommandOn),
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(
					nil,
					"",
					NewSALMonitorTag(
						readWriteModel.NewUnitAddress(13),
						nil,
						1,
					),
					apiModel.SubscriptionEvent,
					0,
				),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: false,
		},
		{
			name: "offer sal tag short wrong application",
			args: args{
				sal: readWriteModel.NewMonitoredSALShortFormBasicMode(
					0,
					nil,
					nil,
					nil,
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					readWriteModel.NewSALDataLighting(
						readWriteModel.NewLightingDataOn(2, readWriteModel.LightingCommandTypeContainer_LightingCommandOn),
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(
					nil,
					"",
					NewSALMonitorTag(
						nil,
						func() *readWriteModel.ApplicationIdContainer {
							a := readWriteModel.ApplicationIdContainer_LIGHTING_3A
							return &a
						}(),
						1,
					),
					apiModel.SubscriptionEvent,
					0,
				),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: false,
		},
		{
			name: "offer sal tag long",
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
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
		{
			name: "offer sal tag long bridged",
			args: args{
				sal: readWriteModel.NewMonitoredSALLongFormSmartMode(
					1,
					nil,
					readWriteModel.NewBridgeAddress(2),
					readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
					nil,
					readWriteModel.NewReplyNetwork(
						readWriteModel.NewNetworkRoute(
							readWriteModel.NewNetworkProtocolControlInformation(1, 1),
							[]readWriteModel.BridgeAddress{
								readWriteModel.NewBridgeAddress(2),
							},
						),
						readWriteModel.NewUnitAddress(0),
					),
					readWriteModel.NewSALDataLighting(
						readWriteModel.NewLightingDataOn(2, readWriteModel.LightingCommandTypeContainer_LightingCommandOn),
						nil,
					),
					0,
					nil,
				),
				subscriptionHandle: NewSubscriptionHandle(nil, "", NewSALMonitorTag(nil, nil, 1), apiModel.SubscriptionEvent, 0),
				consumerProvider: func(t *testing.T) apiModel.PlcSubscriptionEventConsumer {
					return func(event apiModel.PlcSubscriptionEvent) {
						assert.NotNil(t, event)
					}
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Subscriber{
				addSubscriber: tt.fields.addSubscriber,
				consumers:     tt.fields.consumers,
			}
			assert.Equalf(t, tt.want, m.offerSAL(tt.args.sal, tt.args.subscriptionHandle, tt.args.consumerProvider(t)), "offerSAL(\n%v\n, \n%v\n)", tt.args.sal, tt.args.subscriptionHandle)
		})
	}
}

func TestSubscriber_Register(t *testing.T) {
	type fields struct {
		addSubscriber func(subscriber *Subscriber)
		consumers     map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer
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
				addSubscriber: tt.fields.addSubscriber,
				consumers:     tt.fields.consumers,
			}
			assert.NotNilf(t, m.Register(tt.args.consumer, tt.args.handles), "Register(func(), %v)", tt.args.handles)
		})
	}
}

func TestSubscriber_Unregister(t *testing.T) {
	type fields struct {
		addSubscriber func(subscriber *Subscriber)
		consumers     map[*spiModel.DefaultPlcConsumerRegistration]apiModel.PlcSubscriptionEventConsumer
	}
	type args struct {
		registration apiModel.PlcConsumerRegistration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "just do it",
			args: args{
				registration: spiModel.NewDefaultPlcConsumerRegistration(nil, nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Subscriber{
				addSubscriber: tt.fields.addSubscriber,
				consumers:     tt.fields.consumers,
			}
			m.Unregister(tt.args.registration)
		})
	}
}
