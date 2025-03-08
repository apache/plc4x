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
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/local/device"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/netservice"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/service"
)

//go:generate go tool plc4xGenerator -type=BIPSimpleApplication -prefix=app_
type BIPSimpleApplication struct {
	*ApplicationIOController
	*WhoIsIAmServices
	*ReadWritePropertyServices

	localAddress Address                        `stringer:"true"`
	asap         *ApplicationServiceAccessPoint `stringer:"true"`
	smap         *StateMachineAccessPoint       `stringer:"true"`
	nsap         *NetworkServiceAccessPoint     `stringer:"true"`
	nse          *NetworkServiceElement         `stringer:"true"`
	bip          *BIPSimple                     `stringer:"true"`
	annexj       *AnnexJCodec                   `stringer:"true"`
	mux          *UDPMultiplexer                `stringer:"true"`

	log zerolog.Logger
}

func NewBIPSimpleApplication(localLog zerolog.Logger, localDevice LocalDeviceObject, localAddress Address, options ...Option) (*BIPSimpleApplication, error) {
	b := &BIPSimpleApplication{
		log: localLog,
	}
	ApplyAppliers(options, b)
	optionsForParent := AddLeafTypeIfAbundant(options, b)
	var err error
	b.ApplicationIOController, err = NewApplicationIOController(localLog, Combine(optionsForParent, WithApplicationLocalDeviceObject(localDevice))...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating io controller")
	}
	b.WhoIsIAmServices, err = NewWhoIsIAmServices(localLog, b, optionsForParent...)
	if err != nil {
		return nil, errors.Wrap(err, "error WhoIs/IAm services")
	}
	b.ReadWritePropertyServices, err = NewReadWritePropertyServices(optionsForParent...)
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
	b.smap, err = NewStateMachineAccessPoint(localLog, localDevice, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating state machine access point")
	}

	// pass the device object to the state machine access point so it # can know if it should support segmentation
	// Note: deviceInfoCache already passed above, so we don't need to do it again here

	// a network service access point will be needed
	b.nsap, err = NewNetworkServiceAccessPoint(localLog, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}

	// give the NSAP a generic network layer service element
	b.nse, err = NewNetworkServiceElement(localLog, options...)
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
	b.bip, err = NewBIPSimple(localLog, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new bip")
	}
	b.annexj, err = NewAnnexJCodec(localLog, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new annex j codec")
	}
	b.mux, err = NewUDPMultiplexer(localLog, b.localAddress, false, options...)
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

func (b *BIPSimpleApplication) Close() error {
	b.log.Debug().Msg("close socket")

	if b.ApplicationIOController != nil {
		if err := b.ApplicationIOController.Close(); err != nil {
			b.log.Warn().Err(err).Msg("error closing applicationIOController")
		}
	}

	// pass to the multiplexer, then down to the sockets
	return b.mux.Close()
}
