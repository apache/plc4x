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
	"container/list"
	"math/big"
)

type WriteBufferBoxBased interface {
	WriteBuffer
	GetBox() AsciiBox
}

func NewBoxedWriteBuffer() WriteBufferBoxBased {
	return &boxedWriteBuffer{
		List:  list.New(),
		width: 200,
	}
}

type boxedWriteBuffer struct {
	*list.List
	width int
}

func (b *boxedWriteBuffer) GetBox() AsciiBox {
	back := b.Back()
	return back.Value.(AsciiBox)
}

func (b *boxedWriteBuffer) PushContext(logicalName string) error {
	b.PushBack(make([]AsciiBox, 0))
	return nil
}

func (b *boxedWriteBuffer) WriteBit(logicalName string, value bool) error {
	back := b.Back()
	b.Remove(back)
	asciiBoxes := back.Value.([]AsciiBox)
	asciiBoxes = append(asciiBoxes, BoxAnything(logicalName, value, b.width))
	b.PushBack(asciiBoxes)
	return nil
}

func (b *boxedWriteBuffer) WriteUint8(logicalName string, bitLength uint8, value uint8) error {
	back := b.Back()
	b.Remove(back)
	asciiBoxes := back.Value.([]AsciiBox)
	asciiBoxes = append(asciiBoxes, BoxAnything(logicalName, value, b.width))
	b.PushBack(asciiBoxes)
	return nil
}

func (b *boxedWriteBuffer) WriteUint16(logicalName string, bitLength uint8, value uint16) error {
	back := b.Back()
	b.Remove(back)
	asciiBoxes := back.Value.([]AsciiBox)
	asciiBoxes = append(asciiBoxes, BoxAnything(logicalName, value, b.width))
	b.PushBack(asciiBoxes)
	return nil
}

func (b *boxedWriteBuffer) WriteUint32(logicalName string, bitLength uint8, value uint32) error {
	back := b.Back()
	b.Remove(back)
	asciiBoxes := back.Value.([]AsciiBox)
	asciiBoxes = append(asciiBoxes, BoxAnything(logicalName, value, b.width))
	b.PushBack(asciiBoxes)
	return nil
}

func (b *boxedWriteBuffer) WriteUint64(logicalName string, bitLength uint8, value uint64) error {
	back := b.Back()
	b.Remove(back)
	asciiBoxes := back.Value.([]AsciiBox)
	asciiBoxes = append(asciiBoxes, BoxAnything(logicalName, value, b.width))
	b.PushBack(asciiBoxes)
	return nil
}

func (b *boxedWriteBuffer) WriteInt8(logicalName string, bitLength uint8, value int8) error {
	back := b.Back()
	b.Remove(back)
	asciiBoxes := back.Value.([]AsciiBox)
	asciiBoxes = append(asciiBoxes, BoxAnything(logicalName, value, b.width))
	b.PushBack(asciiBoxes)
	return nil
}

func (b *boxedWriteBuffer) WriteInt16(logicalName string, bitLength uint8, value int16) error {
	back := b.Back()
	b.Remove(back)
	asciiBoxes := back.Value.([]AsciiBox)
	asciiBoxes = append(asciiBoxes, BoxAnything(logicalName, value, b.width))
	b.PushBack(asciiBoxes)
	return nil
}

func (b *boxedWriteBuffer) WriteInt32(logicalName string, bitLength uint8, value int32) error {
	back := b.Back()
	b.Remove(back)
	asciiBoxes := back.Value.([]AsciiBox)
	asciiBoxes = append(asciiBoxes, BoxAnything(logicalName, value, b.width))
	b.PushBack(asciiBoxes)
	return nil
}

func (b *boxedWriteBuffer) WriteInt64(logicalName string, bitLength uint8, value int64) error {
	back := b.Back()
	b.Remove(back)
	asciiBoxes := back.Value.([]AsciiBox)
	asciiBoxes = append(asciiBoxes, BoxAnything(logicalName, value, b.width))
	b.PushBack(asciiBoxes)
	return nil
}

func (b *boxedWriteBuffer) WriteBigInt(logicalName string, bitLength uint8, value *big.Int) error {
	back := b.Back()
	b.Remove(back)
	asciiBoxes := back.Value.([]AsciiBox)
	asciiBoxes = append(asciiBoxes, BoxAnything(logicalName, value, b.width))
	b.PushBack(asciiBoxes)
	return nil
}

func (b *boxedWriteBuffer) WriteFloat32(logicalName string, bitLength uint8, value float32) error {
	back := b.Back()
	b.Remove(back)
	asciiBoxes := back.Value.([]AsciiBox)
	asciiBoxes = append(asciiBoxes, BoxAnything(logicalName, value, b.width))
	b.PushBack(asciiBoxes)
	return nil
}

func (b *boxedWriteBuffer) WriteFloat64(logicalName string, bitLength uint8, value float64) error {
	back := b.Back()
	b.Remove(back)
	asciiBoxes := back.Value.([]AsciiBox)
	asciiBoxes = append(asciiBoxes, BoxAnything(logicalName, value, b.width))
	b.PushBack(asciiBoxes)
	return nil
}

func (b *boxedWriteBuffer) WriteString(logicalName string, bitLength uint8, encoding string, value string) error {
	back := b.Back()
	b.Remove(back)
	asciiBoxes := back.Value.([]AsciiBox)
	asciiBoxes = append(asciiBoxes, BoxAnything(logicalName, value, b.width))
	b.PushBack(asciiBoxes)
	return nil
}

func (b *boxedWriteBuffer) PopContext(logicalName string) error {
	finalBoxes := make([]AsciiBox, 0)
findTheBox:
	for back := b.Back(); back != nil; back = b.Back() {
		switch back.Value.(type) {
		case AsciiBox:
			asciiBox := b.Remove(back).(AsciiBox)
			finalBoxes = append(finalBoxes, asciiBox)
		case []AsciiBox:
			b.Remove(back)
			asciiBoxes := b.Remove(back).([]AsciiBox)
			finalBoxes = append(asciiBoxes, finalBoxes...)
			break findTheBox
		default:
			panic("We should never reach this point")
		}
	}
	asciiBox := BoxBox(logicalName, AlignBoxes(finalBoxes, b.width), 0)
	b.PushBack(asciiBox)
	return nil
}
