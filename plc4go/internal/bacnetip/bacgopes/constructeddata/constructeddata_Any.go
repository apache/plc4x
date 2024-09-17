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

package constructeddata

import (
	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
)

type AnyContract interface {
	Encode(taglist TagList) error
	Decode(taglist TagList) error
	castIn(arg Arg) error
	castOut(arg Arg) error
}

type Any struct {
	_contract AnyContract // TODO: finish

	tagList TagList
}

func NewAny(args Args) (*Any, error) {
	a := &Any{}

	// cast the args
	for _, arg := range args {
		if err := a.castIn(arg); err != nil {
			return nil, errors.Wrapf(err, "error casting arg %v", arg)
		}
	}
	return a, nil
}

func WithAnyContract(contract AnyContract) func(*Any) {
	return func(a *Any) {
		a._contract = contract
	}
}

func (a *Any) Encode(tagList TagList) error {
	a.tagList.Extend(tagList.GetTagList()...)
	return nil
}

func (a *Any) Decode(tagList TagList) error {
	lvl := 0
	for len(tagList.GetTagList()) != 0 {
		tag := tagList.Peek()
		if tag.GetTagClass() == TagOpeningTagClass {
			lvl++
		} else if tag.GetTagClass() == TagClosingTagClass {
			lvl--
			if lvl < 0 {
				break
			}
		}
		a.tagList.Append(tagList.Pop())
	}

	// make sure everything balances
	if lvl > 0 {
		return errors.New("mismatched open/close tags")
	}
	return nil
}

func (a *Any) castIn(element Arg) error {
	t := NewTagList(nil)
	switch element.(type) {
	case IsAtomic:
		tag, err := NewTag(nil)
		if err != nil {
			return errors.New("error creating empty tag")
		}
		if err := element.(interface{ Encode(arg Arg) error }).Encode(tag); err != nil {
			return errors.New("error encoding element")
		}
		t.Append(tag)
	case IsAnyAtomic:
		tag, err := NewTag(nil)
		if err != nil {
			return errors.New("error creating empty tag")
		}
		if err := element.(interface{ GetValue() any }).GetValue().(interface{ Encode(Tag) error }).Encode(tag); err != nil {
			return errors.New("error encoding element")
		}
		t.Append(tag)
	default:
		if err := element.(interface{ Encode(arg Arg) error }).Encode(t); err != nil {
			return errors.New("error encoding element")
		}
	}
	a.tagList.Extend(t.GetTagList()...)
	return nil
}

func (a *Any) castOut(element Arg) error {
	panic("implement me")
	return nil
}
