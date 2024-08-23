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
	"context"
	"encoding/binary"
	"fmt"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// BVLPDUTypes is a dictionary of message type values and structs
var BVLPDUTypes map[uint8]func() interface{ Decode(Arg) error }

type BVLCI interface {
	PCI

	Encode(pdu Arg) error
	Decode(pdu Arg) error
}

type _BVLCI struct {
	*_PCI
	*DebugContents
}

var _ BVLCI = (*_BVLCI)(nil)

func NewBVLCI(pdu spi.Message) BVLCI {
	b := &_BVLCI{}
	b._PCI = newPCI(pdu, nil, nil, nil, false, readWriteModel.NPDUNetworkPriority_NORMAL_MESSAGE)
	return b
}

func (b *_BVLCI) Update(bvlci Arg) error {
	if err := b._PCI.Update(bvlci); err != nil {
		return errors.Wrap(err, "Update BVLCI")
	}
	switch bvlci := bvlci.(type) {
	case BVLCI:
		// TODO: update coordinates....
		return nil
	default:
		return errors.Errorf("invalid BVLCI type %T", bvlci)
	}
}

func (b *_BVLCI) Encode(pdu Arg) error {
	if err := pdu.(interface{ Update(Arg) error }).Update(b); err != nil { // TODO: better validate that arg is really PDUData... use switch similar to Update
		return errors.Wrap(err, "error updating pdu")
	}
	// TODO: what should we do here??
	return nil
}

func (b *_BVLCI) Decode(pdu Arg) error {
	if err := b._PCI.Update(pdu); err != nil {
		return errors.Wrap(err, "error updating pdu")
	}
	// TODO: what should we do here??
	return nil
}

func (b *_BVLCI) deepCopy() *_BVLCI {
	return &_BVLCI{_PCI: b._PCI.deepCopy()}
}

type BVLPDU interface {
	readWriteModel.BVLC
	BVLCI
	PDUData

	setBVLC(readWriteModel.BVLC)
	getBVLC() readWriteModel.BVLC
}

type _BVLPDU struct {
	*_BVLCI
	*_PDUData

	bvlc readWriteModel.BVLC
}

var _ BVLPDU = (*_BVLPDU)(nil)

func NewBVLPDU(bvlc readWriteModel.BVLC) BVLPDU {
	b := &_BVLPDU{
		bvlc: bvlc,
	}
	b._BVLCI = NewBVLCI(bvlc).(*_BVLCI)
	b._PDUData = NewPDUData(NoArgs).(*_PDUData)
	return b
}

// Deprecated: check if needed as we do it in update
func (b *_BVLPDU) setBVLC(bvlc readWriteModel.BVLC) {
	b.bvlc = bvlc
}

func (b *_BVLPDU) getBVLC() readWriteModel.BVLC {
	return b.bvlc
}

func (b *_BVLPDU) Update(bvlci Arg) error {
	if err := b._BVLCI.Update(bvlci); err != nil {
		return errors.Wrap(err, "Update BVLCI")
	}
	switch bvlci := bvlci.(type) {
	case BVLCI:
		b.bvlc = b.getBVLC()
		// TODO: update coordinates....
		return nil
	default:
		return errors.Errorf("invalid BVLCI type %T", bvlci)
	}
}

func (b *_BVLPDU) Encode(pdu Arg) error {
	if err := b._BVLCI.Encode(pdu); err != nil {
		return errors.Wrap(err, "error encoding _BVLCI")
	}
	serialize, err := b.bvlc.Serialize()
	if err != nil {
		return errors.Wrap(err, "error serializing BVLC")
	}
	pdu.(interface{ PutData(n ...byte) }).PutData(serialize...) // TODO: ugly cast...
	return nil
}

func (b *_BVLPDU) Decode(pdu Arg) error {
	if err := b._BVLCI.Decode(pdu); err != nil {
		return errors.Wrap(err, "error decoding _BVLCI")
	}
	switch pdu := pdu.(type) {
	case PDUData:
		data := pdu.GetPduData()
		b.PutData(data...)
		var err error
		b.bvlc, err = readWriteModel.BVLCParse(context.Background(), data)
		if err != nil {
			return errors.Wrap(err, "error parsing NPDU")
		}
		b.rootMessage = b.bvlc
	}
	return nil
}

func (b *_BVLPDU) GetBvlcFunction() uint8 {
	if b.bvlc == nil {
		return 0
	}
	return b.bvlc.GetBvlcFunction()
}

func (b *_BVLPDU) GetBvlcPayloadLength() uint16 {
	if b.bvlc == nil {
		return 0
	}
	return b.bvlc.GetBvlcPayloadLength()
}

func (b *_BVLPDU) deepCopy() *_BVLPDU {
	return &_BVLPDU{_BVLCI: b._BVLCI.deepCopy(), _PDUData: b._PDUData.deepCopy(), bvlc: b.bvlc}
}

func (b *_BVLPDU) DeepCopy() PDU {
	return b.deepCopy()
}

func (b *_BVLPDU) String() string {
	return fmt.Sprintf("_BVLPDU{%s, PDUData: %s}", b._BVLCI, b._PDUData)
}

type Result struct {
	*_BVLPDU

	bvlciResultCode readWriteModel.BVLCResultCode
}

var _ BVLPDU = (*Result)(nil)

func NewResult(opts ...func(result *Result)) (*Result, error) {
	b := &Result{}
	for _, opt := range opts {
		opt(b)
	}
	b._BVLPDU = NewBVLPDU(readWriteModel.NewBVLCResult(b.bvlciResultCode)).(*_BVLPDU)
	return b, nil
}

func WithResultBvlciResultCode(code readWriteModel.BVLCResultCode) func(*Result) {
	return func(b *Result) {
		b.bvlciResultCode = code
	}
}

func (n *Result) GetBvlciResultCode() readWriteModel.BVLCResultCode {
	return n.bvlciResultCode
}

func (n *Result) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := bvlpdu.Update(n); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		bvlpdu.PutShort(uint16(n.bvlciResultCode))
		bvlpdu.setBVLC(n.bvlc)
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (n *Result) Decode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := n.Update(bvlpdu); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		switch rm := bvlpdu.GetRootMessage().(type) {
		case readWriteModel.BVLCResultExactly:
			switch bvlc := rm.(type) {
			case readWriteModel.BVLCResult:
				n.setBVLC(bvlc)
				n.bvlciResultCode = bvlc.GetCode()
			}
		}
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (n *Result) String() string {
	return fmt.Sprintf("Result{%v, bvlciResultCode: %v}", n._BVLPDU, n.bvlciResultCode)
}

type WriteBroadcastDistributionTable struct {
	*_BVLPDU

	bvlciBDT []*Address
}

var _ BVLPDU = (*WriteBroadcastDistributionTable)(nil)

func NewWriteBroadcastDistributionTable(opts ...func(*WriteBroadcastDistributionTable)) (*WriteBroadcastDistributionTable, error) {
	b := &WriteBroadcastDistributionTable{}
	for _, opt := range opts {
		opt(b)
	}
	b._BVLPDU = NewBVLPDU(readWriteModel.NewBVLCWriteBroadcastDistributionTable(b.produceBroadcastDistributionTable(), 0)).(*_BVLPDU)
	return b, nil
}

func WithWriteBroadcastDistributionTableBDT(bdt ...*Address) func(*WriteBroadcastDistributionTable) {
	return func(b *WriteBroadcastDistributionTable) {
		b.bvlciBDT = bdt
	}
}

func (w *WriteBroadcastDistributionTable) GetBvlciBDT() []*Address {
	return w.bvlciBDT
}

func (w *WriteBroadcastDistributionTable) produceBroadcastDistributionTable() (entries []readWriteModel.BVLCBroadcastDistributionTableEntry) {
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
		address, _ := NewAddress(zerolog.Nop(), append(addr, portArray...))
		mask := binary.BigEndian.Uint32(entry.GetBroadcastDistributionMap())
		address.AddrMask = &mask
		bvlciBDT = append(bvlciBDT, address)
	}
	return
}

func (w *WriteBroadcastDistributionTable) Encode(bvlpdu Arg) error {
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

func (w *WriteBroadcastDistributionTable) Decode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := w.Update(bvlpdu); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		switch rm := bvlpdu.GetRootMessage().(type) {
		case readWriteModel.BVLCWriteBroadcastDistributionTableExactly:
			switch bvlc := rm.(type) {
			case readWriteModel.BVLCWriteBroadcastDistributionTable:
				w.setBVLC(bvlc)
				w.bvlciBDT = w.produceBvlciBDT(bvlc.GetTable())
			}
		}
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (w *WriteBroadcastDistributionTable) String() string {
	return fmt.Sprintf("WriteBroadcastDistributionTable{%v, bvlciBDT: %v}", w._BVLPDU, w.bvlciBDT)
}

type ReadBroadcastDistributionTable struct {
	*_BVLPDU
}

var _ BVLPDU = (*ReadBroadcastDistributionTable)(nil)

func NewReadBroadcastDistributionTable(opts ...func(*ReadBroadcastDistributionTable)) (*ReadBroadcastDistributionTable, error) {
	b := &ReadBroadcastDistributionTable{}
	for _, opt := range opts {
		opt(b)
	}
	b._BVLPDU = NewBVLPDU(readWriteModel.NewBVLCReadBroadcastDistributionTable()).(*_BVLPDU)
	return b, nil
}

func (w *ReadBroadcastDistributionTable) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := bvlpdu.Update(w); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		bvlpdu.setBVLC(w.bvlc)
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (w *ReadBroadcastDistributionTable) Decode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := w.Update(bvlpdu); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		switch rm := bvlpdu.GetRootMessage().(type) {
		case readWriteModel.BVLCReadBroadcastDistributionTableExactly:
			switch bvlc := rm.(type) {
			case readWriteModel.BVLCReadBroadcastDistributionTable:
				w.setBVLC(bvlc)
			}
		}
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (w *ReadBroadcastDistributionTable) String() string {
	return fmt.Sprintf("ReadBroadcastDistributionTable{%v}", w._BVLPDU)
}

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
	b._BVLPDU = NewBVLPDU(readWriteModel.NewBVLCReadBroadcastDistributionTableAck(b.produceBroadcastDistributionTable(), 0)).(*_BVLPDU)
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

func (w *ReadBroadcastDistributionTableAck) produceBroadcastDistributionTable() (entries []readWriteModel.BVLCBroadcastDistributionTableEntry) {
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

func (w *ReadBroadcastDistributionTableAck) produceBvlciBDT(entries []readWriteModel.BVLCBroadcastDistributionTableEntry) (bvlciBDT []*Address) {
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
		case readWriteModel.BVLCReadBroadcastDistributionTableAckExactly:
			switch bvlc := rm.(type) {
			case readWriteModel.BVLCReadBroadcastDistributionTableAck:
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

type ForwardedNPDU struct {
	*_BVLPDU

	bvlciAddress *Address
}

var _ BVLPDU = (*ForwardedNPDU)(nil)

func NewForwardedNPDU(pdu PDU, opts ...func(*ForwardedNPDU)) (*ForwardedNPDU, error) {
	b := &ForwardedNPDU{}
	for _, opt := range opts {
		opt(b)
	}
	switch npdu := pdu.(type) {
	case readWriteModel.NPDUExactly:
		b._BVLPDU = NewBVLPDU(readWriteModel.NewBVLCForwardedNPDU(b.produceInnerNPDU(npdu))).(*_BVLPDU)
	case nil:
		b._BVLPDU = NewBVLPDU(nil).(*_BVLPDU)
	default:
		// TODO: re-encode seems expensive... check if there is a better option (e.g. only do it on the message bridge)
		data := pdu.GetPduData()
		parse, err := readWriteModel.NPDUParse(context.Background(), data, uint16(len(data)))
		if err != nil {
			return nil, errors.Wrap(err, "error re-encoding")
		}
		b._BVLPDU = NewBVLPDU(readWriteModel.NewBVLCForwardedNPDU(b.produceInnerNPDU(parse))).(*_BVLPDU)
	}
	return b, nil
}

func WithForwardedNPDUAddress(addr *Address) func(*ForwardedNPDU) {
	return func(b *ForwardedNPDU) {
		b.bvlciAddress = addr
	}
}

func (w *ForwardedNPDU) GetBvlciAddress() *Address {
	return w.bvlciAddress
}

func (w *ForwardedNPDU) produceInnerNPDU(inNpdu readWriteModel.NPDU) (ip []uint8, port uint16, npdu readWriteModel.NPDU, bvlcPayloadLength uint16) {
	ip = w.bvlciAddress.AddrAddress[:4]
	port = uint16(47808)
	if w.bvlciAddress.AddrPort != nil {
		port = *w.bvlciAddress.AddrPort
	}
	npdu = inNpdu
	return
}

func (w *ForwardedNPDU) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := bvlpdu.Update(w); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}

		// encode the addrress
		bvlpdu.PutData(w.bvlciAddress.AddrAddress...)

		// encode the rest of the data
		bvlpdu.PutData(w.GetPduData()...)

		bvlpdu.setBVLC(w.bvlc)
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (w *ForwardedNPDU) Decode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := w.Update(bvlpdu); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		switch rm := bvlpdu.GetRootMessage().(type) {
		case readWriteModel.BVLCForwardedNPDUExactly:
			addr := rm.GetIp()
			port := rm.GetPort()
			var portArray = make([]byte, 2)
			binary.BigEndian.PutUint16(portArray, port)
			var err error
			address, err := NewAddress(zerolog.Nop(), append(addr, portArray...))
			if err != nil {
				return errors.Wrap(err, "error creating address")
			}
			w.bvlciAddress = address
			npdu := rm.GetNpdu()
			pduData, err := npdu.Serialize()
			if err != nil {
				return errors.Wrap(err, "error serializing NPDU")
			}
			w.SetPduData(pduData)

			w.setBVLC(rm)
		}
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (w *ForwardedNPDU) String() string {
	return fmt.Sprintf("ForwardedNPDU{%v, bvlciAddress: %v}", w._BVLPDU, w.bvlciAddress)
}

type FDTEntry struct {
	FDAddress *Address
	FDTTL     uint16
	FDRemain  uint16
}

func (f FDTEntry) Equals(other any) bool {
	if f == other {
		return true
	}
	otherEntry, ok := other.(FDTEntry)
	if !ok {
		return false
	}
	return f.FDAddress.Equals(otherEntry.FDAddress) && f.FDTTL == otherEntry.FDTTL && f.FDRemain == otherEntry.FDRemain
}

type RegisterForeignDevice struct {
	*_BVLPDU

	bvlciTimeToLive uint16
}

var _ BVLPDU = (*RegisterForeignDevice)(nil)

func NewRegisterForeignDevice(opts ...func(RegisterForeignDevice *RegisterForeignDevice)) (*RegisterForeignDevice, error) {
	b := &RegisterForeignDevice{}
	for _, opt := range opts {
		opt(b)
	}
	b._BVLPDU = NewBVLPDU(readWriteModel.NewBVLCRegisterForeignDevice(b.bvlciTimeToLive)).(*_BVLPDU)
	return b, nil
}

func WithRegisterForeignDeviceBvlciTimeToLive(ttl uint16) func(*RegisterForeignDevice) {
	return func(b *RegisterForeignDevice) {
		b.bvlciTimeToLive = ttl
	}
}

func (n *RegisterForeignDevice) GetBvlciTimeToLive() uint16 {
	return n.bvlciTimeToLive
}

func (n *RegisterForeignDevice) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := bvlpdu.Update(n); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		bvlpdu.PutShort(n.bvlciTimeToLive)
		bvlpdu.setBVLC(n.bvlc)
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (n *RegisterForeignDevice) Decode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := n.Update(bvlpdu); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		switch rm := bvlpdu.GetRootMessage().(type) {
		case readWriteModel.BVLCRegisterForeignDeviceExactly:
			switch bvlc := rm.(type) {
			case readWriteModel.BVLCRegisterForeignDevice:
				n.setBVLC(bvlc)
				n.bvlciTimeToLive = bvlc.GetTtl()
			}
		}
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (n *RegisterForeignDevice) String() string {
	return fmt.Sprintf("RegisterForeignDevice{%v, bvlciTimeToLive: %v}", n._BVLPDU, n.bvlciTimeToLive)
}

type ReadForeignDeviceTable struct {
	*_BVLPDU
}

var _ BVLPDU = (*ReadForeignDeviceTable)(nil)

func NewReadForeignDeviceTable(opts ...func(*ReadForeignDeviceTable)) (*ReadForeignDeviceTable, error) {
	b := &ReadForeignDeviceTable{}
	for _, opt := range opts {
		opt(b)
	}
	b._BVLPDU = NewBVLPDU(readWriteModel.NewBVLCReadForeignDeviceTable()).(*_BVLPDU)
	return b, nil
}

func (w *ReadForeignDeviceTable) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := bvlpdu.Update(w); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		bvlpdu.setBVLC(w.bvlc)
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (w *ReadForeignDeviceTable) Decode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := w.Update(bvlpdu); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		switch rm := bvlpdu.GetRootMessage().(type) {
		case readWriteModel.BVLCReadForeignDeviceTableExactly:
			switch bvlc := rm.(type) {
			case readWriteModel.BVLCReadForeignDeviceTable:
				w.setBVLC(bvlc)
			}
		}
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (w *ReadForeignDeviceTable) String() string {
	return fmt.Sprintf("ReadForeignDeviceTable{%v}", w._BVLPDU)
}

type ReadForeignDeviceTableAck struct {
	*_BVLPDU

	bvlciFDT []FDTEntry
}

var _ BVLPDU = (*ReadForeignDeviceTableAck)(nil)

func NewReadForeignDeviceTableAck(opts ...func(*ReadForeignDeviceTableAck)) (*ReadForeignDeviceTableAck, error) {
	b := &ReadForeignDeviceTableAck{}
	for _, opt := range opts {
		opt(b)
	}
	b._BVLPDU = NewBVLPDU(readWriteModel.NewBVLCReadForeignDeviceTableAck(b.produceForeignDeviceTable(), 0)).(*_BVLPDU)
	return b, nil
}

func WithReadForeignDeviceTableAckFDT(fdts ...FDTEntry) func(*ReadForeignDeviceTableAck) {
	return func(b *ReadForeignDeviceTableAck) {
		b.bvlciFDT = fdts
	}
}

func (w *ReadForeignDeviceTableAck) GetBvlciFDT() []FDTEntry {
	return w.bvlciFDT
}

func (w *ReadForeignDeviceTableAck) produceForeignDeviceTable() (entries []readWriteModel.BVLCForeignDeviceTableEntry) {
	for _, entry := range w.bvlciFDT {
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

func (w *ReadForeignDeviceTableAck) produceBvlciFDT(entries []readWriteModel.BVLCForeignDeviceTableEntry) (bvlciFDT []FDTEntry) {
	for _, entry := range entries {
		addr := entry.GetIp()
		port := entry.GetPort()
		var portArray = make([]byte, 2)
		binary.BigEndian.PutUint16(portArray, port)
		address, _ := NewAddress(zerolog.Nop(), append(addr, portArray...))
		bvlciFDT = append(bvlciFDT, FDTEntry{
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
		case readWriteModel.BVLCReadForeignDeviceTableAckExactly:
			switch bvlc := rm.(type) {
			case readWriteModel.BVLCReadForeignDeviceTableAck:
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

type DeleteForeignDeviceTableEntry struct {
	*_BVLPDU

	bvlciAddress *Address
}

var _ BVLPDU = (*DeleteForeignDeviceTableEntry)(nil)

func NewDeleteForeignDeviceTableEntry(opts ...func(*DeleteForeignDeviceTableEntry)) (*DeleteForeignDeviceTableEntry, error) {
	d := &DeleteForeignDeviceTableEntry{}
	for _, opt := range opts {
		opt(d)
	}
	d._BVLPDU = NewBVLPDU(readWriteModel.NewBVLCDeleteForeignDeviceTableEntry(d.buildIPArgs())).(*_BVLPDU)
	return d, nil
}

func WithDeleteForeignDeviceTableEntryAddress(address *Address) func(*DeleteForeignDeviceTableEntry) {
	return func(d *DeleteForeignDeviceTableEntry) {
		d.bvlciAddress = address
	}
}

func (d *DeleteForeignDeviceTableEntry) buildIPArgs() (ip []uint8, port uint16) {
	if d.bvlciAddress == nil {
		return
	}
	ip = d.bvlciAddress.AddrAddress[:4]
	port = *d.bvlciAddress.AddrPort
	return
}

func (d *DeleteForeignDeviceTableEntry) buildAddress(ip []uint8, port uint16) *Address {
	var portArray = make([]byte, 2)
	binary.BigEndian.PutUint16(portArray, port)
	address, _ := NewAddress(zerolog.Nop(), append(ip, portArray...))
	return address
}

func (d *DeleteForeignDeviceTableEntry) GetBvlciAddress() *Address {
	return d.bvlciAddress
}

func (d *DeleteForeignDeviceTableEntry) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := bvlpdu.Update(d); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		bvlpdu.PutData(d.bvlciAddress.AddrAddress...)
		bvlpdu.setBVLC(d.bvlc)
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (d *DeleteForeignDeviceTableEntry) Decode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := d.Update(bvlpdu); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		switch rm := bvlpdu.GetRootMessage().(type) {
		case readWriteModel.BVLCDeleteForeignDeviceTableEntryExactly:
			switch bvlc := rm.(type) {
			case readWriteModel.BVLCDeleteForeignDeviceTableEntry:
				d.bvlciAddress = d.buildAddress(bvlc.GetIp(), bvlc.GetPort())
				d.setBVLC(bvlc)
			}
		}
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (d *DeleteForeignDeviceTableEntry) String() string {
	return fmt.Sprintf("DeleteForeignDeviceTableEntry{%v}", d._BVLPDU)
}

type DistributeBroadcastToNetwork struct {
	*_BVLPDU

	// post construct function
	_postConstruct []func()
}

var _ BVLPDU = (*DistributeBroadcastToNetwork)(nil)

func NewDistributeBroadcastToNetwork(pdu PDU, opts ...func(*DistributeBroadcastToNetwork)) (*DistributeBroadcastToNetwork, error) {
	o := &DistributeBroadcastToNetwork{}
	for _, opt := range opts {
		opt(o)
	}
	switch npdu := pdu.(type) {
	case readWriteModel.NPDUExactly:
		o._BVLPDU = NewBVLPDU(readWriteModel.NewBVLCDistributeBroadcastToNetwork(o.produceInnerNPDU(npdu))).(*_BVLPDU)
	case nil:
		o._BVLPDU = NewBVLPDU(nil).(*_BVLPDU)
	default:
		// TODO: re-encode seems expensive... check if there is a better option (e.g. only do it on the message bridge)
		data := pdu.GetPduData()
		parse, err := readWriteModel.NPDUParse(context.Background(), data, uint16(len(data)))
		if err != nil {
			return nil, errors.Wrap(err, "error re-encoding")
		}
		o._BVLPDU = NewBVLPDU(readWriteModel.NewBVLCDistributeBroadcastToNetwork(o.produceInnerNPDU(parse))).(*_BVLPDU)
	}
	// Do a post construct for a bit more easy initialization
	for _, f := range o._postConstruct {
		f()
	}
	o._postConstruct = nil
	return o, nil
}

func WithDistributeBroadcastToNetworkDestination(destination *Address) func(*DistributeBroadcastToNetwork) {
	return func(o *DistributeBroadcastToNetwork) {
		o._postConstruct = append(o._postConstruct, func() {
			o.SetPDUDestination(destination)
		})
	}
}

func WithDistributeBroadcastToNetworkUserData(userData spi.Message) func(*DistributeBroadcastToNetwork) {
	return func(o *DistributeBroadcastToNetwork) {
		o._postConstruct = append(o._postConstruct, func() {
			o.SetPDUUserData(userData)
		})
	}
}

func (o *DistributeBroadcastToNetwork) produceInnerNPDU(inNpdu readWriteModel.NPDU) (npdu readWriteModel.NPDU, bvlcPayloadLength uint16) {
	npdu = inNpdu
	return
}

func (o *DistributeBroadcastToNetwork) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := bvlpdu.Update(o); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}

		bvlpdu.PutData(o.GetPduData()...)

		bvlpdu.setBVLC(o.bvlc)
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (o *DistributeBroadcastToNetwork) Decode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := o.Update(bvlpdu); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		switch rm := bvlpdu.GetRootMessage().(type) {
		case readWriteModel.BVLCDistributeBroadcastToNetworkExactly:
			npdu := rm.GetNpdu()
			pduData, err := npdu.Serialize()
			if err != nil {
				return errors.Wrap(err, "error serializing NPDU")
			}
			o.SetPduData(pduData)
			o.setBVLC(rm)
		}
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (o *DistributeBroadcastToNetwork) String() string {
	return fmt.Sprintf("DistributeBroadcastToNetwork{%s}", o._BVLPDU)
}

type OriginalUnicastNPDU struct {
	*_BVLPDU

	// post construct function
	_postConstruct []func()
}

var _ BVLPDU = (*OriginalUnicastNPDU)(nil)

func NewOriginalUnicastNPDU(pdu PDU, opts ...func(*OriginalUnicastNPDU)) (*OriginalUnicastNPDU, error) {
	o := &OriginalUnicastNPDU{}
	for _, opt := range opts {
		opt(o)
	}
	switch npdu := pdu.(type) {
	case readWriteModel.NPDUExactly:
		o._BVLPDU = NewBVLPDU(readWriteModel.NewBVLCOriginalUnicastNPDU(o.produceInnerNPDU(npdu))).(*_BVLPDU)
	case nil:
		o._BVLPDU = NewBVLPDU(nil).(*_BVLPDU)
	default:
		// TODO: re-encode seems expensive... check if there is a better option (e.g. only do it on the message bridge)
		data := pdu.GetPduData()
		parse, err := readWriteModel.NPDUParse(context.Background(), data, uint16(len(data)))
		if err != nil {
			return nil, errors.Wrap(err, "error re-encoding")
		}
		o._BVLPDU = NewBVLPDU(readWriteModel.NewBVLCOriginalUnicastNPDU(o.produceInnerNPDU(parse))).(*_BVLPDU)
	}
	// Do a post construct for a bit more easy initialization
	for _, f := range o._postConstruct {
		f()
	}
	o._postConstruct = nil
	return o, nil
}

func WithOriginalUnicastNPDUDestination(destination *Address) func(*OriginalUnicastNPDU) {
	return func(o *OriginalUnicastNPDU) {
		o._postConstruct = append(o._postConstruct, func() {
			o.SetPDUDestination(destination)
		})
	}
}

func WithOriginalUnicastNPDUUserData(userData spi.Message) func(*OriginalUnicastNPDU) {
	return func(o *OriginalUnicastNPDU) {
		o._postConstruct = append(o._postConstruct, func() {
			o.SetPDUUserData(userData)
		})
	}
}

func (o *OriginalUnicastNPDU) produceInnerNPDU(inNpdu readWriteModel.NPDU) (npdu readWriteModel.NPDU, bvlcPayloadLength uint16) {
	npdu = inNpdu
	return
}

func (o *OriginalUnicastNPDU) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := bvlpdu.Update(o); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}

		bvlpdu.PutData(o.GetPduData()...)

		bvlpdu.setBVLC(o.bvlc)
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (o *OriginalUnicastNPDU) Decode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := o.Update(bvlpdu); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		switch rm := bvlpdu.GetRootMessage().(type) {
		case readWriteModel.BVLCOriginalUnicastNPDUExactly:
			npdu := rm.GetNpdu()
			pduData, err := npdu.Serialize()
			if err != nil {
				return errors.Wrap(err, "error serializing NPDU")
			}
			o.SetPduData(pduData)
			o.setBVLC(rm)
		}
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (o *OriginalUnicastNPDU) String() string {
	return fmt.Sprintf("OriginalUnicastNPDU{%s}", o._BVLPDU)
}

type OriginalBroadcastNPDU struct {
	*_BVLPDU

	// post construct function
	_postConstruct []func()
}

var _ BVLPDU = (*OriginalBroadcastNPDU)(nil)

func NewOriginalBroadcastNPDU(pdu PDU, opts ...func(*OriginalBroadcastNPDU)) (*OriginalBroadcastNPDU, error) {
	o := &OriginalBroadcastNPDU{}
	for _, opt := range opts {
		opt(o)
	}
	switch npdu := pdu.(type) {
	case readWriteModel.NPDUExactly:
		o._BVLPDU = NewBVLPDU(readWriteModel.NewBVLCOriginalBroadcastNPDU(o.produceInnerNPDU(npdu))).(*_BVLPDU)
	case nil:
		o._BVLPDU = NewBVLPDU(nil).(*_BVLPDU)
	default:
		// TODO: re-encode seems expensive... check if there is a better option (e.g. only do it on the message bridge)
		data := pdu.GetPduData()
		parse, err := readWriteModel.NPDUParse(context.Background(), data, uint16(len(data)))
		if err != nil {
			return nil, errors.Wrap(err, "error re-encoding")
		}
		o._BVLPDU = NewBVLPDU(readWriteModel.NewBVLCOriginalBroadcastNPDU(o.produceInnerNPDU(parse))).(*_BVLPDU)
	}
	// Do a post construct for a bit more easy initialization
	for _, f := range o._postConstruct {
		f()
	}
	o._postConstruct = nil
	return o, nil
}

func WithOriginalBroadcastNPDUDestination(destination *Address) func(*OriginalBroadcastNPDU) {
	return func(o *OriginalBroadcastNPDU) {
		o._postConstruct = append(o._postConstruct, func() {
			o.SetPDUDestination(destination)
		})
	}
}

func WithOriginalBroadcastNPDUUserData(userData spi.Message) func(*OriginalBroadcastNPDU) {
	return func(o *OriginalBroadcastNPDU) {
		o._postConstruct = append(o._postConstruct, func() {
			o.SetPDUUserData(userData)
		})
	}
}

func (o *OriginalBroadcastNPDU) produceInnerNPDU(inNpdu readWriteModel.NPDU) (npdu readWriteModel.NPDU, bvlcPayloadLength uint16) {
	npdu = inNpdu
	return
}

func (o *OriginalBroadcastNPDU) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := bvlpdu.Update(o); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}

		bvlpdu.PutData(o.GetPduData()...)

		bvlpdu.setBVLC(o.bvlc)
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (o *OriginalBroadcastNPDU) Decode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := o.Update(bvlpdu); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		switch rm := bvlpdu.GetRootMessage().(type) {
		case readWriteModel.BVLCOriginalBroadcastNPDUExactly:
			npdu := rm.GetNpdu()
			pduData, err := npdu.Serialize()
			if err != nil {
				return errors.Wrap(err, "error serializing NPDU")
			}
			o.SetPduData(pduData)
			o.setBVLC(rm)
		}
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (o *OriginalBroadcastNPDU) String() string {
	return fmt.Sprintf("OriginalBroadcastNPDU{%s}", o._BVLPDU)
}

func init() {
	BVLPDUTypes = map[uint8]func() interface{ Decode(Arg) error }{
		0x00: func() interface{ Decode(Arg) error } {
			v, _ := NewResult()
			return v
		},
		0x01: func() interface{ Decode(Arg) error } {
			v, _ := NewWriteBroadcastDistributionTable()
			return v
		},
		0x02: func() interface{ Decode(Arg) error } {
			v, _ := NewReadBroadcastDistributionTable()
			return v
		},
		0x03: func() interface{ Decode(Arg) error } {
			v, _ := NewReadBroadcastDistributionTableAck()
			return v
		},
		0x04: func() interface{ Decode(Arg) error } {
			v, _ := NewForwardedNPDU(nil)
			return v
		},
		0x05: func() interface{ Decode(Arg) error } {
			v, _ := NewRegisterForeignDevice()
			return v
		},
		0x06: func() interface{ Decode(Arg) error } {
			v, _ := NewReadForeignDeviceTable()
			return v
		},
		0x07: func() interface{ Decode(Arg) error } {
			v, _ := NewReadForeignDeviceTableAck()
			return v
		},
		0x08: func() interface{ Decode(Arg) error } {
			v, _ := NewDeleteForeignDeviceTableEntry()
			return v
		},
		0x09: func() interface{ Decode(Arg) error } {
			v, _ := NewDistributeBroadcastToNetwork(nil)
			return v
		},
		0x0A: func() interface{ Decode(Arg) error } {
			v, _ := NewOriginalUnicastNPDU(nil)
			return v
		},
		0x0B: func() interface{ Decode(Arg) error } {
			v, _ := NewOriginalBroadcastNPDU(nil)
			return v
		},
	}
}
