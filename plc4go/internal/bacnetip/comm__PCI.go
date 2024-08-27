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

package bacnetip

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/apache/plc4x/plc4go/internal/bacnetip/globals"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type IPCI interface {
	spi.Message
	// GetRootMessage returns this. (TODO: check if type switch works without that, as this is spi.Message which delegates)
	GetRootMessage() spi.Message
	SetPDUUserData(spi.Message)
	GetPDUUserData() spi.Message
	GetPDUSource() *Address
	SetPDUSource(source *Address)
	GetPDUDestination() *Address
	SetPDUDestination(*Address)
	Update(pci Arg) error
}

type __PCI struct {
	rootMessage    spi.Message
	pduUserData    spi.Message
	pduSource      *Address
	pduDestination *Address
}

var _ IPCI = (*__PCI)(nil)

func new__PCI(rootMessage spi.Message, pduUserData spi.Message, pduSource *Address, pduDestination *Address) *__PCI {
	return &__PCI{rootMessage, pduUserData, pduSource, pduDestination}
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
	rootMessageString := "nil"
	if p.rootMessage != nil && globals.ExtendedPDUOutput {
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
	if p.pduUserData != nil && globals.ExtendedPDUOutput {
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
}
