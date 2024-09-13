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

	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
)

type IPRouter struct {
	nodes []*IPRouterNode

	log zerolog.Logger
}

func NewIPRouter(localLog zerolog.Logger) *IPRouter {
	if _debug != nil {
		_debug("__init__")
	}
	return &IPRouter{
		log: localLog,
	}
}

func (n *IPRouter) AddNetwork(addr *Address, lan *IPNetwork) {
	if _debug != nil {
		_debug("add_network %r %r", addr, lan)
	}
	n.log.Debug().Stringer("addr", addr).Stringer("lan", lan).Msg("adding network")

	node, err := NewIPRouterNode(n.log, n, addr, lan)
	if err != nil {
		n.log.Error().Err(err).Msg("error creating IPRouterNode")
		return
	}
	if _debug != nil {
		_debug("    - node: %r", node)
	}
	n.log.Debug().Stringer("node", node).Msg("node")

	n.nodes = append(n.nodes, node)
}

func (n *IPRouter) ProcessPDU(node *IPRouterNode, pdu PDU) {
	if _debug != nil {
		_debug("process_pdu %r %r", node, pdu)
	}
	n.log.Debug().Stringer("node", node).Stringer("pdu", pdu).Msg("processing PDU")

	// unpack the address part of the destination
	addrstr := *pdu.GetPDUDestination().AddrIP //TODO: check if this is the right way here.
	ipaddr := addrstr
	n.log.Debug().Uint32("ipaddr", ipaddr).Msg("ipaddr")
	if _debug != nil {
		_debug("    - ipaddr: %r", ipaddr)
	}

	// loop through the other nodes
	for _, inode := range n.nodes {
		if inode != node {
			if ipaddr&*inode.addrMask == *inode.addrSubnet {
				if _debug != nil {
					_debug("    - inode: %r", inode)
				}
				n.log.Debug().Stringer("inode", inode).Msg("inode")
				if err := inode.ProcessPDU(pdu); err != nil {
					n.log.Debug().Err(err).Msg("error processing inode")
				}
			}
		}
	}
}

func (n *IPRouter) String() string {
	return fmt.Sprintf("IPRouter for %d nodes", len(n.nodes))
}
