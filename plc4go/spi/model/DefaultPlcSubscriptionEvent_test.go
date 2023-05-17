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
	"github.com/stretchr/testify/mock"
	"testing"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"

	"github.com/stretchr/testify/assert"
)

func TestDefaultPlcSubscriptionEvent_GetAddress(t *testing.T) {
	type fields struct {
		DefaultPlcSubscriptionEventRequirements DefaultPlcSubscriptionEventRequirements
		values                                  map[string]*DefaultPlcSubscriptionEventItem
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			fields: fields{
				DefaultPlcSubscriptionEventRequirements: func() DefaultPlcSubscriptionEventRequirements {
					requirements := NewMockDefaultPlcSubscriptionEventRequirements(t)
					requirements.EXPECT().GetAddress(mock.Anything).Return("anything")
					return requirements
				}(),
			},
			name: "get it",
			want: "anything",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcSubscriptionEvent{
				DefaultPlcSubscriptionEventRequirements: tt.fields.DefaultPlcSubscriptionEventRequirements,
				values:                                  tt.fields.values,
			}
			assert.Equalf(t, tt.want, d.GetAddress(tt.args.name), "GetAddress(%v)", tt.args.name)
		})
	}
}

func TestDefaultPlcSubscriptionEvent_GetInterval(t *testing.T) {
	type fields struct {
		DefaultPlcSubscriptionEventRequirements DefaultPlcSubscriptionEventRequirements
		values                                  map[string]*DefaultPlcSubscriptionEventItem
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   time.Duration
	}{
		{
			name: "get it (not found)",
			want: -1,
		},
		{
			name: "get it",
			fields: fields{
				values: map[string]*DefaultPlcSubscriptionEventItem{
					"da field": {interval: 70},
				},
			},
			args: args{
				name: "da field",
			},
			want: 70,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcSubscriptionEvent{
				DefaultPlcSubscriptionEventRequirements: tt.fields.DefaultPlcSubscriptionEventRequirements,
				values:                                  tt.fields.values,
			}
			assert.Equalf(t, tt.want, d.GetInterval(tt.args.name), "GetInterval(%v)", tt.args.name)
		})
	}
}

func TestDefaultPlcSubscriptionEvent_GetResponseCode(t *testing.T) {
	type fields struct {
		DefaultPlcSubscriptionEventRequirements DefaultPlcSubscriptionEventRequirements
		values                                  map[string]*DefaultPlcSubscriptionEventItem
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
		{
			name: "get it (not found)",
			want: apiModel.PlcResponseCode_NOT_FOUND,
		},
		{
			name: "get it",
			fields: fields{
				values: map[string]*DefaultPlcSubscriptionEventItem{
					"da field": {code: apiModel.PlcResponseCode_NOT_FOUND},
				},
			},
			args: args{
				name: "da field",
			},
			want: apiModel.PlcResponseCode_NOT_FOUND,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcSubscriptionEvent{
				DefaultPlcSubscriptionEventRequirements: tt.fields.DefaultPlcSubscriptionEventRequirements,
				values:                                  tt.fields.values,
			}
			assert.Equalf(t, tt.want, d.GetResponseCode(tt.args.name), "GetResponseCode(%v)", tt.args.name)
		})
	}
}

func TestDefaultPlcSubscriptionEvent_GetSource(t *testing.T) {
	type fields struct {
		DefaultPlcSubscriptionEventRequirements DefaultPlcSubscriptionEventRequirements
		values                                  map[string]*DefaultPlcSubscriptionEventItem
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "get it (not found)",
			fields: fields{
				DefaultPlcSubscriptionEventRequirements: func() DefaultPlcSubscriptionEventRequirements {
					requirements := NewMockDefaultPlcSubscriptionEventRequirements(t)
					requirements.EXPECT().GetAddress(mock.Anything).Return("")
					return requirements
				}(),
			},
			want: "",
		},
		{
			name: "get it",
			fields: fields{
				DefaultPlcSubscriptionEventRequirements: func() DefaultPlcSubscriptionEventRequirements {
					requirements := NewMockDefaultPlcSubscriptionEventRequirements(t)
					requirements.EXPECT().GetAddress(mock.Anything).Return("something")
					return requirements
				}(),
			},
			args: args{
				name: "da field",
			},
			want: "something",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcSubscriptionEvent{
				DefaultPlcSubscriptionEventRequirements: tt.fields.DefaultPlcSubscriptionEventRequirements,
				values:                                  tt.fields.values,
			}
			assert.Equalf(t, tt.want, d.GetSource(tt.args.name), "GetSource(%v)", tt.args.name)
		})
	}
}

func TestDefaultPlcSubscriptionEvent_GetTag(t *testing.T) {
	type fields struct {
		DefaultPlcSubscriptionEventRequirements DefaultPlcSubscriptionEventRequirements
		values                                  map[string]*DefaultPlcSubscriptionEventItem
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
			want: nil,
		},
		{
			name: "get it",
			fields: fields{
				values: map[string]*DefaultPlcSubscriptionEventItem{
					"da field": {tag: nil},
				},
			},
			args: args{
				name: "da field",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcSubscriptionEvent{
				DefaultPlcSubscriptionEventRequirements: tt.fields.DefaultPlcSubscriptionEventRequirements,
				values:                                  tt.fields.values,
			}
			assert.Equalf(t, tt.want, d.GetTag(tt.args.name), "GetTag(%v)", tt.args.name)
		})
	}
}

func TestDefaultPlcSubscriptionEvent_GetTagNames(t *testing.T) {
	type fields struct {
		DefaultPlcSubscriptionEventRequirements DefaultPlcSubscriptionEventRequirements
		values                                  map[string]*DefaultPlcSubscriptionEventItem
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "get it (not found)",
			want: nil,
		},
		{
			name: "get it",
			fields: fields{
				values: map[string]*DefaultPlcSubscriptionEventItem{
					"da field":  {tag: nil},
					"da field2": {tag: nil},
				},
			},
			want: []string{"da field", "da field2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcSubscriptionEvent{
				DefaultPlcSubscriptionEventRequirements: tt.fields.DefaultPlcSubscriptionEventRequirements,
				values:                                  tt.fields.values,
			}
			assert.Equalf(t, tt.want, d.GetTagNames(), "GetTagNames()")
		})
	}
}

func TestDefaultPlcSubscriptionEvent_GetType(t *testing.T) {
	type fields struct {
		DefaultPlcSubscriptionEventRequirements DefaultPlcSubscriptionEventRequirements
		values                                  map[string]*DefaultPlcSubscriptionEventItem
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   SubscriptionType
	}{
		{
			name: "get it (not found)",
			want: 0,
		},
		{
			name: "get it",
			fields: fields{
				values: map[string]*DefaultPlcSubscriptionEventItem{
					"da field": {subscriptionType: SubscriptionChangeOfState},
				},
			},
			args: args{name: "da field"},
			want: SubscriptionChangeOfState,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcSubscriptionEvent{
				DefaultPlcSubscriptionEventRequirements: tt.fields.DefaultPlcSubscriptionEventRequirements,
				values:                                  tt.fields.values,
			}
			assert.Equalf(t, tt.want, d.GetType(tt.args.name), "GetType(%v)", tt.args.name)
		})
	}
}

func TestDefaultPlcSubscriptionEvent_GetValue(t *testing.T) {
	type fields struct {
		DefaultPlcSubscriptionEventRequirements DefaultPlcSubscriptionEventRequirements
		values                                  map[string]*DefaultPlcSubscriptionEventItem
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   apiValues.PlcValue
	}{
		{
			name: "get it (not found)",
			want: spiValues.PlcNull{},
		},
		{
			name: "get it",
			fields: fields{
				values: map[string]*DefaultPlcSubscriptionEventItem{
					"da field": {value: spiValues.NewPlcSTRING("yeah")},
				},
			},
			args: args{name: "da field"},
			want: spiValues.NewPlcSTRING("yeah"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcSubscriptionEvent{
				DefaultPlcSubscriptionEventRequirements: tt.fields.DefaultPlcSubscriptionEventRequirements,
				values:                                  tt.fields.values,
			}
			assert.Equalf(t, tt.want, d.GetValue(tt.args.name), "GetValue(%v)", tt.args.name)
		})
	}
}

func TestDefaultPlcSubscriptionEvent_IsAPlcMessage(t *testing.T) {
	type fields struct {
		DefaultPlcSubscriptionEventRequirements DefaultPlcSubscriptionEventRequirements
		values                                  map[string]*DefaultPlcSubscriptionEventItem
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
			d := &DefaultPlcSubscriptionEvent{
				DefaultPlcSubscriptionEventRequirements: tt.fields.DefaultPlcSubscriptionEventRequirements,
				values:                                  tt.fields.values,
			}
			assert.Equalf(t, tt.want, d.IsAPlcMessage(), "IsAPlcMessage()")
		})
	}
}

func TestNewDefaultPlcSubscriptionEvent(t *testing.T) {
	type args struct {
		defaultPlcSubscriptionEventRequirements DefaultPlcSubscriptionEventRequirements
		tags                                    map[string]apiModel.PlcTag
		types                                   map[string]SubscriptionType
		intervals                               map[string]time.Duration
		responseCodes                           map[string]apiModel.PlcResponseCode
		values                                  map[string]apiValues.PlcValue
	}
	tests := []struct {
		name string
		args args
		want apiModel.PlcSubscriptionEvent
	}{
		{
			name: "create it",
			want: &DefaultPlcSubscriptionEvent{values: map[string]*DefaultPlcSubscriptionEventItem{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDefaultPlcSubscriptionEvent(tt.args.defaultPlcSubscriptionEventRequirements, tt.args.tags, tt.args.types, tt.args.intervals, tt.args.responseCodes, tt.args.values), "NewDefaultPlcSubscriptionEvent(%v, %v, %v, %v, %v, %v)", tt.args.defaultPlcSubscriptionEventRequirements, tt.args.tags, tt.args.types, tt.args.intervals, tt.args.responseCodes, tt.args.values)
		})
	}
}
