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

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/quick"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

func TagTuple(tag Tag) (tagClass model.TagClass, tagNumber uint, tagLVT int, tagData []byte) {
	return tag.GetTagClass(), tag.GetTagNumber(), tag.GetTagLvt(), tag.GetTagData()
}

// Build PDU from the string, decode the tag, convert to an object.
func objDecode(blob []byte) any {
	data := quick.PDUData(blob)
	tag := quick.Tag(data)
	obj, err := tag.AppToObject()
	if err != nil {
		panic(err) // TODO
	}
	return obj
}

// Encode the object into a tag, encode it in a PDU, return the data.
func objEncode(obj any) []byte {
	tag := quick.Tag()
	err := obj.(Encoder).Encode(tag)
	if err != nil {
		panic(err)
	}
	data := quick.PDUData()
	tag.Encode(data)
	return data.GetPduData()
}

// Convert the value (a primitive object) to a hex encoded string,
//
//	convert the hex encoded string to and object, and compare the results to
//	each other.
func objEndec(t *testing.T, obj any, x string) {
	t.Helper()
	// convert the hex strings to a blobs
	blob := xtob(x)

	// decode the blob into an object
	obj2 := objDecode(blob)

	// encode the tag into a blob
	blob2 := objEncode(obj)

	// compare the results
	assert.Equal(t, obj, obj2)
	assert.Equal(t, blob, blob2)
}

// Build PDU from the string, decode the tag, convert to an object.
func contextDecode(blob []byte) *ContextTag {
	data := quick.PDUData(blob)
	tag := quick.ContextTag(data)
	return tag
}

// Encode the object into a tag, encode it in a PDU, return the data.
func contextEncode(tag *ContextTag) []byte {
	data := quick.PDUData()
	tag.Encode(data)
	return data.GetPduData()
}

// Convert the value (a primitive object) to a hex encoded string,
// convert the hex encoded string to and object, and compare the results to
// each other.
func contextEndec(t *testing.T, tnum int, x string, y string) {
	t.Helper()
	// convert the hex strings to a blobs
	tdata := xtob(x)
	blob1 := xtob(y)

	// make a context tag
	tag1 := quick.ContextTag(tnum, tdata)

	// decode the blob into a tag
	tag2 := contextDecode(blob1)

	// encode the tag into a blob
	blob2 := contextEncode(tag1)

	// compare the results
	assert.Equal(t, tag1, tag2)
	assert.Equal(t, blob1, blob2)
}

// Build PDU from the string, decode the tag, convert to an object.
func openingDecode(blob []byte) *OpeningTag {
	data := quick.PDUData(blob)
	tag := quick.OpeningTag(data)
	return tag
}

// Encode the object into a tag, encode it in a PDU, return the data.
func openingEncode(tag *OpeningTag) []byte {
	data := quick.PDUData()
	tag.Encode(data)
	return data.GetPduData()
}

// Convert the value (a primitive object) to a hex encoded string,
// convert the hex encoded string to and object, and compare the results to
// each other.
func openingEndec(t *testing.T, tnum int, x string) {

	// convert the hex string to a blob
	blob1 := xtob(x)

	// make a context tag
	tag1 := quick.OpeningTag(tnum)

	// decode the blob into a tag
	tag2 := openingDecode(blob1)

	// encode the tag into a blob
	blob2 := openingEncode(tag1)

	// compare the results
	assert.Equal(t, tag1, tag2)
	assert.Equal(t, blob1, blob2)
}

// Build PDU from the string, decode the tag, convert to an object.
func closingDecode(blob []byte) *ClosingTag {
	data := quick.PDUData(blob)
	tag := quick.ClosingTag(data)
	return tag
}

// Encode the object into a tag, encode it in a PDU, return the data.
func closingEncode(tag *ClosingTag) []byte {
	data := quick.PDUData()
	tag.Encode(data)
	return data.GetPduData()
}

// Convert the value (a primitive object) to a hex encoded string,
// convert the hex encoded string to and object, and compare the results to
// each other.
func closingEndec(t *testing.T, tnum int, x string) {
	t.Helper()
	// convert the hex string to a blob
	blob1 := xtob(x)

	// make a context tag
	tag1 := quick.ClosingTag(tnum)

	// decode the blob into a tag
	tag2 := closingDecode(blob1)

	// encode the tag into a blob
	blob2 := closingEncode(tag1)

	// compare the result
	assert.Equal(t, tag1, tag2)
	assert.Equal(t, blob1, blob2)
}

func TestTag(t *testing.T) {
	tag := quick.Tag()
	assert.Equal(t, model.TagClass(0), tag.GetTagClass())
	assert.Equal(t, uint(0), tag.GetTagNumber())

	// must have a valid encoded tag to extract the from the data
	data := quick.PDUData(xtob(""))
	assert.Panics(t, func() {
		quick.Tag(data)
	})

	// must have two values, class and number
	assert.Panics(t, func() {
		quick.Tag(0)
	})

	tag = quick.Tag(0, 1)
	assert.Equal(t, model.TagClass(0), tag.GetTagClass())
	assert.Equal(t, uint(1), tag.GetTagNumber())

	tag = quick.Tag(0, 1, 2)
	assert.Equal(t, model.TagClass(0), tag.GetTagClass())
	assert.Equal(t, uint(1), tag.GetTagNumber())
	assert.Equal(t, 2, tag.GetTagLvt())

	// tag data must be bytes of bytearray
	assert.Panics(t, func() {
		quick.Tag(0, 1, 2, 3)
	})
}

func TestApplicationTag(t *testing.T) {
	tag := quick.ApplicationTag(0, xtob(""))
	_ = tag

	assert.Panics(t, func() {
		quick.ApplicationTag(0)
	})
}

func TestGenericApplicationToContext(t *testing.T) {
	// create and application
	tag := quick.ApplicationTag(0, xtob("01"))

	// convert it to context tagged, context 0
	ctag, err := tag.AppToContext(0)
	require.NoError(t, err)

	// create a context tag with the same shape
	ttag := quick.ContextTag(0, xtob("01"))

	// check to see if they are the same
	assert.Equal(t, ttag, ctag)

	// convert the context tag back to an application tag
	dtag, err := ctag.ContextToApp(0)
	require.NoError(t, err)

	// check to see it round-tripped
	assert.Equal(t, tag, dtag)
}

func TestBooleanApplicationToContext(t *testing.T) {
	// create and application
	tag := quick.Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_BOOLEAN, 0)

	// convert it to context tagged, context 0
	ctag, err := tag.AppToContext(0)
	require.NoError(t, err)

	// create a context tag with the same shape
	ttag := quick.ContextTag(0, xtob("00"))

	// check to see if they are the same
	assert.Equal(t, ttag, ctag)

	// convert the context tag back to an application tag
	dtag, err := ctag.ContextToApp(uint(model.BACnetDataType_BOOLEAN))
	require.NoError(t, err)

	// check to see it round-tripped
	assert.Equal(t, tag, dtag)
}

func TestBooleanApplicationToObject(t *testing.T) {
	// null
	objEndec(t, quick.Null(), "00")

	// boolean
	objEndec(t, quick.Boolean(true), "11")
	objEndec(t, quick.Boolean(false), "10")

	// unsigned
	objEndec(t, quick.Unsigned(0), "2100")
	objEndec(t, quick.Unsigned(1), "2101")
	objEndec(t, quick.Unsigned(127), "217F")
	objEndec(t, quick.Unsigned(128), "2180")

	// integer
	objEndec(t, quick.Integer(0), "3100")
	objEndec(t, quick.Integer(1), "3101")
	objEndec(t, quick.Integer(-1), "31FF")
	objEndec(t, quick.Integer(128), "320080")
	objEndec(t, quick.Integer(-128), "3180")

	// real
	objEndec(t, quick.Real(0), "4400000000")
	objEndec(t, quick.Real(1), "443F800000")
	objEndec(t, quick.Real(-1), "44BF800000")
	objEndec(t, quick.Real(73.5), "4442930000")

	// double
	objEndec(t, quick.Double(0), "55080000000000000000")
	objEndec(t, quick.Double(1), "55083FF0000000000000")
	objEndec(t, quick.Double(-1), "5508BFF0000000000000")
	objEndec(t, quick.Double(73.5), "55084052600000000000")

	// octet string
	objEndec(t, quick.OctetString(xtob("")), "60")
	objEndec(t, quick.OctetString(xtob("01")), "6101")
	objEndec(t, quick.OctetString(xtob("0102")), "620102")
	objEndec(t, quick.OctetString(xtob("010203040506")), "6506010203040506")

	// character string
	objEndec(t, quick.CharacterString(""), "7100")
	objEndec(t, quick.CharacterString("a"), "720061")
	objEndec(t, quick.CharacterString("abcde"), "7506006162636465")

	// bit string
	objEndec(t, quick.BitString([]bool{}), "8100")
	objEndec(t, quick.BitString([]bool{false}), "820700")
	objEndec(t, quick.BitString([]bool{true}), "820780")
	objEndec(t, quick.BitString([]bool{true, true, true, true, true}), "8203F8")
	objEndec(t, quick.BitString([]bool{true, true, true, true, true, true, true, true, true, true}), "8306FFC0")

	// enumerated
	objEndec(t, quick.Enumerated(0), "9100")
	objEndec(t, quick.Enumerated(1), "9101")
	objEndec(t, quick.Enumerated(127), "917F")
	objEndec(t, quick.Enumerated(128), "9180")

	// date
	objEndec(t, quick.Date(DateTuple{Year: 1, Month: 2, Day: 3, DayOfWeek: 4}), "A401020304")
	objEndec(t, quick.Date(DateTuple{Year: 255, Month: 2, Day: 3, DayOfWeek: 4}), "A4FF020304")
	objEndec(t, quick.Date(DateTuple{Year: 1, Month: 255, Day: 3, DayOfWeek: 4}), "A401FF0304")
	objEndec(t, quick.Date(DateTuple{Year: 1, Month: 2, Day: 255, DayOfWeek: 4}), "A40102FF04")
	objEndec(t, quick.Date(DateTuple{Year: 1, Month: 2, Day: 3, DayOfWeek: 255}), "A4010203FF")

	// time
	objEndec(t, quick.Time(TimeTuple{Hour: 1, Minute: 2, Second: 3, Hundredth: 4}), "B401020304")
	objEndec(t, quick.Time(TimeTuple{Hour: 255, Minute: 2, Second: 3, Hundredth: 4}), "B4FF020304")
	objEndec(t, quick.Time(TimeTuple{Hour: 1, Minute: 255, Second: 3, Hundredth: 4}), "B401FF0304")
	objEndec(t, quick.Time(TimeTuple{Hour: 1, Minute: 2, Second: 255, Hundredth: 4}), "B40102FF04")
	objEndec(t, quick.Time(TimeTuple{Hour: 1, Minute: 2, Second: 3, Hundredth: 255}), "B4010203FF")

	// object identifier
	objEndec(t, quick.ObjectIdentifier(0, 0), "C400000000")
	objEndec(t, quick.ObjectIdentifier(1, 0), "C400400000")
	objEndec(t, quick.ObjectIdentifier(0, 2), "C400000002")
	objEndec(t, quick.ObjectIdentifier(3, 4), "C400C00004")
}

func TestContextTag(t *testing.T) {
	quick.ContextTag(0, xtob(""))

	contextEndec(t, 0, "", "08")
	contextEndec(t, 1, "01", "1901")
	contextEndec(t, 2, "0102", "2A0102")
	contextEndec(t, 3, "010203", "3B010203")
	contextEndec(t, 14, "010203", "EB010203")
	contextEndec(t, 15, "010203", "FB0F010203")
}

func TestOpeningTag(t *testing.T) {
	quick.OpeningTag(0)

	openingEndec(t, 0, "0E")
	openingEndec(t, 1, "1E")
	openingEndec(t, 2, "2E")
	openingEndec(t, 3, "3E")
	openingEndec(t, 14, "EE")
	openingEndec(t, 15, "FE0F")
	openingEndec(t, 254, "FEFE")
}

func TestClosingTag(t *testing.T) {
	quick.ClosingTag(0)

	closingEndec(t, 0, "0f")
	closingEndec(t, 1, "1F")
	closingEndec(t, 2, "2F")
	closingEndec(t, 3, "3F")
	closingEndec(t, 14, "EF")
	closingEndec(t, 15, "FF0F")
	closingEndec(t, 254, "FFFE")
}

func TestTagList(t *testing.T) {
	_ = quick.TagList()
}

func TestPeek(t *testing.T) {
	tag0 := IntegerTag("00")
	taglist := quick.TagList(tag0)

	// peek at the first tag
	assert.Equal(t, tag0, taglist.Peek())

	// pop of the front
	tag1 := taglist.Pop()
	var emptyList = make([]Tag, 0)
	assert.Equal(t, emptyList, taglist.GetTagList())

	// push if back to the front
	taglist.Push(tag1)
	assert.Equal(t, []Tag{tag1}, taglist.GetTagList())
}

func TestGetContext(t *testing.T) {
	tagListData := []Tag{
		quick.ContextTag(0, xtob("00")),
		quick.ContextTag(1, xtob("01")),
		quick.OpeningTag(2),
		IntegerTag("03"),
		quick.OpeningTag(0),
		IntegerTag("04"),
		quick.ClosingTag(0),
		quick.ClosingTag(2),
	}
	taglist := quick.TagList(tagListData...)

	// known to be a simple context encoded element
	context0, err := taglist.GetContext(0)
	require.NoError(t, err)
	assert.Equal(t, tagListData[0], context0)

	// known to be a simple context encoded list of element(s)
	context2, err := taglist.GetContext(2)
	require.NoError(t, err)
	assert.Equal(t, tagListData[3:7], context2.(*TagList).GetTagList())

	// known missing context
	context3, err := taglist.GetContext(3)
	require.NoError(t, err)
	assert.Equal(t, nil, context3)
}

func TestEndec0(t *testing.T) { // Test bracketed application tagged integer encoding and decoding.
	tagList := quick.TagList()

	data := quick.PDUData()
	tagList.Encode(data)
	assert.Nil(t, data.GetPduData())

	tagList = quick.TagList()
	err := tagList.Decode(data)
	assert.NoError(t, err)
	var noItems []Tag
	assert.Equal(t, noItems, tagList.GetTagList())
}

func TestEndec1(t *testing.T) { // Test bracketed application tagged integer encoding and decoding.
	tag0 := IntegerTag("00")
	tag1 := IntegerTag("01")
	tagList := quick.TagList(tag0, tag1)

	data := quick.PDUData()
	tagList.Encode(data)
	assert.Equal(t, xtob("31003101"), data.GetPduData())

	tagList = quick.TagList()
	err := tagList.Decode(data)
	assert.NoError(t, err)
	assert.Equal(t, []Tag{tag0, tag1}, tagList.GetTagList())
}

func TestEndec2(t *testing.T) { // Test bracketed application tagged integer encoding and decoding.
	tag0 := quick.ContextTag(0, xtob("00"))
	tag1 := quick.ContextTag(1, xtob("01"))
	tagList := quick.TagList(tag0, tag1)

	data := quick.PDUData()
	tagList.Encode(data)
	assert.Equal(t, xtob("09001901"), data.GetPduData())

	tagList = quick.TagList()
	err := tagList.Decode(data)
	assert.NoError(t, err)
	assert.Equal(t, []Tag{tag0, tag1}, tagList.GetTagList())
}

func TestEndec3(t *testing.T) { // Test bracketed application tagged integer encoding and decoding.
	tag0 := quick.OpeningTag(0)
	tag1 := IntegerTag("0102")
	tag2 := quick.ClosingTag(0)
	tagList := quick.TagList(tag0, tag1, tag2)

	data := quick.PDUData()
	tagList.Encode(data)
	assert.Equal(t, xtob("0E3201020F"), data.GetPduData())

	tagList = quick.TagList()
	err := tagList.Decode(data)
	assert.NoError(t, err)
	assert.Equal(t, []Tag{tag0, tag1, tag2}, tagList.GetTagList())
}
