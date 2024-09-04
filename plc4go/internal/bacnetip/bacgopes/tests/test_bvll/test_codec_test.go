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

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"

	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/constructors"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

type TestAnnexJCodecSuite struct {
	suite.Suite

	client *tests.TrappedClient
	codec  *bacgopes.AnnexJCodec
	server *tests.TrappedServer

	log zerolog.Logger
}

func (suite *TestAnnexJCodecSuite) SetupTest() {
	suite.log = testutils.ProduceTestingLogger(suite.T())
	// minature trapped stack
	var err error
	suite.codec, err = bacgopes.NewAnnexJCodec(suite.log)
	suite.Require().NoError(err)
	suite.client, err = tests.NewTrappedClient(suite.log)
	suite.Require().NoError(err)
	suite.server, err = tests.NewTrappedServer(suite.log)
	suite.Require().NoError(err)
	err = bacgopes.Bind(suite.log, suite.client, suite.codec, suite.server)
	suite.Require().NoError(err)
}

// Pass the PDU to the client to send down the stack.
func (suite *TestAnnexJCodecSuite) Request(args bacgopes.Args, kwargs bacgopes.KWArgs) error {
	suite.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Request")

	return suite.client.Request(args, kwargs)
}

// Check what the server received.
func (suite *TestAnnexJCodecSuite) Indication(args bacgopes.Args, kwargs bacgopes.KWArgs) error {
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
func (suite *TestAnnexJCodecSuite) Response(args bacgopes.Args, kwargs bacgopes.KWArgs) error {
	suite.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Response")

	return suite.server.Response(args, kwargs)
}

// Check what the server received.
func (suite *TestAnnexJCodecSuite) Confirmation(args bacgopes.Args, kwargs bacgopes.KWArgs) error {
	suite.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Confirmation")

	pduType := args[0].(any)
	pduAttrs := kwargs
	suite.Assert().True(tests.MatchPdu(suite.log, suite.client.GetConfirmationReceived(), pduType, pduAttrs))
	return nil
}

func (suite *TestAnnexJCodecSuite) TestResult() {
	// Request successful
	pduBytes, err := bacgopes.Xtob("81.00.0006.0000")
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(bacgopes.NewArgs(Result(0)), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacgopes.NoArgs, bacgopes.NewKWArgs(bacgopes.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacgopes.NewArgs(bacgopes.NewPDU(bacgopes.NewMessageBridge(pduBytes...))), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacgopes.NewArgs((*bacgopes.Result)(nil)), bacgopes.NewKWArgs(bacgopes.KWBvlciResultCode, readWriteModel.BVLCResultCode(0)))

	// Request error condition
	pduBytes, err = bacgopes.Xtob("81.00.0006.0010") // TODO: check if this is right or if it should be 01 as there is no code for 1
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(bacgopes.NewArgs(Result(0x0010)), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacgopes.NoArgs, bacgopes.NewKWArgs(bacgopes.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacgopes.NewArgs(bacgopes.NewPDU(bacgopes.NewMessageBridge(pduBytes...))), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacgopes.NewArgs((*bacgopes.Result)(nil)), bacgopes.NewKWArgs(bacgopes.KWBvlciResultCode, readWriteModel.BVLCResultCode(0x0010)))
}

func (suite *TestAnnexJCodecSuite) TestWriteBroadcastDistributionTable() {
	// write an empty table
	pduBytes, err := bacgopes.Xtob("81.01.0004")
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(bacgopes.NewArgs(WriteBroadcastDistributionTable()), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacgopes.NoArgs, bacgopes.NewKWArgs(bacgopes.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacgopes.NewArgs(bacgopes.NewPDU(bacgopes.NewMessageBridge(pduBytes...))), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacgopes.NewArgs((*bacgopes.WriteBroadcastDistributionTable)(nil)), bacgopes.NewKWArgs(bacgopes.KWBvlciBDT, []*bacgopes.Address{}))

	// write table with an element
	addr, _ := bacgopes.NewAddress(zerolog.Nop(), "192.168.0.254/24")
	pduBytes, err = bacgopes.Xtob("81.01.000e" +
		"c0.a8.00.fe.ba.c0 ff.ff.ff.00") // address and mask
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(bacgopes.NewArgs(WriteBroadcastDistributionTable(addr)), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacgopes.NoArgs, bacgopes.NewKWArgs(bacgopes.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacgopes.NewArgs(bacgopes.NewPDU(bacgopes.NewMessageBridge(pduBytes...))), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacgopes.NewArgs((*bacgopes.WriteBroadcastDistributionTable)(nil)), bacgopes.NewKWArgs(bacgopes.KWBvlciBDT, []*bacgopes.Address{addr}))
}

func (suite *TestAnnexJCodecSuite) TestReadBroadcastDistributionTable() {
	// Read an empty table
	pduBytes, err := bacgopes.Xtob("81.02.0004")
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(bacgopes.NewArgs(ReadBroadcastDistributionTable()), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacgopes.NoArgs, bacgopes.NewKWArgs(bacgopes.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacgopes.NewArgs(bacgopes.NewPDU(bacgopes.NewMessageBridge(pduBytes...))), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacgopes.NewArgs((*bacgopes.ReadBroadcastDistributionTable)(nil)), bacgopes.NoKWArgs)
}

func (suite *TestAnnexJCodecSuite) TestReadBroadcastDistributionTableAck() {
	// Read an empty TableAck
	pduBytes, err := bacgopes.Xtob("81.03.0004")
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(bacgopes.NewArgs(ReadBroadcastDistributionTableAck()), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacgopes.NoArgs, bacgopes.NewKWArgs(bacgopes.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacgopes.NewArgs(bacgopes.NewPDU(bacgopes.NewMessageBridge(pduBytes...))), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacgopes.NewArgs((*bacgopes.ReadBroadcastDistributionTableAck)(nil)), bacgopes.NewKWArgs(bacgopes.KWBvlciBDT, []*bacgopes.Address{}))

	// Read TableAck with an element
	addr, _ := bacgopes.NewAddress(zerolog.Nop(), "192.168.0.254/24")
	pduBytes, err = bacgopes.Xtob("81.03.000e" + //bvlci
		"c0.a8.00.fe.ba.c0 ff.ff.ff.00") // address and mask
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(bacgopes.NewArgs(ReadBroadcastDistributionTableAck(addr)), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacgopes.NoArgs, bacgopes.NewKWArgs(bacgopes.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacgopes.NewArgs(bacgopes.NewPDU(bacgopes.NewMessageBridge(pduBytes...))), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacgopes.NewArgs((*bacgopes.ReadBroadcastDistributionTableAck)(nil)), bacgopes.NewKWArgs(bacgopes.KWBvlciBDT, []*bacgopes.Address{addr}))
}

func (suite *TestAnnexJCodecSuite) TestForwardNPDU() {
	addr, err := bacgopes.NewAddress(zerolog.Nop(), "192.168.0.1")
	xpdu, err := bacgopes.Xtob(
		// "deadbeef", // forwarded PDU // TODO: this is not a ndpu so we just exploded with that. We use the iartn for that for now
		// TODO: this below is from us as upstream message is not parsable
		"01.80" + // version, network layer message
			"01 0001 0002 0003", // message type and network list
	)
	suite.Require().NoError(err)
	pduBytes, err := bacgopes.Xtob("81.04.0013" + //   bvlci // TODO: length was 0e before
		"c0.a8.00.01.ba.c0" + // original source address
		// "deadbeef", // forwarded PDU // TODO: this is not a ndpu so we just exploded with that. We use the iartn for that for now
		// TODO: this below is from us as upstream message is not parsable
		"01.80" + // version, network layer message
		"01 0001 0002 0003", // message type and network list
	)
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(bacgopes.NewArgs(ForwardedNPDU(addr, xpdu)), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacgopes.NoArgs, bacgopes.NewKWArgs(bacgopes.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacgopes.NewArgs(bacgopes.NewPDU(bacgopes.NewMessageBridge(pduBytes...))), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacgopes.NewArgs((*bacgopes.ForwardedNPDU)(nil)), bacgopes.NewKWArgs(bacgopes.KWBvlciAddress, addr, bacgopes.KWPDUData, xpdu))
	suite.Assert().NoError(err)
}

func (suite *TestAnnexJCodecSuite) TestRegisterForeignDevice() {
	// Request successful
	pduBytes, err := bacgopes.Xtob(
		"81.05.0006" + // bvlci
			"001e", //time-to-live
	)
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(bacgopes.NewArgs(RegisterForeignDevice(30)), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacgopes.NoArgs, bacgopes.NewKWArgs(bacgopes.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacgopes.NewArgs(bacgopes.NewPDU(bacgopes.NewMessageBridge(pduBytes...))), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacgopes.NewArgs((*bacgopes.RegisterForeignDevice)(nil)), bacgopes.NewKWArgs(bacgopes.KWBvlciTimeToLive, uint16(30)))
}

func (suite *TestAnnexJCodecSuite) TestReadForeignDeviceTable() {
	// Read an empty table
	pduBytes, err := bacgopes.Xtob("81.06.0004")
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(bacgopes.NewArgs(ReadForeignDeviceTable()), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacgopes.NoArgs, bacgopes.NewKWArgs(bacgopes.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacgopes.NewArgs(bacgopes.NewPDU(bacgopes.NewMessageBridge(pduBytes...))), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacgopes.NewArgs((*bacgopes.ReadForeignDeviceTable)(nil)), bacgopes.NoKWArgs)
}

func (suite *TestAnnexJCodecSuite) TestReadForeignDeviceTableAck() {
	// Read an empty TableAck
	pduBytes, err := bacgopes.Xtob("81.07.0004")
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(bacgopes.NewArgs(ReadForeignDeviceTableAck()), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacgopes.NoArgs, bacgopes.NewKWArgs(bacgopes.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacgopes.NewArgs(bacgopes.NewPDU(bacgopes.NewMessageBridge(pduBytes...))), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacgopes.NewArgs((*bacgopes.ReadForeignDeviceTableAck)(nil)), bacgopes.NewKWArgs(bacgopes.KWBvlciFDT, []*bacgopes.FDTEntry{}))

	// Read TableAck with one entry
	fdte := FDTEntry()
	fdte.FDAddress, err = bacgopes.NewAddress(suite.log, "192.168.0.10")
	suite.Require().NoError(err)
	fdte.FDTTL = 30
	fdte.FDRemain = 15
	pduBytes, err = bacgopes.Xtob(
		"81.07.000e" + //bvlci
			"c0.a8.00.0a.ba.c0" + // address
			"001e.000f", // ttl and remaining
	)
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(bacgopes.NewArgs(ReadForeignDeviceTableAck(fdte)), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacgopes.NoArgs, bacgopes.NewKWArgs(bacgopes.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacgopes.NewArgs(bacgopes.NewPDU(bacgopes.NewMessageBridge(pduBytes...))), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacgopes.NewArgs((*bacgopes.ReadForeignDeviceTableAck)(nil)), bacgopes.NewKWArgs(bacgopes.KWBvlciFDT, []*bacgopes.FDTEntry{fdte}))
}

func (suite *TestAnnexJCodecSuite) TestDeleteForeignDeviceTableEntry() {
	addr, _ := bacgopes.NewAddress(zerolog.Nop(), "192.168.0.11/24")
	pduBytes, err := bacgopes.Xtob("81.08.000a" + // bvlci
		"c0.a8.00.0b.ba.c0") // address of entry to be deleted
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(bacgopes.NewArgs(DeleteForeignDeviceTableEntry(addr)), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacgopes.NoArgs, bacgopes.NewKWArgs(bacgopes.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacgopes.NewArgs(bacgopes.NewPDU(bacgopes.NewMessageBridge(pduBytes...))), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacgopes.NewArgs((*bacgopes.DeleteForeignDeviceTableEntry)(nil)), bacgopes.NewKWArgs(bacgopes.KWBvlciAddress, addr))
}

func (suite *TestAnnexJCodecSuite) TestDeleteForeignDeviceTableAck() {
	// TODO: implement me
	suite.T().Skip()
}

func (suite *TestAnnexJCodecSuite) TestDistributeBroadcastToNetwork() {
	xpdu, err := bacgopes.Xtob(
		// "deadbeef", // forwarded PDU // TODO: this is not a ndpu so we just exploded with that. We use the iartn for that for now
		// TODO: this below is from us as upstream message is not parsable
		"01.80" + // version, network layer message
			"01 0001 0002 0003", // message type and network list
	)
	suite.Require().NoError(err)
	pduBytes, err := bacgopes.Xtob("81.09.000d" + //   bvlci // TODO: length was 08 before
		// "deadbeef", // forwarded PDU // TODO: this is not a ndpu so we just exploded with that. We use the iartn for that for now
		// TODO: this below is from us as upstream message is not parsable
		"01.80" + // version, network layer message
		"01 0001 0002 0003", // message type and network list
	)
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(bacgopes.NewArgs(DistributeBroadcastToNetwork(xpdu)), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacgopes.NoArgs, bacgopes.NewKWArgs(bacgopes.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacgopes.NewArgs(bacgopes.NewPDU(bacgopes.NewMessageBridge(pduBytes...))), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacgopes.NewArgs((*bacgopes.DistributeBroadcastToNetwork)(nil)), bacgopes.NewKWArgs(bacgopes.KWPDUData, xpdu))
}

func (suite *TestAnnexJCodecSuite) TestOriginalUnicastNPDU() {
	xpdu, err := bacgopes.Xtob(
		// "deadbeef", // forwarded PDU // TODO: this is not a ndpu so we just exploded with that. We use the iartn for that for now
		// TODO: this below is from us as upstream message is not parsable
		"01.80" + // version, network layer message
			"01 0001 0002 0003", // message type and network list
	)
	suite.Require().NoError(err)
	pduBytes, err := bacgopes.Xtob("81.0a.000d" + //   bvlci // TODO: length was 08 before
		// "deadbeef", // forwarded PDU // TODO: this is not a ndpu so we just exploded with that. We use the iartn for that for now
		// TODO: this below is from us as upstream message is not parsable
		"01.80" + // version, network layer message
		"01 0001 0002 0003", // message type and network list
	)
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(bacgopes.NewArgs(OriginalUnicastNPDU(xpdu)), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacgopes.NoArgs, bacgopes.NewKWArgs(bacgopes.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacgopes.NewArgs(bacgopes.NewPDU(bacgopes.NewMessageBridge(pduBytes...))), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacgopes.NewArgs((*bacgopes.OriginalUnicastNPDU)(nil)), bacgopes.NewKWArgs(bacgopes.KWPDUData, xpdu))
}

func (suite *TestAnnexJCodecSuite) TestOriginalBroadcastNPDU() {
	xpdu, err := bacgopes.Xtob(
		// "deadbeef", // forwarded PDU // TODO: this is not a ndpu so we just exploded with that. We use the iartn for that for now
		// TODO: this below is from us as upstream message is not parsable
		"01.80" + // version, network layer message
			"01 0001 0002 0003", // message type and network list
	)
	suite.Require().NoError(err)
	pduBytes, err := bacgopes.Xtob("81.0b.000d" + //   bvlci // TODO: length was 08 before
		// "deadbeef", // forwarded PDU // TODO: this is not a ndpu so we just exploded with that. We use the iartn for that for now
		// TODO: this below is from us as upstream message is not parsable
		"01.80" + // version, network layer message
		"01 0001 0002 0003", // message type and network list
	)
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(bacgopes.NewArgs(OriginalBroadcastNPDU(xpdu)), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(bacgopes.NoArgs, bacgopes.NewKWArgs(bacgopes.KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(bacgopes.NewArgs(bacgopes.NewPDU(bacgopes.NewMessageBridge(pduBytes...))), bacgopes.NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(bacgopes.NewArgs((*bacgopes.OriginalBroadcastNPDU)(nil)), bacgopes.NewKWArgs(bacgopes.KWPDUData, xpdu))
}

func TestAnnexJCodec(t *testing.T) {
	suite.Run(t, new(TestAnnexJCodecSuite))
}
