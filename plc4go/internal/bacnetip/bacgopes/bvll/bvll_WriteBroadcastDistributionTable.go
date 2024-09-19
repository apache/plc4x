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

package bvll

import (
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type WriteBroadcastDistributionTable struct {
	*_BVLPDU

	bvlciBDT []*Address
}

var _ BVLPDU = (*WriteBroadcastDistributionTable)(nil)

// TODO: check this arg desaster
func NewWriteBroadcastDistributionTable(bdt []*Address, args Args, kwArgs KWArgs, options ...Option) (*WriteBroadcastDistributionTable, error) {
	w := &WriteBroadcastDistributionTable{}
	options = AddLeafTypeIfAbundant(options, w)
	w._BVLPDU = NewBVLPDU(args, kwArgs, options...).(*_BVLPDU)
	w.AddDebugContents(w, "bvlciBDT")
	if w.GetRootMessage() == nil {
		w.SetRootMessage(readWriteModel.NewBVLCWriteBroadcastDistributionTable(w.produceBroadcastDistributionTable()))
	}
	w.bvlciFunction = BVLCIWriteBroadcastDistributionTable
	w.bvlciLength = uint16(4 + 10*len(bdt))
	w.bvlciBDT = bdt
	return w, nil
}

func (w *WriteBroadcastDistributionTable) GetDebugAttr(attr string) any {
	switch attr {
	case "bvlciBDT":
		if w.bvlciBDT != nil {
			return w.bvlciBDT
		}
	default:
		return nil
	}
	return nil
}

func (w *WriteBroadcastDistributionTable) GetBvlciBDT() []*Address {
	return w.bvlciBDT
}

func (w *WriteBroadcastDistributionTable) produceBroadcastDistributionTable() (entries []readWriteModel.BVLCBroadcastDistributionTableEntry, _ uint16) {
	for _, address := range w.bvlciBDT {
		addr := address.AddrAddress[:4]
		port := uint16(47808)
		if address.AddrPort != nil {
			port = *address.AddrPort
		}
		mask := make([]byte, 4)
		if address.AddrMask != nil {
			binary.BigEndian.PutUint32(mask, *address.AddrMask)
		}
		entries = append(entries, readWriteModel.NewBVLCBroadcastDistributionTableEntry(addr, port, mask))
	}
	return
}

func (w *WriteBroadcastDistributionTable) produceBvlciBDT(entries []readWriteModel.BVLCBroadcastDistributionTableEntry) (bvlciBDT []*Address) {
	for _, entry := range entries {
		addr := entry.GetIp()
		port := entry.GetPort()
		var portArray = make([]byte, 2)
		binary.BigEndian.PutUint16(portArray, port)
		address, _ := NewAddress(NA(append(addr, portArray...)))
		mask := binary.BigEndian.Uint32(entry.GetBroadcastDistributionMap())
		address.AddrMask = &mask
		bvlciBDT = append(bvlciBDT, address)
	}
	return
}

func (w *WriteBroadcastDistributionTable) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLCI:
		if err := bvlpdu.getBVLCI().Update(w); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
	}
	switch bvlpdu := bvlpdu.(type) {
	case PDUData:
		for _, bdte := range w.bvlciBDT {
			bvlpdu.PutData(bdte.AddrAddress...)
			bvlpdu.PutLong(*bdte.AddrMask)
		}
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
	return nil
}

func (w *WriteBroadcastDistributionTable) Decode(bvlpdu Arg) error {
	if err := w._BVLCI.Update(bvlpdu); err != nil {
		return errors.Wrap(err, "error updating BVLCI")
	}
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		switch rm := bvlpdu.GetRootMessage().(type) {
		case readWriteModel.BVLCWriteBroadcastDistributionTable:
			switch bvlc := rm.(type) {
			case readWriteModel.BVLCWriteBroadcastDistributionTable:
				w.bvlciBDT = w.produceBvlciBDT(bvlc.GetTable())
				w.SetRootMessage(rm)
			}
		}
	}
	switch bvlpdu := bvlpdu.(type) {
	case PDUData:
		w.SetPduData(bvlpdu.GetPduData())
	}
	return nil
}

func (w *WriteBroadcastDistributionTable) String() string {
	if w == nil {
		return "(*WriteBroadcastDistributionTable)(nil)"
	}
	return fmt.Sprintf("WriteBroadcastDistributionTable{%v, bvlciBDT: %v}", w._BVLPDU, w.bvlciBDT)
}
