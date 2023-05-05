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
	"github.com/apache/plc4x/plc4go/spi/transports"
	"net"
	"net/url"
	"reflect"
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTcpTransportInstance(tt.args.remoteAddress, tt.args.connectTimeout, tt.args.transport); !reflect.DeepEqual(got, tt.want) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTransport(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTransport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransportInstance_Close(t *testing.T) {
	type fields struct {
		DefaultBufferedTransportInstance transports.DefaultBufferedTransportInstance
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
		// TODO: Add test cases.
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
		DefaultBufferedTransportInstance transports.DefaultBufferedTransportInstance
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
		// TODO: Add test cases.
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
		DefaultBufferedTransportInstance transports.DefaultBufferedTransportInstance
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
		// TODO: Add test cases.
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
		DefaultBufferedTransportInstance transports.DefaultBufferedTransportInstance
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
		// TODO: Add test cases.
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
			if got := m.GetReader(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetReader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransportInstance_IsConnected(t *testing.T) {
	type fields struct {
		DefaultBufferedTransportInstance transports.DefaultBufferedTransportInstance
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
		// TODO: Add test cases.
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
		DefaultBufferedTransportInstance transports.DefaultBufferedTransportInstance
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
		// TODO: Add test cases.
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
		DefaultBufferedTransportInstance transports.DefaultBufferedTransportInstance
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Transport{}
			got, err := m.CreateTransportInstance(tt.args.transportUrl, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTransportInstance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
