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
	"fmt"

	"github.com/pkg/errors"

	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

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
