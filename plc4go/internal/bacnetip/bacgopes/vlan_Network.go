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
	"fmt"
	"math/rand"

	"github.com/rs/zerolog"
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
	name string

	nodes []NetworkNode

	broadcastAddress *Address
	dropPercent      float32

	trafficLogger TrafficLogger

	log zerolog.Logger
}

func NewNetwork(localLog zerolog.Logger, opts ...func(*Network)) *Network {
	network := &Network{
		log: localLog,
	}
	for _, opt := range opts {
		opt(network)
	}
	return network
}

func WithNetworkName(name string) func(*Network) {
	return func(n *Network) {
		n.name = name
	}
}

func WithNetworkBroadcastAddress(broadcastAddress *Address) func(*Network) {
	return func(n *Network) {
		n.broadcastAddress = broadcastAddress
	}
}

func WithNetworkDropPercent(dropPercent float32) func(*Network) {
	return func(n *Network) {
		n.dropPercent = dropPercent
	}
}

func WithNetworkTrafficLogger(trafficLogger TrafficLogger) func(*Network) {
	return func(n *Network) {
		n.trafficLogger = trafficLogger
	}
}

// AddNode Add a node to this network, let the node know which network it's on.
func (n *Network) AddNode(node NetworkNode) {
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
	n.log.Debug().Stringer("pdu", pdu).Msg("processing pdu")

	// if there is a traffic log call it with the network name and PDU
	if tl := n.trafficLogger; tl != nil {
		tl.Call(NewArgs(n.name, pdu))
	}

	// randomly drop a packet
	if n.dropPercent != 0.0 {
		if rand.Float32()*100 < n.dropPercent {
			n.log.Trace().Msg("Dropping PDU")
			return nil
		}
	}

	if n.broadcastAddress != nil && pdu.GetPDUDestination().Equals(n.broadcastAddress) {
		n.log.Trace().Msg("broadcast")
		for _, node := range n.nodes {
			if !pdu.GetPDUSource().Equals(node.getAddress()) {
				n.log.Debug().Stringer("node", node).Msg("match")
				if err := node.Response(NewArgs(DeepCopy[PDU](pdu)), NoKWArgs); err != nil {
					n.log.Debug().Err(err).Msg("error processing PDU")
				}
			}
		}
	} else {
		n.log.Debug().Msg("unicast")
		for _, node := range n.nodes {
			if node.isPromiscuous() || pdu.GetPDUDestination().Equals(node.getAddress()) {
				n.log.Debug().Stringer("node", node).Msg("match")
				if err := node.Response(NewArgs(DeepCopy[PDU](pdu)), NoKWArgs); err != nil {
					n.log.Debug().Err(err).Msg("error processing PDU")
				}
			}
		}
	}

	return nil
}

func (n *Network) String() string {
	return fmt.Sprintf("<Network name=%s>", n.name)
}
