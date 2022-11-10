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

package bacnetip

import (
	"context"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
	"reflect"
	"testing"
	"time"
)

func TestTransactionStateMachine_Expect(t1 *testing.T) {
	type fields struct {
		MessageCodec          *MessageCodec
		deviceInventory       *DeviceInventory
		retryCount            int
		segmentRetryCount     int
		duplicateCount        int
		sentAllSegments       bool
		lastSequenceNumber    int
		initialSequenceNumber int
		actualWindowSize      int
		proposeWindowSize     int
		segmentTimer          int
		RequestTimer          int
	}
	type args struct {
		ctx            context.Context
		acceptsMessage spi.AcceptsMessage
		handleMessage  spi.HandleMessage
		handleError    spi.HandleError
		ttl            time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TransactionStateMachine{
				MessageCodec:          tt.fields.MessageCodec,
				deviceInventory:       tt.fields.deviceInventory,
				retryCount:            tt.fields.retryCount,
				segmentRetryCount:     tt.fields.segmentRetryCount,
				duplicateCount:        tt.fields.duplicateCount,
				sentAllSegments:       tt.fields.sentAllSegments,
				lastSequenceNumber:    tt.fields.lastSequenceNumber,
				initialSequenceNumber: tt.fields.initialSequenceNumber,
				actualWindowSize:      tt.fields.actualWindowSize,
				proposeWindowSize:     tt.fields.proposeWindowSize,
				segmentTimer:          tt.fields.segmentTimer,
				RequestTimer:          tt.fields.RequestTimer,
			}
			if err := t.Expect(tt.args.ctx, tt.args.acceptsMessage, tt.args.handleMessage, tt.args.handleError, tt.args.ttl); (err != nil) != tt.wantErr {
				t1.Errorf("Expect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransactionStateMachine_GetCodec(t1 *testing.T) {
	type fields struct {
		MessageCodec          *MessageCodec
		deviceInventory       *DeviceInventory
		retryCount            int
		segmentRetryCount     int
		duplicateCount        int
		sentAllSegments       bool
		lastSequenceNumber    int
		initialSequenceNumber int
		actualWindowSize      int
		proposeWindowSize     int
		segmentTimer          int
		RequestTimer          int
	}
	tests := []struct {
		name   string
		fields fields
		want   spi.MessageCodec
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TransactionStateMachine{
				MessageCodec:          tt.fields.MessageCodec,
				deviceInventory:       tt.fields.deviceInventory,
				retryCount:            tt.fields.retryCount,
				segmentRetryCount:     tt.fields.segmentRetryCount,
				duplicateCount:        tt.fields.duplicateCount,
				sentAllSegments:       tt.fields.sentAllSegments,
				lastSequenceNumber:    tt.fields.lastSequenceNumber,
				initialSequenceNumber: tt.fields.initialSequenceNumber,
				actualWindowSize:      tt.fields.actualWindowSize,
				proposeWindowSize:     tt.fields.proposeWindowSize,
				segmentTimer:          tt.fields.segmentTimer,
				RequestTimer:          tt.fields.RequestTimer,
			}
			if got := t.GetCodec(); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("GetCodec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionStateMachine_Send(t1 *testing.T) {
	type fields struct {
		MessageCodec          *MessageCodec
		deviceInventory       *DeviceInventory
		retryCount            int
		segmentRetryCount     int
		duplicateCount        int
		sentAllSegments       bool
		lastSequenceNumber    int
		initialSequenceNumber int
		actualWindowSize      int
		proposeWindowSize     int
		segmentTimer          int
		RequestTimer          int
	}
	type args struct {
		message spi.Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TransactionStateMachine{
				MessageCodec:          tt.fields.MessageCodec,
				deviceInventory:       tt.fields.deviceInventory,
				retryCount:            tt.fields.retryCount,
				segmentRetryCount:     tt.fields.segmentRetryCount,
				duplicateCount:        tt.fields.duplicateCount,
				sentAllSegments:       tt.fields.sentAllSegments,
				lastSequenceNumber:    tt.fields.lastSequenceNumber,
				initialSequenceNumber: tt.fields.initialSequenceNumber,
				actualWindowSize:      tt.fields.actualWindowSize,
				proposeWindowSize:     tt.fields.proposeWindowSize,
				segmentTimer:          tt.fields.segmentTimer,
				RequestTimer:          tt.fields.RequestTimer,
			}
			if err := t.Send(tt.args.message); (err != nil) != tt.wantErr {
				t1.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransactionStateMachine_SendRequest(t1 *testing.T) {
	type fields struct {
		MessageCodec          *MessageCodec
		deviceInventory       *DeviceInventory
		retryCount            int
		segmentRetryCount     int
		duplicateCount        int
		sentAllSegments       bool
		lastSequenceNumber    int
		initialSequenceNumber int
		actualWindowSize      int
		proposeWindowSize     int
		segmentTimer          int
		RequestTimer          int
	}
	type args struct {
		ctx            context.Context
		message        spi.Message
		acceptsMessage spi.AcceptsMessage
		handleMessage  spi.HandleMessage
		handleError    spi.HandleError
		ttl            time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TransactionStateMachine{
				MessageCodec:          tt.fields.MessageCodec,
				deviceInventory:       tt.fields.deviceInventory,
				retryCount:            tt.fields.retryCount,
				segmentRetryCount:     tt.fields.segmentRetryCount,
				duplicateCount:        tt.fields.duplicateCount,
				sentAllSegments:       tt.fields.sentAllSegments,
				lastSequenceNumber:    tt.fields.lastSequenceNumber,
				initialSequenceNumber: tt.fields.initialSequenceNumber,
				actualWindowSize:      tt.fields.actualWindowSize,
				proposeWindowSize:     tt.fields.proposeWindowSize,
				segmentTimer:          tt.fields.segmentTimer,
				RequestTimer:          tt.fields.RequestTimer,
			}
			if err := t.SendRequest(tt.args.ctx, tt.args.message, tt.args.acceptsMessage, tt.args.handleMessage, tt.args.handleError, tt.args.ttl); (err != nil) != tt.wantErr {
				t1.Errorf("SendRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransactionStateMachine_handleOutboundMessage(t1 *testing.T) {
	type fields struct {
		MessageCodec          *MessageCodec
		deviceInventory       *DeviceInventory
		retryCount            int
		segmentRetryCount     int
		duplicateCount        int
		sentAllSegments       bool
		lastSequenceNumber    int
		initialSequenceNumber int
		actualWindowSize      int
		proposeWindowSize     int
		segmentTimer          int
		RequestTimer          int
	}
	type args struct {
		message spi.Message
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantHandled bool
		wantErr     bool
	}{
		{
			name: "message not relevant",
		},
		{
			name: "Normal unsgemented message",
			fields: fields{
				MessageCodec:          NewMessageCodec(test.NewTransportInstance(test.NewTransport())),
				deviceInventory:       nil,
				retryCount:            0,
				segmentRetryCount:     0,
				duplicateCount:        0,
				sentAllSegments:       false,
				lastSequenceNumber:    0,
				initialSequenceNumber: 0,
				actualWindowSize:      0,
				proposeWindowSize:     0,
				segmentTimer:          0,
				RequestTimer:          0,
			},
			args: args{
				message: readWriteModel.NewBVLCOriginalUnicastNPDU(
					readWriteModel.NewNPDU(
						1,
						readWriteModel.NewNPDUControl(false, false, false, false, readWriteModel.NPDUNetworkPriority_NORMAL_MESSAGE),
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						readWriteModel.NewAPDUComplexAck(
							false,
							false,
							13,
							nil,
							nil,
							readWriteModel.NewBACnetServiceAckReadProperty(
								readWriteModel.CreateBACnetContextTagObjectIdentifier(0, 2, 1),
								readWriteModel.CreateBACnetPropertyIdentifierTagged(1, 85),
								nil,
								readWriteModel.NewBACnetConstructedDataAnalogValuePresentValue(
									readWriteModel.CreateBACnetApplicationTagReal(101),
									readWriteModel.CreateBACnetOpeningTag(3),
									readWriteModel.CreateBACnetTagHeaderBalanced(true, 3, 3),
									readWriteModel.CreateBACnetClosingTag(3),
									3,
									nil,
								),
								0,
							),
							nil,
							nil,
							0,
						),
						0,
					),
					0,
				),
			},
		},
		{
			name: "Normal segmented message",
			fields: fields{
				MessageCodec: NewMessageCodec(test.NewTransportInstance(test.NewTransport())),
				deviceInventory: func() *DeviceInventory {
					var deviceInventory = DeviceInventory{
						devices: map[string]DeviceEntry{
							"123": {
								DeviceIdentifier:          nil,
								MaximumApduLengthAccepted: readWriteModel.MaxApduLengthAccepted_NUM_OCTETS_206,
								SegmentationSupported:     readWriteModel.BACnetSegmentation_SEGMENTED_BOTH,
								VendorId:                  0,
								DeviceObjects:             nil,
							},
						},
					}
					return &deviceInventory
				}(),
				retryCount:            0,
				segmentRetryCount:     0,
				duplicateCount:        0,
				sentAllSegments:       false,
				lastSequenceNumber:    0,
				initialSequenceNumber: 0,
				actualWindowSize:      0,
				proposeWindowSize:     0,
				segmentTimer:          0,
				RequestTimer:          0,
			},
			args: args{
				message: readWriteModel.NewBVLCOriginalUnicastNPDU(
					readWriteModel.NewNPDU(
						1,
						readWriteModel.NewNPDUControl(false, true, false, false, readWriteModel.NPDUNetworkPriority_NORMAL_MESSAGE),
						nil,
						func() *uint8 {
							var elements uint8 = 3
							return &elements
						}(),
						[]uint8{0x31, 0x32, 0x33},
						nil,
						nil,
						nil,
						nil,
						nil,
						readWriteModel.NewAPDUComplexAck(
							false,
							false,
							13,
							nil,
							nil,
							readWriteModel.NewBACnetServiceAckReadProperty(
								readWriteModel.CreateBACnetContextTagObjectIdentifier(0, 2, 1),
								readWriteModel.CreateBACnetPropertyIdentifierTagged(1, 85),
								nil,
								readWriteModel.NewBACnetConstructedDataActionText(
									readWriteModel.CreateBACnetApplicationTagUnsignedInteger(100),
									func() []readWriteModel.BACnetApplicationTagCharacterString {
										var characterStrings []readWriteModel.BACnetApplicationTagCharacterString
										for i := 0; i < 100; i++ {
											characterStrings = append(characterStrings, readWriteModel.CreateBACnetApplicationTagCharacterString(readWriteModel.BACnetCharacterEncoding_ISO_10646, "ALAAARM!!"))
										}
										return characterStrings
									}(),
									readWriteModel.CreateBACnetOpeningTag(3),
									readWriteModel.CreateBACnetTagHeaderBalanced(true, 3, 3),
									readWriteModel.CreateBACnetClosingTag(3),
									3,
									nil,
								),
								0,
							),
							nil,
							nil,
							0,
						),
						0,
					),
					0,
				),
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TransactionStateMachine{
				MessageCodec:          tt.fields.MessageCodec,
				deviceInventory:       tt.fields.deviceInventory,
				retryCount:            tt.fields.retryCount,
				segmentRetryCount:     tt.fields.segmentRetryCount,
				duplicateCount:        tt.fields.duplicateCount,
				sentAllSegments:       tt.fields.sentAllSegments,
				lastSequenceNumber:    tt.fields.lastSequenceNumber,
				initialSequenceNumber: tt.fields.initialSequenceNumber,
				actualWindowSize:      tt.fields.actualWindowSize,
				proposeWindowSize:     tt.fields.proposeWindowSize,
				segmentTimer:          tt.fields.segmentTimer,
				RequestTimer:          tt.fields.RequestTimer,
			}
			gotHandled, err := t.handleOutboundMessage(tt.args.message)
			if (err != nil) != tt.wantErr {
				t1.Errorf("handleOutboundMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHandled != tt.wantHandled {
				t1.Errorf("handleOutboundMessage() gotHandled = %v, want %v", gotHandled, tt.wantHandled)
			}
		})
	}
}
