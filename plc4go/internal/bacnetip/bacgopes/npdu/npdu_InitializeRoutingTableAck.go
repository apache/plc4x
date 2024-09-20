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

type InitializeRoutingTableAck struct {
	*_NPDU

	messageType uint8

	irtaTable []*RoutingTableEntry
}

func NewInitializeRoutingTableAck(args Args, kwArgs KWArgs, options ...Option) (*InitializeRoutingTableAck, error) {
	i := &InitializeRoutingTableAck{
		messageType: 0x07,
	}
	ApplyAppliers(options, i)
	options = AddLeafTypeIfAbundant(options, i)
	options = AddNLMIfAbundant(options, model.NewNLMInitializeRoutingTableAck(i.produceNLMInitializeRoutingTableAckPortMapping()))
	npdu, err := NewNPDU(args, kwArgs, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating NPDU")
	}
	i._NPDU = npdu.(*_NPDU)
	i.AddDebugContents(i, "irtaTable++")

	i.npduNetMessage = &i.messageType
	return i, nil
}

// TODO: check if this is rather a KWArgs
func WithInitializeRoutingTableAckIrtaTable(irtaTable ...*RoutingTableEntry) GenericApplier[*InitializeRoutingTableAck] {
	return WrapGenericApplier(func(r *InitializeRoutingTableAck) { r.irtaTable = irtaTable })
}

func (i *InitializeRoutingTableAck) GetDebugAttr(attr string) any {
	switch attr {
	case "irtaTable":
		return i.irtaTable
	}
	return nil
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
