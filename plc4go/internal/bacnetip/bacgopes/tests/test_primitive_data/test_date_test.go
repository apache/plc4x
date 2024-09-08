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

// Convert a hex string to a character_string application tag.
func DateTag(x string) Tag {
	b := xtob(x)
	tag := quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_DATE, len(b), b)
	return tag
}

// Encode a Date object into a tag.
func DateEncode(obj *Date) Tag {
	tag := quick.Tag()
	obj.Encode(tag)
	return tag
}

// Decode a Date application tag into a Date.
func DateDecode(tag Tag) *Date {
	obj := quick.Date(tag)

	return obj
}

// Pass the value to Date, construct a tag from the hex string,
//
//	and compare results of encode and decoding each other.
func DateEndec(t *testing.T, v string, x string) {
	tag := DateTag(x)

	obj := quick.Date(v)

	assert.Equal(t, tag, DateEncode(obj))
	assert.Equal(t, obj, DateDecode(tag))
}

func TestDate(t *testing.T) {
	obj := quick.Date()
	assert.Equal(t, DateTuple{Year: 0xff, Month: 0xff, Day: 0xff, DayOfWeek: 0xff}, obj.GetValue())

	assert.Panics(t, func() {
		quick.Date("some string")
	})
}

func TestDateTuple(t *testing.T) {
	obj := quick.Date(DateTuple{Year: 1, Month: 2, Day: 3, DayOfWeek: 4})
	assert.Equal(t, DateTuple{Year: 1, Month: 2, Day: 3, DayOfWeek: 4}, obj.GetValue())
	assert.Equal(t, "Date(1901-2-3 thu)", obj.String())
}

func TestDateTag(t *testing.T) {
	tag := quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_DATE, 4, xtob("01020304"))
	obj := quick.Date(tag)
	assert.Equal(t, DateTuple{Year: 1, Month: 2, Day: 3, DayOfWeek: 4}, obj.GetValue())

	tag = quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_BOOLEAN, 0, xtob(""))
	assert.Panics(t, func() {
		quick.Date(tag)
	})

	tag = quick.Tag(model.TagClass_CONTEXT_SPECIFIC_TAGS, 0, 1, xtob("ff"))
	assert.Panics(t, func() {
		quick.Date(tag)
	})

	tag = quick.Tag(TagOpeningTagClass, 0)
	assert.Panics(t, func() {
		quick.Date(tag)
	})
}

func TestDateCopy(t *testing.T) {
	obj1 := quick.Date(DateTuple{Year: 1, Month: 2, Day: 3, DayOfWeek: 4})
	obj2 := quick.Date(obj1)
	assert.Equal(t, obj1, obj2)
}

func TestDateNow(t *testing.T) {
	// TODO: upstream doesn't tests this either
}

func TestDateEndec(t *testing.T) {
	assert.Panics(t, func() {
		quick.Date(DateTag(""))
	})
}

func TestDateArgs(t *testing.T) {
	tag := quick.Tag()
	date := quick.Date(nil, 2023, 2, 10)
	date1 := quick.Date(nil, 123, 2, 10)
	assert.Equal(t, date, date1)
	date.Encode(tag)
}
