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
	"reflect"
	"testing"
)

func TestMapEncodedReply(t *testing.T) {
	type args struct {
		transaction     *spi.RequestTransaction
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
				transaction: func() *spi.RequestTransaction {
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
				transaction: func() *spi.RequestTransaction {
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

func TestTagToCBusMessage(t *testing.T) {
	type args struct {
		tag            apiModel.PlcTag
		value          apiValues.PlcValue
		alphaGenerator *AlphaGenerator
		messageCodec   *MessageCodec
	}
	tests := []struct {
		name                  string
		args                  args
		wantCBusMessage       readWriteModel.CBusMessage
		wantSupportsRead      bool
		wantSupportsWrite     bool
		wantSupportsSubscribe bool
		wantErr               bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCBusMessage, gotSupportsRead, gotSupportsWrite, gotSupportsSubscribe, err := TagToCBusMessage(tt.args.tag, tt.args.value, tt.args.alphaGenerator, tt.args.messageCodec)
			if (err != nil) != tt.wantErr {
				t.Errorf("TagToCBusMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCBusMessage, tt.wantCBusMessage) {
				t.Errorf("TagToCBusMessage() gotCBusMessage = %v, want %v", gotCBusMessage, tt.wantCBusMessage)
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
