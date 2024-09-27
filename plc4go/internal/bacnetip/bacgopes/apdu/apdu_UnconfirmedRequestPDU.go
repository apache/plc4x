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

type UnconfirmedRequestPDU struct {
	*___APDU

	serviceRequest readWriteModel.BACnetUnconfirmedServiceRequest
}

var _ readWriteModel.APDUUnconfirmedRequest = (*UnconfirmedRequestPDU)(nil)

func NewUnconfirmedRequestPDU(args Args, kwArgs KWArgs, options ...Option) (*UnconfirmedRequestPDU, error) {
	u := &UnconfirmedRequestPDU{}
	choice, ok := KWO[*readWriteModel.BACnetUnconfirmedServiceChoice](kwArgs, KWUnconfirmedServiceChoice, nil)
	if _debug != nil {
		_debug("__init__ %r %r %r", choice, args, kwArgs)
	}
	options = AddLeafTypeIfAbundant(options, u)
	apdu, err := New_APDU(args, kwArgs, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating _APDU")
	}
	u.___APDU = apdu.(*___APDU)
	u.apduType = ToPtr(readWriteModel.ApduType_UNCONFIRMED_REQUEST_PDU)
	if ok {
		serviceChoice := uint8(*choice)
		u.apduService = &serviceChoice
	}
	switch rm := u.GetRootMessage().(type) {
	case readWriteModel.BACnetUnconfirmedServiceRequest:
		u.serviceRequest = rm
		serviceChoice := rm.GetServiceChoice()
		u.apduService = ToPtr(uint8(serviceChoice))
	}
	u.SetRootMessage(u.buildUnconfirmedServiceRequest(u.serviceRequest))
	return u, nil
}

func (u *UnconfirmedRequestPDU) buildUnconfirmedServiceRequest(serviceRequest readWriteModel.BACnetUnconfirmedServiceRequest) readWriteModel.APDUUnconfirmedRequest {
	if serviceRequest == nil {
		return nil
	}
	return readWriteModel.NewAPDUUnconfirmedRequest(serviceRequest, 0)
}

func (u *UnconfirmedRequestPDU) CreateAPDUUnconfirmedRequestBuilder() readWriteModel.APDUUnconfirmedRequestBuilder {
	panic("implement me")
}

func (u *UnconfirmedRequestPDU) GetServiceRequest() readWriteModel.BACnetUnconfirmedServiceRequest {
	return u.serviceRequest
}

func (u *UnconfirmedRequestPDU) IsAPDUUnconfirmedRequest() {
}
