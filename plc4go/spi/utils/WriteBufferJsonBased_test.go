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
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/big"
	"strings"
	"testing"
)

func TestNewJsonWriteBuffer(t *testing.T) {
	tests := []struct {
		name string
		want WriteBufferJsonBased
	}{
		{
			name: "create it",
			want: func() WriteBufferJsonBased {
				var jsonString strings.Builder
				encoder := json.NewEncoder(&jsonString)
				encoder.SetIndent("", "  ")
				return &jsonWriteBuffer{
					jsonString:   &jsonString,
					Encoder:      encoder,
					doRenderAttr: true,
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewJsonWriteBuffer(), "NewJsonWriteBuffer()")
		})
	}
}

func TestNewJsonWriteBufferWithOptions(t *testing.T) {
	type args struct {
		renderAttr bool
	}
	tests := []struct {
		name string
		args args
		want WriteBufferJsonBased
	}{
		{
			name: "create it",
			want: func() WriteBufferJsonBased {
				var jsonString strings.Builder
				encoder := json.NewEncoder(&jsonString)
				encoder.SetIndent("", "  ")
				return &jsonWriteBuffer{
					jsonString:   &jsonString,
					Encoder:      encoder,
					doRenderAttr: false,
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewJsonWriteBufferWithOptions(tt.args.renderAttr), "NewJsonWriteBufferWithOptions(%v)", tt.args.renderAttr)
		})
	}
}

func Test_jsonWriteBuffer_GetJsonString(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "get it (no content)",
			wantErr: assert.Error,
		},
		// TODO: at other tests
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			got, err := j.GetJsonString()
			if !tt.wantErr(t, err, fmt.Sprintf("GetJsonString()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetJsonString()")
		})
	}
}

func Test_jsonWriteBuffer_GetPos(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
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
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			assert.Equalf(t, tt.want, j.GetPos(), "GetPos()")
		})
	}
}

func Test_jsonWriteBuffer_PopContext(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		logicalName string
		in1         []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "pop it",
			wantErr: assert.Error,
		},
		// TODO: write other tests...
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, j.PopContext(tt.args.logicalName, tt.args.in1...), fmt.Sprintf("PopContext(%v, %v)", tt.args.logicalName, tt.args.in1))
		})
	}
}

func Test_jsonWriteBuffer_PushContext(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		logicalName string
		writerArgs  []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "push it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, j.PushContext(tt.args.logicalName, tt.args.writerArgs...), fmt.Sprintf("PushContext(%v, %v)", tt.args.logicalName, tt.args.writerArgs))
		})
	}
}

func Test_jsonWriteBuffer_WriteBigFloat(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint8
		value       *big.Float
		writerArgs  []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "write it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, j.WriteBigFloat(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteBigFloat(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_jsonWriteBuffer_WriteBigInt(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint8
		value       *big.Int
		writerArgs  []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "write it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, j.WriteBigInt(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteBigInt(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_jsonWriteBuffer_WriteBit(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		logicalName string
		value       bool
		writerArgs  []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "write it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, j.WriteBit(tt.args.logicalName, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteBit(%v, %v, %v)", tt.args.logicalName, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_jsonWriteBuffer_WriteByte(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		logicalName string
		value       byte
		writerArgs  []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "write it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, j.WriteByte(tt.args.logicalName, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteByte(%v, %v, %v)", tt.args.logicalName, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_jsonWriteBuffer_WriteByteArray(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		logicalName string
		data        []byte
		writerArgs  []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "write it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, j.WriteByteArray(tt.args.logicalName, tt.args.data, tt.args.writerArgs...), fmt.Sprintf("WriteByteArray(%v, %v, %v)", tt.args.logicalName, tt.args.data, tt.args.writerArgs))
		})
	}
}

func Test_jsonWriteBuffer_WriteFloat32(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint8
		value       float32
		writerArgs  []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "write it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, j.WriteFloat32(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteFloat32(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_jsonWriteBuffer_WriteFloat64(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint8
		value       float64
		writerArgs  []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "write it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, j.WriteFloat64(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteFloat64(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_jsonWriteBuffer_WriteInt16(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint8
		value       int16
		writerArgs  []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "write it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, j.WriteInt16(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteInt16(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_jsonWriteBuffer_WriteInt32(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint8
		value       int32
		writerArgs  []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "write it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, j.WriteInt32(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteInt32(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_jsonWriteBuffer_WriteInt64(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint8
		value       int64
		writerArgs  []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "write it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, j.WriteInt64(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteInt64(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_jsonWriteBuffer_WriteInt8(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint8
		value       int8
		writerArgs  []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "write it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, j.WriteInt8(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteInt8(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_jsonWriteBuffer_WriteSerializable(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		ctx          context.Context
		serializable Serializable
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "write it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, j.WriteSerializable(tt.args.ctx, tt.args.serializable), fmt.Sprintf("WriteSerializable(%v, %v)", tt.args.ctx, tt.args.serializable))
		})
	}
}

func Test_jsonWriteBuffer_WriteString(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint32
		encoding    string
		value       string
		writerArgs  []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "write it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, j.WriteString(tt.args.logicalName, tt.args.bitLength, tt.args.encoding, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteString(%v, %v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.encoding, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_jsonWriteBuffer_WriteUint16(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint8
		value       uint16
		writerArgs  []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "write it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, j.WriteUint16(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteUint16(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_jsonWriteBuffer_WriteUint32(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint8
		value       uint32
		writerArgs  []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "write it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, j.WriteUint32(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteUint32(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_jsonWriteBuffer_WriteUint64(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint8
		value       uint64
		writerArgs  []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "write it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, j.WriteUint64(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteUint64(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_jsonWriteBuffer_WriteUint8(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint8
		value       uint8
		writerArgs  []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "write it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, j.WriteUint8(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteUint8(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_jsonWriteBuffer_WriteVirtual(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		ctx         context.Context
		logicalName string
		value       any
		writerArgs  []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "write it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, j.WriteVirtual(tt.args.ctx, tt.args.logicalName, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteVirtual(%v, %v, %v, %v)", tt.args.ctx, tt.args.logicalName, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_jsonWriteBuffer_encodeNode(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		logicalName string
		value       any
		attr        map[string]any
		in3         []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "encode it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, j.encodeNode(tt.args.logicalName, tt.args.value, tt.args.attr, tt.args.in3...), fmt.Sprintf("encodeNode(%v, %v, %v, %v)", tt.args.logicalName, tt.args.value, tt.args.attr, tt.args.in3))
		})
	}
}

func Test_jsonWriteBuffer_generateAttr(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		logicalName string
		dataType    string
		bitLength   uint
		writerArgs  []WithWriterArgs
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]any
	}{
		{
			name: "generate it (disabled)",
			want: map[string]any{},
		},
		{
			name: "generate it",
			fields: fields{
				doRenderAttr: true,
			},
			args: args{
				writerArgs: []WithWriterArgs{
					WithAdditionalStringRepresentation("nope"),
				},
			},
			want: map[string]any{
				"value__plc4x_bitLength":            uint(0x0),
				"value__plc4x_dataType":             "",
				"value__plc4x_stringRepresentation": "nope",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			assert.Equalf(t, tt.want, j.generateAttr(tt.args.logicalName, tt.args.dataType, tt.args.bitLength, tt.args.writerArgs...), "generateAttr(%v, %v, %v, %v)", tt.args.logicalName, tt.args.dataType, tt.args.bitLength, tt.args.writerArgs)
		})
	}
}

func Test_jsonWriteBuffer_move(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		Stack         Stack
		Encoder       *json.Encoder
		jsonString    *strings.Builder
		rootNode      any
		doRenderAttr  bool
		pos           uint
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
			j := &jsonWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				Encoder:       tt.fields.Encoder,
				jsonString:    tt.fields.jsonString,
				rootNode:      tt.fields.rootNode,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			j.move(tt.args.bits)
		})
	}
}
