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
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/globals"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

// _APDU masks the Encode() and Decode() functions of the APDU
// so that derived structs use the update function to copy the contents
// between PDU,Otherwise the APCI content would be decoded twice.
type _APDU interface {
	APDU
}

type ___APDU struct {
	*__APDU
}

var _ _APDU = (*___APDU)(nil)

func new_APDU(rootMessage readWriteModel.APDU, opts ...func(*___APDU)) (_APDU, error) {
	i := &___APDU{}
	for _, opt := range opts {
		opt(i)
	}
	var err error
	apdu, err := NewAPDU(rootMessage)
	if err != nil {
		return nil, errors.Wrap(err, "error creating APDU")
	}
	i.__APDU = apdu.(*__APDU)
	return i, nil
}

func (a *___APDU) Encode(pdu Arg) error {
	switch pdu := pdu.(type) {
	case APCI:
		if err := pdu.Update(a); err != nil {
			return errors.Wrap(err, "error updating PDU")
		}
	}
	switch pdu := pdu.(type) {
	case PDUData:
		pdu.SetPduData(a.GetPduData())
	}
	return nil
}

func (a *___APDU) Decode(pdu Arg) error {
	if err := a._APCI.Update(pdu); err != nil {
		return errors.Wrap(err, "error updating pdu")
	}
	a.SetPduData(pdu.(PDUData).GetPduData())
	return nil
}

func (a *___APDU) String() string {
	if ExtendedPDUOutput {
		return fmt.Sprintf("_APDU{%s}", a.__APDU)
	} else {
		sname := fmt.Sprintf("%T", a)

		// the type is the service
		stype := ""
		if a.apduService != nil {
			stype = strconv.Itoa(int(*a.apduService))
		} else {
			stype = "?"
		}

		// add the invoke ID if it has one
		if a.apduInvokeID != nil {
			stype += ", " + strconv.Itoa(int(*a.apduInvokeID))
		}
		// put it together
		return fmt.Sprintf("<%s(%s) instance at %p>", sname, stype, a)
	}
}
