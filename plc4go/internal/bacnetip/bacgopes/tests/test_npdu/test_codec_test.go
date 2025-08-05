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

package test_npdu

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/npdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/quick"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/state_machine"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/trapped_classes"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

type TestNPDUCodecSuite struct {
	suite.Suite

	client *TrappedClient
	codec  *NPDUCodec
	server *TrappedServer

	log zerolog.Logger
}

func (suite *TestNPDUCodecSuite) SetupTest() {
	suite.log = testutils.ProduceTestingLogger(suite.T())
	var err error
	suite.codec, err = NewNPDUCodec(suite.log)
	suite.Require().NoError(err)
	suite.client, err = NewTrappedClient(suite.log)
	suite.Require().NoError(err)
	suite.server, err = NewTrappedServer(suite.log)
	suite.Require().NoError(err)
	err = Bind(suite.log, suite.client, suite.codec, suite.server)
	suite.Require().NoError(err)
}

// Pass the PDU to the client to send down the stack.
func (suite *TestNPDUCodecSuite) Request(args Args, kwArgs KWArgs) error {
	suite.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Request")

	return suite.client.Request(args, kwArgs)
}

// Check what the server received.
func (suite *TestNPDUCodecSuite) Indication(args Args, kwArgs KWArgs) error {
	suite.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Indication")

	var pduType any
	if len(args) > 0 {
		pduType = args[0].(any)
	}
	pduAttrs := kwArgs
	suite.Assert().True(MatchPdu(suite.log, suite.server.GetIndicationReceived(), pduType, pduAttrs))
	return nil
}

// Check what the server received.
func (suite *TestNPDUCodecSuite) Response(args Args, kwArgs KWArgs) error {
	suite.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Response")

	return suite.server.Response(args, kwArgs)
}

// Check what the server received.
func (suite *TestNPDUCodecSuite) Confirmation(args Args, kwArgs KWArgs) error {
	suite.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Confirmation")

	pduType := args[0].(any)
	pduAttrs := kwArgs
	suite.Assert().True(MatchPdu(suite.log, suite.client.GetConfirmationReceived(), pduType, pduAttrs))
	return nil
}

// TODO: This test is failing in something round 50% of all runs.
/*func (suite *TestNPDUCodecSuite) TestWhoIsRouterToNetwork() { // Test the Result encoding and decoding.
	// Request successful
	pduBytes, err := Xtob(
		"01.80" + // version, network layer message
			"00 0001", // message type and network
	)
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.NPDUParse(testutils.TestContext(suite.T()), pduBytes, uint16(len(pduBytes)))
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NA(quick.WhoIsRouterToNetwork(1)), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NKW(KWTestPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NA(NewPDU(NA(pduBytes), NoKWArgs())), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Confirmation(NA(&WhoIsRouterToNetwork{}), NKW(KWWirtnNetwork, uint16(1)))
}*/

func (suite *TestNPDUCodecSuite) TestIAMRouterToNetworkEmpty() { // Test the Result encoding and decoding.
	// Request successful
	networkList := []uint16{}
	pduBytes, err := Xtob(
		"01.80" + // version, network layer message
			"01", // message type, no network
	)
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.NPDUParse(testutils.TestContext(suite.T()), pduBytes, uint16(len(pduBytes)))
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NA(quick.IAmRouterToNetwork(networkList...)), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NKW(KWTestPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NA(NewPDU(NA(pduBytes), NoKWArgs())), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Confirmation(NA(&IAmRouterToNetwork{}), NKW(KWIartnNetworkList, networkList))
}

func (suite *TestNPDUCodecSuite) TestIAMRouterToNetworks() { // Test the Result encoding and decoding.
	// Request successful
	networkList := []uint16{1, 2, 3}
	pduBytes, err := Xtob(
		"01.80" + // version, network layer message
			"01 0001 0002 0003", // message type and network list
	)
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.NPDUParse(testutils.TestContext(suite.T()), pduBytes, uint16(len(pduBytes)))
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NA(quick.IAmRouterToNetwork(networkList...)), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NKW(KWTestPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NA(NewPDU(NA(pduBytes), NoKWArgs())), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Confirmation(NA(&IAmRouterToNetwork{}), NKW(KWIartnNetworkList, networkList))
}

func (suite *TestNPDUCodecSuite) TestICouldBeRouterToNetworks() { // Test the Result encoding and decoding.
	// Request successful
	pduBytes, err := Xtob(
		"01.80" + // version, network layer message
			"02 0001 02", // message type, network, performance
	)
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.NPDUParse(testutils.TestContext(suite.T()), pduBytes, uint16(len(pduBytes)))
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NA(quick.ICouldBeRouterToNetwork(1, 2)), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NKW(KWTestPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NA(NewPDU(NA(pduBytes), NoKWArgs())), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Confirmation(NA(&ICouldBeRouterToNetwork{}), NKW(KWIcbrtnNetwork, uint16(1), KWIcbrtnPerformanceIndex, uint8(2)))
}

func (suite *TestNPDUCodecSuite) TestRejectMessageToNetwork() { // Test the Result encoding and decoding.
	// Request successful
	pduBytes, err := Xtob(
		"01.80" + // version, network layer message
			"03 01 0002", // message type, reason, performance
	)
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.NPDUParse(testutils.TestContext(suite.T()), pduBytes, uint16(len(pduBytes)))
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NA(quick.RejectMessageToNetwork(1, 2)), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NKW(KWTestPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NA(NewPDU(NA(pduBytes), NoKWArgs())), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Confirmation(NA(&RejectMessageToNetwork{}), NKW(KWRmtnRejectionReason, readWriteModel.NLMRejectMessageToNetworkRejectReason(1), KWRmtnDNET, uint16(2)))
}

func (suite *TestNPDUCodecSuite) TestRouterBusyToNetworkEmpty() { // Test the Result encoding and decoding.
	// Request successful
	networkList := []uint16{}
	pduBytes, err := Xtob(
		"01.80" + // version, network layer message
			"04", // message type, no networks
	)
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.NPDUParse(testutils.TestContext(suite.T()), pduBytes, uint16(len(pduBytes)))
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NA(quick.RouterBusyToNetwork(networkList...)), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NKW(KWTestPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NA(NewPDU(NA(pduBytes), NoKWArgs())), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Confirmation(NA(&RouterBusyToNetwork{}), NKW(KWRbtnNetworkList, networkList))
}

func (suite *TestNPDUCodecSuite) TestRouterBusyToNetworkNetworks() { // Test the Result encoding and decoding.
	// Request successful
	networkList := []uint16{1, 2, 3}
	pduBytes, err := Xtob(
		"01.80" + // version, network layer message
			"04 0001 0002 0003", // message type and network list
	)
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.NPDUParse(testutils.TestContext(suite.T()), pduBytes, uint16(len(pduBytes)))
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NA(quick.RouterBusyToNetwork(networkList...)), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NKW(KWTestPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NA(NewPDU(NA(pduBytes), NoKWArgs())), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Confirmation(NA(&RouterBusyToNetwork{}), NKW(KWRbtnNetworkList, networkList))
}

func (suite *TestNPDUCodecSuite) TestRouterAvailableToNetworkEmpty() { // Test the Result encoding and decoding.
	// Request successful
	networkList := []uint16{}
	pduBytes, err := Xtob(
		"01.80" + // version, network layer message
			"05", // message type, no networks
	)
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.NPDUParse(testutils.TestContext(suite.T()), pduBytes, uint16(len(pduBytes)))
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NA(quick.RouterAvailableToNetwork(networkList...)), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NKW(KWTestPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NA(NewPDU(NA(pduBytes), NoKWArgs())), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Confirmation(NA(&RouterAvailableToNetwork{}), NKW(KWRatnNetworkList, networkList))
}

func (suite *TestNPDUCodecSuite) TestRouterAvailableToNetworkNetworks() { // Test the Result encoding and decoding.
	// Request successful
	networkList := []uint16{1, 2, 3}
	pduBytes, err := Xtob(
		"01.80" + // version, network layer message
			"05 0001 0002 0003", // message type and network list
	)
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.NPDUParse(testutils.TestContext(suite.T()), pduBytes, uint16(len(pduBytes)))
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NA(quick.RouterAvailableToNetwork(networkList...)), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NKW(KWTestPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NA(NewPDU(NA(pduBytes), NoKWArgs())), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Confirmation(NA(&RouterAvailableToNetwork{}), NKW(KWRatnNetworkList, networkList))
}

func (suite *TestNPDUCodecSuite) TestInitializeRoutingTableEmpty() { // Test the Result encoding and decoding.
	// Request successful
	pduBytes, err := Xtob(
		"01.80" + // version, network layer message
			"06 00", // message type and list length
	)
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.NPDUParse(testutils.TestContext(suite.T()), pduBytes, uint16(len(pduBytes)))
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NA(quick.InitializeRoutingTable()), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NKW(KWTestPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NA(NewPDU(NA(pduBytes), NoKWArgs())), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Confirmation(NA(&InitializeRoutingTable{}), NKW(KWIrtTable, []*RoutingTableEntry{}))
}

func (suite *TestNPDUCodecSuite) TestInitializeRoutingTable01() { // Test the Result encoding and decoding.
	// Request successful
	xtob, err := Xtob("")
	suite.Require().NoError(err)
	rte := quick.RoutingTableEntry(1, 2, xtob)
	rtEntries := []*RoutingTableEntry{rte}

	// Request successful
	pduBytes, err := Xtob(
		"01.80" + // version, network layer message
			"06 01" + // message type and list length
			"0001 02 00", // network, port number, port info
	)
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.NPDUParse(testutils.TestContext(suite.T()), pduBytes, uint16(len(pduBytes)))
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NA(quick.InitializeRoutingTable(rtEntries...)), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NKW(KWTestPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NA(NewPDU(NA(pduBytes), NoKWArgs())), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Confirmation(NA(&InitializeRoutingTable{}), NKW(KWIrtTable, rtEntries))
}

func (suite *TestNPDUCodecSuite) TestInitializeRoutingTable02() { // Test the Result encoding and decoding.
	// Request successful
	xtob, err := Xtob("deadbeef")
	suite.Require().NoError(err)
	rte := quick.RoutingTableEntry(3, 4, xtob)
	rtEntries := []*RoutingTableEntry{rte}

	// Request successful
	pduBytes, err := Xtob(
		"01.80" + // version, network layer message
			"06 01" + // message type and list length
			"0003 04 04 DEADBEEF", // network, port number, port info
	)
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.NPDUParse(testutils.TestContext(suite.T()), pduBytes, uint16(len(pduBytes)))
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NA(quick.InitializeRoutingTable(rtEntries...)), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NKW(KWTestPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NA(NewPDU(NA(pduBytes), NoKWArgs())), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Confirmation(NA(&InitializeRoutingTable{}), NKW(KWIrtTable, rtEntries))
}

func (suite *TestNPDUCodecSuite) TestInitializeRoutingTableAck01() { // Test the Result encoding and decoding.
	// Request successful
	xtob, err := Xtob("")
	suite.Require().NoError(err)
	rte := quick.RoutingTableEntry(1, 2, xtob)
	rtEntries := []*RoutingTableEntry{rte}

	// Request successful
	pduBytes, err := Xtob(
		"01.80" + // version, network layer message
			"07 01" + // message type and list length
			"0001 02 00", // network, port number, port info
	)
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.NPDUParse(testutils.TestContext(suite.T()), pduBytes, uint16(len(pduBytes)))
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NA(quick.InitializeRoutingTableAck(rtEntries...)), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NKW(KWTestPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NA(NewPDU(NA(pduBytes), NoKWArgs())), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Confirmation(NA(&InitializeRoutingTableAck{}), NKW(KWIrtaTable, rtEntries))
}

func (suite *TestNPDUCodecSuite) TestInitializeRoutingTableAck02() { // Test the Result encoding and decoding.
	// Request successful
	xtob, err := Xtob("deadbeef")
	suite.Require().NoError(err)
	rte := quick.RoutingTableEntry(3, 4, xtob)
	rtEntries := []*RoutingTableEntry{rte}

	// Request successful
	pduBytes, err := Xtob(
		"01.80" + // version, network layer message
			"07 01" + // message type and list length
			"0003 04 04 DEADBEEF", // network, port number, port info
	)
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.NPDUParse(testutils.TestContext(suite.T()), pduBytes, uint16(len(pduBytes)))
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NA(quick.InitializeRoutingTableAck(rtEntries...)), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NKW(KWTestPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NA(NewPDU(NA(pduBytes), NoKWArgs())), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Confirmation(NA(&InitializeRoutingTableAck{}), NKW(KWIrtaTable, rtEntries))
}

func (suite *TestNPDUCodecSuite) TestEstablishConnectionToNetworks() { // Test the Result encoding and decoding.
	// Request successful
	pduBytes, err := Xtob(
		"01.80" + // version, network layer message
			"08 0005 06", // message type, network, terminationTime
	)
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.NPDUParse(testutils.TestContext(suite.T()), pduBytes, uint16(len(pduBytes)))
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NA(quick.EstablishConnectionToNetwork(5, 6)), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NKW(KWTestPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NA(NewPDU(NA(pduBytes), NoKWArgs())), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Confirmation(NA(&EstablishConnectionToNetwork{}), NKW(KWEctnDNET, uint16(5), KWEctnTerminationTime, uint8(6)))
}

func (suite *TestNPDUCodecSuite) TestDisconnectConnectionToNetwork() { // Test the Result encoding and decoding.
	// Request successful
	pduBytes, err := Xtob(
		"01.80" + // version, network layer message
			"09 0007", // message type, network
	)
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.NPDUParse(testutils.TestContext(suite.T()), pduBytes, uint16(len(pduBytes)))
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NA(quick.DisconnectConnectionToNetwork(7)), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NKW(KWTestPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NA(NewPDU(NA(pduBytes), NoKWArgs())), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Confirmation(NA(&DisconnectConnectionToNetwork{}), NKW(KWDctnDNET, uint16(7)))
}

func (suite *TestNPDUCodecSuite) TestWhatIsNetworkNumber() { // Test the Result encoding and decoding.
	// Request successful
	pduBytes, err := Xtob(
		"01.80" + // version, network layer message
			"12", // message type, network
	)
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.NPDUParse(testutils.TestContext(suite.T()), pduBytes, uint16(len(pduBytes)))
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NA(quick.WhatIsNetworkNumber(0)), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NKW(KWTestPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NA(NewPDU(NA(pduBytes), NoKWArgs())), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Confirmation(NA(&WhatIsNetworkNumber{}), NoKWArgs())
}

func (suite *TestNPDUCodecSuite) TestNetworkNumberIs() { // Test the Result encoding and decoding.
	// Request successful
	pduBytes, err := Xtob(
		"01.80" + // version, network layer message
			"13 0008 01", // message type, network, flag
	)
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.NPDUParse(testutils.TestContext(suite.T()), pduBytes, uint16(len(pduBytes)))
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NA(quick.NetworkNumberIs(8, true)), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NKW(KWTestPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NA(NewPDU(NA(pduBytes), NoKWArgs())), NoKWArgs())
	suite.Assert().NoError(err)
	err = suite.Confirmation(NA(&NetworkNumberIs{}), NKW(KWNniNet, uint16(8), KWNniFlag, true))
}

func TestNPDUCodec(t *testing.T) {
	suite.Run(t, new(TestNPDUCodecSuite))
}
