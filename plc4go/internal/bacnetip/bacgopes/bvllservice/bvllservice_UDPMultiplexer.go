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

package bvllservice

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/udp"
)

//go:generate plc4xGenerator -type=UDPMultiplexer -prefix=bvllservice_
type UDPMultiplexer struct {
	address            *Address
	addrTuple          *AddressTuple[string, uint16]
	addrBroadcastTuple *AddressTuple[string, uint16]
	direct             *_MultiplexClient
	directPort         *UDPDirector
	broadcast          *_MultiplexClient
	broadcastPort      *UDPDirector
	annexH             *_MultiplexServer
	AnnexJ             *_MultiplexServer

	log zerolog.Logger
}

func NewUDPMultiplexer(localLog zerolog.Logger, address any, noBroadcast bool) (*UDPMultiplexer, error) {
	localLog.Debug().
		Interface("address", address).
		Bool("noBroadcast", noBroadcast).
		Msg("NewUDPMultiplexer")
	u := &UDPMultiplexer{}

	// check for some options
	specialBroadcast := false
	if address == nil {
		address, _ := NewAddress(NoArgs)
		u.address = address
		u.addrTuple = &AddressTuple[string, uint16]{"", 47808}
		u.addrBroadcastTuple = &AddressTuple[string, uint16]{"255.255.255.255", 47808}
	} else {
		// allow the address to be cast
		if caddress, ok := address.(*Address); ok {
			u.address = caddress
		} else if caddress, ok := address.(Address); ok {
			u.address = &caddress
		} else {
			newAddress, err := NewAddress(NA(address))
			if err != nil {
				return nil, errors.Wrap(err, "error parsing address")
			}
			u.address = newAddress
		}

		// promote the normal and broadcast tuples
		u.addrTuple = u.address.AddrTuple
		u.addrBroadcastTuple = u.address.AddrBroadcastTuple

		// check for no broadcasting (loopback interface)
		if u.addrBroadcastTuple == nil {
			noBroadcast = true
		} else if u.addrTuple == u.addrBroadcastTuple {
			// old school broadcast address
			u.addrBroadcastTuple = &AddressTuple[string, uint16]{"255.255.255.255", u.addrTuple.Right}
		} else {
			specialBroadcast = true
		}
	}

	localLog.Debug().
		Stringer("address", u.address).
		Stringer("addrTuple", u.addrTuple).
		Stringer("addrBroadcastTuple", u.addrBroadcastTuple).
		Bool("route_aware", Settings.RouteAware).
		Msg("working with")

	// create and bind direct address
	var err error
	u.direct, err = _New_MultiplexClient(localLog, u)
	if err != nil {
		return nil, errors.Wrap(err, "error creating multiplex client")
	}
	u.directPort, err = NewUDPDirector(localLog, *u.addrTuple)
	if err := Bind(localLog, u.direct, u.directPort); err != nil {
		return nil, errors.Wrap(err, "error binding ports")
	}

	// create and bind the broadcast address for non-Windows
	if specialBroadcast && !noBroadcast {
		u.broadcast, err = _New_MultiplexClient(localLog, u)
		if err != nil {
			return nil, errors.Wrap(err, "error creating broadcast multiplex client")
		}
		u.broadcastPort, err = NewUDPDirector(localLog, *u.addrBroadcastTuple, WithUDPDirectorReuse(true))
		if err := Bind(localLog, u.direct, u.directPort); err != nil {
			return nil, errors.Wrap(err, "error binding ports")
		}
	}

	// create and bind the Annex H and J servers
	u.annexH, err = _New_MultiplexServer(localLog, u)
	if err != nil {
		return nil, errors.Wrap(err, "error creating annexH")
	}
	u.AnnexJ, err = _New_MultiplexServer(localLog, u)
	if err != nil {
		return nil, errors.Wrap(err, "error creating annexJ")
	}
	return u, nil
}

func (m *UDPMultiplexer) Close() error {
	m.log.Debug().Msg("Close")

	// pass along the close to the director(s)
	if err := m.directPort.Close(); err != nil {
		m.log.Debug().Err(err).Msg("errored")
	}
	if m.broadcastPort != nil {
		if err := m.broadcastPort.Close(); err != nil {
			m.log.Debug().Err(err).Msg("errored")
		}
	}
	return nil
}

func (m *UDPMultiplexer) Indication(args Args, kwargs KWArgs) error {
	m.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")
	server := GA[*_MultiplexServer](args, 0)
	pdu := GA[PDU](args, 1)
	m.log.Debug().
		Stringer("server", server).
		Stringer("pdu", pdu).
		Msg("Indication")

	pduDestination := pdu.GetPDUDestination()

	// broadcast message
	var dest *Address
	if pduDestination.AddrType == LOCAL_BROADCAST_ADDRESS {
		// interface might not support broadcasts
		if m.addrBroadcastTuple == nil {
			return nil
		}

		address, err := NewAddress(NA(*m.addrBroadcastTuple))
		if err != nil {
			return errors.Wrap(err, "error getting address from tuple")
		}
		dest = address
		m.log.Debug().Stringer("dest", dest).Msg("requesting local broadcast")
	} else if pduDestination.AddrType == LOCAL_STATION_ADDRESS {
		dest = pduDestination
	} else {
		return errors.New("invalid destination address type")
	}

	return m.directPort.Indication(NA(NewPDU(NoArgs, NKW(KWCompRootMessage, pdu, KWCPCIDestination, dest))), NoKWArgs)
}

func (m *UDPMultiplexer) Confirmation(args Args, kwargs KWArgs) error {
	m.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Confirmation")
	client := GA[*_MultiplexClient](args, 0)
	pdu := GA[PDU](args, 1)
	m.log.Debug().
		Stringer("client", client).
		Stringer("pdu", pdu).
		Stringer("clientAddress", client.multiplexer.address).
		Msg("Confirmation")

	// if this came from ourselves, dump it
	pduSource := pdu.GetPDUSource()
	if pduSource.Equals(m.address) {
		m.log.Debug().Msg("from us")
		return nil
	}

	// the PDU source is a tuple, convert it to an Address instance
	src, err := NewAddress(NA(pdu.GetPDUSource()))
	if err != nil {
		return errors.Wrap(err, "error creating address")
	}
	var dest *Address

	// match the destination in case the stack needs it
	if client == m.direct {
		m.log.Debug().Msg("direct to us")
		dest = m.address
	} else if client == m.broadcast {
		m.log.Debug().Msg("broadcast to us")
		dest = NewLocalBroadcast(nil)
	} else {
		return errors.New("Confirmation missmatch")
	}
	m.log.Debug().Stringer("dest", dest).Msg("dest")

	// must have at least one octet
	if pdu.GetRootMessage() == nil {
		m.log.Debug().Msg("no data")
		return nil
	}

	// TODO: we only support 0x81 at the moment
	if m.AnnexJ != nil {
		return m.AnnexJ.Response(NA(NewPDU(NoArgs, NKW(KWCompRootMessage, pdu.GetRootMessage(), KWCPCISource, src, KWCPCIDestination, dest))), NoKWArgs)
	}

	return nil
}
