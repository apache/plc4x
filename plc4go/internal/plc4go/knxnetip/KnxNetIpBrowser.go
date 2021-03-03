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
	"encoding/hex"
	"errors"
	"fmt"
	driverModel "github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
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
					Field: NewKnxNetIpDeviceQueryField(
						strconv.Itoa(int(knxAddress.MainGroup)),
						strconv.Itoa(int(knxAddress.MiddleGroup)),
						strconv.Itoa(int(knxAddress.SubGroup)),
					),
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
	readRequestBuilder.AddQuery("groupAddressTableAddress", knxAddressString+"#1/7")
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
		readRequestBuilder.AddQuery("numberOfAddressTableEntries",
			fmt.Sprintf("%s#%X:UINT", knxAddressString, groupAddressTableStartAddress))
	} else {
		readRequestBuilder.AddQuery("numberOfAddressTableEntries",
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
	readRequestBuilder.AddQuery("groupAddressTable",
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
	readRequestBuilder.AddQuery("groupAddressAssociationTableAddress",
		fmt.Sprintf("%s#2/7", knxAddressString))
	readRequest, err = readRequestBuilder.Build()
	if err != nil {
		return nil, errors.New("error creating read request: " + err.Error())
	}
	rrr = readRequest.Execute()
	readResult = <-rrr
	if readResult.Err != nil {
		return nil, errors.New("error reading the group address association table address: " + readResult.Err.Error())
	}
	if (readResult.Response != nil) &&
		(readResult.Response.GetResponseCode("groupAddressAssociationTableAddress") != apiModel.PlcResponseCode_OK) {
		return nil, errors.New("error reading the group address association table address: " +
			readResult.Response.GetResponseCode("groupAddressAssociationTableAddress").GetName())
	}
	groupAddressAssociationTableAddress := readResult.Response.GetValue("groupAddressAssociationTableAddress").GetUint16()

	// Then read one uint16 at the given location.
	// This will return the number of entries in the group address table (each 2 bytes)
	readRequestBuilder = b.connection.ReadRequestBuilder()
	if b.connection.DeviceConnections[*knxAddress].deviceDescriptor == uint16(0x07B0) /* SystemB */ {
		readRequestBuilder.AddQuery("numberOfGroupAddressAssociationTableEntries",
			fmt.Sprintf("%s#%X:UINT", knxAddressString, groupAddressAssociationTableAddress))
	} else {
		readRequestBuilder.AddQuery("numberOfGroupAddressAssociationTableEntries",
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
		readRequestBuilder.AddQuery("groupAddressAssociationTable",
			fmt.Sprintf("%s#%X:UDINT[%d]", knxAddressString, groupAddressAssociationTableAddress+2, numberOfGroupAddressAssociationTableEntries))
	} else {
		readRequestBuilder.AddQuery("groupAddressAssociationTable",
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
	groupAddressComObjectNumberMapping := map[*driverModel.KnxGroupAddress]uint16{}
	if readResult.Response.GetValue("groupAddressAssociationTable").IsList() {
		for _, groupAddressAssociation := range readResult.Response.GetValue("groupAddressAssociationTable").GetList() {
			groupAddress, comObjectNumber := b.parseAssociationTable(b.connection.DeviceConnections[*knxAddress].deviceDescriptor,
				knxGroupAddresses, groupAddressAssociation)
			if groupAddress != nil {
				groupAddressComObjectNumberMapping[groupAddress] = comObjectNumber
			}
		}
	} else {
		groupAddress, comObjectNumber := b.parseAssociationTable(b.connection.DeviceConnections[*knxAddress].deviceDescriptor,
			knxGroupAddresses, readResult.Response.GetValue("groupAddressAssociationTable"))
		if groupAddress != nil {
			groupAddressComObjectNumberMapping[groupAddress] = comObjectNumber
		}
	}

	/////////////////////////////////////////////////////////////////////////////////////////////////////
	// Com Object Table reading (Not supported on all devices)
	// (This part is optional and experimental ...)
	/////////////////////////////////////////////////////////////////////////////////////////////////////

	// In case of System B devices, the com object table is read as a property array
	// In this case we can even read only the com objects we're interested in.
	if b.connection.DeviceConnections[*knxAddress].deviceDescriptor == uint16(0x07B0) /* SystemB */ {
		readRequestBuilder = b.connection.ReadRequestBuilder()
		// Read data for all com objects that are assigned a group address
		for _, comObjectNumber := range groupAddressComObjectNumberMapping {
			readRequestBuilder.AddQuery(strconv.Itoa(int(comObjectNumber)),
				fmt.Sprintf("%s#3/23/%d", knxAddressString, comObjectNumber))
		}
		readRequest, err = readRequestBuilder.Build()
		if err != nil {
			return nil, errors.New("error creating read request: " + err.Error())
		}
		rrr = readRequest.Execute()
		readResult = <-rrr
		for groupAddress, comObjectNumber := range groupAddressComObjectNumberMapping {
			if readResult.Response.GetResponseCode(strconv.Itoa(int(comObjectNumber))) == apiModel.PlcResponseCode_OK {
				comObjectSettings := readResult.Response.GetValue(strconv.Itoa(int(comObjectNumber))).GetUint16()
				data := []uint8{uint8((comObjectSettings >> 8) & 0xFF), uint8(comObjectSettings & 0xFF)}
				rb := utils.NewReadBuffer(data)
				descriptor, err := driverModel.GroupObjectDescriptorRealisationTypeBParse(rb)
				if err != nil {
					log.Infof("error parsing com object descriptor: %s", err.Error())
					continue
				}

				// Assemble a PlcBrowseQueryResult
				var field apiModel.PlcField
				readable := descriptor.CommunicationEnable && descriptor.ReadEnable
				writable := descriptor.CommunicationEnable && descriptor.WriteEnable
				subscribable := descriptor.CommunicationEnable && descriptor.TransmitEnable
				// Find a matching datatype for the given value-type.
				fieldType := b.getFieldTypeForValueType(descriptor.ValueType)
				switch groupAddress.Child.(type) {
				case *driverModel.KnxGroupAddress3Level:
					address3Level := driverModel.CastKnxGroupAddress3Level(groupAddress)
					field = NewKnxNetIpGroupAddress3LevelPlcField(strconv.Itoa(int(address3Level.MainGroup)),
						strconv.Itoa(int(address3Level.MiddleGroup)), strconv.Itoa(int(address3Level.SubGroup)),
						&fieldType)
				case *driverModel.KnxGroupAddress2Level:
					address2Level := driverModel.CastKnxGroupAddress2Level(groupAddress)
					field = NewKnxNetIpGroupAddress2LevelPlcField(strconv.Itoa(int(address2Level.MainGroup)),
						strconv.Itoa(int(address2Level.SubGroup)),
						&fieldType)
				case *driverModel.KnxGroupAddressFreeLevel:
					address1Level := driverModel.CastKnxGroupAddressFreeLevel(groupAddress)
					field = NewKnxNetIpGroupAddress1LevelPlcField(strconv.Itoa(int(address1Level.SubGroup)),
						&fieldType)
				}

				results = append(results, apiModel.PlcBrowseQueryResult{
					Field:             field,
					Name:              fmt.Sprintf("#%d", comObjectNumber),
					Readable:          readable,
					Writable:          writable,
					Subscribable:      subscribable,
					PossibleDataTypes: nil,
				})
			}
		}
	} else if (b.connection.DeviceConnections[*knxAddress].deviceDescriptor & 0xFFF0) == uint16(0x0700) /* System7 */ {
		// For System 7 Devices we unfortunately can't access the information of where the memory address for the
		// Com Object Table is programmatically, so we have to lookup the address which is extracted from the XML data
		// Provided by the manufacturer. Unfortunately in order to be able to do this, we need to get the application
		// version from the device first.

		readRequestBuilder := b.connection.ReadRequestBuilder()
		readRequestBuilder.AddQuery("applicationProgramVersion", knxAddressString+"#3/13")
		readRequestBuilder.AddQuery("interfaceProgramVersion", knxAddressString+"#4/13")
		readRequest, err := readRequestBuilder.Build()
		if err != nil {
			return nil, errors.New("error creating read request: " + err.Error())
		}

		rrr := readRequest.Execute()
		readRequestResult := <-rrr
		readResponse := readRequestResult.Response
		var programVersionData []byte
		if readResponse.GetResponseCode("applicationProgramVersion") == apiModel.PlcResponseCode_OK {
			programVersionData = utils.PlcValueUint8ListToByteArray(readResponse.GetValue("applicationProgramVersion"))
		} else if readResponse.GetResponseCode("interfaceProgramVersion") == apiModel.PlcResponseCode_OK {
			programVersionData = utils.PlcValueUint8ListToByteArray(readResponse.GetValue("interfaceProgramVersion"))
		}
		applicationId := hex.EncodeToString(programVersionData)

		// Lookup the com object table address
		comObjectTableAddresses := driverModel.ComObjectTableAddressesByName("DEV" + strings.ToUpper(applicationId))
		if comObjectTableAddresses == 0 {
			return nil, errors.New("error getting com address table address. No table entry for application id: " + applicationId)
		}

		readRequestBuilder = b.connection.ReadRequestBuilder()
		// Read data for all com objects that are assigned a group address
		groupAddressMap := map[uint16][]*driverModel.KnxGroupAddress{}
		for groupAddress, comObjectNumber := range groupAddressComObjectNumberMapping {
			if groupAddressMap[comObjectNumber] == nil {
				groupAddressMap[comObjectNumber] = []*driverModel.KnxGroupAddress{}
			}
			groupAddressMap[comObjectNumber] = append(groupAddressMap[comObjectNumber], groupAddress)
			entryAddress := comObjectTableAddresses.ComObjectTableAddress() + 3 + (comObjectNumber * 4)
			readRequestBuilder.AddQuery(strconv.Itoa(int(comObjectNumber)),
				fmt.Sprintf("%s#%X:USINT[4]", knxAddressString, entryAddress))
		}
		readRequest, err = readRequestBuilder.Build()
		if err != nil {
			return nil, errors.New("error creating read request: " + err.Error())
		}
		rrr = readRequest.Execute()
		readResult = <-rrr

		for _, fieldName := range readResult.Response.GetFieldNames() {
			array := utils.PlcValueUint8ListToByteArray(readResult.Response.GetValue(fieldName))
			rb := utils.NewReadBuffer(array)
			descriptor, err := driverModel.GroupObjectDescriptorRealisationType7Parse(rb)
			if err != nil {
				return nil, errors.New("error creating read request: " + err.Error())
			}

			// We saved the com object number in the field name.
			comObjectNumber, _ := strconv.Atoi(fieldName)
			groupAddresses := groupAddressMap[uint16(comObjectNumber)]
			readable := descriptor.CommunicationEnable && descriptor.ReadEnable
			writable := descriptor.CommunicationEnable && descriptor.WriteEnable
			subscribable := descriptor.CommunicationEnable && descriptor.TransmitEnable
			// Find a matching datatype for the given value-type.
			fieldType := b.getFieldTypeForValueType(descriptor.ValueType)

			// Create a field for each of the given inputs.
			for _, groupAddress := range groupAddresses {
				field := b.getFieldForGroupAddress(groupAddress, fieldType)

				results = append(results, apiModel.PlcBrowseQueryResult{
					Field:             field,
					Name:              fmt.Sprintf("#%d", comObjectNumber),
					Readable:          readable,
					Writable:          writable,
					Subscribable:      subscribable,
					PossibleDataTypes: nil,
				})
			}
		}
	} else {
		readRequestBuilder = b.connection.ReadRequestBuilder()
		readRequestBuilder.AddQuery("comObjectTableAddress", fmt.Sprintf("%s#3/7", knxAddressString))
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

func (m KnxNetIpBrowser) parseAssociationTable(deviceDescriptor uint16, knxGroupAddresses []*driverModel.KnxGroupAddress, value values.PlcValue) (*driverModel.KnxGroupAddress, uint16) {
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
		return groupAddress, comObjectNumber
	}
	return nil, 0
}

func (m KnxNetIpBrowser) getFieldForGroupAddress(groupAddress *driverModel.KnxGroupAddress, datatype driverModel.KnxDatapointType) apiModel.PlcField {
	switch groupAddress.Child.(type) {
	case *driverModel.KnxGroupAddress3Level:
		groupAddress3Level := driverModel.CastKnxGroupAddress3Level(groupAddress)
		return KnxNetIpGroupAddress3LevelPlcField{
			MainGroup:   strconv.Itoa(int(groupAddress3Level.MainGroup)),
			MiddleGroup: strconv.Itoa(int(groupAddress3Level.MiddleGroup)),
			SubGroup:    strconv.Itoa(int(groupAddress3Level.SubGroup)),
			FieldType:   &datatype,
		}
	case *driverModel.KnxGroupAddress2Level:
		groupAddress2Level := driverModel.CastKnxGroupAddress2Level(groupAddress)
		return KnxNetIpGroupAddress2LevelPlcField{
			MainGroup: strconv.Itoa(int(groupAddress2Level.MainGroup)),
			SubGroup:  strconv.Itoa(int(groupAddress2Level.SubGroup)),
			FieldType: &datatype,
		}
	case *driverModel.KnxGroupAddressFreeLevel:
		groupAddress1Level := driverModel.CastKnxGroupAddressFreeLevel(groupAddress)
		return KnxNetIpGroupAddress1LevelPlcField{
			MainGroup: strconv.Itoa(int(groupAddress1Level.SubGroup)),
			FieldType: &datatype,
		}
	}
	return nil
}

func (m KnxNetIpBrowser) getFieldTypeForValueType(valueType driverModel.ComObjectValueType) driverModel.KnxDatapointType {
	switch valueType {
	case driverModel.ComObjectValueType_BIT1:
		return driverModel.KnxDatapointType_BOOL
	case driverModel.ComObjectValueType_BIT2:
		// Will be an array
		return driverModel.KnxDatapointType_BOOL
	case driverModel.ComObjectValueType_BIT3:
		// Will be an array
		return driverModel.KnxDatapointType_BOOL
	case driverModel.ComObjectValueType_BIT4:
		// Will be an array
		return driverModel.KnxDatapointType_BOOL
	case driverModel.ComObjectValueType_BIT5:
		// Will be an array
		return driverModel.KnxDatapointType_BOOL
	case driverModel.ComObjectValueType_BIT6:
		// Will be an array
		return driverModel.KnxDatapointType_BOOL
	case driverModel.ComObjectValueType_BIT7:
		// Will be an array
		return driverModel.KnxDatapointType_BOOL
	case driverModel.ComObjectValueType_BYTE1:
		return driverModel.KnxDatapointType_USINT
	case driverModel.ComObjectValueType_BYTE2:
		return driverModel.KnxDatapointType_UINT
	case driverModel.ComObjectValueType_BYTE3:
		return driverModel.KnxDatapointType_UDINT
	case driverModel.ComObjectValueType_BYTE4:
		return driverModel.KnxDatapointType_UDINT
	case driverModel.ComObjectValueType_BYTE6:
		// Will be an array
		return driverModel.KnxDatapointType_USINT
	case driverModel.ComObjectValueType_BYTE8:
		// Will be an array
		return driverModel.KnxDatapointType_USINT
	case driverModel.ComObjectValueType_BYTE10:
		// Will be an array
		return driverModel.KnxDatapointType_USINT
	case driverModel.ComObjectValueType_BYTE14:
		// Will be an array
		return driverModel.KnxDatapointType_USINT
	}
	// Just return "byte" in any other case.
	return driverModel.KnxDatapointType_USINT
}
