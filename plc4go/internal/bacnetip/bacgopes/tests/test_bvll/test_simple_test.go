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
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/state_machine"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/time_machine"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/vlan"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

type TNetwork struct {
	*StateMachineGroup
	vlan    *IPNetwork
	td      *BIPSimpleStateMachine
	iut     *BIPSimpleStateMachine
	sniffer *SnifferStateMachine

	t *testing.T

	log zerolog.Logger
}

func NewTNetwork(t *testing.T) *TNetwork {
	localLog := testutils.ProduceTestingLogger(t)
	tn := &TNetwork{
		t:   t,
		log: localLog,
	}
	tn.StateMachineGroup = NewStateMachineGroup(localLog)

	// reset the time machine
	ResetTimeMachine(StartTime)
	localLog.Trace().Msg("time machine reset")

	// make a little LAN
	tn.vlan = NewIPNetwork(localLog)

	// Test devices
	var err error
	tn.td, err = NewBIPSimpleStateMachine(localLog, "192.168.4.1/24", tn.vlan)
	require.NoError(t, err)
	tn.Append(tn.td)

	// implementation under test
	tn.iut, err = NewBIPSimpleStateMachine(localLog, "192.168.4.2/24", tn.vlan)
	require.NoError(t, err)
	tn.Append(tn.iut)

	// sniffer node
	tn.sniffer, err = NewSnifferStateMachine(localLog, "192.168.4.254/24", tn.vlan)
	require.NoError(t, err)
	tn.Append(tn.sniffer)
	return tn
}

func (t *TNetwork) Run(timeLimit time.Duration) {
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
	assert.False(t.t, failed)
}

func TestSimple(t *testing.T) {
	t.Run("test_idle", func(t *testing.T) { //Test an idle network, nothing happens is success.
		ExclusiveGlobalTimeMachine(t)
		tnet := NewTNetwork(t)

		// all start state are successful
		tnet.td.GetStartState().Success("")
		tnet.iut.GetStartState().Success("")
		tnet.sniffer.GetStartState().Success("")

		// run the group
		tnet.Run(0)
	})
	t.Run("test_unicast", func(t *testing.T) { //Test a unicast message from TD to IUT.
		ExclusiveGlobalTimeMachine(t)
		tnet := NewTNetwork(t)

		//make a PDU from node 1 to node 2
		pduData, err := Xtob("dead.beef")
		require.NoError(t, err)
		pdu := NewPDU(NA(pduData), NKW(KWCPCISource, tnet.td.address, KWCPCIDestination, tnet.iut.address))
		t.Logf("pdu: %v", pdu)

		// test device sends it, iut gets it
		tnet.td.GetStartState().Send(pdu, nil).Success("")
		tnet.iut.GetStartState().Receive(NA((PDU)(nil)), NKW(
			KWCPCISource, tnet.td.address,
		)).Success("")

		// sniffer sees message on the wire
		tnet.sniffer.GetStartState().Receive(NA((PDU)(nil)), NKW(
			KWCPCISource, tnet.td.address.AddrTuple,
			KWCPCIDestination, tnet.iut.address.AddrTuple,
			KWTestPDUData, pduData,
		)).Timeout(1.0*time.Millisecond, nil).Success("")

		// run the group
		tnet.Run(0)
	})
	t.Run("test_broadcast", func(t *testing.T) { //Test a broadcast message from TD to IUT.
		ExclusiveGlobalTimeMachine(t)
		tnet := NewTNetwork(t)

		//make a PDU from node 1 to node 2
		pduData, err := Xtob("dead.beef")
		require.NoError(t, err)
		pdu := NewPDU(NA(pduData), NKW(KWCPCISource, tnet.td.address, KWCPCIDestination, NewLocalBroadcast(nil)))
		t.Logf("pdu: %v", pdu)

		// test device sends it, iut gets it
		tnet.td.GetStartState().Send(pdu, nil).Success("")
		tnet.iut.GetStartState().Receive(NA((PDU)(nil)), NKW(
			KWCPCISource, tnet.td.address,
		)).Success("")

		// sniffer sees message on the wire
		tnet.sniffer.GetStartState().Receive(NA((*OriginalBroadcastNPDU)(nil)), NKW(
			KWCPCISource, tnet.td.address.AddrTuple,
			//bacgopes.KWCPCIDestination, tnet.iut.address.AddrTuple,
			KWTestPDUData, pduData,
		)).Timeout(1.0*time.Second, nil).Success("")

		// run the group
		tnet.Run(0)
	})
}
