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
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
	"github.com/stretchr/testify/assert"
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
		addResponseCode func(t *testing.T) func(name string, responseCode apiModel.PlcResponseCode)
		addPlcValue     func(t *testing.T) func(name string, plcValue apiValues.PlcValue)
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
		},
		{
			name: "CALDataStatusExtendedExactly (binary)",
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
		},
		{
			name: "CALDataStatusExtendedExactly (level)",
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
		},
		{
			name: "CALDataIdentifyReplyExactly (sense levels)",
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
		},
		{
			name: "CALDataIdentifyReplyExactly (delays)",
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
		},
		{
			name: "CALDataIdentifyReplyExactly (dsi status)",
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
		},
		{
			name: "CALDataIdentifyReplyExactly (extended diagnostic summary)",
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
		},
		{
			name: "CALDataIdentifyReplyExactly (summary)",
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
		},
		{
			name: "CALDataIdentifyReplyExactly (firmware version)",
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
		},
		{
			name: "CALDataIdentifyReplyExactly (GAV physical addresses)",
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
		},
		{
			name: "CALDataIdentifyReplyExactly (GAV values current)",
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
		},
		{
			name: "CALDataIdentifyReplyExactly (GAV values stored)",
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
		},
		{
			name: "CALDataIdentifyReplyExactly (logical assignment)",
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
		},
		{
			name: "CALDataIdentifyReplyExactly (manufacturer)",
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
		},
		{
			name: "CALDataIdentifyReplyExactly (maximum levels)",
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
		},
		{
			name: "CALDataIdentifyReplyExactly (minimum levels)",
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
		},
		{
			name: "CALDataIdentifyReplyExactly (network terminal levels)",
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
		},
		{
			name: "CALDataIdentifyReplyExactly (network voltage)",
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
		},
		{
			name: "CALDataIdentifyReplyExactly (output unit summary)",
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
		},
		{
			name: "CALDataIdentifyReplyExactly (terminal levels)",
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
		},
		{
			name: "CALDataIdentifyReplyExactly (type)",
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MapEncodedReply(tt.args.transaction, tt.args.encodedReply, tt.args.tagName, tt.args.addResponseCode(t), tt.args.addPlcValue(t)); (err != nil) != tt.wantErr {
				t.Errorf("MapEncodedReply() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
