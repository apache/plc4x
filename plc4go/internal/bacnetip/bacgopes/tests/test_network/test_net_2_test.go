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

package test_network

import (
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/npdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/state_machine"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/time_machine"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/vlan"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

type TNetwork2 struct {
	*StateMachineGroup

	trafficLog *TrafficLog
	iut1       *RouterNode
	iut2       *RouterNode
	vlan1      *Network
	td         *NetworkLayerStateMachine
	sniffer1   *SnifferStateMachine
	vlan2      *Network
	sniffer2   *SnifferStateMachine
	vlan3      *Network
	sniffer3   *SnifferStateMachine

	t *testing.T

	log zerolog.Logger
}

func NewTNetwork2(t *testing.T) *TNetwork2 {
	localLog := testutils.ProduceTestingLogger(t)
	tn := &TNetwork2{
		t:   t,
		log: localLog,
	}
	tn.StateMachineGroup = NewStateMachineGroup(localLog)

	// reset the time machine
	ResetTimeMachine(StartTime)
	localLog.Trace().Msg("time machine reset")

	// create a traffic log
	tn.trafficLog = new(TrafficLog)

	var err error
	// implementation under test
	tn.iut1, err = NewRouterNode(tn.log) // router from vlan1 to vlan2
	require.NoError(t, err)
	tn.iut2, err = NewRouterNode(tn.log) // router from vlan2 to vlan3
	require.NoError(t, err)

	// make a little LAN
	tn.vlan1 = NewNetwork(localLog, WithNetworkName("vlan1"), WithNetworkBroadcastAddress(NewLocalBroadcast(nil)), WithNetworkTrafficLogger(tn.trafficLog))

	// Test devices
	tn.td, err = NewNetworkLayerStateMachine(localLog, "1", tn.vlan1)
	require.NoError(t, err)
	tn.Append(tn.td)

	// sniffer node
	tn.sniffer1, err = NewSnifferStateMachine(localLog, "2", tn.vlan1)
	require.NoError(t, err)
	tn.Append(tn.sniffer1)

	// connect vlan1 to iut1
	err = tn.iut1.AddNetwork("3", tn.vlan1, 1)
	require.NoError(t, err)

	// make another little LAN
	tn.vlan2 = NewNetwork(tn.log, WithNetworkName("vlan2"), WithNetworkBroadcastAddress(NewLocalBroadcast(nil)), WithNetworkTrafficLogger(tn.trafficLog))

	// sniffer node
	tn.sniffer2, err = NewSnifferStateMachine(localLog, "4", tn.vlan2)
	require.NoError(t, err)
	tn.Append(tn.sniffer2)

	// connect vlan2 to both routers
	err = tn.iut1.AddNetwork("5", tn.vlan2, 2)
	require.NoError(t, err)
	err = tn.iut2.AddNetwork("6", tn.vlan2, 2)
	require.NoError(t, err)

	//  make another little LAN
	tn.vlan3 = NewNetwork(tn.log, WithNetworkName("vlan3"), WithNetworkBroadcastAddress(NewLocalBroadcast(nil)), WithNetworkTrafficLogger(tn.trafficLog))

	// sniffer node
	tn.sniffer3, err = NewSnifferStateMachine(localLog, "6", tn.vlan2)
	require.NoError(t, err)
	tn.Append(tn.sniffer3)

	// connect vlan3 to the second router
	err = tn.iut2.AddNetwork("8", tn.vlan3, 3)
	require.NoError(t, err)

	return tn
}

func (t *TNetwork2) Run(timeLimit time.Duration) {
	if timeLimit == 0 {
		timeLimit = 60 * time.Second
	}
	t.log.Debug().Dur("time_limit", timeLimit).Msg("run")

	// run the group
	err := t.StateMachineGroup.Run()
	require.NoError(t.t, err)

	// run it some time
	RunTimeMachine(t.log, timeLimit, time.Time{})
	t.log.Trace().Msg("time machine finished")
	for _, machine := range t.StateMachineGroup.GetStateMachines() {
		t.log.Debug().Stringer("machine", machine).Stringers("transactionLog", ToStringers(machine.GetTransactionLog())).Msg("Machine:")
	}

	// check for success
	success, failed := t.CheckForSuccess()
	assert.True(t.t, success)
	assert.False(t.t, failed)
}

func TestNet2(t *testing.T) {
	t.Run("TestSimple", func(t *testing.T) {
		t.Run("testIdle", func(t *testing.T) {
			// create a network
			ExclusiveGlobalTimeMachine(t)
			tnet := NewTNetwork2(t)

			// all start states are successful
			tnet.td.GetStartState().Success("")
			tnet.sniffer1.GetStartState().Success("")
			tnet.sniffer2.GetStartState().Success("")
			tnet.sniffer3.GetStartState().Success("")

			// run the group
			tnet.Run(0)
		})
	})
	t.Run("TestWhoIsRouterToNetwork", func(t *testing.T) {
		t.Run("test_01", func(t *testing.T) {
			//Test broadcast for any router.
			ExclusiveGlobalTimeMachine(t)

			// create a network
			tnet := NewTNetwork2(t)

			// test device sends request, sees response
			whois, err := NewWhoIsRouterToNetwork(NoArgs, NewKWArgs(KWCPCIDestination, NewLocalBroadcast(nil)))
			require.NoError(t, err)
			tnet.td.GetStartState().Doc("1-1-0").
				Send(whois, nil).Doc("1-1-1").
				Receive(NA((*IAmRouterToNetwork)(nil)), NKW(KWIartnNetworkList, []uint16{2})).Doc("1-1-2").
				Success("")

			// sniffer on network 1 sees the request and the response
			tnet.sniffer1.GetStartState().Doc("1-2-0").
				Receive(NA(PDUMatcher),
					NKW(KWTestPDUData, xtob(
						"01.80"+ //version, network layer
							"00", //message type, no network
					)),
				).Doc("1-2-1").
				Receive(NA(PDUMatcher),
					NKW(KWTestPDUData, xtob(
						"01.80"+ //version, network layer
							"01 0002", //message type and network list
					)),
				).Doc("1-2-2").
				Success("")

			// nothing received on network 2
			tnet.sniffer2.GetStartState().Doc("1-3-0").
				Timeout(3*time.Second, nil).Doc("1-3-1").
				Success("")

			// nothing received on network 3
			tnet.sniffer3.GetStartState().Doc("1-4-0").
				Timeout(3*time.Second, nil).Doc("1-4-1").
				Success("")

			// run the group
			tnet.Run(0)
		})
		t.Run("test_02", func(t *testing.T) {
			//Test broadcast for existing router.
			ExclusiveGlobalTimeMachine(t)
			// create a network
			tnet := NewTNetwork2(t)

			// test device sends request, sees response
			whois, err := NewWhoIsRouterToNetwork(NoArgs, NewKWArgs(KWCPCIDestination, NewLocalBroadcast(nil)), WithWhoIsRouterToNetworkNet(2))
			require.NoError(t, err)
			tnet.td.GetStartState().Doc("2-1-0").
				Send(whois, nil).Doc("2-1-1").
				Receive(NA((*IAmRouterToNetwork)(nil)), NKW(KWIartnNetworkList, []uint16{2})).Doc("2-1-2").
				Success("")

			tnet.sniffer1.GetStartState().Success("")

			// nothing received on network 2
			tnet.sniffer2.GetStartState().Doc("2-2-0").
				Timeout(3*time.Second, nil).Doc("2-2-1").
				Success("")

			// nothing received on network 2
			tnet.sniffer3.GetStartState().Doc("2-3-0").
				Timeout(3*time.Second, nil).Doc("2-3-1").
				Success("")

			// run the group
			tnet.Run(0)
		})
		t.Run("test_03", func(t *testing.T) {
			//Test broadcast for non-existing router.
			ExclusiveGlobalTimeMachine(t)
			// create a network
			tnet := NewTNetwork2(t)

			// test device sends request, sees response
			whois, err := NewWhoIsRouterToNetwork(NoArgs, NewKWArgs(KWCPCIDestination, NewLocalBroadcast(nil)), WithWhoIsRouterToNetworkNet(4))
			require.NoError(t, err)
			tnet.td.GetStartState().Doc("3-1-0").
				Send(whois, nil).Doc("3-1-1").
				Timeout(3*time.Second, nil).Doc("3-1-2").
				Success("")

			// sniffer on network 1 sees the request and the response
			tnet.sniffer1.GetStartState().Doc("3-2-0").
				Receive(NA(PDUMatcher),
					NKW(KWTestPDUData, xtob(
						"01.80"+ //version, network layer
							"00 0004", //message type, and network
					)),
				).Doc("3-2-1").
				Success("")

			// sniffer on network 2 sees request forwarded by router
			tnet.sniffer2.GetStartState().Doc("3-3-0").
				Receive(NA(PDUMatcher),
					NKW(KWTestPDUData, xtob(
						"01.88"+ //version, network layer
							"0001 01 01"+ // snet/slen/sadr
							"00 0004", //message type, and network
					)),
				).Doc("3-3-1").
				Success("")

			tnet.sniffer3.GetStartState().Doc("3-4-0").
				Receive(NA(PDUMatcher),
					NKW(KWTestPDUData, xtob(
						"01.88"+ //version, network layer
							"0001 01 01"+ // snet/slen/sadr
							"00 0004", //message type, and network
					)),
				).Doc("3-4-1").
				Success("")

			// run the group
			tnet.Run(0)
		})
	})
}
