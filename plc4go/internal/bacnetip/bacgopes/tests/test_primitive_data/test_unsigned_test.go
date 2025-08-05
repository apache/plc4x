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

	"github.com/stretchr/testify/assert"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/quick"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

func UnsignedTag(x string) Tag {
	b := xtob(x)
	tag := quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_UNSIGNED_INTEGER, len(b), b)
	return tag
}

// Encode a Unsigned object into a tag.
func UnsignedEncode(obj *Unsigned) Tag {
	tag := quick.Tag()
	obj.Encode(tag)
	return tag
}

// Decode a Unsigned application tag into a Unsigned.
func UnsignedDecode(tag Tag) *Unsigned {
	obj := quick.Unsigned(tag)

	return obj
}

// Pass the value to Unsigned, construct a tag from the hex string,
//
//	and compare results of encode and decoding each other.
func UnsignedEndec(t *testing.T, v uint32, x string) {
	tag := UnsignedTag(x)

	obj := quick.Unsigned(v)

	assert.Equal(t, tag, UnsignedEncode(obj))
	assert.Equal(t, obj, UnsignedDecode(tag))
}

func TestUnsigned(t *testing.T) {
	obj := quick.Unsigned()
	assert.Equal(t, uint32(0), obj.GetValue())

	assert.True(t, obj.IsValid(1))
	assert.True(t, obj.IsValid("1"))
	assert.True(t, obj.IsValid(math.MaxInt64))
	assert.False(t, obj.IsValid(math.MinInt64))

	assert.False(t, obj.IsValid(true))
	assert.False(t, obj.IsValid(-1))
	assert.False(t, obj.IsValid(1.0))
	assert.Panics(t, func() {
		quick.Unsigned("some string")
	})
	assert.Panics(t, func() {
		quick.Unsigned(1.0)
	})
}

func TestUnsignedInt(t *testing.T) {
	obj := quick.Unsigned(1)
	assert.Equal(t, uint32(1), obj.GetValue())
	assert.Equal(t, "Unsigned(1)", obj.String())

	assert.Panics(t, func() {
		quick.Unsigned(-1)
	})
}

func TestUnsignedInt8(t *testing.T) {
	obj := quick.Unsigned8(1)
	assert.Equal(t, uint8(1), obj.GetValue())
	assert.Equal(t, "Unsigned8(1)", obj.String())

	assert.Panics(t, func() {
		quick.Unsigned8(256)
	})
}

func TestUnsignedInt16(t *testing.T) {
	obj := quick.Unsigned16(1)
	assert.Equal(t, uint16(1), obj.GetValue())
	assert.Equal(t, "Unsigned16(1)", obj.String())

	assert.Panics(t, func() {
		quick.Unsigned16(65536)
	})
}

func TestUnsignedTag(t *testing.T) {
	tag := quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_UNSIGNED_INTEGER, 1, xtob("01"))
	obj := quick.Unsigned(tag)
	assert.Equal(t, obj.GetValue(), uint32(1))

	tag = quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_BOOLEAN, 0, xtob(""))
	assert.Panics(t, func() {
		quick.Unsigned(tag)
	})

	tag = quick.Tag(model.TagClass_CONTEXT_SPECIFIC_TAGS, 0, 1, xtob("ff"))
	assert.Panics(t, func() {
		quick.Unsigned(tag)
	})

	tag = quick.Tag(TagOpeningTagClass, 0)
	assert.Panics(t, func() {
		quick.Unsigned(tag)
	})
}

func TestUnsignedCopy(t *testing.T) {
	obj1 := quick.Unsigned(12)
	obj2 := quick.Unsigned(obj1)
	assert.Equal(t, uint32(12), obj2.GetValue())
	assert.Equal(t, obj1, obj2)
}

func TestUnsignedEndec(t *testing.T) {
	assert.Panics(t, func() {
		quick.Unsigned(UnsignedTag(""))
	})
	UnsignedEndec(t, 0, "00")
	UnsignedEndec(t, 1, "01")
	UnsignedEndec(t, 127, "7f")
	UnsignedEndec(t, 128, "80")
	UnsignedEndec(t, 255, "ff")

	UnsignedEndec(t, 32767, "7fff")
	UnsignedEndec(t, 32768, "8000")

	UnsignedEndec(t, 8388607, "7fffff")
	UnsignedEndec(t, 8388608, "800000")

	UnsignedEndec(t, 2147483647, "7fffffff")
	UnsignedEndec(t, 2147483648, "80000000")
}
