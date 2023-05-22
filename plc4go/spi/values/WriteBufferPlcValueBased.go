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
	"math/big"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
)

type WriteBufferPlcValueBased interface {
	utils.WriteBuffer
	GetPlcValue() apiValues.PlcValue
}

func NewWriteBufferPlcValueBased() WriteBufferPlcValueBased {
	return &writeBufferPlcValueBased{}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type writeBufferPlcValueBased struct {
	utils.BufferCommons
	utils.Stack
	rootNode apiValues.PlcValue
	pos      uint
}

type plcValueContext struct {
	logicalName string
	properties  map[string]apiValues.PlcValue
}

type plcListContext struct {
	logicalName string
	list        []apiValues.PlcValue
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (p *writeBufferPlcValueBased) PushContext(logicalName string, writerArgs ...utils.WithWriterArgs) error {
	renderedAsList := p.IsToBeRenderedAsList(utils.UpcastWriterArgs(writerArgs...)...)
	if renderedAsList {
		p.Push(&plcListContext{logicalName, make([]apiValues.PlcValue, 0)})
	} else {
		p.Push(&plcValueContext{logicalName, make(map[string]apiValues.PlcValue)})
	}
	return nil
}

func (p *writeBufferPlcValueBased) GetPos() uint16 {
	return uint16(p.pos / 8)
}

func (p *writeBufferPlcValueBased) WriteBit(logicalName string, value bool, _ ...utils.WithWriterArgs) error {
	p.move(1)
	return p.appendValue(logicalName, NewPlcBOOL(value))
}

func (p *writeBufferPlcValueBased) WriteByte(logicalName string, value byte, _ ...utils.WithWriterArgs) error {
	p.move(8)
	return p.appendValue(logicalName, NewPlcUSINT(value))
}

func (p *writeBufferPlcValueBased) WriteByteArray(logicalName string, data []byte, _ ...utils.WithWriterArgs) error {
	p.move(uint(len(data) * 8))
	return p.appendValue(logicalName, NewPlcRawByteArray(data))
}

func (p *writeBufferPlcValueBased) WriteUint8(logicalName string, bitLength uint8, value uint8, _ ...utils.WithWriterArgs) error {
	p.move(uint(bitLength))
	return p.appendValue(logicalName, NewPlcUSINT(value))
}

func (p *writeBufferPlcValueBased) WriteUint16(logicalName string, bitLength uint8, value uint16, _ ...utils.WithWriterArgs) error {
	p.move(uint(bitLength))
	return p.appendValue(logicalName, NewPlcUINT(value))
}

func (p *writeBufferPlcValueBased) WriteUint32(logicalName string, bitLength uint8, value uint32, _ ...utils.WithWriterArgs) error {
	p.move(uint(bitLength))
	return p.appendValue(logicalName, NewPlcUDINT(value))
}

func (p *writeBufferPlcValueBased) WriteUint64(logicalName string, bitLength uint8, value uint64, _ ...utils.WithWriterArgs) error {
	p.move(uint(bitLength))
	return p.appendValue(logicalName, NewPlcULINT(value))
}

func (p *writeBufferPlcValueBased) WriteInt8(logicalName string, bitLength uint8, value int8, _ ...utils.WithWriterArgs) error {
	p.move(uint(bitLength))
	return p.appendValue(logicalName, NewPlcSINT(value))
}

func (p *writeBufferPlcValueBased) WriteInt16(logicalName string, bitLength uint8, value int16, _ ...utils.WithWriterArgs) error {
	p.move(uint(bitLength))
	return p.appendValue(logicalName, NewPlcINT(value))
}

func (p *writeBufferPlcValueBased) WriteInt32(logicalName string, bitLength uint8, value int32, _ ...utils.WithWriterArgs) error {
	p.move(uint(bitLength))
	return p.appendValue(logicalName, NewPlcDINT(value))
}

func (p *writeBufferPlcValueBased) WriteInt64(logicalName string, bitLength uint8, value int64, _ ...utils.WithWriterArgs) error {
	p.move(uint(bitLength))
	return p.appendValue(logicalName, NewPlcLINT(value))
}

func (p *writeBufferPlcValueBased) WriteBigInt(logicalName string, bitLength uint8, value *big.Int, _ ...utils.WithWriterArgs) error {
	if value == nil {
		return errors.New("value must not be nil")
	}
	p.move(uint(bitLength))
	// TODO: check if we set the type dynamic here...
	return p.appendValue(logicalName, NewPlcRawByteArray(value.Bytes()))
}

func (p *writeBufferPlcValueBased) WriteFloat32(logicalName string, bitLength uint8, value float32, _ ...utils.WithWriterArgs) error {
	p.move(uint(bitLength))
	return p.appendValue(logicalName, NewPlcREAL(value))
}

func (p *writeBufferPlcValueBased) WriteFloat64(logicalName string, bitLength uint8, value float64, _ ...utils.WithWriterArgs) error {
	p.move(uint(bitLength))
	return p.appendValue(logicalName, NewPlcLREAL(value))
}

func (p *writeBufferPlcValueBased) WriteBigFloat(logicalName string, bitLength uint8, value *big.Float, _ ...utils.WithWriterArgs) error {
	if value == nil {
		return errors.New("value must not be nil")
	}
	p.move(uint(bitLength))
	// TODO: check if we set the type dynamic here...
	encode, err := value.GobEncode()
	if err != nil {
		return errors.Wrapf(err, "Error writing %s", logicalName)
	}
	return p.appendValue(logicalName, NewPlcRawByteArray(encode))
}

func (p *writeBufferPlcValueBased) WriteString(logicalName string, bitLength uint32, _ string, value string, _ ...utils.WithWriterArgs) error {
	p.move(uint(bitLength))
	return p.appendValue(logicalName, NewPlcSTRING(value))
}

func (p *writeBufferPlcValueBased) WriteVirtual(ctx context.Context, logicalName string, value any, _ ...utils.WithWriterArgs) error {
	// NO-OP
	return nil
}

func (p *writeBufferPlcValueBased) WriteSerializable(ctx context.Context, serializable utils.Serializable) error {
	if serializable == nil {
		return nil
	}
	return serializable.SerializeWithWriteBuffer(ctx, p)
}

func (p *writeBufferPlcValueBased) PopContext(logicalName string, _ ...utils.WithWriterArgs) error {
	pop := p.Pop()
	var poppedName string
	var unwrapped apiValues.PlcValue
	switch _context := pop.(type) {
	case *plcValueContext:
		poppedName = _context.logicalName
		unwrapped = NewPlcStruct(_context.properties)
	case *plcListContext:
		poppedName = _context.logicalName
		unwrapped = NewPlcList(_context.list)
	default:
		return errors.New("broken context")
	}
	if poppedName != logicalName {
		return errors.Errorf("unexpected closing context %s, expected %s", poppedName, logicalName)
	}
	if p.Empty() {
		p.rootNode = NewPlcStruct(map[string]apiValues.PlcValue{logicalName: unwrapped})
		return nil
	}
	switch _context := p.Peek().(type) {
	case *plcValueContext:
		_context.properties[logicalName] = unwrapped
	case *plcListContext:
		_context.list = append(_context.list, NewPlcStruct(map[string]apiValues.PlcValue{logicalName: unwrapped}))
	default:
		return errors.New("broken context")
	}
	return nil
}

func (p *writeBufferPlcValueBased) GetPlcValue() apiValues.PlcValue {
	return p.rootNode
}

func (p *writeBufferPlcValueBased) appendValue(logicalName string, value apiValues.PlcValue) error {
	logicalName = p.SanitizeLogicalName(logicalName)
	peek := p.Peek()
	switch _context := peek.(type) {
	case *plcValueContext:
		_context.properties[logicalName] = value
		return nil
	case *plcListContext:
		_context.list = append(_context.list, value)
		return nil
	default:
		newContext := &plcValueContext{logicalName, make(map[string]apiValues.PlcValue)}
		newContext.properties[logicalName] = value
		p.Push(newContext)
		return nil
	}
}

func (p *writeBufferPlcValueBased) move(bits uint) {
	p.pos += bits
}
