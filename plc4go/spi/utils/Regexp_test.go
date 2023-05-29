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
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestGetSubgroupMatches(t *testing.T) {
	type args struct {
		r     *regexp.Regexp
		query string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "plc4x",
			args: args{
				r:     regexp.MustCompile("plc(?P<aNumber>4)x"),
				query: "plc4x",
			},
			want: map[string]string{
				"aNumber": "4",
			},
		},
		{
			name: "plc4x not found",
			args: args{
				r:     regexp.MustCompile("plc(?P<aNumber>4)x"),
				query: "plc5x",
			},
			want: nil,
		},
		{
			name: "plc4x not named group",
			args: args{
				r:     regexp.MustCompile("plc(4)(x)"),
				query: "plc4x",
			},
			want: map[string]string{
				"_1": "4",
				"_2": "x",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetSubgroupMatches(tt.args.r, tt.args.query), "GetSubgroupMatches(%v, %v)", tt.args.r, tt.args.query)
		})
	}
}
