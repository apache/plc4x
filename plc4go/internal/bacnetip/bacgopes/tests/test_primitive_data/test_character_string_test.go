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
	"github.com/stretchr/testify/require"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/quick"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

const foxMessage = "the quick brown fox jumped over the lazy dog"

// Convert a hex string to a character_string application tag.
func CharacterStringTag(x string) Tag {
	b := xtob(x)
	tag := quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_CHARACTER_STRING, len(b), b)
	return tag
}

// Encode a CharacterString object into a tag.
func CharacterStringEncode(obj *CharacterString) Tag {
	tag := quick.Tag()
	obj.Encode(tag)
	return tag
}

// Decode a CharacterString application tag into a CharacterString.
func CharacterStringDecode(tag Tag) *CharacterString {
	obj := quick.CharacterString(tag)

	return obj
}

// Pass the value to CharacterString, construct a tag from the hex string,
//
//	and compare results of encode and decoding each other.
func CharacterStringEndec(t *testing.T, v string, x string) {
	tag := CharacterStringTag(x)

	obj := quick.CharacterString(v)

	assert.Equal(t, tag, CharacterStringEncode(obj))
	assert.Equal(t, obj, CharacterStringDecode(tag))
}

func TestCharacterString(t *testing.T) {
	obj := quick.CharacterString()
	assert.Equal(t, "", obj.GetValue())

	assert.Panics(t, func() {
		quick.CharacterString(1)
	})
	assert.Panics(t, func() {
		quick.CharacterString(1.0)
	})
}

func TestCharacterStringStr(t *testing.T) {
	obj := quick.CharacterString("hello")
	assert.Equal(t, "hello", obj.GetValue())
	assert.Equal(t, "CharacterString(0,X'68656c6c6f')", obj.String())
}

func TestCharacterStringStrUnicode(t *testing.T) {
	obj := quick.CharacterString("hello")
	assert.Equal(t, "hello", obj.GetValue())
	assert.Equal(t, "CharacterString(0,X'68656c6c6f')", obj.String())
}

func TestCharacterStringStrUnicodeWithLatin(t *testing.T) {
	// some controllers encoding character string mixing latin-1 and utf-8
	// try to cover those cases without failing
	b := xtob("0030b043") // zero degress celsius
	tag := quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_CHARACTER_STRING, len(b), b)
	obj := quick.CharacterString()
	err := obj.Decode(tag)
	require.NoError(t, err)
	assert.Equal(t, "CharacterString(0,X'30b043')", obj.String())

	assert.Equal(t, "0\xb0C", obj.GetValue())
}

func TestCharacterStringTag(t *testing.T) {
	tag := quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_CHARACTER_STRING, 1, xtob("00"))
	obj := quick.CharacterString(tag)
	assert.Equal(t, obj.GetValue(), "")

	tag = quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_BOOLEAN, 0, xtob(""))
	assert.Panics(t, func() {
		quick.CharacterString(tag)
	})

	tag = quick.Tag(model.TagClass_CONTEXT_SPECIFIC_TAGS, 0, 1, xtob("ff"))
	assert.Panics(t, func() {
		quick.CharacterString(tag)
	})

	tag = quick.Tag(TagOpeningTagClass, 0)
	assert.Panics(t, func() {
		quick.CharacterString(tag)
	})
}

func TestCharacterStringCopy(t *testing.T) {
	obj1 := quick.CharacterString(foxMessage)
	obj2 := quick.CharacterString(obj1)
	assert.Equal(t, obj1, obj2)
}

func TestCharacterStringEndec(t *testing.T) {
	assert.Panics(t, func() {
		quick.CharacterString(CharacterStringTag(""))
	})
	CharacterStringEndec(t, "", "00")
	CharacterStringEndec(t, "abc", "00616263")
}
