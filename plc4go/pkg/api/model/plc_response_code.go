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

	"github.com/apache/plc4x/plc4go/spi/utils"
)

type PlcResponseCode uint8

const (
	PlcResponseCode_OK               PlcResponseCode = 0x01
	PlcResponseCode_NOT_FOUND        PlcResponseCode = 0x02
	PlcResponseCode_ACCESS_DENIED    PlcResponseCode = 0x03
	PlcResponseCode_INVALID_ADDRESS  PlcResponseCode = 0x04
	PlcResponseCode_INVALID_DATATYPE PlcResponseCode = 0x05
	PlcResponseCode_INVALID_DATA     PlcResponseCode = 0x06
	PlcResponseCode_INTERNAL_ERROR   PlcResponseCode = 0x07
	PlcResponseCode_REMOTE_BUSY      PlcResponseCode = 0x08
	PlcResponseCode_REMOTE_ERROR     PlcResponseCode = 0x09
	PlcResponseCode_UNSUPPORTED      PlcResponseCode = 0x10
	PlcResponseCode_RESPONSE_PENDING PlcResponseCode = 0x11
	PlcResponseCode_REQUEST_TIMEOUT  PlcResponseCode = 0x12
)

func (m PlcResponseCode) GetName() string {
	switch m {
	case PlcResponseCode_OK:
		return "OK"
	case PlcResponseCode_NOT_FOUND:
		return "NOT_FOUND"
	case PlcResponseCode_ACCESS_DENIED:
		return "ACCESS_DENIED"
	case PlcResponseCode_INVALID_ADDRESS:
		return "INVALID_ADDRESS"
	case PlcResponseCode_INVALID_DATATYPE:
		return "INVALID_DATATYPE"
	case PlcResponseCode_INVALID_DATA:
		return "INVALID_DATA"
	case PlcResponseCode_INTERNAL_ERROR:
		return "INTERNAL_ERROR"
	case PlcResponseCode_REMOTE_BUSY:
		return "REMOTE_BUSY"
	case PlcResponseCode_REMOTE_ERROR:
		return "REMOTE_ERROR"
	case PlcResponseCode_UNSUPPORTED:
		return "UNSUPPORTED"
	case PlcResponseCode_RESPONSE_PENDING:
		return "RESPONSE_PENDING"
	case PlcResponseCode_REQUEST_TIMEOUT:
		return "REQUEST_TIMEOUT"
	default:
		return ""
	}
}

func (m PlcResponseCode) String() string {
	return m.GetName()
}

func (m PlcResponseCode) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased()
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m PlcResponseCode) SerializeWithWriteBuffer(_ context.Context, writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteUint8("ResponseCode", 8, uint8(m), utils.WithAdditionalStringRepresentation(m.GetName()))
}
