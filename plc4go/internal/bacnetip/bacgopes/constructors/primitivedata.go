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
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes"
)

func BitString(args ...any) *bacgopes.BitString {
	bitString, err := bacgopes.NewBitString(args)
	if err != nil {
		panic(err)
	}
	return bitString
}

func Boolean(arg ...any) *bacgopes.Boolean {
	if len(arg) == 0 {
		boolean, err := bacgopes.NewBoolean(nil)
		if err != nil {
			panic(err)
		}
		return boolean
	}
	boolean, err := bacgopes.NewBoolean(arg[0])
	if err != nil {
		panic(err)
	}
	return boolean
}

func CharacterString(arg ...any) *bacgopes.CharacterString {
	if len(arg) == 0 {
		CharacterString, err := bacgopes.NewCharacterString(nil)
		if err != nil {
			panic(err)
		}
		return CharacterString
	}
	CharacterString, err := bacgopes.NewCharacterString(arg[0])
	if err != nil {
		panic(err)
	}
	return CharacterString
}

func Date(arg ...any) *bacgopes.Date {
	if len(arg) == 0 {
		Date, err := bacgopes.NewDate(nil, nil)
		if err != nil {
			panic(err)
		}
		return Date
	}
	Date, err := bacgopes.NewDate(arg[0], nil)
	if err != nil {
		panic(err)
	}
	return Date
}

func Double(arg ...any) *bacgopes.Double {
	if len(arg) == 0 {
		Double, err := bacgopes.NewDouble(nil)
		if err != nil {
			panic(err)
		}
		return Double
	}
	Double, err := bacgopes.NewDouble(arg[0])
	if err != nil {
		panic(err)
	}
	return Double
}

func Enumerated(arg ...any) *bacgopes.Enumerated {
	Enumerated, err := bacgopes.NewEnumerated(arg...)
	if err != nil {
		panic(err)
	}
	return Enumerated
}

func Integer(arg ...any) *bacgopes.Integer {
	if len(arg) == 0 {
		Integer, err := bacgopes.NewInteger(nil)
		if err != nil {
			panic(err)
		}
		return Integer
	}
	Integer, err := bacgopes.NewInteger(arg[0])
	if err != nil {
		panic(err)
	}
	return Integer
}

func Null(arg ...any) *bacgopes.Null {
	if len(arg) == 0 {
		Null, err := bacgopes.NewNull(nil)
		if err != nil {
			panic(err)
		}
		return Null
	}
	Null, err := bacgopes.NewNull(arg[0])
	if err != nil {
		panic(err)
	}
	return Null
}

func ObjectIdentifier(args ...any) *bacgopes.ObjectIdentifier {
	if len(args) == 0 {
		ObjectIdentifier, err := bacgopes.NewObjectIdentifier(nil)
		if err != nil {
			panic(err)
		}
		return ObjectIdentifier
	}
	ObjectIdentifier, err := bacgopes.NewObjectIdentifier(args)
	if err != nil {
		panic(err)
	}
	return ObjectIdentifier
}

func ObjectType(args ...any) *bacgopes.ObjectType {
	if len(args) == 0 {
		ObjectType, err := bacgopes.NewObjectType(nil)
		if err != nil {
			panic(err)
		}
		return ObjectType
	}
	ObjectType, err := bacgopes.NewObjectType(args)
	if err != nil {
		panic(err)
	}
	return ObjectType
}

func OctetString(args ...any) *bacgopes.OctetString {
	if len(args) == 0 {
		OctetString, err := bacgopes.NewOctetString(nil)
		if err != nil {
			panic(err)
		}
		return OctetString
	}
	OctetString, err := bacgopes.NewOctetString(args[0])
	if err != nil {
		panic(err)
	}
	return OctetString
}

func Real(arg ...any) *bacgopes.Real {
	if len(arg) == 0 {
		Real, err := bacgopes.NewReal(nil)
		if err != nil {
			panic(err)
		}
		return Real
	}
	Real, err := bacgopes.NewReal(arg[0])
	if err != nil {
		panic(err)
	}
	return Real
}

func ApplicationTag(args ...any) *bacgopes.ApplicationTag {
	tag, err := bacgopes.NewApplicationTag(args)
	if err != nil {
		panic(err)
	}
	return tag
}

func ContextTag(args ...any) *bacgopes.ContextTag {
	tag, err := bacgopes.NewContextTag(args)
	if err != nil {
		panic(err)
	}
	return tag
}

func OpeningTag(context any) *bacgopes.OpeningTag {
	openingTag, err := bacgopes.NewOpeningTag(context)
	if err != nil {
		panic(err)
	}
	return openingTag
}

func ClosingTag(context any) *bacgopes.ClosingTag {
	closingTag, err := bacgopes.NewClosingTag(context)
	if err != nil {
		panic(err)
	}
	return closingTag
}

func TagList(tags ...bacgopes.Tag) *bacgopes.TagList {
	return bacgopes.NewTagList(tags)
}

func Time(arg ...any) *bacgopes.Time {
	if len(arg) == 0 {
		Time, err := bacgopes.NewTime(nil, nil)
		if err != nil {
			panic(err)
		}
		return Time
	}
	Time, err := bacgopes.NewTime(arg[0], nil)
	if err != nil {
		panic(err)
	}
	return Time
}

func Unsigned(arg ...any) *bacgopes.Unsigned {
	if len(arg) == 0 {
		unsigned, err := bacgopes.NewUnsigned(nil)
		if err != nil {
			panic(err)
		}
		return unsigned
	}
	unsigned, err := bacgopes.NewUnsigned(arg[0])
	if err != nil {
		panic(err)
	}
	return unsigned
}

func Unsigned8(arg ...any) *bacgopes.Unsigned8 {
	if len(arg) == 0 {
		unsigned, err := bacgopes.NewUnsigned8(nil)
		if err != nil {
			panic(err)
		}
		return unsigned
	}
	unsigned, err := bacgopes.NewUnsigned8(arg[0])
	if err != nil {
		panic(err)
	}
	return unsigned
}

func Unsigned16(arg ...any) *bacgopes.Unsigned16 {
	if len(arg) == 0 {
		unsigned, err := bacgopes.NewUnsigned16(nil)
		if err != nil {
			panic(err)
		}
		return unsigned
	}
	unsigned, err := bacgopes.NewUnsigned16(arg[0])
	if err != nil {
		panic(err)
	}
	return unsigned
}
