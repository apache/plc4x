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

package bacgopes

import (
	"fmt"

	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/globals"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
)

type PDU interface {
	PCI
	PDUData
	DeepCopy() any
}

// PDUContract provides a set of functions which can be overwritten by a sub struct
type PDUContract interface {
	GetName() string
}

type _PDU struct {
	PDUContract
	*_PCI
	*_PDUData
}

func NewPDU(pdu spi.Message, pduOptions ...PDUOption) PDU {
	p := &_PDU{
		_PCI: newPCI(pdu, nil, nil, nil, false, model.NPDUNetworkPriority_NORMAL_MESSAGE),
	}
	p.PDUContract = p
	for _, option := range pduOptions {
		option(p)
	}
	p._PDUData = NewPDUData(NewArgs(pdu)).(*_PDUData)
	return p
}

type PDUOption func(pdu *_PDU)

func WithPDUUserData(message spi.Message) PDUOption {
	return func(pdu *_PDU) {
		pdu.rootMessage = message
	}
}

func WithPDUSource(pduSource *Address) PDUOption {
	return func(pdu *_PDU) {
		pdu.pduSource = pduSource
	}
}

func WithPDUDestination(pduDestination *Address) PDUOption {
	return func(pdu *_PDU) {
		pdu.pduDestination = pduDestination
	}
}

func WithPDUExpectingReply(expectingReply bool) PDUOption {
	return func(pdu *_PDU) {
		pdu.expectingReply = expectingReply
	}
}

func WithPDUNetworkPriority(networkPriority model.NPDUNetworkPriority) PDUOption {
	return func(pdu *_PDU) {
		pdu.networkPriority = networkPriority
	}
}

func (p *_PDU) GetRootMessage() spi.Message {
	return p.rootMessage
}

func (p *_PDU) deepCopy() *_PDU {
	pduCopy := &_PDU{_PCI: p._PCI.deepCopy(), _PDUData: p._PDUData.deepCopy()}
	pduCopy.PDUContract = pduCopy
	return pduCopy
}

func (p *_PDU) DeepCopy() any {
	return p.deepCopy()
}

func (p *_PDU) GetName() string {
	return "PDU"
}

func (p *_PDU) String() string {
	if globals.ExtendedPDUOutput {
		return fmt.Sprintf("_PDU{%s}", p._PCI)
	}
	return fmt.Sprintf("<%s %s -> %s : %s>", p.PDUContract.GetName(), p.GetPDUSource(), p.GetPDUDestination(), p._PDUData)
}
