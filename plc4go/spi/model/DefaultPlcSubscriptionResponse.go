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
	"github.com/pkg/errors"
)

type DefaultPlcSubscriptionResponse struct {
	DefaultResponse
	request model.PlcSubscriptionRequest
	values  map[string]*DefaultPlcSubscriptionResponseItem
}

func NewDefaultPlcSubscriptionResponse(request model.PlcSubscriptionRequest, responseCodes map[string]model.PlcResponseCode, values map[string]model.PlcSubscriptionHandle) *DefaultPlcSubscriptionResponse {
	valueMap := map[string]*DefaultPlcSubscriptionResponseItem{}
	for name, code := range responseCodes {
		value := values[name]
		valueMap[name] = NewSubscriptionResponseItem(code, value)
	}
	plcSubscriptionResponse := DefaultPlcSubscriptionResponse{
		request: request,
		values:  valueMap,
	}
	for subscriptionTagName, consumers := range request.(*DefaultPlcSubscriptionRequest).preRegisteredConsumers {
		subscriptionHandle, err := plcSubscriptionResponse.GetSubscriptionHandle(subscriptionTagName)
		if subscriptionHandle == nil || err != nil {
			panic("PlcSubscriptionHandle for " + subscriptionTagName + " not found")
		}
		for _, consumer := range consumers {
			subscriptionHandle.Register(consumer)
		}
	}
	return &plcSubscriptionResponse
}

func (d *DefaultPlcSubscriptionResponse) GetRequest() model.PlcSubscriptionRequest {
	return d.request
}

func (d *DefaultPlcSubscriptionResponse) GetTagNames() []string {
	var tagNames []string
	// We take the tag names from the request to keep order as map is not ordered
	for _, name := range d.request.(*DefaultPlcSubscriptionRequest).GetTagNames() {
		if _, ok := d.values[name]; ok {
			tagNames = append(tagNames, name)
		}
	}
	return tagNames
}

func (d *DefaultPlcSubscriptionResponse) GetResponseCode(name string) model.PlcResponseCode {
	return d.values[name].GetCode()
}

func (d *DefaultPlcSubscriptionResponse) GetSubscriptionHandle(name string) (model.PlcSubscriptionHandle, error) {
	if d.values[name].GetCode() != model.PlcResponseCode_OK {
		return nil, errors.Errorf("%s failed to subscribe", name)
	}
	return d.values[name].GetSubscriptionHandle(), nil
}

func (d *DefaultPlcSubscriptionResponse) GetSubscriptionHandles() []model.PlcSubscriptionHandle {
	result := make([]model.PlcSubscriptionHandle, 0, len(d.values))
	for _, value := range d.values {
		result = append(result, value.GetSubscriptionHandle())
	}
	return result
}

func (d *DefaultPlcSubscriptionResponse) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := d.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (d *DefaultPlcSubscriptionResponse) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("PlcSubscriptionResponse"); err != nil {
		return err
	}

	if d.request != nil {
		if serializableField, ok := d.request.(utils.Serializable); ok {
			if err := writeBuffer.PushContext("request"); err != nil {
				return err
			}
			if err := serializableField.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
				return err
			}
			if err := writeBuffer.PopContext("request"); err != nil {
				return err
			}
		} else {
			stringValue := fmt.Sprintf("%v", d.request)
			if err := writeBuffer.WriteString("request", uint32(len(stringValue)*8), "UTF-8", stringValue); err != nil {
				return err
			}
		}
	}
	if err := writeBuffer.PushContext("values", utils.WithRenderAsList(true)); err != nil {
		return err
	}
	for name, elem := range d.values {
		_value := fmt.Sprintf("%v", elem)

		if err := writeBuffer.WriteString(name, uint32(len(_value)*8), "UTF-8", _value); err != nil {
			return err
		}
	}
	if err := writeBuffer.PopContext("values", utils.WithRenderAsList(true)); err != nil {
		return err
	}
	if err := writeBuffer.PopContext("PlcSubscriptionResponse"); err != nil {
		return err
	}
	return nil
}

func (d *DefaultPlcSubscriptionResponse) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), d); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
