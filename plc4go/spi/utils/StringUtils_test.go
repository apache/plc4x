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

func TestRandomString(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name       string
		args       args
		wantAssert func(t *testing.T, actual string) bool
	}{
		{
			name: "simple test",
			wantAssert: func(t *testing.T, actual string) bool {
				return assert.Equal(t, 0, len(actual))
			},
		},
		{
			name: "fixed length",
			args: args{length: 60},
			wantAssert: func(t *testing.T, actual string) bool {
				t.Log(actual)
				return assert.Equal(t, 60, len(actual))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RandomString(tt.args.length)
			if tt.wantAssert != nil {
				assert.Truef(t, tt.wantAssert(t, got), "RandomString(%v)", tt.args.length)
			}
		})
	}
}
