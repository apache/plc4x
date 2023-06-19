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
	"context"
	"testing"

	"github.com/apache/plc4x/plc4go/spi/transports"

	"github.com/stretchr/testify/assert"
)

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
			fields: fields{
				connected: true,
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
		until func(pos uint, currentByte byte, reader transports.ExtendedReader) bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "fill it (errors)",
			fields: fields{
				connected: true,
			},
			args: args{
				until: func(pos uint, currentByte byte, reader transports.ExtendedReader) bool {
					return pos < 3
				},
			},
			wantErr: true,
		},
		{
			name: "fill it",
			fields: fields{
				connected:  true,
				readBuffer: []byte{1, 2, 3, 4},
			},
			args: args{
				until: func(pos uint, currentByte byte, reader transports.ExtendedReader) bool {
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
				connected:  true,
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
			fields: fields{
				connected: true,
			},
		},
		{
			name: "get it too",
			fields: fields{
				connected:  true,
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
			fields: fields{
				connected: true,
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
			fields: fields{
				connected: true,
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
			fields: fields{
				connected: true,
			},
			wantErr: true,
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
			fields: fields{
				connected: true,
			},
		},
		{
			name: "write it",
			fields: fields{
				connected: true,
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
