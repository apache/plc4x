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

package pdu

import (
	"fmt"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/globals"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
)

type PCI interface {
	IPCI

	SetExpectingReply(bool)
	GetExpectingReply() bool
	SetNetworkPriority(model.NPDUNetworkPriority)
	GetNetworkPriority() model.NPDUNetworkPriority

	GetPCI() PCI
}

type _PCI struct {
	*__PCI
	expectingReply  bool
	networkPriority model.NPDUNetworkPriority
}

var _ PCI = (*_PCI)(nil)

func NewPCI(rootMessage spi.Message, pduUserData spi.Message, pduSource *Address, pduDestination *Address, expectingReply bool, networkPriority model.NPDUNetworkPriority) *_PCI {
	return &_PCI{
		new__PCI(rootMessage, pduUserData, pduSource, pduDestination),
		expectingReply,
		networkPriority,
	}
}

func (p *_PCI) Update(pci Arg) error {
	if err := p.__PCI.Update(pci); err != nil {
		return errors.Wrap(err, "error updating __PCI")
	}
	switch pci := pci.(type) {
	case PCI:
		p.expectingReply = pci.GetExpectingReply()
		p.networkPriority = pci.GetNetworkPriority()
		return nil
	default:
		return errors.Errorf("invalid PCI type %T", pci)
	}
}

func (p *_PCI) SetExpectingReply(expectingReply bool) {
	p.expectingReply = expectingReply
}

func (p *_PCI) GetExpectingReply() bool {
	return p.expectingReply
}

func (p *_PCI) SetNetworkPriority(priority model.NPDUNetworkPriority) {
	p.networkPriority = priority
}

func (p *_PCI) GetNetworkPriority() model.NPDUNetworkPriority {
	return p.networkPriority
}

func (p *_PCI) GetPCI() PCI {
	return p
}

func (p *_PCI) deepCopy() *_PCI {
	__pci := p.__PCI.deepCopy()
	expectingReply := p.expectingReply
	networkPriority := p.networkPriority // Those are immutable so no copy needed
	return &_PCI{__pci, expectingReply, networkPriority}
}

func (p *_PCI) DeepCopy() any {
	return p.deepCopy()
}

func (p *_PCI) String() string {
	if ExtendedPDUOutput {
		return fmt.Sprintf("_PCI{%s, expectingReply: %t, networkPriority: %s}", p.__PCI, p.expectingReply, p.networkPriority)
	} else {
		return fmt.Sprintf("%s\npduExpectingReply = %t\npduNetworkPriority = %s", p.__PCI, p.expectingReply, p.networkPriority)
	}
}
