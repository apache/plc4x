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

	"github.com/apache/plc4x/plc4go/internal/bacnetip"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/constructors"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

func TagTuple(tag bacnetip.Tag) (tagClass model.TagClass, tagNumber uint, tagLVT int, tagData []byte) {
	return tag.GetTagClass(), tag.GetTagNumber(), tag.GetTagLvt(), tag.GetTagData()
}

// Build PDU from the string, decode the tag, convert to an object.
func objDecode(blob []byte) any {
	data := PDUData(blob)
	tag := Tag(data)
	obj, err := tag.AppToObject()
	if err != nil {
		panic(err) // TODO
	}
	return obj
}

// Encode the object into a tag, encode it in a PDU, return the data.
func objEncode(obj any) []byte {
	tag := Tag()
	err := obj.(interface{ Encode(arg bacnetip.Arg) error }).Encode(tag)
	if err != nil {
		panic(err)
	}
	data := PDUData()
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
func contextDecode(blob []byte) *bacnetip.ContextTag {
	data := PDUData(blob)
	tag := ContextTag(data)
	return tag
}

// Encode the object into a tag, encode it in a PDU, return the data.
func contextEncode(tag *bacnetip.ContextTag) []byte {
	data := PDUData()
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
	tag1 := ContextTag(tnum, tdata)

	// decode the blob into a tag
	tag2 := contextDecode(blob1)

	// encode the tag into a blob
	blob2 := contextEncode(tag1)

	// compare the results
	assert.Equal(t, tag1, tag2)
	assert.Equal(t, blob1, blob2)
}

// Build PDU from the string, decode the tag, convert to an object.
func openingDecode(blob []byte) *bacnetip.OpeningTag {
	data := PDUData(blob)
	tag := OpeningTag(data)
	return tag
}

// Encode the object into a tag, encode it in a PDU, return the data.
func openingEncode(tag *bacnetip.OpeningTag) []byte {
	data := PDUData()
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
	tag1 := OpeningTag(tnum)

	// decode the blob into a tag
	tag2 := openingDecode(blob1)

	// encode the tag into a blob
	blob2 := openingEncode(tag1)

	// compare the results
	assert.Equal(t, tag1, tag2)
	assert.Equal(t, blob1, blob2)
}

// Build PDU from the string, decode the tag, convert to an object.
func closingDecode(blob []byte) *bacnetip.ClosingTag {
	data := PDUData(blob)
	tag := ClosingTag(data)
	return tag
}

// Encode the object into a tag, encode it in a PDU, return the data.
func closingEncode(tag *bacnetip.ClosingTag) []byte {
	data := PDUData()
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
	tag1 := ClosingTag(tnum)

	// decode the blob into a tag
	tag2 := closingDecode(blob1)

	// encode the tag into a blob
	blob2 := closingEncode(tag1)

	// compare the result
	assert.Equal(t, tag1, tag2)
	assert.Equal(t, blob1, blob2)
}

func TestTag(t *testing.T) {
	tag := Tag()
	assert.Equal(t, model.TagClass(0), tag.GetTagClass())
	assert.Equal(t, uint(0), tag.GetTagNumber())

	// must have a valid encoded tag to extract the from the data
	data := PDUData(xtob(""))
	assert.Panics(t, func() {
		Tag(data)
	})

	// must have two values, class and number
	assert.Panics(t, func() {
		Tag(0)
	})

	tag = Tag(0, 1)
	assert.Equal(t, model.TagClass(0), tag.GetTagClass())
	assert.Equal(t, uint(1), tag.GetTagNumber())

	tag = Tag(0, 1, 2)
	assert.Equal(t, model.TagClass(0), tag.GetTagClass())
	assert.Equal(t, uint(1), tag.GetTagNumber())
	assert.Equal(t, 2, tag.GetTagLvt())

	// tag data must be bytes of bytearray
	assert.Panics(t, func() {
		Tag(0, 1, 2, 3)
	})
}

func TestApplicationTag(t *testing.T) {
	tag := ApplicationTag(0, xtob(""))
	_ = tag

	assert.Panics(t, func() {
		ApplicationTag(0)
	})
}

func TestGenericApplicationToContext(t *testing.T) {
	// create and application
	tag := ApplicationTag(0, xtob("01"))

	// convert it to context tagged, context 0
	ctag, err := tag.AppToContext(0)
	require.NoError(t, err)

	// create a context tag with the same shape
	ttag := ContextTag(0, xtob("01"))

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
	tag := Tag(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_BOOLEAN, 0)

	// convert it to context tagged, context 0
	ctag, err := tag.AppToContext(0)
	require.NoError(t, err)

	// create a context tag with the same shape
	ttag := ContextTag(0, xtob("00"))

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
	objEndec(t, Null(), "00")

	// boolean
	objEndec(t, Boolean(true), "11")
	objEndec(t, Boolean(false), "10")

	// unsigned
	objEndec(t, Unsigned(0), "2100")
	objEndec(t, Unsigned(1), "2101")
	objEndec(t, Unsigned(127), "217F")
	objEndec(t, Unsigned(128), "2180")

	// integer
	objEndec(t, Integer(0), "3100")
	objEndec(t, Integer(1), "3101")
	objEndec(t, Integer(-1), "31FF")
	objEndec(t, Integer(128), "320080")
	objEndec(t, Integer(-128), "3180")

	// real
	objEndec(t, Real(0), "4400000000")
	objEndec(t, Real(1), "443F800000")
	objEndec(t, Real(-1), "44BF800000")
	objEndec(t, Real(73.5), "4442930000")

	// double
	objEndec(t, Double(0), "55080000000000000000")
	objEndec(t, Double(1), "55083FF0000000000000")
	objEndec(t, Double(-1), "5508BFF0000000000000")
	objEndec(t, Double(73.5), "55084052600000000000")

	// octet string
	objEndec(t, OctetString(xtob("")), "60")
	objEndec(t, OctetString(xtob("01")), "6101")
	objEndec(t, OctetString(xtob("0102")), "620102")
	objEndec(t, OctetString(xtob("010203040506")), "6506010203040506")

	// character string
	objEndec(t, CharacterString(""), "7100")
	objEndec(t, CharacterString("a"), "720061")
	objEndec(t, CharacterString("abcde"), "7506006162636465")

	// bit string
	objEndec(t, BitString([]bool{}), "8100")
	objEndec(t, BitString([]bool{false}), "820700")
	objEndec(t, BitString([]bool{true}), "820780")
	objEndec(t, BitString([]bool{true, true, true, true, true}), "8203F8")
	objEndec(t, BitString([]bool{true, true, true, true, true, true, true, true, true, true}), "8306FFC0")

	// enumerated
	objEndec(t, Enumerated(0), "9100")
	objEndec(t, Enumerated(1), "9101")
	objEndec(t, Enumerated(127), "917F")
	objEndec(t, Enumerated(128), "9180")

	// date
	objEndec(t, Date(bacnetip.DateTuple{1, 2, 3, 4}), "A401020304")
	objEndec(t, Date(bacnetip.DateTuple{255, 2, 3, 4}), "A4FF020304")
	objEndec(t, Date(bacnetip.DateTuple{1, 255, 3, 4}), "A401FF0304")
	objEndec(t, Date(bacnetip.DateTuple{1, 2, 255, 4}), "A40102FF04")
	objEndec(t, Date(bacnetip.DateTuple{1, 2, 3, 255}), "A4010203FF")

	// time
	objEndec(t, Time(bacnetip.TimeTuple{1, 2, 3, 4}), "B401020304")
	objEndec(t, Time(bacnetip.TimeTuple{255, 2, 3, 4}), "B4FF020304")
	objEndec(t, Time(bacnetip.TimeTuple{1, 255, 3, 4}), "B401FF0304")
	objEndec(t, Time(bacnetip.TimeTuple{1, 2, 255, 4}), "B40102FF04")
	objEndec(t, Time(bacnetip.TimeTuple{1, 2, 3, 255}), "B4010203FF")

	// object identifier
	objEndec(t, ObjectIdentifier(0, 0), "C400000000")
	objEndec(t, ObjectIdentifier(1, 0), "C400400000")
	objEndec(t, ObjectIdentifier(0, 2), "C400000002")
	objEndec(t, ObjectIdentifier(3, 4), "C400C00004")
}

func TestContextTag(t *testing.T) {
	ContextTag(0, xtob(""))

	contextEndec(t, 0, "", "08")
	contextEndec(t, 1, "01", "1901")
	contextEndec(t, 2, "0102", "2A0102")
	contextEndec(t, 3, "010203", "3B010203")
	contextEndec(t, 14, "010203", "EB010203")
	contextEndec(t, 15, "010203", "FB0F010203")
}

func TestOpeningTag(t *testing.T) {
	OpeningTag(0)

	openingEndec(t, 0, "0E")
	openingEndec(t, 1, "1E")
	openingEndec(t, 2, "2E")
	openingEndec(t, 3, "3E")
	openingEndec(t, 14, "EE")
	openingEndec(t, 15, "FE0F")
	openingEndec(t, 254, "FEFE")
}

func TestClosingTag(t *testing.T) {
	ClosingTag(0)

	closingEndec(t, 0, "0f")
	closingEndec(t, 1, "1F")
	closingEndec(t, 2, "2F")
	closingEndec(t, 3, "3F")
	closingEndec(t, 14, "EF")
	closingEndec(t, 15, "FF0F")
	closingEndec(t, 254, "FFFE")
}

func TestTagList(t *testing.T) {
	_ = TagList()
}

func TestPeek(t *testing.T) {
	tag0 := IntegerTag("00")
	taglist := TagList(tag0)

	// peek at the first tag
	assert.Equal(t, tag0, taglist.Peek())

	// pop of the front
	tag1 := taglist.Pop()
	var emptyList = make([]bacnetip.Tag, 0)
	assert.Equal(t, emptyList, taglist.GetTagList())

	// push if back to the front
	taglist.Push(tag1)
	assert.Equal(t, []bacnetip.Tag{tag1}, taglist.GetTagList())
}

func TestGetContext(t *testing.T) {
	tagListData := []bacnetip.Tag{
		ContextTag(0, xtob("00")),
		ContextTag(1, xtob("01")),
		OpeningTag(2),
		IntegerTag("03"),
		OpeningTag(0),
		IntegerTag("04"),
		ClosingTag(0),
		ClosingTag(2),
	}
	taglist := TagList(tagListData...)

	// known to be a simple context encoded element
	context0, err := taglist.GetContext(0)
	require.NoError(t, err)
	assert.Equal(t, tagListData[0], context0)

	// known to be a simple context encoded list of element(s)
	context2, err := taglist.GetContext(2)
	require.NoError(t, err)
	assert.Equal(t, tagListData[3:7], context2.(*bacnetip.TagList).GetTagList())

	// known missing context
	context3, err := taglist.GetContext(3)
	require.NoError(t, err)
	assert.Equal(t, nil, context3)
}

func TestEndec0(t *testing.T) { // Test bracketed application tagged integer encoding and decoding.
	tagList := TagList()

	data := PDUData()
	tagList.Encode(data)
	assert.Equal(t, []byte{}, data.GetPduData())

	tagList = TagList()
	err := tagList.Decode(data)
	assert.NoError(t, err)
	var noItems []bacnetip.Tag
	assert.Equal(t, noItems, tagList.GetTagList())
}

func TestEndec1(t *testing.T) { // Test bracketed application tagged integer encoding and decoding.
	tag0 := IntegerTag("00")
	tag1 := IntegerTag("01")
	tagList := TagList(tag0, tag1)

	data := PDUData()
	tagList.Encode(data)
	assert.Equal(t, xtob("31003101"), data.GetPduData())

	tagList = TagList()
	err := tagList.Decode(data)
	assert.NoError(t, err)
	assert.Equal(t, []bacnetip.Tag{tag0, tag1}, tagList.GetTagList())
}

func TestEndec2(t *testing.T) { // Test bracketed application tagged integer encoding and decoding.
	tag0 := ContextTag(0, xtob("00"))
	tag1 := ContextTag(1, xtob("01"))
	tagList := TagList(tag0, tag1)

	data := PDUData()
	tagList.Encode(data)
	assert.Equal(t, xtob("09001901"), data.GetPduData())

	tagList = TagList()
	err := tagList.Decode(data)
	assert.NoError(t, err)
	assert.Equal(t, []bacnetip.Tag{tag0, tag1}, tagList.GetTagList())
}

func TestEndec3(t *testing.T) { // Test bracketed application tagged integer encoding and decoding.
	tag0 := OpeningTag(0)
	tag1 := IntegerTag("0102")
	tag2 := ClosingTag(0)
	tagList := TagList(tag0, tag1, tag2)

	data := PDUData()
	tagList.Encode(data)
	assert.Equal(t, xtob("0E3201020F"), data.GetPduData())

	tagList = TagList()
	err := tagList.Decode(data)
	assert.NoError(t, err)
	assert.Equal(t, []bacnetip.Tag{tag0, tag1, tag2}, tagList.GetTagList())
}
