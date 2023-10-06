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
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/icza/bitio"
	"github.com/stretchr/testify/mock"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSizingWorksProperly(t *testing.T) {
	t.Run("default sizing", func(t *testing.T) {
		wb := NewWriteBufferByteBased()
		for i := 1; i <= 14; i++ {
			_ = wb.WriteByte("nasd", 12)
		}

		assert.Equal(t, 14, len(wb.GetBytes()))
	})
	t.Run("custom sizing", func(t *testing.T) {
		wb := NewWriteBufferByteBased(WithInitialSizeForByteBasedBuffer(23432342))
		for i := 0; i < 14; i++ {
			_ = wb.WriteByte("nasd", 12)
		}

		assert.Equal(t, 14, len(wb.GetBytes()))
	})
}

func TestNewWriteBufferByteBased(t *testing.T) {
	type args struct {
		options []WriteBufferByteBasedOptions
	}
	tests := []struct {
		name string
		args args
		want WriteBufferByteBased
	}{
		{
			name: "create it,",
			want: &byteWriteBuffer{
				data:      new(bytes.Buffer),
				writer:    bitio.NewWriter(new(bytes.Buffer)),
				byteOrder: binary.BigEndian,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewWriteBufferByteBased(tt.args.options...), "NewWriteBufferByteBased(%v)", tt.args.options)
		})
	}
}

func TestWithByteOrderForByteBasedBuffer(t *testing.T) {
	type args struct {
		byteOrder binary.ByteOrder
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "create it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withByteOrderForByteBasedBuffer := WithByteOrderForByteBasedBuffer(tt.args.byteOrder)
			assert.NotNilf(t, withByteOrderForByteBasedBuffer, "WithByteOrderForByteBasedBuffer(%v)", tt.args.byteOrder)
			withByteOrderForByteBasedBuffer(
				&byteWriteBuffer{
					data:   new(bytes.Buffer),
					writer: bitio.NewWriter(new(bytes.Buffer)),
				},
			)
		})
	}
}

func TestWithCustomBufferForByteBasedBuffer(t *testing.T) {
	type args struct {
		buffer *bytes.Buffer
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "create it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withCustomBufferForByteBasedBuffer := WithCustomBufferForByteBasedBuffer(tt.args.buffer)
			assert.NotNilf(t, withCustomBufferForByteBasedBuffer, "WithCustomBufferForByteBasedBuffer(%v)", tt.args.buffer)
			withCustomBufferForByteBasedBuffer(
				&byteWriteBuffer{
					data:   new(bytes.Buffer),
					writer: bitio.NewWriter(new(bytes.Buffer)),
				},
			)
		})
	}
}

func TestWithInitialSizeForByteBasedBuffer(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "create it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withInitialSizeForByteBasedBuffer := WithInitialSizeForByteBasedBuffer(tt.args.length)
			assert.NotNilf(t, withInitialSizeForByteBasedBuffer, "WithInitialSizeForByteBasedBuffer(%v)", tt.args.length)
			withInitialSizeForByteBasedBuffer(
				&byteWriteBuffer{
					data:   new(bytes.Buffer),
					writer: bitio.NewWriter(new(bytes.Buffer)),
				},
			)
		})
	}
}

func Test_byteWriteBuffer_GetByteOrder(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
	}
	tests := []struct {
		name   string
		fields fields
		want   binary.ByteOrder
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			assert.Equalf(t, tt.want, wb.GetByteOrder(), "GetByteOrder()")
		})
	}
}

func Test_byteWriteBuffer_GetBytes(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "get it",
			fields: fields{
				data: new(bytes.Buffer),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			assert.Equalf(t, tt.want, wb.GetBytes(), "GetBytes()")
		})
	}
}

func Test_byteWriteBuffer_GetPos(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
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
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			assert.Equalf(t, tt.want, wb.GetPos(), "GetPos()")
		})
	}
}

func Test_byteWriteBuffer_GetTotalBytes(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		{
			name: "get it",
			fields: fields{
				data: new(bytes.Buffer),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			assert.Equalf(t, tt.want, wb.GetTotalBytes(), "GetTotalBytes()")
		})
	}
}

func Test_byteWriteBuffer_PopContext(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
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
			name:    "pop it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			tt.wantErr(t, wb.PopContext(tt.args.in0, tt.args.in1...), fmt.Sprintf("PopContext(%v, %v)", tt.args.in0, tt.args.in1))
		})
	}
}

func Test_byteWriteBuffer_PushContext(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
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
			name:    "push it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			tt.wantErr(t, wb.PushContext(tt.args.in0, tt.args.in1...), fmt.Sprintf("PushContext(%v, %v)", tt.args.in0, tt.args.in1))
		})
	}
}

func Test_byteWriteBuffer_SetByteOrder(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
	}
	type args struct {
		byteOrder binary.ByteOrder
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "set it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			wb.SetByteOrder(tt.args.byteOrder)
		})
	}
}

func Test_byteWriteBuffer_WriteBigFloat(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
	}
	type args struct {
		in0       string
		bitLength uint8
		value     *big.Float
		in3       []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "write it",
			wantErr: assert.Error, // TODO: Not yet implemented
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			tt.wantErr(t, wb.WriteBigFloat(tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteBigFloat(%v, %v, %v, %v)", tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3))
		})
	}
}

func Test_byteWriteBuffer_WriteBigInt(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
	}
	type args struct {
		in0       string
		bitLength uint8
		value     *big.Int
		in3       []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "write it",
			wantErr: assert.Error, // TODO: Not yet implemented
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			tt.wantErr(t, wb.WriteBigInt(tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteBigInt(%v, %v, %v, %v)", tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3))
		})
	}
}

func Test_byteWriteBuffer_WriteBit(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
	}
	type args struct {
		in0   string
		value bool
		in2   []WithWriterArgs
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
				writer: bitio.NewWriter(new(bytes.Buffer)),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			tt.wantErr(t, wb.WriteBit(tt.args.in0, tt.args.value, tt.args.in2...), fmt.Sprintf("WriteBit(%v, %v, %v)", tt.args.in0, tt.args.value, tt.args.in2))
		})
	}
}

func Test_byteWriteBuffer_WriteByte(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
	}
	type args struct {
		in0   string
		value byte
		in2   []WithWriterArgs
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
				writer: bitio.NewWriter(new(bytes.Buffer)),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			tt.wantErr(t, wb.WriteByte(tt.args.in0, tt.args.value, tt.args.in2...), fmt.Sprintf("WriteByte(%v, %v, %v)", tt.args.in0, tt.args.value, tt.args.in2))
		})
	}
}

func Test_byteWriteBuffer_WriteByteArray(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
	}
	type args struct {
		in0  string
		data []byte
		in2  []WithWriterArgs
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
				writer: bitio.NewWriter(new(bytes.Buffer)),
			},
			wantErr: assert.NoError,
		},
		{
			name: "write more",
			fields: fields{
				writer: bitio.NewWriter(new(bytes.Buffer)),
			},
			args: args{
				data: []byte{1, 2, 3, 4},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			tt.wantErr(t, wb.WriteByteArray(tt.args.in0, tt.args.data, tt.args.in2...), fmt.Sprintf("WriteByteArray(%v, %v, %v)", tt.args.in0, tt.args.data, tt.args.in2))
		})
	}
}

func Test_byteWriteBuffer_WriteFloat32(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
	}
	type args struct {
		in0       string
		bitLength uint8
		value     float32
		in3       []WithWriterArgs
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
				writer: bitio.NewWriter(new(bytes.Buffer)),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			tt.wantErr(t, wb.WriteFloat32(tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteFloat32(%v, %v, %v, %v)", tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3))
		})
	}
}

func Test_byteWriteBuffer_WriteFloat64(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
	}
	type args struct {
		in0       string
		bitLength uint8
		value     float64
		in3       []WithWriterArgs
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
				writer: bitio.NewWriter(new(bytes.Buffer)),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			tt.wantErr(t, wb.WriteFloat64(tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteFloat64(%v, %v, %v, %v)", tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3))
		})
	}
}

func Test_byteWriteBuffer_WriteInt16(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
	}
	type args struct {
		in0       string
		bitLength uint8
		value     int16
		in3       []WithWriterArgs
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
				writer: bitio.NewWriter(new(bytes.Buffer)),
			},
			wantErr: assert.NoError,
		},
		{
			name: "write it LE",
			fields: fields{
				writer:    bitio.NewWriter(new(bytes.Buffer)),
				byteOrder: binary.LittleEndian,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			tt.wantErr(t, wb.WriteInt16(tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteInt16(%v, %v, %v, %v)", tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3))
		})
	}
}

func Test_byteWriteBuffer_WriteInt32(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
	}
	type args struct {
		in0       string
		bitLength uint8
		value     int32
		in3       []WithWriterArgs
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
				writer: bitio.NewWriter(new(bytes.Buffer)),
			},
			wantErr: assert.NoError,
		},
		{
			name: "write it LE",
			fields: fields{
				writer:    bitio.NewWriter(new(bytes.Buffer)),
				byteOrder: binary.LittleEndian,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			tt.wantErr(t, wb.WriteInt32(tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteInt32(%v, %v, %v, %v)", tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3))
		})
	}
}

func Test_byteWriteBuffer_WriteInt64(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
	}
	type args struct {
		in0       string
		bitLength uint8
		value     int64
		in3       []WithWriterArgs
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
				writer: bitio.NewWriter(new(bytes.Buffer)),
			},
			wantErr: assert.NoError,
		},
		{
			name: "write it LE",
			fields: fields{
				writer:    bitio.NewWriter(new(bytes.Buffer)),
				byteOrder: binary.LittleEndian,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			tt.wantErr(t, wb.WriteInt64(tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteInt64(%v, %v, %v, %v)", tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3))
		})
	}
}

func Test_byteWriteBuffer_WriteInt8(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
	}
	type args struct {
		in0       string
		bitLength uint8
		value     int8
		in3       []WithWriterArgs
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
				writer: bitio.NewWriter(new(bytes.Buffer)),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			tt.wantErr(t, wb.WriteInt8(tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteInt8(%v, %v, %v, %v)", tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3))
		})
	}
}

func Test_byteWriteBuffer_WriteSerializable(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
	}
	type args struct {
		ctx          context.Context
		serializable Serializable
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(t *testing.T, fields *fields, args *args)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "write it",
			fields: fields{
				writer: bitio.NewWriter(new(bytes.Buffer)),
			},
			wantErr: assert.NoError,
		},
		{
			name: "write it",
			fields: fields{
				writer: bitio.NewWriter(new(bytes.Buffer)),
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				serializable := NewMockSerializable(t)
				serializable.EXPECT().SerializeWithWriteBuffer(mock.Anything, mock.Anything).Return(nil)
				args.serializable = serializable
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields, &tt.args)
			}
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			tt.wantErr(t, wb.WriteSerializable(tt.args.ctx, tt.args.serializable), fmt.Sprintf("WriteSerializable(%v, %v)", tt.args.ctx, tt.args.serializable))
		})
	}
}

func Test_byteWriteBuffer_WriteString(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
	}
	type args struct {
		in0       string
		bitLength uint32
		encoding  string
		value     string
		in4       []WithWriterArgs
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
				writer: bitio.NewWriter(new(bytes.Buffer)),
			},
			args: args{
				bitLength: 48,
				value:     "plc4x",
			},
			wantErr: assert.NoError,
		},
		{
			name: "write it UTF-8",
			fields: fields{
				writer: bitio.NewWriter(new(bytes.Buffer)),
			},
			args: args{
				bitLength: 48,
				encoding:  "UTF8",
				value:     "plc4x",
			},
			wantErr: assert.NoError,
		},
		{
			name: "write it UTF-16",
			fields: fields{
				writer: bitio.NewWriter(new(bytes.Buffer)),
			},
			args: args{
				bitLength: 48,
				encoding:  "UTF16",
				value:     "plc4x",
			},
			wantErr: assert.NoError,
		},
		{
			name: "write it UTF-16 BigEndian",
			fields: fields{
				writer: bitio.NewWriter(new(bytes.Buffer)),
			},
			args: args{
				bitLength: 48,
				encoding:  "UTF16BE",
				value:     "plc4x",
			},
			wantErr: assert.NoError,
		},
		{
			name: "write it UTF-16 LittleEndian",
			fields: fields{
				writer: bitio.NewWriter(new(bytes.Buffer)),
			},
			args: args{
				bitLength: 48,
				encoding:  "UTF16LE",
				value:     "plc4x",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			tt.wantErr(t, wb.WriteString(tt.args.in0, tt.args.bitLength, tt.args.encoding, tt.args.value, tt.args.in4...), fmt.Sprintf("WriteString(%v, %v, %v, %v, %v)", tt.args.in0, tt.args.bitLength, tt.args.encoding, tt.args.value, tt.args.in4))
		})
	}
}

func Test_byteWriteBuffer_WriteUint16(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
	}
	type args struct {
		in0       string
		bitLength uint8
		value     uint16
		in3       []WithWriterArgs
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
				writer: bitio.NewWriter(new(bytes.Buffer)),
			},
			wantErr: assert.NoError,
		},
		{
			name: "write it LE",
			fields: fields{
				writer:    bitio.NewWriter(new(bytes.Buffer)),
				byteOrder: binary.LittleEndian,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			tt.wantErr(t, wb.WriteUint16(tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteUint16(%v, %v, %v, %v)", tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3))
		})
	}
}

func Test_byteWriteBuffer_WriteUint32(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
	}
	type args struct {
		in0       string
		bitLength uint8
		value     uint32
		in3       []WithWriterArgs
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
				writer: bitio.NewWriter(new(bytes.Buffer)),
			},
			wantErr: assert.NoError,
		},

		{
			name: "write it LE",
			fields: fields{
				writer:    bitio.NewWriter(new(bytes.Buffer)),
				byteOrder: binary.LittleEndian,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			tt.wantErr(t, wb.WriteUint32(tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteUint32(%v, %v, %v, %v)", tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3))
		})
	}
}

func Test_byteWriteBuffer_WriteUint64(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
	}
	type args struct {
		in0       string
		bitLength uint8
		value     uint64
		in3       []WithWriterArgs
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
				writer: bitio.NewWriter(new(bytes.Buffer)),
			},
			wantErr: assert.NoError,
		},

		{
			name: "write it LE",
			fields: fields{
				writer:    bitio.NewWriter(new(bytes.Buffer)),
				byteOrder: binary.LittleEndian,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			tt.wantErr(t, wb.WriteUint64(tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteUint64(%v, %v, %v, %v)", tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3))
		})
	}
}

func Test_byteWriteBuffer_WriteUint8(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
	}
	type args struct {
		in0       string
		bitLength uint8
		value     uint8
		in3       []WithWriterArgs
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
				writer: bitio.NewWriter(new(bytes.Buffer)),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			tt.wantErr(t, wb.WriteUint8(tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteUint8(%v, %v, %v, %v)", tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3))
		})
	}
}

func Test_byteWriteBuffer_WriteVirtual(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
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
				writer: bitio.NewWriter(new(bytes.Buffer)),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			tt.wantErr(t, wb.WriteVirtual(tt.args.ctx, tt.args.logicalName, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteVirtual(%v, %v, %v, %v)", tt.args.ctx, tt.args.logicalName, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_byteWriteBuffer_move(t *testing.T) {
	type fields struct {
		data      *bytes.Buffer
		writer    *bitio.Writer
		byteOrder binary.ByteOrder
		pos       uint
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
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			wb.move(tt.args.bits)
		})
	}
}
