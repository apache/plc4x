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
    "errors"
    "net"
    "net/url"
    "github.com/apache/plc4x/plc4go/internal/plc4go/transports"
    "github.com/apache/plc4x/plc4go/internal/plc4go/utils"
    "regexp"
    "strconv"
)

type TcpTransport struct {
    transports.Transport
}

func NewTcpTransport() *TcpTransport {
    return &TcpTransport{}
}

func (m TcpTransport) GetTransportCode() string {
    return "tcp"
}

func (m TcpTransport) GetTransportName() string {
    return "TCP/IP Socket Transport"
}

func (m TcpTransport) CreateTransportInstance(transportUrl url.URL, options map[string][]string) (transports.TransportInstance, error) {
    connectionStringRegexp := regexp.MustCompile(`^((?P<ip>[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})|(?P<hostname>[a-zA-Z0-9.\-]+))(:(?P<port>[0-9]{1,5}))?`)
    var address string
    var port int
    if match := utils.GetSubgropMatches(connectionStringRegexp, transportUrl.Host); match != nil {
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
                return nil, errors.New("error setting port: " + err.Error())
            } else {
                port = portVal
            }
        } else if val, ok := options["defaultTcpPort"]; ok && len(val) > 0 {
            portVal, err := strconv.Atoi(val[0])
            if err != nil {
                return nil, errors.New("error setting default tcp port: " + err.Error())
            } else {
                port = portVal
            }
        } else {
            return nil, errors.New("error setting port. No explicit or default port provided")
        }
    }
    var connectTimeout uint32 = 1000
    if val, ok := options["connect-timeout"]; ok {
        ival, err := strconv.Atoi(val[0])
        if err != nil {
            return nil, errors.New("error setting connect-timeout: " + err.Error())
        } else {
            connectTimeout = uint32(ival)
        }
    }

    // Potentially resolve the ip address, if a hostname was provided
    tcpAddr, err := net.ResolveTCPAddr("tcp", address + ":" + strconv.Itoa(port))
    if err != nil {
        return nil, errors.New("error resolving typ address: " + err.Error())
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

type TcpTransportInstance struct {
    RemoteAddress *net.TCPAddr
    LocalAddress *net.TCPAddr
    ConnectTimeout uint32
    transport *TcpTransport
    tcpConn net.Conn
    reader *bufio.Reader
}

func NewTcpTransportInstance(remoteAddress *net.TCPAddr, connectTimeout uint32, transport *TcpTransport) *TcpTransportInstance {
    return &TcpTransportInstance {
        RemoteAddress: remoteAddress,
        ConnectTimeout: connectTimeout,
        transport: transport,
    }
}

func (m *TcpTransportInstance) Connect() error {
    var err error
    m.tcpConn, err = net.Dial("tcp", m.RemoteAddress.String())
    if err != nil {
        return errors.New("error connecting to remote address: " + err.Error())
    }

    m.LocalAddress = m.tcpConn.LocalAddr().(*net.TCPAddr)

    m.reader = bufio.NewReader(m.tcpConn)

    return nil
}

func (m *TcpTransportInstance) Close() error {
    if m.tcpConn != nil {
        err := m.tcpConn.Close()
        if err != nil {
            return errors.New("error closing connection: " + err.Error())
        }
    }
    return nil
}

func (m *TcpTransportInstance) GetNumReadableBytes() (uint32, error) {
    if m.reader != nil {
        _, _ = m.reader.Peek(1)
        return uint32(m.reader.Buffered()), nil
    }
    return 0, errors.New("error getting number of available bytes from transport. No reader available")
}

func (m *TcpTransportInstance) PeekReadableBytes(numBytes uint32) ([]uint8, error) {
    if m.reader != nil {
        return m.reader.Peek(int(numBytes))
    }
    return nil, errors.New("error peeking from transport. No reader available")
}

func (m *TcpTransportInstance) Read(numBytes uint32) ([]uint8, error) {
    if m.reader != nil {
        data := make([]uint8, numBytes)
        for i := uint32(0); i < numBytes; i++ {
            val, err := m.reader.ReadByte()
            if err != nil {
                return nil, errors.New("error reading: " + err.Error())
            }
            data[i] = val
        }
        return data, nil
    }
    return nil, errors.New("error reading from transport. No reader available")
}

func (m *TcpTransportInstance) Write(data []uint8) error {
    if m.tcpConn != nil {
        num, err := m.tcpConn.Write(data)
        if err != nil {
            return errors.New("error writing: " + err.Error())
        }
        if num != len(data) {
            return errors.New("error writing: not all bytes written")
        }
        return nil
    }
    return errors.New("error writing to transport. No writer available")
}
