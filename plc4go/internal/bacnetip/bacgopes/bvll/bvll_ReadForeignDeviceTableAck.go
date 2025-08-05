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

type ReadForeignDeviceTableAck struct {
	*_BVLPDU

	bvlciFDT []*FDTEntry
}

var _ BVLPDU = (*ReadForeignDeviceTableAck)(nil)

// TODO: check this arg desaster
func NewReadForeignDeviceTableAck(fdt []*FDTEntry, args Args, kwArgs KWArgs, options ...Option) (*ReadForeignDeviceTableAck, error) {
	r := &ReadForeignDeviceTableAck{}
	options = AddLeafTypeIfAbundant(options, r)
	r._BVLPDU = NewBVLPDU(args, kwArgs, options...).(*_BVLPDU)
	r.AddDebugContents(r, "bvlciFDT++")
	if r.GetRootMessage() == nil {
		r.SetRootMessage(readWriteModel.NewBVLCReadForeignDeviceTableAck(r.produceForeignDeviceTable()))
	}
	r.bvlciFunction = BVLCIReadForeignDeviceTableAck
	r.bvlciLength = uint16(4 + 10*len(fdt))
	r.bvlciFDT = fdt
	return r, nil
}

func (r *ReadForeignDeviceTableAck) GetDebugAttr(attr string) any {
	switch attr {
	case "bvlciFDT":
		if r.bvlciFDT != nil {
			return r.bvlciFDT
		}
	default:
		return nil
	}
	return nil
}

func (r *ReadForeignDeviceTableAck) GetBvlciFDT() []*FDTEntry {
	return r.bvlciFDT
}

func (r *ReadForeignDeviceTableAck) produceForeignDeviceTable() (entries []readWriteModel.BVLCForeignDeviceTableEntry, _ uint16) {
	for _, entry := range r.bvlciFDT {
		address := entry.FDAddress
		addr := address.AddrAddress[:4]
		port := uint16(47808)
		if address.AddrPort != nil {
			port = *address.AddrPort
		}
		entries = append(entries, readWriteModel.NewBVLCForeignDeviceTableEntry(addr, port, entry.FDTTL, entry.FDRemain))
	}
	return
}

func (r *ReadForeignDeviceTableAck) produceBvlciFDT(entries []readWriteModel.BVLCForeignDeviceTableEntry) (bvlciFDT []*FDTEntry) {
	for _, entry := range entries {
		addr := entry.GetIp()
		port := entry.GetPort()
		var portArray = make([]byte, 2)
		binary.BigEndian.PutUint16(portArray, port)
		address, _ := NewAddress(NA(append(addr, portArray...)))
		bvlciFDT = append(bvlciFDT, &FDTEntry{
			FDAddress: address,
			FDTTL:     entry.GetTtl(),
			FDRemain:  entry.GetSecondRemainingBeforePurge(),
		})
	}
	return
}

func (r *ReadForeignDeviceTableAck) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLCI:
		if err := bvlpdu.getBVLCI().Update(r); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
	}
	switch bvlpdu := bvlpdu.(type) {
	case PDUData:
		for _, fdte := range r.bvlciFDT {
			bvlpdu.PutData(fdte.FDAddress.AddrAddress...)
			bvlpdu.PutShort(fdte.FDTTL)
			bvlpdu.PutShort(fdte.FDRemain)
		}
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
	return nil
}

func (r *ReadForeignDeviceTableAck) Decode(bvlpdu Arg) error {
	if err := r._BVLCI.Update(bvlpdu); err != nil {
		return errors.Wrap(err, "error updating BVLCI")
	}
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		switch rm := bvlpdu.GetRootMessage().(type) {
		case readWriteModel.BVLCReadForeignDeviceTableAck:
			r.bvlciFDT = r.produceBvlciFDT(rm.GetTable())
			r.SetRootMessage(rm)
		}
	}
	switch bvlpdu := bvlpdu.(type) {
	case PDUData:
		r.SetPduData(bvlpdu.GetPduData())
	}
	return nil
}

func (r *ReadForeignDeviceTableAck) String() string {
	if r == nil {
		return "(*ReadForeignDeviceTableAck)(nil)"
	}
	return fmt.Sprintf("ReadForeignDeviceTableAck{%v, bvlciFDT: %v}", r._BVLPDU, r.bvlciFDT)
}
