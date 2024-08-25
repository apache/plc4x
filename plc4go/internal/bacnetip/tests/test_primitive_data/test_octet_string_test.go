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

func OctetStringTag(x string) bacnetip.Tag {
	b := xtob(x)
	tag := Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_OCTET_STRING, len(b), b)
	return tag
}

// Encode a OctetString object into a tag.
func OctetStringEncode(obj *bacnetip.OctetString) bacnetip.Tag {
	tag := Tag()
	obj.Encode(tag)
	return tag
}

// Decode a OctetString application tag into a OctetString.
func OctetStringDecode(tag bacnetip.Tag) *bacnetip.OctetString {
	obj := OctetString(tag)

	return obj
}

// Pass the value to OctetString, construct a tag from the hex string,
//
//	and compare results of encode and decoding each other.
func OctetStringEndec(t *testing.T, x string) {
	tag := OctetStringTag(x)

	obj := OctetString(xtob(x))

	assert.Equal(t, tag, OctetStringEncode(obj))
	assert.Equal(t, obj, OctetStringDecode(tag))
}

func TestOctetString(t *testing.T) {
	obj := OctetString()
	assert.Equal(t, []byte{}, obj.GetValue())

	assert.Panics(t, func() {
		OctetString(1)
	})
}

func TestOctetStringOctetString(t *testing.T) {
	obj := OctetString(xtob("01"))
	assert.Equal(t, xtob("01"), obj.GetValue())
	assert.Equal(t, "OctetString(X'01')", obj.String())

	obj = OctetString(xtob("01020304"))
	assert.Equal(t, xtob("01020304"), obj.GetValue())
	assert.Equal(t, "OctetString(X'01020304')", obj.String())
}

func TestOctetStringTag(t *testing.T) {
	tag := Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_OCTET_STRING, 1, xtob("00"))
	obj := OctetString(tag)
	assert.Equal(t, xtob("00"), obj.GetValue())

	tag = Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_BOOLEAN, 0, xtob(""))
	assert.Panics(t, func() {
		OctetString(tag)
	})

	tag = Tag(model.TagClass_CONTEXT_SPECIFIC_TAGS, 0, 1, xtob("ff"))
	assert.Panics(t, func() {
		OctetString(tag)
	})

	tag = Tag(bacnetip.TagOpeningTagClass, 0)
	assert.Panics(t, func() {
		OctetString(tag)
	})
}

func TestOctetStringCopy(t *testing.T) {
	obj1 := OctetString(xtob("01"))
	obj2 := OctetString(obj1)
	assert.Equal(t, xtob("01"), obj2.GetValue())
	assert.Equal(t, obj1, obj2)
}

func TestOctetStringEndec(t *testing.T) {
	OctetStringEndec(t, "01")
	OctetStringEndec(t, "0101")
	OctetStringEndec(t, "010103")
	OctetStringEndec(t, "01010304")
}
