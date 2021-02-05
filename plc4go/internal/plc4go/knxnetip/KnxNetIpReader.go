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
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi"
    internalModel "github.com/apache/plc4x/plc4go/internal/plc4go/spi/model"
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
    internalValues "github.com/apache/plc4x/plc4go/internal/plc4go/spi/values"
    apiModel "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
    apiValues "github.com/apache/plc4x/plc4go/pkg/plc4go/values"
    "time"
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

func (m KnxNetIpReader) Read(readRequest apiModel.PlcReadRequest) <-chan apiModel.PlcReadRequestResult {
    resultChan := make(chan apiModel.PlcReadRequestResult)
    go func() {
        responseCodes := map[string]apiModel.PlcResponseCode{}
        plcValues := map[string]apiValues.PlcValue{}

        // Sort the fields in direct properties and memory addresses, which will have to be actively
        // read from the devices and group-addresses which will be locally processed from the local cache.
        deviceAddresses := map[driverModel.KnxAddress]map[string]KnxNetIpDeviceField{}
        groupAddresses := map[string]KnxNetIpGroupAddressField{}
        for _, fieldName := range readRequest.GetFieldNames() {
            // Get the knx field
            field, err := CastToKnxNetIpFieldFromPlcField(readRequest.GetField(fieldName))
            if err != nil {
                responseCodes[fieldName] = apiModel.PlcResponseCode_INVALID_ADDRESS
                plcValues[fieldName] = nil
                continue
            }

            switch field.(type) {
            case KnxNetIpDevicePropertyAddressPlcField:
                propertyField := field.(KnxNetIpDevicePropertyAddressPlcField)
                knxAddress := propertyField.toKnxAddress()
                if knxAddress == nil {
                    continue
                }
                if _, ok := deviceAddresses[*knxAddress]; !ok {
                    deviceAddresses[*knxAddress] = map[string]KnxNetIpDeviceField{}
                }
                deviceAddresses[*knxAddress][fieldName] = propertyField
            case KnxNetIpDeviceMemoryAddressPlcField:
                memoryField := field.(KnxNetIpDeviceMemoryAddressPlcField)
                knxAddress := memoryField.toKnxAddress()
                if knxAddress == nil {
                    continue
                }
                if _, ok := deviceAddresses[*knxAddress]; !ok {
                    deviceAddresses[*knxAddress] = map[string]KnxNetIpDeviceField{}
                }
                deviceAddresses[*knxAddress][fieldName] = memoryField
            case KnxNetIpCommunicationObjectQueryField:
                responseCodes[fieldName] = apiModel.PlcResponseCode_INVALID_ADDRESS
                plcValues[fieldName] = nil
                continue
            case KnxNetIpGroupAddressField:
                groupAddressField := field.(KnxNetIpGroupAddressField)
                groupAddresses[fieldName] = groupAddressField
            default:
                responseCodes[fieldName] = apiModel.PlcResponseCode_INVALID_ADDRESS
                plcValues[fieldName] = nil
            }
        }

        // Process the direct properties.
        // Connect to each knx device and read all of the properties on that particular device.
        for deviceAddress, fields := range deviceAddresses {
            // Collect all the properties on this device
            for fieldName, field := range fields {
                switch field.(type) {
                case KnxNetIpDevicePropertyAddressPlcField:
                    propertyField := field.(KnxNetIpDevicePropertyAddressPlcField)

                    results := m.connection.ReadDeviceProperty(deviceAddress, propertyField.ObjectId, propertyField.PropertyId)
                    select {
                    case result := <-results:
                        responseCodes[fieldName] = result.returnCode
                        if result.returnCode == apiModel.PlcResponseCode_OK {
                            plcValues[fieldName] = *result.value
                        } else {
                            plcValues[fieldName] = nil
                        }
                    case <-time.After(m.connection.defaultTtl):
                        responseCodes[fieldName] = apiModel.PlcResponseCode_REMOTE_BUSY
                        plcValues[fieldName] = nil
                    }
                case KnxNetIpDeviceMemoryAddressPlcField:
                    memoryField := field.(KnxNetIpDeviceMemoryAddressPlcField)
                    results := m.connection.ReadDeviceMemory(deviceAddress, memoryField.Address, memoryField.NumberOfElements, memoryField.FieldType)
                    select {
                    case result := <-results:
                        responseCodes[fieldName] = result.returnCode
                        if result.returnCode == apiModel.PlcResponseCode_OK {
                            plcValues[fieldName] = *result.value
                        } else {
                            plcValues[fieldName] = nil
                        }
                    case <-time.After(m.connection.defaultTtl):
                        responseCodes[fieldName] = apiModel.PlcResponseCode_REMOTE_BUSY
                        plcValues[fieldName] = nil
                    }
                }
            }
        }

        // Get the group address values from the cache
        for fieldName, field := range groupAddresses {
            responseCode, plcValue := m.readGroupAddress(field)
            responseCodes[fieldName] = responseCode
            plcValues[fieldName] = plcValue
        }

        // Assemble the results
        result := internalModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues)
        resultChan <- apiModel.PlcReadRequestResult{
            Request:  readRequest,
            Response: result,
            Err:      nil,
        }
    }()
    return resultChan
}

func (m KnxNetIpReader) readGroupAddress(field KnxNetIpGroupAddressField) (apiModel.PlcResponseCode, apiValues.PlcValue) {
    // Pattern fields can match more than one value, therefore we have to handle things differently
    if field.IsPatternField() {
        // Depending on the type of field, get the uint16 ids of all values that match the current field
        matchedAddresses := map[uint16]*driverModel.KnxGroupAddress{}
        switch field.(type) {
        case KnxNetIpGroupAddress3LevelPlcField:
            for key, value := range m.connection.leve3AddressCache {
                if field.matches(value.Parent) {
                    matchedAddresses[key] = value.Parent
                }
            }
        case KnxNetIpGroupAddress2LevelPlcField:
            for key, value := range m.connection.leve2AddressCache {
                if field.matches(value.Parent) {
                    matchedAddresses[key] = value.Parent
                }
            }
        case KnxNetIpGroupAddress1LevelPlcField:
            for key, value := range m.connection.leve1AddressCache {
                if field.matches(value.Parent) {
                    matchedAddresses[key] = value.Parent
                }
            }
        }

        // If not a single match was found, we'll return a "not found" message
        if len(matchedAddresses) == 0 {
            return apiModel.PlcResponseCode_NOT_FOUND, nil
        }

        // Go through all of the values and create a plc-struct from them
        // where the string version of the address becomes the property name
        // and the property value is the corresponding value (Other wise it
        // would be impossible to know which of the fields the pattern matched
        // a given value belongs to)
        values := map[string]apiValues.PlcValue{}
        for numericAddress, address := range matchedAddresses {
            // Get the raw data from the cache
            m.connection.valueCacheMutex.RLock()
            int8s, _ := m.connection.valueCache[numericAddress]
            m.connection.valueCacheMutex.RUnlock()

            // If we don't have any field-type information, add the raw data
            if field.GetTypeName() == "" {
                values[GroupAddressToString(address)] =
                    internalValues.NewPlcByteArray(utils.Int8ArrayToByteArray(int8s))
            } else {
                // Decode the data according to the fields type
                rb := utils.NewReadBuffer(utils.Int8ArrayToUint8Array(int8s))
                if field.GetFieldType() == nil {
                    return apiModel.PlcResponseCode_INVALID_DATATYPE, nil
                }
                // If the size of the field is greater than 6, we have to skip the first byte
                if field.GetFieldType().LengthInBits() > 6 {
                    _, _ = rb.ReadUint8(8)
                }
                plcValue, err := driverModel.KnxDatapointParse(rb, *field.GetFieldType())
                // If any of the values doesn't decode correctly, we can't return any
                if err != nil {
                    return apiModel.PlcResponseCode_INVALID_DATA, nil
                }
                values[GroupAddressToString(address)] = plcValue
            }
        }

        // Add it to the result
        return apiModel.PlcResponseCode_OK, internalValues.NewPlcStruct(values)
    } else {
        // If it's not a pattern field, we can access the cached value a lot simpler

        // Serialize the field to an uint16
        wb := utils.NewWriteBuffer()
        err := field.toGroupAddress().Serialize(*wb)
        if err != nil {
            return apiModel.PlcResponseCode_INVALID_ADDRESS, nil
        }
        rawAddress := wb.GetBytes()
        address := (uint16(rawAddress[0]) << 8) | uint16(rawAddress[1]&0xFF)

        // Get the value form the cache
        m.connection.valueCacheMutex.RLock()
        int8s, ok := m.connection.valueCache[address]
        m.connection.valueCacheMutex.RUnlock()
        if !ok {
            return apiModel.PlcResponseCode_NOT_FOUND, nil
        }

        // Decode the data according to the fields type
        rb := utils.NewReadBuffer(utils.Int8ArrayToUint8Array(int8s))
        if field.GetFieldType() == nil {
            return apiModel.PlcResponseCode_INVALID_DATATYPE, nil
        }
        // If the size of the field is greater than 6, we have to skip the first byte
        if field.GetFieldType().LengthInBits() > 6 {
            _, _ = rb.ReadUint8(8)
        }
        plcValue, err := driverModel.KnxDatapointParse(rb, *field.GetFieldType())
        if err != nil {
            return apiModel.PlcResponseCode_INVALID_DATA, nil
        }

        // Add it to the result
        return apiModel.PlcResponseCode_OK, plcValue
    }
}
