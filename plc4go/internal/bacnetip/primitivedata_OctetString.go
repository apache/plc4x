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
	"bytes"
	"fmt"

	"github.com/pkg/errors"

	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type OctetString struct {
	value []byte
}

func NewOctetString(arg Arg) (*OctetString, error) {
	o := &OctetString{}
	o.value = make([]byte, 0)

	if arg == nil {
		return o, nil
	}
	switch arg := arg.(type) {
	case *tag:
		err := o.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return o, nil
	case []byte:
		if len(arg) == 0 {
			arg = nil
		}
		o.value = arg
	case *OctetString:
		o.value = arg.value
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", arg)
	}

	return o, nil
}

func (o *OctetString) Encode(tag Tag) {
	tag.setAppData(uint(model.BACnetDataType_OCTET_STRING), o.value)
}

func (o *OctetString) Decode(tag Tag) error {
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(model.BACnetDataType_OCTET_STRING) {
		return errors.New("OctetString application tag required")
	}

	o.value = tag.GetTagData()
	return nil
}

func (o *OctetString) Compare(other any) int {
	if _, ok := other.(byte); !ok {
		return -1
	}
	return len(o.value) - len(other.(OctetString).value)
}

func (o *OctetString) LowerThan(other any) bool {
	if _, ok := other.(byte); !ok {
		return false
	}
	return len(o.value) < len(other.(OctetString).value)
}

func (o *OctetString) Equals(other any) bool {
	if _, ok := other.(byte); !ok {
		return false
	}
	return bytes.Equal(o.value, other.([]byte))
}

func (o *OctetString) GetValue() []byte {
	return o.value
}

func (o *OctetString) IsValid(arg any) bool {
	_, ok := arg.([]byte)
	return ok
}

func (o *OctetString) String() string {
	return fmt.Sprintf("OctetString(X'%s')", Btox(o.value, ""))
}
