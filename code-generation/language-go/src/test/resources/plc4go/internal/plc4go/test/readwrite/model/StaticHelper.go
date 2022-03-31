/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package model

import "github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"

func ParseBit(io utils.ReadBuffer) int8 {
	return 0
}

func SerializeBit(io utils.WriteBuffer, data byte) {
}

func ParseByte(io utils.ReadBuffer) int8 {
	return 0
}

func SerializeByte(io utils.WriteBuffer, data byte) {
}

func ParseInt8(io utils.ReadBuffer) int8 {
	return 0
}

func SerializeInt8(io utils.WriteBuffer, data byte) {
}

func ParseUint8(io utils.ReadBuffer) int8 {
	return 0
}

func SerializeUint8(io utils.WriteBuffer, data byte) {
}

func ParseFloat(io utils.ReadBuffer) int8 {
	return 0
}

func SerializeFloat(io utils.WriteBuffer, data byte) {
}

func ParseDouble(io utils.ReadBuffer) int8 {
	return 0
}

func SerializeDouble(io utils.WriteBuffer, data byte) {
}

func ParseString(io utils.ReadBuffer) int8 {
	return 0
}

func SerializeString(io utils.WriteBuffer, data byte) {
}

func Parse(io utils.ReadBuffer) int8 {
	return 0
}

func Serialize(io utils.WriteBuffer, data byte) {
}

func ReadManualField(io utils.ReadBuffer, simpleField uint8) (uint8, error) {
	return 0, nil
}

func CrcInt8(num int) (int8, error) {
	return int8(num), nil
}

func CrcUint8(num int) (uint8, error) {
	return uint8(num), nil
}
