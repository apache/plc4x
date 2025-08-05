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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInlineIf(t *testing.T) {
	type args struct {
		test bool
		a    func() any
		b    func() any
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "take a",
			args: args{
				test: true,
				a: func() any {
					t.Log("alright")
					return "you got it"
				},
			},
			want: "you got it",
		},
		{
			name: "take b",
			args: args{
				test: false,
				b: func() any {
					t.Log("alright")
					return "you got it"
				},
			},
			want: "you got it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, InlineIf(tt.args.test, tt.args.a, tt.args.b), "InlineIf(%v, func(), func())", tt.args.test)
		})
	}
}
