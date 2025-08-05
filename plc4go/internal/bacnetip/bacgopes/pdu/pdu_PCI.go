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
	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
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
	pduExpectingReply  bool
	pduNetworkPriority model.NPDUNetworkPriority
}

var _ PCI = (*_PCI)(nil)

func NewPCI(args Args, kwArgs KWArgs, options ...Option) *_PCI {
	if _debug != nil {
		_debug("__init__ %r %r", args, kwArgs)
	}
	var myKwargs = make(KWArgs)
	var otherKwargs = make(KWArgs)
	for _, element := range []KnownKey{KWPCIExpectingReply, KWPCINetworkPriority} {
		if v, ok := kwArgs[element]; ok {
			myKwargs[element] = v
		}
	}
	for k, v := range kwArgs {
		if _, ok := myKwargs[k]; !ok {
			otherKwargs[k] = v
		}
	}
	if _debug != nil {
		_debug("    - my_kwargs: %r", myKwargs)
	}
	if _debug != nil {
		_debug("    - other_kwargs: %r", otherKwargs)
	}
	expectingReply, _ := KWO(kwArgs, KWPCIExpectingReply, false)
	networkPriority, _ := KWO(kwArgs, KWPCINetworkPriority, model.NPDUNetworkPriority_NORMAL_MESSAGE)
	i := &_PCI{
		new__PCI(args, kwArgs, options...),
		expectingReply,
		networkPriority,
	}
	i.AddDebugContents(i, "pduExpectingReply", "pduNetworkPriority")
	return i
}

func (p *_PCI) GetDebugAttr(attr string) any {
	switch attr {
	case "pduExpectingReply":
		if p.pduExpectingReply {
			return 1
		} else {
			return 0
		}
	case "pduNetworkPriority":
		return p.pduNetworkPriority
	default:
		return nil
	}
}

func (p *_PCI) Update(pci Arg) error {
	if err := p.__PCI.Update(pci); err != nil {
		return errors.Wrap(err, "error updating __PCI")
	}
	switch pci := pci.(type) {
	case PCI:
		p.pduExpectingReply = pci.GetExpectingReply()
		p.pduNetworkPriority = pci.GetNetworkPriority()
		return nil
	default:
		return errors.Errorf("invalid PCI type %T", pci)
	}
}

func (p *_PCI) SetExpectingReply(expectingReply bool) {
	p.pduExpectingReply = expectingReply
}

func (p *_PCI) GetExpectingReply() bool {
	return p.pduExpectingReply
}

func (p *_PCI) SetNetworkPriority(priority model.NPDUNetworkPriority) {
	p.pduNetworkPriority = priority
}

func (p *_PCI) GetNetworkPriority() model.NPDUNetworkPriority {
	return p.pduNetworkPriority
}

func (p *_PCI) GetPCI() PCI {
	return p
}

func (p *_PCI) deepCopy() *_PCI {
	__pci := p.__PCI.deepCopy()
	expectingReply := p.pduExpectingReply
	networkPriority := p.pduNetworkPriority // Those are immutable so no copy needed
	return &_PCI{__pci, expectingReply, networkPriority}
}

func (p *_PCI) DeepCopy() any {
	return p.deepCopy()
}
