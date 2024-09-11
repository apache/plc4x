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

package test_bvll

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/app"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/appservice"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/bvllservice"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/local/device"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/netservice"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/object"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/service"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/quick"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/vlan"
)

type _NetworkServiceElement struct {
	*NetworkServiceElement
}

func new_NetworkServiceElement(localLog zerolog.Logger) (*_NetworkServiceElement, error) {
	i := &_NetworkServiceElement{}

	// This class turns off the deferred startup function call that broadcasts
	// I-Am-Router-To-Network and Network-Number-Is messages.
	var err error
	i.NetworkServiceElement, err = NewNetworkServiceElement(localLog, WithNetworkServiceElementStartupDisabled(true))
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service element")
	}
	return i, nil
}

//go:generate plc4xGenerator -type=FauxMultiplexer -prefix=
type FauxMultiplexer struct {
	Client
	Server

	address        *Address
	unicastTuple   *AddressTuple[string, uint16]
	broadcastTuple *AddressTuple[string, uint16]

	node *IPNode

	log zerolog.Logger
}

func NewFauxMultiplexer(localLog zerolog.Logger, addr *Address, network *IPNetwork) (*FauxMultiplexer, error) {
	f := &FauxMultiplexer{
		address: addr,
		log:     localLog,
	}
	var err error
	f.Client, err = NewClient(localLog, f)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	f.Server, err = NewServer(localLog, f)
	if err != nil {
		return nil, errors.Wrap(err, "error creating server")
	}

	// get the unicast and broadcast tuples
	f.unicastTuple = addr.AddrTuple
	f.broadcastTuple = addr.AddrBroadcastTuple

	// make an internal node and bind to it, this takes the place of
	// both the direct port and broadcast port of the real UDPMultiplexer
	f.node, err = NewIPNode(localLog, addr, network)
	if err != nil {
		return nil, errors.Wrap(err, "error creating ip node")
	}
	if err := Bind(localLog, f, f.node); err != nil {
		return nil, errors.Wrap(err, "error binding")
	}
	return f, nil
}

func (s *FauxMultiplexer) Indication(args Args, kwargs KWArgs) error {
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")

	pdu := GA[PDU](args, 0)

	var dest *Address
	// check for a broadcast message
	if pdu.GetPDUDestination().AddrType == LOCAL_BROADCAST_ADDRESS {
		var err error
		dest, err = NewAddress(NA(NA(s.broadcastTuple)))
		if err != nil {
			return errors.Wrap(err, "error creating address")
		}
		s.log.Debug().Stringer("dest", dest).Msg("Requesting local broadcast")
	} else if pdu.GetPDUDestination().AddrType == LOCAL_STATION_ADDRESS {
		var err error
		dest, err = NewAddress(NA(NA(pdu.GetPDUDestination().AddrAddress)))
		if err != nil {
			return errors.Wrap(err, "error creating address")
		}
		s.log.Debug().Stringer("dest", dest).Msg("Requesting local station")
	} else {
		return errors.New("unknown destination type")
	}

	unicast, err := NewAddress(NA(NA(s.unicastTuple)))
	if err != nil {
		return errors.Wrap(err, "error creating address")
	}
	return s.Request(NA(NewPDU(NoArgs, NKW(KWCompRootMessage, pdu, KWCPCISource, unicast, KWCPCIDestination, dest))), NoKWArgs)
}

func (s *FauxMultiplexer) Confirmation(args Args, kwargs KWArgs) error {
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")
	pdu := GA[PDU](args, 0)

	// the PDU source and destination are tuples, convert them to Address instances
	src := pdu.GetPDUSource()

	broadcast, err := NewAddress(NA(s.broadcastTuple))
	if err != nil {
		return errors.Wrap(err, "error creating address")
	}
	var dest *Address
	// see if the destination was our broadcast address
	if pdu.GetPDUDestination().Equals(broadcast) {
		dest = NewLocalBroadcast(nil)
	} else {
		dest, err = NewAddress(NA(pdu.GetPDUDestination().AddrAddress))
		if err != nil {
			return errors.Wrap(err, "error creating address")
		}
	}

	return s.Response(NA(NewPDU(NoArgs, NKW(KWCompRootMessage, pdu, KWCPCISource, src, KWCPCIDestination, dest))), NoKWArgs)
}

type SnifferStateMachine struct {
	*ClientStateMachine

	address *Address
	annexj  *AnnexJCodec
	mux     *FauxMultiplexer

	log zerolog.Logger
}

func NewSnifferStateMachine(localLog zerolog.Logger, address string, vlan *IPNetwork) (*SnifferStateMachine, error) {
	s := &SnifferStateMachine{
		log: localLog,
	}
	machine, err := NewClientStateMachine(localLog, WithClientStateMachineName(address), WithClientStateMachineExtension(s))
	if err != nil {
		return nil, errors.Wrap(err, "error building client state machine")
	}
	s.ClientStateMachine = machine

	// save the name and address
	s.address, err = NewAddress(NA(address))
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}

	// BACnet/IP interpreter
	s.annexj, err = NewAnnexJCodec(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating annexj")
	}

	// fake multiplexer has a VLAN node in it
	s.mux, err = NewFauxMultiplexer(localLog, s.address, vlan)
	if err != nil {
		return nil, errors.Wrap(err, "error creating faux multiplexer")
	}

	// might receive all packets and allow spoofing
	s.mux.node.SetPromiscuous(true)
	s.mux.node.SetSpoofing(true)

	// bind the stack together
	if err := Bind(localLog, s, s.annexj, s.mux); err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	return s, nil
}

// BIPStateMachine is an application layer for BVLL messages that has no BVLL
//
//	processing like the 'simple', 'foreign', or 'bbmd' versions.  The client
//	state machine sits above and Annex-J codec so the send and receive PDUs are
//	BVLL PDUs.
type BIPStateMachine struct {
	*ClientStateMachine

	address *Address
	annexj  *AnnexJCodec
	mux     *FauxMultiplexer
}

func NewBIPStateMachine(localLog zerolog.Logger, address string, vlan *IPNetwork) (*BIPStateMachine, error) {
	b := &BIPStateMachine{}
	var err error
	b.ClientStateMachine, err = NewClientStateMachine(localLog, WithClientStateMachineName(address), WithClientStateMachineExtension(b))
	if err != nil {
		return nil, errors.Wrap(err, "error building client state machine")
	}

	// save the name and address
	b.address, err = NewAddress(NA(address))
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}

	// BACnet/IP interpreter
	b.annexj, err = NewAnnexJCodec(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating annexj")
	}

	// fake multiplexer has a VLAN node in it
	b.mux, err = NewFauxMultiplexer(localLog, b.address, vlan)

	// bind the stack together
	err = Bind(localLog, b, b.annexj, b.mux)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}
	return b, nil
}

type BIPSimpleStateMachine struct {
	*ClientStateMachine
	name string

	address *Address

	bip    *BIPSimple
	annexj *AnnexJCodec
	mux    *FauxMultiplexer

	log zerolog.Logger
}

func NewBIPSimpleStateMachine(localLog zerolog.Logger, netstring string, vlan *IPNetwork) (*BIPSimpleStateMachine, error) {
	b := &BIPSimpleStateMachine{
		log: localLog,
	}
	var err error
	b.ClientStateMachine, err = NewClientStateMachine(localLog, WithClientStateMachineName(netstring), WithClientStateMachineExtension(b))
	if err != nil {
		return nil, errors.Wrap(err, "error building client state machine")
	}

	// save the name and address
	b.address, err = NewAddress(NA(netstring))
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}

	// BACnet/IP interpreter
	b.bip, err = NewBIPSimple(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building bip simple")
	}
	b.annexj, err = NewAnnexJCodec(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building annexj codec")
	}

	// fake multiplexer has a VLAN node in it
	b.mux, err = NewFauxMultiplexer(localLog, b.address, vlan)
	if err != nil {
		return nil, errors.Wrap(err, "error creating faux")
	}

	// bind the stack together
	if err := Bind(localLog, b, b.bip, b.annexj, b.mux); err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	return b, nil
}

// BIPForeignStateMachine  sits on a BIPForeign instance, the send() and receive()
//
//	parameters are NPDUs.
type BIPForeignStateMachine struct {
	*ClientStateMachine

	address *Address
	bip     *BIPForeign
	annexj  *AnnexJCodec
	mux     *FauxMultiplexer

	log zerolog.Logger
}

func NewBIPForeignStateMachine(localLog zerolog.Logger, address string, vlan *IPNetwork) (*BIPForeignStateMachine, error) {
	b := &BIPForeignStateMachine{
		log: localLog,
	}
	var err error
	b.ClientStateMachine, err = NewClientStateMachine(localLog, WithClientStateMachineName(address), WithClientStateMachineExtension(b))
	if err != nil {
		return nil, errors.New("error building client state machine")
	}

	// save the name and address
	b.address, err = NewAddress(NA(address))
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}

	// BACnet/IP interpreter
	b.bip, err = NewBIPForeign(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating BIPForeign")
	}
	b.annexj, err = NewAnnexJCodec(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating AnnexJCodec")
	}

	// fake multiplexer has a VLAN node in it
	b.mux, err = NewFauxMultiplexer(localLog, b.address, vlan)
	if err != nil {
		return nil, errors.Wrap(err, "error creating FauxMultiplexer")
	}

	// bind the stack together
	err = Bind(b.log, b, b.bip, b.annexj, b.mux)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}
	return b, nil
}

type BIPBBMDStateMachine struct {
	*ClientStateMachine

	address *Address
	bip     *BIPBBMD
	annexj  *AnnexJCodec
	mux     *FauxMultiplexer

	log zerolog.Logger
}

func NewBIPBBMDStateMachine(localLog zerolog.Logger, address string, vlan *IPNetwork) (*BIPBBMDStateMachine, error) {
	b := &BIPBBMDStateMachine{
		log: localLog,
	}
	var err error
	b.ClientStateMachine, err = NewClientStateMachine(localLog, WithClientStateMachineName(address), WithClientStateMachineExtension(b))
	if err != nil {
		return nil, errors.New("error building client state machine")
	}

	// save the name and address
	b.address, err = NewAddress(NA(address))
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}

	// BACnet/IP interpreter
	b.bip, err = NewBIPBBMD(localLog, b.address)
	b.annexj, err = NewAnnexJCodec(localLog)

	// build an address, full mask
	bdtAddress := fmt.Sprintf("%s/32:%d", b.address.AddrTuple.Left, b.address.AddrTuple.Right)
	b.log.Debug().Str("bdtAddress", bdtAddress).Msg("bdtAddress")

	// add itself as the first entry in the BDT
	if err := b.bip.AddPeer(quick.Address(bdtAddress)); err != nil {
		return nil, errors.Wrap(err, "error adding peer")
	}

	// fake multiplexer has a VLAN node in it
	b.mux, err = NewFauxMultiplexer(localLog, b.address, vlan)
	if err != nil {
		return nil, errors.Wrap(err, "error creating faux multiplexer")
	}

	// bind the stack together
	err = Bind(b.log, b, b.bip, b.annexj, b.mux)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}
	return b, nil
}

type BIPSimpleNode struct {
	name    string
	address *Address
	bip     *BIPSimple
	annexj  *AnnexJCodec
	mux     *FauxMultiplexer
}

func NewBIPSimpleNode(localLog zerolog.Logger, address string, vlan *IPNetwork) (*BIPSimpleNode, error) {
	b := &BIPSimpleNode{}

	// save the name and address
	b.name = address
	var err error
	b.address, err = NewAddress(NA(address))
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}

	// BACnet/IP interpreter
	b.bip, err = NewBIPSimple(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building bip simple")
	}
	b.annexj, err = NewAnnexJCodec(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building annexj codec")
	}

	// fake multiplexer has a VLAN node in it
	b.mux, err = NewFauxMultiplexer(localLog, b.address, vlan)

	// bind the stack together
	err = Bind(localLog, b.bip, b.annexj, b.mux)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	return b, nil
}

type BIPBBMDNode struct {
	name    string
	address *Address

	bip    *BIPBBMD
	annexj *AnnexJCodec
	mux    *FauxMultiplexer

	log zerolog.Logger
}

func NewBIPBBMDNode(localLog zerolog.Logger, address string, vlan *IPNetwork) (*BIPBBMDNode, error) {
	b := &BIPBBMDNode{
		log: localLog,
	}

	// build a name, save the address
	b.name = fmt.Sprintf("app @ %s", address)
	var err error
	b.address, err = NewAddress(NA(address))
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}
	b.log.Debug().Str("address", address).Msg("address")

	// BACnet/IP interpreter
	b.bip, err = NewBIPBBMD(b.log, b.address)
	if err != nil {
		return nil, errors.Wrap(err, "error building bip bbmd")
	}
	b.annexj, err = NewAnnexJCodec(b.log)
	if err != nil {
		return nil, errors.Wrap(err, "error building annexj codec")
	}

	// build an address, full mask
	bdtAddress := fmt.Sprintf("%s/32:%d", b.address.AddrTuple.Left, b.address.AddrTuple.Right)
	b.log.Debug().Str("bdtAddress", bdtAddress).Msg("bdtAddress")

	// add itself as the first entry in the BDT
	bbmdAddress, err := NewAddress(NA(bdtAddress))
	if err != nil {
		return nil, errors.Wrap(err, "error creating bbmd address")
	}
	err = b.bip.AddPeer(bbmdAddress)
	if err != nil {
		return nil, errors.Wrap(err, "error adding peer")
	}

	// fake multiplexer has a VLAN node in it
	b.mux, err = NewFauxMultiplexer(b.log, b.address, vlan)
	if err != nil {
		return nil, errors.Wrap(err, "error creating faux multiplexer")
	}

	// bind the stack together
	err = Bind(b.log, b.bip, b.annexj, b.mux)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	return b, nil
}

type TestDeviceObject struct {
	*LocalDeviceObject
}

//go:generate plc4xGenerator -type=BIPSimpleApplicationLayerStateMachine
type BIPSimpleApplicationLayerStateMachine struct {
	ApplicationServiceElementContract
	*ClientStateMachine

	log zerolog.Logger // TODO: move down

	name    string
	address *Address
	asap    *ApplicationServiceAccessPoint
	smap    *StateMachineAccessPoint
	nsap    *NetworkServiceAccessPoint
	nse     *_NetworkServiceElement
	bip     *BIPSimple
	annexj  *AnnexJCodec
	mux     *FauxMultiplexer
}

func NewBIPSimpleApplicationLayerStateMachine(localLog zerolog.Logger, address string, vlan *IPNetwork) (*BIPSimpleApplicationLayerStateMachine, error) {
	b := &BIPSimpleApplicationLayerStateMachine{}
	// build a name, save the address
	b.name = fmt.Sprintf("app @ %s", address)
	var err error
	b.address, err = NewAddress(NA(address))
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}

	// build a local device object
	localDevice := &TestDeviceObject{
		LocalDeviceObject: &LocalDeviceObject{
			ObjectName:       b.name,
			ObjectIdentifier: "device:998",
			VendorIdentifier: 999,
		},
	}

	// continue with initialization
	b.ApplicationServiceElementContract, err = NewApplicationServiceElement(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building application")
	}
	b.ClientStateMachine, err = NewClientStateMachine(localLog, WithClientStateMachineName(b.name), WithClientStateMachineExtension(b))

	// include a application decoder
	b.asap, err = NewApplicationServiceAccessPoint(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building application service access point")
	}

	// pass the device object to the state machine access point so it
	// can know if it should support segmentation
	// the segmentation state machines need access to the same device
	// information cache as the application
	b.smap, err = NewStateMachineAccessPoint(localLog, localDevice.LocalDeviceObject, WithStateMachineAccessPointDeviceInfoCache(NewDeviceInfoCache(localLog))) //TODO: this is a indirection that wasn't intended... we don't use the annotation yet so that might be fine
	if err != nil {
		return nil, errors.Wrap(err, "error building state machine access point")
	}

	// a network service access point will be needed
	b.nsap, err = NewNetworkServiceAccessPoint(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}

	// give the NSAP a generic network layer service element
	b.nse, err = new_NetworkServiceElement(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service element")
	}
	err = Bind(localLog, b.nse, b.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	// bind the top layers
	err = Bind(localLog, b, b.asap, b.smap, b.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	// BACnet/IP interpreter
	b.bip, err = NewBIPSimple(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building bip bbmd")
	}
	b.annexj, err = NewAnnexJCodec(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building annexj codec")
	}

	// fake multiplexer has a VLAN node in it
	b.mux, err = NewFauxMultiplexer(localLog, b.address, vlan)
	if err != nil {
		return nil, errors.Wrap(err, "error building multiplexer")
	}

	// bind the stack together
	err = Bind(localLog, b.bip, b.annexj, b.mux)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	// bind the stack to the local network
	err = b.nsap.Bind(b.bip, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	return b, nil
}

func (b *BIPSimpleApplicationLayerStateMachine) Indication(args Args, kwargs KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")
	return b.Receive(args, NoKWArgs)
}

func (b *BIPSimpleApplicationLayerStateMachine) Confirmation(args Args, kwargs KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Confirmation")
	return b.Receive(args, NoKWArgs)
}

type BIPBBMDApplication struct {
	*Application
	*WhoIsIAmServices
	*ReadWritePropertyServices

	name    string
	address *Address

	asap   *ApplicationServiceAccessPoint
	smap   *StateMachineAccessPoint
	nsap   *NetworkServiceAccessPoint
	nse    *_NetworkServiceElement
	bip    *BIPBBMD
	annexj *AnnexJCodec
	mux    *FauxMultiplexer

	log zerolog.Logger
}

func NewBIPBBMDApplication(localLog zerolog.Logger, address string, vlan *IPNetwork) (*BIPBBMDApplication, error) {
	b := &BIPBBMDApplication{
		log: localLog,
	}

	// build a name, save the address
	b.name = fmt.Sprintf("app @ %s", address)
	var err error
	b.address, err = NewAddress(NA(address))
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}

	// build a local device object
	localDevice := &TestDeviceObject{
		LocalDeviceObject: &LocalDeviceObject{
			ObjectName:       b.name,
			ObjectIdentifier: "device:999",
			VendorIdentifier: 999,
		},
	}

	// continue with initialization
	b.Application, err = NewApplication(localLog, localDevice.LocalDeviceObject) //TODO: this is a indirection that wasn't intended... we don't use the annotation yet so that might be fine
	if err != nil {
		return nil, errors.Wrap(err, "error building application")
	}

	// include a application decoder
	b.asap, err = NewApplicationServiceAccessPoint(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building application service access point")
	}

	// pass the device object to the state machine access point so it
	// can know if it should support segmentation
	// the segmentation state machines need access to the same device
	// information cache as the application
	b.smap, err = NewStateMachineAccessPoint(localLog, localDevice.LocalDeviceObject, WithStateMachineAccessPointDeviceInfoCache(b.GetDeviceInfoCache())) //TODO: this is a indirection that wasn't intended... we don't use the annotation yet so that might be fine
	if err != nil {
		return nil, errors.Wrap(err, "error building state machine access point")
	}

	// a network service access point will be needed
	b.nsap, err = NewNetworkServiceAccessPoint(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}

	// give the NSAP a generic network layer service element
	b.nse, err = new_NetworkServiceElement(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service element")
	}
	err = Bind(localLog, b.nse, b.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	// bind the top layers
	err = Bind(localLog, b, b.asap, b.smap, b.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	// BACnet/IP interpreter
	b.bip, err = NewBIPBBMD(localLog, b.address)
	if err != nil {
		return nil, errors.Wrap(err, "error building bip bbmd")
	}

	b.annexj, err = NewAnnexJCodec(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building annexj codec")
	}

	// build an address, full mask
	bdtAddress := fmt.Sprintf("%s/32:%d", b.address.AddrTuple.Left, b.address.AddrTuple.Right)
	localLog.Debug().Str("bdtAddress", bdtAddress).Msg("bdtAddress")

	// add itself as the first entry in the BDT
	bbmdAddress, err := NewAddress(NA(bdtAddress))
	if err != nil {
		return nil, errors.Wrap(err, "error creating bbmd address")
	}
	err = b.bip.AddPeer(bbmdAddress)
	if err != nil {
		return nil, errors.Wrap(err, "error adding peer")
	}

	// fake multiplexer has a VLAN node in it
	b.mux, err = NewFauxMultiplexer(localLog, b.address, vlan)
	if err != nil {
		return nil, errors.Wrap(err, "error building multiplexer")
	}

	// bind the stack together
	err = Bind(localLog, b.bip, b.annexj, b.mux)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	// bind the stack to the local network
	err = b.nsap.Bind(b.bip, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	return b, nil
}
