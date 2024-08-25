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

// Convert a hex string to a character_string application tag.
func timeTag(x string) bacnetip.Tag {
	b := xtob(x)
	tag := Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_TIME, len(b), b)
	return tag
}

// Encode a Time object into a tag.
func timeEncode(obj *bacnetip.Time) bacnetip.Tag {
	tag := Tag()
	obj.Encode(tag)
	return tag
}

// Decode a Time application tag into a Time.
func timeDecode(tag bacnetip.Tag) *bacnetip.Time {
	obj := Time(tag)

	return obj
}

// Pass the value to Time, construct a tag from the hex string,
//
//	and compare results of encode and decoding each other.
func timeEndec(t *testing.T, v bacnetip.TimeTuple, x string) {
	tag := timeTag(x)

	obj := Time(v)

	assert.Equal(t, tag, timeEncode(obj))
	assert.Equal(t, obj, timeDecode(tag))
}

func TestTime(t *testing.T) {
	obj := Time()
	assert.Equal(t, bacnetip.TimeTuple{Hour: 0xff, Minute: 0xff, Second: 0xff, Hundredth: 0xff}, obj.GetValue())

	assert.Panics(t, func() {
		Time("some string")
	})
}

func TestTimeTuple(t *testing.T) {
	obj := Time(bacnetip.TimeTuple{Hour: 1, Minute: 2, Second: 3, Hundredth: 4})
	assert.Equal(t, bacnetip.TimeTuple{Hour: 1, Minute: 2, Second: 3, Hundredth: 4}, obj.GetValue())
	assert.Equal(t, "Time(01:02:03.04)", obj.String())

	assert.Equal(t, bacnetip.TimeTuple{Hour: 1, Minute: 2, Second: 0, Hundredth: 0}, Time("1:2").GetValue())
	assert.Equal(t, bacnetip.TimeTuple{Hour: 1, Minute: 2, Second: 3, Hundredth: 0}, Time("1:2:3").GetValue())
	assert.Equal(t, bacnetip.TimeTuple{Hour: 1, Minute: 2, Second: 3, Hundredth: 40}, Time("1:2:3.4").GetValue())
	assert.Equal(t, bacnetip.TimeTuple{Hour: 1, Minute: 255, Second: 255, Hundredth: 255}, Time("1:*").GetValue())
	assert.Equal(t, bacnetip.TimeTuple{Hour: 1, Minute: 2, Second: 255, Hundredth: 255}, Time("1:2:*").GetValue())
	assert.Equal(t, bacnetip.TimeTuple{Hour: 1, Minute: 2, Second: 3, Hundredth: 255}, Time("1:2:3.*").GetValue())
}

func TestTimeTag(t *testing.T) {
	tag := Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_TIME, 4, xtob("01020304"))
	obj := Time(tag)
	assert.Equal(t, bacnetip.TimeTuple{Hour: 1, Minute: 2, Second: 3, Hundredth: 4}, obj.GetValue())

	tag = Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_BOOLEAN, 0, xtob(""))
	assert.Panics(t, func() {
		Time(tag)
	})

	tag = Tag(model.TagClass_CONTEXT_SPECIFIC_TAGS, 0, 1, xtob("ff"))
	assert.Panics(t, func() {
		Time(tag)
	})

	tag = Tag(bacnetip.TagOpeningTagClass, 0)
	assert.Panics(t, func() {
		Time(tag)
	})
}

func TestTimeCopy(t *testing.T) {
	obj1 := Time(bacnetip.TimeTuple{Hour: 1, Minute: 2, Second: 3, Hundredth: 4})
	obj2 := Time(obj1)
	assert.Equal(t, obj1, obj2)
}

func TestTimeNow(t *testing.T) {
	// TODO: upstream doesn't tests this either
}

func TestTimeEndec(t *testing.T) {
	assert.Panics(t, func() {
		Time(timeTag(""))
	})

	timeEndec(t, bacnetip.TimeTuple{0, 0, 0, 0}, "00000000")
	timeEndec(t, bacnetip.TimeTuple{1, 0, 0, 0}, "01000000")
	timeEndec(t, bacnetip.TimeTuple{0, 2, 0, 0}, "00020000")
	timeEndec(t, bacnetip.TimeTuple{0, 0, 3, 0}, "00000300")
	timeEndec(t, bacnetip.TimeTuple{0, 0, 0, 4}, "00000004")
}
