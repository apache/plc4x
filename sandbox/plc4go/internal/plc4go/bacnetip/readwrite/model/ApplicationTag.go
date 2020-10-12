//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package model

import "plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"

type ApplicationTag int8

const (
	ApplicationTag_NULL                     ApplicationTag = 0x0
	ApplicationTag_BOOLEAN                  ApplicationTag = 0x1
	ApplicationTag_UNSIGNED_INTEGER         ApplicationTag = 0x2
	ApplicationTag_SIGNED_INTEGER           ApplicationTag = 0x3
	ApplicationTag_REAL                     ApplicationTag = 0x4
	ApplicationTag_DOUBLE                   ApplicationTag = 0x5
	ApplicationTag_OCTET_STRING             ApplicationTag = 0x6
	ApplicationTag_CHARACTER_STRING         ApplicationTag = 0x7
	ApplicationTag_BIT_STRING               ApplicationTag = 0x8
	ApplicationTag_ENUMERATED               ApplicationTag = 0x9
	ApplicationTag_DATE                     ApplicationTag = 0xA
	ApplicationTag_TIME                     ApplicationTag = 0xB
	ApplicationTag_BACNET_OBJECT_IDENTIFIER ApplicationTag = 0xC
)

func CastApplicationTag(structType interface{}) ApplicationTag {
	castFunc := func(typ interface{}) ApplicationTag {
		if sApplicationTag, ok := typ.(ApplicationTag); ok {
			return sApplicationTag
		}
		return 0
	}
	return castFunc(structType)
}

func ApplicationTagParse(io spi.ReadBuffer) (ApplicationTag, error) {
	// TODO: Implement ...
	return 0, nil
}

func (e ApplicationTag) Serialize(io spi.WriteBuffer) {
	// TODO: Implement ...
}
