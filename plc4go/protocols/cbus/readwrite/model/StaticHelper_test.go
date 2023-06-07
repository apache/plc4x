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
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculateChecksum(t *testing.T) {
	type args struct {
		writeBuffer utils.WriteBuffer
		message     spi.Message
		srchk       bool
	}
	tests := []struct {
		name    string
		args    args
		setup   func(t *testing.T, args *args)
		wantErr bool
	}{
		{
			name: "check it (not enabled)",
		},
		{
			name: "check it",
			args: args{
				writeBuffer: utils.NewWriteBufferByteBased(),
				message:     NewZoneStatus(1),
				srchk:       true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.args)
			}
			if err := CalculateChecksum(tt.args.writeBuffer, tt.args.message, tt.args.srchk); (err != nil) != tt.wantErr {
				t.Errorf("CalculateChecksum() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKnowsAccessControlCommandTypeContainer(t *testing.T) {
	type args struct {
		readBuffer utils.ReadBuffer
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no",
			args: args{
				readBuffer: utils.NewReadBufferByteBased([]byte{1, 2, 3, 4}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KnowsAccessControlCommandTypeContainer(tt.args.readBuffer); got != tt.want {
				t.Errorf("KnowsAccessControlCommandTypeContainer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKnowsAirConditioningCommandTypeContainer(t *testing.T) {
	type args struct {
		readBuffer utils.ReadBuffer
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "yes",
			args: args{
				readBuffer: utils.NewReadBufferByteBased([]byte{1, 2, 3, 4}),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KnowsAirConditioningCommandTypeContainer(tt.args.readBuffer); got != tt.want {
				t.Errorf("KnowsAirConditioningCommandTypeContainer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKnowsCALCommandTypeContainer(t *testing.T) {
	type args struct {
		readBuffer utils.ReadBuffer
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no",
			args: args{
				readBuffer: utils.NewReadBufferByteBased([]byte{1, 2, 3, 4}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KnowsCALCommandTypeContainer(tt.args.readBuffer); got != tt.want {
				t.Errorf("KnowsCALCommandTypeContainer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKnowsClockAndTimekeepingCommandTypeContainer(t *testing.T) {
	type args struct {
		readBuffer utils.ReadBuffer
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no",
			args: args{
				readBuffer: utils.NewReadBufferByteBased([]byte{1, 2, 3, 4}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KnowsClockAndTimekeepingCommandTypeContainer(tt.args.readBuffer); got != tt.want {
				t.Errorf("KnowsClockAndTimekeepingCommandTypeContainer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKnowsEnableControlCommandTypeContainer(t *testing.T) {
	type args struct {
		readBuffer utils.ReadBuffer
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no",
			args: args{
				readBuffer: utils.NewReadBufferByteBased([]byte{1, 2, 3, 4}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KnowsEnableControlCommandTypeContainer(tt.args.readBuffer); got != tt.want {
				t.Errorf("KnowsEnableControlCommandTypeContainer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKnowsErrorReportingCommandTypeContainer(t *testing.T) {
	type args struct {
		readBuffer utils.ReadBuffer
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no",
			args: args{
				readBuffer: utils.NewReadBufferByteBased([]byte{1, 2, 3, 4}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KnowsErrorReportingCommandTypeContainer(tt.args.readBuffer); got != tt.want {
				t.Errorf("KnowsErrorReportingCommandTypeContainer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKnowsLightingCommandTypeContainer(t *testing.T) {
	type args struct {
		readBuffer utils.ReadBuffer
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "yes",
			args: args{
				readBuffer: utils.NewReadBufferByteBased([]byte{1, 2, 3, 4}),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KnowsLightingCommandTypeContainer(tt.args.readBuffer); got != tt.want {
				t.Errorf("KnowsLightingCommandTypeContainer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKnowsMeasurementCommandTypeContainer(t *testing.T) {
	type args struct {
		readBuffer utils.ReadBuffer
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no",
			args: args{
				readBuffer: utils.NewReadBufferByteBased([]byte{1, 2, 3, 4}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KnowsMeasurementCommandTypeContainer(tt.args.readBuffer); got != tt.want {
				t.Errorf("KnowsMeasurementCommandTypeContainer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKnowsMediaTransportControlCommandTypeContainer(t *testing.T) {
	type args struct {
		readBuffer utils.ReadBuffer
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "yes",
			args: args{
				readBuffer: utils.NewReadBufferByteBased([]byte{1, 2, 3, 4}),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KnowsMediaTransportControlCommandTypeContainer(tt.args.readBuffer); got != tt.want {
				t.Errorf("KnowsMediaTransportControlCommandTypeContainer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKnowsMeteringCommandTypeContainer(t *testing.T) {
	type args struct {
		readBuffer utils.ReadBuffer
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no",
			args: args{
				readBuffer: utils.NewReadBufferByteBased([]byte{1, 2, 3, 4}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KnowsMeteringCommandTypeContainer(tt.args.readBuffer); got != tt.want {
				t.Errorf("KnowsMeteringCommandTypeContainer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKnowsSecurityCommandTypeContainer(t *testing.T) {
	type args struct {
		readBuffer utils.ReadBuffer
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "yes",
			args: args{
				readBuffer: utils.NewReadBufferByteBased([]byte{1, 2, 3, 4}),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KnowsSecurityCommandTypeContainer(tt.args.readBuffer); got != tt.want {
				t.Errorf("KnowsSecurityCommandTypeContainer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKnowsTelephonyCommandTypeContainer(t *testing.T) {
	type args struct {
		readBuffer utils.ReadBuffer
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no",
			args: args{
				readBuffer: utils.NewReadBufferByteBased([]byte{1, 2, 3, 4}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KnowsTelephonyCommandTypeContainer(tt.args.readBuffer); got != tt.want {
				t.Errorf("KnowsTelephonyCommandTypeContainer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKnowsTemperatureBroadcastCommandTypeContainer(t *testing.T) {
	type args struct {
		readBuffer utils.ReadBuffer
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no",
			args: args{
				readBuffer: utils.NewReadBufferByteBased([]byte{1, 2, 3, 4}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KnowsTemperatureBroadcastCommandTypeContainer(tt.args.readBuffer); got != tt.want {
				t.Errorf("KnowsTemperatureBroadcastCommandTypeContainer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKnowsTriggerControlCommandTypeContainer(t *testing.T) {
	type args struct {
		readBuffer utils.ReadBuffer
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "yes",
			args: args{
				readBuffer: utils.NewReadBufferByteBased([]byte{1, 2, 3, 4}),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KnowsTriggerControlCommandTypeContainer(tt.args.readBuffer); got != tt.want {
				t.Errorf("KnowsTriggerControlCommandTypeContainer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadAndValidateChecksum(t *testing.T) {
	type args struct {
		readBuffer utils.ReadBuffer
		message    spi.Message
		srchk      bool
	}
	tests := []struct {
		name    string
		args    args
		want    Checksum
		wantErr bool
	}{
		{
			name: "don't do it",
		},
		{
			name: "do it wrong message",
			args: args{
				readBuffer: utils.NewReadBufferByteBased([]byte{1, 2, 3, 4}),
				message:    nil,
				srchk:      true,
			},
			wantErr: true,
		},
		{
			name: "do it",
			args: args{
				readBuffer: utils.NewReadBufferByteBased([]byte("00")),
				message:    NewZoneStatus(12),
				srchk:      true,
			},
			want:    NewChecksum(0),
			wantErr: false,
		},
		{
			name: "do it wrong checksum",
			args: args{
				readBuffer: utils.NewReadBufferByteBased([]byte("AFFE")),
				message:    NewZoneStatus(12),
				srchk:      true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadAndValidateChecksum(tt.args.readBuffer, tt.args.message, tt.args.srchk)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadAndValidateChecksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("ReadAndValidateChecksum() got = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func TestReadCALData(t *testing.T) {
	type args struct {
		readBuffer utils.ReadBuffer
	}
	tests := []struct {
		name    string
		args    args
		want    CALData
		wantErr bool
	}{
		{
			name: "failing",
			args: args{
				readBuffer: utils.NewReadBufferByteBased([]byte{1, 2, 3, 4}),
			},
			wantErr: true,
		},
		{
			name: "cal data",
			args: args{
				readBuffer: utils.NewReadBufferByteBased([]byte("2101")),
			},
			want: NewCALDataIdentify(Attribute_Type, CALCommandTypeContainer_CALCommandIdentify, nil, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadCALData(tt.args.readBuffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadCALData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("ReadCALData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadCBusCommand(t *testing.T) {
	type args struct {
		readBuffer  utils.ReadBuffer
		cBusOptions CBusOptions
		srchk       bool
	}
	tests := []struct {
		name    string
		args    args
		want    CBusCommand
		wantErr bool
	}{
		{
			name: "failing",
			args: args{
				readBuffer: utils.NewReadBufferByteBased([]byte{1, 2, 3, 4}),
			},
			wantErr: true,
		},
		{
			name: "cbus command data",
			args: args{
				readBuffer: utils.NewReadBufferByteBased([]byte("46310900A400410600")),
			},
			want: NewCBusCommandPointToPoint(
				NewCBusPointToPointCommandIndirect(
					NewBridgeAddress(49),
					NewNetworkRoute(NewNetworkProtocolControlInformation(1, 1), nil),
					NewUnitAddress(0),
					12553,
					NewCALDataWrite(
						Parameter_UNKNOWN_01,
						65,
						NewParameterValueRaw([]byte{6, 00}, 2),
						CALCommandTypeContainer_CALCommandWrite_4Bytes,
						nil,
						nil,
					),
					nil,
				),
				NewCBusHeader(
					PriorityClass_Class3,
					false,
					0,
					DestinationAddressType_PointToPoint,
				),
				nil,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadCBusCommand(tt.args.readBuffer, tt.args.cBusOptions, tt.args.srchk)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadCBusCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("ReadCBusCommand() got = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func TestReadEncodedReply(t *testing.T) {
	type args struct {
		readBuffer     utils.ReadBuffer
		options        CBusOptions
		requestContext RequestContext
		srchk          bool
	}
	tests := []struct {
		name    string
		args    args
		want    EncodedReply
		wantErr bool
	}{
		{
			name: "failing",
			args: args{
				readBuffer: utils.NewReadBufferByteBased([]byte{1, 2, 3, 4}),
			},
			wantErr: true,
		},
		{
			name: "encoded reply",
			args: args{
				readBuffer:     utils.NewReadBufferByteBased([]byte("8510020000FF6A")),
				options:        NewCBusOptions(false, false, false, false, false, false, false, false, false),
				requestContext: NewRequestContext(false),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadEncodedReply(tt.args.readBuffer, tt.args.options, tt.args.requestContext, tt.args.srchk)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadEncodedReply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("ReadEncodedReply() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteCALData(t *testing.T) {
	type args struct {
		writeBuffer utils.WriteBuffer
		calData     CALData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "write something",
			args: args{
				writeBuffer: utils.NewWriteBufferBoxBased(),
				calData:     NewCALDataReset(0, nil, nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteCALData(tt.args.writeBuffer, tt.args.calData); (err != nil) != tt.wantErr {
				t.Errorf("WriteCALData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWriteCBusCommand(t *testing.T) {
	type args struct {
		writeBuffer utils.WriteBuffer
		cbusCommand CBusCommand
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "write something",
			args: args{
				writeBuffer: utils.NewWriteBufferBoxBased(),
				cbusCommand: NewCBusCommandDeviceManagement(0, 0, NewCBusHeader(0, false, 0, 0), nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteCBusCommand(tt.args.writeBuffer, tt.args.cbusCommand); (err != nil) != tt.wantErr {
				t.Errorf("WriteCBusCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWriteEncodedReply(t *testing.T) {
	type args struct {
		writeBuffer  utils.WriteBuffer
		encodedReply EncodedReply
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "write something",
			args: args{
				writeBuffer: utils.NewWriteBufferBoxBased(),
				encodedReply: NewEncodedReplyCALReply(
					NewCALReplyShort(
						0,
						NewCALDataReset(
							0,
							nil,
							NewRequestContext(false),
						),
						nil,
						NewRequestContext(false),
					),
					0,
					nil,
					NewRequestContext(false),
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteEncodedReply(tt.args.writeBuffer, tt.args.encodedReply); (err != nil) != tt.wantErr {
				t.Errorf("WriteEncodedReply() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_encodeHexUpperCase(t *testing.T) {
	type args struct {
		dst []byte
		src []byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "encode nothing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encodeHexUpperCase(tt.args.dst, tt.args.src); got != tt.want {
				t.Errorf("encodeHexUpperCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findHexEnd(t *testing.T) {
	type args struct {
		readBuffer utils.ReadBuffer
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "sal end",
			args: args{readBuffer: utils.NewReadBufferByteBased([]byte("0500380022AFE012"))},
			want: 16,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findHexEnd(tt.args.readBuffer); got != tt.want {
				t.Errorf("findHexEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getChecksum(t *testing.T) {
	type args struct {
		message spi.Message
	}
	tests := []struct {
		name    string
		args    args
		want    byte
		wantErr bool
	}{
		{
			name: "get it",
			args: args{
				message: NewZoneStatus(0),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getChecksum(tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("getChecksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getChecksum() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readBytesFromHex(t *testing.T) {
	type args struct {
		logicalName string
		readBuffer  utils.ReadBuffer
		srchk       bool
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "read it",
			args: args{
				readBuffer: utils.NewReadBufferByteBased([]byte{1, 2, 3, 4}),
			},
			wantErr: true,
		},
		// TODO: add more tests
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readBytesFromHex(tt.args.logicalName, tt.args.readBuffer, tt.args.srchk)
			if (err != nil) != tt.wantErr {
				t.Errorf("readBytesFromHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("readBytesFromHex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_writeSerializableToHex(t *testing.T) {
	type args struct {
		logicalName  string
		writeBuffer  utils.WriteBuffer
		serializable utils.Serializable
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "write it",
			args: args{
				writeBuffer:  utils.NewWriteBufferBoxBased(),
				serializable: NewZoneStatus(0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeSerializableToHex(tt.args.logicalName, tt.args.writeBuffer, tt.args.serializable); (err != nil) != tt.wantErr {
				t.Errorf("writeSerializableToHex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_writeToHex(t *testing.T) {
	type args struct {
		logicalName  string
		writeBuffer  utils.WriteBuffer
		bytesToWrite []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "write it",
			args: args{
				writeBuffer:  utils.NewWriteBufferBoxBased(),
				bytesToWrite: []byte{1, 2, 3, 4},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeToHex(tt.args.logicalName, tt.args.writeBuffer, tt.args.bytesToWrite); (err != nil) != tt.wantErr {
				t.Errorf("writeToHex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
