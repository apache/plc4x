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
	"strconv"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/appservice"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/netservice"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/state_machine"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/vlan"
)

//go:generate go tool plc4xGenerator -type=ApplicationLayerStateMachine -prefix=helpers_
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
	a.address, err = NewAddress(NA(address))
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}

	// build a local device object
	atoiAdrerss, err := strconv.Atoi(address)
	if err != nil {
		return nil, errors.Wrap(err, "error converting address")
	}
	localDevice, err := NewTestDeviceObject(NoArgs,
		NKW(KWObjectName, a.name,
			KWObjectIdentifier, primitivedata.ObjectIdentifierTuple{Left: "device", Right: atoiAdrerss},
			KWVendorIdentifier, 999,
		))
	if err != nil {
		return nil, errors.Wrap(err, "error creating device")
	}

	if LogTestNetwork {
		a.log.Debug().Stringer("address", a.address).Msg("address")
	}

	// continue with initialization
	a.ApplicationServiceElementContract, err = NewApplicationServiceElement(a.log)
	if err != nil {
		return nil, errors.Wrap(err, "error creating application service")
	}
	a.ClientStateMachine, err = NewClientStateMachine(a.log,
		WithClientStateMachineName(localDevice.GetObjectName()),
		WithClientStateMachineExtension(a),
		WithClientStateMachineName(a.name),
		WithLeafType(a),
	)
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
	a.smap, err = NewStateMachineAccessPoint(a.log, localDevice)
	if err != nil {
		return nil, errors.Wrap(err, "error creating state machine access point")
	}

	// the segmentation state machines need access to some device
	// information cache, usually shared with the application
	a.smap.SetDeviceInfoCache(NewDeviceInfoCache(a.log))

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

func (a *ApplicationLayerStateMachine) Indication(args Args, kwArgs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Indication")
	return a.Receive(args, NoKWArgs())
}

func (a *ApplicationLayerStateMachine) Confirmation(args Args, kwArgs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Confirmation")
	return a.Receive(args, NoKWArgs())
}
