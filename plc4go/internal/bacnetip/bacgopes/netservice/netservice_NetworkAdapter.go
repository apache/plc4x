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

package netservice

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/npdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
)

//go:generate plc4xGenerator -type=NetworkAdapter -prefix=netservice_
type NetworkAdapter struct {
	ClientContract
	adapterSAP           *NetworkServiceAccessPoint `asPtr:"true"`
	adapterNet           *uint16
	adapterAddr          *Address
	adapterNetConfigured *int

	// pass through args
	argCid *int

	log zerolog.Logger
}

func NewNetworkAdapter(localLog zerolog.Logger, sap *NetworkServiceAccessPoint, net *uint16, addr *Address, opts ...func(*NetworkAdapter)) (*NetworkAdapter, error) {
	n := &NetworkAdapter{
		adapterSAP:  sap,
		adapterNet:  net,
		adapterAddr: addr,

		log: localLog,
	}
	for _, opt := range opts {
		opt(n)
	}
	n.log.Trace().Stringer("sap", sap).Interface("net", net).Stringer("addr", addr).Interface("cid", n.argCid).Msg("NewNetworkAdapter")
	var err error
	n.ClientContract, err = NewClient(n.log, OptionalOption2(n.argCid, ToPtr[ClientRequirements](n), WithClientCID))
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	// record if this was 0=learned, 1=configured, None=unknown
	if net != nil {
		var state = 1
		n.adapterNetConfigured = &state
	}
	return n, nil
}

func WithNetworkAdapterCid(cid int) func(*NetworkAdapter) {
	return func(na *NetworkAdapter) {
		na.argCid = &cid
	}
}

// Confirmation Decode upstream PDUs and pass them up to the service access point.
func (n *NetworkAdapter) Confirmation(args Args, kwArgs KWArgs) error {
	n.log.Debug().
		Stringer("Args", args).Stringer("KWArgs", kwArgs).
		Interface("adapterNet", n.adapterNet).
		Msg("confirmation")

	pdu := GA[PDU](args, 0)

	npdu, err := NewNPDU(nil, nil)
	if err != nil {
		return errors.Wrap(err, "error creating NPDU")
	}
	npdu.SetPDUUserData(pdu.GetPDUUserData())
	if err := npdu.Decode(pdu); err != nil {
		return errors.Wrap(err, "error decoding NPDU")
	}
	return n.adapterSAP.ProcessNPDU(n, npdu)
}

// ProcessNPDU Encode NPDUs from the service access point and send them downstream.
func (n *NetworkAdapter) ProcessNPDU(npdu NPDU) error {
	n.log.Debug().
		Stringer("npdu", npdu).
		Interface("adapterNet", n.adapterNet).
		Msg("ProcessNPDU")

	pdu := NewPDU(NoArgs, NKW(KWCPCIUserData, npdu.GetPDUUserData()))
	if err := npdu.Encode(pdu); err != nil {
		return errors.Wrap(err, "error encoding NPDU")
	}
	return n.Request(NA(pdu), NoKWArgs())
}

func (n *NetworkAdapter) EstablishConnectionToNetwork(net any) error {
	panic("not implemented yet")
}

func (n *NetworkAdapter) DisconnectConnectionToNetwork(net any) error {
	panic("not implemented yet")
}
