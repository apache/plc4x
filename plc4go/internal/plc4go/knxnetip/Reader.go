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
	driverModel "github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
	internalModel "github.com/apache/plc4x/plc4go/internal/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	internalValues "github.com/apache/plc4x/plc4go/internal/plc4go/spi/values"
	apiModel "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/plc4go/values"
	"strconv"
	"strings"
	"time"
)

type Reader struct {
	connection *Connection
}

func NewReader(connection *Connection) *Reader {
	return &Reader{
		connection: connection,
	}
}

func (m Reader) Read(readRequest apiModel.PlcReadRequest) <-chan apiModel.PlcReadRequestResult {
	resultChan := make(chan apiModel.PlcReadRequestResult)
	go func() {
		responseCodes := map[string]apiModel.PlcResponseCode{}
		plcValues := map[string]apiValues.PlcValue{}

		// Sort the fields in direct properties and memory addresses, which will have to be actively
		// read from the devices and group-addresses which will be locally processed from the local cache.
		deviceAddresses := map[driverModel.KnxAddress]map[string]DeviceField{}
		groupAddresses := map[string]GroupAddressField{}
		for _, fieldName := range readRequest.GetFieldNames() {
			// Get the knx field
			field, err := CastToFieldFromPlcField(readRequest.GetField(fieldName))
			if err != nil {
				responseCodes[fieldName] = apiModel.PlcResponseCode_INVALID_ADDRESS
				plcValues[fieldName] = nil
				continue
			}

			switch field.(type) {
			case DevicePropertyAddressPlcField:
				propertyField := field.(DevicePropertyAddressPlcField)
				knxAddress := propertyField.toKnxAddress()
				if knxAddress == nil {
					continue
				}
				if _, ok := deviceAddresses[*knxAddress]; !ok {
					deviceAddresses[*knxAddress] = map[string]DeviceField{}
				}
				deviceAddresses[*knxAddress][fieldName] = propertyField
			case DeviceMemoryAddressPlcField:
				memoryField := field.(DeviceMemoryAddressPlcField)
				knxAddress := memoryField.toKnxAddress()
				if knxAddress == nil {
					continue
				}
				if _, ok := deviceAddresses[*knxAddress]; !ok {
					deviceAddresses[*knxAddress] = map[string]DeviceField{}
				}
				deviceAddresses[*knxAddress][fieldName] = memoryField
			case CommunicationObjectQueryField:
				responseCodes[fieldName] = apiModel.PlcResponseCode_INVALID_ADDRESS
				plcValues[fieldName] = nil
				continue
			case GroupAddressField:
				groupAddressField := field.(GroupAddressField)
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
				case DevicePropertyAddressPlcField:
					propertyField := field.(DevicePropertyAddressPlcField)

					results := m.connection.DeviceReadProperty(deviceAddress, propertyField.ObjectId, propertyField.PropertyId, propertyField.PropertyIndex, propertyField.NumElements)
					select {
					case result := <-results:
						if result.err == nil {
							responseCodes[fieldName] = apiModel.PlcResponseCode_OK
							plcValues[fieldName] = *result.value
						} else {
							responseCodes[fieldName] = apiModel.PlcResponseCode_INTERNAL_ERROR
							plcValues[fieldName] = nil
						}
					case <-time.After(m.connection.defaultTtl):
						responseCodes[fieldName] = apiModel.PlcResponseCode_REMOTE_BUSY
						plcValues[fieldName] = nil
					}
				case DeviceMemoryAddressPlcField:
					memoryField := field.(DeviceMemoryAddressPlcField)
					results := m.connection.DeviceReadMemory(deviceAddress, memoryField.Address, memoryField.NumElements, memoryField.FieldType)
					select {
					case result := <-results:
						if result.err == nil {
							responseCodes[fieldName] = apiModel.PlcResponseCode_OK
							plcValues[fieldName] = *result.value
						} else {
							responseCodes[fieldName] = apiModel.PlcResponseCode_INTERNAL_ERROR
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

func (m Reader) readGroupAddress(field GroupAddressField) (apiModel.PlcResponseCode, apiValues.PlcValue) {
	rawAddresses, err := m.resolveAddresses(field)
	if err != nil {
		return apiModel.PlcResponseCode_INVALID_ADDRESS, nil
	}

	// First resolve any pattern to a list of fully qualified group addresses.
	// Then check if any of them is available in the local value cache,
	// if not send a group address read-request.
	values := map[string]apiValues.PlcValue{}
	returnCodes := map[string]apiModel.PlcResponseCode{}
	for _, numericAddress := range rawAddresses {
		// Create a string representation of this numeric address depending on the type of requested address
		stringAddress := NumericGroupAddressToString(numericAddress, field)

		// Try to get a value from the cache
		m.connection.valueCacheMutex.RLock()
		int8s, ok := m.connection.valueCache[numericAddress]
		m.connection.valueCacheMutex.RUnlock()

		// If nothing was found in the cache, try to execute a group address read,
		// Otherwise respond with values from the cache.
		if !ok {
			addr := []int8{int8(numericAddress >> 8), int8(numericAddress & 0xFF)}
			rrc := m.connection.ReadGroupAddress(addr, field.GetFieldType())
			select {
			case readResult := <-rrc:
				if readResult.value != nil {
					if readResult.err == nil {
						returnCodes[stringAddress] = apiModel.PlcResponseCode_OK
						values[stringAddress] = *readResult.value
					} else {
						returnCodes[stringAddress] = apiModel.PlcResponseCode_INTERNAL_ERROR
						values[stringAddress] = nil
					}
				} else {
					returnCodes[stringAddress] = apiModel.PlcResponseCode_NOT_FOUND
					values[stringAddress] = nil
				}
			}
		} else {
			// If we don't have any field-type information, add the raw data
			if field.GetTypeName() == "" {
				values[stringAddress] = internalValues.NewPlcByteArray(utils.Int8ArrayToByteArray(int8s))
			} else {
				// Decode the data according to the fields type
				rb := utils.NewReadBufferByteBased(utils.Int8ArrayToUint8Array(int8s))
				if field.GetFieldType() == nil {
					return apiModel.PlcResponseCode_INVALID_DATATYPE, nil
				}
				// If the size of the field is greater than 6, we have to skip the first byte
				if field.GetFieldType().LengthInBits() > 6 {
					_, _ = rb.ReadUint8("fieldType", 8)
				}
				plcValue, err := driverModel.KnxDatapointParse(rb, *field.GetFieldType())
				// If any of the values doesn't decode correctly, we can't return any
				if err != nil {
					return apiModel.PlcResponseCode_INVALID_DATA, nil
				}
				values[stringAddress] = plcValue
			}
		}
	}

	// If there is only one address to read, return this directly.
	// Otherwise return a struct, with the keys being the string representations of the address.
	if len(rawAddresses) == 1 {
		stringAddress := NumericGroupAddressToString(rawAddresses[0], field)
		return apiModel.PlcResponseCode_OK, values[stringAddress]
	} else if len(rawAddresses) > 1 {
		// Add it to the result
		return apiModel.PlcResponseCode_OK, internalValues.NewPlcStruct(values)
	} else {
		// Add it to the result
		return apiModel.PlcResponseCode_NOT_FOUND, nil
	}
}

// If the given field is a field containing a pattern, resolve to all the possible addresses
// it could be referring to.
func (m Reader) resolveAddresses(field GroupAddressField) ([]uint16, error) {
	// Depending on the type of field, get the uint16 ids of all values that match the current field
	var result []uint16
	switch field.(type) {
	case GroupAddress3LevelPlcField:
		address3LevelPlcField := field.(GroupAddress3LevelPlcField)
		mainSegmentValues, err := m.resoleSegment(address3LevelPlcField.MainGroup, 0, 31)
		if err != nil {
			return []uint16{}, err
		}
		middleSegmentValues, err := m.resoleSegment(address3LevelPlcField.MiddleGroup, 0, 7)
		if err != nil {
			return []uint16{}, err
		}
		subSegmentValues, err := m.resoleSegment(address3LevelPlcField.SubGroup, 0, 255)
		if err != nil {
			return []uint16{}, err
		}
		for _, main := range mainSegmentValues {
			for _, middle := range middleSegmentValues {
				for _, sub := range subSegmentValues {
					result = append(result, main<<11|middle<<8|sub)
				}
			}
		}
	case GroupAddress2LevelPlcField:
		address2LevelPlcField := field.(GroupAddress2LevelPlcField)
		mainSegmentValues, err := m.resoleSegment(address2LevelPlcField.MainGroup, 0, 31)
		if err != nil {
			return []uint16{}, err
		}
		subSegmentValues, err := m.resoleSegment(address2LevelPlcField.SubGroup, 0, 2047)
		if err != nil {
			return []uint16{}, err
		}
		for _, main := range mainSegmentValues {
			for _, sub := range subSegmentValues {
				result = append(result, main<<11|sub)
			}
		}
	case GroupAddress1LevelPlcField:
		address1LevelPlcField := field.(GroupAddress1LevelPlcField)
		mainSegmentValues, err := m.resoleSegment(address1LevelPlcField.MainGroup, 0, 65535)
		if err != nil {
			return []uint16{}, err
		}
		for _, main := range mainSegmentValues {
			result = append(result, main)
		}
	}
	return result, nil
}

func (m Reader) resoleSegment(pattern string, minValue uint16, maxValue uint16) ([]uint16, error) {
	var results []uint16
	// A "*" simply matches everything
	if pattern == "*" {
		for i := minValue; i <= maxValue; i++ {
			results = append(results, i)
		}
	} else if strings.HasPrefix(pattern, "[") && strings.HasSuffix(pattern, "]") {
		// If the pattern starts and ends with square brackets, it's a list of values or range queries
		// Multiple options are separated by ","
		for _, segment := range strings.Split(pattern[1:len(pattern)-1], ",") {
			// If the segment contains a "-", then it's a range query,
			// otherwise it's just a normal value.
			if strings.Contains(segment, "-") {
				split := strings.Split(segment, "-")
				if len(split) == 2 {
					minValue, err := strconv.ParseUint(split[0], 10, 16)
					if err != nil {
						return []uint16{}, errors.New("invalid address")
					}
					maxValue, err := strconv.ParseUint(split[1], 10, 16)
					if err != nil {
						return []uint16{}, errors.New("invalid address")
					}
					for i := uint16(minValue); i <= uint16(maxValue); i++ {
						results = append(results, i)
					}
				} else {
					return []uint16{}, errors.New("invalid address")
				}
			} else {
				value, err := strconv.ParseUint(segment, 10, 16)
				if err != nil {
					return []uint16{}, errors.New("invalid address")
				}
				results = append(results, uint16(value))
			}
		}
	} else {
		value, err := strconv.ParseUint(pattern, 10, 16)
		if err != nil {
			return []uint16{}, errors.New("invalid address")
		}
		results = append(results, uint16(value))
	}

	return results, nil
}
