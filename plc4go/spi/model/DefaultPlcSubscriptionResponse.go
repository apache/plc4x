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
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/pkg/errors"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcSubscriptionResponse
type DefaultPlcSubscriptionResponse struct {
	request apiModel.PlcSubscriptionRequest
	values  map[string]*DefaultPlcSubscriptionResponseItem
}

func NewDefaultPlcSubscriptionResponse(
	request apiModel.PlcSubscriptionRequest,
	responseCodes map[string]apiModel.PlcResponseCode,
	values map[string]apiModel.PlcSubscriptionHandle,
	_options ...options.WithOption,
) apiModel.PlcSubscriptionResponse {
	valueMap := map[string]*DefaultPlcSubscriptionResponseItem{}
	for name, code := range responseCodes {
		value := values[name]
		valueMap[name] = NewDefaultPlcSubscriptionResponseItem(code, value)
	}
	plcSubscriptionResponse := DefaultPlcSubscriptionResponse{
		request: request,
		values:  valueMap,
	}
	localLog := options.ExtractCustomLogger(_options...)
	for subscriptionTagName, consumers := range request.(*DefaultPlcSubscriptionRequest).preRegisteredConsumers {
		subscriptionHandle, err := plcSubscriptionResponse.GetSubscriptionHandle(subscriptionTagName)
		if subscriptionHandle == nil || err != nil {
			localLog.Error().Msgf("PlcSubscriptionHandle for %s not found", subscriptionTagName)
			continue
		}
		for _, consumer := range consumers {
			subscriptionHandle.Register(consumer)
		}
	}
	return &plcSubscriptionResponse
}

func (d *DefaultPlcSubscriptionResponse) IsAPlcMessage() bool {
	return true
}

func (d *DefaultPlcSubscriptionResponse) GetRequest() apiModel.PlcSubscriptionRequest {
	return d.request
}

func (d *DefaultPlcSubscriptionResponse) GetTagNames() []string {
	if d.request == nil {
		// Safety guard
		return nil
	}
	var tagNames []string
	// We take the tag names from the request to keep order as map is not ordered
	for _, name := range d.request.GetTagNames() {
		if _, ok := d.values[name]; ok {
			tagNames = append(tagNames, name)
		}
	}
	return tagNames
}

func (d *DefaultPlcSubscriptionResponse) GetResponseCode(name string) apiModel.PlcResponseCode {
	item, ok := d.values[name]
	if !ok {
		return apiModel.PlcResponseCode_NOT_FOUND
	}
	return item.GetCode()
}

func (d *DefaultPlcSubscriptionResponse) GetSubscriptionHandle(name string) (apiModel.PlcSubscriptionHandle, error) {
	item, ok := d.values[name]
	if !ok {
		return nil, errors.Errorf("item for %s not found", name)
	}
	if item.GetCode() != apiModel.PlcResponseCode_OK {
		return nil, errors.Errorf("%s failed to subscribe", name)
	}
	return item.GetSubscriptionHandle(), nil
}

func (d *DefaultPlcSubscriptionResponse) GetSubscriptionHandles() []apiModel.PlcSubscriptionHandle {
	result := make([]apiModel.PlcSubscriptionHandle, 0, len(d.values))
	for _, value := range d.values {
		result = append(result, value.GetSubscriptionHandle())
	}
	return result
}
