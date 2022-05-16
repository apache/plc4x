/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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
	"fmt"
	"math/big"
)

type WriteBufferBoxBased interface {
	WriteBuffer
	GetBox() AsciiBox
}

func NewBoxedWriteBuffer() WriteBufferBoxBased {
	return NewBoxedWriteBufferWithOptions(false, false)
}

func NewBoxedWriteBufferWithOptions(mergeSingleBoxes bool, omitEmptyBoxes bool) WriteBufferBoxBased {
	return &boxedWriteBuffer{
		List:                list.New(),
		desiredWidth:        120,
		currentWidth:        118,
		mergeSingleBoxes:    mergeSingleBoxes,
		omitEmptyBoxes:      omitEmptyBoxes,
		asciiBoxWriter:      AsciiBoxWriterDefault,
		asciiBoxWriterLight: AsciiBoxWriterLight,
	}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type boxedWriteBuffer struct {
	bufferCommons
	*list.List
	desiredWidth        int
	currentWidth        int
	mergeSingleBoxes    bool
	omitEmptyBoxes      bool
	asciiBoxWriter      AsciiBoxWriter
	asciiBoxWriterLight AsciiBoxWriter
	pos                 uint
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (b *boxedWriteBuffer) GetBox() AsciiBox {
	back := b.Back()
	return back.Value.(AsciiBox)
}

func (b *boxedWriteBuffer) PushContext(_ string, _ ...WithWriterArgs) error {
	b.currentWidth -= boxLineOverheat
	b.PushBack(make([]AsciiBox, 0))
	return nil
}

func (b *boxedWriteBuffer) GetPos() uint16 {
	return uint16(b.pos / 8)
}

func (b *boxedWriteBuffer) WriteBit(logicalName string, value bool, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(upcastWriterArgs(writerArgs...)...)
	asInt := 0
	if value {
		asInt = 1
	}
	b.PushBack(b.asciiBoxWriter.BoxString(logicalName, fmt.Sprintf("b%d %t%s", asInt, value, additionalStringRepresentation), 0))
	b.move(1)
	return nil
}

func (b *boxedWriteBuffer) WriteByte(logicalName string, value byte, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(upcastWriterArgs(writerArgs...)...)
	b.PushBack(b.asciiBoxWriter.BoxString(logicalName, fmt.Sprintf("%#02x%s", value, additionalStringRepresentation), 0))
	b.move(8)
	return nil
}

func (b *boxedWriteBuffer) WriteByteArray(logicalName string, data []byte, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(upcastWriterArgs(writerArgs...)...)
	if additionalStringRepresentation != "" {
		additionalStringRepresentation += "\n"
	}
	b.PushBack(b.asciiBoxWriter.BoxString(logicalName, fmt.Sprintf("%s%s", additionalStringRepresentation, Dump(data)), 0))
	b.move(uint(len(data) * 8))
	return nil
}

func (b *boxedWriteBuffer) WriteUint8(logicalName string, bitLength uint8, value uint8, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(upcastWriterArgs(writerArgs...)...)
	b.PushBack(b.asciiBoxWriter.BoxString(logicalName, fmt.Sprintf("%#0*x %d%s", bitLength/4, value, value, additionalStringRepresentation), 0))
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteUint16(logicalName string, bitLength uint8, value uint16, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(upcastWriterArgs(writerArgs...)...)
	b.PushBack(b.asciiBoxWriter.BoxString(logicalName, fmt.Sprintf("%#0*x %d%s", bitLength/4, value, value, additionalStringRepresentation), 0))
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteUint32(logicalName string, bitLength uint8, value uint32, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(upcastWriterArgs(writerArgs...)...)
	b.PushBack(b.asciiBoxWriter.BoxString(logicalName, fmt.Sprintf("%#0*x %d%s", bitLength/4, value, value, additionalStringRepresentation), 0))
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteUint64(logicalName string, bitLength uint8, value uint64, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(upcastWriterArgs(writerArgs...)...)
	b.PushBack(b.asciiBoxWriter.BoxString(logicalName, fmt.Sprintf("%#0*x %d%s", bitLength/4, value, value, additionalStringRepresentation), 0))
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteInt8(logicalName string, bitLength uint8, value int8, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(upcastWriterArgs(writerArgs...)...)
	b.PushBack(b.asciiBoxWriter.BoxString(logicalName, fmt.Sprintf("%#0*x %d%s", bitLength/4, value, value, additionalStringRepresentation), 0))
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteInt16(logicalName string, bitLength uint8, value int16, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(upcastWriterArgs(writerArgs...)...)
	b.PushBack(b.asciiBoxWriter.BoxString(logicalName, fmt.Sprintf("%#0*x %d%s", bitLength/4, value, value, additionalStringRepresentation), 0))
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteInt32(logicalName string, bitLength uint8, value int32, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(upcastWriterArgs(writerArgs...)...)
	b.PushBack(b.asciiBoxWriter.BoxString(logicalName, fmt.Sprintf("%#0*x %d%s", bitLength/4, value, value, additionalStringRepresentation), 0))
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteInt64(logicalName string, bitLength uint8, value int64, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(upcastWriterArgs(writerArgs...)...)
	b.PushBack(b.asciiBoxWriter.BoxString(logicalName, fmt.Sprintf("%#0*x %d%s", bitLength/4, value, value, additionalStringRepresentation), 0))
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteBigInt(logicalName string, bitLength uint8, value *big.Int, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(upcastWriterArgs(writerArgs...)...)
	b.PushBack(b.asciiBoxWriter.BoxString(logicalName, fmt.Sprintf("%#0*x %d%s", bitLength/4, value, value, additionalStringRepresentation), 0))
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteFloat32(logicalName string, bitLength uint8, value float32, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(upcastWriterArgs(writerArgs...)...)
	b.PushBack(b.asciiBoxWriter.BoxString(logicalName, fmt.Sprintf("%#0*x %f%s", bitLength/4, value, value, additionalStringRepresentation), 0))
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteFloat64(logicalName string, bitLength uint8, value float64, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(upcastWriterArgs(writerArgs...)...)
	b.PushBack(b.asciiBoxWriter.BoxString(logicalName, fmt.Sprintf("%#0*x %f%s", bitLength/4, value, value, additionalStringRepresentation), 0))
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteBigFloat(logicalName string, bitLength uint8, value *big.Float, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(upcastWriterArgs(writerArgs...)...)
	b.PushBack(b.asciiBoxWriter.BoxString(logicalName, fmt.Sprintf("%#0*x %f%s", bitLength/4, value, value, additionalStringRepresentation), 0))
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteString(logicalName string, bitLength uint32, _ string, value string, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(upcastWriterArgs(writerArgs...)...)
	b.PushBack(b.asciiBoxWriter.BoxString(logicalName, fmt.Sprintf("%s%s", value, additionalStringRepresentation), 0))
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteVirtual(logicalName string, value interface{}, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(upcastWriterArgs(writerArgs...)...)
	if value == nil {
		return nil
	}
	var asciiBox AsciiBox
	switch value.(type) {
	case bool:
		asciiBox = b.asciiBoxWriterLight.BoxString(logicalName, fmt.Sprintf("%t", value), 0)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		asciiBox = b.asciiBoxWriterLight.BoxString(logicalName, fmt.Sprintf("%#x %d%s", value, value, additionalStringRepresentation), 0)
	case float32, float64:
		asciiBox = b.asciiBoxWriterLight.BoxString(logicalName, fmt.Sprintf("%x %f%s", value, value, additionalStringRepresentation), 0)
	case Serializable:
		virtualBoxedWriteBuffer := NewBoxedWriteBuffer()
		virtualBoxedWriteBuffer.(*boxedWriteBuffer).asciiBoxWriter = AsciiBoxWriterLight
		if err := value.(Serializable).Serialize(virtualBoxedWriteBuffer); err == nil {
			asciiBox = b.asciiBoxWriterLight.BoxBox(logicalName, virtualBoxedWriteBuffer.GetBox(), 0)
		} else {
			b.asciiBoxWriterLight.BoxString(logicalName, err.Error(), 0)
		}
	default:
		asciiBox = b.asciiBoxWriterLight.BoxString(logicalName, fmt.Sprintf("%v%s", value, additionalStringRepresentation), 0)
	}
	b.PushBack(asciiBox)
	return nil
}

func (b *boxedWriteBuffer) PopContext(logicalName string, _ ...WithWriterArgs) error {
	b.currentWidth += boxLineOverheat
	finalBoxes := make([]AsciiBox, 0)
findTheBox:
	for back := b.Back(); back != nil; back = b.Back() {
		switch back.Value.(type) {
		case AsciiBox:
			asciiBox := b.Remove(back).(AsciiBox)
			if b.omitEmptyBoxes && asciiBox.IsEmpty() {
				continue
			}
			finalBoxes = append([]AsciiBox{asciiBox}, finalBoxes...)
		case []AsciiBox:
			b.Remove(back)
			asciiBoxes := b.Remove(back).([]AsciiBox)
			finalBoxes = append(asciiBoxes, finalBoxes...)
			break findTheBox
		default:
			panic("We should never reach this point")
		}
	}
	if b.mergeSingleBoxes && len(finalBoxes) == 1 {
		onlyChild := finalBoxes[0]
		childName := onlyChild.GetBoxName()
		onlyChild = onlyChild.ChangeBoxName(logicalName + "/" + childName)
		if b.omitEmptyBoxes && onlyChild.IsEmpty() {
			return nil
		}
		b.PushBack(onlyChild)
		return nil
	}
	asciiBox := b.asciiBoxWriter.BoxBox(logicalName, b.asciiBoxWriter.AlignBoxes(finalBoxes, b.currentWidth), 0)
	if b.omitEmptyBoxes && asciiBox.IsEmpty() {
		return nil
	}
	b.PushBack(asciiBox)
	return nil
}

func (b *boxedWriteBuffer) extractAdditionalStringRepresentation(readerWriterArgs ...WithReaderWriterArgs) string {
	representation := b.bufferCommons.extractAdditionalStringRepresentation(readerWriterArgs...)
	if representation != "" {
		return " " + representation
	}
	return ""
}

func (b *boxedWriteBuffer) move(bits uint) {
	b.pos += bits
}
