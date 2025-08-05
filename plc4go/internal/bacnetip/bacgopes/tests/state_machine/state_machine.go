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

package state_machine

import (
	"bytes"
	"fmt"
	"slices"
	"time"

	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/bvll"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/npdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

var _debug = CreateDebugPrinter()

// Transition Instances of this class are transitions betweeen getStates of a state machine.
type Transition struct {
	nextState State
}

func (t Transition) String() string {
	return fmt.Sprintf("Transition{nextState: %s}", t.nextState)
}

type SendTransition struct {
	Transition
	pdu PDU
}

func (t SendTransition) String() string {
	return fmt.Sprintf("SendTransition{Transition: %s, pdu: %s}", t.Transition, t.pdu)
}

type criteria struct {
	pduType  any
	pduAttrs map[KnownKey]any
}

func (c criteria) String() string {
	return fmt.Sprintf("criteria{%T, %v}", c.pduType, c.pduAttrs)
}

type ReceiveTransition struct {
	Transition
	criteria criteria
}

func (t ReceiveTransition) String() string {
	return fmt.Sprintf("ReceiveTransition{Transition: %s, criteria: %s}", t.Transition, t.criteria)
}

type EventTransition struct {
	Transition
	eventId string
}

func (t EventTransition) String() string {
	return fmt.Sprintf("EventTransition{Transition: %s, eventId: %s}", t.Transition, t.eventId)
}

type TimeoutTransition struct {
	Transition
	timeout time.Time
}

func (t TimeoutTransition) String() string {
	return fmt.Sprintf("TimeoutTransition{Transition: %s, timeout: %s}", t.Transition, t.timeout)
}

type fnargs struct {
	fn     GenericFunction
	args   Args
	kwArgs KWArgs
}

func (f fnargs) String() string {
	return fmt.Sprintf("fnargs{fn: %t, args: %s, kwArgs: %s}", f.fn == nil, f.args, f.kwArgs)
}

type CallTransition struct {
	Transition
	fnargs fnargs
}

func (t CallTransition) String() string {
	return fmt.Sprintf("CallTransition{Transition: %s, fnargs: %s}", t.Transition, t.fnargs)
}

func MatchPdu(localLog zerolog.Logger, pdu any, pduType any, pduAttrs map[KnownKey]any) (matches bool) {
	if _debug != nil {
		_debug("match_pdu %r %r %r", pdu, pduType, pduAttrs)
	}
	// check the type
	switch pduType := pduType.(type) {
	case func(any) bool:
		if !pduType(pdu) {
			if _debug != nil {
				_debug("    - failed match, wrong type")
			}
			localLog.Debug().Type("got", pdu).Interface("gotValue", pdu).Msg("failed match, func says no")
			return false
		}
	default:
		if pduType != nil && fmt.Sprintf("%T", pdu) != fmt.Sprintf("%T", pduType) {
			if _debug != nil {
				_debug("    - failed match, wrong type")
			}
			localLog.Debug().Type("got", pdu).Interface("gotValue", pdu).Type("want", pduType).Msg("failed match, wrong type")
			return false
		}
	}

	for attrName, attrValue := range pduAttrs {
		attrLog := localLog.With().Str("attrName", string(attrName)).Interface("attrValue", attrValue).Logger()
		switch attrName {
		case KWCPCISource:
			equals := pdu.(PDU).GetPDUSource().Equals(attrValue)
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case KWCPCIDestination:
			equals := pdu.(PDU).GetPDUDestination().Equals(attrValue)
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case "x": // only used in test cases
			equals := bytes.Equal(pdu.(interface{ X() []byte }).X(), attrValue.([]byte))
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case "y": // only used in test cases
			return false
		case "a": // only used in test cases
			a := pdu.(interface{ A() int }).A()
			if a == 0 {
				return false
			}
			equals := a == attrValue.(int)
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case "b": // only used in test cases
			b := pdu.(interface{ B() int }).B()
			if b == 0 {
				return false
			}
			equals := b == attrValue.(int)
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case KWTestPDUData:
			got := pdu.(PDU).GetPduData()
			var want []byte
			switch attrValue := attrValue.(type) {
			case []byte:
				want = attrValue
			case PDUData:
				want = attrValue.GetPduData()
			default:
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Debug().Type("type", attrValue).Msg("mismatch, attr unhandled")
			}
			equals := bytes.Equal(got, want)
			if !equals {
				attrLog.Debug().Hex("got", got).Hex("want", want).Stringer("diff", utils.DiffHex(want, got)).Msg("mismatch")
			}
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Debug().Msg("pduData doesn't match")
				return false
			}
		case KWWirtnNetwork:
			wirtn, ok := pdu.(*WhoIsRouterToNetwork)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			net := wirtn.GetWirtnNetwork()
			if net == nil {
				return false
			}
			equals := *net == attrValue
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case KWIartnNetworkList:
			iamrtn, ok := pdu.(*IAmRouterToNetwork)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			net := iamrtn.GetIartnNetworkList()
			uint16s, ok := attrValue.([]uint16)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := slices.Equal(net, uint16s)
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case KWIcbrtnNetwork:
			iamrtn, ok := pdu.(*ICouldBeRouterToNetwork)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := iamrtn.GetIcbrtnNetwork() == attrValue
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case KWIcbrtnPerformanceIndex:
			iamrtn, ok := pdu.(*ICouldBeRouterToNetwork)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := iamrtn.GetIcbrtnPerformanceIndex() == attrValue
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case KWRmtnRejectionReason:
			iamrtn, ok := pdu.(*RejectMessageToNetwork)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := iamrtn.GetRmtnRejectionReason() == attrValue
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case KWRmtnDNET:
			iamrtn, ok := pdu.(*RejectMessageToNetwork)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := iamrtn.GetRmtnDNET() == attrValue
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case KWRbtnNetworkList:
			rbtn, ok := pdu.(*RouterBusyToNetwork)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			net := rbtn.GetRbtnNetworkList()
			uint16s, ok := attrValue.([]uint16)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := slices.Equal(net, uint16s)
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case KWRatnNetworkList:
			ratn, ok := pdu.(*RouterAvailableToNetwork)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			net := ratn.GetRatnNetworkList()
			uint16s, ok := attrValue.([]uint16)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := slices.Equal(net, uint16s)
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case KWIrtTable:
			irt, ok := pdu.(*InitializeRoutingTable)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			irts := irt.GetIrtTable()
			oirts, ok := attrValue.([]*RoutingTableEntry)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := slices.EqualFunc(irts, oirts, func(entry *RoutingTableEntry, entry2 *RoutingTableEntry) bool {
				return entry.Equals(entry2)
			})
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case KWIrtaTable:
			irta, ok := pdu.(*InitializeRoutingTableAck)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			irts := irta.GetIrtaTable()
			oirts, ok := attrValue.([]*RoutingTableEntry)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := slices.EqualFunc(irts, oirts, func(entry *RoutingTableEntry, entry2 *RoutingTableEntry) bool {
				return entry.Equals(entry2)
			})
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case KWEctnDNET:
			ectn, ok := pdu.(*EstablishConnectionToNetwork)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := ectn.GetEctnDNET() == attrValue
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case KWEctnTerminationTime:
			ectn, ok := pdu.(*EstablishConnectionToNetwork)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := ectn.GetEctnTerminationTime() == attrValue
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case KWDctnDNET:
			dctn, ok := pdu.(*DisconnectConnectionToNetwork)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := dctn.GetDctnDNET() == attrValue
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case KWNniNet:
			nni, ok := pdu.(*NetworkNumberIs)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := nni.GetNniNet() == attrValue
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case KWNniFlag:
			nni, ok := pdu.(*NetworkNumberIs)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			return nni.GetNniFlag() == attrValue
		case KWBvlciResultCode:
			r, ok := pdu.(*Result)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := r.GetBvlciResultCode() == attrValue
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case KWBvlciBDT:
			var iwbdt []*Address
			switch pdu := pdu.(type) {
			case *WriteBroadcastDistributionTable:
				iwbdt = pdu.GetBvlciBDT()
			case *ReadBroadcastDistributionTableAck:
				iwbdt = pdu.GetBvlciBDT()
			default:
				attrLog.Trace().Type("type", pdu).Msg("doesn't match")
				return false
			}
			owbdt, ok := attrValue.([]*Address)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := slices.EqualFunc(iwbdt, owbdt, func(a *Address, b *Address) bool {
				return a.Equals(b)
			})
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case KWBvlciAddress:
			var address *Address
			switch pdu := pdu.(type) {
			case *ForwardedNPDU:
				address = pdu.GetBvlciAddress()
			case *DeleteForeignDeviceTableEntry:
				address = pdu.GetBvlciAddress()
			default:
				attrLog.Trace().Type("type", pdu).Msg("doesn't match")
				return false
			}
			equals := address.Equals(attrValue)
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case KWFdAddress:
			panic("implement me")
			equals := true // TODO temporary
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case KWFdTTL:
			panic("implement me")
			equals := true // TODO temporary
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case KWFdRemain:
			panic("implement me")
			equals := true // TODO temporary
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case KWBvlciTimeToLive:
			rfd, ok := pdu.(*RegisterForeignDevice)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := rfd.GetBvlciTimeToLive() == attrValue
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case KWBvlciFDT:
			rfdta, ok := pdu.(*ReadForeignDeviceTableAck)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			ifdt := rfdta.GetBvlciFDT()
			oifdt, ok := attrValue.([]*FDTEntry)
			if !ok {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := slices.EqualFunc(ifdt, oifdt, func(a *FDTEntry, b *FDTEntry) bool {
				equals := a.Equals(b)
				if !equals {
					attrLog.Trace().Stringer("a", a).Stringer("b", b).Msg("doesn't match")
				}
				return equals
			})
			if !equals {
				if _debug != nil {
					_debug("    - failed match, attr value: %r, %r", attrName, attrValue)
				}
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		default:
			panic("implement " + attrName)
		}
	}
	localLog.Trace().Msg("successful match")
	return true
}
