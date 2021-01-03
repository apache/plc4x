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
package udp

import (
	"bufio"
    "errors"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports"
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
    "net"
    "net/url"
    "regexp"
    "strconv"
)

type UdpTransport struct {
	transports.Transport
}

func NewUdpTransport() *UdpTransport {
    return &UdpTransport{}
}

func (m UdpTransport) GetTransportCode() string {
    return "udp"
}

func (m UdpTransport) GetTransportName() string {
    return "UDP Datagram Transport"
}

func (m UdpTransport) CreateTransportInstance(transportUrl url.URL, options map[string][]string) (transports.TransportInstance, error) {
    return m.CreateTransportInstanceForLocalAddress(transportUrl, options, nil)
}

func (m UdpTransport) CreateTransportInstanceForLocalAddress(transportUrl url.URL, options map[string][]string, localAddress *net.UDPAddr) (transports.TransportInstance, error) {
    connectionStringRegexp := regexp.MustCompile(`^((?P<ip>[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})|(?P<hostname>[a-zA-Z0-9.\-]+))(:(?P<port>[0-9]{1,5}))?`)
    var remoteAddressString string
    var remotePort int
    if match := utils.GetSubgropMatches(connectionStringRegexp, transportUrl.Host); match != nil {
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
                return nil, errors.New("error setting port: " + err.Error())
            } else {
                remotePort = portVal
            }
        } else if val, ok := options["defaultUdpPort"]; ok && len(val) > 0 {
            portVal, err := strconv.Atoi(val[0])
            if err != nil {
                return nil, errors.New("error setting default udp port: " + err.Error())
            } else {
                remotePort = portVal
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
    remoteAddress, err := net.ResolveUDPAddr("udp", remoteAddressString+":"+strconv.Itoa(remotePort))
    if err != nil {
        return nil, errors.New("error resolving typ address: " + err.Error())
    }

    transportInstance := NewUdpTransportInstance(localAddress, remoteAddress, connectTimeout, &m)

    castFunc := func(typ interface{}) (transports.TransportInstance, error) {
        if transportInstance, ok := typ.(transports.TransportInstance); ok {
            return transportInstance, nil
        }
        return nil, errors.New("couldn't cast to TransportInstance")
    }
    return castFunc(transportInstance)
}

type UdpTransportInstance struct {
    LocalAddress   *net.UDPAddr
    RemoteAddress  *net.UDPAddr
    ConnectTimeout uint32
    transport      *UdpTransport
    udpConn        *net.UDPConn
    reader         *bufio.Reader
}

func NewUdpTransportInstance(localAddress *net.UDPAddr, remoteAddress *net.UDPAddr, connectTimeout uint32, transport *UdpTransport) *UdpTransportInstance {
    return &UdpTransportInstance{
        LocalAddress:   localAddress,
        RemoteAddress:  remoteAddress,
        ConnectTimeout: connectTimeout,
        transport:      transport,
    }
}

func (m *UdpTransportInstance) Connect() error {
    // If we haven't provided a local address, have the system figure it out by dialing
    // the remote address and then using that connections local address as local address.
    if m.LocalAddress == nil {
        udpTest, err := net.Dial("udp", m.RemoteAddress.String())
        if err != nil {
            return errors.New("error connecting to remote address: " + err.Error())
        }
        m.LocalAddress = udpTest.LocalAddr().(*net.UDPAddr)
        err = udpTest.Close()
        if err != nil {
            return errors.New("error closing test-connection: " + err.Error())
        }
    }

    // "connect" to the remote
    var err error
    m.udpConn, err = net.ListenUDP("udp", m.LocalAddress)
    if err != nil {
        return errors.New("error connecting to remote address: " + err.Error())
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

func (m *UdpTransportInstance) Close() error {
    if m.udpConn != nil {
        err := m.udpConn.Close()
        if err != nil {
            return errors.New("error closing connection: " + err.Error())
        }
    }
    return nil
}

func (m *UdpTransportInstance) GetNumReadableBytes() (uint32, error) {
    if m.reader != nil {
        _, _ = m.reader.Peek(1)
        return uint32(m.reader.Buffered()), nil
    }
    return 0, errors.New("error getting number of available bytes from transport. No reader available")
}

func (m *UdpTransportInstance) PeekReadableBytes(numBytes uint32) ([]uint8, error) {
    if m.reader != nil {
        return m.reader.Peek(int(numBytes))
    }
    return nil, errors.New("error peeking from transport. No reader available")
}

func (m *UdpTransportInstance) Read(numBytes uint32) ([]uint8, error) {
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

func (m *UdpTransportInstance) Write(data []uint8) error {
    if m.udpConn != nil {
        num, err := m.udpConn.WriteToUDP(data, m.RemoteAddress)
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
