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
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

//go:generate go tool plc4xGenerator -type=BIPSimple -prefix=bvllservice_
type BIPSimple struct {
	*BIPSAP
	ClientContract
	ServerContract
	*debugging.DefaultRFormatter `ignore:"true"`

	log zerolog.Logger
}

func NewBIPSimple(localLog zerolog.Logger, options ...Option) (*BIPSimple, error) {
	b := &BIPSimple{
		DefaultRFormatter: debugging.NewDefaultRFormatter(),
		log:               localLog,
	}
	ApplyAppliers(options, b)
	optionsForParent := AddLeafTypeIfAbundant(options, b)
	bipsap, err := NewBIPSAP(localLog, b, optionsForParent...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating bisap")
	}
	b.BIPSAP = bipsap
	b.ClientContract, err = NewClient(localLog, optionsForParent...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	b.ServerContract, err = NewServer(localLog, optionsForParent...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating server")
	}
	if _debug != nil {
		_debug("__init__ sapID=%r cid=%r sid=%r", b.GetServiceID(), b.GetClientID(), b.GetServerId())
	}
	localLog.Debug().
		Interface("sapID", b.GetServerId()).
		Interface("cid", b.GetClientID()).
		Interface("sid", b.GetServiceID()).
		Msg("NewBIPSimple")
	return b, nil
}

func (b *BIPSimple) Indication(args Args, kwArgs KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Indication")
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
		xpdu, err := NewOriginalUnicastNPDU(NA(pdu), NKW(KWCPCIDestination, pdu.GetPDUDestination(), KWCPCIUserData, pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error creating original unicastNPDU")
		}
		// TODO: route aware stuff missing here
		if _debug != nil {
			_debug("    - xpdu: %r", xpdu)
		}
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it downstream
		return b.Request(NA(xpdu), NoKWArgs())
	case LOCAL_BROADCAST_ADDRESS:
		// make an original broadcast _PDU
		xpdu, err := NewOriginalBroadcastNPDU(NA(pdu), NKW(KWCPCIDestination, pdu.GetPDUDestination(), KWCPCIUserData, pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error creating original BroadcastNPDU")
		}
		// TODO: route aware stuff missing here
		if _debug != nil {
			_debug("    - xpdu: %r", xpdu)
		}
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it downstream
		return b.Request(NA(xpdu), NoKWArgs())
	default:
		return errors.Errorf("invalid destination address: %s", pdu.GetPDUDestination())
	}
}

func (b *BIPSimple) Confirmation(args Args, kwArgs KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Confirmation")
	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("confirmation %r", pdu)
	}

	switch pdu := pdu.(type) {
	// some kind of response to a request
	case *Result:
		// send this to the service access point
		return b.SapResponse(args, kwArgs)
	case *ReadBroadcastDistributionTableAck:
		// send this to the service access point
		return b.SapResponse(args, kwArgs)
	case *ReadForeignDeviceTableAck:
		// send this to the service access point
		return b.SapResponse(args, kwArgs)
	case *OriginalUnicastNPDU:
		// build a vanilla _PDU
		xpdu := NewPDU(NA(pdu.GetPduData()), NKW(KWCPCISource, pdu.GetPDUSource(), KWCPCIDestination, pdu.GetPDUDestination(), KWCPCIUserData, pdu.GetPDUUserData()))
		if _debug != nil {
			_debug("    - xpdu: %r", xpdu)
		}
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it upstream
		return b.Response(NA(xpdu), kwArgs)
	case *OriginalBroadcastNPDU:
		// build a _PDU with a local broadcast address
		xpdu := NewPDU(NA(pdu.GetPduData()), NKW(KWCPCISource, pdu.GetPDUSource(), KWCPCIDestination, NewLocalBroadcast(nil), KWCPCIUserData, pdu.GetPDUUserData()))
		if _debug != nil {
			_debug("    - xpdu: %r", xpdu)
		}
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it upstream
		return b.Response(NA(xpdu), kwArgs)
	case *ForwardedNPDU:
		// build a _PDU with the source from the real source
		xpdu := NewPDU(NA(pdu.GetPduData()), NKW(KWCPCISource, pdu.GetBvlciAddress(), KWCPCIDestination, NewLocalBroadcast(nil), KWCPCIUserData, pdu.GetPDUUserData()))
		if _debug != nil {
			_debug("    - xpdu: %r", xpdu)
		}
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it upstream
		return b.Response(NA(xpdu), kwArgs)
	case *WriteBroadcastDistributionTable:
		// build a response
		xpdu, err := NewResult(ToPtr(model.BVLCResultCode_WRITE_BROADCAST_DISTRIBUTION_TABLE_NAK), NoArgs, NKW(KWCPCIUserData, pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error building result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NA(xpdu), kwArgs)
	case *ReadBroadcastDistributionTable:
		// build a response
		xpdu, err := NewResult(ToPtr(model.BVLCResultCode_READ_BROADCAST_DISTRIBUTION_TABLE_NAK), NoArgs, NKW(KWCPCIUserData, pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error building result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NA(xpdu), kwArgs)
		// build a response
	case *RegisterForeignDevice:
		// build a response
		xpdu, err := NewResult(ToPtr(model.BVLCResultCode_REGISTER_FOREIGN_DEVICE_NAK), NoArgs, NKW(KWCPCIUserData, pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error building result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NA(xpdu), kwArgs)
	case *ReadForeignDeviceTable:
		// build a response
		xpdu, err := NewResult(ToPtr(model.BVLCResultCode_READ_FOREIGN_DEVICE_TABLE_NAK), NoArgs, NKW(KWCPCIUserData, pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error building result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NA(xpdu), kwArgs)
	case *DeleteForeignDeviceTableEntry:
		// build a response
		xpdu, err := NewResult(ToPtr(model.BVLCResultCode_DELETE_FOREIGN_DEVICE_TABLE_ENTRY_NAK), NoArgs, NKW(KWCPCIUserData, pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error building result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NA(xpdu), kwArgs)
	case *DistributeBroadcastToNetwork:
		// build a response
		xpdu, err := NewResult(ToPtr(model.BVLCResultCode_DISTRIBUTE_BROADCAST_TO_NETWORK_NAK), NoArgs, NKW(KWCPCIUserData, pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error building result")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NA(xpdu), kwArgs)
	default:
		b.log.Warn().Type("pduType", pdu).Msg("invalid pdu type")
		return nil
	}
}
