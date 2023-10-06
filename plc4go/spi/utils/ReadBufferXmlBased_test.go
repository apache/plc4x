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
	"encoding/xml"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"math/big"
	"strings"
	"testing"
)

func TestNewStrictXmlReadBuffer(t *testing.T) {
	type args struct {
		reader       io.Reader
		validateAttr bool
		validateList bool
	}
	tests := []struct {
		name string
		args args
		want ReadBuffer
	}{
		{
			name: "create it",
			args: args{
				reader:       strings.NewReader(""),
				validateAttr: true,
				validateList: true,
			},
			want: &xmlReadBuffer{
				Decoder:        xml.NewDecoder(strings.NewReader("")),
				pos:            1,
				doValidateAttr: true,
				doValidateList: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewStrictXmlReadBuffer(tt.args.reader, tt.args.validateAttr, tt.args.validateList), "NewStrictXmlReadBuffer(%v, %v, %v)", tt.args.reader, tt.args.validateAttr, tt.args.validateList)
		})
	}
}

func TestNewXmlReadBuffer(t *testing.T) {
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
				reader: strings.NewReader(""),
			},
			want: &xmlReadBuffer{
				Decoder:        xml.NewDecoder(strings.NewReader("")),
				pos:            1,
				doValidateAttr: false,
				doValidateList: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewXmlReadBuffer(tt.args.reader), "NewXmlReadBuffer(%v)", tt.args.reader)
		})
	}
}

func Test_xmlReadBuffer_CloseContext(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
	}
	type args struct {
		logicalName string
		in1         []WithReaderArgs
	}
	tests := []struct {
		name                     string
		fields                   fields
		args                     args
		wantErr                  assert.ErrorAssertionFunc
		travelToNextStartElement bool
	}{
		{
			name: "close it EOF",
			fields: fields{
				Decoder: xml.NewDecoder(strings.NewReader("")),
			},
			wantErr: assert.Error,
		},
		{
			name: "close it",
			fields: fields{
				Decoder: xml.NewDecoder(strings.NewReader("<xml></xml>")),
			},
			args: args{
				logicalName: "xml",
			},
			wantErr:                  assert.NoError,
			travelToNextStartElement: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			if tt.travelToNextStartElement {
				element, err := x.travelToNextStartElement()
				assert.NoError(t, err)
				assert.NotNil(t, element)
			}
			tt.wantErr(t, x.CloseContext(tt.args.logicalName, tt.args.in1...), fmt.Sprintf("CloseContext(%v, %v)", tt.args.logicalName, tt.args.in1))
		})
	}
}

func Test_xmlReadBuffer_GetPos(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
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
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			assert.Equalf(t, tt.want, x.GetPos(), "GetPos()")
		})
	}
}

func Test_xmlReadBuffer_HasMore(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
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
		{
			name: "has no more",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			assert.Equalf(t, tt.want, x.HasMore(tt.args.bitLength), "HasMore(%v)", tt.args.bitLength)
		})
	}
}

func Test_xmlReadBuffer_PullContext(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
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
			name: "pull it EOF",
			fields: fields{
				Decoder: xml.NewDecoder(strings.NewReader("")),
			},
			wantErr: assert.Error,
		},
		{
			name: "pull it",
			fields: fields{
				Decoder: xml.NewDecoder(strings.NewReader("<xml></xml>")),
			},
			args: args{
				logicalName: "xml",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			tt.wantErr(t, x.PullContext(tt.args.logicalName, tt.args.readerArgs...), fmt.Sprintf("PullContext(%v, %v)", tt.args.logicalName, tt.args.readerArgs))
		})
	}
}

func Test_xmlReadBuffer_ReadBigFloat(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
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
				Decoder: xml.NewDecoder(strings.NewReader("<xml>123</xml>")),
			},
			args: args{
				logicalName: "xml",
			},
			want:    big.NewFloat(123).SetPrec(64),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			got, err := x.ReadBigFloat(tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadBigFloat(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadBigFloat(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)
		})
	}
}

func Test_xmlReadBuffer_ReadBigInt(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
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
				Decoder: xml.NewDecoder(strings.NewReader("<xml>123</xml>")),
			},
			args: args{
				logicalName: "xml",
			},
			want:    big.NewInt(123),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			got, err := x.ReadBigInt(tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadBigInt(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadBigInt(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)
		})
	}
}

func Test_xmlReadBuffer_ReadBit(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
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
				Decoder: xml.NewDecoder(strings.NewReader("<xml>true</xml>")),
			},
			args: args{
				logicalName: "xml",
			},
			want:    true,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			got, err := x.ReadBit(tt.args.logicalName, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadBit(%v, %v)", tt.args.logicalName, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadBit(%v, %v)", tt.args.logicalName, tt.args.readerArgs)
		})
	}
}

func Test_xmlReadBuffer_ReadByte(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
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
				Decoder: xml.NewDecoder(strings.NewReader("<xml>0x12</xml>")),
			},
			args: args{
				logicalName: "xml",
			},
			want:    0x12,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			got, err := x.ReadByte(tt.args.logicalName, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadByte(%v, %v)", tt.args.logicalName, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadByte(%v, %v)", tt.args.logicalName, tt.args.readerArgs)
		})
	}
}

func Test_xmlReadBuffer_ReadByteArray(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
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
				Decoder: xml.NewDecoder(strings.NewReader("<xml>0xAFFE</xml>")),
			},
			args: args{
				logicalName: "xml",
			},
			want:    []byte{0xAF, 0xFE},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			got, err := x.ReadByteArray(tt.args.logicalName, tt.args.numberOfBytes, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadByteArray(%v, %v, %v)", tt.args.logicalName, tt.args.numberOfBytes, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadByteArray(%v, %v, %v)", tt.args.logicalName, tt.args.numberOfBytes, tt.args.readerArgs)
		})
	}
}

func Test_xmlReadBuffer_ReadFloat32(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
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
				Decoder: xml.NewDecoder(strings.NewReader("<xml>123</xml>")),
			},
			args: args{
				logicalName: "xml",
			},
			want:    123,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			got, err := x.ReadFloat32(tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadFloat32(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadFloat32(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)
		})
	}
}

func Test_xmlReadBuffer_ReadFloat64(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
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
				Decoder: xml.NewDecoder(strings.NewReader("<xml>123</xml>")),
			},
			args: args{
				logicalName: "xml",
			},
			want:    123,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			got, err := x.ReadFloat64(tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadFloat64(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadFloat64(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)
		})
	}
}

func Test_xmlReadBuffer_ReadInt16(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
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
				Decoder: xml.NewDecoder(strings.NewReader("<xml>123</xml>")),
			},
			args: args{
				logicalName: "xml",
			},
			want:    123,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			got, err := x.ReadInt16(tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadInt16(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadInt16(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)
		})
	}
}

func Test_xmlReadBuffer_ReadInt32(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
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
				Decoder: xml.NewDecoder(strings.NewReader("<xml>123</xml>")),
			},
			args: args{
				logicalName: "xml",
			},
			want:    123,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			got, err := x.ReadInt32(tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadInt32(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadInt32(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)
		})
	}
}

func Test_xmlReadBuffer_ReadInt64(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
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
				Decoder: xml.NewDecoder(strings.NewReader("<xml>123</xml>")),
			},
			args: args{
				logicalName: "xml",
			},
			want:    123,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			got, err := x.ReadInt64(tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadInt64(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadInt64(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)
		})
	}
}

func Test_xmlReadBuffer_ReadInt8(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
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
				Decoder: xml.NewDecoder(strings.NewReader("<xml>123</xml>")),
			},
			args: args{
				logicalName: "xml",
			},
			want:    123,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			got, err := x.ReadInt8(tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadInt8(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadInt8(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)
		})
	}
}

func Test_xmlReadBuffer_ReadString(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
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
				Decoder: xml.NewDecoder(strings.NewReader("<xml>123</xml>")),
			},
			args: args{
				logicalName: "xml",
			},
			want:    "123",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			got, err := x.ReadString(tt.args.logicalName, tt.args.bitLength, tt.args.encoding, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadString(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.encoding, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadString(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.encoding, tt.args.readerArgs)
		})
	}
}

func Test_xmlReadBuffer_ReadUint16(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
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
				Decoder: xml.NewDecoder(strings.NewReader("<xml>123</xml>")),
			},
			args: args{
				logicalName: "xml",
			},
			want:    123,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			got, err := x.ReadUint16(tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadUint16(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadUint16(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)
		})
	}
}

func Test_xmlReadBuffer_ReadUint32(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
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
				Decoder: xml.NewDecoder(strings.NewReader("<xml>123</xml>")),
			},
			args: args{
				logicalName: "xml",
			},
			want:    123,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			got, err := x.ReadUint32(tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadUint32(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadUint32(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)
		})
	}
}

func Test_xmlReadBuffer_ReadUint64(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
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
				Decoder: xml.NewDecoder(strings.NewReader("<xml>123</xml>")),
			},
			args: args{
				logicalName: "xml",
			},
			want:    123,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			got, err := x.ReadUint64(tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadUint64(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadUint64(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)
		})
	}
}

func Test_xmlReadBuffer_ReadUint8(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
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
				Decoder: xml.NewDecoder(strings.NewReader("<xml>123</xml>")),
			},
			args: args{
				logicalName: "xml",
			},
			want:    123,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			got, err := x.ReadUint8(tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs...)
			if !tt.wantErr(t, err, fmt.Sprintf("ReadUint8(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ReadUint8(%v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.readerArgs)
		})
	}
}

func Test_xmlReadBuffer_Reset(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
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
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			x.Reset(tt.args.pos)
		})
	}
}

func Test_xmlReadBuffer_decode(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
	}
	type args struct {
		logicalName string
		dataType    string
		bitLength   uint
		readerArgs  []WithReaderArgs
		targetValue any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "decode it",
			args: args{
				logicalName: "xml",
				targetValue: func() any {
					s := ""
					return &s
				}(),
			},
			fields: fields{
				Decoder: xml.NewDecoder(strings.NewReader("<xml>123</xml>")),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			tt.wantErr(t, x.decode(tt.args.logicalName, tt.args.dataType, tt.args.bitLength, tt.args.readerArgs, tt.args.targetValue), fmt.Sprintf("decode(%v, %v, %v, %v, %v)", tt.args.logicalName, tt.args.dataType, tt.args.bitLength, tt.args.readerArgs, tt.args.targetValue))
		})
	}
}

func Test_xmlReadBuffer_move(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
	}
	type args struct {
		bits uint8
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
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			x.move(tt.args.bits)
		})
	}
}

func Test_xmlReadBuffer_travelToNextEndElement(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
	}
	tests := []struct {
		name                     string
		fields                   fields
		want                     xml.EndElement
		wantErr                  assert.ErrorAssertionFunc
		travelToNextStartElement bool
	}{
		{
			name: "travel",
			fields: fields{
				Decoder: xml.NewDecoder(strings.NewReader("<xml></xml>")),
			},
			want:                     xml.EndElement{Name: xml.Name{Local: "xml"}},
			wantErr:                  assert.NoError,
			travelToNextStartElement: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			if tt.travelToNextStartElement {
				element, err := x.travelToNextStartElement()
				assert.NoError(t, err)
				assert.NotNil(t, element)
			}
			got, err := x.travelToNextEndElement()
			if !tt.wantErr(t, err, fmt.Sprintf("travelToNextEndElement()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "travelToNextEndElement()")
		})
	}
}

func Test_xmlReadBuffer_travelToNextStartElement(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    xml.StartElement
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "travel",
			fields: fields{
				Decoder: xml.NewDecoder(strings.NewReader("<xml></xml>")),
			},
			want:    xml.StartElement{Name: xml.Name{Local: "xml"}, Attr: []xml.Attr{}},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			got, err := x.travelToNextStartElement()
			if !tt.wantErr(t, err, fmt.Sprintf("travelToNextStartElement()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "travelToNextStartElement()")
		})
	}
}

func Test_xmlReadBuffer_validateAttr(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
	}
	type args struct {
		attr      []xml.Attr
		dataType  string
		bitLength uint
		in3       []WithReaderArgs
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
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			tt.wantErr(t, x.validateAttr(tt.args.attr, tt.args.dataType, tt.args.bitLength, tt.args.in3...), fmt.Sprintf("validateAttr(%v, %v, %v, %v)", tt.args.attr, tt.args.dataType, tt.args.bitLength, tt.args.in3))
		})
	}
}

func Test_xmlReadBuffer_validateIfList(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
	}
	type args struct {
		readerArgs   []WithReaderArgs
		startElement xml.StartElement
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
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			tt.wantErr(t, x.validateIfList(tt.args.readerArgs, tt.args.startElement), fmt.Sprintf("validateIfList(%v, %v)", tt.args.readerArgs, tt.args.startElement))
		})
	}
}

func Test_xmlReadBuffer_validateStartElement(t *testing.T) {
	type fields struct {
		BufferCommons  BufferCommons
		Decoder        *xml.Decoder
		pos            uint
		doValidateAttr bool
		doValidateList bool
	}
	type args struct {
		startElement xml.StartElement
		logicalName  string
		dataType     string
		bitLength    uint
		readerArgs   []WithReaderArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "validate it",
			args: args{
				startElement: xml.StartElement{Name: xml.Name{Local: "xml"}},
				logicalName:  "xml",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlReadBuffer{
				BufferCommons:  tt.fields.BufferCommons,
				Decoder:        tt.fields.Decoder,
				pos:            tt.fields.pos,
				doValidateAttr: tt.fields.doValidateAttr,
				doValidateList: tt.fields.doValidateList,
			}
			tt.wantErr(t, x.validateStartElement(tt.args.startElement, tt.args.logicalName, tt.args.dataType, tt.args.bitLength, tt.args.readerArgs...), fmt.Sprintf("validateStartElement(%v, %v, %v, %v, %v)", tt.args.startElement, tt.args.logicalName, tt.args.dataType, tt.args.bitLength, tt.args.readerArgs))
		})
	}
}
