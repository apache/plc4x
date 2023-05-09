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

package transports

import (
	"bufio"
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewDefaultBufferedTransportInstance(t *testing.T) {
	type args struct {
		defaultBufferedTransportInstanceRequirements DefaultBufferedTransportInstanceRequirements
	}
	tests := []struct {
		name string
		args args
		want DefaultBufferedTransportInstance
	}{
		{
			name: "create it",
			want: &defaultBufferedTransportInstance{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDefaultBufferedTransportInstance(tt.args.defaultBufferedTransportInstanceRequirements); !assert.Equal(t, tt.want, got) {
				t.Errorf("NewDefaultBufferedTransportInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultBufferedTransportInstance_ConnectWithContext(t *testing.T) {
	type fields struct {
		DefaultBufferedTransportInstanceRequirements DefaultBufferedTransportInstanceRequirements
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func(t *testing.T, fields *fields, args *args)
		wantErr   bool
	}{
		{
			name: "connect",
			args: args{
				func() context.Context {
					ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
					t.Cleanup(cancel)
					return ctx
				}(),
			},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				requirements := NewMockDefaultBufferedTransportInstanceRequirements(t)
				requirements.EXPECT().Connect().Return(nil)
				fields.DefaultBufferedTransportInstanceRequirements = requirements
			},
		},
		{
			name: "connect canceled",
			args: args{
				func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()
					return ctx
				}(),
			},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				requirements := NewMockDefaultBufferedTransportInstanceRequirements(t)
				fields.DefaultBufferedTransportInstanceRequirements = requirements
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields, &tt.args)
			}
			m := &defaultBufferedTransportInstance{
				DefaultBufferedTransportInstanceRequirements: tt.fields.DefaultBufferedTransportInstanceRequirements,
			}
			if err := m.ConnectWithContext(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ConnectWithContext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_defaultBufferedTransportInstance_FillBuffer(t *testing.T) {
	type fields struct {
		DefaultBufferedTransportInstanceRequirements DefaultBufferedTransportInstanceRequirements
	}
	type args struct {
		until func(pos uint, currentByte byte, reader *bufio.Reader) bool
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func(t *testing.T, fields *fields, args *args)
		wantErr   bool
	}{
		{
			name: "fill it",
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				requirements := NewMockDefaultBufferedTransportInstanceRequirements(t)
				expect := requirements.EXPECT()
				expect.GetReader().Return(nil)
				fields.DefaultBufferedTransportInstanceRequirements = requirements
			},
		},
		{
			name: "fill it with reader",
			args: args{func(pos uint, currentByte byte, reader *bufio.Reader) bool {
				return pos < 1
			}},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				requirements := NewMockDefaultBufferedTransportInstanceRequirements(t)
				requirements.EXPECT().GetReader().Return(bufio.NewReader(bytes.NewReader([]byte{0x0, 0x0})))
				fields.DefaultBufferedTransportInstanceRequirements = requirements
			},
		},
		{
			name: "fill it with reader errors",
			args: args{func(pos uint, currentByte byte, reader *bufio.Reader) bool {
				return pos < 2
			}},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				requirements := NewMockDefaultBufferedTransportInstanceRequirements(t)
				requirements.EXPECT().GetReader().Return(bufio.NewReader(bytes.NewReader([]byte{0x0, 0x0})))
				fields.DefaultBufferedTransportInstanceRequirements = requirements
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields, &tt.args)
			}
			m := &defaultBufferedTransportInstance{
				DefaultBufferedTransportInstanceRequirements: tt.fields.DefaultBufferedTransportInstanceRequirements,
			}
			if err := m.FillBuffer(tt.args.until); (err != nil) != tt.wantErr {
				t.Errorf("FillBuffer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_defaultBufferedTransportInstance_GetNumBytesAvailableInBuffer(t *testing.T) {
	type fields struct {
		DefaultBufferedTransportInstanceRequirements DefaultBufferedTransportInstanceRequirements
	}
	tests := []struct {
		name      string
		fields    fields
		want      uint32
		mockSetup func(t *testing.T, fields *fields)
		wantErr   bool
	}{
		{
			name: "get it without reader",
			mockSetup: func(t *testing.T, fields *fields) {
				requirements := NewMockDefaultBufferedTransportInstanceRequirements(t)
				requirements.EXPECT().GetReader().Return(nil)
				fields.DefaultBufferedTransportInstanceRequirements = requirements
			},
		},
		{
			name: "get it with reader",
			mockSetup: func(t *testing.T, fields *fields) {
				requirements := NewMockDefaultBufferedTransportInstanceRequirements(t)
				requirements.EXPECT().GetReader().Return(bufio.NewReader(bytes.NewReader([]byte{0x0, 0x0})))
				fields.DefaultBufferedTransportInstanceRequirements = requirements
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields)
			}
			m := &defaultBufferedTransportInstance{
				DefaultBufferedTransportInstanceRequirements: tt.fields.DefaultBufferedTransportInstanceRequirements,
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

func Test_defaultBufferedTransportInstance_PeekReadableBytes(t *testing.T) {
	type fields struct {
		DefaultBufferedTransportInstanceRequirements DefaultBufferedTransportInstanceRequirements
	}
	type args struct {
		numBytes uint32
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func(t *testing.T, fields *fields, args *args)
		want      []byte
		wantErr   bool
	}{
		{
			name: "peek it without reader",
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				requirements := NewMockDefaultBufferedTransportInstanceRequirements(t)
				requirements.EXPECT().GetReader().Return(nil)
				fields.DefaultBufferedTransportInstanceRequirements = requirements
			},
			wantErr: true,
		},
		{
			name: "peek it with reader",
			args: args{numBytes: 2},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				requirements := NewMockDefaultBufferedTransportInstanceRequirements(t)
				requirements.EXPECT().GetReader().Return(bufio.NewReader(bytes.NewReader([]byte{0x0, 0x0})))
				fields.DefaultBufferedTransportInstanceRequirements = requirements
			},
			want: []byte{0x0, 0x0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields, &tt.args)
			}
			m := &defaultBufferedTransportInstance{
				DefaultBufferedTransportInstanceRequirements: tt.fields.DefaultBufferedTransportInstanceRequirements,
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

func Test_defaultBufferedTransportInstance_Read(t *testing.T) {
	type fields struct {
		DefaultBufferedTransportInstanceRequirements DefaultBufferedTransportInstanceRequirements
	}
	type args struct {
		numBytes uint32
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func(t *testing.T, fields *fields, args *args)
		want      []byte
		wantErr   bool
	}{
		{
			name: "read it without reader",
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				requirements := NewMockDefaultBufferedTransportInstanceRequirements(t)
				requirements.EXPECT().GetReader().Return(nil)
				fields.DefaultBufferedTransportInstanceRequirements = requirements
			},
			wantErr: true,
		},
		{
			name: "read it with reader",
			args: args{numBytes: 2},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				requirements := NewMockDefaultBufferedTransportInstanceRequirements(t)
				requirements.EXPECT().GetReader().Return(bufio.NewReader(bytes.NewReader([]byte{0x0, 0x0})))
				fields.DefaultBufferedTransportInstanceRequirements = requirements
			},
			want: []byte{0x0, 0x0},
		},
		{
			name: "read it with reader errors",
			args: args{numBytes: 2},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				requirements := NewMockDefaultBufferedTransportInstanceRequirements(t)
				requirements.EXPECT().GetReader().Return(bufio.NewReader(bytes.NewReader([]byte{0x0})))
				fields.DefaultBufferedTransportInstanceRequirements = requirements
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields, &tt.args)
			}
			m := &defaultBufferedTransportInstance{
				DefaultBufferedTransportInstanceRequirements: tt.fields.DefaultBufferedTransportInstanceRequirements,
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
