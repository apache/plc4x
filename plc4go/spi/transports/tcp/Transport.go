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

	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	transportUtils "github.com/apache/plc4x/plc4go/spi/transports/utils"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Transport struct {
	log zerolog.Logger
}

func NewTransport(_options ...options.WithOption) *Transport {
	return &Transport{
		log: options.ExtractCustomLogger(_options...),
	}
}

func (m Transport) GetTransportCode() string {
	return "tcp"
}

func (m Transport) GetTransportName() string {
	return "TCP/IP Socket Transport"
}

func (m Transport) CreateTransportInstance(transportUrl url.URL, options map[string][]string, _options ...options.WithOption) (transports.TransportInstance, error) {
	connectionStringRegexp := regexp.MustCompile(`^((?P<ip>[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})|(?P<hostname>[a-zA-Z0-9.\-]+))(:(?P<port>[0-9]{1,5}))?`)
	var address string
	var port int
	if match := utils.GetSubgroupMatches(connectionStringRegexp, transportUrl.Host); match != nil {
		if val, ok := match["ip"]; ok && len(val) > 0 {
			address = val
		} else if val, ok := match["hostname"]; ok && len(val) > 0 {
			address = val
		} // Note: the regex ensures that it is either ip or hostname
		if val, ok := match["port"]; ok && len(val) > 0 {
			port, _ = strconv.Atoi(val) // Note: the regex ensures that this is always a valid number
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

	return NewTcpTransportInstance(tcpAddr, connectTimeout, &m, _options...), nil
}

func (m Transport) String() string {
	return m.GetTransportCode() + "(" + m.GetTransportName() + ")"
}

type TransportInstance struct {
	transportUtils.DefaultBufferedTransportInstance
	RemoteAddress  *net.TCPAddr
	LocalAddress   *net.TCPAddr
	ConnectTimeout uint32
	transport      *Transport
	tcpConn        net.Conn
	reader         *bufio.Reader

	log zerolog.Logger
}

func NewTcpTransportInstance(remoteAddress *net.TCPAddr, connectTimeout uint32, transport *Transport, _options ...options.WithOption) *TransportInstance {
	transportInstance := &TransportInstance{
		RemoteAddress:  remoteAddress,
		ConnectTimeout: connectTimeout,
		transport:      transport,

		log: options.ExtractCustomLogger(_options...),
	}
	transportInstance.DefaultBufferedTransportInstance = transportUtils.NewDefaultBufferedTransportInstance(transportInstance, _options...)
	return transportInstance
}

func (m *TransportInstance) Connect() error {
	return m.ConnectWithContext(context.Background())
}

func (m *TransportInstance) ConnectWithContext(ctx context.Context) error {
	if m.RemoteAddress == nil {
		return errors.New("Required remote address missing")
	}
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
	localAddress := ""
	if m.LocalAddress != nil {
		localAddress = m.LocalAddress.String() + "->"
	}
	return fmt.Sprintf("tcp:%s%s", localAddress, m.RemoteAddress)
}
