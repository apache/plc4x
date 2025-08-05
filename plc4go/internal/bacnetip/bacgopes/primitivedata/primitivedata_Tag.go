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
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

const (
	// Deprecated: use model.TagClass_APPLICATION_TAGS
	TagApplicationTagClass = readWriteModel.TagClass_APPLICATION_TAGS
	// Deprecated: use model.TagClass_CONTEXT_SPECIFIC_TAGS
	TagContextTagClass = readWriteModel.TagClass_CONTEXT_SPECIFIC_TAGS
	TagOpeningTagClass = 2
	TagClosingTagClass = 3

	// Deprecated: use  model.BACnetDataType_NULL
	TagNullAppTag = readWriteModel.BACnetDataType_NULL
	// Deprecated: use  model.BACnetDataType_BOOLEAN
	TagBooleanAppTag = readWriteModel.BACnetDataType_BOOLEAN
	// Deprecated: use  model.BACnetDataType_UNSIGNED_INTEGER
	TagUnsignedAppTag = readWriteModel.BACnetDataType_UNSIGNED_INTEGER
	// Deprecated: use  model.BACnetDataType_SIGNED_INTEGER
	TagIntegerAppTag = readWriteModel.BACnetDataType_SIGNED_INTEGER
	// Deprecated: use  model.BACnetDataType_REAL
	TagRealAppTag = readWriteModel.BACnetDataType_REAL
	// Deprecated: use  model.BACnetDataType_DOUBLE
	TagDoubleAppTag = readWriteModel.BACnetDataType_DOUBLE
	// Deprecated: use  model.BACnetDataType_OCTET_STRING
	TagOctetStringAppTag = readWriteModel.BACnetDataType_OCTET_STRING
	// Deprecated: use  model.BACnetDataType_CHARACTER_STRING
	TagCharacterStringAppTag = readWriteModel.BACnetDataType_CHARACTER_STRING
	// Deprecated: use  model.BACnetDataType_BIT_STRING
	TagBitStringAppTag = readWriteModel.BACnetDataType_BIT_STRING
	// Deprecated: use  model.BACnetDataType_ENUMERATED
	TagEnumeratedAppTag = readWriteModel.BACnetDataType_ENUMERATED
	// Deprecated: use  model.BACnetDataType_DATE
	TagDateAppTag = readWriteModel.BACnetDataType_DATE
	// Deprecated: use  model.BACnetDataType_TIME
	TagTimeAppTag = readWriteModel.BACnetDataType_TIME
	// Deprecated: use  model.BACnetDataType_BACNET_OBJECT_IDENTIFIER
	TagObjectIdentifierAppTag = readWriteModel.BACnetDataType_BACNET_OBJECT_IDENTIFIER
	TagReservedAppTag13       = 13
	TagReservedAppTag14       = 14
	TagReservedAppTag15       = 15
)

type Tag interface {
	DebugContentPrinter
	GetTagClass() readWriteModel.TagClass
	GetTagNumber() uint
	GetTagLvt() int
	GetTagData() []byte
	Encode(pdu PDUData)
	Decode(pdu PDUData) error
	AppToObject() (any, error)
	AppToContext(context uint) (*ContextTag, error)
	ContextToApp(dataType uint) (Tag, error) // TODO: can't be ApplicationTag because boolean gets encoded as Tag???
	setAppData(tagNumber uint, tdata []byte)
	set(args Args)
}

type tag struct {
	tagClass  readWriteModel.TagClass
	tagNumber uint
	tagLVT    int
	tagData   []byte

	appTagName  []string
	appTagClass []any

	_leafName string
}

var _ Tag = (*tag)(nil)

func NewTag(args Args) (Tag, error) {
	t := &tag{
		appTagName: []string{
			"null", "boolean", "unsigned", "integer",
			"real", "double", "octetString", "characterString",
			"bitString", "enumerated", "date", "time",
			"objectIdentifier", "reserved13", "reserved14", "reserved15",
		},
		appTagClass: []any{
			&Null{}, &Boolean{}, &Unsigned{}, &Integer{},
			&Real{}, &Double{}, &OctetString{}, &CharacterString{},
			&BitString{}, &Enumerated{}, &Date{}, &Time{},
			&ObjectIdentifier{}, nil, nil, nil,
		},
		_leafName: StructName(), //TODO: extract from options...
	}
	if len(args) == 0 {
		return t, nil
	}
	if len(args) == 1 {
		if err := t.Decode(args[0].(PDUData)); err != nil {
			return nil, errors.New("error decoding")
		}
	} else if len(args) >= 2 {
		t.set(args)
	} else {
		return nil, errors.Errorf("invalid arguments %v", args)
	}
	return t, nil
}

func (t *tag) Decode(pdu PDUData) error {
	tag, err := pdu.Get()
	if err != nil {
		return errors.Wrap(err, "error decoding tag")
	}

	// extract the type
	t.tagClass = readWriteModel.TagClass(tag >> 3 & 0x01)

	// extract the tag number
	t.tagNumber = uint(tag >> 4)
	if t.tagNumber == 0x0f {
		get, err := pdu.Get()
		if err != nil {
			return errors.Wrap(err, "error decoding get")
		}
		t.tagNumber = uint(get)
	}

	// extract the length
	t.tagLVT = int(tag & 0x07)
	if t.tagLVT == 5 {
		get, err := pdu.Get()
		if err != nil {
			return errors.Wrap(err, "error decoding get")
		}
		t.tagLVT = int(get)
		if t.tagLVT == 254 {
			get, err := pdu.GetShort()
			if err != nil {
				return errors.Wrap(err, "error decoding get")
			}
			t.tagLVT = int(get)
		} else if t.tagLVT == 255 {
			get, err := pdu.GetLong()
			if err != nil {
				return errors.Wrap(err, "error decoding get")
			}
			t.tagLVT = int(get)
		}
	} else if t.tagLVT == 6 {
		t.tagClass = TagOpeningTagClass
		t.tagLVT = 0
	} else if t.tagLVT == 7 {
		t.tagClass = TagClosingTagClass
		t.tagLVT = 0
	}

	if t.tagClass == readWriteModel.TagClass_APPLICATION_TAGS && t.tagNumber == uint(readWriteModel.BACnetDataType_BOOLEAN) {
		// tagLVT contains value
		t.tagData = nil
	} else {
		// tagLVT contains length
		t.tagData, err = pdu.GetData(t.tagLVT)
		if err != nil {
			return errors.Wrap(err, "error decoding tag data")
		}
	}
	return nil
}

func (t *tag) Encode(pdu PDUData) {
	var data byte
	// check for special encoding
	if t.tagClass == readWriteModel.TagClass_CONTEXT_SPECIFIC_TAGS {
		data = 0x08
	} else if t.tagClass == TagOpeningTagClass {
		data = 0x0E
	} else if t.tagClass == TagClosingTagClass {
		data = 0x0F
	} else {
		data = 0x00
	}

	// encode the tag number part
	if t.tagNumber < 15 {
		data += byte(t.tagNumber) << 4
	} else {
		data += 0xF0
	}

	// encode the length/value/type part
	if t.tagLVT < 5 {
		data += byte(t.tagLVT)
	} else {
		data += 0x05
	}

	// save this and the extended tag value
	pdu.Put(data)
	if t.tagNumber >= 15 {
		pdu.Put(byte(t.tagNumber))
	}

	// really short lengths are already done
	if t.tagLVT >= 5 {
		if t.tagLVT <= 253 {
			pdu.Put(byte(t.tagLVT))
		} else if t.tagLVT <= 65535 {
			pdu.Put(254)
			pdu.PutShort(uint16(t.tagLVT))
		} else {
			pdu.Put(255)
			pdu.PutLong(uint32(t.tagLVT))
		}
	}

	// now put the data
	pdu.PutData(t.tagData...)
}

func (t *tag) AppToContext(context uint) (*ContextTag, error) {
	if t.tagClass != readWriteModel.TagClass_APPLICATION_TAGS {
		return nil, errors.New("application tag required")
	}
	if t.tagNumber == uint(readWriteModel.BACnetDataType_BOOLEAN) {
		return NewContextTag(NA(context, []byte{byte(t.tagLVT)}))
	}
	return NewContextTag(NA(context, t.tagData))
}

func (t *tag) ContextToApp(dataType uint) (Tag, error) {
	if t.tagClass != readWriteModel.TagClass_CONTEXT_SPECIFIC_TAGS {
		return nil, errors.New("context tag required")
	}
	if dataType == uint(readWriteModel.BACnetDataType_BOOLEAN) {
		return NewTag(NA(readWriteModel.TagClass_APPLICATION_TAGS, readWriteModel.BACnetDataType_BOOLEAN, t.tagData[0], nil))
	}
	return NewApplicationTag(NA(dataType, t.tagData))
}

func (t *tag) AppToObject() (any, error) {
	if t.tagClass != readWriteModel.TagClass_APPLICATION_TAGS {
		return nil, errors.New("context tag required")
	}

	klass := t.appTagClass[int(t.tagNumber)]
	if klass == nil {
		return nil, nil
	}

	switch klass.(type) {
	case *Null:
		return NewNull(t)
	case *Boolean:
		return NewBoolean(t)
	case *Unsigned:
		return NewUnsigned(t)
	case *Integer:
		return NewInteger(t)
	case *Real:
		return NewReal(t)
	case *Double:
		return NewDouble(t)
	case *OctetString:
		return NewOctetString(t)
	case *CharacterString:
		return NewCharacterString(t)
	case *BitString:
		return NewBitString(NA(t))
	case *Enumerated:
		return NewEnumerated(NA(t))
	case *Date:
		return NewDate(NA(t))
	case *Time:
		return NewTime(NA(t))
	case *ObjectIdentifier:
		return NewObjectIdentifier(NA(t))
	default:
		return nil, errors.Errorf("unknown tag klass %T", klass)
	}
}

func (t *tag) set(args Args) {
	switch tagClass := args[0].(type) {
	case readWriteModel.TagClass:
		t.tagClass = tagClass
	case uint:
		t.tagClass = readWriteModel.TagClass(tagClass)
	case uint8:
		t.tagClass = readWriteModel.TagClass(tagClass)
	case int:
		t.tagClass = readWriteModel.TagClass(tagClass)
	default:
		panic("oh no")
	}
	switch tagNumber := args[1].(type) {
	case readWriteModel.BACnetDataType:
		t.tagNumber = uint(tagNumber)
	case uint:
		t.tagNumber = tagNumber
	case int:
		t.tagNumber = uint(tagNumber)
	default:
		panic("oh no")
	}
	if len(args) == 2 {
		return
	}
	switch tagLVT := args[2].(type) {
	case uint:
		t.tagLVT = int(tagLVT)
	case uint8:
		t.tagLVT = int(tagLVT)
	case int:
		t.tagLVT = tagLVT
	default:
		panic("oh no")
	}
	if len(args) == 3 {
		return
	}
	switch tagData := args[3].(type) {
	case []byte:
		t.tagData = tagData
	case nil:
		// ok
	default:
		panic("oh no")
	}
	if len(args) > 4 {
		panic("oh no")
	}
}

func (t *tag) setAppData(tagNumber uint, tdata []byte) {
	t.tagClass = readWriteModel.TagClass_APPLICATION_TAGS
	t.tagNumber = tagNumber
	t.tagLVT = len(tdata)
	t.tagData = tdata
}

func (t *tag) GetTagClass() readWriteModel.TagClass {
	return t.tagClass
}

func (t *tag) GetTagNumber() uint {
	return t.tagNumber
}

func (t *tag) GetTagLvt() int {
	return t.tagLVT
}

func (t *tag) GetTagData() []byte {
	return t.tagData
}

func (t *tag) Equals(other any) bool {
	if t == nil && other == nil {
		return true
	}
	if other == nil {
		return false
	}
	otherTag, ok := other.(Tag)
	if !ok {
		return false
	}
	return t.tagClass == otherTag.GetTagClass() &&
		t.tagNumber == otherTag.GetTagNumber() &&
		t.tagLVT == otherTag.GetTagLvt() &&
		bytes.Equal(t.tagData, otherTag.GetTagData())
}

func (t *tag) Format(s fmt.State, v rune) {
	switch v {
	case 'r':
		sname := t._leafName
		var desc string
		if t.tagClass == TagOpeningTagClass {
			desc = fmt.Sprintf("(open(%d))", t.tagNumber)
		} else if t.tagClass == TagClosingTagClass {
			desc = fmt.Sprintf("(close(%d))", t.tagNumber)
		} else if t.tagClass == readWriteModel.TagClass_CONTEXT_SPECIFIC_TAGS {
			desc = fmt.Sprintf("(context(%d))", t.tagNumber)
		} else if t.tagClass == readWriteModel.TagClass_APPLICATION_TAGS && t.tagNumber < uint(len(t.appTagName)) {
			desc = fmt.Sprintf("(%s)", t.appTagName[t.tagNumber])
		} else {
			desc = "(?)"
		}

		_, _ = fmt.Fprintf(s, "<%s%s instance at %p", sname, desc, t)
	}
}

func (t *tag) PrintDebugContents(indent int, file io.Writer, _ids []uintptr) {
	// object reference first
	_, _ = fmt.Fprintf(file, "%s%r\n", strings.Repeat("    ", indent), t)
	indent += 1

	// tag class
	msg := fmt.Sprintf("%stagClass = %s ", strings.Repeat("    ", indent), t.tagClass)
	if t.tagClass == readWriteModel.TagClass_APPLICATION_TAGS {
		msg += "application"
	} else if t.tagClass == readWriteModel.TagClass_CONTEXT_SPECIFIC_TAGS {
		msg += "context"
	} else if t.tagClass == TagOpeningTagClass {
		msg += "opening"
	} else if t.tagClass == TagClosingTagClass {
		msg += "closing"
	} else {
		msg += "?"
	}
	_, _ = fmt.Fprint(file, msg+"\n")

	// tag number
	msg = fmt.Sprintf("%stagNumber = %d ", strings.Repeat("    ", indent), t.tagNumber)
	if t.tagClass == readWriteModel.TagClass_APPLICATION_TAGS && t.tagNumber < uint(len(t.appTagName)) {
		msg += t.appTagName[t.tagNumber]
	} else {
		msg += "?"
	}
	_, _ = fmt.Fprint(file, msg+"\n")

	// length, value, type
	_, _ = fmt.Fprintf(file, "%stagLVT = %d\n", strings.Repeat("    ", indent), t.tagLVT)

	// data
	_, _ = fmt.Fprintf(file, "%stagData = '%s'\n", strings.Repeat("    ", indent), Btox(t.tagData, "."))
}
