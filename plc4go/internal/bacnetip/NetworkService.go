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
	"bytes"
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
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

	log zerolog.Logger
}

func NewRouterInfoCache(localLog zerolog.Logger) *RouterInfoCache {
	localLog.Debug().Msg("NewRouterInfoCache")
	return &RouterInfoCache{
		routers:  map[*uint16]*RouterInfo{},
		pathInfo: map[*uint16]*RouterInfo{},

		log: localLog,
	}
}

func (n *RouterInfoCache) GetRouterInfo(*uint16, *uint16) *RouterInfo {
	panic("not implemented yet")
	return nil
}

func (n *RouterInfoCache) UpdateRouterInfo(*uint16, any, any) error {
	panic("not implemented yet")
	return nil
}

func (n *RouterInfoCache) UpdateRouterStatus(*uint16, *Address, []*uint16) {
	panic("not implemented yet")
}

func (n *RouterInfoCache) DeleteRouterInfo(*uint16, any, any) error {
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
	var err error
	n.Client, err = NewClient(localLog, n, func(client *Client) {
		client.clientID = n.argCid
	})
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
func (n *NetworkAdapter) Confirmation(args Args, kwargs KWArgs) error {
	n.log.Debug().
		Stringer("Args", args).Stringer("KWArgs", kwargs).
		Interface("adapterNet", n.adapterNet).
		Msg("confirmation")
	npdu := args.Get0PDU()

	return n.adapterSAP.ProcessNPDU(n, npdu)
}

// ProcessNPDU Encode NPDUs from the service access point and send them downstream.
func (n *NetworkAdapter) ProcessNPDU(npdu PDU) error {
	n.log.Debug().
		Stringer("npdu", npdu).
		Interface("adapterNet", n.adapterNet).
		Msg("ProcessNPDU")
	return n.Request(NewArgs(npdu), NoKWArgs)
}

func (n *NetworkAdapter) EstablishConnectionToNetwork(net any) error {
	panic("not implemented yet")
}

func (n *NetworkAdapter) DisconnectConnectionToNetwork(net any) error {
	panic("not implemented yet")
}

type NetworkServiceAccessPoint struct {
	*ServiceAccessPoint
	*Server
	adapters        map[*uint16]*NetworkAdapter
	routerInfoCache *RouterInfoCache
	pendingNets     map[*uint16][]PDU
	localAdapter    *NetworkAdapter

	// pass through args
	argSapID *int
	argSid   *int

	log zerolog.Logger
}

func NewNetworkServiceAccessPoint(localLog zerolog.Logger, opts ...func(*NetworkServiceAccessPoint)) (*NetworkServiceAccessPoint, error) {
	n := &NetworkServiceAccessPoint{
		log: localLog,
	}
	for _, opt := range opts {
		opt(n)
	}
	var err error
	n.ServiceAccessPoint, err = NewServiceAccessPoint(localLog, n, func(point *ServiceAccessPoint) {
		point.serviceID = n.argSapID
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}
	n.Server, err = NewServer(localLog, n, func(server *Server) {
		server.serverID = n.argSid
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating server")
	}

	// map of directly connected networks
	n.adapters = make(map[*uint16]*NetworkAdapter)

	// use the provided cache or make a default one
	if n.routerInfoCache == nil {
		n.routerInfoCache = NewRouterInfoCache(localLog)
	}

	// map to a list of application layer packets waiting for a path
	n.pendingNets = make(map[*uint16][]PDU)

	return n, nil
}

func WithNetworkServiceAccessPointRouterInfoCache(routerInfoCache *RouterInfoCache) func(*NetworkServiceAccessPoint) {
	return func(n *NetworkServiceAccessPoint) {
		n.routerInfoCache = routerInfoCache
	}
}

func WithNetworkServiceAccessPointRouterSapID(sapID int) func(*NetworkServiceAccessPoint) {
	return func(n *NetworkServiceAccessPoint) {
		n.argSapID = &sapID
	}
}

func WithNetworkServiceAccessPointRouterSID(sid int) func(*NetworkServiceAccessPoint) {
	return func(n *NetworkServiceAccessPoint) {
		n.argSid = &sid
	}
}

func (n *NetworkServiceAccessPoint) String() string {
	return fmt.Sprintf("NetworkServiceAccessPoint(TBD...)") // TODO: fill some info here
}

/*
Bind creates a network adapter object and bind.

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
func (n *NetworkServiceAccessPoint) Bind(server _Server, net *uint16, address *Address) error {
	n.log.Debug().
		Interface("server", server).
		Interface("net", net).
		Stringer("address", address).
		Msg("bind")

	// make sure this hasn't already been called with this network
	if _, ok := n.adapters[net]; ok {
		return errors.Errorf("Allready bound: %v", net)
	}
	// create an adapter object, add it to our map
	adapter, err := NewNetworkAdapter(n.log, n, net, address)
	if err != nil {
		return errors.Wrap(err, "error creating adapter")
	}
	n.adapters[net] = adapter
	n.log.Debug().
		Interface("net", net).
		Stringer("adapter", adapter).
		Msg("adapter")

	// if the address was given, make it the "local" one
	if address != nil {
		n.log.Debug().Msg("setting local adapter")
		n.localAdapter = adapter
	}

	// if the local adapter isn't set yet, make it the first one, and can
	// be overridden by a subsequent call if the address is specified
	if n.localAdapter == nil {
		n.log.Debug().Msg("default local adapter")
		n.localAdapter = adapter
	}

	if n.localAdapter.adapterAddr == nil {
		n.log.Debug().Msg("no local address")
	}

	return Bind(n.log, adapter, server)
}

// UpdateRouterReference Update references to routers.
func (n *NetworkServiceAccessPoint) UpdateRouterReference(snet *uint16, address, dnets any) error {
	n.log.Debug().
		Interface("snet", snet).
		Interface("address", address).
		Interface("dnets", dnets).
		Msg("UpdateRouterReference")

	// see if we have an adapter for the snet
	_, ok := n.adapters[snet]
	if !ok {
		return errors.Errorf("no adapter for network: %d", snet)
	}

	// pass this along to the cache
	return n.routerInfoCache.UpdateRouterInfo(snet, address, dnets)
}

// DeleteRouterReference Delete references to routers/networks.
func (n *NetworkServiceAccessPoint) DeleteRouterReference(snet *uint16, address, dnets any) error {
	n.log.Debug().
		Interface("snet", snet).
		Interface("address", address).
		Interface("dnets", dnets).
		Msg("NetworkServiceAccessPoint")

	// see if we have an adapter for the snet
	_, ok := n.adapters[snet]
	if !ok {
		return errors.Errorf("no adapter for network: %d", snet)
	}

	//pass this along to the cache
	return n.routerInfoCache.DeleteRouterInfo(snet, address, dnets)
}

func (n *NetworkServiceAccessPoint) Indication(args Args, kwargs KWArgs) error {
	n.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")

	// make sure our configuration is OK
	if len(n.adapters) == 0 {
		return errors.New("no adapters")
	}

	// get the local adapter
	localAdapter := n.localAdapter
	n.log.Debug().Stringer("localAdapter", localAdapter).Msg("localAdapter")

	pdu := args.Get0PDU()
	// get the apdu
	apdu := pdu.GetRootMessage().(readWriteModel.APDU)

	// build a npdu
	pduDestination := pdu.GetPDUDestination()

	// the hop count always starts out big
	hopCount := uint8(0xff)

	// if this is route aware, use it for the destination
	if Settings.RouteAware && pduDestination.AddrRoute != nil {
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
		npdu, err := buildNPDU(hopCount, nil, dadr, pdu.GetExpectingReply(), pdu.GetNetworkPriority(), nil, apdu)
		if err != nil {
			return errors.Wrap(err, "error building NPDU")
		}
		return localAdapter.ProcessNPDU(NewPDUFromPDUWithNewMessage(pdu, npdu))
	}

	// local stations given to local adapter
	if pduDestination.AddrType == LOCAL_STATION_ADDRESS {
		npdu, err := buildNPDU(hopCount, nil, nil, pdu.GetExpectingReply(), pdu.GetNetworkPriority(), nil, apdu)
		if err != nil {
			return errors.Wrap(err, "error building NPDU")
		}
		return localAdapter.ProcessNPDU(NewPDUFromPDUWithNewMessage(pdu, npdu))
	}

	// local broadcast given to local adapter
	if pduDestination.AddrType == LOCAL_BROADCAST_ADDRESS {
		npdu, err := buildNPDU(hopCount, nil, nil, pdu.GetExpectingReply(), pdu.GetNetworkPriority(), nil, apdu)
		if err != nil {
			return errors.Wrap(err, "error building NPDU")
		}
		return localAdapter.ProcessNPDU(NewPDUFromPDUWithNewMessage(pdu, npdu))
	}

	// global broadcast
	if pduDestination.AddrType == GLOBAL_BROADCAST_ADDRESS {
		pdu.SetPDUDestination(NewLocalBroadcast(nil))
		npdu, err := buildNPDU(hopCount, nil, pduDestination, pdu.GetExpectingReply(), pdu.GetNetworkPriority(), nil, apdu)
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
	n.log.Debug().Interface("dnet", dnet).Msg("using dnet")

	// if the network matches the local adapter it's local
	if dnet == localAdapter.adapterNet {
		switch pduDestination.AddrType {
		case REMOTE_STATION_ADDRESS:
			n.log.Debug().Msg("mapping remote station to local station")
			localStation, err := NewLocalStation(n.log, pduDestination.AddrAddress, nil)
			if err != nil {
				return errors.Wrap(err, "error building local station")
			}
			pdu.SetPDUDestination(localStation)
		case REMOTE_BROADCAST_ADDRESS:
			pdu.SetPDUDestination(NewLocalBroadcast(nil))
		default:
			return errors.New("Addressing problem")
		}
		npdu, err := buildNPDU(hopCount, nil, pduDestination, pdu.GetExpectingReply(), pdu.GetNetworkPriority(), nil, apdu)
		if err != nil {
			return errors.Wrap(err, "error building NPDU")
		}

		return localAdapter.ProcessNPDU(NewPDUFromPDUWithNewMessage(pdu, npdu))
	}

	// get it ready to send when the path is found
	pdu.SetPDUDestination(nil)

	npdu, err := buildNPDU(hopCount, nil, pduDestination, pdu.GetExpectingReply(), pdu.GetNetworkPriority(), nil, apdu)
	if err != nil {
		return errors.Wrap(err, "error building NPDU")
	}

	// we might already be waiting for a path for this network
	if pendingNet, ok := n.pendingNets[dnet]; ok {
		n.log.Debug().Msg("already waiting for a path")
		var pdu = NewPDUFromPDUWithNewMessage(pdu, npdu)
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
		n.log.Debug().Stringer("routerInfo", routerInfo).Msg("routerInfo found")

		// check the path status
		dnetStatus := routerInfo.dnets[dnet]
		n.log.Debug().Stringer("dnetStatus", dnetStatus).Msg("dnetStatus")

		// fix the destination
		pdu.SetPDUDestination(&routerInfo.address)

		// send it along
		return snetAdapter.ProcessNPDU(NewPDUFromPDUWithNewMessage(pdu, npdu))
	} else {
		n.log.Debug().Msg("no known path to network")

		// add it to the list of packets waiting for the network
		netList := append(n.pendingNets[dnet], NewPDUFromPDUWithNewMessage(pdu, npdu))
		n.pendingNets[dnet] = netList

		// build a request for the network and send it to all the adapters
		whoIsRouterToNetwork := readWriteModel.NewNLMWhoIsRouterToNetwork(dnet, 0)

		// send it to all the adapters
		for _, adapter := range n.adapters {
			if err := n.SapIndication(NewArgs(adapter, NewPDU(whoIsRouterToNetwork, WithPDUDestination(NewLocalBroadcast(nil)))), NoKWArgs); err != nil {
				return errors.Wrap(err, "error doing SapIndication")
			}
		}
	}

	panic("not implemented yet")
}

// TODO: should us the one in NPDU
func buildNPDU(hopCount uint8, source *Address, destination *Address, expectingReply bool, networkPriority readWriteModel.NPDUNetworkPriority, nlm readWriteModel.NLM, apdu readWriteModel.APDU) (readWriteModel.NPDU, error) {
	switch {
	case nlm != nil && apdu != nil:
		return nil, errors.New("either specify a NLM or a APDU exclusive")
	case nlm == nil && apdu == nil:
		return nil, errors.New("either specify a NLM or a APDU")
	}
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
		sourceLengthValueUint8 := sourceLengthValue
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
	var destinationHopCount *uint8
	if destinationSpecified {
		destinationSpecified = true
		destinationNetworkAddress = destination.AddrNet
		destinationLengthValue := *destination.AddrLen
		if destinationLengthValue > math.MaxUint8 {
			return nil, errors.New("source address length overflows")
		}
		destinationLengthValueUint8 := destinationLengthValue
		destinationLength = &destinationLengthValueUint8
		destinationAddress = destination.AddrAddress
		if destinationLengthValueUint8 == 0 {
			// If we define the len 0 we must not send the array
			destinationAddress = nil
		}
		destinationHopCount = &hopCount
	}
	control := readWriteModel.NewNPDUControl(nlm != nil, destinationSpecified, sourceSpecified, expectingReply, networkPriority)
	return readWriteModel.NewNPDU(1, control, destinationNetworkAddress, destinationLength, destinationAddress, sourceNetworkAddress, sourceLength, sourceAddress, destinationHopCount, nlm, apdu, 0), nil
}

func (n *NetworkServiceAccessPoint) ProcessNPDU(adapter *NetworkAdapter, pdu PDU) error {
	n.log.Debug().
		Stringer("adapter", adapter).
		Stringer("pdu", pdu).
		Msg("ProcessNPDU")

	// make sure our configuration is OK
	if len(n.adapters) == 0 {
		return errors.New("no adapters")
	}

	var (
		processLocally bool
		forwardMessage bool
	)

	npdu := pdu.GetRootMessage().(readWriteModel.NPDU)
	sourceAddress := &Address{AddrType: NULL_ADDRESS}
	if npdu.GetControl().GetSourceSpecified() {
		snet := npdu.GetSourceNetworkAddress()
		sadr := npdu.GetSourceAddress()
		var err error
		sourceAddress, err = NewAddress(n.log, snet, sadr)
		if err != nil {
			return errors.Wrapf(err, "error parsing source address %x", sadr)
		}
	}
	destinationAddress := &Address{AddrType: NULL_ADDRESS}
	if npdu.GetControl().GetDestinationSpecified() {
		dnet := npdu.GetDestinationNetworkAddress()
		dadr := npdu.GetDestinationAddress()
		var err error
		destinationAddress, err = NewAddress(n.log, dnet, dadr)
		if err != nil {
			return errors.Wrapf(err, "error parsing destination address %x", dadr)
		}
	}
	switch {
	// check for source routing
	case npdu.GetControl().GetSourceSpecified() && sourceAddress.AddrType != NULL_ADDRESS:
		n.log.Debug().Msg("check source path")

		// see if this is attempting to spoof a directly connected network
		snet := npdu.GetSourceNetworkAddress()
		if _, ok := n.adapters[snet]; !ok {
			n.log.Warn().Msg("path error (1)")
			return nil
		}

		// pass this new path along to the cache
		n.routerInfoCache.UpdateRouterStatus(adapter.adapterNet, pdu.GetPDUSource(), []*uint16{snet})
	// check for destination routing
	case !npdu.GetControl().GetDestinationSpecified() || destinationAddress.AddrType == NULL_ADDRESS:
		n.log.Debug().Msg("no DADR")

		processLocally = adapter == n.localAdapter || npdu.GetControl().GetMessageTypeFieldPresent()
		forwardMessage = false
	case destinationAddress.AddrType == REMOTE_BROADCAST_ADDRESS:
		n.log.Debug().Msg("DADR is remote broadcast")

		if *destinationAddress.AddrNet == *adapter.adapterNet {
			n.log.Warn().Msg("path error (2)")
			return nil
		}

		processLocally = *destinationAddress.AddrNet == *n.localAdapter.adapterNet
		forwardMessage = true
	case destinationAddress.AddrType == REMOTE_STATION_ADDRESS:
		n.log.Debug().Msg("DADR is remote station")

		if *destinationAddress.AddrNet == *adapter.adapterNet {
			n.log.Warn().Msg("path error (3)")
			return nil
		}

		processLocally = *destinationAddress.AddrNet == *n.localAdapter.adapterNet && bytes.Compare(destinationAddress.AddrAddress, n.localAdapter.adapterAddr.AddrAddress) == 0
		forwardMessage = !processLocally
	case destinationAddress.AddrType == GLOBAL_BROADCAST_ADDRESS:
		n.log.Debug().Msg("DADR is global broadcast")

		processLocally = true
		forwardMessage = true
	default:
		n.log.Warn().Stringer("addrType", destinationAddress.AddrType).Msg("invalid destination address type:")
		return nil
	}

	n.log.Debug().Bool("processLocally", processLocally).Msg("processLocally")
	n.log.Debug().Bool("forwardMessage", forwardMessage).Msg("forwardMessage")

	// application or network layer message
	if !npdu.GetControl().GetMessageTypeFieldPresent() {
		n.log.Debug().Msg("application layer message")

		// decode as a generic APDU
		apdu := NewPDU(npdu.GetApdu()).(*_PDU)
		n.log.Debug().Stringer("apdu", apdu).Msg("apdu")

		// see if it needs to look routed
		if len(n.adapters) > 1 && adapter != n.localAdapter {
			// combine the source address
			if !npdu.GetControl().GetSourceSpecified() {
				remoteStationAddress, err := NewAddress(n.log, adapter.adapterNet, pdu.GetPDUSource().AddrAddress)
				if err != nil {
					return errors.Wrap(err, "error creating remote address")
				}
				apdu.pduSource = remoteStationAddress
			} else {
				apdu.pduSource = sourceAddress
			}
			if Settings.RouteAware {
				apdu.pduSource.AddrRoute = pdu.GetPDUSource()
			}

			// map the destination
			if !npdu.GetControl().GetDestinationSpecified() {
				apdu.pduDestination = n.localAdapter.adapterAddr
			} else if destinationAddress.AddrType == GLOBAL_BROADCAST_ADDRESS {
				apdu.pduDestination = NewGlobalBroadcast(nil)
			} else if destinationAddress.AddrType == GLOBAL_BROADCAST_ADDRESS {
				apdu.pduDestination = NewLocalBroadcast(nil)
			} else {
				apdu.pduDestination = n.localAdapter.adapterAddr
			}
		} else {
			// combine the source address
			if npdu.GetControl().GetSourceSpecified() {
				apdu.pduSource = sourceAddress
				if Settings.RouteAware {
					n.log.Debug().Msg("adding route")
					apdu.pduSource = pdu.GetPDUSource()
				}
			} else {
				apdu.pduSource = pdu.GetPDUSource()
			}

			// pass along global broadcast
			if npdu.GetControl().GetDestinationSpecified() && destinationAddress.AddrType == GLOBAL_BROADCAST_ADDRESS {
				apdu.pduDestination = NewGlobalBroadcast(nil)
			} else {
				apdu.pduDestination = pdu.GetPDUDestination()
			}
		}

		n.log.Debug().Stringer("pduSource", apdu.pduSource).Msg("apdu.pduSource")
		n.log.Debug().Stringer("pduDestination", apdu.pduDestination).Msg("apdu.pduDestination")

		if err := n.Response(NewArgs(apdu), NoKWArgs); err != nil {
			return errors.Wrap(err, "error passing response")
		}
	} else {
		n.log.Debug().Msg("network layer message")

		if processLocally {
			n.log.Debug().Msg("processing NPDU locally")

			// pass to the service element
			if err := n.SapRequest(NewArgs(adapter, pdu), NoKWArgs); err != nil {
				return errors.Wrap(err, "error passing sap _request")
			}
		}
	}

	// might not need to forward this to other devices
	if !forwardMessage {
		n.log.Debug().Msg("no forwarding")
		return nil
	}

	// make sure we're really a router
	if len(n.adapters) == 1 {
		n.log.Debug().Msg("not a router")
		return nil
	}

	// make sure it hasn't looped
	if npdu.GetHopCount() != nil && *npdu.GetHopCount() == 0 {
		n.log.Debug().Msg("no more hops")
		return nil
	}

	// build a new NPDU to send to other adapters
	newpdu := NewPDU(pdu).(*_PDU)

	// decrease the hop count
	newNpduHopCount := *npdu.GetHopCount() - 1

	// set the source address
	var newSADR *Address
	if !npdu.GetControl().GetSourceSpecified() {
		var err error
		newSADR, err = NewRemoteStation(n.log, adapter.adapterNet, sourceAddress.AddrAddress, nil)
		if err != nil {
			return errors.Wrap(err, "error creating remote station")
		}
	} else {
		newSADR = destinationAddress
	}

	var newDADR *Address
	// If this is a broadcast it goes everywhere
	if destinationAddress.AddrType == GLOBAL_BROADCAST_ADDRESS {
		n.log.Debug().Msg("global broadcasting")

		newDADR = NewLocalBroadcast(nil)
		newSADRLength := uint8(len(newSADR.AddrAddress))
		newDADRLength := uint8(len(newDADR.AddrAddress))
		newpdu.pduUserData = readWriteModel.NewNPDU(
			npdu.GetProtocolVersionNumber(),
			readWriteModel.NewNPDUControl(
				false,
				true,
				true,
				false,
				readWriteModel.NPDUNetworkPriority_NORMAL_MESSAGE,
			),
			newDADR.AddrNet,
			&newDADRLength,
			newDADR.AddrAddress,
			newSADR.AddrNet,
			&newSADRLength,
			newSADR.AddrAddress,
			&newNpduHopCount,
			nil,
			npdu.GetApdu(),
			0,
		)

		for _, xadapter := range n.adapters {
			if xadapter != adapter {
				if err := xadapter.ProcessNPDU(NewPDU(newpdu)); err != nil {
					n.log.Warn().Err(err).Msg("Error processing npdu")
				}
			}
		}
		return nil
	}

	if destinationAddress.AddrType == REMOTE_BROADCAST_ADDRESS || destinationAddress.AddrType == REMOTE_STATION_ADDRESS {
		dnet := destinationAddress.AddrNet
		n.log.Debug().Msg("remote station/broadcast")

		// see if this a locally connected network
		if xadapter, ok := n.adapters[dnet]; ok {
			if xadapter == adapter {
				n.log.Debug().Msg("path error (4)")
				return nil
			}
			n.log.Debug().Stringer("adapter", adapter).Msg("found path via")

			// if this was a remote broadcast, it's now a local one
			if destinationAddress.AddrType == REMOTE_BROADCAST_ADDRESS {
				newDADR = NewLocalBroadcast(nil)
			} else {
				var err error
				newDADR, err = NewLocalStation(n.log, destinationAddress.AddrAddress, nil)
				if err != nil {
					return errors.Wrap(err, "error building local station")
				}
			}

			// last leg in routing
			newDADR = nil

			// send the packet downstream
			newSADRLength := uint8(len(newSADR.AddrAddress))
			newpdu.pduUserData = readWriteModel.NewNPDU(
				npdu.GetProtocolVersionNumber(),
				readWriteModel.NewNPDUControl(
					false,
					false,
					true,
					false,
					readWriteModel.NPDUNetworkPriority_NORMAL_MESSAGE,
				),
				nil,
				nil,
				nil,
				newSADR.AddrNet,
				&newSADRLength,
				newSADR.AddrAddress,
				&newNpduHopCount,
				nil,
				npdu.GetApdu(),
				0,
			)

			return xadapter.ProcessNPDU(NewPDU(newpdu))
		}

		// look for routing information from the network of one of our adapters to the destination network
		var routerInfo *RouterInfo
		var snetAdapter *NetworkAdapter
		for snet, _snetAdapter := range n.adapters {
			if _routerInfo := n.routerInfoCache.GetRouterInfo(snet, dnet); _routerInfo != nil {
				routerInfo = _routerInfo
				snetAdapter = _snetAdapter
				break
			}
		}

		// found a path
		if routerInfo != nil {
			n.log.Debug().Stringer("routerInfo", routerInfo).Msg("found path via routerInfo")

			// the destination is the address of the router
			pduDestination := routerInfo.address

			//  send the packet downstream
			if snetAdapter == nil {
				return errors.New("snetAdapter nil")
			}
			return snetAdapter.ProcessNPDU(NewPDU(newpdu, WithPDUDestination(&pduDestination)))
		}

		n.log.Debug().Msg("No router info found")

		// try to find a path to the network
		xnpdu := readWriteModel.NewNLMWhoIsRouterToNetwork(dnet, 0)
		pduDestination := NewLocalBroadcast(nil)

		// send it to all the connected adapters
		for _, xadapter := range n.adapters {
			// skip the horse it rode in on
			if xadapter == adapter {
				continue
			}

			// pass this along as if it came from the NSE
			if err := n.SapIndication(NewArgs(xadapter, NewPDU(xnpdu, WithPDUDestination(pduDestination))), NoKWArgs); err != nil {
				return errors.Wrap(err, "error sending indication")
			}
		}

		return nil
	}

	n.log.Debug().
		Interface("destinationNetworkAddress", npdu.GetDestinationNetworkAddress()).
		Interface("destinationAddress", npdu.GetDestinationAddress()).
		Msg("bad DADR")
	return nil
}

func (n *NetworkServiceAccessPoint) SapIndication(args Args, kwargs KWArgs) error {
	n.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("SapIndication")
	adapter := args.Get0NetworkAdapter()
	npdu := args.Get1PDU()

	// encode it as a generic NPDU
	// TODO: we don't need that as a npdu is a npdu

	// tell the adapter to process the NPDU
	return adapter.ProcessNPDU(npdu)
}

func (n *NetworkServiceAccessPoint) SapConfirmation(args Args, kwargs KWArgs) error {
	n.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("SapConfirmation")
	adapter := args.Get0NetworkAdapter()
	npdu := args.Get1PDU()

	// encode it as a generic NPDU
	// TODO: we don't need that as a npdu is a npdu

	return adapter.ProcessNPDU(npdu)
}

type NetworkServiceElement struct {
	*ApplicationServiceElement

	networkNumberIsTask time.Time

	// regular args
	argStartupDisabled bool

	// pass through args
	argEID *int

	log zerolog.Logger
}

func NewNetworkServiceElement(localLog zerolog.Logger, opts ...func(*NetworkServiceElement)) (*NetworkServiceElement, error) {
	n := &NetworkServiceElement{
		log: localLog,
	}
	for _, opt := range opts {
		opt(n)
	}
	var err error
	n.ApplicationServiceElement, err = NewApplicationServiceElement(localLog, n, func(element *ApplicationServiceElement) {
		element.elementID = n.argEID
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating application service element")
	}

	// network number is timeout
	n.networkNumberIsTask = time.Time{}

	// if starting up is enabled defer our startup function
	if !n.argStartupDisabled {
		Deferred(n.Startup, NoArgs, NoKWArgs)
	}
	return n, nil
}

func WithNetworkServiceElementEID(eid int) func(*NetworkServiceElement) {
	return func(n *NetworkServiceElement) {
		n.argEID = &eid
	}
}

func WithNetworkServiceElementStartupDisabled(startupDisabled bool) func(*NetworkServiceElement) {
	return func(n *NetworkServiceElement) {
		n.argStartupDisabled = startupDisabled
	}
}

func (n *NetworkServiceElement) String() string {
	return fmt.Sprintf("NetworkServiceElement(TBD...)") // TODO: fill some info here
}

func (n *NetworkServiceElement) Startup(_ Args, _ KWArgs) error {
	n.log.Debug().Msg("Startup")

	// reference the service access point
	sap := n.elementService.(*NetworkServiceAccessPoint) // TODO: hard cast but seems like adapters appears first in network service access point (so hard binding)
	n.log.Debug().Stringer("sap", sap).Msg("sap")

	// loop through all the adapters
	for _, adapter := range sap.adapters {
		n.log.Debug().Stringer("adapter", adapter).Msg("adapter")

		if adapter.adapterNet == nil {
			n.log.Trace().Msg("skipping, unknown net")
			continue
		}
		if adapter.adapterAddr == nil {
			n.log.Trace().Msg("skipping, unknown addr")
			continue
		}

		// build a list of reachable networks
		var netlist []*uint16

		// loop through the adapters
		for _, xadapter := range sap.adapters {
			if xadapter != adapter {
				if xadapter.adapterNet == nil || xadapter.adapterAddr == nil {
					continue
				}
				netlist = append(netlist, xadapter.adapterNet)
			}
		}

		// skip for an empty list, perhaps they are not yet learned
		if len(netlist) == 0 {
			n.log.Trace().Msg("skipping, no netlist")
			continue
		}

		// pass this along to the cache -- on hold #213
		// sap.router_info_cache.update_router_info(adapter.adapterNet, adapter.adapterAddr, netlist)

		// send an announcement
		if err := n.iamRouterToNetwork(adapter, nil, netlist); err != nil {
			n.log.Debug().Err(err).Msg("I-Am-Router-To-Network failed")
		}
	}
	return nil
}

func (n *NetworkServiceElement) Indication(args Args, kwargs KWArgs) error {
	n.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")

	adapter := args.Get0NetworkAdapter()
	npdu := args.Get1PDU()

	switch message := npdu.GetRootMessage().(type) {
	case readWriteModel.NPDUExactly:
		switch nlm := message.GetNlm().(type) {
		case readWriteModel.NLMWhoIsRouterToNetwork:
			n.WhoIsRouteToNetwork(adapter, nlm)
		case readWriteModel.NLMIAmRouterToNetwork:
			n.IAmRouterToNetwork(adapter, nlm)
		case readWriteModel.NLMICouldBeRouterToNetwork:
			n.ICouldBeRouterToNetwork(adapter, nlm)
		case readWriteModel.NLMRejectMessageToNetwork:
			n.RejectRouterToNetwork(adapter, nlm)
		case readWriteModel.NLMRouterBusyToNetwork:
			n.RouterBusyToNetwork(adapter, nlm)
		case readWriteModel.NLMRouterAvailableToNetwork:
			n.RouterAvailableToNetwork(adapter, nlm)
		case readWriteModel.NLMInitializeRoutingTable:
			n.InitalizeRoutingTable(adapter, nlm)
		case readWriteModel.NLMInitializeRoutingTableAck:
			n.InitalizeRoutingTableAck(adapter, nlm)
		case readWriteModel.NLMEstablishConnectionToNetwork:
			n.EstablishConnectionToNetwork(adapter, nlm)
		case readWriteModel.NLMDisconnectConnectionToNetwork:
			n.DisconnectConnectionToNetwork(adapter, nlm)
		case readWriteModel.NLMWhatIsNetworkNumber:
			n.WhatIsNetworkNumber(adapter, nlm)
		case readWriteModel.NLMNetworkNumberIs:
			n.NetworkNumberIs(adapter, nlm)
		default:
			n.log.Debug().Stringer("nlm", nlm).Msg("Unhandled")
		}
	default:
		n.log.Trace().Msg("can only handle NPDU")
	}
	return nil
}

func (n *NetworkServiceElement) Confirmation(args Args, kwargs KWArgs) error {
	n.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Confirmation")

	adapter := args.Get0NetworkAdapter()
	npdu := args.Get1PDU()

	switch message := npdu.GetRootMessage().(type) {
	case readWriteModel.NPDUExactly:
		switch nlm := message.GetNlm().(type) {
		case readWriteModel.NLMWhoIsRouterToNetwork:
			n.WhoIsRouteToNetwork(adapter, nlm)
		case readWriteModel.NLMIAmRouterToNetwork:
			n.IAmRouterToNetwork(adapter, nlm)
		case readWriteModel.NLMICouldBeRouterToNetwork:
			n.ICouldBeRouterToNetwork(adapter, nlm)
		case readWriteModel.NLMRejectMessageToNetwork:
			n.RejectRouterToNetwork(adapter, nlm)
		case readWriteModel.NLMRouterBusyToNetwork:
			n.RouterBusyToNetwork(adapter, nlm)
		case readWriteModel.NLMRouterAvailableToNetwork:
			n.RouterAvailableToNetwork(adapter, nlm)
		case readWriteModel.NLMInitializeRoutingTable:
			n.InitalizeRoutingTable(adapter, nlm)
		case readWriteModel.NLMInitializeRoutingTableAck:
			n.InitalizeRoutingTableAck(adapter, nlm)
		case readWriteModel.NLMEstablishConnectionToNetwork:
			n.EstablishConnectionToNetwork(adapter, nlm)
		case readWriteModel.NLMDisconnectConnectionToNetwork:
			n.DisconnectConnectionToNetwork(adapter, nlm)
		case readWriteModel.NLMWhatIsNetworkNumber:
			n.WhatIsNetworkNumber(adapter, nlm)
		case readWriteModel.NLMNetworkNumberIs:
			n.NetworkNumberIs(adapter, nlm)
		default:
			n.log.Debug().Stringer("nlm", nlm).Msg("Unhandled")
		}
	default:
		n.log.Trace().Msg("can only handle NPDU")
	}
	return nil
}

func (n *NetworkServiceElement) iamRouterToNetwork(adapter *NetworkAdapter, destination *Address, network []*uint16) error {
	n.log.Debug().Stringer("adapter", adapter).Stringer("destination", destination).Interface("network", network).Msg("IamRouterToNetwork")

	// reference the service access point
	sap := n.elementService.(*NetworkServiceAccessPoint) // TODO: hard cast but seems like adapters appears first in network service access point (so hard binding)
	n.log.Debug().Interface("sap", sap).Msg("SAP")

	// if we're not a router, trouble
	if len(sap.adapters) == 1 {
		return errors.New("not a router")
	}

	if adapter != nil {
		if destination == nil {
			destination = NewLocalBroadcast(nil)
		} else if destination.AddrType == LOCAL_STATION_ADDRESS || destination.AddrType == LOCAL_BROADCAST_ADDRESS {
			return nil
		} else if destination.AddrType == REMOTE_STATION_ADDRESS {
			if *destination.AddrNet != *adapter.adapterNet {
				return errors.New("invalid address, remote station for a different adapter")
			}
			var err error
			destination, err = NewLocalStation(n.log, destination.AddrAddress, nil)
			if err != nil {
				return errors.Wrap(err, "error creating station")
			}
		} else if destination.AddrType == REMOTE_BROADCAST_ADDRESS {
			if *destination.AddrNet != *adapter.adapterNet {
				return errors.New("invalid address, remote station for a different adapter")
			}
			destination = NewLocalBroadcast(nil)
		} else {
			return errors.New("invalid destination address")
		}
	} else {
		if destination == nil {
			destination = NewLocalBroadcast(nil)
		} else if destination.AddrType == LOCAL_STATION_ADDRESS {
			return errors.New("ambiguous destination")
		} else if destination.AddrType == LOCAL_BROADCAST_ADDRESS {
			return nil
		} else if destination.AddrType == REMOTE_STATION_ADDRESS {
			var ok bool
			adapter, ok = sap.adapters[destination.AddrNet]
			if !ok {
				return errors.New("invalid address, no network for remote station")
			}
			var err error
			destination, err = NewLocalStation(n.log, destination.AddrAddress, nil)
			if err != nil {
				return errors.Wrap(err, "error creating station")
			}
		} else if destination.AddrType == REMOTE_BROADCAST_ADDRESS {
			var ok bool
			adapter, ok = sap.adapters[destination.AddrNet]
			if !ok {
				return errors.New("invalid address, no network for remote broadcast")
			}
			destination = NewLocalBroadcast(nil)
		} else {
			return errors.New("invalid destination address")
		}
	}
	n.log.Debug().Stringer("adapter", adapter).Stringer("destination", destination).Interface("network", network).Msg("adapter, destination, network")

	// process a single adapter or all of the adapters
	var adapterList []*NetworkAdapter
	if adapter != nil {
		adapterList = append(adapterList, adapter)
	} else {
		for _, networkAdapter := range sap.adapters {
			adapterList = append(adapterList, networkAdapter)
		}
	}

	// loop through all of the adapters
	for _, adapter := range adapterList {
		// build a list of reachable networks
		var netlist []uint16

		for _, xadapter := range sap.adapters {
			if xadapter != adapter {
				netlist = append(netlist, *xadapter.adapterNet)
			}
		}

		if network == nil {
			return nil
		}
		for _, u := range network {
			if slices.Contains(netlist, *u) {
				netlist = append(netlist, *u)
			}
		}

		// build a response
		iamrtn := readWriteModel.NewNLMIAmRouterToNetwork(netlist, 0)

		npdu, err := buildNPDU(0, nil, destination, false, readWriteModel.NPDUNetworkPriority_NORMAL_MESSAGE, iamrtn, nil)
		if err != nil {
			return errors.Wrap(err, "error building NPDU")
		}

		n.log.Debug().Stringer("npdu", npdu).Msg("adapter, iamrtn")

		// send it back
		if err := n.Request(NewArgs(adapter, npdu), NoKWArgs); err != nil {
			return errors.Wrap(err, "error requesting NPDU")
		}
	}

	return nil
}

func (n *NetworkServiceElement) WhoIsRouteToNetwork(adapter *NetworkAdapter, nlm readWriteModel.NLMWhoIsRouterToNetwork) {
	// TODO: implement me
}

func (n *NetworkServiceElement) IAmRouterToNetwork(adapter *NetworkAdapter, nlm readWriteModel.NLMIAmRouterToNetwork) {
	// TODO: implement me
}

func (n *NetworkServiceElement) ICouldBeRouterToNetwork(adapter *NetworkAdapter, nlm readWriteModel.NLMICouldBeRouterToNetwork) {
	// TODO: implement me
}

func (n *NetworkServiceElement) RejectRouterToNetwork(adapter *NetworkAdapter, nlm readWriteModel.NLMRejectMessageToNetwork) {
	// TODO: implement me
}

func (n *NetworkServiceElement) RouterBusyToNetwork(adapter *NetworkAdapter, nlm readWriteModel.NLMRouterBusyToNetwork) {
	// TODO: implement me
}

func (n *NetworkServiceElement) RouterAvailableToNetwork(adapter *NetworkAdapter, nlm readWriteModel.NLMRouterAvailableToNetwork) {
	// TODO: implement me
}

func (n *NetworkServiceElement) InitalizeRoutingTable(adapter *NetworkAdapter, nlm readWriteModel.NLMInitializeRoutingTable) {
	// TODO: implement me
}

func (n *NetworkServiceElement) InitalizeRoutingTableAck(adapter *NetworkAdapter, nlm readWriteModel.NLMInitializeRoutingTableAck) {
	// TODO: implement me
}

func (n *NetworkServiceElement) EstablishConnectionToNetwork(adapter *NetworkAdapter, nlm readWriteModel.NLMEstablishConnectionToNetwork) {
	// TODO: implement me
}

func (n *NetworkServiceElement) DisconnectConnectionToNetwork(adapter *NetworkAdapter, nlm readWriteModel.NLMDisconnectConnectionToNetwork) {
	// TODO: implement me
}

func (n *NetworkServiceElement) WhatIsNetworkNumber(adapter *NetworkAdapter, nlm readWriteModel.NLMWhatIsNetworkNumber) {
	// TODO: implement me
}

func (n *NetworkServiceElement) NetworkNumberIs(adapter *NetworkAdapter, nlm readWriteModel.NLMNetworkNumberIs) {
	// TODO: implement me
}

// TODO: big WIP implement all other methods...
