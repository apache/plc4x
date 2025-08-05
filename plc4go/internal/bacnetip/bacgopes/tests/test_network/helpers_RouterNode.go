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

package test_network

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/netservice"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/vlan"
)

type RouterNode struct {
	nsap *NetworkServiceAccessPoint
	nse  *_NetworkServiceElement

	log zerolog.Logger
}

func NewRouterNode(localLog zerolog.Logger) (*RouterNode, error) {
	r := &RouterNode{log: localLog}
	var err error
	// a network service access point will be needed
	r.nsap, err = NewNetworkServiceAccessPoint(r.log)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}
	// give the NSAP a generic network layer service element
	r.nse, err = new_NetworkServiceElement(r.log)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service element")
	}
	err = Bind(r.log, r.nse, r.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding")
	}
	if !LogTestNetwork {
		r.log = zerolog.Nop()
	}
	return r, nil
}

func (r *RouterNode) AddNetwork(address string, vlan *Network, net uint16) error {
	r.log.Debug().Str("address", address).Stringer("vlan", vlan).Uint16("net", net).Msg("AddNetwork")

	// convert the address to an Address
	addr, err := NewAddress(NA(address))
	if err != nil {
		return errors.Wrap(err, "error creaing address")
	}

	// create a node, add to the network
	node, err := NewNode(r.log, addr, WithNodeLan(vlan))
	if err != nil {
		return errors.Wrap(err, "error creating node")
	}

	// bind the BIP stack to the local network
	return r.nsap.Bind(node, &net, addr)
}

func (r *RouterNode) String() string {
	return fmt.Sprintf("RouterNode")
}
