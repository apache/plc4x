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
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// TODO: add generator once we support enums
type DefaultPlcBrowseResponseItem struct {
	code    model.PlcResponseCode
	results []model.PlcBrowseItem
}

func NewBrowseResponseItem(code model.PlcResponseCode, results []model.PlcBrowseItem) *DefaultPlcBrowseResponseItem {
	return &DefaultPlcBrowseResponseItem{
		code:    code,
		results: results,
	}
}

func (d *DefaultPlcBrowseResponseItem) GetCode() model.PlcResponseCode {
	return d.code
}

func (d *DefaultPlcBrowseResponseItem) GetResults() []model.PlcBrowseItem {
	return d.results
}

func (d *DefaultPlcBrowseResponseItem) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := d.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (d *DefaultPlcBrowseResponseItem) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("PlcBrowseResponseItem"); err != nil {
		return err
	}

	{
		stringValue := fmt.Sprintf("%v", d.code)
		if err := writeBuffer.WriteString("code", uint32(len(stringValue)*8), "UTF-8", stringValue); err != nil {
			return err
		}
	}
	if err := writeBuffer.PushContext("results", utils.WithRenderAsList(true)); err != nil {
		return err
	}
	for _, elem := range d.results {
		var elem any = elem

		if elem != nil {
			if serializableField, ok := elem.(utils.Serializable); ok {
				if err := writeBuffer.PushContext("value"); err != nil {
					return err
				}
				if err := serializableField.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
					return err
				}
				if err := writeBuffer.PopContext("value"); err != nil {
					return err
				}
			} else {
				stringValue := fmt.Sprintf("%v", elem)
				if err := writeBuffer.WriteString("value", uint32(len(stringValue)*8), "UTF-8", stringValue); err != nil {
					return err
				}
			}
		}
	}
	if err := writeBuffer.PopContext("results", utils.WithRenderAsList(true)); err != nil {
		return err
	}
	if err := writeBuffer.PopContext("PlcBrowseResponseItem"); err != nil {
		return err
	}
	return nil
}

func (d *DefaultPlcBrowseResponseItem) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), d); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
