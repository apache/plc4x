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
	"strconv"
	"strings"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type ObjectIdentifierTuple struct {
	Left  any
	Right int
}

type ObjectIdentifier struct {
	//*Atomic[...] won't work here

	_appTag model.BACnetDataType

	objectTypeClass *ObjectType

	value ObjectIdentifierTuple
}

func NewObjectIdentifier(args Args) (*ObjectIdentifier, error) {
	i := &ObjectIdentifier{
		_appTag: model.BACnetDataType_BACNET_OBJECT_IDENTIFIER,
	}
	var err error
	i.objectTypeClass, err = NewObjectType(nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating object type")
	}
	i.value = ObjectIdentifierTuple{"analogInput", 0}

	if len(args) == 0 || args == nil {
		return i, nil
	}
	if len(args) == 1 {
		arg := args[0]
		switch arg := arg.(type) {
		case *tag:
			err := i.Decode(arg)
			if err != nil {
				return nil, errors.Wrap(err, "error decoding")
			}
			return i, nil
		case int:
			i.setLong(arg)
		case string:
			split := strings.Split(arg, ":")
			var objType, objInstance any = split[0], split[1]
			if objTypeInt, err := strconv.Atoi(fmt.Sprintf("%v", objType)); err == nil {
				objType = objTypeInt
			}
			var err error
			objInstance, err = strconv.Atoi(fmt.Sprintf("%v", objInstance))
			if err != nil {
				return nil, errors.Wrap(err, "error parsing instance")
			}
			if err := i.setTuple(objType, objInstance.(int)); err != nil {
				return nil, errors.Wrap(err, "can't set tuple")
			}
		case ObjectIdentifierTuple:
			if err := i.setTuple(arg.Left, arg.Right); err != nil {
				return nil, errors.Wrap(err, "error setting tuple")
			}
		case *ObjectIdentifier:
			i.value = arg.value
		default:
			return nil, errors.Errorf("invalid constructor datatype: %T", arg)
		}
	} else if len(args) == 2 {
		err := i.setTuple(args[0], args[1].(int))
		if err != nil {
			return nil, errors.Wrap(err, "error setting tuple")
		}
	} else {
		return nil, errors.New("invalid constructor parameters")
	}

	return i, nil
}

func (o *ObjectIdentifier) GetAppTag() model.BACnetDataType {
	return o._appTag
}

func (o *ObjectIdentifier) setTuple(objType any, objInstance int) error {
	switch objType.(type) {
	case int:
		if gotObjType, ok := o.objectTypeClass.GetXlateTable()[uint64(objType.(int))]; ok {
			objType = gotObjType
		}
	case string:
		if _, ok := o.objectTypeClass.GetXlateTable()[objType]; !ok {
			return errors.Errorf("unrecognized object type %s", objType)
		}
	default:
		return errors.Errorf("invalid datatype for object type: %T", objType)
	}

	// check valid instance number
	if objInstance < 0 || objInstance > 0x003FFFFF {
		return errors.Errorf("invalid object instance out of range: %d", objInstance)
	}

	o.value = ObjectIdentifierTuple{objType, objInstance}
	return nil
}

func (o *ObjectIdentifier) getTuple() ObjectIdentifierTuple {
	return o.value
}

func (o *ObjectIdentifier) setLong(value int) {
	// suck out the type
	objTypeInt := (value >> 22) & 0x3ff
	var objType any = objTypeInt

	// try and make it pretty
	if item, ok := o.objectTypeClass.GetItem(uint64(objTypeInt)); ok {
		objType = item
	}

	// suck out the instance
	objInstance := value & 0x003FFFFF

	// save the result
	o.value = ObjectIdentifierTuple{objType, objInstance}
}

func (o *ObjectIdentifier) getLong() int {
	tuple := o.getTuple()
	objType, objInstance := tuple.Left, tuple.Right

	if _, ok := objType.(string); ok {
		if objTypeGot, ok := o.objectTypeClass.GetXlateTable()[objType]; ok {
			objType = int(objTypeGot.(uint64))
		}
	}

	return (objType.(int) << 22) + objInstance
}

func (o *ObjectIdentifier) Encode(arg Arg) error {
	tag, ok := arg.(Tag)
	if !ok {
		return errors.Errorf("%T is not a Tag", arg)
	}
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(o.getLong()))
	tag.setAppData(uint(o._appTag), data)
	return nil
}

func (o *ObjectIdentifier) Decode(arg Arg) error {
	tag, ok := arg.(Tag)
	if !ok {
		return errors.Errorf("%T is not a Tag", arg)
	}
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(o._appTag) {
		return errors.New("ObjectIdentifier application tag required")
	}
	if len(tag.GetTagData()) != 4 {
		return errors.New("invalid tag length")
	}

	tagData := tag.GetTagData()

	o.setLong(int(binary.BigEndian.Uint32(tagData)))
	return nil
}

func (o *ObjectIdentifier) IsValid(arg any) bool {
	switch arg.(type) {
	case ObjectIdentifier:
		return true
	default:
		return false
	}
}

func (o *ObjectIdentifier) Compare(other any) int {
	switch other := other.(type) {
	case *ObjectIdentifier:
		return o.getLong() - other.getLong()
	default:
		return -1
	}
}

func (o *ObjectIdentifier) LowerThan(other any) bool {
	switch other := other.(type) {
	case *ObjectIdentifier:
		return o.getLong() < other.getLong()
	default:
		return false
	}
}

func (o *ObjectIdentifier) Equals(other any) bool {
	return o.value == other
}

func (o *ObjectIdentifier) GetValue() ObjectIdentifierTuple {
	return o.value
}

func (o *ObjectIdentifier) Coerce(arg ObjectIdentifier) ObjectIdentifierTuple {
	return arg.GetValue()
}

func (o *ObjectIdentifier) String() string {
	// rip it apart
	objType, objInstance := o.value.Left, o.value.Right

	var objTypeAsUint64 uint64
	if objTypeAsInt, ok := objType.(int); ok {
		objTypeAsUint64 = uint64(objTypeAsInt)
	}
	var typeString string
	if s, ok := objType.(string); ok {
		typeString = s
	} else if i, intOk := objType.(int); intOk && i < 0 {
		typeString = fmt.Sprintf("Bad %d", i)
	} else if gotType, xlateOk := o.objectTypeClass.GetXlateTable()[objTypeAsUint64]; xlateOk {
		typeString = gotType.(string)
	} else if intOk && i < 128 {
		typeString = fmt.Sprintf("Reserved %d", i)
	} else {
		typeString = fmt.Sprintf("Vendor %s", objType)
	}

	return fmt.Sprintf("ObjectIdentifier(%s,%d)", typeString, objInstance)
}
