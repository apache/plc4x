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
package main

import (
	"encoding/hex"
	"github.com/apache/plc4x/plc4go/pkg/plc4go"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/drivers"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/transports"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

func main() {
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	driverManager := plc4go.NewPlcDriverManager()
	transports.RegisterUdpTransport(driverManager)
	drivers.RegisterKnxDriver(driverManager)

	// Try to auto-find KNX gateways via broadcast-message discovery
	driverManager.Discover(func(event model.PlcDiscoveryEvent) {
		connStr := event.ProtocolCode + /*":" + event.TransportCode +*/ "://" + event.TransportUrl.Host
		crc := driverManager.GetConnection(connStr)

		// Wait for the driver to connect (or not)
		connectionResult := <-crc
		if connectionResult.Err != nil {
			log.Errorf("error connecting to PLC: %s", connectionResult.Err.Error())
			return
		}
		connection := connectionResult.Connection
		defer connection.BlockingClose()

		// Try to find all KNX devices on the current network
		browseRequestBuilder := connection.BrowseRequestBuilder()
		//browseRequestBuilder.AddItem("allDevices", "[1-15].[1-15].[0-255]")
		browseRequestBuilder.AddItem("allDevices", "[1-3].[1-6].[0-60]")
		browseRequest, err := browseRequestBuilder.Build()
		if err != nil {
			log.Errorf("error creating browse request: %s", err.Error())
			return
		}
		brr := browseRequest.ExecuteWithInterceptor(func(result model.PlcBrowseEvent) bool {
			knxAddress := result.Result.Address
			log.Info("Inspecting detected Device at KNX Address: " + knxAddress)

			// Try to get all the com-objects and the group addresses they are attached to.
			browseRequestBuilder = connection.BrowseRequestBuilder()
			browseRequestBuilder.AddItem("comObjects", knxAddress+"#com-obj")
			browseRequest, err := browseRequestBuilder.Build()
			if err != nil {
				log.Errorf("error creating read request: %s", err.Error())
				return false
			}
			brr := browseRequest.Execute()
			browseResult := <-brr
			if browseResult.Err != nil {
				log.Errorf("error executing the browse request for com-objects: %s", browseResult.Err.Error())
				return false
			}
			for _, result := range browseResult.Response.GetQueryResults("comObjects") {
				log.Infof(" - %s", result.Address)
			}

			readRequestBuilder := connection.ReadRequestBuilder()
			readRequestBuilder.AddItem("manufacturerId", knxAddress+"#0/12")
			readRequestBuilder.AddItem("programVersion", knxAddress+"#3/13")
			readRequest, err := readRequestBuilder.Build()
			if err != nil {
				log.Error("Error creating read request for scanning " + knxAddress)
				return false
			}

			rrr := readRequest.Execute()
			readRequestResult := <-rrr

			if readRequestResult.Err != nil {
				log.Error("Error executing read request for reading device identification information " + knxAddress)
				return false
			}
			readResponse := readRequestResult.Response
			manufacturerID := "0"
			if readResponse.GetResponseCode("manufacturerId") == model.PlcResponseCode_OK {
				manufacturerID = strconv.FormatUint(uint64(readResponse.GetValue("manufacturerId").GetUint32()), 10)
			}
			var programVersion []byte
			if readResponse.GetResponseCode("programVersion") == model.PlcResponseCode_OK {
				for _, value := range readResponse.GetValue("programVersion").GetList() {
					programVersion = append(programVersion, value.GetUint8())
				}
			}

			log.Infof("     manufacturer id: %s", manufacturerID)
			log.Infof("     application progra version: %s", hex.EncodeToString(programVersion))

			return true
		})
		if brr == nil {
			log.Errorf("error executing browse request")
			return
		}
		select {
		case browseRequestResult := <-brr:
			log.Info(browseRequestResult)
		}
		return
	})

	time.Sleep(time.Second * 1000000)
}
