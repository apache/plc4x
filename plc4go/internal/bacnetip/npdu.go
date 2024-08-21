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

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/pkg/errors"
)

// NPDUTypes is a dictionary of message type values and structs
var NPDUTypes map[uint8]func() interface{ Decode(Arg) error }

type NPCI interface {
	PCI

	GetNPDUNetMessage() *uint8

	Encode(pdu Arg) error
	Decode(pdu Arg) error

	setNLM(nlm readWriteModel.NLM)
}

type _NPCI struct {
	*_PCI
	*DebugContents

	nlm readWriteModel.NLM
}

var _ NPCI = (*_NPCI)(nil)

func NewNPCI(msg spi.Message, nlm readWriteModel.NLM) NPCI {
	n := &_NPCI{
		nlm: nlm,
	}
	n._PCI = newPCI(msg, nil, nil, false, readWriteModel.NPDUNetworkPriority_NORMAL_MESSAGE)
	return n
}

func (n *_NPCI) GetNPDUNetMessage() *uint8 {
	if n.nlm == nil {
		return nil
	}
	messageType := n.nlm.GetMessageType()
	return &messageType
}

func (n *_NPCI) setNLM(nlm readWriteModel.NLM) {
	n.nlm = nlm
}

func (n *_NPCI) deepCopy() *_NPCI {
	return &_NPCI{_PCI: n._PCI.deepCopy()}
}

func (n *_NPCI) Update(npci Arg) error {
	if err := n._PCI.Update(npci); err != nil {
		return errors.Wrap(err, "error updating _PCI")
	}
	switch pci := npci.(type) {
	case NPCI:
		// TODO: update coordinates....
		return nil
	default:
		return errors.Errorf("invalid APCI type %T", pci)
	}
}

func (n *_NPCI) Encode(pdu Arg) error {
	if err := pdu.(interface{ Update(Arg) error }).Update(n); err != nil { // TODO: better validate that arg is really PDUData... use switch similar to Update
		return errors.Wrap(err, "error updating pdu")
	}
	// TODO: what should we do here??
	return nil
}

func (n *_NPCI) Decode(pdu Arg) error {
	if err := n._PCI.Update(pdu); err != nil {
		return errors.Wrap(err, "error updating pdu")
	}
	// TODO: what should we do here??
	return nil
}

type NPDU interface {
	NPCI
	PDUData

	setNPDU(npdu readWriteModel.NPDU)
	setAPDU(apdu readWriteModel.APDU)
}

type _NPDU struct {
	*_NPCI
	*_PDUData

	npdu readWriteModel.NPDU
	apdu readWriteModel.APDU
}

var _ NPDU = (*_NPDU)(nil)

func NewNPDU(nlm readWriteModel.NLM, apdu readWriteModel.APDU) (NPDU, error) {
	n := &_NPDU{
		apdu: apdu,
	}
	var err error
	n.npdu, err = n.buildNPDU(0, nil, nil, false, readWriteModel.NPDUNetworkPriority_NORMAL_MESSAGE, nlm, apdu)
	if err != nil {
		// TODO: seems to be a legit case, which means we don't build our npdu yet... check that
		//return nil, errors.Wrap(err, "error building NPDU")
	}
	n._NPCI = NewNPCI(n.npdu, nlm).(*_NPCI)
	n._PDUData = newPDUData(n)
	return n, nil
}

func (n *_NPDU) setNPDU(npdu readWriteModel.NPDU) {
	n.npdu = npdu
}

func (n *_NPDU) setAPDU(apdu readWriteModel.APDU) {
	n.apdu = apdu
}

func (n *_NPDU) buildNPDU(hopCount uint8, source *Address, destination *Address, expectingReply bool, networkPriority readWriteModel.NPDUNetworkPriority, nlm readWriteModel.NLM, apdu readWriteModel.APDU) (readWriteModel.NPDU, error) {
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
		sourceLengthValueUint8 := uint8(sourceLengthValue)
		sourceLength = &sourceLengthValueUint8
		sourceAddress = source.AddrAddress
		if sourceLengthValueUint8 == 0 {
			// If we define the len 0 we must not send the array
			sourceAddress = nil
		}
	}
	destinationSpecified := destination != nil
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
		destinationLengthValueUint8 := uint8(destinationLengthValue)
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

func (n *_NPDU) Encode(pdu Arg) error {
	if err := n._NPCI.Encode(pdu); err != nil {
		return errors.Wrap(err, "error encoding _NPCI")
	}
	var err error
	n.npdu, err = n.buildNPDU(0, n.GetPDUSource(), n.GetPDUDestination(), n.GetExpectingReply(), n.GetNetworkPriority(), n.nlm, n.apdu)
	if err != nil {
		return errors.Wrap(err, "error building NPDU")
	}
	n.SetPDUUserData(n.npdu)
	return nil
}

func (n *_NPDU) Decode(pdu Arg) error {
	if err := n._NPCI.Decode(pdu); err != nil {
		return errors.Wrap(err, "error decoding _NPCI")
	}
	switch pdu := pdu.(type) {
	case PDUData:
		data := pdu.GetPduData()
		var err error
		n.npdu, err = readWriteModel.NPDUParse(context.Background(), data, uint16(len(data)))
		if err != nil {
			return errors.Wrap(err, "error parsing NPDU")
		}
		n.pduUserData = n.npdu
		n.nlm = n.npdu.GetNlm()
		n.apdu = n.npdu.GetApdu()
	}
	return nil
}

func (n *_NPDU) GetMessage() spi.Message {
	return n.npdu
}

func (n *_NPDU) getPDUData() []byte {
	if n.GetMessage() == nil {
		return nil
	}
	writeBufferByteBased := utils.NewWriteBufferByteBased()
	if err := n.GetMessage().SerializeWithWriteBuffer(context.Background(), writeBufferByteBased); err != nil {
		panic(err) // TODO: graceful handle
	}
	return writeBufferByteBased.GetBytes()
}

func (n *_NPDU) deepCopy() *_NPDU {
	return &_NPDU{_NPCI: n._NPCI.deepCopy(), _PDUData: n._PDUData.deepCopy()}
}

func (n *_NPDU) DeepCopy() PDU {
	return n.deepCopy()
}

func (n *_NPDU) String() string {
	return fmt.Sprintf("_NPDU{%s}", n._PCI)
}

type WhoIsRouterToNetwork struct {
	*_NPDU

	wirtnNetwork *uint16

	readWriteModel.NLMWhoIsRouterToNetwork
}

func NewWhoIsRouterToNetwork(opts ...func(network *WhoIsRouterToNetwork)) (*WhoIsRouterToNetwork, error) {
	w := &WhoIsRouterToNetwork{}
	for _, opt := range opts {
		opt(w)
	}
	w.NLMWhoIsRouterToNetwork = readWriteModel.NewNLMWhoIsRouterToNetwork(w.wirtnNetwork, 0)
	npdu, err := NewNPDU(w.NLMWhoIsRouterToNetwork, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating NPDU")
	}
	w._NPDU = npdu.(*_NPDU)
	return w, nil
}

func WithWhoIsRouterToNetworkNet(net uint16) func(*WhoIsRouterToNetwork) {
	return func(n *WhoIsRouterToNetwork) {
		n.wirtnNetwork = &net
	}
}

func (n *WhoIsRouterToNetwork) GetWirtnNetwork() *uint16 {
	return n.wirtnNetwork
}

func (n *WhoIsRouterToNetwork) Encode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPDU:
		if err := npdu.Update(n); err != nil {
			return errors.Wrap(err, "error updating _NPCI")
		}
		if n.wirtnNetwork != nil {
			npdu.PutShort(int16(*n.wirtnNetwork))
		}
		npdu.setNPDU(n.npdu)
		npdu.setNLM(n.nlm)
		npdu.setAPDU(n.apdu)
		return nil
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
}

func (n *WhoIsRouterToNetwork) Decode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPDU:
		if err := n.Update(npdu); err != nil {
			return errors.Wrap(err, "error updating _NPCI")
		}
		switch pduUserData := npdu.GetPDUUserData().(type) {
		case readWriteModel.NPDUExactly:
			switch nlm := pduUserData.GetNlm().(type) {
			case readWriteModel.NLMWhoIsRouterToNetworkExactly:
				n.setNLM(nlm)
				n.NLMWhoIsRouterToNetwork = nlm
				n.wirtnNetwork = nlm.GetDestinationNetworkAddress()
			}
		}
		return nil
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
}

func (n *WhoIsRouterToNetwork) String() string {
	return fmt.Sprintf("WhoIsRouterToNetwork{%s, wirtnNetwork: %d}", n._NPDU, n.wirtnNetwork)
}

type IAmRouterToNetwork struct {
	*_NPDU

	iartnNetworkList []uint16

	readWriteModel.NLMIAmRouterToNetwork
}

func NewIAmRouterToNetwork(opts ...func(*IAmRouterToNetwork)) (*IAmRouterToNetwork, error) {
	i := &IAmRouterToNetwork{}
	for _, opt := range opts {
		opt(i)
	}
	i.NLMIAmRouterToNetwork = readWriteModel.NewNLMIAmRouterToNetwork(i.iartnNetworkList, 0)
	npdu, err := NewNPDU(i.NLMIAmRouterToNetwork, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating NPDU")
	}
	i._NPDU = npdu.(*_NPDU)
	return i, nil
}

func WithIAmRouterToNetworkNetworkList(iartnNetworkList ...uint16) func(*IAmRouterToNetwork) {
	return func(n *IAmRouterToNetwork) {
		n.iartnNetworkList = iartnNetworkList
	}
}

func (i *IAmRouterToNetwork) GetIartnNetworkList() []uint16 {
	return i.iartnNetworkList
}

func (i *IAmRouterToNetwork) Encode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPDU:
		if err := npdu.Update(i); err != nil {
			return errors.Wrap(err, "error updating _NPCI")
		}
		for _, net := range i.iartnNetworkList {
			npdu.PutShort(int16(net))
		}
		npdu.setNPDU(i.npdu)
		npdu.setNLM(i.nlm)
		npdu.setAPDU(i.apdu)
		return nil
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
}

func (i *IAmRouterToNetwork) Decode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPDU:
		if err := i.Update(npdu); err != nil {
			return errors.Wrap(err, "error updating _NPCI")
		}
		switch pduUserData := npdu.GetPDUUserData().(type) {
		case readWriteModel.NPDUExactly:
			switch nlm := pduUserData.GetNlm().(type) {
			case readWriteModel.NLMIAmRouterToNetworkExactly:
				i.setNLM(nlm)
				i.NLMIAmRouterToNetwork = nlm
				i.iartnNetworkList = nlm.GetDestinationNetworkAddresses()
			}
		}
		return nil
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
}

func (i *IAmRouterToNetwork) String() string {
	return fmt.Sprintf("IAmRouterToNetwork{%s, iartnNetworkList: %v}", i._NPDU, i.iartnNetworkList)
}

type ICouldBeRouterToNetwork struct {
	*_NPDU

	icbrtnNetwork          uint16
	icbrtnPerformanceIndex uint8

	readWriteModel.NLMICouldBeRouterToNetwork
}

func NewICouldBeRouterToNetwork(opts ...func(*ICouldBeRouterToNetwork)) (*ICouldBeRouterToNetwork, error) {
	i := &ICouldBeRouterToNetwork{}
	for _, opt := range opts {
		opt(i)
	}
	i.NLMICouldBeRouterToNetwork = readWriteModel.NewNLMICouldBeRouterToNetwork(i.icbrtnNetwork, i.icbrtnPerformanceIndex, 0)
	npdu, err := NewNPDU(i.NLMICouldBeRouterToNetwork, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating NPDU")
	}
	i._NPDU = npdu.(*_NPDU)
	return i, nil
}

func WithICouldBeRouterToNetworkNetwork(icbrtnNetwork uint16) func(*ICouldBeRouterToNetwork) {
	return func(n *ICouldBeRouterToNetwork) {
		n.icbrtnNetwork = icbrtnNetwork
	}
}
func WithICouldBeRouterToNetworkPerformanceIndex(icbrtnPerformanceIndex uint8) func(*ICouldBeRouterToNetwork) {
	return func(n *ICouldBeRouterToNetwork) {
		n.icbrtnPerformanceIndex = icbrtnPerformanceIndex
	}
}

func (n *ICouldBeRouterToNetwork) GeticbrtnNetwork() uint16 {
	return n.icbrtnNetwork
}

func (n *ICouldBeRouterToNetwork) GetIcbrtnPerformanceIndex() uint8 {
	return n.icbrtnPerformanceIndex
}

func (n *ICouldBeRouterToNetwork) Encode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPDU:
		if err := npdu.Update(n); err != nil {
			return errors.Wrap(err, "error updating _NPCI")
		}
		npdu.PutShort(int16(n.icbrtnNetwork))
		npdu.Put(n.icbrtnPerformanceIndex)
		npdu.setNPDU(n.npdu)
		npdu.setNLM(n.nlm)
		npdu.setAPDU(n.apdu)
		return nil
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
}

func (n *ICouldBeRouterToNetwork) Decode(npdu Arg) error {
	switch npdu := npdu.(type) {
	case NPDU:
		if err := n.Update(npdu); err != nil {
			return errors.Wrap(err, "error updating _NPCI")
		}
		switch pduUserData := npdu.GetPDUUserData().(type) {
		case readWriteModel.NPDUExactly:
			switch nlm := pduUserData.GetNlm().(type) {
			case readWriteModel.NLMICouldBeRouterToNetworkExactly:
				n.setNLM(nlm)
				n.NLMICouldBeRouterToNetwork = nlm
				n.icbrtnNetwork = nlm.GetDestinationNetworkAddress()
				n.icbrtnPerformanceIndex = nlm.GetPerformanceIndex()
			}
		}
		return nil
	default:
		return errors.Errorf("invalid NPDU type %T", npdu)
	}
}

func (n *ICouldBeRouterToNetwork) String() string {
	return fmt.Sprintf("ICouldBeRouterToNetwork{%s, icbrtnNetwork: %v, icbrtnPerformanceIndex: %v}", n._NPDU, n.icbrtnNetwork, n.icbrtnPerformanceIndex)
}

type RejectMessageToNetwork struct {
	*_NPDU
}

func NewRejectMessageToNetwork() (*RejectMessageToNetwork, error) {
	panic("implement me")
}

type RouterBusyToNetwork struct {
	*_NPDU
}

func NewRouterBusyToNetwork() (*RouterBusyToNetwork, error) {
	panic("implement me")
}

type RouterAvailableToNetwork struct {
	*_NPDU
}

func NewRouterAvailableToNetwork() (*RouterAvailableToNetwork, error) {
	panic("implement me")
}

type RoutingTableEntry struct {
	*_NPDU
}

func NewRoutingTableEntry() (*RoutingTableEntry, error) {
	panic("implement me")
}

type InitializeRoutingTable struct {
	*_NPDU
}

func NewInitializeRoutingTable() (*InitializeRoutingTable, error) {
	panic("implement me")
}

type InitializeRoutingTableAck struct {
	*_NPDU
}

func NewInitializeRoutingTableAck() (*InitializeRoutingTableAck, error) {
	panic("implement me")
}

type EstablishConnectionToNetwork struct {
	*_NPDU
}

func NewEstablishConnectionToNetwork() (*EstablishConnectionToNetwork, error) {
	panic("implement me")
}

type DisconnectConnectionToNetwork struct {
	*_NPDU
}

func NewDisconnectConnectionToNetwork() (*DisconnectConnectionToNetwork, error) {
	panic("implement me")
}

type WhatIsNetworkNumber struct {
	*_NPDU
}

func NewWhatIsNetworkNumber() (*WhatIsNetworkNumber, error) {
	panic("implement me")
}

type NetworkNumberIs struct {
	*_NPDU
}

func NewNetworkNumberIs() (*NetworkNumberIs, error) {
	panic("implement me")
}

func init() {
	NPDUTypes = map[uint8]func() interface{ Decode(Arg) error }{
		0x00: func() interface{ Decode(Arg) error } {
			v, _ := NewWhoIsRouterToNetwork()
			return v
		},
		0x01: func() interface{ Decode(Arg) error } {
			v, _ := NewIAmRouterToNetwork()
			return v
		},
		0x02: func() interface{ Decode(Arg) error } {
			v, _ := NewICouldBeRouterToNetwork()
			return v
		},
		// 0x03: NewRejectRouterToNetwork, // TODO: not present upstream
		0x04: func() interface{ Decode(Arg) error } {
			v, _ := NewRouterBusyToNetwork()
			return v
		},
		0x05: func() interface{ Decode(Arg) error } {
			v, _ := NewRouterBusyToNetwork()
			return v
		},
		0x06: func() interface{ Decode(Arg) error } {
			v, _ := NewInitializeRoutingTable()
			return v
		},
		0x07: func() interface{ Decode(Arg) error } {
			v, _ := NewInitializeRoutingTableAck()
			return v
		},
		0x08: func() interface{ Decode(Arg) error } {
			v, _ := NewEstablishConnectionToNetwork()
			return v
		},
		0x09: func() interface{ Decode(Arg) error } {
			v, _ := NewDisconnectConnectionToNetwork()
			return v
		},
		// 0x0A: NewChallengeRequest, // TODO: not present upstream
		// 0x0B: NewSecurityPayload, // TODO: not present upstream
		// 0x0C: NewSecurityResponse, // TODO: not present upstream
		// 0x0D: NewRequestKeyUpdate, // TODO: not present upstream
		// 0x0E: NewUpdateKeyUpdate, // TODO: not present upstream
		// 0x0F: NewUpdateKeyDistributionKey, // TODO: not present upstream
		// 0x10: NewRequestMasterKey, // TODO: not present upstream
		// 0x11: NewSetMasterKey, // TODO: not present upstream
		0x12: func() interface{ Decode(Arg) error } {
			v, _ := NewWhatIsNetworkNumber()
			return v
		},
		0x13: func() interface{ Decode(Arg) error } {
			v, _ := NewNetworkNumberIs()
			return v
		},
	}
}
