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

package pool

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_future_AwaitCompletion(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		args      args
		completer func(*future)
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "completes with error",
			args: args{ctx: context.TODO()},
			completer: func(f *future) {
				f.Cancel(false, errors.New("Uh oh"))
			},
			wantErr: assert.Error,
		},
		{
			name: "completes regular",
			args: args{ctx: context.TODO()},
			completer: func(f *future) {
				time.Sleep(time.Millisecond * 30)
				f.complete()
			},
			wantErr: assert.NoError,
		},
		{
			name: "completes not int time",
			args: args{ctx: func() context.Context {
				deadline, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Millisecond))
				t.Cleanup(cancel)
				return deadline
			}()},
			completer: func(f *future) {
				time.Sleep(time.Millisecond * 300)
			},
			wantErr: assert.Error,
		},
		{
			name: "completes canceled without error",
			args: args{ctx: context.TODO()},
			completer: func(f *future) {
				time.Sleep(time.Millisecond * 300)
				f.Cancel(true, nil)
			},
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				assert.Same(t, Canceled, err)
				return true
			},
		},
		{
			name: "completes canceled with particular error",
			args: args{ctx: context.TODO()},
			completer: func(f *future) {
				time.Sleep(time.Millisecond * 300)
				f.Cancel(true, errors.New("Uh oh"))
			},
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				assert.Equal(t, "Uh oh", err.Error())
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &future{}
			go tt.completer(f)
			tt.wantErr(t, f.AwaitCompletion(tt.args.ctx), fmt.Sprintf("AwaitCompletion(%v)", tt.args.ctx))
		})
	}
}

func Test_future_Cancel(t *testing.T) {
	type args struct {
		interrupt bool
		err       error
	}
	tests := []struct {
		name     string
		args     args
		verifier func(*testing.T, *future)
	}{
		{
			name: "cancel cancels",
			verifier: func(t *testing.T, f *future) {
				assert.True(t, f.cancelRequested.Load())
			},
		},
		{
			name: "cancel with interrupt",
			args: args{
				interrupt: true,
				err:       nil,
			},
			verifier: func(t *testing.T, f *future) {
				assert.True(t, f.cancelRequested.Load())
				assert.False(t, f.errored.Load())
				assert.Nil(t, f.err.Load())
			},
		},
		{
			name: "cancel with err",
			args: args{
				interrupt: true,
				err:       errors.New("Uh Oh"),
			},
			verifier: func(t *testing.T, f *future) {
				assert.True(t, f.cancelRequested.Load())
				assert.True(t, f.errored.Load())
				assert.NotNil(t, f.err.Load())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &future{}
			f.Cancel(tt.args.interrupt, tt.args.err)
			tt.verifier(t, f)
		})
	}
}

func Test_future_String(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "string it",
			want: "future{\n\tcancelRequested: false,\n\tinterruptRequested: false,\n\tcompleted: false,\n\terrored: false,\n\terr: <nil>,\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &future{}
			assert.Equalf(t, tt.want, f.String(), "String()")
		})
	}
}

func Test_future_complete(t *testing.T) {
	tests := []struct {
		name     string
		verifier func(*testing.T, *future)
	}{
		{
			name: "complete completes",
			verifier: func(t *testing.T, f *future) {
				assert.True(t, f.completed.Load())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &future{}
			f.complete()
			tt.verifier(t, f)
		})
	}
}
