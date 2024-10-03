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

package apdu

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
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
	SetApduInvokeID(*uint8)
	GetApduInvokeID() *uint8
	setApduAbortRejectReason(*uint8)
	getApduAbortRejectReason() *uint8

	getAPCI() APCI
}

type _APCI struct {
	PCI
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

	// Deprecated: hacky workaround
	bytesToDiscard int

	_leafName string
}

var _ APCI = (*_APCI)(nil)

func NewAPCI(args Args, kwArgs KWArgs, options ...Option) (APCI, error) {
	return newAPCI(args, kwArgs, options...)
}

func newAPCI(args Args, kwArgs KWArgs, options ...Option) (*_APCI, error) {
	if _debug != nil {
		_debug("__init__ %r %r", args, kwArgs)
	}
	a := &_APCI{
		_leafName: ExtractLeafName(options, StructName()),
	}
	options = AddLeafTypeIfAbundant(options, a)
	a.DebugContents = NewDebugContents(a, "apduType", "apduSeg", "apduMor", "apduSA", "apduSrv",
		"apduNak", "apduSeq", "apduWin", "apduMaxSegs", "apduMaxResp",
		"apduService", "apduInvokeID", "apduAbortRejectReason")
	a.PCI = NewPCI(args, kwArgs, options...)
	a.AddExtraPrinters(a.PCI.(DebugContentPrinter))
	return a, nil
}

func (a *_APCI) GetDebugAttr(attr string) any {
	switch attr {
	case "apduType":
		if a.apduType != nil {
			return *a.apduType
		}
	case "apduSeg":
		if a.apduSeg {
			return a.apduSeq
		}
	case "apduMor":
		if a.apduMor {
			return a.apduMor
		}
	case "apduSA":
		if a.apduSA {
			return a.apduSA
		}
	case "apduSrv":
		if a.apduSrv {
			return a.apduSrv
		}
	case "apduNak":
		if a.apduNak {
			return a.apduNak
		}
	case "apduSeq":
		if a.apduSeq != nil {
			return *a.apduSeq
		}
	case "apduWin":
		if a.apduWin != nil {
			return *a.apduWin
		}
	case "apduMaxSegs":
		if a.apduMaxSegs != nil {
			return *a.apduMaxSegs
		}
	case "apduMaxResp":
		if a.apduMaxResp != nil {
			return *a.apduMaxResp
		}
	case "apduService":
		if a.apduService != nil {
			return *a.apduService
		}
	case "apduInvokeID":
		if a.apduInvokeID != nil {
			return *a.apduInvokeID
		}
	case "apduAbortRejectReason":
		if a.apduAbortRejectReason != nil {
			return *a.apduAbortRejectReason
		}
	default:
		return nil
	}
	return nil
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
func (a *_APCI) SetApduInvokeID(apduInvokeID *uint8) {
	a.apduInvokeID = apduInvokeID
}
func (a *_APCI) GetApduInvokeID() *uint8 {
	return a.apduInvokeID
}
func (a *_APCI) setApduAbortRejectReason(apduAbortRejectReason *uint8) {
	a.apduAbortRejectReason = apduAbortRejectReason
}
func (a *_APCI) getApduAbortRejectReason() *uint8 {
	return a.apduAbortRejectReason
}

func (a *_APCI) getAPCI() APCI {
	return a
}

func (a *_APCI) Update(apci Arg) error {
	if err := a.PCI.Update(apci); err != nil {
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
		a.apduInvokeID = apci.GetApduInvokeID()
		a.apduAbortRejectReason = apci.getApduAbortRejectReason()
		return nil
	default:
		return errors.Errorf("invalid APCI type %T", apci)
	}
}

func (a *_APCI) Encode(pdu Arg) error {
	if _debug != nil {
		_debug("encode %r", pdu)
	}
	switch pdu := pdu.(type) {
	case PCI:
		if err := pdu.GetPCI().Update(a); err != nil {
			return errors.Wrap(err, "error updating pdu")
		}
	}
	switch pdu := pdu.(type) {
	case PDUData:
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
	if _debug != nil {
		_debug("decode %r", pdu)
	}
	if err := a.PCI.Update(pdu); err != nil {
		return errors.Wrap(err, "error updating pdu")
	}
	// TODO: check if we want to stay with parsing or ditch that for now
	data := pdu.(PDUData).GetPduData()
	parse, err := readWriteModel.APDUParse[readWriteModel.APDU](context.Background(), data, uint16(len(data)))
	if err != nil {
		return errors.Wrap(err, "error parsing apdu")
	}
	a.SetRootMessage(parse)
	readBytes := 0

	// simulate the APCI type decode
	readBytes++
	switch rm := a.GetRootMessage().(type) {
	case readWriteModel.APDU:
		apduType := rm.GetApduType()
		a.apduType = &apduType
	}
	switch rm := a.GetRootMessage().(type) {
	case readWriteModel.APDUConfirmedRequest:
		a.apduSeg = rm.GetSegmentedMessage()
		a.apduMor = rm.GetMoreFollows()
		a.apduSA = rm.GetSegmentedResponseAccepted()
		readBytes++
		maxSegmentsAccepted := uint8(rm.GetMaxSegmentsAccepted())
		a.apduMaxSegs = &maxSegmentsAccepted
		maxApduLengthAccepted := uint8(rm.GetMaxApduLengthAccepted())
		a.apduMaxResp = &maxApduLengthAccepted
		readBytes++
		invokeId := rm.GetInvokeId()
		a.apduInvokeID = &invokeId
		if a.apduSeg {
			readBytes++
			a.apduSeq = rm.GetSequenceNumber()
			readBytes++
			a.apduWin = rm.GetProposedWindowSize()
		}
		readBytes++
		if sr := rm.GetServiceRequest(); sr != nil {
			serviceChoice := uint8(sr.GetServiceChoice())
			a.apduService = &serviceChoice
		}
	case readWriteModel.APDUUnconfirmedRequest:
		readBytes++
		var choice *uint8
		if sr := rm.GetServiceRequest(); sr != nil {
			serviceChoice := uint8(sr.GetServiceChoice())
			choice = &serviceChoice
		}
		a.apduService = choice
	case readWriteModel.APDUSimpleAck:
		readBytes++
		invokeId := rm.GetOriginalInvokeId()
		a.apduInvokeID = &invokeId
	case readWriteModel.APDUComplexAck:
		a.apduSeg = rm.GetSegmentedMessage()
		a.apduMor = rm.GetMoreFollows()
		readBytes++
		invokeId := rm.GetOriginalInvokeId()
		a.apduInvokeID = &invokeId
		if a.apduSeg {
			readBytes++
			a.apduSeq = rm.GetSequenceNumber()
			readBytes++
			a.apduWin = rm.GetProposedWindowSize()
		}
		if sr := rm.GetServiceAck(); sr != nil {
			serviceChoice := uint8(sr.GetServiceChoice())
			a.apduService = &serviceChoice
		}
	case readWriteModel.APDUSegmentAck:
		a.apduNak = rm.GetNegativeAck()
		a.apduSrv = rm.GetServer()
		readBytes++
		invokeId := rm.GetOriginalInvokeId()
		a.apduInvokeID = &invokeId
		readBytes++
		actualWindowSize := rm.GetActualWindowSize()
		a.apduWin = &actualWindowSize
	case readWriteModel.APDUError:
		readBytes++
		invokeId := rm.GetOriginalInvokeId()
		a.apduInvokeID = &invokeId
		serviceChoice := uint8(rm.GetErrorChoice())
		a.apduService = &serviceChoice
	case readWriteModel.APDUReject:
		readBytes++
		invokeId := rm.GetOriginalInvokeId()
		a.apduInvokeID = &invokeId
		apduAbortRejectReason := uint8(rm.GetRejectReason().GetValue())
		a.apduAbortRejectReason = &apduAbortRejectReason
	case readWriteModel.APDUAbort:
		readBytes++
		invokeId := rm.GetOriginalInvokeId()
		a.apduInvokeID = &invokeId
		apduAbortRejectReason := uint8(rm.GetAbortReason().GetValue())
		a.apduAbortRejectReason = &apduAbortRejectReason
	}
	a.bytesToDiscard = readBytes
	return nil
}

func (a *_APCI) deepCopy() *_APCI {
	newA := &_APCI{
		a.PCI.DeepCopy().(PCI),
		nil,
		CopyPtr(a.apduType),
		a.apduSeg,
		a.apduMor,
		a.apduSA,
		a.apduSrv,
		a.apduNak,
		CopyPtr(a.apduSeq),
		CopyPtr(a.apduWin),
		CopyPtr(a.apduMaxSegs),
		CopyPtr(a.apduMaxResp),
		CopyPtr(a.apduService),
		CopyPtr(a.apduInvokeID),
		CopyPtr(a.apduAbortRejectReason),
		a.bytesToDiscard,
		a._leafName,
	}
	newA.DebugContents = NewDebugContents(newA, "apduType", "apduSeg", "apduMor", "apduSA", "apduSrv",
		"apduNak", "apduSeq", "apduWin", "apduMaxSegs", "apduMaxResp",
		"apduService", "apduInvokeID", "apduAbortRejectReason") // TODO: bit ugly to repeat that here again but what are the options...
	newA.AddExtraPrinters(newA.PCI.(DebugContentPrinter))
	return newA
}

func (a *_APCI) String() string {
	sname := a._leafName

	// expand the type if possible
	stype := ""
	if a.apduType != nil {
		if v, ok := APDUTypes[*a.apduType]; ok {
			stype = fmt.Sprintf("%s(%d)", QualifiedTypeName(v()), *a.apduType)
		} else {
			stype = "?"
		}
	}

	// add the invoke ID if it has one
	if a.apduInvokeID != nil {
		stype += ", " + strconv.Itoa(int(*a.apduInvokeID))
	}
	// put it together
	return fmt.Sprintf("<%s(%s instance at %p)>", sname, stype, a)
}
