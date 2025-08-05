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
	"fmt"
	"math/rand"

	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
)

type NetworkNode interface {
	fmt.Stringer
	setLan(lan *Network)
	getName() string
	setName(name string)
	getAddress() *Address
	isPromiscuous() bool
	Response(args Args, kwArgs KWArgs) error
}

type Network struct {
	*DefaultRFormatter

	name string

	nodes []NetworkNode

	broadcastAddress *Address
	dropPercent      float32

	trafficLogger TrafficLogger

	_leafName string

	log zerolog.Logger
}

func NewNetwork(localLog zerolog.Logger, options ...Option) *Network {
	n := &Network{
		DefaultRFormatter: NewDefaultRFormatter(),
		_leafName:         ExtractLeafName(options, StructName()),
		log:               localLog,
	}
	ApplyAppliers(options, n)
	if _debug != nil {
		_debug("__init__ name=%r broadcast_address=%r drop_percent=%r", n.name, n.broadcastAddress, n.dropPercent)
	}
	return n
}

func WithNetworkName(name string) GenericApplier[*Network] {
	return WrapGenericApplier(func(n *Network) { n.name = name })
}

func WithNetworkBroadcastAddress(broadcastAddress *Address) GenericApplier[*Network] {
	return WrapGenericApplier(func(n *Network) { n.broadcastAddress = broadcastAddress })
}

func WithNetworkDropPercent(dropPercent float32) GenericApplier[*Network] {
	return WrapGenericApplier(func(n *Network) { n.dropPercent = dropPercent })
}

func WithNetworkTrafficLogger(trafficLogger TrafficLogger) GenericApplier[*Network] {
	return WrapGenericApplier(func(n *Network) { n.trafficLogger = trafficLogger })
}

// AddNode Add a node to this network, let the node know which network it's on.
func (n *Network) AddNode(node NetworkNode) {
	if _debug != nil {
		_debug("add_node %r", node)
	}
	n.log.Debug().Stringer("node", node).Msg("Adding node")
	n.nodes = append(n.nodes, node)
	node.setLan(n)

	// update the node name
	if node.getName() == "" {
		node.setName(fmt.Sprintf("%s:%s", n.name, node.getAddress()))
	}
}

// RemoveNode Remove a node from this network.
func (n *Network) RemoveNode(node NetworkNode) {
	if _debug != nil {
		_debug("remove_node %r", node)
	}
	n.log.Debug().Stringer("node", node).Msg("Remove node")
	for i, _node := range n.nodes {
		if _node == node {
			n.nodes = append(n.nodes[:i], n.nodes[i+1:]...)
		}
	}
	node.setLan(nil)
}

// ProcessPDU Process a PDU by sending a copy to each node as dictated by the addressing and if a node is promiscuous.
func (n *Network) ProcessPDU(pdu PDU) error {
	if _debug != nil {
		_debug("process_pdu(%s) %r", n.name, pdu)
	}
	n.log.Debug().Stringer("pdu", pdu).Msg("processing pdu")

	// if there is a traffic log call it with the network name and PDU
	if tl := n.trafficLogger; tl != nil {
		tl.Call(NA(n.name, pdu))
	}

	// randomly drop a packet
	if n.dropPercent != 0.0 {
		if rand.Float32()*100 < n.dropPercent {
			if _debug != nil {
				_debug("    - packet dropped")
			}
			n.log.Trace().Msg("Dropping PDU")
			return nil
		}
	}

	if n.broadcastAddress != nil && pdu.GetPDUDestination().Equals(n.broadcastAddress) {
		if _debug != nil {
			_debug("    - broadcast")
		}
		n.log.Trace().Msg("broadcast")
		for _, node := range n.nodes {
			if !pdu.GetPDUSource().Equals(node.getAddress()) {
				if _debug != nil {
					_debug("    - match: %r", node)
				}
				n.log.Debug().Stringer("node", node).Msg("match")
				if err := node.Response(NA(DeepCopy[PDU](pdu)), NoKWArgs()); err != nil {
					n.log.Debug().Err(err).Msg("error processing PDU")
				}
			}
		}
	} else {
		if _debug != nil {
			_debug("    - unicast")
		}
		n.log.Debug().Msg("unicast")
		for _, node := range n.nodes {
			if node.isPromiscuous() || pdu.GetPDUDestination().Equals(node.getAddress()) {
				if _debug != nil {
					_debug("    - match: %r", node)
				}
				n.log.Debug().Stringer("node", node).Msg("match")
				if err := node.Response(NA(DeepCopy[PDU](pdu)), NoKWArgs()); err != nil {
					n.log.Debug().Err(err).Msg("error processing PDU")
				}
			}
		}
	}

	return nil
}

func (n *Network) String() string {
	return fmt.Sprintf("<%s name=%s>", n._leafName, n.name)
}
