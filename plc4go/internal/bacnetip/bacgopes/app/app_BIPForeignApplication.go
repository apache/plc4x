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

package app

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/appservice"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/bvllservice"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/local/device"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/netservice"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/object"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/service"
)

type BIPForeignApplication struct {
	*ApplicationIOController
	*WhoIsIAmServices
	*ReadWritePropertyServices
	localAddress Address
	asap         *ApplicationServiceAccessPoint
	smap         *StateMachineAccessPoint
	nsap         *NetworkServiceAccessPoint
	nse          *NetworkServiceElement
	bip          *BIPForeign
	annexj       *AnnexJCodec
	mux          *UDPMultiplexer

	log zerolog.Logger
}

func NewBIPForeignApplication(localLog zerolog.Logger, localDevice *LocalDeviceObject, localAddress Address, bbmdAddress *Address, bbmdTTL *int, deviceInfoCache *DeviceInfoCache, aseID *int) (*BIPForeignApplication, error) {
	b := &BIPForeignApplication{
		log: localLog,
	}
	var err error
	b.ApplicationIOController, err = NewApplicationIOController(localLog, localDevice, WithApplicationIOControllerDeviceInfoCache(deviceInfoCache), WithApplicationIOControllerAseID(aseID))
	if err != nil {
		return nil, errors.Wrap(err, "error creating io controller")
	}
	b.WhoIsIAmServices, err = NewWhoIsIAmServices(localLog, b)
	if err != nil {
		return nil, errors.Wrap(err, "error WhoIs/IAm services")
	}
	b.ReadWritePropertyServices, err = NewReadWritePropertyServices()
	if err != nil {
		return nil, errors.Wrap(err, "error read write property services")
	}

	b.localAddress = localAddress

	// include a application decoder
	b.asap, err = NewApplicationServiceAccessPoint(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating application service access point")
	}

	// pass the device object to the state machine access point, so it can know if it should support segmentation
	b.smap, err = NewStateMachineAccessPoint(localLog, localDevice, WithStateMachineAccessPointDeviceInfoCache(deviceInfoCache))
	if err != nil {
		return nil, errors.Wrap(err, "error creating state machine access point")
	}

	// pass the device object to the state machine access point so it # can know if it should support segmentation
	// Note: deviceInfoCache already passed above, so we don't need to do it again here

	// a network service access point will be needed
	b.nsap, err = NewNetworkServiceAccessPoint(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}

	// give the NSAP a generic network layer service element
	b.nse, err = NewNetworkServiceElement(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new network service element")
	}
	if err := Bind(localLog, b.nse, b.nsap); err != nil {
		return nil, errors.Wrap(err, "error binding network stack")
	}

	// bind the top layers
	if err := Bind(localLog, b, b.asap, b.smap, b.nsap); err != nil {
		return nil, errors.Wrap(err, "error binding top layers")
	}

	// create a generic BIP stack, bound to the Annex J server on the UDP multiplexer
	b.bip, err = NewBIPForeign(localLog, WithBIPForeignAddress(bbmdAddress), WithBIPForeignTTL(*bbmdTTL))
	if err != nil {
		return nil, errors.Wrap(err, "error creating new bip")
	}
	b.annexj, err = NewAnnexJCodec(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new annex j codec")
	}
	b.mux, err = NewUDPMultiplexer(localLog, b.localAddress, true)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new udp multiplexer")
	}

	// bind the bottom layers
	if err := Bind(localLog, b.bip, b.annexj, b.mux.AnnexJ); err != nil {
		return nil, errors.Wrap(err, "error binding bottom layers")
	}

	// bind the BIP stack to the network, no network number
	if err := b.nsap.Bind(b.bip, nil, &b.localAddress); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *BIPForeignApplication) Close() error {
	b.log.Debug().Msg("close socket")

	if b.ApplicationIOController != nil {
		if err := b.ApplicationIOController.Close(); err != nil {
			b.log.Warn().Err(err).Msg("error closing applicationIOController")
		}
	}

	// pass to the multiplexer, then down to the sockets
	return b.mux.Close()
}
