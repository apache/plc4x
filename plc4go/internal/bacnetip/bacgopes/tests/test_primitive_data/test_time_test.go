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
func timeTag(x string) Tag {
	b := xtob(x)
	tag := quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_TIME, len(b), b)
	return tag
}

// Encode a Time object into a tag.
func timeEncode(obj *Time) Tag {
	tag := quick.Tag()
	obj.Encode(tag)
	return tag
}

// Decode a Time application tag into a Time.
func timeDecode(tag Tag) *Time {
	obj := quick.Time(tag)

	return obj
}

// Pass the value to Time, construct a tag from the hex string,
//
//	and compare results of encode and decoding each other.
func timeEndec(t *testing.T, v TimeTuple, x string) {
	tag := timeTag(x)

	obj := quick.Time(v)

	assert.Equal(t, tag, timeEncode(obj))
	assert.Equal(t, obj, timeDecode(tag))
}

func TestTime(t *testing.T) {
	obj := quick.Time()
	assert.Equal(t, TimeTuple{Hour: 0xff, Minute: 0xff, Second: 0xff, Hundredth: 0xff}, obj.GetValue())

	assert.Panics(t, func() {
		quick.Time("some string")
	})
}

func TestTimeTuple(t *testing.T) {
	obj := quick.Time(TimeTuple{Hour: 1, Minute: 2, Second: 3, Hundredth: 4})
	assert.Equal(t, TimeTuple{Hour: 1, Minute: 2, Second: 3, Hundredth: 4}, obj.GetValue())
	assert.Equal(t, "Time(01:02:03.04)", obj.String())

	assert.Equal(t, TimeTuple{Hour: 1, Minute: 2, Second: 0, Hundredth: 0}, quick.Time("1:2").GetValue())
	assert.Equal(t, TimeTuple{Hour: 1, Minute: 2, Second: 3, Hundredth: 0}, quick.Time("1:2:3").GetValue())
	assert.Equal(t, TimeTuple{Hour: 1, Minute: 2, Second: 3, Hundredth: 40}, quick.Time("1:2:3.4").GetValue())
	assert.Equal(t, TimeTuple{Hour: 1, Minute: 255, Second: 255, Hundredth: 255}, quick.Time("1:*").GetValue())
	assert.Equal(t, TimeTuple{Hour: 1, Minute: 2, Second: 255, Hundredth: 255}, quick.Time("1:2:*").GetValue())
	assert.Equal(t, TimeTuple{Hour: 1, Minute: 2, Second: 3, Hundredth: 255}, quick.Time("1:2:3.*").GetValue())
}

func TestTimeTag(t *testing.T) {
	tag := quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_TIME, 4, xtob("01020304"))
	obj := quick.Time(tag)
	assert.Equal(t, TimeTuple{Hour: 1, Minute: 2, Second: 3, Hundredth: 4}, obj.GetValue())

	tag = quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_BOOLEAN, 0, xtob(""))
	assert.Panics(t, func() {
		quick.Time(tag)
	})

	tag = quick.Tag(model.TagClass_CONTEXT_SPECIFIC_TAGS, 0, 1, xtob("ff"))
	assert.Panics(t, func() {
		quick.Time(tag)
	})

	tag = quick.Tag(TagOpeningTagClass, 0)
	assert.Panics(t, func() {
		quick.Time(tag)
	})
}

func TestTimeCopy(t *testing.T) {
	obj1 := quick.Time(TimeTuple{Hour: 1, Minute: 2, Second: 3, Hundredth: 4})
	obj2 := quick.Time(obj1)
	assert.Equal(t, obj1, obj2)
}

func TestTimeNow(t *testing.T) {
	// TODO: upstream doesn't tests this either
}

func TestTimeEndec(t *testing.T) {
	assert.Panics(t, func() {
		quick.Time(timeTag(""))
	})

	timeEndec(t, TimeTuple{}, "00000000")
	timeEndec(t, TimeTuple{Hour: 1}, "01000000")
	timeEndec(t, TimeTuple{Minute: 2}, "00020000")
	timeEndec(t, TimeTuple{Second: 3}, "00000300")
	timeEndec(t, TimeTuple{Hundredth: 4}, "00000004")
}
