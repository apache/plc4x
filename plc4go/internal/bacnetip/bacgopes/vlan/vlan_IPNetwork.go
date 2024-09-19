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

package vlan

import (
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
)

// IPNetwork instances are Network objects where the addresses on the
//
//	network are tuples that would be used for sockets like ('1.2.3.4', 5).
//	The first node added to the network sets the broadcast address, like
//	('1.2.3.255', 5) and the other nodes must have the same tuple.
type IPNetwork struct {
	*Network
}

func NewIPNetwork(localLog zerolog.Logger, options ...Option) *IPNetwork {
	if _debug != nil {
		_debug("__init__")
	}
	i := &IPNetwork{}
	ApplyAppliers(options, i)
	optionsForParent := AddLeafTypeIfAbundant(options, i)
	i.Network = NewNetwork(localLog, optionsForParent...)
	return i
}

// AddNode Add a node to this network, let the node know which network it's on.
func (n *IPNetwork) AddNode(node NetworkNode) {
	if _debug != nil {
		_debug("add_node %r", node)
	}
	n.log.Debug().Stringer("node", node).Msg("Adding node")

	ipNode := node.(*IPNode)

	address, err := NewAddress(NA(ipNode.addrBroadcastTuple))
	if err != nil {
		panic(err) // TODO: check that we do the right thing here. Originally the tuple gets assigned but that makes trouble downstream
	}

	// first node sets the broadcast tuple, other nodes much match
	if len(n.nodes) == 0 {
		n.broadcastAddress = address
	} else if !address.Equals(n.broadcastAddress) {
		panic("nodes must all have the same broadcast tuple")
	}

	// continue along
	n.Network.AddNode(node)
}
