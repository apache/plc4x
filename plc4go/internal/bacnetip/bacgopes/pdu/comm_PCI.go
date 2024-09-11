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
	"strings"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/globals"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Note: upstream this belongs to comm but that would create a circular dependency

type IPCI interface {
	spi.Message
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
	rootMessage    spi.Message
	pduUserData    spi.Message
	pduSource      *Address
	pduDestination *Address
}

var _ IPCI = (*__PCI)(nil)

func new__PCI(args Args, kwargs KWArgs) *__PCI {
	if _debug != nil {
		_debug("__init__ %r %r", args, kwargs)
	}
	i := &__PCI{
		rootMessage: KWO[spi.Message](kwargs, KWCompRootMessage, nil),
	}
	delete(kwargs, KWCompRootMessage)
	var myKwargs = make(KWArgs)
	var otherKwargs = make(KWArgs)
	for _, element := range []KnownKey{KWCPCIUserData, KWCPCISource, KWCPCIDestination} {
		if v, ok := kwargs[element]; ok {
			myKwargs[element] = v
		}
	}
	for k, v := range kwargs {
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

	i.pduUserData = KWO[spi.Message](kwargs, KWCPCIUserData, nil)
	i.pduSource = KWO[*Address](kwargs, KWCPCISource, nil)
	i.pduDestination = KWO[*Address](kwargs, KWCPCIDestination, nil)
	return i
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
	rootMessage := p.rootMessage // those are immutable so no copy needed
	pduUserData := p.pduUserData // those are immutable so no copy needed
	pduSource := p.pduSource
	if pduSource != nil {
		copyPduSource := *pduSource
		pduSource = &copyPduSource
	}
	pduDestination := p.pduDestination
	if pduDestination != nil {
		copyPduDestination := *pduDestination
		pduDestination = &copyPduDestination
	}
	return &__PCI{rootMessage, pduUserData, pduSource, pduDestination}
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
	if ExtendedPDUOutput {
		rootMessageString := "nil"
		if p.rootMessage != nil && ExtendedPDUOutput {
			rootMessageString = p.rootMessage.String()
			if strings.Contains(rootMessageString, "\n") {
				rootMessageString = "\n" + rootMessageString + "\n"
			}
		} else if p.rootMessage != nil {
			if bytes, err := p.rootMessage.Serialize(); err != nil {
				rootMessageString = err.Error()
			} else {
				rootMessageString = Btox(bytes, ".")
			}
		}
		pduUserDataString := "nil"
		if p.pduUserData != nil && ExtendedPDUOutput {
			pduUserDataString = p.pduUserData.String()
			if strings.Contains(pduUserDataString, "\n") {
				pduUserDataString = "\n" + pduUserDataString + "\n"
			}
		} else if p.pduUserData != nil {
			if bytes, err := p.pduUserData.Serialize(); err != nil {
				pduUserDataString = err.Error()
			} else {
				pduUserDataString = Btox(bytes, ".")
			}
		}
		return fmt.Sprintf("__PCI{rootMessage: %s, pduUserData: %s, pduSource: %s, pduDestination: %s}", rootMessageString, pduUserDataString, p.pduSource, p.pduDestination)
	} else {
		pduSourceStr := ""
		if p.pduSource != nil {
			pduSourceStr = "pduSource = " + p.pduSource.String()
		}
		pduDestinationStr := ""
		if p.pduDestination != nil {
			pduDestinationStr = "\npduDestination = " + p.pduDestination.String()
			if pduSourceStr == "" {
				pduDestinationStr = pduDestinationStr[1:]
			}
		}
		return fmt.Sprintf("%s%s", pduSourceStr, pduDestinationStr)
	}
}
