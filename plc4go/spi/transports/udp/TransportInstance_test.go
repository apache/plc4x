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

package udp

import (
	"bufio"
	"bytes"
	"context"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/rs/zerolog/log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/nettest"
)

func TestNewTransportInstance(t *testing.T) {
	type args struct {
		localAddress   *net.UDPAddr
		remoteAddress  *net.UDPAddr
		connectTimeout uint32
		soReUse        bool
		transport      *Transport
	}
	tests := []struct {
		name string
		args args
		want *TransportInstance
	}{
		{
			name: "create it",
			want: &TransportInstance{
				log: log.Logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTransportInstance(tt.args.localAddress, tt.args.remoteAddress, tt.args.connectTimeout, tt.args.soReUse, tt.args.transport); !assert.Equal(t, tt.want, got) {
				t.Errorf("NewTransportInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransportInstance_Close(t *testing.T) {
	type fields struct {
		LocalAddress   *net.UDPAddr
		RemoteAddress  *net.UDPAddr
		ConnectTimeout uint32
		SoReUse        bool
		transport      *Transport
		udpConn        *net.UDPConn
		reader         *bufio.Reader
	}
	tests := []struct {
		name        string
		fields      fields
		manipulator func(t *testing.T, instance *TransportInstance)
		wantErr     bool
	}{
		{
			name: "close it failing",
			fields: fields{
				udpConn: &net.UDPConn{},
			},
			manipulator: func(t *testing.T, instance *TransportInstance) {
				instance.connected.Store(true)
			},
			wantErr: true,
		},
		{
			name: "close success",
			fields: fields{
				udpConn: func() *net.UDPConn {
					listener, err := nettest.NewLocalPacketListener("udp")
					require.NoError(t, err)
					t.Cleanup(func() {
						assert.Error(t, listener.Close()) // Note: connection should have been closed
					})
					return listener.(*net.UDPConn)
				}(),
			},
			manipulator: func(t *testing.T, instance *TransportInstance) {
				instance.connected.Store(true)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				LocalAddress:   tt.fields.LocalAddress,
				RemoteAddress:  tt.fields.RemoteAddress,
				ConnectTimeout: tt.fields.ConnectTimeout,
				SoReUse:        tt.fields.SoReUse,
				transport:      tt.fields.transport,
				udpConn:        tt.fields.udpConn,
				reader:         tt.fields.reader,
			}
			if tt.manipulator != nil {
				tt.manipulator(t, m)
			}
			if err := m.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransportInstance_Connect(t *testing.T) {
	type fields struct {
		LocalAddress   *net.UDPAddr
		RemoteAddress  *net.UDPAddr
		ConnectTimeout uint32
		SoReUse        bool
		transport      *Transport
		udpConn        *net.UDPConn
		reader         *bufio.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "connect it (error)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				LocalAddress:   tt.fields.LocalAddress,
				RemoteAddress:  tt.fields.RemoteAddress,
				ConnectTimeout: tt.fields.ConnectTimeout,
				SoReUse:        tt.fields.SoReUse,
				transport:      tt.fields.transport,
				udpConn:        tt.fields.udpConn,
				reader:         tt.fields.reader,
			}
			if err := m.Connect(); (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransportInstance_ConnectWithContext(t *testing.T) {
	type fields struct {
		LocalAddress   *net.UDPAddr
		RemoteAddress  *net.UDPAddr
		ConnectTimeout uint32
		SoReUse        bool
		transport      *Transport
		udpConn        *net.UDPConn
		reader         *bufio.Reader
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "connect it",
			args: args{ctx: context.Background()},
		},
		{
			name: "connect",
			fields: fields{
				RemoteAddress: func() *net.UDPAddr {
					listener, err := nettest.NewLocalPacketListener("udp")
					require.NoError(t, err)
					t.Logf("remote listener %#q", listener.LocalAddr())
					t.Cleanup(func() {
						assert.NoError(t, listener.Close())
					})
					return listener.LocalAddr().(*net.UDPAddr)
				}(),
			},
			args: args{ctx: context.Background()},
		},
		{
			name: "connect with wrong address", // TODO: not sure how to tests undialable ips here
			fields: fields{
				RemoteAddress: &net.UDPAddr{IP: net.IPv4(255, 255, 255, 255), Port: 12},
			},
			args: args{ctx: context.Background()},
		},
		{
			name: "connect with localAddress",
			fields: fields{
				LocalAddress: func() *net.UDPAddr {
					listener, err := nettest.NewLocalPacketListener("udp")
					require.NoError(t, err)
					t.Logf("local listener %#q", listener.LocalAddr())
					t.Cleanup(func() {
						assert.NoError(t, listener.Close())
					})
					return listener.LocalAddr().(*net.UDPAddr)
				}(),
				RemoteAddress: func() *net.UDPAddr {
					listener, err := nettest.NewLocalPacketListener("udp")
					require.NoError(t, err)
					t.Logf("remote listener %#q", listener.LocalAddr())
					t.Cleanup(func() {
						assert.NoError(t, listener.Close())
					})
					return listener.LocalAddr().(*net.UDPAddr)
				}(),
			},
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
		{
			name: "connect reuse",
			fields: fields{
				LocalAddress: func() *net.UDPAddr {
					listener, err := nettest.NewLocalPacketListener("udp")
					require.NoError(t, err)
					t.Logf("local listener %#q", listener.LocalAddr())
					assert.NoError(t, listener.Close()) // We close directly again
					return listener.LocalAddr().(*net.UDPAddr)
				}(),
				SoReUse: true,
			},
			args: args{ctx: context.Background()},
		},
		{
			name: "connect reuse (used)",
			fields: fields{
				LocalAddress: func() *net.UDPAddr {
					listener, err := nettest.NewLocalPacketListener("udp")
					require.NoError(t, err)
					t.Logf("local listener %#q", listener.LocalAddr())
					t.Cleanup(func() {
						assert.NoError(t, listener.Close())
					})
					return listener.LocalAddr().(*net.UDPAddr)
				}(),
				SoReUse: true,
			},
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
		{
			name: "connect reuse (used)",
			fields: fields{
				LocalAddress: func() *net.UDPAddr {
					listener, err := nettest.NewLocalPacketListener("udp")
					require.NoError(t, err)
					t.Logf("local listener %#q", listener.LocalAddr())
					t.Cleanup(func() {
						assert.NoError(t, listener.Close())
					})
					return listener.LocalAddr().(*net.UDPAddr)
				}(),
			},
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				LocalAddress:   tt.fields.LocalAddress,
				RemoteAddress:  tt.fields.RemoteAddress,
				ConnectTimeout: tt.fields.ConnectTimeout,
				SoReUse:        tt.fields.SoReUse,
				transport:      tt.fields.transport,
				udpConn:        tt.fields.udpConn,
				reader:         tt.fields.reader,
			}
			if err := m.ConnectWithContext(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ConnectWithContext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransportInstance_FillBuffer(t *testing.T) {
	type fields struct {
		LocalAddress   *net.UDPAddr
		RemoteAddress  *net.UDPAddr
		ConnectTimeout uint32
		SoReUse        bool
		transport      *Transport
		udpConn        *net.UDPConn
		reader         *bufio.Reader
	}
	type args struct {
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
			name:    "do it",
			wantErr: true,
		},
		{
			name: "do it with reader",
			fields: fields{
				reader: bufio.NewReader(bytes.NewReader([]byte{1, 2, 3, 4})),
			},
			args: args{
				until: func(pos uint, currentByte byte, reader transports.ExtendedReader) bool {
					return pos < 2
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
				LocalAddress:   tt.fields.LocalAddress,
				RemoteAddress:  tt.fields.RemoteAddress,
				ConnectTimeout: tt.fields.ConnectTimeout,
				SoReUse:        tt.fields.SoReUse,
				transport:      tt.fields.transport,
				udpConn:        tt.fields.udpConn,
				reader:         tt.fields.reader,
			}
			if tt.manipulator != nil {
				tt.manipulator(t, m)
			}
			if err := m.FillBuffer(tt.args.until); (err != nil) != tt.wantErr {
				t.Errorf("FillBuffer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransportInstance_GetNumBytesAvailableInBuffer(t *testing.T) {
	type fields struct {
		LocalAddress   *net.UDPAddr
		RemoteAddress  *net.UDPAddr
		ConnectTimeout uint32
		SoReUse        bool
		transport      *Transport
		udpConn        *net.UDPConn
		reader         *bufio.Reader
	}
	tests := []struct {
		name        string
		fields      fields
		manipulator func(t *testing.T, instance *TransportInstance)
		want        uint32
		wantErr     bool
	}{
		{
			name: "get em",
			fields: fields{
				reader: bufio.NewReader(bytes.NewReader([]byte{1, 2, 3, 4})),
			},
			manipulator: func(t *testing.T, instance *TransportInstance) {
				instance.connected.Store(true)
			},
			want: 4,
		},
		{
			name: "get em (no reader)",
			manipulator: func(t *testing.T, instance *TransportInstance) {
				instance.connected.Store(true)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				LocalAddress:   tt.fields.LocalAddress,
				RemoteAddress:  tt.fields.RemoteAddress,
				ConnectTimeout: tt.fields.ConnectTimeout,
				SoReUse:        tt.fields.SoReUse,
				transport:      tt.fields.transport,
				udpConn:        tt.fields.udpConn,
				reader:         tt.fields.reader,
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

func TestTransportInstance_IsConnected(t *testing.T) {
	type fields struct {
		LocalAddress   *net.UDPAddr
		RemoteAddress  *net.UDPAddr
		ConnectTimeout uint32
		SoReUse        bool
		transport      *Transport
		udpConn        *net.UDPConn
		reader         *bufio.Reader
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
				LocalAddress:   tt.fields.LocalAddress,
				RemoteAddress:  tt.fields.RemoteAddress,
				ConnectTimeout: tt.fields.ConnectTimeout,
				SoReUse:        tt.fields.SoReUse,
				transport:      tt.fields.transport,
				udpConn:        tt.fields.udpConn,
				reader:         tt.fields.reader,
			}
			if got := m.IsConnected(); got != tt.want {
				t.Errorf("IsConnected() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransportInstance_PeekReadableBytes(t *testing.T) {
	type fields struct {
		LocalAddress   *net.UDPAddr
		RemoteAddress  *net.UDPAddr
		ConnectTimeout uint32
		SoReUse        bool
		transport      *Transport
		udpConn        *net.UDPConn
		reader         *bufio.Reader
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
		wantErr     bool
	}{
		{
			name: "peek it",
			fields: fields{
				reader: bufio.NewReader(bytes.NewReader([]byte{1, 2, 3, 4})),
			},
			manipulator: func(t *testing.T, instance *TransportInstance) {
				instance.connected.Store(true)
			},
			want: []byte{},
		},
		{
			name: "peek it 3",
			fields: fields{
				reader: bufio.NewReader(bytes.NewReader([]byte{1, 2, 3, 4})),
			},
			args: args{
				numBytes: 3,
			},
			manipulator: func(t *testing.T, instance *TransportInstance) {
				instance.connected.Store(true)
			},
			want: []byte{1, 2, 3},
		},
		{
			name:    "peek it (not connected)",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				LocalAddress:   tt.fields.LocalAddress,
				RemoteAddress:  tt.fields.RemoteAddress,
				ConnectTimeout: tt.fields.ConnectTimeout,
				SoReUse:        tt.fields.SoReUse,
				transport:      tt.fields.transport,
				udpConn:        tt.fields.udpConn,
				reader:         tt.fields.reader,
			}
			if tt.manipulator != nil {
				tt.manipulator(t, m)
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
		LocalAddress   *net.UDPAddr
		RemoteAddress  *net.UDPAddr
		ConnectTimeout uint32
		SoReUse        bool
		transport      *Transport
		udpConn        *net.UDPConn
		reader         *bufio.Reader
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
		wantErr     bool
	}{
		{
			name: "read it",
			fields: fields{
				reader: bufio.NewReader(bytes.NewReader([]byte{1, 2, 3, 4})),
			},
			manipulator: func(t *testing.T, instance *TransportInstance) {
				instance.connected.Store(true)
			},
			want: []byte{},
		},
		{
			name: "read it 3",
			fields: fields{
				reader: bufio.NewReader(bytes.NewReader([]byte{1, 2, 3, 4})),
			},
			args: args{
				numBytes: 3,
			},
			manipulator: func(t *testing.T, instance *TransportInstance) {
				instance.connected.Store(true)
			},
			want: []byte{1, 2, 3},
		},
		{
			name: "read it 5",
			fields: fields{
				reader: bufio.NewReader(bytes.NewReader([]byte{1, 2, 3, 4})),
			},
			args: args{
				numBytes: 5,
			},
			manipulator: func(t *testing.T, instance *TransportInstance) {
				instance.connected.Store(true)
			},
			wantErr: true,
		},
		{
			name:    "read it (not connected)",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				LocalAddress:   tt.fields.LocalAddress,
				RemoteAddress:  tt.fields.RemoteAddress,
				ConnectTimeout: tt.fields.ConnectTimeout,
				SoReUse:        tt.fields.SoReUse,
				transport:      tt.fields.transport,
				udpConn:        tt.fields.udpConn,
				reader:         tt.fields.reader,
			}
			if tt.manipulator != nil {
				tt.manipulator(t, m)
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

func TestTransportInstance_String(t *testing.T) {
	type fields struct {
		LocalAddress   *net.UDPAddr
		RemoteAddress  *net.UDPAddr
		ConnectTimeout uint32
		SoReUse        bool
		transport      *Transport
		udpConn        *net.UDPConn
		reader         *bufio.Reader
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "string it",
			want: "udp:<nil>-><nil>",
		},
		{
			name: "string it with content",
			fields: fields{
				LocalAddress:  &net.UDPAddr{IP: net.IPv4(1, 2, 3, 4), Port: 5},
				RemoteAddress: &net.UDPAddr{IP: net.IPv4(6, 7, 8, 9), Port: 10},
			},
			want: "udp:1.2.3.4:5->6.7.8.9:10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				LocalAddress:   tt.fields.LocalAddress,
				RemoteAddress:  tt.fields.RemoteAddress,
				ConnectTimeout: tt.fields.ConnectTimeout,
				SoReUse:        tt.fields.SoReUse,
				transport:      tt.fields.transport,
				udpConn:        tt.fields.udpConn,
				reader:         tt.fields.reader,
			}
			if got := m.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransportInstance_Write(t *testing.T) {
	type fields struct {
		LocalAddress   *net.UDPAddr
		RemoteAddress  *net.UDPAddr
		ConnectTimeout uint32
		SoReUse        bool
		transport      *Transport
		udpConn        *net.UDPConn
		reader         *bufio.Reader
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		setup       func(t *testing.T, fields *fields, args *args)
		manipulator func(t *testing.T, instance *TransportInstance)
		wantErr     bool
	}{
		{
			name:    "write it (no con)",
			wantErr: true,
		},
		{
			name: "write it",
			setup: func(t *testing.T, fields *fields, args *args) {
				listener, err := nettest.NewLocalPacketListener("udp")
				require.NoError(t, err)
				t.Cleanup(func() {
					assert.NoError(t, listener.Close())
				})
				udp, err := net.DialUDP("udp", nil, listener.LocalAddr().(*net.UDPAddr))
				require.NoError(t, err)
				fields.udpConn = udp
			},
			manipulator: func(t *testing.T, instance *TransportInstance) {
				instance.connected.Store(true)
			},
		},
		{
			name: "write it with remote",
			setup: func(t *testing.T, fields *fields, args *args) {
				listener, err := nettest.NewLocalPacketListener("udp")
				require.NoError(t, err)
				t.Cleanup(func() {
					assert.NoError(t, listener.Close())
				})
				remoteAddress := listener.LocalAddr().(*net.UDPAddr)
				fields.RemoteAddress = remoteAddress
				udp, err := net.ListenUDP("udp", nil)
				require.NoError(t, err)
				fields.udpConn = udp
			},
			manipulator: func(t *testing.T, instance *TransportInstance) {
				instance.connected.Store(true)
			},
			args: args{data: []byte{1, 2, 3, 4}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields, &tt.args)
			}
			m := &TransportInstance{
				LocalAddress:   tt.fields.LocalAddress,
				RemoteAddress:  tt.fields.RemoteAddress,
				ConnectTimeout: tt.fields.ConnectTimeout,
				SoReUse:        tt.fields.SoReUse,
				transport:      tt.fields.transport,
				udpConn:        tt.fields.udpConn,
				reader:         tt.fields.reader,
			}
			if tt.manipulator != nil {
				tt.manipulator(t, m)
			}
			if err := m.Write(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
