/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package knxnetip

import (
	"errors"
	"github.com/apache/plc4x/plc4go/internal/spi"
	plc4goModel "github.com/apache/plc4x/plc4go/internal/spi/model"
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
)

type Writer struct {
	messageCodec spi.MessageCodec
}

func NewWriter(messageCodec spi.MessageCodec) Writer {
	return Writer{
		messageCodec: messageCodec,
	}
}

func (m Writer) Write(writeRequest model.PlcWriteRequest) <-chan model.PlcWriteRequestResult {
	result := make(chan model.PlcWriteRequestResult)
	// If we are requesting only one field, use a
	if len(writeRequest.GetFieldNames()) == 1 {
		fieldName := writeRequest.GetFieldNames()[0]

		// Get the KnxNetIp field instance from the request
		field := writeRequest.GetField(fieldName)
		knxNetIpField, err := CastToFieldFromPlcField(field)
		if err != nil {
			result <- &plc4goModel.DefaultPlcWriteRequestResult{
				Request:  writeRequest,
				Response: nil,
				Err:      errors.New("invalid field item type"),
			}
			return result
		}

		// Get the value from the request and serialize it to a byte array
		value := writeRequest.GetValue(fieldName)
		io := utils.NewWriteBufferByteBased()
		fieldType := readWriteModel.KnxDatapointTypeByName(knxNetIpField.GetTypeName())
		if err := readWriteModel.KnxDatapointSerialize(io, value, fieldType); err != nil {
			result <- &plc4goModel.DefaultPlcWriteRequestResult{
				Request:  writeRequest,
				Response: nil,
				Err:      errors.New("error serializing value: " + err.Error()),
			}
			return result
		}
	}
	return result
}
