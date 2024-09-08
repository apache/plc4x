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
	"encoding/binary"
	"fmt"
	"math"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type Integer struct {
	*Atomic[int32]
	*CommonMath

	_appTag model.BACnetDataType
}

func NewInteger(arg Arg) (*Integer, error) {
	i := &Integer{
		_appTag: model.BACnetDataType_SIGNED_INTEGER,
	}
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

func (i *Integer) GetAppTag() model.BACnetDataType {
	return i._appTag
}

func (i *Integer) Encode(arg Arg) error {
	tag, ok := arg.(Tag)
	if !ok {
		return errors.Errorf("%T is not a Tag", arg)
	}
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

	tag.setAppData(uint(i._appTag), data)
	return nil
}

func (i *Integer) Decode(arg Arg) error {
	tag, ok := arg.(Tag)
	if !ok {
		return errors.Errorf("%T is not a Tag", arg)
	}
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(i._appTag) {
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
