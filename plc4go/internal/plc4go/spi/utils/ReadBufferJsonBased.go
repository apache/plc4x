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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"math/big"
	"strings"
)

// NewJsonReadBuffer return as ReadBuffer which doesn't validate attributes and lists
func NewJsonReadBuffer(reader io.Reader) ReadBuffer {
	decoder := json.NewDecoder(reader)
	var rootElement map[string]interface{}
	err := decoder.Decode(&rootElement)
	return &jsonReadBuffer{
		rootElement:    rootElement,
		pos:            1,
		doValidateAttr: true,
		err:            err,
	}
}

// NewStrictJsonReadBuffer return as ReadBuffer which does validate attributes on the setting
func NewStrictJsonReadBuffer(reader io.Reader, validateAttr bool) ReadBuffer {
	decoder := json.NewDecoder(reader)
	var rootElement map[string]interface{}
	err := decoder.Decode(&rootElement)
	return &jsonReadBuffer{
		rootElement:    rootElement,
		pos:            1,
		doValidateAttr: validateAttr,
		err:            err,
	}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type jsonReadBuffer struct {
	bufferCommons
	stack
	rootElement    map[string]interface{}
	pos            uint
	doValidateAttr bool
	err            error
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (j *jsonReadBuffer) GetPos() uint16 {
	return uint16(j.pos / 8)
}

func (j *jsonReadBuffer) Reset(pos uint16) {
	j.pos = uint(pos * 8)
}

func (j *jsonReadBuffer) HasMore(bitLength uint8) bool {
	// TODO: work with x.InputOffset() and check if we are at EOF
	return true
}

func (j *jsonReadBuffer) PullContext(logicalName string, readerArgs ...WithReaderArgs) error {
	if j.err != nil {
		return j.err
	}
	logicalName = j.sanitizeLogicalName(logicalName)
	if j.Empty() {
		if context, ok := j.rootElement[logicalName]; ok {
			j.Push(context)
			return nil
		} else {
			return errors.Errorf("Required context %s not found in %v", logicalName, j.rootElement)
		}
	}
	peek := j.Peek()
	switch peek.(type) {
	case []interface{}:
		pop := j.Pop()
		contextList := pop.([]interface{})
		context := contextList[0].(map[string]interface{})
		if len(contextList) < 2 {
			j.Push(make([]interface{}, 0))
		} else {
			j.Push(contextList[1 : len(contextList)-1])
		}
		if subContext, ok := context[logicalName]; ok {
			j.Push(subContext)
			return nil
		} else {
			return errors.Errorf("Required context %s not found in %v", logicalName, peek)
		}
	}
	if context, ok := peek.(map[string]interface{})[logicalName]; ok {
		j.Push(context)
		return nil
	} else {
		return errors.Errorf("Required context %s not found in %v", logicalName, peek)
	}
}

func (j *jsonReadBuffer) ReadBit(logicalName string, readerArgs ...WithReaderArgs) (bool, error) {
	if j.err != nil {
		return false, j.err
	}
	logicalName = j.sanitizeLogicalName(logicalName)
	j.move(1)
	peek, element := j.getElement(logicalName)
	if err := j.validateAttr(logicalName, element, "bit", 1, readerArgs...); err != nil {
		return false, err
	}
	if value, ok := element[logicalName]; ok {
		return value.(bool), nil
	} else {
		return false, errors.Errorf("Required element %s not found in %v", logicalName, peek)
	}
}

func (j *jsonReadBuffer) ReadByte(logicalName string, readerArgs ...WithReaderArgs) (byte, error) {
	if j.err != nil {
		return 0, j.err
	}
	logicalName = j.sanitizeLogicalName(logicalName)
	j.move(8)
	peek, element := j.getElement(logicalName)
	if err := j.validateAttr(logicalName, element, "byte", 8, readerArgs...); err != nil {
		return 0, err
	}
	if value, ok := element[logicalName]; ok {
		hexString := value.(string)
		if !strings.HasPrefix(hexString, "0x") {
			return 0, errors.Errorf("Hex string should start with 0x. Actual value %s", hexString)
		}
		hexString = strings.Replace(hexString, "0x", "", 1)
		decoded, err := hex.DecodeString(hexString)
		if err != nil {
			return 0, err
		}
		return decoded[0], nil
	} else {
		return 0, errors.Errorf("Required element %s not found in %v", logicalName, peek)
	}
}

func (j *jsonReadBuffer) ReadByteArray(logicalName string, numberOfBytes int, readerArgs ...WithReaderArgs) ([]byte, error) {
	if j.err != nil {
		return nil, j.err
	}
	logicalName = j.sanitizeLogicalName(logicalName)
	j.move(uint(numberOfBytes / 8))
	peek, element := j.getElement(logicalName)
	if err := j.validateAttr(logicalName, element, "byte", uint(numberOfBytes*8), readerArgs...); err != nil {
		return nil, err
	}
	if value, ok := element[logicalName]; ok {
		hexString := value.(string)
		if !strings.HasPrefix(hexString, "0x") {
			return nil, errors.Errorf("Hex string should start with 0x. Actual value %s", hexString)
		}
		hexString = strings.Replace(hexString, "0x", "", 1)
		decoded, err := hex.DecodeString(hexString)
		if err != nil {
			return nil, err
		}
		return decoded, nil
	} else {
		return nil, errors.Errorf("Required element %s not found in %v", logicalName, peek)
	}
}

func (j *jsonReadBuffer) ReadUint8(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (uint8, error) {
	if j.err != nil {
		return 0, j.err
	}
	logicalName = j.sanitizeLogicalName(logicalName)
	j.move(uint(bitLength))
	peek, element := j.getElement(logicalName)
	if err := j.validateAttr(logicalName, element, "uint", uint(bitLength), readerArgs...); err != nil {
		return 0, err
	}
	if value, ok := element[logicalName]; ok {
		return uint8(value.(float64)), nil
	} else {
		return 0, errors.Errorf("Required element %s not found in %v", logicalName, peek)
	}
}

func (j *jsonReadBuffer) ReadUint16(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (uint16, error) {
	if j.err != nil {
		return 0, j.err
	}
	logicalName = j.sanitizeLogicalName(logicalName)
	j.move(uint(bitLength))
	peek, element := j.getElement(logicalName)
	if err := j.validateAttr(logicalName, element, "uint", uint(bitLength), readerArgs...); err != nil {
		return 0, err
	}
	if value, ok := element[logicalName]; ok {
		return uint16(value.(float64)), nil
	} else {
		return 0, errors.Errorf("Required element %s not found in %v", logicalName, peek)
	}
}

func (j *jsonReadBuffer) ReadUint32(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (uint32, error) {
	if j.err != nil {
		return 0, j.err
	}
	logicalName = j.sanitizeLogicalName(logicalName)
	j.move(uint(bitLength))
	peek, element := j.getElement(logicalName)
	if err := j.validateAttr(logicalName, element, "uint", uint(bitLength), readerArgs...); err != nil {
		return 0, err
	}
	if value, ok := element[logicalName]; ok {
		return uint32(value.(float64)), nil
	} else {
		return 0, errors.Errorf("Required element %s not found in %v", logicalName, peek)
	}
}

func (j *jsonReadBuffer) ReadUint64(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (uint64, error) {
	if j.err != nil {
		return 0, j.err
	}
	logicalName = j.sanitizeLogicalName(logicalName)
	j.move(uint(bitLength))
	peek, element := j.getElement(logicalName)
	if err := j.validateAttr(logicalName, element, "uint", uint(bitLength), readerArgs...); err != nil {
		return 0, err
	}
	if value, ok := element[logicalName]; ok {
		return uint64(value.(float64)), nil
	} else {
		return 0, errors.Errorf("Required element %s not found in %v", logicalName, peek)
	}
}

func (j *jsonReadBuffer) ReadInt8(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (int8, error) {
	if j.err != nil {
		return 0, j.err
	}
	logicalName = j.sanitizeLogicalName(logicalName)
	j.move(uint(bitLength))
	peek, element := j.getElement(logicalName)
	if err := j.validateAttr(logicalName, element, "int", uint(bitLength), readerArgs...); err != nil {
		return 0, err
	}
	if value, ok := element[logicalName]; ok {
		return int8(value.(float64)), nil
	} else {
		return 0, errors.Errorf("Required element %s not found in %v", logicalName, peek)
	}
}

func (j *jsonReadBuffer) ReadInt16(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (int16, error) {
	if j.err != nil {
		return 0, j.err
	}
	logicalName = j.sanitizeLogicalName(logicalName)
	j.move(uint(bitLength))
	peek, element := j.getElement(logicalName)
	if err := j.validateAttr(logicalName, element, "int", uint(bitLength), readerArgs...); err != nil {
		return 0, err
	}
	if value, ok := element[logicalName]; ok {
		return int16(value.(float64)), nil
	} else {
		return 0, errors.Errorf("Required element %s not found in %v", logicalName, peek)
	}
}

func (j *jsonReadBuffer) ReadInt32(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (int32, error) {
	if j.err != nil {
		return 0, j.err
	}
	logicalName = j.sanitizeLogicalName(logicalName)
	j.move(uint(bitLength))
	peek, element := j.getElement(logicalName)
	if err := j.validateAttr(logicalName, element, "int", uint(bitLength), readerArgs...); err != nil {
		return 0, err
	}
	if value, ok := element[logicalName]; ok {
		return int32(value.(float64)), nil
	} else {
		return 0, errors.Errorf("Required element %s not found in %v", logicalName, peek)
	}
}

func (j *jsonReadBuffer) ReadInt64(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (int64, error) {
	if j.err != nil {
		return 0, j.err
	}
	logicalName = j.sanitizeLogicalName(logicalName)
	j.move(uint(bitLength))
	peek, element := j.getElement(logicalName)
	if err := j.validateAttr(logicalName, element, "int", uint(bitLength), readerArgs...); err != nil {
		return 0, err
	}
	if value, ok := element[logicalName]; ok {
		return int64(value.(float64)), nil
	} else {
		return 0, errors.Errorf("Required element %s not found in %v", logicalName, peek)
	}
}

func (j *jsonReadBuffer) ReadBigInt(logicalName string, bitLength uint64, readerArgs ...WithReaderArgs) (*big.Int, error) {
	if j.err != nil {
		return nil, j.err
	}
	logicalName = j.sanitizeLogicalName(logicalName)
	logicalName = j.sanitizeLogicalName(logicalName)
	peek, element := j.getElement(logicalName)
	// TODO: not enough bits
	if err := j.validateAttr(logicalName, element, "int", uint(bitLength), readerArgs...); err != nil {
		return nil, err
	}
	newInt := big.NewInt(0)
	if value, ok := element[logicalName]; ok {
		newInt.SetString(value.(string), 10)
		return newInt, nil
	} else {
		return newInt, errors.Errorf("Required element %s not found in %v", logicalName, peek)
	}
}

func (j *jsonReadBuffer) ReadFloat32(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (float32, error) {
	if j.err != nil {
		return 0, j.err
	}
	logicalName = j.sanitizeLogicalName(logicalName)
	j.move(uint(bitLength))
	peek, element := j.getElement(logicalName)
	if err := j.validateAttr(logicalName, element, "float", uint(bitLength), readerArgs...); err != nil {
		return 0, err
	}
	if value, ok := element[logicalName]; ok {
		return float32(value.(float64)), nil
	} else {
		return 0, errors.Errorf("Required element %s not found in %v", logicalName, peek)
	}
}

func (j *jsonReadBuffer) ReadFloat64(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (float64, error) {
	if j.err != nil {
		return 0, j.err
	}
	logicalName = j.sanitizeLogicalName(logicalName)
	j.move(uint(bitLength))
	peek, element := j.getElement(logicalName)
	if err := j.validateAttr(logicalName, element, "float", uint(bitLength), readerArgs...); err != nil {
		return 0, err
	}
	if value, ok := element[logicalName]; ok {
		return float64(value.(float64)), nil
	} else {
		return 0, errors.Errorf("Required element %s not found in %v", logicalName, peek)
	}
}

func (j *jsonReadBuffer) ReadBigFloat(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (*big.Float, error) {
	if j.err != nil {
		return nil, j.err
	}
	logicalName = j.sanitizeLogicalName(logicalName)
	j.move(uint(bitLength))
	peek, element := j.getElement(logicalName)
	if err := j.validateAttr(logicalName, element, "float", uint(bitLength), readerArgs...); err != nil {
		return nil, err
	}
	newFloat := big.NewFloat(0)
	if value, ok := element[logicalName]; ok {
		newFloat.SetString(value.(string))
		return newFloat, nil
	} else {
		return newFloat, errors.Errorf("Required element %s not found in %v", logicalName, peek)
	}
}

func (j *jsonReadBuffer) ReadString(logicalName string, bitLength uint32, readerArgs ...WithReaderArgs) (string, error) {
	if j.err != nil {
		return "", j.err
	}
	logicalName = j.sanitizeLogicalName(logicalName)
	j.move(uint(bitLength))
	peek, element := j.getElement(logicalName)
	if err := j.validateAttr(logicalName, element, "string", uint(bitLength), readerArgs...); err != nil {
		return "", err
	}
	if value, ok := element[logicalName]; ok {
		return value.(string), nil
	} else {
		return "", errors.Errorf("Required element %s not found in %v", logicalName, peek)
	}
}

func (j *jsonReadBuffer) CloseContext(logicalName string, readerArgs ...WithReaderArgs) error {
	if j.err != nil {
		return j.err
	}
	logicalName = j.sanitizeLogicalName(logicalName)
	if j.Empty() {
		return errors.Errorf("Required context close %s not found in %v", logicalName, j.rootElement)
	}
	// Delete us from stack
	_ = j.Pop()
	if j.Empty() {
		return nil
	}
	peek := j.Peek()
	switch peek.(type) {
	case []interface{}:
		return nil
	}
	if _, ok := peek.(map[string]interface{})[logicalName]; ok {
		delete(peek.(map[string]interface{}), logicalName)
		return nil
	} else {
		return errors.Errorf("Required context %s not found in %v", logicalName, peek)
	}
}

func (j *jsonReadBuffer) getElement(logicalName string) (interface{}, map[string]interface{}) {
	logicalName = j.sanitizeLogicalName(logicalName)
	peek := j.Peek()
	var element map[string]interface{}
	switch peek.(type) {
	case []interface{}:
		pop := j.Pop()
		elementList := pop.([]interface{})
		element = elementList[0].(map[string]interface{})
		if len(elementList) < 2 {
			j.Push(make([]interface{}, 0))
		} else {
			j.Push(elementList[1 : len(elementList)-1])
		}
	case map[string]interface{}:
		element = peek.(map[string]interface{})
	default:
		panic(fmt.Sprintf("Invalid state at %s with %v", logicalName, element))
	}
	return peek, element
}

func (j *jsonReadBuffer) move(bits uint) {
	j.pos += bits
}

func (j *jsonReadBuffer) validateAttr(logicalName string, element map[string]interface{}, dataType string, bitLength uint, readerArgs ...WithReaderArgs) error {
	if !j.doValidateAttr {
		return nil
	}
	renderedKeyDataLengthKey := fmt.Sprintf("%s__plc4x_%s", logicalName, rwDataTypeKey)
	if value, ok := element[renderedKeyDataLengthKey]; !ok || value != dataType {
		return errors.Errorf("Unexpected %s :%s. Want %s", renderedKeyDataLengthKey, value, dataType)
	}
	renderedBitLengthKey := fmt.Sprintf("%s__plc4x_%s", logicalName, rwBitLengthKey)
	if value, ok := element[renderedBitLengthKey]; !ok || uint(value.(float64)) != bitLength {
		return errors.Errorf("Unexpected %s :%d. Want %d", renderedBitLengthKey, value, bitLength)
	}
	return nil
}
