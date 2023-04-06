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
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"

	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
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
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return nil, errors.Wrap(err, "error resolving typ address")
	}

	return NewTcpTransportInstance(tcpAddr, connectTimeout, &m), nil
}

func (m Transport) String() string {
	return m.GetTransportCode() + "(" + m.GetTransportName() + ")"
}

type TransportInstance struct {
	transports.DefaultBufferedTransportInstance
	RemoteAddress  *net.TCPAddr
	LocalAddress   *net.TCPAddr
	ConnectTimeout uint32
	transport      *Transport
	tcpConn        net.Conn
	reader         *bufio.Reader
}

func NewTcpTransportInstance(remoteAddress *net.TCPAddr, connectTimeout uint32, transport *Transport) *TransportInstance {
	transportInstance := &TransportInstance{
		RemoteAddress:  remoteAddress,
		ConnectTimeout: connectTimeout,
		transport:      transport,
	}
	transportInstance.DefaultBufferedTransportInstance = transports.NewDefaultBufferedTransportInstance(transportInstance)
	return transportInstance
}

func (m *TransportInstance) Connect() error {
	return m.ConnectWithContext(context.Background())
}

func (m *TransportInstance) ConnectWithContext(ctx context.Context) error {
	var err error
	var d net.Dialer
	m.tcpConn, err = d.DialContext(ctx, "tcp", m.RemoteAddress.String())
	if err != nil {
		return errors.Wrap(err, "error connecting to remote address")
	}

	m.LocalAddress = m.tcpConn.LocalAddr().(*net.TCPAddr)

	m.reader = bufio.NewReaderSize(m.tcpConn, 100000)

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

func (m *TransportInstance) Write(data []byte) error {
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

func (m *TransportInstance) GetReader() *bufio.Reader {
	return m.reader
}

func (m *TransportInstance) String() string {
	return fmt.Sprintf("tcp:%s->%s", m.LocalAddress, m.RemoteAddress)
}
