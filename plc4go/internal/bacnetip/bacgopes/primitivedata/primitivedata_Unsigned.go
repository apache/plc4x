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
	"strconv"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type Unsigned struct {
	*Atomic[uint32]
	*CommonMath

	_appTag readWriteModel.BACnetDataType
}

func NewUnsigned(arg Arg) (*Unsigned, error) {
	i := &Unsigned{
		_appTag: readWriteModel.BACnetDataType_UNSIGNED_INTEGER,
	}
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
	case uint8:
		i.value = uint32(arg)
	case uint16:
		i.value = uint32(arg)
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
	case *readWriteModel.MaxApduLengthAccepted:
		if arg != nil {
			i.value = uint32(*arg)
		}
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", arg)
	}

	return i, nil
}

func (u *Unsigned) GetAppTag() readWriteModel.BACnetDataType {
	return u._appTag
}

func (u *Unsigned) Encode(arg Arg) error {
	tag, ok := arg.(Tag)
	if !ok {
		return errors.Errorf("%T is not a Tag", arg)
	}
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, u.value)

	// reduce the value to the smallest number of bytes
	for len(data) > 1 && data[0] == 0 {
		data = data[1:]
	}

	tag.setAppData(uint(u._appTag), data)
	return nil
}

func (u *Unsigned) Decode(arg Arg) error {
	tag, ok := arg.(Tag)
	if !ok {
		return errors.Errorf("%T is not a Tag", arg)
	}
	if tag.GetTagClass() != readWriteModel.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(u._appTag) {
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
	u.value = rslt
	return nil
}

func (u *Unsigned) IsValid(arg any) bool {
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

func (u *Unsigned) String() string {
	return fmt.Sprintf("Unsigned(%d)", u.value)
}
