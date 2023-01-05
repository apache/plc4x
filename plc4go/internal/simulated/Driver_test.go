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

package simulated

import (
	"net/url"
	"testing"
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

func TestDriver_CheckQuery(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid query",
			args: args{
				query: "STATE/test:UINT[2]",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDriver()
			if err := d.CheckTagAddress(tt.args.query); (err != nil) != tt.wantErr {
				t.Errorf("CheckTagAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDriver_Discover(t *testing.T) {
	type fields struct {
		valueHandler ValueHandler
	}
	type args struct {
		callback         func(event model.PlcDiscoveryItem)
		discoveryOptions []options.WithDiscoveryOption
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "discovery fails",
			fields: fields{
				valueHandler: NewValueHandler(),
			},
			args: args{
				// Can all be nil, as the call is expected to fail
				callback:         nil,
				discoveryOptions: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Driver{
				valueHandler: tt.fields.valueHandler,
			}
			if err := d.Discover(tt.args.callback, tt.args.discoveryOptions...); (err != nil) != tt.wantErr {
				t.Errorf("Discover() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDriver_GetConnection(t *testing.T) {
	type args struct {
		in0     url.URL
		in1     map[string]transports.Transport
		options map[string][]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "simple no options",
			// Input doesn't really matter, as the code simply ignores most of it.
			args: args{
				in0:     url.URL{},
				in1:     nil,
				options: nil,
			},
			wantErr: false,
		},
		{
			name: "simple with options",
			// Input doesn't really matter, as the code simply ignores most of it.
			args: args{
				in0: url.URL{},
				in1: nil,
				options: map[string][]string{
					"testOption": {"testValue"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDriver()
			connectionChan := d.GetConnection(tt.args.in0, tt.args.in1, tt.args.options)
			timeout := time.NewTimer(3 * time.Second)
			defer utils.CleanupTimer(timeout)
			select {
			case connectResult := <-connectionChan:
				if tt.wantErr && (connectResult.GetErr() == nil) {
					t.Errorf("PlcConnectionPool.GetConnection() = %v, wantErr %v", connectResult.GetErr(), tt.wantErr)
				} else if connectResult.GetErr() != nil {
					t.Errorf("PlcConnectionPool.GetConnection() error = %v, wantErr %v", connectResult.GetErr(), tt.wantErr)
				}
			case <-timeout.C:
				t.Errorf("PlcConnectionPool.GetConnection() got timeout")
			}
		})
	}
}

func TestDriver_GetDefaultTransport(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "simple",
			want: "none",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDriver()
			if got := d.GetDefaultTransport(); got != tt.want {
				t.Errorf("GetDefaultTransport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDriver_GetProtocolCode(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "simple",
			want: "simulated",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDriver()
			if got := d.GetProtocolCode(); got != tt.want {
				t.Errorf("GetProtocolCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDriver_GetProtocolName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "simple",
			want: "Simulated PLC4X Datasource",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDriver()
			if got := d.GetProtocolName(); got != tt.want {
				t.Errorf("GetProtocolName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDriver(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "simple",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDriver()
			if got == nil && !tt.wantErr {
				t.Errorf("NewDriver() error creating")
			}
		})
	}
}
