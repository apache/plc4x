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
	"fmt"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type CharacterString struct {
	*Atomic[string]
	*CommonMath

	_appTag model.BACnetDataType

	strEncoding byte
	strValue    []byte
}

func NewCharacterString(arg Arg) (*CharacterString, error) {
	c := &CharacterString{
		_appTag: model.BACnetDataType_CHARACTER_STRING,
	}
	c.Atomic = NewAtomic[string](c)

	if arg == nil {
		return c, nil
	}
	switch arg := arg.(type) {
	case *tag:
		err := c.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return c, nil
	case string:
		c.value = arg
		c.strValue = []byte(c.value)
	case *CharacterString:
		c.value = arg.value
		c.strEncoding = arg.strEncoding
		c.strValue = arg.strValue
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", arg)
	}

	return c, nil
}

func (c *CharacterString) GetAppTag() model.BACnetDataType {
	return c._appTag
}

func (c *CharacterString) Encode(arg Arg) error {
	tag, ok := arg.(Tag)
	if !ok {
		return errors.Errorf("%T is not a Tag", arg)
	}
	tag.setAppData(uint(c._appTag), append([]byte{c.strEncoding}, c.strValue...))
	return nil
}

func (c *CharacterString) Decode(arg Arg) error {
	tag, ok := arg.(Tag)
	if !ok {
		return errors.Errorf("%T is not a Tag", arg)
	}
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(c._appTag) {
		return errors.New("CharacterString application tag required")
	}
	if len(tag.GetTagData()) == 0 {
		return errors.New("invalid tag length")
	}

	tagData := tag.GetTagData()

	// extract the data
	c.strEncoding = tagData[0]
	c.strValue = tagData[1:]

	// normalize the value
	switch c.strEncoding {
	case 0:
		c.value = string(c.strValue)
	case 3: //utf_32be
		panic("implement me") // TODO: implement me
	case 4: //utf_16be
		panic("implement me") // TODO: implement me
	case 5: //latin_1
		panic("implement me") // TODO: implement me
	default:
		c.value = fmt.Sprintf("### unknown encoding: %d ###", c.strEncoding)
	}

	return nil
}

func (c *CharacterString) IsValid(arg any) bool {
	_, ok := arg.(bool)
	return ok
}

func (c *CharacterString) String() string {
	return fmt.Sprintf("CharacterString(%d,X'%s')", c.strEncoding, Btox(c.strValue, ""))
}
