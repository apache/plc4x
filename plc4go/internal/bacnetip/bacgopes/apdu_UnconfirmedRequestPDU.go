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

package bacgopes

import (
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type UnconfirmedRequestPDU struct {
	*___APDU

	serviceRequest readWriteModel.BACnetUnconfirmedServiceRequest
}

var _ readWriteModel.APDUUnconfirmedRequest = (*UnconfirmedRequestPDU)(nil)

func NewUnconfirmedRequestPDU(serviceRequest readWriteModel.BACnetUnconfirmedServiceRequest, opts ...func(*UnconfirmedRequestPDU)) (*UnconfirmedRequestPDU, error) {
	u := &UnconfirmedRequestPDU{
		serviceRequest: serviceRequest,
	}
	for _, opt := range opts {
		opt(u)
	}
	apdu, _ := new_APDU(buildUnconfirmedServiceRequest(serviceRequest))
	u.___APDU = apdu.(*___APDU)
	if serviceRequest != nil {
		serviceChoice := uint8(serviceRequest.GetServiceChoice())
		u.apduService = &serviceChoice
	}
	return u, nil
}

func buildUnconfirmedServiceRequest(serviceRequest readWriteModel.BACnetUnconfirmedServiceRequest) readWriteModel.APDUUnconfirmedRequest {
	if serviceRequest == nil {
		return nil
	}
	return readWriteModel.NewAPDUUnconfirmedRequest(serviceRequest, 0)
}

func (u *UnconfirmedRequestPDU) GetServiceRequest() readWriteModel.BACnetUnconfirmedServiceRequest {
	return u.serviceRequest
}

func (u *UnconfirmedRequestPDU) IsAPDUUnconfirmedRequest() {
}
