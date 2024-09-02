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
	"math"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

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
	setNpduVendorID(*model.BACnetVendorId)
	getNpduVendorID() *model.BACnetVendorId
	setNPDU(model.NPDU)
	getNPDU() model.NPDU
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
	npduVendorID   *model.BACnetVendorId

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
	case model.NLM:
		n.nlm = ud
	case model.APDU:
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

func (n *_NPCI) setNpduVendorID(u *model.BACnetVendorId) {
	n.npduVendorID = u
}

func (n *_NPCI) getNpduVendorID() *model.BACnetVendorId {
	return n.npduVendorID
}

func (n *_NPCI) setNPDU(npdu model.NPDU) {
	n.npdu = npdu
}

func (n *_NPCI) getNPDU() model.NPDU {
	return n.npdu
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

		n.npdu = npci.getNPDU()
		n.nlm = npci.getNLM()
		n.apdu = npci.getAPDU()
		return nil
	default:
		return errors.Errorf("invalid NPCI type %T", npci)
	}
}

// TODO: this needs work as it doesn't do anything any more right now...
func (n *_NPCI) buildNPDU(hopCount uint8, source *Address, destination *Address, expectingReply bool, networkPriority model.NPDUNetworkPriority, nlm model.NLM, apdu model.APDU) (model.NPDU, error) {
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

func (n *_NPCI) Encode(pdu Arg) error {
	switch pdu := pdu.(type) {
	case PDU:
		if err := pdu.Update(n); err != nil {
			return errors.Wrap(err, "error updating NPCI")
		}
		// only version 1 messages supported
		pdu.Put(n.npduVersion)

		// build the flags
		var netLayerMessage uint8
		if n.npduNetMessage != nil {
			netLayerMessage = 0x80
		}

		// map the destination address
		var dnetPresent uint8
		if n.npduDADR != nil {
			dnetPresent = 0x20
		}

		// map the source address
		var snetPresent uint8
		if n.npduSADR != nil {
			snetPresent = 0x08
		}

		// encode the control octet
		control := netLayerMessage | dnetPresent | snetPresent
		if n.GetExpectingReply() {
			control = control | 0x04
		}
		control |= uint8(n.networkPriority) & 0x03
		n.npduControl = control
		pdu.Put(control)

		// make sure expecting reply and priority get passed down
		pdu.SetExpectingReply(true)
		pdu.SetNetworkPriority(n.networkPriority)

		// encode the destination address
		if dnetPresent != 0 {
			if n.npduDADR.AddrType == REMOTE_STATION_ADDRESS {
				pdu.PutShort(*n.npduDADR.AddrNet)
				pdu.Put(*n.npduDADR.AddrLen)
				pdu.PutData(n.npduDADR.AddrAddress...)
			} else if n.npduDADR.AddrType == REMOTE_BROADCAST_ADDRESS {
				pdu.PutShort(*n.npduDADR.AddrNet)
				pdu.Put(0)
			} else if n.npduDADR.AddrType == GLOBAL_BROADCAST_ADDRESS {
				pdu.PutShort(0xFFFF)
				pdu.Put(0)
			}
		}

		// encode the source address
		if snetPresent != 0 {
			pdu.PutShort(*n.npduSADR.AddrNet)
			pdu.Put(*n.npduSADR.AddrLen)
			pdu.PutData(n.npduSADR.AddrAddress...)
		}

		// Put the hop count
		if dnetPresent != 0 {
			pdu.Put(*n.npduHopCount)
		}

		// Put the network layer message type (if present)
		if netLayerMessage != 0 {
			pdu.Put(*n.npduNetMessage)
			// Put the vendor ID
			if *n.npduNetMessage >= 0x80 {
				pdu.PutShort(uint16(*n.npduVendorID))
			}
		}
		switch pdu := pdu.(type) {
		case NPDU:
			pdu.setNPDU(n.npdu)
			pdu.setNLM(n.nlm)
			pdu.setAPDU(n.apdu)
		}
	}
	return nil
}

func (n *_NPCI) Decode(pdu Arg) error {
	if err := n._PCI.Update(pdu); err != nil {
		return errors.Wrap(err, "error updating pdu")
	}
	switch rm := n.rootMessage.(type) {
	case *messageBridge:
		data := rm.GetPduData()
		parse, err := model.NPDUParse(context.Background(), data, uint16(len(data)))
		if err != nil {
			return errors.Wrap(err, "error parsing npdu")
		}
		n.rootMessage = parse
	}
	switch rm := n.rootMessage.(type) {
	case model.NPDU:
		n.npdu = rm
		n.nlm = rm.GetNlm()
		n.apdu = rm.GetApdu()

		n.npduVersion = rm.GetProtocolVersionNumber()
		control := rm.GetControl()
		cs, _ := control.Serialize()
		n.npduControl = cs[0]
		netLayerMessage := control.GetMessageTypeFieldPresent()
		dnetPresent := control.GetDestinationSpecified()
		snetPresent := control.GetSourceSpecified()
		n.SetExpectingReply(control.GetExpectingReply())
		n.SetNetworkPriority(control.GetNetworkPriority())

		// extract the destination address
		if dnetPresent {
			dnet := *rm.GetDestinationNetworkAddress()
			dlen := *rm.GetDestinationLength()
			dadr := rm.GetDestinationAddress()

			if dnet == 0xFFFF {
				n.npduDADR = NewGlobalBroadcast(nil)
			} else if dlen == 0 {
				n.npduDADR = NewRemoteBroadcast(dnet, nil)
			} else {
				var err error
				n.npduDADR, err = NewRemoteStation(zerolog.Nop(), &dnet, dadr, nil)
				if err != nil {
					return errors.Wrap(err, "error creating remote station")
				}
			}

			// extract the source address
			if snetPresent {
				snet := *rm.GetSourceNetworkAddress()
				slen := *rm.GetSourceLength()
				sadr := rm.GetSourceAddress()

				if snet == 0xFFFF {
					return errors.New("SADR can't be a global broadcast")
				} else if slen == 0 {
					return errors.New("SADR can't be a remote broadcast")
				}

				var err error
				n.npduSADR, err = NewRemoteStation(zerolog.Nop(), &snet, sadr, nil)
				if err != nil {
					return errors.Wrap(err, "error creating remote station")
				}
			}

			// extract the hop count
			if dnetPresent {
				n.npduHopCount = rm.GetHopCount()
			}

			// extract the network layer message type (if present)
			if netLayerMessage {
				messageType := rm.GetNlm().GetMessageType()
				n.npduNetMessage = &messageType
				if rm.GetNlm().GetIsVendorProprietaryMessage() {
					// extract the vendor ID
					vendorId := rm.GetNlm().(model.NLMVendorProprietaryMessage).GetVendorId()
					n.npduVendorID = &vendorId
				}
			}

		} else {
			// application layer message
			n.npduNetMessage = nil
		}
	}
	return nil
}

func (n *_NPCI) deepCopy() *_NPCI {
	return &_NPCI{
		_PCI:           n._PCI.deepCopy(),
		npduVersion:    n.npduVersion,
		npduControl:    n.npduControl,
		npduDADR:       n.npduDADR.deepCopy(),
		npduSADR:       n.npduSADR.deepCopy(),
		npduHopCount:   CopyPtr(n.npduHopCount),
		npduNetMessage: CopyPtr(n.npduNetMessage),
		npduVendorID:   CopyPtr(n.npduVendorID),
		npdu:           n.npdu,
		nlm:            n.nlm,
		apdu:           n.apdu,
	}
}
