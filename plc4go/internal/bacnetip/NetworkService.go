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
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// TODO: implement me
type NetworkAdapter struct {
	*Client
	adapterSAP           *NetworkServiceAccessPoint
	adapterNet           interface{}
	adapterAddr          *Address
	adapterNetConfigured *int
}

func NewNetworkAdapter(sap *NetworkServiceAccessPoint, net interface{}, addr *Address, cid *int) (*NetworkAdapter, error) {
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
func (n *NetworkAdapter) Confirmation(npdu _PDU) error {
	log.Debug().Msgf("confirmation\n%s\n%s", npdu, n.adapterNet)

	// TODO: we need generics otherwise this won't work at all here
	return n.adapterSAP.ProcessNPDU(npdu)
}

// ProcessNPDU Encode NPDUs from the service access point and send them downstream.
func (n *NetworkAdapter) ProcessNPDU(npdu _PDU) error {
	log.Debug().Msgf("ProcessNPDU\n%s\n(net=%s)", npdu, n.adapterNet)
	return n.Request(npdu)
}

type NetworkServiceAccessPoint struct {
	*ServiceAccessPoint
	*Server
	adapters        map[string]*NetworkAdapter
	routerInfoCache interface{}
	pendingNets     map[string]interface{}
	localAdapter    *NetworkAdapter
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

/* bind creates a network adapter object and bind.

   bind(s, None, None)
       Called for simple applications, local network unknown, no specific
       address, APDUs sent upstream

   bind(s, net, None)
       Called for routers, bind to the network, (optionally?) drop APDUs

   bind(s, None, address)
       Called for applications or routers, bind to the network (to be
       discovered), send up APDUs with a metching address

   bind(s, net, address)
       Called for applications or routers, bind to the network, send up
       APDUs with a metching address.
*/
func (n *NetworkServiceAccessPoint) bind(server _Server, net interface{}, address *Address) error {
	log.Debug().Msgf("bind %v net=%v address=%v", server, net, address)

	netKey := fmt.Sprintf("%v", net)
	// make sure this hasn't already been called with this network
	if _, ok := n.adapters[netKey]; ok {
		return errors.Errorf("Allready bound: %v", net)
	}
	// create an adapter object, add it to our map
	adapter, err := NewNetworkAdapter(n, net, address, nil)
	if err != nil {
		return errors.Wrap(err, "error creating adapter")
	}
	n.adapters[netKey] = adapter
	log.Debug().Msgf("adapter: %v, %v", netKey, adapter)

	// if the address was given, make it the "local" one
	if address != nil {
		log.Debug().Msg("setting local adapter")
		n.localAdapter = adapter
	}

	// if the local adapter isn't set yet, make it the first one, and can
	// be overridden by a subsequent call if the address is specified
	if n.localAdapter == nil {
		log.Debug().Msg("default local adapter")
		n.localAdapter = adapter
	}

	if n.localAdapter.adapterAddr == nil {
		log.Debug().Msg("no local address")
	}

	return bind(adapter, server)
}

func (n *NetworkServiceAccessPoint) UpdateRouterReference() error {
	panic("not implemented yet")
}

func (n *NetworkServiceAccessPoint) DeleteRouterReference() error {
	panic("not implemented yet")
}

func (n *NetworkServiceAccessPoint) Indication(npdu _PDU) error {
	panic("not implemented yet")
}

func (n *NetworkServiceAccessPoint) ProcessNPDU(npdu _PDU) error {
	panic("not implemented yet")
}

func (n *NetworkServiceAccessPoint) SapIndication(npdu _PDU) error {
	panic("not implemented yet")
}

func (n *NetworkServiceAccessPoint) SapConfirmation(npdu _PDU) error {
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

	Deferred(n.Startup)
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
