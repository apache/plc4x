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
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"math/big"
	"strings"
	"testing"
)

func TestNewJsonReadBuffer(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name string
		args args
		want ReadBuffer
	}{
		{
			name: "create it",
			args: args{
				reader: strings.NewReader("abc"),
			},
			want: func() ReadBuffer {
				decoder := json.NewDecoder(strings.NewReader("abc"))
				var rootElement map[string]any
				err := decoder.Decode(&rootElement)
				return &jsonReadBuffer{
					rootElement:    rootElement,
					pos:            1,
					doValidateAttr: true,
					err:            err,
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewJsonReadBuffer(tt.args.reader), "NewJsonReadBuffer(%v)", tt.args.reader)
		})
	}
}

func TestNewStrictJsonReadBuffer(t *testing.T) {
	type args struct {
		reader       io.Reader
		validateAttr bool
	}
	tests := []struct {
		name string
		args args
		want ReadBuffer
	}{
		{
			name: "create it",
			args: args{
				reader: strings.NewReader("abc"),
			},
			want: func() ReadBuffer {
				decoder := json.NewDecoder(strings.NewReader("abc"))
				var rootElement map[string]any
				err := decoder.Decode(&rootElement)
				return &jsonReadBuffer{
					rootElement:    rootElement,
					pos:            1,
					doValidateAttr: false,
					err:            err,
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewStrictJsonReadBuffer(tt.args.reader, tt.args.validateAttr), "NewStrictJsonReadBuffer(%v, %v)", tt.args.reader, tt.args.validateAttr)
		})
	}
}

func Test_jsonReadBuffer_CloseContext(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Stack          Stack
		rootElement    map[string]any
		pos            uint
		doValidateAttr bool
		err            error
	}
	type args struct {
		logicalName string
		readerArgs  []WithReaderArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "close it (not found)",
			wantErr: assert.Error,
		},
		// TODO: add other tests
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Stack:          tt.fields.Stack,
				rootElement:    tt.fields.rootElement,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				err:            tt.fields.err,
			}
			tt.wantErr(t, j.CloseContext(tt.args.logicalName, tt.args.readerArgs...), fmt.Sprintf("CloseContext(%v, %v)", tt.args.logicalName, tt.args.readerArgs))
		})
	}
}

func Test_jsonReadBuffer_GetPos(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Stack          Stack
		rootElement    map[string]any
		pos            uint
		doValidateAttr bool
		err            error
	}
	tests := []struct {
		name   string
		fields fields
		want   uint16
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Stack:          tt.fields.Stack,
				rootElement:    tt.fields.rootElement,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				err:            tt.fields.err,
			}
			assert.Equalf(t, tt.want, j.GetPos(), "GetPos()")
		})
	}
}

func Test_jsonReadBuffer_HasMore(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Stack          Stack
		rootElement    map[string]any
		pos            uint
		doValidateAttr bool
		err            error
	}
	type args struct {
		bitLength uint8
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Stack:          tt.fields.Stack,
				rootElement:    tt.fields.rootElement,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				err:            tt.fields.err,
			}
			assert.Equalf(t, tt.want, j.HasMore(tt.args.bitLength), "HasMore(%v)", tt.args.bitLength)
		})
	}
}

func Test_jsonReadBuffer_PullContext(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Stack          Stack
		rootElement    map[string]any
		pos            uint
		doValidateAttr bool
		err            error
	}
	type args struct {
		logicalName string
		readerArgs  []WithReaderArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "context not found",
			wantErr: assert.Error,
		},
		// TODO: add other test cases
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Stack:          tt.fields.Stack,
				rootElement:    tt.fields.rootElement,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				err:            tt.fields.err,
			}
			tt.wantErr(t, j.PullContext(tt.args.logicalName, tt.args.readerArgs...), fmt.Sprintf("PullContext(%v, %v)", tt.args.logicalName, tt.args.readerArgs))
		})
	}
}

func Test_jsonReadBuffer_ReadBigFloat(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Stack          Stack
		rootElement    map[string]any
		pos            uint
		doValidateAttr bool
		err            error
	}
	type args struct {
		logicalName string
		bitLength   uint8
		readerArgs  []WithReaderArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *big.Float
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "read it",
			fields: fields{
				Stack: func() Stack {
					s := Stack{}
					s.Push(map[string]any{
						"value": "1337",
					})
					return s
				}(),
			},
			want:    big.NewFloat(1337),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Stack:          tt.fields.Stack,
				rootElement:    tt.fields.rootElement,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				err:            tt.fields.err,
			}
			got, err := j.ReadBigFloat(tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadBigFloat(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadBigFloat(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)
		})
	}
}

func Test_jsonReadBuffer_ReadBigInt(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Stack          Stack
		rootElement    map[string]any
		pos            uint
		doValidateAttr bool
		err            error
	}
	type args struct {
		logicalName string
		bitLength   uint64
		readerArgs  []WithReaderArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *big.Int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "read it",
			fields: fields{
				Stack: func() Stack {
					s := Stack{}
					s.Push(map[string]any{
						"value": "1337",
					})
					return s
				}(),
			},
			want:    big.NewInt(1337),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Stack:          tt.fields.Stack,
				rootElement:    tt.fields.rootElement,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				err:            tt.fields.err,
			}
			got, err := j.ReadBigInt(tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadBigInt(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadBigInt(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)
		})
	}
}

func Test_jsonReadBuffer_ReadBit(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Stack          Stack
		rootElement    map[string]any
		pos            uint
		doValidateAttr bool
		err            error
	}
	type args struct {
		logicalName string
		readerArgs  []WithReaderArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "read it",
			fields: fields{
				Stack: func() Stack {
					s := Stack{}
					s.Push(map[string]any{
						"value": true,
					})
					return s
				}(),
			},
			want:    true,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Stack:          tt.fields.Stack,
				rootElement:    tt.fields.rootElement,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				err:            tt.fields.err,
			}
			got, err := j.ReadBit(tt.args.logicalName, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadBit(%v, %v)", tt.args.logicalName, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadBit(%v, %v)", tt.args.logicalName, tt.args.readerArgs)
		})
	}
}

func Test_jsonReadBuffer_ReadByte(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Stack          Stack
		rootElement    map[string]any
		pos            uint
		doValidateAttr bool
		err            error
	}
	type args struct {
		logicalName string
		readerArgs  []WithReaderArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "read it",
			fields: fields{
				Stack: func() Stack {
					s := Stack{}
					s.Push(map[string]any{
						"value": "0x25",
					})
					return s
				}(),
			},
			want:    0x25,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Stack:          tt.fields.Stack,
				rootElement:    tt.fields.rootElement,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				err:            tt.fields.err,
			}
			got, err := j.ReadByte(tt.args.logicalName, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadByte(%v, %v)", tt.args.logicalName, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadByte(%v, %v)", tt.args.logicalName, tt.args.readerArgs)
		})
	}
}

func Test_jsonReadBuffer_ReadByteArray(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Stack          Stack
		rootElement    map[string]any
		pos            uint
		doValidateAttr bool
		err            error
	}
	type args struct {
		logicalName   string
		numberOfBytes int
		readerArgs    []WithReaderArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "read it",
			fields: fields{
				Stack: func() Stack {
					s := Stack{}
					s.Push(map[string]any{
						"value": "0xAFFE",
					})
					return s
				}(),
			},
			want:    []byte{0xAF, 0xFE},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Stack:          tt.fields.Stack,
				rootElement:    tt.fields.rootElement,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				err:            tt.fields.err,
			}
			got, err := j.ReadByteArray(tt.args.logicalName, tt.args.numberOfBytes, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadByteArray(%v, %v, %v)", tt.args.logicalName, tt.args.numberOfBytes, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadByteArray(%v, %v, %v)", tt.args.logicalName, tt.args.numberOfBytes, tt.args.readerArgs)
		})
	}
}

func Test_jsonReadBuffer_ReadFloat32(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Stack          Stack
		rootElement    map[string]any
		pos            uint
		doValidateAttr bool
		err            error
	}
	type args struct {
		logicalName string
		bitLength   uint8
		readerArgs  []WithReaderArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float32
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "read it",
			fields: fields{
				Stack: func() Stack {
					s := Stack{}
					s.Push(map[string]any{
						"value": 25.3,
					})
					return s
				}(),
			},
			want:    25.3,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Stack:          tt.fields.Stack,
				rootElement:    tt.fields.rootElement,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				err:            tt.fields.err,
			}
			got, err := j.ReadFloat32(tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadFloat32(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadFloat32(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)
		})
	}
}

func Test_jsonReadBuffer_ReadFloat64(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Stack          Stack
		rootElement    map[string]any
		pos            uint
		doValidateAttr bool
		err            error
	}
	type args struct {
		logicalName string
		bitLength   uint8
		readerArgs  []WithReaderArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "read it",
			fields: fields{
				Stack: func() Stack {
					s := Stack{}
					s.Push(map[string]any{
						"value": 25.3,
					})
					return s
				}(),
			},
			want:    25.3,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Stack:          tt.fields.Stack,
				rootElement:    tt.fields.rootElement,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				err:            tt.fields.err,
			}
			got, err := j.ReadFloat64(tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadFloat64(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadFloat64(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)
		})
	}
}

func Test_jsonReadBuffer_ReadInt16(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Stack          Stack
		rootElement    map[string]any
		pos            uint
		doValidateAttr bool
		err            error
	}
	type args struct {
		logicalName string
		bitLength   uint8
		readerArgs  []WithReaderArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int16
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "read it",
			fields: fields{
				Stack: func() Stack {
					s := Stack{}
					s.Push(map[string]any{
						"value": 25.0,
					})
					return s
				}(),
			},
			want:    25,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Stack:          tt.fields.Stack,
				rootElement:    tt.fields.rootElement,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				err:            tt.fields.err,
			}
			got, err := j.ReadInt16(tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadInt16(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadInt16(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)
		})
	}
}

func Test_jsonReadBuffer_ReadInt32(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Stack          Stack
		rootElement    map[string]any
		pos            uint
		doValidateAttr bool
		err            error
	}
	type args struct {
		logicalName string
		bitLength   uint8
		readerArgs  []WithReaderArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int32
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "read it",
			fields: fields{
				Stack: func() Stack {
					s := Stack{}
					s.Push(map[string]any{
						"value": 25.0,
					})
					return s
				}(),
			},
			want:    25,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Stack:          tt.fields.Stack,
				rootElement:    tt.fields.rootElement,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				err:            tt.fields.err,
			}
			got, err := j.ReadInt32(tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadInt32(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadInt32(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)
		})
	}
}

func Test_jsonReadBuffer_ReadInt64(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Stack          Stack
		rootElement    map[string]any
		pos            uint
		doValidateAttr bool
		err            error
	}
	type args struct {
		logicalName string
		bitLength   uint8
		readerArgs  []WithReaderArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "read it",
			fields: fields{
				Stack: func() Stack {
					s := Stack{}
					s.Push(map[string]any{
						"value": 25.0,
					})
					return s
				}(),
			},
			want:    25,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Stack:          tt.fields.Stack,
				rootElement:    tt.fields.rootElement,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				err:            tt.fields.err,
			}
			got, err := j.ReadInt64(tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadInt64(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadInt64(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)
		})
	}
}

func Test_jsonReadBuffer_ReadInt8(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Stack          Stack
		rootElement    map[string]any
		pos            uint
		doValidateAttr bool
		err            error
	}
	type args struct {
		logicalName string
		bitLength   uint8
		readerArgs  []WithReaderArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int8
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "read it",
			fields: fields{
				Stack: func() Stack {
					s := Stack{}
					s.Push(map[string]any{
						"value": 25.0,
					})
					return s
				}(),
			},
			want:    25,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Stack:          tt.fields.Stack,
				rootElement:    tt.fields.rootElement,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				err:            tt.fields.err,
			}
			got, err := j.ReadInt8(tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadInt8(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadInt8(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)
		})
	}
}

func Test_jsonReadBuffer_ReadString(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Stack          Stack
		rootElement    map[string]any
		pos            uint
		doValidateAttr bool
		err            error
	}
	type args struct {
		logicalName string
		bitLength   uint32
		encoding    string
		readerArgs  []WithReaderArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "read it",
			fields: fields{
				Stack: func() Stack {
					s := Stack{}
					s.Push(map[string]any{
						"value": "25.0",
					})
					return s
				}(),
			},
			want:    "25.0",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Stack:          tt.fields.Stack,
				rootElement:    tt.fields.rootElement,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				err:            tt.fields.err,
			}
			got, err := j.ReadString(tt.args.logicalName, tt.args.bitLength, tt.args.encoding, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadString(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.encoding, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadString(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.encoding, tt.args.readerArgs)
		})
	}
}

func Test_jsonReadBuffer_ReadUint16(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Stack          Stack
		rootElement    map[string]any
		pos            uint
		doValidateAttr bool
		err            error
	}
	type args struct {
		logicalName string
		bitLength   uint8
		readerArgs  []WithReaderArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint16
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "read it",
			fields: fields{
				Stack: func() Stack {
					s := Stack{}
					s.Push(map[string]any{
						"value": 25.0,
					})
					return s
				}(),
			},
			want:    25,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Stack:          tt.fields.Stack,
				rootElement:    tt.fields.rootElement,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				err:            tt.fields.err,
			}
			got, err := j.ReadUint16(tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadUint16(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadUint16(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)
		})
	}
}

func Test_jsonReadBuffer_ReadUint32(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Stack          Stack
		rootElement    map[string]any
		pos            uint
		doValidateAttr bool
		err            error
	}
	type args struct {
		logicalName string
		bitLength   uint8
		readerArgs  []WithReaderArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint32
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "read it",
			fields: fields{
				Stack: func() Stack {
					s := Stack{}
					s.Push(map[string]any{
						"value": 25.0,
					})
					return s
				}(),
			},
			want:    25,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Stack:          tt.fields.Stack,
				rootElement:    tt.fields.rootElement,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				err:            tt.fields.err,
			}
			got, err := j.ReadUint32(tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadUint32(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadUint32(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)
		})
	}
}

func Test_jsonReadBuffer_ReadUint64(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Stack          Stack
		rootElement    map[string]any
		pos            uint
		doValidateAttr bool
		err            error
	}
	type args struct {
		logicalName string
		bitLength   uint8
		readerArgs  []WithReaderArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint64
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "read it",
			fields: fields{
				Stack: func() Stack {
					s := Stack{}
					s.Push(map[string]any{
						"value": 25.0,
					})
					return s
				}(),
			},
			want:    25,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Stack:          tt.fields.Stack,
				rootElement:    tt.fields.rootElement,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				err:            tt.fields.err,
			}
			got, err := j.ReadUint64(tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadUint64(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadUint64(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)
		})
	}
}

func Test_jsonReadBuffer_ReadUint8(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Stack          Stack
		rootElement    map[string]any
		pos            uint
		doValidateAttr bool
		err            error
	}
	type args struct {
		logicalName string
		bitLength   uint8
		readerArgs  []WithReaderArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint8
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "read it",
			fields: fields{
				Stack: func() Stack {
					s := Stack{}
					s.Push(map[string]any{
						"value": 25.0,
					})
					return s
				}(),
			},
			want:    25,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Stack:          tt.fields.Stack,
				rootElement:    tt.fields.rootElement,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				err:            tt.fields.err,
			}
			got, err := j.ReadUint8(tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadUint8(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadUint8(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)
		})
	}
}

func Test_jsonReadBuffer_Reset(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Stack          Stack
		rootElement    map[string]any
		pos            uint
		doValidateAttr bool
		err            error
	}
	type args struct {
		pos uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "reset it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Stack:          tt.fields.Stack,
				rootElement:    tt.fields.rootElement,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				err:            tt.fields.err,
			}
			j.Reset(tt.args.pos)
		})
	}
}

func Test_jsonReadBuffer_getElement(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Stack          Stack
		rootElement    map[string]any
		pos            uint
		doValidateAttr bool
		err            error
	}
	type args struct {
		logicalName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   any
		want1  map[string]any
	}{
		{
			name: "get it",
			fields: fields{
				Stack: func() Stack {
					s := Stack{}
					s.Push(map[string]any{})
					return s
				}(),
			},
			want:  map[string]any{},
			want1: map[string]any{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Stack:          tt.fields.Stack,
				rootElement:    tt.fields.rootElement,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				err:            tt.fields.err,
			}
			got, got1 := j.getElement(tt.args.logicalName)
			assert.Equalf(t, tt.want, got, "getElement(%v)", tt.args.logicalName)
			assert.Equalf(t, tt.want1, got1, "getElement(%v)", tt.args.logicalName)
		})
	}
}

func Test_jsonReadBuffer_move(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Stack          Stack
		rootElement    map[string]any
		pos            uint
		doValidateAttr bool
		err            error
	}
	type args struct {
		bits uint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "move it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Stack:          tt.fields.Stack,
				rootElement:    tt.fields.rootElement,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				err:            tt.fields.err,
			}
			j.move(tt.args.bits)
		})
	}
}

func Test_jsonReadBuffer_validateAttr(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Stack          Stack
		rootElement    map[string]any
		pos            uint
		doValidateAttr bool
		err            error
	}
	type args struct {
		logicalName string
		element     map[string]any
		dataType    string
		bitLength   uint
		readerArgs  []WithReaderArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "validate it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Stack:          tt.fields.Stack,
				rootElement:    tt.fields.rootElement,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				err:            tt.fields.err,
			}
			tt.wantErr(t, j.validateAttr(tt.args.logicalName, tt.args.element, tt.args.dataType, tt.args.bitLength, tt.args.readerArgs...), fmt.Sprintf("validateAttr(%v, %v, %v, %v, %v)", tt.args.logicalName, tt.args.element, tt.args.dataType, tt.args.bitLength, tt.args.readerArgs))
		})
	}
}
