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

	"github.com/apache/plc4x/plc4go/spi/testutils"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/bvll"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/deleteme"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/quick"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/vlan"
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
			t.log.Debug().Str("logEntry", s).Msg("logEntry")
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
				return tnet.fd.bip.Register(args[0].(*Address), args[1].(int))
			}, NewArgs(tnet.bbmd.address, 30), NoKWArgs).
			Success("")

		// remote sniffer node
		remoteSniffer, err := NewSnifferStateMachine(testingLogger, "192.168.6.254/24", tnet.vlan6)
		require.NoError(t, err)
		tnet.Append(remoteSniffer)

		// sniffer traffic
		remoteSniffer.GetStartState().Doc("1-1-0").
			Receive(NewArgs((*RegisterForeignDevice)(nil)), NoKWArgs).Doc("1-1-1").
			Receive(NewArgs((*Result)(nil)), NoKWArgs).Doc("1-1-2").
			SetEvent("fd-registered").Doc("1-1-3").
			Success("")

		// the bbmd is idle
		tnet.bbmd.GetStartState().Success("")

		// home snooper node
		homeSnooper, err := NewBIPStateMachine(testingLogger, "192.168.5.2/24", tnet.vlan5)
		tnet.Append(homeSnooper)

		// snooper will read foreign device table
		readForeignDeviceTable := quick.ReadForeignDeviceTable() // TODO: upstream sets this as kwargs, check if we really want to propagate this here...
		readForeignDeviceTable.SetPDUDestination(tnet.bbmd.address)
		homeSnooper.GetStartState().Doc("1-2-0").
			WaitEvent("fd-registered", nil).Doc("1-2-1").
			Send(readForeignDeviceTable, nil).Doc("1-2-2").
			Receive(NewArgs((*ReadForeignDeviceTableAck)(nil)), NoKWArgs).Doc("1-2-3").
			Success("")

		// home sniffer node
		homeSniffer, err := NewSnifferStateMachine(testingLogger, "192.168.5.254/24", tnet.vlan5)
		require.NoError(t, err)
		tnet.Append(homeSniffer)

		// sniffer traffic
		homeSniffer.GetStartState().Doc("1-3-0").
			Receive(NewArgs((*RegisterForeignDevice)(nil)), NoKWArgs).Doc("1-3-1").
			Receive(NewArgs((*Result)(nil)), NoKWArgs).Doc("1-3-2").
			Receive(NewArgs((*ReadForeignDeviceTable)(nil)), NoKWArgs).Doc("1-3-3").
			Receive(NewArgs((*ReadForeignDeviceTableAck)(nil)), NoKWArgs).Doc("1-3-4").
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
				return tnet.fd.bip.Register(args[0].(*Address), args[1].(int))
			}, NewArgs(tnet.bbmd.address, 10), NoKWArgs).
			Success("")

		// the bbmd is idle
		tnet.bbmd.GetStartState().Success("")

		// remote sniffer node
		remoteSniffer, err := NewSnifferStateMachine(testingLogger, "192.168.6.254/24", tnet.vlan6)
		require.NoError(t, err)
		tnet.Append(remoteSniffer)

		// sniffer traffic
		remoteSniffer.GetStartState().Doc("2-1-0").
			Receive(NewArgs((*RegisterForeignDevice)(nil)), NoKWArgs).Doc("2-1-1").
			Receive(NewArgs((*Result)(nil)), NoKWArgs).Doc("2-1-1").
			Receive(NewArgs((*RegisterForeignDevice)(nil)), NoKWArgs).Doc("2-1-3").
			Receive(NewArgs((*Result)(nil)), NoKWArgs).Doc("2-1-4").
			Success("")

		//  run the group
		tnet.Run(0)
	})
	t.Run("test_unicast", func(t *testing.T) { //Test a unicast message from TD to IUT.
		t.Skip("something is broken with routing...") // TODO: fixme...
		ExclusiveGlobalTimeMachine(t)
		testingLogger := testutils.ProduceTestingLogger(t)

		tnet := NewTFNetwork(t)

		//make a PDU from node 1 to node 2
		pduData, err := Xtob(
			//"dead.beef", // TODO: upstream is using invalid data to send around, so we just use a IAm
			"01.80" + // version, network layer message
				"13 0008 01", // message type, network, flag
		)
		require.NoError(t, err)
		pdu := NewPDU(NewMessageBridge(pduData...), WithPDUSource(tnet.fd.address), WithPDUDestination(tnet.bbmd.address))
		t.Logf("pdu: %v", pdu)

		// register, wait for ack, send some beef
		tnet.fd.GetStartState().Doc("3-1-0").
			Call(func(args Args, _ KWArgs) error {
				return tnet.fd.bip.Register(args[0].(*Address), args[1].(int))
			}, NewArgs(tnet.bbmd.address, 60), NoKWArgs).Doc("3-1-1").
			WaitEvent("3-registered", nil).Doc("3-1-2").
			Send(pdu, nil).Doc("3-1-3").
			Success("")

		// the bbmd is happy when it gets the pdu
		tnet.bbmd.GetStartState().
			Receive(NewArgs((PDU)(nil)), NewKWArgs(KWPPDUSource, tnet.fd.address, KWPDUData, pduData)).
			Success("")

		// remote sniffer node
		remoteSniffer, err := NewSnifferStateMachine(testingLogger, "192.168.6.254/24", tnet.vlan6)
		require.NoError(t, err)
		tnet.Append(remoteSniffer)

		// sniffer traffic
		remoteSniffer.GetStartState().Doc("3-2-0").
			Receive(NewArgs((*RegisterForeignDevice)(nil)), NoKWArgs).Doc("3-2-1").
			Receive(NewArgs((*Result)(nil)), NoKWArgs).Doc("3-2-2").
			SetEvent("3-registered").Doc("3-2-3").
			Receive(NewArgs((*OriginalUnicastNPDU)(nil)), NoKWArgs).Doc("3-2-4").
			Success("")

		// run the group
		tnet.Run(0)
	})
	t.Run("test_broadcast", func(t *testing.T) { //Test a broadcast message from TD to IUT.
		t.Skip("something is broken with routing...") // TODO: fixme...
		ExclusiveGlobalTimeMachine(t)
		testingLogger := testutils.ProduceTestingLogger(t)

		tnet := NewTFNetwork(t)

		//make a PDU from node 1 to node 2
		pduData, err := Xtob(
			//"dead.beef", // TODO: upstream is using invalid data to send around, so we just use a IAm
			"01.80" + // version, network layer message
				"13 0008 01", // message type, network, flag
		)
		require.NoError(t, err)
		pdu := NewPDU(NewMessageBridge(pduData...), WithPDUSource(tnet.fd.address), WithPDUDestination(NewLocalBroadcast(nil)))
		t.Logf("pdu: %v", pdu)

		// register, wait for ack, send some beef
		tnet.fd.GetStartState().Doc("4-1-0").
			Call(func(args Args, _ KWArgs) error {
				return tnet.fd.bip.Register(args[0].(*Address), args[1].(int))
			}, NewArgs(tnet.bbmd.address, 60), NoKWArgs).Doc("4-1-1").
			WaitEvent("4-registered", nil).Doc("4-1-2").
			Send(pdu, nil).Doc("4-1-3").
			Success("")

		// the bbmd is happy when it gets the pdu
		tnet.bbmd.GetStartState().
			Receive(NewArgs((PDU)(nil)), NewKWArgs(KWPPDUSource, tnet.fd.address, KWPDUData, pduData)).Doc("4-2-1").
			Success("")

		// home simple node
		homeNode, err := NewBIPSimpleStateMachine(testingLogger, "192.168.5.254/24", tnet.vlan5)
		require.NoError(t, err)

		// home node happy when getting the pdu, broadcast by the bbmd
		homeNode.GetStartState().Doc("4-3-0").
			Receive(NewArgs((PDU)(nil)), NewKWArgs(KWPPDUSource, tnet.fd.address, KWPDUData, pduData)).Doc("4-3-1").
			Success("")

		// remote sniffer node
		remoteSniffer, err := NewSnifferStateMachine(testingLogger, "192.168.6.254/24", tnet.vlan6)
		require.NoError(t, err)
		tnet.Append(remoteSniffer)

		// sniffer traffic
		remoteSniffer.GetStartState().Doc("4-4-0").
			Receive(NewArgs((*RegisterForeignDevice)(nil)), NoKWArgs).Doc("4-4-1").
			Receive(NewArgs((*Result)(nil)), NoKWArgs).Doc("4-4-2").
			SetEvent("4-registered").
			Receive(NewArgs((*DistributeBroadcastToNetwork)(nil)), NoKWArgs).Doc("4-4-3").
			Success("")

		// run the group
		tnet.Run(0)
	})
}
