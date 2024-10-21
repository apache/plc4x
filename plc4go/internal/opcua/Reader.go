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

package opcua

import (
	"context"
	"encoding/binary"
	"runtime/debug"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/opcua/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type Reader struct {
	connection *Connection

	log zerolog.Logger
}

func NewReader(connection *Connection, _options ...options.WithOption) *Reader {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Reader{
		connection: connection,

		log: customLogger,
	}
}

func (m *Reader) Read(ctx context.Context, readRequest apiModel.PlcReadRequest) <-chan apiModel.PlcReadRequestResult {
	m.log.Trace().Msg("Reading")
	result := make(chan apiModel.PlcReadRequestResult, 1)
	go m.readSync(ctx, readRequest, result)
	return result
}

func (m *Reader) readSync(ctx context.Context, readRequest apiModel.PlcReadRequest, result chan apiModel.PlcReadRequestResult) {
	defer func() {
		if err := recover(); err != nil {
			result <- spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()))
		}
	}()

	requestHeader := readWriteModel.NewRequestHeader(
		m.connection.channel.getAuthenticationToken(),
		m.connection.channel.getCurrentDateTime(),
		m.connection.channel.getRequestHandle(),
		0,
		NULL_STRING,
		REQUEST_TIMEOUT_LONG,
		NULL_EXTENSION_OBJECT,
	)
	readValueArray := make([]readWriteModel.ReadValueId, len(readRequest.GetTagNames()))
	for i, tagName := range readRequest.GetTagNames() {
		tag := readRequest.GetTag(tagName).(Tag)

		nodeId, err := generateNodeId(tag)
		if err != nil {
			result <- spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, errors.Wrapf(err, "error generating node id from tag %s", tag))
			return
		}

		readValueArray[i] = readWriteModel.NewReadValueId(nodeId,
			0xD,
			NULL_STRING,
			readWriteModel.NewQualifiedName(0, NULL_STRING),
		)
	}

	opcuaReadRequest := readWriteModel.NewReadRequest(
		requestHeader,
		0.0,
		readWriteModel.TimestampsToReturn_timestampsToReturnNeither,
		readValueArray)

	identifier := opcuaReadRequest.GetExtensionId()
	expandedNodeId := readWriteModel.NewExpandedNodeId(false, //Namespace Uri Specified
		false, //Server Index Specified
		readWriteModel.NewNodeIdFourByte(0, uint16(identifier)),
		nil,
		nil)

	extObject := readWriteModel.NewExtensiblePayload(nil, readWriteModel.NewRootExtensionObject(expandedNodeId, opcuaReadRequest, identifier), 0)

	buffer := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	if err := extObject.SerializeWithWriteBuffer(ctx, buffer); err != nil {
		result <- spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, errors.Wrapf(err, "Unable to serialise the ReadRequest"))
		return
	}

	consumer := func(opcuaResponse []byte) {
		reply, err := readWriteModel.ExtensionObjectParseWithBuffer[readWriteModel.ExtensionObject](ctx, utils.NewReadBufferByteBased(opcuaResponse, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian)), false)
		if err != nil {
			result <- spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, errors.Wrapf(err, "Unable to read the reply"))
			return
		}
		extensionObjectDefinition := reply.GetBody()
		if _readResponse, ok := extensionObjectDefinition.(readWriteModel.ReadResponse); ok {
			result <- spiModel.NewDefaultPlcReadRequestResult(readRequest, spiModel.NewDefaultPlcReadResponse(readResponse(m.log, readRequest, readRequest.GetTagNames(), _readResponse.GetResults())), nil)
			return
		} else {
			if serviceFault, ok := extensionObjectDefinition.(readWriteModel.ServiceFault); ok {
				header := serviceFault.GetResponseHeader()
				m.log.Error().Stringer("header", header).Msg("Read request ended up with ServiceFault")
			} else {
				m.log.Error().Stringer("extensionObjectDefinition", extensionObjectDefinition).Msg("Remote party returned an error")
			}

			responseCodes := map[string]apiModel.PlcResponseCode{}
			for _, tagName := range readRequest.GetTagNames() {
				responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
			}
			result <- spiModel.NewDefaultPlcReadRequestResult(readRequest, spiModel.NewDefaultPlcReadResponse(readRequest, responseCodes, nil), nil)
		}
	}

	errorDispatcher := func(err error) {
		result <- spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, err)
	}

	m.connection.channel.submit(ctx, m.connection.messageCodec, errorDispatcher, consumer, buffer)
}
