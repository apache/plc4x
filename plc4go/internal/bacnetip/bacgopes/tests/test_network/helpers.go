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

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/app"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/appservice"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/deleteme"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/local/device"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/netservice"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/npdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/object"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/service"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests"
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

//go:generate plc4xGenerator -type=NPDUCodec -prefix=
type NPDUCodec struct {
	Client
	Server

	log zerolog.Logger
}

func NewNPDUCodec(localLog zerolog.Logger) (*NPDUCodec, error) {
	n := &NPDUCodec{
		log: localLog,
	}
	var err error
	n.Client, err = NewClient(localLog, n)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	n.Server, err = NewServer(localLog, n)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	if !LogTestNetwork {
		n.log = zerolog.Nop()
	}
	return n, nil
}

func (n *NPDUCodec) Indication(args Args, kwargs KWArgs) error {
	n.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")

	npdu := Get[NPDU](args, 0)

	// first a generic _NPDU
	xpdu, err := NewNPDU(nil, nil)
	if err != nil {
		return errors.Wrap(err, "error creating NPDU")
	}
	if err := npdu.Encode(xpdu); err != nil {
		return errors.Wrap(err, "error encoding xpdu")
	}

	// Now as a vanilla PDU
	ypdu := NewPDU(NewMessageBridge())
	if err := xpdu.Encode(ypdu); err != nil {
		return errors.Wrap(err, "error decoding xpdu")
	}
	n.log.Debug().Stringer("ypdu", ypdu).Msg("encoded")

	// send it downstream
	return n.Request(NewArgs(ypdu), NoKWArgs)
}

func (n *NPDUCodec) Confirmation(args Args, kwargs KWArgs) error {
	n.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")

	pdu := Get[PDU](args, 0)

	// decode as generic _NPDU
	xpdu, err := NewNPDU(nil, nil)
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
	ypdu := NPDUTypes[*xpdu.GetNPDUNetMessage()]()
	if err := ypdu.Decode(xpdu); err != nil {
		return errors.Wrap(err, "error decoding ypdu")
	}

	return n.Response(NewArgs(ypdu), NoKWArgs)
}

type SnifferStateMachine struct {
	*ClientStateMachine

	address *Address
	node    *Node

	log zerolog.Logger
}

func NewSnifferStateMachine(localLog zerolog.Logger, address string, vlan *Network) (*SnifferStateMachine, error) {
	s := &SnifferStateMachine{
		log: localLog,
	}
	var err error
	s.ClientStateMachine, err = NewClientStateMachine(s.log, WithClientStateMachineName(address), WithClientStateMachineExtension(s))
	if err != nil {
		return nil, errors.Wrap(err, "error building client state machine")
	}

	// save the name and address
	s.address, err = NewAddress(s.log, address)
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}

	// create a promiscuous node, added to the network
	s.node, err = NewNode(s.log, s.address, WithNodePromiscuous(true), WithNodeLan(vlan))
	if err != nil {
		return nil, errors.Wrap(err, "error creating node")
	}
	if LogTestNetwork {
		s.log.Debug().Stringer("node", s.node).Msg("node")
	}

	// bind the stack together
	if err := Bind(s.log, s, s.node); err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	if !LogTestNetwork {
		s.log = zerolog.Nop()
	}
	return s, nil
}

type NetworkLayerStateMachine struct {
	*ClientStateMachine

	address *Address

	log   zerolog.Logger
	codec *NPDUCodec
	node  *Node
}

func NewNetworkLayerStateMachine(localLog zerolog.Logger, address string, vlan *Network) (*NetworkLayerStateMachine, error) {
	n := &NetworkLayerStateMachine{
		log: localLog,
	}
	var err error
	n.ClientStateMachine, err = NewClientStateMachine(localLog, WithClientStateMachineName(address), WithClientStateMachineExtension(n))
	if err != nil {
		return nil, errors.Wrap(err, "error building client state machine")
	}

	// save the name and address
	n.address, err = NewAddress(localLog, address)
	if err != nil {
		return nil, errors.Wrap(err, "error creaing address")
	}

	// create a network layer encoder/decoder
	n.codec, err = NewNPDUCodec(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating codec")
	}
	if LogTestNetwork {
		n.log.Debug().Stringer("codec", n.codec).Msg("codec")
	}

	// create a node, added to the network
	n.node, err = NewNode(localLog, n.address, WithNodeLan(vlan))
	if err != nil {
		return nil, errors.Wrap(err, "error creating node")
	}
	if LogTestNetwork {
		n.log.Debug().Stringer("node", n.node).Msg("node")
	}

	// bind this to the node
	if err := Bind(localLog, n, n.codec, n.node); err != nil {
		return nil, errors.Wrap(err, "error binding")
	}
	if !LogTestNetwork {
		n.log = zerolog.Nop()
	}
	return n, nil
}

type RouterNode struct {
	nsap *NetworkServiceAccessPoint
	nse  *_NetworkServiceElement

	log zerolog.Logger
}

func NewRouterNode(localLog zerolog.Logger) (*RouterNode, error) {
	r := &RouterNode{log: localLog}
	var err error
	// a network service access point will be needed
	r.nsap, err = NewNetworkServiceAccessPoint(r.log)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}
	// give the NSAP a generic network layer service element
	r.nse, err = new_NetworkServiceElement(r.log)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service element")
	}
	err = Bind(r.log, r.nse, r.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}
	if !LogTestNetwork {
		r.log = zerolog.Nop()
	}
	return r, nil
}

func (r *RouterNode) AddNetwork(address string, vlan *Network, net uint16) error {
	r.log.Debug().Str("address", address).Stringer("vlan", vlan).Uint16("net", net).Msg("AddNetwork")

	// convert the address to an Address
	addr, err := NewAddress(r.log, address)
	if err != nil {
		return errors.Wrap(err, "error creaing address")
	}

	// create a node, add to the network
	node, err := NewNode(r.log, addr, WithNodeLan(vlan))
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
	StateMachineContract
}

func NewRouterStateMachine(localLog zerolog.Logger) (*RouterStateMachine, error) {
	r := &RouterStateMachine{}
	var err error
	r.RouterNode, err = NewRouterNode(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating router node")
	}
	var initFunc func()
	r.StateMachineContract, initFunc = NewStateMachine(localLog, r)
	initFunc()
	if !LogTestNetwork {
		r.log = zerolog.Nop()
	}
	return r, nil
}

func (r *RouterStateMachine) Send(args Args, kwargs KWArgs) error {
	panic("not available")
}

func (r *RouterStateMachine) String() string {
	return "RouterStateMachine"
}

type TestDeviceObject struct {
	*LocalDeviceObject
}

//go:generate plc4xGenerator -type=ApplicationLayerStateMachine
type ApplicationLayerStateMachine struct {
	ApplicationServiceElementContract
	*ClientStateMachine `ignore:"true"` // TODO: add support

	name    string
	address *Address

	asap *ApplicationServiceAccessPoint
	smap *StateMachineAccessPoint
	nsap *NetworkServiceAccessPoint
	nse  *_NetworkServiceElement
	node *Node

	log zerolog.Logger
}

func NewApplicationLayerStateMachine(localLog zerolog.Logger, address string, vlan *Network) (*ApplicationLayerStateMachine, error) {
	a := &ApplicationLayerStateMachine{
		log: localLog,
	}

	// save the name and address
	a.name = fmt.Sprintf("app @ %s", address)
	var err error
	a.address, err = NewAddress(localLog, address)
	if err != nil {
		return nil, errors.Wrap(err, "error creaing address")
	}

	// build a local device object
	localDevice := TestDeviceObject{
		&LocalDeviceObject{
			ObjectName:       a.name,
			ObjectIdentifier: "device:" + address,
			VendorIdentifier: 999,
		},
	}

	if LogTestNetwork {
		a.log.Debug().Stringer("address", a.address).Msg("address")
	}

	// continue with initialization
	a.ApplicationServiceElementContract, err = NewApplicationServiceElement(a.log)
	if err != nil {
		return nil, errors.Wrap(err, "error creating application service")
	}
	a.ClientStateMachine, err = NewClientStateMachine(a.log, WithClientStateMachineName(localDevice.ObjectName), WithClientStateMachineExtension(a))
	if err != nil {
		return nil, errors.Wrap(err, "error building client state machine")
	}

	// include a application decoder
	a.asap, err = NewApplicationServiceAccessPoint(a.log)
	if err != nil {
		return nil, errors.Wrap(err, "error creating application service access point")
	}

	// pass the device object to the state machine access point so it
	// can know if it should support segmentation
	// the segmentation state machines need access to some device
	// information cache, usually shared with the application
	a.smap, err = NewStateMachineAccessPoint(a.log, localDevice.LocalDeviceObject, WithStateMachineAccessPointDeviceInfoCache(NewDeviceInfoCache(a.log))) // TODO: this is not quite right as we unwrap here
	if err != nil {
		return nil, errors.Wrap(err, "error creating state machine access point")
	}

	//  a network service access point will be needed
	a.nsap, err = NewNetworkServiceAccessPoint(a.log)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}

	//  give the NSAP a generic network layer service element
	a.nse, err = new_NetworkServiceElement(a.log)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service element")
	}
	err = Bind(a.log, a.nse, a.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	//  bind the top layers
	err = Bind(a.log, a, a.asap, a.smap, a.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	//  create a node, added to the network
	a.node, err = NewNode(a.log, a.address, WithNodeLan(vlan))
	if err != nil {
		return nil, errors.Wrap(err, "error creating node")
	}
	if LogTestNetwork {
		a.log.Debug().Stringer("node", a.node).Msg("node")
	}

	//  bind the stack to the local network
	err = a.nsap.Bind(a.node, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	if !LogTestNetwork {
		a.log = zerolog.Nop()
	}
	return a, nil
}

func (a *ApplicationLayerStateMachine) Indication(args Args, kwargs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")
	return a.Receive(args, NoKWArgs)
}

func (a *ApplicationLayerStateMachine) Confirmation(args Args, kwargs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Confirmation")
	return a.Receive(args, NoKWArgs)
}

//go:generate plc4xGenerator -type=ApplicationNode
type ApplicationNode struct {
	*Application
	*WhoIsIAmServices
	*ReadWritePropertyServices

	name    string
	address *Address `directSerialize:"true"`
	asap    *ApplicationServiceAccessPoint
	smap    *StateMachineAccessPoint
	nsap    *NetworkServiceAccessPoint
	nse     *_NetworkServiceElement
	node    *Node

	log zerolog.Logger
}

func NewApplicationNode(localLog zerolog.Logger, address string, vlan *Network) (*ApplicationNode, error) {
	a := &ApplicationNode{
		log: localLog,
	}

	// build a name, save the address
	a.name = fmt.Sprintf("app @ %s", address)
	var err error
	a.address, err = NewAddress(localLog, address)
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}

	// build a local device object
	localDevice := &TestDeviceObject{
		LocalDeviceObject: &LocalDeviceObject{
			ObjectName:       a.name,
			ObjectIdentifier: "device:999",
			VendorIdentifier: 999,
		},
	}

	// continue with initialization
	a.Application, err = NewApplication(localLog, localDevice.LocalDeviceObject) //TODO: this is a indirection that wasn't intended... we don't use the annotation yet so that might be fine
	if err != nil {
		return nil, errors.Wrap(err, "error building application")
	}

	a.WhoIsIAmServices, err = NewWhoIsIAmServices(localLog, a, WithWhoIsIAmServicesLocalDevice(localDevice.LocalDeviceObject)) //TODO: this is a indirection that wasn't intended... we don't use the annotation yet so that might be fine
	if err != nil {
		return nil, errors.Wrap(err, "error building WhoIsIAmServices")
	}

	// include a application decoder
	a.asap, err = NewApplicationServiceAccessPoint(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building application service access point")
	}

	// pass the device object to the state machine access point so it
	// can know if it should support segmentation
	// the segmentation state machines need access to the same device
	// information cache as the application
	a.smap, err = NewStateMachineAccessPoint(localLog, localDevice.LocalDeviceObject, WithStateMachineAccessPointDeviceInfoCache(a.GetDeviceInfoCache())) //TODO: this is a indirection that wasn't intended... we don't use the annotation yet so that might be fine
	if err != nil {
		return nil, errors.Wrap(err, "error building state machine access point")
	}

	// a network service access point will be needed
	a.nsap, err = NewNetworkServiceAccessPoint(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}

	// give the NSAP a generic network layer service element
	a.nse, err = new_NetworkServiceElement(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service element")
	}
	err = Bind(localLog, a.nse, a.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	// bind the top layers
	err = Bind(localLog, a, a.asap, a.smap, a.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	// create a node, added to the network
	a.node, err = NewNode(a.log, a.address, WithNodeLan(vlan))
	if err != nil {
		return nil, errors.Wrap(err, "error creating node")
	}

	// bind the stack to the local network
	err = a.nsap.Bind(a.node, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}
	if !LogTestNetwork {
		a.log = zerolog.Nop()
	}
	return a, nil
}

func xtob(s string) []byte {
	bytes, err := Xtob(s)
	if err != nil {
		panic(err)
	}
	return bytes
}
