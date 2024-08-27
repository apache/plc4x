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

	"github.com/apache/plc4x/plc4go/internal/bacnetip"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/constructors"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/tests"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

type TFNetwork struct {
	*tests.StateMachineGroup

	trafficLog *tests.TrafficLog
	router     *bacnetip.IPRouter
	vlan5      *bacnetip.IPNetwork
	vlan6      *bacnetip.IPNetwork
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
	tfn.StateMachineGroup = tests.NewStateMachineGroup(localLog)

	// reset the time machine
	tests.ResetTimeMachine(tests.StartTime)
	localLog.Trace().Msg("time machine reset")

	// Create a traffic log
	tfn.trafficLog = new(tests.TrafficLog)

	// make a router
	tfn.router = bacnetip.NewIPRouter(localLog)

	// make a home LAN
	tfn.vlan5 = bacnetip.NewIPNetwork(localLog, bacnetip.WithNetworkName("192.168.5.0/24"), bacnetip.WithNetworkTrafficLogger(tfn.trafficLog))
	tfn.router.AddNetwork(Address("192.168.5.1/24"), tfn.vlan5)

	// make a remote LAN
	tfn.vlan6 = bacnetip.NewIPNetwork(localLog, bacnetip.WithNetworkName("192.168.6.0/24"), bacnetip.WithNetworkTrafficLogger(tfn.trafficLog))
	tfn.router.AddNetwork(Address("192.168.6.1/24"), tfn.vlan6)

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
	tests.RunTimeMachine(t.log, timeLimit, time.Time{})
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
	t.Skip("something is completely broken at root level... needs more investigation") // TODO: fix issues
	t.Run("test_idle", func(t *testing.T) {                                            //Test an idle network, nothing happens is success.
		tests.ExclusiveGlobalTimeMachine(t)

		tnet := NewTFNetwork(t)

		// all start state are successful
		tnet.fd.GetStartState().Success("")
		tnet.bbmd.GetStartState().Success("")

		// run the group
		tnet.Run(0)
	})
	t.Run("test_registration", func(t *testing.T) {
		tests.ExclusiveGlobalTimeMachine(t)
		testingLogger := testutils.ProduceTestingLogger(t)

		// create a network
		tnet := NewTFNetwork(t)

		// tell the B/IP layer of the foreign device to register
		tnet.fd.GetStartState().
			Call(func(args bacnetip.Args, _ bacnetip.KWArgs) error {
				return tnet.fd.bip.Register(args[0].(*bacnetip.Address), args[1].(int))
			}, bacnetip.NewArgs(tnet.bbmd.address, 30), bacnetip.NoKWArgs).
			Success("")

		// remote sniffer node
		remoteSniffer, err := NewSnifferStateMachine(testingLogger, "192.168.6.254/24", tnet.vlan6)
		require.NoError(t, err)
		tnet.Append(remoteSniffer)

		// sniffer traffic
		remoteSniffer.GetStartState().Doc("1-1-0").
			Receive(bacnetip.NewArgs((*bacnetip.RegisterForeignDevice)(nil)), bacnetip.NoKWArgs).Doc("1-1-1").
			Receive(bacnetip.NewArgs((*bacnetip.Result)(nil)), bacnetip.NoKWArgs).Doc("1-1-2").
			SetEvent("fd-registered").Doc("1-1-3").
			Success("")

		// the bbmd is idle
		tnet.bbmd.GetStartState().Success("")

		// home snooper node
		homeSnooper, err := NewBIPStateMachine(testingLogger, "192.168.5.2/24", tnet.vlan5)
		tnet.Append(homeSnooper)

		// snooper will read foreign device table
		readForeignDeviceTable := ReadForeignDeviceTable() // TODO: upstream sets this as kwargs, check if we really want to propagate this here...
		readForeignDeviceTable.SetPDUDestination(tnet.bbmd.address)
		homeSnooper.GetStartState().Doc("1-2-0").
			WaitEvent("fd-registered", nil).Doc("1-2-1").
			Send(readForeignDeviceTable, nil).Doc("1-2-2").
			Receive(bacnetip.NewArgs((*bacnetip.ReadForeignDeviceTableAck)(nil)), bacnetip.NoKWArgs).Doc("1-2-3").
			Success("")

		// home sniffer node
		homeSniffer, err := NewSnifferStateMachine(testingLogger, "192.168.5.254/24", tnet.vlan5)
		require.NoError(t, err)
		tnet.Append(homeSniffer)

		// sniffer traffic
		homeSniffer.GetStartState().Doc("1-3-0").
			Receive(bacnetip.NewArgs((*bacnetip.RegisterForeignDevice)(nil)), bacnetip.NoKWArgs).Doc("1-3-1").
			Receive(bacnetip.NewArgs((*bacnetip.Result)(nil)), bacnetip.NoKWArgs).Doc("1-3-2").
			Receive(bacnetip.NewArgs((*bacnetip.ReadForeignDeviceTable)(nil)), bacnetip.NoKWArgs).Doc("1-3-3").
			Receive(bacnetip.NewArgs((*bacnetip.ReadForeignDeviceTableAck)(nil)), bacnetip.NoKWArgs).Doc("1-3-4").
			Success("")

		//  run the group
		tnet.Run(0)
	})
	t.Run("test_refresh_registration", func(t *testing.T) {
		tests.ExclusiveGlobalTimeMachine(t)
		testingLogger := testutils.ProduceTestingLogger(t)

		// create a network
		tnet := NewTFNetwork(t)

		// tell the B/IP layer of the foreign device to register
		tnet.fd.GetStartState().
			Call(func(args bacnetip.Args, _ bacnetip.KWArgs) error {
				return tnet.fd.bip.Register(args[0].(*bacnetip.Address), args[1].(int))
			}, bacnetip.NewArgs(tnet.bbmd.address, 10), bacnetip.NoKWArgs).
			Success("")

		// the bbmd is idle
		tnet.bbmd.GetStartState().Success("")

		// remote sniffer node
		remoteSniffer, err := NewSnifferStateMachine(testingLogger, "192.168.6.254/24", tnet.vlan6)
		require.NoError(t, err)
		tnet.Append(remoteSniffer)

		// sniffer traffic
		remoteSniffer.GetStartState().Doc("2-1-0").
			Receive(bacnetip.NewArgs((*bacnetip.RegisterForeignDevice)(nil)), bacnetip.NoKWArgs).Doc("2-1-1").
			Receive(bacnetip.NewArgs((*bacnetip.Result)(nil)), bacnetip.NoKWArgs).Doc("2-1-1").
			Receive(bacnetip.NewArgs((*bacnetip.RegisterForeignDevice)(nil)), bacnetip.NoKWArgs).Doc("2-1-3").
			Receive(bacnetip.NewArgs((*bacnetip.Result)(nil)), bacnetip.NoKWArgs).Doc("2-1-4").
			Success("")

		//  run the group
		tnet.Run(0)
	})
	t.Run("test_unicast", func(t *testing.T) { //Test a unicast message from TD to IUT.
		tests.ExclusiveGlobalTimeMachine(t)
		testingLogger := testutils.ProduceTestingLogger(t)

		tnet := NewTFNetwork(t)

		//make a PDU from node 1 to node 2
		pduData, err := bacnetip.Xtob(
			//"dead.beef", // TODO: upstream is using invalid data to send around, so we just use a IAm
			"01.80" + // version, network layer message
				"02 0001 02", // message type, network, performance
		)
		require.NoError(t, err)
		pdu := bacnetip.NewPDU(bacnetip.NewMessageBridge(pduData...), bacnetip.WithPDUSource(tnet.fd.address), bacnetip.WithPDUDestination(tnet.bbmd.address))
		t.Logf("pdu: %v", pdu)

		// register, wait for ack, send some beef
		tnet.fd.GetStartState().Doc("3-1-0").
			Call(func(args bacnetip.Args, _ bacnetip.KWArgs) error {
				return tnet.fd.bip.Register(args[0].(*bacnetip.Address), args[1].(int))
			}, bacnetip.NewArgs(tnet.bbmd.address, 60), bacnetip.NoKWArgs).Doc("3-1-1").
			WaitEvent("3-registered", nil).Doc("3-1-2").
			Send(pdu, nil).Doc("3-1-3").
			Success("")

		// the bbmd is happy when it gets the pdu
		tnet.bbmd.GetStartState().
			Receive(bacnetip.NewArgs((bacnetip.PDU)(nil)), bacnetip.NewKWArgs(bacnetip.KWPPDUSource, tnet.fd.address, bacnetip.KWPDUData, pduData)).
			Success("")

		// remote sniffer node
		remoteSniffer, err := NewSnifferStateMachine(testingLogger, "192.168.6.254/24", tnet.vlan6)
		require.NoError(t, err)
		tnet.Append(remoteSniffer)

		// sniffer traffic
		remoteSniffer.GetStartState().Doc("3-2-0").
			Receive(bacnetip.NewArgs((*bacnetip.RegisterForeignDevice)(nil)), bacnetip.NoKWArgs).Doc("3-2-1").
			Receive(bacnetip.NewArgs((*bacnetip.Result)(nil)), bacnetip.NoKWArgs).Doc("3-2-2").
			SetEvent("3-registered").Doc("3-2-3").
			Receive(bacnetip.NewArgs((*bacnetip.OriginalUnicastNPDU)(nil)), bacnetip.NoKWArgs).Doc("3-2-4").
			Success("")

		// run the group
		tnet.Run(0)
	})
	t.Run("test_broadcast", func(t *testing.T) { //Test a broadcast message from TD to IUT.
		tests.ExclusiveGlobalTimeMachine(t)
		testingLogger := testutils.ProduceTestingLogger(t)

		tnet := NewTFNetwork(t)

		//make a PDU from node 1 to node 2
		pduData, err := bacnetip.Xtob(
			//"dead.beef", // TODO: upstream is using invalid data to send around, so we just use a IAm
			"01.80" + // version, network layer message
				"02 0001 02", // message type, network, performance
		)
		require.NoError(t, err)
		pdu := bacnetip.NewPDU(bacnetip.NewMessageBridge(pduData...), bacnetip.WithPDUSource(tnet.fd.address), bacnetip.WithPDUDestination(bacnetip.NewLocalBroadcast(nil)))
		t.Logf("pdu: %v", pdu)

		// register, wait for ack, send some beef
		tnet.fd.GetStartState().Doc("4-1-0").
			Call(func(args bacnetip.Args, _ bacnetip.KWArgs) error {
				return tnet.fd.bip.Register(args[0].(*bacnetip.Address), args[1].(int))
			}, bacnetip.NewArgs(tnet.bbmd.address, 60), bacnetip.NoKWArgs).Doc("4-1-1").
			WaitEvent("4-registered", nil).Doc("4-1-2").
			Send(pdu, nil).Doc("4-1-3").
			Success("")

		// the bbmd is happy when it gets the pdu
		tnet.bbmd.GetStartState().
			Receive(bacnetip.NewArgs((bacnetip.PDU)(nil)), bacnetip.NewKWArgs(bacnetip.KWPPDUSource, tnet.fd.address, bacnetip.KWPDUData, pduData)).Doc("4-2-1").
			Success("")

		// home simple node
		homeNode, err := NewBIPSimpleStateMachine(testingLogger, "192.168.5.254/24", tnet.vlan5)
		require.NoError(t, err)

		// home node happy when getting the pdu, broadcast by the bbmd
		homeNode.GetStartState().Doc("4-3-0").
			Receive(bacnetip.NewArgs((bacnetip.PDU)(nil)), bacnetip.NewKWArgs(bacnetip.KWPPDUSource, tnet.fd.address, bacnetip.KWPDUData, pduData)).Doc("4-3-1").
			Success("")

		// remote sniffer node
		remoteSniffer, err := NewSnifferStateMachine(testingLogger, "192.168.6.254/24", tnet.vlan6)
		require.NoError(t, err)
		tnet.Append(remoteSniffer)

		// sniffer traffic
		remoteSniffer.GetStartState().Doc("4-4-0").
			Receive(bacnetip.NewArgs((*bacnetip.RegisterForeignDevice)(nil)), bacnetip.NoKWArgs).Doc("4-4-1").
			Receive(bacnetip.NewArgs((*bacnetip.Result)(nil)), bacnetip.NoKWArgs).Doc("4-4-2").
			SetEvent("4-registered").
			Receive(bacnetip.NewArgs((*bacnetip.DistributeBroadcastToNetwork)(nil)), bacnetip.NoKWArgs).Doc("4-4-3").
			Success("")

		// run the group
		tnet.Run(0)
	})
}
