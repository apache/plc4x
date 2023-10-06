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
	"testing"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"

	"github.com/stretchr/testify/assert"
)

func TestDefaultPlcTagRequest_GetTag(t *testing.T) {
	type fields struct {
		tags     map[string]apiModel.PlcTag
		tagNames []string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   apiModel.PlcTag
	}{
		{
			name: "get it (not found)",
		},
		{
			name: "get it",
			fields: fields{
				tags: map[string]apiModel.PlcTag{
					"something": nil,
				},
			},
			args: args{
				name: "something",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcTagRequest{
				tags:     tt.fields.tags,
				tagNames: tt.fields.tagNames,
			}
			assert.Equalf(t, tt.want, d.GetTag(tt.args.name), "GetTag(%v)", tt.args.name)
		})
	}
}

func TestDefaultPlcTagRequest_GetTagNames(t *testing.T) {
	type fields struct {
		tags     map[string]apiModel.PlcTag
		tagNames []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "get em",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcTagRequest{
				tags:     tt.fields.tags,
				tagNames: tt.fields.tagNames,
			}
			assert.Equalf(t, tt.want, d.GetTagNames(), "GetTagNames()")
		})
	}
}

func TestDefaultPlcTagRequest_IsAPlcMessage(t *testing.T) {
	type fields struct {
		tags     map[string]apiModel.PlcTag
		tagNames []string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "it is",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcTagRequest{
				tags:     tt.fields.tags,
				tagNames: tt.fields.tagNames,
			}
			assert.Equalf(t, tt.want, d.IsAPlcMessage(), "IsAPlcMessage()")
		})
	}
}

func TestNewDefaultPlcTagRequest(t *testing.T) {
	type args struct {
		tags     map[string]apiModel.PlcTag
		tagNames []string
	}
	tests := []struct {
		name string
		args args
		want *DefaultPlcTagRequest
	}{
		{
			name: "create it",
			want: &DefaultPlcTagRequest{
				tags:     func() map[string]apiModel.PlcTag { return nil }(),
				tagNames: func() []string { return nil }(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDefaultPlcTagRequest(tt.args.tags, tt.args.tagNames), "NewDefaultPlcTagRequest(%v, %v)", tt.args.tags, tt.args.tagNames)
		})
	}
}
