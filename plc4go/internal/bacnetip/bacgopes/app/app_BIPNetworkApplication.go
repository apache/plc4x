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

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/bvllservice"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/netservice"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
)

type BIPNetworkApplication struct {
	*NetworkServiceElement

	localAddress Address
	nsap         *NetworkServiceAccessPoint
	bip          any // BIPSimple or BIPForeign
	annexj       *AnnexJCodec
	mux          *UDPMultiplexer

	// args for constructions
	argBBMDAddress **Address `ignore:"true"`
	argBBMDTTL     *uint16   `ignore:"true"`

	log zerolog.Logger
}

func NewBIPNetworkApplication(localLog zerolog.Logger, localAddress Address, options ...Option) (*BIPNetworkApplication, error) {
	n := &BIPNetworkApplication{
		log: localLog,
	}
	ApplyAppliers(options, n)
	optionsForParent := AddLeafTypeIfAbundant(options, n)
	var err error
	n.NetworkServiceElement, err = NewNetworkServiceElement(localLog, optionsForParent...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new network service element")
	}

	n.localAddress = localAddress

	// a network service access point will be needed
	n.nsap, err = NewNetworkServiceAccessPoint(localLog, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}

	// give the NSAP a generic network layer service element
	if err := Bind(localLog, n, n.nsap); err != nil {
		return nil, errors.New("error binding network layer")
	}

	// create a generic BIP stack, bound to the Annex J server
	// on the UDP multiplexer
	if n.argBBMDAddress == nil && n.argBBMDTTL == nil {
		n.bip, err = NewBIPSimple(localLog, options...)
		if err != nil {
			return nil, errors.Wrap(err, "error creating BIPSimple")
		}
	} else {
		n.bip, err = NewBIPForeign(localLog, Combine(options, OptionalOption(n.argBBMDAddress, WithBIPForeignAddress), OptionalOption(n.argBBMDTTL, WithBIPForeignTTL))...)
		if err != nil {
			return nil, errors.Wrap(err, "error creating BIPForeign")
		}
	}
	n.annexj, err = NewAnnexJCodec(localLog, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new annex j codec")
	}
	n.mux, err = NewUDPMultiplexer(localLog, n.localAddress, true, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new udp multiplexer")
	}

	// bind the bottom layers
	if err := Bind(localLog, n.bip, n.annexj, n.mux.AnnexJ); err != nil {
		return nil, errors.Wrap(err, "error binding bottom layers")
	}

	// bind the BIP stack to the network, no network number
	if err := n.nsap.Bind(n.bip.(Server), nil, &n.localAddress); err != nil {
		return nil, err
	}

	return n, nil
}
