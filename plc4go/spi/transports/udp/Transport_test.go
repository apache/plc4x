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
	"fmt"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			name: "create it",
			want: func() transports.TransportInstance {
				remoteAddress, err := net.ResolveUDPAddr("udp", ":0")
				require.NoError(t, err)
				return &TransportInstance{
					ConnectTimeout: 1000,
					RemoteAddress:  remoteAddress,
					transport:      NewTransport(),
					log:            log.Logger,
				}
			}(),
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

func TestTransport_CreateTransportInstanceForLocalAddress(t *testing.T) {
	type args struct {
		transportUrl url.URL
		options      map[string][]string
		localAddress *net.UDPAddr
	}
	tests := []struct {
		name    string
		args    args
		want    transports.TransportInstance
		wantErr bool
	}{
		{
			name: "Create it",
			want: &TransportInstance{
				transport:      NewTransport(),
				RemoteAddress:  &net.UDPAddr{},
				ConnectTimeout: 1000,
				log:            log.Logger,
			},
		},
		{
			name: "Create it with transport url",
			args: args{
				transportUrl: url.URL{Host: "127.0.0.1:123"},
			},
			want: func() transports.TransportInstance {
				udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", "127.0.0.1", 123))
				assert.NoError(t, err)
				ti := &TransportInstance{
					transport:      NewTransport(),
					RemoteAddress:  udpAddr,
					ConnectTimeout: 1000,
					log:            log.Logger,
				}
				return ti
			}(),
		},
		{
			name: "Create it with transport url (named host)",
			args: args{
				transportUrl: url.URL{Host: "localhost:123"},
			},
			want: func() transports.TransportInstance {
				udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", "localhost", 123))
				assert.NoError(t, err)
				ti := &TransportInstance{
					transport:      NewTransport(),
					RemoteAddress:  udpAddr,
					ConnectTimeout: 1000,
					log:            log.Logger,
				}
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
					"defaultUdpPort": {"123"},
				},
			},
			want: func() transports.TransportInstance {
				udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", "127.0.0.1", 123))
				assert.NoError(t, err)
				ti := &TransportInstance{
					transport:      NewTransport(),
					RemoteAddress:  udpAddr,
					ConnectTimeout: 1000,
					log:            log.Logger,
				}
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
			name: "Create it with transport url (with default port and connection timeout and reuse)",
			args: args{
				transportUrl: url.URL{Host: "127.0.0.1"},
				options: map[string][]string{
					"defaultUdpPort":  {"123"},
					"connect-timeout": {"123"},
					"so-reuse":        {"true"},
				},
			},
			want: func() transports.TransportInstance {
				udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", "127.0.0.1", 123))
				assert.NoError(t, err)
				ti := &TransportInstance{
					transport:      NewTransport(),
					RemoteAddress:  udpAddr,
					ConnectTimeout: 123,
					SoReUse:        true,
					log:            log.Logger,
				}
				return ti
			}(),
		},
		{
			name: "Create it with transport url (with default port and connection timeout and reuse broken)",
			args: args{
				transportUrl: url.URL{Host: "127.0.0.1"},
				options: map[string][]string{
					"defaultUdpPort":  {"123"},
					"connect-timeout": {"123"},
					"so-reuse":        {"banana"},
				},
			},
			wantErr: true,
		},
		{
			name: "Create it with transport url (with default port and connection timeout broken)",
			args: args{
				transportUrl: url.URL{Host: "127.0.0.1"},
				options: map[string][]string{
					"defaultUdpPort":  {"123"},
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
					"defaultUdpPort": {"123"},
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
			got, err := m.CreateTransportInstanceForLocalAddress(tt.args.transportUrl, tt.args.options, tt.args.localAddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTransportInstanceForLocalAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("CreateTransportInstanceForLocalAddress() got = %v, want %v", got, tt.want)
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
			want: "udp",
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
			want: "UDP Datagram Transport",
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
			name: "string it",
			want: "udp(UDP Datagram Transport)",
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
