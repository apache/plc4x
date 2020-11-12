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
package read

import (
	"encoding/json"
	"fmt"
	"github.com/apache/plc4x/plc4go/pkg/plc4go"
)

func main() int {
	// Get a connection to a remote PLC
	crc := plc4go.NewPlcDriverManager().GetConnection("modbus:tcp://192.168.23.30")

	// Wait for the driver to connect (or not)
	connectionResult := <-crc
	if connectionResult.Err != nil {
		fmt.Printf("error connecting to PLC: %s", connectionResult.Err.Error())
		return 1
	}
	connection := connectionResult.Connection

	// Make sure the connection is closed at the end
	defer connection.Close()

	// Prepare a read-request
	rrb := connection.ReadRequestBuilder()
	rrb.AddItem("field", "holding-register:1:REAL[2]")
	readRequest, err := rrb.Build()
	if err != nil {
        fmt.Printf("error preparing read-request: %s", connectionResult.Err.Error())
		return 2
	}

	// Execute a read-request
	rrc := readRequest.Execute()

	// Wait for the response to finish
	rrr := <-rrc
	if rrr.Err != nil {
        fmt.Printf("error executing read-request: %s", rrr.Err.Error())
		return 3
	}

	// Do something with the response
	readResponseJson, err := json.Marshal(rrr.Response)
	if err != nil {
        fmt.Printf("error serializing read-response: %s", err.Error())
		return 4
	}
	fmt.Printf("Result: %s\n", string(readResponseJson))

	return 0
}
