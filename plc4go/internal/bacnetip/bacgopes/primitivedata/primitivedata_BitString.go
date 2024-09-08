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

package primitivedata

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

// BitStringExtension can be used to inherit from BitString
type BitStringExtension interface {
	fmt.Stringer
	GetBitNames() map[string]int
	GetBitLen() int
}

type BitString struct {
	_appTag model.BACnetDataType

	bitStringExtension BitStringExtension
	value              []bool
}

func NewBitString(args Args) (*BitString, error) {
	return NewBitStringWithExtension(nil, args)
}

func NewBitStringWithExtension(bitStringExtension BitStringExtension, args Args) (*BitString, error) {
	b := &BitString{
		_appTag:            model.BACnetDataType_BIT_STRING,
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
	case model.BACnetApplicationTagBitString:
		b.value = arg.GetPayload().GetData()
	default:
		return nil, errors.Errorf("no support for %T yet", arg)
	}
	return b, nil
}

func (b *BitString) GetAppTag() model.BACnetDataType {
	return b._appTag
}

func (b *BitString) Decode(arg Arg) error {
	tag, ok := arg.(Tag)
	if !ok {
		return errors.Errorf("%T is not a Tag", arg)
	}
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(b._appTag) {
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

func (b *BitString) Encode(arg Arg) error {
	tag, ok := arg.(Tag)
	if !ok {
		return errors.Errorf("%T is not a Tag", arg)
	}
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

	tag.setAppData(uint(b._appTag), data)
	return nil
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
