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

	"github.com/stretchr/testify/assert"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/quick"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

func booleanTag(value bool) Tag {
	intValue := 0
	if value {
		intValue = 1
	}
	tag := quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_BOOLEAN, intValue, xtob(""))
	return tag
}

// Encode a Boolean object into a tag.
func booleanEncode(obj *Boolean) Tag {
	tag := quick.Tag()
	obj.Encode(tag)
	return tag
}

// Decode a boolean application tag into a boolean.
func booleanDecode(tag Tag) *Boolean {
	obj := quick.Boolean(tag)

	return obj
}

// Pass the value to Boolean, construct a tag from the hex string,
//
//	and compare results of encode and decoding each other.
func booleanEndec(t *testing.T, v bool, x bool) {
	tag := booleanTag(x)

	obj := quick.Boolean(v)

	assert.Equal(t, tag, booleanEncode(obj))
	assert.Equal(t, obj, booleanDecode(tag))
}

func TestBoolean(t *testing.T) {
	obj := quick.Boolean()
	assert.False(t, obj.GetBoolValue())

	assert.Panics(t, func() {
		quick.Boolean("some string")
	})
	assert.Panics(t, func() {
		quick.Boolean(1.0)
	})
}

func TestBooleanBool(t *testing.T) {
	obj := quick.Boolean(false)
	assert.Equal(t, false, obj.GetBoolValue())
	assert.Equal(t, "Boolean(False)", obj.String())

	obj = quick.Boolean(true)
	assert.Equal(t, true, obj.GetBoolValue())
	assert.Equal(t, "Boolean(True)", obj.String())
}

func TestBooleanTag(t *testing.T) {
	tag := quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_BOOLEAN, 1, xtob("01"))
	obj := quick.Boolean(tag)
	assert.Equal(t, obj.GetValue(), 1)

	tag = quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_SIGNED_INTEGER, 0, xtob(""))
	assert.Panics(t, func() {
		quick.Boolean(tag)
	})

	tag = quick.Tag(model.TagClass_CONTEXT_SPECIFIC_TAGS, 0, 1, xtob(""))
	assert.Panics(t, func() {
		quick.Boolean(tag)
	})

	tag = quick.Tag(TagOpeningTagClass, 0)
	assert.Panics(t, func() {
		quick.Boolean(tag)
	})
}

func TestBooleanCopy(t *testing.T) {
	obj1 := quick.Boolean(true)
	obj2 := quick.Boolean(obj1)
	assert.Equal(t, obj1, obj2)
}

func TestBooleanEndec(t *testing.T) {
	booleanEndec(t, false, false)
	booleanEndec(t, true, true)
}
