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
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type Element interface {
	GetName() string
	GetKlass() func(Args, KWArgs) (ElementKlass, error)
	GetContext() *int
	IsOptional() bool
	Encode(tagList Arg) error
}

type ElementKlass interface {
	Encode(Arg) error
	GetAppTag() readWriteModel.BACnetDataType
}

// TODO: finish
type _Element struct {
	Name     string
	Klass    func(Args, KWArgs) (ElementKlass, error)
	Context  *int
	Optional bool
}

func NewElement(name string, klass func(Args, KWArgs) (ElementKlass, error), opts ...func(*_Element)) Element {
	e := &_Element{
		Name:  name,
		Klass: klass,
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

var _ Element = (*_Element)(nil)

func WithElementOptional(optional bool) func(*_Element) {
	return func(e *_Element) {
		e.Optional = optional
	}
}

func WithElementContext(context int) func(*_Element) {
	return func(e *_Element) {
		e.Context = &context
	}
}

func (e *_Element) GetName() string {
	return e.Name
}

func (e *_Element) GetKlass() func(Args, KWArgs) (ElementKlass, error) {
	return e.Klass
}

func (e *_Element) GetContext() *int {
	return e.Context
}

func (e *_Element) IsOptional() bool {
	return e.Optional
}

func (e *_Element) Encode(tagList Arg) error {
	//TODO implement me
	panic("implement me")
}
