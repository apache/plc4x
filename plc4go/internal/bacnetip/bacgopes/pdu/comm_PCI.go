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
	"context"
	"fmt"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Note: upstream this belongs to comm but that would create a circular dependency

type IPCI interface {
	spi.Message
	GetLeafName() string
	SetRootMessage(spi.Message)
	GetRootMessage() spi.Message
	SetPDUUserData(spi.Message)
	GetPDUUserData() spi.Message
	GetPDUSource() *Address
	SetPDUSource(*Address)
	GetPDUDestination() *Address
	SetPDUDestination(*Address)
	Update(pci Arg) error
	DeepCopy() any
}

type __PCI struct {
	*DebugContents
	rootMessage    spi.Message
	pduUserData    spi.Message
	pduSource      *Address
	pduDestination *Address

	_leafName string
}

var _ IPCI = (*__PCI)(nil)

func new__PCI(args Args, kwArgs KWArgs, options ...Option) *__PCI {
	if _debug != nil {
		_debug("__init__ %r %r", args, kwArgs)
	}
	i := &__PCI{
		rootMessage: ExtractRootMessage(options),
		_leafName:   ExtractLeafName(options, "PCI"),
	}
	i.DebugContents = NewDebugContents(i, "pduUserData+", "pduSource", "pduDestination")
	var myKwargs = make(KWArgs)
	var otherKwargs = make(KWArgs)
	for _, element := range []KnownKey{KWCPCIUserData, KWCPCISource, KWCPCIDestination} {
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

	i.pduUserData, _ = KWO[spi.Message](kwArgs, KWCPCIUserData, nil)
	i.pduSource, _ = KWO[*Address](kwArgs, KWCPCISource, nil)
	i.pduDestination, _ = KWO[*Address](kwArgs, KWCPCIDestination, nil)
	return i
}

func (p *__PCI) GetDebugAttr(attr string) any {
	switch attr {
	case "pduUserData":
		return p.pduUserData
	case "pduSource":
		if p.pduSource != nil {
			return p.pduSource
		}
	case "pduDestination":
		if p.pduDestination != nil {
			return p.pduDestination
		}
	default:
		return nil
	}
	return nil
}

func (p *__PCI) GetLeafName() string {
	return p._leafName
}

func (p *__PCI) SetRootMessage(rootMessage spi.Message) {
	p.rootMessage = rootMessage
}

func (p *__PCI) GetRootMessage() spi.Message {
	return p.rootMessage
}

func (p *__PCI) SetPDUUserData(pduUserData spi.Message) {
	p.pduUserData = pduUserData
}

func (p *__PCI) GetPDUUserData() spi.Message {
	return p.pduUserData
}

func (p *__PCI) GetPDUSource() *Address {
	return p.pduSource
}

func (p *__PCI) SetPDUSource(source *Address) {
	p.pduSource = source
}

func (p *__PCI) GetPDUDestination() *Address {
	return p.pduDestination
}

func (p *__PCI) SetPDUDestination(destination *Address) {
	p.pduDestination = destination
}

func (p *__PCI) Update(pci Arg) error {
	switch pci := pci.(type) {
	case IPCI:
		p.rootMessage = pci.GetRootMessage()
		p.pduUserData = pci.GetPDUUserData()
		p.pduSource = pci.GetPDUSource()
		p.pduDestination = pci.GetPDUDestination()
		return nil
	default:
		return errors.Errorf("invalid IPCI type %T", pci)
	}
}

func (p *__PCI) deepCopy() *__PCI {
	newP := &__PCI{
		nil,
		p.rootMessage, // those are immutable so no copy needed
		p.pduUserData, // those are immutable so no copy needed
		p.pduSource.deepCopy(),
		p.pduDestination.deepCopy(),
		p._leafName,
	}
	newP.DebugContents = NewDebugContents(newP, "pduUserData+", "pduSource", "pduDestination") // TODO: bit ugly to repeat that here again but what are the options...
	return newP
}

func (p *__PCI) DeepCopy() any {
	return p.deepCopy()
}

func (p *__PCI) Serialize() ([]byte, error) {
	if p.rootMessage == nil {
		return nil, errors.New("no pdu userdata")
	}
	return p.rootMessage.Serialize()
}

func (p *__PCI) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if p.rootMessage == nil {
		return errors.New("no pdu userdata")
	}
	return p.rootMessage.SerializeWithWriteBuffer(ctx, writeBuffer)
}

func (p *__PCI) GetLengthInBytes(ctx context.Context) uint16 {
	if p.rootMessage == nil {
		return 0
	}
	return p.rootMessage.GetLengthInBytes(ctx)
}

func (p *__PCI) GetLengthInBits(ctx context.Context) uint16 {
	if p.rootMessage == nil {
		return 0
	}
	return p.rootMessage.GetLengthInBits(ctx)
}

func (p *__PCI) String() string {
	if IsDebuggingActive() {
		return fmt.Sprintf("%s", p) // Delegate
	}
	return fmt.Sprintf("%s", p.rootMessage)
}
