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
	"context"
	"fmt"
	"math"

	"github.com/pkg/errors"

	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
)

type NPDU interface {
	model.NPDU
	NPCI
	PDUData

	setNPDU(npdu model.NPDU)
	setAPDU(apdu model.APDU)
}

type _NPDU struct {
	*_NPCI
	*_PDUData

	// Deprecated: this is the rootMessage so no need to store it here
	npdu model.NPDU
	// Deprecated: this is the rootMessage so no need to store it here
	apdu model.APDU
}

var _ NPDU = (*_NPDU)(nil)

func NewNPDU(nlm model.NLM, apdu model.APDU) (NPDU, error) {
	n := &_NPDU{
		apdu: apdu,
	}
	var err error
	n.npdu, err = n.buildNPDU(0, nil, nil, false, model.NPDUNetworkPriority_NORMAL_MESSAGE, nlm, apdu)
	if err != nil {
		// TODO: seems to be a legit case, which means we don't build our npdu yet... check that
		//return nil, errors.Wrap(err, "error building NPDU")
	}
	n._NPCI = NewNPCI(n.npdu, nlm).(*_NPCI)
	n._PDUData = NewPDUData(NoArgs).(*_PDUData)
	return n, nil
}

// Deprecated: this is the rootMessage so no need to store it here
func (n *_NPDU) setNPDU(npdu model.NPDU) {
	n.npdu = npdu
}

func (n *_NPDU) setAPDU(apdu model.APDU) {
	n.apdu = apdu
}

func (n *_NPDU) buildNPDU(hopCount uint8, source *Address, destination *Address, expectingReply bool, networkPriority model.NPDUNetworkPriority, nlm model.NLM, apdu model.APDU) (model.NPDU, error) {
	switch {
	case nlm != nil && apdu != nil:
		return nil, errors.New("either specify a NLM or a APDU exclusive")
	case nlm == nil && apdu == nil:
		return nil, errors.New("either specify a NLM or a APDU")
	}
	sourceSpecified := source != nil
	var sourceNetworkAddress *uint16
	var sourceLength *uint8
	var sourceAddress []uint8
	if sourceSpecified {
		sourceSpecified = true
		sourceNetworkAddress = source.AddrNet
		sourceLengthValue := *source.AddrLen
		if sourceLengthValue > math.MaxUint8 {
			return nil, errors.New("source address length overflows")
		}
		sourceLengthValueUint8 := sourceLengthValue
		sourceLength = &sourceLengthValueUint8
		sourceAddress = source.AddrAddress
		if sourceLengthValueUint8 == 0 {
			// If we define the len 0 we must not send the array
			sourceAddress = nil
		}
	}
	destinationSpecified := destination != nil && destination.AddrType != LOCAL_BROADCAST_ADDRESS // TODO: check if this is right... (exclude local broadcast)
	var destinationNetworkAddress *uint16
	var destinationLength *uint8
	var destinationAddress []uint8
	var destinationHopCount *uint8
	if destinationSpecified {
		destinationSpecified = true
		destinationNetworkAddress = destination.AddrNet
		destinationLengthValue := *destination.AddrLen
		if destinationLengthValue > math.MaxUint8 {
			return nil, errors.New("source address length overflows")
		}
		destinationLengthValueUint8 := destinationLengthValue
		destinationLength = &destinationLengthValueUint8
		destinationAddress = destination.AddrAddress
		if destinationLengthValueUint8 == 0 {
			// If we define the len 0 we must not send the array
			destinationAddress = nil
		}
		destinationHopCount = &hopCount
	}
	control := model.NewNPDUControl(nlm != nil, destinationSpecified, sourceSpecified, expectingReply, networkPriority)
	return model.NewNPDU(1, control, destinationNetworkAddress, destinationLength, destinationAddress, sourceNetworkAddress, sourceLength, sourceAddress, destinationHopCount, nlm, apdu, 0), nil
}

func (n *_NPDU) Encode(pdu Arg) error {
	if err := n._NPCI.Encode(pdu); err != nil {
		return errors.Wrap(err, "error encoding _NPCI")
	}
	var err error
	n.npdu, err = n.buildNPDU(0, n.GetPDUSource(), n.GetPDUDestination(), n.GetExpectingReply(), n.GetNetworkPriority(), n.nlm, n.apdu)
	if err != nil {
		return errors.Wrap(err, "error building NPDU")
	}
	serialize, err := n.npdu.Serialize()
	if err != nil {
		return errors.Wrap(err, "error serializing NPDU")
	}
	pdu.(interface{ PutData(n ...byte) }).PutData(serialize...) // TODO: ugly cast...
	return nil
}

func (n *_NPDU) Decode(pdu Arg) error {
	if err := n._NPCI.Decode(pdu); err != nil {
		return errors.Wrap(err, "error decoding _NPCI")
	}
	switch pdu := pdu.(type) {
	case PDUData:
		data := pdu.GetPduData()
		n.PutData(data...)
		var err error
		n.npdu, err = model.NPDUParse(context.Background(), data, uint16(len(data)))
		if err != nil {
			return errors.Wrap(err, "error parsing NPDU")
		}
		n.rootMessage = n.npdu
		n.nlm = n.npdu.GetNlm()
		n.apdu = n.npdu.GetApdu()
	}
	return nil
}

func (n *_NPDU) GetRootMessage() spi.Message {
	return n.npdu
}

func (n *_NPDU) GetProtocolVersionNumber() uint8 {
	if n.npdu == nil {
		return 0
	}
	return n.npdu.GetProtocolVersionNumber()
}

func (n *_NPDU) GetControl() model.NPDUControl {
	if n.npdu == nil {
		return nil
	}
	return n.npdu.GetControl()
}

func (n *_NPDU) GetDestinationNetworkAddress() *uint16 {
	if n.npdu == nil {
		return nil
	}
	return n.npdu.GetDestinationNetworkAddress()
}

func (n *_NPDU) GetDestinationLength() *uint8 {
	if n.npdu == nil {
		return nil
	}
	return n.npdu.GetDestinationLength()
}

func (n *_NPDU) GetDestinationAddress() []uint8 {
	if n.npdu == nil {
		return nil
	}
	return n.npdu.GetDestinationAddress()
}

func (n *_NPDU) GetSourceNetworkAddress() *uint16 {
	if n.npdu == nil {
		return nil
	}
	return n.npdu.GetSourceNetworkAddress()
}

func (n *_NPDU) GetSourceLength() *uint8 {
	if n.npdu == nil {
		return nil
	}
	return n.npdu.GetSourceLength()
}

func (n *_NPDU) GetSourceAddress() []uint8 {
	if n.npdu == nil {
		return nil
	}
	return n.npdu.GetSourceAddress()
}

func (n *_NPDU) GetHopCount() *uint8 {
	if n.npdu == nil {
		return nil
	}
	return n.npdu.GetHopCount()
}

func (n *_NPDU) GetNlm() model.NLM {
	if n.npdu == nil {
		return nil
	}
	return n.npdu.GetNlm()
}

func (n *_NPDU) GetApdu() model.APDU {
	if n.npdu == nil {
		return nil
	}
	return n.npdu.GetApdu()
}

func (n *_NPDU) GetDestinationLengthAddon() uint16 {
	if n.npdu == nil {
		return 0
	}
	return n.npdu.GetDestinationLengthAddon()
}

func (n *_NPDU) GetSourceLengthAddon() uint16 {
	if n.npdu == nil {
		return 0
	}
	return n.npdu.GetSourceLengthAddon()
}

func (n *_NPDU) GetPayloadSubtraction() uint16 {
	if n.npdu == nil {
		return 0
	}
	return n.npdu.GetPayloadSubtraction()
}

func (n *_NPDU) deepCopy() *_NPDU {
	return &_NPDU{_NPCI: n._NPCI.deepCopy(), _PDUData: n._PDUData.deepCopy(), npdu: n.npdu, apdu: n.apdu}
}

func (n *_NPDU) DeepCopy() PDU {
	return n.deepCopy()
}

func (n *_NPDU) String() string {
	return fmt.Sprintf("_NPDU{%s}", n._PCI)
}
