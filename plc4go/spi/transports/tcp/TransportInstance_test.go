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

package tcp

import (
	"bufio"
	"context"
	"github.com/rs/zerolog/log"
	"net"
	"testing"

	transportUtils "github.com/apache/plc4x/plc4go/spi/transports/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/nettest"
)

func TestNewTcpTransportInstance(t *testing.T) {
	type args struct {
		remoteAddress  *net.TCPAddr
		connectTimeout uint32
		transport      *Transport
	}
	tests := []struct {
		name string
		args args
		want *TransportInstance
	}{
		{
			name: "create it",
			want: func() *TransportInstance {
				ti := &TransportInstance{
					log: log.Logger,
				}
				ti.DefaultBufferedTransportInstance = transportUtils.NewDefaultBufferedTransportInstance(ti)
				return ti
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTcpTransportInstance(tt.args.remoteAddress, tt.args.connectTimeout, tt.args.transport); !assert.Equal(t, tt.want, got) {
				t.Errorf("NewTcpTransportInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransportInstance_Close(t *testing.T) {
	type fields struct {
		DefaultBufferedTransportInstance transportUtils.DefaultBufferedTransportInstance
		RemoteAddress                    *net.TCPAddr
		LocalAddress                     *net.TCPAddr
		ConnectTimeout                   uint32
		transport                        *Transport
		tcpConn                          net.Conn
		reader                           *bufio.Reader
	}
	tests := []struct {
		name        string
		fields      fields
		manipulator func(t *testing.T, instance *TransportInstance)
		wantErr     bool
	}{
		{
			name: "close it (not connected)",
		},
		{
			name: "close it (broken connection)",
			fields: fields{
				tcpConn: &net.TCPConn{},
			},
			manipulator: func(t *testing.T, instance *TransportInstance) {
				instance.connected.Store(true)
			},
			wantErr: true,
		},
		{
			name: "close it",
			fields: fields{
				tcpConn: func() *net.TCPConn {
					listener, err := nettest.NewLocalListener("tcp")
					require.NoError(t, err)
					t.Cleanup(func() {
						assert.NoError(t, listener.Close())
					})
					go func() {
						_, _ = listener.Accept()
					}()
					tcp, err := net.DialTCP("tcp", nil, listener.Addr().(*net.TCPAddr))
					require.NoError(t, err)
					t.Cleanup(func() {
						// As we already closed the connection with the whole method this should error
						assert.Error(t, tcp.Close())
					})
					return tcp
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
				DefaultBufferedTransportInstance: tt.fields.DefaultBufferedTransportInstance,
				RemoteAddress:                    tt.fields.RemoteAddress,
				LocalAddress:                     tt.fields.LocalAddress,
				ConnectTimeout:                   tt.fields.ConnectTimeout,
				transport:                        tt.fields.transport,
				tcpConn:                          tt.fields.tcpConn,
				reader:                           tt.fields.reader,
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
		DefaultBufferedTransportInstance transportUtils.DefaultBufferedTransportInstance
		RemoteAddress                    *net.TCPAddr
		LocalAddress                     *net.TCPAddr
		ConnectTimeout                   uint32
		transport                        *Transport
		tcpConn                          net.Conn
		reader                           *bufio.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "connect it (failing)",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				DefaultBufferedTransportInstance: tt.fields.DefaultBufferedTransportInstance,
				RemoteAddress:                    tt.fields.RemoteAddress,
				LocalAddress:                     tt.fields.LocalAddress,
				ConnectTimeout:                   tt.fields.ConnectTimeout,
				transport:                        tt.fields.transport,
				tcpConn:                          tt.fields.tcpConn,
				reader:                           tt.fields.reader,
			}
			if err := m.Connect(); (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransportInstance_ConnectWithContext(t *testing.T) {
	type fields struct {
		DefaultBufferedTransportInstance transportUtils.DefaultBufferedTransportInstance
		RemoteAddress                    *net.TCPAddr
		LocalAddress                     *net.TCPAddr
		ConnectTimeout                   uint32
		transport                        *Transport
		tcpConn                          net.Conn
		reader                           *bufio.Reader
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
			fields: fields{
				RemoteAddress: func() *net.TCPAddr {
					listener, err := nettest.NewLocalListener("tcp")
					require.NoError(t, err)
					t.Cleanup(func() {
						assert.NoError(t, listener.Close())
					})
					go func() {
						_, _ = listener.Accept()
					}()
					return listener.Addr().(*net.TCPAddr)
				}(),
			},
			args: args{ctx: context.Background()},
		},
		{
			name: "connect it (non existing address)",
			fields: fields{
				RemoteAddress: &net.TCPAddr{},
			},
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				DefaultBufferedTransportInstance: tt.fields.DefaultBufferedTransportInstance,
				RemoteAddress:                    tt.fields.RemoteAddress,
				LocalAddress:                     tt.fields.LocalAddress,
				ConnectTimeout:                   tt.fields.ConnectTimeout,
				transport:                        tt.fields.transport,
				tcpConn:                          tt.fields.tcpConn,
				reader:                           tt.fields.reader,
			}
			if err := m.ConnectWithContext(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ConnectWithContext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransportInstance_GetReader(t *testing.T) {
	type fields struct {
		DefaultBufferedTransportInstance transportUtils.DefaultBufferedTransportInstance
		RemoteAddress                    *net.TCPAddr
		LocalAddress                     *net.TCPAddr
		ConnectTimeout                   uint32
		transport                        *Transport
		tcpConn                          net.Conn
		reader                           *bufio.Reader
	}
	tests := []struct {
		name   string
		fields fields
		want   *bufio.Reader
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				DefaultBufferedTransportInstance: tt.fields.DefaultBufferedTransportInstance,
				RemoteAddress:                    tt.fields.RemoteAddress,
				LocalAddress:                     tt.fields.LocalAddress,
				ConnectTimeout:                   tt.fields.ConnectTimeout,
				transport:                        tt.fields.transport,
				tcpConn:                          tt.fields.tcpConn,
				reader:                           tt.fields.reader,
			}
			if got := m.GetReader(); !assert.Equal(t, tt.want, got) {
				t.Errorf("GetReader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransportInstance_IsConnected(t *testing.T) {
	type fields struct {
		DefaultBufferedTransportInstance transportUtils.DefaultBufferedTransportInstance
		RemoteAddress                    *net.TCPAddr
		LocalAddress                     *net.TCPAddr
		ConnectTimeout                   uint32
		transport                        *Transport
		tcpConn                          net.Conn
		reader                           *bufio.Reader
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
				DefaultBufferedTransportInstance: tt.fields.DefaultBufferedTransportInstance,
				RemoteAddress:                    tt.fields.RemoteAddress,
				LocalAddress:                     tt.fields.LocalAddress,
				ConnectTimeout:                   tt.fields.ConnectTimeout,
				transport:                        tt.fields.transport,
				tcpConn:                          tt.fields.tcpConn,
				reader:                           tt.fields.reader,
			}
			if got := m.IsConnected(); got != tt.want {
				t.Errorf("IsConnected() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransportInstance_String(t *testing.T) {
	type fields struct {
		DefaultBufferedTransportInstance transportUtils.DefaultBufferedTransportInstance
		RemoteAddress                    *net.TCPAddr
		LocalAddress                     *net.TCPAddr
		ConnectTimeout                   uint32
		transport                        *Transport
		tcpConn                          net.Conn
		reader                           *bufio.Reader
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get it",
			want: "tcp:<nil>",
		},
		{
			name: "get it too",
			fields: fields{
				LocalAddress:  &net.TCPAddr{IP: net.IPv4(1, 2, 3, 4), Port: 5},
				RemoteAddress: &net.TCPAddr{IP: net.IPv4(6, 7, 8, 9), Port: 10},
			},
			want: "tcp:1.2.3.4:5->6.7.8.9:10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				DefaultBufferedTransportInstance: tt.fields.DefaultBufferedTransportInstance,
				RemoteAddress:                    tt.fields.RemoteAddress,
				LocalAddress:                     tt.fields.LocalAddress,
				ConnectTimeout:                   tt.fields.ConnectTimeout,
				transport:                        tt.fields.transport,
				tcpConn:                          tt.fields.tcpConn,
				reader:                           tt.fields.reader,
			}
			if got := m.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransportInstance_Write(t *testing.T) {
	type fields struct {
		DefaultBufferedTransportInstance transportUtils.DefaultBufferedTransportInstance
		RemoteAddress                    *net.TCPAddr
		LocalAddress                     *net.TCPAddr
		ConnectTimeout                   uint32
		transport                        *Transport
		tcpConn                          net.Conn
		reader                           *bufio.Reader
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
			name:    "write it (failing)",
			wantErr: true,
		},
		{
			name: "write it (failing with con)",
			fields: fields{
				tcpConn: &net.TCPConn{},
			},
			wantErr: true,
		},
		{
			name: "write it",
			setup: func(t *testing.T, fields *fields, args *args) {
				listener, err := nettest.NewLocalListener("tcp")
				require.NoError(t, err)
				t.Cleanup(func() {
					assert.NoError(t, listener.Close())
				})
				go func() {
					_, _ = listener.Accept()
				}()
				tcp, err := net.DialTCP("tcp", nil, listener.Addr().(*net.TCPAddr))
				require.NoError(t, err)
				t.Cleanup(func() {
					assert.NoError(t, tcp.Close())
				})
				fields.tcpConn = tcp
			},
			manipulator: func(t *testing.T, instance *TransportInstance) {
				instance.connected.Store(true)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields, &tt.args)
			}
			m := &TransportInstance{
				DefaultBufferedTransportInstance: tt.fields.DefaultBufferedTransportInstance,
				RemoteAddress:                    tt.fields.RemoteAddress,
				LocalAddress:                     tt.fields.LocalAddress,
				ConnectTimeout:                   tt.fields.ConnectTimeout,
				transport:                        tt.fields.transport,
				tcpConn:                          tt.fields.tcpConn,
				reader:                           tt.fields.reader,
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
