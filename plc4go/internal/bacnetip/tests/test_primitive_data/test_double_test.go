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
	"testing"

	"github.com/apache/plc4x/plc4go/internal/bacnetip"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/constructors"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"

	"github.com/stretchr/testify/assert"
)

func DoubleTag(x string) bacnetip.Tag {
	b := xtob(x)
	tag := Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_DOUBLE, len(b), b)
	return tag
}

// Encode a Double object into a tag.
func DoubleEncode(obj *bacnetip.Double) bacnetip.Tag {
	tag := Tag()
	obj.Encode(tag)
	return tag
}

// Decode a Double application tag into a Double.
func DoubleDecode(tag bacnetip.Tag) *bacnetip.Double {
	obj := Double(tag)

	return obj
}

// Pass the value to Double, construct a tag from the hex string,
//
//	and compare results of encode and decoding each other.
func DoubleEndec(t *testing.T, v float64, x string) {
	tag := DoubleTag(x)

	obj := Double(v)

	assert.Equal(t, tag, DoubleEncode(obj))
	assert.Equal(t, obj, DoubleDecode(tag))
}

func TestDouble(t *testing.T) {
	obj := Double()
	assert.Equal(t, 0.0, obj.GetValue())

	assert.Panics(t, func() {
		Double("some string")
	})
}

func TestDoubleDouble(t *testing.T) {
	obj := Double(1.0)
	assert.Equal(t, 1.0, obj.GetValue())
	assert.Equal(t, "Double(1)", obj.String())

	obj = Double(73.5)
	assert.Equal(t, 73.5, obj.GetValue())
	assert.Equal(t, "Double(73.5)", obj.String())
}

func TestDoubleTag(t *testing.T) {
	tag := Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_DOUBLE, 1, xtob("3ff0000000000000"))
	obj := Double(tag)
	assert.Equal(t, obj.GetValue(), 1.0)

	tag = Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_BOOLEAN, 0, xtob(""))
	assert.Panics(t, func() {
		Double(tag)
	})

	tag = Tag(model.TagClass_CONTEXT_SPECIFIC_TAGS, 0, 1, xtob("ff"))
	assert.Panics(t, func() {
		Double(tag)
	})

	tag = Tag(bacnetip.TagOpeningTagClass, 0)
	assert.Panics(t, func() {
		Double(tag)
	})
}

func TestDoubleCopy(t *testing.T) {
	obj1 := Double(12)
	obj2 := Double(obj1)
	assert.Equal(t, 12.0, obj2.GetValue())
	assert.Equal(t, obj1, obj2)
}

func TestDoubleEndec(t *testing.T) {
	assert.Panics(t, func() {
		Double(DoubleTag(""))
	})
	DoubleEndec(t, 0, "0000000000000000")
	DoubleEndec(t, 1, "3ff0000000000000")
}
