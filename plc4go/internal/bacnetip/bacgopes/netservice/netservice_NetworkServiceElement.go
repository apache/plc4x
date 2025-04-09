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
	"fmt"
	"slices"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/core"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/npdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

//go:generate go tool plc4xGenerator -type=NetworkServiceElement -prefix=netservice_
type NetworkServiceElement struct {
	ApplicationServiceElementContract
	*DefaultRFormatter `ignore:"true"`

	networkNumberIsTask time.Time

	// regular args
	argStartupDisabled bool `ignore:"true"`

	log zerolog.Logger
}

func NewNetworkServiceElement(localLog zerolog.Logger, options ...Option) (*NetworkServiceElement, error) {
	n := &NetworkServiceElement{
		DefaultRFormatter: NewDefaultRFormatter(),
		log:               localLog,
	}
	ApplyAppliers(options, n)
	optionsForParent := AddLeafTypeIfAbundant(options, n)
	var err error
	n.ApplicationServiceElementContract, err = NewApplicationServiceElement(localLog, optionsForParent...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating application service element")
	}
	if _debug != nil {
		_debug("__init__ eid=%r", n.GetElementId())
	}

	// network number is timeout
	n.networkNumberIsTask = time.Time{}

	// if starting up is enabled defer our startup function
	if !n.argStartupDisabled {
		Deferred(n.Startup, NoArgs, NoKWArgs())
	}
	return n, nil
}

func WithNetworkServiceElementStartupDisabled(startupDisabled bool) GenericApplier[*NetworkServiceElement] {
	return WrapGenericApplier(func(n *NetworkServiceElement) { n.argStartupDisabled = startupDisabled })
}

func (n *NetworkServiceElement) Startup(_ Args, _ KWArgs) error {
	n.log.Debug().Msg("Startup")
	if _debug != nil {
		_debug("startup")
	}

	// reference the service access point
	sap := n.GetElementService().(*NetworkServiceAccessPoint) // TODO: hard cast but seems like adapters appears first in network service access point (so hard binding)
	n.log.Debug().Stringer("sap", sap).Msg("sap")
	if _debug != nil {
		_debug("    - sap: %r", sap)
	}

	// loop through all the adapters
	for _, adapter := range sap.adapters {
		n.log.Debug().Stringer("adapter", adapter).Msg("adapter")
		if _debug != nil {
			_debug("    - adapter: %r", adapter)
		}

		if adapter.adapterNet == nil {
			n.log.Trace().Msg("skipping, unknown net")
			if _debug != nil {
				_debug("    - skipping, unknown net")
			}
			continue
		}
		if adapter.adapterAddr == nil {
			n.log.Trace().Msg("skipping, unknown addr")
			if _debug != nil {
				_debug("    - skipping, unknown addr")
			}
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
			if _debug != nil {
				_debug("    - skipping, no netlist")
			}
			continue
		}

		// pass this along to the cache -- on hold #213
		// sap.router_info_cache.update_router_info(adapter.adapterNet, adapter.adapterAddr, netlist)

		// send an announcement
		if err := n.iamRouterToNetwork(NA(adapter, nil, netlist), NoKWArgs()); err != nil {
			n.log.Debug().Err(err).Msg("I-Am-Router-To-Network failed")
		}
	}
	return nil
}

func (n *NetworkServiceElement) Indication(args Args, kwArgs KWArgs) error {
	n.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Indication")

	adapter := GA[*NetworkAdapter](args, 0)
	npdu := GA[NPDU](args, 1)
	if _debug != nil {
		_debug("indication %r %r", adapter, npdu)
	}

	switch message := npdu.GetRootMessage().(type) {
	case model.NPDU:
		switch nlm := message.GetNlm().(type) {
		case model.NLMWhoIsRouterToNetwork:
			return n.WhoIsRouterToNetwork(adapter, npdu, nlm)
		case model.NLMIAmRouterToNetwork:
			return n.IAmRouterToNetwork(adapter, npdu, nlm)
		case model.NLMICouldBeRouterToNetwork:
			return n.ICouldBeRouterToNetwork(adapter, npdu, nlm)
		case model.NLMRejectMessageToNetwork:
			return n.RejectRouterToNetwork(adapter, npdu, nlm)
		case model.NLMRouterBusyToNetwork:
			return n.RouterBusyToNetwork(adapter, npdu, nlm)
		case model.NLMRouterAvailableToNetwork:
			return n.RouterAvailableToNetwork(adapter, npdu, nlm)
		case model.NLMInitializeRoutingTable:
			return n.InitalizeRoutingTable(adapter, npdu, nlm)
		case model.NLMInitializeRoutingTableAck:
			return n.InitalizeRoutingTableAck(adapter, npdu, nlm)
		case model.NLMEstablishConnectionToNetwork:
			return n.EstablishConnectionToNetwork(adapter, npdu, nlm)
		case model.NLMDisconnectConnectionToNetwork:
			return n.DisconnectConnectionToNetwork(adapter, npdu, nlm)
		case model.NLMWhatIsNetworkNumber:
			return n.WhatIsNetworkNumber(adapter, npdu, nlm)
		case model.NLMNetworkNumberIs:
			return n.NetworkNumberIs(adapter, npdu, nlm)
		default:
			n.log.Debug().Stringer("nlm", nlm).Msg("Unhandled")
		}
	default:
		n.log.Trace().Msg("can only handle NPDU")
	}
	return nil
}

func (n *NetworkServiceElement) Confirmation(args Args, kwArgs KWArgs) error {
	n.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Confirmation")

	adapter := GA[*NetworkAdapter](args, 0)
	npdu := GA[NPDU](args, 1)
	if _debug != nil {
		_debug("confirmation %r %r", adapter, npdu)
	}

	switch message := npdu.GetRootMessage().(type) {
	case model.NPDU:
		switch nlm := message.GetNlm().(type) {
		case model.NLMWhoIsRouterToNetwork:
			return n.WhoIsRouterToNetwork(adapter, npdu, nlm)
		case model.NLMIAmRouterToNetwork:
			return n.IAmRouterToNetwork(adapter, npdu, nlm)
		case model.NLMICouldBeRouterToNetwork:
			return n.ICouldBeRouterToNetwork(adapter, npdu, nlm)
		case model.NLMRejectMessageToNetwork:
			return n.RejectRouterToNetwork(adapter, npdu, nlm)
		case model.NLMRouterBusyToNetwork:
			return n.RouterBusyToNetwork(adapter, npdu, nlm)
		case model.NLMRouterAvailableToNetwork:
			return n.RouterAvailableToNetwork(adapter, npdu, nlm)
		case model.NLMInitializeRoutingTable:
			return n.InitalizeRoutingTable(adapter, npdu, nlm)
		case model.NLMInitializeRoutingTableAck:
			return n.InitalizeRoutingTableAck(adapter, npdu, nlm)
		case model.NLMEstablishConnectionToNetwork:
			return n.EstablishConnectionToNetwork(adapter, npdu, nlm)
		case model.NLMDisconnectConnectionToNetwork:
			return n.DisconnectConnectionToNetwork(adapter, npdu, nlm)
		case model.NLMWhatIsNetworkNumber:
			return n.WhatIsNetworkNumber(adapter, npdu, nlm)
		case model.NLMNetworkNumberIs:
			return n.NetworkNumberIs(adapter, npdu, nlm)
		default:
			n.log.Debug().Stringer("nlm", nlm).Msg("Unhandled")
		}
	default:
		n.log.Trace().Msg("can only handle NPDU")
	}
	return nil
}

func (n *NetworkServiceElement) iamRouterToNetwork(args Args, _ KWArgs) error {
	adapter, _ := GAO[*NetworkAdapter](args, 0, nil)
	destination, _ := GAO[*Address](args, 1, nil)
	network, _ := GAO[[]*uint16](args, 2, nil)
	n.log.Debug().Stringer("adapter", adapter).Stringer("destination", destination).Interface("network", network).Msg("IamRouterToNetwork")
	if _debug != nil {
		_debug("i_am_router_to_network %r %r %r", adapter, destination, network)
	}

	// reference the service access point
	sap := n.GetElementService().(*NetworkServiceAccessPoint) // TODO: hard cast but seems like adapters appears first in network service access point (so hard binding)
	n.log.Debug().Interface("sap", sap).Msg("SAP")
	if _debug != nil {
		_debug("    - sap: %r", sap)
	}

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
	if _debug != nil {
		_debug("    - adapter, destination, network: %r, %r, %r", adapter, destination, network)
	}

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
		iamrtn, err := NewIAmRouterToNetwork(Nothing())
		if err != nil {
			return errors.Wrap(err, "error creating IAM router to network")
		}
		iamrtn.SetPDUDestination(destination)
		n.log.Debug().Stringer("adapter", adapter).Stringer("iamrtn", iamrtn).Msg("adapter, iamrtn")
		if _debug != nil {
			_debug("    - adapter, iamrtn: %r, %r", adapter, iamrtn)
		}

		// send it back
		if err := n.Request(NA(adapter, iamrtn), NoKWArgs()); err != nil {
			return errors.Wrap(err, "error requesting NPDU")
		}
	}

	return nil
}

func (n *NetworkServiceElement) WhoIsRouterToNetwork(adapter *NetworkAdapter, npdu NPDU, nlm model.NLMWhoIsRouterToNetwork) error {
	n.log.Debug().Stringer("adapter", adapter).Stringer("npdu", npdu).Stringer("nlm", nlm).Msg("WhoIsRouteToNetwork")
	if _debug != nil {
		_debug("WhoIsRouterToNetwork %r %r", adapter, npdu)
	}

	// reference the service access point
	sap := n.GetElementService().(*NetworkServiceAccessPoint) // TODO: check hard cast here...
	n.log.Debug().Stringer("sap", sap).Msg("sap")
	if _debug != nil {
		_debug("    - sap: %r", sap)
	}

	// if we're not a router, skip it
	if len(sap.adapters) == 1 {
		n.log.Trace().Msg("not a router")
		if _debug != nil {
			_debug("    - not a router")
		}
		return nil
	}
	if destinationNetworkAddress := nlm.GetDestinationNetworkAddress(); destinationNetworkAddress == nil {
		// requesting all networks
		n.log.Trace().Msg("requesting all network")
		if _debug != nil {
			_debug("    - requesting all networks")
		}

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
			if _debug != nil {
				_debug("    - found these: %r", netlist)
			}

			// build a response
			iamrtn, err := NewIAmRouterToNetwork(NoArgs, NewKWArgs(KWCPCIUserData, npdu.GetPDUUserData()), WithIAmRouterToNetworkNetworkList(netlist...))
			if err != nil {
				return errors.Wrap(err, "error building IAmRouterToNetwork")
			}
			iamrtn.SetPDUDestination(npdu.GetPDUSource())

			// send it back
			if err := n.Response(NA(adapter, iamrtn), NoKWArgs()); err != nil {
				return errors.Wrap(err, "error sendinf the response")
			}
		}
		return nil
	} else {
		wirtnNetwork := *destinationNetworkAddress
		// requesting a specific network
		n.log.Debug().Uint16("wirtnNetwork", wirtnNetwork).Msg("requesting specific network")
		if _debug != nil {
			_debug("    - requesting specific network: %r", wirtnNetwork)
		}
		dnet := wirtnNetwork

		// check the directly connected networks
		if otherAdapter, ok := sap.adapters[nk(&dnet)]; ok {
			n.log.Trace().Msg("directly connected")
			if _debug != nil {
				_debug("    - directly connected")
			}

			if otherAdapter == adapter {
				n.log.Trace().Msg("same network")
				if _debug != nil {
					_debug("    - same network")
				}
				return nil
			}

			// build a response
			iamrtn, err := NewIAmRouterToNetwork(NoArgs, NewKWArgs(KWCPCIUserData, npdu.GetPDUUserData()), WithIAmRouterToNetworkNetworkList(dnet))
			if err != nil {
				return errors.Wrap(err, "error building IAmRouterToNetwork")
			}
			iamrtn.SetPDUDestination(npdu.GetPDUSource())

			// send it back
			return n.Response(NA(adapter, iamrtn), NoKWArgs())
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
			if _debug != nil {
				_debug("    - router found: %r", routerInfo)
			}

			if snetAdapter == adapter {
				n.log.Trace().Msg("same network")
				if _debug != nil {
					_debug("    - same network")
				}
				return nil
			}

			// build a response
			iamrtn, err := NewIAmRouterToNetwork(NoArgs, NewKWArgs(KWCPCIUserData, npdu.GetPDUUserData()), WithIAmRouterToNetworkNetworkList(dnet))
			if err != nil {
				return errors.Wrap(err, "error building IAmRouterToNetwork")
			}
			iamrtn.SetPDUDestination(npdu.GetPDUSource())

			// send it back
			return n.Response(NA(adapter, iamrtn), NoKWArgs())
		} else {
			n.log.Trace().Msg("forwarding to other adapters")
			if _debug != nil {
				_debug("    - forwarding to other adapters")
			}

			whoisrtn, err := NewWhoIsRouterToNetwork(NoArgs, NewKWArgs(KWCPCIUserData, npdu.GetPDUUserData()), WithWhoIsRouterToNetworkNet(dnet))
			if err != nil {
				return errors.Wrap(err, "error building WhoIsRouterToNetwork")
			}
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
			if _debug != nil {
				_debug("    - whoisrtn: %r", whoisrtn)
			}

			// send it to all (other) adapters
			for _, xadapter := range sap.adapters {
				if xadapter != adapter {
					n.log.Debug().Stringer("xadapter", xadapter).Msg("Sending to adapter")
					if _debug != nil {
						_debug("    - sending on adapter: %r", xadapter)
					}
					if err := n.Request(NA(xadapter, whoisrtn), NoKWArgs()); err != nil {
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
	if _debug != nil {
		_debug("IAmRouterToNetwork %r %r", adapter, npdu)
	}

	// reference the service access point
	sap := n.GetElementService().(*NetworkServiceAccessPoint) // TODO: check hard cast here...
	n.log.Debug().Stringer("sap", sap).Msg("sap")
	if _debug != nil {
		_debug("    - sap: %r", sap)
	}

	// pass along the service access point
	if err := sap.UpdateRouterReference(adapter.adapterNet, npdu.GetPDUSource(), nlm.GetDestinationNetworkAddresses()); err != nil {
		return errors.Wrap(err, "error updating router to network")
	}

	// if we're not a router, skip it
	if len(sap.adapters) == 1 {
		n.log.Trace().Msg("not a router")
		if _debug != nil {
			_debug("    - not a router")
		}
	} else {
		n.log.Trace().Msg("forwarding other adapters")
		if _debug != nil {
			_debug("    - forwarding to other adapters")
		}

		// Build a broadcast announcement
		iamrtn, err := NewIAmRouterToNetwork(NoArgs, NoKWArgs(), WithIAmRouterToNetworkNetworkList(nlm.GetDestinationNetworkAddresses()...))
		if err != nil {
			return errors.Wrap(err, "error building IAmRouterToNetwork")
		}
		iamrtn.SetPDUDestination(NewLocalBroadcast(nil))

		// send it to all the connected adapters
		for _, xadapter := range sap.adapters {
			if xadapter != adapter {
				n.log.Debug().Stringer("xadapter", xadapter).Msg("Sending to adapter")
				if _debug != nil {
					_debug("    - sending on adapter: %r", xadapter)
				}
				if err := n.Request(NA(xadapter, iamrtn), NoKWArgs()); err != nil {
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
			if _debug != nil {
				_debug("    - %d pending to %r", len(pendingNpdus), dnet)
			}

			// delete the references
			delete(sap.pendingNets, nk(&dnet))

			// now reprocess them
			for _, pendingNPDU := range pendingNpdus {
				n.log.Debug().Stringer("pendingNPDU", pendingNPDU).Msg("sending")
				if _debug != nil {
					_debug("    - sending %r", pendingNPDU)
				}

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

func (n *NetworkServiceElement) ICouldBeRouterToNetwork(adapter *NetworkAdapter, npdu NPDU, nlm model.NLMICouldBeRouterToNetwork) error {
	if _debug != nil {
		_debug("ICouldBeRouterToNetwork %r %r", adapter, npdu)
	}
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) RejectRouterToNetwork(adapter *NetworkAdapter, npdu NPDU, nlm model.NLMRejectMessageToNetwork) error {
	if _debug != nil {
		_debug("RejectMessageToNetwork %r %r", adapter, npdu)
	}
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) RouterBusyToNetwork(adapter *NetworkAdapter, npdu NPDU, nlm model.NLMRouterBusyToNetwork) error {
	if _debug != nil {
		_debug("RouterBusyToNetwork %r %r", adapter, npdu)
	}
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) RouterAvailableToNetwork(adapter *NetworkAdapter, npdu NPDU, nlm model.NLMRouterAvailableToNetwork) error {
	if _debug != nil {
		_debug("RouterAvailableToNetwork %r %r", adapter, npdu)
	}
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) InitalizeRoutingTable(adapter *NetworkAdapter, npdu NPDU, nlm model.NLMInitializeRoutingTable) error {
	if _debug != nil {
		_debug("InitializeRoutingTable %r %r", adapter, npdu)
	}
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) InitalizeRoutingTableAck(adapter *NetworkAdapter, npdu NPDU, nlm model.NLMInitializeRoutingTableAck) error {
	if _debug != nil {
		_debug("InitializeRoutingTableAck %r %r", adapter, npdu)
	}
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) EstablishConnectionToNetwork(adapter *NetworkAdapter, npdu NPDU, nlm model.NLMEstablishConnectionToNetwork) error {
	if _debug != nil {
		_debug("EstablishConnectionToNetwork %r %r", adapter, npdu)
	}
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) DisconnectConnectionToNetwork(adapter *NetworkAdapter, npdu NPDU, nlm model.NLMDisconnectConnectionToNetwork) error {
	if _debug != nil {
		_debug("DisconnectConnectionToNetwork %r %r", adapter, npdu)
	}
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) WhatIsNetworkNumber(adapter *NetworkAdapter, npdu NPDU, nlm model.NLMWhatIsNetworkNumber) error {
	if _debug != nil {
		_debug("WhatIsNetworkNumber %r %r", adapter, npdu)
	}
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) NetworkNumberIs(adapter *NetworkAdapter, npdu NPDU, nlm model.NLMNetworkNumberIs) error {
	if _debug != nil {
		_debug("NetworkNumberIs %r %r", adapter, npdu)
	}
	panic("not implemented") // TODO: implement me
}

func (n *NetworkServiceElement) AlternateString() (string, bool) {
	if IsDebuggingActive() {
		return fmt.Sprintf("%r", n), true // Delegate to format
	}
	return "", false
}
