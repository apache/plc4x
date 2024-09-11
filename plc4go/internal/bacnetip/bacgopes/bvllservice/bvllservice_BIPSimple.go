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

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/bvll"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

//go:generate plc4xGenerator -type=BIPSimple -prefix=bvllservice_
type BIPSimple struct {
	*BIPSAP
	ClientContract
	ServerContract

	// pass through args
	argSapID *int `ignore:"true"`
	argCid   *int `ignore:"true"`
	argSid   *int `ignore:"true"`

	log zerolog.Logger
}

func NewBIPSimple(localLog zerolog.Logger, opts ...func(simple *BIPSimple)) (*BIPSimple, error) {
	b := &BIPSimple{
		log: localLog,
	}
	for _, opt := range opts {
		opt(b)
	}
	if _debug != nil {
		_debug("__init__ sapID=%r cid=%r sid=%r", b.argSapID, b.argCid, b.argSid)
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
	b.ClientContract, err = NewClient(localLog, OptionalOption2(b.argCid, ToPtr[ClientRequirements](b), WithClientCID))
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	b.ServerContract, err = NewServer(localLog, OptionalOption2(b.argSid, ToPtr[ServerRequirements](b), WithServerSID))
	if err != nil {
		return nil, errors.Wrap(err, "error creating server")
	}
	return b, nil
}

func (b *BIPSimple) Indication(args Args, kwargs KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")
	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("indication %r", pdu)
	}
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
		if _debug != nil {
			_debug("    - xpdu: %r", xpdu)
		}
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it downstream
		return b.Request(NA(xpdu), NoKWArgs)
	case LOCAL_BROADCAST_ADDRESS:
		// make an original broadcast _PDU
		xpdu, err := NewOriginalBroadcastNPDU(pdu, WithOriginalBroadcastNPDUDestination(pdu.GetPDUDestination()), WithOriginalBroadcastNPDUUserData(pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error creating original BroadcastNPDU")
		}
		// TODO: route aware stuff missing here
		if _debug != nil {
			_debug("    - xpdu: %r", xpdu)
		}
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it downstream
		return b.Request(NA(xpdu), NoKWArgs)
	default:
		return errors.Errorf("invalid destination address: %s", pdu.GetPDUDestination())
	}
}

func (b *BIPSimple) Confirmation(args Args, kwargs KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Confirmation")
	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("confirmation %r", pdu)
	}

	switch msg := pdu.GetRootMessage().(type) {
	// some kind of response to a request
	case model.BVLCResult:
		// send this to the service access point
		return b.SapResponse(args, kwargs)
	case model.BVLCReadBroadcastDistributionTableAck:
		// send this to the service access point
		return b.SapResponse(args, kwargs)
	case model.BVLCReadForeignDeviceTableAck:
		// send this to the service access point
		return b.SapResponse(args, kwargs)
	case model.BVLCOriginalUnicastNPDU:
		// build a vanilla _PDU
		xpdu := NewPDU(NoArgs, NKW(KWCompRootMessage, msg.GetNpdu(), KWCPCISource, pdu.GetPDUSource(), KWCPCIDestination, pdu.GetPDUDestination()))
		if _debug != nil {
			_debug("    - xpdu: %r", xpdu)
		}
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it upstream
		return b.Response(NA(xpdu), kwargs)
	case model.BVLCOriginalBroadcastNPDU:
		// build a _PDU with a local broadcast address
		xpdu := NewPDU(NoArgs, NKW(KWCompRootMessage, msg.GetNpdu(), KWCPCISource, pdu.GetPDUSource(), KWCPCIDestination, NewLocalBroadcast(nil)))
		if _debug != nil {
			_debug("    - xpdu: %r", xpdu)
		}
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it upstream
		return b.Response(NA(xpdu), kwargs)
	case model.BVLCForwardedNPDU:
		// build a _PDU with the source from the real source
		ip := msg.GetIp()
		port := msg.GetPort()
		source, err := NewAddress(NA(append(ip, Uint16ToPort(port)...)))
		if err != nil {
			return errors.Wrap(err, "error building a ip")
		}
		xpdu := NewPDU(NoArgs, NKW(KWCompRootMessage, msg.GetNpdu(), KWCPCISource, source, KWCPCIDestination, NewLocalBroadcast(nil)))
		if _debug != nil {
			_debug("    - xpdu: %r", xpdu)
		}
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it upstream
		return b.Response(NA(xpdu), kwargs)
	case model.BVLCWriteBroadcastDistributionTable:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(model.BVLCResultCode_WRITE_BROADCAST_DISTRIBUTION_TABLE_NAK))
		if err != nil {
			return errors.Wrap(err, "error building result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NA(xpdu), kwargs)
	case model.BVLCReadBroadcastDistributionTable:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(model.BVLCResultCode_READ_BROADCAST_DISTRIBUTION_TABLE_NAK))
		if err != nil {
			return errors.Wrap(err, "error building result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NA(xpdu), kwargs)
		// build a response
	case model.BVLCRegisterForeignDevice:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(model.BVLCResultCode_REGISTER_FOREIGN_DEVICE_NAK))
		if err != nil {
			return errors.Wrap(err, "error building result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NA(xpdu), kwargs)
	case model.BVLCReadForeignDeviceTable:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(model.BVLCResultCode_READ_FOREIGN_DEVICE_TABLE_NAK))
		if err != nil {
			return errors.Wrap(err, "error building result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NA(xpdu), kwargs)
	case model.BVLCDeleteForeignDeviceTableEntry:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(model.BVLCResultCode_DELETE_FOREIGN_DEVICE_TABLE_ENTRY_NAK))
		if err != nil {
			return errors.Wrap(err, "error building result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NA(xpdu), kwargs)
	case model.BVLCDistributeBroadcastToNetwork:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(model.BVLCResultCode_DISTRIBUTE_BROADCAST_TO_NETWORK_NAK))
		if err != nil {
			return errors.Wrap(err, "error building result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NA(xpdu), kwargs)
	default:
		b.log.Warn().Type("msg", msg).Msg("invalid pdu type")
		return nil
	}
}
