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

package ads

import (
	"fmt"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

type ValueHandler struct {
	spiValues.DefaultValueHandler

	driverContext *DriverContext
	tagHandler    TagHandler
}

func NewValueHandler() ValueHandler {
	return ValueHandler{}
}

func NewValueHandlerWithDriverContext(driverContext *DriverContext, tagHandler TagHandler) ValueHandler {
	return ValueHandler{
		driverContext: driverContext,
		tagHandler:    tagHandler,
	}
}

func (t ValueHandler) NewPlcValue(tag model.PlcTag, value interface{}) (apiValues.PlcValue, error) {
	return t.parseType(tag, tag.GetArrayInfo(), value)
}

func (t ValueHandler) parseType(tag model.PlcTag, arrayInfo []model.ArrayInfo, value interface{}) (apiValues.PlcValue, error) {
	// Resolve the symbolic tag to a direct tag, that has all the important information.
	var directTag DirectPlcTag
	switch tag.(type) {
	case SymbolicPlcTag:
		symbolicTag := tag.(SymbolicPlcTag)
		directTagPointer, err := t.driverContext.getDirectTagForSymbolTag(symbolicTag)
		if err != nil {
			return nil, fmt.Errorf("couldn't resolve address %s to a valid tag on the PLC", symbolicTag.SymbolicAddress)
		}
		directTag = *directTagPointer
	case DirectPlcTag:
		directTag = tag.(DirectPlcTag)
	}

	// Do the normal resolution.
	valueType := directTag.GetValueType()
	if (arrayInfo != nil) && (len(arrayInfo) > 0) {
		return t.ParseListType(directTag, arrayInfo, value)
	} else if valueType == apiValues.Struct {
		return t.ParseStructType(directTag, value)
	}
	return t.ParseSimpleType(directTag, value)
}

func (t ValueHandler) ParseStructType(_ model.PlcTag, _ interface{}) (apiValues.PlcValue, error) {
	return nil, nil
}
