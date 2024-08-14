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
	"math/rand"
)

type TrafficLogger interface {
	Call(args Args)
}

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
	name  string
	nodes []NetworkNode

	broadcastAddress *Address
	dropPercent      float32

	trafficLogger TrafficLogger

	log zerolog.Logger
}

func NewNetwork(name string, localLog zerolog.Logger, trafficLogger TrafficLogger, broadcastAddress *Address, dropPercent float32) *Network {
	return &Network{
		name:             name,
		broadcastAddress: broadcastAddress,
		dropPercent:      dropPercent,
		trafficLogger:    trafficLogger,
		log:              localLog,
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
				if err := node.Response(NewArgs(pdu.DeepCopy()), NoKWArgs); err != nil {
					n.log.Debug().Err(err).Msg("error processing PDU")
				}
			}
		}
	} else {
		n.log.Debug().Msg("unicast")
		for _, node := range n.nodes {
			if node.isPromiscuous() || pdu.GetPDUDestination().Equals(node.getAddress()) {
				n.log.Debug().Stringer("node", node).Msg("match")
				if err := node.Response(NewArgs(pdu.DeepCopy()), NoKWArgs); err != nil {
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

// NodeNetworkReference allows Network and IPNetwork to be used from Node.
type NodeNetworkReference interface {
	AddNode(node NetworkNode)
	ProcessPDU(pdu PDU) error
}

type Node struct {
	*Server

	lan     NodeNetworkReference
	address *Address
	name    string

	promiscuous bool
	spoofing    bool

	log zerolog.Logger
}

func NewNode(localLog zerolog.Logger, addr *Address, lan NodeNetworkReference, name string, promiscuous bool, spoofing bool, sid *int) (*Node, error) {
	n := &Node{
		lan:     nil,
		address: addr,
		name:    name,
		log:     localLog.With().Str("name", name).Logger(),
	}
	var err error
	n.Server, err = NewServer(localLog, sid, n)
	if err != nil {
		return nil, errors.Wrap(err, "error creating server")
	}

	// bind to a lan if it was provided
	if lan != nil {
		n.bind(lan)
	}

	// might receive all packets and might spoof
	n.promiscuous = promiscuous
	n.spoofing = spoofing

	return n, nil
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

func (n *Node) String() string {
	return fmt.Sprintf("Node: %s(%v)", n.name, n.serverID)
}

func (n *Node) bind(lan NodeNetworkReference) {
	n.log.Debug().Interface("lan", lan).Msg("binding lan")
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
	pdu := args.Get0PDU()
	if pduSource := pdu.GetPDUSource(); pduSource == nil {
		pdu.SetPDUSource(n.address)
	} else if !n.spoofing && !pduSource.Equals(n.address) {
		return errors.Errorf("spoofing address conflict (pduSource: '%s', nodeAddress: '%s').", pduSource, n.address)
	}

	// actual network delivery is a zero-delay task
	OneShotFunction(func(args Args, kwargs KWArgs) error {
		pdu := args.Get0PDU()
		return n.lan.ProcessPDU(pdu)
	}, args, NoKWArgs)
	return nil
}

// IPNetwork instances are Network objects where the addresses on the
//
//	network are tuples that would be used for sockets like ('1.2.3.4', 5).
//	The first node added to the network sets the broadcast address, like
//	('1.2.3.255', 5) and the other nodes must have the same tuple.
type IPNetwork struct {
	*Network
}

func NewIPNetwork(name string, localLog zerolog.Logger, trafficLogger TrafficLogger, broadcastAddress *Address, dropPercent float32) *IPNetwork {
	return &IPNetwork{
		Network: NewNetwork(name, localLog, trafficLogger, broadcastAddress, dropPercent),
	}
}

// AddNode Add a node to this network, let the node know which network it's on.
func (n *IPNetwork) AddNode(node NetworkNode) {
	n.log.Debug().Stringer("node", node).Msg("Adding node")

	ipNode := node.(*IPNode)

	address, err := NewAddress(n.log, ipNode.addrBroadcastTuple)
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

// An IPNode is a Node where the address is an Address that has an address
//
//	tuple and a broadcast tuple that would be used for socket communications.
type IPNode struct {
	*Node
	addrTuple          *AddressTuple[string, uint16]
	addrBroadcastTuple *AddressTuple[string, uint16]
}

func NewIPNode(localLog zerolog.Logger, addr *Address, lan *IPNetwork, promiscuous bool, spoofing bool, sid *int) (*IPNode, error) {
	i := &IPNode{
		// save the address information
		addrTuple:          addr.AddrTuple,
		addrBroadcastTuple: addr.AddrBroadcastTuple,
	}
	var err error
	i.Node, err = NewNode(localLog, addr, nil, "", promiscuous, spoofing, sid)
	if err != nil {
		return nil, errors.Wrap(err, "error creating node")
	}
	i.bind(lan) // TODO: we bind here otherwise we bind the contained node
	return i, nil
}

func (n *IPNode) bind(lan NodeNetworkReference) { // This is used to preserve the type
	n.log.Debug().Interface("lan", lan).Msg("binding lan")
	lan.AddNode(n)
}

func (n *IPNode) String() string {
	return fmt.Sprintf("IPNode(%v): %s, %v", n.Node, n.addrTuple, n.addrBroadcastTuple)
}

type IPRouterNode struct {
	*Client

	router     *IPRouter
	lan        *IPNetwork
	node       *IPNode
	addrMask   *uint32
	addrSubnet *uint32
}

func NewIPRouterNode(localLog zerolog.Logger, cid *int, router *IPRouter, addr *Address, lan *IPNetwork) (*IPRouterNode, error) {
	i := &IPRouterNode{
		// save the references to the router for packets and the lan for debugging
		router: router,
		lan:    lan,
	}
	var err error
	i.Client, err = NewClient(localLog, cid, i)
	if err != nil {
		return nil, errors.Wrap(err, "error building client")
	}
	// make ourselves an IPNode and bind to it
	i.node, err = NewIPNode(localLog, addr, lan, true, true, nil)
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

type IPRouter struct {
	nodes []*IPRouterNode

	log zerolog.Logger
}

func NewIPRouter(localLog zerolog.Logger) *IPRouter {
	return &IPRouter{
		log: localLog,
	}
}

func (n *IPRouter) AddNetwork(addr *Address, lan *IPNetwork) {
	n.log.Debug().Stringer("addr", addr).Stringer("lan", lan).Msg("adding network")

	node, err := NewIPRouterNode(n.log, nil, n, addr, lan)
	if err != nil {
		n.log.Error().Err(err).Msg("error creating IPRouterNode")
		return
	}
	n.log.Debug().Stringer("node", node).Msg("node")

	n.nodes = append(n.nodes, node)
}

func (n *IPRouter) ProcessPDU(node *IPRouterNode, pdu PDU) {
	n.log.Debug().Stringer("node", node).Stringer("pdu", pdu).Msg("processing PDU")

	// unpack the address part of the destination
	addrstr := *pdu.GetPDUDestination().AddrIP //TODO: check if this is the right way here.
	ipaddr := addrstr
	n.log.Debug().Uint32("ipaddr", ipaddr).Msg("ipaddr")

	// loop through the other nodes
	for _, inode := range n.nodes {
		if inode != node {
			if ipaddr&*inode.addrMask == *inode.addrSubnet {
				n.log.Debug().Stringer("inode", inode).Msg("inode")
				if err := inode.ProcessPDU(pdu); err != nil {
					n.log.Debug().Err(err).Msg("error processing inode")
				}
			}
		}
	}
}
