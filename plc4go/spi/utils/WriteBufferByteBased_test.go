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
	"math/big"
	"reflect"
	"testing"

	"github.com/icza/bitio"
	"github.com/stretchr/testify/assert"
)

func TestNewWriteBufferByteBased(t *testing.T) {
	type args struct {
		options []WriteBufferByteBasedOptions
	}
	tests := []struct {
		name string
		args args
		want WriteBufferByteBased
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWriteBufferByteBased(tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWriteBufferByteBased() = %v, want %v", got, tt.want)
			}
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
		want WriteBufferByteBasedOptions
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithByteOrderForByteBasedBuffer(tt.args.byteOrder); !reflect.DeepEqual(got, tt.want) {
				//t.Errorf("WithByteOrderForByteBasedBuffer() = %v, want %v", got, tt.want)
			}
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
		want WriteBufferByteBasedOptions
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithCustomBufferForByteBasedBuffer(tt.args.buffer); !reflect.DeepEqual(got, tt.want) {
				//	t.Errorf("WithCustomBufferForByteBasedBuffer() = %v, want %v", got, tt.want)
			}
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
		want WriteBufferByteBasedOptions
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithInitialSizeForByteBasedBuffer(tt.args.length); !reflect.DeepEqual(got, tt.want) {
				//t.Errorf("WithInitialSizeForByteBasedBuffer() = %v, want %v", got, tt.want)
			}
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			if got := wb.GetByteOrder(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByteOrder() = %v, want %v", got, tt.want)
			}
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			if got := wb.GetBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBytes() = %v, want %v", got, tt.want)
			}
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			if got := wb.GetPos(); got != tt.want {
				t.Errorf("GetPos() = %v, want %v", got, tt.want)
			}
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			if got := wb.GetTotalBytes(); got != tt.want {
				t.Errorf("GetTotalBytes() = %v, want %v", got, tt.want)
			}
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
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			if err := wb.PopContext(tt.args.in0, tt.args.in1...); (err != nil) != tt.wantErr {
				t.Errorf("PopContext() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			if err := wb.PushContext(tt.args.in0, tt.args.in1...); (err != nil) != tt.wantErr {
				t.Errorf("PushContext() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		// TODO: Add test cases.
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
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			if err := wb.WriteBigFloat(tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3...); (err != nil) != tt.wantErr {
				t.Errorf("WriteBigFloat() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			if err := wb.WriteBigInt(tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3...); (err != nil) != tt.wantErr {
				t.Errorf("WriteBigInt() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			if err := wb.WriteBit(tt.args.in0, tt.args.value, tt.args.in2...); (err != nil) != tt.wantErr {
				t.Errorf("WriteBit() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			if err := wb.WriteByte(tt.args.in0, tt.args.value, tt.args.in2...); (err != nil) != tt.wantErr {
				t.Errorf("WriteByte() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			if err := wb.WriteByteArray(tt.args.in0, tt.args.data, tt.args.in2...); (err != nil) != tt.wantErr {
				t.Errorf("WriteByteArray() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			if err := wb.WriteFloat32(tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3...); (err != nil) != tt.wantErr {
				t.Errorf("WriteFloat32() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			if err := wb.WriteFloat64(tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3...); (err != nil) != tt.wantErr {
				t.Errorf("WriteFloat64() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			if err := wb.WriteInt16(tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3...); (err != nil) != tt.wantErr {
				t.Errorf("WriteInt16() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			if err := wb.WriteInt32(tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3...); (err != nil) != tt.wantErr {
				t.Errorf("WriteInt32() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			if err := wb.WriteInt64(tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3...); (err != nil) != tt.wantErr {
				t.Errorf("WriteInt64() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			if err := wb.WriteInt8(tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3...); (err != nil) != tt.wantErr {
				t.Errorf("WriteInt8() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		serializable Serializable
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			if err := wb.WriteSerializable(context.Background(), tt.args.serializable); (err != nil) != tt.wantErr {
				t.Errorf("WriteSerializable() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			if err := wb.WriteString(tt.args.in0, tt.args.bitLength, tt.args.encoding, tt.args.value, tt.args.in4...); (err != nil) != tt.wantErr {
				t.Errorf("WriteString() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			if err := wb.WriteUint16(tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3...); (err != nil) != tt.wantErr {
				t.Errorf("WriteUint16() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			if err := wb.WriteUint32(tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3...); (err != nil) != tt.wantErr {
				t.Errorf("WriteUint32() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			if err := wb.WriteUint64(tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3...); (err != nil) != tt.wantErr {
				t.Errorf("WriteUint64() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			if err := wb.WriteUint8(tt.args.in0, tt.args.bitLength, tt.args.value, tt.args.in3...); (err != nil) != tt.wantErr {
				t.Errorf("WriteUint8() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		logicalName string
		value       any
		writerArgs  []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := &byteWriteBuffer{
				data:      tt.fields.data,
				writer:    tt.fields.writer,
				byteOrder: tt.fields.byteOrder,
				pos:       tt.fields.pos,
			}
			if err := wb.WriteVirtual(context.Background(), tt.args.logicalName, tt.args.value, tt.args.writerArgs...); (err != nil) != tt.wantErr {
				t.Errorf("WriteVirtual() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		// TODO: Add test cases.
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
