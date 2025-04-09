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
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/netservice"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/service"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/vlan"
)

//go:generate go tool plc4xGenerator -type=ApplicationNode -prefix=helpers_
type ApplicationNode struct {
	*Application
	*WhoIsIAmServices
	*ReadWritePropertyServices
	*DefaultRFormatter `ignore:"true"`

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
		DefaultRFormatter: NewDefaultRFormatter(),
		log:               localLog,
	}

	// build a name, save the address
	a.name = fmt.Sprintf("app @ %s", address)
	var err error
	a.address, err = NewAddress(NA(address))
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}

	// build a local device object
	localDevice, err := NewTestDeviceObject(NoArgs,
		NKW(
			KWObjectName, a.name,
			KWObjectIdentifier, primitivedata.ObjectIdentifierTuple{Left: "device", Right: 999},
			KWVendorIdentifier, 999,
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error creating device")
	}

	// continue with initialization
	a.Application, err = NewApplication(localLog, WithApplicationLocalDeviceObject(localDevice))
	if err != nil {
		return nil, errors.Wrap(err, "error building application")
	}

	a.WhoIsIAmServices, err = NewWhoIsIAmServices(localLog, a, WithWhoIsIAmServicesLocalDevice(localDevice))
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
	a.smap, err = NewStateMachineAccessPoint(localLog, localDevice)
	if err != nil {
		return nil, errors.Wrap(err, "error building state machine access point")
	}

	// the segmentation state machines need access to the same device
	// information cache as the application
	a.smap.SetDeviceInfoCache(a.GetDeviceInfoCache())

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

func (a *ApplicationNode) AlternateString() (string, bool) {
	if IsDebuggingActive() {
		return fmt.Sprintf("%r", a), true
	}
	return "", false
}
