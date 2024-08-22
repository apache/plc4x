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
	"fmt"
	"strconv"
	"testing"

	"github.com/apache/plc4x/plc4go/internal/bacnetip"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"

	"github.com/stretchr/testify/assert"
)

type myObjectType struct {
	*bacnetip.ObjectType

	enumerations map[string]uint64
}

func (m *myObjectType) GetEnumerations() map[string]uint64 {
	return m.enumerations
}

func (m *myObjectType) SetObjectType(objectType *bacnetip.ObjectType) {
	m.ObjectType = objectType
}

func (m *myObjectType) String() string {
	value := strconv.Itoa(int(m.GetValue()))
	if m.GetValueString() != "" {
		value = m.GetValueString()
	}
	return fmt.Sprintf("MyObjectType(%v)", value)
}

func MyObjectType(args ...any) *myObjectType {
	o := &myObjectType{
		enumerations: map[string]uint64{
			"myAnalogInput":  128,
			"myAnalogOutput": 129,
			"myAnalogValue":  130,
		},
	}
	var err error
	o.ObjectType, err = bacnetip.NewObjectType(append([]any{o}, args...))
	if err != nil {
		panic(err)
	}
	return o
}

func ObjectType(args ...any) *bacnetip.ObjectType {
	if len(args) == 0 {
		ObjectType, err := bacnetip.NewObjectType(nil)
		if err != nil {
			panic(err)
		}
		return ObjectType
	}
	ObjectType, err := bacnetip.NewObjectType(args)
	if err != nil {
		panic(err)
	}
	return ObjectType
}

func ObjectTypeTag(x string) bacnetip.Tag {
	b := xtob(x)
	tag := Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_ENUMERATED, len(b), b)
	return tag
}

// Encode a ObjectType object into a tag.
func ObjectTypeEncode(obj *bacnetip.ObjectType) bacnetip.Tag {
	tag := Tag()
	obj.Encode(tag)
	return tag
}

// Decode a ObjectType application tag into a ObjectType.
func ObjectTypeDecode(tag bacnetip.Tag) *bacnetip.ObjectType {
	obj := ObjectType(tag)

	return obj
}

// Pass the value to ObjectType, construct a tag from the hex string,
//
//	and compare results of encode and decoding each other.
func ObjectTypeEndec(t *testing.T, v any, x string) {
	tag := ObjectTypeTag(x)

	obj := ObjectType(v)

	assert.Equal(t, tag, ObjectTypeEncode(obj))
	assert.Equal(t, obj, ObjectTypeDecode(tag))
}

func TestObjectType(t *testing.T) {
	obj := ObjectType()
	assert.Equal(t, uint64(0x0), obj.GetValue())

	assert.Panics(t, func() {
		ObjectType(1.0)
	})
}

func TestObjectTypeInt(t *testing.T) {
	obj := ObjectType(0)
	assert.Equal(t, uint64(0), obj.GetValue())
	assert.Equal(t, "ObjectType(analogInput)", obj.String())

	obj = ObjectType(127)
	assert.Equal(t, uint64(127), obj.GetValue())
	assert.Equal(t, "ObjectType(127)", obj.String())
}

func TestObjectTypeStr(t *testing.T) {
	obj := ObjectType("analogInput")
	assert.Equal(t, uint64(0), obj.GetValue())
	assert.Equal(t, "ObjectType(analogInput)", obj.String())
}

func TestExtendedObjectTypeInt(t *testing.T) {
	obj := MyObjectType(0)
	assert.Equal(t, uint64(0), obj.GetValue())
	assert.Equal(t, "MyObjectType(analogInput)", obj.String())

	obj = MyObjectType(128)
	assert.Equal(t, uint64(128), obj.GetValue())
	assert.Equal(t, "MyObjectType(myAnalogInput)", obj.String())
}

func TestExtendedObjectTypeStr(t *testing.T) {
	obj := MyObjectType("myAnalogInput")
	assert.Equal(t, "myAnalogInput", obj.GetValueString())
	assert.Equal(t, "MyObjectType(myAnalogInput)", obj.String())
}

func TestObjectTypeTag(t *testing.T) {
	tag := Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_ENUMERATED, 1, xtob("01"))
	obj := ObjectType(tag)
	assert.Equal(t, "analogOutput", obj.GetValueString())

	tag = Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_BOOLEAN, 0, xtob(""))
	assert.Panics(t, func() {
		ObjectType(tag)
	})

	tag = Tag(model.TagClass_CONTEXT_SPECIFIC_TAGS, 0, 1, xtob("ff"))
	assert.Panics(t, func() {
		ObjectType(tag)
	})

	tag = Tag(bacnetip.TagOpeningTagClass, 0)
	assert.Panics(t, func() {
		ObjectType(tag)
	})
}

func TestObjectTypeCopy(t *testing.T) {
	obj1 := ObjectType(12)
	obj2 := ObjectType(obj1)
	assert.Equal(t, "loop", obj2.GetValueString())
	assert.True(t, obj1.Equals(obj2))
}

func TestObjectTypeEndec(t *testing.T) {
	assert.Panics(t, func() {
		ObjectType(ObjectTypeTag(""))
	})

	ObjectTypeEndec(t, "analogInput", "00")
	ObjectTypeEndec(t, "analogOutput", "01")

	ObjectTypeEndec(t, 127, "7f")
	ObjectTypeEndec(t, 128, "80")
}
