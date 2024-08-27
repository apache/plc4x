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

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type IPRouterNode struct {
	*Client

	router     *IPRouter
	lan        *IPNetwork
	node       *IPNode
	addrMask   *uint32
	addrSubnet *uint32

	// pass through args
	argCid *int

	log zerolog.Logger
}

func NewIPRouterNode(localLog zerolog.Logger, router *IPRouter, addr *Address, lan *IPNetwork, opts ...func(*IPRouterNode)) (*IPRouterNode, error) {
	i := &IPRouterNode{
		// save the references to the router for packets and the lan for debugging
		router: router,
		lan:    lan,

		log: localLog,
	}
	for _, opt := range opts {
		opt(i)
	}
	var err error
	i.Client, err = NewClient(localLog, i, func(client *Client) {
		client.clientID = i.argCid
	})
	if err != nil {
		return nil, errors.Wrap(err, "error building client")
	}
	// make ourselves an IPNode and bind to it
	i.node, err = NewIPNode(localLog, addr, lan, WithNodePromiscuous(true), WithNodeSpoofing(true))
	if err != nil {
		return nil, errors.Wrap(err, "error building IPNode")
	}
	if err := Bind(localLog, i, i.node); err != nil {
		return nil, errors.Wrap(err, "error binding IPNode")
	}

	// save our mask and subnet
	i.addrMask = addr.AddrMask
	i.addrSubnet = addr.AddrSubnet
	return i, nil
}

func WithIPRouterNodeCid(cid int) func(*IPRouterNode) {
	return func(n *IPRouterNode) {
		n.argCid = &cid
	}
}

func (n *IPRouterNode) Confirmation(args Args, kwargs KWArgs) error {
	pdu := args.Get0PDU()
	n.log.Debug().Stringer("pdu", pdu).Msg("confirmation")
	n.router.ProcessPDU(n, pdu)
	return nil
}

func (n *IPRouterNode) ProcessPDU(pdu PDU) error {
	n.log.Debug().Stringer("pdu", pdu).Msg("ProcessPDU")
	return n.Request(NewArgs(pdu), NoKWArgs)
}

func (n *IPRouterNode) String() string {
	return fmt.Sprintf("IPRouterNode for %s", n.lan.name)
}
