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
	"fmt"

	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type APDU interface {
	model.APDU
	APCI
	PDUData
}

type __APDU struct {
	*_APCI
	*_PDUData
}

var _ APDU = (*__APDU)(nil)

func NewAPDU() (APDU, error) {
	a := &__APDU{}

	a._APCI = NewAPCI(nil).(*_APCI)
	a._PDUData = NewPDUData(NoArgs).(*_PDUData)
	return a, nil
}

func (a *__APDU) Encode(pdu Arg) error {
	panic("implement me")
}

func (a *__APDU) Decode(pdu Arg) error {
	panic("implement me")
}

func (a *__APDU) GetApduLength() uint16 {
	//TODO implement me
	panic("implement me")
}

func (a *__APDU) GetApduType() model.ApduType {
	panic("implement me")
}

func (a *__APDU) deepCopy() *__APDU {
	return &__APDU{_APCI: a._APCI.deepCopy(), _PDUData: a._PDUData.deepCopy()}
}

func (a *__APDU) DeepCopy() any {
	return a.deepCopy()
}

func (a *__APDU) String() string {
	return fmt.Sprintf("__APDU{%s}", a._PCI)
}

func (a *__APDU) IsAPDU() {
	//TODO implement me
	panic("implement me")
}
