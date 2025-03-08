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
	"bytes"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/apdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/npdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

//go:generate go tool plc4xGenerator -type=NetworkServiceAccessPoint -prefix=netservice_
type NetworkServiceAccessPoint struct {
	ServiceAccessPointContract
	ServerContract
	*DebugContents `ignore:"true"`

	adapters        map[netKey]*NetworkAdapter
	routerInfoCache *RouterInfoCache `stringer:"true"`
	pendingNets     map[netKey][]NPDU
	localAdapter    *NetworkAdapter `stringer:"true"`

	log zerolog.Logger
}

func NewNetworkServiceAccessPoint(localLog zerolog.Logger, options ...Option) (*NetworkServiceAccessPoint, error) {
	n := &NetworkServiceAccessPoint{
		log: localLog,
	}
	n.DebugContents = NewDebugContents(n, "adapters++", "pending_nets", "local_adapter-")
	ApplyAppliers(options, n)
	optionsForParent := AddLeafTypeIfAbundant(options, n)
	var err error
	n.ServiceAccessPointContract, err = NewServiceAccessPoint(localLog, optionsForParent...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}
	n.ServerContract, err = NewServer(localLog, optionsForParent...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating server")
	}
	if _debug != nil {
		_debug("__init__ sap=%r sid=%r", n.GetServiceElement(), n.GetServiceID())
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

func WithNetworkServiceAccessPointRouterInfoCache(routerInfoCache *RouterInfoCache) GenericApplier[*NetworkServiceAccessPoint] {
	return WrapGenericApplier(func(n *NetworkServiceAccessPoint) { n.routerInfoCache = routerInfoCache })
}

func (n *NetworkServiceAccessPoint) GetDebugAttr(attr string) any {
	switch attr {
	case "adapters":
		return n.adapters
	case "pending_nets":
		return n.pendingNets
	case "local_adapter":
		return n.localAdapter
	default:
		return nil
	}
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
		Stringer("server", server).
		Interface("net", net).
		Stringer("address", address).
		Msg("bind")
	if _debug != nil {
		_debug("bind %r net=%r address=%r", server, net, address)
	}

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
	if _debug != nil {
		_debug("    - adapter: %r, %r", net, adapter)
	}

	// if the address was given, make it the "local" one
	if address != nil {
		n.log.Debug().Msg("setting local adapter")
		if _debug != nil {
			_debug("    - setting local adapter")
		}
		n.localAdapter = adapter
	}

	// if the local adapter isn't set yet, make it the first one, and can
	// be overridden by a subsequent call if the address is specified
	if n.localAdapter == nil {
		n.log.Debug().Msg("default local adapter")
		if _debug != nil {
			_debug("    - default local adapter")
		}
		n.localAdapter = adapter
	}

	if n.localAdapter.adapterAddr == nil {
		n.log.Debug().Msg("no local address")
		if _debug != nil {
			_debug("    - no local address")
		}
	}

	return Bind(n.log, adapter, server)
}

// UpdateRouterReference Update references to routers.
func (n *NetworkServiceAccessPoint) UpdateRouterReference(snet *uint16, address *Address, dnets []uint16) error {
	n.log.Debug().
		Interface("snet", snet).
		Stringer("address", address).
		Interface("dnets", dnets).
		Msg("UpdateRouterReference")
	if _debug != nil {
		_debug("update_router_references %r %r %r", snet, address, dnets)
	}

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
		Stringer("address", address).
		Interface("dnets", dnets).
		Msg("NetworkServiceAccessPoint")
	if _debug != nil {
		_debug("delete_router_references %r %r %r", snet, address, dnets)
	}

	// see if we have an adapter for the snet
	_, ok := n.adapters[nk(snet)]
	if !ok {
		return errors.Errorf("no adapter for network: %d", snet)
	}

	//pass this along to the cache
	return n.routerInfoCache.DeleteRouterInfo(nk(snet), address, dnets)
}

func (n *NetworkServiceAccessPoint) Indication(args Args, kwArgs KWArgs) error {
	n.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Indication")

	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("indication %r", pdu)
	}

	// make sure our configuration is OK
	if len(n.adapters) == 0 {
		return errors.New("no adapters")
	}

	// get the local adapter
	localAdapter := n.localAdapter
	n.log.Debug().Stringer("localAdapter", localAdapter).Msg("localAdapter")
	if _debug != nil {
		_debug("    - local_adapter: %r", localAdapter)
	}

	// build a generic APDU
	apdu, err := NewAPDU(NoArgs, NKW(KWCPCIUserData, pdu.GetPDUUserData())) // Note: upstream uses _APDU which looks like the _APDU class but is an alias to APDU
	if err != nil {
		return errors.Wrap(err, "error creating _APDU")
	}
	if err := pdu.(Encoder).Encode(apdu); err != nil {
		return errors.Wrap(err, "error encoding APDU")
	}
	n.log.Debug().Stringer("_APDU", apdu).Msg("apdu")
	if _debug != nil {
		_debug("    - apdu: %r", apdu)
	}

	// build an NPDU specific to where it is going
	npdu, err := NewNPDU(NoArgs, NKW(KWCPCIUserData, pdu.GetPDUUserData()))
	if err != nil {
		return errors.Wrap(err, "error creating NPDU")
	}
	if err = apdu.Encode(npdu); err != nil {
		return errors.Wrap(err, "error encoding NPDU")
	}
	n.log.Debug().Stringer("npdu", npdu).Msg("npdu")
	if _debug != nil {
		_debug("    - npdu: %r", npdu)
	}

	// the hop count always starts out big
	hopCount := uint8(0xff)
	npdu.SetNpduHopCount(&hopCount)

	// if this is route aware, use it for the destination
	if Settings.RouteAware && npdu.GetPDUDestination().AddrRoute != nil {
		// always a local station for now, in theory this could also be
		// a local broadcast address, remote station, or remote broadcast
		// but that is not supported by the patterns
		if npdu.GetPDUDestination().AddrRoute.AddrType != LOCAL_STATION_ADDRESS {
			panic("Route must be of type local station")
		}
		if _debug != nil {
			_debug("    - routed: %r", npdu.GetPDUDestination().AddrRoute)
		}

		switch pdu.GetPDUDestination().AddrType {
		case REMOTE_STATION_ADDRESS, REMOTE_BROADCAST_ADDRESS, GLOBAL_BROADCAST_ADDRESS:
			if _debug != nil {
				_debug("    - continue DADR: %r", apdu.GetPDUDestination())
			}
			npdu.SetNpduDADR(apdu.GetPDUDestination())
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
		npdu.SetNpduDADR(apdu.GetPDUDestination())

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
	if _debug != nil {
		_debug("    - dnet: %r", dnet)
	}

	// if the network matches the local adapter it's local
	if dnet == localAdapter.adapterNet {
		switch npdu.GetPDUDestination().AddrType {
		case REMOTE_STATION_ADDRESS:
			n.log.Debug().Msg("mapping remote station to local station")
			if _debug != nil {
				_debug("    - mapping remote station to local station")
			}
			localStation, err := NewLocalStation(npdu.GetPDUDestination().AddrAddress, nil)
			if err != nil {
				return errors.Wrap(err, "error building local station")
			}
			pdu.SetPDUDestination(localStation)
		case REMOTE_BROADCAST_ADDRESS:
			if _debug != nil {
				_debug("    - mapping remote broadcast to local broadcast")
			}
			pdu.SetPDUDestination(NewLocalBroadcast(nil))
		default:
			return errors.New("Addressing problem")
		}
		return localAdapter.ProcessNPDU(npdu)
	}

	// get it ready to send when the path is found
	npdu.SetPDUDestination(nil)
	npdu.SetNpduDADR(apdu.GetPDUDestination())

	// we might already be waiting for a path for this network
	if pendingNet, ok := n.pendingNets[nk(dnet)]; ok {
		n.log.Debug().Msg("already waiting for a path")
		if _debug != nil {
			_debug("    - already waiting for path")
		}
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
		if _debug != nil {
			_debug("    - router_info found: %r", routerInfo)
		}

		// check the path status
		dnetStatus := routerInfo.dnets[nk(dnet)]
		n.log.Debug().Interface("dnetStatus", dnetStatus).Msg("dnetStatus")
		if _debug != nil {
			_debug("    - dnet_status: %r", dnetStatus)
		}

		// fix the destination
		npdu.SetPDUDestination(routerInfo.address)

		// send it along
		return snetAdapter.ProcessNPDU(npdu)
	} else {
		n.log.Debug().Msg("no known path to network")
		if _debug != nil {
			_debug("    - no known path to network")
		}

		// add it to the list of packets waiting for the network
		netList, _ := n.pendingNets[nk(dnet)]
		netList = append(netList, npdu)
		n.pendingNets[nk(dnet)] = netList

		// build a request for the network and send it to all the adapters
		xnpdu, err := NewWhoIsRouterToNetwork(NoArgs, NoKWArgs(), WithWhoIsRouterToNetworkNet(*dnet))
		if err != nil {
			return errors.Wrap(err, "error building WhoIsRouterToNetwork")
		}
		xnpdu.SetPDUDestination(NewLocalBroadcast(nil))

		// send it to all the adapters
		for _, adapter := range n.adapters {
			if err := n.SapIndication(NA(adapter, xnpdu), NoKWArgs()); err != nil {
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
	if _debug != nil {
		_debug("process_npdu %r %r", adapter, npdu)
	}

	// make sure our configuration is OK
	if len(n.adapters) == 0 {
		return errors.New("no adapters")
	}

	var (
		processLocally bool
		forwardMessage bool
	)
	// check for source routing
	if npdu.GetNpduSADR() != nil && npdu.GetNpduSADR().AddrType != NULL_ADDRESS {
		n.log.Debug().Msg("check source path")
		if _debug != nil {
			_debug("    - check source path")
		}

		// see if this is attempting to spoof a directly connected network
		snet := npdu.GetSourceNetworkAddress()
		if _, ok := n.adapters[nk(snet)]; ok {
			n.log.Warn().Msg("path error (1)")
			return nil
		}

		// pass this new path along to the cache
		if err := n.routerInfoCache.UpdateRouterInfo(nk(adapter.adapterNet), npdu.GetPDUSource(), []uint16{*snet}, nil); err != nil {
			return errors.Wrap(err, "error updating router status")
		}
	}
	switch {
	// check for destination routing
	case npdu.GetNpduDADR() == nil || npdu.GetNpduDADR().AddrType == NULL_ADDRESS:
		n.log.Debug().Msg("no DADR")
		if _debug != nil {
			_debug("    - no DADR")
		}

		processLocally = adapter == n.localAdapter || npdu.GetControl().GetMessageTypeFieldPresent()
		forwardMessage = false
	case npdu.GetNpduDADR().AddrType == REMOTE_BROADCAST_ADDRESS:
		n.log.Debug().Msg("DADR is remote broadcast")
		if _debug != nil {
			_debug("    - DADR is remote broadcast")
		}

		if *npdu.GetNpduDADR().AddrNet == *adapter.adapterNet {
			n.log.Warn().Msg("path error (2)")
			return nil
		}

		processLocally = *npdu.GetNpduDADR().AddrNet == *n.localAdapter.adapterNet
		forwardMessage = true
	case npdu.GetNpduDADR().AddrType == REMOTE_STATION_ADDRESS:
		n.log.Debug().Msg("DADR is remote station")
		if _debug != nil {
			_debug("    - DADR is remote station")
		}

		if *npdu.GetNpduDADR().AddrNet == *adapter.adapterNet {
			n.log.Warn().Msg("path error (3)")
			return nil
		}

		processLocally = *npdu.GetNpduDADR().AddrNet == *n.localAdapter.adapterNet && bytes.Compare(npdu.GetNpduDADR().AddrAddress, n.localAdapter.adapterAddr.AddrAddress) == 0
		forwardMessage = !processLocally
	case npdu.GetNpduDADR().AddrType == GLOBAL_BROADCAST_ADDRESS:
		n.log.Debug().Msg("DADR is global broadcast")
		if _debug != nil {
			_debug("    - DADR is global broadcast")
		}

		processLocally = true
		forwardMessage = true
	default:
		n.log.Warn().Stringer("addrType", npdu.GetNpduDADR().AddrType).Msg("invalid destination address type:")
		return nil
	}

	n.log.Debug().Bool("processLocally", processLocally).Msg("processLocally")
	n.log.Debug().Bool("forwardMessage", forwardMessage).Msg("forwardMessage")
	if _debug != nil {
		_debug("    - processLocally: %r", processLocally)
		_debug("    - forwardMessage: %r", forwardMessage)
	}

	// application or network layer message
	if !npdu.GetControl().GetMessageTypeFieldPresent() {
		n.log.Debug().Msg("application layer message")
		if _debug != nil {
			_debug("    - application layer message")
		}

		if processLocally && n.HasServerPeer() {
			n.log.Trace().Msg("processing APDU locally")
			if _debug != nil {
				_debug("    - processing APDU locally")
			}
			// decode as a generic APDU
			apdu, err := NewAPDU(NoArgs, NKW(KWCPCIUserData, npdu.GetPDUUserData())) // Note: upstream uses _APDU which looks like the _APDU class but is an alias to APDU
			if err != nil {
				return errors.Wrap(err, "error creating APDU")
			}
			if err := apdu.Decode(DeepCopy[NPDU](npdu)); err != nil {
				return errors.Wrap(err, "error decoding APDU")
			}
			n.log.Debug().Stringer("apdu", apdu).Msg("apdu")
			if _debug != nil {
				_debug("    - apdu: %r", apdu)
			}

			// see if it needs to look routed
			if len(n.adapters) > 1 && adapter != n.localAdapter {
				// combine the source address
				if !npdu.GetControl().GetSourceSpecified() {
					remoteStationAddress, err := NewAddress(NA(adapter.adapterNet))
					if err != nil {
						return errors.Wrap(err, "error creating remote address")
					}
					apdu.SetPDUSource(remoteStationAddress)
				} else {
					apdu.SetPDUSource(npdu.GetPDUSource())
				}
				if Settings.RouteAware {
					apdu.GetPDUSource().AddrRoute = npdu.GetPDUSource()
				}

				// map the destination
				if !npdu.GetControl().GetDestinationSpecified() {
					apdu.SetPDUDestination(n.localAdapter.adapterAddr)
				} else if npdu.GetNpduDADR().AddrType == GLOBAL_BROADCAST_ADDRESS {
					apdu.SetPDUDestination(NewGlobalBroadcast(nil))
				} else if npdu.GetNpduDADR().AddrType == REMOTE_BROADCAST_ADDRESS {
					apdu.SetPDUDestination(NewLocalBroadcast(nil))
				} else {
					apdu.SetPDUDestination(n.localAdapter.adapterAddr)
				}
			} else {
				// combine the source address
				if npdu.GetControl().GetSourceSpecified() {
					apdu.SetPDUSource(npdu.GetNpduSADR())
					if Settings.RouteAware {
						n.log.Debug().Msg("adding route")
						if _debug != nil {
							_debug("    - adding route")
						}
						apdu.GetPDUSource().AddrRoute = npdu.GetPDUSource()
					}
				} else {
					apdu.SetPDUSource(npdu.GetPDUSource())
				}

				// pass along global broadcast
				if npdu.GetControl().GetDestinationSpecified() && npdu.GetNpduDADR().AddrType == GLOBAL_BROADCAST_ADDRESS {
					apdu.SetPDUDestination(NewGlobalBroadcast(nil))
				} else {
					apdu.SetPDUDestination(npdu.GetPDUDestination())
				}
			}

			n.log.Debug().Stringer("pduSource", apdu.GetPDUSource()).Msg("apdu.pduSource")
			n.log.Debug().Stringer("pduDestination", apdu.GetPDUDestination()).Msg("apdu.pduDestination")
			if _debug != nil {
				_debug("    - apdu.pduSource: %r", apdu.GetPDUSource())
				_debug("    - apdu.pduDestination: %r", apdu.GetPDUDestination())
			}

			// pass upstream to the application layer
			if err := n.Response(NA(apdu), NoKWArgs()); err != nil {
				return errors.Wrap(err, "error passing response")
			}
		}
	} else {
		n.log.Debug().Msg("network layer message")
		if _debug != nil {
			_debug("    - network layer message")
		}

		if processLocally {
			np, ok := NPDUTypes[npdu.GetNlm().GetMessageType()]
			if !ok {
				n.log.Debug().Uint8("messageType", npdu.GetNlm().GetMessageType()).Msg("unknown npdu type")
				if _debug != nil {
					_debug("    - unknown npdu type: %r", npdu.GetNPDUNetMessage())
				}
				return nil
			}
			n.log.Debug().Msg("processing NPDU locally")
			if _debug != nil {
				_debug("    - processing NPDU locally")
			}

			// do a deeper decode of the NPDU
			xpdu, err := np(NoArgs, NKW(KWCPCIUserData, npdu.GetPDUUserData()))
			if err != nil {
				return errors.Wrap(err, "error creating NPDU")
			}
			if err := xpdu.Decode(npdu.DeepCopy()); err != nil {
				return errors.Wrap(err, "error decoding NPDU")
			}

			// pass to the service element
			if err := n.SapRequest(NA(adapter, xpdu), NoKWArgs()); err != nil {
				return errors.Wrap(err, "error passing sap _request")
			}
		}
	}

	// might not need to forward this to other devices
	if !forwardMessage {
		n.log.Debug().Msg("no forwarding")
		if _debug != nil {
			_debug("    - no forwarding")
		}
		return nil
	}

	// make sure we're really a router
	if len(n.adapters) == 1 {
		n.log.Debug().Msg("not a router")
		if _debug != nil {
			_debug("    - not a router")
		}
		return nil
	}

	// make sure it hasn't looped
	if npdu.GetHopCount() != nil && *npdu.GetHopCount() == 0 {
		n.log.Debug().Msg("no more hops")
		if _debug != nil {
			_debug("    - no more hops")
		}
		return nil
	}

	// build a new NPDU to send to other adapters
	newpdu := DeepCopy[NPDU](npdu)
	newpdu.SetRootMessage(nil) // TODO: we would need to update it every modification we do below so we just nuke it for now

	// clear out the source and destination
	newpdu.SetPDUSource(nil)
	newpdu.SetPDUDestination(nil)

	// decrease the hop count
	newNpduHopCount := *npdu.GetHopCount() - 1
	newpdu.SetNpduHopCount(&newNpduHopCount)

	// set the source address
	if !npdu.GetControl().GetSourceSpecified() {
		newSADR, err := NewRemoteStation(adapter.adapterNet, npdu.GetPDUSource().AddrAddress, nil)
		if err != nil {
			return errors.Wrap(err, "error creating remote station")
		}
		newpdu.SetNpduSADR(newSADR)
	} else {
		newpdu.SetNpduSADR(npdu.GetNpduSADR())
	}

	// If this is a broadcast it goes everywhere
	if npdu.GetNpduDADR().AddrType == GLOBAL_BROADCAST_ADDRESS {
		n.log.Debug().Msg("global broadcasting")
		if _debug != nil {
			_debug("    - global broadcasting")
		}

		newpdu.SetPDUDestination(NewLocalBroadcast(nil))

		for _, xadapter := range n.adapters {
			if xadapter != adapter {
				if err := xadapter.ProcessNPDU(DeepCopy[NPDU](newpdu)); err != nil {
					n.log.Warn().Err(err).Msg("Error processing npdu")
				}
			}
		}
		return nil
	}

	if npdu.GetNpduDADR().AddrType == REMOTE_BROADCAST_ADDRESS || npdu.GetNpduDADR().AddrType == REMOTE_STATION_ADDRESS {
		dnet := npdu.GetNpduDADR().AddrNet
		n.log.Debug().Msg("remote station/broadcast")
		if _debug != nil {
			_debug("    - remote station/broadcast")
		}

		// see if this a locally connected network
		if xadapter, ok := n.adapters[nk(dnet)]; ok {
			if xadapter == adapter {
				n.log.Debug().Msg("path error (4)")
				if _debug != nil {
					_debug("    - path error (4)")
				}
				return nil
			}
			n.log.Debug().Stringer("adapter", adapter).Msg("found path via")
			if _debug != nil {
				_debug("    - found path via %r", xadapter)
			}

			// if this was a remote broadcast, it's now a local one
			if npdu.GetNpduDADR().AddrType == REMOTE_BROADCAST_ADDRESS {
				newpdu.SetPDUDestination(NewLocalBroadcast(nil))
			} else {
				newDADR, err := NewLocalStation(npdu.GetNpduDADR().AddrAddress, nil)
				if err != nil {
					return errors.Wrap(err, "error building local station")
				}
				newpdu.SetPDUDestination(newDADR)
			}

			// last leg in routing
			newpdu.SetNpduDADR(nil)

			// send the packet downstream
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
			if _debug != nil {
				_debug("    - found path via %r", routerInfo)
			}

			// the destination is the address of the router
			newpdu.SetPDUDestination(routerInfo.address)

			//  send the packet downstream
			if snetAdapter == nil {
				return errors.New("snetAdapter nil")
			}
			return snetAdapter.ProcessNPDU(DeepCopy[NPDU](newpdu))
		}

		n.log.Debug().Msg("No router info found")
		if _debug != nil {
			_debug("    - no router info found")
		}

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
			if err := n.SapIndication(NA(xadapter, NewPDU(NoArgs, NKW(KWCPCIDestination, pduDestination), WithRootMessage(xnpdu))), NoKWArgs()); err != nil {
				return errors.Wrap(err, "error sending indication")
			}
		}

		return nil
	}

	n.log.Debug().
		Interface("destinationNetworkAddress", npdu.GetDestinationNetworkAddress()).
		Interface("destinationAddress", npdu.GetDestinationAddress()).
		Msg("bad DADR")
	if _debug != nil {
		_debug("    - bad DADR: %r", npdu.GetNpduDADR())
	}
	return nil
}

func (n *NetworkServiceAccessPoint) SapIndication(args Args, kwArgs KWArgs) error {
	n.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("SapIndication")
	adapter := GA[*NetworkAdapter](args, 0)
	npdu := GA[NPDU](args, 1)
	if _debug != nil {
		_debug("sap_indication %r %r", adapter, npdu)
	}

	// encode it as a generic NPDU
	xpdu, err := NewNPDU(NoArgs, NKW(KWCPCIUserData, npdu.GetPDUUserData()))
	if err != nil {
		return errors.Wrap(err, "error building NPDU")
	}
	if err := npdu.Encode(xpdu); err != nil {
		return errors.Wrap(err, "error encoding NPDU")
	}
	// npdu._xpdu = xpdu // TODO: what does that mean?

	// tell the adapter to process the NPDU
	return adapter.ProcessNPDU(xpdu)
}

func (n *NetworkServiceAccessPoint) SapConfirmation(args Args, kwArgs KWArgs) error {
	n.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("SapConfirmation")
	adapter := GA[*NetworkAdapter](args, 0)
	npdu := GA[NPDU](args, 1)
	if _debug != nil {
		_debug("sap_confirmation %r %r", adapter, npdu)
	}

	// encode it as a generic NPDU
	xpdu, err := NewNPDU(NoArgs, NKW(KWCPCIUserData, npdu.GetPDUUserData()))
	if err != nil {
		return errors.Wrap(err, "error building NPDU")
	}
	if err := npdu.Encode(xpdu); err != nil {
		return errors.Wrap(err, "error encoding NPDU")
	}
	// npdu._xpdu = xpdu // TODO: what does this mean

	return adapter.ProcessNPDU(xpdu)
}

func (n *NetworkServiceAccessPoint) AlternateString() (string, bool) {
	if IsDebuggingActive() {
		return fmt.Sprintf("%r", n), true // Delegate to the format method
	}
	return "", false
}
