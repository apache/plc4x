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

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

func expandEnumerations(e EnumeratedContract) {
	xlateTable := e.GetXlateTable()

	for name, value := range e.GetEnumerations() {
		// save the result
		xlateTable[name] = value
		xlateTable[value] = name
	}
}

// EnumeratedContract provides a set of functions which can be overwritten by a sub struct
type EnumeratedContract interface {
	GetEnumerations() map[string]uint64
	GetXlateTable() map[any]any
	// SetEnumerated is required because we do more stuff in the constructor and can't wait for the substruct to finish that
	SetEnumerated(enumerated *Enumerated)
}

type Enumerated struct {
	*Atomic[uint64]
	EnumeratedContract

	_appTag model.BACnetDataType

	_xlateTable map[any]any

	valueString string
}

func NewEnumerated(args Args) (*Enumerated, error) {
	e := &Enumerated{
		_appTag: model.BACnetDataType_ENUMERATED,
	}
	e.EnumeratedContract = e
	e.Atomic = NewAtomic[uint64](e)

	if args == nil || len(args) < 1 {
		return e, nil
	}
	var arg any
	switch len(args) {
	case 1:
		arg = args[0]
	case 2:
		var ok bool
		e.EnumeratedContract, ok = args[0].(EnumeratedContract)
		if !ok {
			return nil, errors.Errorf("%T must be implement EnumeratedContract", args[0])
		}
		e.EnumeratedContract.SetEnumerated(e)
		arg = args[1]
	default:
		return nil, errors.Errorf("invalid arguments %T. %[1]v", args)
	}
	if len(e.EnumeratedContract.GetXlateTable()) == 0 {
		expandEnumerations(e.EnumeratedContract)
	}

	switch arg := arg.(type) {
	case *tag:
		err := e.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return e, nil
	case uint:
		e.value = uint64(arg)

		// convert it to a string if you can
		e.valueString, _ = e.EnumeratedContract.GetXlateTable()[e.value].(string)
	case int:
		if arg < 0 {
			return nil, errors.New("arg must be positive")
		}
		e.value = uint64(arg)

		// convert it to a string if you can
		e.valueString, _ = e.EnumeratedContract.GetXlateTable()[e.value].(string)
	case uint64:
		e.value = arg

		// convert it to a string if you can
		e.valueString, _ = e.EnumeratedContract.GetXlateTable()[e.value].(string)
	case string:
		var ok bool
		var value any
		value, ok = e.EnumeratedContract.GetXlateTable()[arg]
		if !ok {
			return nil, errors.Errorf("undefined enumeration %s", arg)
		}
		e.value = value.(uint64)
		e.valueString = arg
	case *Enumerated:
		e.value = arg.value
		e.valueString = arg.valueString
		e._xlateTable = make(map[any]any)
		for k, v := range arg._xlateTable {
			e._xlateTable[k] = v
		}
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", arg)
	}

	return e, nil
}

func (e *Enumerated) GetAppTag() model.BACnetDataType {
	return e._appTag
}

func (e *Enumerated) GetEnumerations() map[string]uint64 {
	return make(map[string]uint64)
}

func (e *Enumerated) GetXlateTable() map[any]any {
	if e._xlateTable == nil {
		e._xlateTable = make(map[any]any)
	}
	return e._xlateTable
}

func (e *Enumerated) SetEnumerated(_ *Enumerated) {
	panic("must be implemented by substruct")
}

func (e *Enumerated) GetItem(item any) (result any, ok bool) {
	v, ok := e.EnumeratedContract.GetXlateTable()[item]
	return v, ok
}

func (e *Enumerated) GetLong() uint64 {
	if mappedValue, ok := e.EnumeratedContract.GetXlateTable()[e.valueString]; ok {
		return mappedValue.(uint64)
	}
	return e.value
}

func (e *Enumerated) Keylist() []string {
	var result []string
	for key := range e.EnumeratedContract.GetEnumerations() {
		result = append(result, key)
	}
	return result
}

func (e *Enumerated) Compare(other any) int {
	otherEnumerated, ok := other.(Enumerated)
	if !ok {
		return -1
	}

	// get the numeric version
	a := e.GetLong()
	b := otherEnumerated.GetLong()

	// now compare the values
	if a < b {
		return -1
	} else if a > b {
		return 1
	} else {
		return 0
	}
}

func (e *Enumerated) Decode(arg Arg) error {
	tag, ok := arg.(Tag)
	if !ok {
		return errors.Errorf("%T is not a Tag", arg)
	}
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(e._appTag) {
		return errors.New("bit string application tag required")
	}
	if len(tag.GetTagData()) == 0 {
		return errors.New("invalid tag length")
	}

	data := tag.GetTagData()
	if len(data) < 8 {
		data = append(make([]byte, 8-len(data)), data...)
	}
	// get the data
	rslt := binary.BigEndian.Uint64(data)

	if mappedValue, ok := e.EnumeratedContract.GetXlateTable()[rslt]; ok {
		e.valueString = mappedValue.(string)
	}

	// save the result
	e.value = rslt
	return nil
}

func (e *Enumerated) Encode(arg Arg) error {
	tag, ok := arg.(Tag)
	if !ok {
		return errors.Errorf("%T is not a Tag", arg)
	}
	value := e.value
	if mappedValue, ok := e.EnumeratedContract.GetXlateTable()[e.valueString]; ok {
		value = mappedValue.(uint64)
	}

	data := make([]byte, 8)
	// rip apart the number
	binary.BigEndian.PutUint64(data, value)

	// reduce the value to the smallest number of octets
	for len(data) > 1 && data[0] == 0 {
		data = data[1:]
	}

	// encode the tag
	tag.setAppData(uint(e._appTag), data)
	return nil
}

func (e *Enumerated) IsValid(arg any) bool {
	_, ok := arg.(uint64)
	return ok
}

func (e *Enumerated) GetValueString() string {
	return e.valueString
}

func (e *Enumerated) String() string {
	value := strconv.Itoa(int(e.value))
	if e.valueString != "" {
		value = e.valueString
	}
	return fmt.Sprintf("Enumerated(%s)", value)
}
