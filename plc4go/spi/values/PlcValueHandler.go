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

package values

import (
	"errors"
	"fmt"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/rs/zerolog"
	"reflect"
	"strconv"
	"strings"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
)

type DefaultValueHandler struct {
	log zerolog.Logger
}

func NewDefaultValueHandler(_options ...options.WithOption) DefaultValueHandler {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return DefaultValueHandler{
		log: customLogger,
	}
}

func (m DefaultValueHandler) NewPlcValue(tag apiModel.PlcTag, value any) (apiValues.PlcValue, error) {
	return m.parseType(tag, tag.GetArrayInfo(), value)
}

func (m DefaultValueHandler) parseType(tag apiModel.PlcTag, arrayInfo []apiModel.ArrayInfo, value any) (apiValues.PlcValue, error) {
	valueType := tag.GetValueType()
	if (arrayInfo != nil) && (len(arrayInfo) > 0) {
		return m.ParseListType(tag, arrayInfo, value)
	} else if valueType == apiValues.Struct {
		return m.ParseStructType(tag, value)
	}
	return m.ParseSimpleType(tag, value)
}

func (m DefaultValueHandler) ParseListType(tag apiModel.PlcTag, arrayInfo []apiModel.ArrayInfo, value any) (apiValues.PlcValue, error) {
	// We've reached the end of the recursion.
	if len(arrayInfo) == 0 {
		return m.parseType(tag, arrayInfo, value)
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

	// Check that the current slice has enough apiValues.
	if len(curValues) != int(curArrayInfo.GetSize()) {
		return nil, errors.New("number of actual apiValues " + strconv.Itoa(len(curValues)) +
			" doesn't match tag size " + strconv.Itoa(int(curArrayInfo.GetSize())))
	}

	// Actually convert the current array info level.
	var plcValues []apiValues.PlcValue
	for i := uint32(0); i < curArrayInfo.GetSize(); i++ {
		curValue := curValues[i]
		plcValue, err := m.ParseListType(tag, restArrayInfo, curValue)
		if err != nil {
			return nil, errors.New("error parsing PlcValue: " + err.Error())
		}
		plcValues = append(plcValues, plcValue)
	}
	return NewPlcList(plcValues), nil
}

func (m DefaultValueHandler) ParseStructType(_ apiModel.PlcTag, _ any) (apiValues.PlcValue, error) {
	return nil, errors.New("structured types not supported by the base value handler")
}

func (m DefaultValueHandler) ParseSimpleType(tag apiModel.PlcTag, value any) (apiValues.PlcValue, error) {
	plcValue, err := m.NewPlcValueFromType(tag.GetValueType(), value)
	if err != nil && strings.HasPrefix(err.Error(), "couldn't cast") {
		stringValue := fmt.Sprintf("%v", value)
		plcValue, err = m.NewPlcValueFromType(tag.GetValueType(), stringValue)
		if err == nil {
			m.log.Debug().
				Interface("value", value).
				Stringer("plcValue", plcValue).
				Msg("had to convert %v into %v by using string conversion")
		}
	}
	return plcValue, err
}

func (m DefaultValueHandler) NewPlcValueFromType(valueType apiValues.PlcValueType, value any) (apiValues.PlcValue, error) {
	// If the user passed in PLCValues, take a shortcut.
	plcValue, isPlcValue := value.(apiValues.PlcValue)
	if isPlcValue {
		if plcValue.GetPlcValueType() != valueType {
			// TODO: Check if the used PlcValueType can be casted to the target type.
		} else if plcValue.GetPlcValueType() == apiValues.List {
			// TODO: Check all items
		} else if plcValue.GetPlcValueType() == apiValues.Struct {
			// TODO: Check all children
		}
		return plcValue, nil
	}

	stringValue, isString := value.(string)
	switch valueType {
	// Bit & Bit-Strings
	case apiValues.BOOL:
		if isString {
			casted, err := strconv.ParseBool(stringValue)
			if err != nil {
				return nil, errors.New("couldn't parse string value '" + stringValue + "' to bool")
			}
			return NewPlcBOOL(casted), nil
		} else {
			casted, ok := value.(bool)
			if !ok {
				return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to bool")
			}
			return NewPlcBOOL(casted), nil
		}
	case apiValues.BYTE:
		if isString {
			casted, err := strconv.ParseUint(stringValue, 10, 8)
			if err != nil {
				return nil, errors.New("couldn't parse string value '" + stringValue + "' to byte")
			}
			return NewPlcBYTE(uint8(casted)), nil
		} else {
			casted, ok := value.(uint8)
			if !ok {
				return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to uint8")
			}
			return NewPlcBYTE(casted), nil
		}
	case apiValues.WORD:
		if isString {
			casted, err := strconv.ParseUint(stringValue, 10, 16)
			if err != nil {
				return nil, errors.New("couldn't parse string value '" + stringValue + "' to word")
			}
			return NewPlcWORD(uint16(casted)), nil
		} else {
			casted, ok := value.(uint16)
			if !ok {
				return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to uint16")
			}
			return NewPlcWORD(casted), nil
		}
	case apiValues.DWORD:
		if isString {
			casted, err := strconv.ParseUint(stringValue, 10, 32)
			if err != nil {
				return nil, errors.New("couldn't parse string value '" + stringValue + "' to dword")
			}
			return NewPlcDWORD(uint32(casted)), nil
		} else {
			casted, ok := value.(uint32)
			if !ok {
				return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to uint32")
			}
			return NewPlcDWORD(casted), nil
		}
	case apiValues.LWORD:
		if isString {
			casted, err := strconv.ParseUint(stringValue, 10, 64)
			if err != nil {
				return nil, errors.New("couldn't parse string value '" + stringValue + "' to lword")
			}
			return NewPlcLWORD(casted), nil
		} else {
			casted, ok := value.(uint64)
			if !ok {
				return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to uint64")
			}
			return NewPlcLWORD(casted), nil
		}

	// Integers
	case apiValues.USINT:
		if isString {
			casted, err := strconv.ParseUint(stringValue, 10, 8)
			if err != nil {
				return nil, errors.New("couldn't parse string value '" + stringValue + "' to usint")
			}
			return NewPlcUSINT(uint8(casted)), nil
		} else {
			casted, ok := value.(uint8)
			if !ok {
				return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to uint8")
			}
			return NewPlcUSINT(casted), nil
		}
	case apiValues.UINT:
		if isString {
			casted, err := strconv.ParseUint(stringValue, 10, 16)
			if err != nil {
				return nil, errors.New("couldn't parse string value '" + stringValue + "' to uint")
			}
			return NewPlcUINT(uint16(casted)), nil
		} else {
			casted, ok := value.(uint16)
			if !ok {
				return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to uint16")
			}
			return NewPlcUINT(casted), nil
		}
	case apiValues.UDINT:
		if isString {
			casted, err := strconv.ParseUint(stringValue, 10, 32)
			if err != nil {
				return nil, errors.New("couldn't parse string value '" + stringValue + "' to udint")
			}
			return NewPlcUDINT(uint32(casted)), nil
		} else {
			casted, ok := value.(uint32)
			if !ok {
				return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to uint32")
			}
			return NewPlcUDINT(casted), nil
		}
	case apiValues.ULINT:
		if isString {
			casted, err := strconv.ParseUint(stringValue, 10, 64)
			if err != nil {
				return nil, errors.New("couldn't parse string value '" + stringValue + "' to ulint")
			}
			return NewPlcULINT(casted), nil
		} else {
			casted, ok := value.(uint64)
			if !ok {
				return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to uint64")
			}
			return NewPlcULINT(casted), nil
		}
	case apiValues.SINT:
		if isString {
			casted, err := strconv.ParseInt(stringValue, 10, 8)
			if err != nil {
				return nil, errors.New("couldn't parse string value '" + stringValue + "' to sint")
			}
			return NewPlcSINT(int8(casted)), nil
		} else {
			casted, ok := value.(int8)
			if !ok {
				return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to int8")
			}
			return NewPlcSINT(casted), nil
		}
	case apiValues.INT:
		if isString {
			casted, err := strconv.ParseInt(stringValue, 10, 16)
			if err != nil {
				return nil, errors.New("couldn't parse string value '" + stringValue + "' to int")
			}
			return NewPlcINT(int16(casted)), nil
		} else {
			casted, ok := value.(int16)
			if !ok {
				return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to int16")
			}
			return NewPlcINT(casted), nil
		}
	case apiValues.DINT:
		if isString {
			casted, err := strconv.ParseInt(stringValue, 10, 32)
			if err != nil {
				return nil, errors.New("couldn't parse string value '" + stringValue + "' to dint")
			}
			return NewPlcDINT(int32(casted)), nil
		} else {
			casted, ok := value.(int32)
			if !ok {
				return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to int32")
			}
			return NewPlcDINT(casted), nil
		}
	case apiValues.LINT:
		if isString {
			casted, err := strconv.ParseInt(stringValue, 10, 64)
			if err != nil {
				return nil, errors.New("couldn't parse string value '" + stringValue + "' to lint")
			}
			return NewPlcLINT(casted), nil
		} else {
			casted, ok := value.(int64)
			if !ok {
				return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to int64")
			}
			return NewPlcLINT(casted), nil
		}

	// Floating Point Values
	case apiValues.REAL:
		if isString {
			casted, err := strconv.ParseFloat(stringValue, 32)
			if err != nil {
				return nil, errors.New("couldn't parse string value '" + stringValue + "' to real")
			}
			return NewPlcREAL(float32(casted)), nil
		} else {
			casted, ok := value.(float32)
			if !ok {
				return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to float32")
			}
			return NewPlcREAL(casted), nil
		}
	case apiValues.LREAL:
		if isString {
			casted, err := strconv.ParseFloat(stringValue, 64)
			if err != nil {
				return nil, errors.New("couldn't parse string value '" + stringValue + "' to lreal")
			}
			return NewPlcLREAL(casted), nil
		} else {
			casted, ok := value.(float64)
			if !ok {
				return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to float64")
			}
			return NewPlcLREAL(casted), nil
		}

	// Temporal Values
	// - Duration
	case apiValues.TIME:
		if isString {
			return nil, errors.New("string to IEC61131_TIME conversion not implemented")
		} else {
			casted, ok := value.(time.Duration)
			if !ok {
				return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to time.Duration")
			}
			return NewPlcTIME(casted), nil
		}
	// - Date
	case apiValues.DATE:
		if isString {
			return nil, errors.New("string to IEC61131_DATE conversion not implemented")
		} else {
			casted, ok := value.(time.Time)
			if !ok {
				return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to time.Time")
			}
			return NewPlcDATE(casted), nil
		}
	// - Time
	case apiValues.TIME_OF_DAY:
		if isString {
			return nil, errors.New("string to IEC61131_TIME_OF_DAY conversion not implemented")
		} else {
			casted, ok := value.(time.Time)
			if !ok {
				return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to time.Time")
			}
			return NewPlcTIME_OF_DAY(casted), nil
		}
	// - Date and Time
	case apiValues.DATE_AND_TIME:
		if isString {
			return nil, errors.New("string to IEC61131_DATE_AND_TIME conversion not implemented")
		} else {
			casted, ok := value.(time.Time)
			if !ok {
				return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to time.Time")
			}
			return NewPlcDATE_AND_TIME(casted), nil
		}

	// Chars and Strings
	case apiValues.CHAR:
		if !isString {
			return nil, errors.New("non-string to IEC61131_CHAR conversion not implemented")
		} else if len(stringValue) > 1 {
			return nil, errors.New("IEC61131_CHAR can only contain one character")
		} else {
			return NewPlcCHAR(stringValue), nil
		}
	case apiValues.WCHAR:
		if !isString {
			return nil, errors.New("non-string to IEC61131_WCHAR conversion not implemented")
		} else if len(stringValue) > 1 {
			return nil, errors.New("IEC61131_WCHAR can only contain one character")
		} else {
			return NewPlcWCHAR(stringValue), nil
		}
	case apiValues.STRING:
		if !isString {
			return nil, errors.New("non-string to IEC61131_STRING conversion not implemented")
		} else {
			return NewPlcSTRING(stringValue), nil
		}
	case apiValues.WSTRING:
		if !isString {
			return nil, errors.New("non-string to IEC61131_WSTRING conversion not implemented")
		} else {
			return NewPlcWSTRING(stringValue), nil
		}
	}

	return nil, errors.New("Unsupported type " + valueType.String())
}
