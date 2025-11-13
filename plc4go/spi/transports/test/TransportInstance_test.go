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
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/apache/plc4x/plc4go/spi/transports"
)

func TestNewTransportInstance(t *testing.T) {
	type args struct {
		transport *Transport
	}
	tests := []struct {
		name       string
		args       args
		wantAssert func(*testing.T, *TransportInstance) bool
	}{
		{
			name: "create it",
			wantAssert: func(t *testing.T, instance *TransportInstance) bool {
				return assert.NotNil(t, instance)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTransportInstance(tt.args.transport)
			if !assert.True(t, tt.wantAssert(t, got)) {
				t.Errorf("NewTransportInstance() = %v", got)
			}
		})
	}
}

func TestTransportInstance_Close(t *testing.T) {
	type fields struct {
		readBuffer       []byte
		writeBuffer      []byte
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
				transport:        tt.fields.transport,
				writeInterceptor: tt.fields.writeInterceptor,
			}
			if err := m.Connect(tt.args.in0); (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransportInstance_DrainWriteBuffer(t *testing.T) {
	type fields struct {
		readBuffer       []byte
		writeBuffer      []byte
		transport        *Transport
		writeInterceptor func(transportInstance *TransportInstance, data []byte)
	}
	type args struct {
		numBytes uint32
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		manipulator func(t *testing.T, instance *TransportInstance)
		want        []byte
	}{
		{
			name: "drain it",
			manipulator: func(t *testing.T, instance *TransportInstance) {
				instance.connected.Store(true)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				readBuffer:       tt.fields.readBuffer,
				writeBuffer:      tt.fields.writeBuffer,
				transport:        tt.fields.transport,
				writeInterceptor: tt.fields.writeInterceptor,
			}
			if tt.manipulator != nil {
				tt.manipulator(t, m)
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
		transport        *Transport
		writeInterceptor func(transportInstance *TransportInstance, data []byte)
	}
	type args struct {
		ctx   context.Context
		until func(pos uint, currentByte byte, reader transports.ExtendedReader) bool
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		manipulator func(t *testing.T, instance *TransportInstance)
		wantErr     bool
	}{
		{
			name: "fill it (errors)",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
					t.Cleanup(cancel)
					return ctx
				}(),
				until: func(pos uint, currentByte byte, reader transports.ExtendedReader) bool {
					return pos < 3
				},
			},
			manipulator: func(t *testing.T, instance *TransportInstance) {
				instance.connected.Store(true)
			},
			wantErr: true,
		},
		{
			name: "fill it",
			fields: fields{
				readBuffer: []byte{1, 2, 3, 4},
			},
			args: args{
				ctx: t.Context(),
				until: func(pos uint, currentByte byte, reader transports.ExtendedReader) (keepGoing bool) {
					keepGoing = pos < 3
					t.Logf("pos: %d, currentByte: %d: keepGoing: %t", pos, currentByte, keepGoing)
					return keepGoing
				},
			},
			manipulator: func(t *testing.T, instance *TransportInstance) {
				instance.connected.Store(true)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				readChannel:      make(chan []byte, 1),
				simulatedLatency: 10 * time.Millisecond,
				readBuffer:       tt.fields.readBuffer,
				writeBuffer:      tt.fields.writeBuffer,
				transport:        tt.fields.transport,
				writeInterceptor: tt.fields.writeInterceptor,
				log:              zerolog.New(zerolog.NewConsoleWriter(zerolog.ConsoleTestWriter(t))).With().Timestamp().Logger(),
			}
			if tt.manipulator != nil {
				tt.manipulator(t, m)
			}
			if err := m.FillBuffer(tt.args.ctx, tt.args.until); (err != nil) != tt.wantErr {
				t.Errorf("FillBuffer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransportInstance_FillReadBuffer(t *testing.T) {
	type fields struct {
		readBuffer       []byte
		writeBuffer      []byte
		transport        *Transport
		writeInterceptor func(transportInstance *TransportInstance, data []byte)
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		manipulator func(t *testing.T, instance *TransportInstance)
	}{
		{
			name: "fill it",
			fields: fields{
				readBuffer: []byte{1, 2, 3, 4},
			},
			args: args{
				data: []byte{1, 2, 3, 4},
			},
			manipulator: func(t *testing.T, instance *TransportInstance) {
				instance.connected.Store(true)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				readBuffer:       tt.fields.readBuffer,
				writeBuffer:      tt.fields.writeBuffer,
				transport:        tt.fields.transport,
				writeInterceptor: tt.fields.writeInterceptor,
				readChannel:      make(chan []byte, 1),
			}
			if tt.manipulator != nil {
				tt.manipulator(t, m)
			}
			m.FillReadBuffer(tt.args.data)
		})
	}
}

func TestTransportInstance_GetNumBytesAvailableInBuffer(t *testing.T) {
	type fields struct {
		readBuffer       []byte
		writeBuffer      []byte
		transport        *Transport
		writeInterceptor func(transportInstance *TransportInstance, data []byte)
	}
	tests := []struct {
		name        string
		fields      fields
		manipulator func(t *testing.T, instance *TransportInstance)
		want        uint32
		wantErr     bool
	}{
		{
			name: "get it",
			manipulator: func(t *testing.T, instance *TransportInstance) {
				instance.connected.Store(true)
			},
		},
		{
			name: "get it too",
			fields: fields{
				readBuffer: []byte{1, 2, 3, 4},
			},
			manipulator: func(t *testing.T, instance *TransportInstance) {
				instance.connected.Store(true)
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				readBuffer:       tt.fields.readBuffer,
				writeBuffer:      tt.fields.writeBuffer,
				transport:        tt.fields.transport,
				writeInterceptor: tt.fields.writeInterceptor,
			}
			if tt.manipulator != nil {
				tt.manipulator(t, m)
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
		transport        *Transport
		writeInterceptor func(transportInstance *TransportInstance, data []byte)
	}
	tests := []struct {
		name        string
		fields      fields
		manipulator func(t *testing.T, instance *TransportInstance)
		want        uint32
	}{
		{
			name: "get it",
			manipulator: func(t *testing.T, instance *TransportInstance) {
				instance.connected.Store(true)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				readBuffer:       tt.fields.readBuffer,
				writeBuffer:      tt.fields.writeBuffer,
				transport:        tt.fields.transport,
				writeInterceptor: tt.fields.writeInterceptor,
			}
			if tt.manipulator != nil {
				tt.manipulator(t, m)
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
		transport        *Transport
		writeInterceptor func(transportInstance *TransportInstance, data []byte)
	}
	type args struct {
		numBytes uint32
		timeout  time.Duration
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		manipulator func(t *testing.T, instance *TransportInstance)
		want        []byte
		wantErr     bool
	}{
		{
			name: "peek it",
			args: args{
				timeout: 10 * time.Second,
			},
			manipulator: func(t *testing.T, instance *TransportInstance) {
				instance.connected.Store(true)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				readBuffer:       tt.fields.readBuffer,
				writeBuffer:      tt.fields.writeBuffer,
				transport:        tt.fields.transport,
				writeInterceptor: tt.fields.writeInterceptor,
			}
			if tt.manipulator != nil {
				tt.manipulator(t, m)
			}
			got, err := m.PeekReadableBytes(t.Context(), tt.args.numBytes)
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
		transport        *Transport
		writeInterceptor func(transportInstance *TransportInstance, data []byte)
	}
	type args struct {
		numBytes uint32
		timeout  time.Duration
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		manipulator func(t *testing.T, instance *TransportInstance)
		want        []byte
		wantErr     bool
	}{
		{
			name: "read it",
			args: args{
				timeout: 10 * time.Second,
			},
			manipulator: func(t *testing.T, instance *TransportInstance) {
				instance.connected.Store(true)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				readBuffer:       tt.fields.readBuffer,
				writeBuffer:      tt.fields.writeBuffer,
				transport:        tt.fields.transport,
				writeInterceptor: tt.fields.writeInterceptor,
			}
			if tt.manipulator != nil {
				tt.manipulator(t, m)
			}
			got, err := m.Read(t.Context(), tt.args.numBytes)
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
		transport        *Transport
		writeInterceptor func(transportInstance *TransportInstance, data []byte)
	}
	type args struct {
		data    []byte
		timeout time.Duration
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		manipulator func(t *testing.T, instance *TransportInstance)
		wantErr     bool
	}{
		{
			name: "write it",
			args: args{
				timeout: 10 * time.Second,
			},
			manipulator: func(t *testing.T, instance *TransportInstance) {
				instance.connected.Store(true)
			},
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
				data:    []byte{1, 2, 3, 4},
				timeout: 10 * time.Second,
			},
			manipulator: func(t *testing.T, instance *TransportInstance) {
				instance.connected.Store(true)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				readBuffer:       tt.fields.readBuffer,
				writeBuffer:      tt.fields.writeBuffer,
				transport:        tt.fields.transport,
				writeInterceptor: tt.fields.writeInterceptor,
			}
			if tt.manipulator != nil {
				tt.manipulator(t, m)
			}
			if err := m.Write(t.Context(), tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
