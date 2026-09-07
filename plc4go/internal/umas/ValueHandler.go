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

package umas

import (
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

// ValueHandler turns the values a user hands to a write request into plc4x values.
//
// A UMAS tag address is a symbol name and says nothing about the symbol's type - only the data
// dictionary does, and the writer looks the type up there when it serializes. So a tag parsed from an
// address has no value type, and the default handler can't be asked to coerce to it: it would refuse
// every value with "Unsupported type NULL". Instead the value keeps whatever type it naturally has
// and the writer converts it to the symbol's UMAS type, with a range check, in encodeWriteValue.
//
// plc4j ends up in the same place from the other side: its SymbolicUmasTag.getPlcValueType returns
// null rather than PlcValueType.NULL precisely so that the DefaultPlcValueHandler leaves the value
// alone.
type ValueHandler struct {
	spiValues.DefaultValueHandler
}

func NewValueHandler(_options ...options.WithOption) ValueHandler {
	return ValueHandler{
		spiValues.NewDefaultValueHandler(_options...),
	}
}

// NewPlcValue keeps whatever the user handed over when the tag's type is unknown, and defers to the
// default handler once it is known (which is the case for a tag that came out of a browse).
func (m ValueHandler) NewPlcValue(tag apiModel.PlcTag, value any) (apiValues.PlcValue, error) {
	if tag == nil || tag.GetValueType() != apiValues.NULL {
		return m.DefaultValueHandler.NewPlcValue(tag, value)
	}
	return plcValueOfNativeType(value)
}

// plcValueOfNativeType wraps a Go value in the plc4x value of the same shape. Nothing here decides
// what goes on the wire - that is the symbol's UMAS type, applied by encodeWriteValue - this only has
// to preserve the value and its signedness so the range check there is meaningful.
func plcValueOfNativeType(value any) (apiValues.PlcValue, error) {
	switch typedValue := value.(type) {
	case nil:
		return nil, errors.New("there is no value to write")
	case apiValues.PlcValue:
		return typedValue, nil
	case bool:
		return spiValues.NewPlcBOOL(typedValue), nil
	case uint8:
		return spiValues.NewPlcUSINT(typedValue), nil
	case uint16:
		return spiValues.NewPlcUINT(typedValue), nil
	case uint32:
		return spiValues.NewPlcUDINT(typedValue), nil
	case uint64:
		return spiValues.NewPlcULINT(typedValue), nil
	case uint:
		return spiValues.NewPlcULINT(uint64(typedValue)), nil
	case int8:
		return spiValues.NewPlcSINT(typedValue), nil
	case int16:
		return spiValues.NewPlcINT(typedValue), nil
	case int32:
		return spiValues.NewPlcDINT(typedValue), nil
	case int64:
		return spiValues.NewPlcLINT(typedValue), nil
	case int:
		return spiValues.NewPlcLINT(int64(typedValue)), nil
	case float32:
		return spiValues.NewPlcREAL(typedValue), nil
	case float64:
		return spiValues.NewPlcLREAL(typedValue), nil
	case string:
		return spiValues.NewPlcSTRING(typedValue), nil
	case time.Duration:
		return spiValues.NewPlcTIME(typedValue), nil
	case time.Time:
		// A DATE, a TIME_OF_DAY and a DATE_AND_TIME all carry a time.Time and the tag doesn't say
		// which one is meant; the widest of the three keeps every field, and encodeWriteValue picks
		// the ones the symbol's type actually needs.
		return spiValues.NewPlcDATE_AND_TIME(typedValue), nil
	case []byte:
		return spiValues.NewPlcRawByteArray(typedValue), nil
	default:
		return nil, errors.Errorf("Can't turn a %T into a value to write, hand over a PlcValue instead", value)
	}
}
