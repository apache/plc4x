//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package utils

import (
	"encoding/binary"
	"github.com/icza/bitio"
	"math/big"
	"reflect"
	"testing"
)

func TestNewLittleEndianReadBuffer(t *testing.T) {
	type args struct {
		data []uint8
	}
	tests := []struct {
		name string
		args args
		want *ReadBuffer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLittleEndianReadBuffer(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLittleEndianReadBuffer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewReadBuffer(t *testing.T) {
	type args struct {
		data []uint8
	}
	tests := []struct {
		name string
		args args
		want *ReadBuffer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewReadBuffer(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReadBuffer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadBuffer_GetBytes(t *testing.T) {
	type fields struct {
		data      []uint8
		reader    *bitio.Reader
		pos       uint64
		byteOrder binary.ByteOrder
	}
	tests := []struct {
		name   string
		fields fields
		want   []uint8
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := &ReadBuffer{
				data:      tt.fields.data,
				reader:    tt.fields.reader,
				pos:       tt.fields.pos,
				byteOrder: tt.fields.byteOrder,
			}
			if got := rb.GetBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadBuffer_GetPos(t *testing.T) {
	type fields struct {
		data      []uint8
		reader    *bitio.Reader
		pos       uint64
		byteOrder binary.ByteOrder
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
			rb := &ReadBuffer{
				data:      tt.fields.data,
				reader:    tt.fields.reader,
				pos:       tt.fields.pos,
				byteOrder: tt.fields.byteOrder,
			}
			if got := rb.GetPos(); got != tt.want {
				t.Errorf("GetPos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadBuffer_GetTotalBytes(t *testing.T) {
	type fields struct {
		data      []uint8
		reader    *bitio.Reader
		pos       uint64
		byteOrder binary.ByteOrder
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
			rb := &ReadBuffer{
				data:      tt.fields.data,
				reader:    tt.fields.reader,
				pos:       tt.fields.pos,
				byteOrder: tt.fields.byteOrder,
			}
			if got := rb.GetTotalBytes(); got != tt.want {
				t.Errorf("GetTotalBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadBuffer_HasMore(t *testing.T) {
	type fields struct {
		data      []uint8
		reader    *bitio.Reader
		pos       uint64
		byteOrder binary.ByteOrder
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
			rb := &ReadBuffer{
				data:      tt.fields.data,
				reader:    tt.fields.reader,
				pos:       tt.fields.pos,
				byteOrder: tt.fields.byteOrder,
			}
			if got := rb.HasMore(tt.args.bitLength); got != tt.want {
				t.Errorf("HasMore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadBuffer_PeekByte(t *testing.T) {
	type fields struct {
		data      []uint8
		reader    *bitio.Reader
		pos       uint64
		byteOrder binary.ByteOrder
	}
	type args struct {
		offset uint8
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint8
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := &ReadBuffer{
				data:      tt.fields.data,
				reader:    tt.fields.reader,
				pos:       tt.fields.pos,
				byteOrder: tt.fields.byteOrder,
			}
			if got := rb.PeekByte(tt.args.offset); got != tt.want {
				t.Errorf("PeekByte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadBuffer_ReadBigFloat(t *testing.T) {
	type fields struct {
		data      []uint8
		reader    *bitio.Reader
		pos       uint64
		byteOrder binary.ByteOrder
	}
	type args struct {
		signed            bool
		exponentBitLength uint8
		mantissaBitLength uint8
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *big.Float
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := &ReadBuffer{
				data:      tt.fields.data,
				reader:    tt.fields.reader,
				pos:       tt.fields.pos,
				byteOrder: tt.fields.byteOrder,
			}
			got, err := rb.ReadBigFloat(tt.args.signed, tt.args.exponentBitLength, tt.args.mantissaBitLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadBigFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadBigFloat() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadBuffer_ReadBigInt(t *testing.T) {
	type fields struct {
		data      []uint8
		reader    *bitio.Reader
		pos       uint64
		byteOrder binary.ByteOrder
	}
	type args struct {
		bitLength uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *big.Int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := &ReadBuffer{
				data:      tt.fields.data,
				reader:    tt.fields.reader,
				pos:       tt.fields.pos,
				byteOrder: tt.fields.byteOrder,
			}
			got, err := rb.ReadBigInt(tt.args.bitLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadBigInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadBigInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadBuffer_ReadBit(t *testing.T) {
	type fields struct {
		data      []uint8
		reader    *bitio.Reader
		pos       uint64
		byteOrder binary.ByteOrder
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := &ReadBuffer{
				data:      tt.fields.data,
				reader:    tt.fields.reader,
				pos:       tt.fields.pos,
				byteOrder: tt.fields.byteOrder,
			}
			got, err := rb.ReadBit()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadBit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadBit() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadBuffer_ReadFloat32(t *testing.T) {
	type fields struct {
		data      []uint8
		reader    *bitio.Reader
		pos       uint64
		byteOrder binary.ByteOrder
	}
	type args struct {
		signed            bool
		exponentBitLength uint8
		mantissaBitLength uint8
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float32
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := &ReadBuffer{
				data:      tt.fields.data,
				reader:    tt.fields.reader,
				pos:       tt.fields.pos,
				byteOrder: tt.fields.byteOrder,
			}
			got, err := rb.ReadFloat32(tt.args.signed, tt.args.exponentBitLength, tt.args.mantissaBitLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFloat32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadFloat32() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadBuffer_ReadFloat64(t *testing.T) {
	type fields struct {
		data      []uint8
		reader    *bitio.Reader
		pos       uint64
		byteOrder binary.ByteOrder
	}
	type args struct {
		signed            bool
		exponentBitLength uint8
		mantissaBitLength uint8
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := &ReadBuffer{
				data:      tt.fields.data,
				reader:    tt.fields.reader,
				pos:       tt.fields.pos,
				byteOrder: tt.fields.byteOrder,
			}
			got, err := rb.ReadFloat64(tt.args.signed, tt.args.exponentBitLength, tt.args.mantissaBitLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadFloat64() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadBuffer_ReadInt16(t *testing.T) {
	type fields struct {
		data      []uint8
		reader    *bitio.Reader
		pos       uint64
		byteOrder binary.ByteOrder
	}
	type args struct {
		bitLength uint8
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int16
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := &ReadBuffer{
				data:      tt.fields.data,
				reader:    tt.fields.reader,
				pos:       tt.fields.pos,
				byteOrder: tt.fields.byteOrder,
			}
			got, err := rb.ReadInt16(tt.args.bitLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadInt16() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadInt16() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadBuffer_ReadInt32(t *testing.T) {
	type fields struct {
		data      []uint8
		reader    *bitio.Reader
		pos       uint64
		byteOrder binary.ByteOrder
	}
	type args struct {
		bitLength uint8
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int32
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := &ReadBuffer{
				data:      tt.fields.data,
				reader:    tt.fields.reader,
				pos:       tt.fields.pos,
				byteOrder: tt.fields.byteOrder,
			}
			got, err := rb.ReadInt32(tt.args.bitLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadInt32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadInt32() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadBuffer_ReadInt64(t *testing.T) {
	type fields struct {
		data      []uint8
		reader    *bitio.Reader
		pos       uint64
		byteOrder binary.ByteOrder
	}
	type args struct {
		bitLength uint8
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := &ReadBuffer{
				data:      tt.fields.data,
				reader:    tt.fields.reader,
				pos:       tt.fields.pos,
				byteOrder: tt.fields.byteOrder,
			}
			got, err := rb.ReadInt64(tt.args.bitLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadInt64() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadBuffer_ReadInt8(t *testing.T) {
	type fields struct {
		data      []uint8
		reader    *bitio.Reader
		pos       uint64
		byteOrder binary.ByteOrder
	}
	type args struct {
		bitLength uint8
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int8
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := &ReadBuffer{
				data:      tt.fields.data,
				reader:    tt.fields.reader,
				pos:       tt.fields.pos,
				byteOrder: tt.fields.byteOrder,
			}
			got, err := rb.ReadInt8(tt.args.bitLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadInt8() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadInt8() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadBuffer_ReadString(t *testing.T) {
	type fields struct {
		data      []uint8
		reader    *bitio.Reader
		pos       uint64
		byteOrder binary.ByteOrder
	}
	type args struct {
		bitLength uint32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := &ReadBuffer{
				data:      tt.fields.data,
				reader:    tt.fields.reader,
				pos:       tt.fields.pos,
				byteOrder: tt.fields.byteOrder,
			}
			got, err := rb.ReadString(tt.args.bitLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadBuffer_ReadUint16(t *testing.T) {
	type fields struct {
		data      []uint8
		reader    *bitio.Reader
		pos       uint64
		byteOrder binary.ByteOrder
	}
	type args struct {
		bitLength uint8
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint16
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := &ReadBuffer{
				data:      tt.fields.data,
				reader:    tt.fields.reader,
				pos:       tt.fields.pos,
				byteOrder: tt.fields.byteOrder,
			}
			got, err := rb.ReadUint16(tt.args.bitLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadUint16() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadUint16() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadBuffer_ReadUint32(t *testing.T) {
	type fields struct {
		data      []uint8
		reader    *bitio.Reader
		pos       uint64
		byteOrder binary.ByteOrder
	}
	type args struct {
		bitLength uint8
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint32
		wantErr bool
	}{
		{
			name: "1 BE",
			fields: func() fields {
				buffer := NewReadBuffer([]uint8{0x0, 0x0, 0x0, 0x1})
				return fields{
					data:      buffer.data,
					reader:    buffer.reader,
					pos:       buffer.pos,
					byteOrder: buffer.byteOrder,
				}
			}(),
			args:    args{bitLength: 32},
			want:    1,
			wantErr: false,
		},
		{
			name: "16 BE",
			fields: func() fields {
				buffer := NewReadBuffer([]uint8{0x0, 0x0, 0x0, 0x10})
				return fields{
					data:      buffer.data,
					reader:    buffer.reader,
					pos:       buffer.pos,
					byteOrder: buffer.byteOrder,
				}
			}(),
			args:    args{bitLength: 32},
			want:    16,
			wantErr: false,
		},
		{
			name: "256 BE",
			fields: func() fields {
				buffer := NewReadBuffer([]uint8{0x0, 0x0, 0x1, 0x0})
				return fields{
					data:      buffer.data,
					reader:    buffer.reader,
					pos:       buffer.pos,
					byteOrder: buffer.byteOrder,
				}
			}(),
			args:    args{bitLength: 32},
			want:    256,
			wantErr: false,
		},
		{
			name: "65536 BE",
			fields: func() fields {
				buffer := NewReadBuffer([]uint8{0x0, 0x1, 0x0, 0x0})
				return fields{
					data:      buffer.data,
					reader:    buffer.reader,
					pos:       buffer.pos,
					byteOrder: buffer.byteOrder,
				}
			}(),
			args:    args{bitLength: 32},
			want:    65536,
			wantErr: false,
		},
		{
			name: "16777216 BE",
			fields: func() fields {
				buffer := NewReadBuffer([]uint8{0x1, 0x0, 0x0, 0x0})
				return fields{
					data:      buffer.data,
					reader:    buffer.reader,
					pos:       buffer.pos,
					byteOrder: buffer.byteOrder,
				}
			}(),
			args:    args{bitLength: 32},
			want:    16777216,
			wantErr: false,
		},
		{
			name: "1 LE",
			fields: func() fields {
				buffer := NewLittleEndianReadBuffer([]uint8{0x1, 0x0, 0x0, 0x0})
				return fields{
					data:      buffer.data,
					reader:    buffer.reader,
					pos:       buffer.pos,
					byteOrder: buffer.byteOrder,
				}
			}(),
			args:    args{bitLength: 32},
			want:    1,
			wantErr: false,
		},
		{
			name: "16 LE",
			fields: func() fields {
				buffer := NewLittleEndianReadBuffer([]uint8{0x10, 0x0, 0x0, 0x0})
				return fields{
					data:      buffer.data,
					reader:    buffer.reader,
					pos:       buffer.pos,
					byteOrder: buffer.byteOrder,
				}
			}(),
			args:    args{bitLength: 32},
			want:    16,
			wantErr: false,
		},
		{
			name: "256 LE",
			fields: func() fields {
				buffer := NewLittleEndianReadBuffer([]uint8{0x0, 0x1, 0x0, 0x0})
				return fields{
					data:      buffer.data,
					reader:    buffer.reader,
					pos:       buffer.pos,
					byteOrder: buffer.byteOrder,
				}
			}(),
			args:    args{bitLength: 32},
			want:    256,
			wantErr: false,
		},
		{
			name: "65536 LE",
			fields: func() fields {
				buffer := NewLittleEndianReadBuffer([]uint8{0x0, 0x0, 0x1, 0x0})
				return fields{
					data:      buffer.data,
					reader:    buffer.reader,
					pos:       buffer.pos,
					byteOrder: buffer.byteOrder,
				}
			}(),
			args:    args{bitLength: 32},
			want:    65536,
			wantErr: false,
		},
		{
			name: "16777216 LE",
			fields: func() fields {
				buffer := NewLittleEndianReadBuffer([]uint8{0x0, 0x0, 0x0, 0x1})
				return fields{
					data:      buffer.data,
					reader:    buffer.reader,
					pos:       buffer.pos,
					byteOrder: buffer.byteOrder,
				}
			}(),
			args:    args{bitLength: 32},
			want:    16777216,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := &ReadBuffer{
				data:      tt.fields.data,
				reader:    tt.fields.reader,
				pos:       tt.fields.pos,
				byteOrder: tt.fields.byteOrder,
			}
			got, err := rb.ReadUint32(tt.args.bitLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadUint32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadUint32() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadBuffer_ReadUint64(t *testing.T) {
	type fields struct {
		data      []uint8
		reader    *bitio.Reader
		pos       uint64
		byteOrder binary.ByteOrder
	}
	type args struct {
		bitLength uint8
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := &ReadBuffer{
				data:      tt.fields.data,
				reader:    tt.fields.reader,
				pos:       tt.fields.pos,
				byteOrder: tt.fields.byteOrder,
			}
			got, err := rb.ReadUint64(tt.args.bitLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadUint64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadUint64() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadBuffer_ReadUint8(t *testing.T) {
	type fields struct {
		data      []uint8
		reader    *bitio.Reader
		pos       uint64
		byteOrder binary.ByteOrder
	}
	type args struct {
		bitLength uint8
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint8
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := &ReadBuffer{
				data:      tt.fields.data,
				reader:    tt.fields.reader,
				pos:       tt.fields.pos,
				byteOrder: tt.fields.byteOrder,
			}
			got, err := rb.ReadUint8(tt.args.bitLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadUint8() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadUint8() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadBuffer_Reset(t *testing.T) {
	type fields struct {
		data      []uint8
		reader    *bitio.Reader
		pos       uint64
		byteOrder binary.ByteOrder
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := &ReadBuffer{
				data:      tt.fields.data,
				reader:    tt.fields.reader,
				pos:       tt.fields.pos,
				byteOrder: tt.fields.byteOrder,
			}
			rb.Reset()
		})
	}
}
