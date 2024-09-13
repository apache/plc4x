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
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/vlan"
)

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
