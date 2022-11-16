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
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net"
)

// TODO: implement me
type NetworkAdapter struct {
	*Client
	adapterSAP           *NetworkServiceAccessPoint
	adapterNet           interface{}
	adapterAddr          net.Addr
	adapterNetConfigured *int
}

func NewNetworkAdapter(sap *NetworkServiceAccessPoint, net interface{}, addr net.Addr, cid *int) (*NetworkAdapter, error) {
	n := &NetworkAdapter{
		adapterSAP:  sap,
		adapterNet:  net,
		adapterAddr: addr,
	}
	var err error
	n.Client, err = NewClient(cid, n)
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

// Confirmation Decode upstream PDUs and pass them up to the service access point.
func (n *NetworkAdapter) Confirmation(npdu spi.Message) error {
	log.Debug().Msgf("confirmation\n%s\n%s", npdu, n.adapterNet)

	// TODO: we need generics otherwise this won't work at all here
	return n.adapterSAP.ProcessNPDU(npdu)
}

// ProcessNPDU Encode NPDUs from the service access point and send them downstream.
func (n *NetworkAdapter) ProcessNPDU(npdu spi.Message) error {
	log.Debug().Msgf("ProcessNPDU\n%s\n(net=%s)", npdu, n.adapterNet)
	return n.Request(npdu)
}

type NetworkServiceAccessPoint struct {
	*ServiceAccessPoint
	*Server
	adapters        map[string]*NetworkAdapter
	routerInfoCache interface{}
	pendingNets     map[string]interface{}
	localAdapter    interface{}
}

func NewNetworkServiceAccessPoint(routerInfoCache interface{}, sapID *int, sid *int) (*NetworkServiceAccessPoint, error) {
	n := &NetworkServiceAccessPoint{}
	var err error
	n.ServiceAccessPoint, err = NewServiceAccessPoint(sapID, n)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}
	n.Server, err = NewServer(sid, n)
	if err != nil {
		return nil, errors.Wrap(err, "error creating server")
	}

	// map of directly connected networks
	n.adapters = make(map[string]*NetworkAdapter)

	// use the provided cache or make a default one
	if routerInfoCache == nil {
		// TODO: create a new cache
	}
	n.routerInfoCache = routerInfoCache

	// map to a list of application layer packets waiting for a path
	n.pendingNets = make(map[string]interface{})

	return n, nil
}

func (n *NetworkServiceAccessPoint) bind(server _Server, net interface{}, address interface{}) error {
	panic("not implemented yet")
}

func (n *NetworkServiceAccessPoint) UpdateRouterReference() error {
	panic("not implemented yet")
}

func (n *NetworkServiceAccessPoint) DeleteRouterReference() error {
	panic("not implemented yet")
}

func (n *NetworkServiceAccessPoint) Indication(npdu spi.Message) error {
	panic("not implemented yet")
}

func (n *NetworkServiceAccessPoint) ProcessNPDU(npdu spi.Message) error {
	panic("not implemented yet")
}

func (n *NetworkServiceAccessPoint) SapIndication(npdu spi.Message) error {
	// TODO: extract from somewhere
	var pduDestination []byte
	panic("we need pduDestination")
	_ = pduDestination
	panic("not implemented yet")
}

func (n *NetworkServiceAccessPoint) SapConfirmation(npdu spi.Message) error {
	// TODO: extract from somewhere
	var pduDestination []byte
	panic("we need pduDestination")
	_ = pduDestination
	panic("not implemented yet")
}

type NetworkServiceElement struct {
	*ApplicationServiceElement

	// TODO: implement me
}

func NewNetworkServiceElement(eid *int) (*NetworkServiceElement, error) {
	n := &NetworkServiceElement{}
	var err error
	n.ApplicationServiceElement, err = NewApplicationServiceElement(eid, n)
	if err != nil {
		return nil, errors.Wrap(err, "error creating application service element")
	}

	// TODO: we need to use defer from package as this go routine is too early
	go n.Startup()
	return n, nil
}

func (n *NetworkServiceElement) Startup() {
	log.Debug().Msg("Startup")

	// reference the service access point
	sap := n.elementService.(*NetworkServiceAccessPoint) // TODO: hard cast but seems like adapters apears first in network service access point (so hard binding)
	log.Debug().Msgf("sap: %v", sap)

	// loop through all the adapters
	// TODO: no adapters yet
}
