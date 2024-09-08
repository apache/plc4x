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

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/bvll"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/bvllservice"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/deleteme"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/quick"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

type TestAnnexJCodecSuite struct {
	suite.Suite

	client *TrappedClient
	codec  *AnnexJCodec
	server *TrappedServer

	log zerolog.Logger
}

func (suite *TestAnnexJCodecSuite) SetupTest() {
	suite.log = testutils.ProduceTestingLogger(suite.T())
	// minature trapped stack
	var err error
	suite.codec, err = NewAnnexJCodec(suite.log)
	suite.Require().NoError(err)
	suite.client, err = NewTrappedClient(suite.log)
	suite.Require().NoError(err)
	suite.server, err = NewTrappedServer(suite.log)
	suite.Require().NoError(err)
	err = Bind(suite.log, suite.client, suite.codec, suite.server)
	suite.Require().NoError(err)
}

// Pass the PDU to the client to send down the stack.
func (suite *TestAnnexJCodecSuite) Request(args Args, kwargs KWArgs) error {
	suite.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Request")

	return suite.client.Request(args, kwargs)
}

// Check what the server received.
func (suite *TestAnnexJCodecSuite) Indication(args Args, kwargs KWArgs) error {
	suite.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")

	var pduType any
	if len(args) > 0 {
		pduType = args[0].(any)
	}
	pduAttrs := kwargs
	suite.Assert().True(MatchPdu(suite.log, suite.server.GetIndicationReceived(), pduType, pduAttrs))
	return nil
}

// Check what the server received.
func (suite *TestAnnexJCodecSuite) Response(args Args, kwargs KWArgs) error {
	suite.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Response")

	return suite.server.Response(args, kwargs)
}

// Check what the server received.
func (suite *TestAnnexJCodecSuite) Confirmation(args Args, kwargs KWArgs) error {
	suite.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Confirmation")

	pduType := args[0].(any)
	pduAttrs := kwargs
	suite.Assert().True(MatchPdu(suite.log, suite.client.GetConfirmationReceived(), pduType, pduAttrs))
	return nil
}

func (suite *TestAnnexJCodecSuite) TestResult() {
	// Request successful
	pduBytes, err := Xtob("81.00.0006.0000")
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NewArgs(quick.Result(0)), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NewKWArgs(KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NewArgs(NewPDU(NewMessageBridge(pduBytes...))), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(NewArgs((*Result)(nil)), NewKWArgs(KWBvlciResultCode, readWriteModel.BVLCResultCode(0)))

	// Request error condition
	pduBytes, err = Xtob("81.00.0006.0010") // TODO: check if this is right or if it should be 01 as there is no code for 1
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NewArgs(quick.Result(0x0010)), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NewKWArgs(KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NewArgs(NewPDU(NewMessageBridge(pduBytes...))), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(NewArgs((*Result)(nil)), NewKWArgs(KWBvlciResultCode, readWriteModel.BVLCResultCode(0x0010)))
}

func (suite *TestAnnexJCodecSuite) TestWriteBroadcastDistributionTable() {
	// write an empty table
	pduBytes, err := Xtob("81.01.0004")
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NewArgs(quick.WriteBroadcastDistributionTable()), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NewKWArgs(KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NewArgs(NewPDU(NewMessageBridge(pduBytes...))), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(NewArgs((*WriteBroadcastDistributionTable)(nil)), NewKWArgs(KWBvlciBDT, []*Address{}))

	// write table with an element
	addr, _ := NewAddress(zerolog.Nop(), "192.168.0.254/24")
	pduBytes, err = Xtob("81.01.000e" +
		"c0.a8.00.fe.ba.c0 ff.ff.ff.00") // address and mask
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NewArgs(quick.WriteBroadcastDistributionTable(addr)), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NewKWArgs(KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NewArgs(NewPDU(NewMessageBridge(pduBytes...))), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(NewArgs((*WriteBroadcastDistributionTable)(nil)), NewKWArgs(KWBvlciBDT, []*Address{addr}))
}

func (suite *TestAnnexJCodecSuite) TestReadBroadcastDistributionTable() {
	// Read an empty table
	pduBytes, err := Xtob("81.02.0004")
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NewArgs(quick.ReadBroadcastDistributionTable()), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NewKWArgs(KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NewArgs(NewPDU(NewMessageBridge(pduBytes...))), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(NewArgs((*ReadBroadcastDistributionTable)(nil)), NoKWArgs)
}

func (suite *TestAnnexJCodecSuite) TestReadBroadcastDistributionTableAck() {
	// Read an empty TableAck
	pduBytes, err := Xtob("81.03.0004")
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NewArgs(quick.ReadBroadcastDistributionTableAck()), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NewKWArgs(KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NewArgs(NewPDU(NewMessageBridge(pduBytes...))), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(NewArgs((*ReadBroadcastDistributionTableAck)(nil)), NewKWArgs(KWBvlciBDT, []*Address{}))

	// Read TableAck with an element
	addr, _ := NewAddress(zerolog.Nop(), "192.168.0.254/24")
	pduBytes, err = Xtob("81.03.000e" + //bvlci
		"c0.a8.00.fe.ba.c0 ff.ff.ff.00") // address and mask
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NewArgs(quick.ReadBroadcastDistributionTableAck(addr)), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NewKWArgs(KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NewArgs(NewPDU(NewMessageBridge(pduBytes...))), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(NewArgs((*ReadBroadcastDistributionTableAck)(nil)), NewKWArgs(KWBvlciBDT, []*Address{addr}))
}

func (suite *TestAnnexJCodecSuite) TestForwardNPDU() {
	addr, err := NewAddress(zerolog.Nop(), "192.168.0.1")
	xpdu, err := Xtob(
		// "deadbeef", // forwarded PDU // TODO: this is not a ndpu so we just exploded with that. We use the iartn for that for now
		// TODO: this below is from us as upstream message is not parsable
		"01.80" + // version, network layer message
			"01 0001 0002 0003", // message type and network list
	)
	suite.Require().NoError(err)
	pduBytes, err := Xtob("81.04.0013" + //   bvlci // TODO: length was 0e before
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

	err = suite.Request(NewArgs(quick.ForwardedNPDU(addr, xpdu)), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NewKWArgs(KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NewArgs(NewPDU(NewMessageBridge(pduBytes...))), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(NewArgs((*ForwardedNPDU)(nil)), NewKWArgs(KWBvlciAddress, addr, KWPDUData, xpdu))
	suite.Assert().NoError(err)
}

func (suite *TestAnnexJCodecSuite) TestRegisterForeignDevice() {
	// Request successful
	pduBytes, err := Xtob(
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

	err = suite.Request(NewArgs(quick.RegisterForeignDevice(30)), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NewKWArgs(KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NewArgs(NewPDU(NewMessageBridge(pduBytes...))), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(NewArgs((*RegisterForeignDevice)(nil)), NewKWArgs(KWBvlciTimeToLive, uint16(30)))
}

func (suite *TestAnnexJCodecSuite) TestReadForeignDeviceTable() {
	// Read an empty table
	pduBytes, err := Xtob("81.06.0004")
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NewArgs(quick.ReadForeignDeviceTable()), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NewKWArgs(KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NewArgs(NewPDU(NewMessageBridge(pduBytes...))), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(NewArgs((*ReadForeignDeviceTable)(nil)), NoKWArgs)
}

func (suite *TestAnnexJCodecSuite) TestReadForeignDeviceTableAck() {
	// Read an empty TableAck
	pduBytes, err := Xtob("81.07.0004")
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NewArgs(quick.ReadForeignDeviceTableAck()), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NewKWArgs(KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NewArgs(NewPDU(NewMessageBridge(pduBytes...))), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(NewArgs((*ReadForeignDeviceTableAck)(nil)), NewKWArgs(KWBvlciFDT, []*FDTEntry{}))

	// Read TableAck with one entry
	fdte := quick.FDTEntry()
	fdte.FDAddress, err = NewAddress(suite.log, "192.168.0.10")
	suite.Require().NoError(err)
	fdte.FDTTL = 30
	fdte.FDRemain = 15
	pduBytes, err = Xtob(
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

	err = suite.Request(NewArgs(quick.ReadForeignDeviceTableAck(fdte)), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NewKWArgs(KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NewArgs(NewPDU(NewMessageBridge(pduBytes...))), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(NewArgs((*ReadForeignDeviceTableAck)(nil)), NewKWArgs(KWBvlciFDT, []*FDTEntry{fdte}))
}

func (suite *TestAnnexJCodecSuite) TestDeleteForeignDeviceTableEntry() {
	addr, _ := NewAddress(zerolog.Nop(), "192.168.0.11/24")
	pduBytes, err := Xtob("81.08.000a" + // bvlci
		"c0.a8.00.0b.ba.c0") // address of entry to be deleted
	suite.Require().NoError(err)
	{ // Parse with plc4x parser to validate
		parse, err := readWriteModel.BVLCParse[readWriteModel.BVLC](testutils.TestContext(suite.T()), pduBytes)
		suite.Assert().NoError(err)
		if parse != nil {
			suite.T().Log("\n" + parse.String())
		}
	}

	err = suite.Request(NewArgs(quick.DeleteForeignDeviceTableEntry(addr)), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NewKWArgs(KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NewArgs(NewPDU(NewMessageBridge(pduBytes...))), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(NewArgs((*DeleteForeignDeviceTableEntry)(nil)), NewKWArgs(KWBvlciAddress, addr))
}

func (suite *TestAnnexJCodecSuite) TestDeleteForeignDeviceTableAck() {
	// TODO: implement me
	suite.T().Skip()
}

func (suite *TestAnnexJCodecSuite) TestDistributeBroadcastToNetwork() {
	xpdu, err := Xtob(
		// "deadbeef", // forwarded PDU // TODO: this is not a ndpu so we just exploded with that. We use the iartn for that for now
		// TODO: this below is from us as upstream message is not parsable
		"01.80" + // version, network layer message
			"01 0001 0002 0003", // message type and network list
	)
	suite.Require().NoError(err)
	pduBytes, err := Xtob("81.09.000d" + //   bvlci // TODO: length was 08 before
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

	err = suite.Request(NewArgs(quick.DistributeBroadcastToNetwork(xpdu)), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NewKWArgs(KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NewArgs(NewPDU(NewMessageBridge(pduBytes...))), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(NewArgs((*DistributeBroadcastToNetwork)(nil)), NewKWArgs(KWPDUData, xpdu))
}

func (suite *TestAnnexJCodecSuite) TestOriginalUnicastNPDU() {
	xpdu, err := Xtob(
		// "deadbeef", // forwarded PDU // TODO: this is not a ndpu so we just exploded with that. We use the iartn for that for now
		// TODO: this below is from us as upstream message is not parsable
		"01.80" + // version, network layer message
			"01 0001 0002 0003", // message type and network list
	)
	suite.Require().NoError(err)
	pduBytes, err := Xtob("81.0a.000d" + //   bvlci // TODO: length was 08 before
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

	err = suite.Request(NewArgs(quick.OriginalUnicastNPDU(xpdu)), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NewKWArgs(KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NewArgs(NewPDU(NewMessageBridge(pduBytes...))), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(NewArgs((*OriginalUnicastNPDU)(nil)), NewKWArgs(KWPDUData, xpdu))
}

func (suite *TestAnnexJCodecSuite) TestOriginalBroadcastNPDU() {
	xpdu, err := Xtob(
		// "deadbeef", // forwarded PDU // TODO: this is not a ndpu so we just exploded with that. We use the iartn for that for now
		// TODO: this below is from us as upstream message is not parsable
		"01.80" + // version, network layer message
			"01 0001 0002 0003", // message type and network list
	)
	suite.Require().NoError(err)
	pduBytes, err := Xtob("81.0b.000d" + //   bvlci // TODO: length was 08 before
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

	err = suite.Request(NewArgs(quick.OriginalBroadcastNPDU(xpdu)), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Indication(NoArgs, NewKWArgs(KWPDUData, pduBytes))
	suite.Assert().NoError(err)

	err = suite.Response(NewArgs(NewPDU(NewMessageBridge(pduBytes...))), NoKWArgs)
	suite.Assert().NoError(err)
	err = suite.Confirmation(NewArgs((*OriginalBroadcastNPDU)(nil)), NewKWArgs(KWPDUData, xpdu))
}

func TestAnnexJCodec(t *testing.T) {
	suite.Run(t, new(TestAnnexJCodecSuite))
}
