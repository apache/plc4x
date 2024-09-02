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
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

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
	b._BVLPDU = NewBVLPDU(model.NewBVLCReadForeignDeviceTableAck(b.produceForeignDeviceTable(), 0)).(*_BVLPDU)
	return b, nil
}

func WithReadForeignDeviceTableAckFDT(fdts ...*FDTEntry) func(*ReadForeignDeviceTableAck) {
	return func(b *ReadForeignDeviceTableAck) {
		b.bvlciFDT = fdts
	}
}

func (w *ReadForeignDeviceTableAck) GetBvlciFDT() []*FDTEntry {
	return w.bvlciFDT
}

func (w *ReadForeignDeviceTableAck) produceForeignDeviceTable() (entries []model.BVLCForeignDeviceTableEntry) {
	for _, entry := range w.bvlciFDT {
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

func (w *ReadForeignDeviceTableAck) produceBvlciFDT(entries []model.BVLCForeignDeviceTableEntry) (bvlciFDT []*FDTEntry) {
	for _, entry := range entries {
		addr := entry.GetIp()
		port := entry.GetPort()
		var portArray = make([]byte, 2)
		binary.BigEndian.PutUint16(portArray, port)
		address, _ := NewAddress(zerolog.Nop(), append(addr, portArray...))
		bvlciFDT = append(bvlciFDT, &FDTEntry{
			FDAddress: address,
			FDTTL:     entry.GetTtl(),
			FDRemain:  entry.GetSecondRemainingBeforePurge(),
		})
	}
	return
}

func (w *ReadForeignDeviceTableAck) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := bvlpdu.Update(w); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		for _, fdte := range w.bvlciFDT {
			bvlpdu.PutData(fdte.FDAddress.AddrAddress...)
			bvlpdu.PutShort(fdte.FDTTL)
			bvlpdu.PutShort(fdte.FDRemain)
		}
		bvlpdu.setBVLC(w.bvlc)
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (w *ReadForeignDeviceTableAck) Decode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := w.Update(bvlpdu); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		switch rm := bvlpdu.GetRootMessage().(type) {
		case model.BVLCReadForeignDeviceTableAck:
			switch bvlc := rm.(type) {
			case model.BVLCReadForeignDeviceTableAck:
				w.setBVLC(bvlc)
				w.bvlciFDT = w.produceBvlciFDT(bvlc.GetTable())
			}
		}
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (w *ReadForeignDeviceTableAck) String() string {
	return fmt.Sprintf("ReadForeignDeviceTableAck{%v, bvlciFDT: %v}", w._BVLPDU, w.bvlciFDT)
}
