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
	"bufio"
	"fmt"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/transports/udp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net"
	"time"
)

type UDPActor struct {
	director *UDPDirector
	timeout  uint32
	timer    *OneShotFunctionTask
	peer     string
}

func NewUDPActor(director *UDPDirector, peer string) *UDPActor {
	a := &UDPActor{}

	// keep track of the director
	a.director = director

	// associated with a peer
	a.peer = peer

	// Add a timer
	a.timeout = director.timeout
	if a.timeout > 0 {
		a.timer = FunctionTask(a.idleTimeout)
		when := time.Now().Add(time.Duration(a.timeout) * time.Millisecond)
		a.timer.InstallTask(&when, nil)
	}

	// tell the director this is a new actor
	a.director.AddActor(a)
	return a
}

func (a *UDPActor) idleTimeout() error {
	log.Debug().Msg("idleTimeout")

	// tell the director this is gone
	a.director.DelActor(a)
	return nil
}

func (a *UDPActor) Indication(pdu _PDU) error {
	log.Debug().Msgf("Indication %s", pdu)

	// reschedule the timer
	if a.timer != nil {
		when := time.Now().Add(time.Duration(a.timeout) * time.Millisecond)
		a.timer.InstallTask(&when, nil)
	}

	// put it in the outbound queue for the director
	a.director.request <- pdu
	return nil
}

func (a *UDPActor) Response(pdu _PDU) error {
	log.Debug().Msgf("Response %s", pdu)

	// reschedule the timer
	if a.timer != nil {
		when := time.Now().Add(time.Duration(a.timeout) * time.Millisecond)
		a.timer.InstallTask(&when, nil)
	}

	// process this as a response from the director
	return a.director.Response(pdu)
}

func (a *UDPActor) HandleError(err error) {
	log.Debug().Err(err).Msg("HandleError")

	if err != nil {
		a.director.ActorError(err)
	}
}

type UDPDirector struct {
	*Server
	*ServiceAccessPoint

	timeout uint32
	reuse   bool
	address AddressTuple[string, uint16]
	ti      *udp.TransportInstance

	actorClass func(*UDPDirector, string) *UDPActor
	request    chan _PDU
	peers      map[string]*UDPActor
	running    bool
}

func NewUDPDirector(address AddressTuple[string, uint16], timeout *int, reuse *bool, sid *int, sapID *int) (*UDPDirector, error) {
	d := &UDPDirector{}
	var err error
	d.Server, err = NewServer(sid, d)
	if err != nil {
		return nil, errors.Wrap(err, "error creating server")
	}
	d.ServiceAccessPoint, err = NewServiceAccessPoint(sapID, d)
	if err != nil {
		return nil, errors.Wrap(err, "error creating service access point")
	}

	// check the actor class
	d.actorClass = NewUDPActor

	// save the timeout for actors
	if timeout != nil {
		d.timeout = uint32(*timeout)
	}

	if reuse != nil {
		d.reuse = *reuse
	}

	// save the address
	d.address = address

	// ask the dispatcher for a socket
	resolvedAddress, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", address.Left, address.Right))
	if err != nil {
		return nil, errors.Wrap(err, "error resolving udp address")
	}
	d.ti = udp.NewTransportInstance(resolvedAddress, nil, d.timeout, d.reuse, nil)
	if err := d.ti.Connect(); err != nil {
		return nil, errors.Wrap(err, "error connecting transport instance")
	}

	d.running = true
	go func() {
		for d.running {
			d.handleRead()
		}
	}()

	// create the request queue
	d.request = make(chan _PDU)
	go func() {
		// TODO: get requests and send them...
	}()

	// start with an empty peer pool
	d.peers = map[string]*UDPActor{}

	return d, nil
}

// AddActor adds an actor when a new one is connected
func (d *UDPDirector) AddActor(actor *UDPActor) {
	log.Debug().Msgf("AddActor %v", actor)

	d.peers[actor.peer] = actor

	// tell the ASE there is a new client
	if d.serviceElement != nil {
		// TODO: not sure how to realize that
		//d.SapRequest(actor)
	}
}

// DelActor removes an actor when the socket is closed.
func (d *UDPDirector) DelActor(actor *UDPActor) {
	log.Debug().Msgf("DelActor %v", actor)

	delete(d.peers, actor.peer)

	// tell the ASE the client has gone away
	if d.serviceElement != nil {
		// TODO: not sure how to realize that
		//d.SapRequest(actor)
	}
}

func (d *UDPDirector) GetActor(address Address) *UDPActor {
	return d.peers[address.String()]
}

func (d *UDPDirector) ActorError(err error) {
	// tell the ASE the actor had an error
	if d.serviceElement != nil {
		// TODO: not sure how to realize that
		//d.SapRequest(actor, err)
	}
}

func (d *UDPDirector) Close() error {
	d.running = false
	return d.ti.Close()
}

func (d *UDPDirector) handleRead() {
	log.Debug().Msgf("handleRead(%v)", d.address)

	if err := d.ti.FillBuffer(func(pos uint, _ byte, _ *bufio.Reader) bool {
		if pos >= 4 {
			return false
		}
		return true
	}); err != nil {
		// pass along to a handler
		d.handleError(errors.Wrap(err, "error filling buffer"))
		return
	}
	peekedBytes, err := d.ti.PeekReadableBytes(4)
	if err != nil {
		// pass along to a handler
		d.handleError(errors.Wrap(err, "error peeking 4 bytes"))
		return
	}

	length := uint32(peekedBytes[2])<<8 | uint32(peekedBytes[3])
	readBytes, err := d.ti.Read(length)

	bvlc, err := model.BVLCParse(readBytes)
	if err != nil {
		// pass along to a handler
		d.handleError(errors.Wrap(err, "error parsing bvlc"))
		return
	}

	// TODO: how to get the addr? Maybe we ditch the transport instance and use the udp socket directly
	pdu := NewPDU(bvlc)
	// send the PDU up to the client
	go d._response(pdu)
}

func (d *UDPDirector) handleError(err error) {
	log.Debug().Err(err).Msg("handleError")
}

// Indication Client requests are queued for delivery.
func (d *UDPDirector) Indication(pdu _PDU) error {
	log.Debug().Msgf("Indication %s", pdu)

	// get the destination
	addr := pdu.GetPDUDestination()

	// get the peer
	peer, ok := d.peers[addr.String()]
	if !ok {
		peer = d.actorClass(d, (*addr).String())
	}

	// send the message
	return peer.Indication(pdu)
}

// _response Incoming datagrams are routed through an actor.
func (d *UDPDirector) _response(pdu _PDU) error {
	log.Debug().Msgf("_response %s", pdu)

	// get the destination
	addr := pdu.GetPDUDestination()

	// get the peer
	peer, ok := d.peers[addr.String()]
	if !ok {
		peer = d.actorClass(d, (*addr).String())
	}

	// send the message
	return peer.Response(pdu)
}
