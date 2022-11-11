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

package bacnetip

import "github.com/pkg/errors"

// TODO: implement
type Application struct {
	ApplicationServiceElement
	Collector
}

// TODO: implement
type IOController struct {
}

// TODO: implement
type ApplicationIOController struct {
	IOController
	Application
}

func NewApplicationIOController(interface{}, interface{}, interface{}, *int) (*ApplicationIOController, error) {
	return &ApplicationIOController{}, nil
}

type BIPSimpleApplication struct {
	*ApplicationIOController
	*WhoIsIAmServices
	*ReadWritePropertyServices
	localAddress interface{}
	asap         *ApplicationServiceAccessPoint
	smap         *StateMachineAccessPoint
	nsap         *NetworkServiceAccessPoint
	nse          *NetworkServiceElement
}

func NewBIPSimpleApplication(localDevice DeviceEntry, localAddress, deviceInfoCache *DeviceInventory, aseID *int) (*BIPSimpleApplication, error) {
	b := &BIPSimpleApplication{}
	controller, err := NewApplicationIOController(localDevice, localAddress, deviceInfoCache, aseID)
	if err != nil {
		return nil, errors.Wrap(err, "error creating io controller")
	}
	b.ApplicationIOController = controller

	b.localAddress = localAddress

	// include a application decoder
	applicationServiceAccessPoint, err := NewApplicationServiceAccessPoint(nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating application service access point")
	}
	b.asap = applicationServiceAccessPoint

	// pass the device object to the state machine access point, so it can know if it should support segmentation
	stateMachineAccessPoint, err := NewStateMachineAccessPoint(localDevice, deviceInfoCache, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating state machine access point")
	}
	b.smap = stateMachineAccessPoint

	// pass the device object to the state machine access point so it # can know if it should support segmentation
	// Note: deviceInfoCache already passed above so we don't need to do it again here

	// a network service access point will be needed
	networkServiceAccessPoint, err := NewNetworkServiceAccessPoint()
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}
	b.nsap = networkServiceAccessPoint

	// give the NSAP a generic network layer service element
	networkServiceElement, err := NewNetworkServiceElement()
	if err != nil {
		return nil, errors.Wrap(err, "error creating new network service element")
	}
	b.nse = networkServiceElement
	if err := bind(b.nse, b.nsap); err != nil {
		return nil, errors.New("error binding network stack")
	}

	// bind the top layers
	if err := bind(b, b.asap, b.smap, b.nsap); err != nil {
		return nil, errors.New("error binding top layers")
	}

	// TODO: BIP, etc... udp stack binding here

	return b, nil
}
