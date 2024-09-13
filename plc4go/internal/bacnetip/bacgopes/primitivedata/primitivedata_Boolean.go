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

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type Boolean struct {
	*Atomic[int] //Note we need int as bool can't be used

	_appTag model.BACnetDataType
}

func NewBoolean(arg Arg) (*Boolean, error) {
	b := &Boolean{
		_appTag: model.BACnetDataType_BOOLEAN,
	}
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

func (b *Boolean) GetAppTag() model.BACnetDataType {
	return b._appTag
}

func (b *Boolean) Encode(arg Arg) error {
	tag, ok := arg.(Tag)
	if !ok {
		return errors.Errorf("%T is not a Tag", arg)
	} //TODO: move tag number into member variable
	tag.set(NA(model.TagClass_APPLICATION_TAGS, b._appTag, b.value, []byte{}))
	return nil
}

func (b *Boolean) Decode(arg Arg) error {
	tag, ok := arg.(Tag)
	if !ok {
		return errors.Errorf("%T is not a Tag", arg)
	}
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(b._appTag) {
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

// GetValue gives an int value because bool can't be used in constraint. GA convenience method GetBoolValue exists.
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
