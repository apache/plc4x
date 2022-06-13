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

import "github.com/apache/plc4x/plc4go/pkg/plc4go/model"

type DefaultRequest struct {
	fields     map[string]model.PlcField
	fieldNames []string
}

func NewDefaultRequest(Fields map[string]model.PlcField, FieldNames []string) DefaultRequest {
	return DefaultRequest{Fields, FieldNames}
}

func (m DefaultRequest) GetFieldNames() []string {
	return m.fieldNames
}

func (m DefaultRequest) GetField(name string) model.PlcField {
	if field, ok := m.fields[name]; ok {
		return field
	}
	return nil
}
