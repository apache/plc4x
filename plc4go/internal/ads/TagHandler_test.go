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
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/internal/ads/model"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
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
		want    apiModel.PlcTag
		wantErr bool
	}{
		// All tests without any array notation.
		{
			name: "simple direct numeric address",
			args: args{
				query: "1234/5678:BOOL",
			},
			want: model.DirectPlcTag{
				IndexGroup:   1234,
				IndexOffset:  5678,
				ValueType:    apiValues.BOOL,
				StringLength: model.NONE,
			},
		},
		{
			name: "simple direct hex address",
			args: args{
				query: "0x04D2/0x162E:BOOL",
			},
			want: model.DirectPlcTag{
				IndexGroup:   1234,
				IndexOffset:  5678,
				ValueType:    apiValues.BOOL,
				StringLength: model.NONE,
			},
		},
		{
			name: "simple direct numeric string address",
			args: args{
				query: "1234/5678:STRING(80)",
			},
			want: model.DirectPlcTag{
				IndexGroup:   1234,
				IndexOffset:  5678,
				ValueType:    apiValues.STRING,
				StringLength: 80,
			},
		},
		{
			name: "simple direct hex string address",
			args: args{
				query: "0x04D2/0x162E:WSTRING(80)",
			},
			want: model.DirectPlcTag{
				IndexGroup:   1234,
				IndexOffset:  5678,
				ValueType:    apiValues.WSTRING,
				StringLength: 80,
			},
		},
		{
			name: "simple symbolic address",
			args: args{
				query: "MAIN.testVariable",
			},
			want: model.SymbolicPlcTag{
				SymbolicAddress: "MAIN.testVariable",
			},
		},
		// All tests with simple array notation.
		{
			name: "simple array direct numeric address",
			args: args{
				query: "1234/5678[0..41]:BOOL",
			},
			want: model.DirectPlcTag{
				PlcTag: model.PlcTag{
					ArrayInfo: []apiModel.ArrayInfo{
						&spiModel.DefaultArrayInfo{
							LowerBound: 0,
							UpperBound: 41,
							Range:      true,
						},
					},
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				ValueType:    apiValues.BOOL,
				StringLength: model.NONE,
			},
		},
		{
			name: "simple array direct hex address",
			args: args{
				query: "0x04D2/0x162E[0..41]:BOOL",
			},
			want: model.DirectPlcTag{
				PlcTag: model.PlcTag{
					ArrayInfo: []apiModel.ArrayInfo{
						&spiModel.DefaultArrayInfo{
							LowerBound: 0,
							UpperBound: 41,
							Range:      true,
						},
					},
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				ValueType:    apiValues.BOOL,
				StringLength: model.NONE,
			},
		},
		{
			name: "simple array direct numeric string address",
			args: args{
				query: "1234/5678[0..41]:STRING(80)",
			},
			want: model.DirectPlcTag{
				PlcTag: model.PlcTag{
					ArrayInfo: []apiModel.ArrayInfo{
						&spiModel.DefaultArrayInfo{
							LowerBound: 0,
							UpperBound: 41,
							Range:      true,
						},
					},
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				ValueType:    apiValues.STRING,
				StringLength: 80,
			},
		},
		{
			name: "simple array direct hex string address",
			args: args{
				query: "0x04D2/0x162E[0..41]:WSTRING(80)",
			},
			want: model.DirectPlcTag{
				PlcTag: model.PlcTag{
					ArrayInfo: []apiModel.ArrayInfo{
						&spiModel.DefaultArrayInfo{
							LowerBound: 0,
							UpperBound: 41,
							Range:      true,
						},
					},
				},
				IndexGroup:   1234,
				IndexOffset:  5678,
				ValueType:    apiValues.WSTRING,
				StringLength: 80,
			},
		},
		{
			name: "simple array symbolic address",
			args: args{
				query: "MAIN.testVariable[0..41]",
			},
			want: model.SymbolicPlcTag{
				PlcTag: model.PlcTag{
					ArrayInfo: []apiModel.ArrayInfo{
						&spiModel.DefaultArrayInfo{
							LowerBound: 0,
							UpperBound: 41,
							Range:      true,
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
				query: "1234/5678[23..42]:BOOL",
			},
			want: model.DirectPlcTag{
				PlcTag: model.PlcTag{
					ArrayInfo: []apiModel.ArrayInfo{
						&spiModel.DefaultArrayInfo{
							LowerBound: 23,
							UpperBound: 42,
							Range:      true,
						},
					},
				},
				IndexGroup:   1234,
				IndexOffset:  5701,
				ValueType:    apiValues.BOOL,
				StringLength: model.NONE,
			},
		},
		{
			name: "range array direct hex address",
			args: args{
				query: "0x04D2/0x162E[23..42]:BOOL",
			},
			want: model.DirectPlcTag{
				PlcTag: model.PlcTag{
					ArrayInfo: []apiModel.ArrayInfo{
						&spiModel.DefaultArrayInfo{
							LowerBound: 23,
							UpperBound: 42,
							Range:      true,
						},
					},
				},
				IndexGroup:   1234,
				IndexOffset:  5701,
				ValueType:    apiValues.BOOL,
				StringLength: model.NONE,
			},
		},
		{
			name: "range array direct numeric string address",
			args: args{
				query: "1234/5678[23..42]:STRING(80)",
			},
			want: model.DirectPlcTag{
				PlcTag: model.PlcTag{
					ArrayInfo: []apiModel.ArrayInfo{
						&spiModel.DefaultArrayInfo{
							LowerBound: 23,
							UpperBound: 42,
							Range:      true,
						},
					},
				},
				IndexGroup:   1234,
				IndexOffset:  7541,
				ValueType:    apiValues.STRING,
				StringLength: 80,
			},
		},
		{
			name: "range array direct hex string address",
			args: args{
				query: "0x04D2/0x162E[23..42]:WSTRING(80)",
			},
			want: model.DirectPlcTag{
				PlcTag: model.PlcTag{
					ArrayInfo: []apiModel.ArrayInfo{
						&spiModel.DefaultArrayInfo{
							LowerBound: 23,
							UpperBound: 42,
							Range:      true,
						},
					},
				},
				IndexGroup:   1234,
				IndexOffset:  9404,
				ValueType:    apiValues.WSTRING,
				StringLength: 80,
			},
		},
		{
			name: "range array symbolic address",
			args: args{
				query: "MAIN.testVariable[23..42]",
			},
			want: model.SymbolicPlcTag{
				PlcTag: model.PlcTag{
					ArrayInfo: []apiModel.ArrayInfo{
						&spiModel.DefaultArrayInfo{
							LowerBound: 23,
							UpperBound: 42,
							Range:      true,
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
				query: "1234/5678[23..64]:BOOL",
			},
			want: model.DirectPlcTag{
				PlcTag: model.PlcTag{
					ArrayInfo: []apiModel.ArrayInfo{
						&spiModel.DefaultArrayInfo{
							LowerBound: 23,
							UpperBound: 64,
							Range:      true,
						},
					},
				},
				IndexGroup:   1234,
				IndexOffset:  5701,
				ValueType:    apiValues.BOOL,
				StringLength: model.NONE,
			},
		},
		{
			name: "array with offset direct hex address",
			args: args{
				query: "0x04D2/0x162E[23..64]:BOOL",
			},
			want: model.DirectPlcTag{
				PlcTag: model.PlcTag{
					ArrayInfo: []apiModel.ArrayInfo{
						&spiModel.DefaultArrayInfo{
							LowerBound: 23,
							UpperBound: 64,
							Range:      true,
						},
					},
				},
				IndexGroup:   1234,
				IndexOffset:  5701,
				ValueType:    apiValues.BOOL,
				StringLength: model.NONE,
			},
		},
		{
			name: "array with offset direct numeric string address",
			args: args{
				query: "1234/5678[23..64]:STRING(80)",
			},
			want: model.DirectPlcTag{
				PlcTag: model.PlcTag{
					ArrayInfo: []apiModel.ArrayInfo{
						&spiModel.DefaultArrayInfo{
							LowerBound: 23,
							UpperBound: 64,
							Range:      true,
						},
					},
				},
				IndexGroup:   1234,
				IndexOffset:  7541,
				ValueType:    apiValues.STRING,
				StringLength: 80,
			},
		},
		{
			name: "array with offset direct hex string address",
			args: args{
				query: "0x04D2/0x162E[23..64]:WSTRING(80)",
			},
			want: model.DirectPlcTag{
				PlcTag: model.PlcTag{
					ArrayInfo: []apiModel.ArrayInfo{
						&spiModel.DefaultArrayInfo{
							LowerBound: 23,
							UpperBound: 64,
							Range:      true,
						},
					},
				},
				IndexGroup:   1234,
				IndexOffset:  9404,
				ValueType:    apiValues.WSTRING,
				StringLength: 80,
			},
		},
		{
			name: "array with offset symbolic address",
			args: args{
				query: "MAIN.testVariable[23..64]",
			},
			want: model.SymbolicPlcTag{
				PlcTag: model.PlcTag{
					ArrayInfo: []apiModel.ArrayInfo{
						&spiModel.DefaultArrayInfo{
							LowerBound: 23,
							UpperBound: 64,
							Range:      true,
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

// Addresses written before the array notation was unified must not be accepted.
//
// Go's ADS driver diverged from plc4j in two ways: [n] meant a count of n elements rather than
// the element at index n, and [a:b] meant a start and a count. Both are gone, so a Go address
// means what the same Java address means.
func TestTagHandler_LegacyAddressesAreRejected(t *testing.T) {
	handler := NewTagHandler()

	t.Run("the selection may no longer follow the type", func(t *testing.T) {
		// This used to parse as a direct tag. Without a guard it would fall through to the
		// symbolic pattern and silently become a symbol lookup.
		_, err := handler.ParseTag("1234/5678:BOOL[42]")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "0x4020/0[0..3]:DINT", "the message must name the form")
	})

	t.Run("the start-and-count form is gone", func(t *testing.T) {
		_, err := handler.ParseTag("MAIN.testVariable[23:42]")
		require.Error(t, err)
	})

	t.Run("the replacement forms parse", func(t *testing.T) {
		_, err := handler.ParseTag("1234/5678[0..41]:BOOL")
		require.NoError(t, err)
		_, err = handler.ParseTag("MAIN.testVariable[23..64]")
		require.NoError(t, err)
	})
}

// A bare index selects one element and yields a scalar; a range yields an array even when it
// spans one element. This is what a consumer reads GetArrayInfo to decide.
func TestTagHandler_ABareIndexIsAScalar(t *testing.T) {
	handler := NewTagHandler()

	index, err := handler.ParseTag("MAIN.testVariable[4]")
	require.NoError(t, err)
	assert.Empty(t, index.GetArrayInfo(), "one element is a scalar, so there is no dimension to report")
	assert.Equal(t, "MAIN.testVariable[4]", index.GetAddressString(), "and the index it selects is still 4, not a count of four")

	arrayRange, err := handler.ParseTag("MAIN.testVariable[4..4]")
	require.NoError(t, err)
	assert.True(t, arrayRange.GetArrayInfo()[0].IsRange())
	assert.Equal(t, uint32(1), arrayRange.GetArrayInfo()[0].GetSize())
}
