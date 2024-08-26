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

	"github.com/apache/plc4x/plc4go/internal/bacnetip"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/tests"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

type TBNetwork struct {
	*tests.StateMachineGroup

	trafficLog *tests.TrafficLog
	router     *bacnetip.IPRouter
	vlan       []*bacnetip.IPNetwork

	t *testing.T

	log zerolog.Logger
}

func NewTBNetwork(t *testing.T, count int) *TBNetwork {
	localLog := testutils.ProduceTestingLogger(t)
	tbn := &TBNetwork{
		t:   t,
		log: localLog,
	}
	tbn.StateMachineGroup = tests.NewStateMachineGroup(localLog)

	// reset the time machine
	tests.ResetTimeMachine(tests.StartTime)
	localLog.Trace().Msg("time machine reset")

	// Create a traffic log
	tbn.trafficLog = new(tests.TrafficLog)

	// make a router
	tbn.router = bacnetip.NewIPRouter(localLog)

	// make the networks
	for net := range count {
		// make a network and set the traffic log
		ipNetwork := bacnetip.NewIPNetwork(localLog, bacnetip.WithNetworkName(fmt.Sprintf("192.168.%d.0/24", net+1)), bacnetip.WithNetworkTrafficLogger(tbn.trafficLog))

		// make a router
		routerAddress := Address(fmt.Sprintf("192.168.%d.1/24", net+1))
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
	tests.ExclusiveGlobalTimeMachine(t)
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
	writeBroadcastDistributionTable := WriteBroadcastDistributionTable()
	writeBroadcastDistributionTable.SetPDUDestination(suite.iut.address) // TODO: upstream does this inline
	suite.td.GetStartState().Doc("1-1-0").
		Send(writeBroadcastDistributionTable, nil).Doc("1-1-1").
		Receive(bacnetip.NewArgs((*bacnetip.Result)(nil)), bacnetip.NewKWArgs(bacnetip.KWBvlciResultCode, readWriteModel.BVLCResultCode(0x0010))).Doc("1-1-2").
		Success("")

	// run group
	suite.tnet.Run(0)
}

func (suite *NonBBMDSuite) TestReadBDTFail() {
	// read the broadcast distribution table, get a nack
	readBroadcastDistributionTable := ReadBroadcastDistributionTable()
	readBroadcastDistributionTable.SetPDUDestination(suite.iut.address) // TODO: upstream does this inline
	suite.td.GetStartState().Doc("1-2-0").
		Send(readBroadcastDistributionTable, nil).Doc("1-2-1").
		Receive(bacnetip.NewArgs((*bacnetip.Result)(nil)), bacnetip.NewKWArgs(bacnetip.KWBvlciResultCode, readWriteModel.BVLCResultCode(0x0020))).Doc("1-2-2").
		Success("")

	// run group
	suite.tnet.Run(0)
}

func (suite *NonBBMDSuite) TestRegisterFail() {
	// read the broadcast distribution table, get a nack
	registerForeignDevice := RegisterForeignDevice(10)
	registerForeignDevice.SetPDUDestination(suite.iut.address) // TODO: upstream does this inline
	suite.td.GetStartState().Doc("1-3-0").
		Send(registerForeignDevice, nil).Doc("1-3-1").
		Receive(bacnetip.NewArgs((*bacnetip.Result)(nil)), bacnetip.NewKWArgs(bacnetip.KWBvlciResultCode, readWriteModel.BVLCResultCode(0x0030))).Doc("1-3-2").
		Success("")

	// run group
	suite.tnet.Run(0)
}

func (suite *NonBBMDSuite) TestReadFail() {
	// read the broadcast distribution table, get a nack
	readForeignDeviceTable := ReadForeignDeviceTable()
	readForeignDeviceTable.SetPDUDestination(suite.iut.address) // TODO: upstream does this inline
	suite.td.GetStartState().Doc("1-4-0").
		Send(readForeignDeviceTable, nil).Doc("1-4-1").
		Receive(bacnetip.NewArgs((*bacnetip.Result)(nil)), bacnetip.NewKWArgs(bacnetip.KWBvlciResultCode, readWriteModel.BVLCResultCode(0x0040))).Doc("1-4-2").
		Success("")

	// run group
	suite.tnet.Run(0)
}

func (suite *NonBBMDSuite) TestDeleteFail() {
	// read the broadcast distribution table, get a nack
	deleteForeignDeviceTableEntry := DeleteForeignDeviceTableEntry(Address("1.2.3.4"))
	deleteForeignDeviceTableEntry.SetPDUDestination(suite.iut.address) // TODO: upstream does this inline
	suite.td.GetStartState().Doc("1-5-0").
		Send(deleteForeignDeviceTableEntry, nil).Doc("1-5-1").
		Receive(bacnetip.NewArgs((*bacnetip.Result)(nil)), bacnetip.NewKWArgs(bacnetip.KWBvlciResultCode, readWriteModel.BVLCResultCode(0x0050))).Doc("1-5-2").
		Success("")

	// run group
	suite.tnet.Run(0)
}

func (suite *NonBBMDSuite) TestDistributeFail() {
	// read the broadcast distribution table, get a nack
	pduBytes, err := bacnetip.Xtob(
		// "deadbeef", // forwarded PDU // TODO: this is not a ndpu so we just exploded with that. We use the iartn for that for now
		// TODO: this below is from us as upstream message is not parsable
		"01.80" + // version, network layer message
			"01 0001 0002 0003", // message type and network list
	)
	suite.Require().NoError(err)
	distributeBroadcastToNetwork := DistributeBroadcastToNetwork(pduBytes)
	distributeBroadcastToNetwork.SetPDUDestination(suite.iut.address) // TODO: upstream does this inline
	suite.td.GetStartState().Doc("1-6-0").
		Send(distributeBroadcastToNetwork, nil).Doc("1-6-1").
		Receive(bacnetip.NewArgs((*bacnetip.Result)(nil)), bacnetip.NewKWArgs(bacnetip.KWBvlciResultCode, readWriteModel.BVLCResultCode(0x0060))).Doc("1-6-2").
		Success("")

	// run group
	suite.tnet.Run(0)
}

func TestNonBBMD(t *testing.T) {
	suite.Run(t, new(NonBBMDSuite))
}

type BBMDSuite struct {
	suite.Suite
}

func TestBBMD(t *testing.T) {
	suite.Run(t, new(BBMDSuite))
}
