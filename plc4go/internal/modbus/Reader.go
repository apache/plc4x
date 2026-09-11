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

package modbus

import (
	"context"
	"math"
	"runtime/debug"
	"sync"
	"sync/atomic"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/apache/plc4x/plc4go/spi/values"
)

type Reader struct {
	transactionIdentifier int32
	configuration         Configuration
	messageCodec          spi.MessageCodec
	// tm is the connection's request transaction manager, which the blocks go out through so that
	// they queue behind the connection's writes and ping rather than overtaking them.
	tm transactions.RequestTransactionManager

	wg sync.WaitGroup // use to track spawned go routines

	passLogToModel bool
	log            zerolog.Logger
}

func NewReader(configuration Configuration, messageCodec spi.MessageCodec, tm transactions.RequestTransactionManager, _options ...options.WithOption) *Reader {
	passLoggerToModel, _ := options.ExtractPassLoggerToModel(_options...)
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Reader{
		transactionIdentifier: 0,
		configuration:         configuration,
		messageCodec:          messageCodec,
		tm:                    tm,
		passLogToModel:        passLoggerToModel,
		log:                   customLogger,
	}
}

// blockRequestPdu is the request that reads one run of addresses in one area (plc4j
// ModbusTcpConnection.getReadRequestPdu, which builds the same PDUs from a single tag).
//
// quantity counts what the area addresses: coils in the two bit areas, registers in the other
// three. That is the same count a tag's own read carries, which is why a block covering exactly
// one tag produces byte for byte the request that tag was read with before anything was merged.
func blockRequestPdu(tagType TagType, address uint16, quantity uint16) (readWriteModel.ModbusPDU, error) {
	switch tagType {
	case Coil:
		return readWriteModel.NewModbusPDUReadCoilsRequest(address, quantity), nil
	case DiscreteInput:
		return readWriteModel.NewModbusPDUReadDiscreteInputsRequest(address, quantity), nil
	case InputRegister:
		return readWriteModel.NewModbusPDUReadInputRegistersRequest(address, quantity), nil
	case HoldingRegister:
		return readWriteModel.NewModbusPDUReadHoldingRegistersRequest(address, quantity), nil
	case ExtendedRegister:
		// The extended register area is read with FC 0x14, which addresses it as a set of files
		// rather than flat, so the address turns into one item per file it touches (plc4j
		// ModbusTcpConnection.getReadRequestPdu).
		groups := splitExtendedRegister(address, quantity)
		items := make([]readWriteModel.ModbusPDUReadFileRecordRequestItem, len(groups))
		for i, group := range groups {
			items[i] = readWriteModel.NewModbusPDUReadFileRecordRequestItem(
				extendedRegisterReferenceType, group.fileNumber, group.recordNumber, group.lengthWords)
		}
		return readWriteModel.NewModbusPDUReadFileRecordRequest(items), nil
	default:
		return nil, errors.Errorf("unsupported tag type %x", tagType)
	}
}

// readRequestPdu is the request that reads what one tag addresses. Ping is what is left that reads
// a lone tag; a read request goes through the optimizer and asks for blocks instead.
func readRequestPdu(tag modbusTag) (readWriteModel.ModbusPDU, error) {
	// The same count the optimizer measures a block in, so that this and blockRequestPdu can't
	// drift apart. tagSpan doesn't police the 16 bit quantity field, so that is checked here the
	// way lengthWords used to check it.
	span := tagSpan(tag)
	if span > maxWireAddress {
		return nil, errors.Errorf("the requested %d registers don't fit into a request, at most %d do", span, maxWireAddress)
	}
	return blockRequestPdu(tag.TagType, tag.Address, uint16(span))
}

func (m *Reader) Read(ctx context.Context, readRequest apiModel.PlcReadRequest) <-chan apiModel.PlcReadRequestResult {
	m.log.Trace().Msg("Reading")
	result := make(chan apiModel.PlcReadRequestResult, 1)
	m.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack())))
			}
		}()
		response, err := m.read(ctx, readRequest)
		utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(readRequest, response, err))
	})
	return result
}

// read answers the whole request: the tags are merged into the blocks that have to be read, the
// blocks are read one after another, and every block's response is split back apart over the tags
// it covers.
//
// The blocks go out one at a time rather than all at once. Only the TCP flavor carries a
// transaction identifier, so the serial flavors correlate a response by station address and
// function code alone (see rtuAduFactory.acceptsResponse) and two overlapping reads of the same
// area would be indistinguishable. plc4j chains its requests for the same reason.
//
// A response is returned even when some of it failed, with the failures reported both as the
// response code of the tags they took down and as the returned error, the way the single-item
// interceptor reports a partial result. A read of twenty tags where one register is unreadable is
// still nineteen values the caller asked for.
func (m *Reader) read(ctx context.Context, readRequest apiModel.PlcReadRequest) (apiModel.PlcReadResponse, error) {
	tags, invalid := modbusTagsOf(readRequest)
	if len(tags)+len(invalid) == 0 {
		return nil, errors.New("a read request has to ask for at least one tag")
	}

	responseCodes := make(map[string]apiModel.PlcResponseCode, len(tags)+len(invalid))
	plcValues := make(map[string]apiValues.PlcValue, len(tags)+len(invalid))
	var failures []error
	for _, entry := range invalid {
		// A tag this driver can't make sense of fails on its own; there is nothing to ask the
		// device about it, and the rest of the request is still answerable.
		m.log.Debug().Err(entry.err).Str("tagName", entry.name).Msg("Invalid tag")
		responseCodes[entry.name] = apiModel.PlcResponseCode_INVALID_ADDRESS
		// A null rather than nothing, for the same reason a tag whose block failed carries one:
		// the response hands out whatever it was given, so a missing value is a nil to
		// dereference.
		plcValues[entry.name] = values.NewPlcNULL()
		failures = append(failures, entry.err)
	}

	ctxForModel := options.GetLoggerContextForModel(ctx, m.log, options.WithPassLoggerToModel(m.passLogToModel))
	// Blocks go out ONE AT A TIME, and that is the point rather than an oversight.
	//
	// The path this replaces handed every per-tag sub-read to the dispatcher at once, spaced only
	// by a hardcoded 4ms sleep and with no in-flight cap, so a read of N tags put N transactions
	// on the wire together. A direct-TCP device shrugs that off; a TCP<->RTU gateway, which has to
	// serialise each one onto a slow serial line, resets the connection instead - and the reset
	// takes every in-flight request with it, so the read comes back with nothing at all.
	//
	// plc4j reaches the same place from the other direction: its Modbus connections override
	// getMaxConcurrentRequests() to 1, so one PDU is on the wire at a time per connection.
	//
	// The cost is latency for tags the optimizer cannot merge: worst case N round trips where the
	// old path overlapped them. Merging is what buys that back - a profile of 84 tags spanning a
	// few register clusters becomes a handful of blocks, not 84 of anything.
	for _, block := range optimizeReads(tags, m.configuration.maxCoilsPerRequest, m.configuration.maxRegistersPerRequest) {
		// A cancelled read must stop here rather than walk the remaining blocks. Blocks are
		// issued one at a time, so without this a request cancelled after the first block still
		// puts every later block on the wire -- which is worse than the path this replaced,
		// where the sub-reads were already in flight by the time cancellation could be noticed.
		if err := ctx.Err(); err != nil {
			failures = append(failures, err)
			break
		}

		outcome := m.readBlock(ctx, block)
		if outcome.err != nil {
			m.log.Debug().Err(outcome.err).
				Stringer("tagType", block.tagType).
				Uint16("address", block.address).
				Uint16("quantity", block.quantity).
				Msg("Reading a block failed, so every tag in it fails")
			failures = append(failures, outcome.err)
		}
		for _, tagResult := range splitBlockResponse(ctxForModel, block, outcome.responseCode, outcome.data, m.configuration.defaultPayloadByteOrder) {
			responseCodes[tagResult.name] = tagResult.responseCode
			// A tag that failed still gets an entry, so that a caller reading the value of a
			// failed tag finds a null rather than a nil to dereference (as s7's reader does).
			plcValues[tagResult.name] = valueOrNull(tagResult)
			if tagResult.err != nil {
				// The response carries only the code, so this is the one place the reason is
				// known.
				m.log.Debug().Err(tagResult.err).Str("tagName", tagResult.name).Msg("Couldn't answer a tag out of the block's response")
				failures = append(failures, tagResult.err)
			}
		}
	}

	m.log.Trace().Msg("Returning the response")
	return spiModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), errors.Join(failures...)
}

// invalidTag is a tag of the request that isn't a modbus tag at all, kept in the order the request
// named it so that the failures are reported in a stable one.
type invalidTag struct {
	name string
	err  error
}

// modbusTagsOf reads the request's tags, separating the ones this driver can work with from the
// ones it can't.
func modbusTagsOf(readRequest apiModel.PlcReadRequest) ([]namedTag, []invalidTag) {
	tagNames := readRequest.GetTagNames()
	tags := make([]namedTag, 0, len(tagNames))
	var invalid []invalidTag
	for _, tagName := range tagNames {
		tag := readRequest.GetTag(tagName)
		if unparsed, ok := tag.(unparsedTag); ok {
			// An address the builder couldn't parse, kept in the request so that it can be
			// reported rather than taking the request down (see ReadRequestBuilder.go). Its own
			// error says what was wrong with the address, which a failed cast wouldn't.
			invalid = append(invalid, invalidTag{name: tagName, err: errors.Wrapf(unparsed.err, "invalid address for '%s'", tagName)})
			continue
		}
		modbusTagVar, err := castToModbusTagFromPlcTag(tag)
		if err != nil {
			invalid = append(invalid, invalidTag{name: tagName, err: errors.Wrapf(err, "invalid tag item type for '%s'", tagName)})
			continue
		}
		tags = append(tags, namedTag{name: tagName, tag: modbusTagVar})
	}
	return tags, invalid
}

// valueOrNull is the value a tag ended up with; a tag that failed carries a null rather than
// nothing, so that a caller who reads the value before looking at the code gets a PlcValue either
// way.
func valueOrNull(result tagReadResult) apiValues.PlcValue {
	if result.value == nil {
		return values.NewPlcNULL()
	}
	return result.value
}

// blockOutcome is how one block read ended: the payload the device answered with, or the code
// every tag of the block gets instead. err is why the code is not OK, kept for logging - the
// response itself carries only the code.
type blockOutcome struct {
	responseCode apiModel.PlcResponseCode
	data         []byte
	err          error
}

// readBlock sends one block's request and waits for its answer.
func (m *Reader) readBlock(ctx context.Context, block readBlock) blockOutcome {
	pdu, err := blockRequestPdu(block.tagType, block.address, block.quantity)
	if err != nil {
		m.log.Debug().Err(err).Stringer("tagType", block.tagType).Msg("Couldn't build a read request")
		return blockOutcome{responseCode: apiModel.PlcResponseCode_UNSUPPORTED, err: err}
	}

	// Calculate a new transaction identifier
	transactionIdentifier := atomic.AddInt32(&m.transactionIdentifier, 1)
	if transactionIdentifier > math.MaxUint8 {
		transactionIdentifier = 1
		atomic.StoreInt32(&m.transactionIdentifier, 1)
	}
	m.log.Debug().Int32("transactionIdentifier", transactionIdentifier).Msg("Calculated transaction identifier")

	// Assemble the finished ADU. Which one that is depends on the flavor the connection
	// speaks - the TCP one carries the transaction identifier, the RTU one a CRC.
	m.log.Trace().Msg("Assemble ADU")
	adus := m.configuration.adus()
	unitIdentifier := m.configuration.unitIdentifier
	if block.unitId != nil {
		// Every tag of the block named this unit; a block read reaches exactly one device.
		unitIdentifier = *block.unitId
	}
	requestAdu := adus.buildRequest(uint16(transactionIdentifier), unitIdentifier, pdu)

	// A request that nobody answers must not wait forever; the codec derives the lifetime of the
	// expectation it registers from the deadline of the context it is handed. Cancelling once the
	// block is done releases the timer and drops an expectation that is somehow still registered,
	// which is safe here because this waits for the answer rather than returning ahead of it.
	requestCtx, cancelRequest := withRequestTimeout(ctx, m.configuration.requestTimeout)
	defer cancelRequest()

	// Buffered and written to without blocking, so that a handler firing for a block that is
	// already accounted for returns straight away - a codec that reports a send failure and then
	// times the expectation out anyway would otherwise leave that handler blocked forever, and
	// blocked handlers pile up in the codec's WaitGroup and wedge Disconnect.
	outcomes := make(chan blockOutcome, 1)
	complete := func(outcome blockOutcome) {
		select {
		case outcomes <- outcome:
		default:
			m.log.Debug().Msg("dropping a second outcome for a block that is already answered")
		}
	}

	m.log.Trace().Msg("Send ADU")
	if err := sendTransacted(requestCtx, m.log, m.messageCodec, m.tm, "read", requestAdu, func(message spi.Message) bool {
		return adus.acceptsResponse(requestAdu, message)
	}, func(message spi.Message) error {
		// Convert the response into an ADU
		m.log.Trace().Msg("convert response to ADU")
		responsePdu, err := adus.extractPdu(message)
		if err != nil {
			complete(blockOutcome{responseCode: apiModel.PlcResponseCode_INTERNAL_ERROR, err: err})
			return nil
		}
		data, responseCode, err := extractResponseData(responsePdu)
		complete(blockOutcome{responseCode: responseCode, data: data, err: err})
		return nil
	}, func(err error) error {
		complete(blockOutcome{
			responseCode: apiModel.PlcResponseCode_REQUEST_TIMEOUT,
			err:          errors.Wrap(err, "got timeout while waiting for response"),
		})
		return nil
	}); err != nil {
		// Nothing reached the wire. A block that ran out of time while it was still waiting its
		// turn behind another request timed out just as surely as one the device never answered,
		// so it is reported as a timeout; anything else is the codec refusing to send.
		//
		// NOT PINNED BY A TEST. Deleting the choice below and always reporting INTERNAL_ERROR
		// leaves the package green, because the obvious way to provoke it - queue a block behind a
		// slow one until its deadline passes - is noticed by the outer wait further down first,
		// which reports REQUEST_TIMEOUT of its own. Reaching this line needs sendTransacted to
		// return while requestCtx is already expired. Three sites in this function now produce
		// REQUEST_TIMEOUT under conditions that overlap, which is the actual smell; collapsing
		// them into one decision would be worth more than a test that reaches this one.
		responseCode := apiModel.PlcResponseCode_INTERNAL_ERROR
		if requestCtx.Err() != nil {
			responseCode = apiModel.PlcResponseCode_REQUEST_TIMEOUT
		}
		complete(blockOutcome{
			responseCode: responseCode,
			err:          errors.Wrap(err, "error sending message"),
		})
	}

	select {
	case outcome := <-outcomes:
		return outcome
	case <-requestCtx.Done():
		// The codec answers every expectation it registers, so this is the caller giving up rather
		// than the device being slow. An outcome that landed in the same instant is preferred over
		// it - select picks at random between two ready cases, and a delivered answer is the
		// truthful one.
		select {
		case outcome := <-outcomes:
			return outcome
		default:
		}
		return blockOutcome{responseCode: apiModel.PlcResponseCode_REQUEST_TIMEOUT, err: requestCtx.Err()}
	}
}

// extractResponseData pulls the payload out of a read response: the packed bits for the two bit
// areas, the registers for the other three. The response code is what every tag the request
// covered gets; it is OK exactly when there is a payload to answer them out of.
func extractResponseData(responsePdu readWriteModel.ModbusPDU) ([]byte, apiModel.PlcResponseCode, error) {
	switch pdu := responsePdu.(type) {
	case readWriteModel.ModbusPDUReadDiscreteInputsResponse:
		return pdu.GetValue(), apiModel.PlcResponseCode_OK, nil
	case readWriteModel.ModbusPDUReadCoilsResponse:
		return pdu.GetValue(), apiModel.PlcResponseCode_OK, nil
	case readWriteModel.ModbusPDUReadInputRegistersResponse:
		return pdu.GetValue(), apiModel.PlcResponseCode_OK, nil
	case readWriteModel.ModbusPDUReadHoldingRegistersResponse:
		return pdu.GetValue(), apiModel.PlcResponseCode_OK, nil
	case readWriteModel.ModbusPDUReadFileRecordResponse:
		// A request that crossed a file boundary was sent as one item per file, and the response
		// carries one item per request item, so the value is the registers of all of them in
		// order. plc4j's extractResponseData only looks at the first item, which silently drops
		// everything behind the boundary.
		var data []byte
		for _, item := range pdu.GetItems() {
			data = append(data, item.GetData()...)
		}
		return data, apiModel.PlcResponseCode_OK, nil
	case readWriteModel.ModbusPDUError:
		responseCode, mapped := responseCodeOf(pdu.GetExceptionCode())
		switch {
		case !mapped:
			// An exception code the specification doesn't define. The device refused; which way it
			// spelled that is all that is unknown.
			responseCode = apiModel.PlcResponseCode_REMOTE_ERROR
		case responseCode == apiModel.PlcResponseCode_OK:
			// ACKNOWLEDGE is the one exception plc4j maps to OK: the device took the request and
			// will take its time over it. For a write that is a fair answer, for a read it is not
			// - an exception response carries no registers, so OK would promise the caller a value
			// that never arrived.
			responseCode = apiModel.PlcResponseCode_REMOTE_BUSY
		}
		return nil, responseCode, errors.Errorf("got an error from remote. Errorcode %x", pdu.GetExceptionCode())
	default:
		return nil, apiModel.PlcResponseCode_INTERNAL_ERROR, errors.Errorf("unsupported response type %T", pdu)
	}
}

// ToPlc4xReadResponse turns the ADU a device answered with into a PLC4X response. It takes the
// ADU through the little bit of it this needs, so that it serves every flavor - the generated
// ModbusADU parent type carries no PDU accessor of its own.
//
// One ADU answers one block, so this only serves a request the optimizer merges into a single one.
func (m *Reader) ToPlc4xReadResponse(responseAdu aduWithPdu, readRequest apiModel.PlcReadRequest) (apiModel.PlcReadResponse, error) {
	return m.toPlc4xReadResponse(responseAdu.GetPdu(), readRequest)
}

func (m *Reader) toPlc4xReadResponse(responsePdu readWriteModel.ModbusPDU, readRequest apiModel.PlcReadRequest) (apiModel.PlcReadResponse, error) {
	tags, invalid := modbusTagsOf(readRequest)
	if len(invalid) > 0 {
		return nil, invalid[0].err
	}
	blocks := optimizeReads(tags, m.configuration.maxCoilsPerRequest, m.configuration.maxRegistersPerRequest)
	if len(blocks) != 1 {
		return nil, errors.Errorf("this request is read as %d blocks, so a single response can't answer it", len(blocks))
	}
	data, responseCode, err := extractResponseData(responsePdu)
	if err != nil {
		return nil, err
	}

	// Decode the data according to the information from the request
	m.log.Trace().Msg("decode data")
	ctxForModel := options.GetLoggerContextForModel(context.TODO(), m.log, options.WithPassLoggerToModel(m.passLogToModel))
	responseCodes := map[string]apiModel.PlcResponseCode{}
	plcValues := map[string]apiValues.PlcValue{}
	var failures []error
	for _, tagResult := range splitBlockResponse(ctxForModel, blocks[0], responseCode, data, m.configuration.defaultPayloadByteOrder) {
		responseCodes[tagResult.name] = tagResult.responseCode
		plcValues[tagResult.name] = valueOrNull(tagResult)
		if tagResult.err != nil {
			failures = append(failures, tagResult.err)
		}
	}

	// Return the response
	m.log.Trace().Msg("Returning the response")
	return spiModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), errors.Join(failures...)
}
