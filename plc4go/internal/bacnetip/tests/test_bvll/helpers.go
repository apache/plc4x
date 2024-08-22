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

	"github.com/apache/plc4x/plc4go/internal/bacnetip"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/tests"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type _NetworkServiceElement struct {
	*bacnetip.NetworkServiceElement
}

func new_NetworkServiceElement(localLog zerolog.Logger) (*_NetworkServiceElement, error) {
	i := &_NetworkServiceElement{}

	// This class turns off the deferred startup function call that broadcasts
	// I-Am-Router-To-Network and Network-Number-Is messages.
	var err error
	i.NetworkServiceElement, err = bacnetip.NewNetworkServiceElement(localLog, nil, true)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service element")
	}
	return i, nil
}

type FauxMultiplexer struct {
	*bacnetip.Client
	*bacnetip.Server

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

	name    string
	address *bacnetip.Address
	annexj  *bacnetip.AnnexJCodec
	mux     *FauxMultiplexer

	log zerolog.Logger
}

func NewSnifferStateMachine(localLog zerolog.Logger, address string, vlan *bacnetip.IPNetwork) (*SnifferStateMachine, error) {
	s := &SnifferStateMachine{
		log: localLog,
	}
	machine, err := tests.NewClientStateMachine(localLog)
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

type BIPStateMachine struct {
	*tests.ClientStateMachine
}

// TODO: implement BIPStateMachine

type BIPSimpleStateMachine struct {
	*tests.ClientStateMachine
	name string

	address *bacnetip.Address

	bip    *bacnetip.BIPSimple
	annexj *bacnetip.AnnexJCodec
	mux    *FauxMultiplexer

	log zerolog.Logger
}

func NewBIPSimpleStateMachine(name string, localLog zerolog.Logger, netstring string, vlan *bacnetip.IPNetwork) (*BIPSimpleStateMachine, error) {
	address, err := bacnetip.NewAddress(localLog, netstring)
	if err != nil {
		return nil, errors.Wrap(err, "error building address")
	}
	stateMachine := &BIPSimpleStateMachine{
		// save the name and address
		name:    netstring,
		address: address,
		log:     localLog,
	}
	clientStateMachine, err := tests.NewClientStateMachine(localLog, tests.WithClientStateMachineName(name))
	if err != nil {
		return nil, errors.Wrap(err, "error building client state machine")
	}
	stateMachine.ClientStateMachine = clientStateMachine

	// BACnet/IP interpreter
	stateMachine.bip, err = bacnetip.NewBIPSimple(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building bip simple")
	}
	stateMachine.annexj, err = bacnetip.NewAnnexJCodec(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building annexj codec")
	}

	// fake multiplexer has a VLAN node in it
	stateMachine.mux, err = NewFauxMultiplexer(localLog, stateMachine.address, vlan)
	if err != nil {
		return nil, errors.Wrap(err, "error creating faux")
	}

	// bind the stack together
	if err := bacnetip.Bind(localLog, stateMachine, stateMachine.bip, stateMachine.annexj, stateMachine.mux); err != nil {
		return nil, errors.Wrap(err, "error binding")
	}

	return stateMachine, nil
}

type BIPForeignStateMachine struct {
	*tests.ClientStateMachine
}

type BIPBBMDStateMachine struct {
	*tests.ClientStateMachine
}

type BIPSimpleNode struct {
}

type BIPBBMDNode struct {
}

type TestDeviceObject struct {
	*bacnetip.LocalDeviceObject
}

type BIPSimpleApplicationLayerStateMachine struct {
	*bacnetip.ApplicationServiceElement
	*tests.ClientStateMachine
}

type BIPBBMDApplication struct {
	*bacnetip.Application
	*bacnetip.WhoIsIAmServices
	*bacnetip.ReadWritePropertyServices
}
