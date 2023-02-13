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

package model

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type DefaultPlcWriteResponse struct {
	DefaultResponse
	request       model.PlcWriteRequest
	responseCodes map[string]model.PlcResponseCode
}

func NewDefaultPlcWriteResponse(request model.PlcWriteRequest, responseCodes map[string]model.PlcResponseCode) model.PlcWriteResponse {
	return &DefaultPlcWriteResponse{
		request:       request,
		responseCodes: responseCodes,
	}
}

func (d *DefaultPlcWriteResponse) GetTagNames() []string {
	var tagNames []string
	// We take the tag names from the request to keep order as map is not ordered
	for _, name := range d.request.GetTagNames() {
		if _, ok := d.responseCodes[name]; ok {
			tagNames = append(tagNames, name)
		}
	}
	return tagNames
}

func (d *DefaultPlcWriteResponse) GetRequest() model.PlcWriteRequest {
	return d.request
}

func (d *DefaultPlcWriteResponse) GetResponseCode(name string) model.PlcResponseCode {
	return d.responseCodes[name]
}

func (d *DefaultPlcWriteResponse) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := d.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (d *DefaultPlcWriteResponse) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("PlcWriteResponse"); err != nil {
		return err
	}
	if err := d.DefaultResponse.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
		return err
	}

	if d.request != nil {
		if serializableField, ok := d.request.(utils.Serializable); ok {
			if err := serializableField.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
				return err
			}
		} else {
			stringValue := fmt.Sprintf("%v", d.request)
			if err := writeBuffer.WriteString("request", uint32(len(stringValue)*8), "UTF-8", stringValue); err != nil {
				return err
			}
		}
	}
	if err := writeBuffer.PushContext("responseCodes"); err != nil {
		return err
	}
	for name, elem := range d.responseCodes {
		var elem interface{} = elem
		elemAsString := fmt.Sprintf("%v", elem)
		if err := writeBuffer.WriteString(name, uint32(len(elemAsString)*8), "UTF-8", elemAsString); err != nil {
			return err
		}
	}
	if err := writeBuffer.PopContext("responseCodes"); err != nil {
		return err
	}
	if err := writeBuffer.PopContext("PlcWriteResponse"); err != nil {
		return err
	}
	return nil
}

func (d *DefaultPlcWriteResponse) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), d); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
