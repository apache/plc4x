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
	"github.com/apache/plc4x/plc4go/pkg/api/model"
)

func main() {
	driverManager := plc4go.NewPlcDriverManager()
	drivers.RegisterModbusTcpDriver(driverManager)

	// Get a connection to a remote PLC
	crc := driverManager.GetConnection("modbus-tcp://192.168.23.30")

	// Wait for the driver to connect (or not)
	connectionResult := <-crc
	if connectionResult.GetErr() != nil {
		fmt.Printf("error connecting to PLC: %s", connectionResult.GetErr().Error())
		return
	}
	connection := connectionResult.GetConnection()

	// Make sure the connection is closed at the end
	defer connection.BlockingClose()

	// Prepare a write-request
	writeRequest, err := connection.WriteRequestBuilder().
		AddQuery("field", "holding-register:26:REAL", 2.7182818284).
		Build()
	if err != nil {
		fmt.Printf("error preparing read-request: %s", connectionResult.GetErr().Error())
		return
	}

	// Execute a read-request
	wrc := writeRequest.Execute()

	// Wait for the response to finish
	wrr := <-wrc
	if wrr.GetErr() != nil {
		fmt.Printf("error executing write-request: %s", wrr.GetErr().Error())
		return
	}

	if wrr.GetResponse().GetResponseCode("field") != model.PlcResponseCode_OK {
		fmt.Printf("error an non-ok return code: %s", wrr.GetResponse().GetResponseCode("field").GetName())
		return
	}
	fmt.Print("Result: SUCCESS\n")
}
