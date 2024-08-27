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
	"slices"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
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

	// the PDU source is a tuple, convert it to an Address instance
	src, err := NewAddress(m.log, pdu.GetPDUSource())
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

func (b *AnnexJCodec) String() string {
	return fmt.Sprintf("AnnexJCodec(client: %s, server: %s)", b.Client, b.Server)
}

type BIPSAPRequirements interface {
	ServiceAccessPointRequirements
	_Client
}

type BIPSAP struct {
	*ServiceAccessPoint
	rootStruct BIPSAPRequirements

	// pass through args
	argSapID *int

	log zerolog.Logger
}

func NewBIPSAP(localLog zerolog.Logger, bipSapRequirements BIPSAPRequirements, opts ...func(*BIPSAP)) (*BIPSAP, error) {
	b := &BIPSAP{
		log: localLog,
	}
	for _, opt := range opts {
		opt(b)
	}
	localLog.Debug().
		Interface("sapID", b.argSapID).
		Interface("bipSapRequirements", bipSapRequirements).
		Msg("NewBIPSAP")
	serviceAccessPoint, err := NewServiceAccessPoint(localLog, bipSapRequirements, func(point *ServiceAccessPoint) {
		point.serviceID = b.argSapID
	})
	if err != nil {
		return nil, errors.Wrap(err, "Error creating service access point")
	}
	b.ServiceAccessPoint = serviceAccessPoint
	b.rootStruct = bipSapRequirements
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

	switch msg := pdu.GetRootMessage().(type) {
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
		xpdu, err := NewResult(WithResultBvlciResultCode(readWriteModel.BVLCResultCode_WRITE_BROADCAST_DISTRIBUTION_TABLE_NAK))
		if err != nil {
			return errors.Wrap(err, "error building result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NewArgs(xpdu), kwargs)
	case readWriteModel.BVLCReadBroadcastDistributionTableExactly:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(readWriteModel.BVLCResultCode_READ_BROADCAST_DISTRIBUTION_TABLE_NAK))
		if err != nil {
			return errors.Wrap(err, "error building result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NewArgs(xpdu), kwargs)
		// build a response
	case readWriteModel.BVLCRegisterForeignDeviceExactly:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(readWriteModel.BVLCResultCode_REGISTER_FOREIGN_DEVICE_NAK))
		if err != nil {
			return errors.Wrap(err, "error building result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NewArgs(xpdu), kwargs)
	case readWriteModel.BVLCReadForeignDeviceTableExactly:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(readWriteModel.BVLCResultCode_READ_FOREIGN_DEVICE_TABLE_NAK))
		if err != nil {
			return errors.Wrap(err, "error building result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NewArgs(xpdu), kwargs)
	case readWriteModel.BVLCDeleteForeignDeviceTableEntryExactly:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(readWriteModel.BVLCResultCode_DELETE_FOREIGN_DEVICE_TABLE_ENTRY_NAK))
		if err != nil {
			return errors.Wrap(err, "error building result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NewArgs(xpdu), kwargs)
	case readWriteModel.BVLCDistributeBroadcastToNetworkExactly:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(readWriteModel.BVLCResultCode_DISTRIBUTE_BROADCAST_TO_NETWORK_NAK))
		if err != nil {
			return errors.Wrap(err, "error building result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())

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

	// regular args
	argAddr *Address
	argTTL  *int

	// pass through args
	argSapID *int
	argCid   *int
	argSid   *int

	log zerolog.Logger
}

func NewBIPForeign(localLog zerolog.Logger, opts ...func(*BIPForeign)) (*BIPForeign, error) {
	b := &BIPForeign{
		log: localLog,
	}
	for _, opt := range opts {
		opt(b)
	}
	localLog.Debug().
		Stringer("addrs", b.argAddr).
		Interface("ttls", b.argTTL).
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
	b.registrationTimeoutTask = OneShotFunction(b.registrationExpired, NoArgs, NoKWArgs)

	// registration provided
	if b.argAddr != nil {
		// a little error checking
		if b.argTTL == nil {
			return nil, errors.New("BBMD address and time-to-live must both be specified")
		}

		if err := b.Register(b.argAddr, *b.argTTL); err != nil {
			return nil, errors.Wrap(err, "error registering")
		}
	}

	return b, nil
}

func WithBIPForeignAddress(addr *Address) func(*BIPForeign) {
	return func(b *BIPForeign) {
		b.argAddr = addr
	}
}

func WithBIPForeignTTL(ttl int) func(*BIPForeign) {
	return func(b *BIPForeign) {
		b.argTTL = &ttl
	}
}

func (b *BIPForeign) Indication(args Args, kwargs KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")
	pdu := args.Get0PDU()

	// check for local stations
	switch pdu.GetPDUDestination().AddrType {
	case LOCAL_STATION_ADDRESS:
		// make an original unicast _PDU
		xpdu, err := NewOriginalUnicastNPDU(pdu, WithOriginalUnicastNPDUDestination(pdu.GetPDUDestination()), WithOriginalUnicastNPDUUserData(pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error creating original unicast NPDU")
		}
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it downstream
		return b.Request(NewArgs(xpdu), NoKWArgs)
	case LOCAL_BROADCAST_ADDRESS:
		// check the BBMD registration status, we may not be registered
		if b.registrationStatus != 0 {
			b.log.Debug().Msg("packet dropped, unregistered")
			return nil
		}

		// make an original broadcast _PDU
		xpdu, err := NewOriginalBroadcastNPDU(pdu, WithOriginalBroadcastNPDUDestination(b.bbmdAddress), WithOriginalBroadcastNPDUUserData(pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error creating original unicast NPDU")
		}
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it downstream
		return b.Request(NewArgs(NewPDUFromPDUWithNewMessage(pdu, xpdu)), NoKWArgs)
	default:
		return errors.Errorf("invalid destination address: %s", pdu.GetPDUDestination())
	}
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
			b.startTrackRegistration()
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
		xpdu, err := NewResult(WithResultBvlciResultCode(readWriteModel.BVLCResultCode_WRITE_BROADCAST_DISTRIBUTION_TABLE_NAK))
		if err != nil {
			return errors.Wrap(err, "error building a result")
		}
		xpdu.SetPDUUserData(pdu.GetPDUUserData()) // TODO: upstream sets this in the constructor
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NewArgs(xpdu), NoKWArgs)
	case readWriteModel.BVLCReadBroadcastDistributionTableExactly:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(readWriteModel.BVLCResultCode_READ_BROADCAST_DISTRIBUTION_TABLE_NAK))
		if err != nil {
			return errors.Wrap(err, "error building a result")
		}
		xpdu.SetPDUUserData(pdu.GetPDUUserData()) // TODO: upstream sets this in the constructor
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NewArgs(xpdu), NoKWArgs)
	case readWriteModel.BVLCRegisterForeignDeviceExactly:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(readWriteModel.BVLCResultCode_REGISTER_FOREIGN_DEVICE_NAK))
		if err != nil {
			return errors.Wrap(err, "error building a result")
		}
		xpdu.SetPDUUserData(pdu.GetPDUUserData()) // TODO: upstream sets this in the constructor
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NewArgs(xpdu), NoKWArgs)
	case readWriteModel.BVLCReadForeignDeviceTableExactly:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(readWriteModel.BVLCResultCode_READ_FOREIGN_DEVICE_TABLE_NAK))
		if err != nil {
			return errors.Wrap(err, "error building a result")
		}
		xpdu.SetPDUUserData(pdu.GetPDUUserData()) // TODO: upstream sets this in the constructor
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NewArgs(xpdu), NoKWArgs)
	case readWriteModel.BVLCDeleteForeignDeviceTableEntryExactly:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(readWriteModel.BVLCResultCode_DELETE_FOREIGN_DEVICE_TABLE_ENTRY_NAK))
		if err != nil {
			return errors.Wrap(err, "error building a result")
		}
		xpdu.SetPDUUserData(pdu.GetPDUUserData()) // TODO: upstream sets this in the constructor
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NewArgs(xpdu), NoKWArgs)
	case readWriteModel.BVLCDistributeBroadcastToNetworkExactly:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(readWriteModel.BVLCResultCode_DISTRIBUTE_BROADCAST_TO_NETWORK_NAK))
		if err != nil {
			return errors.Wrap(err, "error building a result")
		}
		xpdu.SetPDUUserData(pdu.GetPDUUserData()) // TODO: upstream sets this in the constructor
		xpdu.SetPDUDestination(pdu.GetPDUSource())

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

// Register starts the foreign device registration process with the given BBMD.
//
//	Registration will be renewed periodically according to the ttl value
//	until explicitly stopped by a call to `unregister`.
func (b *BIPForeign) Register(addr *Address, ttl int) error {
	// a little error checking
	if ttl <= 0 {
		return errors.New("time-to-live must be greater than zero")
	}

	// save the BBMD address and time-to-live
	b.bbmdAddress = addr
	b.bbmdTimeToLive = &ttl

	// install this task to do registration renewal according to the TTL
	// and stop tracking any active registration timeouts
	b.InstallTask(WithInstallTaskOptionsWhen(time.Time{}))
	b.stopTrackRegistration()
	return nil
}

// Unregister stops the foreign device registration process.
//
// Immediately drops active foreign device registration and stops further
// registration renewals.
func (b *BIPForeign) Unregister() {
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
	b.stopTrackRegistration()
}

// ProcessTask is called when the registration request should be sent to the BBMD.
func (b *BIPForeign) ProcessTask() error {
	pdu, err := NewRegisterForeignDevice(WithRegisterForeignDeviceBvlciTimeToLive(uint16(*b.bbmdTimeToLive)))
	if err != nil {
		return errors.Wrap(err, "")
	}
	pdu.SetPDUDestination(b.bbmdAddress)

	// send it downstream
	if err := b.Request(NewArgs(pdu), NoKWArgs); err != nil {
		return errors.Wrap(err, "error sending request")
	}

	// schedule the next registration renewal
	b.InstallTask(WithInstallTaskOptionsDelta(time.Duration(*b.bbmdTimeToLive) * time.Second))
	return nil
}

// _start_track_registration From J.5.2.3 Foreign Device Table Operation (paraphrasing): if a
// foreign device does not renew its registration 30 seconds after its
// TTL expired then it will be removed from the BBMD's FDT.
//
// Thus, if we're registered and don't get a response to a subsequent
// renewal request 30 seconds after our TTL expired then we're
// definitely not registered anymore.
func (b *BIPForeign) startTrackRegistration() {
	b.registrationTimeoutTask.InstallTask(WithInstallTaskOptionsDelta(time.Duration(*b.bbmdTimeToLive)*time.Second + (30 * time.Second)))
}

func (b *BIPForeign) stopTrackRegistration() {
	b.registrationTimeoutTask.SuspendTask()
}

// _registration_expired is called when detecting that foreign device registration has definitely expired.
func (b *BIPForeign) registrationExpired(_ Args, _ KWArgs) error {
	b.registrationStatus = -1 // Unregistered
	b.stopTrackRegistration()
	return nil
}

func (b *BIPForeign) String() string {
	return fmt.Sprintf("BIPForeign(taskTime: %v, isScheduled: %t, registrationStatus: %d, bbmdAddress: %s, bbmdTimeToLive: %d)", b.taskTime, b.isScheduled, b.registrationStatus, b.bbmdAddress, b.bbmdTimeToLive)
}

type BIPBBMD struct {
	*BIPSAP
	*Client
	*Server
	*RecurringTask
	*DebugContents

	bbmdAddress *Address
	bbmdBDT     []*Address
	bbmdFDT     []*FDTEntry

	// Pass Through args
	argSapID *int
	argCID   *int
	argSID   *int

	log zerolog.Logger
}

func NewBIPBBMD(localLog zerolog.Logger, addr *Address) (*BIPBBMD, error) {
	b := &BIPBBMD{log: localLog}
	var err error
	b.BIPSAP, err = NewBIPSAP(localLog, b, func(bipsap *BIPSAP) {
		bipsap.argSapID = b.argSapID
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating BIPSAP")
	}
	b.Client, err = NewClient(localLog, b, func(client *Client) {
		client.clientID = b.argCID
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating Client")
	}
	b.Server, err = NewServer(localLog, b, func(Server *Server) {
		Server.serverID = b.argSID
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating Server")
	}
	b.RecurringTask = NewRecurringTask(localLog, b, WithRecurringTaskInterval(1*time.Second))

	b.bbmdAddress = addr

	// install so process_task runs
	b.InstallTask(WithInstallTaskOptionsNone())

	return b, nil
}

func (b *BIPBBMD) Indication(args Args, kwargs KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")

	pdu := args.Get0PDU()

	// check for local stations
	if pdu.GetPDUDestination().AddrType == LOCAL_STATION_ADDRESS {
		// make an original unicast PDU
		var err error
		xpdu, err := NewOriginalUnicastNPDU(pdu, WithOriginalUnicastNPDUDestination(pdu.GetPDUDestination()), WithOriginalUnicastNPDUUserData(pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error creating OriginalUnicastNPDU")
		}
		//           if settings.route_aware and pdu.pduDestination.addrRoute:
		//               xpdu.pduDestination = pdu.pduDestination.addrRoute
		b.log.Debug().Stringer("xpdu", xpdu).Msg("original unicast xpdu")

		// send it downstream
		return b.Request(NewArgs(xpdu), NoKWArgs)
	} else if pdu.GetPDUDestination().AddrType == LOCAL_BROADCAST_ADDRESS { // check for broadcaste
		// make an original unicast PDU
		var xpdu PDU
		var err error
		xpdu, err = NewOriginalBroadcastNPDU(pdu, WithOriginalBroadcastNPDUDestination(pdu.GetPDUDestination()), WithOriginalBroadcastNPDUUserData(pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error creating OriginalUnicastNPDU")
		}
		//           if settings.route_aware and pdu.pduDestination.addrRoute:
		//               xpdu.pduDestination = pdu.pduDestination.addrRoute
		b.log.Debug().Stringer("xpdu", xpdu).Msg("original broadcast xpdu")

		// send it downstream
		err = b.Request(NewArgs(xpdu), NoKWArgs)
		if err != nil {
			return errors.Wrap(err, "error sending request downstream")
		}

		// skip other processing if the route was provided
		//           if settings.route_aware and pdu.pduDestination.addrRoute:
		//               return

		// make a forwarded PDU
		xpdu, err = NewForwardedNPDU(pdu, WithForwardedNPDUAddress(b.bbmdAddress), WithForwardedNPDUUserData(pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error creating ForwardedNPDU")
		}
		b.log.Debug().Stringer("xpdu", xpdu).Msg("forwarded xpdu")

		// send it to the peers
		for _, bdte := range b.bbmdBDT {
			if !bdte.Equals(b.bbmdAddress) {
				dest, err := NewAddress(b.log, &AddressTuple[int, uint16]{int(*bdte.AddrIP | ^*bdte.AddrMask), *bdte.AddrPort})
				if err != nil {
					return errors.Wrap(err, "error creating address tuple")
				}
				xpdu.SetPDUDestination(dest)
				b.log.Debug().Stringer("pduDestination", xpdu.GetPDUDestination()).Msg("sending to peer")
				if err := b.Request(NewArgs(xpdu), NoKWArgs); err != nil {
					return errors.Wrap(err, "error sending request")
				}
			}
		}

		// send it to the registered foreign devies
		for _, fdte := range b.bbmdFDT {
			xpdu.SetPDUDestination(fdte.FDAddress)
			b.log.Debug().Stringer("pduDestination", xpdu.GetPDUDestination()).Msg("sending to foreign device")
			if err := b.Request(NewArgs(xpdu), NoKWArgs); err != nil {
				return errors.Wrap(err, "error sending request")
			}
		}
		return err
	} else {
		return errors.Errorf("invalid destination address: %s", pdu.GetPDUDestination())
	}
}

func (b *BIPBBMD) Confirmation(args Args, kwargs KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Confirmation")

	pdu := args.Get0PDU()
	switch pdu.GetRootMessage().(type) {
	case readWriteModel.BVLCResultExactly: //some kind of response to a request
		// send this to the service access point
		return b.SapResponse(NewArgs(pdu), NoKWArgs)
	case readWriteModel.BVLCWriteBroadcastDistributionTableExactly:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(readWriteModel.BVLCResultCode_WRITE_BROADCAST_DISTRIBUTION_TABLE_NAK))
		if err != nil {
			return errors.Wrap(err, "error creating Result")
		}
		xpdu.SetPDUUserData(pdu.GetPDUUserData())

		// send it downstream
		return b.Request(NewArgs(pdu), NoKWArgs)
	case readWriteModel.BVLCReadBroadcastDistributionTableExactly:
		// build a response
		xpdu, err := NewReadBroadcastDistributionTableAck(WithReadBroadcastDistributionTableAckBDT(b.bbmdBDT...))
		if err != nil {
			return errors.Wrap(err, "error creating Result")
		}
		xpdu.SetPDUUserData(pdu.GetPDUUserData())

		// send it downstream
		return b.Request(NewArgs(pdu), NoKWArgs)
	case readWriteModel.BVLCReadBroadcastDistributionTableAckExactly:
		// send it to the service access point
		return b.SapResponse(NewArgs(pdu), NoKWArgs)
	case readWriteModel.BVLCForwardedNPDUExactly:
		pdu := pdu.(*ForwardedNPDU) // TODO: check if this cast is fine
		// send it upstream if there is a network layer
		if b.serverPeer != nil {
			xpdu := NewPDU(NewMessageBridge(pdu.GetPduData()...), WithPDUSource(pdu.GetBvlciAddress()), WithPDUDestination(NewLocalBroadcast(nil)), WithPDUUserData(pdu.GetPDUUserData()))
			//               if settings.route_aware:
			//                   xpdu.pduSource.addrRoute = pdu.pduSource
			b.log.Debug().Stringer("xpdu", xpdu).Msg("upstream xpdu")

			return b.Response(NewArgs(xpdu), NoKWArgs)
		}

		// build a forwarded NPDU to send out
		xpdu, err := NewForwardedNPDU(pdu, WithForwardedNPDUAddress(pdu.GetBvlciAddress()), WithForwardedNPDUUserData(pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error creating ForwardedNPDU")
		}
		b.log.Debug().Stringer("xpdu", xpdu).Msg("forwarded xpdu")

		// if this was unicast to us, do next hop
		if pdu.GetPDUDestination().AddrType == LOCAL_STATION_ADDRESS {
			b.log.Trace().Msg("unicast message")

			// if this BBMD is listed in its BDT, send a local broadcast
			if slices.ContainsFunc(b.bbmdBDT, func(address *Address) bool {
				return address.Equals(b.bbmdAddress)
			}) {
				b.log.Trace().Msg("local broadcast")
				return b.Request(NewArgs(xpdu), NoKWArgs)
			}
		} else if pdu.GetPDUDestination().AddrType == LOCAL_BROADCAST_ADDRESS {
			b.log.Trace().Msg("directed broadcast message")
		} else {
			b.log.Warn().Stringer("destination", pdu.GetPDUDestination()).Msg("invalid destination address")
		}

		// send it to the registered foreign devices
		for _, fdte := range b.bbmdFDT {
			xpdu.SetPDUDestination(fdte.FDAddress)
			b.log.Warn().Stringer("destination", xpdu.GetPDUDestination()).Msg("sending to foreign device")
			if err := b.Request(NewArgs(xpdu), NoKWArgs); err != nil {
				return errors.Wrapf(err, "error sending request to destination %s", xpdu.GetPDUDestination())
			}
		}
		return nil
	case readWriteModel.BVLCRegisterForeignDeviceExactly:
		pdu := pdu.(*RegisterForeignDevice) // TODO: check if this cast is fine
		// process the request
		stat, err := b.RegisterForeignDevice(pdu.GetPDUSource(), pdu.GetBvlciTimeToLive())
		if err != nil {
			return errors.Wrap(err, "error registering device")
		}

		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(stat))
		if err != nil {
			return errors.Wrap(err, "error creating Result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())
		xpdu.SetPDUUserData(pdu.GetPDUUserData())
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it downstream
		return b.Request(NewArgs(xpdu), NoKWArgs)
	case readWriteModel.BVLCReadForeignDeviceTable:
		// build a response
		xpdu, err := NewReadForeignDeviceTableAck(WithReadForeignDeviceTableAckFDT(b.bbmdFDT...))
		if err != nil {
			return errors.Wrap(err, "error creating ack")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())
		xpdu.SetPDUUserData(pdu.GetPDUUserData())

		// send it downstream
		return b.Request(NewArgs(xpdu), NoKWArgs)
	case readWriteModel.BVLCReadForeignDeviceTableAckExactly:
		// send this to the service access point
		return b.SapResponse(NewArgs(pdu), NoKWArgs)
	case readWriteModel.BVLCDeleteForeignDeviceTableEntryExactly:
		pdu := pdu.(*DeleteForeignDeviceTableEntry) // TODO: check if this cast is fine
		// process the request
		stat, err := b.DeleteForeignDeviceTableEntry(pdu.GetBvlciAddress())
		if err != nil {
			return errors.Wrap(err, "error deleting entry")
		}

		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(stat))
		if err != nil {
			return errors.Wrap(err, "error creating Result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource()) // upstream does this in the constructor

		// send it downstream
		return b.Request(NewArgs(xpdu), NoKWArgs)
	case readWriteModel.BVLCDistributeBroadcastToNetworkExactly:
		// send it upstream if there is a network layer
		if b.serverPeer != nil {
			xpdu := NewPDU(NewMessageBridge(pdu.GetPduData()...), WithPDUSource(pdu.GetPDUSource()), WithPDUDestination(NewLocalBroadcast(nil)), WithPDUUserData(pdu.GetPDUUserData()))
			//               if settings.route_aware:
			//                   xpdu.pduSource.addrRoute = pdu.pduSource
			b.log.Debug().Stringer("xpdu", xpdu).Msg("upstream xpdu")

			return b.Response(NewArgs(xpdu), NoKWArgs)
		}

		// build a forwarded NPDU to send out
		xpdu, err := NewForwardedNPDU(pdu, WithForwardedNPDUAddress(pdu.GetPDUSource()), WithForwardedNPDUUserData(pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error creating ForwardedNPDU")
		}
		b.log.Debug().Stringer("xpdu", xpdu).Msg("forwarded xpdu")

		// send it to the peers
		for _, bdte := range b.bbmdBDT {
			if bdte.Equals(b.bbmdAddress) {
				xpdu.SetPDUDestination(NewLocalBroadcast(nil))
				b.log.Trace().Msg("local broadcast")
				err = b.Request(NewArgs(xpdu), NoKWArgs)
				if err != nil {
					return errors.Wrap(err, "error sending local broadcast")
				}
			} else {
				address, err := NewAddress(b.log, &AddressTuple[int, uint16]{int(*bdte.AddrIP | ^*bdte.AddrMask), *bdte.AddrPort})
				if err != nil {
					return errors.Wrap(err, "error creating address")
				}
				xpdu.SetPDUDestination(address)
				b.log.Debug().Stringer("designation", xpdu.GetPDUDestination()).Msg("sending to peer")
				err = b.Request(NewArgs(xpdu), NoKWArgs)
				if err != nil {
					return errors.Wrap(err, "error sending")
				}
			}
		}

		// send it to the registered foreign devices
		for _, fdte := range b.bbmdFDT {
			if !fdte.Equals(pdu.GetPDUSource()) {
				xpdu.SetPDUDestination(fdte.FDAddress)
				b.log.Warn().Stringer("destination", xpdu.GetPDUDestination()).Msg("sending to foreign device")
				if err := b.Request(NewArgs(xpdu), NoKWArgs); err != nil {
					return errors.Wrapf(err, "error sending request to destination %s", xpdu.GetPDUDestination())
				}
			}
		}
		return nil
	case readWriteModel.BVLCOriginalUnicastNPDUExactly:
		// send it upstream if there is a network layer
		if b.serverPeer != nil {
			// build a PDU
			xpdu := NewPDU(NewMessageBridge(pdu.GetPduData()...), WithPDUSource(pdu.GetPDUSource()), WithPDUDestination(pdu.GetPDUDestination()), WithPDUUserData(pdu.GetPDUUserData()))
			//               if settings.route_aware:
			//                   xpdu.pduSource.addrRoute = pdu.pduSource
			b.log.Debug().Stringer("xpdu", xpdu).Msg("upstream xpdu")

			return b.Response(NewArgs(xpdu), NoKWArgs)
		}
	case readWriteModel.BVLCOriginalBroadcastNPDUExactly:
		// send it upstream if there is a network layer
		if b.serverPeer != nil {
			// build a PDU with a local broadcast address
			xpdu := NewPDU(NewMessageBridge(pdu.GetPduData()...), WithPDUSource(pdu.GetPDUSource()), WithPDUDestination(NewLocalBroadcast(nil)), WithPDUUserData(pdu.GetPDUUserData()))
			//               if settings.route_aware:
			//                   xpdu.pduSource.addrRoute = pdu.pduSource
			b.log.Debug().Stringer("xpdu", xpdu).Msg("upstream xpdu")

			return b.Response(NewArgs(xpdu), NoKWArgs)
		}

		// build a forwarded NPDU
		xpdu, err := NewForwardedNPDU(pdu, WithForwardedNPDUAddress(pdu.GetPDUSource()), WithForwardedNPDUUserData(pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error creating ForwardedNPDU")
		}
		b.log.Debug().Stringer("xpdu", xpdu).Msg("forwarded xpdu")

		// send it to the peers
		for _, bdte := range b.bbmdBDT {
			if !bdte.Equals(b.bbmdAddress) {
				address, err := NewAddress(b.log, &AddressTuple[int, uint16]{int(*bdte.AddrIP | ^*bdte.AddrMask), *bdte.AddrPort})
				if err != nil {
					return errors.Wrap(err, "error creating address")
				}
				xpdu.SetPDUDestination(address)
				b.log.Debug().Stringer("designation", xpdu.GetPDUDestination()).Msg("sending to peer")
				err = b.Request(NewArgs(xpdu), NoKWArgs)
				if err != nil {
					return errors.Wrap(err, "error sending")
				}
			}
		}

		// send it to the registered foreign devices
		for _, fdte := range b.bbmdFDT {
			xpdu.SetPDUDestination(fdte.FDAddress)
			b.log.Warn().Stringer("destination", xpdu.GetPDUDestination()).Msg("sending to foreign device")
			if err := b.Request(NewArgs(xpdu), NoKWArgs); err != nil {
				return errors.Wrapf(err, "error sending request to destination %s", xpdu.GetPDUDestination())
			}
		}
		return nil
	default:
		return errors.Errorf("invalid pdu type: %s", pdu.GetRootMessage())
	}
	return nil
}

func (b *BIPBBMD) RegisterForeignDevice(address Arg, ttl uint16) (readWriteModel.BVLCResultCode, error) {
	b.log.Debug().Interface("address", address).Uint16("ttl", ttl).Msg("registering foreign device")

	var addr *Address
	// see if it is an address or make it one
	switch address := address.(type) {
	case *Address:
		addr = address
	case string:
		var err error
		addr, err = NewAddress(b.log, address)
		if err != nil {
			return readWriteModel.BVLCResultCode_REGISTER_FOREIGN_DEVICE_NAK, errors.Wrap(err, "error creating address")
		}
	default:
		return readWriteModel.BVLCResultCode_REGISTER_FOREIGN_DEVICE_NAK, errors.Errorf("invalid address type: %T", address)
	}

	var fdte *FDTEntry
	for _, fdtEntry := range b.bbmdFDT {
		if addr.Equals(fdtEntry.FDAddress) {
			fdte = fdtEntry
			break
		}
	}
	if fdte == nil {
		fdte = &FDTEntry{
			FDAddress: addr,
		}
		b.bbmdFDT = append(b.bbmdFDT, fdte)
	}

	fdte.FDTTL = ttl
	fdte.FDRemain = ttl + 5
	return readWriteModel.BVLCResultCode_SUCCESSFUL_COMPLETION, nil
}

func (b *BIPBBMD) DeleteForeignDeviceTableEntry(address Arg) (readWriteModel.BVLCResultCode, error) {
	b.log.Debug().Interface("address", address).Msg("delete foreign device")

	var addr *Address
	// see if it is an address or make it one
	switch address := address.(type) {
	case *Address:
		addr = address
	case string:
		var err error
		addr, err = NewAddress(b.log, address)
		if err != nil {
			return readWriteModel.BVLCResultCode_REGISTER_FOREIGN_DEVICE_NAK, errors.Wrap(err, "error creating address")
		}
	default:
		return readWriteModel.BVLCResultCode_REGISTER_FOREIGN_DEVICE_NAK, errors.Errorf("invalid address type: %T", address)
	}

	// find it and delete it
	stat := readWriteModel.BVLCResultCode_SUCCESSFUL_COMPLETION
	deleted := false
	b.bbmdFDT = slices.DeleteFunc(b.bbmdFDT, func(entry *FDTEntry) bool {
		if addr.Equals(entry.FDAddress) {
			deleted = true
			return true
		}
		return false
	})
	if !deleted {
		stat = readWriteModel.BVLCResultCode_DELETE_FOREIGN_DEVICE_TABLE_ENTRY_NAK
	}

	return stat, nil
}

func (b *BIPBBMD) ProcessTask() error {
	// look for foreign device registrations that have expired
	b.bbmdFDT = slices.DeleteFunc(b.bbmdFDT, func(fdte *FDTEntry) bool {
		fdte.FDRemain -= 1

		// delete it if expired
		if fdte.FDRemain <= 0 {
			b.log.Debug().Stringer("addr", fdte.FDAddress).Msg("foreign device expired")
			return true
		}
		return false
	})
	return nil
}

func (b *BIPBBMD) AddPeer(address Arg) error {
	b.log.Debug().Interface("adddress", address).Msg("addr")

	var addr *Address
	// see if it is an address or make it one
	switch address := address.(type) {
	case *Address:
		addr = address
	case string:
		var err error
		addr, err = NewAddress(b.log, address)
		if err != nil {
			return errors.Wrap(err, "error creating address")
		}
	default:
		return errors.Errorf("invalid address type: %T", address)
	}

	// see if it's already there
	found := false
	for _, bdte := range b.bbmdBDT {
		if addr.Equals(bdte) {
			found = true
			break
		}
	}
	if !found {
		b.bbmdBDT = append(b.bbmdBDT, addr)
	}
	return nil
}

func (b *BIPBBMD) DeletePeer(address Arg) error {
	b.log.Debug().Interface("adddress", address).Msg("addr")

	var addr *Address
	// see if it is an address or make it one
	switch address := address.(type) {
	case *Address:
		addr = address
	case string:
		var err error
		addr, err = NewAddress(b.log, address)
		if err != nil {
			return errors.Wrap(err, "error creating address")
		}
	default:
		return errors.Errorf("invalid address type: %T", address)
	}

	// look for the peer address
	b.bbmdBDT = slices.DeleteFunc(b.bbmdBDT, func(bdte *Address) bool {
		return addr.Equals(bdte)
	})
	return nil
}

func (b *BIPBBMD) String() string {
	return "BIPBBMD" // TODO: improve output
}
