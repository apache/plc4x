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
	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type InitializeRoutingTable struct {
	*_NPDU

	messageType uint8

	irtTable []*RoutingTableEntry
}

func NewInitializeRoutingTable(args Args, kwArgs KWArgs, options ...Option) (*InitializeRoutingTable, error) {
	i := &InitializeRoutingTable{
		messageType: 0x06,
	}
	ApplyAppliers(options, i)
	options = AddLeafTypeIfAbundant(options, i)
	options = AddNLMIfAbundant(options, model.NewNLMInitializeRoutingTable(i.produceNLMInitializeRoutingTablePortMapping()))
	npdu, err := NewNPDU(args, kwArgs, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating NPDU")
	}
	i._NPDU = npdu.(*_NPDU)
	i.AddDebugContents(i, "irtTable++")

	i.npduNetMessage = &i.messageType
	return i, nil
}

// TODO: check if this is rather a KWArgs
func WithInitializeRoutingTableIrtTable(irtTable ...*RoutingTableEntry) GenericApplier[*InitializeRoutingTable] {
	return WrapGenericApplier(func(r *InitializeRoutingTable) { r.irtTable = irtTable })
}

func (i *InitializeRoutingTable) GetDebugAttr(attr string) any {
	switch attr {
	case "irtTable":
		return i.irtTable
	}
	return nil
}

func (i *InitializeRoutingTable) GetIrtTable() []*RoutingTableEntry {
	return i.irtTable
}

func (i *InitializeRoutingTable) produceNLMInitializeRoutingTablePortMapping() (numberOfPorts uint8, mappings []model.NLMInitializeRoutingTablePortMapping, _ uint16) {
	numberOfPorts = uint8(len(i.irtTable))
	mappings = make([]model.NLMInitializeRoutingTablePortMapping, numberOfPorts)
	for i, entry := range i.irtTable {
		mappings[i] = model.NewNLMInitializeRoutingTablePortMapping(entry.tuple())
	}
	return
}

func (i *InitializeRoutingTable) produceIRTTable(mappings []model.NLMInitializeRoutingTablePortMapping) (irtTable []*RoutingTableEntry) {
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

func (i *InitializeRoutingTable) Encode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPCI:
		if err := npdu.GetNPCI().Update(i); err != nil {
			return errors.Wrap(err, "error updating NPDU")
		}
	}
	switch npdu := npdu.(type) {
	case PDUData:
		npdu.Put(uint8(len(i.irtTable)))
		for _, rte := range i.irtTable {
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

func (i *InitializeRoutingTable) Decode(npdu Arg) error {
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
			case model.NLMInitializeRoutingTable:
				i.irtTable = i.produceIRTTable(nlm.GetPortMappings())
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
