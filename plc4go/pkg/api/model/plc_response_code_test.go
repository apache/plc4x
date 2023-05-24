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
	"context"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"reflect"
	"testing"
)

func TestPlcResponseCode_GetName(t *testing.T) {
	tests := []struct {
		name string
		m    PlcResponseCode
		want string
	}{
		{m: PlcResponseCode_OK, want: "OK"},
		{m: PlcResponseCode_NOT_FOUND, want: "NOT_FOUND"},
		{m: PlcResponseCode_ACCESS_DENIED, want: "ACCESS_DENIED"},
		{m: PlcResponseCode_INVALID_ADDRESS, want: "INVALID_ADDRESS"},
		{m: PlcResponseCode_INVALID_DATATYPE, want: "INVALID_DATATYPE"},
		{m: PlcResponseCode_INVALID_DATA, want: "INVALID_DATA"},
		{m: PlcResponseCode_INTERNAL_ERROR, want: "INTERNAL_ERROR"},
		{m: PlcResponseCode_REMOTE_BUSY, want: "REMOTE_BUSY"},
		{m: PlcResponseCode_REMOTE_ERROR, want: "REMOTE_ERROR"},
		{m: PlcResponseCode_UNSUPPORTED, want: "UNSUPPORTED"},
		{m: PlcResponseCode_RESPONSE_PENDING, want: "RESPONSE_PENDING"},
		{m: PlcResponseCode_REQUEST_TIMEOUT, want: "REQUEST_TIMEOUT"},
		{m: 0xff},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.GetName(); got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlcResponseCode_Serialize(t *testing.T) {
	tests := []struct {
		name    string
		m       PlcResponseCode
		want    []byte
		wantErr bool
	}{
		{
			name: "serialize it",
			want: []byte{0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Serialize()
			if (err != nil) != tt.wantErr {
				t.Errorf("Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Serialize() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlcResponseCode_SerializeWithWriteBuffer(t *testing.T) {
	type args struct {
		ctx         context.Context
		writeBuffer utils.WriteBuffer
	}
	tests := []struct {
		name    string
		m       PlcResponseCode
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.SerializeWithWriteBuffer(tt.args.ctx, tt.args.writeBuffer); (err != nil) != tt.wantErr {
				t.Errorf("SerializeWithWriteBuffer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlcResponseCode_String(t *testing.T) {
	tests := []struct {
		name string
		m    PlcResponseCode
		want string
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
