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

package knxnetip

import (
	"context"
	"errors"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
)

type Writer struct {
	messageCodec spi.MessageCodec
}

func NewWriter(messageCodec spi.MessageCodec) Writer {
	return Writer{
		messageCodec: messageCodec,
	}
}

func (m Writer) Write(ctx context.Context, writeRequest apiModel.PlcWriteRequest) <-chan apiModel.PlcWriteRequestResult {
	// TODO: handle context
	result := make(chan apiModel.PlcWriteRequestResult, 1)
	// If we are requesting only one tag, use a
	if len(writeRequest.GetTagNames()) == 1 {
		tagName := writeRequest.GetTagNames()[0]

		// Get the KnxNetIp tag instance from the request
		tag := writeRequest.GetTag(tagName)
		groupAddressTag, err := CastToGroupAddressTagFromPlcTag(tag)
		if err != nil {
			result <- spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.New("invalid tag item type"))
			return result
		}

		// Get the value from the request and serialize it to a byte array
		value := writeRequest.GetValue(tagName)
		tagType := groupAddressTag.GetTagType()
		// TODO: why do we ignore the bytes here?
		if _, err := readWriteModel.KnxDatapointSerialize(value, *tagType); err != nil {
			result <- spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.New("error serializing value: "+err.Error()))
			return result
		}
	}
	return result
}
