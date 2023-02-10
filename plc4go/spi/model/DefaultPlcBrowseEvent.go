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

type DefaultPlcBrowseEvent struct {
	Request   model.PlcBrowseRequest
	QueryName string
	Result    model.PlcBrowseItem
	Err       error
}

func (d *DefaultPlcBrowseEvent) IsAPlcMessage() bool {
	return true
}

func (d *DefaultPlcBrowseEvent) GetRequest() model.PlcBrowseRequest {
	return d.Request
}

func (d *DefaultPlcBrowseEvent) GetQueryName() string {
	return d.QueryName
}

func (d *DefaultPlcBrowseEvent) GetResult() model.PlcBrowseItem {
	return d.Result
}

func (d *DefaultPlcBrowseEvent) GetErr() error {
	return d.Err
}

func (d *DefaultPlcBrowseEvent) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := d.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (d *DefaultPlcBrowseEvent) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("PlcBrowseEvent"); err != nil {
		return err
	}

	if d.Request != nil {
		if serializableField, ok := d.Request.(utils.Serializable); ok {
			if err := writeBuffer.PushContext("request"); err != nil {
				return err
			}
			if err := serializableField.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
				return err
			}
			if err := writeBuffer.PopContext("request"); err != nil {
				return err
			}
		} else {
			stringValue := fmt.Sprintf("%v", d.Request)
			if err := writeBuffer.WriteString("request", uint32(len(stringValue)*8), "UTF-8", stringValue); err != nil {
				return err
			}
		}
	}

	if err := writeBuffer.WriteString("queryName", uint32(len(d.QueryName)*8), "UTF-8", d.QueryName); err != nil {
		return err
	}

	if d.Result != nil {
		if serializableField, ok := d.Result.(utils.Serializable); ok {
			if err := writeBuffer.PushContext("result"); err != nil {
				return err
			}
			if err := serializableField.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
				return err
			}
			if err := writeBuffer.PopContext("result"); err != nil {
				return err
			}
		} else {
			stringValue := fmt.Sprintf("%v", d.Result)
			if err := writeBuffer.WriteString("result", uint32(len(stringValue)*8), "UTF-8", stringValue); err != nil {
				return err
			}
		}
	}

	if d.Err != nil {
		if err := writeBuffer.WriteString("err", uint32(len(d.Err.Error())*8), "UTF-8", d.Err.Error()); err != nil {
			return err
		}
	}
	if err := writeBuffer.PopContext("PlcBrowseEvent"); err != nil {
		return err
	}
	return nil
}

func (d *DefaultPlcBrowseEvent) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), d); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
