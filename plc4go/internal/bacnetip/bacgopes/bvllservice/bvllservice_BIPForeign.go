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
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/bvll"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/task"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

//go:generate plc4xGenerator -type=BIPForeign -prefix=bvllservice_
type BIPForeign struct {
	*BIPSAP
	Client
	Server
	*OneShotTask

	registrationStatus      int
	bbmdAddress             *Address
	bbmdTimeToLive          *int
	registrationTimeoutTask *OneShotFunctionTask

	// regular args
	argAddr *Address `ignore:"true"`
	argTTL  *int     `ignore:"true"`

	// pass through args
	argSapID *int `ignore:"true"`
	argCid   *int `ignore:"true"`
	argSid   *int `ignore:"true"`

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
	client, err := NewClient(localLog, b, OptionalOption(b.argCid, WithClientCID))
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	b.Client = client
	server, err := NewServer(localLog, b, OptionalOption(b.argSid, WithServerSID))
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
	pdu := GA[PDU](args, 0)

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
		return b.Request(NA(xpdu), NoKWArgs)
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
		return b.Request(NA(xpdu), NoKWArgs)
	default:
		return errors.Errorf("invalid destination address: %s", pdu.GetPDUDestination())
	}
}

func (b *BIPForeign) Confirmation(args Args, kwargs KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Confirmation")
	pdu := GA[PDU](args, 0)

	switch msg := pdu.GetRootMessage().(type) {
	// check for a registration request result
	case model.BVLCResult:
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
	case model.BVLCOriginalUnicastNPDU:
		// build a vanilla _PDU
		xpdu := NewPDU(NoArgs, NKW(KWCompRootMessage, msg.GetNpdu(), KWCPCISource, pdu.GetPDUSource(), KWCPCIDestination, pdu.GetPDUDestination()))
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it upstream
		return b.Response(NA(xpdu), NoKWArgs)
	case model.BVLCForwardedNPDU:
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
		source, err := NewAddress(NA(append(ip, Uint16ToPort(port)...)))
		if err != nil {
			return errors.Wrap(err, "error building a ip")
		}
		xpdu := NewPDU(NoArgs, NKW(KWCompRootMessage, msg.GetNpdu(), KWCPCISource, source, KWCPCIDestination, NewLocalBroadcast(nil)))
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it upstream
		return b.Response(NA(xpdu), NoKWArgs)
	case model.BVLCReadBroadcastDistributionTableAck:
		// send this to the service access point
		return b.SapResponse(args, NoKWArgs)
	case model.BVLCReadForeignDeviceTableAck:
		// send this to the service access point
		return b.SapResponse(args, NoKWArgs)
	case model.BVLCWriteBroadcastDistributionTable:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(model.BVLCResultCode_WRITE_BROADCAST_DISTRIBUTION_TABLE_NAK))
		if err != nil {
			return errors.Wrap(err, "error building a result")
		}
		xpdu.SetPDUUserData(pdu.GetPDUUserData()) // TODO: upstream sets this in the constructor
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NA(xpdu), NoKWArgs)
	case model.BVLCReadBroadcastDistributionTable:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(model.BVLCResultCode_READ_BROADCAST_DISTRIBUTION_TABLE_NAK))
		if err != nil {
			return errors.Wrap(err, "error building a result")
		}
		xpdu.SetPDUUserData(pdu.GetPDUUserData()) // TODO: upstream sets this in the constructor
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NA(xpdu), NoKWArgs)
	case model.BVLCRegisterForeignDevice:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(model.BVLCResultCode_REGISTER_FOREIGN_DEVICE_NAK))
		if err != nil {
			return errors.Wrap(err, "error building a result")
		}
		xpdu.SetPDUUserData(pdu.GetPDUUserData()) // TODO: upstream sets this in the constructor
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NA(xpdu), NoKWArgs)
	case model.BVLCReadForeignDeviceTable:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(model.BVLCResultCode_READ_FOREIGN_DEVICE_TABLE_NAK))
		if err != nil {
			return errors.Wrap(err, "error building a result")
		}
		xpdu.SetPDUUserData(pdu.GetPDUUserData()) // TODO: upstream sets this in the constructor
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NA(xpdu), NoKWArgs)
	case model.BVLCDeleteForeignDeviceTableEntry:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(model.BVLCResultCode_DELETE_FOREIGN_DEVICE_TABLE_ENTRY_NAK))
		if err != nil {
			return errors.Wrap(err, "error building a result")
		}
		xpdu.SetPDUUserData(pdu.GetPDUUserData()) // TODO: upstream sets this in the constructor
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NA(xpdu), NoKWArgs)
	case model.BVLCDistributeBroadcastToNetwork:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(model.BVLCResultCode_DISTRIBUTE_BROADCAST_TO_NETWORK_NAK))
		if err != nil {
			return errors.Wrap(err, "error building a result")
		}
		xpdu.SetPDUUserData(pdu.GetPDUUserData()) // TODO: upstream sets this in the constructor
		xpdu.SetPDUDestination(pdu.GetPDUSource())

		// send it downstream
		return b.Request(NA(xpdu), NoKWArgs)
	case model.BVLCOriginalBroadcastNPDU:
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
	pdu := NewPDU(NoArgs, NKW(KWCompRootMessage, model.NewBVLCRegisterForeignDevice(0), KWCPCIDestination, b.bbmdAddress))

	// send it downstream
	if err := b.Request(NA(pdu), NoKWArgs); err != nil {
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
	if err := b.Request(NA(pdu), NoKWArgs); err != nil {
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
