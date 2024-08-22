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

package bacnetip

import (
	"bytes"
	"cmp"
	"encoding/binary"
	"fmt"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"

	"github.com/pkg/errors"
)

const (
	// Deprecated: use model.TagClass_APPLICATION_TAGS
	TagApplicationTagClass = model.TagClass_APPLICATION_TAGS
	// Deprecated: use model.TagClass_CONTEXT_SPECIFIC_TAGS
	TagContextTagClass = model.TagClass_CONTEXT_SPECIFIC_TAGS
	TagOpeningTagClass = 2
	TagClosingTagClass = 3

	// Deprecated: use  model.BACnetDataType_NULL
	TagNullAppTag = model.BACnetDataType_NULL
	// Deprecated: use  model.BACnetDataType_BOOLEAN
	TagBooleanAppTag = model.BACnetDataType_BOOLEAN
	// Deprecated: use  model.BACnetDataType_UNSIGNED_INTEGER
	TagUnsignedAppTag = model.BACnetDataType_UNSIGNED_INTEGER
	// Deprecated: use  model.BACnetDataType_SIGNED_INTEGER
	TagIntegerAppTag = model.BACnetDataType_SIGNED_INTEGER
	// Deprecated: use  model.BACnetDataType_REAL
	TagRealAppTag = model.BACnetDataType_REAL
	// Deprecated: use  model.BACnetDataType_DOUBLE
	TagDoubleAppTag = model.BACnetDataType_DOUBLE
	// Deprecated: use  model.BACnetDataType_OCTET_STRING
	TagOctetStringAppTag = model.BACnetDataType_OCTET_STRING
	// Deprecated: use  model.BACnetDataType_CHARACTER_STRING
	TagCharacterStringAppTag = model.BACnetDataType_CHARACTER_STRING
	// Deprecated: use  model.BACnetDataType_BIT_STRING
	TagBitStringAppTag = model.BACnetDataType_BIT_STRING
	// Deprecated: use  model.BACnetDataType_ENUMERATED
	TagEnumeratedAppTag = model.BACnetDataType_ENUMERATED
	// Deprecated: use  model.BACnetDataType_DATE
	TagDateAppTag = model.BACnetDataType_DATE
	// Deprecated: use  model.BACnetDataType_TIME
	TagTimeAppTag = model.BACnetDataType_TIME
	// Deprecated: use  model.BACnetDataType_BACNET_OBJECT_IDENTIFIER
	TagObjectIdentifierAppTag = model.BACnetDataType_BACNET_OBJECT_IDENTIFIER
	TagReservedAppTag13       = 13
	TagReservedAppTag14       = 14
	TagReservedAppTag15       = 15
)

type Tag interface {
	GetTagClass() model.TagClass
	GetTagNumber() uint
	GetTagLvt() int
	GetTagData() []byte
	Encode(pdu PDUData)
	Decode(pdu PDUData) error
	AppToObject() (any, error)
	AppToContext(context uint) (*ContextTag, error)
	setAppData(tagNumber uint, tdata []byte)
	set(args Args)
}

type tag struct {
	tagClass  model.TagClass
	tagNumber uint
	tagLVT    int
	tagData   []byte

	appTagName  []string
	appTagClass []any
}

func NewTag(args Args) (Tag, error) {
	t := &tag{
		appTagName: []string{
			"null", "boolean", "unsigned", "integer",
			"real", "double", "octetString", "characterString",
			"bitString", "enumerated", "date", "time",
			"objectIdentifier", "reserved13", "reserved14", "reserved15",
		},
		appTagClass: []any{
			&Null{}, &Boolean{}, &Unsigned{}, &Integer{},
			&Real{}, &Double{}, &OctetString{}, &CharacterString{},
			&BitString{}, &Enumerated{}, &Date{}, &Time{},
			&ObjectIdentifier{}, nil, nil, nil,
		},
	}
	if len(args) == 0 {
		return t, nil
	}
	if len(args) == 1 {
		if err := t.Decode(args[0].(PDUData)); err != nil {
			return nil, errors.New("error decoding")
		}
	} else if len(args) >= 2 {
		t.set(args)
	} else {
		return nil, errors.Errorf("invalid arguments %v", args)
	}
	return t, nil
}

func (t *tag) Decode(pdu PDUData) error {
	tag, err := pdu.Get()
	if err != nil {
		return errors.Wrap(err, "error decoding tag")
	}

	// extract the type
	t.tagClass = model.TagClass(tag >> 3 & 0x01)

	// extract the tag number
	t.tagNumber = uint(tag >> 4)
	if t.tagNumber == 0x0f {
		get, err := pdu.Get()
		if err != nil {
			return errors.Wrap(err, "error decoding get")
		}
		t.tagNumber = uint(get)
	}

	// extract the length
	t.tagLVT = int(tag & 0x07)
	if t.tagLVT == 5 {
		get, err := pdu.Get()
		if err != nil {
			return errors.Wrap(err, "error decoding get")
		}
		t.tagLVT = int(get)
		if t.tagLVT == 254 {
			get, err := pdu.GetShort()
			if err != nil {
				return errors.Wrap(err, "error decoding get")
			}
			t.tagLVT = int(get)
		} else if t.tagLVT == 255 {
			get, err := pdu.GetLong()
			if err != nil {
				return errors.Wrap(err, "error decoding get")
			}
			t.tagLVT = int(get)
		}
	} else if t.tagLVT == 6 {
		t.tagClass = TagOpeningTagClass
		t.tagLVT = 0
	} else if t.tagLVT == 7 {
		t.tagClass = TagClosingTagClass
		t.tagLVT = 0
	}

	if t.tagClass == model.TagClass_APPLICATION_TAGS && t.tagNumber == uint(model.BACnetDataType_BOOLEAN) {
		// tagLVT contains value
		t.tagData = nil
	} else {
		// tagLVT contains length
		t.tagData, err = pdu.GetData(t.tagLVT)
		if err != nil {
			return errors.Wrap(err, "error decoding tag data")
		}
	}
	return nil
}

func (t *tag) Encode(pdu PDUData) {
	var data byte
	// check for special encoding
	if t.tagClass == model.TagClass_CONTEXT_SPECIFIC_TAGS {
		data = 0x08
	} else if t.tagClass == TagOpeningTagClass {
		data = 0x0E
	} else if t.tagClass == TagClosingTagClass {
		data = 0x0F
	} else {
		data = 0x00
	}

	// encode the tag number part
	if t.tagNumber < 15 {
		data += byte(t.tagNumber) << 4
	} else {
		data += 0xF0
	}

	// encode the length/value/type part
	if t.tagLVT < 5 {
		data += byte(t.tagLVT)
	} else {
		data += 0x05
	}

	// save this and the extended tag value
	pdu.Put(data)
	if t.tagNumber >= 15 {
		pdu.Put(byte(t.tagNumber))
	}

	// really short lengths are already done
	if t.tagLVT >= 5 {
		if t.tagLVT <= 253 {
			pdu.Put(byte(t.tagLVT))
		} else if t.tagLVT <= 65535 {
			pdu.Put(254)
			pdu.PutShort(uint16(t.tagLVT))
		} else {
			pdu.Put(255)
			pdu.PutLong(uint32(t.tagLVT))
		}
	}

	// now put the data
	pdu.PutData(t.tagData...)
}

func (t *tag) AppToContext(context uint) (*ContextTag, error) {
	if t.tagClass != model.TagClass_APPLICATION_TAGS {
		return nil, errors.New("application tag required")
	}
	if t.tagNumber == uint(model.BACnetDataType_BOOLEAN) {
		return NewContextTag(NewArgs(context, []byte{byte(t.tagLVT)}))
	}
	return NewContextTag(NewArgs(context, t.tagData))
}

func (t *tag) ContextToApp(dataType uint) (any, error) {
	if t.tagClass != model.TagClass_CONTEXT_SPECIFIC_TAGS {
		return nil, errors.New("context tag required")
	}
	if dataType == uint(model.BACnetDataType_BOOLEAN) {
		return NewTag(NewArgs(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_BOOLEAN, t.tagData[0], nil))
	}
	return NewApplicationTag(NewArgs(dataType, t.tagData))
}

func (t *tag) AppToObject() (any, error) {
	if t.tagClass != model.TagClass_APPLICATION_TAGS {
		return nil, errors.New("context tag required")
	}

	klass := t.appTagClass[int(t.tagNumber)]
	if klass == nil {
		return nil, nil
	}

	switch klass.(type) {
	case *Null:
		return NewNull(t)
	case *Boolean:
		return NewBoolean(t)
	case *Unsigned:
		return NewUnsigned(t)
	case *Integer:
		return NewInteger(t)
	case *Real:
		return NewReal(t)
	case *Double:
		return NewDouble(t)
	case *OctetString:
		return NewOctetString(t)
	case *CharacterString:
		return NewCharacterString(t)
	case *BitString:
		return NewBitString(NewArgs(t))
	case *Enumerated:
		return NewEnumerated(t)
	case *Date:
		return NewDate(t, NoArgs)
	case *Time:
		return NewTime(t, NoArgs)
	case *ObjectIdentifier:
		return NewObjectIdentifier(NewArgs(t))
	default:
		return nil, errors.Errorf("unknown tag klass %T", klass)
	}
}

func (t *tag) set(args Args) {
	switch tagClass := args[0].(type) {
	case model.TagClass:
		t.tagClass = tagClass
	case uint:
		t.tagClass = model.TagClass(tagClass)
	case uint8:
		t.tagClass = model.TagClass(tagClass)
	case int:
		t.tagClass = model.TagClass(tagClass)
	default:
		panic("oh no")
	}
	switch tagNumber := args[1].(type) {
	case model.BACnetDataType:
		t.tagNumber = uint(tagNumber)
	case uint:
		t.tagNumber = tagNumber
	case int:
		t.tagNumber = uint(tagNumber)
	default:
		panic("oh no")
	}
	if len(args) == 2 {
		return
	}
	switch tagLVT := args[2].(type) {
	case uint:
		t.tagLVT = int(tagLVT)
	case uint8:
		t.tagLVT = int(tagLVT)
	case int:
		t.tagLVT = tagLVT
	default:
		panic("oh no")
	}
	if len(args) == 3 {
		return
	}
	switch tagData := args[3].(type) {
	case []byte:
		t.tagData = tagData
	case nil:
		// ok
	default:
		panic("oh no")
	}
	if len(args) > 4 {
		panic("oh no")
	}
}

func (t *tag) setAppData(tagNumber uint, tdata []byte) {
	t.tagClass = model.TagClass_APPLICATION_TAGS
	t.tagNumber = tagNumber
	t.tagLVT = len(tdata)
	t.tagData = tdata
}

func (t *tag) GetTagClass() model.TagClass {
	return t.tagClass
}

func (t *tag) GetTagNumber() uint {
	return t.tagNumber
}

func (t *tag) GetTagLvt() int {
	return t.tagLVT
}

func (t *tag) GetTagData() []byte {
	return t.tagData
}

func (t *tag) Equals(other any) bool {
	if t == nil && other == nil {
		return true
	}
	if other == nil {
		return false
	}
	otherTag, ok := other.(Tag)
	if !ok {
		return false
	}
	return t.tagClass == otherTag.GetTagClass() &&
		t.tagNumber == otherTag.GetTagNumber() &&
		t.tagLVT == otherTag.GetTagLvt() &&
		bytes.Equal(t.tagData, otherTag.GetTagData())
}

type ApplicationTag struct {
	*tag
}

func NewApplicationTag(args Args) (*ApplicationTag, error) {
	a := &ApplicationTag{}
	switch len(args) {
	case 1:
		if _, ok := args[0].(PDUData); !ok {
			return nil, errors.Errorf("invalid argument %T", args[0])
		}
		_tag, err := NewTag(args)
		if err != nil {
			return nil, errors.New("error creating tag")
		}
		a.tag = _tag.(*tag)
		if a.tagClass != model.TagClass_APPLICATION_TAGS {
			return nil, errors.New("error creating tag: invalid tag class")
		}
		return a, nil
	case 2:
		var tnum any
		tnum, ok := args[0].(uint)
		if !ok {
			tnum, ok = args[0].(int)
			if !ok {
				return nil, errors.New("error creating tag: invalid tag number")
			}
		}
		tdata := args[1].([]byte)
		_tag, err := NewTag(NewArgs(model.TagClass_APPLICATION_TAGS, tnum, len(tdata), tdata))
		if err != nil {
			return nil, errors.New("error creating tag")
		}
		a.tag = _tag.(*tag)
		return a, nil
	default:
		return nil, errors.New("requires type and data or pdu data")
	}
}

type ContextTag struct {
	*tag
}

func NewContextTag(args Args) (*ContextTag, error) {
	c := &ContextTag{}
	switch len(args) {
	case 1:
		if _, ok := args[0].(PDUData); !ok {
			return nil, errors.Errorf("invalid argument %T", args[0])
		}
		_tag, err := NewTag(args)
		if err != nil {
			return nil, errors.New("error creating tag")
		}
		c.tag = _tag.(*tag)
		if c.tagClass != model.TagClass_CONTEXT_SPECIFIC_TAGS {
			return nil, errors.New("error creating tag: invalid tag class")
		}
		return c, nil
	case 2:
		var tnum any
		tnum, ok := args[0].(uint)
		if !ok {
			tnum, ok = args[0].(int)
			if !ok {
				return nil, errors.New("error creating tag: invalid tag number")
			}
		}
		tdata := args[1].([]byte)
		if len(tdata) == 0 {
			tdata = nil
		}
		_tag, err := NewTag(NewArgs(model.TagClass_CONTEXT_SPECIFIC_TAGS, tnum, len(tdata), tdata))
		if err != nil {
			return nil, errors.New("error creating tag")
		}
		c.tag = _tag.(*tag)
		return c, nil
	default:
		return nil, errors.New("requires type and data or pdu data")
	}
}

type OpeningTag struct {
	*tag
}

func NewOpeningTag(context Arg) (*OpeningTag, error) {
	o := &OpeningTag{}
	switch context.(type) {
	case PDUData:
		_tag, err := NewTag(NewArgs(context))
		if err != nil {
			return nil, errors.Wrap(err, "error creating tag")
		}
		o.tag = _tag.(*tag)
		if o.tagClass != TagOpeningTagClass {
			return nil, errors.New("opening tag not decoded")
		}
		return o, nil
	case int, uint:
		_tag, err := NewTag(NewArgs(TagOpeningTagClass, context))
		if err != nil {
			return nil, errors.Wrap(err, "error creating tag")
		}
		o.tag = _tag.(*tag)
		return o, nil
	default:
		return nil, errors.Errorf("invalid argument %T", context)
	}
}

type ClosingTag struct {
	*tag
}

func NewClosingTag(context Arg) (*ClosingTag, error) {
	o := &ClosingTag{}
	switch context.(type) {
	case PDUData:
		_tag, err := NewTag(NewArgs(context))
		if err != nil {
			return nil, errors.Wrap(err, "error creating tag")
		}
		o.tag = _tag.(*tag)
		if o.tagClass != TagClosingTagClass {
			return nil, errors.New("opening tag not decoded")
		}
		return o, nil
	case int, uint:
		_tag, err := NewTag(NewArgs(TagClosingTagClass, context))
		if err != nil {
			return nil, errors.Wrap(err, "error creating tag")
		}
		o.tag = _tag.(*tag)
		return o, nil
	default:
		return nil, errors.Errorf("invalid argument %T", context)
	}
}

type TagList struct {
	tagList []Tag
}

func NewTagList(arg Arg) *TagList {
	t := &TagList{}
	switch arg := arg.(type) {
	case []any:
		args := arg
		for _, a := range args {
			t.tagList = append(t.tagList, a.(Tag))
		}
	case []Tag:
		args := arg
		for _, a := range args {
			t.tagList = append(t.tagList, a.(Tag))
		}
	case Args:
		args := arg
		for _, a := range args {
			t.tagList = append(t.tagList, a.(Tag))
		}
	}
	return t
}

func (b *TagList) Append(tag Tag) {
	b.tagList = append(b.tagList, tag)
}

func (b *TagList) Extend(tags ...Tag) {
	for _, tag := range tags {
		b.tagList = append(b.tagList, tag)
	}
}

func (b *TagList) Peek() Tag {
	if len(b.tagList) < 1 {
		return nil
	}
	return b.tagList[0]
}

func (b *TagList) Push(tag Tag) {
	b.tagList = append([]Tag{tag}, b.tagList...)
}

func (b *TagList) Pop() Tag {
	if len(b.tagList) < 1 {
		return nil
	}
	item := b.tagList[0]
	b.tagList = b.tagList[1:]
	return item
}

// GetContext Return a tag or a list of tags context encoded.
func (b *TagList) GetContext(context uint) (any, error) {
	// forward pass
	i := 0
	for i < len(b.tagList) {
		tag := b.tagList[i]

		switch tag.GetTagClass() {
		case model.TagClass_APPLICATION_TAGS: // skip application stuff
		case model.TagClass_CONTEXT_SPECIFIC_TAGS: // check for context encoded atomic value
			if tag.GetTagNumber() == context {
				return tag, nil
			}
		case TagOpeningTagClass:
			keeper := tag.GetTagNumber() == context
			var rslt []Tag
			i += 1
			lvl := 0
		innerSearch:
			for i < len(b.tagList) {
				tag := b.tagList[i]
				switch tag.GetTagClass() {
				case TagOpeningTagClass:
					lvl += 1
				case TagClosingTagClass:
					lvl -= 1
					if lvl < 0 {
						break innerSearch
					}
				}

				rslt = append(rslt, tag)
				i += 1
			}

			// make sure everything balances
			if lvl >= 1 {
				return nil, errors.New("mismatched open/close tag")
			}

			// get everything we need
			if keeper {
				return NewTagList(rslt), nil
			}
		}
		i += 1
	}
	return nil, nil
}

func (b *TagList) Encode(data PDUData) {
	for _, tag := range b.tagList {
		tag.Encode(data)
	}
}

func (b *TagList) Decode(data PDUData) error {
	for len(data.GetPduData()) != 0 {
		var tag Tag
		tag, err := NewTag(NewArgs(data))
		if err != nil {
			return errors.Wrap(err, "error creating tag")
		}
		switch tag.GetTagClass() {
		case model.TagClass_APPLICATION_TAGS:
		case model.TagClass_CONTEXT_SPECIFIC_TAGS:
			tag, err = NewContextTag(NewArgs(tag.GetTagNumber(), tag.GetTagData()))
			if err != nil {
				panic(err)
			}
		case model.TagClass(TagOpeningTagClass):
			tag, err = NewOpeningTag(tag.GetTagNumber())
			if err != nil {
				panic(err)
			}
		case model.TagClass(TagClosingTagClass):
			tag, err = NewClosingTag(tag.GetTagNumber())
			if err != nil {
				panic(err)
			}
		}
		b.tagList = append(b.tagList, tag)
	}
	return nil
}

func (b *TagList) GetTagList() []Tag {
	return b.tagList
}

type ComparableAndOrdered interface {
	comparable
	cmp.Ordered
}

// AtomicContract provides a set of functions which can be overwritten by a sub struct
type AtomicContract[T ComparableAndOrdered] interface {
	Compare(other any) int
	LowerThan(other any) bool
	Equals(other any) bool
	GetValue() T
}

// AtomicRequirements provides a set of functions which must be overwritten by a sub struct
type AtomicRequirements interface {
	IsValid(arg any) bool
}

// Atomic is an abstract struct
type Atomic[T ComparableAndOrdered] struct {
	AtomicContract[T]
	atomicRequirements AtomicRequirements

	value T
}

func NewAtomic[T ComparableAndOrdered](subStruct interface {
	AtomicContract[T]
	AtomicRequirements
}) *Atomic[T] {
	return &Atomic[T]{
		AtomicContract:     subStruct,
		atomicRequirements: subStruct,
	}
}

func (a *Atomic[T]) Compare(other any) int {
	otherValue := other.(AtomicContract[T]).GetValue()
	// now compare the values
	if a.value < otherValue {
		return -1
	} else if a.value > otherValue {
		return 1
	} else {
		return 0
	}
}

func (a *Atomic[T]) LowerThan(other any) bool {
	otherValue := other.(AtomicContract[T]).GetValue()
	// now compare the values
	return a.value < otherValue
}

func (a *Atomic[T]) Equals(other any) bool {
	otherValue := other.(AtomicContract[T]).GetValue()
	// now compare the values
	return a.value == otherValue
}

func (a *Atomic[T]) GetValue() T {
	return a.value
}

func (a *Atomic[T]) Coerce(arg any) T {
	return arg.(AtomicContract[T]).GetValue()
}

type CommonMath struct {
	// TODO: implement me
}

type Null struct {
	*Atomic[int]
}

func NewNull(arg Arg) (*Null, error) {
	b := &Null{}
	b.Atomic = NewAtomic[int](b)

	if arg == nil {
		return b, nil
	}
	switch arg := arg.(type) {
	case *tag:
		err := b.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return b, nil
	case *Null:
		b.value = arg.value
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", arg)
	}

	return b, nil
}

func (b *Null) Encode(tag Tag) {
	tag.set(NewArgs(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_NULL, b.value, []byte{}))
}

func (b *Null) Decode(tag Tag) error {
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(model.BACnetDataType_NULL) {
		return errors.New("Null application tag required")
	}
	if tag.GetTagLvt() > 1 {
		return errors.New("invalid tag value")
	}

	// get the data
	if tag.GetTagLvt() == 1 {
		b.value = 1
	}
	return nil
}

func (b *Null) IsValid(arg any) bool {
	_, ok := arg.(bool)
	return ok
}

func (b *Null) String() string {
	value := "False"
	if b.value == 1 {
		value = "True"
	}
	return fmt.Sprintf("Null(%s)", value)
}

type Boolean struct {
	*Atomic[int] //Note we need int as bool can't be used
}

func NewBoolean(arg Arg) (*Boolean, error) {
	b := &Boolean{}
	b.Atomic = NewAtomic[int](b)

	if arg == nil {
		return b, nil
	}
	switch arg := arg.(type) {
	case *tag:
		err := b.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return b, nil
	case bool:
		if arg {
			b.value = 1
		}
	case *Boolean:
		b.value = arg.value
	case string:
		switch arg {
		case "True", "true":
			b.value = 1
		case "False", "false":
		default:
			return nil, errors.Errorf("invalid string: %s", arg)
		}
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", arg)
	}

	return b, nil
}

func (b *Boolean) Encode(tag Tag) {
	tag.set(NewArgs(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_BOOLEAN, b.value, []byte{}))
}

func (b *Boolean) Decode(tag Tag) error {
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(model.BACnetDataType_BOOLEAN) {
		return errors.New("boolean application tag required")
	}
	if tag.GetTagLvt() > 1 {
		return errors.New("invalid tag value")
	}

	// get the data
	if tag.GetTagLvt() == 1 {
		b.value = 1
	}
	return nil
}

func (b *Boolean) IsValid(arg any) bool {
	_, ok := arg.(bool)
	return ok
}

// GetValue gives an int value because bool can't be used in constraint. A convenience method GetBoolValue exists.
func (b *Boolean) GetValue() int {
	return b.Atomic.GetValue()
}

func (b *Boolean) GetBoolValue() bool {
	return b.GetValue() == 1
}

func (b *Boolean) String() string {
	value := "False"
	if b.value == 1 {
		value = "True"
	}
	return fmt.Sprintf("Boolean(%s)", value)
}

type Unsigned struct {
	*Atomic[uint32]
	*CommonMath
}

func NewUnsigned(arg Arg) (*Unsigned, error) {
	i := &Unsigned{}
	i.Atomic = NewAtomic[uint32](i)

	if arg == nil {
		return i, nil
	}
	switch arg := arg.(type) {
	case *tag:
		err := i.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return i, nil
	case uint32:
		i.value = arg
	case uint:
		if !i.IsValid(arg) {
			return nil, errors.Errorf("invalid Unsigned: %d", arg)
		}
		i.value = uint32(arg)
	case int32:
		if !i.IsValid(arg) {
			return nil, errors.Errorf("invalid Unsigned: %d", arg)
		}
		i.value = uint32(arg)
	case int:
		if !i.IsValid(arg) {
			return nil, errors.Errorf("invalid Unsigned: %d", arg)
		}
		i.value = uint32(arg)
	case *Unsigned:
		i.value = arg.value
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", arg)
	}

	return i, nil
}

func (i *Unsigned) Encode(tag Tag) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, i.value)

	// reduce the value to the smallest number of bytes
	for len(data) > 1 && data[0] == 0 {
		data = data[1:]
	}

	tag.setAppData(uint(model.BACnetDataType_UNSIGNED_INTEGER), data)
}

func (i *Unsigned) Decode(tag Tag) error {
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(model.BACnetDataType_UNSIGNED_INTEGER) {
		return errors.New("Unsigned application tag required")
	}
	if len(tag.GetTagData()) == 0 {
		return errors.New("invalid tag length")
	}

	tagData := tag.GetTagData()

	// get the data
	rslt := uint32(0)
	for _, c := range tagData {
		rslt = (rslt << 8) + uint32(c)
	}

	// save the result
	i.value = rslt
	return nil
}

func (i *Unsigned) IsValid(arg any) bool {
	switch arg := arg.(type) {
	case string:
		_, err := strconv.Atoi(arg)
		return err == nil
	case int:
		return arg >= 0
	case int32:
		return arg >= 0
	case uint:
		return arg <= math.MaxUint32
	case uint32:
		return true
	default:
		return false
	}
}

func (i *Unsigned) String() string {
	return fmt.Sprintf("Unsigned(%d)", i.value)
}

type Unsigned8 struct {
	*Atomic[uint8]
	*CommonMath
}

func NewUnsigned8(arg Arg) (*Unsigned8, error) {
	i := &Unsigned8{}
	i.Atomic = NewAtomic[uint8](i)

	if arg == nil {
		return i, nil
	}
	switch arg := arg.(type) {
	case *tag:
		err := i.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return i, nil
	case uint8:
		i.value = arg
	case uint:
		if !i.IsValid(arg) {
			return nil, errors.Errorf("invalid Unsigned8: %d", arg)
		}
		i.value = uint8(arg)
	case int:
		if !i.IsValid(arg) {
			return nil, errors.Errorf("invalid Unsigned8: %d", arg)
		}
		i.value = uint8(arg)
	case *Unsigned8:
		i.value = arg.value
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", arg)
	}

	return i, nil
}

func (i *Unsigned8) Encode(tag Tag) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(i.value))

	// reduce the value to the smallest number of bytes
	for len(data) > 1 && data[0] == 0 {
		data = data[1:]
	}

	tag.setAppData(uint(model.BACnetDataType_UNSIGNED_INTEGER), data)
}

func (i *Unsigned8) Decode(tag Tag) error {
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(model.BACnetDataType_UNSIGNED_INTEGER) {
		return errors.New("Unsigned8 application tag required")
	}
	if len(tag.GetTagData()) == 0 {
		return errors.New("invalid tag length")
	}

	tagData := tag.GetTagData()

	// get the data
	rslt := uint8(0)
	for _, c := range tagData {
		rslt = (rslt << 8) + c
	}

	// save the result
	i.value = rslt
	return nil
}

func (i *Unsigned8) IsValid(arg any) bool {
	switch arg := arg.(type) {
	case string:
		_, err := strconv.Atoi(arg)
		return err == nil
	case int:
		return arg > 0 && arg < 256
	case int8:
		return arg > 0
	case uint:
		return arg <= math.MaxUint32
	case uint8:
		return true
	default:
		return false
	}
}

func (i *Unsigned8) String() string {
	return fmt.Sprintf("Unsigned8(%d)", i.value)
}

type Unsigned16 struct {
	*Atomic[uint16]
	*CommonMath
}

func NewUnsigned16(arg Arg) (*Unsigned16, error) {
	i := &Unsigned16{}
	i.Atomic = NewAtomic[uint16](i)

	if arg == nil {
		return i, nil
	}
	switch arg := arg.(type) {
	case *tag:
		err := i.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return i, nil
	case uint16:
		i.value = arg
	case uint:
		if !i.IsValid(arg) {
			return nil, errors.Errorf("invalid Unsigned16: %d", arg)
		}
		i.value = uint16(arg)
	case int:
		if !i.IsValid(arg) {
			return nil, errors.Errorf("invalid Unsigned16: %d", arg)
		}
		i.value = uint16(arg)
	case *Unsigned16:
		i.value = arg.value
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", arg)
	}

	return i, nil
}

func (i *Unsigned16) Encode(tag Tag) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(i.value))

	// reduce the value to the smallest number of bytes
	for len(data) > 1 && data[0] == 0 {
		data = data[1:]
	}

	tag.setAppData(uint(model.BACnetDataType_UNSIGNED_INTEGER), data)
}

func (i *Unsigned16) Decode(tag Tag) error {
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(model.BACnetDataType_UNSIGNED_INTEGER) {
		return errors.New("Unsigned16 application tag required")
	}
	if len(tag.GetTagData()) == 0 {
		return errors.New("invalid tag length")
	}

	tagData := tag.GetTagData()

	// get the data
	rslt := uint16(0)
	for _, c := range tagData {
		rslt = (rslt << 8) + uint16(c)
	}

	// save the result
	i.value = rslt
	return nil
}

func (i *Unsigned16) IsValid(arg any) bool {
	switch arg := arg.(type) {
	case string:
		_, err := strconv.Atoi(arg)
		return err == nil
	case int:
		return arg > 0 && arg < 65536
	case int16:
		return arg > 0
	case uint:
		return arg <= math.MaxUint32
	case uint16:
		return true
	default:
		return false
	}
}

func (i *Unsigned16) String() string {
	return fmt.Sprintf("Unsigned16(%d)", i.value)
}

type Integer struct {
	*Atomic[int32]
	*CommonMath
}

func NewInteger(arg Arg) (*Integer, error) {
	i := &Integer{}
	i.Atomic = NewAtomic[int32](i)

	if arg == nil {
		return i, nil
	}
	switch arg := arg.(type) {
	case *tag:
		err := i.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return i, nil
	case int32:
		i.value = arg
	case int:
		if !i.IsValid(arg) {
			return nil, errors.Errorf("invalid integer: %d", arg)
		}
		i.value = int32(arg)
	case *Integer:
		i.value = arg.value
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", arg)
	}

	return i, nil
}

func (i *Integer) Encode(tag Tag) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(i.value))

	// reduce the value to the smallest number of bytes, be
	// careful about sign extension
	if i.value < 0 {
		for len(data) > 1 {
			if data[0] != 255 {
				break
			}
			if data[1] < 128 {
				break
			}
			data = data[1:]
		}
	} else {
		for len(data) > 1 {
			if data[0] != 0 {
				break
			}
			if data[1] >= 128 {
				break
			}
			data = data[1:]
		}
	}

	tag.setAppData(uint(model.BACnetDataType_SIGNED_INTEGER), data)
}

func (i *Integer) Decode(tag Tag) error {
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(model.BACnetDataType_SIGNED_INTEGER) {
		return errors.New("Integer application tag required")
	}
	if len(tag.GetTagData()) == 0 {
		return errors.New("invalid tag length")
	}

	tagData := tag.GetTagData()

	// get the data
	rslt := int32(tagData[0])
	if rslt&0x80 != 0 {
		rslt = (-1 << 8) | rslt
	}
	if len(tagData) > 1 {
		for _, c := range tagData[1:] {
			rslt = (rslt << 8) | int32(c)
		}
	}

	// save the result
	i.value = rslt
	return nil
}

func (i *Integer) IsValid(arg any) bool {
	switch arg := arg.(type) {
	case int:
		return arg >= math.MinInt32 && arg <= math.MaxInt32
	case int32:
		return true
	default:
		return false
	}
}

func (i *Integer) String() string {
	return fmt.Sprintf("Integer(%d)", i.value)
}

type Real struct {
	*Atomic[float32]
	*CommonMath
}

func NewReal(arg Arg) (*Real, error) {
	b := &Real{}
	b.Atomic = NewAtomic[float32](b)

	if arg == nil {
		return b, nil
	}
	switch arg := arg.(type) {
	case *tag:
		err := b.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return b, nil
	case float32:
		b.value = arg
	case float64:
		b.IsValid(arg)
		b.value = float32(arg)
	case int:
		b.value = float32(arg)
	case *Real:
		b.value = arg.value
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", arg)
	}

	return b, nil
}

func (d *Real) Encode(tag Tag) {
	var _b = make([]byte, 4)
	binary.BigEndian.PutUint32(_b, math.Float32bits(d.value))
	tag.setAppData(uint(model.BACnetDataType_REAL), _b)
}

func (d *Real) Decode(tag Tag) error {
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(model.BACnetDataType_REAL) {
		return errors.New("Real application tag required")
	}
	if len(tag.GetTagData()) != 4 {
		return errors.New("invalid tag length")
	}

	// extract the data
	d.value = math.Float32frombits(binary.BigEndian.Uint32(tag.GetTagData()))
	return nil
}

func (d *Real) IsValid(arg any) bool {
	switch arg := arg.(type) {
	case float32:
		return true
	case float64:
		if arg > math.MaxFloat32 || -arg > math.MaxFloat32 {
			return false
		}
		return true
	default:
		return false
	}
}

func (d *Real) String() string {
	return fmt.Sprintf("Real(%g)", d.value)
}

type Double struct {
	*Atomic[float64]
	*CommonMath
}

func NewDouble(arg Arg) (*Double, error) {
	b := &Double{}
	b.Atomic = NewAtomic[float64](b)

	if arg == nil {
		return b, nil
	}
	switch arg := arg.(type) {
	case *tag:
		err := b.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return b, nil
	case float32:
		b.value = float64(arg)
	case float64:
		b.value = arg
	case int:
		b.value = float64(arg)
	case *Double:
		b.value = arg.value
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", arg)
	}

	return b, nil
}

func (d *Double) Encode(tag Tag) {
	var _b = make([]byte, 8)
	binary.BigEndian.PutUint64(_b, math.Float64bits(d.value))
	tag.setAppData(uint(model.BACnetDataType_DOUBLE), _b)
}

func (d *Double) Decode(tag Tag) error {
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(model.BACnetDataType_DOUBLE) {
		return errors.New("Double application tag required")
	}
	if len(tag.GetTagData()) != 8 {
		return errors.New("invalid tag length")
	}

	// extract the data
	d.value = math.Float64frombits(binary.BigEndian.Uint64(tag.GetTagData()))
	return nil
}

func (d *Double) IsValid(arg any) bool {
	_, ok := arg.(float64)
	return ok
}

func (d *Double) String() string {
	return fmt.Sprintf("Double(%g)", d.value)
}

type OctetString struct {
	value []byte
}

func NewOctetString(arg Arg) (*OctetString, error) {
	o := &OctetString{}
	o.value = make([]byte, 0)

	if arg == nil {
		return o, nil
	}
	switch arg := arg.(type) {
	case *tag:
		err := o.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return o, nil
	case []byte:
		if len(arg) == 0 {
			arg = nil
		}
		o.value = arg
	case *OctetString:
		o.value = arg.value
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", arg)
	}

	return o, nil
}

func (o *OctetString) Encode(tag Tag) {
	tag.setAppData(uint(model.BACnetDataType_OCTET_STRING), o.value)
}

func (o *OctetString) Decode(tag Tag) error {
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(model.BACnetDataType_OCTET_STRING) {
		return errors.New("OctetString application tag required")
	}

	o.value = tag.GetTagData()
	return nil
}

func (o *OctetString) Compare(other any) int {
	if _, ok := other.(byte); !ok {
		return -1
	}
	return len(o.value) - len(other.(OctetString).value)
}

func (o *OctetString) LowerThan(other any) bool {
	if _, ok := other.(byte); !ok {
		return false
	}
	return len(o.value) < len(other.(OctetString).value)
}

func (o *OctetString) Equals(other any) bool {
	if _, ok := other.(byte); !ok {
		return false
	}
	return bytes.Equal(o.value, other.([]byte))
}

func (o *OctetString) GetValue() []byte {
	return o.value
}

func (o *OctetString) IsValid(arg any) bool {
	_, ok := arg.([]byte)
	return ok
}

func (o *OctetString) String() string {
	return fmt.Sprintf("OctetString(X'%s')", Btox(o.value, ""))
}

type CharacterString struct {
	*Atomic[string]
	*CommonMath

	strEncoding byte
	strValue    []byte
}

func NewCharacterString(arg Arg) (*CharacterString, error) {
	c := &CharacterString{}
	c.Atomic = NewAtomic[string](c)

	if arg == nil {
		return c, nil
	}
	switch arg := arg.(type) {
	case *tag:
		err := c.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return c, nil
	case string:
		c.value = arg
		c.strValue = []byte(c.value)
	case *CharacterString:
		c.value = arg.value
		c.strEncoding = arg.strEncoding
		c.strValue = arg.strValue
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", arg)
	}

	return c, nil
}

func (c *CharacterString) Encode(tag Tag) {
	tag.setAppData(uint(model.BACnetDataType_CHARACTER_STRING), append([]byte{c.strEncoding}, c.strValue...))
}

func (c *CharacterString) Decode(tag Tag) error {
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(model.BACnetDataType_CHARACTER_STRING) {
		return errors.New("CharacterString application tag required")
	}
	if len(tag.GetTagData()) == 0 {
		return errors.New("invalid tag length")
	}

	tagData := tag.GetTagData()

	// extract the data
	c.strEncoding = tagData[0]
	c.strValue = tagData[1:]

	// normalize the value
	switch c.strEncoding {
	case 0:
		c.value = string(c.strValue)
	case 3: //utf_32be
		panic("implement me") // TODO: implement me
	case 4: //utf_16be
		panic("implement me") // TODO: implement me
	case 5: //latin_1
		panic("implement me") // TODO: implement me
	default:
		c.value = fmt.Sprintf("### unknown encoding: %d ###", c.strEncoding)
	}

	return nil
}

func (c *CharacterString) IsValid(arg any) bool {
	_, ok := arg.(bool)
	return ok
}

func (c *CharacterString) String() string {
	return fmt.Sprintf("CharacterString(%d,X'%s')", c.strEncoding, Btox(c.strValue, ""))
}

// BitStringExtension can be used to inherit from BitString
type BitStringExtension interface {
	fmt.Stringer
	GetBitNames() map[string]int
	GetBitLen() int
}

type BitString struct {
	bitStringExtension BitStringExtension
	value              []bool
}

func NewBitString(args Args) (*BitString, error) {
	return NewBitStringWithExtension(nil, args)
}

func NewBitStringWithExtension(bitStringExtension BitStringExtension, args Args) (*BitString, error) {
	b := &BitString{
		bitStringExtension: bitStringExtension,
	}
	if len(args) == 0 {
		return b, nil
	}
	if len(args) > 1 {
		return nil, errors.New("too many arguments")
	}
	if bitStringExtension != nil {
		b.value = make([]bool, bitStringExtension.GetBitLen())
	}
	switch arg := args[0].(type) {
	case *tag:
		err := b.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "decoding tag failed")
		}
	case []int:
		b.value = make([]bool, len(arg))
		for i, v := range arg {
			b.value[i] = v != 0
		}
	case []bool:
		b.value = arg
	case []string:
		bitNames := make(map[string]int)
		if bitStringExtension != nil {
			bitNames = bitStringExtension.GetBitNames()
		}
		for _, bit := range arg {
			bit, ok := bitNames[bit]
			if !ok || bit < 0 || bit > len(b.value) {
				return nil, errors.New("constructorElement out of range")
			}
			b.value[bit] = true
		}
	case *BitString:
		b.value = arg.value[:]
	case model.BACnetApplicationTagBitStringExactly:
		b.value = arg.GetPayload().GetData()
	default:
		return nil, errors.Errorf("no support for %T yet", arg)
	}
	return b, nil
}

func (b *BitString) Decode(tag Tag) error {
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(model.BACnetDataType_BIT_STRING) {
		return errors.New("bit string application tag required")
	}
	if len(tag.GetTagData()) == 0 {
		return errors.New("invalid tag length")
	}
	// extract the number of unused bits
	unused := tag.GetTagData()[0]

	// extract the data
	data := make([]bool, 0)
	for _, x := range tag.GetTagData()[1:] {
		for i := range 8 {
			if (x & (1 << (7 - i))) != 0 {
				data = append(data, true)
			} else {
				data = append(data, false)
			}
		}
	}

	// trim off the unused bits
	if unused != 0 && unused != 8 {
		b.value = data[:len(data)-int(unused)]
	} else {
		b.value = data
	}
	return nil
}

func (b *BitString) Encode(tag Tag) {
	used := len(b.value) % 8
	unused := 8 - used
	if unused == 8 {
		unused = 0
	}

	// start with the number of unused bits
	data := []byte{byte(unused)}

	// build and append each packed octet
	bits := append(b.value, make([]bool, unused)...)
	for i := range len(bits) / 8 {
		i = i * 8
		x := byte(0)
		for j := range 8 {
			bit := bits[i+j]
			bitValue := byte(0)
			if bit {
				bitValue = 1
			}
			x |= bitValue << (7 - j)
		}
		data = append(data, x)
	}

	tag.setAppData(uint(model.BACnetDataType_BIT_STRING), data)
}

func (b *BitString) Compare(other any) int {
	switch other := other.(type) {
	case *BitString:
		return len(b.value) - len(other.value)
	default:
		return -1
	}
}

func (b *BitString) LowerThan(other any) bool {
	switch other := other.(type) {
	case *BitString:
		return len(b.value) < len(other.value)
	default:
		return false
	}
}

func (b *BitString) Equals(other any) bool {
	return b == other
}

func (b *BitString) GetValue() []bool {
	return b.value
}

func (b *BitString) String() string {
	// flip the bit names
	bitNames := map[int]string{}
	if b.bitStringExtension != nil {
		for key, value := range b.bitStringExtension.GetBitNames() {
			bitNames[value] = key
		}
	}

	// build a list of values and/or names
	var valueList []string
	for index, value := range b.value {
		if name, ok := bitNames[index]; ok {
			if value == true {
				valueList = append(valueList, name)
			} else {
				valueList = append(valueList, "!"+name)
			}
		} else {
			if value {
				valueList = append(valueList, "1")
			} else {
				valueList = append(valueList, "0")
			}
		}
	}

	// bundle it together
	return fmt.Sprintf("BitString(%v)", strings.Join(valueList, ","))
}

func expandEnumerations(e EnumeratedContract) {
	xlateTable := e.GetXlateTable()

	for name, value := range e.GetEnumerations() {
		// save the result
		xlateTable[name] = value
		xlateTable[value] = name
	}
}

// EnumeratedContract provides a set of functions which can be overwritten by a sub struct
type EnumeratedContract interface {
	GetEnumerations() map[string]uint64
	GetXlateTable() map[any]any
	// SetEnumerated is required because we do more stuff in the constructor and can't wait for the substruct to finish that
	SetEnumerated(enumerated *Enumerated)
}

type Enumerated struct {
	*Atomic[uint64]
	EnumeratedContract

	_xlateTable map[any]any

	valueString string
}

func NewEnumerated(args ...any) (*Enumerated, error) {
	e := &Enumerated{}
	e.EnumeratedContract = e
	e.Atomic = NewAtomic[uint64](e)

	if args == nil {
		return e, nil
	}
	var arg any
	switch len(args) {
	case 1:
		arg = args[0]
	case 2:
		var ok bool
		e.EnumeratedContract, ok = args[0].(EnumeratedContract)
		if !ok {
			return nil, errors.Errorf("%T must be implement EnumeratedContract", args[0])
		}
		e.EnumeratedContract.SetEnumerated(e)
		arg = args[1]
	default:
		return nil, errors.Errorf("invalid arguments %T", args)
	}
	if len(e.EnumeratedContract.GetXlateTable()) == 0 {
		expandEnumerations(e.EnumeratedContract)
	}

	switch arg := arg.(type) {
	case *tag:
		err := e.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return e, nil
	case uint:
		e.value = uint64(arg)

		// convert it to a string if you can
		e.valueString, _ = e.EnumeratedContract.GetXlateTable()[e.value].(string)
	case int:
		if arg < 0 {
			return nil, errors.New("arg must be positive")
		}
		e.value = uint64(arg)

		// convert it to a string if you can
		e.valueString, _ = e.EnumeratedContract.GetXlateTable()[e.value].(string)
	case uint64:
		e.value = arg

		// convert it to a string if you can
		e.valueString, _ = e.EnumeratedContract.GetXlateTable()[e.value].(string)
	case string:
		var ok bool
		var value any
		value, ok = e.EnumeratedContract.GetXlateTable()[arg]
		if !ok {
			return nil, errors.Errorf("undefined enumeration %s", arg)
		}
		e.value = value.(uint64)
		e.valueString = arg
	case *Enumerated:
		e.value = arg.value
		e.valueString = arg.valueString
		e._xlateTable = make(map[any]any)
		for k, v := range arg._xlateTable {
			e._xlateTable[k] = v
		}
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", arg)
	}

	return e, nil
}

func (e *Enumerated) GetEnumerations() map[string]uint64 {
	return make(map[string]uint64)
}

func (e *Enumerated) GetXlateTable() map[any]any {
	if e._xlateTable == nil {
		e._xlateTable = make(map[any]any)
	}
	return e._xlateTable
}

func (e *Enumerated) SetEnumerated(_ *Enumerated) {
	panic("must be implemented by substruct")
}

func (e *Enumerated) GetItem(item any) (result any, ok bool) {
	v, ok := e.EnumeratedContract.GetXlateTable()[item]
	return v, ok
}

func (e *Enumerated) GetLong() uint64 {
	if mappedValue, ok := e.EnumeratedContract.GetXlateTable()[e.valueString]; ok {
		return mappedValue.(uint64)
	}
	return e.value
}

func (e *Enumerated) Keylist() []string {
	var result []string
	for key := range e.EnumeratedContract.GetEnumerations() {
		result = append(result, key)
	}
	return result
}

func (e *Enumerated) Compare(other any) int {
	otherEnumerated, ok := other.(Enumerated)
	if !ok {
		return -1
	}

	// get the numeric version
	a := e.GetLong()
	b := otherEnumerated.GetLong()

	// now compare the values
	if a < b {
		return -1
	} else if a > b {
		return 1
	} else {
		return 0
	}
}

func (e *Enumerated) Decode(tag Tag) error {
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(model.BACnetDataType_ENUMERATED) {
		return errors.New("bit string application tag required")
	}
	if len(tag.GetTagData()) == 0 {
		return errors.New("invalid tag length")
	}

	data := tag.GetTagData()
	if len(data) < 8 {
		data = append(make([]byte, 8-len(data)), data...)
	}
	// get the data
	rslt := binary.BigEndian.Uint64(data)

	if mappedValue, ok := e.EnumeratedContract.GetXlateTable()[rslt]; ok {
		e.valueString = mappedValue.(string)
	}

	// save the result
	e.value = rslt
	return nil
}

func (e *Enumerated) Encode(tag Tag) {
	value := e.value
	if mappedValue, ok := e.EnumeratedContract.GetXlateTable()[e.valueString]; ok {
		value = mappedValue.(uint64)
	}

	data := make([]byte, 8)
	// rip apart the number
	binary.BigEndian.PutUint64(data, value)

	// reduce the value to the smallest number of octets
	for len(data) > 1 && data[0] == 0 {
		data = data[1:]
	}

	// encode the tag
	tag.setAppData(uint(model.BACnetDataType_ENUMERATED), data)
}

func (e *Enumerated) IsValid(arg any) bool {
	_, ok := arg.(uint64)
	return ok
}

func (e *Enumerated) GetValueString() string {
	return e.valueString
}

func (e *Enumerated) String() string {
	value := strconv.Itoa(int(e.value))
	if e.valueString != "" {
		value = e.valueString
	}
	return fmt.Sprintf("Enumerated(%s)", value)
}

const _mm = `(?P<month>0?[1-9]|1[0-4]|odd|even|255|[*])`
const _dd = `(?P<day>[0-3]?\d|last|odd|even|255|[*])`
const _yy = `(?P<year>\d{2}|255|[*])`
const _yyyy = `(?P<year>\d{4}|255|[*])`
const _dow = `(?P<dow>[1-7]|mon|tue|wed|thu|fri|sat|sun|255|[*])`

var _special_mon = map[string]int{"*": 255, "odd": 13, "even": 14, "": 255}
var _special_mon_inv = map[int]string{255: "*", 13: "odd", 14: "even"}

var _special_day = map[string]int{"*": 255, "last": 32, "odd": 33, "even": 34, "": 255}
var _special_day_inv = map[int]string{255: "*", 32: "last", 33: "odd", 34: "even"}

var _special_dow = map[string]int{"*": 255, "mon": 1, "tue": 2, "wed": 3, "thu": 4, "fri": 5, "sat": 6, "sun": 7}
var _special_dow_inv = map[int]string{255: "*", 1: "mon", 2: "tue", 3: "wed", 4: "thu", 5: "fri", 6: "sat", 7: "sun"}

// Create a composite pattern and compile it.
func _merge(args ...string) *regexp.Regexp {
	return regexp.MustCompile(`^` + strings.Join(args, `[/-]`) + `(?:\s+` + _dow + `)?$`)
}

// make a list of compiled patterns
var _date_patterns = []*regexp.Regexp{
	_merge(_yyyy, _mm, _dd),
	_merge(_mm, _dd, _yyyy),
	_merge(_dd, _mm, _yyyy),
	_merge(_yy, _mm, _dd),
	_merge(_mm, _dd, _yy),
	_merge(_dd, _mm, _yy),
}

type DateTuple struct {
	Year      int
	Month     int
	Day       int
	DayOfWeek int
}

type Date struct {
	value DateTuple
}

func NewDate(arg Arg, args Args) (*Date, error) {
	d := &Date{}
	year := 255
	if len(args) > 0 {
		year = args[0].(int)
	}
	if year >= 1900 {
		year = year - 1900
	}
	d.value.Year = year
	month := 0xff
	if len(args) > 1 {
		month = args[1].(int)
	}
	d.value.Month = month
	day := 0xff
	if len(args) > 2 {
		day = args[2].(int)
	}
	d.value.Day = day
	dayOfWeek := 0xff
	if len(args) > 3 {
		dayOfWeek = args[3].(int)
	}
	d.value.DayOfWeek = dayOfWeek

	if arg == nil {
		return d, nil
	}
	switch arg := arg.(type) {
	case *tag:
		err := d.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return d, nil
	case DateTuple:
		d.value = arg
	case string:
		// lower case everything
		arg = strings.ToLower(arg)

		// make a list of the contents from matching patterns
		matches := [][]string{}
		for _, p := range _date_patterns {
			if p.MatchString(arg) {
				groups := combined_pattern.FindStringSubmatch(arg)
				matches = append(matches, groups[1:])
			}
		}
		if len(matches) == 0 {
			return nil, errors.New("unmatched")
		}

		var match []string
		if len(matches) == 1 {
			match = matches[0]
		} else {
			// check to see if they really are the same
			panic("what to do here")
		}

		// extract the year and normalize
		matchedYear := match[0]
		if matchedYear == "*" || matchedYear == "" {
			year = 0xff
		} else {
			yearParse, err := strconv.ParseInt(matchedYear, 10, 64)
			if err != nil {
				return nil, errors.Wrap(err, "error parsing year")
			}
			year = int(yearParse)
			if year == 0xff {
				return d, nil
			}
			if year < 35 {
				year += 2000
			} else if year < 100 {
				year += 1900
			} else if year < 1900 {
				return nil, errors.New("invalid year")
			}
		}

		// extract the month and normalize
		matchedmonth := match[0]
		if specialMonth, ok := _special_mon[matchedmonth]; ok {
			month = specialMonth
		} else {
			monthParse, err := strconv.ParseInt(matchedmonth, 10, 64)
			if err != nil {
				return nil, errors.Wrap(err, "error parsing month")
			}
			month = int(monthParse)
			if month == 0xff {
				return d, nil
			}
			if month == 0 || month > 14 {
				return nil, errors.New("invalid month")
			}
		}

		// extract the day and normalize
		matchedday := match[0]
		if specialday, ok := _special_day[matchedday]; ok {
			day = specialday
		} else {
			dayParse, err := strconv.ParseInt(matchedday, 10, 64)
			if err != nil {
				return nil, errors.Wrap(err, "error parsing day")
			}
			day = int(dayParse)
			if day == 0xff {
				return d, nil
			}
			if day == 0 || day > 34 {
				return nil, errors.New("invalid day")
			}
		}

		// extract the dayOfWeek and normalize
		matcheddayOfWeek := match[0]
		if specialdayOfWeek, ok := _special_dow[matcheddayOfWeek]; ok {
			dayOfWeek = specialdayOfWeek
		} else if matcheddayOfWeek == "" {
			return d, nil
		} else {
			dayOfWeekParse, err := strconv.ParseInt(matcheddayOfWeek, 10, 64)
			if err != nil {
				return nil, errors.Wrap(err, "error parsing dayOfWeek")
			}
			dayOfWeek = int(dayOfWeekParse)
			if dayOfWeek == 0xff {
				return d, nil
			}
			if dayOfWeek > 7 {
				return nil, errors.New("invalid dayOfWeek")
			}
		}

		// year becomes the correct octet
		if year != 0xff {
			year -= 1900
		}

		// save the value
		d.value.Year = year
		d.value.Month = month
		d.value.Day = day
		d.value.DayOfWeek = dayOfWeek

		// calculate the day of the week
		if dayOfWeek == 0 {
			d.calcDayOfWeek()
		}
	case *Date:
		d.value = arg.value
	case float32:
		d.now(arg)
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", arg)
	}

	return d, nil
}

func (d *Date) calcDayOfWeek() {
	year, month, day, dayOfWeek := d.value.Year, d.value.Month, d.value.Day, d.value.DayOfWeek

	// assume the worst
	dayOfWeek = 255

	// check for special values
	if year == 255 {
		return
	} else if _, ok := _special_mon_inv[month]; ok {
		return
	} else if _, ok := _special_day_inv[month]; ok {
		return
	} else {
		var today time.Time
		today = time.Date(year+1900, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		today.Add(24 * time.Hour)
		dayOfWeek = int(today.Weekday())
	}

	// put it back together
	d.value.Year = year
	d.value.Month = month
	d.value.Day = day
	d.value.DayOfWeek = dayOfWeek
}

func (d *Date) now(arg float32) {
	panic("implement me") // TODO
}

func (d *Date) Encode(tag Tag) {
	tag.setAppData(uint(model.BACnetDataType_DATE), []byte{byte(d.value.Year), byte(d.value.Month), byte(d.value.Day), byte(d.value.DayOfWeek)})
}

func (d *Date) Decode(tag Tag) error {
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(model.BACnetDataType_DATE) {
		return errors.New("Date application tag required")
	}
	if len(tag.GetTagData()) != 4 {
		return errors.New("invalid tag length")
	}

	arg := tag.GetTagData()
	year, month, day, dayOfWeek := arg[0], arg[1], arg[2], arg[3]
	d.value.Year, d.value.Month, d.value.Day, d.value.DayOfWeek = int(year), int(month), int(day), int(dayOfWeek)
	return nil
}

func (d *Date) IsValid(arg any) bool {
	_, ok := arg.(bool)
	return ok
}

func (d *Date) Compare(other any) int {
	switch other := other.(type) {
	case *Date:
		_ = other //TODO: implement me
		return -1
	default:
		return -1
	}
}

func (d *Date) LowerThan(other any) bool {
	switch other := other.(type) {
	case *Date:
		// return d.getLong() < other.getLong()
		_ = other // TODO: implement me
		return false
	default:
		return false
	}
}

func (d *Date) Equals(other any) bool {
	return d.value == other
}

func (d *Date) GetValue() DateTuple {
	return d.value
}

func (d *Date) Coerce(arg Date) DateTuple {
	return arg.GetValue()
}

func (d *Date) String() string {
	year, month, day, dayOfWeek := d.value.Year, d.value.Month, d.value.Day, d.value.DayOfWeek
	yearStr := "*"
	if year != 255 {
		yearStr = strconv.Itoa(year + 1900)
	}
	monthStr := strconv.Itoa(month)
	if ms, ok := _special_mon_inv[month]; ok {
		monthStr = ms
	}
	dayStr := strconv.Itoa(day)
	if ms, ok := _special_day_inv[day]; ok {
		dayStr = ms
	}
	dowStr := strconv.Itoa(dayOfWeek)
	if ms, ok := _special_dow_inv[dayOfWeek]; ok {
		dowStr = ms
	}

	return fmt.Sprintf("Date(%s-%s-%s %s)", yearStr, monthStr, dayStr, dowStr)
}

type TimeTuple struct {
	Hour      int
	Minute    int
	Second    int
	Hundredth int
}

type Time struct {
	value TimeTuple
}

func NewTime(arg Arg, args Args) (*Time, error) {
	d := &Time{}
	hour := 255
	if len(args) > 0 {
		hour = args[0].(int)
	}
	d.value.Hour = hour
	minute := 0xff
	if len(args) > 1 {
		minute = args[1].(int)
	}
	d.value.Minute = minute
	second := 0xff
	if len(args) > 2 {
		second = args[2].(int)
	}
	d.value.Second = second
	hundredth := 0xff
	if len(args) > 3 {
		hundredth = args[3].(int)
	}
	d.value.Hundredth = hundredth

	if arg == nil {
		return d, nil
	}
	switch arg := arg.(type) {
	case *tag:
		err := d.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return d, nil
	case TimeTuple:
		d.value = arg
	case string:
		// lower case everything
		arg = strings.ToLower(arg)
		timeRegex := regexp.MustCompile(`^([*]|[0-9]+)[:]([*]|[0-9]+)(?:[:]([*]|[0-9]+)(?:[.]([*]|[0-9]+))?)?$`)

		if !timeRegex.MatchString(arg) {
			return nil, errors.New("invalid time pattern")
		}
		// make a list of the contents from matching patterns
		match := timeRegex.FindStringSubmatch(arg)[1:]
		if len(match) == 0 {
			return nil, errors.New("unmatched")
		}

		var tupList []int
		for _, s := range match {
			if s == "*" {
				tupList = append(tupList, 255)
			} else if s == "" {
				if slices.Contains(match, "*") {
					tupList = append(tupList, 255)
				} else {
					tupList = append(tupList, 0)
				}
			} else {
				i, _ := strconv.Atoi(s)
				tupList = append(tupList, i)
			}
		}
		if tupList[3] != 0xff {
			tupList[3] *= 10
		}
		d.value = TimeTuple{tupList[0], tupList[1], tupList[2], tupList[3]}
	case time.Duration:
		d.value = TimeTuple{int(arg.Hours()), int(arg.Minutes()), int(arg.Seconds()), int(arg.Milliseconds() * 10)}
	case *Time:
		d.value = arg.value
	case float32:
		d.now(arg)
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", arg)
	}

	return d, nil
}

func (t *Time) now(arg float32) {
	panic("implement me") // TODO
}

func (t *Time) Encode(tag Tag) {
	tag.setAppData(uint(model.BACnetDataType_TIME), []byte{byte(t.value.Hour), byte(t.value.Minute), byte(t.value.Second), byte(t.value.Hundredth)})
}

func (t *Time) Decode(tag Tag) error {
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(model.BACnetDataType_TIME) {
		return errors.New("Time application tag required")
	}
	if len(tag.GetTagData()) != 4 {
		return errors.New("invalid tag length")
	}

	tagData := tag.GetTagData()

	t.value = TimeTuple{int(tagData[0]), int(tagData[1]), int(tagData[2]), int(tagData[3])}
	return nil
}

func (t *Time) IsValid(arg any) bool {
	_, ok := arg.(bool)
	return ok
}

func (t *Time) Compare(other any) int {
	switch other := other.(type) {
	case *Time:
		_ = other // TODO: implement
		return -1
	default:
		return -1
	}
}

func (t *Time) LowerThan(other any) bool {
	switch other := other.(type) {
	case *Time:
		_ = other // TODO: implement
		return false
	default:
		return false
	}
}

func (t *Time) Equals(other any) bool {
	return t.value == other
}

func (t *Time) GetValue() TimeTuple {
	return t.value
}

func (t *Time) Coerce(arg Time) TimeTuple {
	return arg.GetValue()
}

func (t *Time) String() string {
	// rip it apart
	hour, minute, second, hundredth := t.value.Hour, t.value.Minute, t.value.Second, t.value.Hundredth

	rslt := "Time("
	if hour == 255 {
		rslt += "*:"
	} else {
		rslt += fmt.Sprintf("%02d:", hour)
	}
	if minute == 255 {
		rslt += "*:"
	} else {
		rslt += fmt.Sprintf("%02d:", minute)
	}
	if second == 255 {
		rslt += "*."
	} else {
		rslt += fmt.Sprintf("%02d.", second)
	}
	if hundredth == 255 {
		rslt += "*)"
	} else {
		rslt += fmt.Sprintf("%02d)", hundredth)
	}
	return rslt
}

// ObjectTypeContract provides a set of functions which can be overwritten by a sub struct
type ObjectTypeContract interface {
	EnumeratedContract
	// SetObjectType is required because we do more stuff in the constructor and can't wait for the substruct to finish that
	SetObjectType(objectType *ObjectType)
}

type ObjectType struct {
	*Enumerated

	enumerations map[string]uint64
}

func NewObjectType(args Args) (*ObjectType, error) {
	o := &ObjectType{
		enumerations: map[string]uint64{
			"accessCredential":      32,
			"accessDoor":            30,
			"accessPoint":           33,
			"accessRights":          34,
			"accessUser":            35,
			"accessZone":            36,
			"accumulator":           23,
			"alertEnrollment":       52,
			"analogInput":           0,
			"analogOutput":          1,
			"analogValue":           2,
			"auditLog":              61,
			"auditReporter":         62,
			"averaging":             18,
			"binaryInput":           3,
			"binaryLightingOutput":  55,
			"binaryOutput":          4,
			"binaryValue":           5,
			"bitstringValue":        39,
			"calendar":              6,
			"channel":               53,
			"characterstringValue":  40,
			"command":               7,
			"credentialDataInput":   37,
			"datePatternValue":      41,
			"dateValue":             42,
			"datetimePatternValue":  43,
			"datetimeValue":         44,
			"device":                8,
			"elevatorGroup":         57,
			"escalator":             58,
			"eventEnrollment":       9,
			"eventLog":              25,
			"file":                  10,
			"globalGroup":           26,
			"group":                 11,
			"integerValue":          45,
			"largeAnalogValue":      46,
			"lifeSafetyPoint":       21,
			"lifeSafetyZone":        22,
			"lift":                  59,
			"lightingOutput":        54,
			"loadControl":           28,
			"loop":                  12,
			"multiStateInput":       13,
			"multiStateOutput":      14,
			"multiStateValue":       19,
			"networkSecurity":       38,
			"networkPort":           56,
			"notificationClass":     15,
			"notificationForwarder": 51,
			"octetstringValue":      47,
			"positiveIntegerValue":  48,
			"program":               16,
			"pulseConverter":        24,
			"schedule":              17,
			"staging":               60,
			"structuredView":        29,
			"timePatternValue":      49,
			"timeValue":             50,
			"timer":                 31,
			"trendLog":              20,
			"trendLogMultiple":      27,
		},
	}
	var enumeratedContract EnumeratedContract = o
	var err error
	var arg0 any = 0
	switch len(args) {
	case 1:
		arg0 = args[0]
		switch arg0 := arg0.(type) {
		case *ObjectType:
			o.Enumerated, _ = NewEnumerated(arg0.Enumerated)
			for k, v := range arg0.enumerations {
				o.enumerations[k] = v
			}
			return o, nil
		}
	case 2:
		switch arg := args[0].(type) {
		case ObjectTypeContract:
			arg.SetObjectType(o)
			enumeratedContract = arg
			argEnumerations := arg.GetEnumerations()
			for k, v := range o.enumerations {
				if _, ok := argEnumerations[k]; !ok {
					argEnumerations[k] = v
				}
			}
			o.enumerations = nil // supper seeded
		default:
			return nil, fmt.Errorf("invalid arg type: %T", arg)
		}
		arg0 = args[1]
	}
	o.Enumerated, err = NewEnumerated(enumeratedContract, arg0)
	if err != nil {
		return nil, errors.Wrap(err, "error creating enumerated")
	}
	return o, nil
}

func (o *ObjectType) GetEnumerations() map[string]uint64 {
	return o.enumerations
}

func (o *ObjectType) SetEnumerated(enumerated *Enumerated) {
	o.Enumerated = enumerated
}

func (o *ObjectType) SetObjectType(_ *ObjectType) {
	panic("must be implemented by substruct")
}

func (o *ObjectType) String() string {
	value := strconv.Itoa(int(o.value))
	if o.valueString != "" {
		value = o.valueString
	}
	return fmt.Sprintf("ObjectType(%v)", value)
}

type ObjectIdentifierTuple struct {
	Left  any
	Right int
}

type ObjectIdentifier struct {
	//*Atomic[...] won't work here

	objectTypeClass *ObjectType

	value ObjectIdentifierTuple
}

func NewObjectIdentifier(args Args) (*ObjectIdentifier, error) {
	i := &ObjectIdentifier{}
	var err error
	i.objectTypeClass, err = NewObjectType(nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating object type")
	}
	i.value = ObjectIdentifierTuple{"analogInput", 0}

	if len(args) == 0 || args == nil {
		return i, nil
	}
	if len(args) == 1 {
		arg := args[0]
		switch arg := arg.(type) {
		case *tag:
			err := i.Decode(arg)
			if err != nil {
				return nil, errors.Wrap(err, "error decoding")
			}
			return i, nil
		case int:
			i.setLong(arg)
		case string:
			split := strings.Split(arg, ":")
			var objType, objInstance any = split[0], split[1]
			if objTypeInt, err := strconv.Atoi(fmt.Sprintf("%v", objType)); err == nil {
				objType = objTypeInt
			}
			var err error
			objInstance, err = strconv.Atoi(fmt.Sprintf("%v", objInstance))
			if err != nil {
				return nil, errors.Wrap(err, "error parsing instance")
			}
			if err := i.setTuple(objType, objInstance.(int)); err != nil {
				return nil, errors.Wrap(err, "can't set tuple")
			}
		case ObjectIdentifierTuple:
			if err := i.setTuple(arg.Left, arg.Right); err != nil {
				return nil, errors.Wrap(err, "error setting tuple")
			}
		case *ObjectIdentifier:
			i.value = arg.value
		default:
			return nil, errors.Errorf("invalid constructor datatype: %T", arg)
		}
	} else if len(args) == 2 {
		err := i.setTuple(args[0], args[1].(int))
		if err != nil {
			return nil, errors.Wrap(err, "error setting tuple")
		}
	} else {
		return nil, errors.New("invalid constructor parameters")
	}

	return i, nil
}

func (o *ObjectIdentifier) setTuple(objType any, objInstance int) error {
	switch objType.(type) {
	case int:
		if gotObjType, ok := o.objectTypeClass.GetXlateTable()[uint64(objType.(int))]; ok {
			objType = gotObjType
		}
	case string:
		if _, ok := o.objectTypeClass.GetXlateTable()[objType]; !ok {
			return errors.Errorf("unrecognized object type %s", objType)
		}
	default:
		return errors.Errorf("invalid datatype for object type: %T", objType)
	}

	// check valid instance number
	if objInstance < 0 || objInstance > 0x003FFFFF {
		return errors.Errorf("invalid object instance out of range: %d", objInstance)
	}

	o.value = ObjectIdentifierTuple{objType, objInstance}
	return nil
}

func (o *ObjectIdentifier) getTuple() ObjectIdentifierTuple {
	return o.value
}

func (o *ObjectIdentifier) setLong(value int) {
	// suck out the type
	objTypeInt := (value >> 22) & 0x3ff
	var objType any = objTypeInt

	// try and make it pretty
	if item, ok := o.objectTypeClass.GetItem(uint64(objTypeInt)); ok {
		objType = item
	}

	// suck out the instance
	objInstance := value & 0x003FFFFF

	// save the result
	o.value = ObjectIdentifierTuple{objType, objInstance}
}

func (o *ObjectIdentifier) getLong() int {
	tuple := o.getTuple()
	objType, objInstance := tuple.Left, tuple.Right

	if _, ok := objType.(string); ok {
		if objTypeGot, ok := o.objectTypeClass.GetXlateTable()[objType]; ok {
			objType = int(objTypeGot.(uint64))
		}
	}

	return (objType.(int) << 22) + objInstance
}

func (o *ObjectIdentifier) Encode(tag Tag) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(o.getLong()))
	tag.setAppData(uint(model.BACnetDataType_BACNET_OBJECT_IDENTIFIER), data)
}

func (o *ObjectIdentifier) Decode(tag Tag) error {
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(model.BACnetDataType_BACNET_OBJECT_IDENTIFIER) {
		return errors.New("ObjectIdentifier application tag required")
	}
	if len(tag.GetTagData()) != 4 {
		return errors.New("invalid tag length")
	}

	tagData := tag.GetTagData()

	o.setLong(int(binary.BigEndian.Uint32(tagData)))
	return nil
}

func (o *ObjectIdentifier) IsValid(arg any) bool {
	switch arg.(type) {
	case ObjectIdentifier:
		return true
	default:
		return false
	}
}

func (o *ObjectIdentifier) Compare(other any) int {
	switch other := other.(type) {
	case *ObjectIdentifier:
		return o.getLong() - other.getLong()
	default:
		return -1
	}
}

func (o *ObjectIdentifier) LowerThan(other any) bool {
	switch other := other.(type) {
	case *ObjectIdentifier:
		return o.getLong() < other.getLong()
	default:
		return false
	}
}

func (o *ObjectIdentifier) Equals(other any) bool {
	return o.value == other
}

func (o *ObjectIdentifier) GetValue() ObjectIdentifierTuple {
	return o.value
}

func (o *ObjectIdentifier) Coerce(arg ObjectIdentifier) ObjectIdentifierTuple {
	return arg.GetValue()
}

func (o *ObjectIdentifier) String() string {
	// rip it apart
	objType, objInstance := o.value.Left, o.value.Right

	var objTypeAsUint64 uint64
	if objTypeAsInt, ok := objType.(int); ok {
		objTypeAsUint64 = uint64(objTypeAsInt)
	}
	var typeString string
	if s, ok := objType.(string); ok {
		typeString = s
	} else if i, intOk := objType.(int); intOk && i < 0 {
		typeString = fmt.Sprintf("Bad %d", i)
	} else if gotType, xlateOk := o.objectTypeClass.GetXlateTable()[objTypeAsUint64]; xlateOk {
		typeString = gotType.(string)
	} else if intOk && i < 128 {
		typeString = fmt.Sprintf("Reserved %d", i)
	} else {
		typeString = fmt.Sprintf("Vendor %s", objType)
	}

	return fmt.Sprintf("ObjectIdentifier(%s,%d)", typeString, objInstance)
}
