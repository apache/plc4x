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
	"math"
	"reflect"
	"strings"
	"testing"
	"time"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/stretchr/testify/assert"
)

func TestCombinations(t *testing.T) {
	tests := []struct {
		name      apiValues.PlcValueType
		arguments []apiValues.PlcValue
	}{
		/*{
			name: apiValues.BINT,
			arguments: []apiValues.PlcValue{
				NewPlcBINT(big.NewInt(0)),
				NewPlcBINT(big.NewInt(64)),
				NewPlcBINT(big.NewInt(255)),
				NewPlcBINT(big.NewInt(math.MinInt64)),
				NewPlcBINT(big.NewInt(math.MaxInt64)),
			},
		},*/
		{
			name: apiValues.BOOL,
			arguments: []apiValues.PlcValue{
				NewPlcBOOL(true),
				NewPlcBOOL(false),
			},
		},
		/*{
			name: apiValues.BREAL,
			arguments: []apiValues.PlcValue{
				NewPlcBREAL(big.NewFloat(0)),
				NewPlcBREAL(big.NewFloat(64)),
				NewPlcBREAL(big.NewFloat(255)),
				NewPlcBREAL(big.NewFloat(math.MinInt64)),
				NewPlcBREAL(big.NewFloat(math.MaxInt64)),
			},
		},*/
		{
			name: apiValues.BYTE,
			arguments: []apiValues.PlcValue{
				NewPlcBYTE(0),
				NewPlcBYTE(64),
				NewPlcBYTE(255),
			},
		},
		{
			name: apiValues.CHAR,
			arguments: []apiValues.PlcValue{
				NewPlcCHAR(""),
			},
		},
		{
			name: apiValues.DATE,
			arguments: []apiValues.PlcValue{
				NewPlcDATE(time.Now()),
			},
		},
		{
			name: apiValues.DATE_AND_TIME,
			arguments: []apiValues.PlcValue{
				NewPlcDATE_AND_TIME(time.Now()),
			},
		},
		{
			name: apiValues.DINT,
			arguments: []apiValues.PlcValue{
				NewPlcDINT(math.MinInt32),
				NewPlcDINT(64),
				NewPlcDINT(255),
				NewPlcDINT(math.MaxInt32),
			},
		},
		{
			name: apiValues.DWORD,
			arguments: []apiValues.PlcValue{
				NewPlcDWORD(0),
				NewPlcDWORD(64),
				NewPlcDWORD(255),
				NewPlcDWORD(math.MaxUint32),
			},
		},
		{
			name: apiValues.INT,
			arguments: []apiValues.PlcValue{
				NewPlcINT(0),
				NewPlcINT(64),
				NewPlcINT(255),
				NewPlcINT(math.MaxInt16),
			},
		},
		{
			name: apiValues.LINT,
			arguments: []apiValues.PlcValue{
				NewPlcLINT(0),
				NewPlcLINT(64),
				NewPlcLINT(255),
				NewPlcLINT(math.MaxInt64),
			},
		},
		{
			name: apiValues.List,
			arguments: []apiValues.PlcValue{
				NewPlcList(nil),
				NewPlcList([]apiValues.PlcValue{
					NewPlcBOOL(true),
					NewPlcBOOL(false),
				}),
			},
		},
		{
			name: apiValues.LREAL,
			arguments: []apiValues.PlcValue{
				NewPlcREAL(0),
				NewPlcREAL(64),
				NewPlcREAL(255),
				NewPlcREAL(math.MinInt64),
				NewPlcREAL(math.MaxInt64),
			},
		},
		{
			name: apiValues.LTIME,
			arguments: []apiValues.PlcValue{
				NewPlcLTIME(0),
				NewPlcLTIME(64),
				NewPlcLTIME(255),
				NewPlcLTIME(math.MaxUint64),
			},
		},
		{
			name: apiValues.LWORD,
			arguments: []apiValues.PlcValue{
				NewPlcLWORD(0),
				NewPlcLWORD(64),
				NewPlcLWORD(255),
				NewPlcLWORD(0),
				NewPlcLWORD(math.MaxUint64),
			},
		},
		{
			name: apiValues.NULL,
			arguments: []apiValues.PlcValue{
				NewPlcNULL(),
			},
		},
		{
			name: apiValues.RAW_BYTE_ARRAY,
			arguments: []apiValues.PlcValue{
				NewPlcRawByteArray(utils.NewReadBufferByteBased([]byte{0x47, 0x11}).GetBytes()),
			},
		},
		{
			name: apiValues.REAL,
			arguments: []apiValues.PlcValue{
				NewPlcREAL(0),
				NewPlcREAL(64),
				NewPlcREAL(255),
				NewPlcREAL(math.MaxInt64),
			},
		},
		{
			name: apiValues.Struct,
			arguments: []apiValues.PlcValue{
				NewPlcStruct(map[string]apiValues.PlcValue{
					"something": NewPlcStruct(map[string]apiValues.PlcValue{
						"more": NewPlcList([]apiValues.PlcValue{
							NewPlcBOOL(true),
							NewPlcBOOL(false),
						}),
						"evenMore": NewPlcBOOL(false),
					}),
					"somethingOther": NewPlcBOOL(true),
				}),
			},
		},
		{
			name: apiValues.SINT,
			arguments: []apiValues.PlcValue{
				NewPlcSINT(-128),
				NewPlcSINT(-64),
				NewPlcSINT(0),
				NewPlcSINT(64),
				NewPlcSINT(127),
			},
		},
		{
			name: apiValues.STRING,
			arguments: []apiValues.PlcValue{
				NewPlcSTRING("Hello Tody"),
				NewPlcSTRING("1"),
				NewPlcSTRING("true"),
				NewPlcSTRING("255"),
				NewPlcSTRING("5.4"),
			},
		},
		{
			name: apiValues.TIME,
			arguments: []apiValues.PlcValue{
				NewPlcTIME(0),
				NewPlcTIME(64),
				NewPlcTIME(255),
				NewPlcTIME(math.MaxUint32),
			},
		},
		{
			name: apiValues.TIME_OF_DAY,
			arguments: []apiValues.PlcValue{
				NewPlcTIME_OF_DAY(time.Now()),
				NewPlcTIME_OF_DAY(0),
				NewPlcTIME_OF_DAY(64),
				NewPlcTIME_OF_DAY(255),
				NewPlcTIME_OF_DAY(math.MaxUint32),
			},
		},
		{
			name: apiValues.UDINT,
			arguments: []apiValues.PlcValue{
				NewPlcUDINT(0),
				NewPlcUDINT(64),
				NewPlcUDINT(255),
				NewPlcUDINT(math.MaxUint32),
			},
		},
		{
			name: apiValues.UINT,
			arguments: []apiValues.PlcValue{
				NewPlcUINT(0),
				NewPlcUINT(64),
				NewPlcUINT(255),
				NewPlcUINT(math.MaxUint16),
			},
		},
		{
			name: apiValues.USINT,
			arguments: []apiValues.PlcValue{
				NewPlcUSINT(0),
				NewPlcUSINT(64),
				NewPlcUSINT(255),
				NewPlcUSINT(math.MaxUint8),
			},
		},
		{
			name: apiValues.ULINT,
			arguments: []apiValues.PlcValue{
				NewPlcULINT(0),
				NewPlcULINT(64),
				NewPlcULINT(math.MaxUint8),
			},
		},
		{
			name: apiValues.WCHAR,
			arguments: []apiValues.PlcValue{
				NewPlcWCHAR("a"),
			},
		},
		{
			name: apiValues.WORD,
			arguments: []apiValues.PlcValue{
				NewPlcWORD(0),
				NewPlcWORD(64),
				NewPlcWORD(255),
				NewPlcWORD(math.MaxUint16),
			},
		},
		{
			name: apiValues.WSTRING,
			arguments: []apiValues.PlcValue{
				NewPlcWSTRING("hurz"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name.String(), func(t *testing.T) {
			for _, argument := range tt.arguments {
				argumentCopy := argument
				t.Run(fmt.Sprintf("%s", argument), func(t *testing.T) {
					PlcValueType := reflect.TypeOf((*apiValues.PlcValue)(nil)).Elem()
					methods := make(map[string]reflect.Method)
					for i := 0; i < PlcValueType.NumMethod(); i++ {
						method := PlcValueType.Method(i)
						methods[method.Name] = method
					}

					for methodName := range methods {
						if strings.HasPrefix(methodName, "Is") {
							queryType := strings.TrimPrefix(methodName, "Is")
							t.Run(queryType, func(t *testing.T) {
								getMethod := methods[fmt.Sprintf("Get%s", queryType)]
								isA := reflect.ValueOf(argumentCopy).MethodByName(methodName).Call([]reflect.Value{})[0].Bool()
								t.Logf("%s() == %t", methodName, isA)
								var emptyMethod reflect.Method
								if isA && getMethod != emptyMethod {
									t.Logf("Calling %s()", getMethod.Name)
									value := reflect.ValueOf(argumentCopy).MethodByName(getMethod.Name).Call([]reflect.Value{})[0]
									assert.True(t, value.IsValid())
									t.Logf("Value: %v", value)
								}
							})
						}
					}
				})
			}
		})
	}
}
