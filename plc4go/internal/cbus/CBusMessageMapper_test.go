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
	"fmt"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/pool"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTagToCBusMessage(t *testing.T) {
	type args struct {
		tag            apiModel.PlcTag
		value          apiValues.PlcValue
		alphaGenerator *AlphaGenerator
		messageCodec   *MessageCodec
	}
	cbusOptions := readWriteModel.NewCBusOptions(false, false, false, false, false, false, false, false, true)
	requestContext := readWriteModel.NewRequestContext(false)
	tests := []struct {
		name                  string
		args                  args
		wantCBusMessage       readWriteModel.CBusMessage
		wantSupportsRead      bool
		wantSupportsWrite     bool
		wantSupportsSubscribe bool
		wantErr               bool
	}{
		{
			name: "direct status binary",
			args: args{
				tag:            NewStatusTag(nil, StatusRequestTypeBinaryState, nil, readWriteModel.ApplicationIdContainer_LIGHTING_3A, 1),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantCBusMessage: readWriteModel.NewCBusMessageToServer(
				readWriteModel.NewRequestCommand(
					readWriteModel.NewCBusCommandPointToMultiPoint(
						readWriteModel.NewCBusPointToMultiPointCommandStatus(
							readWriteModel.NewStatusRequestBinaryState(
								readWriteModel.ApplicationIdContainer_LIGHTING_3A,
								122,
							),
							58,
							cbusOptions,
						),
						readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToMultiPoint),
						cbusOptions,
					),
					nil,
					readWriteModel.NewAlpha(0),
					readWriteModel.RequestType_REQUEST_COMMAND,
					nil,
					nil,
					readWriteModel.RequestType_EMPTY,
					readWriteModel.NewRequestTermination(),
					cbusOptions,
				),
				requestContext,
				cbusOptions,
			),
			wantSupportsRead:      true,
			wantSupportsSubscribe: true,
		},
		{
			name: "direct status level",
			args: args{
				tag: NewStatusTag(nil, StatusRequestTypeLevel, func() *byte {
					var b byte = 13
					return &b
				}(), readWriteModel.ApplicationIdContainer_LIGHTING_3A, 1),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantCBusMessage: readWriteModel.NewCBusMessageToServer(
				readWriteModel.NewRequestCommand(
					readWriteModel.NewCBusCommandPointToMultiPoint(
						readWriteModel.NewCBusPointToMultiPointCommandStatus(
							readWriteModel.NewStatusRequestLevel(
								readWriteModel.ApplicationIdContainer_LIGHTING_3A,
								13,
								115,
							),
							58,
							cbusOptions,
						),
						readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToMultiPoint),
						cbusOptions,
					),
					nil,
					readWriteModel.NewAlpha(0),
					readWriteModel.RequestType_REQUEST_COMMAND,
					nil,
					nil,
					readWriteModel.RequestType_EMPTY,
					readWriteModel.NewRequestTermination(),
					cbusOptions,
				),
				requestContext,
				cbusOptions,
			),
			wantSupportsRead:      true,
			wantSupportsSubscribe: true,
		},
		{
			name: "direct status binary bridged one bridge",
			args: args{
				tag:            NewStatusTag([]readWriteModel.BridgeAddress{readWriteModel.NewBridgeAddress(12)}, StatusRequestTypeBinaryState, nil, readWriteModel.ApplicationIdContainer_LIGHTING_3A, 1),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantCBusMessage: readWriteModel.NewCBusMessageToServer(
				readWriteModel.NewRequestCommand(
					readWriteModel.NewCBusCommandPointToPointToMultiPoint(
						readWriteModel.NewCBusPointToPointToMultiPointCommandStatus(
							readWriteModel.NewStatusRequestBinaryState(
								readWriteModel.ApplicationIdContainer_LIGHTING_3A,
								122,
							),
							readWriteModel.NewBridgeAddress(12),
							readWriteModel.NewNetworkRoute(readWriteModel.NewNetworkProtocolControlInformation(1, 1), []readWriteModel.BridgeAddress{}),
							58,
							cbusOptions,
						),
						readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPointToMultiPoint),
						cbusOptions,
					),
					nil,
					readWriteModel.NewAlpha(0),
					readWriteModel.RequestType_REQUEST_COMMAND,
					nil,
					nil,
					readWriteModel.RequestType_EMPTY,
					readWriteModel.NewRequestTermination(),
					cbusOptions,
				),
				requestContext,
				cbusOptions,
			),
			wantSupportsRead:      true,
			wantSupportsSubscribe: true,
		},
		{
			name: "direct status binary bridged two bridges",
			args: args{
				tag:            NewStatusTag([]readWriteModel.BridgeAddress{readWriteModel.NewBridgeAddress(12), readWriteModel.NewBridgeAddress(13)}, StatusRequestTypeBinaryState, nil, readWriteModel.ApplicationIdContainer_LIGHTING_3A, 1),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantCBusMessage: readWriteModel.NewCBusMessageToServer(
				readWriteModel.NewRequestCommand(
					readWriteModel.NewCBusCommandPointToPointToMultiPoint(
						readWriteModel.NewCBusPointToPointToMultiPointCommandStatus(
							readWriteModel.NewStatusRequestBinaryState(
								readWriteModel.ApplicationIdContainer_LIGHTING_3A,
								122,
							),
							readWriteModel.NewBridgeAddress(12),
							readWriteModel.NewNetworkRoute(readWriteModel.NewNetworkProtocolControlInformation(2, 2), []readWriteModel.BridgeAddress{readWriteModel.NewBridgeAddress(13)}),
							58,
							cbusOptions,
						),
						readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPointToMultiPoint),
						cbusOptions,
					),
					nil,
					readWriteModel.NewAlpha(0),
					readWriteModel.RequestType_REQUEST_COMMAND,
					nil,
					nil,
					readWriteModel.RequestType_EMPTY,
					readWriteModel.NewRequestTermination(),
					cbusOptions,
				),
				requestContext,
				cbusOptions,
			),
			wantSupportsRead:      true,
			wantSupportsSubscribe: true,
		},
		{
			name: "direct status binary bridged 6 bridges",
			args: args{
				tag:            NewStatusTag([]readWriteModel.BridgeAddress{readWriteModel.NewBridgeAddress(12), readWriteModel.NewBridgeAddress(13), readWriteModel.NewBridgeAddress(14), readWriteModel.NewBridgeAddress(15), readWriteModel.NewBridgeAddress(16), readWriteModel.NewBridgeAddress(17)}, StatusRequestTypeBinaryState, nil, readWriteModel.ApplicationIdContainer_LIGHTING_3A, 1),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantCBusMessage: readWriteModel.NewCBusMessageToServer(
				readWriteModel.NewRequestCommand(
					readWriteModel.NewCBusCommandPointToPointToMultiPoint(
						readWriteModel.NewCBusPointToPointToMultiPointCommandStatus(
							readWriteModel.NewStatusRequestBinaryState(
								readWriteModel.ApplicationIdContainer_LIGHTING_3A,
								122,
							),
							readWriteModel.NewBridgeAddress(12),
							readWriteModel.NewNetworkRoute(readWriteModel.NewNetworkProtocolControlInformation(6, 6), []readWriteModel.BridgeAddress{readWriteModel.NewBridgeAddress(13), readWriteModel.NewBridgeAddress(14), readWriteModel.NewBridgeAddress(15), readWriteModel.NewBridgeAddress(16), readWriteModel.NewBridgeAddress(17)}),
							58,
							cbusOptions,
						),
						readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPointToMultiPoint),
						cbusOptions,
					),
					nil,
					readWriteModel.NewAlpha(0),
					readWriteModel.RequestType_REQUEST_COMMAND,
					nil,
					nil,
					readWriteModel.RequestType_EMPTY,
					readWriteModel.NewRequestTermination(),
					cbusOptions,
				),
				requestContext,
				cbusOptions,
			),
			wantSupportsRead:      true,
			wantSupportsSubscribe: true,
		},
		{
			name: "direct status binary bridged 7 bridges",
			args: args{
				tag:            NewStatusTag([]readWriteModel.BridgeAddress{readWriteModel.NewBridgeAddress(12), readWriteModel.NewBridgeAddress(13), readWriteModel.NewBridgeAddress(14), readWriteModel.NewBridgeAddress(15), readWriteModel.NewBridgeAddress(16), readWriteModel.NewBridgeAddress(17), readWriteModel.NewBridgeAddress(18)}, StatusRequestTypeBinaryState, nil, readWriteModel.ApplicationIdContainer_LIGHTING_3A, 1),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantErr: true,
		},
		{
			name: "recall",
			args: args{
				tag:            NewCALRecallTag(readWriteModel.NewUnitAddress(1), nil, readWriteModel.Parameter_BAUD_RATE_SELECTOR, 1, 1),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantCBusMessage: readWriteModel.NewCBusMessageToServer(
				readWriteModel.NewRequestCommand(
					readWriteModel.NewCBusCommandPointToPoint(
						readWriteModel.NewCBusPointToPointCommandDirect(
							readWriteModel.NewUnitAddress(1),
							0,
							readWriteModel.NewCALDataRecall(
								readWriteModel.Parameter_BAUD_RATE_SELECTOR,
								1,
								readWriteModel.CALCommandTypeContainer_CALCommandRecall,
								nil,
								requestContext,
							),
							cbusOptions,
						),
						readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint),
						cbusOptions,
					),
					nil,
					readWriteModel.NewAlpha(0),
					readWriteModel.RequestType_REQUEST_COMMAND,
					nil,
					nil,
					readWriteModel.RequestType_EMPTY,
					readWriteModel.NewRequestTermination(),
					cbusOptions,
				),
				requestContext,
				cbusOptions,
			),
			wantSupportsRead:      true,
			wantSupportsSubscribe: false,
		},
		{
			name: "identify",
			args: args{
				tag:            NewCALIdentifyTag(readWriteModel.NewUnitAddress(1), nil, readWriteModel.Attribute_Manufacturer, 1),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantCBusMessage: readWriteModel.NewCBusMessageToServer(
				readWriteModel.NewRequestCommand(
					readWriteModel.NewCBusCommandPointToPoint(
						readWriteModel.NewCBusPointToPointCommandDirect(
							readWriteModel.NewUnitAddress(1),
							0,
							readWriteModel.NewCALDataIdentify(
								readWriteModel.Attribute_Manufacturer,
								readWriteModel.CALCommandTypeContainer_CALCommandIdentify,
								nil,
								requestContext,
							),
							cbusOptions,
						),
						readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint),
						cbusOptions,
					),
					nil,
					readWriteModel.NewAlpha(0),
					readWriteModel.RequestType_REQUEST_COMMAND,
					nil,
					nil,
					readWriteModel.RequestType_EMPTY,
					readWriteModel.NewRequestTermination(),
					cbusOptions,
				),
				requestContext,
				cbusOptions,
			),
			wantSupportsRead:      true,
			wantSupportsSubscribe: false,
		},
		{
			name: "getStatus",
			args: args{
				tag:            NewCALGetStatusTag(readWriteModel.NewUnitAddress(1), nil, readWriteModel.Parameter_BAUD_RATE_SELECTOR, 1, 1),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantCBusMessage: readWriteModel.NewCBusMessageToServer(
				readWriteModel.NewRequestCommand(
					readWriteModel.NewCBusCommandPointToPoint(
						readWriteModel.NewCBusPointToPointCommandDirect(
							readWriteModel.NewUnitAddress(1),
							0,
							readWriteModel.NewCALDataGetStatus(
								readWriteModel.Parameter_BAUD_RATE_SELECTOR,
								1,
								readWriteModel.CALCommandTypeContainer_CALCommandGetStatus,
								nil,
								requestContext,
							),
							cbusOptions,
						),
						readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint),
						cbusOptions,
					),
					nil,
					readWriteModel.NewAlpha(0),
					readWriteModel.RequestType_REQUEST_COMMAND,
					nil,
					nil,
					readWriteModel.RequestType_EMPTY,
					readWriteModel.NewRequestTermination(),
					cbusOptions,
				),
				requestContext,
				cbusOptions,
			),
			wantSupportsRead:      true,
			wantSupportsSubscribe: false,
		},
		// TODO: test sal free usage
		{
			name: "sal temperature command not found",
			args: args{
				tag:            NewSALTag(nil, readWriteModel.ApplicationIdContainer_TEMPERATURE_BROADCAST_19, "asd", 1),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantErr: true,
		},
		{
			name: "sal temperature command event not enough arguments",
			args: args{
				tag:            NewSALTag(nil, readWriteModel.ApplicationIdContainer_TEMPERATURE_BROADCAST_19, "BROADCAST_EVENT", 1),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantErr: true,
		},
		{
			name: "sal temperature command event wrong arguments",
			args: args{
				tag: NewSALTag(nil, readWriteModel.ApplicationIdContainer_TEMPERATURE_BROADCAST_19, "BROADCAST_EVENT", 1),
				value: spiValues.NewPlcList([]apiValues.PlcValue{
					spiValues.NewPlcSTRING("asd"),
					spiValues.NewPlcBYTE(3),
				}),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantErr: true,
		},
		{
			name: "sal temperature command event",
			args: args{
				tag: NewSALTag(nil, readWriteModel.ApplicationIdContainer_TEMPERATURE_BROADCAST_19, "BROADCAST_EVENT", 1),
				value: spiValues.NewPlcList([]apiValues.PlcValue{
					spiValues.NewPlcBYTE(2),
					spiValues.NewPlcBYTE(3),
				}),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantCBusMessage: readWriteModel.NewCBusMessageToServer(
				readWriteModel.NewRequestCommand(
					readWriteModel.NewCBusCommandPointToMultiPoint(
						readWriteModel.NewCBusPointToMultiPointCommandNormal(
							readWriteModel.ApplicationIdContainer_TEMPERATURE_BROADCAST_19,
							readWriteModel.NewSALDataTemperatureBroadcast(
								readWriteModel.NewTemperatureBroadcastData(
									readWriteModel.TemperatureBroadcastCommandTypeContainer_TemperatureBroadcastCommandSetBroadcastEvent1_2Bytes,
									2,
									3,
								),
								nil,
							),
							0,
							cbusOptions,
						),
						readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint),
						cbusOptions,
					),
					nil,
					readWriteModel.NewAlpha(0),
					readWriteModel.RequestType_REQUEST_COMMAND,
					nil,
					nil,
					readWriteModel.RequestType_EMPTY,
					readWriteModel.NewRequestTermination(),
					cbusOptions,
				),
				requestContext,
				cbusOptions,
			),
			wantSupportsRead:      false,
			wantSupportsWrite:     true,
			wantSupportsSubscribe: false,
		},
		// TODO: test sal room control system
		{
			name: "sal lighting command not found",
			args: args{
				tag:            NewSALTag(nil, readWriteModel.ApplicationIdContainer_LIGHTING_3A, "asd", 1),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantErr: true,
		},
		{
			name: "sal lighting command off not enough arguments",
			args: args{
				tag:            NewSALTag(nil, readWriteModel.ApplicationIdContainer_LIGHTING_3A, "OFF", 1),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantErr: true,
		},
		{
			name: "sal lighting command on not enough arguments",
			args: args{
				tag:            NewSALTag(nil, readWriteModel.ApplicationIdContainer_LIGHTING_3A, "ON", 1),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantErr: true,
		},
		{
			name: "sal lighting command ramp to level not enough arguments",
			args: args{
				tag:            NewSALTag(nil, readWriteModel.ApplicationIdContainer_LIGHTING_3A, "RAMP_TO_LEVEL", 1),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantErr: true,
		},
		{
			name: "sal lighting command terminate ramp not enough arguments",
			args: args{
				tag:            NewSALTag(nil, readWriteModel.ApplicationIdContainer_LIGHTING_3A, "TERMINATE_RAMP", 1),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantErr: true,
		},
		// TODO: add support for label
		{
			name: "sal lighting command off wrong arguments",
			args: args{
				tag: NewSALTag(nil, readWriteModel.ApplicationIdContainer_LIGHTING_3A, "OFF", 1),
				value: spiValues.NewPlcList([]apiValues.PlcValue{
					spiValues.NewPlcSTRING("asd"),
					spiValues.NewPlcBYTE(3),
				}),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantErr: true,
		},
		{
			name: "sal lighting command on wrong arguments",
			args: args{
				tag: NewSALTag(nil, readWriteModel.ApplicationIdContainer_LIGHTING_3A, "ON", 1),
				value: spiValues.NewPlcList([]apiValues.PlcValue{
					spiValues.NewPlcSTRING("asd"),
					spiValues.NewPlcBYTE(3),
				}),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantErr: true,
		},
		{
			name: "sal lighting command ramp to level wrong arguments",
			args: args{
				tag: NewSALTag(nil, readWriteModel.ApplicationIdContainer_LIGHTING_3A, "RAMP_TO_LEVEL", 1),
				value: spiValues.NewPlcList([]apiValues.PlcValue{
					spiValues.NewPlcSTRING("asd"),
					spiValues.NewPlcBYTE(3),
				}),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantErr: true,
		},
		{
			name: "sal lighting command terminate ramp wrong arguments",
			args: args{
				tag: NewSALTag(nil, readWriteModel.ApplicationIdContainer_LIGHTING_3A, "TERMINATE_RAMP", 1),
				value: spiValues.NewPlcList([]apiValues.PlcValue{
					spiValues.NewPlcSTRING("asd"),
					spiValues.NewPlcBYTE(3),
				}),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantErr: true,
		},
		// TODO: implement label support command check
		{
			name: "sal lighting command on",
			args: args{
				tag:            NewSALTag(nil, readWriteModel.ApplicationIdContainer_LIGHTING_3A, "ON", 1),
				value:          spiValues.NewPlcBYTE(2),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantCBusMessage: readWriteModel.NewCBusMessageToServer(
				readWriteModel.NewRequestCommand(
					readWriteModel.NewCBusCommandPointToMultiPoint(
						readWriteModel.NewCBusPointToMultiPointCommandNormal(
							readWriteModel.ApplicationIdContainer_LIGHTING_3A,
							readWriteModel.NewSALDataLighting(
								readWriteModel.NewLightingDataOn(
									2,
									readWriteModel.LightingCommandTypeContainer_LightingCommandOn,
								),
								nil,
							),
							0,
							cbusOptions,
						),
						readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint),
						cbusOptions,
					),
					nil,
					readWriteModel.NewAlpha(0),
					readWriteModel.RequestType_REQUEST_COMMAND,
					nil,
					nil,
					readWriteModel.RequestType_EMPTY,
					readWriteModel.NewRequestTermination(),
					cbusOptions,
				),
				requestContext,
				cbusOptions,
			),
			wantSupportsRead:      false,
			wantSupportsWrite:     true,
			wantSupportsSubscribe: false,
		},
		{
			name: "sal lighting command off",
			args: args{
				tag:            NewSALTag(nil, readWriteModel.ApplicationIdContainer_LIGHTING_3A, "OFF", 1),
				value:          spiValues.NewPlcBYTE(2),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantCBusMessage: readWriteModel.NewCBusMessageToServer(
				readWriteModel.NewRequestCommand(
					readWriteModel.NewCBusCommandPointToMultiPoint(
						readWriteModel.NewCBusPointToMultiPointCommandNormal(
							readWriteModel.ApplicationIdContainer_LIGHTING_3A,
							readWriteModel.NewSALDataLighting(
								readWriteModel.NewLightingDataOff(
									2,
									readWriteModel.LightingCommandTypeContainer_LightingCommandOff,
								),
								nil,
							),
							0,
							cbusOptions,
						),
						readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint),
						cbusOptions,
					),
					nil,
					readWriteModel.NewAlpha(0),
					readWriteModel.RequestType_REQUEST_COMMAND,
					nil,
					nil,
					readWriteModel.RequestType_EMPTY,
					readWriteModel.NewRequestTermination(),
					cbusOptions,
				),
				requestContext,
				cbusOptions,
			),
			wantSupportsRead:      false,
			wantSupportsWrite:     true,
			wantSupportsSubscribe: false,
		},
		{
			name: "sal lighting command ramp to level",
			args: args{
				tag: NewSALTag(nil, readWriteModel.ApplicationIdContainer_LIGHTING_3A, "RAMP_TO_LEVEL", 1),
				value: spiValues.NewPlcList([]apiValues.PlcValue{
					spiValues.NewPlcSTRING("4Second"),
					spiValues.NewPlcBYTE(2),
					spiValues.NewPlcBYTE(3),
				}),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantCBusMessage: readWriteModel.NewCBusMessageToServer(
				readWriteModel.NewRequestCommand(
					readWriteModel.NewCBusCommandPointToMultiPoint(
						readWriteModel.NewCBusPointToMultiPointCommandNormal(
							readWriteModel.ApplicationIdContainer_LIGHTING_3A,
							readWriteModel.NewSALDataLighting(
								readWriteModel.NewLightingDataRampToLevel(
									2,
									3,
									readWriteModel.LightingCommandTypeContainer_LightingCommandRampToLevel_4Second,
								),
								nil,
							),
							0,
							cbusOptions,
						),
						readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint),
						cbusOptions,
					),
					nil,
					readWriteModel.NewAlpha(0),
					readWriteModel.RequestType_REQUEST_COMMAND,
					nil,
					nil,
					readWriteModel.RequestType_EMPTY,
					readWriteModel.NewRequestTermination(),
					cbusOptions,
				),
				requestContext,
				cbusOptions,
			),
			wantSupportsRead:      false,
			wantSupportsWrite:     true,
			wantSupportsSubscribe: false,
		},
		{
			name: "sal lighting command terminate ramp",
			args: args{
				tag:            NewSALTag(nil, readWriteModel.ApplicationIdContainer_LIGHTING_3A, "TERMINATE_RAMP", 1),
				value:          spiValues.NewPlcBYTE(2),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantCBusMessage: readWriteModel.NewCBusMessageToServer(
				readWriteModel.NewRequestCommand(
					readWriteModel.NewCBusCommandPointToMultiPoint(
						readWriteModel.NewCBusPointToMultiPointCommandNormal(
							readWriteModel.ApplicationIdContainer_LIGHTING_3A,
							readWriteModel.NewSALDataLighting(
								readWriteModel.NewLightingDataTerminateRamp(
									2,
									readWriteModel.LightingCommandTypeContainer_LightingCommandTerminateRamp,
								),
								nil,
							),
							0,
							cbusOptions,
						),
						readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint),
						cbusOptions,
					),
					nil,
					readWriteModel.NewAlpha(0),
					readWriteModel.RequestType_REQUEST_COMMAND,
					nil,
					nil,
					readWriteModel.RequestType_EMPTY,
					readWriteModel.NewRequestTermination(),
					cbusOptions,
				),
				requestContext,
				cbusOptions,
			),
			wantSupportsRead:      false,
			wantSupportsWrite:     true,
			wantSupportsSubscribe: false,
		},
		// TODO: implement label support command
		{
			name: "sal ventilation command on", // Note: is based on lighting, so we just test "on" here
			args: args{
				tag:            NewSALTag(nil, readWriteModel.ApplicationIdContainer_VENTILATION_70, "ON", 1),
				value:          spiValues.NewPlcBYTE(2),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantCBusMessage: readWriteModel.NewCBusMessageToServer(
				readWriteModel.NewRequestCommand(
					readWriteModel.NewCBusCommandPointToMultiPoint(
						readWriteModel.NewCBusPointToMultiPointCommandNormal(
							readWriteModel.ApplicationIdContainer_VENTILATION_70,
							readWriteModel.NewSALDataLighting(
								readWriteModel.NewLightingDataOn(
									2,
									readWriteModel.LightingCommandTypeContainer_LightingCommandOn,
								),
								nil,
							),
							0,
							cbusOptions,
						),
						readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint),
						cbusOptions,
					),
					nil,
					readWriteModel.NewAlpha(0),
					readWriteModel.RequestType_REQUEST_COMMAND,
					nil,
					nil,
					readWriteModel.RequestType_EMPTY,
					readWriteModel.NewRequestTermination(),
					cbusOptions,
				),
				requestContext,
				cbusOptions,
			),
			wantSupportsRead:      false,
			wantSupportsWrite:     true,
			wantSupportsSubscribe: false,
		},
		{
			name: "sal irrigation control command on", // Note: is based on lighting, so we just test "on" here
			args: args{
				tag:            NewSALTag(nil, readWriteModel.ApplicationIdContainer_IRRIGATION_CONTROL_71, "ON", 1),
				value:          spiValues.NewPlcBYTE(2),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantCBusMessage: readWriteModel.NewCBusMessageToServer(
				readWriteModel.NewRequestCommand(
					readWriteModel.NewCBusCommandPointToMultiPoint(
						readWriteModel.NewCBusPointToMultiPointCommandNormal(
							readWriteModel.ApplicationIdContainer_IRRIGATION_CONTROL_71,
							readWriteModel.NewSALDataLighting(
								readWriteModel.NewLightingDataOn(
									2,
									readWriteModel.LightingCommandTypeContainer_LightingCommandOn,
								),
								nil,
							),
							0,
							cbusOptions,
						),
						readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint),
						cbusOptions,
					),
					nil,
					readWriteModel.NewAlpha(0),
					readWriteModel.RequestType_REQUEST_COMMAND,
					nil,
					nil,
					readWriteModel.RequestType_EMPTY,
					readWriteModel.NewRequestTermination(),
					cbusOptions,
				),
				requestContext,
				cbusOptions,
			),
			wantSupportsRead:      false,
			wantSupportsWrite:     true,
			wantSupportsSubscribe: false,
		},
		{
			name: "sal pools spas ponds fountains control command on", // Note: is based on lighting, so we just test "on" here
			args: args{
				tag:            NewSALTag(nil, readWriteModel.ApplicationIdContainer_POOLS_SPAS_PONDS_FOUNTAINS_CONTROL_72, "ON", 1),
				value:          spiValues.NewPlcBYTE(2),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantCBusMessage: readWriteModel.NewCBusMessageToServer(
				readWriteModel.NewRequestCommand(
					readWriteModel.NewCBusCommandPointToMultiPoint(
						readWriteModel.NewCBusPointToMultiPointCommandNormal(
							readWriteModel.ApplicationIdContainer_POOLS_SPAS_PONDS_FOUNTAINS_CONTROL_72,
							readWriteModel.NewSALDataLighting(
								readWriteModel.NewLightingDataOn(
									2,
									readWriteModel.LightingCommandTypeContainer_LightingCommandOn,
								),
								nil,
							),
							0,
							cbusOptions,
						),
						readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint),
						cbusOptions,
					),
					nil,
					readWriteModel.NewAlpha(0),
					readWriteModel.RequestType_REQUEST_COMMAND,
					nil,
					nil,
					readWriteModel.RequestType_EMPTY,
					readWriteModel.NewRequestTermination(),
					cbusOptions,
				),
				requestContext,
				cbusOptions,
			),
			wantSupportsRead:      false,
			wantSupportsWrite:     true,
			wantSupportsSubscribe: false,
		},
		{
			name: "sal heating on", // Note: is based on lighting, so we just test "on" here
			args: args{
				tag:            NewSALTag(nil, readWriteModel.ApplicationIdContainer_HEATING_88, "ON", 1),
				value:          spiValues.NewPlcBYTE(2),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantCBusMessage: readWriteModel.NewCBusMessageToServer(
				readWriteModel.NewRequestCommand(
					readWriteModel.NewCBusCommandPointToMultiPoint(
						readWriteModel.NewCBusPointToMultiPointCommandNormal(
							readWriteModel.ApplicationIdContainer_HEATING_88,
							readWriteModel.NewSALDataLighting(
								readWriteModel.NewLightingDataOn(
									2,
									readWriteModel.LightingCommandTypeContainer_LightingCommandOn,
								),
								nil,
							),
							0,
							cbusOptions,
						),
						readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint),
						cbusOptions,
					),
					nil,
					readWriteModel.NewAlpha(0),
					readWriteModel.RequestType_REQUEST_COMMAND,
					nil,
					nil,
					readWriteModel.RequestType_EMPTY,
					readWriteModel.NewRequestTermination(),
					cbusOptions,
				),
				requestContext,
				cbusOptions,
			),
			wantSupportsRead:      false,
			wantSupportsWrite:     true,
			wantSupportsSubscribe: false,
		},
		{
			name: "sal audio and video on", // Note: is based on lighting, so we just test "on" here
			args: args{
				tag:            NewSALTag(nil, readWriteModel.ApplicationIdContainer_AUDIO_AND_VIDEO_CD, "ON", 1),
				value:          spiValues.NewPlcBYTE(2),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantCBusMessage: readWriteModel.NewCBusMessageToServer(
				readWriteModel.NewRequestCommand(
					readWriteModel.NewCBusCommandPointToMultiPoint(
						readWriteModel.NewCBusPointToMultiPointCommandNormal(
							readWriteModel.ApplicationIdContainer_AUDIO_AND_VIDEO_CD,
							readWriteModel.NewSALDataLighting(
								readWriteModel.NewLightingDataOn(
									2,
									readWriteModel.LightingCommandTypeContainer_LightingCommandOn,
								),
								nil,
							),
							0,
							cbusOptions,
						),
						readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint),
						cbusOptions,
					),
					nil,
					readWriteModel.NewAlpha(0),
					readWriteModel.RequestType_REQUEST_COMMAND,
					nil,
					nil,
					readWriteModel.RequestType_EMPTY,
					readWriteModel.NewRequestTermination(),
					cbusOptions,
				),
				requestContext,
				cbusOptions,
			),
			wantSupportsRead:      false,
			wantSupportsWrite:     true,
			wantSupportsSubscribe: false,
		},
		{
			name: "sal hvac actuator on", // Note: is based on lighting, so we just test "on" here
			args: args{
				tag:            NewSALTag(nil, readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74, "ON", 1),
				value:          spiValues.NewPlcBYTE(2),
				alphaGenerator: &AlphaGenerator{},
				messageCodec: &MessageCodec{
					cbusOptions:    cbusOptions,
					requestContext: requestContext,
				},
			},
			wantCBusMessage: readWriteModel.NewCBusMessageToServer(
				readWriteModel.NewRequestCommand(
					readWriteModel.NewCBusCommandPointToMultiPoint(
						readWriteModel.NewCBusPointToMultiPointCommandNormal(
							readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_74,
							readWriteModel.NewSALDataLighting(
								readWriteModel.NewLightingDataOn(
									2,
									readWriteModel.LightingCommandTypeContainer_LightingCommandOn,
								),
								nil,
							),
							0,
							cbusOptions,
						),
						readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint),
						cbusOptions,
					),
					nil,
					readWriteModel.NewAlpha(0),
					readWriteModel.RequestType_REQUEST_COMMAND,
					nil,
					nil,
					readWriteModel.RequestType_EMPTY,
					readWriteModel.NewRequestTermination(),
					cbusOptions,
				),
				requestContext,
				cbusOptions,
			),
			wantSupportsRead:      false,
			wantSupportsWrite:     true,
			wantSupportsSubscribe: false,
		},
		// TODO: implement air conditioning
		// TODO: implement trigger control
		// TODO: implement enable control
		// TODO: implement security
		// TODO: implement metering
		// TODO: implement access control
		// TODO: implement clock and timekeeping
		// TODO: implement telephony status and control
		// TODO: implement measurement
		// TODO: implement testing
		// TODO: implement media transport control
		// TODO: implement error reporting
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCBusMessage, gotSupportsRead, gotSupportsWrite, gotSupportsSubscribe, err := TagToCBusMessage(tt.args.tag, tt.args.value, tt.args.alphaGenerator, tt.args.messageCodec)
			if (err != nil) != tt.wantErr {
				t.Errorf("TagToCBusMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, gotCBusMessage, tt.wantCBusMessage) {
				gotBox := utils.BoxAnything("got", gotCBusMessage, 120)
				wantBox := utils.BoxAnything("want", tt.wantCBusMessage, 120)
				t.Errorf("TagToCBusMessage():\n%s", utils.NewAsciiBoxWriter().BoxSideBySide(gotBox, wantBox))
			}
			if gotSupportsRead != tt.wantSupportsRead {
				t.Errorf("TagToCBusMessage() gotSupportsRead = %v, want %v", gotSupportsRead, tt.wantSupportsRead)
			}
			if gotSupportsWrite != tt.wantSupportsWrite {
				t.Errorf("TagToCBusMessage() gotSupportsWrite = %v, want %v", gotSupportsWrite, tt.wantSupportsWrite)
			}
			if gotSupportsSubscribe != tt.wantSupportsSubscribe {
				t.Errorf("TagToCBusMessage() got SupportsSubscribe = %v, want %v", gotSupportsSubscribe, tt.wantSupportsSubscribe)
			}
		})
	}
}

func Test_producePointToPointCommand(t *testing.T) {
	type args struct {
		unitAddress     readWriteModel.UnitAddress
		bridgeAddresses []readWriteModel.BridgeAddress
		calData         readWriteModel.CALData
		cbusOptions     readWriteModel.CBusOptions
	}
	tests := []struct {
		name    string
		args    args
		want    readWriteModel.CBusCommand
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "no bridge",
			want: readWriteModel.NewCBusCommandPointToPoint(
				readWriteModel.NewCBusPointToPointCommandDirect(
					nil,
					0,
					nil,
					nil,
				),
				readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint),
				nil,
			),
			wantErr: assert.NoError,
		},
		{
			name: "one bridge",
			args: args{
				bridgeAddresses: []readWriteModel.BridgeAddress{
					readWriteModel.NewBridgeAddress(1),
				},
			},
			want: readWriteModel.NewCBusCommandPointToPoint(
				readWriteModel.NewCBusPointToPointCommandIndirect(
					readWriteModel.NewBridgeAddress(1),
					readWriteModel.NewNetworkRoute(
						readWriteModel.NewNetworkProtocolControlInformation(1, 1),
						[]readWriteModel.BridgeAddress{},
					),
					nil,
					0,
					nil,
					nil,
				),
				readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint),
				nil,
			),
			wantErr: assert.NoError,
		},
		{
			name: "6 bridges",
			args: args{
				bridgeAddresses: []readWriteModel.BridgeAddress{
					readWriteModel.NewBridgeAddress(1),
					readWriteModel.NewBridgeAddress(2),
					readWriteModel.NewBridgeAddress(3),
					readWriteModel.NewBridgeAddress(4),
					readWriteModel.NewBridgeAddress(5),
					readWriteModel.NewBridgeAddress(6),
				},
			},
			want: readWriteModel.NewCBusCommandPointToPoint(
				readWriteModel.NewCBusPointToPointCommandIndirect(
					readWriteModel.NewBridgeAddress(1),
					readWriteModel.NewNetworkRoute(
						readWriteModel.NewNetworkProtocolControlInformation(6, 6),
						[]readWriteModel.BridgeAddress{
							readWriteModel.NewBridgeAddress(2),
							readWriteModel.NewBridgeAddress(3),
							readWriteModel.NewBridgeAddress(4),
							readWriteModel.NewBridgeAddress(5),
							readWriteModel.NewBridgeAddress(6),
						},
					),
					nil,
					0,
					nil,
					nil,
				), readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint),
				nil,
			),
			wantErr: assert.NoError,
		},
		{
			name: "7 bridges",

			args: args{
				bridgeAddresses: []readWriteModel.BridgeAddress{
					readWriteModel.NewBridgeAddress(1),
					readWriteModel.NewBridgeAddress(2),
					readWriteModel.NewBridgeAddress(3),
					readWriteModel.NewBridgeAddress(4),
					readWriteModel.NewBridgeAddress(5),
					readWriteModel.NewBridgeAddress(6),
					readWriteModel.NewBridgeAddress(7),
				},
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := producePointToPointCommand(tt.args.unitAddress, tt.args.bridgeAddresses, tt.args.calData, tt.args.cbusOptions)
			if !tt.wantErr(t, err, fmt.Sprintf("producePointToPointCommand(%v, %v, %v, %v)", tt.args.unitAddress, tt.args.bridgeAddresses, tt.args.calData, tt.args.cbusOptions)) {
				return
			}
			assert.Equalf(t, tt.want, got, "producePointToPointCommand(%v, %v, %v, %v)", tt.args.unitAddress, tt.args.bridgeAddresses, tt.args.calData, tt.args.cbusOptions)
		})
	}
}

func Test_producePointToMultiPointCommandStatus(t *testing.T) {
	type args struct {
		bridgeAddresses []readWriteModel.BridgeAddress
		application     readWriteModel.ApplicationIdContainer
		statusRequest   readWriteModel.StatusRequest
		cbusOptions     readWriteModel.CBusOptions
	}
	tests := []struct {
		name    string
		args    args
		want    readWriteModel.CBusCommand
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "no bridge",
			want: readWriteModel.NewCBusCommandPointToMultiPoint(
				readWriteModel.NewCBusPointToMultiPointCommandStatus(nil, 0, nil),
				readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToMultiPoint),
				nil,
			),
			wantErr: assert.NoError,
		},
		{
			name: "one bridge",
			args: args{
				bridgeAddresses: []readWriteModel.BridgeAddress{
					readWriteModel.NewBridgeAddress(1),
				},
			},
			want: readWriteModel.NewCBusCommandPointToPointToMultiPoint(
				readWriteModel.NewCBusPointToPointToMultiPointCommandStatus(
					nil,
					readWriteModel.NewBridgeAddress(1),
					readWriteModel.NewNetworkRoute(
						readWriteModel.NewNetworkProtocolControlInformation(1, 1),
						[]readWriteModel.BridgeAddress{},
					),
					0,
					nil,
				),
				readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPointToMultiPoint),
				nil,
			),
			wantErr: assert.NoError,
		},
		{
			name: "6 bridges",
			args: args{
				bridgeAddresses: []readWriteModel.BridgeAddress{
					readWriteModel.NewBridgeAddress(1),
					readWriteModel.NewBridgeAddress(2),
					readWriteModel.NewBridgeAddress(3),
					readWriteModel.NewBridgeAddress(4),
					readWriteModel.NewBridgeAddress(5),
					readWriteModel.NewBridgeAddress(6),
				},
			},
			want: readWriteModel.NewCBusCommandPointToPointToMultiPoint(
				readWriteModel.NewCBusPointToPointToMultiPointCommandStatus(
					nil,
					readWriteModel.NewBridgeAddress(1),
					readWriteModel.NewNetworkRoute(
						readWriteModel.NewNetworkProtocolControlInformation(6, 6),
						[]readWriteModel.BridgeAddress{
							readWriteModel.NewBridgeAddress(2),
							readWriteModel.NewBridgeAddress(3),
							readWriteModel.NewBridgeAddress(4),
							readWriteModel.NewBridgeAddress(5),
							readWriteModel.NewBridgeAddress(6),
						},
					),
					0,
					nil,
				),
				readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPointToMultiPoint),
				nil,
			),
			wantErr: assert.NoError,
		},
		{
			name: "7 bridges",

			args: args{
				bridgeAddresses: []readWriteModel.BridgeAddress{
					readWriteModel.NewBridgeAddress(1),
					readWriteModel.NewBridgeAddress(2),
					readWriteModel.NewBridgeAddress(3),
					readWriteModel.NewBridgeAddress(4),
					readWriteModel.NewBridgeAddress(5),
					readWriteModel.NewBridgeAddress(6),
					readWriteModel.NewBridgeAddress(7),
				},
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := producePointToMultiPointCommandStatus(tt.args.bridgeAddresses, tt.args.application, tt.args.statusRequest, tt.args.cbusOptions)
			if !tt.wantErr(t, err, fmt.Sprintf("producePointToMultiPointCommandStatus(%v, %v, %v, %v)", tt.args.bridgeAddresses, tt.args.application, tt.args.statusRequest, tt.args.cbusOptions)) {
				return
			}
			assert.Equalf(t, tt.want, got, "producePointToMultiPointCommandStatus(%v, %v, %v, %v)", tt.args.bridgeAddresses, tt.args.application, tt.args.statusRequest, tt.args.cbusOptions)
		})
	}
}

func Test_producePointToMultiPointCommandNormal(t *testing.T) {
	type args struct {
		bridgeAddresses []readWriteModel.BridgeAddress
		application     readWriteModel.ApplicationIdContainer
		salData         readWriteModel.SALData
		cbusOptions     readWriteModel.CBusOptions
	}
	tests := []struct {
		name    string
		args    args
		want    readWriteModel.CBusCommand
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "no bridge",
			want: readWriteModel.NewCBusCommandPointToMultiPoint(
				readWriteModel.NewCBusPointToMultiPointCommandNormal(0, nil, 0, nil),
				readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint),
				nil,
			),
			wantErr: assert.NoError,
		},
		{
			name: "one bridge",
			args: args{
				bridgeAddresses: []readWriteModel.BridgeAddress{
					readWriteModel.NewBridgeAddress(1),
				},
			},
			want: readWriteModel.NewCBusCommandPointToPointToMultiPoint(
				readWriteModel.NewCBusPointToPointToMultiPointCommandNormal(
					0,
					nil,
					readWriteModel.NewBridgeAddress(1),
					readWriteModel.NewNetworkRoute(
						readWriteModel.NewNetworkProtocolControlInformation(1, 1),
						[]readWriteModel.BridgeAddress{},
					),
					0,
					nil,
				),
				readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPointToMultiPoint),
				nil,
			),
			wantErr: assert.NoError,
		},
		{
			name: "6 bridges",
			args: args{
				bridgeAddresses: []readWriteModel.BridgeAddress{
					readWriteModel.NewBridgeAddress(1),
					readWriteModel.NewBridgeAddress(2),
					readWriteModel.NewBridgeAddress(3),
					readWriteModel.NewBridgeAddress(4),
					readWriteModel.NewBridgeAddress(5),
					readWriteModel.NewBridgeAddress(6),
				},
			},
			want: readWriteModel.NewCBusCommandPointToPointToMultiPoint(
				readWriteModel.NewCBusPointToPointToMultiPointCommandNormal(
					0,
					nil,
					readWriteModel.NewBridgeAddress(1),
					readWriteModel.NewNetworkRoute(
						readWriteModel.NewNetworkProtocolControlInformation(6, 6),
						[]readWriteModel.BridgeAddress{
							readWriteModel.NewBridgeAddress(2),
							readWriteModel.NewBridgeAddress(3),
							readWriteModel.NewBridgeAddress(4),
							readWriteModel.NewBridgeAddress(5),
							readWriteModel.NewBridgeAddress(6),
						},
					),
					0,
					nil,
				),
				readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPointToMultiPoint),
				nil,
			),
			wantErr: assert.NoError,
		},
		{
			name: "7 bridges",

			args: args{
				bridgeAddresses: []readWriteModel.BridgeAddress{
					readWriteModel.NewBridgeAddress(1),
					readWriteModel.NewBridgeAddress(2),
					readWriteModel.NewBridgeAddress(3),
					readWriteModel.NewBridgeAddress(4),
					readWriteModel.NewBridgeAddress(5),
					readWriteModel.NewBridgeAddress(6),
					readWriteModel.NewBridgeAddress(7),
				},
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := producePointToMultiPointCommandNormal(tt.args.bridgeAddresses, tt.args.application, tt.args.salData, tt.args.cbusOptions)
			if !tt.wantErr(t, err, fmt.Sprintf("producePointToMultiPointCommandNormal(%v, %v, %v, %v)", tt.args.bridgeAddresses, tt.args.application, tt.args.salData, tt.args.cbusOptions)) {
				return
			}
			assert.Equalf(t, tt.want, got, "producePointToMultiPointCommandNormal(%v, %v, %v, %v)", tt.args.bridgeAddresses, tt.args.application, tt.args.salData, tt.args.cbusOptions)
		})
	}
}

func TestMapEncodedReply(t *testing.T) {
	type args struct {
		transaction     transactions.RequestTransaction
		encodedReply    readWriteModel.EncodedReply
		tagName         string
		addResponseCode func(t *testing.T) func(name string, responseCode apiModel.PlcResponseCode)
		addPlcValue     func(t *testing.T) func(name string, plcValue apiValues.PlcValue)
	}
	tests := []struct {
		name    string
		args    args
		setup   func(t *testing.T, args *args)
		wantErr bool
	}{
		{
			name: "empty input",
			args: args{
				encodedReply: nil,
				tagName:      "",
				addResponseCode: func(t *testing.T) func(name string, responseCode apiModel.PlcResponseCode) {
					return func(name string, responseCode apiModel.PlcResponseCode) {
						// NO-OP
					}
				},
				addPlcValue: func(t *testing.T) func(name string, plcValue apiValues.PlcValue) {
					return func(name string, plcValue apiValues.PlcValue) {
						// NO-OP
					}
				},
			},
			setup: func(t *testing.T, args *args) {
				executor := pool.NewFixedSizeExecutor(10, 50, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				executor.Start()
				transactionManager := transactions.NewRequestTransactionManager(
					1,
					options.WithCustomLogger(testutils.ProduceTestingLogger(t)),
					transactions.WithCustomExecutor(executor),
				)
				t.Cleanup(func() {
					assert.NoError(t, transactionManager.Close())
				})
				transaction := transactionManager.StartTransaction()
				t.Logf("Submitting No-Op to transaction %v", transaction)
				transaction.Submit(func(transaction transactions.RequestTransaction) {
					// NO-OP
				})
				args.transaction = transaction
			},
		},
		{
			name: "CALDataStatus",
			args: args{
				encodedReply: func() readWriteModel.EncodedReplyCALReply {
					statusBytes := []readWriteModel.StatusByte{
						readWriteModel.NewStatusByte(readWriteModel.GAVState_ON, readWriteModel.GAVState_ERROR, readWriteModel.GAVState_DOES_NOT_EXIST, readWriteModel.GAVState_OFF),
						readWriteModel.NewStatusByte(readWriteModel.GAVState_ON, readWriteModel.GAVState_ERROR, readWriteModel.GAVState_DOES_NOT_EXIST, readWriteModel.GAVState_OFF),
					}
					calDataStatus := readWriteModel.NewCALDataStatus(readWriteModel.ApplicationIdContainer_LIGHTING_3A, 0, statusBytes, readWriteModel.CALCommandTypeContainer_CALCommandStatus_0Bytes, nil, nil)
					calReplyShort := readWriteModel.NewCALReplyShort(0, calDataStatus, nil, nil)
					return readWriteModel.NewEncodedReplyCALReply(calReplyShort, 0, nil, nil)
				}(),
				tagName: "someTag",
				addResponseCode: func(t *testing.T) func(string, apiModel.PlcResponseCode) {
					codes := map[string]apiModel.PlcResponseCode{
						"someTag": apiModel.PlcResponseCode_OK,
					}
					return func(name string, responseCode apiModel.PlcResponseCode) {
						if code, ok := codes[name]; ok {
							assert.Equal(t, code, responseCode)
						} else {
							t.Errorf("code for %s not found", name)
						}
					}
				},
				addPlcValue: func(t *testing.T) func(string, apiValues.PlcValue) {
					values := map[string]apiValues.PlcValue{
						"someTag": spiValues.NewPlcStruct(
							map[string]apiValues.PlcValue{
								"application": spiValues.NewPlcSTRING("LIGHTING_3A"),
								"blockStart":  spiValues.NewPlcBYTE(0x0),
								"values": spiValues.NewPlcList([]apiValues.PlcValue{
									spiValues.NewPlcSTRING("OFF"),
									spiValues.NewPlcSTRING("DOES_NOT_EXIST"),
									spiValues.NewPlcSTRING("ERROR"),
									spiValues.NewPlcSTRING("ON"),
									spiValues.NewPlcSTRING("OFF"),
									spiValues.NewPlcSTRING("DOES_NOT_EXIST"),
									spiValues.NewPlcSTRING("ERROR"),
									spiValues.NewPlcSTRING("ON"),
								}),
							},
						),
					}
					return func(name string, plcValue apiValues.PlcValue) {
						if value, ok := values[name]; ok {
							assert.Equal(t, value, plcValue)
						} else {
							t.Errorf("value for %s not found", name)
						}
					}
				},
			},
			setup: func(t *testing.T, args *args) {
				executor := pool.NewFixedSizeExecutor(10, 50, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				executor.Start()
				transactionManager := transactions.NewRequestTransactionManager(
					1,
					options.WithCustomLogger(testutils.ProduceTestingLogger(t)),
					transactions.WithCustomExecutor(executor),
				)
				t.Cleanup(func() {
					assert.NoError(t, transactionManager.Close())
				})
				transaction := transactionManager.StartTransaction()
				t.Logf("Submitting No-Op to transaction %v", transaction)
				transaction.Submit(func(transaction transactions.RequestTransaction) {
					// NO-OP
					t.Log("No op-ing")
				})
				t.Logf("Submitted to transaction %v", transaction)
				args.transaction = transaction
			},
		},
		{
			name: "CALDataStatusExtendedExactly (binary)",
			args: args{
				encodedReply: func() readWriteModel.EncodedReplyCALReply {
					statusBytes := []readWriteModel.StatusByte{
						readWriteModel.NewStatusByte(readWriteModel.GAVState_ON, readWriteModel.GAVState_ERROR, readWriteModel.GAVState_DOES_NOT_EXIST, readWriteModel.GAVState_OFF),
						readWriteModel.NewStatusByte(readWriteModel.GAVState_ON, readWriteModel.GAVState_ERROR, readWriteModel.GAVState_DOES_NOT_EXIST, readWriteModel.GAVState_OFF),
					}
					calDataStatus := readWriteModel.NewCALDataStatusExtended(readWriteModel.StatusCoding_BINARY_BY_THIS_SERIAL_INTERFACE, readWriteModel.ApplicationIdContainer_LIGHTING_3A, 0, statusBytes, nil, readWriteModel.CALCommandTypeContainer_CALCommandStatus_0Bytes, nil, nil)
					calReplyLong := readWriteModel.NewCALReplyLong(0, readWriteModel.NewUnitAddress(1), readWriteModel.NewBridgeAddress(2), readWriteModel.NewSerialInterfaceAddress(3), nil, nil, 0, calDataStatus, nil, nil)
					return readWriteModel.NewEncodedReplyCALReply(calReplyLong, 0, nil, nil)
				}(),
				tagName: "someTag",
				addResponseCode: func(t *testing.T) func(string, apiModel.PlcResponseCode) {
					codes := map[string]apiModel.PlcResponseCode{
						"someTag": apiModel.PlcResponseCode_OK,
					}
					return func(name string, responseCode apiModel.PlcResponseCode) {
						if code, ok := codes[name]; ok {
							assert.Equal(t, code, responseCode)
						} else {
							t.Errorf("code for %s not found", name)
						}
					}
				},
				addPlcValue: func(t *testing.T) func(string, apiValues.PlcValue) {
					values := map[string]apiValues.PlcValue{
						"someTag": spiValues.NewPlcStruct(
							map[string]apiValues.PlcValue{
								"application": spiValues.NewPlcSTRING("LIGHTING_3A"),
								"blockStart":  spiValues.NewPlcBYTE(0x0),
								"values": spiValues.NewPlcList([]apiValues.PlcValue{
									spiValues.NewPlcSTRING("OFF"),
									spiValues.NewPlcSTRING("DOES_NOT_EXIST"),
									spiValues.NewPlcSTRING("ERROR"),
									spiValues.NewPlcSTRING("ON"),
									spiValues.NewPlcSTRING("OFF"),
									spiValues.NewPlcSTRING("DOES_NOT_EXIST"),
									spiValues.NewPlcSTRING("ERROR"),
									spiValues.NewPlcSTRING("ON"),
								}),
							},
						),
					}
					return func(name string, plcValue apiValues.PlcValue) {
						if value, ok := values[name]; ok {
							assert.Equal(t, value, plcValue)
						} else {
							t.Errorf("value for %s not found", name)
						}
					}
				},
			},
			setup: func(t *testing.T, args *args) {
				executor := pool.NewFixedSizeExecutor(10, 50, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				executor.Start()
				transactionManager := transactions.NewRequestTransactionManager(
					1,
					options.WithCustomLogger(testutils.ProduceTestingLogger(t)),
					transactions.WithCustomExecutor(executor),
				)
				t.Cleanup(func() {
					assert.NoError(t, transactionManager.Close())
				})
				transaction := transactionManager.StartTransaction()
				transaction.Submit(func(transaction transactions.RequestTransaction) {
					// NO-OP
				})
				args.transaction = transaction
			},
		},
		{
			name: "CALDataStatusExtendedExactly (level)",
			args: args{
				encodedReply: func() readWriteModel.EncodedReplyCALReply {
					levelInformation := []readWriteModel.LevelInformation{
						readWriteModel.NewLevelInformationNormal(readWriteModel.LevelInformationNibblePair_Value_A, readWriteModel.LevelInformationNibblePair_Value_F, 13),
						readWriteModel.NewLevelInformationAbsent(13),
						readWriteModel.NewLevelInformationCorrupted(13, 14, 15, 16, 17),
					}
					calDataStatus := readWriteModel.NewCALDataStatusExtended(readWriteModel.StatusCoding_LEVEL_BY_THIS_SERIAL_INTERFACE, readWriteModel.ApplicationIdContainer_LIGHTING_3A, 0, nil, levelInformation, readWriteModel.CALCommandTypeContainer_CALCommandStatus_0Bytes, nil, nil)
					calReplyLong := readWriteModel.NewCALReplyLong(0, readWriteModel.NewUnitAddress(1), readWriteModel.NewBridgeAddress(2), readWriteModel.NewSerialInterfaceAddress(3), nil, nil, 0, calDataStatus, nil, nil)
					return readWriteModel.NewEncodedReplyCALReply(calReplyLong, 0, nil, nil)
				}(),
				tagName: "someTag",
				addResponseCode: func(t *testing.T) func(string, apiModel.PlcResponseCode) {
					codes := map[string]apiModel.PlcResponseCode{
						"someTag": apiModel.PlcResponseCode_OK,
					}
					return func(name string, responseCode apiModel.PlcResponseCode) {
						if code, ok := codes[name]; ok {
							assert.Equal(t, code, responseCode)
						} else {
							t.Errorf("code for %s not found", name)
						}
					}
				},
				addPlcValue: func(t *testing.T) func(string, apiValues.PlcValue) {
					values := map[string]apiValues.PlcValue{
						"someTag": spiValues.NewPlcList(
							[]apiValues.PlcValue{
								spiValues.NewPlcUSINT(250),
								spiValues.NewPlcSTRING("is absent"),
								spiValues.NewPlcSTRING("corrupted"),
							},
						),
					}
					return func(name string, plcValue apiValues.PlcValue) {
						if value, ok := values[name]; ok {
							assert.Equal(t, value, plcValue)
						} else {
							t.Errorf("value for %s not found", name)
						}
					}
				},
			},
			setup: func(t *testing.T, args *args) {
				executor := pool.NewFixedSizeExecutor(10, 50, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				executor.Start()
				transactionManager := transactions.NewRequestTransactionManager(
					1,
					options.WithCustomLogger(testutils.ProduceTestingLogger(t)),
					transactions.WithCustomExecutor(executor),
				)
				t.Cleanup(func() {
					assert.NoError(t, transactionManager.Close())
				})
				transaction := transactionManager.StartTransaction()
				transaction.Submit(func(transaction transactions.RequestTransaction) {
					// NO-OP
				})
				args.transaction = transaction
			},
		},
		{
			name: "CALDataIdentifyReplyExactly (sense levels)",
			args: args{
				encodedReply: func() readWriteModel.EncodedReplyCALReply {
					command := readWriteModel.NewIdentifyReplyCommandCurrentSenseLevels([]byte{1, 2, 3, 4}, 4)
					calDataIdentify := readWriteModel.NewCALDataIdentifyReply(readWriteModel.Attribute_CurrentSenseLevels, command, readWriteModel.CALCommandTypeContainer_CALCommandIdentify, nil, nil)
					calReplyLong := readWriteModel.NewCALReplyLong(0, readWriteModel.NewUnitAddress(1), readWriteModel.NewBridgeAddress(2), readWriteModel.NewSerialInterfaceAddress(3), nil, nil, 0, calDataIdentify, nil, nil)
					return readWriteModel.NewEncodedReplyCALReply(calReplyLong, 0, nil, nil)
				}(),
				tagName: "someTag",
				addResponseCode: func(t *testing.T) func(string, apiModel.PlcResponseCode) {
					codes := map[string]apiModel.PlcResponseCode{
						"someTag": apiModel.PlcResponseCode_OK,
					}
					return func(name string, responseCode apiModel.PlcResponseCode) {
						if code, ok := codes[name]; ok {
							assert.Equal(t, code, responseCode)
						} else {
							t.Errorf("code for %s not found", name)
						}
					}
				},
				addPlcValue: func(t *testing.T) func(string, apiValues.PlcValue) {
					values := map[string]apiValues.PlcValue{
						"someTag": spiValues.NewPlcRawByteArray([]byte{1, 2, 3, 4}),
					}
					return func(name string, plcValue apiValues.PlcValue) {
						if value, ok := values[name]; ok {
							assert.Equal(t, value, plcValue)
						} else {
							t.Errorf("value for %s not found", name)
						}
					}
				},
			},
			setup: func(t *testing.T, args *args) {
				executor := pool.NewFixedSizeExecutor(10, 50, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				executor.Start()
				transactionManager := transactions.NewRequestTransactionManager(
					1,
					options.WithCustomLogger(testutils.ProduceTestingLogger(t)),
					transactions.WithCustomExecutor(executor),
				)
				t.Cleanup(func() {
					assert.NoError(t, transactionManager.Close())
				})
				transaction := transactionManager.StartTransaction()
				transaction.Submit(func(transaction transactions.RequestTransaction) {
					// NO-OP
				})
				args.transaction = transaction
			},
		},
		{
			name: "CALDataIdentifyReplyExactly (delays)",
			args: args{
				encodedReply: func() readWriteModel.EncodedReplyCALReply {
					command := readWriteModel.NewIdentifyReplyCommandDelays([]byte{1, 2, 3, 4}, 5, 5)
					calDataIdentify := readWriteModel.NewCALDataIdentifyReply(readWriteModel.Attribute_Delays, command, readWriteModel.CALCommandTypeContainer_CALCommandIdentify, nil, nil)
					calReplyLong := readWriteModel.NewCALReplyLong(0, readWriteModel.NewUnitAddress(1), readWriteModel.NewBridgeAddress(2), readWriteModel.NewSerialInterfaceAddress(3), nil, nil, 0, calDataIdentify, nil, nil)
					return readWriteModel.NewEncodedReplyCALReply(calReplyLong, 0, nil, nil)
				}(),
				tagName: "someTag",
				addResponseCode: func(t *testing.T) func(string, apiModel.PlcResponseCode) {
					codes := map[string]apiModel.PlcResponseCode{
						"someTag": apiModel.PlcResponseCode_OK,
					}
					return func(name string, responseCode apiModel.PlcResponseCode) {
						if code, ok := codes[name]; ok {
							assert.Equal(t, code, responseCode)
						} else {
							t.Errorf("code for %s not found", name)
						}
					}
				},
				addPlcValue: func(t *testing.T) func(string, apiValues.PlcValue) {
					values := map[string]apiValues.PlcValue{
						"someTag": spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
							"reStrikeDelay": spiValues.NewPlcUSINT(5),
							"terminalLevel": spiValues.NewPlcRawByteArray([]byte{1, 2, 3, 4}),
						}),
					}
					return func(name string, plcValue apiValues.PlcValue) {
						if value, ok := values[name]; ok {
							assert.Equal(t, value, plcValue)
						} else {
							t.Errorf("value for %s not found", name)
						}
					}
				},
			},
			setup: func(t *testing.T, args *args) {
				executor := pool.NewFixedSizeExecutor(10, 50, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				executor.Start()
				transactionManager := transactions.NewRequestTransactionManager(
					1,
					options.WithCustomLogger(testutils.ProduceTestingLogger(t)),
					transactions.WithCustomExecutor(executor),
				)
				t.Cleanup(func() {
					assert.NoError(t, transactionManager.Close())
				})
				transaction := transactionManager.StartTransaction()
				transaction.Submit(func(transaction transactions.RequestTransaction) {
					// NO-OP
				})
				args.transaction = transaction
			},
		},
		{
			name: "CALDataIdentifyReplyExactly (dsi status)",
			args: args{
				encodedReply: func() readWriteModel.EncodedReplyCALReply {
					command := readWriteModel.NewIdentifyReplyCommandDSIStatus(
						readWriteModel.ChannelStatus_OK,
						readWriteModel.ChannelStatus_LAMP_FAULT,
						readWriteModel.ChannelStatus_CURRENT_LIMIT_OR_SHORT,
						readWriteModel.ChannelStatus_OK,
						readWriteModel.ChannelStatus_OK,
						readWriteModel.ChannelStatus_OK,
						readWriteModel.ChannelStatus_OK,
						readWriteModel.ChannelStatus_OK,
						readWriteModel.UnitStatus_OK,
						12,
						9,
					)
					calDataIdentify := readWriteModel.NewCALDataIdentifyReply(readWriteModel.Attribute_DSIStatus, command, readWriteModel.CALCommandTypeContainer_CALCommandIdentify, nil, nil)
					calReplyLong := readWriteModel.NewCALReplyLong(0, readWriteModel.NewUnitAddress(1), readWriteModel.NewBridgeAddress(2), readWriteModel.NewSerialInterfaceAddress(3), nil, nil, 0, calDataIdentify, nil, nil)
					return readWriteModel.NewEncodedReplyCALReply(calReplyLong, 0, nil, nil)
				}(),
				tagName: "someTag",
				addResponseCode: func(t *testing.T) func(string, apiModel.PlcResponseCode) {
					codes := map[string]apiModel.PlcResponseCode{
						"someTag": apiModel.PlcResponseCode_OK,
					}
					return func(name string, responseCode apiModel.PlcResponseCode) {
						if code, ok := codes[name]; ok {
							assert.Equal(t, code, responseCode)
						} else {
							t.Errorf("code for %s not found", name)
						}
					}
				},
				addPlcValue: func(t *testing.T) func(string, apiValues.PlcValue) {
					values := map[string]apiValues.PlcValue{
						"someTag": spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
							"channelStatus1":          spiValues.NewPlcSTRING("OK"),
							"channelStatus2":          spiValues.NewPlcSTRING("LAMP_FAULT"),
							"channelStatus3":          spiValues.NewPlcSTRING("CURRENT_LIMIT_OR_SHORT"),
							"channelStatus4":          spiValues.NewPlcSTRING("OK"),
							"channelStatus5":          spiValues.NewPlcSTRING("OK"),
							"channelStatus6":          spiValues.NewPlcSTRING("OK"),
							"channelStatus7":          spiValues.NewPlcSTRING("OK"),
							"channelStatus8":          spiValues.NewPlcSTRING("OK"),
							"dimmingUCRevisionNumber": spiValues.NewPlcUSINT(12),
							"unitStatus":              spiValues.NewPlcSTRING("OK"),
						}),
					}
					return func(name string, plcValue apiValues.PlcValue) {
						if value, ok := values[name]; ok {
							assert.Equal(t, value, plcValue)
						} else {
							t.Errorf("value for %s not found", name)
						}
					}
				},
			},
			setup: func(t *testing.T, args *args) {
				executor := pool.NewFixedSizeExecutor(10, 50, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				executor.Start()
				transactionManager := transactions.NewRequestTransactionManager(
					1,
					options.WithCustomLogger(testutils.ProduceTestingLogger(t)),
					transactions.WithCustomExecutor(executor),
				)
				t.Cleanup(func() {
					assert.NoError(t, transactionManager.Close())
				})
				transaction := transactionManager.StartTransaction()
				transaction.Submit(func(transaction transactions.RequestTransaction) {
					// NO-OP
				})
				args.transaction = transaction
			},
		},
		{
			name: "CALDataIdentifyReplyExactly (extended diagnostic summary)",
			args: args{
				encodedReply: func() readWriteModel.EncodedReplyCALReply {
					command := readWriteModel.NewIdentifyReplyCommandExtendedDiagnosticSummary(
						readWriteModel.ApplicationIdContainer_LIGHTING_3B,
						readWriteModel.ApplicationIdContainer_LIGHTING_3C,
						12,
						13,
						14,
						15,
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						true,
						8,
					)
					calDataIdentify := readWriteModel.NewCALDataIdentifyReply(readWriteModel.Attribute_DSIStatus, command, readWriteModel.CALCommandTypeContainer_CALCommandIdentify, nil, nil)
					calReplyLong := readWriteModel.NewCALReplyLong(0, readWriteModel.NewUnitAddress(1), readWriteModel.NewBridgeAddress(2), readWriteModel.NewSerialInterfaceAddress(3), nil, nil, 0, calDataIdentify, nil, nil)
					return readWriteModel.NewEncodedReplyCALReply(calReplyLong, 0, nil, nil)
				}(),
				tagName: "someTag",
				addResponseCode: func(t *testing.T) func(string, apiModel.PlcResponseCode) {
					codes := map[string]apiModel.PlcResponseCode{
						"someTag": apiModel.PlcResponseCode_OK,
					}
					return func(name string, responseCode apiModel.PlcResponseCode) {
						if code, ok := codes[name]; ok {
							assert.Equal(t, code, responseCode)
						} else {
							t.Errorf("code for %s not found", name)
						}
					}
				},
				addPlcValue: func(t *testing.T) func(string, apiValues.PlcValue) {
					values := map[string]apiValues.PlcValue{
						"someTag": spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
							"lowApplication":         spiValues.NewPlcSTRING("LIGHTING_3B"),
							"highApplication":        spiValues.NewPlcSTRING("LIGHTING_3C"),
							"area":                   spiValues.NewPlcUSINT(12),
							"crc":                    spiValues.NewPlcUINT(13),
							"serialNumber":           spiValues.NewPlcUDINT(14),
							"networkVoltage":         spiValues.NewPlcUSINT(15),
							"unitInLearnMode":        spiValues.NewPlcBOOL(true),
							"networkVoltageLow":      spiValues.NewPlcBOOL(true),
							"networkVoltageMarginal": spiValues.NewPlcBOOL(true),
							"enableChecksumAlarm":    spiValues.NewPlcBOOL(true),
							"outputUnit":             spiValues.NewPlcBOOL(true),
							"installationMMIError":   spiValues.NewPlcBOOL(true),
							"EEWriteError":           spiValues.NewPlcBOOL(true),
							"EEChecksumError":        spiValues.NewPlcBOOL(true),
							"EEDataError":            spiValues.NewPlcBOOL(true),
							"microReset":             spiValues.NewPlcBOOL(true),
							"commsTxError":           spiValues.NewPlcBOOL(true),
							"internalStackOverflow":  spiValues.NewPlcBOOL(true),
							"microPowerReset":        spiValues.NewPlcBOOL(true),
						}),
					}
					return func(name string, plcValue apiValues.PlcValue) {
						if value, ok := values[name]; ok {
							assert.Equal(t, value, plcValue)
						} else {
							t.Errorf("value for %s not found", name)
						}
					}
				},
			},
			setup: func(t *testing.T, args *args) {
				executor := pool.NewFixedSizeExecutor(10, 50, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				executor.Start()
				transactionManager := transactions.NewRequestTransactionManager(
					1,
					options.WithCustomLogger(testutils.ProduceTestingLogger(t)),
					transactions.WithCustomExecutor(executor),
				)
				t.Cleanup(func() {
					assert.NoError(t, transactionManager.Close())
				})
				transaction := transactionManager.StartTransaction()
				transaction.Submit(func(transaction transactions.RequestTransaction) {
					// NO-OP
				})
				args.transaction = transaction
			},
		},
		{
			name: "CALDataIdentifyReplyExactly (summary)",
			args: args{
				encodedReply: func() readWriteModel.EncodedReplyCALReply {
					command := readWriteModel.NewIdentifyReplyCommandSummary("pineapple", 1, "13", 3)
					calDataIdentify := readWriteModel.NewCALDataIdentifyReply(readWriteModel.Attribute_DSIStatus, command, readWriteModel.CALCommandTypeContainer_CALCommandIdentify, nil, nil)
					calReplyLong := readWriteModel.NewCALReplyLong(0, readWriteModel.NewUnitAddress(1), readWriteModel.NewBridgeAddress(2), readWriteModel.NewSerialInterfaceAddress(3), nil, nil, 0, calDataIdentify, nil, nil)
					return readWriteModel.NewEncodedReplyCALReply(calReplyLong, 0, nil, nil)
				}(),
				tagName: "someTag",
				addResponseCode: func(t *testing.T) func(string, apiModel.PlcResponseCode) {
					codes := map[string]apiModel.PlcResponseCode{
						"someTag": apiModel.PlcResponseCode_OK,
					}
					return func(name string, responseCode apiModel.PlcResponseCode) {
						if code, ok := codes[name]; ok {
							assert.Equal(t, code, responseCode)
						} else {
							t.Errorf("code for %s not found", name)
						}
					}
				},
				addPlcValue: func(t *testing.T) func(string, apiValues.PlcValue) {
					values := map[string]apiValues.PlcValue{
						"someTag": spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
							"partName":        spiValues.NewPlcSTRING("pineapple"),
							"unitServiceType": spiValues.NewPlcUSINT(1),
							"version":         spiValues.NewPlcSTRING("13"),
						}),
					}
					return func(name string, plcValue apiValues.PlcValue) {
						if value, ok := values[name]; ok {
							assert.Equal(t, value, plcValue)
						} else {
							t.Errorf("value for %s not found", name)
						}
					}
				},
			},
			setup: func(t *testing.T, args *args) {
				executor := pool.NewFixedSizeExecutor(10, 50, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				executor.Start()
				transactionManager := transactions.NewRequestTransactionManager(
					1,
					options.WithCustomLogger(testutils.ProduceTestingLogger(t)),
					transactions.WithCustomExecutor(executor),
				)
				t.Cleanup(func() {
					assert.NoError(t, transactionManager.Close())
				})
				transaction := transactionManager.StartTransaction()
				transaction.Submit(func(transaction transactions.RequestTransaction) {
					// NO-OP
				})
				args.transaction = transaction
			},
		},
		{
			name: "CALDataIdentifyReplyExactly (firmware version)",
			args: args{
				encodedReply: func() readWriteModel.EncodedReplyCALReply {
					command := readWriteModel.NewIdentifyReplyCommandFirmwareVersion("13", 1)
					calDataIdentify := readWriteModel.NewCALDataIdentifyReply(readWriteModel.Attribute_DSIStatus, command, readWriteModel.CALCommandTypeContainer_CALCommandIdentify, nil, nil)
					calReplyLong := readWriteModel.NewCALReplyLong(0, readWriteModel.NewUnitAddress(1), readWriteModel.NewBridgeAddress(2), readWriteModel.NewSerialInterfaceAddress(3), nil, nil, 0, calDataIdentify, nil, nil)
					return readWriteModel.NewEncodedReplyCALReply(calReplyLong, 0, nil, nil)
				}(),
				tagName: "someTag",
				addResponseCode: func(t *testing.T) func(string, apiModel.PlcResponseCode) {
					codes := map[string]apiModel.PlcResponseCode{
						"someTag": apiModel.PlcResponseCode_OK,
					}
					return func(name string, responseCode apiModel.PlcResponseCode) {
						if code, ok := codes[name]; ok {
							assert.Equal(t, code, responseCode)
						} else {
							t.Errorf("code for %s not found", name)
						}
					}
				},
				addPlcValue: func(t *testing.T) func(string, apiValues.PlcValue) {
					values := map[string]apiValues.PlcValue{
						"someTag": spiValues.NewPlcSTRING("13"),
					}
					return func(name string, plcValue apiValues.PlcValue) {
						if value, ok := values[name]; ok {
							assert.Equal(t, value, plcValue)
						} else {
							t.Errorf("value for %s not found", name)
						}
					}
				},
			},
			setup: func(t *testing.T, args *args) {
				executor := pool.NewFixedSizeExecutor(10, 50, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				executor.Start()
				transactionManager := transactions.NewRequestTransactionManager(
					1,
					options.WithCustomLogger(testutils.ProduceTestingLogger(t)),
					transactions.WithCustomExecutor(executor),
				)
				t.Cleanup(func() {
					assert.NoError(t, transactionManager.Close())
				})
				transaction := transactionManager.StartTransaction()
				transaction.Submit(func(transaction transactions.RequestTransaction) {
					// NO-OP
				})
				args.transaction = transaction
			},
		},
		{
			name: "CALDataIdentifyReplyExactly (GAV physical addresses)",
			args: args{
				encodedReply: func() readWriteModel.EncodedReplyCALReply {
					command := readWriteModel.NewIdentifyReplyCommandGAVPhysicalAddresses([]byte{1, 2, 3, 4}, 4)
					calDataIdentify := readWriteModel.NewCALDataIdentifyReply(readWriteModel.Attribute_DSIStatus, command, readWriteModel.CALCommandTypeContainer_CALCommandIdentify, nil, nil)
					calReplyLong := readWriteModel.NewCALReplyLong(0, readWriteModel.NewUnitAddress(1), readWriteModel.NewBridgeAddress(2), readWriteModel.NewSerialInterfaceAddress(3), nil, nil, 0, calDataIdentify, nil, nil)
					return readWriteModel.NewEncodedReplyCALReply(calReplyLong, 0, nil, nil)
				}(),
				tagName: "someTag",
				addResponseCode: func(t *testing.T) func(string, apiModel.PlcResponseCode) {
					codes := map[string]apiModel.PlcResponseCode{
						"someTag": apiModel.PlcResponseCode_OK,
					}
					return func(name string, responseCode apiModel.PlcResponseCode) {
						if code, ok := codes[name]; ok {
							assert.Equal(t, code, responseCode)
						} else {
							t.Errorf("code for %s not found", name)
						}
					}
				},
				addPlcValue: func(t *testing.T) func(string, apiValues.PlcValue) {
					values := map[string]apiValues.PlcValue{
						"someTag": spiValues.NewPlcRawByteArray([]byte{1, 2, 3, 4}),
					}
					return func(name string, plcValue apiValues.PlcValue) {
						if value, ok := values[name]; ok {
							assert.Equal(t, value, plcValue)
						} else {
							t.Errorf("value for %s not found", name)
						}
					}
				},
			},
			setup: func(t *testing.T, args *args) {
				executor := pool.NewFixedSizeExecutor(10, 50, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				executor.Start()
				transactionManager := transactions.NewRequestTransactionManager(
					1,
					options.WithCustomLogger(testutils.ProduceTestingLogger(t)),
					transactions.WithCustomExecutor(executor),
				)
				t.Cleanup(func() {
					assert.NoError(t, transactionManager.Close())
				})
				transaction := transactionManager.StartTransaction()
				transaction.Submit(func(transaction transactions.RequestTransaction) {
					// NO-OP
				})
				args.transaction = transaction
			},
		},
		{
			name: "CALDataIdentifyReplyExactly (GAV values current)",
			args: args{
				encodedReply: func() readWriteModel.EncodedReplyCALReply {
					command := readWriteModel.NewIdentifyReplyCommandGAVValuesCurrent([]byte{1, 2, 3, 4}, 4)
					calDataIdentify := readWriteModel.NewCALDataIdentifyReply(readWriteModel.Attribute_DSIStatus, command, readWriteModel.CALCommandTypeContainer_CALCommandIdentify, nil, nil)
					calReplyLong := readWriteModel.NewCALReplyLong(0, readWriteModel.NewUnitAddress(1), readWriteModel.NewBridgeAddress(2), readWriteModel.NewSerialInterfaceAddress(3), nil, nil, 0, calDataIdentify, nil, nil)
					return readWriteModel.NewEncodedReplyCALReply(calReplyLong, 0, nil, nil)
				}(),
				tagName: "someTag",
				addResponseCode: func(t *testing.T) func(string, apiModel.PlcResponseCode) {
					codes := map[string]apiModel.PlcResponseCode{
						"someTag": apiModel.PlcResponseCode_OK,
					}
					return func(name string, responseCode apiModel.PlcResponseCode) {
						if code, ok := codes[name]; ok {
							assert.Equal(t, code, responseCode)
						} else {
							t.Errorf("code for %s not found", name)
						}
					}
				},
				addPlcValue: func(t *testing.T) func(string, apiValues.PlcValue) {
					values := map[string]apiValues.PlcValue{
						"someTag": spiValues.NewPlcRawByteArray([]byte{1, 2, 3, 4}),
					}
					return func(name string, plcValue apiValues.PlcValue) {
						if value, ok := values[name]; ok {
							assert.Equal(t, value, plcValue)
						} else {
							t.Errorf("value for %s not found", name)
						}
					}
				},
			},
			setup: func(t *testing.T, args *args) {
				executor := pool.NewFixedSizeExecutor(10, 50, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				executor.Start()
				transactionManager := transactions.NewRequestTransactionManager(
					1,
					options.WithCustomLogger(testutils.ProduceTestingLogger(t)),
					transactions.WithCustomExecutor(executor),
				)
				t.Cleanup(func() {
					assert.NoError(t, transactionManager.Close())
				})
				transaction := transactionManager.StartTransaction()
				transaction.Submit(func(transaction transactions.RequestTransaction) {
					// NO-OP
				})
				args.transaction = transaction
			},
		},
		{
			name: "CALDataIdentifyReplyExactly (GAV values stored)",
			args: args{
				encodedReply: func() readWriteModel.EncodedReplyCALReply {
					command := readWriteModel.NewIdentifyReplyCommandGAVValuesStored([]byte{1, 2, 3, 4}, 4)
					calDataIdentify := readWriteModel.NewCALDataIdentifyReply(readWriteModel.Attribute_DSIStatus, command, readWriteModel.CALCommandTypeContainer_CALCommandIdentify, nil, nil)
					calReplyLong := readWriteModel.NewCALReplyLong(0, readWriteModel.NewUnitAddress(1), readWriteModel.NewBridgeAddress(2), readWriteModel.NewSerialInterfaceAddress(3), nil, nil, 0, calDataIdentify, nil, nil)
					return readWriteModel.NewEncodedReplyCALReply(calReplyLong, 0, nil, nil)
				}(),
				tagName: "someTag",
				addResponseCode: func(t *testing.T) func(string, apiModel.PlcResponseCode) {
					codes := map[string]apiModel.PlcResponseCode{
						"someTag": apiModel.PlcResponseCode_OK,
					}
					return func(name string, responseCode apiModel.PlcResponseCode) {
						if code, ok := codes[name]; ok {
							assert.Equal(t, code, responseCode)
						} else {
							t.Errorf("code for %s not found", name)
						}
					}
				},
				addPlcValue: func(t *testing.T) func(string, apiValues.PlcValue) {
					values := map[string]apiValues.PlcValue{
						"someTag": spiValues.NewPlcRawByteArray([]byte{1, 2, 3, 4}),
					}
					return func(name string, plcValue apiValues.PlcValue) {
						if value, ok := values[name]; ok {
							assert.Equal(t, value, plcValue)
						} else {
							t.Errorf("value for %s not found", name)
						}
					}
				},
			},
			setup: func(t *testing.T, args *args) {
				executor := pool.NewFixedSizeExecutor(10, 50, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				executor.Start()
				transactionManager := transactions.NewRequestTransactionManager(
					1,
					options.WithCustomLogger(testutils.ProduceTestingLogger(t)),
					transactions.WithCustomExecutor(executor),
				)
				t.Cleanup(func() {
					assert.NoError(t, transactionManager.Close())
				})
				transaction := transactionManager.StartTransaction()
				transaction.Submit(func(transaction transactions.RequestTransaction) {
					// NO-OP
				})
				args.transaction = transaction
			},
		},
		{
			name: "CALDataIdentifyReplyExactly (logical assignment)",
			args: args{
				encodedReply: func() readWriteModel.EncodedReplyCALReply {
					command := readWriteModel.NewIdentifyReplyCommandLogicalAssignment([]readWriteModel.LogicAssignment{
						readWriteModel.NewLogicAssignment(true, true, true, true, true, true),
						readWriteModel.NewLogicAssignment(true, true, true, true, true, true),
					}, 4)
					calDataIdentify := readWriteModel.NewCALDataIdentifyReply(readWriteModel.Attribute_DSIStatus, command, readWriteModel.CALCommandTypeContainer_CALCommandIdentify, nil, nil)
					calReplyLong := readWriteModel.NewCALReplyLong(0, readWriteModel.NewUnitAddress(1), readWriteModel.NewBridgeAddress(2), readWriteModel.NewSerialInterfaceAddress(3), nil, nil, 0, calDataIdentify, nil, nil)
					return readWriteModel.NewEncodedReplyCALReply(calReplyLong, 0, nil, nil)
				}(),
				tagName: "someTag",
				addResponseCode: func(t *testing.T) func(string, apiModel.PlcResponseCode) {
					codes := map[string]apiModel.PlcResponseCode{
						"someTag": apiModel.PlcResponseCode_OK,
					}
					return func(name string, responseCode apiModel.PlcResponseCode) {
						if code, ok := codes[name]; ok {
							assert.Equal(t, code, responseCode)
						} else {
							t.Errorf("code for %s not found", name)
						}
					}
				},
				addPlcValue: func(t *testing.T) func(string, apiValues.PlcValue) {
					values := map[string]apiValues.PlcValue{
						"someTag": spiValues.NewPlcList([]apiValues.PlcValue{
							spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
								"greaterOfOrLogic": spiValues.NewPlcBOOL(true),
								"reStrikeDelay":    spiValues.NewPlcBOOL(true),
								"assignedToGav16":  spiValues.NewPlcBOOL(true),
								"assignedToGav15":  spiValues.NewPlcBOOL(true),
								"assignedToGav14":  spiValues.NewPlcBOOL(true),
								"assignedToGav13":  spiValues.NewPlcBOOL(true),
							}),
							spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
								"greaterOfOrLogic": spiValues.NewPlcBOOL(true),
								"reStrikeDelay":    spiValues.NewPlcBOOL(true),
								"assignedToGav16":  spiValues.NewPlcBOOL(true),
								"assignedToGav15":  spiValues.NewPlcBOOL(true),
								"assignedToGav14":  spiValues.NewPlcBOOL(true),
								"assignedToGav13":  spiValues.NewPlcBOOL(true),
							}),
						}),
					}
					return func(name string, plcValue apiValues.PlcValue) {
						if value, ok := values[name]; ok {
							assert.Equal(t, value, plcValue)
						} else {
							t.Errorf("value for %s not found", name)
						}
					}
				},
			},
			setup: func(t *testing.T, args *args) {
				executor := pool.NewFixedSizeExecutor(10, 50, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				executor.Start()
				transactionManager := transactions.NewRequestTransactionManager(
					1,
					options.WithCustomLogger(testutils.ProduceTestingLogger(t)),
					transactions.WithCustomExecutor(executor),
				)
				t.Cleanup(func() {
					assert.NoError(t, transactionManager.Close())
				})
				transaction := transactionManager.StartTransaction()
				transaction.Submit(func(transaction transactions.RequestTransaction) {
					// NO-OP
				})
				args.transaction = transaction
			},
		},
		{
			name: "CALDataIdentifyReplyExactly (manufacturer)",
			args: args{
				encodedReply: func() readWriteModel.EncodedReplyCALReply {
					command := readWriteModel.NewIdentifyReplyCommandManufacturer("Apache", 13)
					calDataIdentify := readWriteModel.NewCALDataIdentifyReply(readWriteModel.Attribute_DSIStatus, command, readWriteModel.CALCommandTypeContainer_CALCommandIdentify, nil, nil)
					calReplyLong := readWriteModel.NewCALReplyLong(0, readWriteModel.NewUnitAddress(1), readWriteModel.NewBridgeAddress(2), readWriteModel.NewSerialInterfaceAddress(3), nil, nil, 0, calDataIdentify, nil, nil)
					return readWriteModel.NewEncodedReplyCALReply(calReplyLong, 0, nil, nil)
				}(),
				tagName: "someTag",
				addResponseCode: func(t *testing.T) func(string, apiModel.PlcResponseCode) {
					codes := map[string]apiModel.PlcResponseCode{
						"someTag": apiModel.PlcResponseCode_OK,
					}
					return func(name string, responseCode apiModel.PlcResponseCode) {
						if code, ok := codes[name]; ok {
							assert.Equal(t, code, responseCode)
						} else {
							t.Errorf("code for %s not found", name)
						}
					}
				},
				addPlcValue: func(t *testing.T) func(string, apiValues.PlcValue) {
					values := map[string]apiValues.PlcValue{
						"someTag": spiValues.NewPlcSTRING("Apache"),
					}
					return func(name string, plcValue apiValues.PlcValue) {
						if value, ok := values[name]; ok {
							assert.Equal(t, value, plcValue)
						} else {
							t.Errorf("value for %s not found", name)
						}
					}
				},
			},
			setup: func(t *testing.T, args *args) {
				executor := pool.NewFixedSizeExecutor(10, 50, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				executor.Start()
				transactionManager := transactions.NewRequestTransactionManager(
					1,
					options.WithCustomLogger(testutils.ProduceTestingLogger(t)),
					transactions.WithCustomExecutor(executor),
				)
				t.Cleanup(func() {
					assert.NoError(t, transactionManager.Close())
				})
				transaction := transactionManager.StartTransaction()
				transaction.Submit(func(transaction transactions.RequestTransaction) {
					// NO-OP
				})
				args.transaction = transaction
			},
		},
		{
			name: "CALDataIdentifyReplyExactly (maximum levels)",
			args: args{
				encodedReply: func() readWriteModel.EncodedReplyCALReply {
					command := readWriteModel.NewIdentifyReplyCommandMaximumLevels([]byte{1, 2, 3, 4}, 1)
					calDataIdentify := readWriteModel.NewCALDataIdentifyReply(readWriteModel.Attribute_DSIStatus, command, readWriteModel.CALCommandTypeContainer_CALCommandIdentify, nil, nil)
					calReplyLong := readWriteModel.NewCALReplyLong(0, readWriteModel.NewUnitAddress(1), readWriteModel.NewBridgeAddress(2), readWriteModel.NewSerialInterfaceAddress(3), nil, nil, 0, calDataIdentify, nil, nil)
					return readWriteModel.NewEncodedReplyCALReply(calReplyLong, 0, nil, nil)
				}(),
				tagName: "someTag",
				addResponseCode: func(t *testing.T) func(string, apiModel.PlcResponseCode) {
					codes := map[string]apiModel.PlcResponseCode{
						"someTag": apiModel.PlcResponseCode_OK,
					}
					return func(name string, responseCode apiModel.PlcResponseCode) {
						if code, ok := codes[name]; ok {
							assert.Equal(t, code, responseCode)
						} else {
							t.Errorf("code for %s not found", name)
						}
					}
				},
				addPlcValue: func(t *testing.T) func(string, apiValues.PlcValue) {
					values := map[string]apiValues.PlcValue{
						"someTag": spiValues.NewPlcRawByteArray([]byte{1, 2, 3, 4}),
					}
					return func(name string, plcValue apiValues.PlcValue) {
						if value, ok := values[name]; ok {
							assert.Equal(t, value, plcValue)
						} else {
							t.Errorf("value for %s not found", name)
						}
					}
				},
			},
			setup: func(t *testing.T, args *args) {
				executor := pool.NewFixedSizeExecutor(10, 50, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				executor.Start()
				transactionManager := transactions.NewRequestTransactionManager(
					1,
					options.WithCustomLogger(testutils.ProduceTestingLogger(t)),
					transactions.WithCustomExecutor(executor),
				)
				t.Cleanup(func() {
					assert.NoError(t, transactionManager.Close())
				})
				transaction := transactionManager.StartTransaction()
				transaction.Submit(func(transaction transactions.RequestTransaction) {
					// NO-OP
				})
				args.transaction = transaction
			},
		},
		{
			name: "CALDataIdentifyReplyExactly (minimum levels)",
			args: args{
				encodedReply: func() readWriteModel.EncodedReplyCALReply {
					command := readWriteModel.NewIdentifyReplyCommandMinimumLevels([]byte{1, 2, 3, 4}, 1)
					calDataIdentify := readWriteModel.NewCALDataIdentifyReply(readWriteModel.Attribute_DSIStatus, command, readWriteModel.CALCommandTypeContainer_CALCommandIdentify, nil, nil)
					calReplyLong := readWriteModel.NewCALReplyLong(0, readWriteModel.NewUnitAddress(1), readWriteModel.NewBridgeAddress(2), readWriteModel.NewSerialInterfaceAddress(3), nil, nil, 0, calDataIdentify, nil, nil)
					return readWriteModel.NewEncodedReplyCALReply(calReplyLong, 0, nil, nil)
				}(),
				tagName: "someTag",
				addResponseCode: func(t *testing.T) func(string, apiModel.PlcResponseCode) {
					codes := map[string]apiModel.PlcResponseCode{
						"someTag": apiModel.PlcResponseCode_OK,
					}
					return func(name string, responseCode apiModel.PlcResponseCode) {
						if code, ok := codes[name]; ok {
							assert.Equal(t, code, responseCode)
						} else {
							t.Errorf("code for %s not found", name)
						}
					}
				},
				addPlcValue: func(t *testing.T) func(string, apiValues.PlcValue) {
					values := map[string]apiValues.PlcValue{
						"someTag": spiValues.NewPlcRawByteArray([]byte{1, 2, 3, 4}),
					}
					return func(name string, plcValue apiValues.PlcValue) {
						if value, ok := values[name]; ok {
							assert.Equal(t, value, plcValue)
						} else {
							t.Errorf("value for %s not found", name)
						}
					}
				},
			},
			setup: func(t *testing.T, args *args) {
				executor := pool.NewFixedSizeExecutor(10, 50, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				executor.Start()
				transactionManager := transactions.NewRequestTransactionManager(
					1,
					options.WithCustomLogger(testutils.ProduceTestingLogger(t)),
					transactions.WithCustomExecutor(executor),
				)
				t.Cleanup(func() {
					assert.NoError(t, transactionManager.Close())
				})
				transaction := transactionManager.StartTransaction()
				transaction.Submit(func(transaction transactions.RequestTransaction) {
					// NO-OP
				})
				args.transaction = transaction
			},
		},
		{
			name: "CALDataIdentifyReplyExactly (network terminal levels)",
			args: args{
				encodedReply: func() readWriteModel.EncodedReplyCALReply {
					command := readWriteModel.NewIdentifyReplyCommandNetworkTerminalLevels([]byte{1, 2, 3, 4}, 1)
					calDataIdentify := readWriteModel.NewCALDataIdentifyReply(readWriteModel.Attribute_DSIStatus, command, readWriteModel.CALCommandTypeContainer_CALCommandIdentify, nil, nil)
					calReplyLong := readWriteModel.NewCALReplyLong(0, readWriteModel.NewUnitAddress(1), readWriteModel.NewBridgeAddress(2), readWriteModel.NewSerialInterfaceAddress(3), nil, nil, 0, calDataIdentify, nil, nil)
					return readWriteModel.NewEncodedReplyCALReply(calReplyLong, 0, nil, nil)
				}(),
				tagName: "someTag",
				addResponseCode: func(t *testing.T) func(string, apiModel.PlcResponseCode) {
					codes := map[string]apiModel.PlcResponseCode{
						"someTag": apiModel.PlcResponseCode_OK,
					}
					return func(name string, responseCode apiModel.PlcResponseCode) {
						if code, ok := codes[name]; ok {
							assert.Equal(t, code, responseCode)
						} else {
							t.Errorf("code for %s not found", name)
						}
					}
				},
				addPlcValue: func(t *testing.T) func(string, apiValues.PlcValue) {
					values := map[string]apiValues.PlcValue{
						"someTag": spiValues.NewPlcRawByteArray([]byte{1, 2, 3, 4}),
					}
					return func(name string, plcValue apiValues.PlcValue) {
						if value, ok := values[name]; ok {
							assert.Equal(t, value, plcValue)
						} else {
							t.Errorf("value for %s not found", name)
						}
					}
				},
			},
			setup: func(t *testing.T, args *args) {
				executor := pool.NewFixedSizeExecutor(10, 50, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				executor.Start()
				transactionManager := transactions.NewRequestTransactionManager(
					1,
					options.WithCustomLogger(testutils.ProduceTestingLogger(t)),
					transactions.WithCustomExecutor(executor),
				)
				t.Cleanup(func() {
					assert.NoError(t, transactionManager.Close())
				})
				transaction := transactionManager.StartTransaction()
				transaction.Submit(func(transaction transactions.RequestTransaction) {
					// NO-OP
				})
				args.transaction = transaction
			},
		},
		{
			name: "CALDataIdentifyReplyExactly (network voltage)",
			args: args{
				encodedReply: func() readWriteModel.EncodedReplyCALReply {
					command := readWriteModel.NewIdentifyReplyCommandNetworkVoltage("13.3", "3", 3)
					calDataIdentify := readWriteModel.NewCALDataIdentifyReply(readWriteModel.Attribute_DSIStatus, command, readWriteModel.CALCommandTypeContainer_CALCommandIdentify, nil, nil)
					calReplyLong := readWriteModel.NewCALReplyLong(0, readWriteModel.NewUnitAddress(1), readWriteModel.NewBridgeAddress(2), readWriteModel.NewSerialInterfaceAddress(3), nil, nil, 0, calDataIdentify, nil, nil)
					return readWriteModel.NewEncodedReplyCALReply(calReplyLong, 0, nil, nil)
				}(),
				tagName: "someTag",
				addResponseCode: func(t *testing.T) func(string, apiModel.PlcResponseCode) {
					codes := map[string]apiModel.PlcResponseCode{
						"someTag": apiModel.PlcResponseCode_OK,
					}
					return func(name string, responseCode apiModel.PlcResponseCode) {
						if code, ok := codes[name]; ok {
							assert.Equal(t, code, responseCode)
						} else {
							t.Errorf("code for %s not found", name)
						}
					}
				},
				addPlcValue: func(t *testing.T) func(string, apiValues.PlcValue) {
					values := map[string]apiValues.PlcValue{
						"someTag": spiValues.NewPlcLREAL(13.600000000000001),
					}
					return func(name string, plcValue apiValues.PlcValue) {
						if value, ok := values[name]; ok {
							assert.Equal(t, value, plcValue)
						} else {
							t.Errorf("value for %s not found", name)
						}
					}
				},
			},
			setup: func(t *testing.T, args *args) {
				executor := pool.NewFixedSizeExecutor(10, 50, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				executor.Start()
				transactionManager := transactions.NewRequestTransactionManager(
					1,
					options.WithCustomLogger(testutils.ProduceTestingLogger(t)),
					transactions.WithCustomExecutor(executor),
				)
				t.Cleanup(func() {
					assert.NoError(t, transactionManager.Close())
				})
				transaction := transactionManager.StartTransaction()
				transaction.Submit(func(transaction transactions.RequestTransaction) {
					// NO-OP
				})
				args.transaction = transaction
			},
		},
		{
			name: "CALDataIdentifyReplyExactly (output unit summary)",
			args: args{
				encodedReply: func() readWriteModel.EncodedReplyCALReply {
					gavStoreEnabledByte1 := byte(2)
					gavStoreEnabledByte2 := byte(3)
					command := readWriteModel.NewIdentifyReplyCommandOutputUnitSummary(readWriteModel.NewIdentifyReplyCommandUnitSummary(true, true, true, true, true, true, true, true), &gavStoreEnabledByte1, &gavStoreEnabledByte2, 13, 13)
					calDataIdentify := readWriteModel.NewCALDataIdentifyReply(readWriteModel.Attribute_DSIStatus, command, readWriteModel.CALCommandTypeContainer_CALCommandIdentify, nil, nil)
					calReplyLong := readWriteModel.NewCALReplyLong(0, readWriteModel.NewUnitAddress(1), readWriteModel.NewBridgeAddress(2), readWriteModel.NewSerialInterfaceAddress(3), nil, nil, 0, calDataIdentify, nil, nil)
					return readWriteModel.NewEncodedReplyCALReply(calReplyLong, 0, nil, nil)
				}(),
				tagName: "someTag",
				addResponseCode: func(t *testing.T) func(string, apiModel.PlcResponseCode) {
					codes := map[string]apiModel.PlcResponseCode{
						"someTag": apiModel.PlcResponseCode_OK,
					}
					return func(name string, responseCode apiModel.PlcResponseCode) {
						if code, ok := codes[name]; ok {
							assert.Equal(t, code, responseCode)
						} else {
							t.Errorf("code for %s not found", name)
						}
					}
				},
				addPlcValue: func(t *testing.T) func(string, apiValues.PlcValue) {
					values := map[string]apiValues.PlcValue{
						"someTag": spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
							"unitFlags": spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
								"assertingNetworkBurden": spiValues.NewPlcBOOL(true),
								"restrikeTimingActive":   spiValues.NewPlcBOOL(true),
								"remoteOFFInputAsserted": spiValues.NewPlcBOOL(true),
								"remoteONInputAsserted":  spiValues.NewPlcBOOL(true),
								"localToggleEnabled":     spiValues.NewPlcBOOL(true),
								"localToggleActiveState": spiValues.NewPlcBOOL(true),
								"clockGenerationEnabled": spiValues.NewPlcBOOL(true),
								"unitGeneratingClock":    spiValues.NewPlcBOOL(true),
							}),
							"timeFromLastRecoverOfMainsInSeconds": spiValues.NewPlcUSINT(13),
							"gavStoreEnabledByte1":                spiValues.NewPlcUSINT(2),
							"gavStoreEnabledByte2":                spiValues.NewPlcUSINT(3),
						}),
					}
					return func(name string, plcValue apiValues.PlcValue) {
						if value, ok := values[name]; ok {
							assert.Equal(t, value, plcValue)
						} else {
							t.Errorf("value for %s not found", name)
						}
					}
				},
			},
			setup: func(t *testing.T, args *args) {
				executor := pool.NewFixedSizeExecutor(10, 50, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				executor.Start()
				transactionManager := transactions.NewRequestTransactionManager(
					1,
					options.WithCustomLogger(testutils.ProduceTestingLogger(t)),
					transactions.WithCustomExecutor(executor),
				)
				t.Cleanup(func() {
					assert.NoError(t, transactionManager.Close())
				})
				transaction := transactionManager.StartTransaction()
				transaction.Submit(func(transaction transactions.RequestTransaction) {
					// NO-OP
				})
				args.transaction = transaction
			},
		},
		{
			name: "CALDataIdentifyReplyExactly (terminal levels)",
			args: args{
				encodedReply: func() readWriteModel.EncodedReplyCALReply {
					command := readWriteModel.NewIdentifyReplyCommandTerminalLevels([]byte{1, 2, 3, 4}, 4)
					calDataIdentify := readWriteModel.NewCALDataIdentifyReply(readWriteModel.Attribute_DSIStatus, command, readWriteModel.CALCommandTypeContainer_CALCommandIdentify, nil, nil)
					calReplyLong := readWriteModel.NewCALReplyLong(0, readWriteModel.NewUnitAddress(1), readWriteModel.NewBridgeAddress(2), readWriteModel.NewSerialInterfaceAddress(3), nil, nil, 0, calDataIdentify, nil, nil)
					return readWriteModel.NewEncodedReplyCALReply(calReplyLong, 0, nil, nil)
				}(),
				tagName: "someTag",
				addResponseCode: func(t *testing.T) func(string, apiModel.PlcResponseCode) {
					codes := map[string]apiModel.PlcResponseCode{
						"someTag": apiModel.PlcResponseCode_OK,
					}
					return func(name string, responseCode apiModel.PlcResponseCode) {
						if code, ok := codes[name]; ok {
							assert.Equal(t, code, responseCode)
						} else {
							t.Errorf("code for %s not found", name)
						}
					}
				},
				addPlcValue: func(t *testing.T) func(string, apiValues.PlcValue) {
					values := map[string]apiValues.PlcValue{
						"someTag": spiValues.NewPlcRawByteArray([]byte{1, 2, 3, 4}),
					}
					return func(name string, plcValue apiValues.PlcValue) {
						if value, ok := values[name]; ok {
							assert.Equal(t, value, plcValue)
						} else {
							t.Errorf("value for %s not found", name)
						}
					}
				},
			},
			setup: func(t *testing.T, args *args) {
				executor := pool.NewFixedSizeExecutor(10, 50, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				executor.Start()
				transactionManager := transactions.NewRequestTransactionManager(
					1,
					options.WithCustomLogger(testutils.ProduceTestingLogger(t)),
					transactions.WithCustomExecutor(executor),
				)
				t.Cleanup(func() {
					assert.NoError(t, transactionManager.Close())
				})
				transaction := transactionManager.StartTransaction()
				transaction.Submit(func(transaction transactions.RequestTransaction) {
					// NO-OP
				})
				args.transaction = transaction
			},
		},
		{
			name: "CALDataIdentifyReplyExactly (type)",
			args: args{
				encodedReply: func() readWriteModel.EncodedReplyCALReply {
					command := readWriteModel.NewIdentifyReplyCommandType("chonkers", 4)
					calDataIdentify := readWriteModel.NewCALDataIdentifyReply(readWriteModel.Attribute_DSIStatus, command, readWriteModel.CALCommandTypeContainer_CALCommandIdentify, nil, nil)
					calReplyLong := readWriteModel.NewCALReplyLong(0, readWriteModel.NewUnitAddress(1), readWriteModel.NewBridgeAddress(2), readWriteModel.NewSerialInterfaceAddress(3), nil, nil, 0, calDataIdentify, nil, nil)
					return readWriteModel.NewEncodedReplyCALReply(calReplyLong, 0, nil, nil)
				}(),
				tagName: "someTag",
				addResponseCode: func(t *testing.T) func(string, apiModel.PlcResponseCode) {
					codes := map[string]apiModel.PlcResponseCode{
						"someTag": apiModel.PlcResponseCode_OK,
					}
					return func(name string, responseCode apiModel.PlcResponseCode) {
						if code, ok := codes[name]; ok {
							assert.Equal(t, code, responseCode)
						} else {
							t.Errorf("code for %s not found", name)
						}
					}
				},
				addPlcValue: func(t *testing.T) func(string, apiValues.PlcValue) {
					values := map[string]apiValues.PlcValue{
						"someTag": spiValues.NewPlcSTRING("chonkers"),
					}
					return func(name string, plcValue apiValues.PlcValue) {
						if value, ok := values[name]; ok {
							assert.Equal(t, value, plcValue)
						} else {
							t.Errorf("value for %s not found", name)
						}
					}
				},
			},
			setup: func(t *testing.T, args *args) {
				executor := pool.NewFixedSizeExecutor(10, 50, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				executor.Start()
				transactionManager := transactions.NewRequestTransactionManager(
					1,
					options.WithCustomLogger(testutils.ProduceTestingLogger(t)),
					transactions.WithCustomExecutor(executor),
				)
				t.Cleanup(func() {
					assert.NoError(t, transactionManager.Close())
				})
				transaction := transactionManager.StartTransaction()
				transaction.Submit(func(transaction transactions.RequestTransaction) {
					// NO-OP
				})
				args.transaction = transaction
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testingLogger := testutils.ProduceTestingLogger(t)
			if tt.setup != nil {
				tt.setup(t, &tt.args)
			}
			if err := MapEncodedReply(testingLogger, tt.args.transaction, tt.args.encodedReply, tt.args.tagName, tt.args.addResponseCode(t), tt.args.addPlcValue(t)); (err != nil) != tt.wantErr {
				t.Errorf("MapEncodedReply() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
