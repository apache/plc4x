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

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/globals"
	"github.com/apache/plc4x/plc4go/spi"
)

type PDU interface {
	PCI
	PDUData
	DeepCopy() any
}

type _PDU struct {
	*_PCI
	*_PDUData
}

func NewPDU(args Args, kwargs KWArgs) PDU {
	if _debug != nil {
		_debug("__init__ %r %r", args, kwargs)
	}
	p := &_PDU{
		_PCI: NewPCI(args, kwargs),
	}
	p._PDUData = NewPDUData(args, kwargs).(*_PDUData)
	return p
}

func (p *_PDU) GetRootMessage() spi.Message {
	return p.rootMessage
}

func (p *_PDU) deepCopy() *_PDU {
	pduCopy := &_PDU{_PCI: p._PCI.deepCopy(), _PDUData: p._PDUData.deepCopy()}
	return pduCopy
}

func (p *_PDU) DeepCopy() any {
	return p.deepCopy()
}

func (p *_PDU) String() string {
	if ExtendedPDUOutput {
		return fmt.Sprintf("_PDU{%s}", p._PCI)
	}
	return fmt.Sprintf("<%T %s -> %s : %s>", p, p.GetPDUSource(), p.GetPDUDestination(), p._PDUData)
}
