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

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/appservice"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/bvllservice"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/netservice"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/state_machine"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/vlan"
)

//go:generate go tool plc4xGenerator -type=BIPSimpleApplicationLayerStateMachine -prefix=helpers_
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
	localDevice, err := NewTestDeviceObject(
		NoArgs,
		NKW(
			KWObjectName, b.name,
			KWObjectIdentifier, "device:998",
			KWVendorIdentifier, 999,
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error creating test device")
	}

	// continue with initialization
	b.ApplicationServiceElementContract, err = NewApplicationServiceElement(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building application")
	}
	b.ClientStateMachine, err = NewClientStateMachine(localLog, WithClientStateMachineName(b.name), WithClientStateMachineExtension(b), WithLeafType(b))

	// include a application decoder
	b.asap, err = NewApplicationServiceAccessPoint(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error building application service access point")
	}

	// pass the device object to the state machine access point so it
	// can know if it should support segmentation
	b.smap, err = NewStateMachineAccessPoint(localLog, localDevice)
	if err != nil {
		return nil, errors.Wrap(err, "error building state machine access point")
	}

	// the segmentation state machines need access to the same device
	// information cache as the application
	b.smap.SetDeviceInfoCache(NewDeviceInfoCache(localLog))

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

func (b *BIPSimpleApplicationLayerStateMachine) Indication(args Args, kwArgs KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Indication")
	return b.Receive(args, NoKWArgs())
}

func (b *BIPSimpleApplicationLayerStateMachine) Confirmation(args Args, kwArgs KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Confirmation")
	return b.Receive(args, NoKWArgs())
}
