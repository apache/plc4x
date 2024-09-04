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

package bacgopes

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type NetworkServiceAccessPoint struct {
	ServiceAccessPointContract
	Server
	adapters        map[netKey]*NetworkAdapter
	routerInfoCache *RouterInfoCache
	pendingNets     map[netKey][]NPDU
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
	n.ServiceAccessPointContract, err = NewServiceAccessPoint(localLog, func(point *serviceAccessPoint) {
		point.serviceID = n.argSapID
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}
	n.Server, err = NewServer(localLog, n, func(server *server) {
		server.serverID = n.argSid
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating server")
	}

	// map of directly connected networks
	n.adapters = make(map[netKey]*NetworkAdapter)

	// use the provided cache or make a default one
	if n.routerInfoCache == nil {
		n.routerInfoCache = NewRouterInfoCache(localLog)
	}

	// map to a list of application layer packets waiting for a path
	n.pendingNets = make(map[netKey][]NPDU)

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
func (n *NetworkServiceAccessPoint) Bind(server Server, net *uint16, address *Address) error {
	n.log.Debug().
		Interface("server", server).
		Interface("net", net).
		Stringer("address", address).
		Msg("bind")

	// make sure this hasn't already been called with this network
	if _, ok := n.adapters[nk(net)]; ok {
		return errors.Errorf("Allready bound: %v", nk(net))
	}
	// create an adapter object, add it to our map
	adapter, err := NewNetworkAdapter(n.log, n, net, address)
	if err != nil {
		return errors.Wrap(err, "error creating adapter")
	}
	n.adapters[nk(net)] = adapter
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
func (n *NetworkServiceAccessPoint) UpdateRouterReference(snet *uint16, address *Address, dnets []uint16) error {
	n.log.Debug().
		Interface("snet", snet).
		Interface("address", address).
		Interface("dnets", dnets).
		Msg("UpdateRouterReference")

	// see if we have an adapter for the snet
	_, ok := n.adapters[nk(snet)]
	if !ok {
		return errors.Errorf("no adapter for network: %d", snet)
	}

	// pass this along to the cache
	return n.routerInfoCache.UpdateRouterInfo(nk(snet), address, dnets, nil)
}

// DeleteRouterReference Delete references to routers/networks.
func (n *NetworkServiceAccessPoint) DeleteRouterReference(snet *uint16, address *Address, dnets []uint16) error {
	n.log.Debug().
		Interface("snet", snet).
		Interface("address", address).
		Interface("dnets", dnets).
		Msg("NetworkServiceAccessPoint")

	// see if we have an adapter for the snet
	_, ok := n.adapters[nk(snet)]
	if !ok {
		return errors.Errorf("no adapter for network: %d", snet)
	}

	//pass this along to the cache
	return n.routerInfoCache.DeleteRouterInfo(snet, address, dnets)
}

func (n *NetworkServiceAccessPoint) Indication(args Args, kwargs KWArgs) error {
	n.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")

	pdu := args.Get0PDU()

	// make sure our configuration is OK
	if len(n.adapters) == 0 {
		return errors.New("no adapters")
	}

	// get the local adapter
	localAdapter := n.localAdapter
	n.log.Debug().Stringer("localAdapter", localAdapter).Msg("localAdapter")

	// build a generic APDU
	apdu, err := NewAPDU(nil, WithAPDUUserData(pdu.GetPDUUserData())) // Note: upstream makes a _APDU instance which looks like a programming error as this class is only useful in an extension context...
	if err != nil {
		return errors.Wrap(err, "error creating _APDU")
	}
	if err := pdu.(Encoder).Encode(apdu); err != nil {
		return errors.Wrap(err, "error encoding APDU")
	}
	n.log.Debug().Stringer("_APDU", apdu).Msg("apdu")

	// build an NPDU specific to where it is going
	npdu, err := NewNPDU(nil, nil)
	if err != nil {
		return errors.Wrap(err, "error creating NPDU")
	}
	if err = apdu.Encode(npdu); err != nil {
		return errors.Wrap(err, "error encoding NPDU")
	}
	n.log.Debug().Stringer("npdu", npdu).Msg("npdu")

	// the hop count always starts out big
	hopCount := uint8(0xff)
	npdu.setNpduHopCount(&hopCount)

	// if this is route aware, use it for the destination
	if Settings.RouteAware && npdu.GetPDUDestination().AddrRoute != nil {
		// always a local station for now, in theory this could also be
		// a local broadcast address, remote station, or remote broadcast
		// but that is not supported by the patterns
		if npdu.GetPDUDestination().AddrRoute.AddrType != LOCAL_STATION_ADDRESS {
			panic("Route must be of type local station")
		}

		switch pdu.GetPDUDestination().AddrType {
		case REMOTE_STATION_ADDRESS, REMOTE_BROADCAST_ADDRESS, GLOBAL_BROADCAST_ADDRESS:
			npdu.setNpduDADR(apdu.GetPDUDestination())
		}

		npdu.SetPDUDestination(npdu.GetPDUDestination().AddrRoute)
		return localAdapter.ProcessNPDU(npdu)
	}

	// local stations given to local adapter
	if npdu.GetPDUDestination().AddrType == LOCAL_STATION_ADDRESS {
		return localAdapter.ProcessNPDU(npdu)
	}

	// local broadcast given to local adapter
	if npdu.GetPDUDestination().AddrType == LOCAL_BROADCAST_ADDRESS {
		return localAdapter.ProcessNPDU(npdu)
	}

	// global broadcast
	if npdu.GetPDUDestination().AddrType == GLOBAL_BROADCAST_ADDRESS {
		npdu.SetPDUDestination(NewLocalBroadcast(nil))
		npdu.setNpduDADR(apdu.GetPDUDestination())

		// send it to all of connected adapters
		for _, xadapter := range n.adapters {
			if err := xadapter.ProcessNPDU(npdu); err != nil {
				return errors.Wrap(err, "error processing NPDU")
			}
		}
		return nil
	}

	// remote broadcast
	switch npdu.GetPDUDestination().AddrType {
	case REMOTE_BROADCAST_ADDRESS, REMOTE_STATION_ADDRESS:
	default:
		return errors.Errorf("invalid destination address type: %s", npdu.GetPDUDestination().AddrType)
	}

	dnet := npdu.GetPDUDestination().AddrNet
	n.log.Debug().Interface("dnet", dnet).Msg("using dnet")

	// if the network matches the local adapter it's local
	if dnet == localAdapter.adapterNet {
		switch npdu.GetPDUDestination().AddrType {
		case REMOTE_STATION_ADDRESS:
			n.log.Debug().Msg("mapping remote station to local station")
			localStation, err := NewLocalStation(n.log, npdu.GetPDUDestination().AddrAddress, nil)
			if err != nil {
				return errors.Wrap(err, "error building local station")
			}
			pdu.SetPDUDestination(localStation)
		case REMOTE_BROADCAST_ADDRESS:
			pdu.SetPDUDestination(NewLocalBroadcast(nil))
		default:
			return errors.New("Addressing problem")
		}
		return localAdapter.ProcessNPDU(npdu)
	}

	// get it ready to send when the path is found
	npdu.SetPDUDestination(nil)
	npdu.setNpduDADR(apdu.GetPDUDestination())

	// we might already be waiting for a path for this network
	if pendingNet, ok := n.pendingNets[nk(dnet)]; ok {
		n.log.Debug().Msg("already waiting for a path")
		n.pendingNets[nk(dnet)] = append(pendingNet, npdu)
		return nil
	}

	// look for routing information from the network of one of our adapters to the destination network
	var routerInfo *RouterInfo
	var snetAdapter *NetworkAdapter
	for snet, adapter := range n.adapters {
		routerInfo = n.routerInfoCache.GetRouterInfo(snet, nk(dnet))
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
		npdu.SetPDUDestination(routerInfo.address)

		// send it along
		return snetAdapter.ProcessNPDU(npdu)
	} else {
		n.log.Debug().Msg("no known path to network")

		// add it to the list of packets waiting for the network
		netList, _ := n.pendingNets[nk(dnet)]
		netList = append(netList, npdu)
		n.pendingNets[nk(dnet)] = netList

		// build a request for the network and send it to all the adapters
		xnpdu, err := NewWhoIsRouterToNetwork(WithWhoIsRouterToNetworkNet(*dnet))
		if err != nil {
			return errors.Wrap(err, "error building WhoIsRouterToNetwork")
		}
		xnpdu.SetPDUDestination(NewLocalBroadcast(nil))

		// send it to all the adapters
		for _, adapter := range n.adapters {
			if err := n.SapIndication(NewArgs(adapter, xnpdu), NoKWArgs); err != nil {
				return errors.Wrap(err, "error doing SapIndication")
			}
		}
		return nil
	}
	return nil
}

func (n *NetworkServiceAccessPoint) ProcessNPDU(adapter *NetworkAdapter, npdu NPDU) error {
	n.log.Debug().
		Stringer("adapter", adapter).
		Stringer("npdu", npdu).
		Msg("ProcessNPDU")

	// make sure our configuration is OK
	if len(n.adapters) == 0 {
		return errors.New("no adapters")
	}

	var (
		processLocally bool
		forwardMessage bool
	)

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
		if _, ok := n.adapters[nk(snet)]; !ok {
			n.log.Warn().Msg("path error (1)")
			return nil
		}

		// pass this new path along to the cache
		n.routerInfoCache.UpdateRouterStatus(adapter.adapterNet, npdu.GetPDUSource(), []*uint16{snet})
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
				remoteStationAddress, err := NewAddress(n.log, adapter.adapterNet, npdu.GetPDUSource().AddrAddress)
				if err != nil {
					return errors.Wrap(err, "error creating remote address")
				}
				apdu.pduSource = remoteStationAddress
			} else {
				apdu.pduSource = sourceAddress
			}
			if Settings.RouteAware {
				apdu.pduSource.AddrRoute = npdu.GetPDUSource()
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
					apdu.pduSource = npdu.GetPDUSource()
				}
			} else {
				apdu.pduSource = npdu.GetPDUSource()
			}

			// pass along global broadcast
			if npdu.GetControl().GetDestinationSpecified() && destinationAddress.AddrType == GLOBAL_BROADCAST_ADDRESS {
				apdu.pduDestination = NewGlobalBroadcast(nil)
			} else {
				apdu.pduDestination = npdu.GetPDUDestination()
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
			np, ok := NPDUTypes[npdu.GetNlm().GetMessageType()]
			if !ok {
				n.log.Debug().Uint8("messageType", npdu.GetNlm().GetMessageType()).Msg("unknown npdu type")
				return nil
			}
			n.log.Debug().Msg("processing NPDU locally")

			// do a deeper decode of the NPDU
			xpdu := np()
			if err := xpdu.Decode(npdu); err != nil {
				return errors.Wrap(err, "error decoding NPDU")
			}

			// pass to the service element
			if err := n.SapRequest(NewArgs(adapter, xpdu), NoKWArgs); err != nil {
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
	newpdu := DeepCopy[NPDU](npdu)

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
		newpdu.SetPDUUserData(
			model.NewNPDU(
				npdu.GetProtocolVersionNumber(),
				model.NewNPDUControl(
					false,
					true,
					true,
					false,
					model.NPDUNetworkPriority_NORMAL_MESSAGE,
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
			),
		)

		for _, xadapter := range n.adapters {
			if xadapter != adapter {
				if err := xadapter.ProcessNPDU(DeepCopy[NPDU](newpdu)); err != nil {
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
		if xadapter, ok := n.adapters[nk(dnet)]; ok {
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
			newpdu.SetPDUUserData(model.NewNPDU(
				npdu.GetProtocolVersionNumber(),
				model.NewNPDUControl(
					false,
					false,
					true,
					false,
					model.NPDUNetworkPriority_NORMAL_MESSAGE,
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
			),
			)

			return xadapter.ProcessNPDU(DeepCopy[NPDU](newpdu))
		}

		// look for routing information from the network of one of our adapters to the destination network
		var routerInfo *RouterInfo
		var snetAdapter *NetworkAdapter
		for snet, _snetAdapter := range n.adapters {
			if _routerInfo := n.routerInfoCache.GetRouterInfo(snet, nk(dnet)); _routerInfo != nil {
				routerInfo = _routerInfo
				snetAdapter = _snetAdapter
				break
			}
		}

		// found a path
		if routerInfo != nil {
			n.log.Debug().Stringer("routerInfo", routerInfo).Msg("found path via routerInfo")

			// the destination is the address of the router
			newpdu.SetPDUDestination(routerInfo.address)

			//  send the packet downstream
			if snetAdapter == nil {
				return errors.New("snetAdapter nil")
			}
			return snetAdapter.ProcessNPDU(DeepCopy[NPDU](newpdu))
		}

		n.log.Debug().Msg("No router info found")

		// try to find a path to the network
		xnpdu := model.NewNLMWhoIsRouterToNetwork(dnet, 0)
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
	npdu := args.Get1NPDU()

	// encode it as a generic NPDU
	xpdu, err := NewNPDU(nil, nil) // TODO: add with user data thingy...
	if err != nil {
		return errors.Wrap(err, "error building NPDU")
	}
	if err := npdu.Encode(xpdu); err != nil {
		return errors.Wrap(err, "error encoding NPDU")
	}
	// npdu._xpdu = xpdu

	// tell the adapter to process the NPDU
	return adapter.ProcessNPDU(xpdu)
}

func (n *NetworkServiceAccessPoint) SapConfirmation(args Args, kwargs KWArgs) error {
	n.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("SapConfirmation")
	adapter := args.Get0NetworkAdapter()
	npdu := args.Get1NPDU()

	// encode it as a generic NPDU
	xpdu, err := NewNPDU(nil, nil) // TODO: add with user data thingy...
	if err != nil {
		return errors.Wrap(err, "error building NPDU")
	}
	if err := npdu.Encode(xpdu); err != nil {
		return errors.Wrap(err, "error encoding NPDU")
	}
	// npdu._xpdu = xpdu

	return adapter.ProcessNPDU(xpdu)
}
