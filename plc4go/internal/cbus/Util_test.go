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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
)

func TestCreateRequestContext(t *testing.T) {
	type args struct {
		cBusMessage model.CBusMessage
	}
	tests := []struct {
		name string
		args args
		want model.RequestContext
	}{
		{
			name: "just call it",
			want: readWriteModel.NewRequestContext(false),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, CreateRequestContext(tt.args.cBusMessage), "CreateRequestContext(%v)", tt.args.cBusMessage)
		})
	}
}

func TestCreateRequestContextWithInfoCallback(t *testing.T) {
	type args struct {
		cBusMessage  model.CBusMessage
		infoCallBack func(string)
	}
	tests := []struct {
		name string
		args args
		want model.RequestContext
	}{
		{
			name: "just call it",
			want: readWriteModel.NewRequestContext(false),
		},
		{
			name: "just call it with info callback",
			args: args{
				infoCallBack: func(_ string) {},
			},
			want: readWriteModel.NewRequestContext(false),
		},
		{
			name: "request context server direct command access",
			args: args{
				cBusMessage: readWriteModel.NewCBusMessageToServer(
					readWriteModel.NewRequestDirectCommandAccess(
						0,
						nil,
						nil,
						0,
						readWriteModel.NewRequestTermination(),
						readWriteModel.NewCALDataReset(readWriteModel.CALCommandTypeContainer_CALCommandReset, nil, nil),
						nil,
						nil,
					),
					nil,
					nil,
				),
			},
			want: readWriteModel.NewRequestContext(false),
		},
		{
			name: "request context server direct command access identify",
			args: args{
				cBusMessage: readWriteModel.NewCBusMessageToServer(
					readWriteModel.NewRequestDirectCommandAccess(
						0,
						nil,
						nil,
						0,
						readWriteModel.NewRequestTermination(),
						readWriteModel.NewCALDataIdentify(
							0,
							nil,
							0,
							nil,
						),
						nil,
						nil,
					),
					nil,
					nil,
				),
			},
			want: readWriteModel.NewRequestContext(true),
		},
		{
			name: "request context server command access",
			args: args{
				cBusMessage: readWriteModel.NewCBusMessageToServer(
					readWriteModel.NewRequestCommand(
						0,
						nil,
						nil,
						0,
						readWriteModel.NewRequestTermination(),
						readWriteModel.NewCBusCommandPointToPoint(
							readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class2, false, 0, readWriteModel.DestinationAddressType_PointToMultiPoint),
							readWriteModel.NewCBusPointToPointCommandDirect(
								0,
								readWriteModel.NewCALDataReset(readWriteModel.CALCommandTypeContainer_CALCommandReset, nil, nil),
								readWriteModel.NewUnitAddress(0),
								nil,
							),
							nil,
						),
						nil,
						nil,
						nil,
					),
					nil,
					nil,
				),
			},
			want: readWriteModel.NewRequestContext(false),
		},
		{
			name: "request context server command access identify",
			args: args{
				cBusMessage: readWriteModel.NewCBusMessageToServer(
					readWriteModel.NewRequestCommand(
						0,
						nil,
						nil,
						0,
						readWriteModel.NewRequestTermination(),
						readWriteModel.NewCBusCommandPointToPoint(
							readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class2, false, 0, readWriteModel.DestinationAddressType_PointToMultiPoint),
							readWriteModel.NewCBusPointToPointCommandDirect(
								0,
								readWriteModel.NewCALDataIdentify(
									0,
									nil,
									0,
									nil,
								),
								readWriteModel.NewUnitAddress(0),
								nil,
							),
							nil,
						),
						nil,
						nil,
						nil,
					),
					nil,
					nil,
				),
			},
			want: readWriteModel.NewRequestContext(true),
		},
		{
			name: "request context server direct command access obsolete",
			args: args{
				cBusMessage: readWriteModel.NewCBusMessageToServer(
					readWriteModel.NewRequestObsolete(
						0,
						nil,
						nil,
						0,
						readWriteModel.NewRequestTermination(),
						readWriteModel.NewCALDataReset(readWriteModel.CALCommandTypeContainer_CALCommandReset, nil, nil),
						nil,
						nil,
					),
					nil,
					nil,
				),
			},
			want: readWriteModel.NewRequestContext(false),
		},
		{
			name: "request context server direct command access identify obsolete",
			args: args{
				cBusMessage: readWriteModel.NewCBusMessageToServer(
					readWriteModel.NewRequestObsolete(
						0,
						nil,
						nil,
						0,
						readWriteModel.NewRequestTermination(),
						readWriteModel.NewCALDataIdentify(
							0,
							nil,
							0,
							nil,
						),
						nil,
						nil,
					),
					nil,
					nil,
				),
			},
			want: readWriteModel.NewRequestContext(true),
		},
		{
			name: "request context server direct command access identify obsolete",
			args: args{
				cBusMessage: readWriteModel.NewCBusMessageToClient(
					readWriteModel.NewReplyOrConfirmationConfirmation(
						0,
						readWriteModel.NewConfirmation(
							readWriteModel.NewAlpha(0),
							nil,
							readWriteModel.ConfirmationType_CONFIRMATION_SUCCESSFUL,
						),
						nil,
						nil,
						nil,
					),
					nil,
					nil,
				),
			},
			want: readWriteModel.NewRequestContext(false),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, CreateRequestContextWithInfoCallback(tt.args.cBusMessage, tt.args.infoCallBack), "CreateRequestContextWithInfoCallback(%v, fun())", tt.args.cBusMessage)
		})
	}
}
