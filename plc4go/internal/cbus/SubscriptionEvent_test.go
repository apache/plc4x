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
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
)

func TestNewSubscriptionEvent(t *testing.T) {
	type args struct {
		tags          map[string]apiModel.PlcTag
		types         map[string]spiModel.SubscriptionType
		intervals     map[string]time.Duration
		responseCodes map[string]apiModel.PlcResponseCode
		address       map[string]string
		sources       map[string]string
		values        map[string]apiValues.PlcValue
	}
	tests := []struct {
		name string
		args args
		want SubscriptionEvent
	}{
		{
			name: "empty",
			want: func() SubscriptionEvent {
				subscriptionEvent := SubscriptionEvent{}
				event := spiModel.NewDefaultPlcSubscriptionEvent(&subscriptionEvent, nil, nil, nil, nil, nil)
				subscriptionEvent.DefaultPlcSubscriptionEvent = event.(*spiModel.DefaultPlcSubscriptionEvent)
				return subscriptionEvent
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewSubscriptionEvent(tt.args.tags, tt.args.types, tt.args.intervals, tt.args.responseCodes, tt.args.address, tt.args.sources, tt.args.values), "NewSubscriptionEvent(%v, %v, %v, %v, %v, %v, %v)", tt.args.tags, tt.args.types, tt.args.intervals, tt.args.responseCodes, tt.args.address, tt.args.sources, tt.args.values)
		})
	}
}

func TestSubscriptionEvent_GetAddress(t *testing.T) {
	type fields struct {
		DefaultPlcSubscriptionEvent *spiModel.DefaultPlcSubscriptionEvent
		address                     map[string]string
		sources                     map[string]string
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
			name: "just get",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := SubscriptionEvent{
				DefaultPlcSubscriptionEvent: tt.fields.DefaultPlcSubscriptionEvent,
				address:                     tt.fields.address,
				sources:                     tt.fields.sources,
			}
			assert.Equalf(t, tt.want, m.GetAddress(tt.args.name), "GetAddress(%v)", tt.args.name)
		})
	}
}

func TestSubscriptionEvent_GetSource(t *testing.T) {
	type fields struct {
		DefaultPlcSubscriptionEvent *spiModel.DefaultPlcSubscriptionEvent
		address                     map[string]string
		sources                     map[string]string
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
			name: "just get",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := SubscriptionEvent{
				DefaultPlcSubscriptionEvent: tt.fields.DefaultPlcSubscriptionEvent,
				address:                     tt.fields.address,
				sources:                     tt.fields.sources,
			}
			assert.Equalf(t, tt.want, m.GetSource(tt.args.name), "GetSource(%v)", tt.args.name)
		})
	}
}
