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

package eip

import (
	"context"
	"encoding/binary"
	"fmt"
	"regexp"
	"runtime/debug"
	"strconv"
	"sync"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type Reader struct {
	messageCodec  spi.MessageCodec
	tm            transactions.RequestTransactionManager
	configuration Configuration
	sessionState  *SessionState

	wg sync.WaitGroup // use to track spawned go routines

	log zerolog.Logger
}

func NewReader(messageCodec spi.MessageCodec, tm transactions.RequestTransactionManager, configuration Configuration, sessionState *SessionState, _options ...options.WithOption) *Reader {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Reader{
		messageCodec:  messageCodec,
		tm:            tm,
		configuration: configuration,
		sessionState:  sessionState,

		log: customLogger,
	}
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
		var response apiModel.PlcReadResponse
		var err error
		switch {
		case m.configuration.forceUnconnectedOperation || (!m.sessionState.useMessageRouter && !m.sessionState.useConnectionManager):
			response, err = m.readWithoutMessageRouter(ctx, readRequest)
		case m.sessionState.useMessageRouter && !m.sessionState.useConnectionManager:
			response, err = m.readWithoutConnectionManager(ctx, readRequest)
		default:
			response, err = m.readWithConnectionManager(ctx, readRequest)
		}
		if err != nil {
			utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, err))
			return
		}
		utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(readRequest, response, nil))
	})
	return result
}

// readWithoutMessageRouter sends one unconnected request per tag, sequentially.
func (m *Reader) readWithoutMessageRouter(ctx context.Context, readRequest apiModel.PlcReadRequest) (apiModel.PlcReadResponse, error) {
	classSegment := readWriteModel.NewLogicalSegment(readWriteModel.NewClassID(0, 6))
	instanceSegment := readWriteModel.NewLogicalSegment(readWriteModel.NewInstanceID(0, 1))
	responseCodes := map[string]apiModel.PlcResponseCode{}
	plcValues := map[string]values.PlcValue{}
	for _, tagName := range readRequest.GetTagNames() {
		tag := readRequest.GetTag(tagName).(PlcTag)
		ansi, err := toAnsi(tag.GetTag())
		if err != nil {
			responseCodes[tagName] = apiModel.PlcResponseCode_INVALID_ADDRESS
			continue
		}
		request := readWriteModel.NewCipRRData(
			m.sessionState.sessionHandle, uint32(readWriteModel.CIPStatus_Success), m.sessionState.senderContext, 0,
			EmptyInterfaceHandle, 0,
			[]readWriteModel.TypeId{
				readWriteModel.NewNullAddressItem(),
				readWriteModel.NewUnConnectedDataItem(readWriteModel.NewCipUnconnectedRequest(
					classSegment, instanceSegment,
					readWriteModel.NewCipReadRequest(ansi, elementCount(tag)),
					m.configuration.backplane, m.configuration.slot,
				)),
			},
		)
		message, err := sendTransactedRequestAndWait(ctx, m.log, m.messageCodec, m.tm, "read", request, m.acceptsCipRRData)
		if err != nil {
			responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
			continue
		}
		service, ok := unconnectedService(message)
		if !ok {
			responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
			continue
		}
		if connectedResponse, ok := service.(readWriteModel.CipConnectedResponse); ok && connectedResponse.GetStatus() == 0x03 {
			responseCodes[tagName] = apiModel.PlcResponseCode_INVALID_ADDRESS
			continue
		}
		cipReadResponse, ok := service.(readWriteModel.CipReadResponse)
		if !ok {
			responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
			continue
		}
		code, value := decodeSingleReadResponse(m.log, tag, cipReadResponse)
		responseCodes[tagName] = code
		plcValues[tagName] = value
	}
	return spiModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil
}

// readWithoutConnectionManager sends all tags in one unconnected message-router request.
func (m *Reader) readWithoutConnectionManager(ctx context.Context, readRequest apiModel.PlcReadRequest) (apiModel.PlcReadResponse, error) {
	outerService, err := m.buildReadService(readRequest)
	if err != nil {
		return nil, err
	}
	request := readWriteModel.NewCipRRData(
		m.sessionState.sessionHandle, uint32(readWriteModel.CIPStatus_Success), m.sessionState.senderContext, 0,
		EmptyInterfaceHandle, 0,
		[]readWriteModel.TypeId{
			readWriteModel.NewNullAddressItem(),
			readWriteModel.NewUnConnectedDataItem(readWriteModel.NewCipUnconnectedRequest(
				readWriteModel.NewLogicalSegment(readWriteModel.NewClassID(0, 6)),
				readWriteModel.NewLogicalSegment(readWriteModel.NewInstanceID(0, 1)),
				outerService,
				m.configuration.backplane, m.configuration.slot,
			)),
		},
	)
	message, err := sendTransactedRequestAndWait(ctx, m.log, m.messageCodec, m.tm, "read", request, m.acceptsCipRRData)
	if err != nil {
		return nil, err
	}
	service, ok := unconnectedService(message)
	if !ok {
		return nil, errors.New("read response contains no unconnected data item")
	}
	return m.decodeReadResponse(service, readRequest)
}

// readWithConnectionManager sends all tags over the connected channel opened by ForwardOpen.
func (m *Reader) readWithConnectionManager(ctx context.Context, readRequest apiModel.PlcReadRequest) (apiModel.PlcReadResponse, error) {
	outerService, err := m.buildReadService(readRequest)
	if err != nil {
		return nil, err
	}
	request := readWriteModel.NewSendUnitData(
		m.sessionState.sessionHandle, uint32(readWriteModel.CIPStatus_Success), m.sessionState.senderContext, 0, 0,
		[]readWriteModel.TypeId{
			readWriteModel.NewConnectedAddressItem(m.sessionState.connectionId),
			readWriteModel.NewConnectedDataItem(m.sessionState.nextSequenceCount(), outerService),
		},
	)
	message, err := sendTransactedRequestAndWait(ctx, m.log, m.messageCodec, m.tm, "read", request,
		func(message spi.Message) bool {
			sendUnitData, ok := message.(readWriteModel.SendUnitData)
			return ok && sendUnitData.GetSessionHandle() == m.sessionState.sessionHandle
		})
	if err != nil {
		return nil, err
	}
	sendUnitData := message.(readWriteModel.SendUnitData)
	if len(sendUnitData.GetTypeIds()) < 2 {
		return nil, errors.New("read response contains no connected data item")
	}
	dataItem, ok := sendUnitData.GetTypeIds()[1].(readWriteModel.ConnectedDataItem)
	if !ok {
		return nil, errors.Errorf("unexpected type id in connected read response: %T", sendUnitData.GetTypeIds()[1])
	}
	return m.decodeReadResponse(dataItem.GetService(), readRequest)
}

// buildReadService assembles the per-tag CipReadRequests, wrapping more than
// one in a MultipleServiceRequest with plc4j's offset layout.
func (m *Reader) buildReadService(readRequest apiModel.PlcReadRequest) (readWriteModel.CipService, error) {
	tagNames := readRequest.GetTagNames()
	requests := make([]readWriteModel.CipService, 0, len(tagNames))
	for _, tagName := range tagNames {
		tag := readRequest.GetTag(tagName).(PlcTag)
		ansi, err := toAnsi(tag.GetTag())
		if err != nil {
			return nil, errors.Wrapf(err, "error encoding eip ansi for tag %s", tagName)
		}
		requests = append(requests, readWriteModel.NewCipReadRequest(ansi, elementCount(tag)))
	}
	return wrapMultipleServices(requests), nil
}

func (m *Reader) acceptsCipRRData(message spi.Message) bool {
	cipRRData, ok := message.(readWriteModel.CipRRData)
	return ok && cipRRData.GetSessionHandle() == m.sessionState.sessionHandle
}

// unconnectedService extracts the CipService from a CipRRData's second type id.
func unconnectedService(message spi.Message) (readWriteModel.CipService, bool) {
	cipRRData, ok := message.(readWriteModel.CipRRData)
	if !ok || len(cipRRData.GetTypeIds()) < 2 {
		return nil, false
	}
	dataItem, ok := cipRRData.GetTypeIds()[1].(readWriteModel.UnConnectedDataItem)
	if !ok {
		return nil, false
	}
	return dataItem.GetService(), true
}

func (m *Reader) decodeReadResponse(service readWriteModel.CipService, readRequest apiModel.PlcReadRequest) (apiModel.PlcReadResponse, error) {
	responseCodes := map[string]apiModel.PlcResponseCode{}
	plcValues := map[string]values.PlcValue{}
	tagNames := readRequest.GetTagNames()
	switch service := service.(type) {
	case readWriteModel.CipReadResponse:
		tagName := tagNames[0]
		tag := readRequest.GetTag(tagName).(PlcTag)
		code, value := decodeSingleReadResponse(m.log, tag, service)
		responseCodes[tagName] = code
		plcValues[tagName] = value
	case readWriteModel.MultipleServiceResponse:
		nb := int(service.GetServiceNb())
		servicesData := service.GetServicesData()
		offsets := service.GetOffsets()
		for i := 0; i < nb && i < len(tagNames); i++ {
			tagName := tagNames[i]
			tag := readRequest.GetTag(tagName).(PlcTag)
			singleService, ok := sliceService(servicesData, offsets, i, nb)
			if !ok {
				responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
				continue
			}
			cipReadResponse, ok := singleService.(readWriteModel.CipReadResponse)
			if !ok {
				responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
				continue
			}
			code, value := decodeSingleReadResponse(m.log, tag, cipReadResponse)
			responseCodes[tagName] = code
			plcValues[tagName] = value
		}
	default:
		return nil, errors.Errorf("unsupported response type %T", service)
	}
	return spiModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil
}

// sliceService re-parses service i out of a MultipleServiceResponse payload,
// bounds-checked so a lying offset table cannot panic.
func sliceService(servicesData []byte, offsets []uint16, i, nb int) (readWriteModel.CipService, bool) {
	if len(offsets) < nb {
		return nil, false
	}
	offset := int(offsets[i]) - int(offsets[0])
	var length int
	if i == nb-1 {
		length = len(servicesData) - offset
	} else {
		length = int(offsets[i+1]) - int(offsets[0]) - offset
	}
	if offset < 0 || length <= 0 || offset+length > len(servicesData) {
		return nil, false
	}
	serviceBuf := utils.NewReadBufferByteBased(servicesData[offset:offset+length],
		utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
	service, err := readWriteModel.CipServiceParseWithBuffer[readWriteModel.CipService](
		context.Background(), serviceBuf, false, uint16(length))
	if err != nil {
		return nil, false
	}
	return service, true
}

func decodeSingleReadResponse(log zerolog.Logger, tag PlcTag, response readWriteModel.CipReadResponse) (apiModel.PlcResponseCode, values.PlcValue) {
	code := decodeResponseCode(response.GetStatus())
	if code != apiModel.PlcResponseCode_OK {
		return code, nil
	}
	data := response.GetData()
	if data == nil {
		return apiModel.PlcResponseCode_INTERNAL_ERROR, nil
	}
	value, err := parsePlcValue(tag, data.GetData(), data.GetDataType())
	if err != nil {
		// Undecodable payload - don't report that as a successful read.
		log.Warn().Err(err).Str("tag", tag.GetTag()).Msg("undecodable read payload")
		return apiModel.PlcResponseCode_INTERNAL_ERROR, nil
	}
	return code, value
}

func toAnsi(tag string) ([]byte, error) {
	ctx := context.TODO()
	resourceAddressPattern := regexp.MustCompile("([.\\[\\]])*([A-Za-z_0-9]+){1}")

	segments := make([]readWriteModel.PathSegment, 0)
	lengthInBytes := uint16(0)
	submatch := resourceAddressPattern.FindAllStringSubmatch(tag, -1)
	for _, match := range submatch {
		identifier := match[2]
		qualifier := match[1]

		var newSegment readWriteModel.PathSegment
		if len(qualifier) > 0 {
			if qualifier == "[" {
				numericIdentifier, err := strconv.Atoi(identifier)
				if err != nil {
					return nil, fmt.Errorf("error parsing address %s, identifier %s couldn't be parsed to an integer", tag, identifier)
				}
				newSegment = readWriteModel.NewLogicalSegment(readWriteModel.NewMemberID(0, uint8(numericIdentifier)))
			} else {
				var pad *uint8
				if len(identifier)%2 != 0 {
					paddingValue := uint8(0)
					pad = &paddingValue
				}
				newSegment = readWriteModel.NewDataSegment(readWriteModel.NewAnsiExtendedSymbolSegment(identifier, pad))
			}
		} else {
			var pad *uint8
			if len(identifier)%2 != 0 {
				paddingValue := uint8(0)
				pad = &paddingValue
			}
			newSegment = readWriteModel.NewDataSegment(readWriteModel.NewAnsiExtendedSymbolSegment(identifier, pad))
		}
		lengthInBytes += uint16(newSegment.GetLengthInBytes(ctx))
		segments = append(segments, newSegment)
	}
	buffer := utils.NewWriteBufferByteBased(
		utils.WithInitialSizeForByteBasedBuffer(int(lengthInBytes)),
		utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	for _, segment := range segments {
		if err := segment.SerializeWithWriteBuffer(context.Background(), buffer); err != nil {
			return nil, errors.Wrap(err, "error converting tag to ansi")
		}
	}
	return buffer.GetBytes(), nil
}
