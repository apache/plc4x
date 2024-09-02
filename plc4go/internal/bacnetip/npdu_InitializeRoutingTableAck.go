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
	"fmt"

	"github.com/pkg/errors"

	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type InitializeRoutingTableAck struct {
	*_NPDU

	messageType uint8

	irtaTable []*RoutingTableEntry
}

func NewInitializeRoutingTableAck(opts ...func(*InitializeRoutingTableAck)) (*InitializeRoutingTableAck, error) {
	i := &InitializeRoutingTableAck{
		messageType: 0x07,
	}
	for _, opt := range opts {
		opt(i)
	}
	npdu, err := NewNPDU(model.NewNLMInitializeRoutingTableAck(i.produceNLMInitializeRoutingTableAckPortMapping()), nil)
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

func (r *InitializeRoutingTableAck) GetIrtaTable() []*RoutingTableEntry {
	return r.irtaTable
}

func (r *InitializeRoutingTableAck) produceNLMInitializeRoutingTableAckPortMapping() (numberOfPorts uint8, mappings []model.NLMInitializeRoutingTablePortMapping, _ uint16) {
	numberOfPorts = uint8(len(r.irtaTable))
	mappings = make([]model.NLMInitializeRoutingTablePortMapping, numberOfPorts)
	for i, entry := range r.irtaTable {
		mappings[i] = model.NewNLMInitializeRoutingTablePortMapping(entry.tuple())
	}
	return
}

func (r *InitializeRoutingTableAck) produceIRTTable(mappings []model.NLMInitializeRoutingTablePortMapping) (irtTable []*RoutingTableEntry) {
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

func (r *InitializeRoutingTableAck) Encode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPDU:
		if err := npdu.Update(r); err != nil {
			return errors.Wrap(err, "error updating NPDU")
		}
		npdu.Put(byte(len(r.irtaTable)))
		for _, rte := range r.irtaTable {
			npdu.PutShort(rte.rtDNET)
			npdu.Put(rte.rtPortId)
			npdu.Put(byte(len(rte.rtPortInfo)))
			npdu.PutData(rte.rtPortInfo...)
		}
		npdu.setNLM(r.nlm)
		npdu.setAPDU(r.apdu)
		return nil
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
}

func (r *InitializeRoutingTableAck) Decode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPDU:
		if err := r.Update(npdu); err != nil {
			return errors.Wrap(err, "error updating NPDU")
		}
		switch pduUserData := npdu.GetRootMessage().(type) {
		case model.NPDU:
			switch nlm := pduUserData.GetNlm().(type) {
			case model.NLMInitializeRoutingTableAck:
				r.setNLM(nlm)
				r.irtaTable = r.produceIRTTable(nlm.GetPortMappings())
			}
		}
		return nil
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
}

func (r *InitializeRoutingTableAck) String() string {
	return fmt.Sprintf("InitializeRoutingTableAck{%s, irtaTable: %v}", r._NPDU, r.irtaTable)
}
