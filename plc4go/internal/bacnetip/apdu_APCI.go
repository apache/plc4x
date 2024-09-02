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

	"github.com/pkg/errors"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
)

type APCI interface {
	PCI

	Encode(pdu Arg) error
	Decode(pdu Arg) error

	setApduType(*readWriteModel.ApduType)
	getApduType() *readWriteModel.ApduType
	setApduSeg(bool)
	getApduSeg() bool
	setApduMor(bool)
	getApduMor() bool
	setApduSA(bool)
	getApduSA() bool
	setApduSrv(bool)
	getApduSrv() bool
	setApduNak(bool)
	getApduNak() bool
	setApduSeq(*uint8)
	getApduSeq() *uint8
	setApduWin(*uint8)
	getApduWin() *uint8
	setApduMaxSegs(*uint8)
	getApduMaxSegs() *uint8
	setApduMaxResp(*uint8)
	getApduMaxResp() *uint8
	setApduService(*uint8)
	getApduService() *uint8
	setApduInvokeID(*uint8)
	getApduInvokeID() *uint8
	setApduAbortRejectReason(*uint8)
	getApduAbortRejectReason() *uint8
}

type _APCI struct {
	*_PCI
	*DebugContents

	apduType              *readWriteModel.ApduType
	apduSeg               bool   // segmented
	apduMor               bool   // more follows
	apduSA                bool   // segmented response accepted
	apduSrv               bool   // sent by server
	apduNak               bool   // negative acknowledgement
	apduSeq               *uint8 // sequence number
	apduWin               *uint8 // actual/proposed window size
	apduMaxSegs           *uint8 // maximum segments accepted (decoded)
	apduMaxResp           *uint8 // max response accepted (decoded)
	apduService           *uint8
	apduInvokeID          *uint8
	apduAbortRejectReason *uint8

	// TODO: check if we store plc4go payload here
}

var _ APCI = (*_APCI)(nil)

func NewAPCI(pduUserData spi.Message) APCI {
	a := &_APCI{}
	a._PCI = newPCI(pduUserData, nil, nil, nil, false, readWriteModel.NPDUNetworkPriority_NORMAL_MESSAGE)
	switch ud := pduUserData.(type) {
	case readWriteModel.APDU:
		_ = ud // TODO: WIP
	}
	return a
}

func (a *_APCI) setApduType(apduType *readWriteModel.ApduType) {
	a.apduType = apduType
}
func (a *_APCI) getApduType() *readWriteModel.ApduType {
	return a.apduType
}
func (a *_APCI) setApduSeg(apduSeg bool) {
	a.apduSeg = apduSeg
}
func (a *_APCI) getApduSeg() bool {
	return a.apduSeg
}
func (a *_APCI) setApduMor(apduMor bool) {
	a.apduMor = apduMor
}
func (a *_APCI) getApduMor() bool {
	return a.apduMor
}
func (a *_APCI) setApduSA(apduSA bool) {
	a.apduSA = apduSA
}
func (a *_APCI) getApduSA() bool {
	return a.apduSA
}
func (a *_APCI) setApduSrv(apduSrv bool) {
	a.apduSrv = apduSrv
}
func (a *_APCI) getApduSrv() bool {
	return a.apduSrv
}
func (a *_APCI) setApduNak(apduNak bool) {
	a.apduNak = apduNak
}
func (a *_APCI) getApduNak() bool {
	return a.apduNak
}
func (a *_APCI) setApduSeq(apduSeq *uint8) {
	a.apduSeq = apduSeq
}
func (a *_APCI) getApduSeq() *uint8 {
	return a.apduSeq
}
func (a *_APCI) setApduWin(apduWin *uint8) {
	a.apduWin = apduWin
}
func (a *_APCI) getApduWin() *uint8 {
	return a.apduWin
}
func (a *_APCI) setApduMaxSegs(apduMaxSegs *uint8) {
	a.apduMaxSegs = apduMaxSegs
}
func (a *_APCI) getApduMaxSegs() *uint8 {
	return a.apduMaxSegs
}
func (a *_APCI) setApduMaxResp(apduMaxResp *uint8) {
	a.apduMaxResp = apduMaxResp
}
func (a *_APCI) getApduMaxResp() *uint8 {
	return a.apduMaxResp
}
func (a *_APCI) setApduService(apduService *uint8) {
	a.apduService = apduService
}
func (a *_APCI) getApduService() *uint8 {
	return a.apduService
}
func (a *_APCI) setApduInvokeID(apduInvokeID *uint8) {
	a.apduInvokeID = apduInvokeID
}
func (a *_APCI) getApduInvokeID() *uint8 {
	return a.apduInvokeID
}
func (a *_APCI) setApduAbortRejectReason(apduAbortRejectReason *uint8) {
	a.apduAbortRejectReason = apduAbortRejectReason
}
func (a *_APCI) getApduAbortRejectReason() *uint8 {
	return a.apduAbortRejectReason
}

func (a *_APCI) Update(apci Arg) error {
	if err := a._PCI.Update(apci); err != nil {
		return errors.Wrap(err, "error updating _PCI")
	}
	switch apci := apci.(type) {
	case APCI:
		a.apduType = apci.getApduType()
		a.apduSeg = apci.getApduSeg()
		a.apduMor = apci.getApduMor()
		a.apduSA = apci.getApduSA()
		a.apduSrv = apci.getApduSrv()
		a.apduNak = apci.getApduNak()
		a.apduSeq = apci.getApduSeq()
		a.apduWin = apci.getApduWin()
		a.apduMaxSegs = apci.getApduMaxSegs()
		a.apduMaxResp = apci.getApduMaxResp()
		a.apduService = apci.getApduService()
		a.apduInvokeID = apci.getApduInvokeID()
		a.apduAbortRejectReason = apci.getApduAbortRejectReason()
		return nil
	default:
		return errors.Errorf("invalid APCI type %T", apci)
	}
}

func (a *_APCI) Encode(pdu Arg) error {
	switch pdu := pdu.(type) {
	case PDU:
		if err := pdu.(interface{ Update(Arg) error }).Update(a); err != nil { // TODO: better validate that arg is really PDUData... use switch similar to Update
			return errors.Wrap(err, "error updating pdu")
		}

		if a.apduType == nil {
			return errors.New("APCI does not have APDU type")
		}
		switch *a.apduType {
		case readWriteModel.ApduType_CONFIRMED_REQUEST_PDU:
			// PDU type
			buff := byte(*a.apduType << 4)
			if a.apduSeg {
				buff += 0x08
			}
			if a.apduMor {
				buff += 0x04
			}
			if a.apduSA {
				buff += 0x02
			}
			pdu.Put(buff)
			pdu.Put((*a.apduMaxSegs << 4) + *a.apduMaxResp)
			pdu.Put(*a.apduInvokeID)
			if a.apduSeg {
				pdu.Put(*a.apduSeq)
				pdu.Put(*a.apduWin)
			}
			pdu.Put(byte(*a.apduService))

		case readWriteModel.ApduType_UNCONFIRMED_REQUEST_PDU:
			pdu.Put(uint8(*a.apduType) << 4)
			pdu.Put(uint8(*a.apduService))

		case readWriteModel.ApduType_SIMPLE_ACK_PDU:
			pdu.Put(uint8(*a.apduType) << 4)
			pdu.Put(*a.apduInvokeID)
			pdu.Put(uint8(*a.apduService))

		case readWriteModel.ApduType_COMPLEX_ACK_PDU:
			// PDU type
			buff := uint8(*a.apduType << 4)
			if a.apduSeg {
				buff += 0x08
			}
			if a.apduMor {
				buff += 0x04
			}
			pdu.Put(buff)
			pdu.Put(*a.apduInvokeID)
			if a.apduSeg {
				pdu.Put(*a.apduSeq)
				pdu.Put(*a.apduWin)
			}
			pdu.Put(uint8(*a.apduService))

		case readWriteModel.ApduType_SEGMENT_ACK_PDU:
			// PDU type
			buff := uint8(*a.apduType << 4)
			if a.apduNak {
				buff += 0x02
			}
			if a.apduSrv {
				buff += 0x01
			}
			pdu.Put(buff)
			pdu.Put(*a.apduInvokeID)
			pdu.Put(*a.apduSeq)
			pdu.Put(*a.apduWin)

		case readWriteModel.ApduType_ERROR_PDU:
			pdu.Put(uint8(*a.apduType << 4))
			pdu.Put(*a.apduInvokeID)
			pdu.Put(uint8(*a.apduService))

		case readWriteModel.ApduType_REJECT_PDU:
			pdu.Put(uint8(*a.apduType << 4))
			pdu.Put(*a.apduInvokeID)
			pdu.Put(*a.apduAbortRejectReason)

		case readWriteModel.ApduType_ABORT_PDU:
			// PDU type
			buff := uint8(*a.apduType << 4)
			if a.apduSrv {
				buff += 0x01
			}
			pdu.Put(buff)
			pdu.Put(*a.apduInvokeID)
			pdu.Put(*a.apduAbortRejectReason)
		}
	}
	return nil
}

func (a *_APCI) Decode(pdu Arg) error {
	if err := a._PCI.Update(pdu); err != nil {
		return errors.Wrap(err, "error updating pdu")
	}
	switch rm := a.rootMessage.(type) {
	case *messageBridge:
		data := rm.GetPduData()
		parse, err := readWriteModel.APDUParse[readWriteModel.APDU](context.Background(), data, uint16(len(data)))
		if err != nil {
			return errors.Wrap(err, "error parsing apdu")
		}
		a.rootMessage = parse
	}
	switch rm := a.rootMessage.(type) {
	case readWriteModel.APDU:
		_ = rm
		panic("implement me")
	}
	return nil
}

func (a *_APCI) deepCopy() *_APCI {
	// TODO: fix deepcopy
	return &_APCI{
		_PCI:                  a._PCI.deepCopy(),
		apduType:              a.apduType,
		apduSeg:               a.apduSeg,
		apduMor:               a.apduMor,
		apduSA:                a.apduSA,
		apduSrv:               a.apduSrv,
		apduNak:               a.apduNak,
		apduSeq:               a.apduSeq,
		apduWin:               a.apduWin,
		apduMaxSegs:           a.apduMaxSegs,
		apduMaxResp:           a.apduMaxResp,
		apduService:           a.apduService,
		apduInvokeID:          a.apduInvokeID,
		apduAbortRejectReason: a.apduAbortRejectReason,
	}
}

func (a *_APCI) String() string {
	return ""
}
