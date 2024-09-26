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
	"fmt"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type Element interface {
	fmt.Formatter
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

type _Element struct {
	name     string
	klass    func(Args, KWArgs) (ElementKlass, error)
	context  *int
	optional bool

	_leafName string
}

func NewElement(name string, klass func(Args, KWArgs) (ElementKlass, error), options ...Option) Element {
	e := &_Element{
		name:      name,
		klass:     klass,
		_leafName: ExtractLeafName(options, StructName()),
	}
	ApplyAppliers(options, e)
	return e
}

var _ Element = (*_Element)(nil)

func WithElementOptional(optional bool) GenericApplier[*_Element] {
	return WrapGenericApplier(func(e *_Element) { e.optional = optional })
}

func WithElementContext(context int) GenericApplier[*_Element] {
	return WrapGenericApplier(func(e *_Element) { e.context = &context })
}

func (e *_Element) GetName() string {
	return e.name
}

func (e *_Element) GetKlass() func(Args, KWArgs) (ElementKlass, error) {
	return e.klass
}

func (e *_Element) GetContext() *int {
	return e.context
}

func (e *_Element) IsOptional() bool {
	return e.optional
}

func (e *_Element) Encode(tagList Arg) error {
	//TODO implement me
	panic("implement me")
}

func (e *_Element) Format(s fmt.State, verb rune) {
	switch verb {
	case 'r':
		desc := fmt.Sprintf("%s(%s", e._leafName, e.name)
		elementKlass, _ := e.klass(Nothing())
		desc += " " + QualifiedTypeName(elementKlass)
		if e.context != nil {
			desc += fmt.Sprintf(", context=%d", *e.context)
		}
		if e.optional {
			desc += ", optional"
		}
		desc += ")"

		_, _ = fmt.Fprintf(s, "<%s instance at %p>", desc, e)
	}
}
