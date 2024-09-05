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
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

func WhoIsRouterToNetwork(net uint16) *bacgopes.WhoIsRouterToNetwork {
	network, err := bacgopes.NewWhoIsRouterToNetwork(bacgopes.WithWhoIsRouterToNetworkNet(net))
	if err != nil {
		panic(err)
	}
	return network
}

func IAmRouterToNetwork(netList ...uint16) *bacgopes.IAmRouterToNetwork {
	network, err := bacgopes.NewIAmRouterToNetwork(bacgopes.WithIAmRouterToNetworkNetworkList(netList...))
	if err != nil {
		panic(err)
	}
	return network
}

func ICouldBeRouterToNetwork(net uint16, perf uint8) *bacgopes.ICouldBeRouterToNetwork {
	network, err := bacgopes.NewICouldBeRouterToNetwork(bacgopes.WithICouldBeRouterToNetworkNetwork(net), bacgopes.WithICouldBeRouterToNetworkPerformanceIndex(perf))
	if err != nil {
		panic(err)
	}
	return network
}

func RejectMessageToNetwork(reason uint8, dnet uint16) *bacgopes.RejectMessageToNetwork {
	network, err := bacgopes.NewRejectMessageToNetwork(bacgopes.WithRejectMessageToNetworkRejectionReason(readWriteModel.NLMRejectMessageToNetworkRejectReason(reason)), bacgopes.WithRejectMessageToNetworkDnet(dnet))
	if err != nil {
		panic(err)
	}
	return network
}

func RouterBusyToNetwork(netList ...uint16) *bacgopes.RouterBusyToNetwork {
	network, err := bacgopes.NewRouterBusyToNetwork(bacgopes.WithRouterBusyToNetworkDnet(netList))
	if err != nil {
		panic(err)
	}
	return network
}

func RouterAvailableToNetwork(netList ...uint16) *bacgopes.RouterAvailableToNetwork {
	network, err := bacgopes.NewRouterAvailableToNetwork(bacgopes.WithRouterAvailableToNetworkDnet(netList))
	if err != nil {
		panic(err)
	}
	return network
}

func InitializeRoutingTable(irtTable ...*bacgopes.RoutingTableEntry) *bacgopes.InitializeRoutingTable {
	network, err := bacgopes.NewInitializeRoutingTable(bacgopes.WithInitializeRoutingTableIrtTable(irtTable...))
	if err != nil {
		panic(err)
	}
	return network
}

func RoutingTableEntry(address uint16, portId uint8, portInfo []byte) *bacgopes.RoutingTableEntry {
	return bacgopes.NewRoutingTableEntry(
		bacgopes.WithRoutingTableEntryDestinationNetworkAddress(address),
		bacgopes.WithRoutingTableEntryPortId(portId),
		bacgopes.WithRoutingTableEntryPortInfo(portInfo),
	)
}

func InitializeRoutingTableAck(irtaTable ...*bacgopes.RoutingTableEntry) *bacgopes.InitializeRoutingTableAck {
	network, err := bacgopes.NewInitializeRoutingTableAck(bacgopes.WithInitializeRoutingTableAckIrtaTable(irtaTable...))
	if err != nil {
		panic(err)
	}
	return network
}

func EstablishConnectionToNetwork(dnet uint16, terminationTime uint8) *bacgopes.EstablishConnectionToNetwork {
	network, err := bacgopes.NewEstablishConnectionToNetwork(bacgopes.WithEstablishConnectionToNetworkDNET(dnet), bacgopes.WithEstablishConnectionToNetworkTerminationTime(terminationTime))
	if err != nil {
		panic(err)
	}
	return network
}

func DisconnectConnectionToNetwork(dnet uint16) *bacgopes.DisconnectConnectionToNetwork {
	network, err := bacgopes.NewDisconnectConnectionToNetwork(bacgopes.WithDisconnectConnectionToNetworkDNET(dnet))
	if err != nil {
		panic(err)
	}
	return network
}

func WhatIsNetworkNumber(dnet uint16) *bacgopes.WhatIsNetworkNumber {
	network, err := bacgopes.NewWhatIsNetworkNumber()
	if err != nil {
		panic(err)
	}
	return network
}

func NetworkNumberIs(net uint16, flag bool) *bacgopes.NetworkNumberIs {
	network, err := bacgopes.NewNetworkNumberIs(bacgopes.WithNetworkNumberIsNET(net), bacgopes.WithNetworkNumberIsTerminationConfigured(flag))
	if err != nil {
		panic(err)
	}
	return network
}
