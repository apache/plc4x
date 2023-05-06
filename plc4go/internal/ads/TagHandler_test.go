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

package ads

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"

	model2 "github.com/apache/plc4x/plc4go/internal/ads/model"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	model3 "github.com/apache/plc4x/plc4go/spi/model"
)

func TestTagHandler_ParseQuery(t *testing.T) {
	type tags struct {
		directAdsStringTag *regexp.Regexp
		directAdsTag       *regexp.Regexp
		symbolicAdsTag     *regexp.Regexp
	}
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		tags    tags
		args    args
		want    model.PlcTag
		wantErr bool
	}{
		// All tests without any array notation.
		{
			name: "simple direct numeric address",
			args: args{
				query: "1234/5678:BOOL",
			},
			want: model2.DirectPlcTag{
				IndexGroup:   1234,
				IndexOffset:  5678,
				ValueType:    values.BOOL,
				StringLength: model2.NONE,
			},
		},
		{
			name: "simple direct hex address",
			args: args{
				query: "0x04D2/0x162E:BOOL",
			},
			want: model2.DirectPlcTag{
				IndexGroup:   1234,
				IndexOffset:  5678,
				ValueType:    values.BOOL,
				StringLength: model2.NONE,
			},
		},
		{
			name: "simple direct numeric string address",
			args: args{
				query: "1234/5678:STRING(80)",
			},
			want: model2.DirectPlcTag{
				IndexGroup:   1234,
				IndexOffset:  5678,
				ValueType:    values.STRING,
				StringLength: 80,
			},
		},
		{
			name: "simple direct hex string address",
			args: args{
				query: "0x04D2/0x162E:WSTRING(80)",
			},
			want: model2.DirectPlcTag{
				IndexGroup:   1234,
				IndexOffset:  5678,
				ValueType:    values.WSTRING,
				StringLength: 80,
			},
		},
		{
			name: "simple symbolic address",
			args: args{
				query: "MAIN.testVariable",
			},
			want: model2.SymbolicPlcTag{
				SymbolicAddress: "MAIN.testVariable",
			},
		},
		// All tests with simple array notation.
		{
			name: "simple array direct numeric address",
			args: args{
				query: "1234/5678:BOOL[42]",
			},
			want: model2.DirectPlcTag{
				PlcTag: model2.PlcTag{
					ArrayInfo: []model.ArrayInfo{
						&model3.DefaultArrayInfo{
							LowerBound: 0,
							UpperBound: 42,
						},
					},
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				ValueType:    values.BOOL,
				StringLength: model2.NONE,
			},
		},
		{
			name: "simple array direct hex address",
			args: args{
				query: "0x04D2/0x162E:BOOL[42]",
			},
			want: model2.DirectPlcTag{
				PlcTag: model2.PlcTag{
					ArrayInfo: []model.ArrayInfo{
						&model3.DefaultArrayInfo{
							LowerBound: 0,
							UpperBound: 42,
						},
					},
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				ValueType:    values.BOOL,
				StringLength: model2.NONE,
			},
		},
		{
			name: "simple array direct numeric string address",
			args: args{
				query: "1234/5678:STRING(80)[42]",
			},
			want: model2.DirectPlcTag{
				PlcTag: model2.PlcTag{
					ArrayInfo: []model.ArrayInfo{
						&model3.DefaultArrayInfo{
							LowerBound: 0,
							UpperBound: 42,
						},
					},
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				ValueType:    values.STRING,
				StringLength: 80,
			},
		},
		{
			name: "simple array direct hex string address",
			args: args{
				query: "0x04D2/0x162E:WSTRING(80)[42]",
			},
			want: model2.DirectPlcTag{
				PlcTag: model2.PlcTag{
					ArrayInfo: []model.ArrayInfo{
						&model3.DefaultArrayInfo{
							LowerBound: 0,
							UpperBound: 42,
						},
					},
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				ValueType:    values.WSTRING,
				StringLength: 80,
			},
		},
		{
			name: "simple array symbolic address",
			args: args{
				query: "MAIN.testVariable[42]",
			},
			want: model2.SymbolicPlcTag{
				PlcTag: model2.PlcTag{
					ArrayInfo: []model.ArrayInfo{
						&model3.DefaultArrayInfo{
							LowerBound: 0,
							UpperBound: 42,
						},
					},
				},
				SymbolicAddress: "MAIN.testVariable",
			},
		},
		// All tests with range array notation.
		{
			name: "range array direct numeric address",
			args: args{
				query: "1234/5678:BOOL[23..42]",
			},
			want: model2.DirectPlcTag{
				PlcTag: model2.PlcTag{
					ArrayInfo: []model.ArrayInfo{
						&model3.DefaultArrayInfo{
							LowerBound: 23,
							UpperBound: 42,
						},
					},
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				ValueType:    values.BOOL,
				StringLength: model2.NONE,
			},
		},
		{
			name: "range array direct hex address",
			args: args{
				query: "0x04D2/0x162E:BOOL[23..42]",
			},
			want: model2.DirectPlcTag{
				PlcTag: model2.PlcTag{
					ArrayInfo: []model.ArrayInfo{
						&model3.DefaultArrayInfo{
							LowerBound: 23,
							UpperBound: 42,
						},
					},
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				ValueType:    values.BOOL,
				StringLength: model2.NONE,
			},
		},
		{
			name: "range array direct numeric string address",
			args: args{
				query: "1234/5678:STRING(80)[23..42]",
			},
			want: model2.DirectPlcTag{
				PlcTag: model2.PlcTag{
					ArrayInfo: []model.ArrayInfo{
						&model3.DefaultArrayInfo{
							LowerBound: 23,
							UpperBound: 42,
						},
					},
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				ValueType:    values.STRING,
				StringLength: 80,
			},
		},
		{
			name: "range array direct hex string address",
			args: args{
				query: "0x04D2/0x162E:WSTRING(80)[23..42]",
			},
			want: model2.DirectPlcTag{
				PlcTag: model2.PlcTag{
					ArrayInfo: []model.ArrayInfo{
						&model3.DefaultArrayInfo{
							LowerBound: 23,
							UpperBound: 42,
						},
					},
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				ValueType:    values.WSTRING,
				StringLength: 80,
			},
		},
		{
			name: "range array symbolic address",
			args: args{
				query: "MAIN.testVariable[23..42]",
			},
			want: model2.SymbolicPlcTag{
				PlcTag: model2.PlcTag{
					ArrayInfo: []model.ArrayInfo{
						&model3.DefaultArrayInfo{
							LowerBound: 23,
							UpperBound: 42,
						},
					},
				},
				SymbolicAddress: "MAIN.testVariable",
			},
		},
		// All tests with array with offset notation.
		{
			name: "array with offset direct numeric address",
			args: args{
				query: "1234/5678:BOOL[23:42]",
			},
			want: model2.DirectPlcTag{
				PlcTag: model2.PlcTag{
					ArrayInfo: []model.ArrayInfo{
						&model3.DefaultArrayInfo{
							LowerBound: 23,
							UpperBound: 65,
						},
					},
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				ValueType:    values.BOOL,
				StringLength: model2.NONE,
			},
		},
		{
			name: "array with offset direct hex address",
			args: args{
				query: "0x04D2/0x162E:BOOL[23:42]",
			},
			want: model2.DirectPlcTag{
				PlcTag: model2.PlcTag{
					ArrayInfo: []model.ArrayInfo{
						&model3.DefaultArrayInfo{
							LowerBound: 23,
							UpperBound: 65,
						},
					},
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				ValueType:    values.BOOL,
				StringLength: model2.NONE,
			},
		},
		{
			name: "array with offset direct numeric string address",
			args: args{
				query: "1234/5678:STRING(80)[23:42]",
			},
			want: model2.DirectPlcTag{
				PlcTag: model2.PlcTag{
					ArrayInfo: []model.ArrayInfo{
						&model3.DefaultArrayInfo{
							LowerBound: 23,
							UpperBound: 65,
						},
					},
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				ValueType:    values.STRING,
				StringLength: 80,
			},
		},
		{
			name: "array with offset direct hex string address",
			args: args{
				query: "0x04D2/0x162E:WSTRING(80)[23:42]",
			},
			want: model2.DirectPlcTag{
				PlcTag: model2.PlcTag{
					ArrayInfo: []model.ArrayInfo{
						&model3.DefaultArrayInfo{
							LowerBound: 23,
							UpperBound: 65,
						},
					},
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				ValueType:    values.WSTRING,
				StringLength: 80,
			},
		},
		{
			name: "array with offset symbolic address",
			args: args{
				query: "MAIN.testVariable[23:42]",
			},
			want: model2.SymbolicPlcTag{
				PlcTag: model2.PlcTag{
					ArrayInfo: []model.ArrayInfo{
						&model3.DefaultArrayInfo{
							LowerBound: 23,
							UpperBound: 65,
						},
					},
				},
				SymbolicAddress: "MAIN.testVariable",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewTagHandler()
			got, err := m.ParseTag(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("ParseQuery() got = %v, want %v", got, tt.want)
			}
		})
	}
}
