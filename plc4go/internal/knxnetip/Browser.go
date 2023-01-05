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
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	_default "github.com/apache/plc4x/plc4go/spi/default"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Browser struct {
	_default.DefaultBrowser
	connection      *Connection
	messageCodec    spi.MessageCodec
	sequenceCounter uint8
}

func NewBrowser(connection *Connection, messageCodec spi.MessageCodec) *Browser {
	browser := Browser{
		connection:      connection,
		messageCodec:    messageCodec,
		sequenceCounter: 0,
	}
	browser.DefaultBrowser = _default.NewDefaultBrowser(browser)
	return &browser
}

func (m Browser) BrowseQuery(ctx context.Context, browseRequest apiModel.PlcBrowseRequest, interceptor func(result apiModel.PlcBrowseItem) bool, queryName string, query apiModel.PlcQuery) (apiModel.PlcResponseCode, []apiModel.PlcBrowseItem) {
	switch query.(type) {
	case DeviceQuery:
		queryResults, err := m.executeDeviceQuery(ctx, query.(DeviceQuery), browseRequest, queryName, interceptor)
		if err != nil {
			log.Warn().Err(err).Msg("Error executing device query")
			return apiModel.PlcResponseCode_INTERNAL_ERROR, nil
		} else {
			return apiModel.PlcResponseCode_OK, queryResults
		}
	case CommunicationObjectQuery:
		queryResults, err := m.executeCommunicationObjectQuery(ctx, query.(CommunicationObjectQuery))
		if err != nil {
			log.Warn().Err(err).Msg("Error executing device query")
			return apiModel.PlcResponseCode_INTERNAL_ERROR, nil
		} else {
			return apiModel.PlcResponseCode_OK, queryResults
		}
	default:
		return apiModel.PlcResponseCode_INTERNAL_ERROR, nil
	}
}

func (m Browser) executeDeviceQuery(ctx context.Context, query DeviceQuery, browseRequest apiModel.PlcBrowseRequest, queryName string, interceptor func(result apiModel.PlcBrowseItem) bool) ([]apiModel.PlcBrowseItem, error) {
	// Create a list of address strings, which doesn't contain any ranges, lists or wildcards
	knxAddresses, err := m.calculateAddresses(query)
	if err != nil {
		return nil, err
	}
	if len(knxAddresses) == 0 {
		return nil, errors.New("query resulted in not a single valid address")
	}

	var queryResults []apiModel.PlcBrowseItem
	// Parse each of these expanded addresses and handle them accordingly.
	for _, knxAddress := range knxAddresses {
		// Send a connection request to the device
		connectTtlTimer := time.NewTimer(m.connection.defaultTtl)
		deviceConnections := m.connection.DeviceConnect(ctx, knxAddress)
		select {
		case deviceConnection := <-deviceConnections:
			if !connectTtlTimer.Stop() {
				<-connectTtlTimer.C
			}
			// If the request returned a connection, process it,
			// otherwise just ignore it.
			if deviceConnection.connection != nil {
				queryResult := &model.DefaultPlcBrowseItem{
					Tag: NewDeviceQuery(
						strconv.Itoa(int(knxAddress.GetMainGroup())),
						strconv.Itoa(int(knxAddress.GetMiddleGroup())),
						strconv.Itoa(int(knxAddress.GetSubGroup())),
					),
				}

				// Pass it to the callback
				add := true
				if interceptor != nil {
					add = interceptor(queryResult)
				}

				// If the interceptor opted for adding it to the result, do so
				if add {
					queryResults = append(queryResults, queryResult)
				}

				disconnectTtlTimer := time.NewTimer(m.connection.defaultTtl * 10)
				deviceDisconnections := m.connection.DeviceDisconnect(ctx, knxAddress)
				select {
				case _ = <-deviceDisconnections:
					if !disconnectTtlTimer.Stop() {
						<-disconnectTtlTimer.C
					}
				case <-disconnectTtlTimer.C:
					disconnectTtlTimer.Stop()
					// Just ignore this case ...
				}
			}
		case <-connectTtlTimer.C:
			connectTtlTimer.Stop()
			// In this case the remote was just not responding.
		}
		// Just to slow things down a bit (This way we can't exceed the max number of requests per minute)
		//time.Sleep(time.Millisecond * 20)
	}
	return queryResults, nil
}

func (m Browser) executeCommunicationObjectQuery(ctx context.Context, query CommunicationObjectQuery) ([]apiModel.PlcBrowseItem, error) {
	var results []apiModel.PlcBrowseItem

	knxAddress := query.toKnxAddress()
	knxAddressString := KnxAddressToString(knxAddress)

	// If we have a building Key, try that to login in order to access protected
	if m.connection.buildingKey != nil {
		arr := m.connection.DeviceAuthenticate(ctx, knxAddress, m.connection.buildingKey)
		<-arr
	}

	/////////////////////////////////////////////////////////////////////////////////////////////////////
	// Group Address Table reading
	/////////////////////////////////////////////////////////////////////////////////////////////////////

	// First, request the starting address of the group address table
	readRequestBuilder := m.connection.ReadRequestBuilder()
	readRequestBuilder.AddTagAddress("groupAddressTableAddress", knxAddressString+"#1/7")
	readRequest, err := readRequestBuilder.Build()
	if err != nil {
		return nil, errors.Wrap(err, "error creating read request")
	}
	rrr := readRequest.Execute()
	readResult := <-rrr
	if readResult.GetErr() != nil {
		return nil, errors.Wrap(readResult.GetErr(), "error reading the group address table starting address:")
	}
	if readResult.GetResponse().GetResponseCode("groupAddressTableAddress") != apiModel.PlcResponseCode_OK {
		return nil, errors.Errorf("error reading group address table starting address: %s",
			readResult.GetResponse().GetResponseCode("groupAddressTableAddress").GetName())
	}
	groupAddressTableStartAddress := readResult.GetResponse().GetValue("groupAddressTableAddress").GetUint32()

	// Then read one byte at the given location.
	// This will return the number of entries in the group address table (each 2 bytes)
	readRequestBuilder = m.connection.ReadRequestBuilder()
	// Depending on the type of device, query an USINT (1 byte) or UINT (2 bytes)
	// TODO: Do this correctly depending on the device connection device-descriptor
	if m.connection.DeviceConnections[knxAddress].deviceDescriptor == uint16(0x07B0) /* SystemB */ {
		readRequestBuilder.AddTagAddress("numberOfAddressTableEntries",
			fmt.Sprintf("%s#%X:UINT", knxAddressString, groupAddressTableStartAddress))
	} else {
		readRequestBuilder.AddTagAddress("numberOfAddressTableEntries",
			fmt.Sprintf("%s#%X:USINT", knxAddressString, groupAddressTableStartAddress))
	}
	readRequest, err = readRequestBuilder.Build()
	if err != nil {
		return nil, errors.Wrap(err, "error creating read request")
	}
	rrr = readRequest.Execute()
	readResult = <-rrr
	if readResult.GetErr() != nil {
		return nil, errors.Wrap(readResult.GetErr(), "error reading the number of group address table entries")
	}
	if readResult.GetResponse().GetResponseCode("numberOfAddressTableEntries") != apiModel.PlcResponseCode_OK {
		return nil, errors.Errorf("error reading the number of group address table entries: %s",
			readResult.GetResponse().GetResponseCode("numberOfAddressTableEntries").GetName())
	}
	numGroupAddresses := readResult.GetResponse().GetValue("numberOfAddressTableEntries").GetUint16()

	if m.connection.DeviceConnections[knxAddress].deviceDescriptor == uint16(0x07B0) /* SystemB */ {
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
	readRequest, err = m.connection.ReadRequestBuilder().
		AddTagAddress("groupAddressTable",
			fmt.Sprintf("%s#%X:UINT[%d]", knxAddressString, groupAddressTableStartAddress, numGroupAddresses)).
		Build()
	if err != nil {
		return nil, errors.Wrap(err, "error creating read request")
	}
	rrr = readRequest.Execute()
	readResult = <-rrr
	if readResult.GetErr() != nil {
		return nil, errors.Wrap(readResult.GetErr(), "error reading the group address table content")
	}
	if (readResult.GetResponse() == nil) ||
		(readResult.GetResponse().GetResponseCode("groupAddressTable") != apiModel.PlcResponseCode_OK) {
		return nil, errors.Errorf("error reading the group address table content: %s",
			readResult.GetResponse().GetResponseCode("groupAddressTable").GetName())
	}
	var knxGroupAddresses []driverModel.KnxGroupAddress
	if readResult.GetResponse().GetValue("groupAddressTable").IsList() {
		for _, groupAddress := range readResult.GetResponse().GetValue("groupAddressTable").GetList() {
			groupAddress := Uint16ToKnxGroupAddress(groupAddress.GetUint16(), 3)
			knxGroupAddresses = append(knxGroupAddresses, groupAddress)
		}
	} else {
		groupAddress := Uint16ToKnxGroupAddress(readResult.GetResponse().GetValue("groupAddressTable").GetUint16(), 3)
		knxGroupAddresses = append(knxGroupAddresses, groupAddress)
	}

	/////////////////////////////////////////////////////////////////////////////////////////////////////
	// Group Address Association Table reading
	/////////////////////////////////////////////////////////////////////////////////////////////////////

	// Now we read the group address association table address
	readRequestBuilder = m.connection.ReadRequestBuilder()
	readRequestBuilder.AddTagAddress("groupAddressAssociationTableAddress",
		fmt.Sprintf("%s#2/7", knxAddressString))
	readRequest, err = readRequestBuilder.Build()
	if err != nil {
		return nil, errors.Wrap(err, "error creating read request")
	}
	rrr = readRequest.Execute()
	readResult = <-rrr
	if readResult.GetErr() != nil {
		return nil, errors.Wrap(readResult.GetErr(), "error reading the group address association table address")
	}
	if (readResult.GetResponse() != nil) &&
		(readResult.GetResponse().GetResponseCode("groupAddressAssociationTableAddress") != apiModel.PlcResponseCode_OK) {
		return nil, errors.Errorf("error reading the group address association table address: %s",
			readResult.GetResponse().GetResponseCode("groupAddressAssociationTableAddress").GetName())
	}
	groupAddressAssociationTableAddress := readResult.GetResponse().GetValue("groupAddressAssociationTableAddress").GetUint16()

	// Then read one uint16 at the given location.
	// This will return the number of entries in the group address table (each 2 bytes)
	readRequestBuilder = m.connection.ReadRequestBuilder()
	if m.connection.DeviceConnections[knxAddress].deviceDescriptor == uint16(0x07B0) /* SystemB */ {
		readRequestBuilder.AddTagAddress("numberOfGroupAddressAssociationTableEntries",
			fmt.Sprintf("%s#%X:UINT", knxAddressString, groupAddressAssociationTableAddress))
	} else {
		readRequestBuilder.AddTagAddress("numberOfGroupAddressAssociationTableEntries",
			fmt.Sprintf("%s#%X:USINT", knxAddressString, groupAddressAssociationTableAddress))
	}
	readRequest, err = readRequestBuilder.Build()
	if err != nil {
		return nil, errors.Wrap(err, "error creating read request")
	}
	rrr = readRequest.Execute()
	readResult = <-rrr
	if readResult.GetErr() != nil {
		return nil, errors.Wrap(readResult.GetErr(), "error reading the number of group address association table entries")
	}
	if (readResult.GetResponse() != nil) &&
		(readResult.GetResponse().GetResponseCode("numberOfGroupAddressAssociationTableEntries") != apiModel.PlcResponseCode_OK) {
		return nil, errors.Errorf("error reading the number of group address association table entries: %s",
			readResult.GetResponse().GetResponseCode("numberOfGroupAddressAssociationTableEntries").GetName())
	}
	numberOfGroupAddressAssociationTableEntries := readResult.GetResponse().GetValue("numberOfGroupAddressAssociationTableEntries").GetUint16()

	// Read the data in the group address table
	readRequestBuilder = m.connection.ReadRequestBuilder()
	// TODO: This request needs to be automatically split up into multiple requests.
	// Reasons for splitting up:
	// - Max APDU Size exceeded
	// - Max 63 bytes readable in one request, due to max of count tag
	if m.connection.DeviceConnections[knxAddress].deviceDescriptor == uint16(0x07B0) /* SystemB */ {
		readRequestBuilder.AddTagAddress("groupAddressAssociationTable",
			fmt.Sprintf("%s#%X:UDINT[%d]", knxAddressString, groupAddressAssociationTableAddress+2, numberOfGroupAddressAssociationTableEntries))
	} else {
		readRequestBuilder.AddTagAddress("groupAddressAssociationTable",
			fmt.Sprintf("%s#%X:UINT[%d]", knxAddressString, groupAddressAssociationTableAddress+1, numberOfGroupAddressAssociationTableEntries))
	}
	readRequest, err = readRequestBuilder.Build()
	if err != nil {
		return nil, errors.Wrap(err, "error creating read request")
	}
	rrr = readRequest.Execute()
	readResult = <-rrr
	if readResult.GetErr() != nil {
		return nil, errors.Wrap(readResult.GetErr(), "error reading the group address association table content")
	}
	if (readResult.GetResponse() != nil) &&
		(readResult.GetResponse().GetResponseCode("groupAddressAssociationTable") != apiModel.PlcResponseCode_OK) {
		return nil, errors.Errorf("error reading the group address association table content: %s",
			readResult.GetResponse().GetResponseCode("groupAddressAssociationTable").GetName())
	}
	// Output the group addresses
	groupAddressComObjectNumberMapping := map[driverModel.KnxGroupAddress]uint16{}
	if readResult.GetResponse().GetValue("groupAddressAssociationTable").IsList() {
		for _, groupAddressAssociation := range readResult.GetResponse().GetValue("groupAddressAssociationTable").GetList() {
			groupAddress, comObjectNumber := m.parseAssociationTable(m.connection.DeviceConnections[knxAddress].deviceDescriptor,
				knxGroupAddresses, groupAddressAssociation)
			if groupAddress != nil {
				groupAddressComObjectNumberMapping[groupAddress] = comObjectNumber
			}
		}
	} else {
		groupAddress, comObjectNumber := m.parseAssociationTable(m.connection.DeviceConnections[knxAddress].deviceDescriptor,
			knxGroupAddresses, readResult.GetResponse().GetValue("groupAddressAssociationTable"))
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
	if m.connection.DeviceConnections[knxAddress].deviceDescriptor == uint16(0x07B0) /* SystemB */ {
		readRequestBuilder = m.connection.ReadRequestBuilder()
		// Read data for all com objects that are assigned a group address
		for _, comObjectNumber := range groupAddressComObjectNumberMapping {
			readRequestBuilder.AddTagAddress(strconv.Itoa(int(comObjectNumber)),
				fmt.Sprintf("%s#3/23/%d", knxAddressString, comObjectNumber))
		}
		readRequest, err = readRequestBuilder.Build()
		if err != nil {
			return nil, errors.Wrap(err, "error creating read request")
		}
		rrr = readRequest.Execute()
		readResult = <-rrr
		for groupAddress, comObjectNumber := range groupAddressComObjectNumberMapping {
			if readResult.GetResponse().GetResponseCode(strconv.Itoa(int(comObjectNumber))) != apiModel.PlcResponseCode_OK {
				continue
			}
			comObjectSettings := readResult.GetResponse().GetValue(strconv.Itoa(int(comObjectNumber))).GetUint16()
			data := []uint8{uint8((comObjectSettings >> 8) & 0xFF), uint8(comObjectSettings & 0xFF)}
			descriptor, err := driverModel.GroupObjectDescriptorRealisationTypeBParse(data)
			if err != nil {
				log.Info().Err(err).Msg("error parsing com object descriptor")
				continue
			}

			// Assemble a PlcBrowseFoundTag
			var tag apiModel.PlcTag
			communicationEnable := descriptor.GetCommunicationEnable()
			readable := communicationEnable && descriptor.GetReadEnable()
			writable := communicationEnable && descriptor.GetWriteEnable()
			subscribable := communicationEnable && descriptor.GetTransmitEnable()
			// Find a matching datatype for the given value-type.
			tagType := m.getTagTypeForValueType(descriptor.GetValueType())
			switch groupAddress := groupAddress.(type) {
			case driverModel.KnxGroupAddress3Level:
				address3Level := groupAddress
				tag = NewGroupAddress3LevelPlcTag(strconv.Itoa(int(address3Level.GetMainGroup())),
					strconv.Itoa(int(address3Level.GetMiddleGroup())), strconv.Itoa(int(address3Level.GetSubGroup())),
					&tagType)
			case driverModel.KnxGroupAddress2Level:
				address2Level := groupAddress
				tag = NewGroupAddress2LevelPlcTag(strconv.Itoa(int(address2Level.GetMainGroup())),
					strconv.Itoa(int(address2Level.GetSubGroup())),
					&tagType)
			case driverModel.KnxGroupAddressFreeLevel:
				address1Level := groupAddress
				tag = NewGroupAddress1LevelPlcTag(strconv.Itoa(int(address1Level.GetSubGroup())),
					&tagType)
			}

			results = append(results, &model.DefaultPlcBrowseItem{
				Tag:          tag,
				Name:         fmt.Sprintf("#%d", comObjectNumber),
				Readable:     readable,
				Writable:     writable,
				Subscribable: subscribable,
			})
		}
	} else if (m.connection.DeviceConnections[knxAddress].deviceDescriptor & 0xFFF0) == uint16(0x0700) /* System7 */ {
		// For System 7 Devices we unfortunately can't access the information of where the memory address for the
		// Com Object Table is programmatically, so we have to look up the address which is extracted from the XML data
		// Provided by the manufacturer. Unfortunately in order to be able to do this, we need to get the application
		// version from the device first.

		readRequestBuilder = m.connection.ReadRequestBuilder()
		readRequestBuilder.AddTagAddress("applicationProgramVersion", knxAddressString+"#3/13")
		readRequestBuilder.AddTagAddress("interfaceProgramVersion", knxAddressString+"#4/13")
		readRequest, err = readRequestBuilder.Build()
		if err != nil {
			return nil, errors.Wrap(err, "error creating read request")
		}

		rrr = readRequest.Execute()
		readRequestResult := <-rrr
		readResponse := readRequestResult.GetResponse()
		var programVersionData []byte
		if readResponse.GetResponseCode("applicationProgramVersion") == apiModel.PlcResponseCode_OK {
			programVersionData = utils.PlcValueUint8ListToByteArray(readResponse.GetValue("applicationProgramVersion"))
		} else if readResponse.GetResponseCode("interfaceProgramVersion") == apiModel.PlcResponseCode_OK {
			programVersionData = utils.PlcValueUint8ListToByteArray(readResponse.GetValue("interfaceProgramVersion"))
		}
		applicationId := hex.EncodeToString(programVersionData)

		// Lookup the com object table address
		comObjectTableAddresses, _ := driverModel.ComObjectTableAddressesByName("DEV" + strings.ToUpper(applicationId))
		if comObjectTableAddresses == 0 {
			return nil, errors.Errorf("error getting com address table address. No table entry for application id: %s", applicationId)
		}

		readRequestBuilder = m.connection.ReadRequestBuilder()
		// Read data for all com objects that are assigned a group address
		groupAddressMap := map[uint16][]driverModel.KnxGroupAddress{}
		for groupAddress, comObjectNumber := range groupAddressComObjectNumberMapping {
			if groupAddressMap[comObjectNumber] == nil {
				groupAddressMap[comObjectNumber] = []driverModel.KnxGroupAddress{}
			}
			groupAddressMap[comObjectNumber] = append(groupAddressMap[comObjectNumber], groupAddress)
			entryAddress := comObjectTableAddresses.ComObjectTableAddress() + 3 + (comObjectNumber * 4)
			readRequestBuilder.AddTagAddress(strconv.Itoa(int(comObjectNumber)),
				fmt.Sprintf("%s#%X:USINT[4]", knxAddressString, entryAddress))
		}
		readRequest, err = readRequestBuilder.Build()
		if err != nil {
			return nil, errors.Wrap(err, "error creating read request")
		}
		rrr = readRequest.Execute()
		readResult = <-rrr

		for _, tagName := range readResult.GetResponse().GetTagNames() {
			array := utils.PlcValueUint8ListToByteArray(readResult.GetResponse().GetValue(tagName))
			descriptor, err := driverModel.GroupObjectDescriptorRealisationType7Parse(array)
			if err != nil {
				return nil, errors.Wrap(err, "error creating read request")
			}

			// We saved the com object number in the tag name.
			comObjectNumber, _ := strconv.ParseUint(tagName, 10, 16)
			groupAddresses := groupAddressMap[uint16(comObjectNumber)]
			communicationEnable := descriptor.GetCommunicationEnable()
			readable := communicationEnable && descriptor.GetReadEnable()
			writable := communicationEnable && descriptor.GetWriteEnable()
			subscribable := communicationEnable && descriptor.GetTransmitEnable()
			// Find a matching datatype for the given value-type.
			tagType := m.getTagTypeForValueType(descriptor.GetValueType())

			// Create a tag for each of the given inputs.
			for _, groupAddress := range groupAddresses {
				tag := m.getTagForGroupAddress(groupAddress, tagType)

				results = append(results, &model.DefaultPlcBrowseItem{
					Tag:          tag,
					Name:         fmt.Sprintf("#%d", comObjectNumber),
					Readable:     readable,
					Writable:     writable,
					Subscribable: subscribable,
				})
			}
		}
	} else {
		readRequestBuilder = m.connection.ReadRequestBuilder()
		readRequestBuilder.AddTagAddress("comObjectTableAddress", fmt.Sprintf("%s#3/7", knxAddressString))
		readRequest, err = readRequestBuilder.Build()
		if err != nil {
			return nil, errors.Wrap(err, "error creating read request")
		}
		rrr = readRequest.Execute()
		readResult = <-rrr
		if readResult.GetResponse().GetResponseCode("comObjectTableAddress") == apiModel.PlcResponseCode_OK {
			comObjectTableAddress := readResult.GetResponse().GetValue("comObjectTableAddress").GetUint16()
			log.Info().Msgf("Com Object Table Address: %x", comObjectTableAddress)
		}
	}

	return results, nil
}

func (m Browser) calculateAddresses(query DeviceQuery) ([]driverModel.KnxAddress, error) {
	var explodedAddresses []driverModel.KnxAddress
	mainGroupOptions, err := m.explodeSegment(query.MainGroup, 1, 15)
	if err != nil {
		return nil, err
	}
	middleGroupOptions, err := m.explodeSegment(query.MiddleGroup, 1, 15)
	if err != nil {
		return nil, err
	}
	subGroupOptions, err := m.explodeSegment(query.SubGroup, 0, 255)
	if err != nil {
		return nil, err
	}
	for _, mainOption := range mainGroupOptions {
		for _, middleOption := range middleGroupOptions {
			for _, subOption := range subGroupOptions {
				// Don't try connecting to ourselves.
				if m.connection.ClientKnxAddress != nil {
					currentAddress := driverModel.NewKnxAddress(
						mainOption,
						middleOption,
						subOption,
					)
					explodedAddresses = append(explodedAddresses, currentAddress)
				}
			}
		}
	}
	return explodedAddresses, nil
}

func (m Browser) explodeSegment(segment string, min uint8, max uint8) ([]uint8, error) {
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
				localMin, err := strconv.ParseUint(split[0], 10, 8)
				if err != nil {
					return nil, err
				}
				localMax, err := strconv.ParseUint(split[1], 10, 8)
				if err != nil {
					return nil, err
				}
				for i := localMin; i <= localMax; i++ {
					options = append(options, uint8(i))
				}
			} else {
				option, err := strconv.ParseUint(segment, 10, 8)
				if err != nil {
					return nil, err
				}
				options = append(options, uint8(option))
			}
		}
	} else {
		value, err := strconv.ParseUint(segment, 10, 8)
		if err != nil {
			return nil, err
		}
		if uint8(value) >= min && uint8(value) <= max {
			options = append(options, uint8(value))
		}
	}
	return options, nil
}

func (m Browser) parseAssociationTable(deviceDescriptor uint16, knxGroupAddresses []driverModel.KnxGroupAddress, value values.PlcValue) (driverModel.KnxGroupAddress, uint16) {
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

func (m Browser) getTagForGroupAddress(groupAddress driverModel.KnxGroupAddress, datatype driverModel.KnxDatapointType) apiModel.PlcTag {
	switch groupAddress := groupAddress.(type) {
	case driverModel.KnxGroupAddress3Level:
		groupAddress3Level := groupAddress
		return GroupAddress3LevelPlcTag{
			MainGroup:   strconv.Itoa(int(groupAddress3Level.GetMainGroup())),
			MiddleGroup: strconv.Itoa(int(groupAddress3Level.GetMiddleGroup())),
			SubGroup:    strconv.Itoa(int(groupAddress3Level.GetSubGroup())),
			TagType:     &datatype,
		}
	case driverModel.KnxGroupAddress2Level:
		groupAddress2Level := groupAddress
		return GroupAddress2LevelPlcTag{
			MainGroup: strconv.Itoa(int(groupAddress2Level.GetMainGroup())),
			SubGroup:  strconv.Itoa(int(groupAddress2Level.GetSubGroup())),
			TagType:   &datatype,
		}
	case driverModel.KnxGroupAddressFreeLevel:
		groupAddress1Level := groupAddress
		return GroupAddress1LevelPlcTag{
			MainGroup: strconv.Itoa(int(groupAddress1Level.GetSubGroup())),
			TagType:   &datatype,
		}
	}
	return nil
}

func (m Browser) getTagTypeForValueType(valueType driverModel.ComObjectValueType) driverModel.KnxDatapointType {
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
