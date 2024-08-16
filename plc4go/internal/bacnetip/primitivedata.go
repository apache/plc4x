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
	"fmt"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/pkg/errors"
	"strings"
)

const (
	// Deprecated: use model.TagClass_APPLICATION_TAGS
	TagApplicationTagClass = model.TagClass_APPLICATION_TAGS
	// Deprecated: use model.TagClass_CONTEXT_SPECIFIC_TAGS
	TagContextTagClass = model.TagClass_CONTEXT_SPECIFIC_TAGS
	TagOpeningTagClass = 2
	TagClosingTagClass = 3

	// Deprecated: use  model.BACnetDataType_NULL
	Tag_nullAppTag = model.BACnetDataType_NULL
	// Deprecated: use  model.BACnetDataType_BOOLEAN
	Tag_booleanAppTag = model.BACnetDataType_BOOLEAN
	// Deprecated: use  model.BACnetDataType_UNSIGNED_INTEGER
	Tag_unsignedAppTag = model.BACnetDataType_UNSIGNED_INTEGER
	// Deprecated: use  model.BACnetDataType_SIGNED_INTEGER
	Tag_integerAppTag = model.BACnetDataType_SIGNED_INTEGER
	// Deprecated: use  model.BACnetDataType_REAL
	Tag_realAppTag = model.BACnetDataType_REAL
	// Deprecated: use  model.BACnetDataType_DOUBLE
	Tag_doubleAppTag = model.BACnetDataType_DOUBLE
	// Deprecated: use  model.BACnetDataType_OCTET_STRING
	Tag_octetStringAppTag = model.BACnetDataType_OCTET_STRING
	// Deprecated: use  model.BACnetDataType_CHARACTER_STRING
	Tag_characterStringAppTag = model.BACnetDataType_CHARACTER_STRING
	// Deprecated: use  model.BACnetDataType_BIT_STRING
	Tag_bitStringAppTag = model.BACnetDataType_BIT_STRING
	// Deprecated: use  model.BACnetDataType_ENUMERATED
	Tag_enumeratedAppTag = model.BACnetDataType_ENUMERATED
	// Deprecated: use  model.BACnetDataType_DATE
	Tag_dateAppTag = model.BACnetDataType_DATE
	// Deprecated: use  model.BACnetDataType_TIME
	Tag_timeAppTag = model.BACnetDataType_TIME
	// Deprecated: use  model.BACnetDataType_BACNET_OBJECT_IDENTIFIER
	Tag_objectIdentifierAppTag = model.BACnetDataType_BACNET_OBJECT_IDENTIFIER
	Tag_reservedAppTag13       = 13
	Tag_reservedAppTag14       = 14
	Tag_reservedAppTag15       = 15
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
	t.tagClass = args[0].(model.TagClass)
	switch arg1 := args[1].(type) {
	case model.BACnetDataType:
		t.tagNumber = uint(arg1)
	case uint:
		t.tagNumber = arg1
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
}

type ContextTag struct {
	model.BACnetContextTag
}

// BitStringExtension can be used to inherit from BitString
type BitStringExtension interface {
	fmt.Stringer
	GetBitNames() map[string]int
	GetBitLen() int
}

type BitString struct {
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
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(Tag_bitStringAppTag) {
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
