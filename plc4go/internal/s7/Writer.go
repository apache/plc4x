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
	"runtime/debug"
	"sync"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

type Writer struct {
	tpduGenerator *TpduGenerator
	messageCodec  spi.MessageCodec
	tm            transactions.RequestTransactionManager
	driverContext *DriverContext

	wg sync.WaitGroup // use to track spawned go routines

	log zerolog.Logger
}

func NewWriter(tpduGenerator *TpduGenerator, messageCodec spi.MessageCodec, tm transactions.RequestTransactionManager, driverContext *DriverContext, _options ...options.WithOption) *Writer {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Writer{
		tpduGenerator: tpduGenerator,
		messageCodec:  messageCodec,
		tm:            tm,
		driverContext: driverContext,
		log:           customLogger,
	}
}

func (m *Writer) Write(ctx context.Context, writeRequest apiModel.PlcWriteRequest) <-chan apiModel.PlcWriteRequestResult {
	// TODO: handle context
	result := make(chan apiModel.PlcWriteRequestResult, 1)
	m.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack())))
			}
		}()
		var items []writeChunkItem
		for _, item := range writeItemsOf(writeRequest) {
			s7Address, err := encodeS7Address(item.tag)
			if err != nil {
				utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Wrapf(err, "Error encoding s7 address for tag %s", item.tagName)))
				return
			}
			plcValue, err := item.valueOf(writeRequest.GetValue(item.tagName))
			if err != nil {
				utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Wrapf(err, "Error encoding value for tag %s", item.tagName)))
				return
			}
			value, err := serializePlcValue(item.tag, plcValue, m.driverContext.ControllerType)
			if err != nil {
				utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Wrapf(err, "Error encoding value for tag %s", item.tagName)))
				return
			}
			items = append(items, writeChunkItem{
				item:      item,
				parameter: readWriteModel.NewS7VarRequestParameterItemAddress(s7Address),
				payload:   value,
			})
		}

		responseCodes := map[string]apiModel.PlcResponseCode{}
		chunks, tooBig := splitWrite(items, m.driverContext.PduSize)
		for _, tagName := range tooBig {
			mergeResponseCode(responseCodes, tagName, apiModel.PlcResponseCode_INVALID_DATA)
		}
		for _, chunk := range chunks {
			// The pieces of a tag that is already known to fail are not written.
			var parameterItems []readWriteModel.S7VarRequestParameterItem
			var payloadItems []readWriteModel.S7VarPayloadDataItem
			var sent []writeChunkItem
			for _, item := range chunk {
				if code, ok := responseCodes[item.item.tagName]; ok && code == apiModel.PlcResponseCode_INVALID_DATA {
					continue
				}
				parameterItems = append(parameterItems, item.parameter)
				payloadItems = append(payloadItems, item.payload)
				sent = append(sent, item)
			}
			if len(sent) == 0 {
				continue
			}
			m.log.Trace().Int("items", len(sent)).Msg("Send write request")
			response, err := exchange(ctx, m.log, m.tpduGenerator, m.messageCodec, m.tm, "write",
				readWriteModel.NewS7ParameterWriteVarRequest(parameterItems),
				readWriteModel.NewS7PayloadWriteVarRequest(payloadItems))
			if err != nil {
				utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, err))
				return
			}
			if err := m.applyWriteResponse(response, sent, responseCodes); err != nil {
				utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Wrap(err, "Error decoding response")))
				return
			}
		}
		m.log.Trace().Msg("Returning the response")
		utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcWriteRequestResult(writeRequest,
			spiModel.NewDefaultPlcWriteResponse(writeRequest, responseCodes), nil))
	})
	return result
}

// applyWriteResponse records what the response to one request of a write says about its tags.
func (m *Writer) applyWriteResponse(response readWriteModel.S7Message, sent []writeChunkItem, responseCodes map[string]apiModel.PlcResponseCode) error {
	errorClass, errorCode, err := headerError(response)
	if err != nil {
		return err
	}
	// If the result contains any form of non-null error code, handle this instead.
	if (errorClass != 0) || (errorCode != 0) {
		responseCode := mapPlcErrorCode(m.log, errorClass, errorCode)
		for _, item := range sent {
			mergeResponseCode(responseCodes, item.item.tagName, responseCode)
		}
		return nil
	}

	payload, ok := response.GetPayload().(readWriteModel.S7PayloadWriteVarResponse)
	if !ok {
		return errors.Errorf("unexpected payload %T", response.GetPayload())
	}
	// If the numbers of items don't match, we're in big trouble as the only
	// way to know how to interpret the responses is by aligning them with the
	// items from the request as this information is not returned by the PLC.
	if len(sent) != len(payload.GetItems()) {
		return errors.New("The number of requested items doesn't match the number of returned items")
	}
	for i, item := range sent {
		mergeResponseCode(responseCodes, item.item.tagName, decodeResponseCode(payload.GetItems()[i].GetReturnCode()))
	}
	return nil
}

// writeItem is one protocol item of a write.
type writeItem struct {
	// tagName is the user tag the item writes to.
	tagName string
	// tag is the part of that tag the item writes - the tag itself unless it was split.
	tag apiModel.PlcTag
	// firstElement is the index of the tag's first element the item carries.
	firstElement uint16
	// split is set for a piece of a BOOL array, whose value is a part of the tag's list.
	split bool
}

// writeItemsOf is every protocol item the request is written with, in order.
func writeItemsOf(writeRequest apiModel.PlcWriteRequest) []writeItem {
	var items []writeItem
	for _, tagName := range writeRequest.GetTagNames() {
		items = append(items, writeItemsFor(tagName, writeRequest.GetTag(tagName))...)
	}
	return items
}

// writeItemsFor is the protocol items a tag is written with. A write of whole bytes overwrites
// every bit in them, so a BOOL array is only written packed for the bytes it covers completely;
// each bit it shares a byte with something else goes out as an item of its own, leaving the rest
// of that byte alone. Any other tag is a single item.
func writeItemsFor(tagName string, tag apiModel.PlcTag) []writeItem {
	s7Tag, ok := tag.(PlcTag)
	if !ok || s7Tag.GetDataType() != readWriteModel.TransportSize_BOOL || s7Tag.GetNumElements() <= 1 {
		return []writeItem{{tagName: tagName, tag: tag}}
	}
	var items []writeItem
	numElements := uint32(s7Tag.GetNumElements())
	firstBit := uint32(s7Tag.GetByteOffset())*8 + uint32(s7Tag.GetBitOffset())
	for element := uint32(0); element < numElements; {
		bit := firstBit + element
		count := uint32(1)
		if wholeBytes := (numElements - element) / 8; bit%8 == 0 && wholeBytes > 0 {
			count = wholeBytes * 8
		}
		piece := NewTag(s7Tag.GetMemoryArea(), s7Tag.GetBlockNumber(), uint16(bit/8), uint8(bit%8),
			uint16(count), readWriteModel.TransportSize_BOOL)
		items = append(items, writeItem{tagName: tagName, tag: piece, firstElement: uint16(element), split: true})
		element += count
	}
	return items
}

// valueOf is the part of the tag's value the item carries: all of it, unless the item is a piece
// of a BOOL array, in which case the piece's elements of the list - a single one unwrapped.
func (m writeItem) valueOf(value apiValues.PlcValue) (apiValues.PlcValue, error) {
	if !m.split {
		return value, nil
	}
	if !value.IsList() {
		return nil, errors.New("expected a list value for a multi-element BOOL tag")
	}
	list := value.GetList()
	count := m.tag.(PlcTag).GetNumElements()
	if len(list) < int(m.firstElement)+int(count) {
		return nil, errors.Errorf("expected at least %d list elements, got %d", int(m.firstElement)+int(count), len(list))
	}
	if count == 1 {
		return list[m.firstElement], nil
	}
	return spiValues.NewPlcList(list[m.firstElement : m.firstElement+count]), nil
}

func serializePlcValue(tag apiModel.PlcTag, plcValue apiValues.PlcValue, controllerType readWriteModel.ControllerType) (readWriteModel.S7VarPayloadDataItem, error) {
	s7Tag, ok := tag.(PlcTag)
	if !ok {
		return nil, errors.Errorf("Unsupported address type %t", tag)
	}
	transportSize := s7Tag.GetDataType().DataTransportSize()
	stringLength := uint16(254)
	if s7StringTag, ok := tag.(PlcStringTag); ok {
		stringLength = s7StringTag.stringLength
	}
	numElements := s7Tag.GetNumElements()
	dataProtocolId := s7Tag.GetDataType().DataProtocolId()

	if s7Tag.GetDataType() == readWriteModel.TransportSize_BYTE && numElements > 1 {
		return readWriteModel.NewS7VarPayloadDataItem(
			readWriteModel.DataTransportErrorCode_OK,
			transportSize, plcValue.GetRaw(),
		), nil
	}
	if s7Tag.GetDataType() == readWriteModel.TransportSize_BOOL && numElements > 1 {
		if !plcValue.IsList() {
			return nil, errors.New("expected a list value for a multi-element BOOL tag")
		}
		list := plcValue.GetList()
		if len(list) < int(numElements) {
			return nil, errors.Errorf("expected %d list elements, got %d", numElements, len(list))
		}
		data := make([]byte, (numElements+7)/8)
		for i := range numElements {
			if list[i].GetBool() {
				data[i/8] |= 1 << (i % 8)
			}
		}
		return readWriteModel.NewS7VarPayloadDataItem(
			readWriteModel.DataTransportErrorCode_OK,
			readWriteModel.DataTransportSize_BYTE_WORD_DWORD, data,
		), nil
	}

	var data []byte
	if numElements == 1 {
		serialized, err := readWriteModel.DataItemSerialize(plcValue, dataProtocolId, controllerType, int32(stringLength))
		if err != nil {
			return nil, errors.Wrapf(err, "Error serializing tag item of type: '%v'", s7Tag.GetDataType())
		}
		data = serialized
	} else {
		for i := range numElements {
			serialized, err := readWriteModel.DataItemSerialize(plcValue.GetIndex(uint32(i)), dataProtocolId, controllerType, int32(stringLength))
			if err != nil {
				return nil, errors.Wrapf(err, "Error serializing element %d of type: '%v'", i, s7Tag.GetDataType())
			}
			data = append(data, serialized...)
		}
	}
	return readWriteModel.NewS7VarPayloadDataItem(
		readWriteModel.DataTransportErrorCode_OK,
		transportSize, data,
	), nil
}
