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
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
)

// TODO: big WIP
type ObjectIdentifierProperty interface {
	ReadableProperty
}

type _ObjectIdentifierProperty struct {
	ReadableProperty
}

func NewObjectIdentifierProperty(name string, klass func(Args, KWArgs) (PropertyKlass, error), options ...Option) ObjectIdentifierProperty {
	o := &_ObjectIdentifierProperty{}
	o.ReadableProperty = NewReadableProperty(name, klass, options...)
	return o
}

func (o *_ObjectIdentifierProperty) WriteProperty(args Args, kwArgs KWArgs) error {
	obj := GA[Object](args, 0)
	value := GA[any](args, 1)
	arrayIndex, _ := KWO[int](kwArgs, "arrayIndex", 0)
	priority, _ := KWO[int](kwArgs, "priority", 0)
	if _debug != nil {
		_debug("WriteProperty %r %r arrayIndex=%r priority=%r", obj, value, arrayIndex, priority)
	}

	// make it easy to default
	if value == nil {
		return nil
	}
	switch castedValue := value.(type) {
	case int:
		value = ObjectIdentifierTuple{Left: obj.GetObjectType(), Right: castedValue}
	case ObjectIdentifierTuple:
		if castedValue.Left != obj.GetObjectType() {
			return ValueError{Message: fmt.Sprintf("%s required", obj.GetObjectType())}
		}
	default:
		return TypeError{Message: "object required"}
	}

	args[1] = value
	return o.ReadableProperty.WriteProperty(args, kwArgs)
}
