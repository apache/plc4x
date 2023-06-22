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
	"fmt"
	"github.com/apache/plc4x/plc4go/spi/transports"
	transportUtils "github.com/apache/plc4x/plc4go/spi/transports/utils"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"net"
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
				log: log.Logger,
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
					log:            log.Logger,
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
					log:            log.Logger,
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
					log:            log.Logger,
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
					log:            log.Logger,
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
					log:            log.Logger,
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
			m := Transport{
				log: log.Logger,
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
