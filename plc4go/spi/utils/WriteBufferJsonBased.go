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
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/pkg/errors"
)

type WriteBufferJsonBased interface {
	WriteBuffer
	GetJsonString() (string, error)
}

func NewJsonWriteBuffer() WriteBufferJsonBased {
	return NewJsonWriteBufferWithOptions(true)
}

func NewJsonWriteBufferWithOptions(renderAttr bool) WriteBufferJsonBased {
	var jsonString strings.Builder
	encoder := json.NewEncoder(&jsonString)
	encoder.SetIndent("", "  ")
	return &jsonWriteBuffer{
		jsonString:   &jsonString,
		Encoder:      encoder,
		doRenderAttr: renderAttr,
	}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type jsonWriteBuffer struct {
	BufferCommons
	Stack
	*json.Encoder
	jsonString   *strings.Builder
	rootNode     any
	doRenderAttr bool
	pos          uint
}

type elementContext struct {
	logicalName string
	properties  map[string]any
}

type listContext struct {
	logicalName string
	list        []any
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (j *jsonWriteBuffer) PushContext(logicalName string, writerArgs ...WithWriterArgs) error {
	renderedAsList := j.IsToBeRenderedAsList(UpcastWriterArgs(writerArgs...)...)
	if renderedAsList {
		j.Push(&listContext{logicalName, make([]any, 0)})
	} else {
		j.Push(&elementContext{logicalName, make(map[string]any)})
	}
	return nil
}

func (j *jsonWriteBuffer) GetPos() uint16 {
	return uint16(j.pos / 8)
}

func (j *jsonWriteBuffer) WriteBit(logicalName string, value bool, writerArgs ...WithWriterArgs) error {
	j.move(1)
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwBitKey, 1, writerArgs...))
}

func (j *jsonWriteBuffer) WriteByte(logicalName string, value byte, writerArgs ...WithWriterArgs) error {
	j.move(8)
	return j.encodeNode(logicalName, fmt.Sprintf("%#02x", value), j.generateAttr(logicalName, rwByteKey, 8, writerArgs...))
}

func (j *jsonWriteBuffer) WriteByteArray(logicalName string, data []byte, writerArgs ...WithWriterArgs) error {
	hexString := fmt.Sprintf("%#x", data)
	if hexString == "00" {
		// golang does mess up the formatting on empty arrays
		hexString = "0x"
	}
	j.move(uint(len(data) * 8))
	return j.encodeNode(logicalName, hexString, j.generateAttr(logicalName, rwByteKey, uint(len(data)*8), writerArgs...))
}

func (j *jsonWriteBuffer) WriteUint8(logicalName string, bitLength uint8, value uint8, writerArgs ...WithWriterArgs) error {
	j.move(uint(bitLength))
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwUintKey, uint(bitLength), writerArgs...))
}

func (j *jsonWriteBuffer) WriteUint16(logicalName string, bitLength uint8, value uint16, writerArgs ...WithWriterArgs) error {
	j.move(uint(bitLength))
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwUintKey, uint(bitLength), writerArgs...))
}

func (j *jsonWriteBuffer) WriteUint32(logicalName string, bitLength uint8, value uint32, writerArgs ...WithWriterArgs) error {
	j.move(uint(bitLength))
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwUintKey, uint(bitLength), writerArgs...))
}

func (j *jsonWriteBuffer) WriteUint64(logicalName string, bitLength uint8, value uint64, writerArgs ...WithWriterArgs) error {
	j.move(uint(bitLength))
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwUintKey, uint(bitLength), writerArgs...))
}

func (j *jsonWriteBuffer) WriteInt8(logicalName string, bitLength uint8, value int8, writerArgs ...WithWriterArgs) error {
	j.move(uint(bitLength))
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwIntKey, uint(bitLength), writerArgs...))
}

func (j *jsonWriteBuffer) WriteInt16(logicalName string, bitLength uint8, value int16, writerArgs ...WithWriterArgs) error {
	j.move(uint(bitLength))
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwIntKey, uint(bitLength), writerArgs...))
}

func (j *jsonWriteBuffer) WriteInt32(logicalName string, bitLength uint8, value int32, writerArgs ...WithWriterArgs) error {
	j.move(uint(bitLength))
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwIntKey, uint(bitLength), writerArgs...))
}

func (j *jsonWriteBuffer) WriteInt64(logicalName string, bitLength uint8, value int64, writerArgs ...WithWriterArgs) error {
	j.move(uint(bitLength))
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwIntKey, uint(bitLength), writerArgs...))
}

func (j *jsonWriteBuffer) WriteBigInt(logicalName string, bitLength uint8, value *big.Int, writerArgs ...WithWriterArgs) error {
	j.move(uint(bitLength))
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwIntKey, uint(bitLength), writerArgs...))
}

func (j *jsonWriteBuffer) WriteFloat32(logicalName string, bitLength uint8, value float32, writerArgs ...WithWriterArgs) error {
	j.move(uint(bitLength))
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwFloatKey, uint(bitLength), writerArgs...))
}

func (j *jsonWriteBuffer) WriteFloat64(logicalName string, bitLength uint8, value float64, writerArgs ...WithWriterArgs) error {
	j.move(uint(bitLength))
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwFloatKey, uint(bitLength), writerArgs...))
}

func (j *jsonWriteBuffer) WriteBigFloat(logicalName string, bitLength uint8, value *big.Float, writerArgs ...WithWriterArgs) error {
	j.move(uint(bitLength))
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwFloatKey, uint(bitLength), writerArgs...))
}

func (j *jsonWriteBuffer) WriteString(logicalName string, bitLength uint32, encoding string, value string, writerArgs ...WithWriterArgs) error {
	attr := j.generateAttr(logicalName, rwStringKey, uint(bitLength), writerArgs...)
	attr[fmt.Sprintf("%s__plc4x_%s", logicalName, rwEncodingKey)] = encoding
	j.move(uint(bitLength))
	return j.encodeNode(logicalName, value, attr)
}

func (j *jsonWriteBuffer) WriteVirtual(ctx context.Context, logicalName string, value any, writerArgs ...WithWriterArgs) error {
	// NO-OP
	return nil
}

func (j *jsonWriteBuffer) WriteSerializable(ctx context.Context, serializable Serializable) error {
	if serializable == nil {
		return nil
	}
	return serializable.SerializeWithWriteBuffer(ctx, j)
}

func (j *jsonWriteBuffer) PopContext(logicalName string, _ ...WithWriterArgs) error {
	pop := j.Pop()
	var poppedName string
	var unwrapped any
	switch _context := pop.(type) {
	case *elementContext:
		poppedName = _context.logicalName
		unwrapped = _context.properties
	case *listContext:
		poppedName = _context.logicalName
		unwrapped = _context.list
	default:
		return errors.New("broken context")
	}
	if poppedName != logicalName {
		return errors.Errorf("unexpected closing context %s, expected %s", poppedName, logicalName)
	}
	if j.Empty() {
		lastElement := make(map[string]any)
		lastElement[logicalName] = unwrapped
		j.rootNode = lastElement
		return nil
	}
	j.rootNode = j.Peek()
	switch _context := j.rootNode.(type) {
	case *elementContext:
		_context.properties[logicalName] = unwrapped
	case *listContext:
		wrappedWrap := make(map[string]any)
		wrappedWrap[logicalName] = unwrapped
		_context.list = append(_context.list, wrappedWrap)
	default:
		return errors.New("broken context")
	}
	return nil
}

func (j *jsonWriteBuffer) GetJsonString() (string, error) {
	if j.rootNode == nil {
		return "", errors.New("No content available")
	}
	err := j.Encode(j.rootNode)
	if err != nil {
		return "", err
	}
	return j.jsonString.String(), nil
}

func (j *jsonWriteBuffer) encodeNode(logicalName string, value any, attr map[string]any, _ ...WithWriterArgs) error {
	logicalName = j.SanitizeLogicalName(logicalName)
	peek := j.Peek()
	switch _context := peek.(type) {
	case *elementContext:
		_context.properties[logicalName] = value
		for key, attrValue := range attr {
			_context.properties[key] = attrValue
		}
		return nil
	case *listContext:
		m := make(map[string]any)
		m[logicalName] = value
		for attrKey, attrValue := range attr {
			m[attrKey] = attrValue
		}
		_context.list = append(_context.list, m)
		return nil
	default:
		newContext := &elementContext{logicalName, make(map[string]any)}
		newContext.properties[logicalName] = value
		for key, attrValue := range attr {
			newContext.properties[key] = attrValue
		}
		j.Push(newContext)
		return nil
	}
}

func (j *jsonWriteBuffer) generateAttr(logicalName string, dataType string, bitLength uint, writerArgs ...WithWriterArgs) map[string]any {
	attr := make(map[string]any)
	if !j.doRenderAttr {
		return attr
	}
	logicalName = j.SanitizeLogicalName(logicalName)
	attr[fmt.Sprintf("%s__plc4x_%s", logicalName, rwDataTypeKey)] = dataType
	attr[fmt.Sprintf("%s__plc4x_%s", logicalName, rwBitLengthKey)] = bitLength
	additionalStringRepresentation := j.ExtractAdditionalStringRepresentation(UpcastWriterArgs(writerArgs...)...)
	if additionalStringRepresentation != "" {
		attr[fmt.Sprintf("%s__plc4x_%s", logicalName, rwStringRepresentationKey)] = additionalStringRepresentation
	}
	return attr
}

func (j *jsonWriteBuffer) move(bits uint) {
	j.pos += bits
}
