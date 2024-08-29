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
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/apache/plc4x/plc4go/internal/bacnetip"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/constructors"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/tests"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

type TIPNetwork struct {
	*tests.StateMachineGroup

	vlan *bacnetip.IPNetwork

	t *testing.T

	log zerolog.Logger
}

func NewTIPNetwork(t *testing.T, nodeCount int, addressPattern string, promiscuous bool, spoofing bool) *TIPNetwork {
	localLog := testutils.ProduceTestingLogger(t)
	tn := &TIPNetwork{
		t:   t,
		log: localLog,
	}
	tn.StateMachineGroup = tests.NewStateMachineGroup(localLog)

	// make a little LAN
	tn.vlan = bacnetip.NewIPNetwork(localLog)

	for i := range nodeCount {
		nodeAddress, err := bacnetip.NewAddress(localLog, fmt.Sprintf(addressPattern, i+1))
		require.NoError(t, err)
		node, err := bacnetip.NewIPNode(localLog, nodeAddress, tn.vlan, bacnetip.WithNodePromiscuous(promiscuous), bacnetip.WithNodeSpoofing(spoofing), bacnetip.WithNodeName("node"+strconv.Itoa(i+1)))
		require.NoError(t, err)

		// bind a client state machine to the ndoe
		csm, err := tests.NewClientStateMachine(localLog)
		require.NoError(t, err)

		err = bacnetip.Bind(localLog, csm, node)
		require.NoError(t, err)

		// add it to this group
		tn.Append(csm)
	}

	return tn
}

func (t *TIPNetwork) Run(timeLimit time.Duration) error {
	if timeLimit == 0 {
		timeLimit = 60 * time.Second
	}
	t.log.Debug().Dur("time_limit", timeLimit).Msg("run")

	// reset the time machine
	tests.ResetTimeMachine(tests.StartTime)
	t.log.Trace().Msg("time machine reset")

	// run the group
	if err := t.StateMachineGroup.Run(); err != nil {
		return err
	}

	// run it some time
	tests.RunTimeMachine(t.log, timeLimit, time.Time{})
	t.log.Trace().Msg("time machine finished")

	// check for success
	success, failed := t.CheckForSuccess()
	if !success {
		return errors.New("not all succeeded")
	}
	_ = failed
	return nil
}

func TestIPVLAN(t *testing.T) {
	t.Run("test_idle", func(t *testing.T) { // Test that a very quiet network can exist. This is not a network test so much as a state machine group test
		tests.ExclusiveGlobalTimeMachine(t)

		// two element network
		tnet := NewTIPNetwork(t, 2, "192.168.1.%d/24", false, false)

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
		tests.ExclusiveGlobalTimeMachine(t)

		// two element network
		tnet := NewTIPNetwork(t, 2, "192.168.2.%d/24", false, false)

		stateMachines := tnet.GetStateMachines()
		tnode1, tnode2 := stateMachines[0], stateMachines[1]

		// make a PDU from node 1 to node 2
		pdu := bacnetip.NewPDU(tests.NewDummyMessage([]byte("data")...),
			bacnetip.WithPDUSource(Address("192.168.2.1:47808")),
			bacnetip.WithPDUDestination(Address("192.168.2.2:47808")),
		)
		t.Log(pdu)

		// node 1 sends the pdu, mode 2 gets it
		tnode1.GetStartState().Send(pdu, nil).Success("")
		tnode2.GetStartState().Receive(bacnetip.NewArgs((bacnetip.PDU)(nil)), bacnetip.NewKWArgs(
			bacnetip.KWPPDUSource, AddressTuple("192.168.2.1", 47808),
		)).Success("")

		// run the group
		err := tnet.Run(0)
		assert.NoError(t, err)
	})
	t.Run("test_broadcast", func(t *testing.T) { // Test that a node can send out a 'local broadcast' message which will be received by every other node.
		tests.ExclusiveGlobalTimeMachine(t)

		// three element network
		tnet := NewTIPNetwork(t, 3, "192.168.3.%d/24", false, false)

		stateMachines := tnet.GetStateMachines()
		tnode1, tnode2, tnode3 := stateMachines[0], stateMachines[1], stateMachines[2]

		// make a broadcast PDU
		pdu := bacnetip.NewPDU(tests.NewDummyMessage([]byte("data")...),
			bacnetip.WithPDUSource(Address("192.168.3.1:47808")),
			bacnetip.WithPDUDestination(Address("192.168.3.255:47808")),
		)
		t.Log(pdu)

		// node 1 sends the pdu, node 2 and 3 each get it
		tnode1.GetStartState().Send(pdu, nil).Success("")
		tnode2.GetStartState().Receive(bacnetip.NewArgs(bacnetip.NewPDU(nil)), bacnetip.NewKWArgs(
			bacnetip.KWPPDUSource, AddressTuple("192.168.3.1", 47808),
		)).Success("")
		tnode3.GetStartState().Receive(bacnetip.NewArgs(bacnetip.NewPDU(nil)), bacnetip.NewKWArgs(
			bacnetip.KWPPDUSource, AddressTuple("192.168.3.1", 47808),
		)).Success("")

		// run the group
		err := tnet.Run(0)
		assert.NoError(t, err)
	})
	t.Run("test_spoof_fail", func(t *testing.T) { // Test verifying that a node cannot send out packets with a source address other than its own, see also test_spoof_pass().
		tests.ExclusiveGlobalTimeMachine(t)

		// one element network
		tnet := NewTIPNetwork(t, 1, "192.168.4.%d/24", false, false)

		stateMachines := tnet.GetStateMachines()
		tnode1 := stateMachines[0]

		// make an unicast PDU with the wrong source
		pdu := bacnetip.NewPDU(tests.NewDummyMessage([]byte("data")...),
			bacnetip.WithPDUSource(Address("192.168.4.2:47808")),
			bacnetip.WithPDUDestination(Address("192.168.4.3:47808")),
		)
		t.Log(pdu)

		// when the node attempts to send it raises an error
		tnode1.GetStartState().Send(pdu, nil).Success("")

		// run the group
		err := tnet.Run(0)
		assert.Error(t, err)
	})
	t.Run("test_spoof_pass", func(t *testing.T) { // Test allowing a node to send out packets with a source address other than its own, see also test_spoof_fail().
		tests.ExclusiveGlobalTimeMachine(t)

		// one element network
		tnet := NewTIPNetwork(t, 1, "192.168.5.%d/24", false, true)

		stateMachines := tnet.GetStateMachines()
		tnode1 := stateMachines[0]

		// make an unicast PDU from a fictitious node
		pdu := bacnetip.NewPDU(tests.NewDummyMessage([]byte("data")...),
			bacnetip.WithPDUSource(Address("192.168.5.3:47808")),
			bacnetip.WithPDUDestination(Address("192.168.5.1:47808")),
		)
		t.Log(pdu)

		// node 1 sends the pdu, but gets it back as if it was from node 3
		tnode1.GetStartState().
			Send(pdu, nil).
			Receive(bacnetip.NewArgs(bacnetip.NewPDU(nil)), bacnetip.NewKWArgs(
				bacnetip.KWPPDUSource, AddressTuple("192.168.5.3", 47808),
			)).
			Success("")

		// run the group
		err := tnet.Run(0)
		assert.NoError(t, err)
	})
	t.Run("test_promiscuous_pass", func(t *testing.T) { // Test 'promiscuous mode' of a node which allows it to receive every packet sent on the network.  This is like the network is a hub, or the node is connected to a 'monitor' port on a managed switch.
		testingLogger := testutils.ProduceTestingLogger(t)
		tests.ExclusiveGlobalTimeMachine(t)

		// three element network
		tnet := NewTIPNetwork(t, 3, "192.168.6.%d/24", true, false)

		stateMachines := tnet.GetStateMachines()
		tnode1, tnode2, tnode3 := stateMachines[0], stateMachines[1], stateMachines[2]

		// make a PDU from node 1 to node 2
		src, err := bacnetip.NewAddress(testingLogger, "192.168.6.1:47808")
		require.NoError(t, err)
		dest, err := bacnetip.NewAddress(testingLogger, "192.168.6.2:47808")
		require.NoError(t, err)
		pdu := bacnetip.NewPDU(nil, bacnetip.WithPDUSource(src), bacnetip.WithPDUDestination(dest))
		t.Log(pdu)

		// node 1 sends the pdu, node 2 and 3 each get it
		tnode1.GetStartState().Send(pdu, nil).Success("")
		tnode2.GetStartState().Receive(bacnetip.NewArgs(bacnetip.NewPDU(nil)), bacnetip.NewKWArgs(
			bacnetip.KWPPDUSource, src,
		)).Success("")
		tnode3.GetStartState().Receive(bacnetip.NewArgs(bacnetip.NewPDU(nil)), bacnetip.NewKWArgs(
			bacnetip.KWPDUDestination, dest,
		)).Success("")

		// run the group
		err = tnet.Run(0)
		assert.NoError(t, err)
	})
	t.Run("test_promiscuous_fail", func(t *testing.T) {
		tests.ExclusiveGlobalTimeMachine(t)

		// three element network
		tnet := NewTIPNetwork(t, 3, "192.168.7.%d/24", false, false)

		stateMachines := tnet.GetStateMachines()
		tnode1, tnode2, tnode3 := stateMachines[0], stateMachines[1], stateMachines[2]

		// make a PDU from node 1 to node 2
		pdu := bacnetip.NewPDU(nil, bacnetip.WithPDUSource(Address("192.168.7.1:47808")), bacnetip.WithPDUDestination(Address("192.168.7.2:47808")))
		t.Log(pdu)

		// node 1 sends the pdu to node 2, node 3 waits and gets nothing
		tnode1.GetStartState().Send(pdu, nil).Success("")
		tnode2.GetStartState().Receive(bacnetip.NewArgs(bacnetip.NewPDU(nil)), bacnetip.NewKWArgs(
			bacnetip.KWPPDUSource, AddressTuple("192.168.7.1", 47808),
		)).Success("")

		// if node 3 receives anything it will trigger unexpected receive and fail
		tnode3.GetStartState().Timeout(1*time.Millisecond, nil).Success("")

		// run the group
		err := tnet.Run(0)
		assert.NoError(t, err)
	})
}

type RouterSuite struct {
	suite.Suite

	smg *tests.StateMachineGroup

	log zerolog.Logger
}

func (suite *RouterSuite) SetupTest() {
	t := suite.T()
	suite.log = testutils.ProduceTestingLogger(t)
	tests.ExclusiveGlobalTimeMachine(t)
	// create a state machine group that has all nodes on all networks
	suite.smg = tests.NewStateMachineGroup(suite.log)

	// make some networks
	vlan10 := bacnetip.NewIPNetwork(suite.log)
	vlan20 := bacnetip.NewIPNetwork(suite.log)

	// make a router and add the networks
	trouter := bacnetip.NewIPRouter(suite.log)
	trouter.AddNetwork(Address("192.168.10.1/24"), vlan10)
	trouter.AddNetwork(Address("192.168.20.1/24"), vlan20)

	for pattern, lan := range map[string]*bacnetip.IPNetwork{
		"192.168.10.%d/24": vlan10,
		"192.168.20.%d/24": vlan20,
	} {
		for i := range 2 {
			nodeAddress, err := bacnetip.NewAddress(suite.log, fmt.Sprintf(pattern, i+2))
			suite.NoError(err)
			node, err := bacnetip.NewIPNode(suite.log, nodeAddress, lan)
			suite.NoError(err)
			t.Logf("Node: %v", node)

			// bind a client state machine to the node
			csm, err := tests.NewClientStateMachine(suite.log)
			suite.NoError(err)
			err = bacnetip.Bind(suite.log, csm, node)
			suite.NoError(err)

			// add it to the group
			suite.smg.Append(csm)
		}
	}
}

func (suite *RouterSuite) TearDownTest() {
	// reset the time machine
	tests.ResetTimeMachine(tests.StartTime)
	suite.T().Log("time machine reset")

	// run the group
	err := suite.smg.Run()
	suite.NoError(err)

	// run it for some time
	tests.RunTimeMachine(suite.log, 60*time.Second, time.Time{})
	suite.T().Log("time machine finished")

	// check for success
	success, failed := suite.smg.CheckForSuccess()
	suite.True(success)
	_ = failed
}

func (suite *RouterSuite) TestIdle() {
	// all success
	for _, csm := range suite.smg.GetStateMachines() {
		csm.GetStartState().Success("")
	}
}

func (suite *RouterSuite) TestSendReceive() { // Test that a node can send a message to another node.
	//unpack the state machines
	stateMachines := suite.smg.GetStateMachines()
	csm_10_2, csm_10_3, csm_20_2, csm_20_3 := stateMachines[0], stateMachines[1], stateMachines[2], stateMachines[3]

	// make a PDU from network 10 node 1 to network 20 node 2
	pdu := bacnetip.NewPDU(tests.NewDummyMessage([]byte("data")...),
		bacnetip.WithPDUSource(Address("192.168.10.2:47808")),
		bacnetip.WithPDUDestination(Address("192.168.20.3:47808")))
	suite.T().Log(pdu)

	// node 1 sends the pdu, mode 2 gets it
	csm_10_2.GetStartState().Send(pdu, nil).Success("")
	csm_20_3.GetStartState().Receive(bacnetip.NewArgs(bacnetip.NewPDU(nil)), bacnetip.NewKWArgs(
		bacnetip.KWPPDUSource, AddressTuple("192.168.10.2", 47808),
	)).Success("")

	// other nodes get nothing
	csm_10_3.GetStartState().Timeout(1*time.Second, nil).Success("")
	csm_20_2.GetStartState().Timeout(1*time.Second, nil).Success("")
}

func (suite *RouterSuite) TestLocalBroadcast() { // Test that a node can send a message to all of the other nodes on the same network.
	stateMachines := suite.smg.GetStateMachines()
	csm_10_2, csm_10_3, csm_20_2, csm_20_3 := stateMachines[0], stateMachines[1], stateMachines[2], stateMachines[3]

	// make a PDU from network 10 node 1 to network 20 node 2
	src, err := bacnetip.NewAddress(suite.log, "192.168.10.2:47808")
	suite.Require().NoError(err)
	dest, err := bacnetip.NewAddress(suite.log, "192.168.10.255:47808")
	suite.Require().NoError(err)
	pdu := bacnetip.NewPDU(nil, bacnetip.WithPDUSource(src), bacnetip.WithPDUDestination(dest))
	suite.T().Log(pdu)

	//  node 10-2 sends the pdu, node 10-3 gets pdu, nodes 20-2 and 20-3 dont
	csm_10_2.GetStartState().Send(pdu, nil).Success("")
	csm_10_3.GetStartState().Receive(bacnetip.NewArgs(bacnetip.NewPDU(nil)), bacnetip.NewKWArgs(
		bacnetip.KWPPDUSource, src,
	)).Success("")
	csm_20_3.GetStartState().Timeout(1*time.Second, nil).Success("")
	csm_20_2.GetStartState().Timeout(1*time.Second, nil).Success("")
}

func (suite *RouterSuite) TestRemoteBroadcast() { // Test that a node can send a message to all of the other nodes on a different network.
	t := suite.T()

	stateMachines := suite.smg.GetStateMachines()
	csm_10_2, csm_10_3, csm_20_2, csm_20_3 := stateMachines[0], stateMachines[1], stateMachines[2], stateMachines[3]

	//  make a PDU from network 10 node 1 to network 20 node 2
	src, err := bacnetip.NewAddress(suite.log, "192.168.10.2:47808")
	require.NoError(t, err)
	dest, err := bacnetip.NewAddress(suite.log, "192.168.20.255:47808")
	require.NoError(t, err)
	pdu := bacnetip.NewPDU(nil, bacnetip.WithPDUSource(src), bacnetip.WithPDUDestination(dest))
	t.Log(pdu)

	//  node 10-2 sends the pdu, node 10-3 gets pdu, nodes 20-2 and 20-3 dont
	csm_10_2.GetStartState().Send(pdu, nil).Success("")
	csm_10_3.GetStartState().Timeout(1*time.Second, nil).Success("")
	csm_20_2.GetStartState().Receive(bacnetip.NewArgs(bacnetip.NewPDU(nil)), bacnetip.NewKWArgs(
		bacnetip.KWPPDUSource, src,
	)).Success("")
	csm_20_3.GetStartState().Receive(bacnetip.NewArgs(bacnetip.NewPDU(nil)), bacnetip.NewKWArgs(
		bacnetip.KWPPDUSource, src,
	)).Success("")
}

func TestRouter(t *testing.T) {
	suite.Run(t, new(RouterSuite))
}
