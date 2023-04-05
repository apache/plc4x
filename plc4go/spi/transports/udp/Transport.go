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
	"bufio"
	"context"
	"fmt"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/libp2p/go-reuseport"
	"github.com/pkg/errors"
	"net"
	"net/url"
	"regexp"
	"strconv"
)

type Transport struct {
}

func NewTransport() *Transport {
	return &Transport{}
}

func (m Transport) GetTransportCode() string {
	return "udp"
}

func (m Transport) GetTransportName() string {
	return "UDP Datagram Transport"
}

func (m Transport) CreateTransportInstance(transportUrl url.URL, options map[string][]string) (transports.TransportInstance, error) {
	return m.CreateTransportInstanceForLocalAddress(transportUrl, options, nil)
}

func (m Transport) CreateTransportInstanceForLocalAddress(transportUrl url.URL, options map[string][]string, localAddress *net.UDPAddr) (transports.TransportInstance, error) {
	connectionStringRegexp := regexp.MustCompile(`^((?P<ip>[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})|(?P<hostname>[a-zA-Z0-9.\-]+))(:(?P<port>[0-9]{1,5}))?`)
	var remoteAddressString string
	var remotePort int
	if match := utils.GetSubgroupMatches(connectionStringRegexp, transportUrl.Host); match != nil {
		if val, ok := match["ip"]; ok && len(val) > 0 {
			remoteAddressString = val
		} else if val, ok := match["hostname"]; ok && len(val) > 0 {
			remoteAddressString = val
		} else {
			return nil, errors.New("missing hostname or ip to connect")
		}

		if val, ok := match["port"]; ok && len(val) > 0 {
			portVal, err := strconv.Atoi(val)
			if err != nil {
				return nil, errors.Wrap(err, "error setting port")
			}
			remotePort = portVal
		} else if val, ok := options["defaultUdpPort"]; ok && len(val) > 0 {
			portVal, err := strconv.Atoi(val[0])
			if err != nil {
				return nil, errors.Wrap(err, "error setting default udp port")
			}
			remotePort = portVal
		} else {
			return nil, errors.New("error setting port. No explicit or default port provided")
		}
	}
	var connectTimeout uint32 = 1000
	if val, ok := options["connect-timeout"]; ok {
		if parsedConnectTimeout, err := strconv.ParseUint(val[0], 10, 32); err != nil {
			return nil, errors.Wrap(err, "error setting connect-timeout")
		} else {
			connectTimeout = uint32(parsedConnectTimeout)
		}
	}

	var soReUse bool
	if val, ok := options["so-reuse"]; ok {
		if parseBool, err := strconv.ParseBool(val[0]); err != nil {
			return nil, errors.Wrap(err, "error setting so-reuse")
		} else {
			soReUse = parseBool
		}
	}

	// Potentially resolve the ip address, if a hostname was provided
	remoteAddress, err := net.ResolveUDPAddr("udp", remoteAddressString+":"+strconv.Itoa(remotePort))
	if err != nil {
		return nil, errors.Wrap(err, "error resolving typ address")
	}

	return NewTransportInstance(localAddress, remoteAddress, connectTimeout, soReUse, &m), nil
}

func (m Transport) String() string {
	return m.GetTransportCode() + "(" + m.GetTransportName() + ")"
}

type TransportInstance struct {
	LocalAddress   *net.UDPAddr
	RemoteAddress  *net.UDPAddr
	ConnectTimeout uint32
	SoReUse        bool
	transport      *Transport
	udpConn        *net.UDPConn
	reader         *bufio.Reader
}

func NewTransportInstance(localAddress *net.UDPAddr, remoteAddress *net.UDPAddr, connectTimeout uint32, soReUse bool, transport *Transport) *TransportInstance {
	return &TransportInstance{
		LocalAddress:   localAddress,
		RemoteAddress:  remoteAddress,
		ConnectTimeout: connectTimeout,
		SoReUse:        soReUse,
		transport:      transport,
	}
}

func (m *TransportInstance) Connect() error {
	return m.ConnectWithContext(context.Background())
}

func (m *TransportInstance) ConnectWithContext(ctx context.Context) error {
	// If we haven't provided a local address, have the system figure it out by dialing
	// the remote address and then using that connections local address as local address.
	if m.LocalAddress == nil {
		var d net.Dialer
		udpTest, err := d.DialContext(ctx, "udp", m.RemoteAddress.String())
		if err != nil {
			return errors.Wrap(err, "error connecting to remote address")
		}
		m.LocalAddress = udpTest.LocalAddr().(*net.UDPAddr)
		err = udpTest.Close()
		if err != nil {
			return errors.Wrap(err, "error closing test-connection")
		}
	}

	// "connect" to the remote
	var err error
	if m.RemoteAddress != nil {
		if m.udpConn, err = net.DialUDP("udp", m.LocalAddress, m.RemoteAddress); err != nil {
			return errors.Wrap(err, "error connecting to remote address")
		}
	} else if m.SoReUse {
		if packetConn, err := reuseport.ListenPacket("udp", m.LocalAddress.String()); err != nil {
			return errors.Wrap(err, "error connecting to local address")
		} else {
			m.udpConn = packetConn.(*net.UDPConn)
		}
	} else {
		if m.udpConn, err = net.ListenUDP("udp", m.LocalAddress); err != nil {
			return errors.Wrap(err, "error connecting to local address")
		}
	}

	// TODO: Start a worker that uses m.udpConn.ReadFromUDP() to fill a buffer
	/*go func() {
	    buf := make([]byte, 1024)
	    for {
	        rsize, raddr, err := m.udpConn.ReadFromUDP(buf)
	        if err != nil {
	            fmt.Printf("Got %d bytes from %v: %v", rsize, raddr, buf)
	        }
	    }
	}()*/
	m.reader = bufio.NewReader(m.udpConn)

	return nil
}

func (m *TransportInstance) Close() error {
	if m.udpConn == nil {
		return nil
	}
	err := m.udpConn.Close()
	if err != nil {
		return errors.Wrap(err, "error closing connection")
	}
	m.udpConn = nil
	return nil
}

func (m *TransportInstance) IsConnected() bool {
	return m.udpConn != nil
}

func (m *TransportInstance) GetNumBytesAvailableInBuffer() (uint32, error) {
	if m.reader == nil {
		return 0, nil
	}
	_, _ = m.reader.Peek(1)
	return uint32(m.reader.Buffered()), nil
}

func (m *TransportInstance) FillBuffer(until func(pos uint, currentByte byte, reader *bufio.Reader) bool) error {
	nBytes := uint32(1)
	for {
		_bytes, err := m.PeekReadableBytes(nBytes)
		if err != nil {
			return errors.Wrap(err, "Error while peeking")
		}
		if keepGoing := until(uint(nBytes-1), _bytes[len(_bytes)-1], m.reader); !keepGoing {
			return nil
		}
		nBytes++
	}
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

func (m *TransportInstance) Write(data []uint8) error {
	if m.udpConn == nil {
		return errors.New("error writing to transport. No writer available")
	}
	var num int
	var err error
	if m.RemoteAddress != nil {

		// TODO: usually this happens on the dial port... is there a better way to catch that?
		num, err = m.udpConn.Write(data)
	} else {
		num, err = m.udpConn.WriteToUDP(data, m.RemoteAddress)
	}

	if err != nil {
		return errors.Wrap(err, "error writing")
	}
	if num != len(data) {
		return errors.New("error writing: not all bytes written")
	}
	return nil
}

func (m *TransportInstance) String() string {
	return fmt.Sprintf("udp:%s->%s", m.LocalAddress, m.RemoteAddress)
}
