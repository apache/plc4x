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

package umas

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/umas/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transactions"
)

const (
	testHostname        = "M340"
	testPlcModel        = uint16(0x0141)
	testComVersion      = uint16(0x0207)
	testMaxFrameSize    = uint16(260)
	testFirmwareVersion = uint16(0x0304)
	testHardwareId      = uint32(0xDEADBEEF)
	testFirstHash       = uint32(0x11111111)
	testSecondHash      = uint32(0x22222222)
	// testProjectCrc is the sum of the two hashes, which is how plc4j derives the CRC that every
	// variable read and write has to carry.
	testProjectCrc = testFirstHash + testSecondHash
)

func newTestConnection(t *testing.T, configuration Configuration) (*Connection, *stubCodec) {
	t.Helper()
	_options := testutils.EnrichOptionsWithOptionsForTesting(t)
	codec := newStubCodec()
	tm := transactions.NewRequestTransactionManager(maxConcurrentRequests, _options...)
	t.Cleanup(func() {
		assert.NoError(t, tm.Close())
	})
	driverContext, err := NewDriverContext(configuration)
	require.NoError(t, err)
	connection := NewConnection(codec, configuration, driverContext, NewTagHandler(), tm, map[string][]string{}, _options...)
	return connection, codec
}

// plcFixture is the PLC the handshake talks to: what it says about itself and what its data
// dictionary holds.
type plcFixture struct {
	datatypes []datatypeRecord
	// typeDefinitions is the payload answered for the DD02 request of one custom type id.
	typeDefinitions map[uint16][]byte
	symbols         []symbolRecord
}

// defaultFixture is a PLC with one primitive symbol of every interesting kind, one struct type and
// one array type.
func defaultFixture() plcFixture {
	return plcFixture{
		datatypes: []datatypeRecord{
			// customTypeIdBase: an alias of a primitive, so no definition is fetched for it.
			{name: "MY_STRING", dataSize: 20, classIdentifier: 0, dataType: uint8(typeIdString)},
			// customTypeIdBase+1: a struct, whose definition is fetched.
			{name: "MY_STRUCT", dataSize: 8, classIdentifier: 2},
			// customTypeIdBase+2: an array, whose definition is fetched.
			{name: "MY_ARRAY", dataSize: 40, classIdentifier: 4},
		},
		typeDefinitions: map[uint16][]byte{
			customTypeIdBase + 1: udtDefinitionPayload(
				udtMember{name: "meta", dataType: typeIdDint, offset: 0},
				udtMember{name: "r32", dataType: typeIdReal, offset: 4}),
			customTypeIdBase + 2: arrayTypeDefinitionPayload(typeIdDint,
				arrayDimension{startIndex: 0, upperBound: 9}),
		},
		symbols: []symbolRecord{
			{name: "g_r32", dataType: typeIdReal, block: 2, offset: 0x00001234},
			{name: "g_b16", dataType: typeIdInt, block: 2, offset: 0x00001240},
			{name: "g_string", dataType: typeIdString, block: 2, offset: 0x00001250},
			{name: "g_plant", dataType: customTypeIdBase + 1, block: 3, offset: 0x00000000},
			{name: "g_arrInt", dataType: customTypeIdBase + 2, block: 3, offset: 0x00000100},
			// An offset big enough to show that a read and a write split it differently.
			{name: "g_far", dataType: typeIdDint, block: 4, offset: 0x00012345},
			// And one no read reference can address: its high bits don't fit the 16 bit baseOffset.
			{name: "g_toofar", dataType: typeIdDint, block: 5, offset: 0x01000000},
		},
	}
}

// runHandshake plays the PLC side of the connect handshake, asserting that every request is the one
// plc4j's doHandshake sends at that point and in that order. The sequence is capture-derived, so this
// is the only place it is written down and checked.
func runHandshake(t *testing.T, codec *stubCodec, fixture plcFixture) {
	t.Helper()

	// 1. PlcIdent: who are you?
	request := codec.nextRequest(t)
	identRequest, ok := request.item(t).(readWriteModel.UmasPDUPlcIdentRequest)
	require.True(t, ok, "the handshake has to start with a PlcIdent, got %T", request.item(t))
	assert.Equal(t, uint8(0x02), identRequest.GetUmasFunctionKey())
	assert.Equal(t, uint8(0), identRequest.GetPairingKey(), "plc4j never assigns a pairing key")
	codec.answerWith(t, request, readWriteModel.NewUmasPDUPlcIdentResponse(
		0, 0, 0, testPlcModel, testComVersion, 0, 0, 0, 0, 0,
		uint8(len(testHostname)), testHostname, 0, nil))

	// 2. InitComms: negotiate the frame size. plc4j sends subCode 0x00.
	request = codec.nextRequest(t)
	initRequest, ok := request.item(t).(readWriteModel.UmasInitCommsRequest)
	require.True(t, ok, "InitComms has to follow PlcIdent, got %T", request.item(t))
	assert.Equal(t, uint8(0x01), initRequest.GetUmasFunctionKey())
	assert.Equal(t, uint8(0x00), initRequest.GetSubCode())
	codec.answerWith(t, request, readWriteModel.NewUmasInitCommsResponse(
		0, testMaxFrameSize, testFirmwareVersion, 0, 0, uint8(len(testHostname)), testHostname))

	// 3. Repeat: echo a payload of the negotiated frame size minus three, first byte zero and the
	// rest 0x54.
	request = codec.nextRequest(t)
	repeatRequest, ok := request.item(t).(readWriteModel.UmasPDURepeatRequest)
	require.True(t, ok, "Repeat has to follow InitComms, got %T", request.item(t))
	assert.Equal(t, uint8(0x0A), repeatRequest.GetUmasFunctionKey())
	payload := repeatRequest.GetData()
	require.Len(t, payload, int(testMaxFrameSize-repeatPayloadOverhead),
		"the echo payload is sized from the frame size the PLC reported, not the configured one")
	assert.Equal(t, byte(0x00), payload[0], "the first byte of the echo payload stays zero")
	for i := 1; i < len(payload); i++ {
		require.Equal(t, repeatFillerByte, payload[i], "byte %d of the echo payload", i)
	}
	codec.answerWith(t, request, readWriteModel.NewUmasPDURepeatResponse(0, payload))

	// 4. Memory block 0x30: the hardware id and the two hashes whose sum is the project CRC.
	request = codec.nextRequest(t)
	blockRequest, ok := request.item(t).(readWriteModel.UmasPDUReadMemoryBlockRequest)
	require.True(t, ok, "a memory block read has to follow Repeat, got %T", request.item(t))
	assert.Equal(t, uint8(0x20), blockRequest.GetUmasFunctionKey())
	assert.Equal(t, projectMemoryBlockNumber, blockRequest.GetBlockNumber())
	assert.Equal(t, handshakeMemoryBlockRange, blockRequest.GetRange())
	assert.Equal(t, uint16(0), blockRequest.GetOffset())
	assert.Equal(t, handshakeMemoryBlockLength, blockRequest.GetNumberOfBytes())
	projectBlock := projectMemoryBlockPayload(testHardwareId, testFirstHash, testSecondHash)
	codec.answerWith(t, request, readWriteModel.NewUmasPDUReadMemoryBlockResponse(
		0, 0, uint16(len(projectBlock)), projectBlock))

	// 5. ProjectInfo(1).
	answerProjectInfo(t, codec, 1)

	// 6. Memory block 0x13. plc4j reads it, logs its length and looks at nothing in it.
	request = codec.nextRequest(t)
	blockRequest, ok = request.item(t).(readWriteModel.UmasPDUReadMemoryBlockRequest)
	require.True(t, ok, "the second memory block read has to follow ProjectInfo(1), got %T", request.item(t))
	assert.Equal(t, secondaryMemoryBlockNumber, blockRequest.GetBlockNumber())
	codec.answerWith(t, request, readWriteModel.NewUmasPDUReadMemoryBlockResponse(0, 0, 1, []byte{0x00}))

	// 7. ProjectInfo(0), (4), (1) again and (3). Subcode 1 really is sent twice.
	for _, subcode := range []uint8{0, 4, 1, 3} {
		answerProjectInfo(t, codec, subcode)
	}

	runDictionaryDownload(t, codec, fixture)
}

func answerProjectInfo(t *testing.T, codec *stubCodec, wantSubcode uint8) {
	t.Helper()
	request := codec.nextRequest(t)
	projectInfoRequest, ok := request.item(t).(readWriteModel.UmasPDUProjectInfoRequest)
	require.True(t, ok, "expected a ProjectInfo request, got %T", request.item(t))
	assert.Equal(t, uint8(0x03), projectInfoRequest.GetUmasFunctionKey())
	assert.Equal(t, wantSubcode, projectInfoRequest.GetSubcode())
	codec.answerWith(t, request, readWriteModel.NewUmasPDUProjectInfoResponse(0, []byte{0x00}))
}

// runDictionaryDownload plays the PLC side of the data dictionary download: the datatype dictionary
// first, then one request per custom type which has a definition, then the symbol table.
func runDictionaryDownload(t *testing.T, codec *stubCodec, fixture plcFixture) {
	t.Helper()

	// The datatype dictionary, record type DD03, which carries no trailing padding.
	request := codec.nextRequest(t)
	dictionaryRequest, ok := request.item(t).(readWriteModel.UmasPDUReadUnlocatedVariableNamesRequest)
	require.True(t, ok, "expected a data dictionary request, got %T", request.item(t))
	assert.Equal(t, uint8(0x26), dictionaryRequest.GetUmasFunctionKey())
	assert.Equal(t, recordTypeDatatypes, dictionaryRequest.GetRecordType())
	assert.Equal(t, dictionaryRequestIndex, dictionaryRequest.GetIndex())
	assert.Equal(t, testHardwareId, dictionaryRequest.GetHardwareId(),
		"the dictionary request needs the hardware id from memory block 0x30")
	assert.Equal(t, uint16(0), dictionaryRequest.GetBlockNo())
	assert.Nil(t, dictionaryRequest.GetBlank(), "a DD03 request has no trailing padding")
	codec.answerWith(t, request, readWriteModel.NewUmasPDUReadUnlocatedVariableResponse(
		0, datatypeDictionaryPayload(fixture.datatypes...)))

	// One request per custom type whose class identifier is non-zero, in dictionary order.
	for i, datatype := range fixture.datatypes {
		if datatype.classIdentifier == 0 {
			continue
		}
		typeId := customTypeIdBase + uint16(i)
		request = codec.nextRequest(t)
		typeRequest, ok := request.item(t).(readWriteModel.UmasPDUReadUnlocatedVariableNamesRequest)
		require.True(t, ok, "expected a type definition request, got %T", request.item(t))
		assert.Equal(t, recordTypeSymbols, typeRequest.GetRecordType())
		assert.Equal(t, typeId, typeRequest.GetBlockNo(),
			"the type is addressed by its dictionary index offset by the custom type threshold")
		require.NotNil(t, typeRequest.GetBlank(), "a DD02 request carries the trailing padding")
		assert.Equal(t, uint16(0), *typeRequest.GetBlank())
		codec.answerWith(t, request, readWriteModel.NewUmasPDUReadUnlocatedVariableResponse(
			0, fixture.typeDefinitions[typeId]))
	}

	// And the symbol table, addressed by the pseudo block number 0xFFFF.
	request = codec.nextRequest(t)
	symbolRequest, ok := request.item(t).(readWriteModel.UmasPDUReadUnlocatedVariableNamesRequest)
	require.True(t, ok, "expected the symbol table request, got %T", request.item(t))
	assert.Equal(t, recordTypeSymbols, symbolRequest.GetRecordType())
	assert.Equal(t, symbolTableBlockNumber, symbolRequest.GetBlockNo())
	codec.answerWith(t, request, readWriteModel.NewUmasPDUReadUnlocatedVariableResponse(
		0, symbolTablePayload(fixture.symbols...)))
}

// connect runs Connect on its own go routine while the test plays the PLC.
func connect(t *testing.T, connection *Connection) <-chan error {
	t.Helper()
	connectResult := make(chan error, 1)
	go func() {
		connectResult <- connection.Connect(testutils.TestContext(t))
	}()
	t.Cleanup(func() {
		assert.NoError(t, connection.Close())
	})
	return connectResult
}

func requireConnected(t *testing.T, connectResult <-chan error) {
	t.Helper()
	select {
	case err := <-connectResult:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Connect never returned")
	}
}

// newConnectedConnection is the fixture the read, write and browse tests start from: a connection
// whose handshake and dictionary download are done.
func newConnectedConnection(t *testing.T) (*Connection, *stubCodec) {
	t.Helper()
	connection, codec := newTestConnection(t, DefaultConfiguration())
	connectResult := connect(t, connection)
	runHandshake(t, codec, defaultFixture())
	requireConnected(t, connectResult)
	return connection, codec
}

func TestConnection_ConnectRunsTheHandshakeAndDownloadsTheDictionary(t *testing.T) {
	connection, _ := newConnectedConnection(t)

	assert.True(t, connection.IsConnected())
	// What the handshake learned, which every later request depends on.
	assert.Equal(t, testHardwareId, connection.session.getHardwareId())
	assert.Equal(t, uint32(testProjectCrc), connection.session.getProjectCrc())
	assert.Equal(t, testMaxFrameSize, connection.session.getMaxFrameSize())

	// And the dictionary.
	assert.True(t, connection.session.hasSymbols())
	symbol, ok := connection.session.lookupSymbol("g_r32")
	require.True(t, ok)
	assert.Equal(t, typeIdReal, symbol.GetDataType())
	assert.Equal(t, uint32(0x00001234), symbol.GetOffset())

	structType, known := connection.session.customType(customTypeIdBase + 1)
	require.True(t, known, "the struct type should have been resolved")
	assert.Equal(t, "MY_STRUCT", structType.name)
	require.Len(t, structType.fields, 2)

	arrayType, known := connection.session.customType(customTypeIdBase + 2)
	require.True(t, known, "the array type should have been resolved")
	assert.True(t, arrayType.isArray)
	assert.Equal(t, typeIdDint, arrayType.elementTypeId)

	// The alias of a primitive has no definition of its own, so no request was sent for it - but its
	// declared size is what a STRING read of a symbol of that type asks for.
	_, known = connection.session.customType(customTypeIdBase)
	assert.False(t, known, "an alias of a primitive needs no definition")
	assert.Equal(t, uint16(20), connection.session.dataTypeSizes[customTypeIdBase])
}

// Every request of the handshake carries a transaction identifier of its own, because that identifier
// is the only thing tying a response to its request.
func TestConnection_EveryRequestGetsItsOwnTransactionIdentifier(t *testing.T) {
	_, codec := newConnectedConnection(t)

	codec.mutex.Lock()
	defer codec.mutex.Unlock()
	seen := map[uint16]bool{}
	for _, message := range codec.sent {
		adu, ok := message.(readWriteModel.ModbusTcpADU)
		require.True(t, ok)
		identifier := adu.GetTransactionIdentifier()
		assert.False(t, seen[identifier], "the transaction identifier %d was used twice", identifier)
		assert.NotEqual(t, uint16(0), identifier, "zero is left out so it never looks uninitialized")
		seen[identifier] = true
	}
	assert.NotEmpty(t, seen)
}

func TestConnection_ConnectFailsWhenThePlcDoesNotIdentifyItself(t *testing.T) {
	connection, codec := newTestConnection(t, DefaultConfiguration())
	connectResult := connect(t, connection)

	request := codec.nextRequest(t)
	codec.failRequest(t, request, errors.New("the PLC hung up"))

	select {
	case err := <-connectResult:
		assert.Error(t, err, "a failed handshake has to fail the connect")
	case <-time.After(5 * time.Second):
		t.Fatal("Connect never returned")
	}
	assert.False(t, connection.IsConnected())
}

func TestConnection_ConnectFailsWhenTheHandshakeIsAnsweredWithTheWrongThing(t *testing.T) {
	connection, codec := newTestConnection(t, DefaultConfiguration())
	connectResult := connect(t, connection)

	// A success response instead of the PlcIdent answer: the PLC said something, but not what the
	// handshake needs.
	request := codec.nextRequest(t)
	codec.answerWith(t, request, readWriteModel.NewUmasPDUSuccessResponse(0, nil))

	select {
	case err := <-connectResult:
		assert.Error(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Connect never returned")
	}
}

func TestConnection_ConnectFailsOnAModbusException(t *testing.T) {
	connection, codec := newTestConnection(t, DefaultConfiguration())
	connectResult := connect(t, connection)

	request := codec.nextRequest(t)
	codec.answerWithPdu(t, request,
		readWriteModel.NewModbusPDUError(readWriteModel.ModbusErrorCode_ILLEGAL_FUNCTION))

	select {
	case err := <-connectResult:
		assert.Error(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Connect never returned")
	}
}

func TestConnection_ConnectFailsWhenNothingCanBeSent(t *testing.T) {
	connection, codec := newTestConnection(t, DefaultConfiguration())
	codec.failSends()
	assert.Error(t, connection.Connect(testutils.TestContext(t)))
	assert.NoError(t, connection.Close())
}

// A dictionary download which fails leaves a usable connection: it can still ping, and a read of an
// unresolvable symbol answers NOT_FOUND rather than hanging. plc4j warns and carries on in the same
// situation.
func TestConnection_ConnectSurvivesAFailedDictionaryDownload(t *testing.T) {
	connection, codec := newTestConnection(t, DefaultConfiguration())
	connectResult := connect(t, connection)

	// Everything up to the dictionary, then a refusal.
	runHandshakeWithoutDictionary(t, codec)
	request := codec.nextRequest(t)
	codec.answerWith(t, request, readWriteModel.NewUmasPDUErrorResponse(0, []byte{0x01}))

	requireConnected(t, connectResult)
	assert.True(t, connection.IsConnected())
	assert.False(t, connection.session.hasSymbols())
}

// runHandshakeWithoutDictionary is runHandshake up to, but not including, the dictionary download.
func runHandshakeWithoutDictionary(t *testing.T, codec *stubCodec) {
	t.Helper()
	request := codec.nextRequest(t)
	require.IsType(t, readWriteModel.NewUmasPDUPlcIdentRequest(0), request.item(t))
	codec.answerWith(t, request, readWriteModel.NewUmasPDUPlcIdentResponse(
		0, 0, 0, testPlcModel, testComVersion, 0, 0, 0, 0, 0,
		uint8(len(testHostname)), testHostname, 0, nil))

	request = codec.nextRequest(t)
	codec.answerWith(t, request, readWriteModel.NewUmasInitCommsResponse(
		0, testMaxFrameSize, testFirmwareVersion, 0, 0, uint8(len(testHostname)), testHostname))

	request = codec.nextRequest(t)
	repeatRequest := request.item(t).(readWriteModel.UmasPDURepeatRequest)
	codec.answerWith(t, request, readWriteModel.NewUmasPDURepeatResponse(0, repeatRequest.GetData()))

	request = codec.nextRequest(t)
	projectBlock := projectMemoryBlockPayload(testHardwareId, testFirstHash, testSecondHash)
	codec.answerWith(t, request, readWriteModel.NewUmasPDUReadMemoryBlockResponse(
		0, 0, uint16(len(projectBlock)), projectBlock))

	answerProjectInfo(t, codec, 1)

	request = codec.nextRequest(t)
	codec.answerWith(t, request, readWriteModel.NewUmasPDUReadMemoryBlockResponse(0, 0, 1, []byte{0x00}))

	for _, subcode := range []uint8{0, 4, 1, 3} {
		answerProjectInfo(t, codec, subcode)
	}
}

// A block 0x30 too short to hold the project identity is logged and tolerated - the reads that need
// the CRC are what would fail, and they say so themselves.
func TestConnection_ConnectSurvivesAShortProjectBlock(t *testing.T) {
	connection, codec := newTestConnection(t, DefaultConfiguration())
	connectResult := connect(t, connection)

	request := codec.nextRequest(t)
	codec.answerWith(t, request, readWriteModel.NewUmasPDUPlcIdentResponse(
		0, 0, 0, testPlcModel, testComVersion, 0, 0, 0, 0, 0,
		uint8(len(testHostname)), testHostname, 0, nil))
	request = codec.nextRequest(t)
	codec.answerWith(t, request, readWriteModel.NewUmasInitCommsResponse(
		0, testMaxFrameSize, testFirmwareVersion, 0, 0, uint8(len(testHostname)), testHostname))
	request = codec.nextRequest(t)
	codec.answerWith(t, request, readWriteModel.NewUmasPDURepeatResponse(0,
		request.item(t).(readWriteModel.UmasPDURepeatRequest).GetData()))

	// Ten bytes where seventeen are needed.
	request = codec.nextRequest(t)
	codec.answerWith(t, request, readWriteModel.NewUmasPDUReadMemoryBlockResponse(0, 0, 10, make([]byte, 10)))

	answerProjectInfo(t, codec, 1)
	request = codec.nextRequest(t)
	codec.answerWith(t, request, readWriteModel.NewUmasPDUReadMemoryBlockResponse(0, 0, 1, []byte{0x00}))
	for _, subcode := range []uint8{0, 4, 1, 3} {
		answerProjectInfo(t, codec, subcode)
	}
	// The dictionary download goes ahead with a hardware id of zero.
	request = codec.nextRequest(t)
	codec.answerWith(t, request, readWriteModel.NewUmasPDUReadUnlocatedVariableResponse(0, nil))
	request = codec.nextRequest(t)
	codec.answerWith(t, request, readWriteModel.NewUmasPDUReadUnlocatedVariableResponse(0, nil))

	requireConnected(t, connectResult)
	assert.Equal(t, uint32(0), connection.session.getProjectCrc())
}

func TestConnection_Ping(t *testing.T) {
	connection, codec := newConnectedConnection(t)

	pingResult := make(chan error, 1)
	go func() { pingResult <- connection.Ping(testutils.TestContext(t)) }()

	request := codec.nextRequest(t)
	statusRequest, ok := request.item(t).(readWriteModel.UmasPDUPlcStatusRequest)
	require.True(t, ok, "a ping is a status request, got %T", request.item(t))
	assert.Equal(t, uint8(0x04), statusRequest.GetUmasFunctionKey())
	codec.answerWith(t, request, readWriteModel.NewUmasPDUSuccessResponse(0, nil))

	select {
	case err := <-pingResult:
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Ping never returned")
	}
}

// The PLC refusing the status request still proves it is there, which is all a ping asks.
func TestConnection_PingAcceptsARefusal(t *testing.T) {
	connection, codec := newConnectedConnection(t)

	pingResult := make(chan error, 1)
	go func() { pingResult <- connection.Ping(testutils.TestContext(t)) }()

	request := codec.nextRequest(t)
	codec.answerWith(t, request, readWriteModel.NewUmasPDUErrorResponse(0, []byte{0x01}))

	select {
	case err := <-pingResult:
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Ping never returned")
	}
}

func TestConnection_PingFailsWhenThePlcIsGone(t *testing.T) {
	connection, codec := newConnectedConnection(t)

	pingResult := make(chan error, 1)
	go func() { pingResult <- connection.Ping(testutils.TestContext(t)) }()

	codec.failRequest(t, codec.nextRequest(t), errors.New("the PLC hung up"))

	select {
	case err := <-pingResult:
		assert.Error(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Ping never returned")
	}
}

func TestConnection_String(t *testing.T) {
	connection, _ := newTestConnection(t, DefaultConfiguration())
	assert.Contains(t, connection.String(), "umas.Connection")
	assert.NoError(t, connection.Close())
}

// Close without a preceding Connect, and twice, has to be harmless.
func TestConnection_CloseIsIdempotent(t *testing.T) {
	connection, _ := newTestConnection(t, DefaultConfiguration())
	assert.NoError(t, connection.Close())
	assert.NoError(t, connection.Close())
}

// The echo payload is sized from the frame size the PLC reports, but a PLC reporting the biggest
// frame size there is would push the Modbus/TCP length field one byte past what a uint16 can hold -
// and the frame would go out claiming a length of zero.
func TestConnection_TheEchoPayloadStaysInsideTheLengthField(t *testing.T) {
	connection, _ := newTestConnection(t, DefaultConfiguration())
	t.Cleanup(func() { assert.NoError(t, connection.Close()) })

	// The largest frame size a PLC can report.
	connection.session.setCommsParameters(0xFFFF, 0)
	payload := connection.buildRepeatPayload()
	assert.Len(t, payload, maxRepeatPayloadSize)

	// The whole frame still serializes to a length field which counts every byte after it.
	request := readWriteModel.NewModbusTcpADU(1, 0,
		readWriteModel.NewUmasPDU(readWriteModel.NewUmasPDURepeatRequest(0, payload)))
	theBytes, err := request.Serialize()
	require.NoError(t, err)
	length := uint16(theBytes[4])<<8 | uint16(theBytes[5])
	assert.Equal(t, len(theBytes), int(length)+6, "the length field has to count the rest of the frame")

	// A frame size that leaves room is used as it is.
	connection.session.setCommsParameters(260, 0)
	assert.Len(t, connection.buildRepeatPayload(), int(260-repeatPayloadOverhead))
}
