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
	*bacnetip.Client
	*bacnetip.Server

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

	name    string
	address *bacnetip.Address
	node    *bacnetip.Node

	log zerolog.Logger
}

func NewSnifferStateMachine(localLog zerolog.Logger, address string, vlan *bacnetip.IPNetwork) (*SnifferStateMachine, error) {
	s := &SnifferStateMachine{
		log: localLog,
	}
	machine, err := tests.NewClientStateMachine(localLog, tests.WithClientStateMachineName(address))
	if err != nil {
		return nil, errors.Wrap(err, "error building client state machine")
	}
	s.ClientStateMachine = machine

	// save the name and address
	s.name = address
	s.address, err = bacnetip.NewAddress(localLog, address)
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}

	// create a promiscuous node, added to the network
	s.node, err = bacnetip.NewNode(s.log, s.address, bacnetip.WithNodePromiscuous(true))
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

func (s *SnifferStateMachine) String() string {
	return fmt.Sprintf("SnifferStateMachine(%s)", s.name)
}

type NetworkLayerStateMachine struct {
	*tests.ClientStateMachine

	name    string
	address *bacnetip.Address

	log   zerolog.Logger
	codec *NPDUCodec
	node  *bacnetip.Node
}

func NewNetworkLayerStateMachine(localLog zerolog.Logger, address string, vlan *bacnetip.IPNetwork) (*NetworkLayerStateMachine, error) {
	n := &NetworkLayerStateMachine{
		log: localLog,
	}

	// save the name and address
	n.name = fmt.Sprintf("app @ %s", address)
	n.address = Address(address)

	var err error
	n.ClientStateMachine, err = tests.NewClientStateMachine(localLog, tests.WithClientStateMachineName(n.name))
	if err != nil {
		return nil, errors.Wrap(err, "error building client state machine")
	}
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

//
