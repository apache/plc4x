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
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"io"
	"net/url"
	"strconv"
	"sync"
	"time"
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
	var speedFactor float32 = 1.0
	if val, ok := options["speed-factor"]; ok {
		if parsedSpeedFactory, err := strconv.ParseFloat(val[0], 32); err != nil {
			return nil, errors.Wrap(err, "error parsing speed-factor")
		} else {
			speedFactor = float32(parsedSpeedFactory)
		}
	}

	transportInstance := NewPcapTransportInstance(transportUrl.Path, transportType, portRange, speedFactor, &m)

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
	speedFactor   float32
	connected     bool
	transport     *Transport
	reader        *bufio.Reader
	handle        *pcap.Handle
	mutex         sync.Mutex
}

func NewPcapTransportInstance(transportFile string, transportType TransportType, portRange string, speedFactor float32, transport *Transport) *TransportInstance {
	return &TransportInstance{
		transportFile: transportFile,
		transportType: transportType,
		portRange:     portRange,
		speedFactor:   speedFactor,
		transport:     transport,
	}
}

func (m *TransportInstance) Connect() error {
	if m.connected {
		return errors.New("Already connected")
	}
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
	buffer := new(bytes.Buffer)
	m.reader = bufio.NewReader(buffer)

	go func(m *TransportInstance, buffer *bytes.Buffer) {
		var lastPacketTime *time.Time
		for m.connected {
			packetData, captureInfo, err := m.handle.ReadPacketData()
			log.Info().Msgf("Read new package %#v", captureInfo)
			if err != nil {
				if err == io.EOF {
					log.Info().Msg("Done reading pcap")
					break
				}
				log.Warn().Err(err).Msg("Error reading")
				panic(err)
			}
			if lastPacketTime != nil && m.speedFactor != 0 {
				timeToSleep := captureInfo.Timestamp.Sub(*lastPacketTime)
				timeToSleep = time.Duration(int64(float64(timeToSleep) / float64(m.speedFactor)))
				log.Debug().Msgf("Sleeping for %v (Speed factor %fx)", timeToSleep, m.speedFactor)
				time.Sleep(timeToSleep)
			}

			// Decode a packet
			packet := gopacket.NewPacket(packetData, layers.LayerTypeEthernet, gopacket.Default)
			var payload []byte
			switch m.transportType {
			case TCP:
				if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
					tcp, _ := tcpLayer.(*layers.TCP)
					payload = tcp.Payload
					log.Debug().Msgf("TCP: From src port %d to dst port %d", tcp.SrcPort, tcp.DstPort)
				} else {
					continue
				}
			case UDP:
				if tcpLayer := packet.Layer(layers.LayerTypeUDP); tcpLayer != nil {
					udp, _ := tcpLayer.(*layers.UDP)
					payload = udp.Payload
					log.Debug().Msgf("UDP: From src port %d to dst port %d", udp.SrcPort, udp.DstPort)
				} else {
					continue
				}
			}
			buffer.Grow(len(payload))
			buffer.Write(payload)
			lastPacketTime = &captureInfo.Timestamp
		}
	}(m, buffer)

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
