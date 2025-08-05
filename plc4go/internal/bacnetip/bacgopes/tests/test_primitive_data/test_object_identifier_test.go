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

func ObjectIdentifierTag(x string) Tag {
	b := xtob(x)
	tag := quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_BACNET_OBJECT_IDENTIFIER, len(b), b)
	return tag
}

// Encode a ObjectIdentifier object into a tag.
func ObjectIdentifierEncode(obj *ObjectIdentifier) Tag {
	tag := quick.Tag()
	obj.Encode(tag)
	return tag
}

// Decode a ObjectIdentifier application tag into a ObjectIdentifier.
func ObjectIdentifierDecode(tag Tag) *ObjectIdentifier {
	obj := quick.ObjectIdentifier(tag)

	return obj
}

// Pass the value to ObjectIdentifier, construct a tag from the hex string,
//
//	and compare results of encode and decoding each other.
func ObjectIdentifierEndec(t *testing.T, v any, x string) {
	tag := ObjectIdentifierTag(x)

	obj := quick.ObjectIdentifier(v)

	assert.Equal(t, tag, ObjectIdentifierEncode(obj))
	assert.Equal(t, obj, ObjectIdentifierDecode(tag))
}

func TestObjectIdentifier(t *testing.T) {
	obj := quick.ObjectIdentifier()
	assert.Equal(t, ObjectIdentifierTuple{Left: "analogInput"}, obj.GetValue())

	assert.Panics(t, func() {
		quick.ObjectIdentifier(1.0)
	})
}

func TestObjectIdentifierInt(t *testing.T) {
	obj := quick.ObjectIdentifier(1)
	assert.Equal(t, ObjectIdentifierTuple{Left: "analogInput", Right: 1}, obj.GetValue())
	assert.Equal(t, "ObjectIdentifier(analogInput,1)", obj.String())

	obj = quick.ObjectIdentifier(0x0400002)
	assert.Equal(t, ObjectIdentifierTuple{Left: "analogOutput", Right: 2}, obj.GetValue())
	assert.Equal(t, "ObjectIdentifier(analogOutput,2)", obj.String())
}

func TestObjectIdentifierStr(t *testing.T) {
	obj := quick.ObjectIdentifier("analogInput:1")
	assert.Equal(t, ObjectIdentifierTuple{Left: "analogInput", Right: 1}, obj.GetValue())
	assert.Equal(t, "ObjectIdentifier(analogInput,1)", obj.String())

	obj = quick.ObjectIdentifier("8:123")
	assert.Equal(t, ObjectIdentifierTuple{Left: "device", Right: 123}, obj.GetValue())
	assert.Equal(t, "ObjectIdentifier(device,123)", obj.String())

	assert.Panics(t, func() {
		quick.ObjectIdentifier("x")
	})
	assert.Panics(t, func() {
		quick.ObjectIdentifier(":1")
	})
	assert.Panics(t, func() {
		quick.ObjectIdentifier("1:")
	})
	assert.Panics(t, func() {
		quick.ObjectIdentifier("1:b")
	})
}

func TestObjectIdentifierTuple(t *testing.T) {
	obj := quick.ObjectIdentifier(ObjectIdentifierTuple{Left: "analogInput", Right: 1})
	assert.Equal(t, ObjectIdentifierTuple{Left: "analogInput", Right: 1}, obj.GetValue())
	assert.Equal(t, "ObjectIdentifier(analogInput,1)", obj.String())

	assert.Panics(t, func() {
		quick.ObjectIdentifier(ObjectIdentifierTuple{Left: 0, Right: -1})
	})
	assert.Panics(t, func() {
		quick.ObjectIdentifier(ObjectIdentifierTuple{Left: 0, Right: 0x003FFFFF + 1})
	})
}

func TestObjectIdentifierTag(t *testing.T) {
	tag := quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_BACNET_OBJECT_IDENTIFIER, 1, xtob("06000003"))
	obj := quick.ObjectIdentifier(tag)
	assert.Equal(t, ObjectIdentifierTuple{Left: "pulseConverter", Right: 3}, obj.GetValue())

	tag = quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_BOOLEAN, 0, xtob(""))
	assert.Panics(t, func() {
		quick.ObjectIdentifier(tag)
	})

	tag = quick.Tag(model.TagClass_CONTEXT_SPECIFIC_TAGS, 0, 1, xtob("ff"))
	assert.Panics(t, func() {
		quick.ObjectIdentifier(tag)
	})

	tag = quick.Tag(TagOpeningTagClass, 0)
	assert.Panics(t, func() {
		quick.ObjectIdentifier(tag)
	})
}

func TestObjectIdentifierCopy(t *testing.T) {
	obj1 := quick.ObjectIdentifier(ObjectIdentifierTuple{Left: "pulseConverter", Right: 3})
	obj2 := quick.ObjectIdentifier(obj1)
	assert.Equal(t, ObjectIdentifierTuple{Left: "pulseConverter", Right: 3}, obj2.GetValue())
	assert.Equal(t, obj1, obj2)
}

func TestObjectIdentifierEndec(t *testing.T) {
	assert.Panics(t, func() {
		quick.ObjectIdentifier(ObjectIdentifierTag(""))
	})
	ObjectIdentifierEndec(t, ObjectIdentifierTuple{Left: "analogInput", Right: 0}, "00000000")
}
