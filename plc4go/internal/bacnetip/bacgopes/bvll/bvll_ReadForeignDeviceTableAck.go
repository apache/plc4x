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
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type ReadForeignDeviceTableAck struct {
	*_BVLPDU

	bvlciFDT []*FDTEntry
}

var _ BVLPDU = (*ReadForeignDeviceTableAck)(nil)

func NewReadForeignDeviceTableAck(opts ...func(*ReadForeignDeviceTableAck)) (*ReadForeignDeviceTableAck, error) {
	b := &ReadForeignDeviceTableAck{}
	for _, opt := range opts {
		opt(b)
	}
	b._BVLPDU = NewBVLPDU(NoArgs, NewKWArgs(KWCompRootMessage, model.NewBVLCReadForeignDeviceTableAck(b.produceForeignDeviceTable(), 0))).(*_BVLPDU)
	return b, nil
}

func WithReadForeignDeviceTableAckFDT(fdts ...*FDTEntry) func(*ReadForeignDeviceTableAck) {
	return func(b *ReadForeignDeviceTableAck) {
		b.bvlciFDT = fdts
	}
}

func (r *ReadForeignDeviceTableAck) GetBvlciFDT() []*FDTEntry {
	return r.bvlciFDT
}

func (r *ReadForeignDeviceTableAck) produceForeignDeviceTable() (entries []model.BVLCForeignDeviceTableEntry) {
	for _, entry := range r.bvlciFDT {
		address := entry.FDAddress
		addr := address.AddrAddress[:4]
		port := uint16(47808)
		if address.AddrPort != nil {
			port = *address.AddrPort
		}
		entries = append(entries, model.NewBVLCForeignDeviceTableEntry(addr, port, entry.FDTTL, entry.FDRemain))
	}
	return
}

func (r *ReadForeignDeviceTableAck) produceBvlciFDT(entries []model.BVLCForeignDeviceTableEntry) (bvlciFDT []*FDTEntry) {
	for _, entry := range entries {
		addr := entry.GetIp()
		port := entry.GetPort()
		var portArray = make([]byte, 2)
		binary.BigEndian.PutUint16(portArray, port)
		address, _ := NewAddress(NewArgs(append(addr, portArray...)))
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
		case model.BVLCReadForeignDeviceTableAck:
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
