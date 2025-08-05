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
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
)

// _APDU masks the Encode() and Decode() functions of the APDU
// so that derived structs use the update function to copy the contents
// between PDU,Otherwise the APCI content would be decoded twice.
type _APDU interface {
	APDU
}

type ___APDU struct {
	*__APDU

	_leafName string
}

var _ _APDU = (*___APDU)(nil)

func New_APDU(args Args, kwArgs KWArgs, options ...Option) (_APDU, error) {
	i := &___APDU{
		_leafName: ExtractLeafName(options, StructName()),
	}
	options = AddLeafTypeIfAbundant(options, i)
	var err error
	apdu, err := NewAPDU(args, kwArgs, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating APDU")
	}
	i.__APDU = apdu.(*__APDU)
	return i, nil
}

func (a *___APDU) Encode(pdu Arg) error {
	if _debug != nil {
		_debug("encode %r", pdu)
	}
	switch pdu := pdu.(type) {
	case PCI:
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
	if _debug != nil {
		_debug("decode %r", pdu)
	}
	if err := a._APCI.Update(pdu); err != nil {
		return errors.Wrap(err, "error updating pdu")
	}
	switch pdu := pdu.(type) {
	case PDUData:
		data, err := pdu.GetData(len(pdu.GetPduData()))
		if err != nil {
			return errors.Wrap(err, "error getting data")
		}
		a.SetPduData(data)
	}
	return nil
}

func (a *___APDU) SetContext(context APDU) {
	if _debug != nil {
		_debug("set_context %r", context)
	}
	a.SetPDUUserData(context.GetPDUUserData())
	a.SetPDUDestination(context.GetPDUSource())
	a.SetExpectingReply(false)
	a.SetNetworkPriority(context.GetNetworkPriority())
	a.apduInvokeID = context.GetApduInvokeID()
}

func (a *___APDU) Format(s fmt.State, v rune) {
	switch v {
	case 'v', 's', 'r':
		sname := a._leafName

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
		_, _ = fmt.Fprintf(s, "<%s(%s) instance at %p>\n", sname, stype, a)
	}
	a.PrintDebugContents(2, s, nil)
}
