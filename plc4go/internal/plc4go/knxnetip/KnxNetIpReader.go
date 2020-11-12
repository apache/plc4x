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
    driverModel "github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
    internalModel "github.com/apache/plc4x/plc4go/internal/plc4go/model"
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi"
    "github.com/apache/plc4x/plc4go/internal/plc4go/utils"
    "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
    "github.com/apache/plc4x/plc4go/pkg/plc4go/values"
)

type KnxNetIpReader struct {
    connection *KnxNetIpConnection
    spi.PlcReader
}

func NewKnxNetIpReader(connection *KnxNetIpConnection) *KnxNetIpReader {
    return &KnxNetIpReader{
        connection: connection,
    }
}

func (m KnxNetIpReader) Read(readRequest model.PlcReadRequest) <-chan model.PlcReadRequestResult {
    resultChan := make(chan model.PlcReadRequestResult)
    go func() {
        responseCodes := map[string]model.PlcResponseCode{}
        plcValues := map[string]values.PlcValue{}
        for _, fieldName := range readRequest.GetFieldNames() {
            field, err := CastToKnxNetIpFieldFromPlcField(readRequest.GetField(fieldName))
            if err != nil {
                responseCodes[fieldName] = model.PlcResponseCode_INVALID_ADDRESS
                plcValues[fieldName] = nil
                continue
            }

            // Serialize the field to an uint16
            wb := utils.NewWriteBuffer()
            err = field.toGroupAddress().Serialize(*wb)
            if err != nil {
                responseCodes[fieldName] = model.PlcResponseCode_INVALID_ADDRESS
                plcValues[fieldName] = nil
                continue
            }
            rawAddress := wb.GetBytes()
            address := (uint16(rawAddress[0]) << 8) | uint16(rawAddress[1] & 0xFF)

            // Get the value form the cache
            int8s, ok := m.connection.valueCache[address]
            if !ok {
                responseCodes[fieldName] = model.PlcResponseCode_NOT_FOUND
                plcValues[fieldName] = nil
                continue
            }

            // Decode the data according to the fields type
            rb := utils.NewReadBuffer(utils.Int8ToUint8(int8s))
            plcValue, err := driverModel.KnxDatapointParse(rb, field.GetTypeName())
            if err != nil {
                responseCodes[fieldName] = model.PlcResponseCode_INVALID_DATA
                plcValues[fieldName] = nil
                continue
            }

            // Add it to the result
            responseCodes[fieldName] = model.PlcResponseCode_OK
            plcValues[fieldName] = plcValue
        }
        result := internalModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues)
        resultChan <- model.PlcReadRequestResult{
            Request:  readRequest,
            Response: result,
            Err:      nil,
        }
    }()
    return resultChan
}





