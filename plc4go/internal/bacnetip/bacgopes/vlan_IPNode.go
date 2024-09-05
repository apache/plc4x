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

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// An IPNode is a Node where the address is an Address that has an address
//
//	tuple and a broadcast tuple that would be used for socket communications.
type IPNode struct {
	*Node
	addrTuple          *AddressTuple[string, uint16]
	addrBroadcastTuple *AddressTuple[string, uint16]
}

func NewIPNode(localLog zerolog.Logger, addr *Address, lan *IPNetwork, opts ...func(*Node)) (*IPNode, error) {
	i := &IPNode{
		// save the address information
		addrTuple:          addr.AddrTuple,
		addrBroadcastTuple: addr.AddrBroadcastTuple,
	}
	var err error
	i.Node, err = NewNode(localLog, addr, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating node")
	}
	i.bind(lan) // bind here otherwise we bind the contained node
	return i, nil
}

func (n *IPNode) bind(lan NodeNetworkReference) { // This is used to preserve the type
	n.log.Debug().Interface("lan", lan).Msg("binding lan")
	lan.AddNode(n)
}

func (n *IPNode) String() string {
	return fmt.Sprintf("IPNode(%v): %s, %v", n.Node, n.addrTuple, n.addrBroadcastTuple)
}
