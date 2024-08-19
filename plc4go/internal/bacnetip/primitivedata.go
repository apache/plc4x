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

type Tag struct {
	tagClass  model.TagClass
	tagNumber uint
	tagLVT    int
	tagData   []byte
}

func NewTag(args Args) (*Tag, error) {
	t := &Tag{}
	if len(args) == 0 {
		return t, nil
	}
	if len(args) == 1 {
		t.decode(args[0])
	} else if len(args) >= 2 {
		t.set(args)
	} else {
		return nil, errors.New("invalid arguments")
	}
	return t, nil
}

func (t *Tag) decode(arg any) {
	panic("implement me")
}

func (t *Tag) set(args Args) {
	switch arg0 := args[0].(type) {
	case model.TagClass:
		t.tagClass = arg0
	case uint:
		t.tagClass = model.TagClass(arg0)
	case int:
		t.tagClass = model.TagClass(arg0)
	default:
		panic("oh no")
	}
	switch arg1 := args[1].(type) {
	case model.BACnetDataType:
		t.tagNumber = uint(arg1)
	case uint:
		t.tagNumber = arg1
	case int:
		t.tagNumber = uint(arg1)
	default:
		panic("oh no")
	}
	if len(args) == 2 {
		return
	}
	t.tagLVT = args[2].(int)
	if len(args) == 3 {
		return
	}
	t.tagData = args[3].([]byte)
}

func (t *Tag) setAppData(tagNumber uint, tdata []byte) {
	t.tagClass = model.TagClass_APPLICATION_TAGS
	t.tagNumber = tagNumber
	t.tagLVT = len(tdata)
	t.tagData = tdata
}

func (t *Tag) GetTagClass() model.TagClass {
	return t.tagClass
}

func (t *Tag) GetTagNumber() uint {
	return t.tagNumber
}

func (t *Tag) GetTagData() []byte {
	return t.tagData
}

func (t *Tag) Equals(other any) bool {
	if t == nil && other == nil {
		return true
	}
	if other == nil {
		return false
	}
	otherTag, ok := other.(*Tag)
	if !ok {
		return false
	}
	return t.tagClass == otherTag.tagClass &&
		t.tagNumber == otherTag.tagNumber &&
		t.tagLVT == otherTag.tagLVT &&
		bytes.Equal(t.tagData, otherTag.tagData)
}

type ApplicationTag struct {
	model.BACnetApplicationTag
	// TODO: implement me
}

type ContextTag struct {
	model.BACnetContextTag
	// TODO: implement me
}

type OpeningTag struct {
	// TODO: implement me
}

type ClosingTag struct {
	// TODO: implement me
}

type TagList struct {
	// TODO: implement me
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

// TODO: finish
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
	case *Tag:
		err := b.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return b, nil
	case bool:
		if arg {
			b.value = 1
		}
	case *Null:
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

func (b *Null) Encode(tag *Tag) {
	tag.set(NewArgs(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_NULL, b.value, []byte{}))
}

func (b *Null) Decode(tag *Tag) error {
	if tag.tagClass != model.TagClass_APPLICATION_TAGS || tag.tagNumber != uint(model.BACnetDataType_NULL) {
		return errors.New("Null application tag required")
	}
	if tag.tagLVT > 1 {
		return errors.New("invalid tag value")
	}

	// get the data
	if tag.tagLVT == 1 {
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
	case *Tag:
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

func (b *Boolean) Encode(tag *Tag) {
	tag.set(NewArgs(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_BOOLEAN, b.value, []byte{}))
}

func (b *Boolean) Decode(tag *Tag) error {
	if tag.tagClass != model.TagClass_APPLICATION_TAGS || tag.tagNumber != uint(model.BACnetDataType_BOOLEAN) {
		return errors.New("boolean application tag required")
	}
	if tag.tagLVT > 1 {
		return errors.New("invalid tag value")
	}

	// get the data
	if tag.tagLVT == 1 {
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

// TODO: finish
type Unsigned struct {
	*Atomic[uint]
	*CommonMath
}

func NewUnsigned(arg Arg) (*Unsigned, error) {
	b := &Unsigned{}
	b.Atomic = NewAtomic[uint](b)

	if arg == nil {
		return b, nil
	}
	switch arg := arg.(type) {
	case *Tag:
		err := b.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return b, nil
	case bool:
		if arg {
			b.value = 1
		}
	case *Unsigned:
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

func (b *Unsigned) Encode(tag *Tag) {
	tag.set(NewArgs(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_UNSIGNED_INTEGER, b.value, []byte{}))
}

func (b *Unsigned) Decode(tag *Tag) error {
	if tag.tagClass != model.TagClass_APPLICATION_TAGS || tag.tagNumber != uint(model.BACnetDataType_UNSIGNED_INTEGER) {
		return errors.New("Unsigned application tag required")
	}
	if tag.tagLVT > 1 {
		return errors.New("invalid tag value")
	}

	// get the data
	if tag.tagLVT == 1 {
		b.value = 1
	}
	return nil
}

func (b *Unsigned) IsValid(arg any) bool {
	_, ok := arg.(bool)
	return ok
}

func (b *Unsigned) String() string {
	value := "False"
	if b.value == 1 {
		value = "True"
	}
	return fmt.Sprintf("Unsigned(%s)", value)
}

// TODO: finish
type Unsigned8 struct {
	*Atomic[uint8]
}

func NewUnsigned8(arg Arg) (*Unsigned8, error) {
	b := &Unsigned8{}
	b.Atomic = NewAtomic[uint8](b)

	if arg == nil {
		return b, nil
	}
	switch arg := arg.(type) {
	case *Tag:
		err := b.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return b, nil
	case bool:
		if arg {
			b.value = 1
		}
	case *Unsigned8:
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

func (b *Unsigned8) Encode(tag *Tag) {
	tag.set(NewArgs(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_UNSIGNED_INTEGER, b.value, []byte{}))
}

func (b *Unsigned8) Decode(tag *Tag) error {
	if tag.tagClass != model.TagClass_APPLICATION_TAGS || tag.tagNumber != uint(model.BACnetDataType_UNSIGNED_INTEGER) {
		return errors.New("Unsigned8 application tag required")
	}
	if tag.tagLVT > 1 {
		return errors.New("invalid tag value")
	}

	// get the data
	if tag.tagLVT == 1 {
		b.value = 1
	}
	return nil
}

func (b *Unsigned8) IsValid(arg any) bool {
	_, ok := arg.(bool)
	return ok
}

func (b *Unsigned8) String() string {
	value := "False"
	if b.value == 1 {
		value = "True"
	}
	return fmt.Sprintf("Unsigned8(%s)", value)
}

// TODO: finish
type Unsigned16 struct {
	*Atomic[uint16]
}

func NewUnsigned16(arg Arg) (*Unsigned16, error) {
	b := &Unsigned16{}
	b.Atomic = NewAtomic[uint16](b)

	if arg == nil {
		return b, nil
	}
	switch arg := arg.(type) {
	case *Tag:
		err := b.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return b, nil
	case bool:
		if arg {
			b.value = 1
		}
	case *Unsigned16:
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

func (b *Unsigned16) Encode(tag *Tag) {
	tag.set(NewArgs(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_UNSIGNED_INTEGER, b.value, []byte{}))
}

func (b *Unsigned16) Decode(tag *Tag) error {
	if tag.tagClass != model.TagClass_APPLICATION_TAGS || tag.tagNumber != uint(model.BACnetDataType_UNSIGNED_INTEGER) {
		return errors.New("Unsigned16 application tag required")
	}
	if tag.tagLVT > 1 {
		return errors.New("invalid tag value")
	}

	// get the data
	if tag.tagLVT == 1 {
		b.value = 1
	}
	return nil
}

func (b *Unsigned16) IsValid(arg any) bool {
	_, ok := arg.(bool)
	return ok
}

func (b *Unsigned16) String() string {
	value := "False"
	if b.value == 1 {
		value = "True"
	}
	return fmt.Sprintf("Unsigned16(%s)", value)
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
	case *Tag:
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

func (i *Integer) Encode(tag *Tag) {
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

func (i *Integer) Decode(tag *Tag) error {
	if tag.tagClass != model.TagClass_APPLICATION_TAGS || tag.tagNumber != uint(model.BACnetDataType_SIGNED_INTEGER) {
		return errors.New("Integer application tag required")
	}
	if len(tag.tagData) == 0 {
		return errors.New("invalid tag length")
	}

	tagData := tag.tagData

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

// TODO: finish
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
	case *Tag:
		err := b.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return b, nil
	case bool:
		if arg {
			b.value = 1
		}
	case *Real:
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

func (b *Real) Encode(tag *Tag) {
	tag.set(NewArgs(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_REAL, b.value, []byte{}))
}

func (b *Real) Decode(tag *Tag) error {
	if tag.tagClass != model.TagClass_APPLICATION_TAGS || tag.tagNumber != uint(model.BACnetDataType_REAL) {
		return errors.New("Real application tag required")
	}
	if tag.tagLVT > 1 {
		return errors.New("invalid tag value")
	}

	// get the data
	if tag.tagLVT == 1 {
		b.value = 1
	}
	return nil
}

func (b *Real) IsValid(arg any) bool {
	_, ok := arg.(bool)
	return ok
}

func (b *Real) String() string {
	value := "False"
	if b.value == 1 {
		value = "True"
	}
	return fmt.Sprintf("Real(%s)", value)
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
	case *Tag:
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

func (d *Double) Encode(tag *Tag) {
	var _b = make([]byte, 8)
	binary.BigEndian.PutUint64(_b, math.Float64bits(d.value))
	tag.setAppData(uint(model.BACnetDataType_DOUBLE), _b)
}

func (d *Double) Decode(tag *Tag) error {
	if tag.tagClass != model.TagClass_APPLICATION_TAGS || tag.tagNumber != uint(model.BACnetDataType_DOUBLE) {
		return errors.New("Double application tag required")
	}
	if len(tag.tagData) != 8 {
		return errors.New("invalid tag length")
	}

	// extract the data
	d.value = math.Float64frombits(binary.BigEndian.Uint64(tag.tagData))
	return nil
}

func (d *Double) IsValid(arg any) bool {
	_, ok := arg.(float64)
	return ok
}

func (d *Double) String() string {
	return fmt.Sprintf("Double(%g)", d.value)
}

// TODO: finish
type OctetString struct {
	*Atomic[string]
	*CommonMath
}

func NewOctetString(arg Arg) (*OctetString, error) {
	b := &OctetString{}
	b.Atomic = NewAtomic[string](b)

	if arg == nil {
		return b, nil
	}
	switch arg := arg.(type) {
	case *Tag:
		err := b.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return b, nil
	case bool:
		if arg {
			b.value = "1"
		}
	case *OctetString:
		b.value = arg.value
	case string:
		switch arg {
		case "True", "true":
			b.value = "1"
		case "False", "false":
		default:
			return nil, errors.Errorf("invalid string: %s", arg)
		}
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", arg)
	}

	return b, nil
}

func (b *OctetString) Encode(tag *Tag) {
	tag.set(NewArgs(model.TagClass_APPLICATION_TAGS, model.BACnetDataType_OCTET_STRING, b.value, []byte{}))
}

func (b *OctetString) Decode(tag *Tag) error {
	if tag.tagClass != model.TagClass_APPLICATION_TAGS || tag.tagNumber != uint(model.BACnetDataType_OCTET_STRING) {
		return errors.New("OctetString application tag required")
	}
	if tag.tagLVT > 1 {
		return errors.New("invalid tag value")
	}

	// get the data
	if tag.tagLVT == 1 {
		b.value = "1"
	}
	return nil
}

func (b *OctetString) IsValid(arg any) bool {
	_, ok := arg.(bool)
	return ok
}

func (b *OctetString) String() string {
	value := "False"
	if b.value == "1" {
		value = "True"
	}
	return fmt.Sprintf("OctetString(%s)", value)
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
	case *Tag:
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

func (c *CharacterString) Encode(tag *Tag) {
	tag.setAppData(uint(model.BACnetDataType_CHARACTER_STRING), append([]byte{c.strEncoding}, c.strValue...))
}

func (c *CharacterString) Decode(tag *Tag) error {
	if tag.tagClass != model.TagClass_APPLICATION_TAGS || tag.tagNumber != uint(model.BACnetDataType_CHARACTER_STRING) {
		return errors.New("CharacterString application tag required")
	}
	if len(tag.tagData) == 0 {
		return errors.New("invalid tag length")
	}

	tagData := tag.tagData

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
	return fmt.Sprintf("CharacterString(%d,X'%s')", c.strEncoding, Btox(c.strValue))
}

// BitStringExtension can be used to inherit from BitString
type BitStringExtension interface {
	fmt.Stringer
	GetBitNames() map[string]int
	GetBitLen() int
}

type BitString struct {
	*Atomic[int] // TODO: implement properly

	bitStringExtension BitStringExtension
	Value              []bool
}

func NewBitString(arg ...any) (*BitString, error) {
	return NewBitStringWithExtension(nil, arg...)
}

func NewBitStringWithExtension(bitStringExtension BitStringExtension, arg ...any) (*BitString, error) {
	b := &BitString{
		bitStringExtension: bitStringExtension,
	}
	if len(arg) == 0 {
		return b, nil
	}
	if len(arg) > 1 {
		return nil, errors.New("too many arguments")
	}
	if bitStringExtension != nil {
		b.Value = make([]bool, bitStringExtension.GetBitLen())
	}
	switch arg := arg[0].(type) {
	case *Tag:
		err := b.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "decoding tag failed")
		}
	case []int:
		b.Value = make([]bool, len(arg))
		for i, v := range arg {
			b.Value[i] = v != 0
		}
	case []string:
		bitNames := make(map[string]int)
		if bitStringExtension != nil {
			bitNames = bitStringExtension.GetBitNames()
		}
		for _, bit := range arg {
			bit, ok := bitNames[bit]
			if !ok || bit < 0 || bit > len(b.Value) {
				return nil, errors.New("constructorElement out of range")
			}
			b.Value[bit] = true
		}
	case *BitString:
		b.Value = arg.Value[:]
	case model.BACnetApplicationTagBitStringExactly:
		b.Value = arg.GetPayload().GetData()
	default:
		return nil, errors.Errorf("no support for %T yet", arg)
	}
	return b, nil
}

func (b *BitString) Decode(tag *Tag) error {
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(TagBitStringAppTag) {
		return errors.New("bit string application tag required")
	}
	if len(tag.GetTagData()) == 0 {
		return errors.New("invalid tag length")
	}
	// extract the number of unused bits
	unused := tag.tagData[0]

	// extract the data
	data := make([]bool, 0)
	for _, x := range tag.tagData[1:] {
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
		b.Value = data[:len(data)-int(unused)]
	} else {
		b.Value = data
	}
	return nil
}

func (b *BitString) Encode(tag *Tag) {
	used := len(b.Value) % 8
	unused := 8 - used
	if unused == 8 {
		unused = 0
	}

	// start with the number of unused bits
	data := []byte{byte(unused)}

	// build and append each packed octet
	bits := append(b.Value, make([]bool, unused)...)
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
	for index, value := range b.Value {
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

	_enumerations map[string]uint64
	_xlateTable   map[any]any

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
		return nil, errors.New("invalid argument")
	}
	if len(e.EnumeratedContract.GetXlateTable()) == 0 {
		expandEnumerations(e.EnumeratedContract)
	}

	switch arg := arg.(type) {
	case *Tag:
		err := e.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return e, nil
	case uint:
		e.value = uint64(arg)
	case int:
		if arg < 0 {
			return nil, errors.New("arg must be positive")
		}
		e.value = uint64(arg)
	case uint64:
		e.value = arg

		// convert it to a string if you can
		e.valueString, _ = e.EnumeratedContract.GetXlateTable()[arg].(string)
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
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", arg)
	}

	return e, nil
}

func (e *Enumerated) GetEnumerations() map[string]uint64 {
	if e._enumerations == nil {
		e._enumerations = make(map[string]uint64)
	}
	return e._enumerations
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

func (e *Enumerated) GetItem(item any) any {
	return e.EnumeratedContract.GetXlateTable()[item]
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

func (e *Enumerated) Decode(tag *Tag) error {
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(model.BACnetDataType_ENUMERATED) {
		return errors.New("bit string application tag required")
	}
	if len(tag.GetTagData()) == 0 {
		return errors.New("invalid tag length")
	}

	data := tag.tagData
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

func (e *Enumerated) Encode(tag *Tag) {
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
	*Atomic[int64]

	year      int
	month     int
	day       int
	dayOfWeek int
}

func NewDate(arg Arg, args Args) (*Date, error) {
	d := &Date{}
	d.Atomic = NewAtomic[int64](d)
	year := 255
	if len(args) > 0 {
		year = args[0].(int)
	}
	if year >= 1900 {
		year = year - 1900
	}
	d.year = year
	month := 0xff
	if len(args) > 1 {
		month = args[1].(int)
	}
	d.month = month
	day := 0xff
	if len(args) > 2 {
		day = args[2].(int)
	}
	d.day = day
	dayOfWeek := 0xff
	if len(args) > 3 {
		dayOfWeek = args[3].(int)
	}
	d.dayOfWeek = dayOfWeek

	if arg == nil {
		return d, nil
	}
	switch arg := arg.(type) {
	case *Tag:
		err := d.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return d, nil
	case DateTuple:
		d.year, d.month, d.day, d.dayOfWeek = arg.Year, arg.Month, arg.Day, arg.DayOfWeek
		var tempTime time.Time
		tempTime.AddDate(d.year, d.month, d.day)
		d.value = tempTime.UnixNano() - (time.Time{}.UnixNano()) // TODO: check this
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
		d.year = year
		d.month = month
		d.day = day
		d.dayOfWeek = dayOfWeek

		var tempTime time.Time
		tempTime.AddDate(year, month, day)
		d.value = tempTime.UnixNano() - (time.Time{}.UnixNano()) // TODO: check this

		// calculate the day of the week
		if dayOfWeek == 0 {
			d.calcDayOfWeek()
		}
	case *Date:
		d.value = arg.value
		d.year = arg.year
		d.month = arg.month
		d.day = arg.day
		d.dayOfWeek = arg.dayOfWeek
	case float32:
		d.now(arg)
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", arg)
	}

	return d, nil
}

func (d *Date) GetTupleValue() (year int, month int, day int, dayOfWeek int) {
	return d.year, d.month, d.day, d.dayOfWeek
}

func (d *Date) calcDayOfWeek() {
	year, month, day, dayOfWeek := d.year, d.month, d.day, d.dayOfWeek

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
		panic(today) // TODO: implement me
	}

	// put it back together
	d.year = year
	d.month = month
	d.day = day
	d.dayOfWeek = dayOfWeek
}

func (d *Date) now(arg float32) {
	panic("implement me") // TODO
}

func (d *Date) Encode(tag *Tag) {
	var b []byte
	binary.BigEndian.AppendUint64(b, uint64(d.value))
	tag.setAppData(uint(model.BACnetDataType_DATE), b)
}

func (d *Date) Decode(tag *Tag) error {
	if tag.tagClass != model.TagClass_APPLICATION_TAGS || tag.tagNumber != uint(model.BACnetDataType_DATE) {
		return errors.New("Date application tag required")
	}
	if len(tag.tagData) != 4 {
		return errors.New("invalid tag length")
	}

	arg := tag.tagData
	year, month, day, dayOfWeek := arg[0], arg[1], arg[2], arg[3]
	var tempTime time.Time
	tempTime.AddDate(int(year), int(month), int(day))
	d.value = tempTime.UnixNano() - (time.Time{}.UnixNano()) // TODO: check this
	d.year, d.month, d.day, d.dayOfWeek = int(year), int(month), int(day), int(dayOfWeek)
	return nil
}

func (d *Date) IsValid(arg any) bool {
	_, ok := arg.(bool)
	return ok
}

func (d *Date) String() string {
	year, month, day, dayOfWeek := d.year, d.month, d.day, d.dayOfWeek
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
