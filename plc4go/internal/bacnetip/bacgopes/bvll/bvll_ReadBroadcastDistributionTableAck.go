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

type ReadBroadcastDistributionTableAck struct {
	*_BVLPDU

	bvlciBDT []*Address
}

var _ BVLPDU = (*ReadBroadcastDistributionTableAck)(nil)

// TODO: check this arg desaster
func NewReadBroadcastDistributionTableAck(bdt []*Address, args Args, kwArgs KWArgs, options ...Option) (*ReadBroadcastDistributionTableAck, error) {
	r := &ReadBroadcastDistributionTableAck{}
	options = AddLeafTypeIfAbundant(options, r)
	r._BVLPDU = NewBVLPDU(args, kwArgs, options...).(*_BVLPDU)
	r.AddDebugContents(r, "bvlciBDT")
	if r.GetRootMessage() == nil {
		r.SetRootMessage(readWriteModel.NewBVLCReadBroadcastDistributionTableAck(r.produceBroadcastDistributionTable()))
	}
	r.bvlciFunction = BVLCIReadBroadcastDistributionTableAck
	r.bvlciLength = uint16(4 + 10*len(bdt))
	r.bvlciBDT = bdt
	return r, nil
}

func (r *ReadBroadcastDistributionTableAck) GetDebugAttr(attr string) any {
	switch attr {
	case "bvlciBDT":
		if r.bvlciBDT != nil {
			return r.bvlciBDT
		}
	default:
		return nil
	}
	return nil
}

func (r *ReadBroadcastDistributionTableAck) GetBvlciBDT() []*Address {
	return r.bvlciBDT
}

func (r *ReadBroadcastDistributionTableAck) produceBroadcastDistributionTable() (entries []readWriteModel.BVLCBroadcastDistributionTableEntry, _ uint16) {
	for _, address := range r.bvlciBDT {
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

func (r *ReadBroadcastDistributionTableAck) produceBvlciBDT(entries []readWriteModel.BVLCBroadcastDistributionTableEntry) (bvlciBDT []*Address) {
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

func (r *ReadBroadcastDistributionTableAck) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLCI:
		if err := bvlpdu.getBVLCI().Update(r); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
	}
	switch bvlpdu := bvlpdu.(type) {
	case PDUData:
		for _, bdte := range r.bvlciBDT {
			bvlpdu.PutData(bdte.AddrAddress...)
			bvlpdu.PutLong(*bdte.AddrMask)
		}
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
	return nil
}

func (r *ReadBroadcastDistributionTableAck) Decode(bvlpdu Arg) error {
	if err := r._BVLCI.Update(bvlpdu); err != nil {
		return errors.Wrap(err, "error updating BVLCI")
	}
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		switch rm := bvlpdu.GetRootMessage().(type) {
		case readWriteModel.BVLCReadBroadcastDistributionTableAck:
			switch bvlc := rm.(type) {
			case readWriteModel.BVLCReadBroadcastDistributionTableAck:
				r.bvlciBDT = r.produceBvlciBDT(bvlc.GetTable())
				r.SetRootMessage(rm)
			}
		}
	}
	switch bvlpdu := bvlpdu.(type) {
	case PDUData:
		r.SetPduData(bvlpdu.GetPduData())
	}
	return nil
}

func (r *ReadBroadcastDistributionTableAck) String() string {
	if r == nil {
		return "(*ReadBroadcastDistributionTableAck)(nil)"
	}
	return fmt.Sprintf("ReadBroadcastDistributionTableAck{%v, bvlciBDT: %v}", r._BVLPDU, r.bvlciBDT)
}
