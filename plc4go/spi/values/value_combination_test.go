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
	type argument interface {
		apiValues.PlcValue
		utils.Serializable
	}
	tests := []struct {
		name      apiValues.PlcValueType
		arguments []argument
	}{
		/*{
			name: apiValues.BINT,
			arguments: []argument{
				NewPlcBINT(big.NewInt(0)),
				NewPlcBINT(big.NewInt(64)),
				NewPlcBINT(big.NewInt(255)),
				NewPlcBINT(big.NewInt(math.MinInt64)),
				NewPlcBINT(big.NewInt(math.MaxInt64)),
			},
		},*/
		{
			name: apiValues.BOOL,
			arguments: []argument{
				NewPlcBOOL(true),
				NewPlcBOOL(false),
			},
		},
		/*{
			name: apiValues.BREAL,
			arguments: []argument{
				NewPlcBREAL(big.NewFloat(0)),
				NewPlcBREAL(big.NewFloat(64)),
				NewPlcBREAL(big.NewFloat(255)),
				NewPlcBREAL(big.NewFloat(math.MinInt64)),
				NewPlcBREAL(big.NewFloat(math.MaxInt64)),
			},
		},*/
		{
			name: apiValues.BYTE,
			arguments: []argument{
				NewPlcBYTE(0),
				NewPlcBYTE(1),
				NewPlcBYTE(64),
				NewPlcBYTE(255),
			},
		},
		{
			name: apiValues.CHAR,
			arguments: []argument{
				NewPlcCHAR(""),
			},
		},
		{
			name: apiValues.DATE,
			arguments: []argument{
				NewPlcDATE(time.Now()),
			},
		},
		{
			name: apiValues.DATE_AND_TIME,
			arguments: []argument{
				NewPlcDATE_AND_TIME(time.Now()),
			},
		},
		{
			name: apiValues.DINT,
			arguments: []argument{
				NewPlcDINT(math.MinInt32),
				NewPlcDINT(0),
				NewPlcDINT(1),
				NewPlcDINT(64),
				NewPlcDINT(255),
				NewPlcDINT(math.MaxInt32),
			},
		},
		{
			name: apiValues.DWORD,
			arguments: []argument{
				NewPlcDWORD(0),
				NewPlcDWORD(1),
				NewPlcDWORD(64),
				NewPlcDWORD(255),
				NewPlcDWORD(math.MaxUint32),
			},
		},
		{
			name: apiValues.INT,
			arguments: []argument{
				NewPlcINT(0),
				NewPlcINT(1),
				NewPlcINT(64),
				NewPlcINT(255),
				NewPlcINT(math.MaxInt16),
			},
		},
		{
			name: apiValues.LDATE,
			arguments: []argument{
				NewPlcLDATE(time.Now()),
			},
		},
		{
			name: apiValues.LDATE_AND_TIME,
			arguments: []argument{
				NewPlcLDATE_AND_TIME(time.Now()),
			},
		},
		{
			name: apiValues.LINT,
			arguments: []argument{
				NewPlcLINT(0),
				NewPlcLINT(1),
				NewPlcLINT(64),
				NewPlcLINT(255),
				NewPlcLINT(math.MaxInt64),
			},
		},
		{
			name: apiValues.List,
			arguments: []argument{
				NewPlcList(nil),
				NewPlcList([]apiValues.PlcValue{
					NewPlcBOOL(true),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcBYTE(1),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcCHAR("a"),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcDATE(time.Now()),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcDATE_AND_TIME(time.Now()),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcDINT(1),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcDWORD(1),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcINT(1),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcLDATE(time.Now()),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcLDATE_AND_TIME(time.Now()),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcLINT(1),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcList([]apiValues.PlcValue{}),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcLREAL(1),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcLTIME(time.Now().Sub(time.Now())),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcLTIME_OF_DAY(time.Now()),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcLWORD(1),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcNULL(),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcRawByteArray([]byte{1, 2, 3}),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcREAL(1),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcSINT(1),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcSTRING("asd"),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcStruct(map[string]apiValues.PlcValue{}),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcTIME(time.Now().Sub(time.Now())),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcTIME_OF_DAY(time.Now().Sub(time.Now())),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcUDINT(1),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcUINT(1),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcUSINT(1),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcWCHAR("a"),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcWORD(1),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcWSTRING("a"),
				}),
				NewPlcList([]apiValues.PlcValue{
					NewPlcBOOL(true),
					NewPlcBOOL(false),
				}),
			},
		},
		{
			name: apiValues.LREAL,
			arguments: []argument{
				NewPlcLREAL(0),
				NewPlcLREAL(1),
				NewPlcLREAL(64),
				NewPlcLREAL(255),
				NewPlcLREAL(math.MinInt64),
				NewPlcLREAL(math.MaxInt64),
			},
		},
		{
			name: apiValues.LTIME,
			arguments: []argument{
				NewPlcLTIME(0),
				NewPlcLTIME(1),
				NewPlcLTIME(64),
				NewPlcLTIME(255),
				//NewPlcLTIME(math.MaxUint64),
			},
		},
		{
			name: apiValues.LWORD,
			arguments: []argument{
				NewPlcLWORD(0),
				NewPlcLWORD(1),
				NewPlcLWORD(64),
				NewPlcLWORD(255),
				NewPlcLWORD(0),
				NewPlcLWORD(math.MaxUint64),
			},
		},
		{
			name: apiValues.NULL,
			arguments: []argument{
				NewPlcNULL(),
			},
		},
		{
			name: apiValues.RAW_BYTE_ARRAY,
			arguments: []argument{
				NewPlcRawByteArray(utils.NewReadBufferByteBased([]byte{0x47, 0x11}).GetBytes()),
			},
		},
		{
			name: apiValues.REAL,
			arguments: []argument{
				NewPlcREAL(0),
				NewPlcREAL(1),
				NewPlcREAL(64),
				NewPlcREAL(255),
				NewPlcREAL(math.MaxInt64),
			},
		},
		{
			name: apiValues.Struct,
			arguments: []argument{
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
			arguments: []argument{
				NewPlcSINT(-128),
				NewPlcSINT(-64),
				NewPlcSINT(0),
				NewPlcSINT(1),
				NewPlcSINT(64),
				NewPlcSINT(127),
			},
		},
		{
			name: apiValues.STRING,
			arguments: []argument{
				NewPlcSTRING("Hello Tody"),
				NewPlcSTRING("1"),
				NewPlcSTRING("true"),
				NewPlcSTRING("255"),
				NewPlcSTRING("5.4"),
			},
		},
		{
			name: apiValues.TIME,
			arguments: []argument{
				NewPlcTIME(0),
				NewPlcTIME(1),
				NewPlcTIME(64),
				NewPlcTIME(255),
				NewPlcTIME(math.MaxUint32),
			},
		},
		{
			name: apiValues.TIME_OF_DAY,
			arguments: []argument{
				NewPlcTIME_OF_DAY(time.Now()),
				NewPlcTIME_OF_DAY(0),
				NewPlcTIME_OF_DAY(1),
				NewPlcTIME_OF_DAY(64),
				NewPlcTIME_OF_DAY(255),
				NewPlcTIME_OF_DAY(math.MaxUint32),
			},
		},
		{
			name: apiValues.LTIME_OF_DAY,
			arguments: []argument{
				NewPlcLTIME_OF_DAY(time.Now()),
				NewPlcLTIME_OF_DAY(0),
				NewPlcLTIME_OF_DAY(1),
				NewPlcLTIME_OF_DAY(64),
				NewPlcLTIME_OF_DAY(255),
				NewPlcLTIME_OF_DAY(math.MaxUint32),
			},
		},
		{
			name: apiValues.UDINT,
			arguments: []argument{
				NewPlcUDINT(0),
				NewPlcUDINT(1),
				NewPlcUDINT(64),
				NewPlcUDINT(255),
				NewPlcUDINT(math.MaxUint32),
			},
		},
		{
			name: apiValues.UINT,
			arguments: []argument{
				NewPlcUINT(0),
				NewPlcUINT(1),
				NewPlcUINT(64),
				NewPlcUINT(255),
				NewPlcUINT(math.MaxUint16),
			},
		},
		{
			name: apiValues.USINT,
			arguments: []argument{
				NewPlcUSINT(0),
				NewPlcUSINT(1),
				NewPlcUSINT(64),
				NewPlcUSINT(255),
				NewPlcUSINT(math.MaxUint8),
			},
		},
		{
			name: apiValues.ULINT,
			arguments: []argument{
				NewPlcULINT(0),
				NewPlcULINT(1),
				NewPlcULINT(64),
				NewPlcULINT(math.MaxUint8),
			},
		},
		{
			name: apiValues.WCHAR,
			arguments: []argument{
				NewPlcWCHAR("a"),
			},
		},
		{
			name: apiValues.WORD,
			arguments: []argument{
				NewPlcWORD(0),
				NewPlcWORD(1),
				NewPlcWORD(64),
				NewPlcWORD(255),
				NewPlcWORD(math.MaxUint16),
			},
		},
		{
			name: apiValues.WSTRING,
			arguments: []argument{
				NewPlcWSTRING("hurz"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name.String(), func(t *testing.T) {
			for _, _argument := range tt.arguments {
				argument := _argument
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
								isA := reflect.ValueOf(argument).MethodByName(methodName).Call([]reflect.Value{})[0].Bool()
								t.Logf("%s() == %t", methodName, isA)
								var emptyMethod reflect.Method
								if isA && getMethod != emptyMethod {
									t.Logf("Calling %s()", getMethod.Name)
									value := reflect.ValueOf(argument).MethodByName(getMethod.Name).Call([]reflect.Value{})[0]
									assert.True(t, value.IsValid())
									t.Logf("Value: %v", value)
								}
							})
						}
					}

					t.Run("Serialize", func(t *testing.T) {
						serialize, err := argument.Serialize()
						assert.NoError(t, err)
						_ = serialize
					})

					t.Run("Raw", func(t *testing.T) {
						defer func() {
							if err := recover(); err != nil {
								t.Log(err)
								t.Skip("this panics and there is no was to check if it is supported") // TODO: fix that
							}
						}()
						argument.GetRaw()
					})

					t.Run("Bools", func(t *testing.T) {
						defer func() {
							if err := recover(); err != nil {
								t.Log(err)
								t.Skip("this panics and there is no was to check if it is supported") // TODO: fix that
							}
						}()
						argument.GetBoolLength()
						argument.GetBoolAt(0)
						argument.GetBoolAt(1)
						argument.GetBoolArray()
					})

					t.Run("Value Type", func(t *testing.T) {
						defer func() {
							if err := recover(); err != nil {
								t.Log(err)
								t.Skip("this panics and there is no was to check if it is supported") // TODO: fix that
							}
						}()
						argument.GetPlcValueType()
					})
				})
			}
		})
	}
}

func TestExtraTests(t *testing.T) {
	// These tests are for extra methods not covered by TestCombinations
	t.Run("PLCDate", func(t *testing.T) {
		t.Run("constructors", func(t *testing.T) {
			t.Run("uint16", func(t *testing.T) {
				assert.NotNil(t, NewPlcDATE(uint16(1)))
			})
			t.Run("uint32", func(t *testing.T) {
				assert.NotNil(t, NewPlcDATE(uint32(1)))
			})
			t.Run("FromSecondsSinceEpoch", func(t *testing.T) {
				assert.NotNil(t, NewPlcDATEFromSecondsSinceEpoch(1))
			})
			t.Run("FromDaysSinceEpoch", func(t *testing.T) {
				assert.NotNil(t, NewPlcDATEFromDaysSinceEpoch(1))
			})
		})
		t.Run("GetSecondsSinceEpoch", func(t *testing.T) {
			NewPlcDATE(time.Now()).GetSecondsSinceEpoch()
		})
		t.Run("GetDaysSinceEpoch", func(t *testing.T) {
			NewPlcDATE(time.Now()).GetDaysSinceEpoch()
		})
		t.Run("GetDaysSinceSiemensEpoch", func(t *testing.T) {
			NewPlcDATE(time.Now()).GetDaysSinceSiemensEpoch()
		})
	})
	t.Run("PlcLDATE_AND_TIME", func(t *testing.T) {
		t.Run("constructors", func(t *testing.T) {
			t.Run("uint16", func(t *testing.T) {
				assert.NotNil(t, NewPlcLDATE_AND_TIME(uint16(1)))
			})
			t.Run("uint32", func(t *testing.T) {
				assert.NotNil(t, NewPlcLDATE_AND_TIME(uint32(1)))
			})
			t.Run("FromNanosecondsSinceEpoch", func(t *testing.T) {
				assert.NotNil(t, NewPlcLDATE_AND_TIMEFromNanosecondsSinceEpoch(1))
			})
		})
		t.Run("GetNanosecondsSinceEpoch", func(t *testing.T) {
			NewPlcLDATE_AND_TIME(time.Now()).GetNanosecondsSinceEpoch()
		})
	})
	t.Run("PlcLDATE", func(t *testing.T) {
		t.Run("constructors", func(t *testing.T) {
			t.Run("uint16", func(t *testing.T) {
				assert.NotNil(t, NewPlcLDATE(uint16(1)))
			})
			t.Run("uint64", func(t *testing.T) {
				assert.NotNil(t, NewPlcLDATE(uint64(1)))
			})
			t.Run("FromNanosecondsSinceEpoch", func(t *testing.T) {
				assert.NotNil(t, NewPlcLDATEFromNanosecondsSinceEpoch(1))
			})
		})
		t.Run("GetNanosecondsSinceEpoch", func(t *testing.T) {
			NewPlcLDATE(time.Now()).GetNanosecondsSinceEpoch()
		})
	})
	t.Run("PlcList", func(t *testing.T) {
		t.Run("GetLength", func(t *testing.T) {
			NewPlcList([]apiValues.PlcValue{NewPlcNULL()}).GetLength()
		})
		t.Run("GetIndex", func(t *testing.T) {
			NewPlcList([]apiValues.PlcValue{NewPlcNULL()}).GetIndex(0)
		})
		t.Run("GetKeys", func(t *testing.T) {
			NewPlcList([]apiValues.PlcValue{NewPlcStruct(map[string]apiValues.PlcValue{})}).GetKeys()
		})
		t.Run("HasKey", func(t *testing.T) {
			NewPlcList([]apiValues.PlcValue{NewPlcStruct(map[string]apiValues.PlcValue{})}).HasKey("")
		})
		t.Run("GetValue", func(t *testing.T) {
			NewPlcList([]apiValues.PlcValue{NewPlcStruct(map[string]apiValues.PlcValue{})}).GetValue("")
		})
	})
	t.Run("PlcLTIME", func(t *testing.T) {
		t.Run("FromNanoseconds", func(t *testing.T) {
			assert.NotNil(t, NewPlcLTIMEFromNanoseconds(0))
		})
		t.Run("GetNanoseconds", func(t *testing.T) {
			NewPlcLTIMEFromNanoseconds(0).GetNanoseconds()
		})
	})
	t.Run("PlcLTIME_OF_DAY", func(t *testing.T) {
		t.Run("constructors", func(t *testing.T) {
			t.Run("uint16", func(t *testing.T) {
				assert.NotNil(t, NewPlcLTIME_OF_DAY(uint16(1)))
			})
			t.Run("uint64", func(t *testing.T) {
				assert.NotNil(t, NewPlcLTIME_OF_DAY(uint64(1)))
			})
			t.Run("FromNanosecondsSinceMidnight", func(t *testing.T) {
				assert.NotNil(t, NewPlcLTIME_OF_DAYFromNanosecondsSinceMidnight(1))
			})
		})
		t.Run("GetNanosecondsSinceEpoch", func(t *testing.T) {
			NewPlcLTIME_OF_DAY(time.Now()).GetNanosecondsSinceMidnight()
		})
	})
	t.Run("PlcRawByteArray", func(t *testing.T) {
		t.Run("GetLength", func(t *testing.T) {
			NewPlcRawByteArray([]byte{1, 2, 3, 4}).GetLength()
		})
		t.Run("GetIndex", func(t *testing.T) {
			NewPlcRawByteArray([]byte{1, 2, 3, 4}).GetIndex(0)
		})
	})
	t.Run("PlcStruct", func(t *testing.T) {
		t.Run("GetKeys", func(t *testing.T) {
			NewPlcStruct(map[string]apiValues.PlcValue{
				"a": NewPlcNULL(),
				"b": NewPlcNULL(),
				"c": NewPlcNULL(),
			}).GetKeys()
		})
		t.Run("HasKey", func(t *testing.T) {
			NewPlcStruct(map[string]apiValues.PlcValue{
				"a": NewPlcNULL(),
				"b": NewPlcNULL(),
				"c": NewPlcNULL(),
			}).HasKey("a")
		})
		t.Run("GetValue", func(t *testing.T) {
			NewPlcStruct(map[string]apiValues.PlcValue{
				"a": NewPlcNULL(),
				"b": NewPlcNULL(),
				"c": NewPlcNULL(),
			}).GetValue("a")
		})
	})
	t.Run("PlcTIME", func(t *testing.T) {
		t.Run("FromMilliseconds", func(t *testing.T) {
			assert.NotNil(t, NewPlcTIMEFromMilliseconds(1))
		})
		t.Run("GetMilliseconds", func(t *testing.T) {
			NewPlcTIMEFromMilliseconds(1).GetMilliseconds()
		})
	})
	t.Run("PlcTIME_OF_DAY", func(t *testing.T) {
		t.Run("constructors", func(t *testing.T) {
			t.Run("uint16", func(t *testing.T) {
				assert.NotNil(t, NewPlcTIME_OF_DAY(uint16(1)))
			})
			t.Run("uint32", func(t *testing.T) {
				assert.NotNil(t, NewPlcTIME_OF_DAY(uint32(1)))
			})
			t.Run("FromMillisecondsSinceMidnight(", func(t *testing.T) {
				assert.NotNil(t, NewPlcTIME_OF_DAYFromMillisecondsSinceMidnight(1))
			})
		})
		t.Run("GetMillisecondsSinceMidnight", func(t *testing.T) {
			NewPlcTIME_OF_DAY(time.Now()).GetMillisecondsSinceMidnight()
		})
	})
}
