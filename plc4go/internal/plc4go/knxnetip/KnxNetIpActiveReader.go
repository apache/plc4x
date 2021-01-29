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
package knxnetip

import (
    "errors"
    "fmt"
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi"
    "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
)

type KnxNetIpActiveReader struct {
	messageCodec spi.MessageCodec
	spi.PlcWriter
}

func NewKnxNetIpActiveReader(messageCodec spi.MessageCodec) KnxNetIpActiveReader {
	return KnxNetIpActiveReader{
		messageCodec: messageCodec,
	}
}

func (m KnxNetIpActiveReader) Read(readRequest model.PlcReadRequest) <-chan model.PlcReadRequestResult {
	result := make(chan model.PlcReadRequestResult)
	// If we are requesting only one field, use a
	if len(readRequest.GetFieldNames()) == 1 {
		fieldName := readRequest.GetFieldNames()[0]

		// Get the KnxNetIp field instance from the request
		field := readRequest.GetField(fieldName)
		knxNetIpField, err := CastToKnxNetIpFieldFromPlcField(field)
		if err != nil {
			result <- model.PlcReadRequestResult{
				Request:  readRequest,
				Response: nil,
				Err:      errors.New("invalid field item type"),
			}
			return result
		}
		fmt.Printf("%v", knxNetIpField)

        /*knxnetipModel.LDataReq{
            AdditionalInformationLength: 0,
            AdditionalInformation:       nil,
            DataFrame: nil                  ,
            Parent:                      nil,
            ILDataReq:                   nil,
        }*/
	}
	return result
}
