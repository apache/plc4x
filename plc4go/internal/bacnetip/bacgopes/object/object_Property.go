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

package object

import (
	"fmt"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/errors"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

// TODO: big WIP
type Property interface {
	GetIdentifier() string
	GetDataType() any
	IsOptional() bool
	Get_Default() any
	IsMutable() bool
	ReadProperty(args Args, kwArgs KWArgs) error
	WriteProperty(args Args, kwArgs KWArgs) error
}

type PropertyKlass interface {
	Encode(Arg) error
}

func NewProperty(identifier string, datatype func(Args, KWArgs) (PropertyKlass, error), options ...Option) Property {
	i := &_Property{
		// keep the arguments
		identifier: identifier,
		dataType:   datatype,
	}
	ApplyAppliers(options, i)

	// check the datatype
	dataType, _ := i.dataType(Nothing())
	switch dataType.(type) {
	default:
		// TODO: check for Atomic, Sequence, Choice, Array, List, AnyAtomic
	}

	return i
}

type _Property struct {
	identifier string
	dataType   func(Args, KWArgs) (PropertyKlass, error)
	optional   bool
	_default   any
	mutable    bool
}

var _ Property = (*_Property)(nil)

func WithPropertyOptional(optional bool) GenericApplier[*_Property] {
	return WrapGenericApplier(func(e *_Property) { e.optional = optional })
}

func WithPropertyDefault(_default any) GenericApplier[*_Property] {
	return WrapGenericApplier(func(e *_Property) { e._default = _default })
}

func WithPropertyMutable(mutable bool) GenericApplier[*_Property] {
	return WrapGenericApplier(func(e *_Property) { e.mutable = mutable })
}

func (p *_Property) GetIdentifier() string {
	return p.identifier
}

func (p *_Property) GetDataType() any {
	return p.dataType
}

func (p *_Property) IsOptional() bool {
	return p.optional
}

func (p *_Property) Get_Default() any {
	return p._default
}

func (p *_Property) IsMutable() bool {
	return p.mutable
}

func (p *_Property) ReadProperty(Args, KWArgs) error {
	//TODO implement me
	panic("implement me")
}

func (p *_Property) WriteProperty(args Args, kwArgs KWArgs) error {
	obj := GA[Object](args, 0)
	value := GA[any](args, 1)
	arrayIndex, hasArrayIndex := KWO[int](kwArgs, "arrayIndex", 0)
	priority, _ := KWO[int](kwArgs, "priority", 0)
	direct, _ := KWO[bool](kwArgs, "direct", false)
	if _debug != nil {
		_debug("WriteProperty %r %r arrayIndex=%r priority=%r", obj, value, arrayIndex, priority)
	}
	if direct {
		if _debug != nil {
			_debug("    - direct write")
		}
	} else {
		// see if it must be provided
		if !p.optional && IsNil(value) {
			return ValueError{Message: fmt.Sprintf("%v value required", p.identifier)}
		}

		// see if it can be changed
		if !p.mutable {
			return errors.ExecutionError{ErrorClass: readWriteModel.ErrorClass_PROPERTY, ErrorCode: readWriteModel.ErrorCode_WRITE_ACCESS_DENIED}
		}

		// if changing the length if the array, the value is unsigned
		if hasArrayIndex && arrayIndex == 0 {
			panic("implement unsigned check") // TODO
		}
		panic("finish me") // TODO
	}

	// local check if the property is monitored
	_, isMonitored := obj.Get_PropertiesMonitors()[p.GetIdentifier()]

	if hasArrayIndex {
		panic("finish me") // TODO
	} else {
		var oldValue any
		if isMonitored {
			oldValue = obj.Get_Values()[p.GetIdentifier()]
		}

		// seems to be OK
		obj.Get_Values()[p.GetIdentifier()] = value

		// check for monitors, call each one with the old and new value
		if isMonitored {
			for _, fn := range obj.Get_PropertiesMonitors()[p.GetIdentifier()] {
				if _debug != nil {
					_debug("    - monitor: %r", fn)
				}
				fn(oldValue, value)
			}
		}
	}
	return nil
}
