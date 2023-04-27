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
	"math"
	"strconv"
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api/values"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	values2 "github.com/apache/plc4x/plc4go/spi/values"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////
// KNX Specific Operations used by the driver internally
//
// These functions all provide access to some of the internal KNX operations
// They provide this functionality to other parts of the KNX driver, which are
// not part of the PLC4Go API.
//
// Remarks about these functions:
// They expect the called private functions to handle timeouts, so these will not.
///////////////////////////////////////////////////////////////////////////////////////////////////////

func (m *Connection) ReadGroupAddress(ctx context.Context, groupAddress []byte, datapointType *driverModel.KnxDatapointType) <-chan KnxReadResult {
	result := make(chan KnxReadResult)

	sendResponse := func(value values.PlcValue, numItems uint8, err error) {
		timeout := time.NewTimer(time.Millisecond * 10)
		select {
		case result <- KnxReadResult{
			value:    value,
			numItems: numItems,
			err:      err,
		}:
			if !timeout.Stop() {
				<-timeout.C
			}
		case <-timeout.C:
			timeout.Stop()
		}
	}

	go func() {
		groupAddressReadResponse, err := m.sendGroupAddressReadRequest(ctx, groupAddress)
		if err != nil {
			sendResponse(nil, 0, errors.Wrap(err, "error reading group address"))
			return
		}

		var payload []byte
		// TODO: maybe groupAddressReadResponse.DataFirstByte can be written as uint 6 so the we wouldn't need to cast
		payload = append(payload, byte(groupAddressReadResponse.GetDataFirstByte()))
		payload = append(payload, groupAddressReadResponse.GetData()...)

		// Parse the response data.
		rb := utils.NewReadBufferByteBased(payload)
		// If the size of the tag is greater than 6, we have to skip the first byte
		if datapointType.DatapointMainType().SizeInBits() > 6 {
			_, _ = rb.ReadUint8("datapointType", 8)
		}
		// Set a default datatype if none is provided
		if *datapointType == driverModel.KnxDatapointType_DPT_UNKNOWN {
			defaultDatapointType := driverModel.KnxDatapointType_USINT
			datapointType = &defaultDatapointType
		}
		// Parse the value
		plcValue, err := driverModel.KnxDatapointParseWithBuffer(context.Background(), rb, *datapointType)
		if err != nil {
			sendResponse(nil, 0, errors.Wrap(err, "error parsing group address response"))
			return
		}

		// Return the value
		sendResponse(plcValue, 1, nil)
	}()

	return result
}

func (m *Connection) DeviceConnect(ctx context.Context, targetAddress driverModel.KnxAddress) <-chan KnxDeviceConnectResult {
	result := make(chan KnxDeviceConnectResult)

	sendResponse := func(connection *KnxDeviceConnection, err error) {
		timeout := time.NewTimer(time.Millisecond * 10)
		select {
		case result <- KnxDeviceConnectResult{
			connection: connection,
			err:        err,
		}:
			if !timeout.Stop() {
				<-timeout.C
			}
		case <-timeout.C:
			timeout.Stop()
		}
	}

	go func() {
		// If we're already connected, use that connection instead.
		if connection, ok := m.DeviceConnections[targetAddress]; ok {
			sendResponse(connection, nil)
			return
		}

		// First send a connection request
		controlConnectResponse, err := m.sendDeviceConnectionRequest(ctx, targetAddress)
		if err != nil {
			sendResponse(nil, errors.Wrap(err, "error creating device connection"))
			return
		}
		if controlConnectResponse == nil {
			sendResponse(nil, errors.New("error creating device connection"))
			return
		}

		// Create the new connection object.
		connection := &KnxDeviceConnection{
			counter: 0,
			// I was told this value on the knx-forum.
			// Seems the max payload is 3 bytes less ...
			maxApdu: 0, // This is the default max APDU Size
		}
		m.DeviceConnections[targetAddress] = connection

		// If the connection request was successful, try to read the device-descriptor
		deviceDescriptorResponse, err := m.sendDeviceDeviceDescriptorReadRequest(ctx, targetAddress)
		if err != nil {
			sendResponse(nil, errors.New(
				"error reading device descriptor: "+err.Error()))
			return
		}
		// Save the device-descriptor value
		deviceDescriptor := uint16(deviceDescriptorResponse.GetData()[0])<<8 | (uint16(deviceDescriptorResponse.GetData()[1]) & 0xFF)
		connection.deviceDescriptor = deviceDescriptor

		// Last, not least, read the max APDU size
		// If we were able to read the max APDU size, then use the minimum of
		// the connection APDU size and the device APDU size, otherwise use the
		// default APDU Size of 15
		// Defined in: 03_05_01 Resources v01.09.03 AS Page 40
		deviceApduSize := uint16(15)
		propertyValueResponse, err := m.sendDevicePropertyReadRequest(ctx, targetAddress, 0, 56, 1, 1)
		if err == nil {
			// If the count is 0, then this property doesn't exist or the user has no permission to read it.
			// In all other cases we expect the response to contain the value.
			if propertyValueResponse.GetCount() > 0 {
				dataLength := uint8(len(propertyValueResponse.GetData()))
				data := propertyValueResponse.GetData()
				plcValue, err := driverModel.KnxPropertyParse(context.Background(), data,
					driverModel.KnxInterfaceObjectProperty_PID_DEVICE_MAX_APDULENGTH.PropertyDataType(), dataLength)

				// Return the result
				if err == nil {
					deviceApduSize = plcValue.GetUint16()
				} else {
					log.Debug().Err(err).Msgf("Error parsing knx property")
				}
			}
		}

		// Set the max apdu size for this connection.
		connection.maxApdu = uint16(math.Min(float64(deviceApduSize), 240))

		sendResponse(connection, nil)
	}()

	return result
}

func (m *Connection) DeviceDisconnect(ctx context.Context, targetAddress driverModel.KnxAddress) <-chan KnxDeviceDisconnectResult {
	result := make(chan KnxDeviceDisconnectResult)

	sendResponse := func(connection *KnxDeviceConnection, err error) {
		timeout := time.NewTimer(time.Millisecond * 10)
		select {
		case result <- KnxDeviceDisconnectResult{
			connection: connection,
			err:        err,
		}:
			if !timeout.Stop() {
				<-timeout.C
			}
		case <-timeout.C:
			timeout.Stop()
		}
	}

	go func() {
		if connection, ok := m.DeviceConnections[targetAddress]; ok {
			_, err := m.sendDeviceDisconnectionRequest(ctx, targetAddress)

			// Remove the connection from the list.
			delete(m.DeviceConnections, targetAddress)

			sendResponse(connection, err)
		} else {
			sendResponse(connection, nil)
		}
	}()

	return result
}

func (m *Connection) DeviceAuthenticate(ctx context.Context, targetAddress driverModel.KnxAddress, buildingKey []byte) <-chan KnxDeviceAuthenticateResult {
	result := make(chan KnxDeviceAuthenticateResult)

	sendResponse := func(err error) {
		timeout := time.NewTimer(time.Millisecond * 10)
		select {
		case result <- KnxDeviceAuthenticateResult{
			err: err,
		}:
			if !timeout.Stop() {
				<-timeout.C
			}
		case <-timeout.C:
			timeout.Stop()
		}
	}

	go func() {
		// Check if there is already a connection available,
		// if not, create a new one.
		connection, ok := m.DeviceConnections[targetAddress]
		if !ok {
			connections := m.DeviceConnect(ctx, targetAddress)
			deviceConnectionResult := <-connections
			// If we didn't get a connect, abort
			if deviceConnectionResult.err != nil {
				sendResponse(errors.Wrapf(deviceConnectionResult.err, "error connecting to device at: %s", KnxAddressToString(targetAddress)))
			}
		}

		// If we successfully got a connection, read the property
		if connection == nil {
			sendResponse(errors.New("unable to connect to device"))
			return
		}
		authenticationLevel := uint8(0)
		authenticationResponse, err := m.sendDeviceAuthentication(ctx, targetAddress, authenticationLevel, buildingKey)
		if err == nil {
			if authenticationResponse.GetLevel() == authenticationLevel {
				sendResponse(nil)
			} else {
				// We authenticated correctly but not to the level requested.
				sendResponse(errors.Errorf("got error authenticating at device %s",
					KnxAddressToString(targetAddress)))
			}
		} else {
			sendResponse(errors.Errorf("got error authenticating at device %s", KnxAddressToString(targetAddress)))
		}
	}()

	return result
}

func (m *Connection) DeviceReadProperty(ctx context.Context, targetAddress driverModel.KnxAddress, objectId uint8, propertyId uint8, propertyIndex uint16, numElements uint8) <-chan KnxReadResult {
	result := make(chan KnxReadResult)

	sendResponse := func(value values.PlcValue, numItems uint8, err error) {
		timeout := time.NewTimer(time.Millisecond * 10)
		select {
		case result <- KnxReadResult{
			value:    value,
			numItems: numItems,
			err:      err,
		}:
			if !timeout.Stop() {
				<-timeout.C
			}
		case <-timeout.C:
			timeout.Stop()
		}
	}

	go func() {
		// Check if there is already a connection available,
		// if not, create a new one.
		connection, ok := m.DeviceConnections[targetAddress]
		if !ok {
			connections := m.DeviceConnect(ctx, targetAddress)
			deviceConnectionResult := <-connections
			// If we didn't get a connect, abort
			if deviceConnectionResult.err != nil {
				sendResponse(nil,
					0,
					errors.Wrapf(deviceConnectionResult.err, "error connecting to device at: %s", KnxAddressToString(targetAddress)),
				)
			}
		}

		// If we successfully got a connection, read the property
		if connection == nil {
			sendResponse(nil, 0, errors.New("unable to connect to device"))
			return
		}
		propertyValueResponse, err := m.sendDevicePropertyReadRequest(ctx, targetAddress, objectId, propertyId, propertyIndex, numElements)
		if err != nil {
			sendResponse(nil, 0, err)
			return
		}

		// Find out the type of the property
		var objectType *driverModel.KnxInterfaceObjectType
		for curObjectType := driverModel.KnxInterfaceObjectType_OT_UNKNOWN; curObjectType <= driverModel.KnxInterfaceObjectType_OT_SUNBLIND_SENSOR_BASIC; curObjectType++ {
			if curObjectType.Code() == strconv.Itoa(int(objectId)) {
				objectType = &curObjectType
				break
			}
		}
		property := driverModel.KnxInterfaceObjectProperty_PID_UNKNOWN
		if objectType != nil {
			for curProperty := driverModel.KnxInterfaceObjectProperty_PID_UNKNOWN; curProperty <= driverModel.KnxInterfaceObjectProperty_PID_SUNBLIND_SENSOR_BASIC_ENABLE_TOGGLE_MODE; curProperty++ {
				if curProperty.PropertyId() == propertyId &&
					(curProperty.ObjectType() == driverModel.KnxInterfaceObjectType_OT_GENERAL ||
						curProperty.ObjectType() == *objectType) {
					property = curProperty
					break
				}
			}
		}

		dataLength := uint8(len(propertyValueResponse.GetData()))
		data := propertyValueResponse.GetData()
		plcValue, err := driverModel.KnxPropertyParse(context.Background(), data, property.PropertyDataType(), dataLength)
		if err != nil {
			sendResponse(nil, 0, err)
		} else {
			sendResponse(plcValue, 1, err)
		}
	}()

	return result
}

func (m *Connection) DeviceReadPropertyDescriptor(ctx context.Context, targetAddress driverModel.KnxAddress, objectId uint8, propertyId uint8) <-chan KnxReadResult {
	result := make(chan KnxReadResult)

	sendResponse := func(value values.PlcValue, numItems uint8, err error) {
		timeout := time.NewTimer(time.Millisecond * 10)
		select {
		case result <- KnxReadResult{
			value:    value,
			numItems: numItems,
			err:      err,
		}:
			if !timeout.Stop() {
				<-timeout.C
			}
		case <-timeout.C:
			timeout.Stop()
		}
	}

	go func() {
		// Check if there is already a connection available,
		// if not, create a new one.
		connection, ok := m.DeviceConnections[targetAddress]
		if !ok {
			connections := m.DeviceConnect(ctx, targetAddress)
			deviceConnectionResult := <-connections
			// If we didn't get a connect, abort
			if deviceConnectionResult.err != nil {
				sendResponse(
					nil,
					0,
					errors.Wrapf(deviceConnectionResult.err, "error connecting to device at: %s", KnxAddressToString(targetAddress)),
				)
			}
		}

		if connection == nil {
			sendResponse(nil, 0, errors.New("unable to connect to device"))
			return
		}
		// If we successfully got a connection, read the property
		propertyDescriptionResponse, err := m.sendDevicePropertyDescriptionReadRequest(ctx, targetAddress, objectId, propertyId)
		if err != nil {
			sendResponse(nil, 0, err)
			return
		}

		val := map[string]values.PlcValue{}
		val["writable"] = values2.NewPlcBOOL(propertyDescriptionResponse.GetWriteEnabled())
		val["dataType"] = values2.NewPlcSTRING(propertyDescriptionResponse.GetPropertyDataType().Name())
		val["maxElements"] = values2.NewPlcUINT(propertyDescriptionResponse.GetMaxNrOfElements())
		val["readLevel"] = values2.NewPlcSTRING(propertyDescriptionResponse.GetReadLevel().String())
		val["writeLevel"] = values2.NewPlcSTRING(propertyDescriptionResponse.GetWriteLevel().String())
		str := values2.NewPlcStruct(val)
		sendResponse(&str, 1, nil)
	}()

	return result
}

func (m *Connection) DeviceReadMemory(ctx context.Context, targetAddress driverModel.KnxAddress, address uint16, numElements uint8, datapointType *driverModel.KnxDatapointType) <-chan KnxReadResult {
	result := make(chan KnxReadResult)

	sendResponse := func(value values.PlcValue, numItems uint8, err error) {
		timeout := time.NewTimer(time.Millisecond * 10)
		select {
		case result <- KnxReadResult{
			value:    value,
			numItems: numItems,
			err:      err,
		}:
			if !timeout.Stop() {
				<-timeout.C
			}
		case <-timeout.C:
			timeout.Stop()
		}
	}

	go func() {
		// Set a default datatype, if none is specified
		if datapointType == nil {
			dpt := driverModel.KnxDatapointType_USINT
			datapointType = &dpt
		}

		// Check if there is already a connection available,
		// if not, create a new one.
		connection, ok := m.DeviceConnections[targetAddress]
		if !ok {
			connections := m.DeviceConnect(ctx, targetAddress)
			deviceConnectionResult := <-connections
			// If we didn't get a connect, abort
			if deviceConnectionResult.err != nil {
				sendResponse(
					nil,
					0,
					errors.Wrapf(deviceConnectionResult.err, "error connecting to device at: %s", KnxAddressToString(targetAddress)),
				)
			}
		}

		if connection == nil {
			// TODO: do we need to send a response here
			return
		}
		// If we successfully got a connection, read the property
		// Depending on the gateway Max APDU and the device Max APDU, split this up into multiple requests.
		// An APDU starts with the last 6 bits of the first data byte containing the count
		// followed by the 16-bit address, so these are already used.
		elementSize := datapointType.DatapointMainType().SizeInBits() / 8
		remainingRequestElements := numElements
		curStartingAddress := address
		var results []values.PlcValue
		for remainingRequestElements > 0 {
			// As the maxApdu can change, we have to do this in the loop.
			maxNumBytes := uint8(math.Min(float64(connection.maxApdu-3), float64(63)))
			maxNumElementsPerRequest := uint8(math.Floor(float64(maxNumBytes / elementSize)))
			numElements := uint8(math.Min(float64(remainingRequestElements), float64(maxNumElementsPerRequest)))
			numBytes := numElements * uint8(math.Max(float64(1), float64(datapointType.DatapointMainType().SizeInBits()/8)))
			memoryReadResponse, err := m.sendDeviceMemoryReadRequest(ctx, targetAddress, curStartingAddress, numBytes)
			if err != nil {
				// TODO: do we need to send a response here
				return
			}

			// If the number of bytes read is less than expected,
			// Update the connection.maxApdu value. This is required
			// as some devices seem to be sending back less than the
			// number of bytes specified than the maxApdu.
			if uint8(len(memoryReadResponse.GetData())) < numBytes {
				connection.maxApdu = uint16(len(memoryReadResponse.GetData()) + 3)
			}

			// Parse the data according to the property type information
			rb := utils.NewReadBufferByteBased(memoryReadResponse.GetData())
			for rb.HasMore(datapointType.DatapointMainType().SizeInBits()) {
				plcValue, err := driverModel.KnxDatapointParseWithBuffer(context.Background(), rb, *datapointType)
				// Return the result
				if err != nil {
					sendResponse(nil, 0, err)
					return
				}
				results = append(results, plcValue)

				// Update the counters and addresses.
				remainingRequestElements--
				curStartingAddress = curStartingAddress + uint16(elementSize)
			}
			// If there are still remaining bytes, keep them for the next time.
		}
		if len(results) > 1 {
			var plcList values.PlcValue
			plcList = values2.NewPlcList(results)
			sendResponse(plcList, 1, nil)
		} else if len(results) == 1 {
			sendResponse(results[0], 1, nil)
		}
	}()

	return result
}
