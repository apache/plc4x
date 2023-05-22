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
	"testing"
)

func TestUpcastReaderArgs(t *testing.T) {
	type args struct {
		args []WithReaderArgs
	}
	tests := []struct {
		name string
		args args
		want []WithReaderWriterArgs
	}{
		{
			name: "nothing results in nothing",
			want: []WithReaderWriterArgs{},
		},
		{
			name: "only a reader arg",
			args: args{
				[]WithReaderArgs{readerArg{}},
			},
			want: []WithReaderWriterArgs{readerWriterArg{readerArg{}, writerArg{}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, UpcastReaderArgs(tt.args.args...), "UpcastReaderArgs(%v)", tt.args.args)
		})
	}
}

func TestUpcastWriterArgs(t *testing.T) {
	type args struct {
		args []WithWriterArgs
	}
	tests := []struct {
		name string
		args args
		want []WithReaderWriterArgs
	}{
		{
			name: "nothing results in nothing",
			want: []WithReaderWriterArgs{},
		},
		{
			name: "only a reader arg",
			args: args{
				[]WithWriterArgs{writerArg{}},
			},
			want: []WithReaderWriterArgs{readerWriterArg{readerArg{}, writerArg{}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, UpcastWriterArgs(tt.args.args...), "UpcastWriterArgs(%v)", tt.args.args)
		})
	}
}

func TestWithAdditionalStringRepresentation(t *testing.T) {
	type args struct {
		stringRepresentation string
	}
	tests := []struct {
		name string
		args args
		want WithReaderWriterArgs
	}{
		{
			name: "some string representation",
			want: withAdditionalStringRepresentation{readerWriterArg: readerWriterArg{readerArg{}, writerArg{}}, stringRepresentation: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, WithAdditionalStringRepresentation(tt.args.stringRepresentation), "WithAdditionalStringRepresentation(%v)", tt.args.stringRepresentation)
		})
	}
}

func TestWithRenderAsList(t *testing.T) {
	type args struct {
		renderAsList bool
	}
	tests := []struct {
		name string
		args args
		want WithReaderWriterArgs
	}{
		{
			name: "render as list",
			args: args{renderAsList: true},
			want: withRenderAsList{readerWriterArg: readerWriterArg{readerArg{}, writerArg{}}, renderAsList: true},
		},
		{
			name: "render not as list",
			args: args{renderAsList: false},
			want: withRenderAsList{readerWriterArg: readerWriterArg{readerArg{}, writerArg{}}, renderAsList: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, WithRenderAsList(tt.args.renderAsList), "WithRenderAsList(%v)", tt.args.renderAsList)
		})
	}
}

func Test_readerWriterArg_isReaderArgs(t *testing.T) {
	type fields struct {
		readerArg WithReaderArgs
		writerArg WithWriterArgs
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "it is",
			fields: fields{
				readerArg: readerArg{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := readerWriterArg{
				WithReaderArgs: tt.fields.readerArg,
				WithWriterArgs: tt.fields.writerArg,
			}
			assert.Equalf(t, tt.want, re.isReaderArgs(), "isReaderArgs()")
		})
	}
}

func Test_readerWriterArg_isWriterArgs(t *testing.T) {
	type fields struct {
		readerArg WithReaderArgs
		writerArg WithWriterArgs
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "it is",
			fields: fields{
				writerArg: writerArg{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := readerWriterArg{
				WithReaderArgs: tt.fields.readerArg,
				WithWriterArgs: tt.fields.writerArg,
			}
			assert.Equalf(t, tt.want, re.isWriterArgs(), "isWriterArgs()")
		})
	}
}
