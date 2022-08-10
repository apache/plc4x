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

package model

import (
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
)

func ParseAmsString(io utils.ReadBuffer, stringLength int32, encoding string) (string, error) {
	var multiplier int32
	switch encoding {
	case "UTF-8":
		multiplier = 8
	case "UTF-16":
		multiplier = 16
	}
	return io.ReadString("", uint32(stringLength*multiplier))
}

func SerializeAmsString(io utils.WriteBuffer, value values.PlcValue, stringLength int32, encoding string) error {
	var multiplier int32
	switch encoding {
	case "UTF-8":
		multiplier = 8
	case "UTF-16":
		multiplier = 16
	}
	return io.WriteString("", uint32(stringLength*multiplier), encoding, value.GetString())
}

func STR_LEN(str string) int {
	return len(str)
}
