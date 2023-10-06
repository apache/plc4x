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

package utils

import (
	"strconv"

	"github.com/apache/plc4x/plc4go/pkg/api/values"
)

func PlcValueUint8ListToByteArray(value values.PlcValue) []byte {
	if value == nil || !value.IsList() {
		return []byte{}
	}
	list := value.GetList()
	result := make([]byte, len(list))
	for i, valueItem := range list {
		// TODO: we should sanity check if this is indeed a uint8
		result[i] = valueItem.GetUint8()
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

func StrToInt32(str string) (int32, error) {
	intVal, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(intVal), nil
}

func StrToString(s string) (string, error) {
	return s, nil
}
