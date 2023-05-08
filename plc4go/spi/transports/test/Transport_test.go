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
	"bufio"
	"context"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
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

func TestNewTransportInstance(t *testing.T) {
	type args struct {
		transport *Transport
	}
	tests := []struct {
		name string
		args args
		want *TransportInstance
	}{
		{
			name: "create it",
			want: &TransportInstance{
				readBuffer:  []byte{},
				writeBuffer: []byte{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTransportInstance(tt.args.transport); !assert.Equal(t, tt.want, got) {
				t.Errorf("NewTransportInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransportInstance_Close(t *testing.T) {
	type fields struct {
		readBuffer       []byte
		writeBuffer      []byte
		connected        bool
		transport        *Transport
		writeInterceptor func(transportInstance *TransportInstance, data []byte)
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "close it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				readBuffer:       tt.fields.readBuffer,
				writeBuffer:      tt.fields.writeBuffer,
				connected:        tt.fields.connected,
				transport:        tt.fields.transport,
				writeInterceptor: tt.fields.writeInterceptor,
			}
			if err := m.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransportInstance_Connect(t *testing.T) {
	type fields struct {
		readBuffer       []byte
		writeBuffer      []byte
		connected        bool
		transport        *Transport
		writeInterceptor func(transportInstance *TransportInstance, data []byte)
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "connect it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				readBuffer:       tt.fields.readBuffer,
				writeBuffer:      tt.fields.writeBuffer,
				connected:        tt.fields.connected,
				transport:        tt.fields.transport,
				writeInterceptor: tt.fields.writeInterceptor,
			}
			if err := m.Connect(); (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransportInstance_ConnectWithContext(t *testing.T) {
	type fields struct {
		readBuffer       []byte
		writeBuffer      []byte
		connected        bool
		transport        *Transport
		writeInterceptor func(transportInstance *TransportInstance, data []byte)
	}
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "connect it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				readBuffer:       tt.fields.readBuffer,
				writeBuffer:      tt.fields.writeBuffer,
				connected:        tt.fields.connected,
				transport:        tt.fields.transport,
				writeInterceptor: tt.fields.writeInterceptor,
			}
			if err := m.ConnectWithContext(tt.args.in0); (err != nil) != tt.wantErr {
				t.Errorf("ConnectWithContext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransportInstance_DrainWriteBuffer(t *testing.T) {
	type fields struct {
		readBuffer       []byte
		writeBuffer      []byte
		connected        bool
		transport        *Transport
		writeInterceptor func(transportInstance *TransportInstance, data []byte)
	}
	type args struct {
		numBytes uint32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name: "drain it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				readBuffer:       tt.fields.readBuffer,
				writeBuffer:      tt.fields.writeBuffer,
				connected:        tt.fields.connected,
				transport:        tt.fields.transport,
				writeInterceptor: tt.fields.writeInterceptor,
			}
			if got := m.DrainWriteBuffer(tt.args.numBytes); !assert.Equal(t, tt.want, got) {
				t.Errorf("DrainWriteBuffer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransportInstance_FillBuffer(t *testing.T) {
	type fields struct {
		readBuffer       []byte
		writeBuffer      []byte
		connected        bool
		transport        *Transport
		writeInterceptor func(transportInstance *TransportInstance, data []byte)
	}
	type args struct {
		until func(pos uint, currentByte byte, reader *bufio.Reader) bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "fill it (errors)",
			args: args{
				until: func(pos uint, currentByte byte, reader *bufio.Reader) bool {
					return pos < 3
				},
			},
			wantErr: true,
		},
		{
			name: "fill it",
			fields: fields{
				readBuffer: []byte{1, 2, 3, 4},
			},
			args: args{
				until: func(pos uint, currentByte byte, reader *bufio.Reader) bool {
					return pos < 3
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				readBuffer:       tt.fields.readBuffer,
				writeBuffer:      tt.fields.writeBuffer,
				connected:        tt.fields.connected,
				transport:        tt.fields.transport,
				writeInterceptor: tt.fields.writeInterceptor,
			}
			if err := m.FillBuffer(tt.args.until); (err != nil) != tt.wantErr {
				t.Errorf("FillBuffer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransportInstance_FillReadBuffer(t *testing.T) {
	type fields struct {
		readBuffer       []byte
		writeBuffer      []byte
		connected        bool
		transport        *Transport
		writeInterceptor func(transportInstance *TransportInstance, data []byte)
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "fill it",
			fields: fields{
				readBuffer: []byte{1, 2, 3, 4},
			},
			args: args{
				data: []byte{1, 2, 3, 4},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				readBuffer:       tt.fields.readBuffer,
				writeBuffer:      tt.fields.writeBuffer,
				connected:        tt.fields.connected,
				transport:        tt.fields.transport,
				writeInterceptor: tt.fields.writeInterceptor,
			}
			m.FillReadBuffer(tt.args.data)
		})
	}
}

func TestTransportInstance_GetNumBytesAvailableInBuffer(t *testing.T) {
	type fields struct {
		readBuffer       []byte
		writeBuffer      []byte
		connected        bool
		transport        *Transport
		writeInterceptor func(transportInstance *TransportInstance, data []byte)
	}
	tests := []struct {
		name    string
		fields  fields
		want    uint32
		wantErr bool
	}{
		{
			name: "get it",
		},
		{
			name: "get it",
			fields: fields{
				readBuffer: []byte{1, 2, 3, 4},
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				readBuffer:       tt.fields.readBuffer,
				writeBuffer:      tt.fields.writeBuffer,
				connected:        tt.fields.connected,
				transport:        tt.fields.transport,
				writeInterceptor: tt.fields.writeInterceptor,
			}
			got, err := m.GetNumBytesAvailableInBuffer()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNumBytesAvailableInBuffer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetNumBytesAvailableInBuffer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransportInstance_GetNumDrainableBytes(t *testing.T) {
	type fields struct {
		readBuffer       []byte
		writeBuffer      []byte
		connected        bool
		transport        *Transport
		writeInterceptor func(transportInstance *TransportInstance, data []byte)
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				readBuffer:       tt.fields.readBuffer,
				writeBuffer:      tt.fields.writeBuffer,
				connected:        tt.fields.connected,
				transport:        tt.fields.transport,
				writeInterceptor: tt.fields.writeInterceptor,
			}
			if got := m.GetNumDrainableBytes(); got != tt.want {
				t.Errorf("GetNumDrainableBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransportInstance_IsConnected(t *testing.T) {
	type fields struct {
		readBuffer       []byte
		writeBuffer      []byte
		connected        bool
		transport        *Transport
		writeInterceptor func(transportInstance *TransportInstance, data []byte)
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "check it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				readBuffer:       tt.fields.readBuffer,
				writeBuffer:      tt.fields.writeBuffer,
				connected:        tt.fields.connected,
				transport:        tt.fields.transport,
				writeInterceptor: tt.fields.writeInterceptor,
			}
			if got := m.IsConnected(); got != tt.want {
				t.Errorf("IsConnected() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransportInstance_PeekReadableBytes(t *testing.T) {
	type fields struct {
		readBuffer       []byte
		writeBuffer      []byte
		connected        bool
		transport        *Transport
		writeInterceptor func(transportInstance *TransportInstance, data []byte)
	}
	type args struct {
		numBytes uint32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "peek it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				readBuffer:       tt.fields.readBuffer,
				writeBuffer:      tt.fields.writeBuffer,
				connected:        tt.fields.connected,
				transport:        tt.fields.transport,
				writeInterceptor: tt.fields.writeInterceptor,
			}
			got, err := m.PeekReadableBytes(tt.args.numBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("PeekReadableBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("PeekReadableBytes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransportInstance_Read(t *testing.T) {
	type fields struct {
		readBuffer       []byte
		writeBuffer      []byte
		connected        bool
		transport        *Transport
		writeInterceptor func(transportInstance *TransportInstance, data []byte)
	}
	type args struct {
		numBytes uint32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "read it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				readBuffer:       tt.fields.readBuffer,
				writeBuffer:      tt.fields.writeBuffer,
				connected:        tt.fields.connected,
				transport:        tt.fields.transport,
				writeInterceptor: tt.fields.writeInterceptor,
			}
			got, err := m.Read(tt.args.numBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransportInstance_SetWriteInterceptor(t *testing.T) {
	type fields struct {
		readBuffer       []byte
		writeBuffer      []byte
		connected        bool
		transport        *Transport
		writeInterceptor func(transportInstance *TransportInstance, data []byte)
	}
	type args struct {
		writeInterceptor func(transportInstance *TransportInstance, data []byte)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "set it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				readBuffer:       tt.fields.readBuffer,
				writeBuffer:      tt.fields.writeBuffer,
				connected:        tt.fields.connected,
				transport:        tt.fields.transport,
				writeInterceptor: tt.fields.writeInterceptor,
			}
			m.SetWriteInterceptor(tt.args.writeInterceptor)
		})
	}
}

func TestTransportInstance_String(t *testing.T) {
	type fields struct {
		readBuffer       []byte
		writeBuffer      []byte
		connected        bool
		transport        *Transport
		writeInterceptor func(transportInstance *TransportInstance, data []byte)
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "string it",
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				readBuffer:       tt.fields.readBuffer,
				writeBuffer:      tt.fields.writeBuffer,
				connected:        tt.fields.connected,
				transport:        tt.fields.transport,
				writeInterceptor: tt.fields.writeInterceptor,
			}
			if got := m.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransportInstance_Write(t *testing.T) {
	type fields struct {
		readBuffer       []byte
		writeBuffer      []byte
		connected        bool
		transport        *Transport
		writeInterceptor func(transportInstance *TransportInstance, data []byte)
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "write it",
		},
		{
			name: "write it",
			fields: fields{
				writeInterceptor: func(transportInstance *TransportInstance, data []byte) {
					assert.NotNil(t, transportInstance)
					assert.Equal(t, []byte{1, 2, 3, 4}, data)
				},
			},
			args: args{
				data: []byte{1, 2, 3, 4},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				readBuffer:       tt.fields.readBuffer,
				writeBuffer:      tt.fields.writeBuffer,
				connected:        tt.fields.connected,
				transport:        tt.fields.transport,
				writeInterceptor: tt.fields.writeInterceptor,
			}
			if err := m.Write(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
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
