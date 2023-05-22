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
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateArrayContext(t *testing.T) {
	type args struct {
		ctx      context.Context
		numItems int
		curItem  int
	}
	tests := []struct {
		name         string
		args         args
		wantVerifier func(t *testing.T, ctx context.Context) bool
	}{
		{
			name: "Create one",
			args: args{
				ctx: context.Background(),
			},
			wantVerifier: func(t *testing.T, ctx context.Context) bool {
				_arrayInfo := ctx.Value(keyArrayInfo)
				assert.NotNil(t, _arrayInfo)
				assert.IsType(t, arrayInfo{}, _arrayInfo)
				assert.Equal(t, 0, _arrayInfo.(arrayInfo).numItems)
				assert.Equal(t, 0, _arrayInfo.(arrayInfo).curItem)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateArrayContext(tt.args.ctx, tt.args.numItems, tt.args.curItem)
			assert.Truef(t, tt.wantVerifier(t, got), "CreateArrayContext(%v, %v, %v)", tt.args.ctx, tt.args.numItems, tt.args.curItem)
		})
	}
}

func TestGetCurItemFromContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		args      args
		want      int
		wantPanic bool
	}{
		{
			name: "key not set",
			args: args{
				ctx: context.Background(),
			},
			wantPanic: true,
		},
		{
			name: "key present but wrong value",
			args: args{
				ctx: context.WithValue(context.Background(), keyArrayInfo, nil),
			},
			wantPanic: true,
		},
		{
			name: "key present",
			args: args{
				ctx: context.WithValue(context.Background(), keyArrayInfo, arrayInfo{
					numItems: 1,
					curItem:  2,
				}),
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if tt.wantPanic {
					assert.NotNil(t, recover(), "we expected a panic")
				} else {
					assert.Nil(t, recover(), "we don't expected a panic")
				}
			}()
			assert.Equalf(t, tt.want, GetCurItemFromContext(tt.args.ctx), "GetCurItemFromContext(%v)", tt.args.ctx)
		})
	}
}

func TestGetLastItemFromContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		args      args
		want      bool
		wantPanic bool
	}{
		{
			name: "key not set",
			args: args{
				ctx: context.Background(),
			},
			wantPanic: true,
		},
		{
			name: "key present but wrong value",
			args: args{
				ctx: context.WithValue(context.Background(), keyArrayInfo, nil),
			},
			wantPanic: true,
		},
		{
			name: "key present",
			args: args{
				ctx: context.WithValue(context.Background(), keyArrayInfo, arrayInfo{
					numItems: 2,
					curItem:  1,
				}),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if tt.wantPanic {
					assert.NotNil(t, recover(), "we expected a panic")
				} else {
					assert.Nil(t, recover(), "we don't expected a panic")
				}
			}()
			assert.Equalf(t, tt.want, GetLastItemFromContext(tt.args.ctx), "GetLastItemFromContext(%v)", tt.args.ctx)
		})
	}
}

func TestGetNumItemsFromContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		args      args
		want      int
		wantPanic bool
	}{
		{
			name: "key not set",
			args: args{
				ctx: context.Background(),
			},
			wantPanic: true,
		},
		{
			name: "key present but wrong value",
			args: args{
				ctx: context.WithValue(context.Background(), keyArrayInfo, nil),
			},
			wantPanic: true,
		},
		{
			name: "key present",
			args: args{
				ctx: context.WithValue(context.Background(), keyArrayInfo, arrayInfo{
					numItems: 1,
					curItem:  2,
				}),
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if tt.wantPanic {
					assert.NotNil(t, recover(), "we expected a panic")
				} else {
					assert.Nil(t, recover(), "we don't expected a panic")
				}
			}()
			assert.Equalf(t, tt.want, GetNumItemsFromContext(tt.args.ctx), "GetNumItemsFromContext(%v)", tt.args.ctx)
		})
	}
}
