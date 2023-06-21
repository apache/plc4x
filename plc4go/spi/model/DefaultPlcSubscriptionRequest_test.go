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
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"

	"github.com/stretchr/testify/assert"
)

func TestDefaultPlcSubscriptionRequestBuilder_AddChangeOfStateTag(t *testing.T) {
	type fields struct {
		subscriber             spi.PlcSubscriber
		tagHandler             spi.PlcTagHandler
		valueHandler           spi.PlcValueHandler
		tagNames               []string
		tagAddresses           map[string]string
		tags                   map[string]apiModel.PlcTag
		types                  map[string]SubscriptionType
		intervals              map[string]time.Duration
		preRegisteredConsumers map[string][]apiModel.PlcSubscriptionEventConsumer
	}
	type args struct {
		name string
		tag  apiModel.PlcTag
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   apiModel.PlcSubscriptionRequestBuilder
	}{
		{
			name: "add it",
			fields: fields{
				tagAddresses:           map[string]string{},
				tags:                   map[string]apiModel.PlcTag{},
				types:                  map[string]SubscriptionType{},
				intervals:              map[string]time.Duration{},
				preRegisteredConsumers: map[string][]apiModel.PlcSubscriptionEventConsumer{},
			},
			want: &DefaultPlcSubscriptionRequestBuilder{
				tagNames:               []string{""},
				tagAddresses:           map[string]string{},
				tags:                   map[string]apiModel.PlcTag{"": nil},
				types:                  map[string]SubscriptionType{"": 2},
				intervals:              map[string]time.Duration{},
				preRegisteredConsumers: map[string][]apiModel.PlcSubscriptionEventConsumer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcSubscriptionRequestBuilder{
				subscriber:             tt.fields.subscriber,
				tagHandler:             tt.fields.tagHandler,
				valueHandler:           tt.fields.valueHandler,
				tagNames:               tt.fields.tagNames,
				tagAddresses:           tt.fields.tagAddresses,
				tags:                   tt.fields.tags,
				types:                  tt.fields.types,
				intervals:              tt.fields.intervals,
				preRegisteredConsumers: tt.fields.preRegisteredConsumers,
			}
			assert.Equalf(t, tt.want, d.AddChangeOfStateTag(tt.args.name, tt.args.tag), "AddChangeOfStateTag(%v, %v)", tt.args.name, tt.args.tag)
		})
	}
}

func TestDefaultPlcSubscriptionRequestBuilder_AddChangeOfStateTagAddress(t *testing.T) {
	type fields struct {
		subscriber             spi.PlcSubscriber
		tagHandler             spi.PlcTagHandler
		valueHandler           spi.PlcValueHandler
		tagNames               []string
		tagAddresses           map[string]string
		tags                   map[string]apiModel.PlcTag
		types                  map[string]SubscriptionType
		intervals              map[string]time.Duration
		preRegisteredConsumers map[string][]apiModel.PlcSubscriptionEventConsumer
	}
	type args struct {
		name       string
		tagAddress string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   apiModel.PlcSubscriptionRequestBuilder
	}{
		{
			name: "add it",
			fields: fields{
				tagAddresses:           map[string]string{},
				tags:                   map[string]apiModel.PlcTag{},
				types:                  map[string]SubscriptionType{},
				intervals:              map[string]time.Duration{},
				preRegisteredConsumers: map[string][]apiModel.PlcSubscriptionEventConsumer{},
			},
			want: &DefaultPlcSubscriptionRequestBuilder{
				tagNames:               []string{""},
				tagAddresses:           map[string]string{"": ""},
				tags:                   map[string]apiModel.PlcTag{},
				types:                  map[string]SubscriptionType{"": 2},
				intervals:              map[string]time.Duration{},
				preRegisteredConsumers: map[string][]apiModel.PlcSubscriptionEventConsumer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcSubscriptionRequestBuilder{
				subscriber:             tt.fields.subscriber,
				tagHandler:             tt.fields.tagHandler,
				valueHandler:           tt.fields.valueHandler,
				tagNames:               tt.fields.tagNames,
				tagAddresses:           tt.fields.tagAddresses,
				tags:                   tt.fields.tags,
				types:                  tt.fields.types,
				intervals:              tt.fields.intervals,
				preRegisteredConsumers: tt.fields.preRegisteredConsumers,
			}
			assert.Equalf(t, tt.want, d.AddChangeOfStateTagAddress(tt.args.name, tt.args.tagAddress), "AddChangeOfStateTagAddress(%v, %v)", tt.args.name, tt.args.tagAddress)
		})
	}
}

func TestDefaultPlcSubscriptionRequestBuilder_AddCyclicTag(t *testing.T) {
	type fields struct {
		subscriber             spi.PlcSubscriber
		tagHandler             spi.PlcTagHandler
		valueHandler           spi.PlcValueHandler
		tagNames               []string
		tagAddresses           map[string]string
		tags                   map[string]apiModel.PlcTag
		types                  map[string]SubscriptionType
		intervals              map[string]time.Duration
		preRegisteredConsumers map[string][]apiModel.PlcSubscriptionEventConsumer
	}
	type args struct {
		name     string
		tag      apiModel.PlcTag
		interval time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   apiModel.PlcSubscriptionRequestBuilder
	}{
		{
			name: "add it",
			fields: fields{
				tagAddresses:           map[string]string{},
				tags:                   map[string]apiModel.PlcTag{},
				types:                  map[string]SubscriptionType{},
				intervals:              map[string]time.Duration{},
				preRegisteredConsumers: map[string][]apiModel.PlcSubscriptionEventConsumer{},
			},
			want: &DefaultPlcSubscriptionRequestBuilder{
				tagNames:               []string{""},
				tagAddresses:           map[string]string{},
				tags:                   map[string]apiModel.PlcTag{"": nil},
				types:                  map[string]SubscriptionType{"": 1},
				intervals:              map[string]time.Duration{"": 0},
				preRegisteredConsumers: map[string][]apiModel.PlcSubscriptionEventConsumer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcSubscriptionRequestBuilder{
				subscriber:             tt.fields.subscriber,
				tagHandler:             tt.fields.tagHandler,
				valueHandler:           tt.fields.valueHandler,
				tagNames:               tt.fields.tagNames,
				tagAddresses:           tt.fields.tagAddresses,
				tags:                   tt.fields.tags,
				types:                  tt.fields.types,
				intervals:              tt.fields.intervals,
				preRegisteredConsumers: tt.fields.preRegisteredConsumers,
			}
			assert.Equalf(t, tt.want, d.AddCyclicTag(tt.args.name, tt.args.tag, tt.args.interval), "AddCyclicTag(%v, %v, %v)", tt.args.name, tt.args.tag, tt.args.interval)
		})
	}
}

func TestDefaultPlcSubscriptionRequestBuilder_AddCyclicTagAddress(t *testing.T) {
	type fields struct {
		subscriber             spi.PlcSubscriber
		tagHandler             spi.PlcTagHandler
		valueHandler           spi.PlcValueHandler
		tagNames               []string
		tagAddresses           map[string]string
		tags                   map[string]apiModel.PlcTag
		types                  map[string]SubscriptionType
		intervals              map[string]time.Duration
		preRegisteredConsumers map[string][]apiModel.PlcSubscriptionEventConsumer
	}
	type args struct {
		name       string
		tagAddress string
		interval   time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   apiModel.PlcSubscriptionRequestBuilder
	}{
		{
			name: "add it",
			fields: fields{
				tagAddresses:           map[string]string{},
				tags:                   map[string]apiModel.PlcTag{},
				types:                  map[string]SubscriptionType{},
				intervals:              map[string]time.Duration{},
				preRegisteredConsumers: map[string][]apiModel.PlcSubscriptionEventConsumer{},
			},
			want: &DefaultPlcSubscriptionRequestBuilder{
				tagNames:               []string{""},
				tagAddresses:           map[string]string{"": ""},
				tags:                   map[string]apiModel.PlcTag{},
				types:                  map[string]SubscriptionType{"": 1},
				intervals:              map[string]time.Duration{"": 0},
				preRegisteredConsumers: map[string][]apiModel.PlcSubscriptionEventConsumer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcSubscriptionRequestBuilder{
				subscriber:             tt.fields.subscriber,
				tagHandler:             tt.fields.tagHandler,
				valueHandler:           tt.fields.valueHandler,
				tagNames:               tt.fields.tagNames,
				tagAddresses:           tt.fields.tagAddresses,
				tags:                   tt.fields.tags,
				types:                  tt.fields.types,
				intervals:              tt.fields.intervals,
				preRegisteredConsumers: tt.fields.preRegisteredConsumers,
			}
			assert.Equalf(t, tt.want, d.AddCyclicTagAddress(tt.args.name, tt.args.tagAddress, tt.args.interval), "AddCyclicTagAddress(%v, %v, %v)", tt.args.name, tt.args.tagAddress, tt.args.interval)
		})
	}
}

func TestDefaultPlcSubscriptionRequestBuilder_AddEventTag(t *testing.T) {
	type fields struct {
		subscriber             spi.PlcSubscriber
		tagHandler             spi.PlcTagHandler
		valueHandler           spi.PlcValueHandler
		tagNames               []string
		tagAddresses           map[string]string
		tags                   map[string]apiModel.PlcTag
		types                  map[string]SubscriptionType
		intervals              map[string]time.Duration
		preRegisteredConsumers map[string][]apiModel.PlcSubscriptionEventConsumer
	}
	type args struct {
		name string
		tag  apiModel.PlcTag
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   apiModel.PlcSubscriptionRequestBuilder
	}{
		{
			name: "add it",
			fields: fields{
				tagAddresses:           map[string]string{},
				tags:                   map[string]apiModel.PlcTag{},
				types:                  map[string]SubscriptionType{},
				intervals:              map[string]time.Duration{},
				preRegisteredConsumers: map[string][]apiModel.PlcSubscriptionEventConsumer{},
			},
			want: &DefaultPlcSubscriptionRequestBuilder{
				tagNames:               []string{""},
				tagAddresses:           map[string]string{},
				tags:                   map[string]apiModel.PlcTag{"": nil},
				types:                  map[string]SubscriptionType{"": 3},
				intervals:              map[string]time.Duration{},
				preRegisteredConsumers: map[string][]apiModel.PlcSubscriptionEventConsumer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcSubscriptionRequestBuilder{
				subscriber:             tt.fields.subscriber,
				tagHandler:             tt.fields.tagHandler,
				valueHandler:           tt.fields.valueHandler,
				tagNames:               tt.fields.tagNames,
				tagAddresses:           tt.fields.tagAddresses,
				tags:                   tt.fields.tags,
				types:                  tt.fields.types,
				intervals:              tt.fields.intervals,
				preRegisteredConsumers: tt.fields.preRegisteredConsumers,
			}
			assert.Equalf(t, tt.want, d.AddEventTag(tt.args.name, tt.args.tag), "AddEventTag(%v, %v)", tt.args.name, tt.args.tag)
		})
	}
}

func TestDefaultPlcSubscriptionRequestBuilder_AddEventTagAddress(t *testing.T) {
	type fields struct {
		subscriber             spi.PlcSubscriber
		tagHandler             spi.PlcTagHandler
		valueHandler           spi.PlcValueHandler
		tagNames               []string
		tagAddresses           map[string]string
		tags                   map[string]apiModel.PlcTag
		types                  map[string]SubscriptionType
		intervals              map[string]time.Duration
		preRegisteredConsumers map[string][]apiModel.PlcSubscriptionEventConsumer
	}
	type args struct {
		name       string
		tagAddress string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   apiModel.PlcSubscriptionRequestBuilder
	}{
		{
			name: "add it",
			fields: fields{
				tagAddresses:           map[string]string{},
				tags:                   map[string]apiModel.PlcTag{},
				types:                  map[string]SubscriptionType{},
				intervals:              map[string]time.Duration{},
				preRegisteredConsumers: map[string][]apiModel.PlcSubscriptionEventConsumer{},
			},
			want: &DefaultPlcSubscriptionRequestBuilder{
				tagNames:               []string{""},
				tagAddresses:           map[string]string{"": ""},
				tags:                   map[string]apiModel.PlcTag{},
				types:                  map[string]SubscriptionType{"": 3},
				intervals:              map[string]time.Duration{},
				preRegisteredConsumers: map[string][]apiModel.PlcSubscriptionEventConsumer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcSubscriptionRequestBuilder{
				subscriber:             tt.fields.subscriber,
				tagHandler:             tt.fields.tagHandler,
				valueHandler:           tt.fields.valueHandler,
				tagNames:               tt.fields.tagNames,
				tagAddresses:           tt.fields.tagAddresses,
				tags:                   tt.fields.tags,
				types:                  tt.fields.types,
				intervals:              tt.fields.intervals,
				preRegisteredConsumers: tt.fields.preRegisteredConsumers,
			}
			assert.Equalf(t, tt.want, d.AddEventTagAddress(tt.args.name, tt.args.tagAddress), "AddEventTagAddress(%v, %v)", tt.args.name, tt.args.tagAddress)
		})
	}
}

func TestDefaultPlcSubscriptionRequestBuilder_AddPreRegisteredConsumer(t *testing.T) {
	type fields struct {
		subscriber             spi.PlcSubscriber
		tagHandler             spi.PlcTagHandler
		valueHandler           spi.PlcValueHandler
		tagNames               []string
		tagAddresses           map[string]string
		tags                   map[string]apiModel.PlcTag
		types                  map[string]SubscriptionType
		intervals              map[string]time.Duration
		preRegisteredConsumers map[string][]apiModel.PlcSubscriptionEventConsumer
	}
	type args struct {
		name     string
		consumer apiModel.PlcSubscriptionEventConsumer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   apiModel.PlcSubscriptionRequestBuilder
	}{
		{
			name: "add it",
			fields: fields{
				preRegisteredConsumers: map[string][]apiModel.PlcSubscriptionEventConsumer{},
			},
			want: &DefaultPlcSubscriptionRequestBuilder{
				preRegisteredConsumers: map[string][]apiModel.PlcSubscriptionEventConsumer{
					"": {func() apiModel.PlcSubscriptionEventConsumer { return nil }()},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcSubscriptionRequestBuilder{
				subscriber:             tt.fields.subscriber,
				tagHandler:             tt.fields.tagHandler,
				valueHandler:           tt.fields.valueHandler,
				tagNames:               tt.fields.tagNames,
				tagAddresses:           tt.fields.tagAddresses,
				tags:                   tt.fields.tags,
				types:                  tt.fields.types,
				intervals:              tt.fields.intervals,
				preRegisteredConsumers: tt.fields.preRegisteredConsumers,
			}
			assert.Equalf(t, tt.want, d.AddPreRegisteredConsumer(tt.args.name, tt.args.consumer), "AddPreRegisteredConsumer(%v, func(%t))", tt.args.name, tt.args.consumer != nil)
		})
	}
}

func TestDefaultPlcSubscriptionRequestBuilder_Build(t *testing.T) {
	type fields struct {
		subscriber             spi.PlcSubscriber
		tagHandler             spi.PlcTagHandler
		valueHandler           spi.PlcValueHandler
		tagNames               []string
		tagAddresses           map[string]string
		tags                   map[string]apiModel.PlcTag
		types                  map[string]SubscriptionType
		intervals              map[string]time.Duration
		preRegisteredConsumers map[string][]apiModel.PlcSubscriptionEventConsumer
	}
	tests := []struct {
		name    string
		fields  fields
		setup   func(t *testing.T, fields *fields)
		want    apiModel.PlcSubscriptionRequest
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "build it",
			want: &DefaultPlcSubscriptionRequest{
				DefaultPlcTagRequest: NewDefaultPlcTagRequest(nil, nil),
			},
			wantErr: assert.NoError,
		},
		{
			name: "build it (with addresses)",
			fields: fields{
				tagNames:     []string{"a"},
				tagAddresses: map[string]string{"a": ""},
				tags:         map[string]apiModel.PlcTag{},
			},
			setup: func(t *testing.T, fields *fields) {
				handler := NewMockPlcTagHandler(t)
				handler.EXPECT().ParseTag(mock.Anything).Return(nil, nil)
				fields.tagHandler = handler
			},
			want: &DefaultPlcSubscriptionRequest{
				DefaultPlcTagRequest: NewDefaultPlcTagRequest(map[string]apiModel.PlcTag{"a": nil}, []string{"a"}),
			},
			wantErr: assert.NoError,
		},
		{
			name: "build it (with addresses failing)",
			fields: fields{
				tagNames:     []string{"a"},
				tagAddresses: map[string]string{"a": ""},
				tags:         map[string]apiModel.PlcTag{},
			},
			setup: func(t *testing.T, fields *fields) {
				handler := NewMockPlcTagHandler(t)
				handler.EXPECT().ParseTag(mock.Anything).Return(nil, errors.New("nope"))
				fields.tagHandler = handler
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			d := &DefaultPlcSubscriptionRequestBuilder{
				subscriber:             tt.fields.subscriber,
				tagHandler:             tt.fields.tagHandler,
				valueHandler:           tt.fields.valueHandler,
				tagNames:               tt.fields.tagNames,
				tagAddresses:           tt.fields.tagAddresses,
				tags:                   tt.fields.tags,
				types:                  tt.fields.types,
				intervals:              tt.fields.intervals,
				preRegisteredConsumers: tt.fields.preRegisteredConsumers,
			}
			got, err := d.Build()
			if !tt.wantErr(t, err, fmt.Sprintf("Build()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "Build()")
		})
	}
}

func TestDefaultPlcSubscriptionRequest_Execute(t *testing.T) {
	type fields struct {
		DefaultPlcTagRequest   *DefaultPlcTagRequest
		types                  map[string]SubscriptionType
		intervals              map[string]time.Duration
		preRegisteredConsumers map[string][]apiModel.PlcSubscriptionEventConsumer
		subscriber             spi.PlcSubscriber
	}
	tests := []struct {
		name   string
		fields fields
		setup  func(t *testing.T, fields *fields)
		want   <-chan apiModel.PlcSubscriptionRequestResult
	}{
		{
			name: "execute it",
			setup: func(t *testing.T, fields *fields) {
				subscriber := NewMockPlcSubscriber(t)
				subscriber.EXPECT().Subscribe(mock.Anything, mock.Anything).Return(nil)
				fields.subscriber = subscriber
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			d := &DefaultPlcSubscriptionRequest{
				DefaultPlcTagRequest:   tt.fields.DefaultPlcTagRequest,
				types:                  tt.fields.types,
				intervals:              tt.fields.intervals,
				preRegisteredConsumers: tt.fields.preRegisteredConsumers,
				subscriber:             tt.fields.subscriber,
			}
			assert.Equalf(t, tt.want, d.Execute(), "Execute()")
		})
	}
}

func TestDefaultPlcSubscriptionRequest_ExecuteWithContext(t *testing.T) {
	type fields struct {
		DefaultPlcTagRequest   *DefaultPlcTagRequest
		types                  map[string]SubscriptionType
		intervals              map[string]time.Duration
		preRegisteredConsumers map[string][]apiModel.PlcSubscriptionEventConsumer
		subscriber             spi.PlcSubscriber
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func(t *testing.T, fields *fields)
		want   <-chan apiModel.PlcSubscriptionRequestResult
	}{
		{
			name: "execute it",
			setup: func(t *testing.T, fields *fields) {
				subscriber := NewMockPlcSubscriber(t)
				subscriber.EXPECT().Subscribe(mock.Anything, mock.Anything).Return(nil)
				fields.subscriber = subscriber
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			d := &DefaultPlcSubscriptionRequest{
				DefaultPlcTagRequest:   tt.fields.DefaultPlcTagRequest,
				types:                  tt.fields.types,
				intervals:              tt.fields.intervals,
				preRegisteredConsumers: tt.fields.preRegisteredConsumers,
				subscriber:             tt.fields.subscriber,
			}
			assert.Equalf(t, tt.want, d.ExecuteWithContext(tt.args.ctx), "ExecuteWithContext(%v)", tt.args.ctx)
		})
	}
}

func TestDefaultPlcSubscriptionRequest_GetInterval(t *testing.T) {
	type fields struct {
		DefaultPlcTagRequest   *DefaultPlcTagRequest
		types                  map[string]SubscriptionType
		intervals              map[string]time.Duration
		preRegisteredConsumers map[string][]apiModel.PlcSubscriptionEventConsumer
		subscriber             spi.PlcSubscriber
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
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcSubscriptionRequest{
				DefaultPlcTagRequest:   tt.fields.DefaultPlcTagRequest,
				types:                  tt.fields.types,
				intervals:              tt.fields.intervals,
				preRegisteredConsumers: tt.fields.preRegisteredConsumers,
				subscriber:             tt.fields.subscriber,
			}
			assert.Equalf(t, tt.want, d.GetInterval(tt.args.name), "GetInterval(%v)", tt.args.name)
		})
	}
}

func TestDefaultPlcSubscriptionRequest_GetPreRegisteredConsumers(t *testing.T) {
	type fields struct {
		DefaultPlcTagRequest   *DefaultPlcTagRequest
		types                  map[string]SubscriptionType
		intervals              map[string]time.Duration
		preRegisteredConsumers map[string][]apiModel.PlcSubscriptionEventConsumer
		subscriber             spi.PlcSubscriber
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []apiModel.PlcSubscriptionEventConsumer
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcSubscriptionRequest{
				DefaultPlcTagRequest:   tt.fields.DefaultPlcTagRequest,
				types:                  tt.fields.types,
				intervals:              tt.fields.intervals,
				preRegisteredConsumers: tt.fields.preRegisteredConsumers,
				subscriber:             tt.fields.subscriber,
			}
			assert.Equalf(t, tt.want, d.GetPreRegisteredConsumers(tt.args.name), "GetPreRegisteredConsumers(%v)", tt.args.name)
		})
	}
}

func TestDefaultPlcSubscriptionRequest_GetType(t *testing.T) {
	type fields struct {
		DefaultPlcTagRequest   *DefaultPlcTagRequest
		types                  map[string]SubscriptionType
		intervals              map[string]time.Duration
		preRegisteredConsumers map[string][]apiModel.PlcSubscriptionEventConsumer
		subscriber             spi.PlcSubscriber
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
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcSubscriptionRequest{
				DefaultPlcTagRequest:   tt.fields.DefaultPlcTagRequest,
				types:                  tt.fields.types,
				intervals:              tt.fields.intervals,
				preRegisteredConsumers: tt.fields.preRegisteredConsumers,
				subscriber:             tt.fields.subscriber,
			}
			assert.Equalf(t, tt.want, d.GetType(tt.args.name), "GetType(%v)", tt.args.name)
		})
	}
}

func TestNewDefaultPlcSubscriptionRequest(t *testing.T) {
	type args struct {
		subscriber             spi.PlcSubscriber
		tagNames               []string
		tags                   map[string]apiModel.PlcTag
		types                  map[string]SubscriptionType
		intervals              map[string]time.Duration
		preRegisteredConsumers map[string][]apiModel.PlcSubscriptionEventConsumer
	}
	tests := []struct {
		name string
		args args
		want apiModel.PlcSubscriptionRequest
	}{
		{
			name: "create it",
			want: &DefaultPlcSubscriptionRequest{
				DefaultPlcTagRequest: NewDefaultPlcTagRequest(nil, nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDefaultPlcSubscriptionRequest(tt.args.subscriber, tt.args.tagNames, tt.args.tags, tt.args.types, tt.args.intervals, tt.args.preRegisteredConsumers), "NewDefaultPlcSubscriptionRequest(%v, %v, %v, %v, %v, %v)", tt.args.subscriber, tt.args.tagNames, tt.args.tags, tt.args.types, tt.args.intervals, tt.args.preRegisteredConsumers)
		})
	}
}

func TestNewDefaultPlcSubscriptionRequestBuilder(t *testing.T) {
	type args struct {
		tagHandler   spi.PlcTagHandler
		valueHandler spi.PlcValueHandler
		subscriber   spi.PlcSubscriber
	}
	tests := []struct {
		name string
		args args
		want apiModel.PlcSubscriptionRequestBuilder
	}{
		{
			name: "create it",
			want: &DefaultPlcSubscriptionRequestBuilder{
				tagNames:               []string{},
				tagAddresses:           map[string]string{},
				tags:                   map[string]apiModel.PlcTag{},
				types:                  map[string]SubscriptionType{},
				intervals:              map[string]time.Duration{},
				preRegisteredConsumers: map[string][]apiModel.PlcSubscriptionEventConsumer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDefaultPlcSubscriptionRequestBuilder(tt.args.tagHandler, tt.args.valueHandler, tt.args.subscriber), "NewDefaultPlcSubscriptionRequestBuilder(%v, %v, %v)", tt.args.tagHandler, tt.args.valueHandler, tt.args.subscriber)
		})
	}
}

func TestSubscriptionType_String(t *testing.T) {
	tests := []struct {
		name string
		s    SubscriptionType
		want string
	}{
		{
			name: "SubscriptionCyclic",
			s:    SubscriptionCyclic,
			want: "SubscriptionCyclic",
		},
		{
			name: "SubscriptionChangeOfState",
			s:    SubscriptionChangeOfState,
			want: "SubscriptionChangeOfState",
		},
		{
			name: "SubscriptionEvent",
			s:    SubscriptionEvent,
			want: "SubscriptionEvent",
		},
		{
			name: "Unknown",
			s:    255,
			want: "Unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.s.String(), "String()")
		})
	}
}
