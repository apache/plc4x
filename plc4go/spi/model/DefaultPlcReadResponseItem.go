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

package model

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type DefaultPlcReadResponseItem struct {
	code  model.PlcResponseCode
	value values.PlcValue
}

func NewReadResponseItem(code model.PlcResponseCode, value values.PlcValue) *DefaultPlcReadResponseItem {
	return &DefaultPlcReadResponseItem{
		code:  code,
		value: value,
	}
}

func (r *DefaultPlcReadResponseItem) GetCode() model.PlcResponseCode {
	return r.code
}

func (r *DefaultPlcReadResponseItem) GetValue() values.PlcValue {
	return r.value
}

func (d *DefaultPlcReadResponseItem) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := d.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (d *DefaultPlcReadResponseItem) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("ResponseItem"); err != nil {
		return err
	}

	{
		stringValue := fmt.Sprintf("%v", d.code)
		if err := writeBuffer.WriteString("result", uint32(len(stringValue)*8), "UTF-8", stringValue); err != nil {
			return err
		}
	}

	if d.value != nil {
		var elem interface{} = d.value
		if err := elem.(utils.Serializable).SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
			return err
		}
	}

	if err := writeBuffer.PopContext("ResponseItem"); err != nil {
		return err
	}
	return nil
}

func (d *DefaultPlcReadResponseItem) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), d); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
