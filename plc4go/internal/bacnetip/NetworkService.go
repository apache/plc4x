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
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"math"
)

type RouterStatus uint8

const (
	ROUTER_AVAILABLE    RouterStatus = iota // normal
	ROUTER_BUSY                             // router is busy
	ROUTER_DISCONNECTED                     // could make a connection, but hasn't
	ROUTER_UNREACHABLE                      // temporarily unreachable
)

func (r RouterStatus) String() string {
	switch r {
	case ROUTER_AVAILABLE:
		return "ROUTER_AVAILABLE"
	case ROUTER_BUSY:
		return "ROUTER_BUSY"
	case ROUTER_DISCONNECTED:
		return "ROUTER_DISCONNECTED"
	case ROUTER_UNREACHABLE:
		return "ROUTER_UNREACHABLE"
	default:
		return "Unknown"
	}
}

type RouterInfo struct {
	snet    *uint16
	address Address
	dnets   map[*uint16]RouterStatus
}

func (r RouterInfo) String() string {
	return fmt.Sprintf("%#q", r)
}

type RouterInfoCache struct {
	routers  map[*uint16]*RouterInfo // TODO: snet -> {Address: RouterInfo}
	pathInfo map[*uint16]*RouterInfo // TODO: (snet, dnet) -> RouterInfo
}

func NewRouterInfoCache() *RouterInfoCache {
	log.Debug().Msg("NewRouterInfoCache")
	return &RouterInfoCache{
		routers:  map[*uint16]*RouterInfo{},
		pathInfo: map[*uint16]*RouterInfo{},
	}
}

func (n *RouterInfoCache) GetRouterInfo(*uint16, *uint16) *RouterInfo {
	panic("not implemented yet")
	return nil
}

func (n *RouterInfoCache) UpdateRouterInfo(*uint16, interface{}, interface{}) error {
	panic("not implemented yet")
	return nil
}

func (n *RouterInfoCache) UpdateRouterStatus(*uint16, *Address, []*uint16) {
	panic("not implemented yet")
}

func (n *RouterInfoCache) DeleteRouterInfo(*uint16, interface{}, interface{}) error {
	panic("not implemented yet")
	return nil
}

func (n *RouterInfoCache) UpdateSourceNetwork() {
	panic("not implemented yet")
}

// TODO: implement me
type NetworkAdapter struct {
	*Client
	adapterSAP           *NetworkServiceAccessPoint
	adapterNet           *uint16
	adapterAddr          *Address
	adapterNetConfigured *int
}

func NewNetworkAdapter(sap *NetworkServiceAccessPoint, net *uint16, addr *Address, cid *int) (*NetworkAdapter, error) {
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
	log.Debug().Msgf("confirmation\n%s\n%d", npdu, n.adapterNet)

	return n.adapterSAP.ProcessNPDU(n, npdu)
}

// ProcessNPDU Encode NPDUs from the service access point and send them downstream.
func (n *NetworkAdapter) ProcessNPDU(npdu _PDU) error {
	log.Debug().Msgf("ProcessNPDU\n%s\n(net=%d)", npdu, n.adapterNet)
	return n.Request(npdu)
}

func (n *NetworkAdapter) EstablishConnectionToNetwork(net interface{}) error {
	panic("not implemented yet")
}

func (n *NetworkAdapter) DisconnectConnectionToNetwork(net interface{}) error {
	panic("not implemented yet")
}

type NetworkServiceAccessPoint struct {
	*ServiceAccessPoint
	*Server
	adapters        map[*uint16]*NetworkAdapter
	routerInfoCache *RouterInfoCache
	pendingNets     map[*uint16][]_PDU
	localAdapter    *NetworkAdapter
}

func NewNetworkServiceAccessPoint(routerInfoCache *RouterInfoCache, sapID *int, sid *int) (*NetworkServiceAccessPoint, error) {
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
	n.adapters = make(map[*uint16]*NetworkAdapter)

	// use the provided cache or make a default one
	if routerInfoCache == nil {
		routerInfoCache = NewRouterInfoCache()
	}
	n.routerInfoCache = routerInfoCache

	// map to a list of application layer packets waiting for a path
	n.pendingNets = make(map[*uint16][]_PDU)

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
func (n *NetworkServiceAccessPoint) bind(server _Server, net *uint16, address *Address) error {
	log.Debug().Msgf("bind %v net=%v address=%v", server, net, address)

	// make sure this hasn't already been called with this network
	if _, ok := n.adapters[net]; ok {
		return errors.Errorf("Allready bound: %v", net)
	}
	// create an adapter object, add it to our map
	adapter, err := NewNetworkAdapter(n, net, address, nil)
	if err != nil {
		return errors.Wrap(err, "error creating adapter")
	}
	n.adapters[net] = adapter
	log.Debug().Msgf("adapter: %v, %v", net, adapter)

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

// UpdateRouterReference Update references to routers.
func (n *NetworkServiceAccessPoint) UpdateRouterReference(snet *uint16, address, dnets interface{}) error {
	log.Debug().Msgf("UpdateRouterReference %d %s %d", snet, address, dnets)

	// see if we have an adapter for the snet
	_, ok := n.adapters[snet]
	if !ok {
		return errors.Errorf("no adapter for network: %d", snet)
	}

	// pass this along to the cache
	return n.routerInfoCache.UpdateRouterInfo(snet, address, dnets)
}

// DeleteRouterReference Delete references to routers/networks.
func (n *NetworkServiceAccessPoint) DeleteRouterReference(snet *uint16, address, dnets interface{}) error {
	log.Debug().Msgf("NetworkServiceAccessPoint %d %s %s", snet, address, dnets)

	// see if we have an adapter for the snet
	_, ok := n.adapters[snet]
	if !ok {
		return errors.Errorf("no adapter for network: %d", snet)
	}

	//pass this along to the cache
	return n.routerInfoCache.DeleteRouterInfo(snet, address, dnets)
}

func (n *NetworkServiceAccessPoint) Indication(pdu _PDU) error {
	log.Debug().Msgf("Indication:\n%s", pdu)

	// make sure our configuration is OK
	if len(n.adapters) == 0 {
		return errors.New("no adapters")
	}

	// get the local adapter
	localAdapter := n.localAdapter
	log.Debug().Msgf("localAdapter: %s", localAdapter)

	// get the apdu
	apdu := pdu.GetMessage().(readWriteModel.APDU)

	// build a npdu
	pduDestination := pdu.GetPDUDestination()

	// the hop count always starts out big
	hopCount := uint8(0xff)

	// if this is route aware, use it for the destination
	if settings.RouteAware && pduDestination.AddrRoute != nil {
		// always a local station for now, in theory this could also be
		// a local broadcast address, remote station, or remote broadcast
		// but that is not supported by the patterns
		if pduDestination.AddrRoute.AddrType != LOCAL_STATION_ADDRESS {
			panic("Route must be of type local station")
		}

		var dadr *Address
		switch pduDestination.AddrType {
		case REMOTE_STATION_ADDRESS, REMOTE_BROADCAST_ADDRESS, GLOBAL_BROADCAST_ADDRESS:
			dadr = pduDestination.AddrRoute
		}

		pdu.SetPDUDestination(pduDestination.AddrRoute)
		npdu, err := buildNPDU(hopCount, nil, dadr, pdu.GetExpectingReply(), pdu.GetNetworkPriority(), apdu)
		if err != nil {
			return errors.Wrap(err, "error building NPDU")
		}
		return localAdapter.ProcessNPDU(NewPDUFromPDUWithNewMessage(pdu, npdu))
	}

	// local stations given to local adapter
	if pduDestination.AddrType == LOCAL_STATION_ADDRESS {
		npdu, err := buildNPDU(hopCount, nil, nil, pdu.GetExpectingReply(), pdu.GetNetworkPriority(), apdu)
		if err != nil {
			return errors.Wrap(err, "error building NPDU")
		}
		return localAdapter.ProcessNPDU(NewPDUFromPDUWithNewMessage(pdu, npdu))
	}

	// local broadcast given to local adapter
	if pduDestination.AddrType == LOCAL_BROADCAST_ADDRESS {
		npdu, err := buildNPDU(hopCount, nil, nil, pdu.GetExpectingReply(), pdu.GetNetworkPriority(), apdu)
		if err != nil {
			return errors.Wrap(err, "error building NPDU")
		}
		return localAdapter.ProcessNPDU(NewPDUFromPDUWithNewMessage(pdu, npdu))
	}

	// global broadcast
	if pduDestination.AddrType == GLOBAL_BROADCAST_ADDRESS {
		pdu.SetPDUDestination(NewLocalBroadcast(nil))
		npdu, err := buildNPDU(hopCount, nil, pduDestination, pdu.GetExpectingReply(), pdu.GetNetworkPriority(), apdu)
		if err != nil {
			return errors.Wrap(err, "error building NPDU")
		}

		// send it to all of connected adapters
		for _, xadapter := range n.adapters {
			if err := xadapter.ProcessNPDU(NewPDUFromPDUWithNewMessage(pdu, npdu)); err != nil {
				return errors.Wrap(err, "error processing NPDU")
			}
		}
		return nil
	}

	// remote broadcast
	switch pduDestination.AddrType {
	case REMOTE_BROADCAST_ADDRESS, REMOTE_STATION_ADDRESS:
	default:
		return errors.Errorf("invalid destination address type: %s", pduDestination.AddrType)
	}

	dnet := pduDestination.AddrNet
	log.Debug().Msgf("dnet %d", dnet)

	// if the network matches the local adapter it's local
	if dnet == localAdapter.adapterNet {
		switch pduDestination.AddrType {
		case REMOTE_STATION_ADDRESS:
			log.Debug().Msg("mapping remote station to local station")
			localStation, err := NewLocalStation(pduDestination.AddrAddress, nil)
			if err != nil {
				return errors.Wrap(err, "error building local station")
			}
			pdu.SetPDUDestination(localStation)
		case REMOTE_BROADCAST_ADDRESS:
			pdu.SetPDUDestination(NewLocalBroadcast(nil))
		default:
			return errors.New("Addressing problem")
		}
		npdu, err := buildNPDU(hopCount, nil, pduDestination, pdu.GetExpectingReply(), pdu.GetNetworkPriority(), apdu)
		if err != nil {
			return errors.Wrap(err, "error building NPDU")
		}

		return localAdapter.ProcessNPDU(NewPDUFromPDUWithNewMessage(pdu, npdu))
	}

	// get it ready to send when the path is found
	pdu.SetPDUDestination(nil)

	npdu, err := buildNPDU(hopCount, nil, pduDestination, pdu.GetExpectingReply(), pdu.GetNetworkPriority(), apdu)
	if err != nil {
		return errors.Wrap(err, "error building NPDU")
	}

	// we might already be waiting for a path for this network
	if pendingNet, ok := n.pendingNets[dnet]; ok {
		log.Debug().Msg("already waiting for a path")
		var pdu _PDU = NewPDUFromPDUWithNewMessage(pdu, npdu)
		n.pendingNets[dnet] = append(pendingNet, pdu)
		return nil
	}

	// look for routing information from the network of one of our adapters to the destination network
	var routerInfo *RouterInfo
	var snetAdapter *NetworkAdapter
	for snet, adapter := range n.adapters {
		routerInfo = n.routerInfoCache.GetRouterInfo(snet, dnet)
		if routerInfo != nil {
			snetAdapter = adapter
			break
		}
	}

	// if there is info, we have a path
	if routerInfo != nil {
		log.Debug().Msgf("routerInfo found %s", routerInfo)

		// check the path status
		dnetStatus := routerInfo.dnets[dnet]
		log.Debug().Msgf("dnetStatus %s", dnetStatus)

		// fix the destination
		pdu.SetPDUDestination(&routerInfo.address)

		// send it along
		return snetAdapter.ProcessNPDU(NewPDUFromPDUWithNewMessage(pdu, npdu))
	} else {
		log.Debug().Msg("no known path to network")

		// add it to the list of packets waiting for the network
		netList := append(n.pendingNets[dnet], NewPDUFromPDUWithNewMessage(pdu, npdu))
		n.pendingNets[dnet] = netList

		// build a request for the network and send it to all the adapters
		whoIsRouterToNetwork := readWriteModel.NewNLMWhoIsRouterToNetwork(dnet, 0)

		// send it to all the adapters
		for _, adapter := range n.adapters {
			if err := n.SapIndicationWithAdapter(adapter, NewPDU(whoIsRouterToNetwork, WithPDUDestination(NewLocalBroadcast(nil)))); err != nil {
				return errors.Wrap(err, "error doing SapIndication")
			}
		}
	}

	panic("not implemented yet")
}

func buildNPDU(hopCount uint8, source *Address, destination *Address, expectingReply bool, networkPriority readWriteModel.NPDUNetworkPriority, apdu readWriteModel.APDU) (readWriteModel.NPDU, error) {
	sourceSpecified := source != nil
	var sourceNetworkAddress *uint16
	var sourceLength *uint8
	var sourceAddress []uint8
	if sourceSpecified {
		sourceSpecified = true
		sourceNetworkAddress = source.AddrNet
		sourceLengthValue := *source.AddrLen
		if sourceLengthValue > math.MaxUint8 {
			return nil, errors.New("source address length overflows")
		}
		sourceLengthValueUint8 := uint8(sourceLengthValue)
		sourceLength = &sourceLengthValueUint8
		sourceAddress = source.AddrAddress
		if sourceLengthValueUint8 == 0 {
			// If we define the len 0 we must not send the array
			sourceAddress = nil
		}
	}
	destinationSpecified := destination != nil
	var destinationNetworkAddress *uint16
	var destinationLength *uint8
	var destinationAddress []uint8
	if destinationSpecified {
		destinationSpecified = true
		destinationNetworkAddress = destination.AddrNet
		destinationLengthValue := *destination.AddrLen
		if destinationLengthValue > math.MaxUint8 {
			return nil, errors.New("source address length overflows")
		}
		destinationLengthValueUint8 := uint8(destinationLengthValue)
		destinationLength = &destinationLengthValueUint8
		destinationAddress = destination.AddrAddress
		if destinationLengthValueUint8 == 0 {
			// If we define the len 0 we must not send the array
			destinationAddress = nil
		}
	}
	control := readWriteModel.NewNPDUControl(false, destinationSpecified, sourceSpecified, expectingReply, networkPriority)
	return readWriteModel.NewNPDU(1, control, destinationNetworkAddress, destinationLength, destinationAddress, sourceNetworkAddress, sourceLength, sourceAddress, &hopCount, nil, apdu, 0), nil
}

func (n *NetworkServiceAccessPoint) ProcessNPDU(adapter *NetworkAdapter, pdu _PDU) error {
	log.Debug().Msgf("ProcessNPDU %s, %s", adapter, pdu)

	// make sure our configuration is OK
	if len(n.adapters) == 0 {
		return errors.New("no adapters")
	}

	npdu := pdu.GetMessage().(readWriteModel.NPDU)

	var snet *uint16
	// check for source routing
	if npdu.GetControl().GetSourceSpecified() {
		log.Debug().Msg("check source path")

		// see if this is attempting to spoof a directly connected network
		snet = npdu.GetSourceNetworkAddress()
		if _, ok := n.adapters[snet]; !ok {
			log.Warn().Msg("path error (1)")
			return nil
		}

		// pass this new path along to the cache
		n.routerInfoCache.UpdateRouterStatus(adapter.adapterNet, pdu.GetPDUSource(), []*uint16{snet})
	}

	var (
		processLocally bool
		forwardMessage bool
	)
	// check for destination routing
	if npdu.GetControl().GetDestinationSpecified() {
		log.Debug().Msg("no DADR")

		processLocally = adapter == n.localAdapter || npdu.GetControl().GetMessageTypeFieldPresent()
		forwardMessage = false
	}
	// TODO: we need the type from the DADR which we don't have in our readwrite.NPDU so we might need a special _NPDU

	panic("implement me")
	_ = processLocally
	_ = forwardMessage
	return nil
}

func (n *NetworkServiceAccessPoint) SapIndication(npdu _PDU) error {
	panic("unused")
}

func (n *NetworkServiceAccessPoint) SapIndicationWithAdapter(adapter *NetworkAdapter, npdu _PDU) error {
	log.Debug().Msgf("SapIndication %s %s", adapter, npdu)

	// encode it as a generic NPDU
	// TODO: we don't need that as a npdu is a npdu

	// tell the adapter to process the NPDU
	return adapter.ProcessNPDU(npdu)
}

func (n *NetworkServiceAccessPoint) SapConfirmation(npdu _PDU) error {
	panic("unused")
}

func (n *NetworkServiceAccessPoint) SapConfirmationWithAdapter(adapter *NetworkAdapter, npdu _PDU) error {
	log.Debug().Msgf("SapConfirmationWithAdapter %s %s", adapter, npdu)

	// encode it as a generic NPDU
	// TODO: we don't need that as a npdu is a npdu

	return adapter.ProcessNPDU(npdu)
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

func (n *NetworkServiceElement) Startup() error {
	log.Debug().Msg("Startup")

	// reference the service access point
	sap := n.elementService.(*NetworkServiceAccessPoint) // TODO: hard cast but seems like adapters appears first in network service access point (so hard binding)
	log.Debug().Msgf("sap: %v", sap)

	// loop through all the adapters
	// TODO: no adapters yet
	return nil
}
