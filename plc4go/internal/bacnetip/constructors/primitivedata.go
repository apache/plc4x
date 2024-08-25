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

package constructors

import (
	"github.com/apache/plc4x/plc4go/internal/bacnetip"
)

func BitString(args ...any) *bacnetip.BitString {
	bitString, err := bacnetip.NewBitString(args)
	if err != nil {
		panic(err)
	}
	return bitString
}

func Boolean(arg ...any) *bacnetip.Boolean {
	if len(arg) == 0 {
		boolean, err := bacnetip.NewBoolean(nil)
		if err != nil {
			panic(err)
		}
		return boolean
	}
	boolean, err := bacnetip.NewBoolean(arg[0])
	if err != nil {
		panic(err)
	}
	return boolean
}

func CharacterString(arg ...any) *bacnetip.CharacterString {
	if len(arg) == 0 {
		CharacterString, err := bacnetip.NewCharacterString(nil)
		if err != nil {
			panic(err)
		}
		return CharacterString
	}
	CharacterString, err := bacnetip.NewCharacterString(arg[0])
	if err != nil {
		panic(err)
	}
	return CharacterString
}

func Date(arg ...any) *bacnetip.Date {
	if len(arg) == 0 {
		Date, err := bacnetip.NewDate(nil, nil)
		if err != nil {
			panic(err)
		}
		return Date
	}
	Date, err := bacnetip.NewDate(arg[0], nil)
	if err != nil {
		panic(err)
	}
	return Date
}

func Double(arg ...any) *bacnetip.Double {
	if len(arg) == 0 {
		Double, err := bacnetip.NewDouble(nil)
		if err != nil {
			panic(err)
		}
		return Double
	}
	Double, err := bacnetip.NewDouble(arg[0])
	if err != nil {
		panic(err)
	}
	return Double
}

func Enumerated(arg ...any) *bacnetip.Enumerated {
	Enumerated, err := bacnetip.NewEnumerated(arg...)
	if err != nil {
		panic(err)
	}
	return Enumerated
}

func Integer(arg ...any) *bacnetip.Integer {
	if len(arg) == 0 {
		Integer, err := bacnetip.NewInteger(nil)
		if err != nil {
			panic(err)
		}
		return Integer
	}
	Integer, err := bacnetip.NewInteger(arg[0])
	if err != nil {
		panic(err)
	}
	return Integer
}

func Null(arg ...any) *bacnetip.Null {
	if len(arg) == 0 {
		Null, err := bacnetip.NewNull(nil)
		if err != nil {
			panic(err)
		}
		return Null
	}
	Null, err := bacnetip.NewNull(arg[0])
	if err != nil {
		panic(err)
	}
	return Null
}

func ObjectIdentifier(args ...any) *bacnetip.ObjectIdentifier {
	if len(args) == 0 {
		ObjectIdentifier, err := bacnetip.NewObjectIdentifier(nil)
		if err != nil {
			panic(err)
		}
		return ObjectIdentifier
	}
	ObjectIdentifier, err := bacnetip.NewObjectIdentifier(args)
	if err != nil {
		panic(err)
	}
	return ObjectIdentifier
}

func ObjectType(args ...any) *bacnetip.ObjectType {
	if len(args) == 0 {
		ObjectType, err := bacnetip.NewObjectType(nil)
		if err != nil {
			panic(err)
		}
		return ObjectType
	}
	ObjectType, err := bacnetip.NewObjectType(args)
	if err != nil {
		panic(err)
	}
	return ObjectType
}

func OctetString(args ...any) *bacnetip.OctetString {
	if len(args) == 0 {
		OctetString, err := bacnetip.NewOctetString(nil)
		if err != nil {
			panic(err)
		}
		return OctetString
	}
	OctetString, err := bacnetip.NewOctetString(args[0])
	if err != nil {
		panic(err)
	}
	return OctetString
}

func Real(arg ...any) *bacnetip.Real {
	if len(arg) == 0 {
		Real, err := bacnetip.NewReal(nil)
		if err != nil {
			panic(err)
		}
		return Real
	}
	Real, err := bacnetip.NewReal(arg[0])
	if err != nil {
		panic(err)
	}
	return Real
}

func ApplicationTag(args ...any) *bacnetip.ApplicationTag {
	tag, err := bacnetip.NewApplicationTag(args)
	if err != nil {
		panic(err)
	}
	return tag
}

func ContextTag(args ...any) *bacnetip.ContextTag {
	tag, err := bacnetip.NewContextTag(args)
	if err != nil {
		panic(err)
	}
	return tag
}

func OpeningTag(context any) *bacnetip.OpeningTag {
	openingTag, err := bacnetip.NewOpeningTag(context)
	if err != nil {
		panic(err)
	}
	return openingTag
}

func ClosingTag(context any) *bacnetip.ClosingTag {
	closingTag, err := bacnetip.NewClosingTag(context)
	if err != nil {
		panic(err)
	}
	return closingTag
}

func TagList(tags ...bacnetip.Tag) *bacnetip.TagList {
	return bacnetip.NewTagList(tags)
}

func Time(arg ...any) *bacnetip.Time {
	if len(arg) == 0 {
		Time, err := bacnetip.NewTime(nil, nil)
		if err != nil {
			panic(err)
		}
		return Time
	}
	Time, err := bacnetip.NewTime(arg[0], nil)
	if err != nil {
		panic(err)
	}
	return Time
}

func Unsigned(arg ...any) *bacnetip.Unsigned {
	if len(arg) == 0 {
		unsigned, err := bacnetip.NewUnsigned(nil)
		if err != nil {
			panic(err)
		}
		return unsigned
	}
	unsigned, err := bacnetip.NewUnsigned(arg[0])
	if err != nil {
		panic(err)
	}
	return unsigned
}

func Unsigned8(arg ...any) *bacnetip.Unsigned8 {
	if len(arg) == 0 {
		unsigned, err := bacnetip.NewUnsigned8(nil)
		if err != nil {
			panic(err)
		}
		return unsigned
	}
	unsigned, err := bacnetip.NewUnsigned8(arg[0])
	if err != nil {
		panic(err)
	}
	return unsigned
}

func Unsigned16(arg ...any) *bacnetip.Unsigned16 {
	if len(arg) == 0 {
		unsigned, err := bacnetip.NewUnsigned16(nil)
		if err != nil {
			panic(err)
		}
		return unsigned
	}
	unsigned, err := bacnetip.NewUnsigned16(arg[0])
	if err != nil {
		panic(err)
	}
	return unsigned
}
