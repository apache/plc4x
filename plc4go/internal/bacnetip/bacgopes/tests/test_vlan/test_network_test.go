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

package test_vlan

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/quick"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/state_machine"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/time_machine"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/vlan"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

type TNetwork struct {
	*StateMachineGroup

	vlan *Network

	t *testing.T

	log zerolog.Logger
}

func NewTNetwork(t *testing.T, nodeCount int, promiscuous bool, spoofing bool) *TNetwork {
	localLog := testutils.ProduceTestingLogger(t)
	tn := &TNetwork{
		t:   t,
		log: localLog,
	}
	tn.StateMachineGroup = NewStateMachineGroup(localLog)

	broadcastAddress, err := NewAddress(NA(0))
	require.NoError(t, err)
	// make a little LAN
	tn.vlan = NewNetwork(localLog, WithNetworkBroadcastAddress(broadcastAddress))

	for i := range nodeCount {
		nodeAddress, err := NewAddress(NA(i + 1))
		require.NoError(t, err)
		node, err := NewNode(localLog, nodeAddress, WithNodeLan(tn.vlan), WithNodePromiscuous(promiscuous), WithNodeSpoofing(spoofing))
		require.NoError(t, err)

		// bind a client state machine to the ndoe
		csm, err := NewClientStateMachine(localLog)
		require.NoError(t, err)

		err = Bind(localLog, csm, node)
		require.NoError(t, err)

		// add it to this group
		tn.Append(csm)
	}

	return tn
}

func (t *TNetwork) Run(timeLimit time.Duration) error {
	if timeLimit == 0 {
		timeLimit = 60 * time.Second
	}
	t.log.Debug().Dur("time_limit", timeLimit).Msg("run")

	// reset the time machine
	ResetTimeMachine(StartTime)
	t.log.Trace().Msg("time machine reset")

	// run the group
	if err := t.StateMachineGroup.Run(); err != nil {
		return err
	}

	// run it some time
	RunTimeMachine(t.log, timeLimit, time.Time{})
	t.log.Trace().Msg("time machine finished")

	// check for success
	success, failed := t.CheckForSuccess()
	if !success {
		return errors.New("not all succeeded")
	}
	if failed {
		return errors.New("some failed")
	}
	return nil
}

func TestVLAN(t *testing.T) {
	t.Run("test_idle", func(t *testing.T) { // Test that a very quiet network can exist. This is not a network test so much as a state machine group test
		ExclusiveGlobalTimeMachine(t)

		// two element network
		tnet := NewTNetwork(t, 2, false, false)

		stateMachines := tnet.GetStateMachines()
		tnode1, tnode2 := stateMachines[0], stateMachines[1]

		// set the start states of both machines to success
		tnode1.GetStartState().Success("")
		tnode2.GetStartState().Success("")

		// run the group
		err := tnet.Run(0)
		assert.NoError(t, err)
	})
	t.Run("test_send_receive", func(t *testing.T) { // Test that a node can send a message to another node.
		ExclusiveGlobalTimeMachine(t)

		// two element network
		tnet := NewTNetwork(t, 2, false, false)

		stateMachines := tnet.GetStateMachines()
		tnode1, tnode2 := stateMachines[0], stateMachines[1]

		// make a PDU from node 1 to node 2
		src, err := NewAddress(NA(1))
		require.NoError(t, err)
		dest, err := NewAddress(NA(2))
		require.NoError(t, err)
		pdu := NewPDU(NoArgs, NKW(KWCPCISource, src, KWCPCIDestination, dest))
		t.Log(pdu)

		// node 1 sends the pdu, mode 2 gets it
		tnode1.GetStartState().Send(pdu, nil).Success("")
		tnode2.GetStartState().Receive(NA(NewPDU(Nothing())), NKW(
			KWCPCISource, src,
		)).Success("")

		// run the group
		err = tnet.Run(0)
		assert.NoError(t, err)
	})
	t.Run("test_broadcast", func(t *testing.T) { // Test that a node can send out a 'local broadcast' message which will be received by every other node.
		ExclusiveGlobalTimeMachine(t)

		// three element network
		tnet := NewTNetwork(t, 3, false, false)

		stateMachines := tnet.GetStateMachines()
		tnode1, tnode2, tnode3 := stateMachines[0], stateMachines[1], stateMachines[2]

		// make a PDU from node 1 to node 2
		src, err := NewAddress(NA(1))
		require.NoError(t, err)
		dest, err := NewAddress(NA(0))
		require.NoError(t, err)
		pdu := NewPDU(NoArgs, NKW(KWCPCISource, src, KWCPCIDestination, dest))
		t.Log(pdu)

		// node 1 sends the pdu, node 2 and 3 each get it
		tnode1.GetStartState().Send(pdu, nil).Success("")
		tnode2.GetStartState().Receive(NA(NewPDU(Nothing())), NKW(
			KWCPCISource, src,
		)).Success("")
		tnode3.GetStartState().Receive(NA(NewPDU(Nothing())), NKW(
			KWCPCISource, src,
		)).Success("")

		// run the group
		err = tnet.Run(0)
		assert.NoError(t, err)
	})
	t.Run("test_spoof_fail", func(t *testing.T) { // Test verifying that a node cannot send out packets with a source address other than its own, see also test_spoof_pass().
		ExclusiveGlobalTimeMachine(t)

		// one element network
		tnet := NewTNetwork(t, 1, false, false)

		stateMachines := tnet.GetStateMachines()
		tnode1 := stateMachines[0]

		// make an unicast PDU with the wrong source
		src, err := NewAddress(NA(2))
		require.NoError(t, err)
		dest, err := NewAddress(NA(3))
		require.NoError(t, err)
		pdu := NewPDU(NoArgs, NKW(KWCPCISource, src, KWCPCIDestination, dest))
		t.Log(pdu)

		// node 1 sends the pdu, node 2 and 3 each get it
		tnode1.GetStartState().Send(pdu, nil).Success("")

		// run the group
		err = tnet.Run(0)
		assert.Error(t, err)
	})
	t.Run("test_spoof_pass", func(t *testing.T) { // Test allowing a node to send out packets with a source address other than its own, see also test_spoof_fail().
		ExclusiveGlobalTimeMachine(t)

		// one element network
		tnet := NewTNetwork(t, 1, false, true)

		stateMachines := tnet.GetStateMachines()
		tnode1 := stateMachines[0]

		// make an unicast PDU with the wrong source
		src, err := NewAddress(NA(3))
		require.NoError(t, err)
		dest, err := NewAddress(NA(1))
		require.NoError(t, err)
		pdu := NewPDU(NoArgs, NKW(KWCPCISource, src, KWCPCIDestination, dest))
		t.Log(pdu)

		// node 1 sends the pdu, but gets it back as if it was from node 3
		tnode1.GetStartState().
			Send(pdu, nil).
			Receive(NA(NewPDU(Nothing())), NKW(
				KWCPCISource, src,
			)).
			Success("")

		// run the group
		err = tnet.Run(0)
		assert.NoError(t, err)
	})
	t.Run("test_promiscuous_pass", func(t *testing.T) { // Test 'promiscuous mode' of a node which allows it to receive every packet sent on the network.  This is like the network is a hub, or the node is connected to a 'monitor' port on a managed switch.
		ExclusiveGlobalTimeMachine(t)

		// three element network
		tnet := NewTNetwork(t, 3, true, false)

		stateMachines := tnet.GetStateMachines()
		tnode1, tnode2, tnode3 := stateMachines[0], stateMachines[1], stateMachines[2]

		// make a PDU from node 1 to node 2
		src, err := NewAddress(NA(1))
		require.NoError(t, err)
		dest, err := NewAddress(NA(2))
		require.NoError(t, err)
		pdu := NewPDU(NoArgs, NKW(KWCPCISource, src, KWCPCIDestination, dest))
		t.Log(pdu)

		// node 1 sends the pdu, node 2 and 3 each get it
		tnode1.GetStartState().Send(pdu, nil).Success("")
		tnode2.GetStartState().Receive(NA(NewPDU(Nothing())), NKW(
			KWCPCISource, src,
		)).Success("")
		tnode3.GetStartState().Receive(NA(NewPDU(Nothing())), NKW(
			KWCPCIDestination, dest,
		)).Success("")

		// run the group
		err = tnet.Run(0)
		assert.NoError(t, err)
	})
	t.Run("test_promiscuous_fail", func(t *testing.T) {
		ExclusiveGlobalTimeMachine(t)

		// three element network
		tnet := NewTNetwork(t, 3, false, false)

		stateMachines := tnet.GetStateMachines()
		tnode1, tnode2, tnode3 := stateMachines[0], stateMachines[1], stateMachines[2]

		// make a PDU from node 1 to node 2
		pdu := NewPDU(NoArgs, NKW(KWCPCISource, quick.Address(1), KWCPCIDestination, quick.Address(1)))
		t.Log(pdu)

		// node 1 sends the pdu to node 2, node 3 waits and gets nothing
		tnode1.GetStartState().Send(pdu, nil).Success("")
		tnode2.GetStartState().Receive(NA(NewPDU(Nothing())), NKW(
			KWCPCISource, quick.Address(1),
		)).Success("")

		// if node 3 receives anything it will trigger unexpected receive and fail
		tnode3.GetStartState().Timeout(500*time.Millisecond, nil).Success("")

		// run the group
		err := tnet.Run(0)
		assert.Error(t, err)
	})
}

func TestVLANEvents(t *testing.T) {
	t.Run("test_send_receive", func(t *testing.T) { // Test that a node can send a message to another node and use events to continue with the messages.
		ExclusiveGlobalTimeMachine(t)

		// two element network
		tnet := NewTNetwork(t, 2, false, false)

		stateMachines := tnet.GetStateMachines()
		tnode1, tnode2 := stateMachines[0], stateMachines[1]

		// make a PDU from node 1 to node 2
		src, err := NewAddress(NA(1))
		require.NoError(t, err)
		dest, err := NewAddress(NA(2))
		require.NoError(t, err)

		deadPDU := NewPDU(NA([]byte{0xde, 0xad}), NKW(KWCPCISource, src, KWCPCIDestination, dest))
		t.Log(deadPDU)

		// make a PDU from node 1 to node 2
		beefPDU := NewPDU(NA([]byte{0xbe, 0xef}), NKW(KWCPCISource, src, KWCPCIDestination, dest))
		t.Log(beefPDU)

		//  node 1 sends dead_pdu, waits for event, sends beef_pdu
		tnode1.GetStartState().
			Send(deadPDU, nil).WaitEvent("e", nil).
			Send(beefPDU, nil).Success("")

		// node 2 receives dead_pdu, sets event, waits for beef_pdu
		tnode2.GetStartState().
			Receive(NA(NewPDU(Nothing())), NKW(
				KWTestPDUData, []byte{0xde, 0xad},
			)).SetEvent("e").
			Receive(NA(NewPDU(Nothing())), NKW(
				KWTestPDUData, []byte{0xbe, 0xef},
			)).Success("")

		// run the group
		err = tnet.Run(0)
		assert.NoError(t, err)
	})
}
