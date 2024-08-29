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

	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type BIPSimple struct {
	*BIPSAP
	Client
	Server

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
	client, err := NewClient(localLog, b, func(client *client) {
		client.clientID = b.argCid
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	b.Client = client
	server, err := NewServer(localLog, b, func(server *server) {
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
	case model.BVLCResultExactly:
		// send this to the service access point
		return b.SapResponse(args, kwargs)
	case model.BVLCReadBroadcastDistributionTableAckExactly:
		// send this to the service access point
		return b.SapResponse(args, kwargs)
	case model.BVLCReadForeignDeviceTableAckExactly:
		// send this to the service access point
		return b.SapResponse(args, kwargs)
	case model.BVLCOriginalUnicastNPDUExactly:
		// build a vanilla _PDU
		xpdu := NewPDU(msg.GetNpdu(), WithPDUSource(pdu.GetPDUSource()), WithPDUDestination(pdu.GetPDUDestination()))
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it upstream
		return b.Response(NewArgs(xpdu), kwargs)
	case model.BVLCOriginalBroadcastNPDUExactly:
		// build a _PDU with a local broadcast address
		xpdu := NewPDU(msg.GetNpdu(), WithPDUSource(pdu.GetPDUSource()), WithPDUDestination(NewLocalBroadcast(nil)))
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it upstream
		return b.Response(NewArgs(xpdu), kwargs)
	case model.BVLCForwardedNPDUExactly:
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
	case model.BVLCWriteBroadcastDistributionTableExactly:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(model.BVLCResultCode_WRITE_BROADCAST_DISTRIBUTION_TABLE_NAK))
		if err != nil {
			return errors.Wrap(err, "error building result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NewArgs(xpdu), kwargs)
	case model.BVLCReadBroadcastDistributionTableExactly:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(model.BVLCResultCode_READ_BROADCAST_DISTRIBUTION_TABLE_NAK))
		if err != nil {
			return errors.Wrap(err, "error building result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NewArgs(xpdu), kwargs)
		// build a response
	case model.BVLCRegisterForeignDeviceExactly:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(model.BVLCResultCode_REGISTER_FOREIGN_DEVICE_NAK))
		if err != nil {
			return errors.Wrap(err, "error building result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NewArgs(xpdu), kwargs)
	case model.BVLCReadForeignDeviceTableExactly:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(model.BVLCResultCode_READ_FOREIGN_DEVICE_TABLE_NAK))
		if err != nil {
			return errors.Wrap(err, "error building result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NewArgs(xpdu), kwargs)
	case model.BVLCDeleteForeignDeviceTableEntryExactly:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(model.BVLCResultCode_DELETE_FOREIGN_DEVICE_TABLE_ENTRY_NAK))
		if err != nil {
			return errors.Wrap(err, "error building result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NewArgs(xpdu), kwargs)
	case model.BVLCDistributeBroadcastToNetworkExactly:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(model.BVLCResultCode_DISTRIBUTE_BROADCAST_TO_NETWORK_NAK))
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
