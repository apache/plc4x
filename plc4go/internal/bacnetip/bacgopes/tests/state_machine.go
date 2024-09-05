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

	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes"
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
	pdu bacgopes.PDU
}

func (t SendTransition) String() string {
	return fmt.Sprintf("SendTransition{Transition: %s, pdu: %s}", t.Transition, t.pdu)
}

type criteria struct {
	pduType  any
	pduAttrs map[bacgopes.KnownKey]any
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
	fn     func(args bacgopes.Args, kwargs bacgopes.KWArgs) error
	args   bacgopes.Args
	kwargs bacgopes.KWArgs
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

func MatchPdu(localLog zerolog.Logger, pdu any, pduType any, pduAttrs map[bacgopes.KnownKey]any) (matches bool) {
	// check the type
	switch pduType := pduType.(type) {
	case func(any) bool:
		if !pduType(pdu) {
			localLog.Debug().Type("got", pdu).Interface("gotValue", pdu).Msg("failed match, func says no")
			return false
		}
	default:
		if pduType != nil && fmt.Sprintf("%T", pdu) != fmt.Sprintf("%T", pduType) {
			localLog.Debug().Type("got", pdu).Interface("gotValue", pdu).Type("want", pduType).Msg("failed match, wrong type")
			return false
		}
	}

	for attrName, attrValue := range pduAttrs {
		attrLog := localLog.With().Str("attrName", string(attrName)).Interface("attrValue", attrValue).Logger()
		switch attrName {
		case bacgopes.KWPPDUSource:
			equals := pdu.(bacgopes.PDU).GetPDUSource().Equals(attrValue)
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacgopes.KWPDUDestination:
			equals := pdu.(bacgopes.PDU).GetPDUDestination().Equals(attrValue)
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
		case bacgopes.KWPDUData:
			got := pdu.(bacgopes.PDU).GetPduData()
			var want []byte
			switch attrValue := attrValue.(type) {
			case []byte:
				want = attrValue
			case bacgopes.PDUData:
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
		case bacgopes.KWWirtnNetwork:
			wirtn, ok := pdu.(*bacgopes.WhoIsRouterToNetwork)
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
		case bacgopes.KWIartnNetworkList:
			iamrtn, ok := pdu.(*bacgopes.IAmRouterToNetwork)
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
		case bacgopes.KWIcbrtnNetwork:
			iamrtn, ok := pdu.(*bacgopes.ICouldBeRouterToNetwork)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := iamrtn.GetIcbrtnNetwork() == attrValue
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacgopes.KWIcbrtnPerformanceIndex:
			iamrtn, ok := pdu.(*bacgopes.ICouldBeRouterToNetwork)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := iamrtn.GetIcbrtnPerformanceIndex() == attrValue
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacgopes.KWRmtnRejectionReason:
			iamrtn, ok := pdu.(*bacgopes.RejectMessageToNetwork)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := iamrtn.GetRmtnRejectionReason() == attrValue
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacgopes.KWRmtnDNET:
			iamrtn, ok := pdu.(*bacgopes.RejectMessageToNetwork)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := iamrtn.GetRmtnDNET() == attrValue
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacgopes.KWRbtnNetworkList:
			rbtn, ok := pdu.(*bacgopes.RouterBusyToNetwork)
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
		case bacgopes.KWRatnNetworkList:
			ratn, ok := pdu.(*bacgopes.RouterAvailableToNetwork)
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
		case bacgopes.KWIrtTable:
			irt, ok := pdu.(*bacgopes.InitializeRoutingTable)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			irts := irt.GetIrtTable()
			oirts, ok := attrValue.([]*bacgopes.RoutingTableEntry)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := slices.EqualFunc(irts, oirts, func(entry *bacgopes.RoutingTableEntry, entry2 *bacgopes.RoutingTableEntry) bool {
				return entry.Equals(entry2)
			})
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacgopes.KWIrtaTable:
			irta, ok := pdu.(*bacgopes.InitializeRoutingTableAck)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			irts := irta.GetIrtaTable()
			oirts, ok := attrValue.([]*bacgopes.RoutingTableEntry)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := slices.EqualFunc(irts, oirts, func(entry *bacgopes.RoutingTableEntry, entry2 *bacgopes.RoutingTableEntry) bool {
				return entry.Equals(entry2)
			})
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacgopes.KWEctnDNET:
			ectn, ok := pdu.(*bacgopes.EstablishConnectionToNetwork)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := ectn.GetEctnDNET() == attrValue
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacgopes.KWEctnTerminationTime:
			ectn, ok := pdu.(*bacgopes.EstablishConnectionToNetwork)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := ectn.GetEctnTerminationTime() == attrValue
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacgopes.KWDctnDNET:
			dctn, ok := pdu.(*bacgopes.DisconnectConnectionToNetwork)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := dctn.GetDctnDNET() == attrValue
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacgopes.KWNniNet:
			nni, ok := pdu.(*bacgopes.NetworkNumberIs)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := nni.GetNniNet() == attrValue
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacgopes.KWNniFlag:
			nni, ok := pdu.(*bacgopes.NetworkNumberIs)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			return nni.GetNniFlag() == attrValue
		case bacgopes.KWBvlciResultCode:
			r, ok := pdu.(*bacgopes.Result)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := r.GetBvlciResultCode() == attrValue
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacgopes.KWBvlciBDT:
			var iwbdt []*bacgopes.Address
			switch pdu := pdu.(type) {
			case *bacgopes.WriteBroadcastDistributionTable:
				iwbdt = pdu.GetBvlciBDT()
			case *bacgopes.ReadBroadcastDistributionTableAck:
				iwbdt = pdu.GetBvlciBDT()
			default:
				attrLog.Trace().Type("type", pdu).Msg("doesn't match")
				return false
			}
			owbdt, ok := attrValue.([]*bacgopes.Address)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := slices.EqualFunc(iwbdt, owbdt, func(a *bacgopes.Address, b *bacgopes.Address) bool {
				return a.Equals(b)
			})
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacgopes.KWBvlciAddress:
			var address *bacgopes.Address
			switch pdu := pdu.(type) {
			case *bacgopes.ForwardedNPDU:
				address = pdu.GetBvlciAddress()
			case *bacgopes.DeleteForeignDeviceTableEntry:
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
		case bacgopes.KWFdAddress:
			panic("implement me")
			equals := true // TODO temporary
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacgopes.KWFdTTL:
			panic("implement me")
			equals := true // TODO temporary
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacgopes.KWFdRemain:
			panic("implement me")
			equals := true // TODO temporary
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacgopes.KWBvlciTimeToLive:
			rfd, ok := pdu.(*bacgopes.RegisterForeignDevice)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := rfd.GetBvlciTimeToLive() == attrValue
			if !equals {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
		case bacgopes.KWBvlciFDT:
			rfdta, ok := pdu.(*bacgopes.ReadForeignDeviceTableAck)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			ifdt := rfdta.GetBvlciFDT()
			oifdt, ok := attrValue.([]*bacgopes.FDTEntry)
			if !ok {
				attrLog.Trace().Msg("doesn't match")
				return false
			}
			equals := slices.EqualFunc(ifdt, oifdt, func(a *bacgopes.FDTEntry, b *bacgopes.FDTEntry) bool {
				equals := a.Equals(b)
				if !equals {
					attrLog.Trace().Stringer("a", a).Stringer("b", b).Msg("doesn't match")
				}
				return equals
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

//go:generate plc4xGenerator -type=TimeoutTask -prefix=task_
type TimeoutTask struct {
	*bacgopes.OneShotTask

	fn     func(args bacgopes.Args, kwargs bacgopes.KWArgs) error `ignore:"true"`
	args   bacgopes.Args
	kwargs bacgopes.KWArgs
}

func NewTimeoutTask(fn func(args bacgopes.Args, kwargs bacgopes.KWArgs) error, args bacgopes.Args, kwargs bacgopes.KWArgs, when *time.Time) *TimeoutTask {
	task := &TimeoutTask{
		fn:     fn,
		args:   args,
		kwargs: kwargs,
	}
	task.OneShotTask = bacgopes.NewOneShotTask(task, when)
	return task
}

func (t *TimeoutTask) ProcessTask() error {
	return t.fn(t.args, t.kwargs)
}
