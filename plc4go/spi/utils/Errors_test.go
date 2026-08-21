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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTimeoutError(t *testing.T) {
	type args struct {
		timeout time.Duration
	}
	tests := []struct {
		name string
		args args
		want TimeoutError
	}{
		{
			name: "create it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewTimeoutError(tt.args.timeout), "NewTimeoutError(%v)", tt.args.timeout)
		})
	}
}

func TestParseAssertError_Error(t *testing.T) {
	type fields struct {
		Message string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "no message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ParseAssertError{
				Message: tt.fields.Message,
			}
			assert.Equalf(t, tt.want, e.Error(), "Error()")
		})
	}
}

func TestParseAssertError_Is(t *testing.T) {
	type fields struct {
		Message string
	}
	type args struct {
		target error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "nay",
		},
		{
			name: "yay",
			args: args{
				target: ParseAssertError{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ParseAssertError{
				Message: tt.fields.Message,
			}
			assert.Equalf(t, tt.want, e.Is(tt.args.target), "Is(%v)", tt.args.target)
		})
	}
}

func TestParseValidationError_Error(t *testing.T) {
	type fields struct {
		Message string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "no message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ParseValidationError{
				Message: tt.fields.Message,
			}
			assert.Equalf(t, tt.want, e.Error(), "Error()")
		})
	}
}

func TestParseValidationError_Is(t *testing.T) {
	type fields struct {
		Message string
	}
	type args struct {
		target error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "nay",
		},
		{
			name: "yay",
			args: args{
				target: ParseValidationError{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ParseValidationError{
				Message: tt.fields.Message,
			}
			assert.Equalf(t, tt.want, e.Is(tt.args.target), "Is(%v)", tt.args.target)
		})
	}
}

func TestTimeoutError_Error(t1 *testing.T) {
	type fields struct {
		timeout time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "no time",
			want: "got timeout after 0s",
		},
		{
			name: "some time",
			fields: fields{
				timeout: 3 * time.Hour,
			},
			want: "got timeout after 3h0m0s",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := TimeoutError{
				timeout: tt.fields.timeout,
			}
			assert.Equalf(t1, tt.want, t.Error(), "Error()")
		})
	}
}

func TestTimeoutError_Is(t1 *testing.T) {
	type fields struct {
		timeout time.Duration
	}
	type args struct {
		target error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "nay",
		},
		{
			name: "yay",
			args: args{
				target: ParseValidationError{},
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := TimeoutError{
				timeout: tt.fields.timeout,
			}
			assert.Equalf(t1, tt.want, t.Is(tt.args.target), "Is(%v)", tt.args.target)
		})
	}
}

// A request that runs out of time reports it either as the codec's expectation expiring or as the
// request context reaching its deadline, and which of the two a caller sees is a race. Both have to
// read as a timeout, and nothing else may.
func TestIsTimeoutError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"the codec's expectation expiring", NewTimeoutError(250 * time.Millisecond), true},
		{"the request context's deadline", context.DeadlineExceeded, true},
		{"either one, wrapped by a caller", fmt.Errorf("error waiting for read response: %w", NewTimeoutError(time.Second)), true},
		{"a deadline wrapped by a caller", fmt.Errorf("read aborted: %w", context.DeadlineExceeded), true},
		{"a cancelled request is not a timeout", context.Canceled, false},
		{"an unrelated failure is not a timeout", fmt.Errorf("connection reset"), false},
		{"no error at all", nil, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, IsTimeoutError(test.err))
		})
	}
}
