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

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// run interprets the instructions this package emits, the way the kernel would. Loads past the end
// of the frame reject it, as in the kernel.
func run(t *testing.T, program []Instruction, frame []byte) bool {
	t.Helper()
	var a, x uint32
	for pc := 0; pc < len(program); pc++ {
		in := program[pc]
		load := func(offset uint32, size uint32) (uint32, bool) {
			if offset+size > uint32(len(frame)) {
				return 0, false
			}
			if size == 1 {
				return uint32(frame[offset]), true
			}
			return uint32(binary.BigEndian.Uint16(frame[offset:])), true
		}
		var ok bool
		switch in.Code {
		case bpfLdHAbs:
			if a, ok = load(in.K, 2); !ok {
				return false
			}
		case bpfLdBAbs:
			if a, ok = load(in.K, 1); !ok {
				return false
			}
		case bpfLdHInd:
			if a, ok = load(x+in.K, 2); !ok {
				return false
			}
		case bpfLdxBMsh:
			var b uint32
			if b, ok = load(in.K, 1); !ok {
				return false
			}
			x = 4 * (b & 0xF)
		case bpfJeqK:
			if a == in.K {
				pc += int(in.Jt)
			} else {
				pc += int(in.Jf)
			}
		case bpfJsetK:
			if a&in.K != 0 {
				pc += int(in.Jt)
			} else {
				pc += int(in.Jf)
			}
		case bpfRetK:
			return in.K != 0
		default:
			t.Fatalf("unexpected opcode 0x%02x", in.Code)
		}
	}
	t.Fatal("program fell off its end")
	return false
}

var (
	destination = []byte{0x01, 0x0E, 0xCF, 0x00, 0x00, 0x00}
	source      = []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x01}
)

func ethernet(etherType uint16, payload []byte) []byte {
	frame := append(append([]byte{}, destination...), source...)
	frame = binary.BigEndian.AppendUint16(frame, etherType)
	return append(frame, payload...)
}

func tagged(etherType uint16, payload []byte) []byte {
	frame := append(append([]byte{}, destination...), source...)
	frame = binary.BigEndian.AppendUint16(frame, etherTypeVLAN)
	frame = binary.BigEndian.AppendUint16(frame, 0x0000) // priority 0, VLAN 0
	frame = binary.BigEndian.AppendUint16(frame, etherType)
	return append(frame, payload...)
}

// ipv4 builds an IPv4 header (with options when ihl > 5) followed by a UDP or other header.
func ipv4(protocol byte, ihl byte, fragment uint16, sourcePort, destinationPort uint16) []byte {
	header := make([]byte, int(ihl)*4)
	header[0] = 0x40 | ihl
	binary.BigEndian.PutUint16(header[6:], fragment)
	header[9] = protocol
	transport := make([]byte, 8)
	binary.BigEndian.PutUint16(transport[0:], sourcePort)
	binary.BigEndian.PutUint16(transport[2:], destinationPort)
	return append(header, transport...)
}

func TestFilterProgram(t *testing.T) {
	filter := Filter{EtherTypes: []uint16{0x8892, 0x88CC}, UDPPorts: []uint16{34964, 50000}}
	program, err := filter.Program()
	require.NoError(t, err)
	for name, c := range map[string]struct {
		frame []byte
		want  bool
	}{
		"profinet":                   {ethernet(0x8892, make([]byte, 40)), true},
		"lldp":                       {ethernet(0x88CC, make([]byte, 40)), true},
		"tagged profinet":            {tagged(0x8892, make([]byte, 40)), true},
		"arp":                        {ethernet(0x0806, make([]byte, 28)), false},
		"udp to 34964":               {ethernet(0x0800, ipv4(17, 5, 0, 49152, 34964)), true},
		"udp from 34964":             {ethernet(0x0800, ipv4(17, 5, 0, 34964, 49152)), true},
		"udp to 50000, ip options":   {ethernet(0x0800, ipv4(17, 6, 0, 1234, 50000)), true},
		"tagged udp to 34964":        {tagged(0x0800, ipv4(17, 5, 0, 49152, 34964)), true},
		"udp to another port":        {ethernet(0x0800, ipv4(17, 5, 0, 1234, 53)), false},
		"tcp to 34964":               {ethernet(0x0800, ipv4(6, 5, 0, 1234, 34964)), false},
		"udp fragment":               {ethernet(0x0800, ipv4(17, 5, 0x0010, 1234, 34964)), false},
		"udp first fragment":         {ethernet(0x0800, ipv4(17, 5, 0x2000, 1234, 34964)), true},
		"truncated ipv4":             {ethernet(0x0800, []byte{0x45}), false},
		"ipv6":                       {ethernet(0x86DD, make([]byte, 48)), false},
		"tagged with another inside": {tagged(0x0806, make([]byte, 28)), false},
	} {
		assert.Equal(t, c.want, run(t, program, c.frame), name)
		assert.Equal(t, c.want, filter.Matches(c.frame), "Matches: "+name)
	}
}

func TestFilterEtherTypesOnly(t *testing.T) {
	program, err := Filter{EtherTypes: []uint16{0x8892}}.Program()
	require.NoError(t, err)
	assert.True(t, run(t, program, ethernet(0x8892, make([]byte, 40))))
	assert.False(t, run(t, program, ethernet(0x0800, ipv4(17, 5, 0, 1, 34964))))
}

func TestEmptyFilterAcceptsAll(t *testing.T) {
	program, err := Filter{}.Program()
	require.NoError(t, err)
	assert.True(t, run(t, program, ethernet(0x0806, make([]byte, 28))))
}

func TestParseFilter(t *testing.T) {
	filter, err := ParseFilter(map[string][]string{
		EtherTypesOption: {"0x8892, 0x88cc"},
		UDPPortsOption:   {"34964,0xC350"},
	})
	require.NoError(t, err)
	assert.Equal(t, Filter{EtherTypes: []uint16{0x8892, 0x88CC}, UDPPorts: []uint16{34964, 50000}}, filter)
	_, err = ParseFilter(map[string][]string{UDPPortsOption: {"70000"}})
	assert.Error(t, err)
}

func TestRecords(t *testing.T) {
	frame := ethernet(0x8892, make([]byte, 46))
	record, err := EncodeRecord(frame)
	require.NoError(t, err)
	length, err := RecordFrameLength(record)
	require.NoError(t, err)
	assert.Equal(t, len(frame), length)
	assert.True(t, bytes.Equal(frame, record[RecordHeaderSize:]))
	_, err = EncodeRecord(make([]byte, 10))
	assert.Error(t, err)
	_, err = EncodeRecord(make([]byte, MaxFrameSize+1))
	assert.Error(t, err)
	_, err = RecordFrameLength([]byte{0xFF, 0xFF})
	assert.Error(t, err)
}
