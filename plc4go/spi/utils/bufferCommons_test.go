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
	"container/list"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBufferCommons_ExtractAdditionalStringRepresentation(t *testing.T) {
	type args struct {
		readerWriterArgs []WithReaderWriterArgs
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "extract nothing",
		},
		{
			name: "extract the argument",
			args: args{
				readerWriterArgs: []WithReaderWriterArgs{
					withAdditionalStringRepresentation{readerWriterArg: readerWriterArg{readerArg{}, writerArg{}}, stringRepresentation: "plc4xftw"},
				},
			},
			want: "plc4xftw",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BufferCommons{}
			assert.Equalf(t, tt.want, b.ExtractAdditionalStringRepresentation(tt.args.readerWriterArgs...), "ExtractAdditionalStringRepresentation(%v)", tt.args.readerWriterArgs)
		})
	}
}

func TestBufferCommons_IsToBeRenderedAsList(t *testing.T) {
	type args struct {
		readerWriterArgs []WithReaderWriterArgs
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "nope",
		},
		{
			name: "it is not",
			args: args{
				readerWriterArgs: []WithReaderWriterArgs{
					withRenderAsList{readerWriterArg: readerWriterArg{readerArg{}, writerArg{}}},
				},
			},
		},
		{
			name: "it is",
			args: args{
				readerWriterArgs: []WithReaderWriterArgs{
					withRenderAsList{readerWriterArg: readerWriterArg{readerArg{}, writerArg{}}, renderAsList: true},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BufferCommons{}
			assert.Equalf(t, tt.want, b.IsToBeRenderedAsList(tt.args.readerWriterArgs...), "IsToBeRenderedAsList(%v)", tt.args.readerWriterArgs)
		})
	}
}

func TestBufferCommons_SanitizeLogicalName(t *testing.T) {
	type args struct {
		logicalName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "nothing get's a value",
			want: "value",
		},
		{
			name: "something stays something",
			args: args{
				logicalName: "something",
			},
			want: "something",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BufferCommons{}
			assert.Equalf(t, tt.want, b.SanitizeLogicalName(tt.args.logicalName), "SanitizeLogicalName(%v)", tt.args.logicalName)
		})
	}
}

func TestStack_Empty(t *testing.T) {
	type fields struct {
		List list.List
	}
	tests := []struct {
		name        string
		fields      fields
		want        bool
		stackAssert func(t *testing.T, stack *Stack)
	}{
		{
			name: "it is empty",
			want: true,
		},
		{
			name: "it is not empty",
			fields: fields{
				List: func() list.List {
					l := list.List{}
					l.PushBack("boink")
					return l
				}(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Stack{
				List: tt.fields.List,
			}
			assert.Equalf(t, tt.want, s.Empty(), "Empty()")
			if tt.stackAssert != nil {
				tt.stackAssert(t, s)
			}
		})
	}
}

func TestStack_Peek(t *testing.T) {
	type fields struct {
		List list.List
	}
	tests := []struct {
		name        string
		fields      fields
		want        any
		stackAssert func(t *testing.T, stack *Stack)
	}{
		{
			name: "empty",
		},
		{
			name: "not empty",
			fields: fields{
				List: func() list.List {
					l := list.List{}
					l.PushBack("boink")
					return l
				}(),
			},
			want: "boink",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Stack{
				List: tt.fields.List,
			}
			assert.Equalf(t, tt.want, s.Peek(), "Peek()")
			if tt.stackAssert != nil {
				tt.stackAssert(t, s)
			}
		})
	}
}

func TestStack_Pop(t *testing.T) {
	type fields struct {
		List list.List
	}
	tests := []struct {
		name        string
		fields      fields
		setupAssert func(t *testing.T, stack *Stack)
		want        any
		stackAssert func(t *testing.T, stack *Stack)
	}{
		{
			name: "nothing to pop",
			stackAssert: func(t *testing.T, stack *Stack) {
				assert.Equal(t, stack.Len(), 0)
			},
		},
		{
			name: "something to pop",
			setupAssert: func(t *testing.T, stack *Stack) {
				stack.PushBack("boink")
			},
			want: "boink",
			stackAssert: func(t *testing.T, stack *Stack) {
				assert.Equal(t, stack.Len(), 0)
			},
		},
		{
			name: "something to pop with more to pop",
			setupAssert: func(t *testing.T, stack *Stack) {
				stack.PushBack("boink")
				stack.PushBack("boink")
			},
			want: "boink",
			stackAssert: func(t *testing.T, stack *Stack) {
				assert.Equal(t, stack.Len(), 1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Stack{
				List: tt.fields.List,
			}
			if tt.setupAssert != nil {
				tt.setupAssert(t, s)
			}
			assert.Equalf(t, tt.want, s.Pop(), "Pop()")
			if tt.stackAssert != nil {
				tt.stackAssert(t, s)
			}
		})
	}
}

func TestStack_Push(t *testing.T) {
	type fields struct {
		List list.List
	}
	type args struct {
		value any
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        any
		stackAssert func(t *testing.T, stack *Stack)
	}{
		{
			name: "push nothing",
			stackAssert: func(t *testing.T, stack *Stack) {
				assert.Equal(t, stack.Len(), 1)
			},
		},
		{
			name: "push something",
			args: args{
				value: "boink",
			},
			want: "boink",
			stackAssert: func(t *testing.T, stack *Stack) {
				assert.Equal(t, stack.Len(), 1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Stack{
				List: tt.fields.List,
			}
			assert.Equalf(t, tt.want, s.Push(tt.args.value), "Push(%v)", tt.args.value)
			if tt.stackAssert != nil {
				tt.stackAssert(t, s)
			}
		})
	}
}
