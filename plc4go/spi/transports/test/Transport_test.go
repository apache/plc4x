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

package test

import (
	"github.com/rs/zerolog/log"
	"net/url"
	"testing"

	"github.com/apache/plc4x/plc4go/spi/transports"

	"github.com/stretchr/testify/assert"
)

func TestNewTransport(t *testing.T) {
	tests := []struct {
		name string
		want *Transport
	}{
		{
			name: "create it",
			want: &Transport{
				preregisteredInstances: map[url.URL]transports.TransportInstance{},
				log:                    log.Logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTransport(); !assert.Equal(t, tt.want, got) {
				t.Errorf("NewTransport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransport_AddPreregisteredInstances(t *testing.T) {
	type fields struct {
		preregisteredInstances map[url.URL]transports.TransportInstance
	}
	type args struct {
		transportUrl          url.URL
		preregisteredInstance transports.TransportInstance
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "add it",
			fields: fields{
				preregisteredInstances: map[url.URL]transports.TransportInstance{},
			},
		},
		{
			name: "add it (existing)",
			fields: fields{
				preregisteredInstances: map[url.URL]transports.TransportInstance{
					url.URL{Host: "abcdefg"}: nil,
				},
			},
			args: args{
				transportUrl: url.URL{Host: "abcdefg"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Transport{
				preregisteredInstances: tt.fields.preregisteredInstances,
			}
			if err := m.AddPreregisteredInstances(tt.args.transportUrl, tt.args.preregisteredInstance); (err != nil) != tt.wantErr {
				t.Errorf("AddPreregisteredInstances() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransport_CreateTransportInstance(t *testing.T) {
	type fields struct {
		preregisteredInstances map[url.URL]transports.TransportInstance
	}
	type args struct {
		transportUrl url.URL
		options      map[string][]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    transports.TransportInstance
		wantErr bool
	}{
		{
			name: "create it",
			fields: fields{
				preregisteredInstances: map[url.URL]transports.TransportInstance{},
			},
			want: &TransportInstance{
				readBuffer:  []byte{},
				writeBuffer: []byte{},
				transport:   NewTransport(),
				log:         log.Logger,
			},
		},
		{
			name: "create it (pre registered",
			fields: fields{
				preregisteredInstances: map[url.URL]transports.TransportInstance{
					url.URL{Host: "abcdefg"}: nil,
				},
			},
			args: args{
				transportUrl: url.URL{Host: "abcdefg"},
			},
			want: nil,
		},
		{
			name: "fail it on purpose",
			args: args{
				options: map[string][]string{
					"failTestTransport": {"yes please"},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Transport{
				preregisteredInstances: tt.fields.preregisteredInstances,
				log:                    log.Logger,
			}
			got, err := m.CreateTransportInstance(tt.args.transportUrl, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTransportInstance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("CreateTransportInstance() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransport_GetTransportCode(t *testing.T) {
	type fields struct {
		preregisteredInstances map[url.URL]transports.TransportInstance
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get it",
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Transport{
				preregisteredInstances: tt.fields.preregisteredInstances,
			}
			if got := m.GetTransportCode(); got != tt.want {
				t.Errorf("GetTransportCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransport_GetTransportName(t *testing.T) {
	type fields struct {
		preregisteredInstances map[url.URL]transports.TransportInstance
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get it",
			want: "Test Transport",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Transport{
				preregisteredInstances: tt.fields.preregisteredInstances,
			}
			if got := m.GetTransportName(); got != tt.want {
				t.Errorf("GetTransportName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransport_String(t *testing.T) {
	type fields struct {
		preregisteredInstances map[url.URL]transports.TransportInstance
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "string it",
			want: "test(Test Transport)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Transport{
				preregisteredInstances: tt.fields.preregisteredInstances,
			}
			if got := m.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
