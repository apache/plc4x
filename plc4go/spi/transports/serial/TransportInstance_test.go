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
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/spi/transports/serial/serialport"
	transportUtils "github.com/apache/plc4x/plc4go/spi/transports/utils"
)

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
		serialPort                       serialport.Port
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
		serialPort                       serialport.Port
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
			if err := m.Connect(t.Context()); (err != nil) != tt.wantErr {
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
		serialPort                       serialport.Port
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
		serialPort                       serialport.Port
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
		serialPort                       serialport.Port
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
		serialPort                       serialport.Port
		reader                           *bufio.Reader
	}
	type args struct {
		data    []byte
		timeout time.Duration
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
			if err := m.Write(t.Context(), tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// stubPort is a minimal serialport.Port for white-box state tests.
type stubPort struct{}

func (stubPort) Read([]byte) (int, error)         { return 0, nil }
func (stubPort) Write(p []byte) (int, error)      { return len(p), nil }
func (stubPort) Close() error                     { return nil }
func (stubPort) SetReadDeadline(time.Time) error  { return nil }
func (stubPort) SetWriteDeadline(time.Time) error { return nil }

// recordingWritePort records SetWriteDeadline calls.
type recordingWritePort struct {
	stubPort
	writeDeadlines []time.Time
}

func (r *recordingWritePort) SetWriteDeadline(t time.Time) error {
	r.writeDeadlines = append(r.writeDeadlines, t)
	return nil
}

func TestTransportInstance_IsConnected_keysOnConnectedFlag(t *testing.T) {
	m := &TransportInstance{}
	m.serialPort = stubPort{}
	assert.False(t, m.IsConnected(), "a port that never connected must not report connected")
	m.connected.Store(true)
	assert.True(t, m.IsConnected())
}

func TestParseAndCreate_OptionsReachInstance(t *testing.T) {
	transport := NewTransport()
	instance, err := transport.CreateTransportInstance(
		url.URL{Scheme: "serial", Path: "/dev/ttyTest0"},
		map[string][]string{
			"baud-rate":    {"19200"},
			"data-bits":    {"7"},
			"parity":       {"even"},
			"stop-bits":    {"2"},
			"read-timeout": {"250"},
		},
	)
	require.NoError(t, err)
	serialInstance := instance.(*TransportInstance)
	assert.Equal(t, "/dev/ttyTest0", serialInstance.SerialPortName)
	assert.Equal(t, uint(19200), serialInstance.BaudRate, "legacy mirror field")
	assert.Equal(t, uint(19200), serialInstance.cfg.port.BaudRate)
	assert.Equal(t, uint(7), serialInstance.cfg.port.DataBits)
	assert.Equal(t, serialport.ParityEven, serialInstance.cfg.port.Parity)
	assert.Equal(t, serialport.StopBitsTwo, serialInstance.cfg.port.StopBits)
	assert.Equal(t, 250*time.Millisecond, serialInstance.cfg.readTimeout)
}

func TestParseAndCreate_InvalidOptionFailsFast(t *testing.T) {
	transport := NewTransport()
	_, err := transport.CreateTransportInstance(
		url.URL{Scheme: "serial", Path: "/dev/ttyTest0"},
		map[string][]string{"parity": {"strong"}},
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), `"parity"`)
}

func TestWrite_FallbackWriteDeadline(t *testing.T) {
	port := &recordingWritePort{}
	m := NewTransportInstanceWithConfig("test", func() serialConfig {
		c := defaultSerialConfig()
		c.writeTimeout = 300 * time.Millisecond
		return c
	}(), NewTransport())
	m.serialPort = port
	m.connected.Store(true)

	// No ctx deadline: fallback must be armed.
	before := time.Now()
	require.NoError(t, m.Write(context.Background(), []byte{0x01}))
	require.Len(t, port.writeDeadlines, 1)
	assert.False(t, port.writeDeadlines[0].Before(before.Add(300*time.Millisecond)))

	// Explicit ctx deadline wins.
	deadline := time.Now().Add(42 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	require.NoError(t, m.Write(ctx, []byte{0x02}))
	require.Len(t, port.writeDeadlines, 2)
	assert.Equal(t, deadline, port.writeDeadlines[1])
}

func TestWrite_ZeroWriteTimeoutClearsDeadline(t *testing.T) {
	port := &recordingWritePort{}
	m := NewTransportInstanceWithConfig("test", func() serialConfig {
		c := defaultSerialConfig()
		c.writeTimeout = 0
		return c
	}(), NewTransport())
	m.serialPort = port
	m.connected.Store(true)

	require.NoError(t, m.Write(context.Background(), []byte{0x01}))
	require.Len(t, port.writeDeadlines, 1)
	assert.True(t, port.writeDeadlines[0].IsZero())
}
