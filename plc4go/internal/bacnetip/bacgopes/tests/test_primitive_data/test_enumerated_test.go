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

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/quick"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type quickBrownFox struct {
	*Enumerated
}

func (q *quickBrownFox) GetEnumerations() map[string]uint64 {
	return map[string]uint64{
		"quick": 0,
		"brown": 1,
		"fox":   2,
	}
}

func (q *quickBrownFox) SetEnumerated(enumerated *Enumerated) {
	q.Enumerated = enumerated
}

func QuickBrownFox(args ...any) *quickBrownFox {
	q := &quickBrownFox{}
	var err error
	q.Enumerated, err = NewEnumerated(NA(append([]any{q}, args...)...))
	if err != nil {
		panic(err)
	}
	return q
}

func EnumeratedTag(x string) Tag {
	b := xtob(x)
	tag := quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_ENUMERATED, len(b), b)
	return tag
}

// Encode a Enumerated object into a tag.
func EnumeratedEncode(obj *Enumerated) Tag {
	tag := quick.Tag()
	obj.Encode(tag)
	return tag
}

// Decode a Enumerated application tag into a Enumerated.
func EnumeratedDecode(tag Tag) *Enumerated {
	obj := quick.Enumerated(tag)

	return obj
}

// Pass the value to Enumerated, construct a tag from the hex string,
//
//	and compare results of encode and decoding each other.
func EnumeratedEndec(t *testing.T, v uint64, x string) {
	tag := EnumeratedTag(x)

	obj := quick.Enumerated(v)

	assert.Equal(t, tag, EnumeratedEncode(obj))
	assert.Equal(t, obj, EnumeratedDecode(tag))
}

func TestEnumerated(t *testing.T) {
	obj := quick.Enumerated()
	assert.Equal(t, uint64(0), obj.GetValue())

	assert.Panics(t, func() {
		quick.Enumerated("label")
	})
	assert.Panics(t, func() {
		quick.Enumerated(1.0)
	})
}

func TestEnumeratedInt(t *testing.T) {
	obj := quick.Enumerated(1)
	assert.Equal(t, uint64(1), obj.GetValue())
	assert.Equal(t, "Enumerated(1)", obj.String())

	assert.Panics(t, func() {
		quick.Enumerated(-1)
	})
}

func TestEnumeratedStr(t *testing.T) {
	obj := QuickBrownFox("quick")
	assert.Equal(t, "Enumerated(quick)", obj.String())

	assert.Panics(t, func() {
		quick.Enumerated(-1)
	})
	assert.Panics(t, func() {
		quick.Enumerated("lazyDog")
	})

	tag := quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_ENUMERATED, 1, xtob("01"))
	obj = QuickBrownFox(tag)
	assert.Equal(t, "Enumerated(brown)", obj.String())
}

func TestEnumeratedTag(t *testing.T) {
	tag := quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_ENUMERATED, 1, xtob("01"))
	obj := quick.Enumerated(tag)
	assert.Equal(t, obj.GetValue(), uint64(1))

	tag = quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_BOOLEAN, 0, xtob(""))
	assert.Panics(t, func() {
		quick.Enumerated(tag)
	})

	tag = quick.Tag(model.TagClass_CONTEXT_SPECIFIC_TAGS, 0, 1, xtob("ff"))
	assert.Panics(t, func() {
		quick.Enumerated(tag)
	})

	tag = quick.Tag(TagOpeningTagClass, 0)
	assert.Panics(t, func() {
		quick.Enumerated(tag)
	})
}

func TestEnumeratedCopy(t *testing.T) {
	obj1 := quick.Enumerated(12)
	obj2 := quick.Enumerated(obj1)
	assert.Equal(t, uint64(12), obj2.GetValue())
	assert.Equal(t, obj1, obj2)
}

func TestEnumeratedEndec(t *testing.T) {
	assert.Panics(t, func() {
		quick.Enumerated(EnumeratedTag(""))
	})
	EnumeratedEndec(t, 0, "00")
	EnumeratedEndec(t, 1, "01")
	EnumeratedEndec(t, 127, "7f")
	EnumeratedEndec(t, 128, "80")
	EnumeratedEndec(t, 255, "ff")
	EnumeratedEndec(t, 32767, "7fff")
	EnumeratedEndec(t, 32768, "8000")
	EnumeratedEndec(t, 8388607, "7fffff")
	EnumeratedEndec(t, 8388608, "800000")
	EnumeratedEndec(t, 2147483647, "7fffffff")
	EnumeratedEndec(t, 2147483648, "80000000")
}
