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
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"reflect"
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCBusMessage, gotSupportsRead, gotSupportsWrite, gotSupportsSubscribe, err := TagToCBusMessage(tt.args.tag, tt.args.value, tt.args.alphaGenerator, tt.args.messageCodec)
			if (err != nil) != tt.wantErr {
				t.Errorf("TagToCBusMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCBusMessage, tt.wantCBusMessage) {
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
				t.Errorf("TagToCBusMessage() gotSupportsSubscribe = %v, want %v", gotSupportsSubscribe, tt.wantSupportsSubscribe)
			}
		})
	}
}

func TestMapEncodedReply(t *testing.T) {
	type args struct {
		transaction     spi.RequestTransaction
		encodedReply    readWriteModel.EncodedReply
		tagName         string
		addResponseCode func(name string, responseCode apiModel.PlcResponseCode)
		addPlcValue     func(name string, plcValue apiValues.PlcValue)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty input",
			args: args{
				transaction: func() spi.RequestTransaction {
					transactionManager := spi.NewRequestTransactionManager(1)
					transaction := transactionManager.StartTransaction()
					transaction.Submit(func() {
						// NO-OP
					})
					return transaction
				}(),
				encodedReply:    nil,
				tagName:         "",
				addResponseCode: nil,
				addPlcValue:     nil,
			},
		},
		{
			name: "CALDataStatus",
			args: args{
				transaction: func() spi.RequestTransaction {
					transactionManager := spi.NewRequestTransactionManager(1)
					transaction := transactionManager.StartTransaction()
					transaction.Submit(func() {
						// NO-OP
					})
					return transaction
				}(),
				encodedReply: func() readWriteModel.EncodedReplyCALReply {
					calDataStatus := readWriteModel.NewCALDataStatus(readWriteModel.ApplicationIdContainer_LIGHTING_3A, 0, nil, readWriteModel.CALCommandTypeContainer_CALCommandStatus_0Bytes, nil, nil)
					calReplyShort := readWriteModel.NewCALReplyShort(0, calDataStatus, nil, nil)
					return readWriteModel.NewEncodedReplyCALReply(calReplyShort, 0, nil, nil)
				}(),
				tagName: "someTag",
				addResponseCode: func(name string, responseCode apiModel.PlcResponseCode) {
					// TODO: add assertions
				},
				addPlcValue: func(name string, plcValue apiValues.PlcValue) {
					// TODO: add assertions
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MapEncodedReply(tt.args.transaction, tt.args.encodedReply, tt.args.tagName, tt.args.addResponseCode, tt.args.addPlcValue); (err != nil) != tt.wantErr {
				t.Errorf("MapEncodedReply() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
