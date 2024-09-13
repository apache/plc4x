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
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/apdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/bvll"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/quick"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/state_machine"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/time_machine"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/vlan"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

type TBNetwork struct {
	*StateMachineGroup

	trafficLog *TrafficLog
	router     *IPRouter
	vlan       []*IPNetwork

	t *testing.T

	log zerolog.Logger
}

func NewTBNetwork(t *testing.T, count int) *TBNetwork {
	localLog := testutils.ProduceTestingLogger(t)
	tbn := &TBNetwork{
		t:   t,
		log: localLog,
	}
	tbn.StateMachineGroup = NewStateMachineGroup(localLog)

	// reset the time machine
	ResetTimeMachine(StartTime)
	localLog.Trace().Msg("time machine reset")

	// Create a traffic log
	tbn.trafficLog = new(TrafficLog)

	// make a router
	tbn.router = NewIPRouter(localLog)

	// make the networks
	for net := range count {
		// make a network and set the traffic log
		ipNetwork := NewIPNetwork(localLog, WithNetworkName(fmt.Sprintf("192.168.%d.0/24", net+1)), WithNetworkTrafficLogger(tbn.trafficLog))

		// make a router
		routerAddress := quick.Address(fmt.Sprintf("192.168.%d.1/24", net+1))
		tbn.router.AddNetwork(routerAddress, ipNetwork)

		tbn.vlan = append(tbn.vlan, ipNetwork)
	}

	return tbn
}

func (t *TBNetwork) Run(timeLimit time.Duration) {
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

type NonBBMDSuite struct {
	suite.Suite

	tnet *TBNetwork
	td   *BIPStateMachine
	iut  *BIPSimpleNode

	log zerolog.Logger
}

func (suite *NonBBMDSuite) SetupTest() {
	t := suite.T()
	ExclusiveGlobalTimeMachine(t)
	suite.log = testutils.ProduceTestingLogger(t)

	// create a network
	suite.tnet = NewTBNetwork(t, 1)

	// test device
	var err error
	suite.td, err = NewBIPStateMachine(suite.log, "192.168.1.2/24", suite.tnet.vlan[0])
	suite.Require().NoError(err)
	suite.tnet.Append(suite.td)

	// implementation under test
	suite.iut, err = NewBIPSimpleNode(suite.log, "192.168.1.3/24", suite.tnet.vlan[0])
	suite.Require().NoError(err)
}

func (suite *NonBBMDSuite) TestWriteBDTFail() {
	// read the broadcast distribution table, get a nack
	writeBroadcastDistributionTable := quick.WriteBroadcastDistributionTable()
	writeBroadcastDistributionTable.SetPDUDestination(suite.iut.address) // TODO: upstream does this inline
	suite.td.GetStartState().Doc("1-1-0").
		Send(writeBroadcastDistributionTable, nil).Doc("1-1-1").
		Receive(NA((*Result)(nil)), NKW(KWBvlciResultCode, readWriteModel.BVLCResultCode(0x0010))).Doc("1-1-2").
		Success("")

	// run group
	suite.tnet.Run(0)
}

func (suite *NonBBMDSuite) TestReadBDTFail() {
	// read the broadcast distribution table, get a nack
	readBroadcastDistributionTable := quick.ReadBroadcastDistributionTable()
	readBroadcastDistributionTable.SetPDUDestination(suite.iut.address) // TODO: upstream does this inline
	suite.td.GetStartState().Doc("1-2-0").
		Send(readBroadcastDistributionTable, nil).Doc("1-2-1").
		Receive(NA((*Result)(nil)), NKW(KWBvlciResultCode, readWriteModel.BVLCResultCode(0x0020))).Doc("1-2-2").
		Success("")

	// run group
	suite.tnet.Run(0)
}

func (suite *NonBBMDSuite) TestRegisterFail() {
	// read the broadcast distribution table, get a nack
	registerForeignDevice := quick.RegisterForeignDevice(10)
	registerForeignDevice.SetPDUDestination(suite.iut.address) // TODO: upstream does this inline
	suite.td.GetStartState().Doc("1-3-0").
		Send(registerForeignDevice, nil).Doc("1-3-1").
		Receive(NA((*Result)(nil)), NKW(KWBvlciResultCode, readWriteModel.BVLCResultCode(0x0030))).Doc("1-3-2").
		Success("")

	// run group
	suite.tnet.Run(0)
}

func (suite *NonBBMDSuite) TestReadFail() {
	// read the broadcast distribution table, get a nack
	readForeignDeviceTable := quick.ReadForeignDeviceTable(nil)
	readForeignDeviceTable.SetPDUDestination(suite.iut.address) // TODO: upstream does this inline
	suite.td.GetStartState().Doc("1-4-0").
		Send(readForeignDeviceTable, nil).Doc("1-4-1").
		Receive(NA((*Result)(nil)), NKW(KWBvlciResultCode, readWriteModel.BVLCResultCode(0x0040))).Doc("1-4-2").
		Success("")

	// run group
	suite.tnet.Run(0)
}

func (suite *NonBBMDSuite) TestDeleteFail() {
	// read the broadcast distribution table, get a nack
	deleteForeignDeviceTableEntry := quick.DeleteForeignDeviceTableEntry(quick.Address("1.2.3.4"))
	deleteForeignDeviceTableEntry.SetPDUDestination(suite.iut.address) // TODO: upstream does this inline
	suite.td.GetStartState().Doc("1-5-0").
		Send(deleteForeignDeviceTableEntry, nil).Doc("1-5-1").
		Receive(NA((*Result)(nil)), NKW(KWBvlciResultCode, readWriteModel.BVLCResultCode(0x0050))).Doc("1-5-2").
		Success("")

	// run group
	suite.tnet.Run(0)
}

func (suite *NonBBMDSuite) TestDistributeFail() {
	// read the broadcast distribution table, get a nack
	pduBytes, err := Xtob(
		"deadbeef", // forwarded PDU
	)
	suite.Require().NoError(err)
	distributeBroadcastToNetwork := quick.DistributeBroadcastToNetwork(pduBytes)
	distributeBroadcastToNetwork.SetPDUDestination(suite.iut.address) // TODO: upstream does this inline
	suite.td.GetStartState().Doc("1-6-0").
		Send(distributeBroadcastToNetwork, nil).Doc("1-6-1").
		Receive(NA((*Result)(nil)), NKW(KWBvlciResultCode, readWriteModel.BVLCResultCode(0x0060))).Doc("1-6-2").
		Success("")

	// run group
	suite.tnet.Run(0)
}

func TestNonBBMD(t *testing.T) {
	suite.Run(t, new(NonBBMDSuite))
}

func TestBBMD(t *testing.T) {
	t.Run("test_14_2_1_1", func(t *testing.T) { //14.2.1.1 Execute Forwarded-NPDU (One-hop Distribution).
		t.Skip("needs more work before it can do something") // TODO: implement me
		ExclusiveGlobalTimeMachine(t)
		testLogger := testutils.ProduceTestingLogger(t)

		// Create a network
		tnet := NewTBNetwork(t, 2)

		// implementation under test
		iut, err := NewBIPBBMDApplication(testLogger, "192.168.1.2/24", tnet.vlan[0])
		require.NoError(t, err)
		testLogger.Debug().Stringer("iutbip", iut.bip).Msg("iut.bip")

		// BBMD on net 2
		bbmd1, err := NewBIPBBMDNode(testLogger, "192.168.2.2/24", tnet.vlan[1])
		require.NoError(t, err)

		// add the IUT as a one-hop peer
		err = bbmd1.bip.AddPeer(quick.Address("192.168.1.2/24"))
		require.NoError(t, err)
		testLogger.Debug().Stringer("bbmd1bip", bbmd1.bip).Msg("bbmd1.bip")

		// test device
		td, err := NewBIPSimpleApplicationLayerStateMachine(testLogger, "192.168.2.3/24", tnet.vlan[1])
		require.NoError(t, err)
		tnet.Append(td)

		// listener looks for extra traffic
		listener, err := NewBIPStateMachine(testLogger, "192.168.1.3/24", tnet.vlan[0])
		listener.mux.node.SetPromiscuous(true)
		tnet.Append(listener)

		// broadcast a forwarded NPDU
		td.GetStartState().Doc("2-1-0").
			Send(quick.WhoIsRequest(NKW(KWCPCIDestination, NewLocalBroadcast(nil))), nil).Doc("2-1-1").
			Receive(NA((*IAmRequest)(nil)), NoKWArgs()).Doc("2-1-2").
			Success("")

		// listen for the directed broadcast, then the original unicast,
		// then fail if there's anything else
		listener.GetStartState().Doc("2-2-0").
			Receive(NA((*ForwardedNPDU)(nil)), NoKWArgs()).Doc("2-2-1").
			Receive(NA((*OriginalUnicastNPDU)(nil)), NoKWArgs()).Doc("2-2-2").
			Timeout(3*time.Second, nil).Doc("2-2-3").
			Success("")

		// run the group
		tnet.Run(0)

	})
	t.Run("test_14_2_1_1", func(t *testing.T) { // 14.2.1.1 Execute Forwarded-NPDU (Two-hop Distribution).
		t.Skip("needs more work before it can do something") // TODO: implement me
		ExclusiveGlobalTimeMachine(t)
		testLogger := testutils.ProduceTestingLogger(t)

		// Create a network
		tnet := NewTBNetwork(t, 2)

		// implementation under test
		iut, err := NewBIPBBMDApplication(testLogger, "192.168.1.2/24", tnet.vlan[0])
		require.NoError(t, err)
		testLogger.Debug().Stringer("iutbip", iut.bip).Msg("iut.bip")

		// BBMD on net 2
		bbmd1, err := NewBIPBBMDNode(testLogger, "192.168.2.2/24", tnet.vlan[1])
		require.NoError(t, err)

		// add the IUT as a two-hop peer
		err = bbmd1.bip.AddPeer(quick.Address("192.168.1.2/32"))
		require.NoError(t, err)
		testLogger.Debug().Stringer("bbmd1bip", bbmd1.bip).Msg("bbmd1.bip")

		// test device
		td, err := NewBIPSimpleApplicationLayerStateMachine(testLogger, "192.168.2.3/24", tnet.vlan[1])
		require.NoError(t, err)
		tnet.Append(td)

		// listener looks for extra traffic
		listener, err := NewBIPStateMachine(testLogger, "192.168.1.3/24", tnet.vlan[0])
		listener.mux.node.SetPromiscuous(true)
		tnet.Append(listener)

		// broadcast a forwarded NPDU
		td.GetStartState().Doc("2-3-0").
			Send(quick.WhoIsRequest(NKW(KWCPCIDestination, NewLocalBroadcast(nil))), nil).Doc("2-3-1").
			Receive(NA((*IAmRequest)(nil)), NoKWArgs()).Doc("2-3-2").
			Success("")

		// listen for the forwarded NPDU.  The packet will be sent upstream which
		// will generate the original unicast going back, then it will be
		// re-broadcast on the local LAN.  Fail if there's anything after that.
		s241 := listener.GetStartState().Doc("2-4-0").
			Receive(NA((*ForwardedNPDU)(nil)), NoKWArgs()).Doc("2-4-1")

		// look for the original unicast going back, followed by the rebroadcast
		// of the forwarded NPDU on the local LAN
		both := s241.
			Receive(NA((*OriginalUnicastNPDU)(nil)), NoKWArgs()).Doc("2-4-1-a").
			Receive(NA((*ForwardedNPDU)(nil)), NoKWArgs()).Doc("2-4-1-b")

		// fail if anything is received after both packets
		both.Timeout(3*time.Second, nil).Doc("2-4-4").
			Success("")

		// allow the two packets in either order
		s241.Receive(NA((*ForwardedNPDU)(nil)), NoKWArgs()).Doc("2-4-2-a").
			Receive(NA((*OriginalUnicastNPDU)(nil)), NKW("nextState", both)).Doc("2-4-2-b")

		// run the group
		tnet.Run(0)

	})
}
