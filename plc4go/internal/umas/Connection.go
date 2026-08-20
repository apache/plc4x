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
	"context"
	"encoding/binary"
	"fmt"
	"runtime/debug"
	"slices"
	"sync"

	"github.com/rs/zerolog"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/umas/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/tracer"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

const (
	// handshakeTimeoutFactor is how many request timeouts the whole connect handshake may take.
	// plc4j waits requestTimeout * 10 for the same chain of exchanges.
	handshakeTimeoutFactor = 10

	// repeatFillerByte is what the echo request of the handshake fills its payload with, from the
	// second byte on - the first stays zero. Both come straight from plc4j's performRepeat, which
	// took them from a capture; what makes the PLC care about the content is not documented
	// anywhere, so the value is named after where it appears rather than after a meaning.
	repeatFillerByte = byte(0x54)
	// repeatPayloadOverhead is how much smaller than the negotiated frame size the echo payload is.
	// plc4j subtracts exactly 3 with no explanation.
	repeatPayloadOverhead = uint16(3)
	// maxRepeatPayloadSize is the largest echo payload whose frame still fits the Modbus/TCP length
	// field. That field is 16 bits and counts the unit identifier plus the whole PDU, and the echo
	// PDU is the function code, the pairing key, the UMAS function key and then the payload. A PLC
	// which reports the biggest frame size there is would otherwise push the length field one byte
	// past what it can hold, and the frame would go out with a length of zero - plc4j builds the
	// payload from maxFrameSize - 3 unconditionally and has the same hole.
	maxRepeatPayloadSize = 0xFFFF - 1 - 3

	// projectMemoryBlockNumber is the memory block the hardware id and the project hashes are read
	// out of.
	projectMemoryBlockNumber = uint16(0x30)
	// secondaryMemoryBlockNumber is a second block the handshake reads. plc4j reads it, logs its
	// length and looks at nothing in it; leaving it out has not been tried against a device, so it
	// is sent to keep the handshake byte-identical.
	secondaryMemoryBlockNumber = uint16(0x13)
	// handshakeMemoryBlockLength is how many bytes each of the two handshake block reads asks for.
	handshakeMemoryBlockLength = uint16(33)
	// handshakeMemoryBlockRange is the range byte of a memory block read. plc4j sends 0x01.
	handshakeMemoryBlockRange = uint8(0x01)

	// projectMemoryBlockMinimumLength is how much of block 0x30 has to be there for the hardware id
	// and both project hashes to be readable: a 9 byte header (range, notSure, index, hardwareId)
	// followed by two 32 bit hashes.
	projectMemoryBlockMinimumLength = 17

	// recordTypeSymbols selects the symbol and custom-type records of the data dictionary, and
	// recordTypeDatatypes the datatype dictionary. plc4j calls them RECORD_TYPE_DD02 and
	// RECORD_TYPE_DD03.
	recordTypeSymbols   = uint16(0xDD02)
	recordTypeDatatypes = uint16(0xDD03)
	// dictionaryRequestIndex is the index byte of every data-dictionary request. plc4j sends 0x03.
	dictionaryRequestIndex = uint8(0x03)
	// symbolTableBlockNumber is the pseudo block number which asks for the whole symbol table
	// rather than for one type. plc4j calls it SYMBOL_TABLE_BLOCK.
	symbolTableBlockNumber = uint16(0xFFFF)
	// dictionaryRequestOffset is the offset every data-dictionary request this driver sends starts
	// at. Continuing a dictionary at the nextAddress the response reports is not implemented, here
	// or in plc4j: both read the dictionary in a single request.
	dictionaryRequestOffset = uint16(0)
)

// projectInfoSubcodes is the sequence of project-info subcodes the handshake sends, in order, with a
// memory block read in the middle - see handshake. Subcode 1 really is sent twice; plc4j does the
// same and the sequence is capture-derived, so it is reproduced rather than deduplicated.
var projectInfoSubcodes = struct {
	beforeSecondaryBlock []uint8
	afterSecondaryBlock  []uint8
}{
	beforeSecondaryBlock: []uint8{1},
	afterSecondaryBlock:  []uint8{0, 4, 1, 3},
}

// Connection is a connection to a Schneider Modicon PLC speaking UMAS, the port of plc4j's
// UmasConnection.
//
// UMAS is tunneled inside Modbus/TCP as function code 0x5A. On connect the driver runs a multi-step
// handshake (PlcIdent, InitComms, an echo, two memory block reads and a series of project-info
// queries) and then downloads the project's data dictionary, because UMAS addresses are symbol names
// and only the dictionary says which memory block and offset a symbol lives at.
type Connection struct {
	_default.DefaultConnection

	configuration Configuration
	driverContext DriverContext
	messageCodec  spi.MessageCodec
	tm            transactions.RequestTransactionManager
	session       *session
	requester     *requester
	subscriber    _default.DefaultPollingSubscriber
	options       map[string][]string

	connectionId string
	tracer       tracer.Tracer

	// lifecycleMutex guards the stray-message drain's handles so Close can be called concurrently
	// with, and without, a preceding Connect.
	lifecycleMutex sync.Mutex
	drainCancel    context.CancelFunc
	drainDone      chan struct{}

	// dictionaryMutex serializes the data dictionary download, so a burst of browse requests on a
	// connection whose handshake couldn't download it doesn't download it several times over.
	dictionaryMutex sync.Mutex

	wg sync.WaitGroup // use to track spawned go routines

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

var (
	_ spi.TransportInstanceExposer                  = (*Connection)(nil)
	_ _default.DefaultPollingSubscriberRequirements = (*Connection)(nil)
	_ spi.PlcBrowser                                = (*Connection)(nil)
)

func NewConnection(
	messageCodec spi.MessageCodec,
	configuration Configuration,
	driverContext DriverContext,
	tagHandler spi.PlcTagHandler,
	tm transactions.RequestTransactionManager,
	connectionOptions map[string][]string,
	_options ...options.WithOption,
) *Connection {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	// Clipping the option list is what keeps the reader, the writer and the polling subscriber from
	// racing: each of them hands its own logger option downstream with append, and on a slice that
	// still has spare capacity those appends write into the very same backing array.
	_options = slices.Clip(_options)
	connection := &Connection{
		configuration: configuration,
		driverContext: driverContext,
		messageCodec:  messageCodec,
		tm:            tm,
		session:       newSession(configuration.maxFrameSize),
		options:       connectionOptions,
		log:           customLogger,
		_options:      _options,
	}
	connection.requester = newRequester(messageCodec, tm, connection.session, configuration,
		append(_options, options.WithCustomLogger(customLogger))...)
	if traceEnabledOption, ok := connectionOptions["traceEnabled"]; ok {
		if len(traceEnabledOption) == 1 {
			connection.tracer = tracer.NewTracer(connection.connectionId, _options...)
		}
	}
	connection.DefaultConnection = _default.NewDefaultConnection(connection,
		append(_options,
			_default.WithPlcTagHandler(tagHandler),
			_default.WithPlcValueHandler(NewValueHandler(_options...)),
		)...,
	)
	// UMAS has no subscription mechanism of its own; the subscriber emulates one by polling the read
	// path, which is exactly what plc4j's PollingSubscriptionConnectionBase does for this driver.
	connection.subscriber = _default.NewDefaultPollingSubscriber(connection, _options...)
	return connection
}

func (c *Connection) GetConnectionId() string {
	return c.connectionId
}

func (c *Connection) IsTraceEnabled() bool {
	return c.tracer != nil
}

func (c *Connection) GetTracer() tracer.Tracer {
	return c.tracer
}

func (c *Connection) GetConnection() plc4go.PlcConnection {
	return c
}

func (c *Connection) GetMessageCodec() spi.MessageCodec {
	return c.messageCodec
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Connect / close
////////////////////////////////////////////////////////////////////////////////////////////////////

// Connect opens the transport, runs the connect handshake and downloads the data dictionary. Without
// the handshake nothing else works: the project CRC and the hardware id it yields are required by
// every variable read, variable write and dictionary request.
func (c *Connection) Connect(ctx context.Context) error {
	c.log.Trace().Msg("Connecting")
	if err := c.messageCodec.Connect(ctx); err != nil {
		return errors.Wrap(err, "error connecting message codec")
	}
	c.startStrayMessageDrain()

	// For testing purposes we can skip waiting for a complete connection.
	if !c.driverContext.awaitSetupComplete {
		// The caller's ctx is canceled as soon as GetConnection returns, so the background
		// handshake must not inherit that cancellation.
		setupCtx := context.WithoutCancel(ctx)
		c.wg.Go(func() {
			setupCtx, cancel := utils.WithNamedTimeout(setupCtx, "umas setup timeout",
				c.configuration.requestTimeout*handshakeTimeoutFactor)
			defer cancel()
			if err := c.handshake(setupCtx); err != nil {
				c.log.Error().Err(err).Msg("error during the umas connect handshake")
			}
		})
		c.log.Warn().Msg("Connection used in an unsafe way. !!!DON'T USE IN PRODUCTION!!!")
		// Here we mark connected without waiting for the handshake to have happened.
		c.SetConnected(true)
		return nil
	}

	handshakeCtx, cancelHandshake := utils.WithNamedTimeout(ctx, "umas setup timeout",
		c.configuration.requestTimeout*handshakeTimeoutFactor)
	defer cancelHandshake()
	if err := c.handshake(handshakeCtx); err != nil {
		// The caller gets an error instead of a connection and will never call Close on it, so the
		// drain has to be reeled back in here or its goroutine outlives the failed attempt.
		c.stopStrayMessageDrain()
		if disconnectErr := c.messageCodec.Disconnect(); disconnectErr != nil {
			c.log.Debug().Err(disconnectErr).Msg("error disconnecting after a failed handshake")
		}
		return errors.Wrap(err, "error during the umas connect handshake")
	}
	return nil
}

// handshake runs the exchanges plc4j's UmasConnection.doHandshake chains together, in the same order.
// The order is capture-derived: the sequence of project-info subcodes and the two memory block reads
// are what a Schneider programming tool sends, and no document says which of them the PLC insists on.
func (c *Connection) handshake(ctx context.Context) error {
	if err := c.performPlcIdent(ctx); err != nil {
		return err
	}
	if err := c.performInitComms(ctx); err != nil {
		return err
	}
	if err := c.performRepeat(ctx); err != nil {
		return err
	}
	if err := c.performReadMemoryBlock(ctx, projectMemoryBlockNumber); err != nil {
		return err
	}
	for _, subcode := range projectInfoSubcodes.beforeSecondaryBlock {
		if err := c.performProjectInfo(ctx, subcode); err != nil {
			return err
		}
	}
	if err := c.performReadMemoryBlock(ctx, secondaryMemoryBlockNumber); err != nil {
		return err
	}
	for _, subcode := range projectInfoSubcodes.afterSecondaryBlock {
		if err := c.performProjectInfo(ctx, subcode); err != nil {
			return err
		}
	}

	// The connection counts as up once the handshake is through. plc4j marks it connected at the
	// same point and only warns if the dictionary download fails, because a connection which can
	// ping and answer "symbol not found" is more useful than no connection at all.
	c.SetConnected(true)
	c.log.Info().Stringer("session", c.session).Msg("Connected to a UMAS PLC")

	if err := c.loadDataDictionary(ctx); err != nil {
		c.log.Warn().Err(err).Msg("Failed to load the data dictionary during connect, reads and writes will not resolve symbols until a browse retries it")
	}
	return nil
}

func (c *Connection) performPlcIdent(ctx context.Context) error {
	item, err := c.requester.exchange(ctx, "PlcIdent",
		readWriteModel.NewUmasPDUPlcIdentRequest(c.session.getPairingKey()))
	if err != nil {
		return err
	}
	response, err := expectItem[readWriteModel.UmasPDUPlcIdentResponse]("PlcIdent", item)
	if err != nil {
		return err
	}
	c.session.setIdentity(response.GetHostname(), response.GetModel(), response.GetComVersion())
	c.log.Info().
		Str("hostname", response.GetHostname()).
		Uint16("model", response.GetModel()).
		Uint16("comVersion", response.GetComVersion()).
		Msg("PlcIdent")
	return nil
}

func (c *Connection) performInitComms(ctx context.Context) error {
	// plc4j sends subCode 0x00 here and the mspec has no other value for it.
	item, err := c.requester.exchange(ctx, "InitComms",
		readWriteModel.NewUmasInitCommsRequest(c.session.getPairingKey(), 0x00))
	if err != nil {
		return err
	}
	response, err := expectItem[readWriteModel.UmasInitCommsResponse]("InitComms", item)
	if err != nil {
		return err
	}
	c.session.setCommsParameters(response.GetMaxFrameSize(), response.GetFirmwareVersion())
	c.log.Info().
		Uint16("maxFrameSize", response.GetMaxFrameSize()).
		Uint16("firmwareVersion", response.GetFirmwareVersion()).
		Msg("InitComms")
	return nil
}

// performRepeat sends the echo request the handshake ends its negotiation with: a payload of
// maxFrameSize - 3 bytes which the PLC sends straight back. Ported from plc4j's performRepeat,
// filler byte and all.
func (c *Connection) performRepeat(ctx context.Context) error {
	item, err := c.requester.exchange(ctx, "Repeat",
		readWriteModel.NewUmasPDURepeatRequest(c.session.getPairingKey(), c.buildRepeatPayload()))
	if err != nil {
		return err
	}
	response, err := expectItem[readWriteModel.UmasPDURepeatResponse]("Repeat", item)
	if err != nil {
		return err
	}
	c.log.Info().Int("echoedBytes", len(response.GetBlock())).Msg("Repeat")
	return nil
}

// buildRepeatPayload is the echo payload: the first byte stays zero and every following one is the
// filler. The negotiated frame size is never below minMaxFrameSize, so the subtraction can't wrap.
func (c *Connection) buildRepeatPayload() []byte {
	payloadSize := int(c.session.getMaxFrameSize() - repeatPayloadOverhead)
	if payloadSize > maxRepeatPayloadSize {
		payloadSize = maxRepeatPayloadSize
	}
	payload := make([]byte, payloadSize)
	for i := 1; i < len(payload); i++ {
		payload[i] = repeatFillerByte
	}
	return payload
}

// performReadMemoryBlock reads one of the two memory blocks the handshake looks at. Block 0x30
// carries the hardware id and the project hashes; a response which doesn't parse or is too short is
// logged and tolerated, the way plc4j tolerates it - the reads that follow are what would fail.
func (c *Connection) performReadMemoryBlock(ctx context.Context, blockNumber uint16) error {
	name := fmt.Sprintf("ReadMemoryBlock(0x%02X)", blockNumber)
	item, err := c.requester.exchange(ctx, name, readWriteModel.NewUmasPDUReadMemoryBlockRequest(
		c.session.getPairingKey(), handshakeMemoryBlockRange, blockNumber, 0, 0, handshakeMemoryBlockLength))
	if err != nil {
		return err
	}
	response, ok := item.(readWriteModel.UmasPDUReadMemoryBlockResponse)
	if !ok {
		c.log.Warn().Str("step", name).Type("responseType", item).Msg("Unexpected response type, carrying on")
		return nil
	}
	block := response.GetBlock()
	c.log.Info().Str("step", name).Int("bytes", len(block)).Msg("ReadMemoryBlock")
	if blockNumber != projectMemoryBlockNumber {
		return nil
	}
	if len(block) < projectMemoryBlockMinimumLength {
		c.log.Warn().
			Str("step", name).
			Int("bytes", len(block)).
			Int("required", projectMemoryBlockMinimumLength).
			Msg("Block 0x30 is too short to carry the hardware id and the project hashes")
		return nil
	}
	// Layout of block 0x30, all little endian: the 9 byte UmasMemoryBlockBasicInfo header
	// (range(2) + notSure(2) + index(1) + hardwareId(4)) followed by two 32 bit hashes.
	hardwareId := binary.LittleEndian.Uint32(block[5:9])
	firstHash := binary.LittleEndian.Uint32(block[9:13])
	secondHash := binary.LittleEndian.Uint32(block[13:17])
	// The project CRC that variable reads and writes have to carry is the sum of the two hashes.
	// plc4j documents this as discovered by comparing working Schneider OPC UA Server traffic
	// against the raw block; there is no specification of it, and the wrap-around is what the 32
	// bit addition does on the wire.
	projectCrc := firstHash + secondHash
	c.session.setProjectIdentity(hardwareId, projectCrc)
	c.log.Info().
		Str("step", name).
		Str("hardwareId", fmt.Sprintf("0x%08X", hardwareId)).
		Str("firstHash", fmt.Sprintf("0x%08X", firstHash)).
		Str("secondHash", fmt.Sprintf("0x%08X", secondHash)).
		Str("projectCrc", fmt.Sprintf("0x%08X", projectCrc)).
		Msg("Read the project identity out of block 0x30")
	return nil
}

// performProjectInfo sends one of the project-info queries of the handshake. Nothing in plc4j looks
// at the payload of the response - it only logs its length - and neither does this driver; what the
// subcodes mean is not documented anywhere, only that a programming tool sends this sequence.
func (c *Connection) performProjectInfo(ctx context.Context, subcode uint8) error {
	name := fmt.Sprintf("ProjectInfo(subcode=%d)", subcode)
	item, err := c.requester.exchange(ctx, name,
		readWriteModel.NewUmasPDUProjectInfoRequest(c.session.getPairingKey(), subcode))
	if err != nil {
		return err
	}
	response, ok := item.(readWriteModel.UmasPDUProjectInfoResponse)
	if !ok {
		c.log.Warn().Str("step", name).Type("responseType", item).Msg("Unexpected response type, carrying on")
		return nil
	}
	c.log.Info().Str("step", name).Int("bytes", len(response.GetBlock())).Msg("ProjectInfo")
	return nil
}

// Close stops the pollers and drops the transport. UMAS has a "release PLC reservation" request
// (FC 0x11), but this driver never takes a reservation in the first place - plc4j's UmasConnection
// doesn't either - so there is nothing to hand back and closing the socket is the whole disconnect.
func (c *Connection) Close() error {
	c.log.Trace().Msg("Closing")
	c.stopStrayMessageDrain()
	if c.subscriber != nil {
		c.subscriber.Close()
	}
	err := c.DefaultConnection.Close()
	c.wg.Wait()
	return err
}

// startStrayMessageDrain keeps the codec's default incoming message channel empty. UMAS is strictly
// request/response and every response is matched to its request by transaction identifier, so
// anything the codec can't match - a response which arrived after its request timed out, a
// duplicate, an unsolicited packet - is pushed into that 100 slot buffer and would stay there for the
// life of the connection. Once the buffer is full the codec logs a warning per further packet, so on
// a flaky link a handful of late responses turn into permanent log noise.
func (c *Connection) startStrayMessageDrain() {
	// A second Connect on the same connection must not leave the first drain running.
	c.stopStrayMessageDrain()

	drainCtx, drainCancel := context.WithCancel(context.Background())
	drainDone := make(chan struct{})
	c.lifecycleMutex.Lock()
	c.drainCancel = drainCancel
	c.drainDone = drainDone
	c.lifecycleMutex.Unlock()

	incomingMessageChannel := c.messageCodec.GetDefaultIncomingMessageChannel()
	c.wg.Go(func() {
		// Closing this last, after the recover below, is what lets stopStrayMessageDrain wait for
		// the goroutine to really be gone.
		defer close(drainDone)
		defer func() {
			if err := recover(); err != nil {
				c.log.Error().
					Str("stack", string(debug.Stack())).
					Interface("err", err).
					Msg("panic-ed")
			}
		}()
		for {
			select {
			case <-drainCtx.Done():
				c.log.Trace().Msg("Ending the stray message drain")
				return
			case message := <-incomingMessageChannel:
				c.log.Debug().
					Type("message", message).
					Msg("Discarding a packet no request was waiting for")
			}
		}
	})
}

// stopStrayMessageDrain ends the drain and waits for its goroutine to be gone, so that a caller which
// stops the drain can rely on nothing reading the channel any more. Calling it without a running
// drain, or twice, is harmless.
func (c *Connection) stopStrayMessageDrain() {
	c.lifecycleMutex.Lock()
	drainCancel, drainDone := c.drainCancel, c.drainDone
	c.drainCancel, c.drainDone = nil, nil
	c.lifecycleMutex.Unlock()
	if drainCancel == nil {
		return
	}
	drainCancel()
	<-drainDone
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Data dictionary
////////////////////////////////////////////////////////////////////////////////////////////////////

// loadDataDictionary downloads the datatype dictionary, resolves the project's own types and then
// downloads the symbol table. Ported from plc4j's UmasConnection.loadDataDictionary.
func (c *Connection) loadDataDictionary(ctx context.Context) error {
	c.dictionaryMutex.Lock()
	defer c.dictionaryMutex.Unlock()

	c.log.Info().Msg("Loading the UMAS data dictionary (types and symbols)")
	datatypeReferences, err := c.downloadDatatypeNames(ctx)
	if err != nil {
		return errors.Wrap(err, "error downloading the datatype dictionary")
	}
	c.session.setDataTypes(datatypeReferences)
	c.log.Info().Int("datatypes", len(datatypeReferences)).Msg("Downloaded the datatype dictionary")

	// A type whose class identifier is zero is a plain alias of a primitive and needs no second
	// request; everything else is a struct or an array whose definition has to be fetched.
	for i, reference := range datatypeReferences {
		if reference.GetClassIdentifier() == 0 {
			continue
		}
		typeId := customTypeIdBase + uint16(i)
		if err := c.resolveCustomType(ctx, typeId, reference); err != nil {
			// One unresolvable type costs its own browse detail, not the whole dictionary.
			c.log.Warn().Err(err).
				Uint16("typeId", typeId).
				Str("typeName", reference.GetValue()).
				Msg("Failed to resolve a custom type")
		}
	}

	symbols, err := c.downloadSymbolTable(ctx)
	if err != nil {
		return errors.Wrap(err, "error downloading the symbol table")
	}
	c.session.setSymbols(symbols)
	c.log.Info().Int("symbols", len(symbols)).Msg("Loaded the UMAS data dictionary")
	return nil
}

func (c *Connection) downloadDatatypeNames(ctx context.Context) ([]readWriteModel.UmasDatatypeReference, error) {
	// A DD03 request carries no trailing padding, which is what the nil blank stands for; the mspec
	// makes that field present for DD02 only.
	item, err := c.requester.exchange(ctx, "BrowseDatatypeNames",
		readWriteModel.NewUmasPDUReadUnlocatedVariableNamesRequest(
			c.session.getPairingKey(), recordTypeDatatypes, dictionaryRequestIndex,
			c.session.getHardwareId(), 0, dictionaryRequestOffset, nil))
	if err != nil {
		return nil, err
	}
	response, err := expectItem[readWriteModel.UmasPDUReadUnlocatedVariableResponse]("BrowseDatatypeNames", item)
	if err != nil {
		return nil, err
	}
	if len(response.GetBlock()) == 0 {
		return nil, nil
	}
	return parseDatatypeNames(response.GetBlock())
}

// resolveCustomType fetches the definition of one of the project's own types. The first byte of the
// payload says which kind it is: arrayClassId marks an array, anything else a struct.
func (c *Connection) resolveCustomType(ctx context.Context, typeId uint16, reference readWriteModel.UmasDatatypeReference) error {
	name := "ResolveType(" + reference.GetValue() + ")"
	blank := uint16(0)
	item, err := c.requester.exchange(ctx, name,
		readWriteModel.NewUmasPDUReadUnlocatedVariableNamesRequest(
			c.session.getPairingKey(), recordTypeSymbols, dictionaryRequestIndex,
			c.session.getHardwareId(), typeId, dictionaryRequestOffset, &blank))
	if err != nil {
		return err
	}
	response, err := expectItem[readWriteModel.UmasPDUReadUnlocatedVariableResponse](name, item)
	if err != nil {
		return err
	}
	block := response.GetBlock()
	if len(block) < 2 {
		return errors.Errorf("%s: the type definition is %d bytes long", name, len(block))
	}
	if block[0] == arrayClassId {
		arrayDefinition, err := readWriteModel.UmasArrayTypeDefinitionParse(ctx, block)
		if err != nil {
			return errors.Wrapf(err, "%s: error parsing the array definition", name)
		}
		c.session.setArrayType(typeId, reference.GetValue(),
			arrayDefinition.GetElementTypeId(), arrayDefinition.GetDimensions())
		return nil
	}
	fields, err := parseUdtDefinition(block)
	if err != nil {
		return errors.Wrapf(err, "%s: error parsing the struct definition", name)
	}
	c.session.setStructType(typeId, reference.GetValue(), fields)
	return nil
}

func (c *Connection) downloadSymbolTable(ctx context.Context) ([]readWriteModel.UmasUnlocatedVariableReference, error) {
	blank := uint16(0)
	item, err := c.requester.exchange(ctx, "BrowseSymbolTable",
		readWriteModel.NewUmasPDUReadUnlocatedVariableNamesRequest(
			c.session.getPairingKey(), recordTypeSymbols, dictionaryRequestIndex,
			c.session.getHardwareId(), symbolTableBlockNumber, dictionaryRequestOffset, &blank))
	if err != nil {
		return nil, err
	}
	response, err := expectItem[readWriteModel.UmasPDUReadUnlocatedVariableResponse]("BrowseSymbolTable", item)
	if err != nil {
		return nil, err
	}
	if len(response.GetBlock()) == 0 {
		return nil, nil
	}
	return parseSymbolTable(response.GetBlock())
}

// ensureDataDictionary downloads the dictionary if the handshake couldn't. plc4j's browse does the
// same lazy retry.
func (c *Connection) ensureDataDictionary(ctx context.Context) error {
	if c.session.hasSymbols() {
		return nil
	}
	return c.loadDataDictionary(ctx)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Ping
////////////////////////////////////////////////////////////////////////////////////////////////////

// Ping asks the PLC for its status. Any answer means it is reachable, which is all a ping asks;
// plc4j's onPing treats the exchange the same way.
func (c *Connection) Ping(ctx context.Context) error {
	if c.DefaultConnection.IsInvalidated() {
		return errors.New("connection has been invalidated")
	}
	c.log.Trace().Msg("Pinging")
	item, err := c.requester.exchange(ctx, "Ping",
		readWriteModel.NewUmasPDUPlcStatusRequest(c.session.getPairingKey()))
	if err != nil {
		return errors.Wrap(err, "error pinging the PLC")
	}
	if isUmasErrorResponse(item) {
		// The PLC refused the status request but it answered, so it is very much alive. plc4j's
		// ping reports OK for anything that isn't an exception either.
		c.log.Debug().Msg("The PLC refused the status request, but answering at all means it is reachable")
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Metadata and builders
////////////////////////////////////////////////////////////////////////////////////////////////////

func (c *Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return &_default.DefaultConnectionMetadata{
		// Reading, writing and browsing are what plc4j's UmasDriver advertises (canRead, canWrite,
		// canBrowse and canPing). Subscribing works because it is emulated by polling the read path,
		// the same way plc4j's PollingSubscriptionConnectionBase does it.
		ProvidesReading:     true,
		ProvidesWriting:     true,
		ProvidesSubscribing: true,
		ProvidesBrowsing:    true,
	}
}

func (c *Connection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return spiModel.NewDefaultPlcReadRequestBuilder(
		c.GetPlcTagHandler(),
		NewReader(
			c.requester,
			c.session,
			append(c._options, options.WithCustomLogger(c.log))...,
		),
	)
}

func (c *Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return spiModel.NewDefaultPlcWriteRequestBuilder(
		c.GetPlcTagHandler(),
		c.GetPlcValueHandler(),
		NewWriter(
			c.requester,
			c.session,
			append(c._options, options.WithCustomLogger(c.log))...,
		),
	)
}

// PollingReadRequestBuilder feeds the polling subscriber, which turns every poll cycle into an
// ordinary read request.
func (c *Connection) PollingReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return c.ReadRequestBuilder()
}

func (c *Connection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	return spiModel.NewDefaultPlcSubscriptionRequestBuilder(
		c.GetPlcTagHandler(),
		c.GetPlcValueHandler(),
		c.subscriber,
	)
}

// UnsubscriptionRequestBuilder hands out the default builder: the handles carry the subscriber they
// belong to, so nothing driver specific is needed to take them back.
func (c *Connection) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	return spiModel.NewDefaultPlcUnsubscriptionRequestBuilder()
}

func (c *Connection) BrowseRequestBuilder() apiModel.PlcBrowseRequestBuilder {
	return spiModel.NewDefaultPlcBrowseRequestBuilder(c.GetPlcTagHandler(), c)
}

func (c *Connection) String() string {
	return fmt.Sprintf("umas.Connection{%s, %s}", c.configuration, c.session)
}
