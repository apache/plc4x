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
	"encoding/binary"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type DefaultPlcRequest struct {
}

func NewDefaultPlcRequest() model.PlcRequest {
	return &DefaultPlcRequest{}
}

func (d *DefaultPlcRequest) IsAPlcMessage() bool {
	return true
}

func (d *DefaultPlcRequest) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := d.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (d *DefaultPlcRequest) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("DefaultPlcRequest"); err != nil {
		return err
	}
	if err := writeBuffer.PopContext("DefaultPlcRequest"); err != nil {
		return err
	}
	return nil
}

func (d *DefaultPlcRequest) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(d); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
