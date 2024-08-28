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

type NPCI interface {
	PCI

	GetNPDUNetMessage() *uint8

	Encode(pdu Arg) error
	Decode(pdu Arg) error

	setNpduVersion(uint8)
	getNpduVersion() uint8
	setNpduControl(uint8)
	getNpduControl() uint8
	setNpduDADR(*Address)
	getNpduDADR() *Address
	setNpduSADR(*Address)
	getNpduSADR() *Address
	setNpduHopCount(*uint8)
	getNpduHopCount() *uint8
	setNpduNetMessage(*uint8)
	getNpduNetMessage() *uint8
	setNpduVendorID(*uint16)
	getNpduVendorID() *uint16
	setNLM(model.NLM)
	getNLM() model.NLM
	setAPDU(model.APDU)
	getAPDU() model.APDU
}

type _NPCI struct {
	*_PCI
	*DebugContents

	npduVersion    uint8
	npduControl    uint8
	npduDADR       *Address
	npduSADR       *Address
	npduHopCount   *uint8
	npduNetMessage *uint8
	npduVendorID   *uint16

	npdu model.NPDU
	nlm  model.NLM
	apdu model.APDU
}

var _ NPCI = (*_NPCI)(nil)

func NewNPCI(pduUserData spi.Message, nlm model.NLM, apdu model.APDU) NPCI {
	n := &_NPCI{
		nlm:  nlm,
		apdu: apdu,

		npduVersion: 1,
	}
	n._PCI = newPCI(pduUserData, nil, nil, nil, false, model.NPDUNetworkPriority_NORMAL_MESSAGE)
	switch ud := pduUserData.(type) {
	case model.NLMExactly:
		n.nlm = ud
	case model.APDUExactly:
		n.apdu = ud
	}
	return n
}

func (n *_NPCI) GetNPDUNetMessage() *uint8 {
	if n.nlm == nil {
		return nil
	}
	messageType := n.nlm.GetMessageType()
	return &messageType
}

func (n *_NPCI) setNpduVersion(u uint8) {
	n.npduVersion = u
}

func (n *_NPCI) getNpduVersion() uint8 {
	return n.npduVersion
}

func (n *_NPCI) setNpduControl(u uint8) {
	n.npduControl = u
}

func (n *_NPCI) getNpduControl() uint8 {
	return n.npduControl
}

func (n *_NPCI) setNpduDADR(address *Address) {
	n.npduDADR = address
}

func (n *_NPCI) getNpduDADR() *Address {
	return n.npduDADR
}

func (n *_NPCI) setNpduSADR(address *Address) {
	n.npduSADR = address
}

func (n *_NPCI) getNpduSADR() *Address {
	return n.npduSADR
}

func (n *_NPCI) setNpduHopCount(u *uint8) {
	n.npduHopCount = u
}

func (n *_NPCI) getNpduHopCount() *uint8 {
	return n.npduHopCount
}

func (n *_NPCI) setNpduNetMessage(u *uint8) {
	n.npduNetMessage = u
}

func (n *_NPCI) getNpduNetMessage() *uint8 {
	return n.npduNetMessage
}

func (n *_NPCI) setNpduVendorID(u *uint16) {
	n.npduVendorID = u
}

func (n *_NPCI) getNpduVendorID() *uint16 {
	return n.npduVendorID
}

func (n *_NPCI) setNLM(nlm model.NLM) {
	n.nlm = nlm
}

func (n *_NPCI) getNLM() model.NLM {
	return n.nlm
}

func (n *_NPCI) setAPDU(apdu model.APDU) {
	n.apdu = apdu
}

func (n *_NPCI) getAPDU() model.APDU {
	return n.apdu
}

func (n *_NPCI) Update(npci Arg) error {
	if err := n._PCI.Update(npci); err != nil {
		return errors.Wrap(err, "error updating _PCI")
	}
	switch npci := npci.(type) {
	case NPCI:
		n.npduVersion = npci.getNpduVersion()
		n.npduControl = npci.getNpduControl()
		n.npduDADR = npci.getNpduDADR()
		n.npduSADR = npci.getNpduSADR()
		n.npduHopCount = npci.getNpduHopCount()
		n.npduNetMessage = npci.getNpduNetMessage()
		n.npduVendorID = npci.getNpduVendorID()

		n.nlm = npci.getNLM()
		n.apdu = npci.getAPDU()
		return nil
	default:
		return errors.Errorf("invalid NPCI type %T", npci)
	}
}

func (n *_NPCI) Encode(pdu Arg) error {
	if err := pdu.(interface{ Update(Arg) error }).Update(n); err != nil { // TODO: better validate that arg is really PDUData... use switch similar to Update
		return errors.Wrap(err, "error updating pdu")
	}
	switch pdu := pdu.(type) {
	case NPCI:
		pdu.setNLM(n.nlm)
	}
	return nil
}

func (n *_NPCI) Decode(pdu Arg) error {
	if err := n._PCI.Update(pdu); err != nil {
		return errors.Wrap(err, "error updating pdu")
	}
	switch rm := n.rootMessage.(type) {
	case model.NPDUExactly:
		n.nlm = rm.GetNlm()
	}
	return nil
}

func (n *_NPCI) deepCopy() *_NPCI {
	return &_NPCI{_PCI: n._PCI.deepCopy(), nlm: n.nlm}
}
