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
	"strconv"
	"strings"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	internalModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	internalValues "github.com/apache/plc4x/plc4go/spi/values"
)

type Reader struct {
	connection *Connection
}

func NewReader(connection *Connection) *Reader {
	return &Reader{
		connection: connection,
	}
}

func (m Reader) Read(ctx context.Context, readRequest apiModel.PlcReadRequest) <-chan apiModel.PlcReadRequestResult {
	// TODO: handle ctx
	resultChan := make(chan apiModel.PlcReadRequestResult)
	go func() {
		responseCodes := map[string]apiModel.PlcResponseCode{}
		plcValues := map[string]apiValues.PlcValue{}

		// Sort the tags in direct properties and memory addresses, which will have to be actively
		// read from the devices and group-addresses which will be locally processed from the local cache.
		deviceAddresses := map[driverModel.KnxAddress]map[string]DeviceTag{}
		groupAddresses := map[string]GroupAddressTag{}
		for _, tagName := range readRequest.GetTagNames() {
			// Get the knx knxTag
			knxTag, err := CastToKnxTagFromPlcTag(readRequest.GetTag(tagName))
			if err != nil {
				responseCodes[tagName] = apiModel.PlcResponseCode_INVALID_ADDRESS
				plcValues[tagName] = nil
				continue
			}

			switch knxTag.(type) {
			case DevicePropertyAddressPlcTag:
				propertyTag := knxTag.(DevicePropertyAddressPlcTag)
				knxAddress := propertyTag.toKnxAddress()
				if knxAddress == nil {
					continue
				}
				if _, ok := deviceAddresses[knxAddress]; !ok {
					deviceAddresses[knxAddress] = map[string]DeviceTag{}
				}
				deviceAddresses[knxAddress][tagName] = propertyTag
			case DeviceMemoryAddressPlcTag:
				memoryTag := knxTag.(DeviceMemoryAddressPlcTag)
				knxAddress := memoryTag.toKnxAddress()
				if knxAddress == nil {
					continue
				}
				if _, ok := deviceAddresses[knxAddress]; !ok {
					deviceAddresses[knxAddress] = map[string]DeviceTag{}
				}
				deviceAddresses[knxAddress][tagName] = memoryTag
			case GroupAddressTag:
				groupAddressTag := knxTag.(GroupAddressTag)
				groupAddresses[tagName] = groupAddressTag
			default:
				responseCodes[tagName] = apiModel.PlcResponseCode_INVALID_ADDRESS
				plcValues[tagName] = nil
			}
		}

		// Process the direct properties.
		// Connect to each knx device and read all of the properties on that particular device.
		for deviceAddress, tags := range deviceAddresses {
			// Collect all the properties on this device
			for tagName, tag := range tags {
				switch tag.(type) {
				case DevicePropertyAddressPlcTag:
					propertyTag := tag.(DevicePropertyAddressPlcTag)

					timeout := time.NewTimer(m.connection.defaultTtl)
					results := m.connection.DeviceReadProperty(ctx, deviceAddress, propertyTag.ObjectId, propertyTag.PropertyId, propertyTag.PropertyIndex, propertyTag.NumElements)
					select {
					case result := <-results:
						if !timeout.Stop() {
							<-timeout.C
						}
						if result.err == nil {
							responseCodes[tagName] = apiModel.PlcResponseCode_OK
							plcValues[tagName] = *result.value
						} else {
							responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
							plcValues[tagName] = nil
						}
					case <-timeout.C:
						timeout.Stop()
						responseCodes[tagName] = apiModel.PlcResponseCode_REMOTE_BUSY
						plcValues[tagName] = nil
					}
				case DeviceMemoryAddressPlcTag:
					timeout := time.NewTimer(m.connection.defaultTtl)
					memoryTag := tag.(DeviceMemoryAddressPlcTag)
					results := m.connection.DeviceReadMemory(ctx, deviceAddress, memoryTag.Address, memoryTag.NumElements, memoryTag.TagType)
					select {
					case result := <-results:
						if !timeout.Stop() {
							<-timeout.C
						}
						if result.err == nil {
							responseCodes[tagName] = apiModel.PlcResponseCode_OK
							plcValues[tagName] = *result.value
						} else {
							responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
							plcValues[tagName] = nil
						}
					case <-timeout.C:
						timeout.Stop()
						responseCodes[tagName] = apiModel.PlcResponseCode_REMOTE_BUSY
						plcValues[tagName] = nil
					}
				}
			}
		}

		// Get the group address values from the cache
		for tagName, tag := range groupAddresses {
			responseCode, plcValue := m.readGroupAddress(ctx, tag)
			responseCodes[tagName] = responseCode
			plcValues[tagName] = plcValue
		}

		// Assemble the results
		result := internalModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues)
		resultChan <- &internalModel.DefaultPlcReadRequestResult{
			Request:  readRequest,
			Response: result,
			Err:      nil,
		}
	}()
	return resultChan
}

func (m Reader) readGroupAddress(ctx context.Context, tag GroupAddressTag) (apiModel.PlcResponseCode, apiValues.PlcValue) {
	rawAddresses, err := m.resolveAddresses(tag)
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
		stringAddress := NumericGroupAddressToString(numericAddress, tag)

		// Try to get a value from the cache
		m.connection.valueCacheMutex.RLock()
		int8s, ok := m.connection.valueCache[numericAddress]
		m.connection.valueCacheMutex.RUnlock()

		// If nothing was found in the cache, try to execute a group address read,
		// Otherwise respond with values from the cache.
		if !ok {
			addr := []byte{byte(numericAddress >> 8), byte(numericAddress & 0xFF)}
			rrc := m.connection.ReadGroupAddress(ctx, addr, tag.GetTagType())
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
				// TODO: Do we need a "default" case here?
			}
		} else {
			// If we don't have any tag-type information, add the raw data
			if tag.GetTagType() == nil {
				values[stringAddress] = internalValues.NewPlcRawByteArray(int8s)
			} else {
				// Decode the data according to the tags type
				rb := utils.NewReadBufferByteBased(int8s)
				if tag.GetTagType() == nil {
					return apiModel.PlcResponseCode_INVALID_DATATYPE, nil
				}
				// If the size of the tag is greater than 6, we have to skip the first byte
				if tag.GetTagType().GetLengthInBits() > 6 {
					_, _ = rb.ReadUint8("tagType", 8)
				}
				plcValue, err := driverModel.KnxDatapointParseWithBuffer(rb, *tag.GetTagType())
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
		stringAddress := NumericGroupAddressToString(rawAddresses[0], tag)
		return apiModel.PlcResponseCode_OK, values[stringAddress]
	} else if len(rawAddresses) > 1 {
		// Add it to the result
		return apiModel.PlcResponseCode_OK, internalValues.NewPlcStruct(values)
	} else {
		// Add it to the result
		return apiModel.PlcResponseCode_NOT_FOUND, nil
	}
}

// If the given tag is a tag containing a pattern, resolve to all the possible addresses
// it could be referring to.
func (m Reader) resolveAddresses(tag GroupAddressTag) ([]uint16, error) {
	// Depending on the type of tag, get the uint16 ids of all values that match the current tag
	var result []uint16
	switch tag.(type) {
	case GroupAddress3LevelPlcTag:
		address3LevelPlcTag := tag.(GroupAddress3LevelPlcTag)
		mainSegmentValues, err := m.resoleSegment(address3LevelPlcTag.MainGroup, 0, 31)
		if err != nil {
			return []uint16{}, err
		}
		middleSegmentValues, err := m.resoleSegment(address3LevelPlcTag.MiddleGroup, 0, 7)
		if err != nil {
			return []uint16{}, err
		}
		subSegmentValues, err := m.resoleSegment(address3LevelPlcTag.SubGroup, 0, 255)
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
	case GroupAddress2LevelPlcTag:
		address2LevelPlcTag := tag.(GroupAddress2LevelPlcTag)
		mainSegmentValues, err := m.resoleSegment(address2LevelPlcTag.MainGroup, 0, 31)
		if err != nil {
			return []uint16{}, err
		}
		subSegmentValues, err := m.resoleSegment(address2LevelPlcTag.SubGroup, 0, 2047)
		if err != nil {
			return []uint16{}, err
		}
		for _, main := range mainSegmentValues {
			for _, sub := range subSegmentValues {
				result = append(result, main<<11|sub)
			}
		}
	case GroupAddress1LevelPlcTag:
		address1LevelPlcTag := tag.(GroupAddress1LevelPlcTag)
		mainSegmentValues, err := m.resoleSegment(address1LevelPlcTag.MainGroup, 0, 65535)
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
