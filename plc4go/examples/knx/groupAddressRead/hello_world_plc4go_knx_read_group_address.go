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
	"fmt"
	"github.com/apache/plc4x/plc4go/pkg/plc4go"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/drivers"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/logging"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
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
	if connectionResult.Err != nil {
		fmt.Printf("error connecting to PLC: %s", connectionResult.Err.Error())
		return
	}
	connection := connectionResult.Connection

	// Make sure the connection is closed at the end
	defer connection.BlockingClose()

	// Prepare a read-request
	readRequest, err := connection.ReadRequestBuilder().
		AddQuery("firstFlorTemperatures", "2/[1,2,4,6]/10:DPT_Value_Temp").
		AddQuery("secondFlorTemperatures", "3/[2,3,4,6]/10:DPT_Value_Temp").
		Build()
	if err != nil {
		fmt.Printf("error preparing read-request: %s", connectionResult.Err.Error())
		return
	}

	// Execute a read-request
	rrc := readRequest.Execute()

	// Wait for the response to finish
	rrr := <-rrc
	if rrr.Err != nil {
		fmt.Printf("error executing read-request: %s", rrr.Err.Error())
		return
	}

	// Do something with the response
	for _, fieldName := range rrr.Response.GetFieldNames() {
		if rrr.Response.GetResponseCode(fieldName) != model.PlcResponseCode_OK {
			fmt.Printf("error an non-ok return code for field %s: %s\n", fieldName, rrr.Response.GetResponseCode(fieldName).GetName())
			continue
		}

		value := rrr.Response.GetValue(fieldName)
		if value == nil {
			fmt.Printf("Got nil for field %s\n", fieldName)
		} else if value.GetStruct() != nil {
			for address, structValue := range value.GetStruct() {
				fmt.Printf("Got result for field %s with address: %s: %s °C\n", fieldName, address, structValue.GetString())
			}
		} else {
			fmt.Printf("Got result for field %s: %s °C\n", fieldName, value.GetString())
		}
	}
}
