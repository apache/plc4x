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
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net"
)

type _MultiplexClient struct {
	*Client
	multiplexer *UDPMultiplexer
}

func _New_MultiplexClient(multiplexer *UDPMultiplexer) (*_MultiplexClient, error) {
	m := &_MultiplexClient{
		multiplexer: multiplexer,
	}
	var err error
	m.Client, err = NewClient(nil, m)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	return m, nil
}

func (m *_MultiplexClient) Confirmation(pdu spi.Message) error {
	return m.multiplexer.Confirmation(pdu)
}

type _MultiplexServer struct {
	*Server
	multiplexer *UDPMultiplexer
}

func _New_MultiplexServer(multiplexer *UDPMultiplexer) (*_MultiplexServer, error) {
	m := &_MultiplexServer{
		multiplexer: multiplexer,
	}
	var err error
	m.Server, err = NewServer(nil, m)
	if err != nil {
		return nil, errors.Wrap(err, "error creating server")
	}
	return m, nil
}

func (m *_MultiplexServer) Indication(pdu spi.Message) error {
	return m.multiplexer.Indication(pdu)
}

type UDPMultiplexer struct {
	address              *net.UDPAddr
	selfBroadcastAddress *net.UDPAddr
	direct               *_MultiplexClient
	directPort           *UDPDirector
	broadcast            *_MultiplexClient
	broadcastPort        *UDPDirector
	annexH               *_MultiplexServer
	annexJ               *_MultiplexServer
}

func NewUDPMultiplexer(address net.Addr, noBroadcast bool) (*UDPMultiplexer, error) {
	log.Debug().Msgf("NewUDPMultiplexer %v noBroadcast=%t", address, noBroadcast)
	u := &UDPMultiplexer{}

	// check for some options
	specialBroadcast := false
	if address == nil {
		u.address = &net.UDPAddr{IP: nil, Port: 46808}
		u.selfBroadcastAddress = &net.UDPAddr{IP: net.IPv4(255, 255, 255, 255), Port: 47808}
	} else {
		udpAddr, err := net.ResolveUDPAddr("udp", address.String())
		if err != nil {
			return nil, errors.Wrap(err, "error resolving upd")
		}
		u.address = udpAddr
		// TODO: we need to find a way to resolved broadcast
		u.selfBroadcastAddress = &net.UDPAddr{IP: net.IPv4(255, 255, 255, 255), Port: 47808}
	}

	log.Debug().Msgf("address %s, broadcast %s", u.address, u.selfBroadcastAddress)

	// create and bind direct address
	var err error
	u.direct, err = _New_MultiplexClient(u)
	if err != nil {
		return nil, errors.Wrap(err, "error creating multiplex client")
	}
	u.directPort, err = NewUDPDirector(u.address, nil, nil, nil, nil)
	if err := bind(u.direct, u.directPort); err != nil {
		return nil, errors.Wrap(err, "error binding ports")
	}

	// create and bind the broadcast address for non-Windows
	if specialBroadcast && !noBroadcast {
		u.broadcast, err = _New_MultiplexClient(u)
		if err != nil {
			return nil, errors.Wrap(err, "error creating broadcast multiplex client")
		}
		u.broadcastPort, err = NewUDPDirector(u.selfBroadcastAddress, nil, nil, nil, nil)
		if err := bind(u.direct, u.directPort); err != nil {
			return nil, errors.Wrap(err, "error binding ports")
		}
	}

	// create and bind the Annex H and J servers
	u.annexH, err = _New_MultiplexServer(u)
	if err != nil {
		return nil, errors.Wrap(err, "error creating annexH")
	}
	u.annexJ, err = _New_MultiplexServer(u)
	if err != nil {
		return nil, errors.Wrap(err, "error creating annexJ")
	}
	return u, nil
}

func (m *UDPMultiplexer) Close() error {
	panic("implement me")
}

func (m *UDPMultiplexer) Confirmation(pdu spi.Message) error {
	panic("implement me")
}

func (m *UDPMultiplexer) Indication(pdu spi.Message) error {
	panic("implement me")
}

type AnnexJCodec struct {
	*Client
	*Server
}

func NewAnnexJCodec(cid *int, sid *int) (*AnnexJCodec, error) {
	log.Debug().Msgf("NewAnnexJCodec cid=%d sid=%d", cid, sid)
	a := &AnnexJCodec{}
	client, err := NewClient(cid, a)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	a.Client = client
	server, err := NewServer(sid, a)
	if err != nil {
		return nil, errors.Wrap(err, "error creating server")
	}
	a.Server = server
	return a, nil
}

func (b *AnnexJCodec) Indication(apdu spi.Message) error {
	panic("we need to implement this with  generics as we handle npdu not apdu here")
}

func (b *AnnexJCodec) Confirmation(apdu spi.Message) error {
	panic("we need to implement this with  generics as we handle npdu not apdu here")
}

type _BIPSAP interface {
	_ServiceAccessPoint
	_Client
}

type BIPSAP struct {
	*ServiceAccessPoint
	rootStruct _BIPSAP
}

func NewBIPSAP(sapID *int, rootStruct _BIPSAP) (*BIPSAP, error) {
	log.Debug().Msgf("NewBIPSAP sapID=%d", sapID)
	b := &BIPSAP{}
	serviceAccessPoint, err := NewServiceAccessPoint(sapID, rootStruct)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating service access point")
	}
	b.ServiceAccessPoint = serviceAccessPoint
	b.rootStruct = rootStruct
	return b, nil
}

func (b *BIPSAP) SapIndication(pdu spi.Message) error {
	// TODO: extract from somewhere
	var pduDestination []byte
	panic("we need pduDestination")
	log.Debug().Msgf("SapIndication\n%s\n%s", pdu, pduDestination)
	// TODO: what to do with the destination?
	// this is a request initiated by the ASE, send this downstream
	return b.rootStruct.Request(pdu)
}

func (b *BIPSAP) SapConfirmation(pdu spi.Message) error {
	// TODO: extract from somewhere
	var pduDestination []byte
	panic("we need pduDestination")
	log.Debug().Msgf("SapConfirmation\n%s\n%s", pdu, pduDestination)
	// TODO: what to do with the destination?
	// this is a response from the ASE, send this downstream
	return b.rootStruct.Request(pdu)
}

type BIPSimple struct {
	*BIPSAP
	*Client
	*Server
}

func NewBIPSimple(sapID *int, cid *int, sid *int) (*BIPSimple, error) {
	log.Debug().Msgf("NewBIPSimple sapID=%d cid=%d sid=%d", sapID, cid, sid)
	b := &BIPSimple{}
	bipsap, err := NewBIPSAP(sapID, b)
	if err != nil {
		return nil, errors.Wrap(err, "error creating bisap")
	}
	b.BIPSAP = bipsap
	client, err := NewClient(cid, b)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	b.Client = client
	server, err := NewServer(sid, b)
	if err != nil {
		return nil, errors.Wrap(err, "error creating server")
	}
	b.Server = server
	return b, nil
}

func (b *BIPSimple) Indication(apdu spi.Message) error {
	panic("we need to implement this with  generics as we handle npdu not apdu here")
}

func (b *BIPSimple) Response(apdu spi.Message) error {
	panic("we need to implement this with  generics as we handle npdu not apdu here")
}
