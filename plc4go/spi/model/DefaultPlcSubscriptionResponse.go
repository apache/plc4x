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
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/pkg/errors"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcSubscriptionResponse
type DefaultPlcSubscriptionResponse struct {
	DefaultResponse
	request model.PlcSubscriptionRequest
	values  map[string]model.PlcSubscriptionHandle
}

func NewDefaultPlcSubscriptionResponse(request model.PlcSubscriptionRequest, responseCodes map[string]model.PlcResponseCode, values map[string]model.PlcSubscriptionHandle) *DefaultPlcSubscriptionResponse {
	plcSubscriptionResponse := DefaultPlcSubscriptionResponse{
		DefaultResponse: NewDefaultResponse(responseCodes),
		request:         request,
		values:          values,
	}
	for subscriptionFieldName, consumers := range request.(*DefaultPlcSubscriptionRequest).preRegisteredConsumers {
		subscriptionHandle, err := plcSubscriptionResponse.GetSubscriptionHandle(subscriptionFieldName)
		if subscriptionHandle == nil || err != nil {
			panic("PlcSubscriptionHandle for " + subscriptionFieldName + " not found")
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

func (d *DefaultPlcSubscriptionResponse) GetFieldNames() []string {
	var fieldNames []string
	// We take the field names from the request to keep order as map is not ordered
	for _, name := range d.request.(*DefaultPlcSubscriptionRequest).GetFieldNames() {
		if _, ok := d.responseCodes[name]; ok {
			fieldNames = append(fieldNames, name)
		}
	}
	return fieldNames
}

func (d *DefaultPlcSubscriptionResponse) GetSubscriptionHandle(name string) (model.PlcSubscriptionHandle, error) {
	if d.responseCodes[name] != model.PlcResponseCode_OK {
		return nil, errors.Errorf("%s failed to subscribe", name)
	}
	return d.values[name], nil
}

func (d *DefaultPlcSubscriptionResponse) GetSubscriptionHandles() []model.PlcSubscriptionHandle {
	result := make([]model.PlcSubscriptionHandle, 0, len(d.values))

	for _, value := range d.values {
		result = append(result, value)
	}
	return result
}
