//go:build linux

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

package rawsocket

// Runs over a real interface pair when RAWSOCKET_TEST_PAIR names one (needs CAP_NET_RAW):
//
//	ip link add veth0 type veth peer name veth1 && ip link set up veth0 && ip link set up veth1
//	RAWSOCKET_TEST_PAIR=veth0,veth1 go test ./spi/transports/rawsocket/

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func open(t *testing.T, name string, filter Filter) *TransportInstance {
	t.Helper()
	instance := NewTransportInstance(name, filter, NewTransport())
	require.NoError(t, instance.Connect(t.Context()))
	t.Cleanup(func() { _ = instance.Close() })
	return instance
}

func receive(t *testing.T, instance *TransportInstance, timeout time.Duration) ([]byte, error) {
	ctx, cancel := context.WithTimeout(t.Context(), timeout)
	defer cancel()
	header, err := instance.PeekReadableBytes(ctx, RecordHeaderSize)
	if err != nil {
		return nil, err
	}
	length, err := RecordFrameLength(header)
	if err != nil {
		return nil, err
	}
	record, err := instance.Read(ctx, uint32(RecordHeaderSize+length))
	if err != nil {
		return nil, err
	}
	return record[RecordHeaderSize:], nil
}

func TestRawSocketPair(t *testing.T) {
	pair := strings.Split(os.Getenv("RAWSOCKET_TEST_PAIR"), ",")
	if len(pair) != 2 {
		t.Skip("set RAWSOCKET_TEST_PAIR to two connected interfaces (a veth pair) to run")
	}
	sender := open(t, pair[0], Filter{})
	receiver := open(t, pair[1], Filter{EtherTypes: []uint16{0x8892}, UDPPorts: []uint16{34964}})
	require.Len(t, sender.HardwareAddress(), 6)

	profinet := ethernet(0x8892, append([]byte{0xFE, 0xFE, 0x05, 0x00}, make([]byte, 42)...))
	arp := ethernet(0x0806, make([]byte, 46))
	udp := ethernet(0x0800, append(ipv4(17, 5, 0, 49152, 34964), make([]byte, 20)...))
	var records []byte
	for _, frame := range [][]byte{arp, profinet, udp} {
		record, err := EncodeRecord(frame)
		require.NoError(t, err)
		records = append(records, record...)
	}
	require.NoError(t, sender.Write(t.Context(), records))

	got, err := receive(t, receiver, 2*time.Second)
	require.NoError(t, err)
	assert.Equal(t, profinet, got, "ARP is filtered out, Profinet arrives first")
	got, err = receive(t, receiver, 2*time.Second)
	require.NoError(t, err)
	assert.Equal(t, udp, got)

	_, err = receive(t, receiver, 100*time.Millisecond)
	assert.Error(t, err, "nothing else arrives")
	assert.Equal(t, 0, sender.reader.Buffered(), "the sender does not hear its own frames")
}
