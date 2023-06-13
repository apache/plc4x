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

package transactions

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_completedFuture_AwaitCompletion(t *testing.T) {
	type fields struct {
		err error
	}
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "does nothing",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := completedFuture{
				err: tt.fields.err,
			}
			tt.wantErr(t, c.AwaitCompletion(tt.args.in0), fmt.Sprintf("AwaitCompletion(%v)", tt.args.in0))
		})
	}
}

func Test_completedFuture_Cancel(t *testing.T) {
	type fields struct {
		err error
	}
	type args struct {
		in0 bool
		in1 error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "does nothing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			co := completedFuture{
				err: tt.fields.err,
			}
			co.Cancel(tt.args.in0, tt.args.in1)
		})
	}
}

func Test_completedFuture_String(t *testing.T) {
	type fields struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "gives the error",
			want: "completedFuture{\n\terr: <nil>,\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := completedFuture{
				err: tt.fields.err,
			}
			assert.Equalf(t, tt.want, c.String(), "String()")
		})
	}
}
