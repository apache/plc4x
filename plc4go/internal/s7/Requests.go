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

package s7

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/transactions"
)

// The size of the S7 messages around the items, as the PDU size counts them (ported from plc4j's
// S7Optimizer).
var (
	emptyReadRequestSize = readWriteModel.NewS7MessageRequest(0,
		readWriteModel.NewS7ParameterReadVarRequest(nil), nil).GetLengthInBytes(context.Background())
	emptyReadResponseSize = readWriteModel.NewS7MessageResponseData(0,
		readWriteModel.NewS7ParameterReadVarResponse(0),
		readWriteModel.NewS7PayloadReadVarResponse(nil), 0, 0).GetLengthInBytes(context.Background())
	emptyWriteRequestSize = readWriteModel.NewS7MessageRequest(0,
		readWriteModel.NewS7ParameterWriteVarRequest(nil),
		readWriteModel.NewS7PayloadWriteVarRequest(nil)).GetLengthInBytes(context.Background())
	emptyWriteResponseSize = readWriteModel.NewS7MessageResponseData(0,
		readWriteModel.NewS7ParameterWriteVarResponse(0),
		readWriteModel.NewS7PayloadWriteVarResponse(nil), 0, 0).GetLengthInBytes(context.Background())
	addressItemSize = readWriteModel.NewS7VarRequestParameterItemAddress(
		readWriteModel.NewS7AddressAny(readWriteModel.TransportSize_INT, 1, 1, readWriteModel.MemoryArea_DATA_BLOCKS, 1, 0),
	).GetLengthInBytes(context.Background())
)

// readItem is one item of a read request.
type readItem struct {
	// tagName is the user tag the item reads.
	tagName string
	address readWriteModel.S7AddressAny
	// split is set when the tag did not fit into one PDU and is read in byte fragments; the
	// fragments are reassembled into the tag's payload before it is decoded.
	split bool
	// fragmentByteOffset is where the item's data goes in that payload.
	fragmentByteOffset uint32
}

// payloadBytes is how many bytes of data a read of the address returns.
func payloadBytes(address readWriteModel.S7AddressAny) uint32 {
	return uint32(address.GetNumberOfElements()) * uint32(address.GetTransportSize().SizeInBytes())
}

func roundUpEven(n uint64) uint64 {
	return n + n%2
}

// splitRead groups the tags of a read into requests that fit the PDU, in request order. A tag
// that does not fit into one PDU on its own is read as byte fragments of its payload, one request
// each.
func splitRead(readRequest apiModel.PlcReadRequest, pduSize uint16) ([][]readItem, error) {
	var chunks [][]readItem
	var current []readItem
	requestSize := emptyReadRequestSize
	responseSize := emptyReadResponseSize
	flush := func() {
		if len(current) > 0 {
			chunks = append(chunks, current)
			current = nil
			requestSize = emptyReadRequestSize
			responseSize = emptyReadResponseSize
		}
	}
	for _, tagName := range readRequest.GetTagNames() {
		tag, ok := readRequest.GetTag(tagName).(PlcTag)
		if !ok {
			return nil, errors.Errorf("Error encoding s7 address for tag %s: unsupported address type %T", tagName, readRequest.GetTag(tagName))
		}
		address := readAddressOf(tag)
		total := payloadBytes(address)
		responseItemSize := roundUpEven(4 + uint64(total))
		if responseItemSize <= uint64(pduSize)-emptyReadResponseSize {
			if requestSize+addressItemSize > uint64(pduSize) || responseSize+responseItemSize > uint64(pduSize) {
				flush()
			}
			current = append(current, readItem{tagName: tagName, address: address})
			requestSize += addressItemSize
			responseSize += responseItemSize
			continue
		}
		flush()
		// Four bytes of item header, and the payload is padded to an even length.
		bytesPerRequest := max(int64(pduSize)-int64(emptyReadResponseSize)-5, 1)
		for offset := uint32(0); offset < total; offset += uint32(bytesPerRequest) {
			n := min(total-offset, uint32(bytesPerRequest))
			fragment := readWriteModel.NewS7AddressAny(readWriteModel.TransportSize_BYTE, uint16(n),
				address.GetDbNumber(), address.GetArea(), address.GetByteAddress()+uint16(offset), 0)
			chunks = append(chunks, []readItem{{tagName: tagName, address: fragment, split: true, fragmentByteOffset: offset}})
		}
	}
	flush()
	return chunks, nil
}

// writeChunkItem is one item of a write request, as it goes on the wire.
type writeChunkItem struct {
	item      writeItem
	parameter readWriteModel.S7VarRequestParameterItem
	payload   readWriteModel.S7VarPayloadDataItem
}

// splitWrite groups the items of a write into requests that fit the PDU, in request order. An
// item that does not fit into one PDU on its own is returned in tooBig instead, by tag name: a
// write cannot be split without changing what it writes.
func splitWrite(items []writeChunkItem, pduSize uint16) (chunks [][]writeChunkItem, tooBig []string) {
	var current []writeChunkItem
	requestSize := emptyWriteRequestSize
	responseSize := emptyWriteResponseSize
	for _, item := range items {
		requestItemSize := roundUpEven(item.parameter.GetLengthInBytes(context.Background()) +
			4 + uint64(len(item.payload.GetData())))
		const responseItemSize = 1
		if emptyWriteRequestSize+requestItemSize > uint64(pduSize) {
			tooBig = append(tooBig, item.item.tagName)
			continue
		}
		if requestSize+requestItemSize > uint64(pduSize) || responseSize+responseItemSize > uint64(pduSize) {
			chunks = append(chunks, current)
			current = nil
			requestSize = emptyWriteRequestSize
			responseSize = emptyWriteResponseSize
		}
		current = append(current, item)
		requestSize += requestItemSize
		responseSize += responseItemSize
	}
	if len(current) > 0 {
		chunks = append(chunks, current)
	}
	return chunks, tooBig
}

// mergeResponseCode records the code of one item for its tag. A tag read or written in several
// items is only OK if all of them were, so the first failure sticks.
func mergeResponseCode(responseCodes map[string]apiModel.PlcResponseCode, tagName string, code apiModel.PlcResponseCode) {
	if current, ok := responseCodes[tagName]; !ok || current == apiModel.PlcResponseCode_OK {
		responseCodes[tagName] = code
	}
}

// exchange sends one S7 request and waits for the response carrying its TPDU reference, so that a
// request split into several can be sent one part after the other.
func exchange(ctx context.Context, log zerolog.Logger, tpduGenerator *TpduGenerator, messageCodec spi.MessageCodec,
	tm transactions.RequestTransactionManager, name string, parameter readWriteModel.S7Parameter,
	payload readWriteModel.S7Payload) (readWriteModel.S7Message, error) {
	tpduId := tpduGenerator.getAndIncrement()
	// TODO: why do we use a uint16 above and the cotp a uint8?
	tpktPacket := readWriteModel.NewTPKTPacket(
		readWriteModel.NewCOTPPacketData(
			nil,
			readWriteModel.NewS7MessageRequest(tpduId, parameter, payload),
			true,
			uint8(tpduId),
		),
	)
	type outcome struct {
		message readWriteModel.S7Message
		err     error
	}
	result := make(chan outcome, 1)
	// Start a new request-transaction (Is ended in the response-handler)
	transaction := tm.StartTransaction(name)
	transaction.Submit(name+"Operation", func(transactionContext context.Context, transaction transactions.RequestTransaction) {
		ctx, cancel := context.WithCancel(ctx)
		context.AfterFunc(transactionContext, cancel)
		if err := messageCodec.SendRequest(ctx, name, tpktPacket, func(message spi.Message) bool {
			tpktPacket, ok := message.(readWriteModel.TPKTPacket)
			if !ok {
				return false
			}
			cotpPacketData, ok := tpktPacket.GetPayload().(readWriteModel.COTPPacketData)
			if !ok {
				return false
			}
			payload := cotpPacketData.GetPayload()
			if payload == nil {
				return false
			}
			return payload.GetTpduReference() == tpduId
		}, func(message spi.Message) error {
			cotpPacketData := message.(readWriteModel.TPKTPacket).GetPayload().(readWriteModel.COTPPacketData)
			result <- outcome{message: cotpPacketData.GetPayload()}
			return transaction.EndRequest()
		}, func(err error) error {
			result <- outcome{err: errors.Wrap(err, "got timeout while waiting for response")}
			return transaction.EndRequest()
		}); err != nil {
			result <- outcome{err: errors.Wrap(err, "error sending message")}
			if err := transaction.FailRequest(errors.Errorf("timeout after %s", 1*time.Second)); err != nil {
				log.Debug().Err(err).Msg("Error failing request")
			}
		}
	})
	select {
	case r := <-result:
		return r.message, r.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// headerError is the error class and code of a response, which refer to the whole request.
func headerError(response readWriteModel.S7Message) (errorClass uint8, errorCode uint8, err error) {
	switch messageResponseData := response.(type) {
	case readWriteModel.S7MessageResponseData:
		return messageResponseData.GetErrorClass(), messageResponseData.GetErrorCode(), nil
	case readWriteModel.S7MessageResponse:
		return messageResponseData.GetErrorClass(), messageResponseData.GetErrorCode(), nil
	default:
		return 0, 0, errors.Errorf("unsupported response type %T", response)
	}
}
