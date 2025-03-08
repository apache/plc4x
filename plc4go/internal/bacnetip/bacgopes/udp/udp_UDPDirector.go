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

package udp

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/libp2p/go-reuseport"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/options"
)

//go:generate go tool plc4xGenerator -type=UDPDirector -prefix=udp_
type UDPDirector struct {
	ServerContract
	ServiceAccessPointContract

	timeout uint32
	reuse   bool
	address AddressTuple[string, uint16]
	udpConn *net.UDPConn `asPtr:"true"`

	actorClass func(zerolog.Logger, *UDPDirector, string) *UDPActor
	request    chan PDU
	peers      map[string]*UDPActor
	running    bool

	wg sync.WaitGroup

	passLogToModel bool
	log            zerolog.Logger
}

func NewUDPDirector(localLog zerolog.Logger, address AddressTuple[string, uint16], options ...Option) (*UDPDirector, error) {
	d := &UDPDirector{}
	ApplyAppliers(options, d)
	optionsForParent := AddLeafTypeIfAbundant(options, d)
	var err error
	d.ServerContract, err = NewServer(localLog, optionsForParent...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating server")
	}
	d.ServiceAccessPointContract, err = NewServiceAccessPoint(localLog, optionsForParent...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating service access point")
	}

	// check the actor class
	d.actorClass = NewUDPActor

	// save the address
	d.address = address

	// ask the dispatcher for a socket
	resolvedAddress, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", address.Left, address.Right))
	if err != nil {
		return nil, errors.Wrap(err, "error resolving udp address")
	}
	if d.reuse {
		if packetConn, err := reuseport.ListenPacket("udp", resolvedAddress.String()); err != nil {
			return nil, errors.Wrap(err, "error connecting to local address")
		} else {
			d.udpConn = packetConn.(*net.UDPConn)
		}
	} else {
		if d.udpConn, err = net.ListenUDP("udp", resolvedAddress); err != nil {
			return nil, errors.Wrap(err, "error connecting to local address")
		}
	}

	d.running = true
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		for d.running {
			d.handleRead()
		}
	}()

	// create the request queue
	d.request = make(chan PDU)
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		for d.running {
			pdu := <-d.request
			serialize, err := pdu.GetRootMessage().Serialize()
			if err != nil {
				localLog.Error().Err(err).Msg("Error building message")
				continue
			}
			// TODO: wonky address object
			destination := pdu.GetPDUDestination()
			addr := net.IPv4(destination.AddrAddress[0], destination.AddrAddress[1], destination.AddrAddress[2], destination.AddrAddress[3])
			udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", addr, *destination.AddrPort))
			if err != nil {
				localLog.Error().Err(err).Msg("Error resolving address")
				continue
			}
			writtenBytes, err := d.udpConn.WriteToUDP(serialize, udpAddr)
			if err != nil {
				localLog.Error().Err(err).Msg("Error writing bytes")
				continue
			}
			localLog.Debug().Int("writtenBytes", writtenBytes).Msg("written bytes")
		}
	}()

	// start with an empty peer pool
	d.peers = map[string]*UDPActor{}

	return d, nil
}

func WithUDPDirectorReuse(reuse bool) GenericApplier[*UDPDirector] {
	return WrapGenericApplier(func(d *UDPDirector) { d.reuse = reuse })
}

// AddActor adds an actor when a new one is connected
func (d *UDPDirector) AddActor(actor *UDPActor) {
	d.log.Debug().Stringer("actor", actor).Msg("AddActor %v")

	d.peers[actor.peer] = actor

	// tell the ASE there is a new client
	if d.GetServiceElement() != nil {
		if err := d.SapRequest(NoArgs, NKW(KWAddActor, actor)); err != nil {
			d.log.Error().Err(err).Msg("Error in add actor")
		}
	}
}

// DelActor removes an actor when the socket is closed.
func (d *UDPDirector) DelActor(actor *UDPActor) {
	d.log.Debug().Stringer("actor", actor).Msg("DelActor")

	delete(d.peers, actor.peer)

	// tell the ASE the client has gone away
	if d.GetServiceElement() != nil {
		if err := d.SapRequest(NoArgs, NKW(KWDelActor, actor)); err != nil {
			d.log.Error().Err(err).Msg("Error in del actor")
		}
	}
}

func (d *UDPDirector) ActorError(actor *UDPActor, err error) {
	// tell the ASE the actor had an error
	if d.GetServiceElement() != nil {
		if err := d.SapRequest(NoArgs, NKW(KWActorError, actor, KWError, err)); err != nil {
			d.log.Error().Err(err).Msg("Error in actor error")
		}
	}
}

func (d *UDPDirector) GetActor(address Address) *UDPActor {
	return d.peers[address.String()]
}

func (d *UDPDirector) handleConnect() {
	panic("implement me") // TODO: implement me
}

func (d *UDPDirector) readable() {
	panic("implement me") // TODO: implement me
}

func (d *UDPDirector) handleRead() {
	d.log.Debug().Stringer("address", &d.address).Msg("handleRead")

	readBytes := make([]byte, 1500) // TODO: check if that is sufficient
	var sourceAddr *net.UDPAddr
	if _, addr, err := d.udpConn.ReadFromUDP(readBytes); err != nil {
		d.log.Error().Err(err).Msg("error reading")
		return
	} else {
		sourceAddr = addr
	}

	ctxForModel := options.GetLoggerContextForModel(context.TODO(), d.log, options.WithPassLoggerToModel(d.passLogToModel))
	bvlc, err := model.BVLCParse[model.BVLC](ctxForModel, readBytes)
	if err != nil {
		// pass along to a handler
		d.handleError(errors.Wrap(err, "error parsing bvlc"))
		return
	}

	saddr, err := NewAddress(NA(sourceAddr))
	if err != nil {
		// pass along to a handler
		d.handleError(errors.Wrap(err, "error parsing source address"))
		return
	}
	daddr, err := NewAddress(NA(d.udpConn.LocalAddr()))
	if err != nil {
		// pass along to a handler
		d.handleError(errors.Wrap(err, "error parsing destination address"))
		return
	}
	pdu := NewCPDU(readBytes, NKW(KWCPCISource, saddr, KWCPCIDestination, daddr), WithRootMessage(bvlc)) // TODO: why do we set the destination here??? This might be completely wrong
	// send the _PDU up to the client
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		if err := d._response(pdu); err != nil {
			d.log.Debug().Err(err).Msg("errored")
		}
	}()
}

func (d *UDPDirector) writeable() {
	panic("implement me") // TODO: implement me
}

func (d *UDPDirector) handleWrite() {
	panic("implement me") // TODO: implement me
}

func (d *UDPDirector) Close() error {
	d.log.Debug().Msg("UDPDirector closing")
	defer func() {
		d.log.Debug().Msg("waiting for running tasks to finnish")
		d.wg.Wait()
		d.log.Debug().Msg("waiting done")
	}()
	d.running = false
	return d.udpConn.Close()
}

func (d *UDPDirector) handleClose() {
	panic("implement me") // TODO: implement me
}

func (d *UDPDirector) handleError(err error) {
	d.log.Debug().Err(err).Msg("handleError")
}

// Indication Client requests are queued for delivery.
func (d *UDPDirector) Indication(args Args, kwArgs KWArgs) error {
	d.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Indication")
	pdu := GA[PDU](args, 0)

	// get the destination
	addr := pdu.GetPDUDestination()

	// get the peer
	peer, ok := d.peers[addr.String()]
	if !ok {
		peer = d.actorClass(d.log, d, (*addr).String())
	}

	// send the message
	return peer.Indication(args, kwArgs)
}

// _response Incoming datagrams are routed through an actor.
func (d *UDPDirector) _response(pdu PDU) error {
	d.log.Debug().Stringer("pdu", pdu).Msg("_response")

	// get the destination
	addr := pdu.GetPDUSource()

	// get the peer
	peer, ok := d.peers[addr.String()]
	if !ok {
		peer = d.actorClass(d.log, d, addr.String())
	}

	// send the message
	return peer.Response(NA(pdu), NoKWArgs())
}
