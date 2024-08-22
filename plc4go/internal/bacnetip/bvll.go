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
	"fmt"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/pkg/errors"
)

// BCLPDUTypes is a dictionary of message type values and structs
var BCLPDUTypes map[uint8]func() interface{ Decode(Arg) error }

type BVLCI interface {
	PCI
	Update(bvlci Arg) error
	Encode(pdu Arg) error
	Decode(pdu Arg) error

	setBVLC(apdu readWriteModel.BVLC)
}

// BVLCIContract provides a set of functions which can be overwritten by a sub struct
type BVLCIContract interface {
}

// BVLCIRequirements provides a set of functions which need to be overwritten by a sub struct
type BVLCIRequirements interface {
	BVLCIContract
}

type _BVLCI struct {
	*_PCI
	*DebugContents
	BVLCIContract
	requirements BVLCIRequirements

	bvlc readWriteModel.BVLC
}

var _ BVLCI = (*_BVLCI)(nil)

func NewBVLCI(pduUserData spi.Message, requirements BVLCIRequirements) BVLCI {
	b := &_BVLCI{
		requirements:  requirements,
		BVLCIContract: requirements,
	}
	b._PCI = newPCI(pduUserData, nil, nil, false, readWriteModel.NPDUNetworkPriority_NORMAL_MESSAGE)
	return b
}

func (b *_BVLCI) setBVLC(bvlc readWriteModel.BVLC) {
	b.bvlc = bvlc
}

func (b *_BVLCI) Update(bvlci Arg) error {
	if err := b._PCI.Update(bvlci); err != nil {
		return errors.Wrap(err, "Update BVLCI")
	}
	// TODO
	return nil
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

func (b *_BVLCI) GetMessage() spi.Message {
	return b.bvlc
}

func (b *_BVLCI) getPDUData() []byte {
	if b.GetMessage() == nil {
		return nil
	}
	writeBufferByteBased := utils.NewWriteBufferByteBased()
	if err := b.GetMessage().SerializeWithWriteBuffer(context.Background(), writeBufferByteBased); err != nil {
		panic(err) // TODO: graceful handle
	}
	return writeBufferByteBased.GetBytes()
}

func (b *_BVLCI) deepCopy() *_BVLCI {
	return &_BVLCI{_PCI: b._PCI.deepCopy()}
}

type BVLPDU interface {
	BVLCI
	PDUData
}

type _BVLPDU struct {
	*_BVLCI
	*_PDUData
}

var _ BVLPDU = (*_BVLPDU)(nil)

func NewBVLPDU(bvlc readWriteModel.BVLC) BVLPDU {
	b := &_BVLPDU{}
	b._BVLCI = NewBVLCI(bvlc, b).(*_BVLCI)
	return b
}

func (b *_BVLPDU) Encode(pdu Arg) error {
	if err := b._BVLCI.Encode(pdu); err != nil {
		return errors.Wrap(err, "error encoding _BVLCI")
	}
	b.SetPDUUserData(b.bvlc)
	return nil
}

func (b *_BVLPDU) Decode(pdu Arg) error {
	if err := b._BVLCI.Decode(pdu); err != nil {
		return errors.Wrap(err, "error decoding _BVLCI")
	}
	switch pdu := pdu.(type) {
	case PDUData:
		data := pdu.GetPduData()
		var err error
		b.bvlc, err = readWriteModel.BVLCParse(context.Background(), data)
		if err != nil {
			return errors.Wrap(err, "error parsing NPDU")
		}
		b.pduUserData = b.bvlc
	}
	return nil
}

func (b *_BVLPDU) deepCopy() *_BVLPDU {
	return &_BVLPDU{_BVLCI: b._BVLCI.deepCopy(), _PDUData: b._PDUData.deepCopy()}
}

func (b *_BVLPDU) DeepCopy() PDU {
	return b.deepCopy()
}

type Result struct {
	*_BVLPDU

	bvlciResultCode readWriteModel.BVLCResultCode
}

var _ BVLPDU = (*Result)(nil)

func NewResult() (BVLPDU, error) {
	b := &Result{}
	b._BVLPDU = NewBVLPDU(readWriteModel.NewBVLCResult(0)).(*_BVLPDU)
	return b, nil
}

func WithResultBvlciResultCode(code readWriteModel.BVLCResultCode) func(*Result) {
	return func(b *Result) {
		b.bvlciResultCode = code
	}
}

func (n *Result) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := bvlpdu.Update(n); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		bvlpdu.PutShort(int16(n.bvlciResultCode))
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
		switch pduUserData := bvlpdu.GetPDUUserData().(type) {
		case readWriteModel.BVLCExactly:
			switch bvlc := pduUserData.(type) {
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
	return fmt.Sprintf("Result{%s, bvlciResultCode: %v}", n._BVLPDU, n.bvlciResultCode)
}

// TODO: finish
type WriteBroadcastDistributionTable struct {
	*_BVLPDU
}

var _ BVLPDU = (*WriteBroadcastDistributionTable)(nil)

func NewWriteBroadcastDistributionTable() (BVLPDU, error) {
	b := &WriteBroadcastDistributionTable{}
	b._BVLPDU = NewBVLPDU(nil).(*_BVLPDU)
	return b, nil
}

func (b *WriteBroadcastDistributionTable) Encode(pdu Arg) error {
	// TODO: finish
	return nil
}

func (b *WriteBroadcastDistributionTable) Decode(pdu Arg) error {
	// TODO: finish
	return nil
}

// TODO: finish
type ReadBroadcastDistributionTable struct {
	*_BVLPDU
}

var _ BVLPDU = (*ReadBroadcastDistributionTable)(nil)

func NewReadBroadcastDistributionTable() (BVLPDU, error) {
	b := &ReadBroadcastDistributionTable{}
	b._BVLPDU = NewBVLPDU(nil).(*_BVLPDU)
	return b, nil
}

func (b *ReadBroadcastDistributionTable) Encode(pdu Arg) error {
	// TODO: finish
	return nil
}

func (b *ReadBroadcastDistributionTable) Decode(pdu Arg) error {
	// TODO: finish
	return nil
}

// TODO: finish
type ReadBroadcastDistributionTableAck struct {
	*_BVLPDU
}

var _ BVLPDU = (*ReadBroadcastDistributionTableAck)(nil)

func NewReadBroadcastDistributionTableAck() (BVLPDU, error) {
	b := &ReadBroadcastDistributionTableAck{}
	b._BVLPDU = NewBVLPDU(nil).(*_BVLPDU)
	return b, nil
}

func (b *ReadBroadcastDistributionTableAck) Encode(pdu Arg) error {
	// TODO: finish
	return nil
}

func (b *ReadBroadcastDistributionTableAck) Decode(pdu Arg) error {
	// TODO: finish
	return nil
}

// TODO: finish
type ForwardedNPDU struct {
	*_BVLPDU
}

var _ BVLPDU = (*ForwardedNPDU)(nil)

func NewForwardedNPDU() (BVLPDU, error) {
	b := &ForwardedNPDU{}
	b._BVLPDU = NewBVLPDU(nil).(*_BVLPDU)
	return b, nil
}

func (b *ForwardedNPDU) Encode(pdu Arg) error {
	// TODO: finish
	return nil
}

func (b *ForwardedNPDU) Decode(pdu Arg) error {
	// TODO: finish
	return nil
}

// TODO: finish
type RegisterForeignDevice struct {
	*_BVLPDU
}

var _ BVLPDU = (*RegisterForeignDevice)(nil)

func NewRegisterForeignDevice() (BVLPDU, error) {
	b := &RegisterForeignDevice{}
	b._BVLPDU = NewBVLPDU(nil).(*_BVLPDU)
	return b, nil
}

func (b *RegisterForeignDevice) Encode(pdu Arg) error {
	// TODO: finish
	return nil
}

func (b *RegisterForeignDevice) Decode(pdu Arg) error {
	// TODO: finish
	return nil
}

// TODO: finish
type ReadForeignDeviceTable struct {
	*_BVLPDU
}

var _ BVLPDU = (*ReadForeignDeviceTable)(nil)

func NewReadForeignDeviceTable() (BVLPDU, error) {
	b := &ReadForeignDeviceTable{}
	b._BVLPDU = NewBVLPDU(nil).(*_BVLPDU)
	return b, nil
}

func (b *ReadForeignDeviceTable) Encode(pdu Arg) error {
	// TODO: finish
	return nil
}

func (b *ReadForeignDeviceTable) Decode(pdu Arg) error {
	// TODO: finish
	return nil
}

// TODO: finish
type ReadForeignDeviceTableAck struct {
	*_BVLPDU
}

var _ BVLPDU = (*ReadForeignDeviceTableAck)(nil)

func NewReadForeignDeviceTableAck() (BVLPDU, error) {
	b := &ReadForeignDeviceTableAck{}
	b._BVLPDU = NewBVLPDU(nil).(*_BVLPDU)
	return b, nil
}

func (b *ReadForeignDeviceTableAck) Encode(pdu Arg) error {
	// TODO: finish
	return nil
}

func (b *ReadForeignDeviceTableAck) Decode(pdu Arg) error {
	// TODO: finish
	return nil
}

// TODO: finish
type DeleteForeignDeviceTableEntry struct {
	*_BVLPDU
}

var _ BVLPDU = (*DeleteForeignDeviceTableEntry)(nil)

func NewDeleteForeignDeviceTableEntry() (BVLPDU, error) {
	b := &DeleteForeignDeviceTableEntry{}
	b._BVLPDU = NewBVLPDU(nil).(*_BVLPDU)
	return b, nil
}

func (b *DeleteForeignDeviceTableEntry) Encode(pdu Arg) error {
	// TODO: finish
	return nil
}

func (b *DeleteForeignDeviceTableEntry) Decode(pdu Arg) error {
	// TODO: finish
	return nil
}

// TODO: finish
type DistributeBroadcastToNetwork struct {
	*_BVLPDU
}

var _ BVLPDU = (*DistributeBroadcastToNetwork)(nil)

func NewDistributeBroadcastToNetwork() (BVLPDU, error) {
	b := &DistributeBroadcastToNetwork{}
	b._BVLPDU = NewBVLPDU(nil).(*_BVLPDU)
	return b, nil
}

func (b *DistributeBroadcastToNetwork) Encode(pdu Arg) error {
	// TODO: finish
	return nil
}

func (b *DistributeBroadcastToNetwork) Decode(pdu Arg) error {
	// TODO: finish
	return nil
}

type OriginalUnicastNPDU struct {
	*_BVLPDU
}

var _ BVLPDU = (*OriginalUnicastNPDU)(nil)

func NewOriginalUnicastNPDU(npdu NPDU, opts ...func(*OriginalUnicastNPDU)) (BVLPDU, error) {
	b := &OriginalUnicastNPDU{}
	for _, opt := range opts {
		opt(b)
	}
	b._BVLPDU = NewBVLPDU(readWriteModel.NewBVLCOriginalUnicastNPDU(npdu, npdu.GetLengthInBytes(context.Background()))).(*_BVLPDU)
	return b, nil
}

func WithOriginalUnicastNPDUDestination(destination *Address) func(*OriginalUnicastNPDU) {
	return func(o *OriginalUnicastNPDU) {
		o.pduDestination = destination
	}
}

func WithOriginalUnicastNPDUUserData(userData spi.Message) func(*OriginalUnicastNPDU) {
	return func(o *OriginalUnicastNPDU) {
		o.pduUserData = userData
	}
}

func (n *OriginalUnicastNPDU) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := bvlpdu.Update(n); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		bvlpdu.setBVLC(n.bvlc)
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (n *OriginalUnicastNPDU) Decode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := n.Update(bvlpdu); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		switch pduUserData := bvlpdu.GetPDUUserData().(type) {
		case readWriteModel.BVLCExactly:
			switch bvlc := pduUserData.(type) {
			case readWriteModel.BVLCOriginalUnicastNPDU:
				n.setBVLC(bvlc)
			}
		}
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (n *OriginalUnicastNPDU) String() string {
	return fmt.Sprintf("OriginalUnicastNPDU{%s}", n._BVLPDU)
}

type OriginalBroadcastNPDU struct {
	*_BVLPDU
}

func NewOriginalBroadcastNPDU(npdu NPDU, opts ...func(*OriginalBroadcastNPDU)) (BVLPDU, error) {
	b := &OriginalBroadcastNPDU{}
	for _, opt := range opts {
		opt(b)
	}
	b._BVLPDU = NewBVLPDU(readWriteModel.NewBVLCOriginalBroadcastNPDU(npdu, npdu.GetLengthInBytes(context.Background()))).(*_BVLPDU)
	return b, nil
}

func WithOriginalBroadcastNPDUDestination(destination *Address) func(*OriginalBroadcastNPDU) {
	return func(o *OriginalBroadcastNPDU) {
		o.pduDestination = destination
	}
}

func WithOriginalBroadcastNPDUUserData(userData spi.Message) func(*OriginalBroadcastNPDU) {
	return func(o *OriginalBroadcastNPDU) {
		o.pduUserData = userData
	}
}

func (n *OriginalBroadcastNPDU) Encode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := bvlpdu.Update(n); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		bvlpdu.setBVLC(n.bvlc)
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (n *OriginalBroadcastNPDU) Decode(bvlpdu Arg) error {
	switch bvlpdu := bvlpdu.(type) {
	case BVLPDU:
		if err := n.Update(bvlpdu); err != nil {
			return errors.Wrap(err, "error updating BVLPDU")
		}
		switch pduUserData := bvlpdu.GetPDUUserData().(type) {
		case readWriteModel.BVLCExactly:
			switch bvlc := pduUserData.(type) {
			case readWriteModel.BVLCOriginalBroadcastNPDU:
				n.setBVLC(bvlc)
			}
		}
		return nil
	default:
		return errors.Errorf("invalid BVLPDU type %T", bvlpdu)
	}
}

func (n *OriginalBroadcastNPDU) String() string {
	return fmt.Sprintf("OriginalBroadcastNPDU{%s}", n._BVLPDU)
}

func init() {
	BCLPDUTypes = map[uint8]func() interface{ Decode(Arg) error }{
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
			v, _ := NewForwardedNPDU()
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
			v, _ := NewDistributeBroadcastToNetwork()
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
