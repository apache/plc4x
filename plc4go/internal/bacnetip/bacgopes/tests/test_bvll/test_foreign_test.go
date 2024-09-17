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

package test_bvll

import (
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/bvll"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/quick"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/state_machine"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/time_machine"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/vlan"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

type TFNetwork struct {
	*StateMachineGroup

	trafficLog *TrafficLog
	router     *IPRouter
	vlan5      *IPNetwork
	vlan6      *IPNetwork
	fd         *BIPForeignStateMachine
	bbmd       *BIPBBMDStateMachine

	t *testing.T

	log zerolog.Logger
}

func NewTFNetwork(t *testing.T) *TFNetwork {
	localLog := testutils.ProduceTestingLogger(t)
	tfn := &TFNetwork{
		t:   t,
		log: localLog,
	}
	tfn.StateMachineGroup = NewStateMachineGroup(localLog)

	// reset the time machine
	ResetTimeMachine(StartTime)
	localLog.Trace().Msg("time machine reset")

	// Create a traffic log
	tfn.trafficLog = new(TrafficLog)

	// make a router
	tfn.router = NewIPRouter(localLog)

	// make a home LAN
	tfn.vlan5 = NewIPNetwork(localLog, WithNetworkName("192.168.5.0/24"), WithNetworkTrafficLogger(tfn.trafficLog))
	tfn.router.AddNetwork(quick.Address("192.168.5.1/24"), tfn.vlan5)

	// make a remote LAN
	tfn.vlan6 = NewIPNetwork(localLog, WithNetworkName("192.168.6.0/24"), WithNetworkTrafficLogger(tfn.trafficLog))
	tfn.router.AddNetwork(quick.Address("192.168.6.1/24"), tfn.vlan6)

	var err error
	// the foreign device
	tfn.fd, err = NewBIPForeignStateMachine(localLog, "192.168.6.2/24", tfn.vlan6)
	require.NoError(t, err)
	tfn.Append(tfn.fd)

	// bbmd
	tfn.bbmd, err = NewBIPBBMDStateMachine(localLog, "192.168.5.3/24", tfn.vlan5)
	require.NoError(t, err)
	tfn.Append(tfn.bbmd)
	return tfn
}

func (t *TFNetwork) Run(timeLimit time.Duration) {
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
		t.log.Debug().Stringer("machine", machine).Msg("Machine:")
		for _, s := range machine.GetTransactionLog() {
			t.log.Debug().Stringer("logEntry", s).Msg("logEntry")
		}
	}

	// check for success
	success, failed := t.CheckForSuccess()
	assert.True(t.t, success)
	_ = failed
}

func TestForeign(t *testing.T) {
	t.Run("test_idle", func(t *testing.T) { //Test an idle network, nothing happens is success.
		ExclusiveGlobalTimeMachine(t)

		tnet := NewTFNetwork(t)

		// all start state are successful
		tnet.fd.GetStartState().Success("")
		tnet.bbmd.GetStartState().Success("")

		// run the group
		tnet.Run(0)
	})
	t.Run("test_registration", func(t *testing.T) {
		ExclusiveGlobalTimeMachine(t)
		testingLogger := testutils.ProduceTestingLogger(t)

		// create a network
		tnet := NewTFNetwork(t)

		// tell the B/IP layer of the foreign device to register
		tnet.fd.GetStartState().
			Call(func(args Args, _ KWArgs) error {
				return tnet.fd.bip.Register(args[0].(*Address), args[1].(uint16))
			}, NA(tnet.bbmd.address, uint16(30)), NoKWArgs()).
			Success("")

		// remote sniffer node
		remoteSniffer, err := NewSnifferStateMachine(testingLogger, "192.168.6.254/24", tnet.vlan6)
		require.NoError(t, err)
		tnet.Append(remoteSniffer)

		// sniffer traffic
		remoteSniffer.GetStartState().Doc("1-1-0").
			Receive(NA((*RegisterForeignDevice)(nil)), NoKWArgs()).Doc("1-1-1").
			Receive(NA((*Result)(nil)), NoKWArgs()).Doc("1-1-2").
			SetEvent("fd-registered").Doc("1-1-3").
			Success("")

		// the bbmd is idle
		tnet.bbmd.GetStartState().Success("")

		// home snooper node
		homeSnooper, err := NewBIPStateMachine(testingLogger, "192.168.5.2/24", tnet.vlan5)
		tnet.Append(homeSnooper)

		// snooper will read foreign device table
		homeSnooper.GetStartState().Doc("1-2-0").
			WaitEvent("fd-registered", nil).Doc("1-2-1").
			Send(quick.ReadForeignDeviceTable(tnet.bbmd.address), nil).Doc("1-2-2").
			Receive(NA((*ReadForeignDeviceTableAck)(nil)), NoKWArgs()).Doc("1-2-3").
			Success("")

		// home sniffer node
		homeSniffer, err := NewSnifferStateMachine(testingLogger, "192.168.5.254/24", tnet.vlan5)
		require.NoError(t, err)
		tnet.Append(homeSniffer)

		// sniffer traffic
		homeSniffer.GetStartState().Doc("1-3-0").
			Receive(NA((*RegisterForeignDevice)(nil)), NoKWArgs()).Doc("1-3-1").
			Receive(NA((*Result)(nil)), NoKWArgs()).Doc("1-3-2").
			Receive(NA((*ReadForeignDeviceTable)(nil)), NoKWArgs()).Doc("1-3-3").
			Receive(NA((*ReadForeignDeviceTableAck)(nil)), NoKWArgs()).Doc("1-3-4").
			Success("")

		//  run the group
		tnet.Run(0)
	})
	t.Run("test_refresh_registration", func(t *testing.T) {
		ExclusiveGlobalTimeMachine(t)
		testingLogger := testutils.ProduceTestingLogger(t)

		// create a network
		tnet := NewTFNetwork(t)

		// tell the B/IP layer of the foreign device to register
		tnet.fd.GetStartState().
			Call(func(args Args, _ KWArgs) error {
				return tnet.fd.bip.Register(args[0].(*Address), args[1].(uint16))
			}, NA(tnet.bbmd.address, uint16(10)), NoKWArgs()).
			Success("")

		// the bbmd is idle
		tnet.bbmd.GetStartState().Success("")

		// remote sniffer node
		remoteSniffer, err := NewSnifferStateMachine(testingLogger, "192.168.6.254/24", tnet.vlan6)
		require.NoError(t, err)
		tnet.Append(remoteSniffer)

		// sniffer traffic
		remoteSniffer.GetStartState().Doc("2-1-0").
			Receive(NA((*RegisterForeignDevice)(nil)), NoKWArgs()).Doc("2-1-1").
			Receive(NA((*Result)(nil)), NoKWArgs()).Doc("2-1-1").
			Receive(NA((*RegisterForeignDevice)(nil)), NoKWArgs()).Doc("2-1-3").
			Receive(NA((*Result)(nil)), NoKWArgs()).Doc("2-1-4").
			Success("")

		//  run the group
		tnet.Run(0)
	})
	t.Run("test_unicast", func(t *testing.T) { //Test a unicast message from TD to IUT.
		ExclusiveGlobalTimeMachine(t)
		testingLogger := testutils.ProduceTestingLogger(t)

		tnet := NewTFNetwork(t)

		//make a PDU from node 1 to node 2
		pduData, err := Xtob("dead.beef")
		require.NoError(t, err)
		pdu := NewPDU(NA(pduData), NKW(KWCPCISource, tnet.fd.address, KWCPCIDestination, tnet.bbmd.address))
		t.Logf("    - pdu: %s", pdu)

		// register, wait for ack, send some beef
		tnet.fd.GetStartState().Doc("3-1-0").
			Call(func(args Args, _ KWArgs) error {
				return tnet.fd.bip.Register(args[0].(*Address), args[1].(uint16))
			}, NA(tnet.bbmd.address, uint16(60)), NoKWArgs()).Doc("3-1-1").
			WaitEvent("3-registered", nil).Doc("3-1-2").
			Send(pdu, nil).Doc("3-1-3").
			Success("")

		// the bbmd is happy when it gets the pdu
		tnet.bbmd.GetStartState().
			Receive(NA((PDU)(nil)), NKW(KWCPCISource, tnet.fd.address, KWTestPDUData, pduData)).
			Success("")

		// remote sniffer node
		remoteSniffer, err := NewSnifferStateMachine(testingLogger, "192.168.6.254/24", tnet.vlan6)
		require.NoError(t, err)
		tnet.Append(remoteSniffer)

		// sniffer traffic
		remoteSniffer.GetStartState().Doc("3-2-0").
			Receive(NA((*RegisterForeignDevice)(nil)), NoKWArgs()).Doc("3-2-1").
			Receive(NA((*Result)(nil)), NoKWArgs()).Doc("3-2-2").
			SetEvent("3-registered").Doc("3-2-3").
			Receive(NA((*OriginalUnicastNPDU)(nil)), NoKWArgs()).Doc("3-2-4").
			Success("")

		// run the group
		tnet.Run(0)
	})
	t.Run("test_broadcast", func(t *testing.T) { //Test a broadcast message from TD to IUT.
		ExclusiveGlobalTimeMachine(t)
		testingLogger := testutils.ProduceTestingLogger(t)

		tnet := NewTFNetwork(t)

		//make a PDU from node 1 to node 2
		pduData, err := Xtob("dead.beef")
		require.NoError(t, err)
		pdu := NewPDU(NA(pduData), NKW(KWCPCISource, tnet.fd.address, KWCPCIDestination, NewLocalBroadcast(nil)))
		t.Logf("pdu: %v", pdu)

		// register, wait for ack, send some beef
		tnet.fd.GetStartState().Doc("4-1-0").
			Call(func(args Args, _ KWArgs) error {
				return tnet.fd.bip.Register(args[0].(*Address), args[1].(uint16))
			}, NA(tnet.bbmd.address, uint16(60)), NoKWArgs()).Doc("4-1-1").
			WaitEvent("4-registered", nil).Doc("4-1-2").
			Send(pdu, nil).Doc("4-1-3").
			Success("")

		// the bbmd is happy when it gets the pdu
		tnet.bbmd.GetStartState().
			Receive(NA((PDU)(nil)), NKW(KWCPCISource, tnet.fd.address, KWTestPDUData, pduData)).Doc("4-2-1").
			Success("")

		// home simple node
		homeNode, err := NewBIPSimpleStateMachine(testingLogger, "192.168.5.254/24", tnet.vlan5)
		require.NoError(t, err)

		// home node happy when getting the pdu, broadcast by the bbmd
		homeNode.GetStartState().Doc("4-3-0").
			Receive(NA((PDU)(nil)), NKW(KWCPCISource, tnet.fd.address, KWTestPDUData, pduData)).Doc("4-3-1").
			Success("")

		// remote sniffer node
		remoteSniffer, err := NewSnifferStateMachine(testingLogger, "192.168.6.254/24", tnet.vlan6)
		require.NoError(t, err)
		tnet.Append(remoteSniffer)

		// sniffer traffic
		remoteSniffer.GetStartState().Doc("4-4-0").
			Receive(NA((*RegisterForeignDevice)(nil)), NoKWArgs()).Doc("4-4-1").
			Receive(NA((*Result)(nil)), NoKWArgs()).Doc("4-4-2").
			SetEvent("4-registered").
			Receive(NA((*DistributeBroadcastToNetwork)(nil)), NoKWArgs()).Doc("4-4-3").
			Success("")

		// run the group
		tnet.Run(0)
	})
}
