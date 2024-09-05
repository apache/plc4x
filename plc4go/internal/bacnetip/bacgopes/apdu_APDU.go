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
	"context"
	"fmt"

	"github.com/pkg/errors"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
)

type APDU interface {
	readWriteModel.APDU
	APCI
	PDUData
}

type __APDU struct {
	*_APCI
	*_PDUData

	// post construct function
	_postConstruct []func()
}

var _ APDU = (*__APDU)(nil)

// TODO: optimize with options and smart non-recoding...
func NewAPDU(apdu readWriteModel.APDU, opts ...func(*__APDU)) (APDU, error) {
	a := &__APDU{}
	for _, opt := range opts {
		opt(a)
	}
	a._APCI = NewAPCI(apdu).(*_APCI)
	a._PDUData = NewPDUData(NoArgs).(*_PDUData)
	// Do a post construct for a bit more easy initialization
	for _, f := range a._postConstruct {
		f()
	}
	a._postConstruct = nil
	if a.rootMessage != nil {
		a.data, _ = a.rootMessage.Serialize()
	}
	return a, nil
}

func WithAPDUUserData(userData spi.Message) func(*__APDU) {
	return func(apdu *__APDU) {
		apdu._postConstruct = append(apdu._postConstruct, func() {
			apdu.pduUserData = userData
		})
	}
}

func (a *__APDU) Encode(pdu Arg) error {
	if err := a._APCI.Encode(pdu); err != nil {
		return errors.Wrap(err, "error encoding APCI")
	}
	switch pdu := pdu.(type) {
	case PDUData:
		pdu.PutData(a.data...)
	}
	return nil
}

func (a *__APDU) Decode(pdu Arg) error {
	var rootMessage spi.Message
	switch pdu := pdu.(type) { // Save a root message as long as we have enough data
	case PDUData:
		data := pdu.GetPduData()
		rootMessage, _ = readWriteModel.APDUParse[readWriteModel.APDU](context.Background(), data, uint16(len(data)))
	}
	if err := a._APCI.Decode(pdu); err != nil {
		return errors.Wrap(err, "error decoding APCI")
	}
	switch pdu := pdu.(type) {
	case PDUData:
		a.PutData(pdu.GetPduData()...)
	}
	if rootMessage != nil {
		// Overwrite the root message again so we can use it for matching
		a.rootMessage = rootMessage
	}
	return nil
}

func (a *__APDU) GetApduType() readWriteModel.ApduType {
	switch rm := a.rootMessage.(type) {
	case readWriteModel.APDU:
		return rm.GetApduType()
	default:
		return 0
	}
}

func (a *__APDU) GetApduLength() uint16 {
	switch rm := a.rootMessage.(type) {
	case readWriteModel.APDU:
		return rm.GetApduLength()
	default:
		return 0
	}
}

func (a *__APDU) IsAPDU() {
}

func (a *__APDU) deepCopy() *__APDU {
	return &__APDU{_APCI: a._APCI.deepCopy(), _PDUData: a._PDUData.deepCopy()}
}

func (a *__APDU) DeepCopy() any {
	return a.deepCopy()
}

func (a *__APDU) String() string {
	return fmt.Sprintf("APDU{%s}", a._PCI)
}
