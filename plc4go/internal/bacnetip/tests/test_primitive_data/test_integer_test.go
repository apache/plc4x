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

package test_primitive_data

import (
	"math"
	"testing"

	"github.com/apache/plc4x/plc4go/internal/bacnetip"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/constructors"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"

	"github.com/stretchr/testify/assert"
)

func IntegerTag(x string) bacnetip.Tag {
	b, err := bacnetip.Xtob(x)
	if err != nil {
		panic(err)
	}
	tag := Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_SIGNED_INTEGER, len(b), b)
	return tag
}

// Encode a Integer object into a tag.
func IntegerEncode(obj *bacnetip.Integer) bacnetip.Tag {
	tag := Tag()
	obj.Encode(tag)
	return tag
}

// Decode a Integer application tag into a Integer.
func IntegerDecode(tag bacnetip.Tag) *bacnetip.Integer {
	obj := Integer(tag)

	return obj
}

// Pass the value to Integer, construct a tag from the hex string,
//
//	and compare results of encode and decoding each other.
func IntegerEndec(t *testing.T, v int32, x string) {
	tag := IntegerTag(x)

	obj := Integer(v)

	assert.Equal(t, tag, IntegerEncode(obj))
	assert.Equal(t, obj, IntegerDecode(tag))
}

func TestInteger(t *testing.T) {
	obj := Integer()
	assert.Equal(t, int32(0), obj.GetValue())

	assert.True(t, obj.IsValid(1))
	assert.True(t, obj.IsValid(-1))
	assert.False(t, obj.IsValid(math.MaxInt64))
	assert.False(t, obj.IsValid(math.MinInt64))

	assert.False(t, obj.IsValid(true))
	assert.False(t, obj.IsValid(1.0))
	assert.Panics(t, func() {
		Integer("some string")
	})
	assert.Panics(t, func() {
		Integer(1.0)
	})
}

func TestIntegerInt(t *testing.T) {
	obj := Integer(1)
	assert.Equal(t, int32(1), obj.GetValue())
	assert.Equal(t, "Integer(1)", obj.String())

	obj = Integer(-1)
	assert.Equal(t, int32(-1), obj.GetValue())
	assert.Equal(t, "Integer(-1)", obj.String())
}

func TestIntegerTag(t *testing.T) {
	tag := Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_SIGNED_INTEGER, 1, xtob("01"))
	obj := Integer(tag)
	assert.Equal(t, obj.GetValue(), int32(1))

	tag = Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_BOOLEAN, 0, xtob(""))
	assert.Panics(t, func() {
		Integer(tag)
	})

	tag = Tag(model.TagClass_CONTEXT_SPECIFIC_TAGS, 0, 1, xtob("ff"))
	assert.Panics(t, func() {
		Integer(tag)
	})

	tag = Tag(bacnetip.TagOpeningTagClass, 0)
	assert.Panics(t, func() {
		Integer(tag)
	})
}

func TestIntegerCopy(t *testing.T) {
	obj1 := Integer(12)
	obj2 := Integer(obj1)
	assert.Equal(t, int32(12), obj2.GetValue())
	assert.Equal(t, obj1, obj2)
}

func TestIntegerEndec(t *testing.T) {
	assert.Panics(t, func() {
		Integer(IntegerTag(""))
	})
	IntegerEndec(t, 0, "00")
	IntegerEndec(t, 1, "01")
	IntegerEndec(t, 127, "7f")
	IntegerEndec(t, -128, "80")
	IntegerEndec(t, -1, "ff")

	IntegerEndec(t, 32767, "7fff")
	IntegerEndec(t, -32768, "8000")

	IntegerEndec(t, 8388607, "7fffff")
	IntegerEndec(t, -8388608, "800000")

	IntegerEndec(t, 2147483647, "7fffffff")
	IntegerEndec(t, -2147483648, "80000000")
}
