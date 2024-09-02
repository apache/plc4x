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

	"github.com/pkg/errors"

	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type NPDU interface {
	model.NPDU
	NPCI
	PDUData

	DeepCopy() any
}

type _NPDU struct {
	*_NPCI
	*_PDUData
}

var _ NPDU = (*_NPDU)(nil)

// TODO: optimize with options and smart non-recoding...
func NewNPDU(nlm model.NLM, apdu model.APDU) (NPDU, error) {
	n := &_NPDU{}
	npdu, _ := n.buildNPDU(0, nil, nil, false, model.NPDUNetworkPriority_NORMAL_MESSAGE, nlm, apdu)
	n._NPCI = NewNPCI(npdu, nlm, apdu).(*_NPCI)
	n._PDUData = NewPDUData(NoArgs).(*_PDUData)
	if npdu != nil {
		n.data, _ = npdu.Serialize()
	}
	return n, nil
}

func (n *_NPDU) Encode(pdu Arg) error {
	if err := n._NPCI.Encode(pdu); err != nil {
		return errors.Wrap(err, "error encoding _NPCI")
	}
	switch pdu := pdu.(type) {
	case PDU:
		pdu.PutData(n.GetPduData()...)
	}
	return nil
}

func (n *_NPDU) Decode(pdu Arg) error {
	if err := n._NPCI.Decode(pdu); err != nil {
		return errors.Wrap(err, "error decoding _NPCI")
	}
	switch pdu := pdu.(type) {
	case PDUData:
		n.PutData(pdu.GetPduData()...)
	}
	return nil
}

func (n *_NPDU) getNPDUModel() (model.NPDU, bool) {
	rm := n.GetRootMessage()
	npdu, ok := rm.(model.NPDU)
	return npdu, ok
}

func (n *_NPDU) GetProtocolVersionNumber() uint8 {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return 0
	}
	return npdu.GetProtocolVersionNumber()
}

func (n *_NPDU) GetControl() model.NPDUControl {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return nil
	}
	return npdu.GetControl()
}

func (n *_NPDU) GetDestinationNetworkAddress() *uint16 {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return nil
	}
	return npdu.GetDestinationNetworkAddress()
}

func (n *_NPDU) GetDestinationLength() *uint8 {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return nil
	}
	return npdu.GetDestinationLength()
}

func (n *_NPDU) GetDestinationAddress() []uint8 {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return nil
	}
	return npdu.GetDestinationAddress()
}

func (n *_NPDU) GetSourceNetworkAddress() *uint16 {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return nil
	}
	return npdu.GetSourceNetworkAddress()
}

func (n *_NPDU) GetSourceLength() *uint8 {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return nil
	}
	return npdu.GetSourceLength()
}

func (n *_NPDU) GetSourceAddress() []uint8 {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return nil
	}
	return npdu.GetSourceAddress()
}

func (n *_NPDU) GetHopCount() *uint8 {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return nil
	}
	return npdu.GetHopCount()
}

func (n *_NPDU) GetNlm() model.NLM {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return nil
	}
	return npdu.GetNlm()
}

func (n *_NPDU) GetApdu() model.APDU {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return nil
	}
	return npdu.GetApdu()
}

func (n *_NPDU) GetDestinationLengthAddon() uint16 {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return 0
	}
	return npdu.GetDestinationLengthAddon()
}

func (n *_NPDU) GetSourceLengthAddon() uint16 {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return 0
	}
	return npdu.GetSourceLengthAddon()
}

func (n *_NPDU) GetPayloadSubtraction() uint16 {
	npdu, ok := n.getNPDUModel()
	if !ok {
		return 0
	}
	return npdu.GetPayloadSubtraction()
}

func (n *_NPDU) IsNPDU() {
}

func (n *_NPDU) deepCopy() *_NPDU {
	return &_NPDU{_NPCI: n._NPCI.deepCopy(), _PDUData: n._PDUData.deepCopy()}
}

func (n *_NPDU) DeepCopy() any {
	return n.deepCopy()
}

func (n *_NPDU) String() string {
	return fmt.Sprintf("_NPDU{%s}", n._PCI)
}
