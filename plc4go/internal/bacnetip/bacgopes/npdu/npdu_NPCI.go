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

package npdu

import (
	"context"
	"math"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type NPCI interface {
	PCI

	GetNPDUNetMessage() *uint8

	Encode(pdu Arg) error
	Decode(pdu Arg) error

	SetNpduVersion(uint8)
	GetNpduVersion() uint8
	SetNpduControl(uint8)
	GetNpduControl() uint8
	SetNpduDADR(*Address)
	GetNpduDADR() *Address
	SetNpduSADR(*Address)
	GetNpduSADR() *Address
	SetNpduHopCount(*uint8)
	GetNpduHopCount() *uint8
	SetNpduNetMessage(*uint8)
	GetNpduNetMessage() *uint8
	SetNpduVendorID(*readWriteModel.BACnetVendorId)
	GetNpduVendorID() *readWriteModel.BACnetVendorId

	GetNPCI() NPCI
}

type _NPCI struct {
	PCI
	*DebugContents

	npduVersion    uint8
	npduControl    uint8
	npduDADR       *Address
	npduSADR       *Address
	npduHopCount   *uint8
	npduNetMessage *uint8
	npduVendorID   *readWriteModel.BACnetVendorId `directSerialize:"true"`

	// Deprecated: hacky workaround
	bytesToDiscard int `ignore:"true"`
}

var _ NPCI = (*_NPCI)(nil)

func NewNPCI(args Args, kwArgs KWArgs, options ...Option) NPCI {
	n := &_NPCI{
		npduVersion: 1,
	}
	n.DebugContents = NewDebugContents(n, "npduVersion", "npduControl", "npduDADR", "npduSADR",
		"npduHopCount", "npduNetMessage", "npduVendorID")
	options = AddLeafTypeIfAbundant(options, n)
	n.PCI = NewPCI(args, kwArgs, options...)
	n.AddExtraPrinters(n.PCI.(DebugContentPrinter))
	if n.GetRootMessage() == nil {
		nlm := ExtractNLM(options)
		apdu := ExtractAPDU(options)
		if nlm != nil || apdu != nil {
			npdu, _ := n.buildNPDU(0, nil, nil, false, readWriteModel.NPDUNetworkPriority_NORMAL_MESSAGE, nlm, apdu)
			n.SetRootMessage(npdu)
		}
	}
	return n
}

func (n *_NPCI) GetDebugAttr(attr string) any {
	switch attr {
	case "npduVersion":
		return n.npduVersion
	case "npduControl":
		return n.npduControl
	case "npduDADR":
		if n.npduDADR != nil {
			return *n.npduDADR
		}
	case "npduSADR":
		if n.npduSADR != nil {
			return *n.npduSADR
		}
	case "npduHopCount":
		if n.npduHopCount != nil {
			return *n.npduHopCount
		}
	case "npduNetMessage":
		if n.npduNetMessage != nil {
			return *n.npduNetMessage
		}
	case "npduVendorID":
		if n.npduVendorID != nil {
			return *n.npduVendorID
		}
	default:
		return nil
	}
	return nil
}

func (n *_NPCI) Update(npci Arg) error {
	if err := n.PCI.Update(npci); err != nil {
		return errors.Wrap(err, "error updating _PCI")
	}
	switch npci := npci.(type) {
	case NPCI:
		n.npduVersion = npci.GetNpduVersion()
		n.npduControl = npci.GetNpduControl()
		n.npduDADR = npci.GetNpduDADR()
		n.npduSADR = npci.GetNpduSADR()
		n.npduHopCount = npci.GetNpduHopCount()
		n.npduNetMessage = npci.GetNpduNetMessage()
		n.npduVendorID = npci.GetNpduVendorID()
		return nil
	default:
		return errors.Errorf("invalid NPCI type %T", npci)
	}
}

func (n *_NPCI) Encode(pdu Arg) error {
	switch pdu := pdu.(type) {
	case PCI:
		if err := pdu.GetPCI().Update(n); err != nil {
			return errors.Wrap(err, "error updating pdu")
		}
	}
	switch pdu := pdu.(type) {
	case PDU:
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
		control |= uint8(n.GetNetworkPriority()) & 0x03
		n.npduControl = control
		pdu.Put(control)

		// make sure expecting reply and priority get passed down
		pdu.SetExpectingReply(n.GetExpectingReply())
		pdu.SetNetworkPriority(n.GetNetworkPriority())

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
	}
	return nil
}

func (n *_NPCI) Decode(pdu Arg) error {
	if err := n.PCI.Update(pdu); err != nil {
		return errors.Wrap(err, "error updating pdu")
	}
	// TODO: check if we want to stay with parsing or ditch that for now
	data := pdu.(PDUData).GetPduData()
	parse, err := readWriteModel.NPDUParse(context.Background(), data, uint16(len(data)))
	if err != nil {
		return errors.Wrap(err, "error parsing npdu")
	}
	n.SetRootMessage(parse)
	readBytes := 0
	switch rm := n.GetRootMessage().(type) {
	case readWriteModel.NPDU:
		readBytes += 1
		n.npduVersion = rm.GetProtocolVersionNumber()
		readBytes += 1
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
			readBytes += 2
			dnet := *rm.GetDestinationNetworkAddress()
			readBytes += 1
			dlen := *rm.GetDestinationLength()
			readBytes += int(dlen)
			dadr := rm.GetDestinationAddress()

			if dnet == 0xFFFF {
				n.npduDADR = NewGlobalBroadcast(nil)
			} else if dlen == 0 {
				n.npduDADR = NewRemoteBroadcast(dnet, nil)
			} else {
				var err error
				n.npduDADR, err = NewRemoteStation(&dnet, dadr, nil)
				if err != nil {
					return errors.Wrap(err, "error creating remote station")
				}
			}
		}

		// extract the source address
		if snetPresent {
			readBytes += 2
			snet := *rm.GetSourceNetworkAddress()
			readBytes += 1
			slen := *rm.GetSourceLength()
			readBytes += int(slen)
			sadr := rm.GetSourceAddress()

			if snet == 0xFFFF {
				return errors.New("SADR can't be a global broadcast")
			} else if slen == 0 {
				return errors.New("SADR can't be a remote broadcast")
			}

			var err error
			n.npduSADR, err = NewRemoteStation(&snet, sadr, nil)
			if err != nil {
				return errors.Wrap(err, "error creating remote station")
			}
		}

		// extract the hop count
		if dnetPresent {
			readBytes += 1
			n.npduHopCount = rm.GetHopCount()
		}

		// extract the network layer message type (if present)
		if netLayerMessage {
			readBytes += 1
			messageType := rm.GetNlm().GetMessageType()
			n.npduNetMessage = &messageType
			if rm.GetNlm().GetIsVendorProprietaryMessage() {
				// extract the vendor ID
				vendorId := rm.GetNlm().(readWriteModel.NLMVendorProprietaryMessage).GetVendorId()
				n.npduVendorID = &vendorId
			}
		} else {
			// application layer message
			n.npduNetMessage = nil
		}
	}
	// TODO: this is ugly but we need to read away data otherwise downstream code just fails... maybe better rip out all plc4x parsing till we get all tests running
	n.bytesToDiscard = readBytes
	return nil
}

func (n *_NPCI) GetNPDUNetMessage() *uint8 {
	switch rm := n.GetRootMessage().(type) {
	case readWriteModel.NPDU:
		if nlm := rm.GetNlm(); nlm != nil {
			messageType := nlm.GetMessageType()
			return &messageType
		}
		return nil
	default:
		return nil
	}
}

func (n *_NPCI) SetNpduVersion(u uint8) {
	n.npduVersion = u
}

func (n *_NPCI) GetNpduVersion() uint8 {
	return n.npduVersion
}

func (n *_NPCI) SetNpduControl(u uint8) {
	n.npduControl = u
}

func (n *_NPCI) GetNpduControl() uint8 {
	return n.npduControl
}

func (n *_NPCI) SetNpduDADR(address *Address) {
	n.npduDADR = address
}

func (n *_NPCI) GetNpduDADR() *Address {
	return n.npduDADR
}

func (n *_NPCI) SetNpduSADR(address *Address) {
	n.npduSADR = address
}

func (n *_NPCI) GetNpduSADR() *Address {
	return n.npduSADR
}

func (n *_NPCI) SetNpduHopCount(u *uint8) {
	n.npduHopCount = u
}

func (n *_NPCI) GetNpduHopCount() *uint8 {
	return n.npduHopCount
}

func (n *_NPCI) SetNpduNetMessage(u *uint8) {
	n.npduNetMessage = u
}

func (n *_NPCI) GetNpduNetMessage() *uint8 {
	return n.npduNetMessage
}

func (n *_NPCI) SetNpduVendorID(u *readWriteModel.BACnetVendorId) {
	n.npduVendorID = u
}

func (n *_NPCI) GetNpduVendorID() *readWriteModel.BACnetVendorId {
	return n.npduVendorID
}

func (n *_NPCI) GetNPCI() NPCI {
	return n
}

func (n *_NPCI) GetProtocolVersionNumber() uint8 {
	return n.npduVersion
}

func (n *_NPCI) GetControl() readWriteModel.NPDUControl {
	ctl, _ := readWriteModel.NPDUControlParse(context.Background(), []byte{n.npduControl})
	return ctl
}

func (n *_NPCI) GetDestinationNetworkAddress() *uint16 {
	if n.npduDADR == nil {
		return nil
	}
	return n.npduDADR.AddrNet
}

func (n *_NPCI) GetDestinationLength() *uint8 {
	if n.npduDADR == nil {
		return nil
	}
	return n.npduDADR.AddrLen
}

func (n *_NPCI) GetDestinationAddress() []uint8 {
	if n.npduDADR == nil {
		return nil
	}
	return n.npduDADR.AddrAddress
}

func (n *_NPCI) GetSourceNetworkAddress() *uint16 {
	if n.npduSADR == nil {
		return nil
	}
	return n.npduSADR.AddrNet
}

func (n *_NPCI) GetSourceLength() *uint8 {
	if n.npduSADR == nil {
		return nil
	}
	return n.npduSADR.AddrLen
}

func (n *_NPCI) GetSourceAddress() []uint8 {
	if n.npduSADR == nil {
		return nil
	}
	return n.npduSADR.AddrAddress
}

func (n *_NPCI) GetHopCount() *uint8 {
	return n.npduHopCount
}

// TODO: this needs work as it doesn't do anything any more right now... // we could hook it to update
func (n *_NPCI) buildNPDU(hopCount uint8, source *Address, destination *Address, expectingReply bool, networkPriority readWriteModel.NPDUNetworkPriority, nlm readWriteModel.NLM, apdu readWriteModel.APDU) (readWriteModel.NPDU, error) {
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
	control := readWriteModel.NewNPDUControl(nlm != nil, destinationSpecified, sourceSpecified, expectingReply, networkPriority)
	return readWriteModel.NewNPDU(1, control, destinationNetworkAddress, destinationLength, destinationAddress, sourceNetworkAddress, sourceLength, sourceAddress, destinationHopCount, nlm, apdu, 0), nil
}

func (n *_NPCI) deepCopy() *_NPCI {
	newN := &_NPCI{
		n.PCI.DeepCopy().(PCI),
		nil,
		n.npduVersion,
		n.npduControl,
		n.npduDADR.DeepCopy().(*Address),
		n.npduSADR.DeepCopy().(*Address),
		CopyPtr(n.npduHopCount),
		CopyPtr(n.npduNetMessage),
		CopyPtr(n.npduVendorID),
		n.bytesToDiscard,
	}
	newN.DebugContents = NewDebugContents(newN, "npduVersion", "npduControl", "npduDADR", "npduSADR",
		"npduHopCount", "npduNetMessage", "npduVendorID")
	newN.AddExtraPrinters(newN.PCI.(DebugContentPrinter))
	return newN
}

func (n *_NPCI) DeepCopy() any {
	return n.deepCopy()
}
