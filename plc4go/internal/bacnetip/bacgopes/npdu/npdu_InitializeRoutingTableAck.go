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

package npdu

import (
	"fmt"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type InitializeRoutingTableAck struct {
	*_NPDU

	messageType uint8

	irtaTable []*RoutingTableEntry
}

func NewInitializeRoutingTableAck(args Args, kwArgs KWArgs, opts ...func(*InitializeRoutingTableAck)) (*InitializeRoutingTableAck, error) {
	i := &InitializeRoutingTableAck{
		messageType: 0x07,
	}
	for _, opt := range opts {
		opt(i)
	}
	if _, ok := kwArgs[KWCompNLM]; ok {
		kwArgs[KWCompNLM] = model.NewNLMInitializeRoutingTableAck(i.produceNLMInitializeRoutingTableAckPortMapping())
	}
	npdu, err := NewNPDU(args, kwArgs)
	if err != nil {
		return nil, errors.Wrap(err, "error creating NPDU")
	}
	i._NPDU = npdu.(*_NPDU)

	i.npduNetMessage = &i.messageType
	return i, nil
}

func WithInitializeRoutingTableAckIrtaTable(irtaTable ...*RoutingTableEntry) func(*InitializeRoutingTableAck) {
	return func(r *InitializeRoutingTableAck) {
		r.irtaTable = irtaTable
	}
}

func (i *InitializeRoutingTableAck) GetIrtaTable() []*RoutingTableEntry {
	return i.irtaTable
}

func (i *InitializeRoutingTableAck) produceNLMInitializeRoutingTableAckPortMapping() (numberOfPorts uint8, mappings []model.NLMInitializeRoutingTablePortMapping, _ uint16) {
	numberOfPorts = uint8(len(i.irtaTable))
	mappings = make([]model.NLMInitializeRoutingTablePortMapping, numberOfPorts)
	for i, entry := range i.irtaTable {
		mappings[i] = model.NewNLMInitializeRoutingTablePortMapping(entry.tuple())
	}
	return
}

func (i *InitializeRoutingTableAck) produceIRTTable(mappings []model.NLMInitializeRoutingTablePortMapping) (irtTable []*RoutingTableEntry) {
	irtTable = make([]*RoutingTableEntry, len(mappings))
	for i, entry := range mappings {
		irtTable[i] = NewRoutingTableEntry(
			WithRoutingTableEntryDestinationNetworkAddress(entry.GetDestinationNetworkAddress()),
			WithRoutingTableEntryPortId(entry.GetPortId()),
			WithRoutingTableEntryPortInfo(entry.GetPortInfo()),
		)
	}
	return
}

func (i *InitializeRoutingTableAck) Encode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPCI:
		if err := npdu.GetNPCI().Update(i); err != nil {
			return errors.Wrap(err, "error updating NPDU")
		}
	}
	switch npdu := npdu.(type) {
	case PDUData:
		npdu.Put(byte(len(i.irtaTable)))
		for _, rte := range i.irtaTable {
			npdu.PutShort(rte.rtDNET)
			npdu.Put(rte.rtPortId)
			npdu.Put(byte(len(rte.rtPortInfo)))
			npdu.PutData(rte.rtPortInfo...)
		}
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
	return nil
}

func (i *InitializeRoutingTableAck) Decode(npdu Arg) error {
	if err := i.GetNPCI().Update(npdu); err != nil {
		return errors.Wrap(err, "error updating NPCI")
	}
	switch npdu := npdu.(type) {
	case NPDU:
		if err := i.Update(npdu); err != nil {
			return errors.Wrap(err, "error updating NPDU")
		}
		switch rm := npdu.GetRootMessage().(type) {
		case model.NPDU:
			switch nlm := rm.GetNlm().(type) {
			case model.NLMInitializeRoutingTableAck:
				i.irtaTable = i.produceIRTTable(nlm.GetPortMappings())
				i.SetRootMessage(rm)
			}
		}
	}
	switch npdu := npdu.(type) {
	case PDUData:
		i.SetPduData(npdu.GetPduData())
	}
	return nil
}

func (i *InitializeRoutingTableAck) String() string {
	return fmt.Sprintf("InitializeRoutingTableAck{%s, irtaTable: %v}", i._NPDU, i.irtaTable)
}
