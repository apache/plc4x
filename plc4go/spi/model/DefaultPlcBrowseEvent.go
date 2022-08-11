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
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
)

type DefaultPlcBrowseEvent struct {
	Request   model.PlcBrowseRequest
	FieldName string
	Result    model.PlcBrowseFoundField
	Err       error
}

func (d *DefaultPlcBrowseEvent) IsAPlcMessage() bool {
	return true
}

func (d *DefaultPlcBrowseEvent) GetRequest() model.PlcBrowseRequest {
	return d.Request
}

func (d *DefaultPlcBrowseEvent) GetFieldName() string {
	return d.FieldName
}

func (d *DefaultPlcBrowseEvent) GetResult() model.PlcBrowseFoundField {
	return d.Result
}

func (d *DefaultPlcBrowseEvent) GetErr() error {
	return d.Err
}

func (d *DefaultPlcBrowseEvent) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("PlcBrowseEvent"); err != nil {
		return err
	}

	if d.Request != nil {
		if err := writeBuffer.PushContext("Request"); err != nil {
			return err
		}
		if serializableField, ok := d.Request.(utils.Serializable); ok {
			if err := serializableField.Serialize(writeBuffer); err != nil {
				return err
			}
		} else {
			return errors.Errorf("Error serializing. Field %T doesn't implement Serializable", d.Request)
		}
		if err := writeBuffer.PopContext("Request"); err != nil {
			return err
		}
	}
	if err := writeBuffer.WriteString("fieldName", uint32(len(d.FieldName)*8), "UTF-8", d.FieldName); err != nil {
		return err
	}
	if d.Result != nil {
		if err := writeBuffer.PushContext("Result"); err != nil {
			return err
		}
		if serializableField, ok := d.Result.(utils.Serializable); ok {
			if err := serializableField.Serialize(writeBuffer); err != nil {
				return err
			}
		} else {
			return errors.Errorf("Error serializing. Field %T doesn't implement Serializable", d.Result)
		}
		if err := writeBuffer.PopContext("Result"); err != nil {
			return err
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
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(d); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
