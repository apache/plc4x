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
	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type ErrorPDU struct {
	*___APDU

	bacnetError readWriteModel.BACnetError
}

var _ readWriteModel.APDUError = (*ErrorPDU)(nil)

func NewErrorPDU(args Args, kwArgs KWArgs, options ...Option) (*ErrorPDU, error) {
	e := &ErrorPDU{}
	choice, ok := KWO[*readWriteModel.BACnetConfirmedServiceChoice](kwArgs, KWConfirmedServiceChoice, nil)
	invokeID, ok := KWO[*uint8](kwArgs, KWInvokedID, nil)
	context, ok := KWO[APDU](kwArgs, KWContext, nil)
	if _debug != nil {
		_debug("__init__ %r %r %r %r %r", choice, invokeID, context, args, kwArgs)
	}
	options = AddLeafTypeIfAbundant(options, e)
	apdu, err := New_APDU(args, kwArgs, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating _APDU")
	}
	e.___APDU = apdu.(*___APDU)
	e.apduType = ToPtr(readWriteModel.ApduType_ERROR_PDU)
	if ok {
		serviceChoice := uint8(*choice)
		e.apduService = &serviceChoice
	}
	e.apduInvokeID = invokeID

	// use the context to fill in most of the fields
	if context != nil {
		e.apduService = context.getApduService()
		e.SetContext(context)
	}

	switch rm := e.GetRootMessage().(type) {
	case readWriteModel.BACnetError:
		e.bacnetError = rm
		serviceChoice := rm.GetErrorChoice()
		e.apduService = ToPtr(uint8(serviceChoice))
	}
	e.SetRootMessage(e.buildConfirmedRequest(e.bacnetError))

	return e, nil
}

func (e *ErrorPDU) buildConfirmedRequest(bacnetError readWriteModel.BACnetError) readWriteModel.APDUError {
	if bacnetError == nil {
		return nil
	}
	invokeID := uint8(0)
	if e.apduInvokeID != nil {
		invokeID = *e.apduInvokeID
	}
	return readWriteModel.NewAPDUError(
		invokeID,
		e.GetErrorChoice(),
		e.bacnetError,
		0,
	)
}

func (e *ErrorPDU) CreateAPDUErrorBuilder() readWriteModel.APDUErrorBuilder {
	//TODO implement me
	panic("implement me")
}

func (e *ErrorPDU) GetSegmentedMessage() bool {
	return e.apduSeg
}

func (e *ErrorPDU) GetMoreFollows() bool {
	return e.apduMor
}

func (e *ErrorPDU) GetSegmentedResponseAccepted() bool {
	return e.apduSA
}

func (e *ErrorPDU) GetMaxSegmentsAccepted() readWriteModel.MaxSegmentsAccepted {
	if e.apduMaxSegs != nil {
		return readWriteModel.MaxSegmentsAccepted(*e.apduMaxSegs)
	}
	return readWriteModel.MaxSegmentsAccepted_MORE_THAN_64_SEGMENTS
}

func (e *ErrorPDU) GetMaxApduLengthAccepted() readWriteModel.MaxApduLengthAccepted {
	if e.apduMaxResp != nil {
		return readWriteModel.MaxApduLengthAccepted(*e.apduMaxResp)
	}
	return readWriteModel.MaxApduLengthAccepted_NUM_OCTETS_1476
}

func (e *ErrorPDU) GetInvokeId() uint8 {
	if e.apduInvokeID != nil {
		return *e.apduInvokeID
	}
	return 0
}

func (e *ErrorPDU) GetSequenceNumber() *uint8 {
	return e.apduSeq
}

func (e *ErrorPDU) GetProposedWindowSize() *uint8 {
	return e.apduWin
}

func (e *ErrorPDU) GetSegment() []byte {
	return nil
}

func (e *ErrorPDU) GetApduHeaderReduction() uint16 {
	return 0
}

func (e *ErrorPDU) GetSegmentReduction() uint16 {
	return 0
}

func (e *ErrorPDU) GetOriginalInvokeId() uint8 {
	invokeID := e.apduInvokeID
	if invokeID != nil {
		return *invokeID
	}
	return 0
}

func (e *ErrorPDU) GetErrorChoice() readWriteModel.BACnetConfirmedServiceChoice {
	service := e.apduService
	if service != nil {
		return readWriteModel.BACnetConfirmedServiceChoice(*service)
	}
	return 0
}

func (e *ErrorPDU) GetError() readWriteModel.BACnetError {
	return e.bacnetError
}

func (e *ErrorPDU) IsAPDUError() {
}
