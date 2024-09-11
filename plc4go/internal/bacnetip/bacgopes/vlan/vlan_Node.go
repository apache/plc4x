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
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/globals"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/task"
)

// NodeNetworkReference allows Network and IPNetwork to be used from Node.
type NodeNetworkReference interface {
	fmt.Stringer
	AddNode(node NetworkNode)
	ProcessPDU(pdu PDU) error
}

//go:generate plc4xGenerator -type=Node -prefix=vlan_
type Node struct {
	Server

	lan     NodeNetworkReference
	address *Address
	name    string

	promiscuous bool
	spoofing    bool

	// pass through args
	argSid *int `ignore:"true"`

	log zerolog.Logger
}

func NewNode(localLog zerolog.Logger, addr *Address, opts ...func(*Node)) (*Node, error) {
	n := &Node{
		address: addr,
		log:     localLog,
	}
	for _, opt := range opts {
		opt(n)
	}
	if n.name != "" {
		n.log = n.log.With().Str("name", n.name).Logger()
	}
	if LogVlan {
		n.log.Trace().
			Stringer("address", n.address).
			Stringer("lan", n.lan).
			Bool("promiscuous", n.promiscuous).
			Bool("spoofing", n.spoofing).
			Interface("sid", n.argSid).
			Msg("NewNode")
	}
	var err error
	n.Server, err = NewServer(localLog, n, OptionalOption(n.argSid, WithServerSID))
	if err != nil {
		return nil, errors.Wrap(err, "error creating server")
	}

	// bind to a lan if it was provided
	if n.lan != nil {
		n.bind(n.lan)
	}
	if !LogVlan {
		n.log = zerolog.Nop()
	}
	return n, nil
}

func WithNodeName(name string) func(*Node) {
	return func(n *Node) {
		n.name = name
	}
}

func WithNodeLan(lan NodeNetworkReference) func(*Node) {
	return func(n *Node) {
		n.lan = lan
	}
}

func WithNodePromiscuous(promiscuous bool) func(*Node) {
	return func(n *Node) {
		n.promiscuous = promiscuous
	}
}

func WithNodeSpoofing(spoofing bool) func(*Node) {
	return func(n *Node) {
		n.spoofing = spoofing
	}
}

func WithNodeSid(sid int) func(*Node) {
	return func(n *Node) {
		n.argSid = &sid
	}
}

func (n *Node) setLan(lan *Network) {
	n.lan = lan
}

func (n *Node) getName() string {
	return n.name
}

func (n *Node) setName(name string) {
	n.name = name
}

func (n *Node) getAddress() *Address {
	return n.address
}

func (n *Node) isPromiscuous() bool {
	return n.promiscuous
}

func (n *Node) SetPromiscuous(promiscuous bool) {
	n.promiscuous = promiscuous
}

func (n *Node) SetSpoofing(spoofing bool) {
	n.spoofing = spoofing
}

func (n *Node) bind(lan NodeNetworkReference) {
	n.log.Debug().Stringer("lan", lan).Msg("binding lan")
	lan.AddNode(n)
}

func (n *Node) Indication(args Args, kwargs KWArgs) error {
	n.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Indication")

	// Make sure we are connected
	if n.lan == nil {
		return errors.New("unbound node")
	}

	// if the pduSource is unset, fill in our address, otherwise
	// leave it alone to allow for simulated spoofing
	pdu := GA[PDU](args, 0)
	if pduSource := pdu.GetPDUSource(); pduSource == nil {
		pdu.SetPDUSource(n.address)
	} else if !n.spoofing && !pduSource.Equals(n.address) {
		return errors.Errorf("spoofing address conflict (pduSource: '%s', nodeAddress: '%s').", pduSource, n.address)
	}

	// actual network delivery is a zero-delay task
	OneShotFunction(func(args Args, kwargs KWArgs) error {
		return n.lan.ProcessPDU(pdu)
	}, args, NoKWArgs)
	return nil
}

func (n *Node) AlternateString() (string, bool) {
	if ExtendedGeneralOutput {
		return fmt.Sprintf("Node: %s(%v)", n.name, n.GetServerId()), true
	} else {
		return "", false
	}
}
