//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package tcp

import (
	"bufio"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/pkg/errors"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

type Transport struct {
}

func NewTransport() *Transport {
	return &Transport{}
}

func (m Transport) GetTransportCode() string {
	return "tcp"
}

func (m Transport) GetTransportName() string {
	return "TCP/IP Socket Transport"
}

func (m Transport) CreateTransportInstance(transportUrl url.URL, options map[string][]string) (transports.TransportInstance, error) {
	connectionStringRegexp := regexp.MustCompile(`^((?P<ip>[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})|(?P<hostname>[a-zA-Z0-9.\-]+))(:(?P<port>[0-9]{1,5}))?`)
	var address string
	var port int
	if match := utils.GetSubgroupMatches(connectionStringRegexp, transportUrl.Host); match != nil {
		if val, ok := match["ip"]; ok && len(val) > 0 {
			address = val
		} else if val, ok := match["hostname"]; ok && len(val) > 0 {
			address = val
		} else {
			return nil, errors.New("missing hostname or ip to connect")
		}
		if val, ok := match["port"]; ok && len(val) > 0 {
			portVal, err := strconv.Atoi(val)
			if err != nil {
				return nil, errors.Wrap(err, "error setting port")
			} else {
				port = portVal
			}
		} else if val, ok := options["defaultTcpPort"]; ok && len(val) > 0 {
			portVal, err := strconv.Atoi(val[0])
			if err != nil {
				return nil, errors.Wrap(err, "error setting default tcp port")
			}
			port = portVal
		} else {
			return nil, errors.New("error setting port. No explicit or default port provided")
		}
	}
	var connectTimeout uint32 = 1000
	if val, ok := options["connect-timeout"]; ok {
		parsedConnectTimeout, err := strconv.ParseUint(val[0], 10, 32)
		if err != nil {
			return nil, errors.Wrap(err, "error setting connect-timeout")
		}
		connectTimeout = uint32(parsedConnectTimeout)
	}

	// Potentially resolve the ip address, if a hostname was provided
	tcpAddr, err := net.ResolveTCPAddr("tcp", address+":"+strconv.Itoa(port))
	if err != nil {
		return nil, errors.Wrap(err, "error resolving typ address")
	}

	transportInstance := NewTcpTransportInstance(tcpAddr, connectTimeout, &m)

	castFunc := func(typ interface{}) (transports.TransportInstance, error) {
		if transportInstance, ok := typ.(transports.TransportInstance); ok {
			return transportInstance, nil
		}
		return nil, errors.New("couldn't cast to TransportInstance")
	}
	return castFunc(transportInstance)
}

type TransportInstance struct {
	RemoteAddress  *net.TCPAddr
	LocalAddress   *net.TCPAddr
	ConnectTimeout uint32
	transport      *Transport
	tcpConn        net.Conn
	reader         *bufio.Reader
}

func NewTcpTransportInstance(remoteAddress *net.TCPAddr, connectTimeout uint32, transport *Transport) *TransportInstance {
	return &TransportInstance{
		RemoteAddress:  remoteAddress,
		ConnectTimeout: connectTimeout,
		transport:      transport,
	}
}

func (m *TransportInstance) Connect() error {
	var err error
	m.tcpConn, err = net.Dial("tcp", m.RemoteAddress.String())
	if err != nil {
		return errors.Wrap(err, "error connecting to remote address")
	}

	m.LocalAddress = m.tcpConn.LocalAddr().(*net.TCPAddr)

	m.reader = bufio.NewReader(m.tcpConn)

	return nil
}

func (m *TransportInstance) Close() error {
	if m.tcpConn == nil {
		return nil
	}
	err := m.tcpConn.Close()
	if err != nil {
		return errors.Wrap(err, "error closing connection")
	}
	m.tcpConn = nil
	return nil
}

func (m *TransportInstance) IsConnected() bool {
	return m.tcpConn != nil
}

func (m *TransportInstance) GetNumReadableBytes() (uint32, error) {
	if m.reader == nil {
		return 0, nil
	}
	peekChan := make(chan bool)
	go func() {
		_, _ = m.reader.Peek(1)
		peekChan <- true
	}()
	select {
	case <-peekChan:
		return uint32(m.reader.Buffered()), nil
	case <-time.After(10 * time.Millisecond):
		return 0, nil
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
	if m.tcpConn == nil {
		return errors.New("error writing to transport. No writer available")
	}
	num, err := m.tcpConn.Write(data)
	if err != nil {
		return errors.Wrap(err, "error writing")
	}
	if num != len(data) {
		return errors.New("error writing: not all bytes written")
	}
	return nil
}
