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

	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/constructors"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type sampleBitString struct {
	*bacgopes.BitString
	bitLen   int
	bitNames map[string]int
}

func SampleBitString(args ...any) *sampleBitString { //nolint:all
	s := &sampleBitString{
		bitLen: 13,
		bitNames: map[string]int{
			"b0":  0,
			"b1":  1,
			"b4":  4,
			"b7":  7,
			"b8":  8,
			"b12": 12,
		},
	}
	var err error
	s.BitString, err = bacgopes.NewBitStringWithExtension(s, args)
	if err != nil {
		panic(err)
	}
	return s
}

func (s *sampleBitString) GetBitNames() map[string]int {
	return s.bitNames
}

func (s *sampleBitString) GetBitLen() int {
	return s.bitLen
}

func (s *sampleBitString) String() string {
	return s.BitString.String()
}

// Convert a hex string to a bit_string application tag.
func bitStringTag(t *testing.T, x string) bacgopes.Tag {
	t.Helper()
	b, err := bacgopes.Xtob(x)
	require.NoError(t, err)
	tag, err := bacgopes.NewTag(bacgopes.NewArgs(model.TagClass_APPLICATION_TAGS, bacgopes.TagBitStringAppTag, len(b), b))
	require.NoError(t, err)
	return tag
}

func bitStringEncode(obj *bacgopes.BitString) bacgopes.Tag {
	tag := Tag()
	obj.Encode(tag)
	return tag
}

func bitStringDecode(tag bacgopes.Tag) *bacgopes.BitString {
	obj := BitString(tag)
	return obj
}

// Pass the value to BitString, construct a tag from the hex string,
//
//	and compare results of encode and decoding each other.
func bitStringEndec(t *testing.T, v []int, x string) {
	t.Helper()
	tag := bitStringTag(t, x)

	obj := BitString(v)

	assert.Equal(t, tag, bitStringEncode(obj), "encoded obj should match the tag")
	assert.Equal(t, obj, bitStringDecode(tag), "decoded tag should match the obj")
}

func TestBitString(t *testing.T) {
	obj := BitString()
	assert.Len(t, obj.GetValue(), 0)
	assert.Equal(t, `BitString()`, obj.String())

	obj = BitString([]int{0})
	assert.Equal(t, []bool{false}, obj.GetValue())
	assert.Equal(t, `BitString(0)`, obj.String())

	obj = BitString([]int{0, 1})
	assert.Equal(t, []bool{false, true}, obj.GetValue())
	assert.Equal(t, `BitString(0,1)`, obj.String())

	assert.Panics(t, func() {
		BitString("some string")
	})

	assert.Panics(t, func() {
		BitString("1.0")
	})
}

func TestBitStringSample(t *testing.T) {
	obj := SampleBitString()
	assert.Equal(t, []bool(nil), obj.GetValue())

	obj = SampleBitString([]int{1})
	assert.Equal(t, `BitString(b0)`, obj.String())

	obj = SampleBitString([]string{"b4"})
	assert.Equal(t, `BitString(!b0,!b1,0,0,b4,0,0,!b7,!b8,0,0,0,!b12)`, obj.String())

	assert.Panics(t, func() {
		BitString("x1")
	})
}

func TestBitStringTag(t *testing.T) {
	tag := Tag(bacgopes.TagApplicationTagClass, bacgopes.TagBitStringAppTag, 1, xtob("08"))
	obj := BitString(tag)
	assert.Len(t, obj.GetValue(), 0)

	tag = Tag(bacgopes.TagApplicationTagClass, bacgopes.TagBitStringAppTag, 1, xtob("0102"))
	obj = BitString(tag)
	assert.Equal(t, []bool{false, false, false, false, false, false, true}, obj.GetValue())

	tag = Tag(bacgopes.TagApplicationTagClass, bacgopes.TagBitStringAppTag, 1, xtob(""))
	assert.Panics(t, func() {
		BitString(tag)
	})

	tag = Tag(bacgopes.TagContextTagClass, 0, 1, xtob("ff"))
	assert.Panics(t, func() {
		BitString(tag)
	})

	tag = Tag(bacgopes.TagApplicationTagClass, 0)
	assert.Panics(t, func() {
		BitString(tag)
	})
}

func TestBitStringCopy(t *testing.T) {
	sampleValue := []int{0, 1, 0, 1}
	obj1 := BitString(sampleValue)
	obj2 := BitString(obj1)
	assert.Equal(t, obj1, obj2)
}

func TestBitStringEndec(t *testing.T) {
	bitStringEndec(t, []int{}, "00")
	bitStringEndec(t, []int{0}, "0700")
	bitStringEndec(t, []int{1}, "0780")
	bitStringEndec(t, []int{0, 0}, "0600")
	bitStringEndec(t, []int{1, 1}, "06c0")
	bitStringEndec(t, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, "060000")
	bitStringEndec(t, []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, "06ffc0")
}
