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

package values

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValues(t *testing.T) {
	values := []PlcValueType{
		NULL,
		BOOL,
		BYTE,
		WORD,
		DWORD,
		LWORD,
		USINT,
		UINT,
		UDINT,
		ULINT,
		SINT,
		INT,
		DINT,
		LINT,
		REAL,
		LREAL,
		CHAR,
		WCHAR,
		STRING,
		WSTRING,
		TIME,
		LTIME,
		DATE,
		LDATE,
		TIME_OF_DAY,
		LTIME_OF_DAY,
		DATE_AND_TIME,
		LDATE_AND_TIME,
		Struct,
		List,
		RAW_BYTE_ARRAY,
		0xFF,
	}
	for _, value := range values {
		t.Run(fmt.Sprintf("value %d", value), func(t *testing.T) {
			valueString := value.String()
			assert.NotEmpty(t, valueString)
			valueType, ok := PlcValueByName(valueString)
			expectedValue := value
			expectedOk := true
			if value == 0xff {
				expectedValue = NULL
				expectedOk = false
			}
			assert.Equal(t, expectedValue, valueType, "value type not equals")
			assert.Equal(t, expectedOk, ok, "ok state doesn't match")
		})
	}
}
