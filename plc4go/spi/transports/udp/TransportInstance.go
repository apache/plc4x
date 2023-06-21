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
	"net"
	"sync"
	"sync/atomic"

	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"

	"github.com/libp2p/go-reuseport"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type TransportInstance struct {
	LocalAddress   *net.UDPAddr
	RemoteAddress  *net.UDPAddr
	ConnectTimeout uint32
	SoReUse        bool

	transport *Transport
	udpConn   *net.UDPConn
	reader    *bufio.Reader

	connected        atomic.Bool
	stateChangeMutex sync.RWMutex

	log zerolog.Logger
}

func NewTransportInstance(localAddress *net.UDPAddr, remoteAddress *net.UDPAddr, connectTimeout uint32, soReUse bool, transport *Transport, _options ...options.WithOption) *TransportInstance {
	logger, _ := options.ExtractCustomLogger(_options...)
	return &TransportInstance{
		LocalAddress:   localAddress,
		RemoteAddress:  remoteAddress,
		ConnectTimeout: connectTimeout,
		SoReUse:        soReUse,
		transport:      transport,

		log: logger,
	}
}

func (m *TransportInstance) Connect() error {
	return m.ConnectWithContext(context.Background())
}

func (m *TransportInstance) ConnectWithContext(ctx context.Context) error {
	if m.connected.Load() {
		return errors.New("already connected")
	}
	m.stateChangeMutex.Lock()
	defer m.stateChangeMutex.Unlock()
	// If we haven't provided a local address, have the system figure it out by dialing
	// the remote address and then using that connections local address as local address.
	if m.LocalAddress == nil && m.RemoteAddress != nil {
		var d net.Dialer
		udpTest, err := d.DialContext(ctx, "udp", m.RemoteAddress.String())
		if err != nil {
			return errors.Wrapf(err, "error connecting to remote address '%s'", m.RemoteAddress)
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
			return errors.Wrapf(err, "error connecting to remote address '%s'", m.RemoteAddress)
		}
	} else if m.SoReUse && m.LocalAddress != nil {
		if packetConn, err := reuseport.ListenPacket("udp", m.LocalAddress.String()); err != nil {
			return errors.Wrapf(err, "error connecting to local address '%s'", m.LocalAddress)
		} else {
			m.udpConn = packetConn.(*net.UDPConn)
		}
	} else {
		if m.udpConn, err = net.ListenUDP("udp", m.LocalAddress); err != nil {
			return errors.Wrapf(err, "error connecting to local address '%s'", m.LocalAddress)
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

	m.connected.Store(true)

	return nil
}

func (m *TransportInstance) Close() error {
	m.stateChangeMutex.Lock()
	defer m.stateChangeMutex.Unlock()
	if !m.connected.Load() {
		return nil
	}
	err := m.udpConn.Close()
	if err != nil {
		return errors.Wrap(err, "error closing connection")
	}
	m.connected.Store(false)
	return nil
}

func (m *TransportInstance) IsConnected() bool {
	return m.connected.Load()
}

func (m *TransportInstance) GetNumBytesAvailableInBuffer() (uint32, error) {
	if !m.IsConnected() {
		return 0, errors.New("working on a unconnected connection")
	}
	if m.reader == nil {
		return 0, nil
	}
	_, _ = m.reader.Peek(1)
	return uint32(m.reader.Buffered()), nil
}

func (m *TransportInstance) FillBuffer(until func(pos uint, currentByte byte, reader transports.ExtendedReader) bool) error {
	if !m.IsConnected() {
		return errors.New("working on a unconnected connection")
	}
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

func (m *TransportInstance) PeekReadableBytes(numBytes uint32) ([]byte, error) {
	if !m.IsConnected() {
		return nil, errors.New("working on a unconnected connection")
	}
	return m.reader.Peek(int(numBytes))
}

func (m *TransportInstance) Read(numBytes uint32) ([]byte, error) {
	if !m.IsConnected() {
		return nil, errors.New("working on a unconnected connection")
	}
	data := make([]byte, numBytes)
	for i := uint32(0); i < numBytes; i++ {
		val, err := m.reader.ReadByte()
		if err != nil {
			return nil, errors.Wrap(err, "error reading")
		}
		data[i] = val
	}
	return data, nil
}

func (m *TransportInstance) Write(data []byte) error {
	if !m.IsConnected() {
		return errors.New("working on a unconnected connection")
	}
	var num int
	var err error
	if m.RemoteAddress == nil {
		// TODO: usually this happens on the dial port... is there a better way to catch that?
		num, err = m.udpConn.Write(data)
	} else {
		num, err = m.udpConn.WriteToUDP(data, m.RemoteAddress)
	}
	if err != nil {
		return errors.Wrapf(err, "error writing (remote address: %s)", m.RemoteAddress)
	}
	if num != len(data) {
		return errors.Errorf("error writing: not all bytes written (Expected %d, Actual %d)", len(data), num)
	}
	return nil
}

func (m *TransportInstance) String() string {
	return fmt.Sprintf("udp:%s->%s", m.LocalAddress, m.RemoteAddress)
}
