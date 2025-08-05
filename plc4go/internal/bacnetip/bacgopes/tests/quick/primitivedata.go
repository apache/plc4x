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

package quick

import "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"

func Tag(args ...any) primitivedata.Tag {
	tag, err := primitivedata.NewTag(args)
	if err != nil {
		panic(err)
	}
	return tag
}

func BitString(args ...any) *primitivedata.BitString {
	bitString, err := primitivedata.NewBitString(args)
	if err != nil {
		panic(err)
	}
	return bitString
}

func Boolean(arg ...any) *primitivedata.Boolean {
	if len(arg) == 0 {
		boolean, err := primitivedata.NewBoolean(nil)
		if err != nil {
			panic(err)
		}
		return boolean
	}
	boolean, err := primitivedata.NewBoolean(arg[0])
	if err != nil {
		panic(err)
	}
	return boolean
}

func CharacterString(arg ...any) *primitivedata.CharacterString {
	if len(arg) == 0 {
		CharacterString, err := primitivedata.NewCharacterString(nil)
		if err != nil {
			panic(err)
		}
		return CharacterString
	}
	CharacterString, err := primitivedata.NewCharacterString(arg[0])
	if err != nil {
		panic(err)
	}
	return CharacterString
}

func Date(args ...any) *primitivedata.Date {
	Date, err := primitivedata.NewDate(args)
	if err != nil {
		panic(err)
	}
	return Date
}

func Double(arg ...any) *primitivedata.Double {
	if len(arg) == 0 {
		Double, err := primitivedata.NewDouble(nil)
		if err != nil {
			panic(err)
		}
		return Double
	}
	Double, err := primitivedata.NewDouble(arg[0])
	if err != nil {
		panic(err)
	}
	return Double
}

func Enumerated(args ...any) *primitivedata.Enumerated {
	Enumerated, err := primitivedata.NewEnumerated(args)
	if err != nil {
		panic(err)
	}
	return Enumerated
}

func Integer(arg ...any) *primitivedata.Integer {
	if len(arg) == 0 {
		Integer, err := primitivedata.NewInteger(nil)
		if err != nil {
			panic(err)
		}
		return Integer
	}
	Integer, err := primitivedata.NewInteger(arg[0])
	if err != nil {
		panic(err)
	}
	return Integer
}

func Null(arg ...any) *primitivedata.Null {
	if len(arg) == 0 {
		Null, err := primitivedata.NewNull(nil)
		if err != nil {
			panic(err)
		}
		return Null
	}
	Null, err := primitivedata.NewNull(arg[0])
	if err != nil {
		panic(err)
	}
	return Null
}

func ObjectIdentifier(args ...any) *primitivedata.ObjectIdentifier {
	if len(args) == 0 {
		ObjectIdentifier, err := primitivedata.NewObjectIdentifier(nil)
		if err != nil {
			panic(err)
		}
		return ObjectIdentifier
	}
	ObjectIdentifier, err := primitivedata.NewObjectIdentifier(args)
	if err != nil {
		panic(err)
	}
	return ObjectIdentifier
}

func ObjectType(args ...any) *primitivedata.ObjectType {
	if len(args) == 0 {
		ObjectType, err := primitivedata.NewObjectType(nil)
		if err != nil {
			panic(err)
		}
		return ObjectType
	}
	ObjectType, err := primitivedata.NewObjectType(args)
	if err != nil {
		panic(err)
	}
	return ObjectType
}

func OctetString(args ...any) *primitivedata.OctetString {
	if len(args) == 0 {
		OctetString, err := primitivedata.NewOctetString(nil)
		if err != nil {
			panic(err)
		}
		return OctetString
	}
	OctetString, err := primitivedata.NewOctetString(args[0])
	if err != nil {
		panic(err)
	}
	return OctetString
}

func Real(arg ...any) *primitivedata.Real {
	if len(arg) == 0 {
		Real, err := primitivedata.NewReal(nil)
		if err != nil {
			panic(err)
		}
		return Real
	}
	Real, err := primitivedata.NewReal(arg[0])
	if err != nil {
		panic(err)
	}
	return Real
}

func ApplicationTag(args ...any) *primitivedata.ApplicationTag {
	tag, err := primitivedata.NewApplicationTag(args)
	if err != nil {
		panic(err)
	}
	return tag
}

func ContextTag(args ...any) *primitivedata.ContextTag {
	tag, err := primitivedata.NewContextTag(args)
	if err != nil {
		panic(err)
	}
	return tag
}

func OpeningTag(context any) *primitivedata.OpeningTag {
	openingTag, err := primitivedata.NewOpeningTag(context)
	if err != nil {
		panic(err)
	}
	return openingTag
}

func ClosingTag(context any) *primitivedata.ClosingTag {
	closingTag, err := primitivedata.NewClosingTag(context)
	if err != nil {
		panic(err)
	}
	return closingTag
}

func TagList(tags ...primitivedata.Tag) *primitivedata.TagList {
	return primitivedata.NewTagList(tags)
}

func Time(args ...any) *primitivedata.Time {
	t, err := primitivedata.NewTime(args)
	if err != nil {
		panic(err)
	}
	return t
}

func Unsigned(arg ...any) *primitivedata.Unsigned {
	if len(arg) == 0 {
		unsigned, err := primitivedata.NewUnsigned(nil)
		if err != nil {
			panic(err)
		}
		return unsigned
	}
	unsigned, err := primitivedata.NewUnsigned(arg[0])
	if err != nil {
		panic(err)
	}
	return unsigned
}

func Unsigned8(arg ...any) *primitivedata.Unsigned8 {
	if len(arg) == 0 {
		unsigned, err := primitivedata.NewUnsigned8(nil)
		if err != nil {
			panic(err)
		}
		return unsigned
	}
	unsigned, err := primitivedata.NewUnsigned8(arg[0])
	if err != nil {
		panic(err)
	}
	return unsigned
}

func Unsigned16(arg ...any) *primitivedata.Unsigned16 {
	if len(arg) == 0 {
		unsigned, err := primitivedata.NewUnsigned16(nil)
		if err != nil {
			panic(err)
		}
		return unsigned
	}
	unsigned, err := primitivedata.NewUnsigned16(arg[0])
	if err != nil {
		panic(err)
	}
	return unsigned
}
