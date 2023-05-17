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
	"fmt"

	"github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/drivers"
	"github.com/apache/plc4x/plc4go/pkg/api/logging"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
)

func main() {
	// Set logging to INFO
	logging.InfoLevel()

	driverManager := plc4go.NewPlcDriverManager()
	drivers.RegisterKnxDriver(driverManager)

	// Get a connection to a remote PLC
	crc := driverManager.GetConnection("knxnet-ip:udp://192.168.42.11")

	// Wait for the driver to connect (or not)
	connectionResult := <-crc
	if connectionResult.GetErr() != nil {
		fmt.Printf("error connecting to PLC: %s", connectionResult.GetErr().Error())
		return
	}
	connection := connectionResult.GetConnection()

	// Make sure the connection is closed at the end
	defer connection.BlockingClose()

	// Prepare a read-request
	readRequest, err := connection.ReadRequestBuilder().
		AddTagAddress("firstFlorTemperatures", "2/[1,2,4,6]/10:DPT_Value_Temp").
		AddTagAddress("secondFlorTemperatures", "3/[2,3,4,6]/10:DPT_Value_Temp").
		Build()
	if err != nil {
		fmt.Printf("error preparing read-request: %s", connectionResult.GetErr().Error())
		return
	}

	// Execute a read-request
	rrc := readRequest.Execute()

	// Wait for the response to finish
	rrr := <-rrc
	if rrr.GetErr() != nil {
		fmt.Printf("error executing read-request: %s", rrr.GetErr().Error())
		return
	}

	// Do something with the response
	for _, tagName := range rrr.GetResponse().GetTagNames() {
		if rrr.GetResponse().GetResponseCode(tagName) != apiModel.PlcResponseCode_OK {
			fmt.Printf("error an non-ok return code for tag %s: %s\n", tagName, rrr.GetResponse().GetResponseCode(tagName).GetName())
			continue
		}

		value := rrr.GetResponse().GetValue(tagName)
		if value == nil {
			fmt.Printf("Got nil for tag %s\n", tagName)
		} else if value.GetStruct() != nil {
			for address, structValue := range value.GetStruct() {
				fmt.Printf("Got result for tag %s with address: %s: %s °C\n", tagName, address, structValue.GetString())
			}
		} else {
			fmt.Printf("Got result for tag %s: %s °C\n", tagName, value.GetString())
		}
	}
}
