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
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestNewWriteBufferBoxBased(t *testing.T) {
	tests := []struct {
		name string
		want WriteBufferBoxBased
	}{
		{
			name: "create it",
			want: &boxedWriteBuffer{
				List:                list.New(),
				desiredWidth:        120,
				currentWidth:        118,
				mergeSingleBoxes:    false,
				omitEmptyBoxes:      false,
				asciiBoxWriter:      AsciiBoxWriterDefault,
				asciiBoxWriterLight: AsciiBoxWriterLight,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewWriteBufferBoxBased(), "NewWriteBufferBoxBased()")
		})
	}
}

func TestNewWriteBufferBoxBasedWithOptions(t *testing.T) {
	type args struct {
		mergeSingleBoxes bool
		omitEmptyBoxes   bool
	}
	tests := []struct {
		name string
		args args
		want WriteBufferBoxBased
	}{
		{
			name: "create it",
			args: args{
				mergeSingleBoxes: true,
				omitEmptyBoxes:   true,
			},
			want: &boxedWriteBuffer{
				List:                list.New(),
				desiredWidth:        120,
				currentWidth:        118,
				mergeSingleBoxes:    true,
				omitEmptyBoxes:      true,
				asciiBoxWriter:      AsciiBoxWriterDefault,
				asciiBoxWriterLight: AsciiBoxWriterLight,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewWriteBufferBoxBasedWithOptions(tt.args.mergeSingleBoxes, tt.args.omitEmptyBoxes), "NewWriteBufferBoxBasedWithOptions(%v, %v)", tt.args.mergeSingleBoxes, tt.args.omitEmptyBoxes)
		})
	}
}

func Test_boxedWriteBuffer_GetBox(t *testing.T) {
	type fields struct {
		BufferCommons       BufferCommons
		List                *list.List
		desiredWidth        int
		currentWidth        int
		mergeSingleBoxes    bool
		omitEmptyBoxes      bool
		asciiBoxWriter      AsciiBoxWriter
		asciiBoxWriterLight AsciiBoxWriter
		pos                 uint
	}
	tests := []struct {
		name   string
		fields fields
		want   AsciiBox
	}{
		{
			name: "get it",
			fields: fields{
				List: list.New(),
			},
			want: AsciiBox{data: "<nil>"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &boxedWriteBuffer{
				BufferCommons:       tt.fields.BufferCommons,
				List:                tt.fields.List,
				desiredWidth:        tt.fields.desiredWidth,
				currentWidth:        tt.fields.currentWidth,
				mergeSingleBoxes:    tt.fields.mergeSingleBoxes,
				omitEmptyBoxes:      tt.fields.omitEmptyBoxes,
				asciiBoxWriter:      tt.fields.asciiBoxWriter,
				asciiBoxWriterLight: tt.fields.asciiBoxWriterLight,
				pos:                 tt.fields.pos,
			}
			assert.Equalf(t, tt.want, b.GetBox(), "GetBox()")
		})
	}
}

func Test_boxedWriteBuffer_GetPos(t *testing.T) {
	type fields struct {
		BufferCommons       BufferCommons
		List                *list.List
		desiredWidth        int
		currentWidth        int
		mergeSingleBoxes    bool
		omitEmptyBoxes      bool
		asciiBoxWriter      AsciiBoxWriter
		asciiBoxWriterLight AsciiBoxWriter
		pos                 uint
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
			b := &boxedWriteBuffer{
				BufferCommons:       tt.fields.BufferCommons,
				List:                tt.fields.List,
				desiredWidth:        tt.fields.desiredWidth,
				currentWidth:        tt.fields.currentWidth,
				mergeSingleBoxes:    tt.fields.mergeSingleBoxes,
				omitEmptyBoxes:      tt.fields.omitEmptyBoxes,
				asciiBoxWriter:      tt.fields.asciiBoxWriter,
				asciiBoxWriterLight: tt.fields.asciiBoxWriterLight,
				pos:                 tt.fields.pos,
			}
			assert.Equalf(t, tt.want, b.GetPos(), "GetPos()")
		})
	}
}

func Test_boxedWriteBuffer_PopContext(t *testing.T) {
	type fields struct {
		BufferCommons       BufferCommons
		List                *list.List
		desiredWidth        int
		currentWidth        int
		mergeSingleBoxes    bool
		omitEmptyBoxes      bool
		asciiBoxWriter      AsciiBoxWriter
		asciiBoxWriterLight AsciiBoxWriter
		pos                 uint
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
			name: "pop it",
			fields: fields{
				List:           list.New(),
				asciiBoxWriter: NewAsciiBoxWriter(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &boxedWriteBuffer{
				BufferCommons:       tt.fields.BufferCommons,
				List:                tt.fields.List,
				desiredWidth:        tt.fields.desiredWidth,
				currentWidth:        tt.fields.currentWidth,
				mergeSingleBoxes:    tt.fields.mergeSingleBoxes,
				omitEmptyBoxes:      tt.fields.omitEmptyBoxes,
				asciiBoxWriter:      tt.fields.asciiBoxWriter,
				asciiBoxWriterLight: tt.fields.asciiBoxWriterLight,
				pos:                 tt.fields.pos,
			}
			tt.wantErr(t, b.PopContext(tt.args.logicalName, tt.args.in1...), fmt.Sprintf("PopContext(%v, %v)", tt.args.logicalName, tt.args.in1))
		})
	}
}

func Test_boxedWriteBuffer_PushContext(t *testing.T) {
	type fields struct {
		BufferCommons       BufferCommons
		List                *list.List
		desiredWidth        int
		currentWidth        int
		mergeSingleBoxes    bool
		omitEmptyBoxes      bool
		asciiBoxWriter      AsciiBoxWriter
		asciiBoxWriterLight AsciiBoxWriter
		pos                 uint
	}
	type args struct {
		in0 string
		in1 []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "push it",
			fields: fields{
				List:           list.New(),
				asciiBoxWriter: NewAsciiBoxWriter(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &boxedWriteBuffer{
				BufferCommons:       tt.fields.BufferCommons,
				List:                tt.fields.List,
				desiredWidth:        tt.fields.desiredWidth,
				currentWidth:        tt.fields.currentWidth,
				mergeSingleBoxes:    tt.fields.mergeSingleBoxes,
				omitEmptyBoxes:      tt.fields.omitEmptyBoxes,
				asciiBoxWriter:      tt.fields.asciiBoxWriter,
				asciiBoxWriterLight: tt.fields.asciiBoxWriterLight,
				pos:                 tt.fields.pos,
			}
			tt.wantErr(t, b.PushContext(tt.args.in0, tt.args.in1...), fmt.Sprintf("PushContext(%v, %v)", tt.args.in0, tt.args.in1))
		})
	}
}

func Test_boxedWriteBuffer_WriteBigFloat(t *testing.T) {
	type fields struct {
		BufferCommons       BufferCommons
		List                *list.List
		desiredWidth        int
		currentWidth        int
		mergeSingleBoxes    bool
		omitEmptyBoxes      bool
		asciiBoxWriter      AsciiBoxWriter
		asciiBoxWriterLight AsciiBoxWriter
		pos                 uint
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
			name: "write it",
			fields: fields{
				List:           list.New(),
				asciiBoxWriter: NewAsciiBoxWriter(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &boxedWriteBuffer{
				BufferCommons:       tt.fields.BufferCommons,
				List:                tt.fields.List,
				desiredWidth:        tt.fields.desiredWidth,
				currentWidth:        tt.fields.currentWidth,
				mergeSingleBoxes:    tt.fields.mergeSingleBoxes,
				omitEmptyBoxes:      tt.fields.omitEmptyBoxes,
				asciiBoxWriter:      tt.fields.asciiBoxWriter,
				asciiBoxWriterLight: tt.fields.asciiBoxWriterLight,
				pos:                 tt.fields.pos,
			}
			tt.wantErr(t, b.WriteBigFloat(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteBigFloat(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_boxedWriteBuffer_WriteBigInt(t *testing.T) {
	type fields struct {
		BufferCommons       BufferCommons
		List                *list.List
		desiredWidth        int
		currentWidth        int
		mergeSingleBoxes    bool
		omitEmptyBoxes      bool
		asciiBoxWriter      AsciiBoxWriter
		asciiBoxWriterLight AsciiBoxWriter
		pos                 uint
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
			name: "write it",
			fields: fields{
				List:           list.New(),
				asciiBoxWriter: NewAsciiBoxWriter(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &boxedWriteBuffer{
				BufferCommons:       tt.fields.BufferCommons,
				List:                tt.fields.List,
				desiredWidth:        tt.fields.desiredWidth,
				currentWidth:        tt.fields.currentWidth,
				mergeSingleBoxes:    tt.fields.mergeSingleBoxes,
				omitEmptyBoxes:      tt.fields.omitEmptyBoxes,
				asciiBoxWriter:      tt.fields.asciiBoxWriter,
				asciiBoxWriterLight: tt.fields.asciiBoxWriterLight,
				pos:                 tt.fields.pos,
			}
			tt.wantErr(t, b.WriteBigInt(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteBigInt(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_boxedWriteBuffer_WriteBit(t *testing.T) {
	type fields struct {
		BufferCommons       BufferCommons
		List                *list.List
		desiredWidth        int
		currentWidth        int
		mergeSingleBoxes    bool
		omitEmptyBoxes      bool
		asciiBoxWriter      AsciiBoxWriter
		asciiBoxWriterLight AsciiBoxWriter
		pos                 uint
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
			name: "write it",
			fields: fields{
				List:           list.New(),
				asciiBoxWriter: NewAsciiBoxWriter(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &boxedWriteBuffer{
				BufferCommons:       tt.fields.BufferCommons,
				List:                tt.fields.List,
				desiredWidth:        tt.fields.desiredWidth,
				currentWidth:        tt.fields.currentWidth,
				mergeSingleBoxes:    tt.fields.mergeSingleBoxes,
				omitEmptyBoxes:      tt.fields.omitEmptyBoxes,
				asciiBoxWriter:      tt.fields.asciiBoxWriter,
				asciiBoxWriterLight: tt.fields.asciiBoxWriterLight,
				pos:                 tt.fields.pos,
			}
			tt.wantErr(t, b.WriteBit(tt.args.logicalName, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteBit(%v, %v, %v)", tt.args.logicalName, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_boxedWriteBuffer_WriteByte(t *testing.T) {
	type fields struct {
		BufferCommons       BufferCommons
		List                *list.List
		desiredWidth        int
		currentWidth        int
		mergeSingleBoxes    bool
		omitEmptyBoxes      bool
		asciiBoxWriter      AsciiBoxWriter
		asciiBoxWriterLight AsciiBoxWriter
		pos                 uint
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
			name: "write it",
			fields: fields{
				List:           list.New(),
				asciiBoxWriter: NewAsciiBoxWriter(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &boxedWriteBuffer{
				BufferCommons:       tt.fields.BufferCommons,
				List:                tt.fields.List,
				desiredWidth:        tt.fields.desiredWidth,
				currentWidth:        tt.fields.currentWidth,
				mergeSingleBoxes:    tt.fields.mergeSingleBoxes,
				omitEmptyBoxes:      tt.fields.omitEmptyBoxes,
				asciiBoxWriter:      tt.fields.asciiBoxWriter,
				asciiBoxWriterLight: tt.fields.asciiBoxWriterLight,
				pos:                 tt.fields.pos,
			}
			tt.wantErr(t, b.WriteByte(tt.args.logicalName, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteByte(%v, %v, %v)", tt.args.logicalName, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_boxedWriteBuffer_WriteByteArray(t *testing.T) {
	type fields struct {
		BufferCommons       BufferCommons
		List                *list.List
		desiredWidth        int
		currentWidth        int
		mergeSingleBoxes    bool
		omitEmptyBoxes      bool
		asciiBoxWriter      AsciiBoxWriter
		asciiBoxWriterLight AsciiBoxWriter
		pos                 uint
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
			name: "write it",
			fields: fields{
				List:           list.New(),
				asciiBoxWriter: NewAsciiBoxWriter(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &boxedWriteBuffer{
				BufferCommons:       tt.fields.BufferCommons,
				List:                tt.fields.List,
				desiredWidth:        tt.fields.desiredWidth,
				currentWidth:        tt.fields.currentWidth,
				mergeSingleBoxes:    tt.fields.mergeSingleBoxes,
				omitEmptyBoxes:      tt.fields.omitEmptyBoxes,
				asciiBoxWriter:      tt.fields.asciiBoxWriter,
				asciiBoxWriterLight: tt.fields.asciiBoxWriterLight,
				pos:                 tt.fields.pos,
			}
			tt.wantErr(t, b.WriteByteArray(tt.args.logicalName, tt.args.data, tt.args.writerArgs...), fmt.Sprintf("WriteByteArray(%v, %v, %v)", tt.args.logicalName, tt.args.data, tt.args.writerArgs))
		})
	}
}

func Test_boxedWriteBuffer_WriteFloat32(t *testing.T) {
	type fields struct {
		BufferCommons       BufferCommons
		List                *list.List
		desiredWidth        int
		currentWidth        int
		mergeSingleBoxes    bool
		omitEmptyBoxes      bool
		asciiBoxWriter      AsciiBoxWriter
		asciiBoxWriterLight AsciiBoxWriter
		pos                 uint
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
			name: "write it",
			fields: fields{
				List:           list.New(),
				asciiBoxWriter: NewAsciiBoxWriter(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &boxedWriteBuffer{
				BufferCommons:       tt.fields.BufferCommons,
				List:                tt.fields.List,
				desiredWidth:        tt.fields.desiredWidth,
				currentWidth:        tt.fields.currentWidth,
				mergeSingleBoxes:    tt.fields.mergeSingleBoxes,
				omitEmptyBoxes:      tt.fields.omitEmptyBoxes,
				asciiBoxWriter:      tt.fields.asciiBoxWriter,
				asciiBoxWriterLight: tt.fields.asciiBoxWriterLight,
				pos:                 tt.fields.pos,
			}
			tt.wantErr(t, b.WriteFloat32(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteFloat32(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_boxedWriteBuffer_WriteFloat64(t *testing.T) {
	type fields struct {
		BufferCommons       BufferCommons
		List                *list.List
		desiredWidth        int
		currentWidth        int
		mergeSingleBoxes    bool
		omitEmptyBoxes      bool
		asciiBoxWriter      AsciiBoxWriter
		asciiBoxWriterLight AsciiBoxWriter
		pos                 uint
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
			name: "write it",
			fields: fields{
				List:           list.New(),
				asciiBoxWriter: NewAsciiBoxWriter(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &boxedWriteBuffer{
				BufferCommons:       tt.fields.BufferCommons,
				List:                tt.fields.List,
				desiredWidth:        tt.fields.desiredWidth,
				currentWidth:        tt.fields.currentWidth,
				mergeSingleBoxes:    tt.fields.mergeSingleBoxes,
				omitEmptyBoxes:      tt.fields.omitEmptyBoxes,
				asciiBoxWriter:      tt.fields.asciiBoxWriter,
				asciiBoxWriterLight: tt.fields.asciiBoxWriterLight,
				pos:                 tt.fields.pos,
			}
			tt.wantErr(t, b.WriteFloat64(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteFloat64(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_boxedWriteBuffer_WriteInt16(t *testing.T) {
	type fields struct {
		BufferCommons       BufferCommons
		List                *list.List
		desiredWidth        int
		currentWidth        int
		mergeSingleBoxes    bool
		omitEmptyBoxes      bool
		asciiBoxWriter      AsciiBoxWriter
		asciiBoxWriterLight AsciiBoxWriter
		pos                 uint
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
			name: "write it",
			fields: fields{
				List:           list.New(),
				asciiBoxWriter: NewAsciiBoxWriter(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &boxedWriteBuffer{
				BufferCommons:       tt.fields.BufferCommons,
				List:                tt.fields.List,
				desiredWidth:        tt.fields.desiredWidth,
				currentWidth:        tt.fields.currentWidth,
				mergeSingleBoxes:    tt.fields.mergeSingleBoxes,
				omitEmptyBoxes:      tt.fields.omitEmptyBoxes,
				asciiBoxWriter:      tt.fields.asciiBoxWriter,
				asciiBoxWriterLight: tt.fields.asciiBoxWriterLight,
				pos:                 tt.fields.pos,
			}
			tt.wantErr(t, b.WriteInt16(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteInt16(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_boxedWriteBuffer_WriteInt32(t *testing.T) {
	type fields struct {
		BufferCommons       BufferCommons
		List                *list.List
		desiredWidth        int
		currentWidth        int
		mergeSingleBoxes    bool
		omitEmptyBoxes      bool
		asciiBoxWriter      AsciiBoxWriter
		asciiBoxWriterLight AsciiBoxWriter
		pos                 uint
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
			name: "write it",
			fields: fields{
				List:           list.New(),
				asciiBoxWriter: NewAsciiBoxWriter(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &boxedWriteBuffer{
				BufferCommons:       tt.fields.BufferCommons,
				List:                tt.fields.List,
				desiredWidth:        tt.fields.desiredWidth,
				currentWidth:        tt.fields.currentWidth,
				mergeSingleBoxes:    tt.fields.mergeSingleBoxes,
				omitEmptyBoxes:      tt.fields.omitEmptyBoxes,
				asciiBoxWriter:      tt.fields.asciiBoxWriter,
				asciiBoxWriterLight: tt.fields.asciiBoxWriterLight,
				pos:                 tt.fields.pos,
			}
			tt.wantErr(t, b.WriteInt32(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteInt32(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_boxedWriteBuffer_WriteInt64(t *testing.T) {
	type fields struct {
		BufferCommons       BufferCommons
		List                *list.List
		desiredWidth        int
		currentWidth        int
		mergeSingleBoxes    bool
		omitEmptyBoxes      bool
		asciiBoxWriter      AsciiBoxWriter
		asciiBoxWriterLight AsciiBoxWriter
		pos                 uint
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
			name: "write it",
			fields: fields{
				List:           list.New(),
				asciiBoxWriter: NewAsciiBoxWriter(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &boxedWriteBuffer{
				BufferCommons:       tt.fields.BufferCommons,
				List:                tt.fields.List,
				desiredWidth:        tt.fields.desiredWidth,
				currentWidth:        tt.fields.currentWidth,
				mergeSingleBoxes:    tt.fields.mergeSingleBoxes,
				omitEmptyBoxes:      tt.fields.omitEmptyBoxes,
				asciiBoxWriter:      tt.fields.asciiBoxWriter,
				asciiBoxWriterLight: tt.fields.asciiBoxWriterLight,
				pos:                 tt.fields.pos,
			}
			tt.wantErr(t, b.WriteInt64(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteInt64(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_boxedWriteBuffer_WriteInt8(t *testing.T) {
	type fields struct {
		BufferCommons       BufferCommons
		List                *list.List
		desiredWidth        int
		currentWidth        int
		mergeSingleBoxes    bool
		omitEmptyBoxes      bool
		asciiBoxWriter      AsciiBoxWriter
		asciiBoxWriterLight AsciiBoxWriter
		pos                 uint
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
			name: "write it",
			fields: fields{
				List:           list.New(),
				asciiBoxWriter: NewAsciiBoxWriter(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &boxedWriteBuffer{
				BufferCommons:       tt.fields.BufferCommons,
				List:                tt.fields.List,
				desiredWidth:        tt.fields.desiredWidth,
				currentWidth:        tt.fields.currentWidth,
				mergeSingleBoxes:    tt.fields.mergeSingleBoxes,
				omitEmptyBoxes:      tt.fields.omitEmptyBoxes,
				asciiBoxWriter:      tt.fields.asciiBoxWriter,
				asciiBoxWriterLight: tt.fields.asciiBoxWriterLight,
				pos:                 tt.fields.pos,
			}
			tt.wantErr(t, b.WriteInt8(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteInt8(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_boxedWriteBuffer_WriteSerializable(t *testing.T) {
	type fields struct {
		BufferCommons       BufferCommons
		List                *list.List
		desiredWidth        int
		currentWidth        int
		mergeSingleBoxes    bool
		omitEmptyBoxes      bool
		asciiBoxWriter      AsciiBoxWriter
		asciiBoxWriterLight AsciiBoxWriter
		pos                 uint
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
			name: "write it",
			fields: fields{
				List:           list.New(),
				asciiBoxWriter: NewAsciiBoxWriter(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &boxedWriteBuffer{
				BufferCommons:       tt.fields.BufferCommons,
				List:                tt.fields.List,
				desiredWidth:        tt.fields.desiredWidth,
				currentWidth:        tt.fields.currentWidth,
				mergeSingleBoxes:    tt.fields.mergeSingleBoxes,
				omitEmptyBoxes:      tt.fields.omitEmptyBoxes,
				asciiBoxWriter:      tt.fields.asciiBoxWriter,
				asciiBoxWriterLight: tt.fields.asciiBoxWriterLight,
				pos:                 tt.fields.pos,
			}
			tt.wantErr(t, b.WriteSerializable(tt.args.ctx, tt.args.serializable), fmt.Sprintf("WriteSerializable(%v, %v)", tt.args.ctx, tt.args.serializable))
		})
	}
}

func Test_boxedWriteBuffer_WriteString(t *testing.T) {
	type fields struct {
		BufferCommons       BufferCommons
		List                *list.List
		desiredWidth        int
		currentWidth        int
		mergeSingleBoxes    bool
		omitEmptyBoxes      bool
		asciiBoxWriter      AsciiBoxWriter
		asciiBoxWriterLight AsciiBoxWriter
		pos                 uint
	}
	type args struct {
		logicalName string
		bitLength   uint32
		in2         string
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
			name: "write it",
			fields: fields{
				List:           list.New(),
				asciiBoxWriter: NewAsciiBoxWriter(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &boxedWriteBuffer{
				BufferCommons:       tt.fields.BufferCommons,
				List:                tt.fields.List,
				desiredWidth:        tt.fields.desiredWidth,
				currentWidth:        tt.fields.currentWidth,
				mergeSingleBoxes:    tt.fields.mergeSingleBoxes,
				omitEmptyBoxes:      tt.fields.omitEmptyBoxes,
				asciiBoxWriter:      tt.fields.asciiBoxWriter,
				asciiBoxWriterLight: tt.fields.asciiBoxWriterLight,
				pos:                 tt.fields.pos,
			}
			tt.wantErr(t, b.WriteString(tt.args.logicalName, tt.args.bitLength, tt.args.in2, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteString(%v, %v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.in2, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_boxedWriteBuffer_WriteUint16(t *testing.T) {
	type fields struct {
		BufferCommons       BufferCommons
		List                *list.List
		desiredWidth        int
		currentWidth        int
		mergeSingleBoxes    bool
		omitEmptyBoxes      bool
		asciiBoxWriter      AsciiBoxWriter
		asciiBoxWriterLight AsciiBoxWriter
		pos                 uint
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
			name: "write it",
			fields: fields{
				List:           list.New(),
				asciiBoxWriter: NewAsciiBoxWriter(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &boxedWriteBuffer{
				BufferCommons:       tt.fields.BufferCommons,
				List:                tt.fields.List,
				desiredWidth:        tt.fields.desiredWidth,
				currentWidth:        tt.fields.currentWidth,
				mergeSingleBoxes:    tt.fields.mergeSingleBoxes,
				omitEmptyBoxes:      tt.fields.omitEmptyBoxes,
				asciiBoxWriter:      tt.fields.asciiBoxWriter,
				asciiBoxWriterLight: tt.fields.asciiBoxWriterLight,
				pos:                 tt.fields.pos,
			}
			tt.wantErr(t, b.WriteUint16(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteUint16(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_boxedWriteBuffer_WriteUint32(t *testing.T) {
	type fields struct {
		BufferCommons       BufferCommons
		List                *list.List
		desiredWidth        int
		currentWidth        int
		mergeSingleBoxes    bool
		omitEmptyBoxes      bool
		asciiBoxWriter      AsciiBoxWriter
		asciiBoxWriterLight AsciiBoxWriter
		pos                 uint
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
			name: "write it",
			fields: fields{
				List:           list.New(),
				asciiBoxWriter: NewAsciiBoxWriter(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &boxedWriteBuffer{
				BufferCommons:       tt.fields.BufferCommons,
				List:                tt.fields.List,
				desiredWidth:        tt.fields.desiredWidth,
				currentWidth:        tt.fields.currentWidth,
				mergeSingleBoxes:    tt.fields.mergeSingleBoxes,
				omitEmptyBoxes:      tt.fields.omitEmptyBoxes,
				asciiBoxWriter:      tt.fields.asciiBoxWriter,
				asciiBoxWriterLight: tt.fields.asciiBoxWriterLight,
				pos:                 tt.fields.pos,
			}
			tt.wantErr(t, b.WriteUint32(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteUint32(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_boxedWriteBuffer_WriteUint64(t *testing.T) {
	type fields struct {
		BufferCommons       BufferCommons
		List                *list.List
		desiredWidth        int
		currentWidth        int
		mergeSingleBoxes    bool
		omitEmptyBoxes      bool
		asciiBoxWriter      AsciiBoxWriter
		asciiBoxWriterLight AsciiBoxWriter
		pos                 uint
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
			name: "write it",
			fields: fields{
				List:           list.New(),
				asciiBoxWriter: NewAsciiBoxWriter(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &boxedWriteBuffer{
				BufferCommons:       tt.fields.BufferCommons,
				List:                tt.fields.List,
				desiredWidth:        tt.fields.desiredWidth,
				currentWidth:        tt.fields.currentWidth,
				mergeSingleBoxes:    tt.fields.mergeSingleBoxes,
				omitEmptyBoxes:      tt.fields.omitEmptyBoxes,
				asciiBoxWriter:      tt.fields.asciiBoxWriter,
				asciiBoxWriterLight: tt.fields.asciiBoxWriterLight,
				pos:                 tt.fields.pos,
			}
			tt.wantErr(t, b.WriteUint64(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteUint64(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_boxedWriteBuffer_WriteUint8(t *testing.T) {
	type fields struct {
		BufferCommons       BufferCommons
		List                *list.List
		desiredWidth        int
		currentWidth        int
		mergeSingleBoxes    bool
		omitEmptyBoxes      bool
		asciiBoxWriter      AsciiBoxWriter
		asciiBoxWriterLight AsciiBoxWriter
		pos                 uint
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
			name: "write it",
			fields: fields{
				List:           list.New(),
				asciiBoxWriter: NewAsciiBoxWriter(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &boxedWriteBuffer{
				BufferCommons:       tt.fields.BufferCommons,
				List:                tt.fields.List,
				desiredWidth:        tt.fields.desiredWidth,
				currentWidth:        tt.fields.currentWidth,
				mergeSingleBoxes:    tt.fields.mergeSingleBoxes,
				omitEmptyBoxes:      tt.fields.omitEmptyBoxes,
				asciiBoxWriter:      tt.fields.asciiBoxWriter,
				asciiBoxWriterLight: tt.fields.asciiBoxWriterLight,
				pos:                 tt.fields.pos,
			}
			tt.wantErr(t, b.WriteUint8(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteUint8(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_boxedWriteBuffer_WriteVirtual(t *testing.T) {
	type fields struct {
		BufferCommons       BufferCommons
		List                *list.List
		desiredWidth        int
		currentWidth        int
		mergeSingleBoxes    bool
		omitEmptyBoxes      bool
		asciiBoxWriter      AsciiBoxWriter
		asciiBoxWriterLight AsciiBoxWriter
		pos                 uint
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
			name: "write it",
			fields: fields{
				List:           list.New(),
				asciiBoxWriter: NewAsciiBoxWriter(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &boxedWriteBuffer{
				BufferCommons:       tt.fields.BufferCommons,
				List:                tt.fields.List,
				desiredWidth:        tt.fields.desiredWidth,
				currentWidth:        tt.fields.currentWidth,
				mergeSingleBoxes:    tt.fields.mergeSingleBoxes,
				omitEmptyBoxes:      tt.fields.omitEmptyBoxes,
				asciiBoxWriter:      tt.fields.asciiBoxWriter,
				asciiBoxWriterLight: tt.fields.asciiBoxWriterLight,
				pos:                 tt.fields.pos,
			}
			tt.wantErr(t, b.WriteVirtual(tt.args.ctx, tt.args.logicalName, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteVirtual(%v, %v, %v, %v)", tt.args.ctx, tt.args.logicalName, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_boxedWriteBuffer_extractAdditionalStringRepresentation(t *testing.T) {
	type fields struct {
		BufferCommons       BufferCommons
		List                *list.List
		desiredWidth        int
		currentWidth        int
		mergeSingleBoxes    bool
		omitEmptyBoxes      bool
		asciiBoxWriter      AsciiBoxWriter
		asciiBoxWriterLight AsciiBoxWriter
		pos                 uint
	}
	type args struct {
		readerWriterArgs []WithReaderWriterArgs
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "extract it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &boxedWriteBuffer{
				BufferCommons:       tt.fields.BufferCommons,
				List:                tt.fields.List,
				desiredWidth:        tt.fields.desiredWidth,
				currentWidth:        tt.fields.currentWidth,
				mergeSingleBoxes:    tt.fields.mergeSingleBoxes,
				omitEmptyBoxes:      tt.fields.omitEmptyBoxes,
				asciiBoxWriter:      tt.fields.asciiBoxWriter,
				asciiBoxWriterLight: tt.fields.asciiBoxWriterLight,
				pos:                 tt.fields.pos,
			}
			assert.Equalf(t, tt.want, b.extractAdditionalStringRepresentation(tt.args.readerWriterArgs...), "extractAdditionalStringRepresentation(%v)", tt.args.readerWriterArgs)
		})
	}
}

func Test_boxedWriteBuffer_move(t *testing.T) {
	type fields struct {
		BufferCommons       BufferCommons
		List                *list.List
		desiredWidth        int
		currentWidth        int
		mergeSingleBoxes    bool
		omitEmptyBoxes      bool
		asciiBoxWriter      AsciiBoxWriter
		asciiBoxWriterLight AsciiBoxWriter
		pos                 uint
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
			b := &boxedWriteBuffer{
				BufferCommons:       tt.fields.BufferCommons,
				List:                tt.fields.List,
				desiredWidth:        tt.fields.desiredWidth,
				currentWidth:        tt.fields.currentWidth,
				mergeSingleBoxes:    tt.fields.mergeSingleBoxes,
				omitEmptyBoxes:      tt.fields.omitEmptyBoxes,
				asciiBoxWriter:      tt.fields.asciiBoxWriter,
				asciiBoxWriterLight: tt.fields.asciiBoxWriterLight,
				pos:                 tt.fields.pos,
			}
			b.move(tt.args.bits)
		})
	}
}
