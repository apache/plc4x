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

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/apdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/state_machine"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/time_machine"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/vlan"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

type TNetwork3 struct {
	*StateMachineGroup

	iut      *RouterNode
	vlan1    *Network
	td       *ApplicationLayerStateMachine
	sniffer1 *SnifferStateMachine
	vlan2    *Network
	app2     *ApplicationNode
	sniffer2 *SnifferStateMachine

	t *testing.T

	log zerolog.Logger
}

func NewTNetwork3(t *testing.T) *TNetwork3 {
	localLog := testutils.ProduceTestingLogger(t)
	tn := &TNetwork3{
		t:   t,
		log: localLog,
	}
	tn.StateMachineGroup = NewStateMachineGroup(localLog)

	// reset the time machine
	ResetTimeMachine(StartTime)
	localLog.Trace().Msg("time machine reset")

	var err error
	// implementation under test
	tn.iut, err = NewRouterNode(tn.log)
	require.NoError(t, err)

	// make a little LAN
	tn.vlan1 = NewNetwork(localLog, WithNetworkName("vlan1"), WithNetworkBroadcastAddress(NewLocalBroadcast(nil)))

	// Test devices
	tn.td, err = NewApplicationLayerStateMachine(localLog, "1", tn.vlan1)
	require.NoError(t, err)
	tn.Append(tn.td)

	// sniffer node
	tn.sniffer1, err = NewSnifferStateMachine(localLog, "2", tn.vlan1)
	require.NoError(t, err)
	tn.Append(tn.sniffer1)

	// add the network
	err = tn.iut.AddNetwork("3", tn.vlan1, 1)
	require.NoError(t, err)

	//  make another little LAN
	tn.vlan2 = NewNetwork(tn.log, WithNetworkName("vlan3"), WithNetworkBroadcastAddress(NewLocalBroadcast(nil)))

	// application node, not a state machine
	tn.app2, err = NewApplicationNode(tn.log, "4", tn.vlan2)
	require.NoError(t, err)

	//  sniffer node
	tn.sniffer2, err = NewSnifferStateMachine(localLog, "5", tn.vlan2)
	require.NoError(t, err)
	tn.Append(tn.sniffer2)

	//  add the network
	err = tn.iut.AddNetwork("6", tn.vlan2, 2)
	require.NoError(t, err)

	return tn
}

func (t *TNetwork3) Run(timeLimit time.Duration) {
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

func TestNet3(t *testing.T) {
	t.Skip("Needs more testing") // TODO: fix it
	t.Run("TestSimple", func(t *testing.T) {
		t.Run("testIdle", func(t *testing.T) {
			// create a network
			ExclusiveGlobalTimeMachine(t)
			tnet := NewTNetwork3(t)

			// all start states are successful
			tnet.td.GetStartState().Success("")
			tnet.sniffer1.GetStartState().Success("")
			tnet.sniffer2.GetStartState().Success("")

			// run the group
			tnet.Run(0)
		})
	})
	t.Run("TestUnconfirmedRequests", func(t *testing.T) {
		t.Run("test_local_broadcast", func(t *testing.T) {
			//Test broadcast for any router.
			ExclusiveGlobalTimeMachine(t)

			// create a network
			tnet := NewTNetwork3(t)

			// test device sends request, no response
			whois, err := NewWhoIsRequest(NoArgs, NKW(KWCPCIDestination, NewLocalBroadcast(nil)))
			require.NoError(t, err)
			tnet.td.GetStartState().Doc("1-1-0").
				Send(whois, nil).Doc("1-1-1").
				Timeout(3*time.Second, nil).Doc("1-1-2").
				Success("")

			// sniffer on network 1 sees the request and the response
			tnet.sniffer1.GetStartState().Doc("1-2-0").
				Receive(NA(PDUMatcher),
					NKW(KWTestPDUData, xtob(
						"01.00"+ //version, network layer
							"10 08", // unconfirmed Who-Is
					),
					),
				).Doc("1-2-1").
				Timeout(3*time.Second, nil).Doc("1-2-2").
				Success("")

			// nothing received on network 2
			tnet.sniffer2.GetStartState().Success("")

			// run the group
			tnet.Run(0)
		})
		t.Run("test_remote_broadcast_2", func(t *testing.T) {
			//Test broadcast, matching device.
			ExclusiveGlobalTimeMachine(t)

			// create a network
			tnet := NewTNetwork3(t)

			// test device sends request, no response
			whois, err := NewWhoIsRequest(NoArgs, NKW(KWCPCIDestination, NewRemoteBroadcast(2, nil)))
			require.NoError(t, err)
			tnet.td.GetStartState().Doc("2-1-0").
				Send(whois, nil).Doc("2-1-1").
				Success("")

			// sniffer on network 1 sees the request and the response
			tnet.sniffer1.GetStartState().Doc("2-2-0").
				Receive(NA(PDUMatcher),
					NKW(KWTestPDUData, xtob("01.80.00.00.02")), // who is router to network

				).Doc("2-2-1").
				Receive(NA(PDUMatcher),
					NKW(KWTestPDUData, xtob("01.80.01.00.02")), // I am router to network

				).Doc("2-2-1").
				Receive(NA(PDUMatcher),
					NKW(KWTestPDUData, xtob("01.20.00.02.00.ff"+ // remote broadcast goes out
						"10.08",
					)),
				).Doc("2-2-1").
				Receive(NA(PDUMatcher),
					NKW(KWTestPDUData, xtob("01.08.00.02.01.04"+ // unicast response
						"10.00.c4.02.00.00.04.22.04.00.91.00.22.03.e7",
					)),
				).Doc("2-2-2").
				Timeout(3*time.Second, nil).Doc("2-2-3").
				Success("")

			// network 2 sees local broadcast request and unicast response
			tnet.sniffer2.GetStartState().Doc("2-3-0").
				Receive(NA(PDUMatcher),
					NKW(KWTestPDUData, xtob("01.08.00.01.01.01"+ // local broadcast
						"10.08",
					)),
				).Doc("2-3-1").
				Receive(NA(PDUMatcher),
					NKW(KWTestPDUData, xtob("01.20.00.01.01.01.ff"+ // unicast response
						"10.00.c4.02.00.00.04.22.04.00.91.00.22.03.e7",
					)),
				).Doc("2-3-1").
				Timeout(3*time.Second, nil).Doc("2-3-2").
				Success("")

			// run the group
			tnet.Run(0)
		})
		t.Run("test_remote_broadcast_3", func(t *testing.T) {

			//Test broadcast, matching device.
			ExclusiveGlobalTimeMachine(t)

			// create a network
			tnet := NewTNetwork3(t)

			// test device sends request, no response
			whois, err := NewWhoIsRequest(NoArgs, NKW(KWCPCIDestination, NewRemoteBroadcast(3, nil)))
			require.NoError(t, err)
			tnet.td.GetStartState().Doc("3-1-0").
				Send(whois, nil).Doc("3-1-1").
				Success("")

			// sniffer on network 1 sees the request and the response
			tnet.sniffer1.GetStartState().Doc("3-2-0").
				Receive(NA(PDUMatcher),
					NKW(KWTestPDUData, xtob(
						"01.80.00.00.03", // who is router to network
					)),
				).Doc("3-2-1").
				Timeout(3*time.Second, nil).Doc("3-2-3").
				Success("")

			// network 2 sees local broadcast looking for network 3
			tnet.sniffer2.GetStartState().
				Receive(NA(PDUMatcher),
					NKW(KWTestPDUData, xtob(
						"01.88.00.01.01.01.00.00.03",
					)),
				).Doc("3-3-1").
				Timeout(3*time.Second, nil).Doc("3-3-2").
				Success("")

			// run the group
			tnet.Run(0)
		})
		t.Run("test_global_broadcast", func(t *testing.T) {
			//Test broadcast, matching device.
			ExclusiveGlobalTimeMachine(t)

			// create a network
			tnet := NewTNetwork3(t)

			// test device sends request, no response
			whois, err := NewWhoIsRequest(NoArgs, NKW(KWCPCIDestination, NewGlobalBroadcast(nil)))
			require.NoError(t, err)
			tnet.td.GetStartState().Doc("4-1-0").
				Send(whois, nil).Doc("4-1-1").
				Receive(NA((*IAmRequest)(nil)), NoKWArgs()).Doc("4-1-2").
				Success("")

			// sniffer on network 1 sees the request and the response
			tnet.sniffer1.GetStartState().Doc("3-2-0").
				Receive(NA(PDUMatcher),
					NKW(KWTestPDUData, xtob(
						"01.20.ff.ff.00.ff"+
							"10.08",
					)),
				).Doc("4-2-1").
				Receive(NA(PDUMatcher),
					NKW(KWTestPDUData, xtob(
						"01.08.00.02.01.04"+
							"10.00.c4.02.00.00.04.22.04.00.91.00.22.03.e7",
					)),
				).Doc("4-2-2").
				Timeout(3*time.Second, nil).Doc("4-2-3").
				Success("")

			// network 2 sees local broadcast v
			tnet.sniffer2.GetStartState().Doc("4-3-0").
				Receive(NA(PDUMatcher),
					NKW(KWTestPDUData, xtob(
						"01.28.ff.ff.00.00.01.01.01.fe"+
							"10.08",
					)),
				).Doc("4-3-1").
				Receive(NA(PDUMatcher),
					NKW(KWTestPDUData, xtob(
						"01.20.00.01.01.01.ff"+
							"10.00.c4.02.00.00.04.22.04.00.91.00.22.03.e7",
					)),
				).Doc("4-3-3").
				Timeout(3*time.Second, nil).Doc("4-3-3").
				Success("")

			// run the group
			tnet.Run(0)
		})
	})
	t.Run("TestConfirmedRequests", func(t *testing.T) {
		t.Skip("implement me") // TODO
	})
}
