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

package utils

import (
	"github.com/apache/plc4x/plc4go/pkg/plc4go/values"
	"strconv"
	"strings"
)

func Int8ArrayToUint8Array(input []int8) []uint8 {
	output := make([]uint8, len(input))
	if input != nil {
		for i, _val := range input {
			output[i] = uint8(_val)
		}
	}
	return output
}

func Int8ArrayToString(data []int8, separator string) string {
	var sb strings.Builder
	if data != nil {
		for i, element := range data {
			sb.WriteString(strconv.Itoa(int(uint8(element))))
			if i < (len(data) - 1) {
				sb.WriteString(separator)
			}
		}
	}
	return sb.String()
}

func Uint8ArrayToInt8Array(input []uint8) []int8 {
	output := make([]int8, len(input))
	if input != nil {
		for i, _val := range input {
			output[i] = int8(_val)
		}
	}
	return output
}

func Int8ArrayToByteArray(input []int8) []byte {
	output := make([]byte, len(input))
	if input != nil {
		for i, _val := range input {
			output[i] = byte(_val)
		}
	}
	return output
}

func ByteArrayToInt8Array(input []byte) []int8 {
	output := make([]int8, len(input))
	if input != nil {
		for i, _val := range input {
			output[i] = int8(_val)
		}
	}
	return output
}

func ByteArrayToUint8Array(input []byte) []uint8 {
	output := make([]uint8, len(input))
	if input != nil {
		for i, _val := range input {
			output[i] = _val
		}
	}
	return output
}

func PlcValueUint8ListToByteArray(value values.PlcValue) []byte {
	var result []byte
	if value != nil {
		for _, valueItem := range value.GetList() {
			result = append(result, valueItem.GetUint8())
		}
	}
	return result
}

func StrToBool(str string) (bool, error) {
	boolVal, err := strconv.ParseBool(str)
	if err != nil {
		return false, err
	}
	return boolVal, nil
}

func StrToUint8(str string) (uint8, error) {
	intVal, err := strconv.ParseInt(str, 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(intVal), nil
}

func StrToUint16(str string) (uint16, error) {
	intVal, err := strconv.ParseInt(str, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(intVal), nil
}

func StrToUint32(str string) (uint32, error) {
	intVal, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(intVal), nil
}
