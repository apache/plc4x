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
	"encoding/binary"
	"slices"

	"github.com/apache/plc4x/plc4go/spi/errors"
)

// Classic BPF opcodes this package uses (linux/filter.h).
const (
	bpfLdHAbs  = 0x28 // A = P[k:2]
	bpfLdBAbs  = 0x30 // A = P[k:1]
	bpfLdHInd  = 0x48 // A = P[X+k:2]
	bpfLdxBMsh = 0xb1 // X = 4*(P[k]&0xf)
	bpfJeqK    = 0x15 // pc += (A == k) ? jt : jf
	bpfJsetK   = 0x45 // pc += (A & k) ? jt : jf
	bpfRetK    = 0x06 // return k
)

const (
	etherTypeVLAN = 0x8100
	etherTypeIPv4 = 0x0800
	ipProtocolUDP = 17
	// acceptLength is how much of an accepted frame the filter keeps: all of it.
	acceptLength = 0x40000
)

// Instruction is one classic BPF instruction, the layout of struct sock_filter.
type Instruction struct {
	Code uint16
	Jt   uint8
	Jf   uint8
	K    uint32
}

// Filter says which frames an instance receives. The kernel applies it, so frames nobody asked for
// never wake the process - on a plant network most traffic is somebody else's.
type Filter struct {
	// EtherTypes are accepted whole, tagged with a VLAN or not.
	EtherTypes []uint16
	// UDPPorts accepts unfragmented IPv4 UDP datagrams whose source or destination port is listed.
	UDPPorts []uint16
}

// IsEmpty is a filter that accepts everything.
func (f Filter) IsEmpty() bool {
	return len(f.EtherTypes) == 0 && len(f.UDPPorts) == 0
}

// assembly is a BPF program with symbolic jump targets.
type assembly struct {
	instructions []Instruction
	jumps        []symbolicJump
	labels       map[string]int
}

type symbolicJump struct {
	at     int
	jt, jf string
}

func (a *assembly) emit(code uint16, k uint32) {
	a.instructions = append(a.instructions, Instruction{Code: code, K: k})
}

// jump emits a conditional jump; an empty label means "the next instruction".
func (a *assembly) jump(code uint16, k uint32, jt, jf string) {
	a.jumps = append(a.jumps, symbolicJump{at: len(a.instructions), jt: jt, jf: jf})
	a.emit(code, k)
}

func (a *assembly) label(name string) {
	a.labels[name] = len(a.instructions)
}

func (a *assembly) resolve() ([]Instruction, error) {
	offset := func(at int, label string) (uint8, error) {
		if label == "" {
			return 0, nil
		}
		target, ok := a.labels[label]
		if !ok {
			return 0, errors.Errorf("undefined label %s", label)
		}
		distance := target - at - 1
		if distance < 0 || distance > 255 {
			return 0, errors.Errorf("jump to %s is %d instructions, classic BPF jumps forward at most 255", label, distance)
		}
		return uint8(distance), nil
	}
	for _, j := range a.jumps {
		var err error
		if a.instructions[j.at].Jt, err = offset(j.at, j.jt); err != nil {
			return nil, err
		}
		if a.instructions[j.at].Jf, err = offset(j.at, j.jf); err != nil {
			return nil, err
		}
	}
	return a.instructions, nil
}

// Program compiles the filter. Frames carrying an 802.1Q tag in the data are handled too, for
// interfaces without VLAN offload; with offload the kernel strips the tag before the filter runs.
func (f Filter) Program() ([]Instruction, error) {
	if f.IsEmpty() {
		return []Instruction{{Code: bpfRetK, K: acceptLength}}, nil
	}
	a := &assembly{labels: map[string]int{}}
	a.emit(bpfLdHAbs, 12)
	a.jump(bpfJeqK, etherTypeVLAN, "tagged", "")
	f.block(a, 0, "untagged")
	a.label("tagged")
	a.emit(bpfLdHAbs, 16)
	f.block(a, 4, "tagged")
	a.label("accept")
	a.emit(bpfRetK, acceptLength)
	a.label("reject")
	a.emit(bpfRetK, 0)
	return a.resolve()
}

// block matches with the EtherType in A and the IPv4 header at 14+base.
func (f Filter) block(a *assembly, base uint32, name string) {
	for _, etherType := range f.EtherTypes {
		a.jump(bpfJeqK, uint32(etherType), "accept", "")
	}
	if len(f.UDPPorts) == 0 {
		a.jump(bpfJeqK, 0xFFFFFFFF, "reject", "reject")
		return
	}
	a.jump(bpfJeqK, etherTypeIPv4, "", "reject")
	a.emit(bpfLdBAbs, base+14+9) // IPv4 protocol
	a.jump(bpfJeqK, ipProtocolUDP, "", "reject")
	a.emit(bpfLdHAbs, base+14+6) // flags and fragment offset
	a.jump(bpfJsetK, 0x1FFF, "reject", "")
	a.emit(bpfLdxBMsh, base+14) // X = IPv4 header length
	a.emit(bpfLdHInd, base+14)  // UDP source port
	for _, port := range f.UDPPorts {
		a.jump(bpfJeqK, uint32(port), "accept", "")
	}
	a.emit(bpfLdHInd, base+14+2) // UDP destination port
	for _, port := range f.UDPPorts {
		a.jump(bpfJeqK, uint32(port), "accept", "")
	}
	a.jump(bpfJeqK, 0xFFFFFFFF, "reject", "reject")
}

// Matches applies the filter in Go, for transports without a kernel to do it (the virtual bus). It
// accepts exactly what Program accepts.
func (f Filter) Matches(frame []byte) bool {
	if f.IsEmpty() {
		return true
	}
	if len(frame) < MinFrameSize {
		return false
	}
	base := 0
	etherType := binary.BigEndian.Uint16(frame[12:])
	if etherType == etherTypeVLAN {
		if len(frame) < 18 {
			return false
		}
		base = 4
		etherType = binary.BigEndian.Uint16(frame[16:])
	}
	if slices.Contains(f.EtherTypes, etherType) {
		return true
	}
	if len(f.UDPPorts) == 0 || etherType != etherTypeIPv4 {
		return false
	}
	ip := frame[14+base:]
	if len(ip) < 20 || ip[9] != ipProtocolUDP || binary.BigEndian.Uint16(ip[6:])&0x1FFF != 0 {
		return false
	}
	headerLength := int(ip[0]&0xF) * 4
	if len(ip) < headerLength+4 {
		return false
	}
	return slices.Contains(f.UDPPorts, binary.BigEndian.Uint16(ip[headerLength:])) ||
		slices.Contains(f.UDPPorts, binary.BigEndian.Uint16(ip[headerLength+2:]))
}
