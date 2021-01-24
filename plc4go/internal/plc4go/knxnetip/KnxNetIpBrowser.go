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
	"fmt"
	driverModel "github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/model"
	apiModel "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
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
	result := make(chan apiModel.PlcBrowseRequestResult)
	sendResult := func(browseResponse apiModel.PlcBrowseResponse, err error) {
		select {
		case result <- apiModel.PlcBrowseRequestResult{
			Request:  browseRequest,
			Response: browseResponse,
			Err:      err,
		}:
		default:
		}
	}

	go func() {
		results := map[string][]apiModel.PlcBrowseQueryResult{}
		for _, queryName := range browseRequest.GetQueryNames() {
			queryString := browseRequest.GetQueryString(queryName)
			field, err := b.connection.fieldHandler.ParseQuery(queryString)
			if err != nil {
				sendResult(nil, err)
			}

			// Create a list of address strings, which doesn't contain any ranges, lists or wildcards
			options, err := b.calculateAddresses(field)
			if err != nil {
				sendResult(nil, err)
			}

			var queryResults []apiModel.PlcBrowseQueryResult
			// Parse each of these expanded addresses and handle them accordingly.
			for _, option := range options {
				field, err = b.connection.fieldHandler.ParseQuery(option)
				if err != nil {
					sendResult(nil, err)
				}

				// The following browse actions could be required:
				switch field.(type) {
				// - A Device Address
				//   - A Device has to be detected (This is done in every case)
				//      TODO: Send a Connect to the physical knx address
				//   - If an object-id is provided, check if this object id exists
				//   - If a property-id is provided, check if this property exists and try to get more information about it
				case KnxNetIpDevicePropertyAddressPlcField:
					individualAddress := field.(KnxNetIpDevicePropertyAddressPlcField)
					sourceAddress := &driverModel.KnxAddress{
						MainGroup:   0,
						MiddleGroup: 0,
						SubGroup:    0,
					}

					// Serialize the target address to a 2-byte value
					targetAddress := make([]int8, 2)
					main, _ := strconv.Atoi(individualAddress.MainGroup)
					middle, _ := strconv.Atoi(individualAddress.MiddleGroup)
					sub, _ := strconv.Atoi(individualAddress.SubGroup)
					targetAddress[0] = int8((main&0xF)<<4 | (middle & 0xF))
					targetAddress[1] = int8(sub)

					curSequenceCounter := b.sequenceCounter
					b.sequenceCounter++
					deviceConnectionRequest := driverModel.NewTunnelingRequest(
						driverModel.NewTunnelingRequestDataBlock(
							b.connection.CommunicationChannelId,
							curSequenceCounter),
						driverModel.NewLDataReq(0, nil,
							driverModel.NewLDataFrameDataExt(false, 6, uint8(0),
								sourceAddress, targetAddress,
								driverModel.NewApduControlContainer(driverModel.NewApduControlConnect(), 0, false, 0),
								true, true, driverModel.CEMIPriority_SYSTEM, false, false)))

					// Send the request
					done := make(chan bool)
					err = b.connection.SendRequest(
						deviceConnectionRequest,
						// The Gateway is now supposed to send an Ack to this request.
						func(message interface{}) bool {
							tunnelingRequest := driverModel.CastTunnelingRequest(message)
							if tunnelingRequest == nil ||
							    tunnelingRequest.TunnelingRequestDataBlock.CommunicationChannelId != b.connection.CommunicationChannelId {
								return false
							}
							lDataCon := driverModel.CastLDataCon(tunnelingRequest.Cemi)
							return lDataCon != nil
						},
						func(message interface{}) error {
							tunnelingRequest := driverModel.CastTunnelingRequest(message)
							lDataCon := driverModel.CastLDataCon(tunnelingRequest.Cemi)
							// If the error flag is not set, we've found a device
							if !lDataCon.DataFrame.ErrorFlag {
								queryResult := apiModel.PlcBrowseQueryResult{
									Address: fmt.Sprintf("%s.%s.%s",
										individualAddress.MainGroup,
										individualAddress.MiddleGroup,
										individualAddress.SubGroup),
									PossibleDataTypes: nil,
								}
								queryResults = append(queryResults, queryResult)
							}
							done <- true
							return nil
						},
						time.Second*1)
					select {
					case <-done:
					case <-time.After(time.Second * 2):
					}

					// Just to slow things down a bit (This way we can't exceed the max number of requests per minute)
					time.Sleep(time.Millisecond * 20)
				// - A Group Address
				//   - Check the cache of known group addresses. If there is data available from that group-id, it's returned
				case KnxNetIpGroupAddress3LevelPlcField:
				// - A Group Address
				//   - Check the cache of known group addresses. If there is data available from that group-id, it's returned
				case KnxNetIpGroupAddress2LevelPlcField:
				// - A Group Address
				//   - Check the cache of known group addresses. If there is data available from that group-id, it's returned
				case KnxNetIpGroupAddress1LevelPlcField:
				}
			}
			results[queryName] = queryResults
		}
		sendResult(model.NewDefaultPlcBrowseResponse(browseRequest, results), nil)
	}()
	return result
}

func (b KnxNetIpBrowser) calculateAddresses(field apiModel.PlcField) ([]string, error) {
	var addresses []string
	switch field.(type) {
	case KnxNetIpDevicePropertyAddressPlcField:
		propertyAddressField := field.(KnxNetIpDevicePropertyAddressPlcField)
		mainGroupOptions, err := b.explodeSegment(propertyAddressField.MainGroup, 1, 15)
		if err != nil {
			return nil, err
		}
		middleGroupOptions, err := b.explodeSegment(propertyAddressField.MiddleGroup, 1, 15)
		if err != nil {
			return nil, err
		}
		subGroupOptions, err := b.explodeSegment(propertyAddressField.SubGroup, 0, 255)
		if err != nil {
			return nil, err
		}
		for _, mainOption := range mainGroupOptions {
			for _, middleOption := range middleGroupOptions {
				for _, subOption := range subGroupOptions {
					address := fmt.Sprintf("%d.%d.%d", mainOption, middleOption, subOption)
					addresses = append(addresses, address)
				}
			}
		}
	case KnxNetIpGroupAddress3LevelPlcField:
	case KnxNetIpGroupAddress2LevelPlcField:
	case KnxNetIpGroupAddress1LevelPlcField:

	}
	return addresses, nil
}

func (b KnxNetIpBrowser) explodeSegment(segment string, min uint16, max uint16) ([]uint16, error) {
	var options []uint16
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
					options = append(options, uint16(i))
				}
			} else {
				option, err := strconv.Atoi(segment)
				if err != nil {
					return nil, err
				}
				options = append(options, uint16(option))
			}
		}
	}
	return options, nil
}
