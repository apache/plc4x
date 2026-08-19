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
	"runtime/debug"
	"sync"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type Writer struct {
	messageCodec  spi.MessageCodec
	tm            transactions.RequestTransactionManager
	configuration Configuration
	sessionState  *SessionState

	wg sync.WaitGroup // use to track spawned go routines

	log zerolog.Logger
}

func NewWriter(messageCodec spi.MessageCodec, tm transactions.RequestTransactionManager, configuration Configuration, sessionState *SessionState, _options ...options.WithOption) *Writer {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Writer{
		messageCodec:  messageCodec,
		tm:            tm,
		configuration: configuration,
		sessionState:  sessionState,
		log:           customLogger,
	}
}

func (m *Writer) Write(ctx context.Context, writeRequest apiModel.PlcWriteRequest) <-chan apiModel.PlcWriteRequestResult {
	m.log.Trace().Msg("Writing")
	result := make(chan apiModel.PlcWriteRequestResult, 1)
	m.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack())))
			}
		}()
		var response apiModel.PlcWriteResponse
		var err error
		switch {
		case m.configuration.forceUnconnectedOperation || (!m.sessionState.useMessageRouter && !m.sessionState.useConnectionManager):
			response, err = m.writeWithoutMessageRouter(ctx, writeRequest)
		case m.sessionState.useMessageRouter && !m.sessionState.useConnectionManager:
			response, err = m.writeWithoutConnectionManager(ctx, writeRequest)
		default:
			response, err = m.writeWithConnectionManager(ctx, writeRequest)
		}
		if err != nil {
			utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, err))
			return
		}
		utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcWriteRequestResult(writeRequest, response, nil))
	})
	return result
}

// buildWriteRequest assembles the CipWriteRequest for a single tag.
func (m *Writer) buildWriteRequest(writeRequest apiModel.PlcWriteRequest, tagName string) (readWriteModel.CipService, error) {
	tag := writeRequest.GetTag(tagName).(PlcTag)
	value := writeRequest.GetValue(tagName)
	data, err := encodeValue(value, tag.GetType())
	if err != nil {
		return nil, errors.Wrapf(err, "error encoding value for tag %s", tagName)
	}
	ansi, err := toAnsi(tag.GetTag())
	if err != nil {
		return nil, errors.Wrapf(err, "error encoding eip ansi for tag %s", tagName)
	}
	return readWriteModel.NewCipWriteRequest(ansi, tag.GetType(), elementCount(tag), data), nil
}

// writeWithoutMessageRouter sends one unconnected request per tag, sequentially.
func (m *Writer) writeWithoutMessageRouter(ctx context.Context, writeRequest apiModel.PlcWriteRequest) (apiModel.PlcWriteResponse, error) {
	classSegment := readWriteModel.NewLogicalSegment(readWriteModel.NewClassID(0, 6))
	instanceSegment := readWriteModel.NewLogicalSegment(readWriteModel.NewInstanceID(0, 1))
	responseCodes := map[string]apiModel.PlcResponseCode{}
	for _, tagName := range writeRequest.GetTagNames() {
		service, err := m.buildWriteRequest(writeRequest, tagName)
		if err != nil {
			responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
			continue
		}
		request := readWriteModel.NewCipRRData(
			m.sessionState.sessionHandle, uint32(readWriteModel.CIPStatus_Success), m.sessionState.senderContext, 0,
			EmptyInterfaceHandle, 0,
			[]readWriteModel.TypeId{
				readWriteModel.NewNullAddressItem(),
				readWriteModel.NewUnConnectedDataItem(readWriteModel.NewCipUnconnectedRequest(
					classSegment, instanceSegment,
					service,
					m.configuration.backplane, m.configuration.slot,
				)),
			},
		)
		message, err := sendTransactedRequestAndWait(ctx, m.log, m.messageCodec, m.tm, "write", request, m.acceptsCipRRData)
		if err != nil {
			responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
			continue
		}
		responseService, ok := unconnectedService(message)
		if !ok {
			responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
			continue
		}
		response, ok := responseService.(readWriteModel.CipWriteResponse)
		if !ok {
			responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
			continue
		}
		responseCodes[tagName] = decodeResponseCode(response.GetStatus())
	}
	return spiModel.NewDefaultPlcWriteResponse(writeRequest, responseCodes), nil
}

// writeWithoutConnectionManager sends all tags in one unconnected message-router request.
func (m *Writer) writeWithoutConnectionManager(ctx context.Context, writeRequest apiModel.PlcWriteRequest) (apiModel.PlcWriteResponse, error) {
	outerService, err := m.buildWriteService(writeRequest)
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
	message, err := sendTransactedRequestAndWait(ctx, m.log, m.messageCodec, m.tm, "write", request, m.acceptsCipRRData)
	if err != nil {
		return nil, err
	}
	service, ok := unconnectedService(message)
	if !ok {
		return nil, errors.New("write response contains no unconnected data item")
	}
	return decodeWriteResponse(service, writeRequest)
}

// writeWithConnectionManager sends all tags over the connected channel opened by ForwardOpen.
func (m *Writer) writeWithConnectionManager(ctx context.Context, writeRequest apiModel.PlcWriteRequest) (apiModel.PlcWriteResponse, error) {
	outerService, err := m.buildWriteService(writeRequest)
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
	message, err := sendTransactedRequestAndWait(ctx, m.log, m.messageCodec, m.tm, "write", request,
		func(message spi.Message) bool {
			sendUnitData, ok := message.(readWriteModel.SendUnitData)
			return ok && sendUnitData.GetSessionHandle() == m.sessionState.sessionHandle
		})
	if err != nil {
		return nil, err
	}
	sendUnitData := message.(readWriteModel.SendUnitData)
	if len(sendUnitData.GetTypeIds()) < 2 {
		return nil, errors.New("write response contains no connected data item")
	}
	dataItem, ok := sendUnitData.GetTypeIds()[1].(readWriteModel.ConnectedDataItem)
	if !ok {
		return nil, errors.Errorf("unexpected type id in connected write response: %T", sendUnitData.GetTypeIds()[1])
	}
	return decodeWriteResponse(dataItem.GetService(), writeRequest)
}

// buildWriteService assembles the per-tag CipWriteRequests, wrapping more than
// one in a MultipleServiceRequest with plc4j's offset layout.
func (m *Writer) buildWriteService(writeRequest apiModel.PlcWriteRequest) (readWriteModel.CipService, error) {
	tagNames := writeRequest.GetTagNames()
	requests := make([]readWriteModel.CipService, 0, len(tagNames))
	for _, tagName := range tagNames {
		service, err := m.buildWriteRequest(writeRequest, tagName)
		if err != nil {
			return nil, err
		}
		requests = append(requests, service)
	}
	return wrapMultipleServices(requests), nil
}

func (m *Writer) acceptsCipRRData(message spi.Message) bool {
	cipRRData, ok := message.(readWriteModel.CipRRData)
	return ok && cipRRData.GetSessionHandle() == m.sessionState.sessionHandle
}

func decodeWriteResponse(service readWriteModel.CipService, writeRequest apiModel.PlcWriteRequest) (apiModel.PlcWriteResponse, error) {
	responseCodes := map[string]apiModel.PlcResponseCode{}
	tagNames := writeRequest.GetTagNames()
	switch service := service.(type) {
	case readWriteModel.CipWriteResponse:
		tagName := tagNames[0]
		responseCodes[tagName] = decodeResponseCode(service.GetStatus())
	case readWriteModel.MultipleServiceResponse:
		nb := int(service.GetServiceNb())
		servicesData := service.GetServicesData()
		offsets := service.GetOffsets()
		for i := 0; i < nb && i < len(tagNames); i++ {
			tagName := tagNames[i]
			singleService, ok := sliceService(servicesData, offsets, i, nb)
			if !ok {
				responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
				continue
			}
			cipWriteResponse, ok := singleService.(readWriteModel.CipWriteResponse)
			if !ok {
				responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
				continue
			}
			responseCodes[tagName] = decodeResponseCode(cipWriteResponse.GetStatus())
		}
	default:
		return nil, errors.Errorf("unsupported response type %T", service)
	}
	return spiModel.NewDefaultPlcWriteResponse(writeRequest, responseCodes), nil
}
