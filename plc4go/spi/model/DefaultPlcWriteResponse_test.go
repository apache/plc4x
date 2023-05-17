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
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultPlcWriteResponse_GetRequest(t *testing.T) {
	type fields struct {
		request       apiModel.PlcWriteRequest
		responseCodes map[string]apiModel.PlcResponseCode
	}
	tests := []struct {
		name   string
		fields fields
		want   apiModel.PlcWriteRequest
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcWriteResponse{
				request:       tt.fields.request,
				responseCodes: tt.fields.responseCodes,
			}
			assert.Equalf(t, tt.want, d.GetRequest(), "GetRequest()")
		})
	}
}

func TestDefaultPlcWriteResponse_GetResponseCode(t *testing.T) {
	type fields struct {
		request       apiModel.PlcWriteRequest
		responseCodes map[string]apiModel.PlcResponseCode
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   apiModel.PlcResponseCode
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcWriteResponse{
				request:       tt.fields.request,
				responseCodes: tt.fields.responseCodes,
			}
			assert.Equalf(t, tt.want, d.GetResponseCode(tt.args.name), "GetResponseCode(%v)", tt.args.name)
		})
	}
}

func TestDefaultPlcWriteResponse_GetTagNames(t *testing.T) {
	type fields struct {
		request       apiModel.PlcWriteRequest
		responseCodes map[string]apiModel.PlcResponseCode
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcWriteResponse{
				request:       tt.fields.request,
				responseCodes: tt.fields.responseCodes,
			}
			assert.Equalf(t, tt.want, d.GetTagNames(), "GetTagNames()")
		})
	}
}

func TestDefaultPlcWriteResponse_IsAPlcMessage(t *testing.T) {
	type fields struct {
		request       apiModel.PlcWriteRequest
		responseCodes map[string]apiModel.PlcResponseCode
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcWriteResponse{
				request:       tt.fields.request,
				responseCodes: tt.fields.responseCodes,
			}
			assert.Equalf(t, tt.want, d.IsAPlcMessage(), "IsAPlcMessage()")
		})
	}
}

func TestNewDefaultPlcWriteResponse(t *testing.T) {
	type args struct {
		request       apiModel.PlcWriteRequest
		responseCodes map[string]apiModel.PlcResponseCode
	}
	tests := []struct {
		name string
		args args
		want apiModel.PlcWriteResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDefaultPlcWriteResponse(tt.args.request, tt.args.responseCodes), "NewDefaultPlcWriteResponse(%v, %v)", tt.args.request, tt.args.responseCodes)
		})
	}
}
