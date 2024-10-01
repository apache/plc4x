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
	"encoding/binary"
	"fmt"
	"math/big"
	"strconv"

	"github.com/pkg/errors"
)

type WriteBufferBoxBased interface {
	WriteBuffer
	GetBox() AsciiBox
}

func NewWriteBufferBoxBased(opts ...func(buffer *boxedWriteBuffer)) WriteBufferBoxBased {
	wb := &boxedWriteBuffer{
		List:                list.New(),
		desiredWidth:        120,
		currentWidth:        118,
		asciiBoxWriter:      AsciiBoxWriterDefault,
		asciiBoxWriterLight: AsciiBoxWriterLight,
	}
	for _, opt := range opts {
		opt(wb)
	}
	return wb
}

func WithWriteBufferBoxBasedMergeSingleBoxes() func(*boxedWriteBuffer) {
	return func(wb *boxedWriteBuffer) {
		wb.mergeSingleBoxes = true
	}
}

func WithWriteBufferBoxBasedOmitEmptyBoxes() func(*boxedWriteBuffer) {
	return func(wb *boxedWriteBuffer) {
		wb.omitEmptyBoxes = true
	}
}

func WithWriteBufferBoxBasedPrintPosLengthFooter() func(*boxedWriteBuffer) {
	return func(wb *boxedWriteBuffer) {
		wb.printPosLengthFooter = true
	}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type boxedWriteBuffer struct {
	BufferCommons
	*list.List
	desiredWidth         int
	currentWidth         int
	mergeSingleBoxes     bool
	omitEmptyBoxes       bool
	printPosLengthFooter bool
	asciiBoxWriter       AsciiBoxWriter
	asciiBoxWriterLight  AsciiBoxWriter
	pos                  uint
}

var _ WriteBuffer = (*boxedWriteBuffer)(nil)

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (*boxedWriteBuffer) GetByteOrder() binary.ByteOrder {
	return binary.BigEndian
}

func (*boxedWriteBuffer) SetByteOrder(_ binary.ByteOrder) {
}

func (b *boxedWriteBuffer) GetBox() AsciiBox {
	back := b.Back()
	if back == nil {
		return AsciiBox{data: "<nil>"}
	}
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
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(UpcastWriterArgs(writerArgs...)...)
	asInt := 0
	if value {
		asInt = 1
	}
	stringValue := fmt.Sprintf("b%d %t%s", asInt, value, additionalStringRepresentation)
	box := b.asciiBoxWriter.BoxString(stringValue, WithAsciiBoxName(logicalName), WithAsciiBoxFooter(b.getPosFooter(1)))
	b.PushBack(box)
	b.move(1)
	return nil
}

func (b *boxedWriteBuffer) WriteByte(logicalName string, value byte, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(UpcastWriterArgs(writerArgs...)...)
	printSafeChar := value
	if value < 32 || value > 126 {
		printSafeChar = '.'
	}
	stringValue := fmt.Sprintf("%#02x '%c'%s", value, printSafeChar, additionalStringRepresentation)
	box := b.asciiBoxWriter.BoxString(stringValue, WithAsciiBoxName(logicalName), WithAsciiBoxFooter(b.getPosFooter(8)))
	b.PushBack(box)
	b.move(8)
	return nil
}

func (b *boxedWriteBuffer) WriteByteArray(logicalName string, data []byte, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(UpcastWriterArgs(writerArgs...)...)
	if additionalStringRepresentation != "" {
		additionalStringRepresentation += "\n"
	}
	stringValue := fmt.Sprintf("%s%s", additionalStringRepresentation, Dump(data))
	box := b.asciiBoxWriter.BoxString(stringValue, WithAsciiBoxName(logicalName), WithAsciiBoxFooter(b.getPosFooter(len(data)*8)))
	b.PushBack(box)
	b.move(uint(len(data) * 8))
	return nil
}

func (b *boxedWriteBuffer) WriteUint8(logicalName string, bitLength uint8, value uint8, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(UpcastWriterArgs(writerArgs...)...)
	stringValue := fmt.Sprintf("%#0*x %d%s", bitLength/4, value, value, additionalStringRepresentation)
	box := b.asciiBoxWriter.BoxString(stringValue, WithAsciiBoxName(logicalName), WithAsciiBoxFooter(b.getPosFooter(int(bitLength))))
	b.PushBack(box)
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteUint16(logicalName string, bitLength uint8, value uint16, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(UpcastWriterArgs(writerArgs...)...)
	stringValue := fmt.Sprintf("%#0*x %d%s", bitLength/4, value, value, additionalStringRepresentation)
	box := b.asciiBoxWriter.BoxString(stringValue, WithAsciiBoxName(logicalName), WithAsciiBoxFooter(b.getPosFooter(int(bitLength))))
	b.PushBack(box)
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteUint32(logicalName string, bitLength uint8, value uint32, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(UpcastWriterArgs(writerArgs...)...)
	stringValue := fmt.Sprintf("%#0*x %d%s", bitLength/4, value, value, additionalStringRepresentation)
	box := b.asciiBoxWriter.BoxString(stringValue, WithAsciiBoxName(logicalName), WithAsciiBoxFooter(b.getPosFooter(int(bitLength))))
	b.PushBack(box)
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteUint64(logicalName string, bitLength uint8, value uint64, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(UpcastWriterArgs(writerArgs...)...)
	stringValue := fmt.Sprintf("%#0*x %d%s", bitLength/4, value, value, additionalStringRepresentation)
	box := b.asciiBoxWriter.BoxString(stringValue, WithAsciiBoxName(logicalName), WithAsciiBoxFooter(b.getPosFooter(int(bitLength))))
	b.PushBack(box)
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteInt8(logicalName string, bitLength uint8, value int8, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(UpcastWriterArgs(writerArgs...)...)
	stringValue := fmt.Sprintf("%#0*x %d%s", bitLength/4, value, value, additionalStringRepresentation)
	box := b.asciiBoxWriter.BoxString(stringValue, WithAsciiBoxName(logicalName), WithAsciiBoxFooter(b.getPosFooter(int(bitLength))))
	b.PushBack(box)
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteInt16(logicalName string, bitLength uint8, value int16, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(UpcastWriterArgs(writerArgs...)...)
	stringValue := fmt.Sprintf("%#0*x %d%s", bitLength/4, value, value, additionalStringRepresentation)
	box := b.asciiBoxWriter.BoxString(stringValue, WithAsciiBoxName(logicalName), WithAsciiBoxFooter(b.getPosFooter(int(bitLength))))
	b.PushBack(box)
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteInt32(logicalName string, bitLength uint8, value int32, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(UpcastWriterArgs(writerArgs...)...)
	stringValue := fmt.Sprintf("%#0*x %d%s", bitLength/4, value, value, additionalStringRepresentation)
	box := b.asciiBoxWriter.BoxString(stringValue, WithAsciiBoxName(logicalName), WithAsciiBoxFooter(b.getPosFooter(int(bitLength))))
	b.PushBack(box)
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteInt64(logicalName string, bitLength uint8, value int64, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(UpcastWriterArgs(writerArgs...)...)
	stringValue := fmt.Sprintf("%#0*x %d%s", bitLength/4, value, value, additionalStringRepresentation)
	box := b.asciiBoxWriter.BoxString(stringValue, WithAsciiBoxName(logicalName), WithAsciiBoxFooter(b.getPosFooter(int(bitLength))))
	b.PushBack(box)
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteBigInt(logicalName string, bitLength uint8, value *big.Int, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(UpcastWriterArgs(writerArgs...)...)
	stringValue := fmt.Sprintf("%#0*x %d%s", bitLength/4, value, value, additionalStringRepresentation)
	box := b.asciiBoxWriter.BoxString(stringValue, WithAsciiBoxName(logicalName), WithAsciiBoxFooter(b.getPosFooter(int(bitLength))))
	b.PushBack(box)
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteFloat32(logicalName string, bitLength uint8, value float32, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(UpcastWriterArgs(writerArgs...)...)
	stringValue := fmt.Sprintf("%#0*x %f%s", bitLength/4, value, value, additionalStringRepresentation)
	box := b.asciiBoxWriter.BoxString(stringValue, WithAsciiBoxName(logicalName), WithAsciiBoxFooter(b.getPosFooter(int(bitLength))))
	b.PushBack(box)
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteFloat64(logicalName string, bitLength uint8, value float64, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(UpcastWriterArgs(writerArgs...)...)
	stringValue := fmt.Sprintf("%#0*x %f%s", bitLength/4, value, value, additionalStringRepresentation)
	box := b.asciiBoxWriter.BoxString(stringValue, WithAsciiBoxName(logicalName), WithAsciiBoxFooter(b.getPosFooter(int(bitLength))))
	b.PushBack(box)
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteBigFloat(logicalName string, bitLength uint8, value *big.Float, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(UpcastWriterArgs(writerArgs...)...)
	stringValue := fmt.Sprintf("%#0*x %f%s", bitLength/4, value, value, additionalStringRepresentation)
	box := b.asciiBoxWriter.BoxString(stringValue, WithAsciiBoxName(logicalName), WithAsciiBoxFooter(b.getPosFooter(int(bitLength))))
	b.PushBack(box)
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteString(logicalName string, bitLength uint32, value string, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(UpcastWriterArgs(writerArgs...)...)
	stringValue := fmt.Sprintf("%s%s", value, additionalStringRepresentation)
	box := b.asciiBoxWriter.BoxString(stringValue, WithAsciiBoxName(logicalName), WithAsciiBoxFooter(b.getPosFooter(int(bitLength))))
	b.PushBack(box)
	b.move(uint(bitLength))
	return nil
}

func (b *boxedWriteBuffer) WriteVirtual(ctx context.Context, logicalName string, value any, writerArgs ...WithWriterArgs) error {
	additionalStringRepresentation := b.extractAdditionalStringRepresentation(UpcastWriterArgs(writerArgs...)...)
	if value == nil {
		return nil
	}
	var asciiBox AsciiBox
	switch value.(type) {
	case bool:
		stringValue := fmt.Sprintf("%t", value)
		asciiBox = b.asciiBoxWriterLight.BoxString(stringValue, WithAsciiBoxName(logicalName))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		stringValue := fmt.Sprintf("%#x %d%s", value, value, additionalStringRepresentation)
		asciiBox = b.asciiBoxWriterLight.BoxString(stringValue, WithAsciiBoxName(logicalName))
	case float32, float64:
		stringValue := fmt.Sprintf("%x %f%s", value, value, additionalStringRepresentation)
		asciiBox = b.asciiBoxWriterLight.BoxString(stringValue, WithAsciiBoxName(logicalName))
	case Serializable:
		virtualBoxedWriteBuffer := NewWriteBufferBoxBased().(*boxedWriteBuffer)
		virtualBoxedWriteBuffer.mergeSingleBoxes = b.mergeSingleBoxes
		virtualBoxedWriteBuffer.omitEmptyBoxes = b.omitEmptyBoxes
		virtualBoxedWriteBuffer.printPosLengthFooter = b.printPosLengthFooter
		if err := value.(Serializable).SerializeWithWriteBuffer(ctx, virtualBoxedWriteBuffer); err == nil {
			asciiBox = b.asciiBoxWriterLight.BoxBox(virtualBoxedWriteBuffer.GetBox(), WithAsciiBoxName(logicalName))
		} else {
			b.asciiBoxWriterLight.BoxString(err.Error(), WithAsciiBoxName(logicalName))
		}
	default:
		asciiBox = b.asciiBoxWriterLight.BoxString(fmt.Sprintf("%v%s", value, additionalStringRepresentation), WithAsciiBoxName(logicalName))
	}
	b.PushBack(asciiBox)
	return nil
}

func (b *boxedWriteBuffer) WriteSerializable(ctx context.Context, serializable Serializable) error {
	if serializable == nil {
		return nil
	}
	currentPos := int(b.pos) // used for footer so we remember that before we advance
	if err := serializable.SerializeWithWriteBuffer(ctx, b); err != nil {
		return err
	}
	back := b.Back()
	if back == nil {
		return nil
	}
	if ab, ok := back.Value.(AsciiBox); ok {
		if la, ok := serializable.(LengthAware); ok {
			bitLength := int(la.GetLengthInBits(ctx))
			back.Value = ab.ChangeBoxFooter(b.getPosFooterWithCurrentPost(currentPos, bitLength))
		}
	}
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
			return errors.New("We should never reach this point")
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
	asciiBox := b.asciiBoxWriter.BoxBox(b.asciiBoxWriter.AlignBoxes(finalBoxes, b.currentWidth), WithAsciiBoxName(logicalName))
	if b.omitEmptyBoxes && asciiBox.IsEmpty() {
		return nil
	}
	b.PushBack(asciiBox)
	return nil
}

func (b *boxedWriteBuffer) extractAdditionalStringRepresentation(readerWriterArgs ...WithReaderWriterArgs) string {
	representation := b.BufferCommons.ExtractAdditionalStringRepresentation(readerWriterArgs...)
	if representation != "" {
		return " " + representation
	}
	return ""
}

func (b *boxedWriteBuffer) move(bits uint) {
	b.pos += bits
}

func (b *boxedWriteBuffer) getPosFooter(bitLength int) string {
	return b.getPosFooterWithCurrentPost(int(b.pos), bitLength)
}

func (b *boxedWriteBuffer) getPosFooterWithCurrentPost(currentPos, bitLength int) string {
	if !b.printPosLengthFooter {
		return ""
	}
	bytePos := currentPos / 8
	bitRemainder := currentPos % 8
	pos := strconv.Itoa(bytePos)
	if bitRemainder != 0 {
		pos += "." + strconv.Itoa(bitRemainder)
	}
	byteLength := bitLength / 8
	bitLengthRemainder := bitLength % 8
	length := strconv.Itoa(byteLength)
	if bitLengthRemainder != 0 {
		length += "." + strconv.Itoa(bitLengthRemainder)
	}
	return fmt.Sprintf("%s/%s", pos, length)
}
