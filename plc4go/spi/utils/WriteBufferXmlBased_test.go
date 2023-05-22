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
	"encoding/xml"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/big"
	"strings"
	"testing"
)

func TestNewConfiguredXmlWriteBuffer(t *testing.T) {
	type args struct {
		renderLists bool
		renderAttr  bool
	}
	tests := []struct {
		name string
		args args
		want WriteBufferXmlBased
	}{
		{
			name: "create it",
			want: func() WriteBufferXmlBased {
				var xmlString strings.Builder
				encoder := xml.NewEncoder(&xmlString)
				encoder.Indent("", "  ")
				return &xmlWriteBuffer{
					xmlString:     &xmlString,
					Encoder:       encoder,
					doRenderLists: false,
					doRenderAttr:  false,
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewConfiguredXmlWriteBuffer(tt.args.renderLists, tt.args.renderAttr), "NewConfiguredXmlWriteBuffer(%v, %v)", tt.args.renderLists, tt.args.renderAttr)
		})
	}
}

func TestNewXmlWriteBuffer(t *testing.T) {
	tests := []struct {
		name string
		want WriteBufferXmlBased
	}{
		{
			name: "create it",
			want: func() WriteBufferXmlBased {
				var xmlString strings.Builder
				encoder := xml.NewEncoder(&xmlString)
				encoder.Indent("", "  ")
				return &xmlWriteBuffer{
					xmlString:     &xmlString,
					Encoder:       encoder,
					doRenderLists: true,
					doRenderAttr:  true,
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewXmlWriteBuffer(), "NewXmlWriteBuffer()")
		})
	}
}

func Test_xmlWriteBuffer_GetPos(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
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
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			assert.Equalf(t, tt.want, x.GetPos(), "GetPos()")
		})
	}
}

func Test_xmlWriteBuffer_GetXmlString(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
		doRenderAttr  bool
		pos           uint
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get it",
			fields: fields{
				xmlString: &strings.Builder{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			assert.Equalf(t, tt.want, x.GetXmlString(), "GetXmlString()")
		})
	}
}

func Test_xmlWriteBuffer_PopContext(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
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
			name: "pop it (not context)",
			fields: fields{
				Encoder: func() *xml.Encoder {
					var xmlString strings.Builder
					return xml.NewEncoder(&xmlString)
				}(),
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, x.PopContext(tt.args.logicalName, tt.args.in1...), fmt.Sprintf("PopContext(%v, %v)", tt.args.logicalName, tt.args.in1))
		})
	}
}

func Test_xmlWriteBuffer_PushContext(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
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
			name: "push it",
			fields: fields{
				Encoder: func() *xml.Encoder {
					var xmlString strings.Builder
					return xml.NewEncoder(&xmlString)
				}(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, x.PushContext(tt.args.logicalName, tt.args.writerArgs...), fmt.Sprintf("PushContext(%v, %v)", tt.args.logicalName, tt.args.writerArgs))
		})
	}
}

func Test_xmlWriteBuffer_WriteBigFloat(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
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
			name: "write it",
			fields: fields{
				Encoder: func() *xml.Encoder {
					var xmlString strings.Builder
					return xml.NewEncoder(&xmlString)
				}(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, x.WriteBigFloat(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteBigFloat(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_xmlWriteBuffer_WriteBigInt(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
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
			name: "write it",
			fields: fields{
				Encoder: func() *xml.Encoder {
					var xmlString strings.Builder
					return xml.NewEncoder(&xmlString)
				}(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, x.WriteBigInt(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteBigInt(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_xmlWriteBuffer_WriteBit(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
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
			name: "write it",
			fields: fields{
				Encoder: func() *xml.Encoder {
					var xmlString strings.Builder
					return xml.NewEncoder(&xmlString)
				}(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, x.WriteBit(tt.args.logicalName, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteBit(%v, %v, %v)", tt.args.logicalName, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_xmlWriteBuffer_WriteByte(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
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
			name: "write it",
			fields: fields{
				Encoder: func() *xml.Encoder {
					var xmlString strings.Builder
					return xml.NewEncoder(&xmlString)
				}(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, x.WriteByte(tt.args.logicalName, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteByte(%v, %v, %v)", tt.args.logicalName, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_xmlWriteBuffer_WriteByteArray(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
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
			name: "write it",
			fields: fields{
				Encoder: func() *xml.Encoder {
					var xmlString strings.Builder
					return xml.NewEncoder(&xmlString)
				}(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, x.WriteByteArray(tt.args.logicalName, tt.args.data, tt.args.writerArgs...), fmt.Sprintf("WriteByteArray(%v, %v, %v)", tt.args.logicalName, tt.args.data, tt.args.writerArgs))
		})
	}
}

func Test_xmlWriteBuffer_WriteFloat32(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
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
			name: "write it",
			fields: fields{
				Encoder: func() *xml.Encoder {
					var xmlString strings.Builder
					return xml.NewEncoder(&xmlString)
				}(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, x.WriteFloat32(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteFloat32(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_xmlWriteBuffer_WriteFloat64(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
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
			name: "write it",
			fields: fields{
				Encoder: func() *xml.Encoder {
					var xmlString strings.Builder
					return xml.NewEncoder(&xmlString)
				}(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, x.WriteFloat64(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteFloat64(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_xmlWriteBuffer_WriteInt16(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
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
			name: "write it",
			fields: fields{
				Encoder: func() *xml.Encoder {
					var xmlString strings.Builder
					return xml.NewEncoder(&xmlString)
				}(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, x.WriteInt16(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteInt16(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_xmlWriteBuffer_WriteInt32(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
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
			name: "write it",
			fields: fields{
				Encoder: func() *xml.Encoder {
					var xmlString strings.Builder
					return xml.NewEncoder(&xmlString)
				}(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, x.WriteInt32(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteInt32(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_xmlWriteBuffer_WriteInt64(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
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
			name: "write it",
			fields: fields{
				Encoder: func() *xml.Encoder {
					var xmlString strings.Builder
					return xml.NewEncoder(&xmlString)
				}(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, x.WriteInt64(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteInt64(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_xmlWriteBuffer_WriteInt8(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
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
			name: "write it",
			fields: fields{
				Encoder: func() *xml.Encoder {
					var xmlString strings.Builder
					return xml.NewEncoder(&xmlString)
				}(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, x.WriteInt8(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteInt8(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_xmlWriteBuffer_WriteSerializable(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
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
			name: "write it",
			fields: fields{
				Encoder: func() *xml.Encoder {
					var xmlString strings.Builder
					return xml.NewEncoder(&xmlString)
				}(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, x.WriteSerializable(tt.args.ctx, tt.args.serializable), fmt.Sprintf("WriteSerializable(%v, %v)", tt.args.ctx, tt.args.serializable))
		})
	}
}

func Test_xmlWriteBuffer_WriteString(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
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
			name: "write it",
			fields: fields{
				Encoder: func() *xml.Encoder {
					var xmlString strings.Builder
					return xml.NewEncoder(&xmlString)
				}(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, x.WriteString(tt.args.logicalName, tt.args.bitLength, tt.args.encoding, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteString(%v, %v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.encoding, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_xmlWriteBuffer_WriteUint16(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
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
			name: "write it",
			fields: fields{
				Encoder: func() *xml.Encoder {
					var xmlString strings.Builder
					return xml.NewEncoder(&xmlString)
				}(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, x.WriteUint16(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteUint16(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_xmlWriteBuffer_WriteUint32(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
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
			name: "write it",
			fields: fields{
				Encoder: func() *xml.Encoder {
					var xmlString strings.Builder
					return xml.NewEncoder(&xmlString)
				}(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, x.WriteUint32(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteUint32(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_xmlWriteBuffer_WriteUint64(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
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
			name: "write it",
			fields: fields{
				Encoder: func() *xml.Encoder {
					var xmlString strings.Builder
					return xml.NewEncoder(&xmlString)
				}(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, x.WriteUint64(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteUint64(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_xmlWriteBuffer_WriteUint8(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
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
			name: "write it",
			fields: fields{
				Encoder: func() *xml.Encoder {
					var xmlString strings.Builder
					return xml.NewEncoder(&xmlString)
				}(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, x.WriteUint8(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs...), fmt.Sprintf("WriteUint8(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.writerArgs))
		})
	}
}

func Test_xmlWriteBuffer_WriteVirtual(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		in0 context.Context
		in1 string
		in2 any
		in3 []WithWriterArgs
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
				Encoder: func() *xml.Encoder {
					var xmlString strings.Builder
					return xml.NewEncoder(&xmlString)
				}(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, x.WriteVirtual(tt.args.in0, tt.args.in1, tt.args.in2, tt.args.in3...), fmt.Sprintf("WriteVirtual(%v, %v, %v, %v)", tt.args.in0, tt.args.in1, tt.args.in2, tt.args.in3))
		})
	}
}

func Test_xmlWriteBuffer_encodeElement(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		logicalName string
		value       any
		attr        []xml.Attr
		in3         []WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "encode it",
			fields: fields{
				Encoder: func() *xml.Encoder {
					var xmlString strings.Builder
					return xml.NewEncoder(&xmlString)
				}(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, x.encodeElement(tt.args.logicalName, tt.args.value, tt.args.attr, tt.args.in3...), fmt.Sprintf("encodeElement(%v, %v, %v, %v)", tt.args.logicalName, tt.args.value, tt.args.attr, tt.args.in3))
		})
	}
}

func Test_xmlWriteBuffer_generateAttr(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		dataType   string
		bitLength  uint
		writerArgs []WithWriterArgs
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []xml.Attr
	}{
		{
			name: "generate it",
			want: []xml.Attr{{Name: xml.Name{Space: "", Local: ""}, Value: ""}, {Name: xml.Name{Space: "", Local: ""}, Value: ""}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			assert.Equalf(t, tt.want, x.generateAttr(tt.args.dataType, tt.args.bitLength, tt.args.writerArgs...), "generateAttr(%v, %v, %v)", tt.args.dataType, tt.args.bitLength, tt.args.writerArgs)
		})
	}
}

func Test_xmlWriteBuffer_markAsListIfRequired(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
		doRenderAttr  bool
		pos           uint
	}
	type args struct {
		writerArgs []WithWriterArgs
		attrs      []xml.Attr
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []xml.Attr
	}{
		{
			name: "mark it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			assert.Equalf(t, tt.want, x.markAsListIfRequired(tt.args.writerArgs, tt.args.attrs), "markAsListIfRequired(%v, %v)", tt.args.writerArgs, tt.args.attrs)
		})
	}
}

func Test_xmlWriteBuffer_move(t *testing.T) {
	type fields struct {
		BufferCommons BufferCommons
		xmlString     *strings.Builder
		Encoder       *xml.Encoder
		doRenderLists bool
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
			x := &xmlWriteBuffer{
				BufferCommons: tt.fields.BufferCommons,
				xmlString:     tt.fields.xmlString,
				Encoder:       tt.fields.Encoder,
				doRenderLists: tt.fields.doRenderLists,
				doRenderAttr:  tt.fields.doRenderAttr,
				pos:           tt.fields.pos,
			}
			x.move(tt.args.bits)
		})
	}
}
