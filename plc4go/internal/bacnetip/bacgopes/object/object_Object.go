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
	"io"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/basetypes"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
)

// TODO: big WIP
type Object interface {
	DebugContentPrinter
	GetAttr(name string) (any, bool)
	SetAttr(name string, value any)
	AddProperty(prop Property)
	DeleteProperty(prop string)
	ReadProperty(Args, KWArgs) error
	WriteProperty(Args, KWArgs) error
}

type _Object struct {
	_objectSupportsCov bool

	properties        []Property
	_properties       map[string]Property
	_propertyMonitors map[string]any

	_app    any
	_values map[string]any

	_leafName string
}

func NewObject(kwArgs KWArgs, options ...Option) (Object, error) {
	o := &_Object{
		properties: []Property{
			NewObjectIdentifierProperty("objectIdentifier", Vs2P(NewObjectIdentifier), WithPropertyOptional(false)),
			NewReadableProperty("objectName", V2P(NewCharacterString), WithPropertyOptional(false)),
			NewOptionalProperty("description", V2P(NewCharacterString)),
			NewOptionalProperty("profileName", V2P(NewCharacterString)),
			NewReadableProperty("propertyList", ArrayOfP(NewPropertyIdentifier, 0, 0)),
			NewOptionalProperty("auditLevel", V2P(NewAuditLevel)),
			NewOptionalProperty("auditableOperations", V2P(NewAuditOperationFlags)),
			NewOptionalProperty("tags", ArrayOfsP(NewNameValue, 0, 0)),
			NewOptionalProperty("profileLocation", V2P(NewCharacterString)),
			NewOptionalProperty("profileName", V2P(NewCharacterString)),
		},
		_properties: map[string]Property{},
		_leafName:   ExtractLeafName(options, StructName()),
	}
	ApplyAppliers(options, o)
	if _debug != nil {
		_debug("__init__(%s) %r", o._leafName, kwArgs)
	}
	// map the golang names into property names and make sure they
	// are appropriate for this object
	var initArgs = make(KWArgs)
	for key, value := range kwArgs {
		if _, ok := o._properties[string(key)]; ok {
			return nil, PropertyError{string(key)}
		}
		initArgs[key] = value
	}

	// object is detached from an application until it is added
	o._app = nil

	// start with a clean dict of values
	o._values = make(map[string]any)

	// empty list of property monitors
	o._propertyMonitors = make(map[string]any)

	// initialize the object
	for propid, prop := range o._properties {
		if _, ok := initArgs[KnownKey(propid)]; ok {
			if _debug != nil {
				_debug("    - setting %s from initargs", propid)
			}

			// defer to the property object for error checking
			if err := prop.WriteProperty(NA(o, initArgs[KnownKey(propid)]), NKW(KnownKey("direct"), "true")); err != nil {
				return nil, errors.Wrap(err, "error writing property")
			}
		} else if prop.Get_Default() != nil {
			if _debug != nil {
				_debug("    - setting %s from default", propid)
			}

			// default values bypass property interface
			o._values[propid] = prop.Get_Default()
		} else {
			if !prop.IsOptional() {
				if _debug != nil {
					_debug("    - %s value required", propid)
				}
			}

			o._values[propid] = nil
		}
	}

	if _debug != nil {
		_debug("    - done __init__")
	}

	return o, nil
}

func WithObject_Properties(_properties map[string]Property) GenericApplier[*_Object] {
	return WrapGenericApplier(func(o *_Object) {
		o._properties = _properties
	})
}

func (o *_Object) PrintDebugContents(indent int, file io.Writer, _ids []uintptr) {
	//TODO implement me
	panic("implement me")
}

func (o *_Object) GetAttr(name string) (any, bool) {
	v, ok := o._values[name]
	return v, ok
}

func (o *_Object) SetAttr(name string, value any) {
	//TODO implement me
	panic("implement me")
}

func (o *_Object) AddProperty(prop Property) {
	//TODO implement me
	panic("implement me")
}

func (o *_Object) DeleteProperty(prop string) {
	//TODO implement me
	panic("implement me")
}

func (o *_Object) ReadProperty(Args, KWArgs) error {
	//TODO implement me
	panic("implement me")
}

func (o *_Object) WriteProperty(Args, KWArgs) error {
	//TODO implement me
	panic("implement me")
}
