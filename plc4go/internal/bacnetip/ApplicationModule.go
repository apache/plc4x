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

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

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
	bip          *BIPSimple
	annexj       *AnnexJCodec
	mux          *UDPMultiplexer
}

func NewBIPSimpleApplication(localDevice DeviceEntry, localAddress, deviceInfoCache *DeviceInventory, aseID *int) (*BIPSimpleApplication, error) {
	b := &BIPSimpleApplication{}
	var err error
	b.ApplicationIOController, err = NewApplicationIOController(localDevice, localAddress, deviceInfoCache, aseID)
	if err != nil {
		return nil, errors.Wrap(err, "error creating io controller")
	}

	b.localAddress = localAddress

	// include a application decoder
	b.asap, err = NewApplicationServiceAccessPoint(nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating application service access point")
	}

	// pass the device object to the state machine access point, so it can know if it should support segmentation
	b.smap, err = NewStateMachineAccessPoint(localDevice, deviceInfoCache, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating state machine access point")
	}

	// pass the device object to the state machine access point so it # can know if it should support segmentation
	// Note: deviceInfoCache already passed above, so we don't need to do it again here

	// a network service access point will be needed
	b.nsap, err = NewNetworkServiceAccessPoint()
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}

	// give the NSAP a generic network layer service element
	b.nse, err = NewNetworkServiceElement()
	if err != nil {
		return nil, errors.Wrap(err, "error creating new network service element")
	}
	if err := bind(b.nse, b.nsap); err != nil {
		return nil, errors.New("error binding network stack")
	}

	// bind the top layers
	if err := bind(b, b.asap, b.smap, b.nsap); err != nil {
		return nil, errors.New("error binding top layers")
	}

	// create a generic BIP stack, bound to the Annex J server on the UDP multiplexer
	b.bip, err = NewBIPSimple(nil, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new bip")
	}
	b.annexj, err = NewAnnexJCodec(nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new annex j codec")
	}
	b.mux, err = NewUDPMultiplexer(b.localAddress, false)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new udp multiplexer")
	}

	// bind the bottom layers
	if err := bind(b.bip, b.annexj, b.mux.annexJ); err != nil {
		return nil, errors.New("error binding bottom layers")
	}

	// bind the BIP stack to the network, no network number
	if err := b.nsap.bind(b.bip, nil, b.localAddress); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *BIPSimpleApplication) Close() error {
	log.Debug().Msg("close socket")
	// pass to the multiplexer, then down to the sockets
	return b.mux.Close()
}
