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
	"reflect"
	"strconv"

	model2 "github.com/apache/plc4x/plc4go/internal/ads/model"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	driverModel "github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
	"github.com/pkg/errors"
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

func (t ValueHandler) NewPlcValue(tag model.PlcTag, value any) (apiValues.PlcValue, error) {
	return t.parseType(tag, value)
}

func (t ValueHandler) parseType(tag model.PlcTag, value any) (apiValues.PlcValue, error) {
	// Resolve the symbolic tag to a direct tag, that has all the important information.
	var directTag model2.DirectPlcTag
	switch tag.(type) {
	case model2.SymbolicPlcTag:
		symbolicTag := tag.(model2.SymbolicPlcTag)
		directTagPointer, err := t.driverContext.getDirectTagForSymbolTag(symbolicTag)
		if err != nil {
			return nil, fmt.Errorf("couldn't resolve address %s to a valid tag on the PLC", symbolicTag.SymbolicAddress)
		}
		directTag = *directTagPointer
	case model2.DirectPlcTag:
		directTag = tag.(model2.DirectPlcTag)
	}

	return t.AdsParseType(directTag.DataType, directTag.DataType.GetArrayInfo(), value)
}

func (t ValueHandler) AdsParseType(datatype driverModel.AdsDataTypeTableEntry, arrayInfo []driverModel.AdsDataTypeArrayInfo, value any) (apiValues.PlcValue, error) {
	// Do the normal resolution.
	if (arrayInfo != nil) && (len(arrayInfo) > 0) {
		return t.AdsParseListType(datatype, arrayInfo, value)
	} else if datatype.GetNumChildren() > 0 {
		return t.AdsParseStructType(datatype, value)
	}
	return t.AdsParseSimpleType(datatype, value)
}

func (t ValueHandler) AdsParseListType(dataType driverModel.AdsDataTypeTableEntry, arrayInfo []driverModel.AdsDataTypeArrayInfo, value any) (apiValues.PlcValue, error) {
	// We've reached the end of the recursion.
	if len(arrayInfo) == 0 {
		return t.AdsParseType(dataType, arrayInfo, value)
	}

	s := reflect.ValueOf(value)
	if s.Kind() != reflect.Slice {
		return nil, errors.New("couldn't cast value to []any")
	}
	curValues := make([]any, s.Len())
	for i := 0; i < s.Len(); i++ {
		curValues[i] = s.Index(i).Interface()
	}

	curArrayInfo := arrayInfo[0]
	restArrayInfo := arrayInfo[1:]

	// Check that the current slice has enough values.
	if len(curValues) != int(curArrayInfo.GetNumElements()) {
		return nil, errors.New("number of actual values " + strconv.Itoa(len(curValues)) +
			" doesn't match tag size " + strconv.Itoa(int(curArrayInfo.GetNumElements())))
	}

	// Actually convert the current array info level.
	var plcValues []apiValues.PlcValue
	for i := uint32(0); i < curArrayInfo.GetNumElements(); i++ {
		curValue := curValues[i]
		plcValue, err := t.AdsParseListType(dataType, restArrayInfo, curValue)
		if err != nil {
			return nil, errors.New("error parsing PlcValue: " + err.Error())
		}
		plcValues = append(plcValues, plcValue)
	}

	return spiValues.NewPlcList(plcValues), nil
}

func (t ValueHandler) AdsParseStructType(dataType driverModel.AdsDataTypeTableEntry, value any) (apiValues.PlcValue, error) {
	// Unfortunately it seems impossible to cast map[string]apiValues.PlcValue to map[string]any
	if plcStruct, ok := value.(spiValues.PlcStruct); ok {
		parsedValues := map[string]apiValues.PlcValue{}
		childValues := plcStruct.GetStruct()

		for _, childTypeEntry := range dataType.GetChildren() {
			childName := childTypeEntry.GetPropertyName()
			childType := t.driverContext.dataTypeTable[childTypeEntry.GetDataTypeName()]
			childArrayInfo := childType.GetArrayInfo()
			childValue, ok := childValues[childTypeEntry.GetPropertyName()]
			if !ok {
				return nil, fmt.Errorf("missing child value named %s", childName)
			}
			parsedChildValue, err := t.AdsParseType(childType, childArrayInfo, childValue)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("error parsing child %s", childName))
			}
			parsedValues[childName] = parsedChildValue
		}

		return spiValues.NewPlcStruct(parsedValues), nil
	} else if simpleMap, ok := value.(map[string]any); ok {
		parsedValues := map[string]apiValues.PlcValue{}

		for _, childTypeEntry := range dataType.GetChildren() {
			childName := childTypeEntry.GetPropertyName()
			childType := t.driverContext.dataTypeTable[childTypeEntry.GetDataTypeName()]
			childArrayInfo := childType.GetArrayInfo()
			childValue, ok := simpleMap[childTypeEntry.GetPropertyName()]
			if !ok {
				return nil, fmt.Errorf("missing child value named %s", childName)
			}
			parsedChildValue, err := t.AdsParseType(childType, childArrayInfo, childValue)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("error parsing child %s", childName))
			}
			parsedValues[childName] = parsedChildValue
		}

		return spiValues.NewPlcStruct(parsedValues), nil
	}

	return nil, nil
}

func (t ValueHandler) AdsParseSimpleType(dataType driverModel.AdsDataTypeTableEntry, value any) (apiValues.PlcValue, error) {
	// Get the PlcValue type for this ads-datatype.
	plcValueType := t.driverContext.getDataTypeForDataTypeTableEntry(dataType)
	return t.NewPlcValueFromType(plcValueType, value)
}
