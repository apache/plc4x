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
	"slices"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/core"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/npdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

//go:generate plc4xGenerator -type=NetworkServiceElement -prefix=netservice_
type NetworkServiceElement struct {
	ApplicationServiceElementContract

	networkNumberIsTask time.Time

	// regular args
	argStartupDisabled bool `ignore:"true"`

	// pass through args
	argEID *int                       `ignore:"true"`
	argAse *ApplicationServiceElement `ignore:"true"`

	log zerolog.Logger
}

func NewNetworkServiceElement(localLog zerolog.Logger, opts ...func(*NetworkServiceElement)) (*NetworkServiceElement, error) {
	n := &NetworkServiceElement{
		log: localLog,
	}
	for _, opt := range opts {
		opt(n)
	}
	n.log.Trace().Interface("eid", n.argEID).Msg("NewNetworkServiceElement")
	var err error
	n.ApplicationServiceElementContract, err = NewApplicationServiceElement(localLog, OptionalOptionDual(n.argEID, n.argAse, WithApplicationServiceElementAseID))
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

func WithNetworkServiceElementEID(eid int, ase ApplicationServiceElement) func(*NetworkServiceElement) {
	return func(n *NetworkServiceElement) {
		n.argEID = &eid
		n.argAse = &ase
	}
}

func WithNetworkServiceElementStartupDisabled(startupDisabled bool) func(*NetworkServiceElement) {
	return func(n *NetworkServiceElement) {
		n.argStartupDisabled = startupDisabled
	}
}

func (n *NetworkServiceElement) Startup(_ Args, _ KWArgs) error {
	n.log.Debug().Msg("Startup")

	// reference the service access point
	sap := n.GetElementService().(*NetworkServiceAccessPoint) // TODO: hard cast but seems like adapters appears first in network service access point (so hard binding)
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
		if err := n.iamRouterToNetwork(NA(adapter, nil, netlist), NoKWArgs); err != nil {
			n.log.Debug().Err(err).Msg("I-Am-Router-To-Network failed")
		}
	}
	return nil
}

func (n *NetworkServiceElement) Indication(args Args, kwargs KWArgs) error {
	n.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")

	adapter := GA[*NetworkAdapter](args, 0)
	npdu := GA[NPDU](args, 1)

	switch message := npdu.GetRootMessage().(type) {
	case model.NPDU:
		switch nlm := message.GetNlm().(type) {
		case model.NLMWhoIsRouterToNetwork:
			return n.WhoIsRouteToNetwork(adapter, npdu, nlm)
		case model.NLMIAmRouterToNetwork:
			return n.IAmRouterToNetwork(adapter, npdu, nlm)
		case model.NLMICouldBeRouterToNetwork:
			return n.ICouldBeRouterToNetwork(adapter, nlm)
		case model.NLMRejectMessageToNetwork:
			return n.RejectRouterToNetwork(adapter, nlm)
		case model.NLMRouterBusyToNetwork:
			return n.RouterBusyToNetwork(adapter, nlm)
		case model.NLMRouterAvailableToNetwork:
			return n.RouterAvailableToNetwork(adapter, nlm)
		case model.NLMInitializeRoutingTable:
			return n.InitalizeRoutingTable(adapter, nlm)
		case model.NLMInitializeRoutingTableAck:
			return n.InitalizeRoutingTableAck(adapter, nlm)
		case model.NLMEstablishConnectionToNetwork:
			return n.EstablishConnectionToNetwork(adapter, nlm)
		case model.NLMDisconnectConnectionToNetwork:
			return n.DisconnectConnectionToNetwork(adapter, nlm)
		case model.NLMWhatIsNetworkNumber:
			return n.WhatIsNetworkNumber(adapter, nlm)
		case model.NLMNetworkNumberIs:
			return n.NetworkNumberIs(adapter, nlm)
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

	adapter := GA[*NetworkAdapter](args, 0)
	npdu := GA[NPDU](args, 1)

	switch message := npdu.GetRootMessage().(type) {
	case model.NPDU:
		switch nlm := message.GetNlm().(type) {
		case model.NLMWhoIsRouterToNetwork:
			return n.WhoIsRouteToNetwork(adapter, npdu, nlm)
		case model.NLMIAmRouterToNetwork:
			return n.IAmRouterToNetwork(adapter, npdu, nlm)
		case model.NLMICouldBeRouterToNetwork:
			return n.ICouldBeRouterToNetwork(adapter, nlm)
		case model.NLMRejectMessageToNetwork:
			return n.RejectRouterToNetwork(adapter, nlm)
		case model.NLMRouterBusyToNetwork:
			return n.RouterBusyToNetwork(adapter, nlm)
		case model.NLMRouterAvailableToNetwork:
			return n.RouterAvailableToNetwork(adapter, nlm)
		case model.NLMInitializeRoutingTable:
			return n.InitalizeRoutingTable(adapter, nlm)
		case model.NLMInitializeRoutingTableAck:
			return n.InitalizeRoutingTableAck(adapter, nlm)
		case model.NLMEstablishConnectionToNetwork:
			return n.EstablishConnectionToNetwork(adapter, nlm)
		case model.NLMDisconnectConnectionToNetwork:
			return n.DisconnectConnectionToNetwork(adapter, nlm)
		case model.NLMWhatIsNetworkNumber:
			return n.WhatIsNetworkNumber(adapter, nlm)
		case model.NLMNetworkNumberIs:
			return n.NetworkNumberIs(adapter, nlm)
		default:
			n.log.Debug().Stringer("nlm", nlm).Msg("Unhandled")
		}
	default:
		n.log.Trace().Msg("can only handle NPDU")
	}
	return nil
}

func (n *NetworkServiceElement) iamRouterToNetwork(args Args, _ KWArgs) error {
	adapter := GAO[*NetworkAdapter](args, 0, nil)
	destination := GAO[*Address](args, 1, nil)
	network := GAO[[]*uint16](args, 2, nil)
	n.log.Debug().Stringer("adapter", adapter).Stringer("destination", destination).Interface("network", network).Msg("IamRouterToNetwork")

	// reference the service access point
	sap := n.GetElementService().(*NetworkServiceAccessPoint) // TODO: hard cast but seems like adapters appears first in network service access point (so hard binding)
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
			destination, err = NewLocalStation(destination.AddrAddress, nil)
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
			destination, err = NewLocalStation(destination.AddrAddress, nil)
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

		if len(network) == 0 {
			return nil
		}
		for _, u := range network {
			if slices.Contains(netlist, *u) {
				netlist = append(netlist, *u)
			}
		}

		// build a response
		iamrtn, err := NewIAmRouterToNetwork()
		if err != nil {
			return errors.Wrap(err, "error creating IAM router to network")
		}
		iamrtn.SetPDUDestination(destination)
		n.log.Debug().Stringer("adapter", adapter).Stringer("iamrtn", iamrtn).Msg("adapter, iamrtn")

		// send it back
		if err := n.Request(NA(adapter, iamrtn), NoKWArgs); err != nil {
			return errors.Wrap(err, "error requesting NPDU")
		}
	}

	return nil
}

func (n *NetworkServiceElement) WhoIsRouteToNetwork(adapter *NetworkAdapter, npdu NPDU, nlm model.NLMWhoIsRouterToNetwork) error {
	n.log.Debug().Stringer("adapter", adapter).Stringer("npdu", npdu).Stringer("nlm", nlm).Msg("WhoIsRouteToNetwork")

	// reference the service access point
	sap := n.GetElementService().(*NetworkServiceAccessPoint) // TODO: check hard cast here...
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
			iamrtn.SetPDUDestination(npdu.GetPDUSource())

			// send it back
			if err := n.Response(NA(adapter, iamrtn), NoKWArgs); err != nil {
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
			iamrtn.SetPDUDestination(npdu.GetPDUSource())

			// send it back
			return n.Response(NA(adapter, iamrtn), NoKWArgs)
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
			iamrtn.SetPDUDestination(npdu.GetPDUSource())

			// send it back
			return n.Response(NA(adapter, iamrtn), NoKWArgs)
		} else {
			n.log.Trace().Msg("forwarding to other adapters")

			whoisrtn, err := NewWhoIsRouterToNetwork(WithWhoIsRouterToNetworkNet(dnet))
			if err != nil {
				return errors.Wrap(err, "error building WhoIsRouterToNetwork")
			}
			whoisrtn.SetPDUUserData(npdu.GetPDUUserData()) // TODO: upstream does this inline
			whoisrtn.SetPDUDestination(NewLocalBroadcast(nil))

			// if the request had a source forward it along
			if npdu.GetSourceNetworkAddress() != nil {
				whoisrtn.SetNpduSADR(npdu.GetNpduSADR())
			} else {
				station, err := NewRemoteStation(adapter.adapterNet, npdu.GetPDUSource().AddrAddress, nil)
				if err != nil {
					return errors.Wrap(err, "error building RemoteStation")
				}
				whoisrtn.SetNpduSADR(station)
			}

			// send it to all (other) adapters
			for _, xadapter := range sap.adapters {
				if xadapter != adapter {
					n.log.Debug().Stringer("xadapter", xadapter).Msg("Sending to adapter")
					if err := n.Request(NA(xadapter, whoisrtn), NoKWArgs); err != nil {
						return errors.Wrap(err, "error sending Who is router to network")
					}
				}
			}
			return nil
		}
	}
}

func (n *NetworkServiceElement) IAmRouterToNetwork(adapter *NetworkAdapter, npdu NPDU, nlm model.NLMIAmRouterToNetwork) error {
	n.log.Debug().Stringer("adapter", adapter).Stringer("npdu", npdu).Stringer("nlm", nlm).Msg("IAmRouterToNetwork")

	// reference the service access point
	sap := n.GetElementService().(*NetworkServiceAccessPoint) // TODO: check hard cast here...
	n.log.Debug().Stringer("sap", sap).Msg("sap")

	// pass along the service access point
	if err := sap.UpdateRouterReference(adapter.adapterNet, npdu.GetPDUSource(), nlm.GetDestinationNetworkAddresses()); err != nil {
		return errors.Wrap(err, "error updating router to network")
	}

	// if we're not a router, skip it
	if len(sap.adapters) == 1 {
		n.log.Trace().Msg("not a router")
	} else {
		n.log.Trace().Msg("forwarding other adapters")

		// Build a broadcast announcement
		iamrtn, err := NewIAmRouterToNetwork(WithIAmRouterToNetworkNetworkList(nlm.GetDestinationNetworkAddresses()...))
		if err != nil {
			return errors.Wrap(err, "error building IAmRouterToNetwork")
		}
		iamrtn.SetPDUDestination(NewLocalBroadcast(nil))

		// send it to all the connected adapters
		for _, xadapter := range sap.adapters {
			if xadapter != adapter {
				n.log.Debug().Stringer("xadapter", xadapter).Msg("Sending to adapter")
				if err := n.Request(NA(xadapter, iamrtn), NoKWArgs); err != nil {
					return errors.Wrap(err, "error sending I am router to network")
				}
			}
		}
	}

	// look for pending NPDUs for the networks
	for _, dnet := range nlm.GetDestinationNetworkAddresses() {
		pendingNpdus, ok := sap.pendingNets[nk(&dnet)]
		if ok {
			n.log.Debug().Int("pendingNpdus", len(pendingNpdus)).Uint16("dnet", dnet).Msg("pending NPDUs to dnet")

			// delete the references
			delete(sap.pendingNets, nk(&dnet))

			// now reprocess them
			for _, pendingNPDU := range pendingNpdus {
				n.log.Debug().Stringer("pendingNPDU", pendingNPDU).Msg("sending")

				// the destination is the address of the router
				pendingNPDU.SetPDUDestination(npdu.GetPDUSource())

				// sent the packet downstream
				if err := adapter.ProcessNPDU(pendingNPDU); err != nil {
					return errors.Wrap(err, "error processing pending NPDU")
				}
			}
		}
	}
	return nil
}

func (n *NetworkServiceElement) ICouldBeRouterToNetwork(adapter *NetworkAdapter, nlm model.NLMICouldBeRouterToNetwork) error {
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) RejectRouterToNetwork(adapter *NetworkAdapter, nlm model.NLMRejectMessageToNetwork) error {
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) RouterBusyToNetwork(adapter *NetworkAdapter, nlm model.NLMRouterBusyToNetwork) error {
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) RouterAvailableToNetwork(adapter *NetworkAdapter, nlm model.NLMRouterAvailableToNetwork) error {
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) InitalizeRoutingTable(adapter *NetworkAdapter, nlm model.NLMInitializeRoutingTable) error {
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) InitalizeRoutingTableAck(adapter *NetworkAdapter, nlm model.NLMInitializeRoutingTableAck) error {
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) EstablishConnectionToNetwork(adapter *NetworkAdapter, nlm model.NLMEstablishConnectionToNetwork) error {
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) DisconnectConnectionToNetwork(adapter *NetworkAdapter, nlm model.NLMDisconnectConnectionToNetwork) error {
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) WhatIsNetworkNumber(adapter *NetworkAdapter, nlm model.NLMWhatIsNetworkNumber) error {
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) NetworkNumberIs(adapter *NetworkAdapter, nlm model.NLMNetworkNumberIs) error {
	panic("not implemented") // TODO: implement me
}
