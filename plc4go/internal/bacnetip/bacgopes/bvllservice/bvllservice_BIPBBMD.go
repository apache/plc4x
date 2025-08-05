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

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/bvll"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/task"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

//go:generate go tool plc4xGenerator -type=BIPBBMD -prefix=bvllservice_
type BIPBBMD struct {
	*BIPSAP
	ClientContract
	ServerContract
	*RecurringTask
	*DebugContents `ignore:"true"`

	bbmdAddress *Address
	bbmdBDT     []*Address
	bbmdFDT     []*FDTEntry

	log zerolog.Logger
}

func NewBIPBBMD(localLog zerolog.Logger, addr *Address, options ...Option) (*BIPBBMD, error) {
	b := &BIPBBMD{log: localLog}
	ApplyAppliers(options, b)
	optionsForParent := AddLeafTypeIfAbundant(options, b)
	b.DebugContents = NewDebugContents(b, "bbmdAddress", "bbmdBDT+", "bbmdFDT+")
	var err error
	b.BIPSAP, err = NewBIPSAP(localLog, b, optionsForParent...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating BIPSAP")
	}
	b.ClientContract, err = NewClient(localLog, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating Client")
	}
	b.ServerContract, err = NewServer(localLog, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating Server")
	}
	b.RecurringTask = NewRecurringTask(localLog, b, WithRecurringTaskInterval(1*time.Second))
	b.AddExtraPrinters(b.RecurringTask)

	if _debug != nil {
		_debug("__init__ %r sapID=%r cid=%r sid=%r", addr, b.GetServiceID(), b.GetClientID(), b.GetServiceID())
	}

	b.bbmdAddress = addr

	// install so process_task runs
	b.InstallTask(WithInstallTaskOptionsNone())

	return b, nil
}

func (b *BIPBBMD) GetDebugAttr(attr string) any {
	switch attr {
	case "bbmdAddress":
		if b.bbmdAddress != nil {
			return b.bbmdAddress
		}
	case "bbmdBDT":
		return b.bbmdBDT
	case "bbmdFDT":
		return b.bbmdFDT
	default:
		return nil
	}
	return nil
}

func (b *BIPBBMD) Indication(args Args, kwArgs KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Indication")
	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("indication %r", pdu)
	}

	// check for local stations
	if pdu.GetPDUDestination().AddrType == LOCAL_STATION_ADDRESS {
		// make an original unicast PDU
		var err error
		xpdu, err := NewOriginalUnicastNPDU(NA(pdu), NKW(KWCPCIDestination, pdu.GetPDUDestination(), KWCPCIUserData, pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error creating OriginalUnicastNPDU")
		}
		//           if settings.route_aware and PDUDestination.addrRoute:
		//               xpdu.pduDestination = PDUDestination.addrRoute
		if _debug != nil {
			_debug("    - original unicast xpdu: %r", xpdu)
		}
		b.log.Debug().Stringer("xpdu", xpdu).Msg("original unicast xpdu")

		// send it downstream
		return b.Request(NA(xpdu), NoKWArgs())
	} else if pdu.GetPDUDestination().AddrType == LOCAL_BROADCAST_ADDRESS { // check for broadcaste
		// make an original unicast PDU
		var xpdu PDU
		var err error
		xpdu, err = NewOriginalBroadcastNPDU(NA(pdu), NKW(KWCPCIDestination, pdu.GetPDUDestination(), KWCPCIUserData, pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error creating OriginalUnicastNPDU")
		}
		//           if settings.route_aware and PDUDestination.addrRoute:
		//               xpdu.pduDestination = PDUDestination.addrRoute
		if _debug != nil {
			_debug("    - original broadcast xpdu: %r", xpdu)
		}
		b.log.Debug().Stringer("xpdu", xpdu).Msg("original broadcast xpdu")

		// send it downstream
		err = b.Request(NA(xpdu), NoKWArgs())
		if err != nil {
			return errors.Wrap(err, "error sending request downstream")
		}

		// skip other processing if the route was provided
		//           if settings.route_aware and PDUDestination.addrRoute:
		//               return

		// make a forwarded PDU
		xpdu, err = NewForwardedNPDU(b.bbmdAddress, NA(pdu), NKW(KWCPCIUserData, pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error creating ForwardedNPDU")
		}
		if _debug != nil {
			_debug("    - forwarded xpdu: %r", xpdu)
		}
		b.log.Debug().Stringer("xpdu", xpdu).Msg("forwarded xpdu")

		// send it to the peers
		for _, bdte := range b.bbmdBDT {
			if !bdte.Equals(b.bbmdAddress) {
				dest, err := NewAddress(NA(&AddressTuple[int, uint16]{Left: int(*bdte.AddrIP | ^*bdte.AddrMask), Right: *bdte.AddrPort}))
				if err != nil {
					return errors.Wrap(err, "error creating address tuple")
				}
				xpdu.SetPDUDestination(dest)
				b.log.Debug().Stringer("pduDestination", xpdu.GetPDUDestination()).Msg("sending to peer")
				if err := b.Request(NA(xpdu), NoKWArgs()); err != nil {
					return errors.Wrap(err, "error sending request")
				}
			}
		}

		// send it to the registered foreign devies
		for _, fdte := range b.bbmdFDT {
			xpdu.SetPDUDestination(fdte.FDAddress)
			if _debug != nil {
				_debug("    - sending to foreign device: %r", xpdu.GetPDUDestination())
			}
			b.log.Debug().Stringer("pduDestination", xpdu.GetPDUDestination()).Msg("sending to foreign device")
			if err := b.Request(NA(xpdu), NoKWArgs()); err != nil {
				return errors.Wrap(err, "error sending request")
			}
		}
		return err
	} else {
		return errors.Errorf("invalid destination address: %s", pdu.GetPDUDestination())
	}
}

func (b *BIPBBMD) Confirmation(args Args, kwArgs KWArgs) error {
	b.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Confirmation")
	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("confirmation %r", pdu)
	}

	switch pdu := pdu.(type) {
	case *Result: //some kind of response to a request
		// send this to the service access point
		return b.SapResponse(NA(pdu), NoKWArgs())
	case *WriteBroadcastDistributionTable:
		// build a response
		xpdu, err := NewResult(ToPtr(model.BVLCResultCode_WRITE_BROADCAST_DISTRIBUTION_TABLE_NAK), NoArgs, NoKWArgs())
		if err != nil {
			return errors.Wrap(err, "error creating Result")
		}
		xpdu.SetPDUUserData(pdu.GetPDUUserData())
		if _debug != nil {
			_debug("    - xpdu: %r", xpdu)
		}

		// send it downstream
		return b.Request(NA(pdu), NoKWArgs())
	case *ReadBroadcastDistributionTable:
		// build a response
		xpdu, err := NewReadBroadcastDistributionTableAck(b.bbmdBDT, NoArgs, NoKWArgs())
		if err != nil {
			return errors.Wrap(err, "error creating Result")
		}
		xpdu.SetPDUUserData(pdu.GetPDUUserData())
		if _debug != nil {
			_debug("    - xpdu: %r", xpdu)
		}

		// send it downstream
		return b.Request(NA(pdu), NoKWArgs())
	case *ReadBroadcastDistributionTableAck:
		// send it to the service access point
		return b.SapResponse(NA(pdu), NoKWArgs())
	case *ForwardedNPDU:
		// send it upstream if there is a network layer
		if b.HasServerPeer() {
			xpdu := NewPDU(NA(pdu.GetPduData()), NKW(
				KWCPCISource, pdu.GetBvlciAddress(),
				KWCPCIDestination, NewLocalBroadcast(nil),
				KWCPCIUserData, pdu.GetPDUUserData(),
			),
			)
			//               if settings.route_aware:
			//                   xpdu.pduSource.addrRoute = PDUSource
			if _debug != nil {
				_debug("    - upstream xpdu: %r", xpdu)
			}
			b.log.Debug().Stringer("xpdu", xpdu).Msg("upstream xpdu")

			return b.Response(NA(xpdu), NoKWArgs())
		}

		// build a forwarded NPDU to send out
		xpdu, err := NewForwardedNPDU(pdu.GetBvlciAddress(), NA(pdu), NKW(KWCPCISource, pdu.GetBvlciAddress(), KWCPCIDestination, pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error creating ForwardedNPDU")
		}
		if _debug != nil {
			_debug("    - forwarded xpdu: %r", xpdu)
		}
		b.log.Debug().Stringer("xpdu", xpdu).Msg("forwarded xpdu")

		// if this was unicast to us, do next hop
		if pdu.GetPDUDestination().AddrType == LOCAL_STATION_ADDRESS {
			if _debug != nil {
				_debug("    - unicast message")
			}
			b.log.Trace().Msg("unicast message")

			// if this BBMD is listed in its BDT, send a local broadcast
			if slices.ContainsFunc(b.bbmdBDT, func(address *Address) bool {
				return address.Equals(b.bbmdAddress)
			}) {
				if _debug != nil {
					_debug("    - local broadcast")
				}
				b.log.Trace().Msg("local broadcast")
				return b.Request(NA(xpdu), NoKWArgs())
			}
		} else if pdu.GetPDUDestination().AddrType == LOCAL_BROADCAST_ADDRESS {
			if _debug != nil {
				_debug("    - directed broadcast message")
			}
			b.log.Trace().Msg("directed broadcast message")
		} else {
			b.log.Warn().Stringer("destination", pdu.GetPDUDestination()).Msg("invalid destination address")
		}

		// send it to the registered foreign devices
		for _, fdte := range b.bbmdFDT {
			xpdu.SetPDUDestination(fdte.FDAddress)
			if _debug != nil {
				_debug("    - sending to foreign device: %r", xpdu.GetPDUDestination())
			}
			b.log.Warn().Stringer("destination", xpdu.GetPDUDestination()).Msg("sending to foreign device")
			if err := b.Request(NA(xpdu), NoKWArgs()); err != nil {
				return errors.Wrapf(err, "error sending request to destination %s", xpdu.GetPDUDestination())
			}
		}
		return nil
	case *RegisterForeignDevice:
		// process the request
		stat, err := b.RegisterForeignDevice(pdu.GetPDUSource(), pdu.GetBvlciTimeToLive())
		if err != nil {
			return errors.Wrap(err, "error registering device")
		}

		// build a response
		xpdu, err := NewResult(&stat, NoArgs, NKW(KWCPCIDestination, pdu.GetPDUSource(), KWCPCIUserData, pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error creating Result")
		}
		if _debug != nil {
			_debug("    - xpdu: %r", xpdu)
		}
		b.log.Debug().Stringer("xpdu", xpdu).Msg("xpdu")

		// send it downstream
		return b.Request(NA(xpdu), NoKWArgs())
	case *ReadForeignDeviceTable:
		// build a response
		xpdu, err := NewReadForeignDeviceTableAck(b.bbmdFDT, NoArgs, NKW(KWCPCIDestination, pdu.GetPDUSource(), KWCPCIUserData, pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error creating ack")
		}
		if _debug != nil {
			_debug("    - xpdu: %r", xpdu)
		}

		// send it downstream
		return b.Request(NA(xpdu), NoKWArgs())
	case *ReadForeignDeviceTableAck:
		// send this to the service access point
		return b.SapResponse(NA(pdu), NoKWArgs())
	case *DeleteForeignDeviceTableEntry:
		// process the request
		stat, err := b.DeleteForeignDeviceTableEntry(pdu.GetBvlciAddress())
		if err != nil {
			return errors.Wrap(err, "error deleting entry")
		}

		// build a response
		xpdu, err := NewResult(&stat, NoArgs, NKW(KWCPCIDestination, pdu.GetPDUSource(), KWCPCIUserData, pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error creating Result")
		}
		if _debug != nil {
			_debug("    - xpdu: %r", xpdu)
		}

		// send it downstream
		return b.Request(NA(xpdu), NoKWArgs())
	case *DistributeBroadcastToNetwork:
		// send it upstream if there is a network layer
		if b.HasServerPeer() {
			xpdu := NewPDU(NA(pdu.GetPduData()), NKW(
				KWCPCISource, pdu.GetPDUSource(),
				KWCPCIDestination, NewLocalBroadcast(nil),
				KWCPCIUserData, pdu.GetPDUUserData(),
			),
			)
			if _debug != nil {
				_debug("    - upstream xpdu: %r", xpdu)
			}
			b.log.Debug().Stringer("xpdu", xpdu).Msg("upstream xpdu")

			return b.Response(NA(xpdu), NoKWArgs())
		}

		// build a forwarded NPDU to send out
		xpdu, err := NewForwardedNPDU(pdu.GetPDUSource(), NA(pdu), NKW(KWCPCIUserData, pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error creating ForwardedNPDU")
		}
		if _debug != nil {
			_debug("    - forwarded xpdu: %r", xpdu)
		}
		b.log.Debug().Stringer("xpdu", xpdu).Msg("forwarded xpdu")

		// send it to the peers
		for _, bdte := range b.bbmdBDT {
			if bdte.Equals(b.bbmdAddress) {
				xpdu.SetPDUDestination(NewLocalBroadcast(nil))
				if _debug != nil {
					_debug("    - local broadcast")
				}
				b.log.Trace().Msg("local broadcast")
				err = b.Request(NA(xpdu), NoKWArgs())
				if err != nil {
					return errors.Wrap(err, "error sending local broadcast")
				}
			} else {
				address, err := NewAddress(NA(&AddressTuple[int, uint16]{Left: int(*bdte.AddrIP | ^*bdte.AddrMask), Right: *bdte.AddrPort}))
				if err != nil {
					return errors.Wrap(err, "error creating address")
				}
				xpdu.SetPDUDestination(address)
				if _debug != nil {
					_debug("    - sending to peer: %r", xpdu.GetPDUDestination())
				}
				b.log.Debug().Stringer("designation", xpdu.GetPDUDestination()).Msg("sending to peer")
				err = b.Request(NA(xpdu), NoKWArgs())
				if err != nil {
					return errors.Wrap(err, "error sending")
				}
			}
		}

		// send it to the registered foreign devices
		for _, fdte := range b.bbmdFDT {
			if !fdte.Equals(pdu.GetPDUSource()) {
				xpdu.SetPDUDestination(fdte.FDAddress)
				if _debug != nil {
					_debug("    - sending to foreign device: %r", xpdu.GetPDUDestination())
				}
				b.log.Warn().Stringer("destination", xpdu.GetPDUDestination()).Msg("sending to foreign device")
				if err := b.Request(NA(xpdu), NoKWArgs()); err != nil {
					return errors.Wrapf(err, "error sending request to destination %s", xpdu.GetPDUDestination())
				}
			}
		}
		return nil
	case *OriginalUnicastNPDU:
		// send it upstream if there is a network layer
		if b.HasServerPeer() {
			// build a PDU
			xpdu := NewPDU(NA(pdu.GetPduData()), NKW(
				KWCPCISource, pdu.GetPDUSource(),
				KWCPCIDestination, pdu.GetPDUDestination(),
				KWCPCIUserData, pdu.GetPDUUserData(),
			),
			)
			//               if settings.route_aware:
			//                   xpdu.pduSource.addrRoute = PDUSource
			b.log.Debug().Stringer("xpdu", xpdu).Msg("upstream xpdu")

			return b.Response(NA(xpdu), NoKWArgs())
		}
	case *OriginalBroadcastNPDU:
		// send it upstream if there is a network layer
		if b.HasServerPeer() {
			// build a PDU with a local broadcast address
			xpdu := NewPDU(NA(pdu.GetPduData()), NKW(
				KWCPCISource, pdu.GetPDUSource(),
				KWCPCIDestination, NewLocalBroadcast(nil),
				KWCPCIUserData, pdu.GetPDUUserData(),
			),
			)
			//               if settings.route_aware:
			//                   xpdu.pduSource.addrRoute = PDUSource
			if _debug != nil {
				_debug("    - upstream xpdu: %r", xpdu)
			}
			b.log.Debug().Stringer("xpdu", xpdu).Msg("upstream xpdu")

			if err := b.Response(NA(xpdu), NoKWArgs()); err != nil {
				return errors.Wrap(err, "error sending local broadcast")
			}
		}

		// build a forwarded NPDU
		xpdu, err := NewForwardedNPDU(pdu.GetPDUSource(), NA(pdu), NKW(KWCPCIUserData, pdu.GetPDUUserData()))
		if err != nil {
			return errors.Wrap(err, "error creating ForwardedNPDU")
		}
		if _debug != nil {
			_debug("    - forwarded xpdu: %r", xpdu)
		}
		b.log.Debug().Stringer("xpdu", xpdu).Msg("forwarded xpdu")

		// send it to the peers
		for _, bdte := range b.bbmdBDT {
			if !bdte.Equals(b.bbmdAddress) {
				address, err := NewAddress(NA(&AddressTuple[int, uint16]{Left: int(*bdte.AddrIP | ^*bdte.AddrMask), Right: *bdte.AddrPort}))
				if err != nil {
					return errors.Wrap(err, "error creating address")
				}
				xpdu.SetPDUDestination(address)
				if _debug != nil {
					_debug("    - sending to peer: %r", xpdu.GetPDUDestination())
				}
				b.log.Debug().Stringer("designation", xpdu.GetPDUDestination()).Msg("sending to peer")
				err = b.Request(NA(xpdu), NoKWArgs())
				if err != nil {
					return errors.Wrap(err, "error sending")
				}
			}
		}

		// send it to the registered foreign devices
		for _, fdte := range b.bbmdFDT {
			xpdu.SetPDUDestination(fdte.FDAddress)
			if _debug != nil {
				_debug("    - sending to foreign device: %r", xpdu.GetPDUDestination())
			}
			b.log.Warn().Stringer("destination", xpdu.GetPDUDestination()).Msg("sending to foreign device")
			if err := b.Request(NA(xpdu), NoKWArgs()); err != nil {
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
	if _debug != nil {
		_debug("register_foreign_device %r %r", address, ttl)
	}
	b.log.Debug().Interface("address", address).Uint16("ttl", ttl).Msg("registering foreign device")

	var addr *Address
	// see if it is an address or make it one
	switch address := address.(type) {
	case *Address:
		addr = address
	case string:
		var err error
		addr, err = NewAddress(NA(address))
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
	if _debug != nil {
		_debug("delete_foreign_device_table_entry %r", address)
	}
	b.log.Debug().Interface("address", address).Msg("delete foreign device")

	var addr *Address
	// see if it is an address or make it one
	switch address := address.(type) {
	case *Address:
		addr = address
	case string:
		var err error
		addr, err = NewAddress(NA(address))
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
			if _debug != nil {
				_debug("foreign device expired: %r", fdte.FDAddress)
			}
			b.log.Debug().Stringer("addr", fdte.FDAddress).Msg("foreign device expired")
			return true
		}
		return false
	})
	return nil
}

func (b *BIPBBMD) AddPeer(address Arg) error {
	if _debug != nil {
		_debug("add_peer %r", address)
	}
	b.log.Debug().Interface("adddress", address).Msg("addr")

	var addr *Address
	// see if it is an address or make it one
	switch address := address.(type) {
	case *Address:
		addr = address
	case string:
		var err error
		addr, err = NewAddress(NA(address))
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
	if _debug != nil {
		_debug("delete_peer %r", address)
	}
	b.log.Debug().Interface("adddress", address).Msg("addr")

	var addr *Address
	// see if it is an address or make it one
	switch address := address.(type) {
	case *Address:
		addr = address
	case string:
		var err error
		addr, err = NewAddress(NA(address))
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
