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
	"slices"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

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
	npdu := args.Get1NPDU()

	switch message := npdu.GetRootMessage().(type) {
	case model.NPDUExactly:
		switch nlm := message.GetNlm().(type) {
		case model.NLMWhoIsRouterToNetwork:
			return n.WhoIsRouteToNetwork(adapter, npdu, nlm)
		case model.NLMIAmRouterToNetwork:
			n.IAmRouterToNetwork(adapter, nlm)
		case model.NLMICouldBeRouterToNetwork:
			n.ICouldBeRouterToNetwork(adapter, nlm)
		case model.NLMRejectMessageToNetwork:
			n.RejectRouterToNetwork(adapter, nlm)
		case model.NLMRouterBusyToNetwork:
			n.RouterBusyToNetwork(adapter, nlm)
		case model.NLMRouterAvailableToNetwork:
			n.RouterAvailableToNetwork(adapter, nlm)
		case model.NLMInitializeRoutingTable:
			n.InitalizeRoutingTable(adapter, nlm)
		case model.NLMInitializeRoutingTableAck:
			n.InitalizeRoutingTableAck(adapter, nlm)
		case model.NLMEstablishConnectionToNetwork:
			n.EstablishConnectionToNetwork(adapter, nlm)
		case model.NLMDisconnectConnectionToNetwork:
			n.DisconnectConnectionToNetwork(adapter, nlm)
		case model.NLMWhatIsNetworkNumber:
			n.WhatIsNetworkNumber(adapter, nlm)
		case model.NLMNetworkNumberIs:
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
	npdu := args.Get1NPDU()

	switch message := npdu.GetRootMessage().(type) {
	case model.NPDUExactly:
		switch nlm := message.GetNlm().(type) {
		case model.NLMWhoIsRouterToNetwork:
			return n.WhoIsRouteToNetwork(adapter, npdu, nlm)
		case model.NLMIAmRouterToNetwork:
			n.IAmRouterToNetwork(adapter, nlm)
		case model.NLMICouldBeRouterToNetwork:
			n.ICouldBeRouterToNetwork(adapter, nlm)
		case model.NLMRejectMessageToNetwork:
			n.RejectRouterToNetwork(adapter, nlm)
		case model.NLMRouterBusyToNetwork:
			n.RouterBusyToNetwork(adapter, nlm)
		case model.NLMRouterAvailableToNetwork:
			n.RouterAvailableToNetwork(adapter, nlm)
		case model.NLMInitializeRoutingTable:
			n.InitalizeRoutingTable(adapter, nlm)
		case model.NLMInitializeRoutingTableAck:
			n.InitalizeRoutingTableAck(adapter, nlm)
		case model.NLMEstablishConnectionToNetwork:
			n.EstablishConnectionToNetwork(adapter, nlm)
		case model.NLMDisconnectConnectionToNetwork:
			n.DisconnectConnectionToNetwork(adapter, nlm)
		case model.NLMWhatIsNetworkNumber:
			n.WhatIsNetworkNumber(adapter, nlm)
		case model.NLMNetworkNumberIs:
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
			adapter, ok = sap.adapters[nk(destination.AddrNet)]
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
			adapter, ok = sap.adapters[nk(destination.AddrNet)]
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
		iamrtn := model.NewNLMIAmRouterToNetwork(netlist, 0)

		npdu, err := buildNPDU(0, nil, destination, false, model.NPDUNetworkPriority_NORMAL_MESSAGE, iamrtn, nil)
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

func (n *NetworkServiceElement) WhoIsRouteToNetwork(adapter *NetworkAdapter, npdu NPDU, nlm model.NLMWhoIsRouterToNetwork) error {
	n.log.Debug().Stringer("adapter", adapter).Stringer("nlm", nlm).Msg("WhoIsRouteToNetwork")

	// reference the service access point
	sap := n.elementService.(*NetworkServiceAccessPoint) // TODO: check hard cast here...
	n.log.Debug().Stringer("sap", sap).Msg("sap")

	// if we're not a router, skip it
	if len(sap.adapters) == 1 {
		n.log.Trace().Msg("not a router")
		return nil
	}
	if destinationNetworkAddress := nlm.GetDestinationNetworkAddress(); destinationNetworkAddress == nil {
		// requesting all networks
		n.log.Trace().Msg("requesting all network")

		// build a list of reachable networks
		var netlist []uint16

		for _, xadapter := range sap.adapters {
			if xadapter == adapter {
				continue
			}

			// add the dire network
			netlist = append(netlist, *xadapter.adapterNet)

			// add the other reachable networks? // TODO: upstream todo...
		}

		if len(netlist) > 0 {
			n.log.Debug().Uints16("netlist", netlist).Msg("found these")

			// build a response
			iamrtn, err := NewIAmRouterToNetwork(WithIAmRouterToNetworkNetworkList(netlist...))
			if err != nil {
				return errors.Wrap(err, "error building IAmRouterToNetwork")
			}
			iamrtn.SetPDUUserData(npdu.GetPDUUserData()) // TODO: upstream does this inline
			iamrtn.pduDestination = npdu.GetPDUSource()

			// send it back
			if err := n.Response(NewArgs(adapter, iamrtn), NoKWArgs); err != nil {
				return errors.Wrap(err, "error sendinf the response")
			}
		}
		return nil
	} else {
		wirtnNetwork := *destinationNetworkAddress
		// requesting a specific network
		n.log.Debug().Uint16("wirtnNetwork", wirtnNetwork).Msg("requesting specific network")
		dnet := wirtnNetwork

		// check the directly connected networks
		if otherAdapter, ok := sap.adapters[nk(&dnet)]; ok {
			n.log.Trace().Msg("directly connected")

			if otherAdapter == adapter {
				n.log.Trace().Msg("same network")
				return nil
			}

			// build a response
			iamrtn, err := NewIAmRouterToNetwork(WithIAmRouterToNetworkNetworkList(dnet))
			if err != nil {
				return errors.Wrap(err, "error building IAmRouterToNetwork")
			}
			iamrtn.SetPDUUserData(npdu.GetPDUUserData()) // TODO: upstream does this inline
			iamrtn.pduDestination = npdu.GetPDUSource()

			// send it back
			return n.Response(NewArgs(adapter, iamrtn), NoKWArgs)
		}

		// look for routing information from the network of one of our
		// adapters to the destination network
		var snetAdapter *NetworkAdapter
		var routerInfo *RouterInfo
		for snet, _snetAdapter := range sap.adapters {
			snetAdapter = _snetAdapter
			routerInfo = sap.routerInfoCache.GetRouterInfo(snet, nk(&dnet))
			if routerInfo != nil {
				break
			}
		}

		// found a path
		if routerInfo != nil {
			n.log.Debug().Stringer("routerInfo", routerInfo).Msg("router round")

			if snetAdapter == adapter {
				n.log.Trace().Msg("same network")
				return nil
			}

			// build a response
			iamrtn, err := NewIAmRouterToNetwork(WithIAmRouterToNetworkNetworkList(dnet))
			if err != nil {
				return errors.Wrap(err, "error building IAmRouterToNetwork")
			}
			iamrtn.SetPDUUserData(npdu.GetPDUUserData()) // TODO: upstream does this inline
			iamrtn.pduDestination = npdu.GetPDUSource()

			// send it back
			return n.Response(NewArgs(adapter, iamrtn), NoKWArgs)
		} else {
			n.log.Trace().Msg("forwarding to other adapters")

			whoisrtn, err := NewWhoIsRouterToNetwork(WithWhoIsRouterToNetworkNet(dnet))
			if err != nil {
				return errors.Wrap(err, "error building WhoIsRouterToNetwork")
			}
			whoisrtn.SetPDUUserData(npdu.GetPDUUserData()) // TODO: upstream does this inline
			whoisrtn.pduDestination = NewLocalBroadcast(nil)

			// if the request had a source forward it along
			if npdu.GetSourceNetworkAddress() != nil {
				panic("we need to forward SADR on NPDU as this it different from pdu source")
				// whoisrtn.npduSADR = npdu.npduSADR
			} else {
				panic("we need to forward SADR on NPDU as this it different from pdu source")
				// whoisrtn.npduSADR = RemoteStation(adapter.adapterNet, npdu.pduSource.addrAddr)
			}

			// send it to all (other) adapters
			for _, xadapter := range sap.adapters {
				if xadapter != adapter {
					n.log.Debug().Stringer("xadapter", xadapter).Msg("Sending to adapter")
					if err := n.Request(NewArgs(xadapter, whoisrtn), NoKWArgs); err != nil {
						return errors.Wrap(err, "error sending Who is router to network")
					}
				}
			}
			return nil
		}
	}
	return nil
}

func (n *NetworkServiceElement) IAmRouterToNetwork(adapter *NetworkAdapter, nlm model.NLMIAmRouterToNetwork) {
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) ICouldBeRouterToNetwork(adapter *NetworkAdapter, nlm model.NLMICouldBeRouterToNetwork) {
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) RejectRouterToNetwork(adapter *NetworkAdapter, nlm model.NLMRejectMessageToNetwork) {
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) RouterBusyToNetwork(adapter *NetworkAdapter, nlm model.NLMRouterBusyToNetwork) {
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) RouterAvailableToNetwork(adapter *NetworkAdapter, nlm model.NLMRouterAvailableToNetwork) {
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) InitalizeRoutingTable(adapter *NetworkAdapter, nlm model.NLMInitializeRoutingTable) {
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) InitalizeRoutingTableAck(adapter *NetworkAdapter, nlm model.NLMInitializeRoutingTableAck) {
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) EstablishConnectionToNetwork(adapter *NetworkAdapter, nlm model.NLMEstablishConnectionToNetwork) {
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) DisconnectConnectionToNetwork(adapter *NetworkAdapter, nlm model.NLMDisconnectConnectionToNetwork) {
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) WhatIsNetworkNumber(adapter *NetworkAdapter, nlm model.NLMWhatIsNetworkNumber) {
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) NetworkNumberIs(adapter *NetworkAdapter, nlm model.NLMNetworkNumberIs) {
	panic("not implemented") // TODO: implement me
}
