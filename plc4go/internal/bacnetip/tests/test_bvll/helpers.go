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

	"github.com/apache/plc4x/plc4go/internal/bacnetip"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/constructors"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/tests"
)

type _NetworkServiceElement struct {
	*bacnetip.NetworkServiceElement
}

func new_NetworkServiceElement(localLog zerolog.Logger) (*_NetworkServiceElement, error) {
	i := &_NetworkServiceElement{}

	// This class turns off the deferred startup function call that broadcasts
	// I-Am-Router-To-Network and Network-Number-Is messages.
	var err error
	i.NetworkServiceElement, err = bacnetip.NewNetworkServiceElement(localLog, bacnetip.WithNetworkServiceElementStartupDisabled(true))
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service element")
	}
	return i, nil
}

type FauxMultiplexer struct {
	bacnetip.Client
	bacnetip.Server

	address        *bacnetip.Address
	unicastTuple   *bacnetip.AddressTuple[string, uint16]
	broadcastTuple *bacnetip.AddressTuple[string, uint16]

	node *bacnetip.IPNode

	log zerolog.Logger
}

func NewFauxMultiplexer(localLog zerolog.Logger, addr *bacnetip.Address, network *bacnetip.IPNetwork) (*FauxMultiplexer, error) {
	f := &FauxMultiplexer{
		address: addr,
		log:     localLog,
	}
	var err error
	f.Client, err = bacnetip.NewClient(localLog, f)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	f.Server, err = bacnetip.NewServer(localLog, f)
	if err != nil {
		return nil, errors.Wrap(err, "error creating server")
	}

	// get the unicast and broadcast tuples
	f.unicastTuple = addr.AddrTuple
	f.broadcastTuple = addr.AddrBroadcastTuple

	// make an internal node and bind to it, this takes the place of
	// both the direct port and broadcast port of the real UDPMultiplexer
	f.node, err = bacnetip.NewIPNode(localLog, addr, network)
	if err != nil {
		return nil, errors.Wrap(err, "error creating ip node")
	}
	if err := bacnetip.Bind(localLog, f, f.node); err != nil {
		return nil, errors.Wrap(err, "error binding")
	}
	return f, nil
}

func (s *FauxMultiplexer) String() string {
	return fmt.Sprintf("FauxMultiplexer(TBD...)") // TODO: fill some info here
}

func (s *FauxMultiplexer) Indication(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")

	pdu := args.Get0PDU()

	var dest *bacnetip.Address
	// check for a broadcast message
	if pdu.GetPDUDestination().AddrType == bacnetip.LOCAL_BROADCAST_ADDRESS {
		var err error
		dest, err = bacnetip.NewAddress(s.log, s.broadcastTuple)
		if err != nil {
			return errors.Wrap(err, "error creating address")
		}
		s.log.Debug().Stringer("dest", dest).Msg("Requesting local broadcast")
	} else if pdu.GetPDUDestination().AddrType == bacnetip.LOCAL_STATION_ADDRESS {
		var err error
		dest, err = bacnetip.NewAddress(s.log, pdu.GetPDUDestination().AddrAddress)
		if err != nil {
			return errors.Wrap(err, "error creating address")
		}
		s.log.Debug().Stringer("dest", dest).Msg("Requesting local station")
	} else {
		return errors.New("unknown destination type")
	}

	unicast, err := bacnetip.NewAddress(s.log, s.unicastTuple)
	if err != nil {
		return errors.Wrap(err, "error creating address")
	}
	return s.Request(bacnetip.NewArgs(bacnetip.NewPDU(pdu, bacnetip.WithPDUSource(unicast), bacnetip.WithPDUDestination(dest))), bacnetip.NoKWArgs)
}

func (s *FauxMultiplexer) Confirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")
	pdu := args.Get0PDU()

	// the PDU source and destination are tuples, convert them to Address instances
	src := pdu.GetPDUSource()

	broadcast, err := bacnetip.NewAddress(s.log, s.broadcastTuple)
	if err != nil {
		return errors.Wrap(err, "error creating address")
	}
	var dest *bacnetip.Address
	// see if the destination was our broadcast address
	if pdu.GetPDUDestination().Equals(broadcast) {
		dest = bacnetip.NewLocalBroadcast(nil)
	} else {
		dest, err = bacnetip.NewAddress(s.log, pdu.GetPDUDestination().AddrAddress)
		if err != nil {
			return errors.Wrap(err, "error creating address")
		}
	}

	return s.Response(bacnetip.NewArgs(bacnetip.NewPDU(pdu, bacnetip.WithPDUSource(src), bacnetip.WithPDUDestination(dest))), bacnetip.NoKWArgs)
}

type SnifferStateMachine struct {
	*tests.ClientStateMachine

	address *bacnetip.Address
	annexj  *bacnetip.AnnexJCodec
	mux     *FauxMultiplexer

	log zerolog.Logger
}

func NewSnifferStateMachine(localLog zerolog.Logger, address string, vlan *bacnetip.IPNetwork) (*SnifferStateMachine, error) {
	s := &SnifferStateMachine{
		log: localLog,
	}
	machine, err := tests.NewClientStateMachine(localLog, tests.WithClientStateMachineName(address), tests.WithClientStateMachineExtension(s))
	if err != nil {
		return nil, errors.Wrap(err, "error building client state machine")
	}
	s.ClientStateMachine = machine

	// save the name and address
	s.address, err = bacnetip.NewAddress(localLog, address)
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}

	// BACnet/IP interpreter
	s.annexj, err = bacnetip.NewAnnexJCodec(localLog)
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
	if err := bacnetip.Bind(localLog, s, s.annexj, s.mux); err != nil {
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
	*tests.ClientStateMachine

	address *bacnetip.Address
	annexj  *bacnetip.AnnexJCodec
	mux     *FauxMultiplexer
}

func NewBIPStateMachine(localLog zerolog.Logger, address string, vlan *bacnetip.IPNetwork) (*BIPStateMachine, error) {
	b := &BIPStateMachine{}
	var err error
	b.ClientStateMachine, err = tests.NewClientStateMachine(localLog, tests.WithClientStateMachineName(address), tests.WithClientStateMachineExtension(b))
	if err != nil {
		return nil, errors.Wrap(err, "error building client state machine")
	}

	// save the name and address
	b.address = Address(address)

	// BACnet/IP interpreter
	b.annexj, err = bacnetip.NewAnnexJCodec(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating annexj")
	}

	// fake multiplexer has a VLAN node in it
	b.mux, err = NewFauxMultiplexer(localLog, b.address, vlan)

	// bind the stack together
	err = bacnetip.Bind(localLog, b, b.annexj, b.mux)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}
	return b, nil
}

type BIPSimpleStateMachine struct {
	*tests.ClientStateMachine
	name string

	address *bacnetip.Address

	bip    *bacnetip.BIPSimple
	annexj *bacnetip.AnnexJCodec
	mux    *FauxMultiplexer

	log zerolog.Logger
}

func NewBIPSimpleStateMachine(localLog zerolog.Logger, netstring string, vlan *bacnetip.IPNetwork) (*BIPSimpleStateMachine, error) {
	b := &BIPSimpleStateMachine{
		log: localLog,
	}
	var err error
	b.ClientStateMachine, err = tests.NewClientStateMachine(localLog, tests.WithClientStateMachineName(netstring), tests.WithClientStateMachineExtension(b))
	if err != nil {
		return nil, errors.Wrap(err, "error building client state machine")
	}

	// save the name and address
	b.address = Address(netstring)

	// BACnet/IP interpreter
	b.bip, err = bacnetip.NewBIPSimple(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building bip simple")
	}
	b.annexj, err = bacnetip.NewAnnexJCodec(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building annexj codec")
	}

	// fake multiplexer has a VLAN node in it
	b.mux, err = NewFauxMultiplexer(localLog, b.address, vlan)
	if err != nil {
		return nil, errors.Wrap(err, "error creating faux")
	}

	// bind the stack together
	if err := bacnetip.Bind(localLog, b, b.bip, b.annexj, b.mux); err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	return b, nil
}

// BIPForeignStateMachine  sits on a BIPForeign instance, the send() and receive()
//
//	parameters are NPDUs.
type BIPForeignStateMachine struct {
	*tests.ClientStateMachine

	address *bacnetip.Address
	bip     *bacnetip.BIPForeign
	annexj  *bacnetip.AnnexJCodec
	mux     *FauxMultiplexer

	log zerolog.Logger
}

func NewBIPForeignStateMachine(localLog zerolog.Logger, address string, vlan *bacnetip.IPNetwork) (*BIPForeignStateMachine, error) {
	b := &BIPForeignStateMachine{
		log: localLog,
	}
	var err error
	b.ClientStateMachine, err = tests.NewClientStateMachine(localLog, tests.WithClientStateMachineName(address), tests.WithClientStateMachineExtension(b))
	if err != nil {
		return nil, errors.New("error building client state machine")
	}

	// save the name and address
	b.address = Address(address)

	// BACnet/IP interpreter
	b.bip, err = bacnetip.NewBIPForeign(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating BIPForeign")
	}
	b.annexj, err = bacnetip.NewAnnexJCodec(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating AnnexJCodec")
	}

	// fake multiplexer has a VLAN node in it
	b.mux, err = NewFauxMultiplexer(localLog, b.address, vlan)
	if err != nil {
		return nil, errors.Wrap(err, "error creating FauxMultiplexer")
	}

	// bind the stack together
	err = bacnetip.Bind(b.log, b, b.bip, b.annexj, b.mux)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}
	return b, nil
}

type BIPBBMDStateMachine struct {
	*tests.ClientStateMachine

	address *bacnetip.Address
	bip     *bacnetip.BIPBBMD
	annexj  *bacnetip.AnnexJCodec
	mux     *FauxMultiplexer

	log zerolog.Logger
}

func NewBIPBBMDStateMachine(localLog zerolog.Logger, address string, vlan *bacnetip.IPNetwork) (*BIPBBMDStateMachine, error) {
	b := &BIPBBMDStateMachine{
		log: localLog,
	}
	var err error
	b.ClientStateMachine, err = tests.NewClientStateMachine(localLog, tests.WithClientStateMachineName(address), tests.WithClientStateMachineExtension(b))
	if err != nil {
		return nil, errors.New("error building client state machine")
	}

	// save the name and address
	b.address = Address(address)

	// BACnet/IP interpreter
	b.bip, err = bacnetip.NewBIPBBMD(localLog, b.address)
	b.annexj, err = bacnetip.NewAnnexJCodec(localLog)

	// build an address, full mask
	bdtAddress := fmt.Sprintf("%s/32:%d", b.address.AddrTuple.Left, b.address.AddrTuple.Right)
	b.log.Debug().Str("bdtAddress", bdtAddress).Msg("bdtAddress")

	// add itself as the first entry in the BDT
	if err := b.bip.AddPeer(Address(bdtAddress)); err != nil {
		return nil, errors.Wrap(err, "error adding peer")
	}

	// fake multiplexer has a VLAN node in it
	b.mux, err = NewFauxMultiplexer(localLog, b.address, vlan)
	if err != nil {
		return nil, errors.Wrap(err, "error creating faux multiplexer")
	}

	// bind the stack together
	err = bacnetip.Bind(b.log, b, b.bip, b.annexj, b.mux)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}
	return b, nil
}

type BIPSimpleNode struct {
	name    string
	address *bacnetip.Address
	bip     *bacnetip.BIPSimple
	annexj  *bacnetip.AnnexJCodec
	mux     *FauxMultiplexer
}

func NewBIPSimpleNode(localLog zerolog.Logger, address string, vlan *bacnetip.IPNetwork) (*BIPSimpleNode, error) {
	b := &BIPSimpleNode{}

	// save the name and address
	b.name = address
	b.address = Address(address)

	var err error
	// BACnet/IP interpreter
	b.bip, err = bacnetip.NewBIPSimple(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building bip simple")
	}
	b.annexj, err = bacnetip.NewAnnexJCodec(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building annexj codec")
	}

	// fake multiplexer has a VLAN node in it
	b.mux, err = NewFauxMultiplexer(localLog, b.address, vlan)

	// bind the stack together
	err = bacnetip.Bind(localLog, b.bip, b.annexj, b.mux)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	return b, nil
}

type BIPBBMDNode struct {
	name    string
	address *bacnetip.Address

	bip    *bacnetip.BIPBBMD
	annexj *bacnetip.AnnexJCodec
	mux    *FauxMultiplexer

	log zerolog.Logger
}

func NewBIPBBMDNode(localLog zerolog.Logger, address string, vlan *bacnetip.IPNetwork) (*BIPBBMDNode, error) {
	b := &BIPBBMDNode{
		log: localLog,
	}

	// build a name, save the address
	b.name = fmt.Sprintf("app @ %s", address)
	b.address = Address(address)
	b.log.Debug().Str("address", address).Msg("address")

	var err error
	// BACnet/IP interpreter
	b.bip, err = bacnetip.NewBIPBBMD(b.log, b.address)
	if err != nil {
		return nil, errors.Wrap(err, "error building bip bbmd")
	}
	b.annexj, err = bacnetip.NewAnnexJCodec(b.log)
	if err != nil {
		return nil, errors.Wrap(err, "error building annexj codec")
	}

	// build an address, full mask
	bdtAddress := fmt.Sprintf("%s/32:%d", b.address.AddrTuple.Left, b.address.AddrTuple.Right)
	b.log.Debug().Str("bdtAddress", bdtAddress).Msg("bdtAddress")

	// add itself as the first entry in the BDT
	err = b.bip.AddPeer(Address(bdtAddress))
	if err != nil {
		return nil, errors.Wrap(err, "error adding peer")
	}

	// fake multiplexer has a VLAN node in it
	b.mux, err = NewFauxMultiplexer(b.log, b.address, vlan)
	if err != nil {
		return nil, errors.Wrap(err, "error creating faux multiplexer")
	}

	// bind the stack together
	err = bacnetip.Bind(b.log, b.bip, b.annexj, b.mux)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	return b, nil
}

type TestDeviceObject struct {
	*bacnetip.LocalDeviceObject
}

type BIPSimpleApplicationLayerStateMachine struct {
	*bacnetip.ApplicationServiceElement
	*tests.ClientStateMachine

	log zerolog.Logger // TODO: move down

	name    string
	address *bacnetip.Address
	asap    *bacnetip.ApplicationServiceAccessPoint
	smap    *bacnetip.StateMachineAccessPoint
	nsap    *bacnetip.NetworkServiceAccessPoint
	nse     *_NetworkServiceElement
	bip     *bacnetip.BIPSimple
	annexj  *bacnetip.AnnexJCodec
	mux     *FauxMultiplexer
}

func NewBIPSimpleApplicationLayerStateMachine(localLog zerolog.Logger, address string, vlan *bacnetip.IPNetwork) (*BIPSimpleApplicationLayerStateMachine, error) {
	b := &BIPSimpleApplicationLayerStateMachine{}
	// build a name, save the address
	b.name = fmt.Sprintf("app @ %s", address)
	b.address = Address(address)

	// build a local device object
	localDevice := &TestDeviceObject{
		LocalDeviceObject: &bacnetip.LocalDeviceObject{
			ObjectName:       b.name,
			ObjectIdentifier: "device:998",
			VendorIdentifier: 999,
		},
	}

	var err error
	// continue with initialization
	b.ApplicationServiceElement, err = bacnetip.NewApplicationServiceElement(localLog, b)
	if err != nil {
		return nil, errors.Wrap(err, "error building application")
	}
	b.ClientStateMachine, err = tests.NewClientStateMachine(localLog, tests.WithClientStateMachineName(b.name), tests.WithClientStateMachineExtension(b))

	// include a application decoder
	b.asap, err = bacnetip.NewApplicationServiceAccessPoint(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building application service access point")
	}

	// pass the device object to the state machine access point so it
	// can know if it should support segmentation
	// the segmentation state machines need access to the same device
	// information cache as the application
	b.smap, err = bacnetip.NewStateMachineAccessPoint(localLog, localDevice.LocalDeviceObject, bacnetip.WithStateMachineAccessPointDeviceInfoCache(bacnetip.NewDeviceInfoCache(localLog))) //TODO: this is a indirection that wasn't intended... we don't use the annotation yet so that might be fine
	if err != nil {
		return nil, errors.Wrap(err, "error building state machine access point")
	}

	// a network service access point will be needed
	b.nsap, err = bacnetip.NewNetworkServiceAccessPoint(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}

	// give the NSAP a generic network layer service element
	b.nse, err = new_NetworkServiceElement(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service element")
	}
	err = bacnetip.Bind(localLog, b.nse, b.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	// bind the top layers
	err = bacnetip.Bind(localLog, b, b.asap, b.smap, b.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	// BACnet/IP interpreter
	b.bip, err = bacnetip.NewBIPSimple(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building bip bbmd")
	}
	b.annexj, err = bacnetip.NewAnnexJCodec(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building annexj codec")
	}

	// fake multiplexer has a VLAN node in it
	b.mux, err = NewFauxMultiplexer(localLog, b.address, vlan)
	if err != nil {
		return nil, errors.Wrap(err, "error building multiplexer")
	}

	// bind the stack together
	err = bacnetip.Bind(localLog, b.bip, b.annexj, b.mux)
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

func (b *BIPSimpleApplicationLayerStateMachine) Indication(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")
	return b.Receive(args, bacnetip.NoKWArgs)
}

func (b *BIPSimpleApplicationLayerStateMachine) Confirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Confirmation")
	return b.Receive(args, bacnetip.NoKWArgs)
}

type BIPBBMDApplication struct {
	*bacnetip.Application
	*bacnetip.WhoIsIAmServices
	*bacnetip.ReadWritePropertyServices

	name    string
	address *bacnetip.Address

	asap   *bacnetip.ApplicationServiceAccessPoint
	smap   *bacnetip.StateMachineAccessPoint
	nsap   *bacnetip.NetworkServiceAccessPoint
	nse    *_NetworkServiceElement
	bip    *bacnetip.BIPBBMD
	annexj *bacnetip.AnnexJCodec
	mux    *FauxMultiplexer

	log zerolog.Logger
}

func NewBIPBBMDApplication(localLog zerolog.Logger, address string, vlan *bacnetip.IPNetwork) (*BIPBBMDApplication, error) {
	b := &BIPBBMDApplication{
		log: localLog,
	}

	// build a name, save the address
	b.name = fmt.Sprintf("app @ %s", address)
	b.address = Address(address)

	// build a local device object
	localDevice := &TestDeviceObject{
		LocalDeviceObject: &bacnetip.LocalDeviceObject{
			ObjectName:       b.name,
			ObjectIdentifier: "device:999",
			VendorIdentifier: 999,
		},
	}

	var err error
	// continue with initialization
	b.Application, err = bacnetip.NewApplication(localLog, localDevice.LocalDeviceObject, b) //TODO: this is a indirection that wasn't intended... we don't use the annotation yet so that might be fine
	if err != nil {
		return nil, errors.Wrap(err, "error building application")
	}

	// include a application decoder
	b.asap, err = bacnetip.NewApplicationServiceAccessPoint(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building application service access point")
	}

	// pass the device object to the state machine access point so it
	// can know if it should support segmentation
	// the segmentation state machines need access to the same device
	// information cache as the application
	b.smap, err = bacnetip.NewStateMachineAccessPoint(localLog, localDevice.LocalDeviceObject, bacnetip.WithStateMachineAccessPointDeviceInfoCache(b.GetDeviceInfoCache())) //TODO: this is a indirection that wasn't intended... we don't use the annotation yet so that might be fine
	if err != nil {
		return nil, errors.Wrap(err, "error building state machine access point")
	}

	// a network service access point will be needed
	b.nsap, err = bacnetip.NewNetworkServiceAccessPoint(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}

	// give the NSAP a generic network layer service element
	b.nse, err = new_NetworkServiceElement(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service element")
	}
	err = bacnetip.Bind(localLog, b.nse, b.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	// bind the top layers
	err = bacnetip.Bind(localLog, b, b.asap, b.smap, b.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	// BACnet/IP interpreter
	b.bip, err = bacnetip.NewBIPBBMD(localLog, b.address)
	if err != nil {
		return nil, errors.Wrap(err, "error building bip bbmd")
	}

	b.annexj, err = bacnetip.NewAnnexJCodec(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building annexj codec")
	}

	// build an address, full mask
	bdtAddress := fmt.Sprintf("%s/32:%d", b.address.AddrTuple.Left, b.address.AddrTuple.Right)
	localLog.Debug().Str("bdtAddress", bdtAddress).Msg("bdtAddress")

	// add itself as the first entry in the BDT
	err = b.bip.AddPeer(Address(bdtAddress))
	if err != nil {
		return nil, errors.Wrap(err, "error adding peer")
	}

	// fake multiplexer has a VLAN node in it
	b.mux, err = NewFauxMultiplexer(localLog, b.address, vlan)
	if err != nil {
		return nil, errors.Wrap(err, "error building multiplexer")
	}

	// bind the stack together
	err = bacnetip.Bind(localLog, b.bip, b.annexj, b.mux)
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
