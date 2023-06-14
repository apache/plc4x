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

package pcap

import (
	"bufio"
	"github.com/apache/plc4x/plc4go/spi/transports"
	transportUtils "github.com/apache/plc4x/plc4go/spi/transports/utils"
	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
	"github.com/gopacket/gopacket/pcap"
	"github.com/gopacket/gopacket/pcapgo"
	"github.com/stretchr/testify/assert"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestNewPcapTransportInstance(t *testing.T) {
	type args struct {
		transportFile string
		transportType TransportType
		portRange     string
		speedFactor   float32
		transport     *Transport
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
			if got := NewPcapTransportInstance(tt.args.transportFile, tt.args.transportType, tt.args.portRange, tt.args.speedFactor, tt.args.transport); !assert.Equal(t, tt.want, got) {
				t.Errorf("NewPcapTransportInstance() = %v, want %v", got, tt.want)
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
		transportFile                    string
		transportType                    TransportType
		portRange                        string
		speedFactor                      float32
		transport                        *Transport
		handle                           *pcap.Handle
		reader                           *bufio.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "close it",
			fields: fields{
				handle: &pcap.Handle{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				DefaultBufferedTransportInstance: tt.fields.DefaultBufferedTransportInstance,
				transportFile:                    tt.fields.transportFile,
				transportType:                    tt.fields.transportType,
				portRange:                        tt.fields.portRange,
				speedFactor:                      tt.fields.speedFactor,
				transport:                        tt.fields.transport,
				handle:                           tt.fields.handle,
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
		transportFile                    string
		transportType                    TransportType
		portRange                        string
		speedFactor                      float32
		transport                        *Transport
		handle                           *pcap.Handle
		reader                           *bufio.Reader
	}
	tests := []struct {
		name        string
		fields      fields
		mockSetup   func(t *testing.T, fields *fields)
		manipulator func(t *testing.T, transportInstance *TransportInstance)
		wantErr     bool
	}{
		{
			name: "already connected",
			manipulator: func(t *testing.T, transportInstance *TransportInstance) {
				transportInstance.connected.Store(true)
			},
			wantErr: true,
		},
		{
			name:    "connect no file",
			wantErr: true,
		},
		{
			name: "connect with file",
			mockSetup: func(t *testing.T, fields *fields) {
				fields.transportFile = createPcap(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields)
			}
			m := &TransportInstance{
				DefaultBufferedTransportInstance: tt.fields.DefaultBufferedTransportInstance,
				transportFile:                    tt.fields.transportFile,
				transportType:                    tt.fields.transportType,
				portRange:                        tt.fields.portRange,
				speedFactor:                      tt.fields.speedFactor,
				transport:                        tt.fields.transport,
				handle:                           tt.fields.handle,
				reader:                           tt.fields.reader,
			}
			if tt.manipulator != nil {
				tt.manipulator(t, m)
			}
			if err := m.Connect(); (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Cleanup(func() {
				if err := m.Close(); err != nil {
					t.Logf("Error during close %v", err)
				}
			})
		})
	}
}

func TestTransportInstance_GetReader(t *testing.T) {
	type fields struct {
		DefaultBufferedTransportInstance transportUtils.DefaultBufferedTransportInstance
		transportFile                    string
		transportType                    TransportType
		portRange                        string
		speedFactor                      float32
		transport                        *Transport
		handle                           *pcap.Handle
		reader                           *bufio.Reader
	}
	tests := []struct {
		name   string
		fields fields
		want   *bufio.Reader
	}{
		{
			name: "get reader",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				DefaultBufferedTransportInstance: tt.fields.DefaultBufferedTransportInstance,
				transportFile:                    tt.fields.transportFile,
				transportType:                    tt.fields.transportType,
				portRange:                        tt.fields.portRange,
				speedFactor:                      tt.fields.speedFactor,
				transport:                        tt.fields.transport,
				handle:                           tt.fields.handle,
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
		transportFile                    string
		transportType                    TransportType
		portRange                        string
		speedFactor                      float32
		transport                        *Transport
		handle                           *pcap.Handle
		reader                           *bufio.Reader
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "is connected",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				DefaultBufferedTransportInstance: tt.fields.DefaultBufferedTransportInstance,
				transportFile:                    tt.fields.transportFile,
				transportType:                    tt.fields.transportType,
				portRange:                        tt.fields.portRange,
				speedFactor:                      tt.fields.speedFactor,
				transport:                        tt.fields.transport,
				handle:                           tt.fields.handle,
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
		transportFile                    string
		transportType                    TransportType
		portRange                        string
		speedFactor                      float32
		transport                        *Transport
		handle                           *pcap.Handle
		reader                           *bufio.Reader
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get the string",
			fields: fields{
				transportFile: "plc4x",
				portRange:     "is-the-best",
				speedFactor:   3.14,
			},
			want: "pcap:plc4x(is-the-best)x3.140000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				DefaultBufferedTransportInstance: tt.fields.DefaultBufferedTransportInstance,
				transportFile:                    tt.fields.transportFile,
				transportType:                    tt.fields.transportType,
				portRange:                        tt.fields.portRange,
				speedFactor:                      tt.fields.speedFactor,
				transport:                        tt.fields.transport,
				handle:                           tt.fields.handle,
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
		transportFile                    string
		transportType                    TransportType
		portRange                        string
		speedFactor                      float32
		transport                        *Transport
		handle                           *pcap.Handle
		reader                           *bufio.Reader
	}
	type args struct {
		in0 []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "we can't write",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TransportInstance{
				DefaultBufferedTransportInstance: tt.fields.DefaultBufferedTransportInstance,
				transportFile:                    tt.fields.transportFile,
				transportType:                    tt.fields.transportType,
				portRange:                        tt.fields.portRange,
				speedFactor:                      tt.fields.speedFactor,
				transport:                        tt.fields.transport,
				handle:                           tt.fields.handle,
				reader:                           tt.fields.reader,
			}
			if err := m.Write(tt.args.in0); (err != nil) != tt.wantErr {
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
			name: "create it",
			args: args{
				options: map[string][]string{
					"transport-type":       {"pcap"},
					"transport-port-range": {"1-3"},
					"speed-factor":         {"1.5"},
				},
			},
			want: func() transports.TransportInstance {
				ti := &TransportInstance{
					transportType: "pcap",
					speedFactor:   1.5,
					transport:     NewTransport(),
					portRange:     "1-3",
				}
				ti.DefaultBufferedTransportInstance = transportUtils.NewDefaultBufferedTransportInstance(ti)
				return ti
			}(),
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
			want: "pcap",
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
			want: "PCAP(NG) Playback Transport",
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
			name: "get it",
			want: "pcap(PCAP(NG) Playback Transport)",
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

func createPcap(t *testing.T) string {
	tempFile, err := os.CreateTemp(t.TempDir(), "some*.pcap")
	if err != nil {
		t.Fatal(err)
	}

	w := pcapgo.NewWriter(tempFile)
	if err := w.WriteFileHeader(65536, layers.LinkTypeEthernet); err != nil {
		t.Fatal(err)
	}
	ci := gopacket.CaptureInfo{
		Timestamp:     time.Unix(0x01020304, 0xAA*1000),
		Length:        0xABCD,
		CaptureLength: 10,
	}
	data := []byte{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	if err := w.WritePacket(ci, data); err != nil {
		t.Fatal(err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatal(err)
	}
	return tempFile.Name()
}
