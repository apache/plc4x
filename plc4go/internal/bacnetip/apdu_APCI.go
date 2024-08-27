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
	"github.com/pkg/errors"

	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
)

type APCI interface {
	PCI

	GetApduInvokeID() *uint8

	Encode(pdu Arg) error
	Decode(pdu Arg) error

	setAPDU(model.APDU)
	getAPDU() model.APDU
}

type _APCI struct {
	*_PCI
	*DebugContents

	apdu model.APDU // TODO: check if this is part of the _APCI or _APDU
}

var _ APCI = (*_APCI)(nil)

func NewAPCI(pduUserData spi.Message, apdu model.APDU) APCI {
	a := &_APCI{
		apdu: apdu,
	}
	a._PCI = newPCI(pduUserData, nil, nil, nil, false, model.NPDUNetworkPriority_NORMAL_MESSAGE)
	switch apdu := pduUserData.(type) {
	case model.APDUExactly:
		a.apdu = apdu
	}
	return a
}

func (n *_APCI) GetApduInvokeID() *uint8 {
	if n.apdu == nil {
		return nil
	}
	switch apdu := n.apdu.(type) {
	case model.APDUConfirmedRequestExactly:
		invokeId := apdu.GetInvokeId()
		return &invokeId
	default:
		return nil
	}
}

// Deprecated: check if needed as we do it in update
func (n *_APCI) setAPDU(apdu model.APDU) {
	n.apdu = apdu
}

func (n *_APCI) getAPDU() model.APDU {
	return n.apdu
}

func (n *_APCI) Update(apci Arg) error {
	if err := n._PCI.Update(apci); err != nil {
		return errors.Wrap(err, "error updating _PCI")
	}
	switch apci := apci.(type) {
	case APDU:
		n.apdu = apci.getAPDU()
		// TODO: update coordinates...
		return nil
	default:
		return errors.Errorf("invalid APCI type %T", apci)
	}
}

func (n *_APCI) Encode(pdu Arg) error {
	if err := pdu.(interface{ Update(Arg) error }).Update(n); err != nil { // TODO: better validate that arg is really PDUData... use switch similar to Update
		return errors.Wrap(err, "error updating pdu")
	}
	// TODO: what should we do here??
	return nil
}

func (n *_APCI) Decode(pdu Arg) error {
	if err := n._PCI.Update(pdu); err != nil {
		return errors.Wrap(err, "error updating pdu")
	}
	// TODO: what should we do here??
	return nil
}

func (n *_APCI) deepCopy() *_APCI {
	return &_APCI{_PCI: n._PCI.deepCopy(), apdu: n.apdu}
}
