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
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/rs/zerolog/log"
)

const (
	////
	// Bit & Bit-Strings

	IEC61131_BOOL  string = "BOOL"
	IEC61131_BYTE  string = "BYTE"
	IEC61131_WORD  string = "WORD"
	IEC61131_DWORD string = "DWORD"
	IEC61131_LWORD string = "LWORD"
	//
	////

	////
	// Integers

	IEC61131_USINT string = "USINT"
	IEC61131_UINT  string = "UINT"
	IEC61131_UDINT string = "UDINT"
	IEC61131_ULINT string = "ULINT"
	IEC61131_SINT  string = "SINT"
	IEC61131_INT   string = "INT"
	IEC61131_DINT  string = "DINT"
	IEC61131_LINT  string = "LINT"
	//
	////

	////
	// Floating Point Values

	IEC61131_REAL  string = "REAL"
	IEC61131_LREAL string = "LREAL"
	//
	////

	////
	// Temporal Values

	IEC61131_TIME          string = "TIME"
	IEC61131_DATE          string = "DATE"
	IEC61131_TIME_OF_DAY   string = "TIME_OF_DAY"
	IEC61131_DATE_AND_TIME string = "DATE_AND_TIME"
	//
	////

	////
	// Chars and Strings

	IEC61131_CHAR    string = "CHAR"
	IEC61131_WCHAR   string = "WCHAR"
	IEC61131_STRING  string = "STRING"
	IEC61131_WSTRING string = "WSTRING"
	//
	////
)

type DefaultValueHandler struct {
}

func (m DefaultValueHandler) NewPlcValue(tag model.PlcTag, value interface{}) (values.PlcValue, error) {
	return m.parseType(tag, tag.GetArrayInfo(), value)
}

func (m DefaultValueHandler) parseType(tag model.PlcTag, arrayInfo []model.ArrayInfo, value interface{}) (values.PlcValue, error) {
	valueType := tag.GetValueType()
	if (arrayInfo != nil) && (len(arrayInfo) > 0) {
		return m.ParseListType(tag, arrayInfo, value)
	} else if valueType == values.Struct {
		return m.ParseStructType(tag, value)
	}
	return m.ParseSimpleType(tag, value)
}

func (m DefaultValueHandler) ParseListType(tag model.PlcTag, arrayInfo []model.ArrayInfo, value interface{}) (values.PlcValue, error) {
	// We've reached the end of the recursion.
	if len(arrayInfo) == 0 {
		return m.parseType(tag, arrayInfo, value)
	}

	s := reflect.ValueOf(value)
	if s.Kind() != reflect.Slice {
		return nil, errors.New("couldn't cast value to []interface{}")
	}
	curValues := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		curValues[i] = s.Index(i).Interface()
	}

	curArrayInfo := arrayInfo[0]
	restArrayInfo := arrayInfo[1:]

	// Check that the current slice has enough values.
	if len(curValues) != int(curArrayInfo.GetSize()) {
		return nil, errors.New("number of actual values " + strconv.Itoa(len(curValues)) +
			" doesn't match tag size " + strconv.Itoa(int(curArrayInfo.GetSize())))
	}

	// Actually convert the current array info level.
	var plcValues []values.PlcValue
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

func (m DefaultValueHandler) ParseStructType(_ model.PlcTag, _ interface{}) (values.PlcValue, error) {
	return nil, errors.New("structured types not supported by the base value handler")
}

func (m DefaultValueHandler) ParseSimpleType(tag model.PlcTag, value interface{}) (values.PlcValue, error) {
	plcValue, err := m.NewPlcValueFromType(tag.GetValueType(), value)
	if err != nil && strings.HasPrefix(err.Error(), "couldn't cast") {
		stringValue := fmt.Sprintf("%v", value)
		plcValue, err = m.NewPlcValueFromType(tag.GetValueType(), stringValue)
		if err == nil {
			log.Debug().Msgf("had to convert %v into %v by using string conversion", value, plcValue)
		}
	}
	return plcValue, err
}

func (m DefaultValueHandler) NewPlcValueFromType(valueType values.PlcValueType, value interface{}) (values.PlcValue, error) {
	// If the user passed in PLCValues, take a shortcut.
	plcValue, isPlcValue := value.(values.PlcValue)
	if isPlcValue {
		if plcValue.GetPlcValueType() != valueType {
			// TODO: Check if the used PlcValueType can be casted to the target type.
		} else if plcValue.GetPlcValueType() == values.List {
			// TODO: Check all items
		} else if plcValue.GetPlcValueType() == values.Struct {
			// TODO: Check all children
		}
		return plcValue, nil
	}

	stringValue, isString := value.(string)
	switch valueType {
	// Bit & Bit-Strings
	case values.BOOL:
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
	case values.BYTE:
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
	case values.WORD:
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
	case values.DWORD:
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
	case values.LWORD:
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
	case values.USINT:
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
	case values.UINT:
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
	case values.UDINT:
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
	case values.ULINT:
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
	case values.SINT:
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
	case values.INT:
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
	case values.DINT:
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
	case values.LINT:
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
	case values.REAL:
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
	case values.LREAL:
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
	case values.TIME:
		if isString {
			return nil, errors.New("string to IEC61131_TIME conversion not implemented")
		} else {
			casted, ok := value.(uint32)
			if !ok {
				return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to uint32")
			}
			return NewPlcTIME(casted), nil
		}
	case values.DATE:
		if isString {
			return nil, errors.New("string to IEC61131_DATE conversion not implemented")
		} else {
			casted, ok := value.(time.Time)
			if !ok {
				return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to time.Time")
			}
			return NewPlcDATE(casted), nil
		}
	case values.TIME_OF_DAY:
		if isString {
			return nil, errors.New("string to IEC61131_TIME_OF_DAY conversion not implemented")
		} else {
			casted, ok := value.(time.Time)
			if !ok {
				return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to time.Time")
			}
			return NewPlcTIME_OF_DAY(casted), nil
		}
	case values.DATE_AND_TIME:
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
	case values.CHAR:
		if !isString {
			return nil, errors.New("non-string to IEC61131_CHAR conversion not implemented")
		} else if len(stringValue) > 1 {
			return nil, errors.New("IEC61131_CHAR can only contain one character")
		} else {
			return NewPlcCHAR(stringValue), nil
		}
	case values.WCHAR:
		if !isString {
			return nil, errors.New("non-string to IEC61131_WCHAR conversion not implemented")
		} else if len(stringValue) > 1 {
			return nil, errors.New("IEC61131_WCHAR can only contain one character")
		} else {
			return NewPlcWCHAR(stringValue), nil
		}
	case values.STRING:
		if !isString {
			return nil, errors.New("non-string to IEC61131_STRING conversion not implemented")
		} else {
			return NewPlcSTRING(stringValue), nil
		}
	case values.WSTRING:
		if !isString {
			return nil, errors.New("non-string to IEC61131_WSTRING conversion not implemented")
		} else {
			return NewPlcSTRING(stringValue), nil
		}
	}

	return nil, errors.New("Unsupported type " + valueType.String())
}
