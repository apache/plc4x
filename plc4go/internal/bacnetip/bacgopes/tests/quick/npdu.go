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

package quick

import (
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"

	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/npdu"
)

func WhoIsRouterToNetwork(net uint16) *npdu.WhoIsRouterToNetwork {
	network, err := npdu.NewWhoIsRouterToNetwork(npdu.WithWhoIsRouterToNetworkNet(net))
	if err != nil {
		panic(err)
	}
	return network
}

func IAmRouterToNetwork(netList ...uint16) *npdu.IAmRouterToNetwork {
	network, err := npdu.NewIAmRouterToNetwork(npdu.WithIAmRouterToNetworkNetworkList(netList...))
	if err != nil {
		panic(err)
	}
	return network
}

func ICouldBeRouterToNetwork(net uint16, perf uint8) *npdu.ICouldBeRouterToNetwork {
	network, err := npdu.NewICouldBeRouterToNetwork(npdu.WithICouldBeRouterToNetworkNetwork(net), npdu.WithICouldBeRouterToNetworkPerformanceIndex(perf))
	if err != nil {
		panic(err)
	}
	return network
}

func RejectMessageToNetwork(reason uint8, dnet uint16) *npdu.RejectMessageToNetwork {
	network, err := npdu.NewRejectMessageToNetwork(npdu.WithRejectMessageToNetworkRejectionReason(readWriteModel.NLMRejectMessageToNetworkRejectReason(reason)), npdu.WithRejectMessageToNetworkDnet(dnet))
	if err != nil {
		panic(err)
	}
	return network
}

func RouterBusyToNetwork(netList ...uint16) *npdu.RouterBusyToNetwork {
	network, err := npdu.NewRouterBusyToNetwork(npdu.WithRouterBusyToNetworkDnet(netList))
	if err != nil {
		panic(err)
	}
	return network
}

func RouterAvailableToNetwork(netList ...uint16) *npdu.RouterAvailableToNetwork {
	network, err := npdu.NewRouterAvailableToNetwork(npdu.WithRouterAvailableToNetworkDnet(netList))
	if err != nil {
		panic(err)
	}
	return network
}

func InitializeRoutingTable(irtTable ...*npdu.RoutingTableEntry) *npdu.InitializeRoutingTable {
	network, err := npdu.NewInitializeRoutingTable(npdu.WithInitializeRoutingTableIrtTable(irtTable...))
	if err != nil {
		panic(err)
	}
	return network
}

func RoutingTableEntry(address uint16, portId uint8, portInfo []byte) *npdu.RoutingTableEntry {
	return npdu.NewRoutingTableEntry(
		npdu.WithRoutingTableEntryDestinationNetworkAddress(address),
		npdu.WithRoutingTableEntryPortId(portId),
		npdu.WithRoutingTableEntryPortInfo(portInfo),
	)
}

func InitializeRoutingTableAck(irtaTable ...*npdu.RoutingTableEntry) *npdu.InitializeRoutingTableAck {
	network, err := npdu.NewInitializeRoutingTableAck(npdu.WithInitializeRoutingTableAckIrtaTable(irtaTable...))
	if err != nil {
		panic(err)
	}
	return network
}

func EstablishConnectionToNetwork(dnet uint16, terminationTime uint8) *npdu.EstablishConnectionToNetwork {
	network, err := npdu.NewEstablishConnectionToNetwork(npdu.WithEstablishConnectionToNetworkDNET(dnet), npdu.WithEstablishConnectionToNetworkTerminationTime(terminationTime))
	if err != nil {
		panic(err)
	}
	return network
}

func DisconnectConnectionToNetwork(dnet uint16) *npdu.DisconnectConnectionToNetwork {
	network, err := npdu.NewDisconnectConnectionToNetwork(npdu.WithDisconnectConnectionToNetworkDNET(dnet))
	if err != nil {
		panic(err)
	}
	return network
}

func WhatIsNetworkNumber(dnet uint16) *npdu.WhatIsNetworkNumber {
	network, err := npdu.NewWhatIsNetworkNumber()
	if err != nil {
		panic(err)
	}
	return network
}

func NetworkNumberIs(net uint16, flag bool) *npdu.NetworkNumberIs {
	network, err := npdu.NewNetworkNumberIs(npdu.WithNetworkNumberIsNET(net), npdu.WithNetworkNumberIsTerminationConfigured(flag))
	if err != nil {
		panic(err)
	}
	return network
}
