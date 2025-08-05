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
type SimpleAckPDU struct {
	*___APDU
}

var _ readWriteModel.APDUComplexAck = (*SimpleAckPDU)(nil)

func (s *SimpleAckPDU) CreateAPDUComplexAckBuilder() readWriteModel.APDUComplexAckBuilder {
	panic("implement me")
}

func (s *SimpleAckPDU) GetSegmentedMessage() bool {
	//TODO implement me
	panic("implement me")
}

func (s *SimpleAckPDU) GetMoreFollows() bool {
	//TODO implement me
	panic("implement me")
}

func (s *SimpleAckPDU) GetOriginalInvokeId() uint8 {
	//TODO implement me
	panic("implement me")
}

func (s *SimpleAckPDU) GetSequenceNumber() *uint8 {
	//TODO implement me
	panic("implement me")
}

func (s *SimpleAckPDU) GetProposedWindowSize() *uint8 {
	//TODO implement me
	panic("implement me")
}

func (s *SimpleAckPDU) GetServiceAck() readWriteModel.BACnetServiceAck {
	//TODO implement me
	panic("implement me")
}

func (s *SimpleAckPDU) GetSegmentServiceChoice() *readWriteModel.BACnetConfirmedServiceChoice {
	//TODO implement me
	panic("implement me")
}

func (s *SimpleAckPDU) GetSegment() []byte {
	//TODO implement me
	panic("implement me")
}

func (s *SimpleAckPDU) GetApduHeaderReduction() uint16 {
	//TODO implement me
	panic("implement me")
}

func (s *SimpleAckPDU) GetSegmentReduction() uint16 {
	//TODO implement me
	panic("implement me")
}

func (s *SimpleAckPDU) IsAPDUComplexAck() {
}
