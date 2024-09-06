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
	"slices"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/bvll"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/deleteme"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/task"
)

//go:generate plc4xGenerator -type=BIPBBMD -prefix=bvllservice_
type BIPBBMD struct {
	*BIPSAP
	Client
	Server
	*RecurringTask
	*DebugContents `ignore:"true"`

	bbmdAddress *Address
	bbmdBDT     []*Address
	bbmdFDT     []*FDTEntry

	// Pass Through args
	argSapID *int `ignore:"true"`
	argCID   *int `ignore:"true"`
	argSID   *int `ignore:"true"`

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
	b.Client, err = NewClient(localLog, b, OptionalOption(b.argCID, WithClientCID))
	if err != nil {
		return nil, errors.Wrap(err, "error creating Client")
	}
	b.Server, err = NewServer(localLog, b, OptionalOption(b.argSID, WithServerSID))
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

	pdu := Get[PDU](args, 0)

	// check for local stations
	if pdu.GetPDUDestination().AddrType == LOCAL_STATION_ADDRESS {
		// make an original unicast PDU
		var err error
		xpdu, err := NewOriginalUnicastNPDU(pdu, WithOriginalUnicastNPDUDestination(pdu.GetPDUDestination()), WithOriginalUnicastNPDUUserData(pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error creating OriginalUnicastNPDU")
		}
		//           if settings.route_aware and PDUDestination.addrRoute:
		//               xpdu.pduDestination = PDUDestination.addrRoute
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
		//           if settings.route_aware and PDUDestination.addrRoute:
		//               xpdu.pduDestination = PDUDestination.addrRoute
		b.log.Debug().Stringer("xpdu", xpdu).Msg("original broadcast xpdu")

		// send it downstream
		err = b.Request(NewArgs(xpdu), NoKWArgs)
		if err != nil {
			return errors.Wrap(err, "error sending request downstream")
		}

		// skip other processing if the route was provided
		//           if settings.route_aware and PDUDestination.addrRoute:
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

	pdu := Get[PDU](args, 0)
	switch pdu.GetRootMessage().(type) {
	case model.BVLCResult: //some kind of response to a request
		// send this to the service access point
		return b.SapResponse(NewArgs(pdu), NoKWArgs)
	case model.BVLCWriteBroadcastDistributionTable:
		// build a response
		xpdu, err := NewResult(WithResultBvlciResultCode(model.BVLCResultCode_WRITE_BROADCAST_DISTRIBUTION_TABLE_NAK))
		if err != nil {
			return errors.Wrap(err, "error creating Result")
		}
		xpdu.SetPDUUserData(pdu.GetPDUUserData())

		// send it downstream
		return b.Request(NewArgs(pdu), NoKWArgs)
	case model.BVLCReadBroadcastDistributionTable:
		// build a response
		xpdu, err := NewReadBroadcastDistributionTableAck(WithReadBroadcastDistributionTableAckBDT(b.bbmdBDT...))
		if err != nil {
			return errors.Wrap(err, "error creating Result")
		}
		xpdu.SetPDUUserData(pdu.GetPDUUserData())

		// send it downstream
		return b.Request(NewArgs(pdu), NoKWArgs)
	case model.BVLCReadBroadcastDistributionTableAck:
		// send it to the service access point
		return b.SapResponse(NewArgs(pdu), NoKWArgs)
	case model.BVLCForwardedNPDU:
		pdu := pdu.(*ForwardedNPDU) // TODO: check if this cast is fine
		// send it upstream if there is a network layer
		if b.HasServerPeer() {
			xpdu := NewPDU(NewMessageBridge(pdu.GetPduData()...), WithPDUSource(pdu.GetBvlciAddress()), WithPDUDestination(NewLocalBroadcast(nil)), WithPDUUserData(pdu.GetPDUUserData()))
			//               if settings.route_aware:
			//                   xpdu.pduSource.addrRoute = PDUSource
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
	case model.BVLCRegisterForeignDevice:
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
	case model.BVLCReadForeignDeviceTable:
		// build a response
		xpdu, err := NewReadForeignDeviceTableAck(WithReadForeignDeviceTableAckFDT(b.bbmdFDT...))
		if err != nil {
			return errors.Wrap(err, "error creating ack")
		}
		xpdu.SetPDUDestination(pdu.GetPDUSource())
		xpdu.SetPDUUserData(pdu.GetPDUUserData())

		// send it downstream
		return b.Request(NewArgs(xpdu), NoKWArgs)
	case model.BVLCReadForeignDeviceTableAck:
		// send this to the service access point
		return b.SapResponse(NewArgs(pdu), NoKWArgs)
	case model.BVLCDeleteForeignDeviceTableEntry:
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
	case model.BVLCDistributeBroadcastToNetwork:
		// send it upstream if there is a network layer
		if b.HasServerPeer() {
			xpdu := NewPDU(NewMessageBridge(pdu.GetPduData()...), WithPDUSource(pdu.GetPDUSource()), WithPDUDestination(NewLocalBroadcast(nil)), WithPDUUserData(pdu.GetPDUUserData()))
			//               if settings.route_aware:
			//                   xpdu.pduSource.addrRoute = PDUSource
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
	case model.BVLCOriginalUnicastNPDU:
		// send it upstream if there is a network layer
		if b.HasServerPeer() {
			// build a PDU
			xpdu := NewPDU(NewMessageBridge(pdu.GetPduData()...), WithPDUSource(pdu.GetPDUSource()), WithPDUDestination(pdu.GetPDUDestination()), WithPDUUserData(pdu.GetPDUUserData()))
			//               if settings.route_aware:
			//                   xpdu.pduSource.addrRoute = PDUSource
			b.log.Debug().Stringer("xpdu", xpdu).Msg("upstream xpdu")

			return b.Response(NewArgs(xpdu), NoKWArgs)
		}
	case model.BVLCOriginalBroadcastNPDU:
		// send it upstream if there is a network layer
		if b.HasServerPeer() {
			// build a PDU with a local broadcast address
			xpdu := NewPDU(NewMessageBridge(pdu.GetPduData()...), WithPDUSource(pdu.GetPDUSource()), WithPDUDestination(NewLocalBroadcast(nil)), WithPDUUserData(pdu.GetPDUUserData()))
			//               if settings.route_aware:
			//                   xpdu.pduSource.addrRoute = PDUSource
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

func (b *BIPBBMD) RegisterForeignDevice(address Arg, ttl uint16) (model.BVLCResultCode, error) {
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
			return model.BVLCResultCode_REGISTER_FOREIGN_DEVICE_NAK, errors.Wrap(err, "error creating address")
		}
	default:
		return model.BVLCResultCode_REGISTER_FOREIGN_DEVICE_NAK, errors.Errorf("invalid address type: %T", address)
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
	return model.BVLCResultCode_SUCCESSFUL_COMPLETION, nil
}

func (b *BIPBBMD) DeleteForeignDeviceTableEntry(address Arg) (model.BVLCResultCode, error) {
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
			return model.BVLCResultCode_REGISTER_FOREIGN_DEVICE_NAK, errors.Wrap(err, "error creating address")
		}
	default:
		return model.BVLCResultCode_REGISTER_FOREIGN_DEVICE_NAK, errors.Errorf("invalid address type: %T", address)
	}

	// find it and delete it
	stat := model.BVLCResultCode_SUCCESSFUL_COMPLETION
	deleted := false
	b.bbmdFDT = slices.DeleteFunc(b.bbmdFDT, func(entry *FDTEntry) bool {
		if addr.Equals(entry.FDAddress) {
			deleted = true
			return true
		}
		return false
	})
	if !deleted {
		stat = model.BVLCResultCode_DELETE_FOREIGN_DEVICE_TABLE_ENTRY_NAK
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
