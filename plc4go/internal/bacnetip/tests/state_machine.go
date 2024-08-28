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

package tests

import (
	"bytes"
	"fmt"
	"slices"
	"time"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/internal/bacnetip"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Transition Instances of this class are transitions betweeen getStates of a state machine.
type Transition struct {
	nextState State
}

func (t Transition) String() string {
	return fmt.Sprintf("Transition{nextState: %s}", t.nextState)
}

type SendTransition struct {
	Transition
	pdu bacnetip.PDU
}

func (t SendTransition) String() string {
	return fmt.Sprintf("SendTransition{Transition: %s, pdu: %s}", t.Transition, t.pdu)
}

type criteria struct {
	pduType  any
	pduAttrs map[bacnetip.KnownKey]any
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
	fn     func(args bacnetip.Args, kwargs bacnetip.KWArgs) error
	args   bacnetip.Args
	kwargs bacnetip.KWArgs
}

func (f fnargs) String() string {
	return fmt.Sprintf("fnargs{fn: %t, args: %s, kwargs: %s}", f.fn == nil, f.args, f.kwargs)
}

type CallTransition struct {
	Transition
	fnargs fnargs
}

func (t CallTransition) String() string {
	return fmt.Sprintf("CallTransition{Transition: %s, fnargs: %s}", t.Transition, t.fnargs)
}

func MatchPdu(localLog zerolog.Logger, pdu any, pduType any, pduAttrs map[bacnetip.KnownKey]any) (matches bool) {
	// check the type
	if pduType != nil && fmt.Sprintf("%T", pdu) != fmt.Sprintf("%T", pduType) {
		localLog.Debug().Type("got", pdu).Type("want", pduType).Msg("failed match, wrong type")
		return false
	}
	for attrName, attrValue := range pduAttrs {
		attrLog := localLog.With().Str("attrName", string(attrName)).Interface("attrValue", attrValue).Logger()
		switch attrName {
		case bacnetip.KWPPDUSource:
			equals := pdu.(bacnetip.PDU).GetPDUSource().Equals(attrValue)
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacnetip.KWPDUDestination:
			equals := pdu.(bacnetip.PDU).GetPDUDestination().Equals(attrValue)
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case "x": // only used in test cases
			equals := bytes.Equal(pdu.(interface{ X() []byte }).X(), attrValue.([]byte))
			if !equals {
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
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacnetip.KWPDUData:
			got := pdu.(bacnetip.PDU).GetPduData()
			var want []byte
			switch attrValue := attrValue.(type) {
			case []byte:
				want = attrValue
			case bacnetip.PDUData:
				want = attrValue.GetPduData()
			default:
				attrLog.Debug().Type("type", attrValue).Msg("mismatch, attr unhandled")
			}
			equals := bytes.Equal(got, want)
			if !equals {
				attrLog.Debug().Hex("got", got).Hex("want", want).Stringer("diff", utils.DiffHex(want, got)).Msg("mismatch")
			}
			if !equals {
				attrLog.Debug().Msg("pduData doesn't match")
				return false
			}
		case bacnetip.KWWirtnNetwork:
			wirtn, ok := pdu.(*bacnetip.WhoIsRouterToNetwork)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			net := wirtn.GetWirtnNetwork()
			if net == nil {
				return false
			}
			equals := *net == attrValue
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacnetip.KWIartnNetworkList:
			iamrtn, ok := pdu.(*bacnetip.IAmRouterToNetwork)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			net := iamrtn.GetIartnNetworkList()
			uint16s, ok := attrValue.([]uint16)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := slices.Equal(net, uint16s)
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacnetip.KWIcbrtnNetwork:
			iamrtn, ok := pdu.(*bacnetip.ICouldBeRouterToNetwork)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := iamrtn.GetIcbrtnNetwork() == attrValue
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacnetip.KWIcbrtnPerformanceIndex:
			iamrtn, ok := pdu.(*bacnetip.ICouldBeRouterToNetwork)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := iamrtn.GetIcbrtnPerformanceIndex() == attrValue
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacnetip.KWRmtnRejectionReason:
			iamrtn, ok := pdu.(*bacnetip.RejectMessageToNetwork)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := iamrtn.GetRmtnRejectionReason() == attrValue
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacnetip.KWRmtnDNET:
			iamrtn, ok := pdu.(*bacnetip.RejectMessageToNetwork)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := iamrtn.GetRmtnDNET() == attrValue
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacnetip.KWRbtnNetworkList:
			rbtn, ok := pdu.(*bacnetip.RouterBusyToNetwork)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			net := rbtn.GetRbtnNetworkList()
			uint16s, ok := attrValue.([]uint16)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := slices.Equal(net, uint16s)
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacnetip.KWRatnNetworkList:
			ratn, ok := pdu.(*bacnetip.RouterAvailableToNetwork)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			net := ratn.GetRatnNetworkList()
			uint16s, ok := attrValue.([]uint16)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := slices.Equal(net, uint16s)
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacnetip.KWIrtTable:
			irt, ok := pdu.(*bacnetip.InitializeRoutingTable)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			irts := irt.GetIrtTable()
			oirts, ok := attrValue.([]*bacnetip.RoutingTableEntry)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := slices.EqualFunc(irts, oirts, func(entry *bacnetip.RoutingTableEntry, entry2 *bacnetip.RoutingTableEntry) bool {
				return entry.Equals(entry2)
			})
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacnetip.KWIrtaTable:
			irta, ok := pdu.(*bacnetip.InitializeRoutingTableAck)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			irts := irta.GetIrtaTable()
			oirts, ok := attrValue.([]*bacnetip.RoutingTableEntry)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := slices.EqualFunc(irts, oirts, func(entry *bacnetip.RoutingTableEntry, entry2 *bacnetip.RoutingTableEntry) bool {
				return entry.Equals(entry2)
			})
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacnetip.KWEctnDNET:
			ectn, ok := pdu.(*bacnetip.EstablishConnectionToNetwork)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := ectn.GetEctnDNET() == attrValue
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacnetip.KWEctnTerminationTime:
			ectn, ok := pdu.(*bacnetip.EstablishConnectionToNetwork)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := ectn.GetEctnTerminationTime() == attrValue
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacnetip.KWDctnDNET:
			dctn, ok := pdu.(*bacnetip.DisconnectConnectionToNetwork)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := dctn.GetDctnDNET() == attrValue
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacnetip.KWNniNet:
			nni, ok := pdu.(*bacnetip.NetworkNumberIs)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := nni.GetNniNet() == attrValue
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacnetip.KWNniFlag:
			nni, ok := pdu.(*bacnetip.NetworkNumberIs)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			return nni.GetNniFlag() == attrValue
		case bacnetip.KWBvlciResultCode:
			r, ok := pdu.(*bacnetip.Result)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := r.GetBvlciResultCode() == attrValue
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacnetip.KWBvlciBDT:
			var iwbdt []*bacnetip.Address
			switch pdu := pdu.(type) {
			case *bacnetip.WriteBroadcastDistributionTable:
				iwbdt = pdu.GetBvlciBDT()
			case *bacnetip.ReadBroadcastDistributionTableAck:
				iwbdt = pdu.GetBvlciBDT()
			default:
				attrLog.Trace().Type("type", pdu).Msg("doesn't match")
				return false
			}
			owbdt, ok := attrValue.([]*bacnetip.Address)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := slices.EqualFunc(iwbdt, owbdt, func(a *bacnetip.Address, b *bacnetip.Address) bool {
				return a.Equals(b)
			})
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacnetip.KWBvlciAddress:
			var address *bacnetip.Address
			switch pdu := pdu.(type) {
			case *bacnetip.ForwardedNPDU:
				address = pdu.GetBvlciAddress()
			case *bacnetip.DeleteForeignDeviceTableEntry:
				address = pdu.GetBvlciAddress()
			default:
				attrLog.Trace().Type("type", pdu).Msg("doesn't match")
				return false
			}
			equals := address.Equals(attrValue)
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacnetip.KWFdAddress:
			panic("implement me")
			equals := true // TODO temporary
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacnetip.KWFdTTL:
			panic("implement me")
			equals := true // TODO temporary
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacnetip.KWFdRemain:
			panic("implement me")
			equals := true // TODO temporary
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacnetip.KWBvlciTimeToLive:
			rfd, ok := pdu.(*bacnetip.RegisterForeignDevice)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := rfd.GetBvlciTimeToLive() == attrValue
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacnetip.KWBvlciFDT:
			rfdta, ok := pdu.(*bacnetip.ReadForeignDeviceTableAck)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			ifdt := rfdta.GetBvlciFDT()
			oifdt, ok := attrValue.([]*bacnetip.FDTEntry)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := slices.EqualFunc(ifdt, oifdt, func(a *bacnetip.FDTEntry, b *bacnetip.FDTEntry) bool {
				return a.Equals(b)
			})
			if !equals {
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

type TimeoutTask struct {
	*bacnetip.OneShotTask

	fn     func(args bacnetip.Args, kwargs bacnetip.KWArgs) error
	args   bacnetip.Args
	kwargs bacnetip.KWArgs
}

func NewTimeoutTask(fn func(args bacnetip.Args, kwargs bacnetip.KWArgs) error, args bacnetip.Args, kwargs bacnetip.KWArgs, when *time.Time) *TimeoutTask {
	task := &TimeoutTask{
		fn:     fn,
		args:   args,
		kwargs: kwargs,
	}
	task.OneShotTask = bacnetip.NewOneShotTask(task, when)
	return task
}

func (t *TimeoutTask) ProcessTask() error {
	return t.fn(t.args, t.kwargs)
}

func (t *TimeoutTask) String() string {
	return fmt.Sprintf("TimeoutTask(%v, fn: %t, args: %s, kwargs: %s)", t.Task, t.fn != nil, t.args, t.kwargs)
}
