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
	"math"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type Double struct {
	*Atomic[float64]
	*CommonMath

	_appTag model.BACnetDataType
}

func NewDouble(arg Arg) (*Double, error) {
	b := &Double{
		_appTag: model.BACnetDataType_DOUBLE,
	}
	b.Atomic = NewAtomic[float64](b)

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
	case float32:
		b.value = float64(arg)
	case float64:
		b.value = arg
	case int:
		b.value = float64(arg)
	case *Double:
		b.value = arg.value
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", arg)
	}

	return b, nil
}

func (d *Double) GetAppTag() model.BACnetDataType {
	return d._appTag
}

func (d *Double) Encode(arg Arg) error {
	tag, ok := arg.(Tag)
	if !ok {
		return errors.Errorf("%T is not a Tag", arg)
	}
	var _b = make([]byte, 8)
	binary.BigEndian.PutUint64(_b, math.Float64bits(d.value))
	tag.setAppData(uint(d._appTag), _b)
	return nil
}

func (d *Double) Decode(arg Arg) error {
	tag, ok := arg.(Tag)
	if !ok {
		return errors.Errorf("%T is not a Tag", arg)
	}
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(d._appTag) {
		return errors.New("Double application tag required")
	}
	if len(tag.GetTagData()) != 8 {
		return errors.New("invalid tag length")
	}

	// extract the data
	d.value = math.Float64frombits(binary.BigEndian.Uint64(tag.GetTagData()))
	return nil
}

func (d *Double) IsValid(arg any) bool {
	_, ok := arg.(float64)
	return ok
}

func (d *Double) String() string {
	return fmt.Sprintf("Double(%g)", d.value)
}
