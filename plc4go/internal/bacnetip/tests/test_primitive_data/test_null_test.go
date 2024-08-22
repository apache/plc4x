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
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"

	"github.com/stretchr/testify/assert"
)

func Null(arg ...any) *bacnetip.Null {
	if len(arg) == 0 {
		Null, err := bacnetip.NewNull(nil)
		if err != nil {
			panic(err)
		}
		return Null
	}
	Null, err := bacnetip.NewNull(arg[0])
	if err != nil {
		panic(err)
	}
	return Null
}

func NullTag(x string) bacnetip.Tag {
	b := xtob(x)
	tag := Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_NULL, len(b), b)
	return tag
}

// Encode a Null object into a tag.
func NullEncode(obj *bacnetip.Null) bacnetip.Tag {
	tag := Tag()
	obj.Encode(tag)
	return tag
}

// Decode a Null application tag into a Null.
func NullDecode(tag bacnetip.Tag) *bacnetip.Null {
	obj := Null(tag)

	return obj
}

// Pass the value to Null, construct a tag from the hex string,
//
//	and compare results of encode and decoding each other.
func NullEndec(t *testing.T, v any, x string) {
	tag := NullTag(x)

	obj := Null(v)

	assert.Equal(t, tag, NullEncode(obj))
	assert.Equal(t, obj, NullDecode(tag))
}

func TestNull(t *testing.T) {
	obj := Null()
	assert.Equal(t, 0, obj.GetValue())
	assert.Panics(t, func() {
		Null("some string")
	})
	assert.Panics(t, func() {
		Null(1.0)
	})
}

func TestNullNull(t *testing.T) {
	obj := Null()
	assert.Equal(t, 0, obj.GetValue())
}

func TestNullTag(t *testing.T) {
	tag := Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_NULL, 0, xtob(""))
	obj := Null(tag)
	assert.Equal(t, obj.GetValue(), 0)

	tag = Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_BOOLEAN, 0, xtob(""))
	assert.Panics(t, func() {
		Null(tag)
	})

	tag = Tag(model.TagClass_CONTEXT_SPECIFIC_TAGS, 0, 1, xtob("ff"))
	assert.Panics(t, func() {
		Null(tag)
	})

	tag = Tag(bacnetip.TagOpeningTagClass, 0)
	assert.Panics(t, func() {
		Null(tag)
	})
}

func TestNullCopy(t *testing.T) {
	obj1 := Null()
	obj2 := Null(obj1)
	assert.Equal(t, 0, obj2.GetValue())
	assert.Equal(t, obj1, obj2)
}

func TestNullEndec(t *testing.T) {
	NullEndec(t, nil, "")
}
