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

import readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"

// TODO: implement it...
type ComplexAckPDU struct {
	*___APDU
}

var _ readWriteModel.APDUComplexAck = (*ComplexAckPDU)(nil)

func (c *ComplexAckPDU) CreateAPDUComplexAckBuilder() readWriteModel.APDUComplexAckBuilder {
	panic("implement me")
}

func (c *ComplexAckPDU) GetSegmentedMessage() bool {
	//TODO implement me
	panic("implement me")
}

func (c *ComplexAckPDU) GetMoreFollows() bool {
	//TODO implement me
	panic("implement me")
}

func (c *ComplexAckPDU) GetOriginalInvokeId() uint8 {
	//TODO implement me
	panic("implement me")
}

func (c *ComplexAckPDU) GetSequenceNumber() *uint8 {
	//TODO implement me
	panic("implement me")
}

func (c *ComplexAckPDU) GetProposedWindowSize() *uint8 {
	//TODO implement me
	panic("implement me")
}

func (c *ComplexAckPDU) GetServiceAck() readWriteModel.BACnetServiceAck {
	//TODO implement me
	panic("implement me")
}

func (c *ComplexAckPDU) GetSegmentServiceChoice() *readWriteModel.BACnetConfirmedServiceChoice {
	//TODO implement me
	panic("implement me")
}

func (c *ComplexAckPDU) GetSegment() []byte {
	//TODO implement me
	panic("implement me")
}

func (c *ComplexAckPDU) GetApduHeaderReduction() uint16 {
	//TODO implement me
	panic("implement me")
}

func (c *ComplexAckPDU) GetSegmentReduction() uint16 {
	//TODO implement me
	panic("implement me")
}

func (c *ComplexAckPDU) IsAPDUComplexAck() {
	//TODO implement me
	panic("implement me")
}
