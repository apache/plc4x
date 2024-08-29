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

package test_network

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

type NPDUCodec struct {
	bacnetip.Client
	bacnetip.Server

	log zerolog.Logger
}

func NewNPDUCodec(localLog zerolog.Logger) (*NPDUCodec, error) {
	n := &NPDUCodec{
		log: localLog,
	}
	var err error
	n.Client, err = bacnetip.NewClient(localLog, n)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	n.Server, err = bacnetip.NewServer(localLog, n)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	return n, nil
}

func (n *NPDUCodec) Indication(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	n.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")

	npdu := args.Get0NPDU()

	// first a generic _NPDU
	xpdu, err := bacnetip.NewNPDU(nil, nil)
	if err != nil {
		return errors.Wrap(err, "error creating NPDU")
	}
	if err := npdu.Encode(xpdu); err != nil {
		return errors.Wrap(err, "error encoding xpdu")
	}

	// Now as a vanilla PDU
	ypdu := bacnetip.NewPDU(bacnetip.NewMessageBridge())
	if err := xpdu.Encode(ypdu); err != nil {
		return errors.Wrap(err, "error decoding xpdu")
	}
	n.log.Debug().Stringer("ypdu", ypdu).Msg("encoded")

	// send it downstream
	return n.Request(bacnetip.NewArgs(ypdu), bacnetip.NoKWArgs)
}

func (n *NPDUCodec) Confirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	n.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")

	pdu := args.Get0PDU()

	// decode as generic _NPDU
	xpdu, err := bacnetip.NewNPDU(nil, nil)
	if err != nil {
		return errors.Wrap(err, "error creating NPDU")
	}
	if err := xpdu.Decode(pdu); err != nil {
		return errors.Wrap(err, "error decoding xpdu")
	}

	// drop application layer message
	if xpdu.GetNPDUNetMessage() == nil {
		n.log.Trace().Msg("drop message")
		return nil
	}

	// do a deeper decode of the _NPDU
	ypdu := bacnetip.NPDUTypes[*xpdu.GetNPDUNetMessage()]()
	if err := ypdu.Decode(xpdu); err != nil {
		return errors.Wrap(err, "error decoding ypdu")
	}

	return n.Response(bacnetip.NewArgs(ypdu), bacnetip.NoKWArgs)
}

func (n *NPDUCodec) String() string {
	return "NPDUCodec"
}

type SnifferStateMachine struct {
	*tests.ClientStateMachine

	address *bacnetip.Address
	node    *bacnetip.Node

	log zerolog.Logger
}

func NewSnifferStateMachine(localLog zerolog.Logger, address string, vlan *bacnetip.Network) (*SnifferStateMachine, error) {
	s := &SnifferStateMachine{
		log: localLog,
	}
	var err error
	s.ClientStateMachine, err = tests.NewClientStateMachine(localLog, tests.WithClientStateMachineName(address), tests.WithClientStateMachineExtension(s))
	if err != nil {
		return nil, errors.Wrap(err, "error building client state machine")
	}

	// save the name and address
	s.address, err = bacnetip.NewAddress(localLog, address)
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}

	// create a promiscuous node, added to the network
	s.node, err = bacnetip.NewNode(s.log, s.address, bacnetip.WithNodePromiscuous(true), bacnetip.WithNodeLan(vlan))
	if err != nil {
		return nil, errors.Wrap(err, "error creating node")
	}
	s.log.Debug().Stringer("node", s.node).Msg("node")

	// bind the stack together
	if err := bacnetip.Bind(localLog, s, s.node); err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	return s, nil
}

type NetworkLayerStateMachine struct {
	*tests.ClientStateMachine

	address *bacnetip.Address

	log   zerolog.Logger
	codec *NPDUCodec
	node  *bacnetip.Node
}

func NewNetworkLayerStateMachine(localLog zerolog.Logger, address string, vlan *bacnetip.Network) (*NetworkLayerStateMachine, error) {
	n := &NetworkLayerStateMachine{
		log: localLog,
	}
	var err error
	n.ClientStateMachine, err = tests.NewClientStateMachine(localLog, tests.WithClientStateMachineName(address), tests.WithClientStateMachineExtension(n))
	if err != nil {
		return nil, errors.Wrap(err, "error building client state machine")
	}

	// save the name and address
	n.address = Address(address)

	// create a network layer encoder/decoder
	n.codec, err = NewNPDUCodec(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating codec")
	}
	n.log.Debug().Stringer("codec", n.codec).Msg("codec")

	// create a node, added to the network
	n.node, err = bacnetip.NewNode(localLog, n.address, bacnetip.WithNodeLan(vlan))
	if err != nil {
		return nil, errors.Wrap(err, "error creating node")
	}
	n.log.Debug().Stringer("node", n.node).Msg("node")

	// bind this to the node
	if err := bacnetip.Bind(localLog, n, n.codec, n.node); err != nil {
		return nil, errors.Wrap(err, "error binding")
	}
	return n, nil
}

type RouterNode struct {
	nsap *bacnetip.NetworkServiceAccessPoint
	nse  *_NetworkServiceElement

	log zerolog.Logger
}

func NewRouterNode(localLog zerolog.Logger) (*RouterNode, error) {
	r := &RouterNode{log: localLog}
	var err error
	// a network service access point will be needed
	r.nsap, err = bacnetip.NewNetworkServiceAccessPoint(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}
	// give the NSAP a generic network layer service element
	r.nse, err = new_NetworkServiceElement(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service element")
	}
	err = bacnetip.Bind(localLog, r.nse, r.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}
	return r, nil
}

func (r *RouterNode) AddNetwork(address string, vlan *bacnetip.Network, net uint16) error {
	r.log.Debug().Str("address", address).Stringer("vlan", vlan).Uint16("net", net).Msg("AddNetwork")

	// convert the address to an Address
	addr := Address(address)

	// create a node, add to the network
	node, err := bacnetip.NewNode(r.log, addr, bacnetip.WithNodeLan(vlan))
	if err != nil {
		return errors.Wrap(err, "error creating node")
	}

	// bind the BIP stack to the local network
	return r.nsap.Bind(node, &net, addr)
}

func (r *RouterNode) String() string {
	return fmt.Sprintf("RouterNode")
}

type RouterStateMachine struct {
	*RouterNode
	tests.StateMachineContract
}

func NewRouterStateMachine(localLog zerolog.Logger) (*RouterStateMachine, error) {
	r := &RouterStateMachine{}
	var err error
	r.RouterNode, err = NewRouterNode(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating router node")
	}
	var initFunc func()
	r.StateMachineContract, initFunc = tests.NewStateMachine(localLog, r)
	initFunc()
	return r, nil
}

func (r *RouterStateMachine) Send(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	panic("not available")
}

func (r *RouterStateMachine) String() string {
	return "RouterStateMachine"
}

type TestDeviceObject struct {
	*bacnetip.LocalDeviceObject
}

type ApplicationLayerStateMachine struct {
	bacnetip.ApplicationServiceElementContract
	*tests.ClientStateMachine

	name    string
	address *bacnetip.Address

	asap *bacnetip.ApplicationServiceAccessPoint
	smap *bacnetip.StateMachineAccessPoint
	nsap *bacnetip.NetworkServiceAccessPoint
	nse  *_NetworkServiceElement
	node *bacnetip.Node

	log zerolog.Logger
}

func NewApplicationLayerStateMachine(localLog zerolog.Logger, address string, vlan *bacnetip.Network) (*ApplicationLayerStateMachine, error) {
	a := &ApplicationLayerStateMachine{
		log: localLog,
	}

	// save the name and address
	a.name = fmt.Sprintf("app @ %s", address)
	a.address = Address(address)

	// build a local device object
	localDevice := TestDeviceObject{
		&bacnetip.LocalDeviceObject{
			ObjectName:       a.name,
			ObjectIdentifier: "device:" + address,
			VendorIdentifier: 999,
		},
	}

	a.log.Debug().Stringer("address", a.address).Msg("address")

	var err error
	// continue with initialization
	a.ApplicationServiceElementContract, err = bacnetip.NewApplicationServiceElement(a.log)
	if err != nil {
		return nil, errors.Wrap(err, "error creating application service")
	}
	a.ClientStateMachine, err = tests.NewClientStateMachine(a.log, tests.WithClientStateMachineName(localDevice.ObjectName), tests.WithClientStateMachineExtension(a))
	if err != nil {
		return nil, errors.Wrap(err, "error building client state machine")
	}

	// include a application decoder
	a.asap, err = bacnetip.NewApplicationServiceAccessPoint(a.log)
	if err != nil {
		return nil, errors.Wrap(err, "error creating application service access point")
	}

	// pass the device object to the state machine access point so it
	// can know if it should support segmentation
	// the segmentation state machines need access to some device
	// information cache, usually shared with the application
	a.smap, err = bacnetip.NewStateMachineAccessPoint(a.log, localDevice.LocalDeviceObject, bacnetip.WithStateMachineAccessPointDeviceInfoCache(bacnetip.NewDeviceInfoCache(a.log))) // TODO: this is not quite right as we unwrap here
	if err != nil {
		return nil, errors.Wrap(err, "error creating state machine access point")
	}

	//  a network service access point will be needed
	a.nsap, err = bacnetip.NewNetworkServiceAccessPoint(a.log)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}

	//  give the NSAP a generic network layer service element
	a.nse, err = new_NetworkServiceElement(a.log)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service element")
	}
	err = bacnetip.Bind(a.log, a.nse, a.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	//  bind the top layers
	err = bacnetip.Bind(a.log, a, a.asap, a.smap, a.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	//  create a node, added to the network
	a.node, err = bacnetip.NewNode(a.log, a.address, bacnetip.WithNodeLan(vlan))
	if err != nil {
		return nil, errors.Wrap(err, "error creating node")
	}
	a.log.Debug().Stringer("node", a.node).Msg("node")

	//  bind the stack to the local network
	err = a.nsap.Bind(a.node, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	return a, nil
}

func (a *ApplicationLayerStateMachine) Indication(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")
	return a.Receive(args, bacnetip.NoKWArgs)
}

func (a *ApplicationLayerStateMachine) Confirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Confirmation")
	return a.Receive(args, bacnetip.NoKWArgs)
}

type ApplicationNode struct {
	*bacnetip.Application
	*bacnetip.WhoIsIAmServices
	*bacnetip.ReadWritePropertyServices

	name    string
	address *bacnetip.Address
	asap    *bacnetip.ApplicationServiceAccessPoint
	smap    *bacnetip.StateMachineAccessPoint
	nsap    *bacnetip.NetworkServiceAccessPoint
	nse     *_NetworkServiceElement
	node    *bacnetip.Node

	log zerolog.Logger
}

func NewApplicationNode(localLog zerolog.Logger, address string, vlan *bacnetip.Network) (*ApplicationNode, error) {
	a := &ApplicationNode{
		log: localLog,
	}

	// build a name, save the address
	a.name = fmt.Sprintf("app @ %s", address)
	a.address = Address(address)

	// build a local device object
	localDevice := &TestDeviceObject{
		LocalDeviceObject: &bacnetip.LocalDeviceObject{
			ObjectName:       a.name,
			ObjectIdentifier: "device:999",
			VendorIdentifier: 999,
		},
	}

	var err error
	// continue with initialization
	a.Application, err = bacnetip.NewApplication(localLog, localDevice.LocalDeviceObject) //TODO: this is a indirection that wasn't intended... we don't use the annotation yet so that might be fine
	if err != nil {
		return nil, errors.Wrap(err, "error building application")
	}

	// include a application decoder
	a.asap, err = bacnetip.NewApplicationServiceAccessPoint(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building application service access point")
	}

	// pass the device object to the state machine access point so it
	// can know if it should support segmentation
	// the segmentation state machines need access to the same device
	// information cache as the application
	a.smap, err = bacnetip.NewStateMachineAccessPoint(localLog, localDevice.LocalDeviceObject, bacnetip.WithStateMachineAccessPointDeviceInfoCache(a.GetDeviceInfoCache())) //TODO: this is a indirection that wasn't intended... we don't use the annotation yet so that might be fine
	if err != nil {
		return nil, errors.Wrap(err, "error building state machine access point")
	}

	// a network service access point will be needed
	a.nsap, err = bacnetip.NewNetworkServiceAccessPoint(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}

	// give the NSAP a generic network layer service element
	a.nse, err = new_NetworkServiceElement(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service element")
	}
	err = bacnetip.Bind(localLog, a.nse, a.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	// bind the top layers
	err = bacnetip.Bind(localLog, a, a.asap, a.smap, a.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	// create a node, added to the network
	a.node, err = bacnetip.NewNode(a.log, a.address, bacnetip.WithNodeLan(vlan))
	if err != nil {
		return nil, errors.Wrap(err, "error creating node")
	}

	// bind the stack to the local network
	err = a.nsap.Bind(a.node, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	return a, nil
}

func xtob(s string) []byte {
	bytes, err := bacnetip.Xtob(s)
	if err != nil {
		panic(err)
	}
	return bytes
}
