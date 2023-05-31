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
	"fmt"
	"github.com/apache/plc4x/plc4go/spi/transports"
	transportUtils "github.com/apache/plc4x/plc4go/spi/transports/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/nettest"
	"net"
	"net/url"
	"testing"
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
				ti := &TransportInstance{}
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

func TestNewTransport(t *testing.T) {
	tests := []struct {
		name string
		want *Transport
	}{
		{
			name: "create it",
			want: &Transport{},
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
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "close it (no conn)",
		},
		{
			name: "close it (broken connection)",
			fields: fields{
				tcpConn: &net.TCPConn{},
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
		name    string
		fields  fields
		args    args
		wantErr bool
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
						assert.NoError(t, tcp.Close())
					})
					return tcp
				}(),
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
			if err := m.Write(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransport_CreateTransportInstance(t *testing.T) {
	type args struct {
		transportUrl url.URL
		options      map[string][]string
	}
	tests := []struct {
		name    string
		args    args
		want    transports.TransportInstance
		wantErr bool
	}{
		{
			name: "Create it",
			want: func() transports.TransportInstance {
				tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", "", 0))
				assert.NoError(t, err)
				ti := &TransportInstance{
					transport:      NewTransport(),
					RemoteAddress:  tcpAddr,
					ConnectTimeout: 1000,
				}
				ti.DefaultBufferedTransportInstance = transportUtils.NewDefaultBufferedTransportInstance(ti)
				return ti
			}(),
		},
		{
			name: "Create it with transport url",
			args: args{
				transportUrl: url.URL{Host: "127.0.0.1:123"},
			},
			want: func() transports.TransportInstance {
				tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", 123))
				assert.NoError(t, err)
				ti := &TransportInstance{
					transport:      NewTransport(),
					RemoteAddress:  tcpAddr,
					ConnectTimeout: 1000,
				}
				ti.DefaultBufferedTransportInstance = transportUtils.NewDefaultBufferedTransportInstance(ti)
				return ti
			}(),
		},
		{
			name: "Create it with transport url (named host)",
			args: args{
				transportUrl: url.URL{Host: "localhost:123"},
			},
			want: func() transports.TransportInstance {
				tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", "localhost", 123))
				assert.NoError(t, err)
				ti := &TransportInstance{
					transport:      NewTransport(),
					RemoteAddress:  tcpAddr,
					ConnectTimeout: 1000,
				}
				ti.DefaultBufferedTransportInstance = transportUtils.NewDefaultBufferedTransportInstance(ti)
				return ti
			}(),
		},
		{
			name: "Create it with transport url (without port)",
			args: args{
				transportUrl: url.URL{Host: "127.0.0.1"},
			},
			wantErr: true,
		},
		{
			name: "Create it with transport url (with nonsense port)",
			args: args{
				transportUrl: url.URL{Host: "127.0.0.1:banana"},
			},
			wantErr: true,
		},
		{
			name: "Create it with transport url (with default port)",
			args: args{
				transportUrl: url.URL{Host: "127.0.0.1"},
				options: map[string][]string{
					"defaultTcpPort": {"123"},
				},
			},
			want: func() transports.TransportInstance {
				tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", 123))
				assert.NoError(t, err)
				ti := &TransportInstance{
					transport:      NewTransport(),
					RemoteAddress:  tcpAddr,
					ConnectTimeout: 1000,
				}
				ti.DefaultBufferedTransportInstance = transportUtils.NewDefaultBufferedTransportInstance(ti)
				return ti
			}(),
		},
		{
			name: "Create it with transport url (with broken default port)",
			args: args{
				transportUrl: url.URL{Host: "127.0.0.1"},
				options: map[string][]string{
					"defaultTcpPort": {"default"},
				},
			},
			wantErr: true,
		},
		{
			name: "Create it with transport url (with default port and connection timeout)",
			args: args{
				transportUrl: url.URL{Host: "127.0.0.1"},
				options: map[string][]string{
					"defaultTcpPort":  {"123"},
					"connect-timeout": {"123"},
				},
			},
			want: func() transports.TransportInstance {
				tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", 123))
				assert.NoError(t, err)
				ti := &TransportInstance{
					transport:      NewTransport(),
					RemoteAddress:  tcpAddr,
					ConnectTimeout: 123,
				}
				ti.DefaultBufferedTransportInstance = transportUtils.NewDefaultBufferedTransportInstance(ti)
				return ti
			}(),
		},
		{
			name: "Create it with transport url (with default port and connection timeout broken)",
			args: args{
				transportUrl: url.URL{Host: "127.0.0.1"},
				options: map[string][]string{
					"defaultTcpPort":  {"123"},
					"connect-timeout": {"banana"},
				},
			},
			wantErr: true,
		},
		{
			name: "Create it with unresolvable host",
			args: args{
				transportUrl: url.URL{Host: "plc4xhostnothere"},
				options: map[string][]string{
					"defaultTcpPort": {"123"},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Transport{}
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
	tests := []struct {
		name string
		want string
	}{
		{
			name: "get it",
			want: "tcp",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Transport{}
			if got := m.GetTransportCode(); got != tt.want {
				t.Errorf("GetTransportCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransport_GetTransportName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "get it",
			want: "TCP/IP Socket Transport",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Transport{}
			if got := m.GetTransportName(); got != tt.want {
				t.Errorf("GetTransportName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransport_String(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "get the string",
			want: "tcp(TCP/IP Socket Transport)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Transport{}
			if got := m.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
