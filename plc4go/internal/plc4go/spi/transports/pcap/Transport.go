/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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
	"bytes"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports"
	"github.com/google/gopacket/pcap"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/url"
	"sync"
)

type TransportType string

const (
	UDP TransportType = "udp"
	TCP TransportType = "tcp"
)

type Transport struct {
}

func NewTransport() *Transport {
	return &Transport{}
}

func (m Transport) GetTransportCode() string {
	return "pcap"
}

func (m Transport) GetTransportName() string {
	return "PCAP(NG) Playback Transport"
}

func (m Transport) CreateTransportInstance(transportUrl url.URL, options map[string][]string) (transports.TransportInstance, error) {
	var transportType = TCP
	if val, ok := options["transport-type"]; ok {
		transportType = TransportType(val[0])
	}
	var portRange = ""
	if val, ok := options["transport-port-range"]; ok {
		portRange = val[0]
	}

	transportInstance := NewPcapTransportInstance(transportUrl.Path, transportType, portRange, &m)

	castFunc := func(typ interface{}) (transports.TransportInstance, error) {
		if transportInstance, ok := typ.(transports.TransportInstance); ok {
			return transportInstance, nil
		}
		return nil, errors.New("couldn't cast to TransportInstance")
	}
	return castFunc(transportInstance)
}

type TransportInstance struct {
	transportFile string
	transportType TransportType
	portRange     string
	connected     bool
	transport     *Transport
	reader        *bufio.Reader
	handle        *pcap.Handle
	mutex         sync.Mutex
}

func NewPcapTransportInstance(transportFile string, transportType TransportType, portRange string, transport *Transport) *TransportInstance {
	return &TransportInstance{
		transportFile: transportFile,
		transportType: transportType,
		portRange:     portRange,
		transport:     transport,
	}
}

func (m *TransportInstance) Connect() error {
	handle, err := pcap.OpenOffline(m.transportFile)
	if err != nil {
		return err
	}
	filter := string(m.transportType)
	if m.portRange != "" {
		filter += " dst port " + m.portRange
	}
	if err := handle.SetBPFFilter(filter); err != nil {
		return err
	}
	m.handle = handle
	m.connected = true
	return nil
}

func (m *TransportInstance) Close() error {
	m.handle.Close()
	m.connected = false
	return nil
}

func (m *TransportInstance) IsConnected() bool {
	return m.connected
}

func (m *TransportInstance) GetNumReadableBytes() (uint32, error) {
	if err := m.checkForNextPackage(); err != nil {
		return 0, err
	}
	if m.reader == nil {
		return 0, nil
	}
	_, _ = m.reader.Peek(1)
	return uint32(m.reader.Buffered()), nil
}

func (m *TransportInstance) PeekReadableBytes(numBytes uint32) ([]uint8, error) {
	if m.reader == nil {
		return nil, errors.New("error peeking from transport. No reader available")
	}
	return m.reader.Peek(int(numBytes))
}

func (m *TransportInstance) Read(numBytes uint32) ([]uint8, error) {
	if m.reader == nil {
		return nil, errors.New("error reading from transport. No reader available")
	}
	data := make([]uint8, numBytes)
	for i := uint32(0); i < numBytes; i++ {
		val, err := m.reader.ReadByte()
		if err != nil {
			return nil, errors.Wrap(err, "error reading")
		}
		data[i] = val
	}
	return data, nil
}

func (m *TransportInstance) Write(_ []uint8) error {
	panic("Write to pcap not supported")
}

func (m *TransportInstance) checkForNextPackage() error {
	m.mutex.Lock()
	defer m.mutex.Lock()
	// TODO: this will only work with the first packet and after this we are done
	if m.reader == nil {
		packetData, captureInfo, err := m.handle.ReadPacketData()
		log.Info().Msgf("Read new package %v", captureInfo)
		if err != nil {
			return err
		}
		m.reader = bufio.NewReader(bytes.NewBuffer(packetData))
	}
	return nil
}
