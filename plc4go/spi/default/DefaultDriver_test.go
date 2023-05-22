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

package _default

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/url"
	"testing"

	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

func TestNewDefaultDriver(t *testing.T) {
	type args struct {
		defaultDriverRequirements DefaultDriverRequirements
		protocolCode              string
		protocolName              string
		defaultTransport          string
		plcTagHandler             spi.PlcTagHandler
	}
	tests := []struct {
		name string
		args args
		want DefaultDriver
	}{
		{
			name: "create a new one",
			want: &defaultDriver{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDefaultDriver(tt.args.defaultDriverRequirements, tt.args.protocolCode, tt.args.protocolName, tt.args.defaultTransport, tt.args.plcTagHandler), "NewDefaultDriver(%v, %v, %v, %v, %v)", tt.args.defaultDriverRequirements, tt.args.protocolCode, tt.args.protocolName, tt.args.defaultTransport, tt.args.plcTagHandler)
		})
	}
}

func Test_defaultDriver_CheckQuery(t *testing.T) {
	type fields struct {
		DefaultDriverRequirements DefaultDriverRequirements
		protocolCode              string
		protocolName              string
		defaultTransport          string
		plcTagHandler             spi.PlcTagHandler
	}
	type args struct {
		query string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func(t *testing.T, fields *fields, args *args)
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "check it",
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				handler := NewMockPlcTagHandler(t)
				handler.EXPECT().ParseQuery(mock.Anything).Return(nil, nil)
				fields.plcTagHandler = handler
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields, &tt.args)
			}
			d := &defaultDriver{
				DefaultDriverRequirements: tt.fields.DefaultDriverRequirements,
				protocolCode:              tt.fields.protocolCode,
				protocolName:              tt.fields.protocolName,
				defaultTransport:          tt.fields.defaultTransport,
				plcTagHandler:             tt.fields.plcTagHandler,
			}
			tt.wantErr(t, d.CheckQuery(tt.args.query), fmt.Sprintf("CheckQuery(%v)", tt.args.query))
		})
	}
}

func Test_defaultDriver_CheckTagAddress(t *testing.T) {
	type fields struct {
		DefaultDriverRequirements DefaultDriverRequirements
		protocolCode              string
		protocolName              string
		defaultTransport          string
		plcTagHandler             spi.PlcTagHandler
	}
	type args struct {
		query string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func(t *testing.T, fields *fields, args *args)
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "check it",
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				handler := NewMockPlcTagHandler(t)
				handler.EXPECT().ParseTag(mock.Anything).Return(nil, nil)
				fields.plcTagHandler = handler
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields, &tt.args)
			}
			d := &defaultDriver{
				DefaultDriverRequirements: tt.fields.DefaultDriverRequirements,
				protocolCode:              tt.fields.protocolCode,
				protocolName:              tt.fields.protocolName,
				defaultTransport:          tt.fields.defaultTransport,
				plcTagHandler:             tt.fields.plcTagHandler,
			}
			tt.wantErr(t, d.CheckTagAddress(tt.args.query), fmt.Sprintf("CheckTagAddress(%v)", tt.args.query))
		})
	}
}

func Test_defaultDriver_Discover(t *testing.T) {
	type fields struct {
		DefaultDriverRequirements DefaultDriverRequirements
		protocolCode              string
		protocolName              string
		defaultTransport          string
		plcTagHandler             spi.PlcTagHandler
	}
	type args struct {
		callback         func(event apiModel.PlcDiscoveryItem)
		discoveryOptions []options.WithDiscoveryOption
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func(t *testing.T, fields *fields, args *args)
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "discover it",
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				requirements := NewMockDefaultDriverRequirements(t)
				requirements.EXPECT().DiscoverWithContext(mock.Anything, mock.Anything).Return(nil)
				fields.DefaultDriverRequirements = requirements
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields, &tt.args)
			}
			d := &defaultDriver{
				DefaultDriverRequirements: tt.fields.DefaultDriverRequirements,
				protocolCode:              tt.fields.protocolCode,
				protocolName:              tt.fields.protocolName,
				defaultTransport:          tt.fields.defaultTransport,
				plcTagHandler:             tt.fields.plcTagHandler,
			}
			tt.wantErr(t, d.Discover(tt.args.callback, tt.args.discoveryOptions...), fmt.Sprintf("Discover(func(), %v)", tt.args.discoveryOptions))
		})
	}
}

func Test_defaultDriver_DiscoverWithContext(t *testing.T) {
	type fields struct {
		DefaultDriverRequirements DefaultDriverRequirements
		protocolCode              string
		protocolName              string
		defaultTransport          string
		plcTagHandler             spi.PlcTagHandler
	}
	type args struct {
		in0 context.Context
		in1 func(event apiModel.PlcDiscoveryItem)
		in2 []options.WithDiscoveryOption
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "discover it",
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultDriver{
				DefaultDriverRequirements: tt.fields.DefaultDriverRequirements,
				protocolCode:              tt.fields.protocolCode,
				protocolName:              tt.fields.protocolName,
				defaultTransport:          tt.fields.defaultTransport,
				plcTagHandler:             tt.fields.plcTagHandler,
			}
			tt.wantErr(t, d.DiscoverWithContext(tt.args.in0, tt.args.in1, tt.args.in2...), fmt.Sprintf("DiscoverWithContext(%v, func(), %v)", tt.args.in0, tt.args.in2))
		})
	}
}

func Test_defaultDriver_GetConnection(t *testing.T) {
	type fields struct {
		DefaultDriverRequirements DefaultDriverRequirements
		protocolCode              string
		protocolName              string
		defaultTransport          string
		plcTagHandler             spi.PlcTagHandler
	}
	type args struct {
		transportUrl url.URL
		transports   map[string]transports.Transport
		options      map[string][]string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func(t *testing.T, fields *fields, args *args, want *<-chan plc4go.PlcConnectionConnectResult)
		want      <-chan plc4go.PlcConnectionConnectResult
	}{
		{
			name: "get a connection",
			mockSetup: func(t *testing.T, fields *fields, args *args, want *<-chan plc4go.PlcConnectionConnectResult) {
				requirements := NewMockDefaultDriverRequirements(t)
				results := make(chan plc4go.PlcConnectionConnectResult, 1)
				*want = results
				results <- NewMockPlcConnectionConnectResult(t)
				requirements.EXPECT().GetConnectionWithContext(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(results)
				fields.DefaultDriverRequirements = requirements
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields, &tt.args, &tt.want)
			}
			d := &defaultDriver{
				DefaultDriverRequirements: tt.fields.DefaultDriverRequirements,
				protocolCode:              tt.fields.protocolCode,
				protocolName:              tt.fields.protocolName,
				defaultTransport:          tt.fields.defaultTransport,
				plcTagHandler:             tt.fields.plcTagHandler,
			}
			assert.Equalf(t, tt.want, d.GetConnection(tt.args.transportUrl, tt.args.transports, tt.args.options), "GetConnection(%v, %v, %v)", tt.args.transportUrl, tt.args.transports, tt.args.options)
		})
	}
}

func Test_defaultDriver_GetDefaultTransport(t *testing.T) {
	type fields struct {
		DefaultDriverRequirements DefaultDriverRequirements
		protocolCode              string
		protocolName              string
		defaultTransport          string
		plcTagHandler             spi.PlcTagHandler
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultDriver{
				DefaultDriverRequirements: tt.fields.DefaultDriverRequirements,
				protocolCode:              tt.fields.protocolCode,
				protocolName:              tt.fields.protocolName,
				defaultTransport:          tt.fields.defaultTransport,
				plcTagHandler:             tt.fields.plcTagHandler,
			}
			assert.Equalf(t, tt.want, d.GetDefaultTransport(), "GetDefaultTransport()")
		})
	}
}

func Test_defaultDriver_GetPlcTagHandler(t *testing.T) {
	type fields struct {
		DefaultDriverRequirements DefaultDriverRequirements
		protocolCode              string
		protocolName              string
		defaultTransport          string
		plcTagHandler             spi.PlcTagHandler
	}
	tests := []struct {
		name   string
		fields fields
		want   spi.PlcTagHandler
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultDriver{
				DefaultDriverRequirements: tt.fields.DefaultDriverRequirements,
				protocolCode:              tt.fields.protocolCode,
				protocolName:              tt.fields.protocolName,
				defaultTransport:          tt.fields.defaultTransport,
				plcTagHandler:             tt.fields.plcTagHandler,
			}
			assert.Equalf(t, tt.want, d.GetPlcTagHandler(), "GetPlcTagHandler()")
		})
	}
}

func Test_defaultDriver_GetProtocolCode(t *testing.T) {
	type fields struct {
		DefaultDriverRequirements DefaultDriverRequirements
		protocolCode              string
		protocolName              string
		defaultTransport          string
		plcTagHandler             spi.PlcTagHandler
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultDriver{
				DefaultDriverRequirements: tt.fields.DefaultDriverRequirements,
				protocolCode:              tt.fields.protocolCode,
				protocolName:              tt.fields.protocolName,
				defaultTransport:          tt.fields.defaultTransport,
				plcTagHandler:             tt.fields.plcTagHandler,
			}
			assert.Equalf(t, tt.want, d.GetProtocolCode(), "GetProtocolCode()")
		})
	}
}

func Test_defaultDriver_GetProtocolName(t *testing.T) {
	type fields struct {
		DefaultDriverRequirements DefaultDriverRequirements
		protocolCode              string
		protocolName              string
		defaultTransport          string
		plcTagHandler             spi.PlcTagHandler
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultDriver{
				DefaultDriverRequirements: tt.fields.DefaultDriverRequirements,
				protocolCode:              tt.fields.protocolCode,
				protocolName:              tt.fields.protocolName,
				defaultTransport:          tt.fields.defaultTransport,
				plcTagHandler:             tt.fields.plcTagHandler,
			}
			assert.Equalf(t, tt.want, d.GetProtocolName(), "GetProtocolName()")
		})
	}
}

func Test_defaultDriver_String(t *testing.T) {
	type fields struct {
		DefaultDriverRequirements DefaultDriverRequirements
		protocolCode              string
		protocolName              string
		defaultTransport          string
		plcTagHandler             spi.PlcTagHandler
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get it",
			fields: fields{
				protocolName:     "abc",
				protocolCode:     "def",
				defaultTransport: "ghi",
			},
			want: "abc (def) [ghi]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultDriver{
				DefaultDriverRequirements: tt.fields.DefaultDriverRequirements,
				protocolCode:              tt.fields.protocolCode,
				protocolName:              tt.fields.protocolName,
				defaultTransport:          tt.fields.defaultTransport,
				plcTagHandler:             tt.fields.plcTagHandler,
			}
			assert.Equalf(t, tt.want, d.String(), "String()")
		})
	}
}

func Test_defaultDriver_SupportsDiscovery(t *testing.T) {
	type fields struct {
		DefaultDriverRequirements DefaultDriverRequirements
		protocolCode              string
		protocolName              string
		defaultTransport          string
		plcTagHandler             spi.PlcTagHandler
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultDriver{
				DefaultDriverRequirements: tt.fields.DefaultDriverRequirements,
				protocolCode:              tt.fields.protocolCode,
				protocolName:              tt.fields.protocolName,
				defaultTransport:          tt.fields.defaultTransport,
				plcTagHandler:             tt.fields.plcTagHandler,
			}
			assert.Equalf(t, tt.want, d.SupportsDiscovery(), "SupportsDiscovery()")
		})
	}
}
