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
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	"github.com/apache/plc4x/plc4go/spi"
)

type PDU interface {
	PCI
	PDUData
	DeepCopy() any
}

type _PDU struct {
	*DefaultRFormatter
	*_PCI
	*_PDUData
}

func NewPDU(args Args, kwArgs KWArgs) PDU {
	if _debug != nil {
		_debug("__init__ %r %r", args, kwArgs)
	}
	p := &_PDU{}
	p._PCI = NewPCI(args, kwArgs)
	p._PDUData = NewPDUData(args, kwArgs).(*_PDUData)
	p.DefaultRFormatter = NewDefaultRFormatter(p._PCI, p._PDUData)
	return p
}

func (p *_PDU) GetRootMessage() spi.Message {
	return p.rootMessage
}

func (p *_PDU) deepCopy() *_PDU {
	pduCopy := &_PDU{p.DefaultRFormatter, p._PCI.deepCopy(), p._PDUData.deepCopy()}
	return pduCopy
}

func (p *_PDU) DeepCopy() any {
	return p.deepCopy()
}

func (p *_PDU) Format(s fmt.State, v rune) {
	switch v {
	case 's':
		_, _ = fmt.Fprint(s, p.String())
	default:
		p.DefaultRFormatter.Format(s, v)
	}
}

func (p *_PDU) String() string {
	return fmt.Sprintf("<%T %s -> %s : %s>", p, p.GetPDUSource(), p.GetPDUDestination(), p._PDUData)
}
