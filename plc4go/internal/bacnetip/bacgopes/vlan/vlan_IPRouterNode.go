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

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
)

//go:generate go tool plc4xGenerator -type=IPRouterNode -prefix=vlan_
type IPRouterNode struct {
	ClientContract

	router     *IPRouter
	lan        *IPNetwork
	node       *IPNode
	addrMask   *uint32
	addrSubnet *uint32

	_leafName string

	log zerolog.Logger
}

func NewIPRouterNode(localLog zerolog.Logger, router *IPRouter, addr *Address, lan *IPNetwork, options ...Option) (*IPRouterNode, error) {
	i := &IPRouterNode{
		// save the references to the router for packets and the lan for debugging
		router: router,
		lan:    lan,

		_leafName: ExtractLeafName(options, StructName()),

		log: localLog,
	}
	ApplyAppliers(options, i)
	optionsForParent := AddLeafTypeIfAbundant(options, i)
	if _debug != nil {
		_debug("__init__ %r %r lan=%r", router, addr, lan)
	}
	var err error
	i.ClientContract, err = NewClient(localLog, optionsForParent...)
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

func (n *IPRouterNode) Confirmation(args Args, kwArgs KWArgs) error {
	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("confirmation %r", pdu)
	}
	n.log.Debug().Stringer("pdu", pdu).Msg("confirmation")
	n.router.ProcessPDU(n, pdu)
	return nil
}

func (n *IPRouterNode) ProcessPDU(pdu PDU) error {
	if _debug != nil {
		_debug("process_pdu %r", pdu)
	}
	n.log.Debug().Stringer("pdu", pdu).Msg("ProcessPDU")
	return n.Request(NA(pdu), NoKWArgs())
}

func (n *IPRouterNode) Format(s fmt.State, v rune) {
	switch v {
	case 's', 'v', 'r':
		_, _ = fmt.Fprintf(s, "<%s for %s>", n._leafName, n.lan.name)
	}
}
