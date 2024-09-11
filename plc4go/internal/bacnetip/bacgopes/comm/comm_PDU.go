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

package comm

import (
	"fmt"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/globals"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	"github.com/apache/plc4x/plc4go/spi"
)

// TODO: this is named PDU usually. If we get rid of the . import for comm all over the place this might work again

type CPDU interface {
	PCI
	PDUData
	DeepCopy() any
}

type _PDU struct {
	PCI
	PDUData
}

func NewCPDU(args Args, kwargs KWArgs) CPDU {
	data := GAO[any](args, 0, nil)
	if _debug != nil {
		_debug("__init__ %r %r", args, kwargs)
	}

	// pick up some optional kwargs
	userData := KWO[spi.Message](kwargs, KWCPCIUserData, nil)
	source := KWO[*Address](kwargs, KWCPCISource, nil)
	destination := KWO[*Address](kwargs, KWCPCIDestination, nil)

	// carry source and destination from another PDU
	// so this can act like a copy constructor
	if data, ok := data.(PDU); ok {
		// allow parameters to override values
		userData = OR(userData, data.GetPDUUserData())
		source = OR(source, data.GetPDUSource())
		destination = OR(destination, data.GetPDUDestination())
	}

	// now continue on
	p := &_PDU{
		PCI: NewPCI(args, NKW(KWCPCIUserData, userData, KWCPCISource, source, destination, KWCPCIDestination, destination)),
	}
	p.PDUData = NewPDUData(NA(KWO[spi.Message](kwargs, KWCompRootMessage, nil)), NoKWArgs)
	return p
}

func (p *_PDU) deepCopy() *_PDU {
	pduCopy := &_PDU{PCI: p.PCI.DeepCopy().(PCI), PDUData: p.PDUData.DeepCopy().(PDUData)}
	return pduCopy
}

func (p *_PDU) DeepCopy() any {
	return p.deepCopy()
}

func (p *_PDU) GetName() string {
	return "PDU"
}

func (p *_PDU) String() string {
	if ExtendedPDUOutput {
		return fmt.Sprintf("_PDU{%s}", p.PCI)
	}
	return fmt.Sprintf("<%T %s -> %s : %s>", p, p.GetPDUSource(), p.GetPDUDestination(), p.PDUData)
}
