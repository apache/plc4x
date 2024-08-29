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

type ReadBroadcastDistributionTableAck struct {
	*_BVLPDU

	bvlciBDT []*Address
}

var _ BVLPDU = (*ReadBroadcastDistributionTableAck)(nil)

func NewReadBroadcastDistributionTableAck(opts ...func(*ReadBroadcastDistributionTableAck)) (*ReadBroadcastDistributionTableAck, error) {
	b := &ReadBroadcastDistributionTableAck{}
	for _, opt := range opts {
		opt(b)
	}
	b._BVLPDU = NewBVLPDU(model.NewBVLCReadBroadcastDistributionTableAck(b.produceBroadcastDistributionTable(), 0)).(*_BVLPDU)
	return b, nil
}

func WithReadBroadcastDistributionTableAckBDT(bdt ...*Address) func(*ReadBroadcastDistributionTableAck) {
	return func(b *ReadBroadcastDistributionTableAck) {
		b.bvlciBDT = bdt
	}
}

func (w *ReadBroadcastDistributionTableAck) GetBvlciBDT() []*Address {
	return w.bvlciBDT
}

func (w *ReadBroadcastDistributionTableAck) produceBroadcastDistributionTable() (entries []model.BVLCBroadcastDistributionTableEntry) {
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
		entries = append(entries, model.NewBVLCBroadcastDistributionTableEntry(addr, port, mask))
	}
	return
}

func (w *ReadBroadcastDistributionTableAck) produceBvlciBDT(entries []model.BVLCBroadcastDistributionTableEntry) (bvlciBDT []*Address) {
	for _, entry := range entries {
		addr := entry.GetIp()
		port := entry.GetPort()
		var portArray = make([]byte, 2)
		binary.BigEndian.PutUint16(portArray, port)
		address, _ := NewAddress(zerolog.Nop(), append(addr, portArray...))
		mask := binary.BigEndian.Uint32(entry.GetBroadcastDistributionMap())
		address.AddrMask = &mask
		bvlciBDT = append(bvlciBDT, address)
	}
	return
}

func (w *ReadBroadcastDistributionTableAck) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := bvlpdu.Update(w); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		for _, bdte := range w.bvlciBDT {
			bvlpdu.PutData(bdte.AddrAddress...)
			bvlpdu.PutLong(*bdte.AddrMask)
		}
		bvlpdu.setBVLC(w.bvlc)
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (w *ReadBroadcastDistributionTableAck) Decode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := w.Update(bvlpdu); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		switch rm := bvlpdu.GetRootMessage().(type) {
		case model.BVLCReadBroadcastDistributionTableAckExactly:
			switch bvlc := rm.(type) {
			case model.BVLCReadBroadcastDistributionTableAck:
				w.setBVLC(bvlc)
				w.bvlciBDT = w.produceBvlciBDT(bvlc.GetTable())
			}
		}
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (w *ReadBroadcastDistributionTableAck) String() string {
	return fmt.Sprintf("ReadBroadcastDistributionTableAck{%v, bvlciBDT: %v}", w._BVLPDU, w.bvlciBDT)
}
