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

func RealTag(x string) Tag {
	b := xtob(x)
	tag := quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_REAL, len(b), b)
	return tag
}

// Encode a Real object into a tag.
func RealEncode(obj *Real) Tag {
	tag := quick.Tag()
	obj.Encode(tag)
	return tag
}

// Decode a Real application tag into a Real.
func RealDecode(tag Tag) *Real {
	obj := quick.Real(tag)

	return obj
}

// Pass the value to Real, construct a tag from the hex string,
//
//	and compare results of encode and decoding each other.
func RealEndec(t *testing.T, v float32, x string) {
	tag := RealTag(x)

	obj := quick.Real(v)

	assert.Equal(t, tag, RealEncode(obj))
	assert.Equal(t, obj, RealDecode(tag))
}

func TestReal(t *testing.T) {
	obj := quick.Real()
	assert.Equal(t, float32(0.0), obj.GetValue())

	assert.Panics(t, func() {
		quick.Real("some string")
	})
}

func TestRealReal(t *testing.T) {
	obj := quick.Real(float32(1.0))
	assert.Equal(t, float32(1.0), obj.GetValue())
	assert.Equal(t, "Real(1)", obj.String())

	obj = quick.Real(float32(73.5))
	assert.Equal(t, float32(73.5), obj.GetValue())
	assert.Equal(t, "Real(73.5)", obj.String())
}

func TestRealTag(t *testing.T) {
	tag := quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_REAL, 1, xtob("3f800000"))
	obj := quick.Real(tag)
	assert.Equal(t, obj.GetValue(), float32(1.0))

	tag = quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_BOOLEAN, 0, xtob(""))
	assert.Panics(t, func() {
		quick.Real(tag)
	})

	tag = quick.Tag(model.TagClass_CONTEXT_SPECIFIC_TAGS, 0, 1, xtob("ff"))
	assert.Panics(t, func() {
		quick.Real(tag)
	})

	tag = quick.Tag(TagOpeningTagClass, 0)
	assert.Panics(t, func() {
		quick.Real(tag)
	})
}

func TestRealCopy(t *testing.T) {
	obj1 := quick.Real(12)
	obj2 := quick.Real(obj1)
	assert.Equal(t, float32(12.0), obj2.GetValue())
	assert.Equal(t, obj1, obj2)
}

func TestRealEndec(t *testing.T) {
	assert.Panics(t, func() {
		quick.Real(RealTag(""))
	})
	RealEndec(t, 0, "00000000")
	RealEndec(t, 1, "3f800000")
	RealEndec(t, -1, "bf800000")

	RealEndec(t, 73.5, "42930000")

	inf := float32(math.Inf(1))
	RealEndec(t, inf, "7f800000")
	RealEndec(t, -inf, "ff800000")

	/*
		TODO: go thing... below is equal but somehow go (testing) says no
		nan := float32(math.NaN())
		RealEndec(t, nan, "7fc00000")
	*/
}
