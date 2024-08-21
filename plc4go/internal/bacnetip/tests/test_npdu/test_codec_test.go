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

	"github.com/apache/plc4x/plc4go/internal/bacnetip"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/tests"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

func WhoIsRouterToNetwork(net uint16) *bacnetip.WhoIsRouterToNetwork {
	network, err := bacnetip.NewWhoIsRouterToNetwork(bacnetip.WithWhoIsRouterToNetworkNet(net))
	if err != nil {
		panic(err)
	}
	return network
}

func IAmRouterToNetwork(netList ...uint16) *bacnetip.IAmRouterToNetwork {
	network, err := bacnetip.NewIAmRouterToNetwork(bacnetip.WithIAmRouterToNetworkNetworkList(netList...))
	if err != nil {
		panic(err)
	}
	return network
}

func ICouldBeRouterToNetwork(net uint16, perf uint8) *bacnetip.ICouldBeRouterToNetwork {
	network, err := bacnetip.NewICouldBeRouterToNetwork(bacnetip.WithICouldBeRouterToNetworkNetwork(net), bacnetip.WithICouldBeRouterToNetworkPerformanceIndex(perf))
	if err != nil {
		panic(err)
	}
	return network
}

func RejectMessageToNetwork(reason uint8, dnet uint16) *bacnetip.RejectMessageToNetwork {
	network, err := bacnetip.NewRejectMessageToNetwork(bacnetip.WithRejectMessageToNetworkRejectionReason(readWriteModel.NLMRejectMessageToNetworkRejectReason(reason)), bacnetip.WithRejectMessageToNetworkDnet(dnet))
	if err != nil {
		panic(err)
	}
	return network
}

func RouterBusyToNetwork(netList ...uint16) *bacnetip.RouterBusyToNetwork {
	network, err := bacnetip.NewRouterBusyToNetwork(bacnetip.WithRouterBusyToNetworkDnet(netList))
	if err != nil {
		panic(err)
	}
	return network
}

func RouterAvailableToNetwork(netList ...uint16) *bacnetip.RouterAvailableToNetwork {
	network, err := bacnetip.NewRouterAvailableToNetwork(bacnetip.WithRouterAvailableToNetworkDnet(netList))
	if err != nil {
		panic(err)
	}
	return network
}

func InitializeRoutingTable(irtTable ...*bacnetip.RoutingTableEntry) *bacnetip.InitializeRoutingTable {
	network, err := bacnetip.NewInitializeRoutingTable(bacnetip.WithInitializeRoutingTableIrtTable(irtTable...))
	if err != nil {
		panic(err)
	}
	return network
}

func RoutingTableEntry(address uint16, portId uint8, portInfo []byte) *bacnetip.RoutingTableEntry {
	return bacnetip.NewRoutingTableEntry(
		bacnetip.WithRoutingTableEntryDestinationNetworkAddress(address),
		bacnetip.WithRoutingTableEntryPortId(portId),
		bacnetip.WithRoutingTableEntryPortInfo(portInfo),
	)
}

func InitializeRoutingTableAck(irtaTable ...*bacnetip.RoutingTableEntry) *bacnetip.InitializeRoutingTableAck {
	network, err := bacnetip.NewInitializeRoutingTableAck(bacnetip.WithInitializeRoutingTableAckIrtaTable(irtaTable...))
	if err != nil {
		panic(err)
	}
	return network
}

func EstablishConnectionToNetwork(dnet uint16, terminationTime uint8) *bacnetip.EstablishConnectionToNetwork {
	network, err := bacnetip.NewEstablishConnectionToNetwork(bacnetip.WithEstablishConnectionToNetworkDNET(dnet), bacnetip.WithEstablishConnectionToNetworkTerminationTime(terminationTime))
	if err != nil {
		panic(err)
	}
	return network
}

func DisconnectConnectionToNetwork(dnet uint16) *bacnetip.DisconnectConnectionToNetwork {
	network, err := bacnetip.NewDisconnectConnectionToNetwork(bacnetip.WithDisconnectConnectionToNetworkDNET(dnet))
	if err != nil {
		panic(err)
	}
	return network
}

func WhatIsNetworkNumber(dnet uint16) *bacnetip.WhatIsNetworkNumber {
	network, err := bacnetip.NewWhatIsNetworkNumber()
	if err != nil {
		panic(err)
	}
	return network
}

type TestNPDUCodecSuite struct {
	suite.Suite

	client *tests.TrappedClient
	codec  *NPDUCodec
	server *tests.TrappedServer

	log zerolog.Logger
}

func (suite *TestNPDUCodecSuite) SetupSuite() {
	t := suite.T()
	suite.log = testutils.ProduceTestingLogger(t)
}

func (suite *TestNPDUCodecSuite) SetupTest() {
	var err error
	suite.codec, err = NewNPDUCodec(suite.log)
	suite.Require().NoError(err)
	suite.client, err = tests.NewTrappedClient(suite.log)
	suite.Require().NoError(err)
	suite.server, err = tests.NewTrappedServer(suite.log)
	suite.Require().NoError(err)
	err = bacnetip.Bind(suite.log, suite.client, suite.codec, suite.server)
	suite.Require().NoError(err)
}

// Pass the PDU to the client to send down the stack.
func (suite *TestNPDUCodecSuite) Request(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	suite.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Request")

	return suite.client.Request(args, kwargs)
}

// Check what the server received.
func (suite *TestNPDUCodecSuite) Indication(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	suite.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")

	var pduType any
	if len(args) > 0 {
		pduType = args[0].(any)
	}
	pduAttrs := kwargs
	suite.Assert().True(tests.MatchPdu(suite.log, suite.server.GetIndicationReceived(), pduType, pduAttrs))
	return nil
}

// Check what the server received.
func (suite *TestNPDUCodecSuite) Response(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	suite.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Response")

	return suite.server.Response(args, kwargs)
}

// Check what the server received.
func (suite *TestNPDUCodecSuite) Confirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	suite.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Confirmation")

	pduType := args[0].(any)
	pduAttrs := kwargs
	suite.Assert().True(tests.MatchPdu(suite.log, suite.client.GetConfirmationReceived(), pduType, pduAttrs))
	return nil
}

func (suite *TestNPDUCodecSuite) TestWhoIsRouterToNetwork() { // Test the Result encoding and decoding.
	// Request successful
	pduBytes, err := bacnetip.Xtob(
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

	err = suite.Request(bacnetip.NewArgs(WhoIsRouterToNetwork(1)), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacnetip.NoArgs, bacnetip.NewKWArgs(bacnetip.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacnetip.NewArgs(bacnetip.NewPDU(&bacnetip.MessageBridge{Bytes: pduBytes})), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacnetip.NewArgs(&bacnetip.WhoIsRouterToNetwork{}), bacnetip.NewKWArgs(bacnetip.KWWirtnNetwork, uint16(1)))
}

func (suite *TestNPDUCodecSuite) TestIAMRouterToNetworkEmpty() { // Test the Result encoding and decoding.
	// Request successful
	networkList := []uint16{}
	pduBytes, err := bacnetip.Xtob(
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

	err = suite.Request(bacnetip.NewArgs(IAmRouterToNetwork(networkList...)), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacnetip.NoArgs, bacnetip.NewKWArgs(bacnetip.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacnetip.NewArgs(bacnetip.NewPDU(&bacnetip.MessageBridge{Bytes: pduBytes})), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacnetip.NewArgs(&bacnetip.IAmRouterToNetwork{}), bacnetip.NewKWArgs(bacnetip.KWIartnNetworkList, networkList))
}

func (suite *TestNPDUCodecSuite) TestIAMRouterToNetworks() { // Test the Result encoding and decoding.
	// Request successful
	networkList := []uint16{1, 2, 3}
	pduBytes, err := bacnetip.Xtob(
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

	err = suite.Request(bacnetip.NewArgs(IAmRouterToNetwork(networkList...)), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacnetip.NoArgs, bacnetip.NewKWArgs(bacnetip.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacnetip.NewArgs(bacnetip.NewPDU(&bacnetip.MessageBridge{Bytes: pduBytes})), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacnetip.NewArgs(&bacnetip.IAmRouterToNetwork{}), bacnetip.NewKWArgs(bacnetip.KWIartnNetworkList, networkList))
}

func (suite *TestNPDUCodecSuite) TestICouldBeRouterToNetworks() { // Test the Result encoding and decoding.
	// Request successful
	pduBytes, err := bacnetip.Xtob(
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

	err = suite.Request(bacnetip.NewArgs(ICouldBeRouterToNetwork(1, 2)), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacnetip.NoArgs, bacnetip.NewKWArgs(bacnetip.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacnetip.NewArgs(bacnetip.NewPDU(&bacnetip.MessageBridge{Bytes: pduBytes})), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacnetip.NewArgs(&bacnetip.ICouldBeRouterToNetwork{}), bacnetip.NewKWArgs(bacnetip.KWIcbrtnNetwork, uint16(1), bacnetip.KWIcbrtnPerformanceIndex, uint8(2)))
}

func (suite *TestNPDUCodecSuite) TestRejectMessageToNetwork() { // Test the Result encoding and decoding.
	// Request successful
	pduBytes, err := bacnetip.Xtob(
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

	err = suite.Request(bacnetip.NewArgs(RejectMessageToNetwork(1, 2)), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacnetip.NoArgs, bacnetip.NewKWArgs(bacnetip.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacnetip.NewArgs(bacnetip.NewPDU(&bacnetip.MessageBridge{Bytes: pduBytes})), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacnetip.NewArgs(&bacnetip.RejectMessageToNetwork{}), bacnetip.NewKWArgs(bacnetip.KWRmtnRejectionReason, readWriteModel.NLMRejectMessageToNetworkRejectReason(1), bacnetip.KWRmtnDNET, uint16(2)))
}

func (suite *TestNPDUCodecSuite) TestRouterBusyToNetworkEmpty() { // Test the Result encoding and decoding.
	// Request successful
	networkList := []uint16{}
	pduBytes, err := bacnetip.Xtob(
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

	err = suite.Request(bacnetip.NewArgs(RouterBusyToNetwork(networkList...)), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacnetip.NoArgs, bacnetip.NewKWArgs(bacnetip.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacnetip.NewArgs(bacnetip.NewPDU(&bacnetip.MessageBridge{Bytes: pduBytes})), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacnetip.NewArgs(&bacnetip.RouterBusyToNetwork{}), bacnetip.NewKWArgs(bacnetip.KWRbtnNetworkList, networkList))
}

func (suite *TestNPDUCodecSuite) TestRouterBusyToNetworkNetworks() { // Test the Result encoding and decoding.
	// Request successful
	networkList := []uint16{1, 2, 3}
	pduBytes, err := bacnetip.Xtob(
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

	err = suite.Request(bacnetip.NewArgs(RouterBusyToNetwork(networkList...)), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacnetip.NoArgs, bacnetip.NewKWArgs(bacnetip.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacnetip.NewArgs(bacnetip.NewPDU(&bacnetip.MessageBridge{Bytes: pduBytes})), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacnetip.NewArgs(&bacnetip.RouterBusyToNetwork{}), bacnetip.NewKWArgs(bacnetip.KWRbtnNetworkList, networkList))
}

func (suite *TestNPDUCodecSuite) TestRouterAvailableToNetworkEmpty() { // Test the Result encoding and decoding.
	// Request successful
	networkList := []uint16{}
	pduBytes, err := bacnetip.Xtob(
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

	err = suite.Request(bacnetip.NewArgs(RouterAvailableToNetwork(networkList...)), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacnetip.NoArgs, bacnetip.NewKWArgs(bacnetip.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacnetip.NewArgs(bacnetip.NewPDU(&bacnetip.MessageBridge{Bytes: pduBytes})), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacnetip.NewArgs(&bacnetip.RouterAvailableToNetwork{}), bacnetip.NewKWArgs(bacnetip.KWRatnNetworkList, networkList))
}

func (suite *TestNPDUCodecSuite) TestRouterAvailableToNetworkNetworks() { // Test the Result encoding and decoding.
	// Request successful
	networkList := []uint16{1, 2, 3}
	pduBytes, err := bacnetip.Xtob(
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

	err = suite.Request(bacnetip.NewArgs(RouterAvailableToNetwork(networkList...)), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacnetip.NoArgs, bacnetip.NewKWArgs(bacnetip.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacnetip.NewArgs(bacnetip.NewPDU(&bacnetip.MessageBridge{Bytes: pduBytes})), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacnetip.NewArgs(&bacnetip.RouterAvailableToNetwork{}), bacnetip.NewKWArgs(bacnetip.KWRatnNetworkList, networkList))
}

func (suite *TestNPDUCodecSuite) TestInitializeRoutingTableEmpty() { // Test the Result encoding and decoding.
	// Request successful
	pduBytes, err := bacnetip.Xtob(
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

	err = suite.Request(bacnetip.NewArgs(InitializeRoutingTable()), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacnetip.NoArgs, bacnetip.NewKWArgs(bacnetip.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacnetip.NewArgs(bacnetip.NewPDU(&bacnetip.MessageBridge{Bytes: pduBytes})), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacnetip.NewArgs(&bacnetip.InitializeRoutingTable{}), bacnetip.NewKWArgs(bacnetip.KWIrtTable, []*bacnetip.RoutingTableEntry{}))
}

func (suite *TestNPDUCodecSuite) TestInitializeRoutingTable01() { // Test the Result encoding and decoding.
	// Request successful
	xtob, err := bacnetip.Xtob("")
	suite.Require().NoError(err)
	rte := RoutingTableEntry(1, 2, xtob)
	rtEntries := []*bacnetip.RoutingTableEntry{rte}

	// Request successful
	pduBytes, err := bacnetip.Xtob(
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

	err = suite.Request(bacnetip.NewArgs(InitializeRoutingTable(rtEntries...)), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacnetip.NoArgs, bacnetip.NewKWArgs(bacnetip.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacnetip.NewArgs(bacnetip.NewPDU(&bacnetip.MessageBridge{Bytes: pduBytes})), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacnetip.NewArgs(&bacnetip.InitializeRoutingTable{}), bacnetip.NewKWArgs(bacnetip.KWIrtTable, rtEntries))
}

func (suite *TestNPDUCodecSuite) TestInitializeRoutingTable02() { // Test the Result encoding and decoding.
	// Request successful
	xtob, err := bacnetip.Xtob("deadbeef")
	suite.Require().NoError(err)
	rte := RoutingTableEntry(3, 4, xtob)
	rtEntries := []*bacnetip.RoutingTableEntry{rte}

	// Request successful
	pduBytes, err := bacnetip.Xtob(
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

	err = suite.Request(bacnetip.NewArgs(InitializeRoutingTable(rtEntries...)), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacnetip.NoArgs, bacnetip.NewKWArgs(bacnetip.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacnetip.NewArgs(bacnetip.NewPDU(&bacnetip.MessageBridge{Bytes: pduBytes})), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacnetip.NewArgs(&bacnetip.InitializeRoutingTable{}), bacnetip.NewKWArgs(bacnetip.KWIrtTable, rtEntries))
}

func (suite *TestNPDUCodecSuite) TestInitializeRoutingTableAck01() { // Test the Result encoding and decoding.
	// Request successful
	xtob, err := bacnetip.Xtob("")
	suite.Require().NoError(err)
	rte := RoutingTableEntry(1, 2, xtob)
	rtEntries := []*bacnetip.RoutingTableEntry{rte}

	// Request successful
	pduBytes, err := bacnetip.Xtob(
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

	err = suite.Request(bacnetip.NewArgs(InitializeRoutingTableAck(rtEntries...)), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacnetip.NoArgs, bacnetip.NewKWArgs(bacnetip.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacnetip.NewArgs(bacnetip.NewPDU(&bacnetip.MessageBridge{Bytes: pduBytes})), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacnetip.NewArgs(&bacnetip.InitializeRoutingTableAck{}), bacnetip.NewKWArgs(bacnetip.KWIrtaTable, rtEntries))
}

func (suite *TestNPDUCodecSuite) TestInitializeRoutingTableAck02() { // Test the Result encoding and decoding.
	// Request successful
	xtob, err := bacnetip.Xtob("deadbeef")
	suite.Require().NoError(err)
	rte := RoutingTableEntry(3, 4, xtob)
	rtEntries := []*bacnetip.RoutingTableEntry{rte}

	// Request successful
	pduBytes, err := bacnetip.Xtob(
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

	err = suite.Request(bacnetip.NewArgs(InitializeRoutingTableAck(rtEntries...)), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacnetip.NoArgs, bacnetip.NewKWArgs(bacnetip.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacnetip.NewArgs(bacnetip.NewPDU(&bacnetip.MessageBridge{Bytes: pduBytes})), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacnetip.NewArgs(&bacnetip.InitializeRoutingTableAck{}), bacnetip.NewKWArgs(bacnetip.KWIrtaTable, rtEntries))
}

func (suite *TestNPDUCodecSuite) TestEstablishConnectionToNetworks() { // Test the Result encoding and decoding.
	// Request successful
	pduBytes, err := bacnetip.Xtob(
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

	err = suite.Request(bacnetip.NewArgs(EstablishConnectionToNetwork(5, 6)), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacnetip.NoArgs, bacnetip.NewKWArgs(bacnetip.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacnetip.NewArgs(bacnetip.NewPDU(&bacnetip.MessageBridge{Bytes: pduBytes})), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacnetip.NewArgs(&bacnetip.EstablishConnectionToNetwork{}), bacnetip.NewKWArgs(bacnetip.KWEctnDNET, uint16(5), bacnetip.KWEctnTerminationTime, uint8(6)))
}

func (suite *TestNPDUCodecSuite) TestDisconnectConnectionToNetwork() { // Test the Result encoding and decoding.
	// Request successful
	pduBytes, err := bacnetip.Xtob(
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

	err = suite.Request(bacnetip.NewArgs(DisconnectConnectionToNetwork(7)), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacnetip.NoArgs, bacnetip.NewKWArgs(bacnetip.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacnetip.NewArgs(bacnetip.NewPDU(&bacnetip.MessageBridge{Bytes: pduBytes})), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacnetip.NewArgs(&bacnetip.DisconnectConnectionToNetwork{}), bacnetip.NewKWArgs(bacnetip.KWDctnDNET, uint16(7)))
}

func (suite *TestNPDUCodecSuite) TestWhatIsNetworkNumber() { // Test the Result encoding and decoding.
	// Request successful
	pduBytes, err := bacnetip.Xtob(
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

	err = suite.Request(bacnetip.NewArgs(WhatIsNetworkNumber(0)), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacnetip.NoArgs, bacnetip.NewKWArgs(bacnetip.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacnetip.NewArgs(bacnetip.NewPDU(&bacnetip.MessageBridge{Bytes: pduBytes})), bacnetip.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacnetip.NewArgs(&bacnetip.WhatIsNetworkNumber{}), bacnetip.NoKWArgs)
}

func TestNPDUCodec(t *testing.T) {
	suite.Run(t, new(TestNPDUCodecSuite))
}
