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

package serial

import (
	"bufio"
	"github.com/apache/plc4x/plc4go/spi/transports"
	transportUtils "github.com/apache/plc4x/plc4go/spi/transports/utils"
	"github.com/stretchr/testify/assert"
	"io"
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

func TestNewTransportInstance(t *testing.T) {
	type args struct {
		serialPortName string
		baudRate       uint
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
			if got := NewTransportInstance(tt.args.serialPortName, tt.args.baudRate, tt.args.connectTimeout, tt.args.transport); !assert.Equal(t, tt.want, got) {
				t.Errorf("NewTransportInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransportInstance_Close(t *testing.T) {
	type fields struct {
		DefaultBufferedTransportInstance transportUtils.DefaultBufferedTransportInstance
		SerialPortName                   string
		BaudRate                         uint
		ConnectTimeout                   uint32
		transport                        *Transport
		serialPort                       io.ReadWriteCloser
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
				SerialPortName:                   tt.fields.SerialPortName,
				BaudRate:                         tt.fields.BaudRate,
				ConnectTimeout:                   tt.fields.ConnectTimeout,
				transport:                        tt.fields.transport,
				serialPort:                       tt.fields.serialPort,
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
		SerialPortName                   string
		BaudRate                         uint
		ConnectTimeout                   uint32
		transport                        *Transport
		serialPort                       io.ReadWriteCloser
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
				SerialPortName:                   tt.fields.SerialPortName,
				BaudRate:                         tt.fields.BaudRate,
				ConnectTimeout:                   tt.fields.ConnectTimeout,
				transport:                        tt.fields.transport,
				serialPort:                       tt.fields.serialPort,
				reader:                           tt.fields.reader,
			}
			if err := m.Connect(); (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransportInstance_GetReader(t *testing.T) {
	type fields struct {
		DefaultBufferedTransportInstance transportUtils.DefaultBufferedTransportInstance
		SerialPortName                   string
		BaudRate                         uint
		ConnectTimeout                   uint32
		transport                        *Transport
		serialPort                       io.ReadWriteCloser
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
				SerialPortName:                   tt.fields.SerialPortName,
				BaudRate:                         tt.fields.BaudRate,
				ConnectTimeout:                   tt.fields.ConnectTimeout,
				transport:                        tt.fields.transport,
				serialPort:                       tt.fields.serialPort,
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
		SerialPortName                   string
		BaudRate                         uint
		ConnectTimeout                   uint32
		transport                        *Transport
		serialPort                       io.ReadWriteCloser
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
				SerialPortName:                   tt.fields.SerialPortName,
				BaudRate:                         tt.fields.BaudRate,
				ConnectTimeout:                   tt.fields.ConnectTimeout,
				transport:                        tt.fields.transport,
				serialPort:                       tt.fields.serialPort,
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
		SerialPortName                   string
		BaudRate                         uint
		ConnectTimeout                   uint32
		transport                        *Transport
		serialPort                       io.ReadWriteCloser
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
				SerialPortName:                   tt.fields.SerialPortName,
				BaudRate:                         tt.fields.BaudRate,
				ConnectTimeout:                   tt.fields.ConnectTimeout,
				transport:                        tt.fields.transport,
				serialPort:                       tt.fields.serialPort,
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
		SerialPortName                   string
		BaudRate                         uint
		ConnectTimeout                   uint32
		transport                        *Transport
		serialPort                       io.ReadWriteCloser
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
				SerialPortName:                   tt.fields.SerialPortName,
				BaudRate:                         tt.fields.BaudRate,
				ConnectTimeout:                   tt.fields.ConnectTimeout,
				transport:                        tt.fields.transport,
				serialPort:                       tt.fields.serialPort,
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
		in2          *net.UDPAddr
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
			got, err := m.CreateTransportInstanceForLocalAddress(tt.args.transportUrl, tt.args.options, tt.args.in2)
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
