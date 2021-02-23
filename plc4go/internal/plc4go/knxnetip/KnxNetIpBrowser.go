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
	driverModel "github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/model"
	apiModel "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/values"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type KnxNetIpBrowser struct {
	connection      *KnxNetIpConnection
	messageCodec    spi.MessageCodec
	sequenceCounter uint8
	spi.PlcBrowser
}

func NewKnxNetIpBrowser(connection *KnxNetIpConnection, messageCodec spi.MessageCodec) *KnxNetIpBrowser {
	return &KnxNetIpBrowser{
		connection:      connection,
		messageCodec:    messageCodec,
		sequenceCounter: 0,
	}
}

func (b KnxNetIpBrowser) Browse(browseRequest apiModel.PlcBrowseRequest) <-chan apiModel.PlcBrowseRequestResult {
	return b.BrowseWithInterceptor(browseRequest, func(result apiModel.PlcBrowseEvent) bool {
		return true
	})
}

func (b KnxNetIpBrowser) BrowseWithInterceptor(browseRequest apiModel.PlcBrowseRequest, interceptor func(result apiModel.PlcBrowseEvent) bool) <-chan apiModel.PlcBrowseRequestResult {
	result := make(chan apiModel.PlcBrowseRequestResult)
	sendResult := func(browseResponse apiModel.PlcBrowseResponse, err error) {
		result <- apiModel.PlcBrowseRequestResult{
			Request:  browseRequest,
			Response: browseResponse,
			Err:      err,
		}
	}

	go func() {
		results := map[string][]apiModel.PlcBrowseQueryResult{}
		for _, queryName := range browseRequest.GetQueryNames() {
			queryString := browseRequest.GetQueryString(queryName)
			field, err := b.connection.fieldHandler.ParseQuery(queryString)
			if err != nil {
				sendResult(nil, err)
				return
			}

			switch field.(type) {
			case KnxNetIpDeviceQueryField:
				queryResults, err := b.executeDeviceQuery(field.(KnxNetIpDeviceQueryField), browseRequest, queryName, interceptor)
				if err != nil {
					// TODO: Return some sort of return code like with the read and write APIs
					results[queryName] = nil
				} else {
					results[queryName] = queryResults
				}
			case KnxNetIpCommunicationObjectQueryField:
				queryResults, err := b.executeCommunicationObjectQuery(field.(KnxNetIpCommunicationObjectQueryField))
				if err != nil {
					// TODO: Return some sort of return code like with the read and write APIs
					results[queryName] = nil
				} else {
					results[queryName] = queryResults
				}
			default:
				// TODO: Return some sort of return code like with the read and write APIs
				results[queryName] = nil
			}
		}
		sendResult(model.NewDefaultPlcBrowseResponse(browseRequest, results), nil)
	}()
	return result
}

func (b KnxNetIpBrowser) executeDeviceQuery(field KnxNetIpDeviceQueryField, browseRequest apiModel.PlcBrowseRequest, queryName string, interceptor func(result apiModel.PlcBrowseEvent) bool) ([]apiModel.PlcBrowseQueryResult, error) {
	// Create a list of address strings, which doesn't contain any ranges, lists or wildcards
	knxAddresses, err := b.calculateAddresses(field)
	if err != nil {
		return nil, err
	}
	if len(knxAddresses) == 0 {
		return nil, errors.New("query resulted in not a single valid address")
	}

	var queryResults []apiModel.PlcBrowseQueryResult
	// Parse each of these expanded addresses and handle them accordingly.
	for _, knxAddress := range knxAddresses {
		// Send a connection request to the device
		deviceConnections := b.connection.ConnectToDevice(knxAddress)
		select {
		case deviceConnection := <-deviceConnections:
			if deviceConnection != nil {
				queryResult := apiModel.PlcBrowseQueryResult{
					Address: fmt.Sprintf("%d.%d.%d",
						knxAddress.MainGroup,
						knxAddress.MiddleGroup,
						knxAddress.SubGroup),
					PossibleDataTypes: nil,
				}

				// Pass it to the callback
				add := true
				if interceptor != nil {
					add = interceptor(apiModel.PlcBrowseEvent{
						Request:   browseRequest,
						QueryName: queryName,
						Result:    &queryResult,
						Err:       nil,
					})
				}

				// If the interceptor opted for adding it to the result, do so
				if add {
					queryResults = append(queryResults, queryResult)
				}

				deviceDisconnections := b.connection.DisconnectFromDevice(knxAddress)
				select {
				case _ = <-deviceDisconnections:
				case <-time.After(b.connection.defaultTtl * 10):
					// Just ignore this case ...
				}
			}
		case <-time.After(b.connection.defaultTtl):
			// In this case the remote was just not responding.
		}
		// Just to slow things down a bit (This way we can't exceed the max number of requests per minute)
		//time.Sleep(time.Millisecond * 20)
	}
	return queryResults, nil
}

func (b KnxNetIpBrowser) executeCommunicationObjectQuery(field KnxNetIpCommunicationObjectQueryField) ([]apiModel.PlcBrowseQueryResult, error) {
	var results []apiModel.PlcBrowseQueryResult

	knxAddress := field.toKnxAddress()
	knxAddressString := KnxAddressToString(knxAddress)

	// If we have a building Key, try that to login in order to access protected
	if b.connection.buildingKey != nil {
		arr := b.connection.AuthenticateDevice(*knxAddress, b.connection.buildingKey)
		<-arr
	}

	/////////////////////////////////////////////////////////////////////////////////////////////////////
	// Group Address Table reading
	/////////////////////////////////////////////////////////////////////////////////////////////////////

	// First of all, request the starting address of the group address table
	readRequestBuilder := b.connection.ReadRequestBuilder()
	readRequestBuilder.AddItem("groupAddressTableAddress", knxAddressString+"#1/7")
	readRequest, err := readRequestBuilder.Build()
	if err != nil {
		return nil, errors.New("error creating read request: " + err.Error())
	}
	rrr := readRequest.Execute()
	readResult := <-rrr
	if readResult.Err != nil {
		return nil, errors.New("error reading the group address table starting address: " + readResult.Err.Error())
	}
	if readResult.Response.GetResponseCode("groupAddressTableAddress") != apiModel.PlcResponseCode_OK {
		return nil, errors.New("error reading group address table starting address: " +
			readResult.Response.GetResponseCode("groupAddressTableAddress").GetName())
	}
	groupAddressTableStartAddress := readResult.Response.GetValue("groupAddressTableAddress").GetUint32()

	// Then read one byte at the given location.
	// This will return the number of entries in the group address table (each 2 bytes)
	readRequestBuilder = b.connection.ReadRequestBuilder()
	// Depending on the type of device, query an USINT (1 byte) or UINT (2 bytes)
	// TODO: Do this correctly depending on the device connection device-descriptor
	if b.connection.DeviceConnections[*knxAddress].deviceDescriptor == uint16(0x07B0) /* SystemB */ {
		readRequestBuilder.AddItem("numberOfAddressTableEntries",
			fmt.Sprintf("%s#%X:UINT", knxAddressString, groupAddressTableStartAddress))
	} else {
		readRequestBuilder.AddItem("numberOfAddressTableEntries",
			fmt.Sprintf("%s#%X:USINT", knxAddressString, groupAddressTableStartAddress))
	}
	readRequest, err = readRequestBuilder.Build()
	if err != nil {
		return nil, errors.New("error creating read request: " + err.Error())
	}
	rrr = readRequest.Execute()
	readResult = <-rrr
	if readResult.Err != nil {
		return nil, errors.New("error reading the number of group address table entries: " + readResult.Err.Error())
	}
	if readResult.Response.GetResponseCode("numberOfAddressTableEntries") != apiModel.PlcResponseCode_OK {
		return nil, errors.New("error reading the number of group address table entries: " +
			readResult.Response.GetResponseCode("numberOfAddressTableEntries").GetName())
	}
	numGroupAddresses := readResult.Response.GetValue("numberOfAddressTableEntries").GetUint16()

	if b.connection.DeviceConnections[*knxAddress].deviceDescriptor == uint16(0x07B0) /* SystemB */ {
		groupAddressTableStartAddress += 2
	} else {
		groupAddressTableStartAddress += 3
		numGroupAddresses--
	}

	// Abort, if there aren't any addresses to read.
	if numGroupAddresses == 0 {
		return results, nil
	}

	// Read the data in the group address table
	readRequestBuilder = b.connection.ReadRequestBuilder()
	readRequestBuilder.AddItem("groupAddressTable",
		fmt.Sprintf("%s#%X:UINT[%d]", knxAddressString, groupAddressTableStartAddress, numGroupAddresses))
	readRequest, err = readRequestBuilder.Build()
	if err != nil {
		return nil, errors.New("error creating read request: " + err.Error())
	}
	rrr = readRequest.Execute()
	readResult = <-rrr
	if readResult.Err != nil {
		return nil, errors.New("error reading the group address table content: " + readResult.Err.Error())
	}
	if (readResult.Response == nil) ||
		(readResult.Response.GetResponseCode("groupAddressTable") != apiModel.PlcResponseCode_OK) {
		return nil, errors.New("error reading the group address table content: " +
			readResult.Response.GetResponseCode("groupAddressTable").GetName())
	}
	var knxGroupAddresses []*driverModel.KnxGroupAddress
	if readResult.Response.GetValue("groupAddressTable").IsList() {
		for _, groupAddress := range readResult.Response.GetValue("groupAddressTable").GetList() {
			groupAddress := Uint16ToKnxGroupAddress(groupAddress.GetUint16(), 3)
			knxGroupAddresses = append(knxGroupAddresses, groupAddress)
		}
	} else {
		groupAddress := Uint16ToKnxGroupAddress(readResult.Response.GetValue("groupAddressTable").GetUint16(), 3)
		knxGroupAddresses = append(knxGroupAddresses, groupAddress)
	}

	/////////////////////////////////////////////////////////////////////////////////////////////////////
	// Group Address Association Table reading
	/////////////////////////////////////////////////////////////////////////////////////////////////////

	// Now we read the group address association table address
	readRequestBuilder = b.connection.ReadRequestBuilder()
	readRequestBuilder.AddItem("groupAddressAssociationTableAddress",
		fmt.Sprintf("%s#2/7", knxAddressString))
	readRequest, err = readRequestBuilder.Build()
	if err != nil {
		return nil, errors.New("error creating read request: " + err.Error())
	}
	rrr = readRequest.Execute()
	readResult = <-rrr
	if readResult.Err != nil {
		return nil, errors.New("error reading the group address association table content: " + readResult.Err.Error())
	}
	if (readResult.Response != nil) &&
		(readResult.Response.GetResponseCode("groupAddressAssociationTableAddress") != apiModel.PlcResponseCode_OK) {
		return nil, errors.New("error reading the group address association table content: " +
			readResult.Response.GetResponseCode("groupAddressAssociationTableAddress").GetName())
	}
	groupAddressAssociationTableAddress := readResult.Response.GetValue("groupAddressAssociationTableAddress").GetUint16()

	// Then read one uint16 at the given location.
	// This will return the number of entries in the group address table (each 2 bytes)
	readRequestBuilder = b.connection.ReadRequestBuilder()
	if b.connection.DeviceConnections[*knxAddress].deviceDescriptor == uint16(0x07B0) /* SystemB */ {
		readRequestBuilder.AddItem("numberOfGroupAddressAssociationTableEntries",
			fmt.Sprintf("%s#%X:UINT", knxAddressString, groupAddressAssociationTableAddress))
	} else {
		readRequestBuilder.AddItem("numberOfGroupAddressAssociationTableEntries",
			fmt.Sprintf("%s#%X:USINT", knxAddressString, groupAddressAssociationTableAddress))
	}
	readRequest, err = readRequestBuilder.Build()
	if err != nil {
		return nil, errors.New("error creating read request: " + err.Error())
	}
	rrr = readRequest.Execute()
	readResult = <-rrr
	if readResult.Err != nil {
		return nil, errors.New("error reading the number of group address association table entries: " + readResult.Err.Error())
	}
	if (readResult.Response != nil) &&
		(readResult.Response.GetResponseCode("numberOfGroupAddressAssociationTableEntries") != apiModel.PlcResponseCode_OK) {
		return nil, errors.New("error reading the number of group address association table entries: " +
			readResult.Response.GetResponseCode("numberOfGroupAddressAssociationTableEntries").GetName())
	}
	numberOfGroupAddressAssociationTableEntries := readResult.Response.GetValue("numberOfGroupAddressAssociationTableEntries").GetUint16()

	// Read the data in the group address table
	readRequestBuilder = b.connection.ReadRequestBuilder()
	// TODO: This request needs to be automatically split up into multiple requests.
	// Reasons for splitting up:
	// - Max APDU Size exceeded
	// - Max 63 bytes readable in one request, due to max of count field
	if b.connection.DeviceConnections[*knxAddress].deviceDescriptor == uint16(0x07B0) /* SystemB */ {
		readRequestBuilder.AddItem("groupAddressAssociationTable",
			fmt.Sprintf("%s#%X:UDINT[%d]", knxAddressString, groupAddressAssociationTableAddress+2, numberOfGroupAddressAssociationTableEntries))
	} else {
		readRequestBuilder.AddItem("groupAddressAssociationTable",
			fmt.Sprintf("%s#%X:UINT[%d]", knxAddressString, groupAddressAssociationTableAddress+1, numberOfGroupAddressAssociationTableEntries))
	}
	readRequest, err = readRequestBuilder.Build()
	if err != nil {
		return nil, errors.New("error creating read request: " + err.Error())
	}
	rrr = readRequest.Execute()
	readResult = <-rrr
	if readResult.Err != nil {
		return nil, errors.New("error reading the group address association table content: " + readResult.Err.Error())
	}
	if (readResult.Response != nil) &&
		(readResult.Response.GetResponseCode("groupAddressAssociationTable") != apiModel.PlcResponseCode_OK) {
		return nil, errors.New("error reading the group address association table content: " +
			readResult.Response.GetResponseCode("groupAddressAssociationTable").GetName())
	}
	// Output the group addresses
	if readResult.Response.GetValue("groupAddressAssociationTable").IsList() {
		for _, groupAddressAssociation := range readResult.Response.GetValue("groupAddressAssociationTable").GetList() {
			result := b.parseAssociationTable(knxAddressString, b.connection.DeviceConnections[*knxAddress].deviceDescriptor,
				knxGroupAddresses, groupAddressAssociation)
			if result != nil {
				results = append(results, *result)
			}
		}
	} else {
		result := b.parseAssociationTable(knxAddressString, b.connection.DeviceConnections[*knxAddress].deviceDescriptor,
			knxGroupAddresses, readResult.Response.GetValue("groupAddressAssociationTable"))
		if result != nil {
			results = append(results, *result)
		}
	}

	/////////////////////////////////////////////////////////////////////////////////////////////////////
	// Com Object Table reading (Not supported on all devices)
	// (This part is optional and experimental ...)
	/////////////////////////////////////////////////////////////////////////////////////////////////////

	// Now we read the group address association table address
	readRequestBuilder = b.connection.ReadRequestBuilder()
	readRequestBuilder.AddItem("comObjectTableAddress", fmt.Sprintf("%s#3/7", knxAddressString))
	readRequest, err = readRequestBuilder.Build()
	if err != nil {
		return nil, errors.New("error creating read request: " + err.Error())
	}
	rrr = readRequest.Execute()
	readResult = <-rrr
	if readResult.Response.GetResponseCode("comObjectTableAddress") == apiModel.PlcResponseCode_OK {
		comObjectTableAddress := readResult.Response.GetValue("comObjectTableAddress").GetUint16()
		log.Infof("Com Object Table Address: %x", comObjectTableAddress)
	}

	return results, nil
}

func (b KnxNetIpBrowser) calculateAddresses(field KnxNetIpDeviceQueryField) ([]driverModel.KnxAddress, error) {
	var explodedAddresses []driverModel.KnxAddress
	mainGroupOptions, err := b.explodeSegment(field.MainGroup, 1, 15)
	if err != nil {
		return nil, err
	}
	middleGroupOptions, err := b.explodeSegment(field.MiddleGroup, 1, 15)
	if err != nil {
		return nil, err
	}
	subGroupOptions, err := b.explodeSegment(field.SubGroup, 0, 255)
	if err != nil {
		return nil, err
	}
	for _, mainOption := range mainGroupOptions {
		for _, middleOption := range middleGroupOptions {
			for _, subOption := range subGroupOptions {
				// Don't try connecting to ourselves.
				if b.connection.ClientKnxAddress != nil {
					currentAddress := driverModel.KnxAddress{
						MainGroup:   mainOption,
						MiddleGroup: middleOption,
						SubGroup:    subOption,
					}
					explodedAddresses = append(explodedAddresses, currentAddress)
				}
			}
		}
	}
	return explodedAddresses, nil
}

func (b KnxNetIpBrowser) explodeSegment(segment string, min uint8, max uint8) ([]uint8, error) {
	var options []uint8
	if strings.Contains(segment, "*") {
		for i := min; i <= max; i++ {
			options = append(options, i)
		}
	} else if strings.HasPrefix(segment, "[") && strings.HasSuffix(segment, "]") {
		segment = strings.TrimPrefix(segment, "[")
		segment = strings.TrimSuffix(segment, "]")
		for _, segment := range strings.Split(segment, ",") {
			if strings.Contains(segment, "-") {
				split := strings.Split(segment, "-")
				localMin, err := strconv.Atoi(split[0])
				if err != nil {
					return nil, err
				}
				localMax, err := strconv.Atoi(split[1])
				if err != nil {
					return nil, err
				}
				for i := localMin; i <= localMax; i++ {
					options = append(options, uint8(i))
				}
			} else {
				option, err := strconv.Atoi(segment)
				if err != nil {
					return nil, err
				}
				options = append(options, uint8(option))
			}
		}
	} else {
		value, err := strconv.Atoi(segment)
		if err != nil {
			return nil, err
		}
		if uint8(value) >= min && uint8(value) <= max {
			options = append(options, uint8(value))
		}
	}
	return options, nil
}

func (m KnxNetIpBrowser) parseAssociationTable(knxAddressString string, deviceDescriptor uint16, knxGroupAddresses []*driverModel.KnxGroupAddress, value values.PlcValue) *apiModel.PlcBrowseQueryResult {
	var addressIndex uint16
	var comObjectNumber uint16
	if deviceDescriptor == uint16(0x07B0) /* SystemB */ {
		addressIndex = uint16((value.GetUint32()>>16)&0xFFFF) - 1
		comObjectNumber = uint16(value.GetUint32() & 0xFFFF)
	} else {
		addressIndex = ((value.GetUint16() >> 8) & 0xFF) - 1
		comObjectNumber = value.GetUint16() & 0xFF
	}
	if addressIndex < uint16(len(knxGroupAddresses)) {
		groupAddress := knxGroupAddresses[addressIndex]
		return &apiModel.PlcBrowseQueryResult{
			Address: fmt.Sprintf(
				"%s#%s %d", knxAddressString, GroupAddressToString(groupAddress), comObjectNumber),
			PossibleDataTypes: nil,
		}
	}
	return nil
}
