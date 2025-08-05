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

type ConfirmedRequestPDU struct {
	*___APDU

	serviceRequest readWriteModel.BACnetConfirmedServiceRequest
}

var _ readWriteModel.APDUConfirmedRequest = (*ConfirmedRequestPDU)(nil)

func NewConfirmedRequestPDU(args Args, kwArgs KWArgs, options ...Option) (*ConfirmedRequestPDU, error) {
	c := &ConfirmedRequestPDU{}
	choice, ok := KWO[*readWriteModel.BACnetConfirmedServiceChoice](kwArgs, KWConfirmedServiceChoice, nil)
	if _debug != nil {
		_debug("__init__ %r %r %r", choice, args, kwArgs)
	}
	apdu, err := New_APDU(args, kwArgs, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating _APDU")
	}
	c.___APDU = apdu.(*___APDU)
	c.apduType = ToPtr(readWriteModel.ApduType_CONFIRMED_REQUEST_PDU)
	if ok {
		serviceChoice := uint8(*choice)
		c.apduService = &serviceChoice
	}
	c.SetExpectingReply(true)
	switch rm := c.GetRootMessage().(type) {
	case readWriteModel.BACnetConfirmedServiceRequest:
		c.serviceRequest = rm
		serviceChoice := rm.GetServiceChoice()
		c.apduService = ToPtr(uint8(serviceChoice))
	}
	c.SetRootMessage(c.buildConfirmedRequest(c.serviceRequest))

	return c, nil
}

func (c *ConfirmedRequestPDU) buildConfirmedRequest(serviceRequest readWriteModel.BACnetConfirmedServiceRequest) readWriteModel.APDUConfirmedRequest {
	if serviceRequest == nil {
		return nil
	}
	invokeID := uint8(0)
	if c.apduInvokeID != nil {
		invokeID = *c.apduInvokeID
	}
	return readWriteModel.NewAPDUConfirmedRequest(
		c.apduSeg,
		c.apduMor,
		c.apduSA,
		readWriteModel.MaxSegmentsAccepted_MORE_THAN_64_SEGMENTS,
		readWriteModel.MaxApduLengthAccepted_NUM_OCTETS_1476,
		invokeID,
		c.apduSeq,
		c.apduWin,
		serviceRequest,
		nil, // TODO: where to get from
		nil, // TODO: where to get from
		0,
	)
}

func (c *ConfirmedRequestPDU) CreateAPDUConfirmedRequestBuilder() readWriteModel.APDUConfirmedRequestBuilder {
	//TODO implement me
	panic("implement me")
}

func (c *ConfirmedRequestPDU) GetSegmentedMessage() bool {
	return c.apduSeg
}

func (c *ConfirmedRequestPDU) GetMoreFollows() bool {
	return c.apduMor
}

func (c *ConfirmedRequestPDU) GetSegmentedResponseAccepted() bool {
	return c.apduSA
}

func (c *ConfirmedRequestPDU) GetMaxSegmentsAccepted() readWriteModel.MaxSegmentsAccepted {
	if c.apduMaxSegs != nil {
		return readWriteModel.MaxSegmentsAccepted(*c.apduMaxSegs)
	}
	return readWriteModel.MaxSegmentsAccepted_MORE_THAN_64_SEGMENTS
}

func (c *ConfirmedRequestPDU) GetMaxApduLengthAccepted() readWriteModel.MaxApduLengthAccepted {
	if c.apduMaxResp != nil {
		return readWriteModel.MaxApduLengthAccepted(*c.apduMaxResp)
	}
	return readWriteModel.MaxApduLengthAccepted_NUM_OCTETS_1476
}

func (c *ConfirmedRequestPDU) GetInvokeId() uint8 {
	if c.apduInvokeID != nil {
		return *c.apduInvokeID
	}
	return 0
}

func (c *ConfirmedRequestPDU) GetSequenceNumber() *uint8 {
	return c.apduSeq
}

func (c *ConfirmedRequestPDU) GetProposedWindowSize() *uint8 {
	return c.apduWin
}

func (c *ConfirmedRequestPDU) GetServiceRequest() readWriteModel.BACnetConfirmedServiceRequest {
	return c.serviceRequest
}

func (c *ConfirmedRequestPDU) GetSegmentServiceChoice() *readWriteModel.BACnetConfirmedServiceChoice {
	serviceChoice := c.serviceRequest.GetServiceChoice()
	return &serviceChoice
}

func (c *ConfirmedRequestPDU) GetSegment() []byte {
	return nil
}

func (c *ConfirmedRequestPDU) GetApduHeaderReduction() uint16 {
	return 0
}

func (c *ConfirmedRequestPDU) GetSegmentReduction() uint16 {
	return 0
}

func (c *ConfirmedRequestPDU) IsAPDUConfirmedRequest() {
}
