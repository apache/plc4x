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
	"context"
	"fmt"
	"time"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type _MultiplexClient struct {
	*Client
	multiplexer *UDPMultiplexer
}

func _New_MultiplexClient(localLog zerolog.Logger, multiplexer *UDPMultiplexer) (*_MultiplexClient, error) {
	m := &_MultiplexClient{
		multiplexer: multiplexer,
	}
	var err error
	m.Client, err = NewClient(localLog, m)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	return m, nil
}

func (m *_MultiplexClient) Confirmation(args Args, kwargs KWArgs) error {
	return m.multiplexer.Confirmation(NewArgs(m, args), NoKWArgs)
}

type _MultiplexServer struct {
	*Server
	multiplexer *UDPMultiplexer
}

func _New_MultiplexServer(localLog zerolog.Logger, multiplexer *UDPMultiplexer) (*_MultiplexServer, error) {
	m := &_MultiplexServer{
		multiplexer: multiplexer,
	}
	var err error
	m.Server, err = NewServer(localLog, m)
	if err != nil {
		return nil, errors.Wrap(err, "error creating server")
	}
	return m, nil
}

func (m *_MultiplexServer) Indication(args Args, kwargs KWArgs) error {
	return m.multiplexer.Indication(NewArgs(m, args), NoKWArgs)
}

type UDPMultiplexer struct {
	address            *Address
	addrTuple          *AddressTuple[string, uint16]
	addrBroadcastTuple *AddressTuple[string, uint16]
	direct             *_MultiplexClient
	directPort         *UDPDirector
	broadcast          *_MultiplexClient
	broadcastPort      *UDPDirector
	annexH             *_MultiplexServer
	annexJ             *_MultiplexServer

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
		address, _ := NewAddress(localLog)
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
			newAddress, err := NewAddress(localLog, address)
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
	u.directPort, err = NewUDPDirector(localLog, *u.addrTuple, nil, nil, nil, nil)
	if err := Bind(localLog, u.direct, u.directPort); err != nil {
		return nil, errors.Wrap(err, "error binding ports")
	}

	// create and bind the broadcast address for non-Windows
	if specialBroadcast && !noBroadcast {
		u.broadcast, err = _New_MultiplexClient(localLog, u)
		if err != nil {
			return nil, errors.Wrap(err, "error creating broadcast multiplex client")
		}
		reuse := true
		u.broadcastPort, err = NewUDPDirector(localLog, *u.addrBroadcastTuple, nil, &reuse, nil, nil)
		if err := Bind(localLog, u.direct, u.directPort); err != nil {
			return nil, errors.Wrap(err, "error binding ports")
		}
	}

	// create and bind the Annex H and J servers
	u.annexH, err = _New_MultiplexServer(localLog, u)
	if err != nil {
		return nil, errors.Wrap(err, "error creating annexH")
	}
	u.annexJ, err = _New_MultiplexServer(localLog, u)
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
	server := args.Get0MultiplexServer()
	pdu := args.Get1PDU()
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

		address, err := NewAddress(m.log, *m.addrBroadcastTuple)
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

	return m.directPort.Indication(NewArgs(NewPDU(pdu, WithPDUDestination(dest))), NoKWArgs)
}

func (m *UDPMultiplexer) Confirmation(args Args, kwargs KWArgs) error {
	m.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Confirmation")
	client := args.Get0MultiplexClient()
	pdu := args.Get1PDU()
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

	// TODO: upstream this is a tuple but we don't have that here so we can work with what we have
	src := pduSource
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
	if m.annexJ != nil {
		return m.annexJ.Response(NewArgs(NewPDU(pdu.GetRootMessage(), WithPDUSource(src), WithPDUDestination(dest))), NoKWArgs)
	}

	return nil
}

type AnnexJCodec struct {
	*Client
	*Server

	// pass through args
	argCid *int
	argSid *int

	log zerolog.Logger
}

func NewAnnexJCodec(localLog zerolog.Logger, opts ...func(*AnnexJCodec)) (*AnnexJCodec, error) {
	a := &AnnexJCodec{
		log: localLog,
	}
	for _, opt := range opts {
		opt(a)
	}
	localLog.Debug().
		Interface("cid", a.argCid).
		Interface("sid", a.argSid).
		Msg("NewAnnexJCodec")
	client, err := NewClient(localLog, a, func(client *Client) {
		client.clientID = a.argSid
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	a.Client = client
	server, err := NewServer(localLog, a, func(server *Server) {
		server.serverID = a.argSid
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating server")
	}
	a.Server = server
	return a, nil
}

func WithAnnexJCodecCid(cid int) func(*AnnexJCodec) {
	return func(a *AnnexJCodec) {
		a.argCid = &cid
	}
}
func WithAnnexJCodecSid(sid int) func(*AnnexJCodec) {
	return func(a *AnnexJCodec) {
		a.argSid = &sid
	}
}

func (b *AnnexJCodec) String() string {
	return fmt.Sprintf("AnnexJCodec(client: %s, server: %s)", b.Client, b.Server)
}

func (b *AnnexJCodec) Indication(args Args, kwargs KWArgs) error {
	b.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Indication")

	rpdu := args.Get0PDU()

	// encode it as a generic BVLL PDU
	bvlpdu := NewBVLPDU(nil)
	if err := rpdu.(interface{ Encode(Arg) error }).Encode(bvlpdu); err != nil {
		return errors.Wrap(err, "error encoding PDU")
	}

	// encode it as a PDU
	pdu := NewPDU(nil)
	if err := bvlpdu.Encode(pdu); err != nil {
		return errors.Wrap(err, "error encoding PDU")
	}

	// send it downstream
	return b.Request(NewArgs(pdu), NoKWArgs)
}

func (b *AnnexJCodec) Confirmation(args Args, kwargs KWArgs) error {
	b.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Confirmation")

	pdu := args.Get0PDU()

	// interpret as a BVLL PDU
	bvlpdu := NewBVLPDU(nil)
	if err := bvlpdu.Decode(pdu); err != nil {
		return errors.Wrap(err, "error decoding pdu")
	}

	// get the class related to the function
	rpdu := BVLPDUTypes[bvlpdu.GetBvlcFunction()]()
	if err := rpdu.Decode(bvlpdu); err != nil {
		return errors.Wrap(err, "error decoding PDU")
	}

	// send it upstream
	return b.Response(NewArgs(rpdu), NoKWArgs)
}

type _BIPSAP interface {
	_ServiceAccessPoint
	_Client
}

type BIPSAP struct {
	*ServiceAccessPoint
	rootStruct _BIPSAP

	// pass through args
	argSapID *int

	log zerolog.Logger
}

func NewBIPSAP(localLog zerolog.Logger, rootStruct _BIPSAP, opts ...func(*BIPSAP)) (*BIPSAP, error) {
	b := &BIPSAP{
		log: localLog,
	}
	for _, opt := range opts {
		opt(b)
	}
	localLog.Debug().
		Interface("sapID", b.argSapID).
		Interface("rootStruct", rootStruct).
		Msg("NewBIPSAP")
	serviceAccessPoint, err := NewServiceAccessPoint(localLog, rootStruct, func(point *ServiceAccessPoint) {
		point.serviceID = b.argSapID
	})
	if err != nil {
		return nil, errors.Wrap(err, "Error creating service access point")
	}
	b.ServiceAccessPoint = serviceAccessPoint
	b.rootStruct = rootStruct
	return b, nil
}

func WithBIPSAPSapID(sapID int) func(*BIPSAP) {
	return func(point *BIPSAP) {
		point.argSapID = &sapID
	}
}

func (b *BIPSAP) String() string {
	return fmt.Sprintf("BIPSAP(SAP: %s)", b.ServiceAccessPoint)
}

func (b *BIPSAP) SapIndication(args Args, kwargs KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("SapIndication")
	// this is a request initiated by the ASE, send this downstream
	return b.rootStruct.Request(args, kwargs)
}

func (b *BIPSAP) SapConfirmation(args Args, kwargs KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("SapConfirmation")
	// this is a response from the ASE, send this downstream
	return b.rootStruct.Request(args, kwargs)
}

type BIPSimple struct {
	*BIPSAP
	*Client
	*Server

	// pass through args
	argSapID *int
	argCid   *int
	argSid   *int

	log zerolog.Logger
}

func NewBIPSimple(localLog zerolog.Logger, opts ...func(simple *BIPSimple)) (*BIPSimple, error) {
	b := &BIPSimple{
		log: localLog,
	}
	for _, opt := range opts {
		opt(b)
	}
	localLog.Debug().
		Interface("sapID", b.argSapID).
		Interface("cid", b.argCid).
		Interface("sid", b.argSid).
		Msg("NewBIPSimple")
	bipsap, err := NewBIPSAP(localLog, b, func(bipsap *BIPSAP) {
		bipsap.argSapID = b.argSapID
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating bisap")
	}
	b.BIPSAP = bipsap
	client, err := NewClient(localLog, b, func(client *Client) {
		client.clientID = b.argCid
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	b.Client = client
	server, err := NewServer(localLog, b, func(server *Server) {
		server.serverID = b.argSid
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating server")
	}
	b.Server = server
	return b, nil
}

func (b *BIPSimple) String() string {
	return fmt.Sprintf("BIPSimple(BIPSAP: %s, Client: %s, Server: %s)", b.BIPSAP, b.Client, b.Server)
}

func (b *BIPSimple) Indication(args Args, kwargs KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")
	pdu := args.Get0PDU()
	if pdu == nil {
		return errors.New("no pdu")
	}
	if pdu.GetPDUDestination() == nil {
		return errors.New("no pdu destination")
	}

	// check for local stations
	switch pdu.GetPDUDestination().AddrType {
	case LOCAL_STATION_ADDRESS:
		// make an original unicast _PDU
		xpdu, err := NewOriginalUnicastNPDU(pdu, WithOriginalUnicastNPDUDestination(pdu.GetPDUDestination()), WithOriginalUnicastNPDUUserData(pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error creating original unicastNPDU")
		}
		// TODO: route aware stuff missing here
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it downstream
		return b.Request(NewArgs(xpdu), NoKWArgs)
	case LOCAL_BROADCAST_ADDRESS:
		// make an original broadcast _PDU
		xpdu, err := NewOriginalBroadcastNPDU(pdu, WithOriginalBroadcastNPDUDestination(pdu.GetPDUDestination()), WithOriginalBroadcastNPDUUserData(pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error creating original BroadcastNPDU")
		}
		// TODO: route aware stuff missing here
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it downstream
		return b.Request(NewArgs(xpdu), NoKWArgs)
	default:
		return errors.Errorf("invalid destination address: %s", pdu.GetPDUDestination())
	}
}

func (b *BIPSimple) Confirmation(args Args, kwargs KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Confirmation")
	pdu := args.Get0PDU()

	// TODO: come up with a better way to check that... this is hugely inefficient
	_data := pdu.GetPDUUserData()
	_ = _data
	data := pdu.GetPduData()
	bvlcParse, err := readWriteModel.NPDUParse(context.Background(), data, uint16(len(data)))
	if err != nil {
		panic(err)
	}

	// TODO: we need to work with the inner types here....
	panic("todo")

	switch msg := bvlcParse.(type) {
	// some kind of response to a request
	case readWriteModel.BVLCResultExactly:
		// send this to the service access point
		return b.SapResponse(args, kwargs)
	case readWriteModel.BVLCReadBroadcastDistributionTableAckExactly:
		// send this to the service access point
		return b.SapResponse(args, kwargs)
	case readWriteModel.BVLCReadForeignDeviceTableAckExactly:
		// send this to the service access point
		return b.SapResponse(args, kwargs)
	case readWriteModel.BVLCOriginalUnicastNPDUExactly:
		// build a vanilla _PDU
		xpdu := NewPDU(msg.GetNpdu(), WithPDUSource(pdu.GetPDUSource()), WithPDUDestination(pdu.GetPDUDestination()))
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it upstream
		return b.Response(NewArgs(xpdu), kwargs)
	case readWriteModel.BVLCOriginalBroadcastNPDUExactly:
		// build a _PDU with a local broadcast address
		xpdu := NewPDU(msg.GetNpdu(), WithPDUSource(pdu.GetPDUSource()), WithPDUDestination(NewLocalBroadcast(nil)))
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it upstream
		return b.Response(NewArgs(xpdu), kwargs)
	case readWriteModel.BVLCForwardedNPDUExactly:
		// build a _PDU with the source from the real source
		ip := msg.GetIp()
		port := msg.GetPort()
		source, err := NewAddress(b.log, append(ip, uint16ToPort(port)...))
		if err != nil {
			return errors.Wrap(err, "error building a ip")
		}
		xpdu := NewPDU(msg.GetNpdu(), WithPDUSource(source), WithPDUDestination(NewLocalBroadcast(nil)))
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it upstream
		return b.Response(NewArgs(xpdu), kwargs)
	case readWriteModel.BVLCWriteBroadcastDistributionTableExactly:
		// build a response
		result := readWriteModel.NewBVLCResult(readWriteModel.BVLCResultCode_WRITE_BROADCAST_DISTRIBUTION_TABLE_NAK)
		xpdu := NewPDU(result, WithPDUDestination(pdu.GetPDUSource()))

		// send it downstream
		return b.Request(NewArgs(xpdu), kwargs)
	case readWriteModel.BVLCReadBroadcastDistributionTableExactly:
		// build a response
		result := readWriteModel.NewBVLCResult(readWriteModel.BVLCResultCode_READ_BROADCAST_DISTRIBUTION_TABLE_NAK)
		xpdu := NewPDU(result, WithPDUDestination(pdu.GetPDUSource()))

		// send it downstream
		return b.Request(NewArgs(xpdu), kwargs)
		// build a response
	case readWriteModel.BVLCRegisterForeignDeviceExactly:
		// build a response
		result := readWriteModel.NewBVLCResult(readWriteModel.BVLCResultCode_REGISTER_FOREIGN_DEVICE_NAK)
		xpdu := NewPDU(result, WithPDUDestination(pdu.GetPDUSource()))

		// send it downstream
		return b.Request(NewArgs(xpdu), kwargs)
	case readWriteModel.BVLCReadForeignDeviceTableExactly:
		// build a response
		result := readWriteModel.NewBVLCResult(readWriteModel.BVLCResultCode_READ_FOREIGN_DEVICE_TABLE_NAK)
		xpdu := NewPDU(result, WithPDUDestination(pdu.GetPDUSource()))

		// send it downstream
		return b.Request(NewArgs(xpdu), kwargs)
	case readWriteModel.BVLCDeleteForeignDeviceTableEntryExactly:
		// build a response
		result := readWriteModel.NewBVLCResult(readWriteModel.BVLCResultCode_DELETE_FOREIGN_DEVICE_TABLE_ENTRY_NAK)
		xpdu := NewPDU(result, WithPDUDestination(pdu.GetPDUSource()))

		// send it downstream
		return b.Request(NewArgs(xpdu), kwargs)
	case readWriteModel.BVLCDistributeBroadcastToNetworkExactly:
		// build a response
		result := readWriteModel.NewBVLCResult(readWriteModel.BVLCResultCode_DISTRIBUTE_BROADCAST_TO_NETWORK_NAK)
		xpdu := NewPDU(result, WithPDUDestination(pdu.GetPDUSource()))

		// send it downstream
		return b.Request(NewArgs(xpdu), kwargs)
	default:
		b.log.Warn().Type("msg", msg).Msg("invalid pdu type")
		return nil
	}
}

type BIPForeign struct {
	*BIPSAP
	*Client
	*Server
	*OneShotTask

	registrationStatus      int
	bbmdAddress             *Address
	bbmdTimeToLive          *int
	registrationTimeoutTask *OneShotFunctionTask

	// pass through args
	argSapID *int
	argCid   *int
	argSid   *int

	log zerolog.Logger
}

func NewBIPForeign(localLog zerolog.Logger, addr *Address, ttl *int, opts ...func(*BIPForeign)) (*BIPForeign, error) {
	b := &BIPForeign{
		log: localLog,
	}
	for _, opt := range opts {
		opt(b)
	}
	localLog.Debug().
		Stringer("addrs", addr).
		Interface("ttls", ttl).
		Interface("sapID", b.argSapID).
		Interface("cid", b.argCid).
		Interface("sid", b.argSid).
		Msg("NewBIPForeign")
	bipsap, err := NewBIPSAP(localLog, b, func(bipsap *BIPSAP) {
		bipsap.argSapID = b.argSapID
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating bisap")
	}
	b.BIPSAP = bipsap
	client, err := NewClient(localLog, b, func(client *Client) {
		client.clientID = b.argCid
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	b.Client = client
	server, err := NewServer(localLog, b, func(server *Server) {
		server.serverID = b.argSid
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating server")
	}
	b.Server = server
	b.OneShotTask = NewOneShotTask(b, nil)

	// -2=unregistered, -1=not attempted or no ack, 0=OK, >0 error
	b.registrationStatus = -1

	// clear the BBMD address and time-to-live
	b.bbmdAddress = nil
	b.bbmdTimeToLive = nil

	// used in tracking active registration timeouts
	b.registrationTimeoutTask = OneShotFunction(b._registration_expired, NoArgs, NoKWArgs)

	// registration provided
	if addr != nil {
		// a little error checking
		if ttl == nil {
			return nil, errors.New("BBMD address and time-to-live must both be specified")
		}

		if err := b.register(*addr, *ttl); err != nil {
			return nil, errors.Wrap(err, "error registering")
		}
	}

	return b, nil
}

func (b *BIPForeign) Indication(args Args, kwargs KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")
	pdu := args.Get0PDU()

	// check for local stations
	switch pdu.GetPDUDestination().AddrType {
	case LOCAL_STATION_ADDRESS:
		// make an original unicast _PDU
		xpdu := readWriteModel.NewBVLCOriginalUnicastNPDU(pdu.GetRootMessage().(readWriteModel.NPDU), 0)
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it downstream
		return b.Request(NewArgs(NewPDUFromPDUWithNewMessage(pdu, xpdu)), NoKWArgs)
	case LOCAL_BROADCAST_ADDRESS:
		// check the BBMD registration status, we may not be registered
		if b.registrationStatus != 0 {
			b.log.Debug().Msg("packet dropped, unregistered")
			return nil
		}

		// make an original broadcast _PDU
		xpdu := readWriteModel.NewBVLCOriginalBroadcastNPDU(pdu.GetRootMessage().(readWriteModel.NPDU), 0)

		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it downstream
		return b.Request(NewArgs(NewPDUFromPDUWithNewMessage(pdu, xpdu)), NoKWArgs)
	default:
		return errors.Errorf("invalid destination address: %s", pdu.GetPDUDestination())
	}
}

func (b *BIPForeign) String() string {
	return fmt.Sprintf("BIPForeign(TBD...)") // TODO: fill some info here
}

func (b *BIPForeign) Confirmation(args Args, kwargs KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Confirmation")
	pdu := args.Get0PDU()

	switch msg := pdu.GetRootMessage().(type) {
	// check for a registration request result
	case readWriteModel.BVLCResultExactly:
		// if we are unbinding, do nothing
		if b.registrationStatus == -2 {
			return nil
		}

		// make sure we have a bind request in process

		// make sure the result is from the bbmd

		if !pdu.GetPDUSource().Equals(b.bbmdAddress) {
			b.log.Debug().Msg("packet dropped, not from the BBMD")
			return nil
		}
		// save the result code as the status
		b.registrationStatus = int(msg.GetCode())

		// If successful, track registration timeout
		if b.registrationStatus == 0 {
			b._start_track_registration()
		}

		return nil
	case readWriteModel.BVLCOriginalUnicastNPDUExactly:
		// build a vanilla _PDU
		xpdu := NewPDU(msg.GetNpdu(), WithPDUSource(pdu.GetPDUSource()), WithPDUDestination(pdu.GetPDUDestination()))
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it upstream
		return b.Response(NewArgs(xpdu), NoKWArgs)
	case readWriteModel.BVLCForwardedNPDUExactly:
		// check the BBMD registration status, we may not be registered
		if b.registrationStatus != 0 {
			b.log.Debug().Msg("packet dropped, unregistered")
			return nil
		}

		// make sure the forwarded _PDU from the bbmd
		if !pdu.GetPDUSource().Equals(b.bbmdAddress) {
			b.log.Debug().Msg("packet dropped, not from the BBMD")
			return nil
		}

		// build a _PDU with the source from the real source
		ip := msg.GetIp()
		port := msg.GetPort()
		source, err := NewAddress(b.log, append(ip, uint16ToPort(port)...))
		if err != nil {
			return errors.Wrap(err, "error building a ip")
		}
		xpdu := NewPDU(msg.GetNpdu(), WithPDUSource(source), WithPDUDestination(NewLocalBroadcast(nil)))
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it upstream
		return b.Response(NewArgs(xpdu), NoKWArgs)
	case readWriteModel.BVLCReadBroadcastDistributionTableAckExactly:
		// send this to the service access point
		return b.SapResponse(args, NoKWArgs)
	case readWriteModel.BVLCReadForeignDeviceTableAckExactly:
		// send this to the service access point
		return b.SapResponse(args, NoKWArgs)
	case readWriteModel.BVLCWriteBroadcastDistributionTableExactly:
		// build a response
		result := readWriteModel.NewBVLCResult(readWriteModel.BVLCResultCode_WRITE_BROADCAST_DISTRIBUTION_TABLE_NAK)
		xpdu := NewPDU(result, WithPDUDestination(pdu.GetPDUSource()))

		// send it downstream
		return b.Request(NewArgs(xpdu), NoKWArgs)
	case readWriteModel.BVLCReadBroadcastDistributionTableExactly:
		// build a response
		result := readWriteModel.NewBVLCResult(readWriteModel.BVLCResultCode_READ_BROADCAST_DISTRIBUTION_TABLE_NAK)
		xpdu := NewPDU(result, WithPDUDestination(pdu.GetPDUSource()))

		// send it downstream
		return b.Request(NewArgs(xpdu), NoKWArgs)
	case readWriteModel.BVLCRegisterForeignDeviceExactly:
		// build a response
		result := readWriteModel.NewBVLCResult(readWriteModel.BVLCResultCode_REGISTER_FOREIGN_DEVICE_NAK)
		xpdu := NewPDU(result, WithPDUDestination(pdu.GetPDUSource()))

		// send it downstream
		return b.Request(NewArgs(xpdu), NoKWArgs)
	case readWriteModel.BVLCReadForeignDeviceTableExactly:
		// build a response
		result := readWriteModel.NewBVLCResult(readWriteModel.BVLCResultCode_READ_FOREIGN_DEVICE_TABLE_NAK)
		xpdu := NewPDU(result, WithPDUDestination(pdu.GetPDUSource()))

		// send it downstream
		return b.Request(NewArgs(xpdu), NoKWArgs)
	case readWriteModel.BVLCDeleteForeignDeviceTableEntryExactly:
		// build a response
		result := readWriteModel.NewBVLCResult(readWriteModel.BVLCResultCode_DELETE_FOREIGN_DEVICE_TABLE_ENTRY_NAK)
		xpdu := NewPDU(result, WithPDUDestination(pdu.GetPDUSource()))

		// send it downstream
		return b.Request(NewArgs(xpdu), NoKWArgs)
	case readWriteModel.BVLCDistributeBroadcastToNetworkExactly:
		// build a response
		result := readWriteModel.NewBVLCResult(readWriteModel.BVLCResultCode_DISTRIBUTE_BROADCAST_TO_NETWORK_NAK)
		xpdu := NewPDU(result, WithPDUDestination(pdu.GetPDUSource()))

		// send it downstream
		return b.Request(NewArgs(xpdu), NoKWArgs)
	case readWriteModel.BVLCOriginalBroadcastNPDUExactly:
		b.log.Debug().Msg("packet dropped")
		return nil
	default:
		b.log.Warn().Type("msg", msg).Msg("invalid pdu type")
		return nil
	}
}

// register starts the foreign device registration process with the given BBMD.
//
//	Registration will be renewed periodically according to the ttl value
//	until explicitly stopped by a call to `unregister`.
func (b *BIPForeign) register(addr Address, ttl int) error {
	// a little error checking
	if ttl <= 0 {
		return errors.New("time-to-live must be greater than zero")
	}

	// save the BBMD address and time-to-live
	b.bbmdAddress = &addr
	b.bbmdTimeToLive = &ttl

	// install this task to do registration renewal according to the TTL
	// and stop tracking any active registration timeouts
	var taskTime time.Time
	b.InstallTask(InstallTaskOptions{When: &taskTime})
	b._stop_track_registration()
	return nil
}

// unregister stops the foreign device registration process.
//
// Immediately drops active foreign device registration and stops further
// registration renewals.
func (b *BIPForeign) unregister() {
	pdu := NewPDU(readWriteModel.NewBVLCRegisterForeignDevice(0), WithPDUDestination(b.bbmdAddress))

	// send it downstream
	if err := b.Request(NewArgs(pdu), NoKWArgs); err != nil {
		b.log.Debug().Err(err).Msg("error sending request")
		return
	}

	// change the status to unregistered
	b.registrationStatus = -2

	// clear the BBMD address and time-to-live
	b.bbmdAddress = nil
	b.bbmdTimeToLive = nil

	// unschedule registration renewal & timeout tracking if previously
	// scheduled
	b.SuspendTask()
	b._stop_track_registration()
}

// ProcessTask is called when the registration request should be sent to the BBMD.
func (b *BIPForeign) ProcessTask() error {
	pdu := NewPDU(readWriteModel.NewBVLCRegisterForeignDevice(uint16(*b.bbmdTimeToLive)), WithPDUDestination(b.bbmdAddress))

	// send it downstream
	if err := b.Request(NewArgs(pdu), NoKWArgs); err != nil {
		return errors.Wrap(err, "error sending request")
	}

	// schedule the next registration renewal
	var delta = time.Duration(*b.bbmdTimeToLive) * time.Second
	b.InstallTask(InstallTaskOptions{Delta: &delta})
	return nil
}

// _start_track_registration From J.5.2.3 Foreign Device Table Operation (paraphrasing): if a
// foreign device does not renew its registration 30 seconds after its
// TTL expired then it will be removed from the BBMD's FDT.
//
// Thus, if we're registered and don't get a response to a subsequent
// renewal request 30 seconds after our TTL expired then we're
// definitely not registered anymore.
func (b *BIPForeign) _start_track_registration() {
	var delta = time.Duration(*b.bbmdTimeToLive)*time.Second + (30 * time.Second)
	b.registrationTimeoutTask.InstallTask(InstallTaskOptions{Delta: &delta})
}

func (b *BIPForeign) _stop_track_registration() {
	b.registrationTimeoutTask.SuspendTask()
}

// _registration_expired is called when detecting that foreign device registration has definitely expired.
func (b *BIPForeign) _registration_expired(_ Args, _ KWArgs) error {
	b.registrationStatus = -1 // Unregistered
	b._stop_track_registration()
	return nil
}
