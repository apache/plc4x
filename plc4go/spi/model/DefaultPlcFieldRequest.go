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
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcFieldRequest
type DefaultPlcFieldRequest struct {
	fields     map[string]model.PlcField
	fieldNames []string
}

func (d *DefaultPlcFieldRequest) IsAPlcMessage() bool {
	return true
}

func NewDefaultPlcFieldRequest(fields map[string]model.PlcField, fieldNames []string) DefaultPlcFieldRequest {
	return DefaultPlcFieldRequest{fields: fields, fieldNames: fieldNames}
}

func (d *DefaultPlcFieldRequest) GetFieldNames() []string {
	return d.fieldNames
}

func (d *DefaultPlcFieldRequest) GetField(name string) model.PlcField {
	if field, ok := d.fields[name]; ok {
		return field
	}
	return nil
}
