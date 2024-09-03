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

package constructors

import (
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"

	"github.com/apache/plc4x/plc4go/internal/bacnetip"
)

func WhoIsRouterToNetwork(net uint16) *bacnetip.WhoIsRouterToNetwork {
	network, err := bacnetip.NewWhoIsRouterToNetwork(bacnetip.WithWhoIsRouterToNetworkNet(net))
	if err != nil {
		panic(err)
	}
	return network
}

func IAmRouterToNetwork(netList ...uint16) *bacnetip.IAmRouterToNetwork {
	network, err := bacnetip.NewIAmRouterToNetwork(bacnetip.WithIAmRouterToNetworkNetworkList(netList...))
	if err != nil {
		panic(err)
	}
	return network
}

func ICouldBeRouterToNetwork(net uint16, perf uint8) *bacnetip.ICouldBeRouterToNetwork {
	network, err := bacnetip.NewICouldBeRouterToNetwork(bacnetip.WithICouldBeRouterToNetworkNetwork(net), bacnetip.WithICouldBeRouterToNetworkPerformanceIndex(perf))
	if err != nil {
		panic(err)
	}
	return network
}

func RejectMessageToNetwork(reason uint8, dnet uint16) *bacnetip.RejectMessageToNetwork {
	network, err := bacnetip.NewRejectMessageToNetwork(bacnetip.WithRejectMessageToNetworkRejectionReason(readWriteModel.NLMRejectMessageToNetworkRejectReason(reason)), bacnetip.WithRejectMessageToNetworkDnet(dnet))
	if err != nil {
		panic(err)
	}
	return network
}

func RouterBusyToNetwork(netList ...uint16) *bacnetip.RouterBusyToNetwork {
	network, err := bacnetip.NewRouterBusyToNetwork(bacnetip.WithRouterBusyToNetworkDnet(netList))
	if err != nil {
		panic(err)
	}
	return network
}

func RouterAvailableToNetwork(netList ...uint16) *bacnetip.RouterAvailableToNetwork {
	network, err := bacnetip.NewRouterAvailableToNetwork(bacnetip.WithRouterAvailableToNetworkDnet(netList))
	if err != nil {
		panic(err)
	}
	return network
}

func InitializeRoutingTable(irtTable ...*bacnetip.RoutingTableEntry) *bacnetip.InitializeRoutingTable {
	network, err := bacnetip.NewInitializeRoutingTable(bacnetip.WithInitializeRoutingTableIrtTable(irtTable...))
	if err != nil {
		panic(err)
	}
	return network
}

func RoutingTableEntry(address uint16, portId uint8, portInfo []byte) *bacnetip.RoutingTableEntry {
	return bacnetip.NewRoutingTableEntry(
		bacnetip.WithRoutingTableEntryDestinationNetworkAddress(address),
		bacnetip.WithRoutingTableEntryPortId(portId),
		bacnetip.WithRoutingTableEntryPortInfo(portInfo),
	)
}

func InitializeRoutingTableAck(irtaTable ...*bacnetip.RoutingTableEntry) *bacnetip.InitializeRoutingTableAck {
	network, err := bacnetip.NewInitializeRoutingTableAck(bacnetip.WithInitializeRoutingTableAckIrtaTable(irtaTable...))
	if err != nil {
		panic(err)
	}
	return network
}

func EstablishConnectionToNetwork(dnet uint16, terminationTime uint8) *bacnetip.EstablishConnectionToNetwork {
	network, err := bacnetip.NewEstablishConnectionToNetwork(bacnetip.WithEstablishConnectionToNetworkDNET(dnet), bacnetip.WithEstablishConnectionToNetworkTerminationTime(terminationTime))
	if err != nil {
		panic(err)
	}
	return network
}

func DisconnectConnectionToNetwork(dnet uint16) *bacnetip.DisconnectConnectionToNetwork {
	network, err := bacnetip.NewDisconnectConnectionToNetwork(bacnetip.WithDisconnectConnectionToNetworkDNET(dnet))
	if err != nil {
		panic(err)
	}
	return network
}

func WhatIsNetworkNumber(dnet uint16) *bacnetip.WhatIsNetworkNumber {
	network, err := bacnetip.NewWhatIsNetworkNumber()
	if err != nil {
		panic(err)
	}
	return network
}

func NetworkNumberIs(net uint16, flag bool) *bacnetip.NetworkNumberIs {
	network, err := bacnetip.NewNetworkNumberIs(bacnetip.WithNetworkNumberIsNET(net), bacnetip.WithNetworkNumberIsTerminationConfigured(flag))
	if err != nil {
		panic(err)
	}
	return network
}
