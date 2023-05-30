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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlcSimpleNumericValueAdapter_IsBool(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsBool())
}

func TestPlcSimpleNumericValueAdapter_IsByte(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsByte())
}

func TestPlcSimpleNumericValueAdapter_IsFloat32(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsFloat32())
}

func TestPlcSimpleNumericValueAdapter_IsFloat64(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsFloat64())
}

func TestPlcSimpleNumericValueAdapter_IsInt16(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsInt16())
}

func TestPlcSimpleNumericValueAdapter_IsInt32(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsInt32())
}

func TestPlcSimpleNumericValueAdapter_IsInt64(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsInt64())
}

func TestPlcSimpleNumericValueAdapter_IsInt8(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsInt8())
}

func TestPlcSimpleNumericValueAdapter_IsString(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsString())
}

func TestPlcSimpleNumericValueAdapter_IsUint16(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsUint16())
}

func TestPlcSimpleNumericValueAdapter_IsUint32(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsUint32())
}

func TestPlcSimpleNumericValueAdapter_IsUint64(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsUint64())
}

func TestPlcSimpleNumericValueAdapter_IsUint8(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsUint8())
}

func TestPlcSimpleNumericValueAdapter_String(t *testing.T) {
	assert.Equal(t, "not implemented", PlcSimpleNumericValueAdapter{}.String())
}

func TestPlcSimpleValueAdapter_GetLength(t *testing.T) {
	assert.Equal(t, uint32(1), PlcSimpleNumericValueAdapter{}.GetLength())
}

func TestPlcSimpleValueAdapter_IsSimple(t *testing.T) {
	assert.True(t, PlcSimpleNumericValueAdapter{}.IsSimple())
}

func TestPlcSimpleValueAdapter_String(t *testing.T) {
	assert.Equal(t, "not implemented", PlcSimpleValueAdapter{}.String())
}

func TestPlcValueAdapter_GetBool(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.GetBool())
}

func TestPlcValueAdapter_GetBoolArray(t *testing.T) {
	assert.Nil(t, PlcValueAdapter{}.GetBoolArray())
}

func TestPlcValueAdapter_GetBoolAt(t *testing.T) {
	t.Run("index 0", func(t *testing.T) {
		assert.False(t, PlcValueAdapter{}.GetBoolAt(0))
	})
	t.Run("index 1", func(t *testing.T) {
		assert.False(t, PlcValueAdapter{}.GetBoolAt(1))
	})
}

func TestPlcValueAdapter_GetBoolLength(t *testing.T) {
	assert.Equal(t, uint32(1), PlcValueAdapter{}.GetBoolLength())
}

func TestPlcValueAdapter_GetByte(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.GetByte()
	t.Error("above should panic")
}

func TestPlcValueAdapter_GetDate(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.GetDate()
	t.Error("above should panic")
}

func TestPlcValueAdapter_GetDateTime(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.GetDateTime()
	t.Error("above should panic")
}

func TestPlcValueAdapter_GetDuration(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.GetDuration()
	t.Error("above should panic")
}

func TestPlcValueAdapter_GetFloat32(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.GetFloat32()
	t.Error("above should panic")
}

func TestPlcValueAdapter_GetFloat64(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.GetFloat64()
	t.Error("above should panic")
}

func TestPlcValueAdapter_GetIndex(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.GetIndex(0)
	t.Error("above should panic")
}

func TestPlcValueAdapter_GetInt16(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.GetInt16()
	t.Error("above should panic")
}

func TestPlcValueAdapter_GetInt32(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.GetInt32()
	t.Error("above should panic")
}

func TestPlcValueAdapter_GetInt64(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.GetInt64()
	t.Error("above should panic")
}

func TestPlcValueAdapter_GetInt8(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.GetInt8()
	t.Error("above should panic")
}

func TestPlcValueAdapter_GetKeys(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.GetKeys()
	t.Error("above should panic")
}

func TestPlcValueAdapter_GetLength(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.GetLength()
	t.Error("above should panic")
}

func TestPlcValueAdapter_GetList(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.GetList()
	t.Error("above should panic")
}

func TestPlcValueAdapter_GetPlcValueType(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.GetPlcValueType()
	t.Error("above should panic")
}

func TestPlcValueAdapter_GetRaw(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.GetRaw()
	t.Error("above should panic")
}

func TestPlcValueAdapter_GetString(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.GetString()
	t.Error("above should panic")
}

func TestPlcValueAdapter_GetStruct(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.GetStruct()
	t.Error("above should panic")
}

func TestPlcValueAdapter_GetTime(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.GetTime()
	t.Error("above should panic")
}

func TestPlcValueAdapter_GetUint16(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.GetUint16()
	t.Error("above should panic")
}

func TestPlcValueAdapter_GetUint32(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.GetUint32()
	t.Error("above should panic")
}

func TestPlcValueAdapter_GetUint64(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.GetUint64()
	t.Error("above should panic")
}

func TestPlcValueAdapter_GetUint8(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.GetUint8()
	t.Error("above should panic")
}

func TestPlcValueAdapter_GetValue(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.GetValue("")
	t.Error("above should panic")
}

func TestPlcValueAdapter_HasKey(t *testing.T) {
	defer func() {
		recover()
	}()
	PlcValueAdapter{}.HasKey("")
	t.Error("above should panic")
}

func TestPlcValueAdapter_IsBool(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsBool())
}

func TestPlcValueAdapter_IsByte(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsByte())
}

func TestPlcValueAdapter_IsDate(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsDate())
}

func TestPlcValueAdapter_IsDateTime(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsDateTime())
}

func TestPlcValueAdapter_IsDuration(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsDuration())
}

func TestPlcValueAdapter_IsFloat32(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsFloat32())
}

func TestPlcValueAdapter_IsFloat64(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsFloat64())
}

func TestPlcValueAdapter_IsInt16(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsInt16())
}

func TestPlcValueAdapter_IsInt32(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsInt32())
}

func TestPlcValueAdapter_IsInt64(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsInt64())
}

func TestPlcValueAdapter_IsInt8(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsInt8())
}

func TestPlcValueAdapter_IsList(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsList())
}

func TestPlcValueAdapter_IsNull(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsNull())
}

func TestPlcValueAdapter_IsNullable(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsNullable())
}

func TestPlcValueAdapter_IsRaw(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsRaw())
}

func TestPlcValueAdapter_IsSimple(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsSimple())
}

func TestPlcValueAdapter_IsString(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsString())
}

func TestPlcValueAdapter_IsStruct(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsStruct())
}

func TestPlcValueAdapter_IsTime(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsTime())
}

func TestPlcValueAdapter_IsUint16(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsUint16())
}

func TestPlcValueAdapter_IsUint32(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsUint32())
}

func TestPlcValueAdapter_IsUint64(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsUint64())
}

func TestPlcValueAdapter_IsUint8(t *testing.T) {
	assert.False(t, PlcValueAdapter{}.IsUint8())
}

func TestPlcValueAdapter_String(t *testing.T) {
	assert.Equal(t, "not implemented", PlcValueAdapter{}.String())
}
