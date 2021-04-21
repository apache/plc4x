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
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"math/big"
	"strings"
)

type WriteBufferJsonBased interface {
	WriteBuffer
	GetJsonString() (string, error)
}

func NewJsonWriteBuffer() WriteBufferJsonBased {
	var jsonString strings.Builder
	encoder := json.NewEncoder(&jsonString)
	encoder.SetIndent("", "  ")
	return &jsonWriteBuffer{
		jsonString:   &jsonString,
		Encoder:      encoder,
		doRenderAttr: true,
	}
}

func NewConfiguredJsonWriteBuffer(renderAttr bool) WriteBufferJsonBased {
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
	bufferCommons
	stack
	*json.Encoder
	jsonString    *strings.Builder
	rootNode      interface{}
	doRenderLists bool
	doRenderAttr  bool
}

type elementContext struct {
	logicalName string
	properties  map[string]interface{}
}

type listContext struct {
	logicalName string
	list        []interface{}
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (j *jsonWriteBuffer) PushContext(logicalName string, writerArgs ...WithWriterArgs) error {
	renderedAsList := j.isToBeRenderedAsList(upcastWriterArgs(writerArgs...)...)
	if renderedAsList {
		j.Push(&listContext{logicalName, make([]interface{}, 0)})
	} else {
		j.Push(&elementContext{logicalName, make(map[string]interface{})})
	}
	return nil
}

func (j *jsonWriteBuffer) WriteBit(logicalName string, value bool, writerArgs ...WithWriterArgs) error {
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwBitKey, 1, writerArgs...))
}

func (j *jsonWriteBuffer) WriteUint8(logicalName string, bitLength uint8, value uint8, writerArgs ...WithWriterArgs) error {
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwUintKey, bitLength, writerArgs...))
}

func (j *jsonWriteBuffer) WriteUint16(logicalName string, bitLength uint8, value uint16, writerArgs ...WithWriterArgs) error {
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwUintKey, bitLength, writerArgs...))
}

func (j *jsonWriteBuffer) WriteUint32(logicalName string, bitLength uint8, value uint32, writerArgs ...WithWriterArgs) error {
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwUintKey, bitLength, writerArgs...))
}

func (j *jsonWriteBuffer) WriteUint64(logicalName string, bitLength uint8, value uint64, writerArgs ...WithWriterArgs) error {
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwUintKey, bitLength, writerArgs...))
}

func (j *jsonWriteBuffer) WriteInt8(logicalName string, bitLength uint8, value int8, writerArgs ...WithWriterArgs) error {
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwIntKey, bitLength, writerArgs...))
}

func (j *jsonWriteBuffer) WriteInt16(logicalName string, bitLength uint8, value int16, writerArgs ...WithWriterArgs) error {
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwIntKey, bitLength, writerArgs...))
}

func (j *jsonWriteBuffer) WriteInt32(logicalName string, bitLength uint8, value int32, writerArgs ...WithWriterArgs) error {
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwIntKey, bitLength, writerArgs...))
}

func (j *jsonWriteBuffer) WriteInt64(logicalName string, bitLength uint8, value int64, writerArgs ...WithWriterArgs) error {
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwIntKey, bitLength, writerArgs...))
}

func (j *jsonWriteBuffer) WriteBigInt(logicalName string, bitLength uint8, value *big.Int, writerArgs ...WithWriterArgs) error {
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwIntKey, bitLength, writerArgs...))
}

func (j *jsonWriteBuffer) WriteFloat32(logicalName string, bitLength uint8, value float32, writerArgs ...WithWriterArgs) error {
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwFloatKey, bitLength, writerArgs...))
}

func (j *jsonWriteBuffer) WriteFloat64(logicalName string, bitLength uint8, value float64, writerArgs ...WithWriterArgs) error {
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwFloatKey, bitLength, writerArgs...))
}

func (j *jsonWriteBuffer) WriteBigFloat(logicalName string, bitLength uint8, value *big.Float, writerArgs ...WithWriterArgs) error {
	return j.encodeNode(logicalName, value, j.generateAttr(logicalName, rwFloatKey, bitLength, writerArgs...))
}

func (j *jsonWriteBuffer) WriteString(logicalName string, bitLength uint8, encoding string, value string, writerArgs ...WithWriterArgs) error {
	attr := j.generateAttr(logicalName, rwStringKey, bitLength, writerArgs...)
	attr[fmt.Sprintf("__plc4x_%s", rwEncodingKey)] = encoding
	return j.encodeNode(logicalName, value, attr)
}

func (j *jsonWriteBuffer) PopContext(logicalName string, _ ...WithWriterArgs) error {
	pop := j.Pop()
	var poppedName string
	var unwrapped interface{}
	switch pop.(type) {
	case *elementContext:
		context := pop.(*elementContext)
		poppedName = context.logicalName
		unwrapped = context.properties
	case *listContext:
		context := pop.(*listContext)
		poppedName = context.logicalName
		unwrapped = context.list
	default:
		panic("broken context")
	}
	if poppedName != logicalName {
		return errors.Errorf("unexpected closing context %s, expected %s", poppedName, logicalName)
	}
	if j.Empty() {
		lastElement := make(map[string]interface{})
		lastElement[logicalName] = unwrapped
		j.rootNode = lastElement
		return nil
	}
	j.rootNode = j.Peek()
	switch j.rootNode.(type) {
	case *elementContext:
		context := j.rootNode.(*elementContext)
		context.properties[logicalName] = unwrapped
	case *listContext:
		context := j.rootNode.(*listContext)
		wrappedWrap := make(map[string]interface{})
		wrappedWrap[logicalName] = unwrapped
		context.list = append(context.list, wrappedWrap)
	default:
		panic("broken context")
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

func (j *jsonWriteBuffer) encodeNode(logicalName string, value interface{}, attr map[string]interface{}, _ ...WithWriterArgs) error {
	logicalName = j.sanitizeLogicalName(logicalName)
	peek := j.Peek()
	switch peek.(type) {
	case *elementContext:
		context := peek.(*elementContext)
		context.properties[logicalName] = value
		for key, attrValue := range attr {
			context.properties[key] = attrValue
		}
		return nil
	case *listContext:
		context := peek.(*listContext)
		m := make(map[string]interface{})
		m[logicalName] = value
		context.list = append(context.list, m)
		return nil
	default:
		panic("broken context")
	}
}

func (j *jsonWriteBuffer) generateAttr(logicalName string, dataType string, bitLength uint8, writerArgs ...WithWriterArgs) map[string]interface{} {
	attr := make(map[string]interface{})
	if !j.doRenderAttr {
		return attr
	}
	logicalName = j.sanitizeLogicalName(logicalName)
	attr[fmt.Sprintf("%s__plc4x_%s", logicalName, rwDataTypeKey)] = dataType
	attr[fmt.Sprintf("%s__plc4x_%s", logicalName, rwBitLengthKey)] = bitLength
	for _, arg := range writerArgs {
		if !arg.isWriterArgs() {
			panic("not a writer arg")
		}
		switch arg.(type) {
		case withAdditionalStringRepresentation:
			attr[fmt.Sprintf("%s__plc4x_%s", logicalName, rwStringRepresentationKey)] = arg.(withAdditionalStringRepresentation).stringRepresentation
		}
	}
	return attr
}
