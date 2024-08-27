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

type InitializeRoutingTable struct {
	*_NPDU

	irtTable []*RoutingTableEntry
}

func NewInitializeRoutingTable(opts ...func(*InitializeRoutingTable)) (*InitializeRoutingTable, error) {
	i := &InitializeRoutingTable{}
	for _, opt := range opts {
		opt(i)
	}
	npdu, err := NewNPDU(model.NewNLMInitializeRoutingTable(i.produceNLMInitializeRoutingTablePortMapping()), nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating NPDU")
	}
	i._NPDU = npdu.(*_NPDU)
	return i, nil
}

func WithInitializeRoutingTableIrtTable(irtTable ...*RoutingTableEntry) func(*InitializeRoutingTable) {
	return func(r *InitializeRoutingTable) {
		r.irtTable = irtTable
	}
}

func (r *InitializeRoutingTable) GetIrtTable() []*RoutingTableEntry {
	return r.irtTable
}

func (r *InitializeRoutingTable) produceNLMInitializeRoutingTablePortMapping() (numberOfPorts uint8, mappings []model.NLMInitializeRoutingTablePortMapping, _ uint16) {
	numberOfPorts = uint8(len(r.irtTable))
	mappings = make([]model.NLMInitializeRoutingTablePortMapping, numberOfPorts)
	for i, entry := range r.irtTable {
		mappings[i] = model.NewNLMInitializeRoutingTablePortMapping(entry.tuple())
	}
	return
}

func (r *InitializeRoutingTable) produceIRTTable(mappings []model.NLMInitializeRoutingTablePortMapping) (irtTable []*RoutingTableEntry) {
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

func (r *InitializeRoutingTable) Encode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPDU:
		if err := npdu.Update(r); err != nil {
			return errors.Wrap(err, "error updating _NPCI")
		}
		for _, rte := range r.irtTable {
			npdu.PutShort(rte.rtDNET)
			npdu.Put(rte.rtPortId)
			npdu.Put(byte(len(rte.rtPortInfo)))
			npdu.PutData(rte.rtPortInfo...)
		}
		npdu.setNPDU(r.npdu)
		npdu.setNLM(r.nlm)
		npdu.setAPDU(r.apdu)
		return nil
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
}

func (r *InitializeRoutingTable) Decode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPDU:
		if err := r.Update(npdu); err != nil {
			return errors.Wrap(err, "error updating _NPCI")
		}
		switch pduUserData := npdu.GetRootMessage().(type) {
		case model.NPDUExactly:
			switch nlm := pduUserData.GetNlm().(type) {
			case model.NLMInitializeRoutingTable:
				r.setNLM(nlm)
				r.irtTable = r.produceIRTTable(nlm.GetPortMappings())
			}
		}
		return nil
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
}

func (r *InitializeRoutingTable) String() string {
	return fmt.Sprintf("InitializeRoutingTable{%s, irtTable: %v}", r._NPDU, r.irtTable)
}
