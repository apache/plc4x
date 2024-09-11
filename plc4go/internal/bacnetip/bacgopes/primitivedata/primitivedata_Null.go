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

type Null struct {
	*Atomic[int]

	_appTag model.BACnetDataType
}

func NewNull(arg Arg) (*Null, error) {
	b := &Null{
		_appTag: model.BACnetDataType_NULL,
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
	case *Null:
		b.value = arg.value
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", arg)
	}

	return b, nil
}

func (n *Null) GetAppTag() model.BACnetDataType {
	return n._appTag
}

func (n *Null) Encode(arg Arg) error {
	tag, ok := arg.(Tag)
	if !ok {
		return errors.Errorf("%T is not a Tag", arg)
	}
	tag.set(NA(model.TagClass_APPLICATION_TAGS, n._appTag, n.value, []byte{}))
	return nil
}

func (n *Null) Decode(arg Arg) error {
	tag, ok := arg.(Tag)
	if !ok {
		return errors.Errorf("%T is not a Tag", arg)
	}
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(n._appTag) {
		return errors.New("Null application tag required")
	}
	if tag.GetTagLvt() > 1 {
		return errors.New("invalid tag value")
	}

	// get the data
	if tag.GetTagLvt() == 1 {
		n.value = 1
	}
	return nil
}

func (n *Null) IsValid(arg any) bool {
	_, ok := arg.(bool)
	return ok
}

func (n *Null) String() string {
	value := "False"
	if n.value == 1 {
		value = "True"
	}
	return fmt.Sprintf("Null(%s)", value)
}
