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

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
)

type APDU interface {
	Copyable
	readWriteModel.APDU
	APCI
	PDUData
}

type __APDU struct {
	*_APCI
	PDUData
}

var _ APDU = (*__APDU)(nil)

func NewAPDU(args Args, kwArgs KWArgs, options ...Option) (APDU, error) {
	if _debug != nil {
		_debug("__init__ %r %r", args, kwArgs)
	}
	a := &__APDU{}
	options = AddLeafTypeIfAbundant(options, a)
	var err error
	a._APCI, err = CreateSharedSuperIfAbundant[_APCI](options, newAPCI, args, kwArgs, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating APCI")
	}
	a.PDUData = NewPDUData(args, kwArgs, options...)
	a.AddExtraPrinters(a.PDUData.(DebugContentPrinter))
	return a, nil
}

func (a *__APDU) Encode(pdu Arg) error {
	if _debug != nil {
		_debug("encode %s", pdu.(fmt.Stringer).String())
	}
	if err := a._APCI.Encode(pdu); err != nil {
		return errors.Wrap(err, "error encoding APCI")
	}
	switch pdu := pdu.(type) {
	case PDUData:
		pdu.PutData(a.GetPduData()...)
	}
	return nil
}

func (a *__APDU) Decode(pdu Arg) error {
	if _debug != nil {
		_debug("decode %s", pdu.(fmt.Stringer).String())
	}
	var rootMessage spi.Message
	switch pdu := pdu.(type) { // Save a root message as long as we have enough data
	case PDUData:
		data := pdu.GetPduData()
		rootMessage, _ = Try1(func() (readWriteModel.APDU, error) {
			return readWriteModel.APDUParse[readWriteModel.APDU](context.Background(), data, uint16(len(data)))
		})
	}
	switch pdu := pdu.(type) {
	case IPCI:
		if rootMessage != nil { // in this case we are good and want to parse from that
			pdu.SetRootMessage(rootMessage)
		}
	}
	if err := a._APCI.Decode(pdu); err != nil {
		return errors.Wrap(err, "error decoding APCI")
	}
	switch pdu := pdu.(type) {
	case PDUData:
		a.PutData(pdu.GetPduData()[a.bytesToDiscard:]...)
	}
	if rootMessage != nil {
		// Overwrite the root message again so we can use it for matching
		a.SetRootMessage(rootMessage)
	}
	return nil
}

func (a *__APDU) CreateAPDUBuilder() readWriteModel.APDUBuilder {
	switch rm := a.GetRootMessage().(type) {
	case readWriteModel.APDU:
		return rm.CreateAPDUBuilder()
	default:
		return readWriteModel.NewAPDUBuilder()
	}
}

func (a *__APDU) GetApduType() readWriteModel.ApduType {
	switch rm := a.GetRootMessage().(type) {
	case readWriteModel.APDU:
		return rm.GetApduType()
	default:
		return *a.apduType
	}
}

func (a *__APDU) GetApduLength() uint16 {
	switch rm := a.GetRootMessage().(type) {
	case readWriteModel.APDU:
		return rm.GetApduLength()
	default:
		return 0 // TODO: what is a good fallback
	}
}

func (a *__APDU) IsAPDU() {
}

func (a *__APDU) deepCopy() *__APDU {
	return &__APDU{_APCI: a._APCI.deepCopy(), PDUData: a.PDUData.DeepCopy().(PDUData)}
}

func (a *__APDU) DeepCopy() any {
	return a.deepCopy()
}

func (a *__APDU) String() string {
	return a._APCI.String()
}
