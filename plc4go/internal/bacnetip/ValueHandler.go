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

package bacnetip

import (
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/values"
)

// ValueHandler maps user-supplied Go values (the third arg to AddTagAddress
// in a write request) onto plc4go PlcValues. Because BACnet properties have
// many possible underlying types and our plcTag reports PlcValueType_Struct,
// the default handler's struct path bails out — we override NewPlcValue
// directly to dispatch on the user's Go type.
//
// We do NOT embed DefaultValueHandler. Embedded-method dispatch in Go is
// static: DefaultValueHandler.parseType would call DefaultValueHandler's own
// ParseStructType, not our override. Implementing the spi.PlcValueHandler
// interface (just NewPlcValue) directly keeps control of the dispatch.
type ValueHandler struct{}

func NewValueHandler() ValueHandler {
	return ValueHandler{}
}

// NewPlcValue wraps a raw Go primitive into the matching PlcValue. The
// Writer's ValueEncoder later inverts this onto a BACnet ApplicationTag, so
// here we only care about preserving the type information losslessly.
func (h ValueHandler) NewPlcValue(_ apiModel.PlcTag, value any) (apiValues.PlcValue, error) {
	if v, ok := value.(apiValues.PlcValue); ok {
		return v, nil
	}
	switch v := value.(type) {
	case nil:
		return values.NewPlcNULL(), nil
	case bool:
		return values.NewPlcBOOL(v), nil
	case float32:
		return values.NewPlcREAL(v), nil
	case float64:
		return values.NewPlcLREAL(v), nil
	case int8:
		return values.NewPlcSINT(v), nil
	case int16:
		return values.NewPlcINT(v), nil
	case int32:
		return values.NewPlcDINT(v), nil
	case int64:
		return values.NewPlcLINT(v), nil
	case int:
		return values.NewPlcLINT(int64(v)), nil
	case uint8:
		return values.NewPlcUSINT(v), nil
	case uint16:
		return values.NewPlcUINT(v), nil
	case uint32:
		return values.NewPlcUDINT(v), nil
	case uint64:
		return values.NewPlcULINT(v), nil
	case uint:
		return values.NewPlcULINT(uint64(v)), nil
	case string:
		return values.NewPlcSTRING(v), nil
	case []byte:
		return values.NewPlcRawByteArray(v), nil
	default:
		return nil, errors.Errorf("BACnet value handler can't encode Go type %T", value)
	}
}
