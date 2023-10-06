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

package values

import (
	"context"
	"fmt"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestNewWriteBufferPlcValueBased(t *testing.T) {
	tests := []struct {
		name string
		want WriteBufferPlcValueBased
	}{
		{
			name: "create it",
			want: &writeBufferPlcValueBased{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewWriteBufferPlcValueBased(), "NewWriteBufferPlcValueBased()")
		})
	}
}

func Test_writeBufferPlcValueBased_GetPlcValue(t *testing.T) {
	type fields struct {
		BufferCommons utils.BufferCommons
		Stack         utils.Stack
		rootNode      values.PlcValue
		pos           uint
	}
	tests := []struct {
		name   string
		fields fields
		want   values.PlcValue
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &writeBufferPlcValueBased{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				rootNode:      tt.fields.rootNode,
				pos:           tt.fields.pos,
			}
			assert.Equalf(t, tt.want, p.GetPlcValue(), "GetPlcValue()")
		})
	}
}

func Test_writeBufferPlcValueBased_GetPos(t *testing.T) {
	type fields struct {
		BufferCommons utils.BufferCommons
		Stack         utils.Stack
		rootNode      values.PlcValue
		pos           uint
	}
	tests := []struct {
		name   string
		fields fields
		want   uint16
	}{
		{
			name: "get it",
			want: 0,
		},
		{
			name: "get it 7",
			fields: fields{
				pos: 7,
			},
			want: 0,
		},
		{
			name: "get it 8",
			fields: fields{
				pos: 8,
			},
			want: 1,
		},
		{
			name: "get it 16",
			fields: fields{
				pos: 16,
			},
			want: 2,
		},
		{
			name: "get it 17",
			fields: fields{
				pos: 17,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &writeBufferPlcValueBased{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				rootNode:      tt.fields.rootNode,
				pos:           tt.fields.pos,
			}
			assert.Equalf(t, tt.want, p.GetPos(), "GetPos()")
		})
	}
}

func Test_writeBufferPlcValueBased_PopContext(t *testing.T) {
	type fields struct {
		BufferCommons utils.BufferCommons
		Stack         utils.Stack
		rootNode      values.PlcValue
		pos           uint
	}
	type args struct {
		logicalName string
		in1         []utils.WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(t *testing.T, wb *writeBufferPlcValueBased)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "pop it without something on the stack",
			wantErr: assert.Error,
		},
		{
			name: "pop value with unexpected name",
			setup: func(t *testing.T, wb *writeBufferPlcValueBased) {
				wb.Push(&plcValueContext{logicalName: "doesn't matter"})
			},
			wantErr: assert.Error,
		},
		{
			name: "pop value with name",
			args: args{
				logicalName: "doesn't matter",
			},
			setup: func(t *testing.T, wb *writeBufferPlcValueBased) {
				wb.Push(&plcValueContext{
					logicalName: "doesn't matter",
					properties:  map[string]values.PlcValue{},
				})
			},
			wantErr: assert.NoError,
		},
		{
			name: "pop value with name and more elements",
			args: args{
				logicalName: "doesn't matter",
			},
			setup: func(t *testing.T, wb *writeBufferPlcValueBased) {
				wb.Push(&plcValueContext{
					logicalName: "doesn't matter 2",
					properties:  map[string]values.PlcValue{},
				})
				wb.Push(&plcValueContext{
					logicalName: "doesn't matter",
					properties:  map[string]values.PlcValue{},
				})
			},
			wantErr: assert.NoError,
		},
		{
			name: "pop list with name and more elements",
			args: args{
				logicalName: "doesn't matter",
			},
			setup: func(t *testing.T, wb *writeBufferPlcValueBased) {
				wb.Push(&plcListContext{
					logicalName: "doesn't matter 2",
				})
				wb.Push(&plcListContext{
					logicalName: "doesn't matter",
				})
			},
			wantErr: assert.NoError,
		},
		{
			name: "pop list with name and more elements",
			args: args{
				logicalName: "doesn't matter",
			},
			setup: func(t *testing.T, wb *writeBufferPlcValueBased) {
				wb.Push("wat")
				wb.Push(&plcListContext{
					logicalName: "doesn't matter",
				})
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &writeBufferPlcValueBased{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				rootNode:      tt.fields.rootNode,
				pos:           tt.fields.pos,
			}
			if tt.setup != nil {
				tt.setup(t, p)
			}
			tt.wantErr(t, p.PopContext(tt.args.logicalName, tt.args.in1...), fmt.Sprintf("PopContext(%v, %v)", tt.args.logicalName, tt.args.in1))
		})
	}
}

func Test_writeBufferPlcValueBased_PushContext(t *testing.T) {
	type fields struct {
		BufferCommons utils.BufferCommons
		Stack         utils.Stack
		rootNode      values.PlcValue
		pos           uint
	}
	type args struct {
		logicalName string
		writerArgs  []utils.WithWriterArgs
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
		{
			name: "push it as list",
			args: args{
				writerArgs: []utils.WithWriterArgs{
					utils.WithRenderAsList(true),
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &writeBufferPlcValueBased{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				rootNode:      tt.fields.rootNode,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, p.PushContext(tt.args.logicalName, tt.args.writerArgs...), fmt.Sprintf("PushContext(%v, %v)", tt.args.logicalName, tt.args.writerArgs))
		})
	}
}

func Test_writeBufferPlcValueBased_WriteBigFloat(t *testing.T) {
	type fields struct {
		BufferCommons utils.BufferCommons
		Stack         utils.Stack
		rootNode      values.PlcValue
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint8
		value       *big.Float
		in3         []utils.WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "write it (failing)",
			wantErr: assert.Error,
		},
		{
			name: "write it",
			args: args{
				value: func() *big.Float {
					return big.NewFloat(4)
				}(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &writeBufferPlcValueBased{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				rootNode:      tt.fields.rootNode,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, p.WriteBigFloat(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteBigFloat(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.in3))
		})
	}
}

func Test_writeBufferPlcValueBased_WriteBigInt(t *testing.T) {
	type fields struct {
		BufferCommons utils.BufferCommons
		Stack         utils.Stack
		rootNode      values.PlcValue
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint8
		value       *big.Int
		in3         []utils.WithWriterArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "write it (failing)",
			wantErr: assert.Error,
		},
		{
			name: "write it",
			args: args{
				value: func() *big.Int {
					return big.NewInt(4)
				}(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &writeBufferPlcValueBased{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				rootNode:      tt.fields.rootNode,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, p.WriteBigInt(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteBigInt(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.in3))
		})
	}
}

func Test_writeBufferPlcValueBased_WriteBit(t *testing.T) {
	type fields struct {
		BufferCommons utils.BufferCommons
		Stack         utils.Stack
		rootNode      values.PlcValue
		pos           uint
	}
	type args struct {
		logicalName string
		value       bool
		in2         []utils.WithWriterArgs
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
			p := &writeBufferPlcValueBased{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				rootNode:      tt.fields.rootNode,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, p.WriteBit(tt.args.logicalName, tt.args.value, tt.args.in2...), fmt.Sprintf("WriteBit(%v, %v, %v)", tt.args.logicalName, tt.args.value, tt.args.in2))
		})
	}
}

func Test_writeBufferPlcValueBased_WriteByte(t *testing.T) {
	type fields struct {
		BufferCommons utils.BufferCommons
		Stack         utils.Stack
		rootNode      values.PlcValue
		pos           uint
	}
	type args struct {
		logicalName string
		value       byte
		in2         []utils.WithWriterArgs
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
			p := &writeBufferPlcValueBased{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				rootNode:      tt.fields.rootNode,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, p.WriteByte(tt.args.logicalName, tt.args.value, tt.args.in2...), fmt.Sprintf("WriteByte(%v, %v, %v)", tt.args.logicalName, tt.args.value, tt.args.in2))
		})
	}
}

func Test_writeBufferPlcValueBased_WriteByteArray(t *testing.T) {
	type fields struct {
		BufferCommons utils.BufferCommons
		Stack         utils.Stack
		rootNode      values.PlcValue
		pos           uint
	}
	type args struct {
		logicalName string
		data        []byte
		in2         []utils.WithWriterArgs
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
			p := &writeBufferPlcValueBased{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				rootNode:      tt.fields.rootNode,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, p.WriteByteArray(tt.args.logicalName, tt.args.data, tt.args.in2...), fmt.Sprintf("WriteByteArray(%v, %v, %v)", tt.args.logicalName, tt.args.data, tt.args.in2))
		})
	}
}

func Test_writeBufferPlcValueBased_WriteFloat32(t *testing.T) {
	type fields struct {
		BufferCommons utils.BufferCommons
		Stack         utils.Stack
		rootNode      values.PlcValue
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint8
		value       float32
		in3         []utils.WithWriterArgs
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
			p := &writeBufferPlcValueBased{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				rootNode:      tt.fields.rootNode,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, p.WriteFloat32(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteFloat32(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.in3))
		})
	}
}

func Test_writeBufferPlcValueBased_WriteFloat64(t *testing.T) {
	type fields struct {
		BufferCommons utils.BufferCommons
		Stack         utils.Stack
		rootNode      values.PlcValue
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint8
		value       float64
		in3         []utils.WithWriterArgs
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
			p := &writeBufferPlcValueBased{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				rootNode:      tt.fields.rootNode,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, p.WriteFloat64(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteFloat64(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.in3))
		})
	}
}

func Test_writeBufferPlcValueBased_WriteInt16(t *testing.T) {
	type fields struct {
		BufferCommons utils.BufferCommons
		Stack         utils.Stack
		rootNode      values.PlcValue
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint8
		value       int16
		in3         []utils.WithWriterArgs
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
			p := &writeBufferPlcValueBased{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				rootNode:      tt.fields.rootNode,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, p.WriteInt16(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteInt16(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.in3))
		})
	}
}

func Test_writeBufferPlcValueBased_WriteInt32(t *testing.T) {
	type fields struct {
		BufferCommons utils.BufferCommons
		Stack         utils.Stack
		rootNode      values.PlcValue
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint8
		value       int32
		in3         []utils.WithWriterArgs
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
			p := &writeBufferPlcValueBased{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				rootNode:      tt.fields.rootNode,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, p.WriteInt32(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteInt32(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.in3))
		})
	}
}

func Test_writeBufferPlcValueBased_WriteInt64(t *testing.T) {
	type fields struct {
		BufferCommons utils.BufferCommons
		Stack         utils.Stack
		rootNode      values.PlcValue
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint8
		value       int64
		in3         []utils.WithWriterArgs
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
			p := &writeBufferPlcValueBased{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				rootNode:      tt.fields.rootNode,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, p.WriteInt64(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteInt64(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.in3))
		})
	}
}

func Test_writeBufferPlcValueBased_WriteInt8(t *testing.T) {
	type fields struct {
		BufferCommons utils.BufferCommons
		Stack         utils.Stack
		rootNode      values.PlcValue
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint8
		value       int8
		in3         []utils.WithWriterArgs
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
			p := &writeBufferPlcValueBased{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				rootNode:      tt.fields.rootNode,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, p.WriteInt8(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteInt8(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.in3))
		})
	}
}

func Test_writeBufferPlcValueBased_WriteSerializable(t *testing.T) {
	type fields struct {
		BufferCommons utils.BufferCommons
		Stack         utils.Stack
		rootNode      values.PlcValue
		pos           uint
	}
	type args struct {
		ctx          context.Context
		serializable utils.Serializable
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "write it (no serializable)",
			wantErr: assert.NoError,
		},
		{
			name: "write it",
			args: args{
				ctx:          context.Background(),
				serializable: NewPlcBOOL(true),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &writeBufferPlcValueBased{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				rootNode:      tt.fields.rootNode,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, p.WriteSerializable(tt.args.ctx, tt.args.serializable), fmt.Sprintf("WriteSerializable(%v, %v)", tt.args.ctx, tt.args.serializable))
		})
	}
}

func Test_writeBufferPlcValueBased_WriteString(t *testing.T) {
	type fields struct {
		BufferCommons utils.BufferCommons
		Stack         utils.Stack
		rootNode      values.PlcValue
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint32
		in2         string
		value       string
		in4         []utils.WithWriterArgs
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
			p := &writeBufferPlcValueBased{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				rootNode:      tt.fields.rootNode,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, p.WriteString(tt.args.logicalName, tt.args.bitLength, tt.args.in2, tt.args.value, tt.args.in4...), fmt.Sprintf("WriteString(%v, %v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.in2, tt.args.value, tt.args.in4))
		})
	}
}

func Test_writeBufferPlcValueBased_WriteUint16(t *testing.T) {
	type fields struct {
		BufferCommons utils.BufferCommons
		Stack         utils.Stack
		rootNode      values.PlcValue
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint8
		value       uint16
		in3         []utils.WithWriterArgs
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
			p := &writeBufferPlcValueBased{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				rootNode:      tt.fields.rootNode,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, p.WriteUint16(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteUint16(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.in3))
		})
	}
}

func Test_writeBufferPlcValueBased_WriteUint32(t *testing.T) {
	type fields struct {
		BufferCommons utils.BufferCommons
		Stack         utils.Stack
		rootNode      values.PlcValue
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint8
		value       uint32
		in3         []utils.WithWriterArgs
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
			p := &writeBufferPlcValueBased{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				rootNode:      tt.fields.rootNode,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, p.WriteUint32(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteUint32(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.in3))
		})
	}
}

func Test_writeBufferPlcValueBased_WriteUint64(t *testing.T) {
	type fields struct {
		BufferCommons utils.BufferCommons
		Stack         utils.Stack
		rootNode      values.PlcValue
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint8
		value       uint64
		in3         []utils.WithWriterArgs
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
			p := &writeBufferPlcValueBased{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				rootNode:      tt.fields.rootNode,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, p.WriteUint64(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteUint64(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.in3))
		})
	}
}

func Test_writeBufferPlcValueBased_WriteUint8(t *testing.T) {
	type fields struct {
		BufferCommons utils.BufferCommons
		Stack         utils.Stack
		rootNode      values.PlcValue
		pos           uint
	}
	type args struct {
		logicalName string
		bitLength   uint8
		value       uint8
		in3         []utils.WithWriterArgs
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
			p := &writeBufferPlcValueBased{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				rootNode:      tt.fields.rootNode,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, p.WriteUint8(tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteUint8(%v, %v, %v, %v)", tt.args.logicalName, tt.args.bitLength, tt.args.value, tt.args.in3))
		})
	}
}

func Test_writeBufferPlcValueBased_WriteVirtual(t *testing.T) {
	type fields struct {
		BufferCommons utils.BufferCommons
		Stack         utils.Stack
		rootNode      values.PlcValue
		pos           uint
	}
	type args struct {
		ctx         context.Context
		logicalName string
		value       any
		in3         []utils.WithWriterArgs
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
			p := &writeBufferPlcValueBased{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				rootNode:      tt.fields.rootNode,
				pos:           tt.fields.pos,
			}
			tt.wantErr(t, p.WriteVirtual(tt.args.ctx, tt.args.logicalName, tt.args.value, tt.args.in3...), fmt.Sprintf("WriteVirtual(%v, %v, %v, %v)", tt.args.ctx, tt.args.logicalName, tt.args.value, tt.args.in3))
		})
	}
}

func Test_writeBufferPlcValueBased_appendValue(t *testing.T) {
	type fields struct {
		BufferCommons utils.BufferCommons
		Stack         utils.Stack
		rootNode      values.PlcValue
		pos           uint
	}
	type args struct {
		logicalName string
		value       values.PlcValue
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(t *testing.T, wb *writeBufferPlcValueBased)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "append it",
			wantErr: assert.NoError,
		},
		{
			name: "append it value",
			args: args{
				value: NewPlcBOOL(true),
			},
			setup: func(t *testing.T, wb *writeBufferPlcValueBased) {
				wb.Push(&plcValueContext{
					properties: map[string]values.PlcValue{},
				})
			},
			wantErr: assert.NoError,
		},
		{
			name: "append it list",
			args: args{
				value: NewPlcBOOL(true),
			},
			setup: func(t *testing.T, wb *writeBufferPlcValueBased) {
				wb.Push(&plcListContext{})
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &writeBufferPlcValueBased{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				rootNode:      tt.fields.rootNode,
				pos:           tt.fields.pos,
			}
			if tt.setup != nil {
				tt.setup(t, p)
			}
			tt.wantErr(t, p.appendValue(tt.args.logicalName, tt.args.value), fmt.Sprintf("appendValue(%v, %v)", tt.args.logicalName, tt.args.value))
		})
	}
}

func Test_writeBufferPlcValueBased_move(t *testing.T) {
	type fields struct {
		BufferCommons utils.BufferCommons
		Stack         utils.Stack
		rootNode      values.PlcValue
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
			p := &writeBufferPlcValueBased{
				BufferCommons: tt.fields.BufferCommons,
				Stack:         tt.fields.Stack,
				rootNode:      tt.fields.rootNode,
				pos:           tt.fields.pos,
			}
			p.move(tt.args.bits)
		})
	}
}
