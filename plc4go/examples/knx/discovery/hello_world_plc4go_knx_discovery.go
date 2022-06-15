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

package main

import (
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/drivers"
	"github.com/apache/plc4x/plc4go/pkg/api/logging"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func main() {
	// Set logging to INFO
	logging.InfoLevel()

	driverManager := plc4go.NewPlcDriverManager()
	drivers.RegisterKnxDriver(driverManager)

	var connectionStrings []string
	if len(os.Args) < 2 {
		// Try to auto-find KNX gateways via broadcast-message discovery
		_ = driverManager.Discover(func(event model.PlcDiscoveryEvent) {
			connStr := event.GetProtocolCode() + "://" + event.GetTransportUrl().Host
			log.Info().Str("connection string", connStr).Msg("Found KNX Gateway")

			connectionStrings = append(connectionStrings, connStr)
		})
		// Wait for 5 seconds for incoming responses
		time.Sleep(time.Second * 5)
	} else {
		connStr := "knxnet-ip://" + os.Args[1] + ":3671"
		log.Info().Str("connection string", connStr).Msg("Using manually provided KNX Gateway")
		connectionStrings = append(connectionStrings, connStr)
	}

	for _, connStr := range connectionStrings {
		log.Info().Str("connection string", connStr).Msg("Connecting")
		crc := driverManager.GetConnection(connStr)

		// Wait for the driver to connect (or not)
		connectionResult := <-crc
		if connectionResult.GetErr() != nil {
			log.Error().Msgf("error connecting to PLC: %s", connectionResult.GetErr().Error())
			return
		}
		log.Info().Str("connection string", connStr).Msg("Connected")
		connection := connectionResult.GetConnection()
		defer connection.BlockingClose()

		// Try to find all KNX devices on the current network
		browseRequest, err := connection.BrowseRequestBuilder().
			AddItem("allDevices", "[1-15].[1-15].[0-255]").
			//AddItem("allMyDevices", "[1-3].[1-6].[0-60]").
			//AddItem("onlyOneDevice", "1.1.20")
			Build()
		if err != nil {
			log.Error().Err(err).Msg("error creating browse request")
			return
		}
		brr := browseRequest.ExecuteWithInterceptor(func(result model.PlcBrowseEvent) bool {
			knxField := result.GetResult().GetField()
			knxAddress := knxField.GetAddressString()
			log.Info().Msgf("Inspecting detected Device at KNX Address: %s", knxAddress)

			// Try to get all the com-objects and the group addresses they are attached to.
			browseRequest, err := connection.BrowseRequestBuilder().
				AddItem("comObjects", knxAddress+"#com-obj").
				Build()
			if err != nil {
				log.Error().Err(err).Msg("error creating read request")
				return false
			}
			brr := browseRequest.Execute()
			browseResult := <-brr
			if browseResult.GetErr() != nil {
				log.Error().Err(browseResult.GetErr()).Msg("error executing the browse request for com-objects")
				return false
			}
			for _, result := range browseResult.GetResponse().GetQueryResults("comObjects") {
				permissions := ""
				if result.IsReadable() {
					permissions += "R"
				} else {
					permissions += " "
				}
				if result.IsWritable() {
					permissions += "W"
				} else {
					permissions += " "
				}
				if result.IsSubscribable() {
					permissions += "S"
				} else {
					permissions += " "
				}
				log.Info().Msgf(" - %15s (%s) %s", result.GetField().GetAddressString(), permissions, result.GetName())
			}

			readRequest, err := connection.ReadRequestBuilder().
				AddQuery("applicationProgramVersion", knxAddress+"#3/13").
				AddQuery("interfaceProgramVersion", knxAddress+"#4/13").
				Build()
			if err != nil {
				log.Error().Msgf("Error creating read request for scanning %s", knxAddress)
				return false
			}

			rrr := readRequest.Execute()
			readRequestResult := <-rrr

			if readRequestResult.GetErr() != nil {
				log.Error().Msgf("Error executing read request for reading device identification information %s", knxAddress)
				return false
			}
			readResponse := readRequestResult.GetResponse()
			var programVersion []byte
			if readResponse.GetResponseCode("applicationProgramVersion") == model.PlcResponseCode_OK {
				programVersion = utils.PlcValueUint8ListToByteArray(readResponse.GetValue("applicationProgramVersion"))
			} else if readResponse.GetResponseCode("interfaceProgramVersion") == model.PlcResponseCode_OK {
				programVersion = utils.PlcValueUint8ListToByteArray(readResponse.GetValue("interfaceProgramVersion"))
			}
			rb := utils.NewReadBufferByteBased(utils.ByteArrayToUint8Array(programVersion))
			manufacturerId := uint16(0)
			applicationId := uint16(0)
			applicationVersionMajor := uint8(0)
			applicationVersionMinor := uint8(0)
			if rb.GetTotalBytes() == 5 {
				manufacturerId, err = rb.ReadUint16("manufacturerId", 16)
				if err != nil {
					log.Error().Err(err).Msgf("Error reading manufacturer id from")
					return false
				}
				applicationId, err = rb.ReadUint16("applicationId", 16)
				if err != nil {
					log.Error().Err(err).Msgf("Error reading application id from")
					return false
				}
				applicationVersionMajor, err = rb.ReadUint8("applicationVersionMajor", 4)
				if err != nil {
					log.Error().Err(err).Msgf("Error reading application version major from %s", knxAddress)
					return false
				}
				applicationVersionMinor, err = rb.ReadUint8("applicationVersionMinor", 4)
				if err != nil {
					log.Error().Err(err).Msgf("Error reading application version minor from %s", knxAddress)
					return false
				}
			}

			log.Info().Msgf("     manufacturer id: %d", manufacturerId)
			log.Info().Msgf("     program id: %d version %d.%d", applicationId, applicationVersionMajor, applicationVersionMinor)

			return true
		})
		if brr == nil {
			log.Error().Msg("error executing browse request")
			return
		}
		select {
		case browseRequestResult := <-brr:
			log.Info().Msgf("Browse Request Result:\n%v", browseRequestResult)
		}
		return
	}
}
