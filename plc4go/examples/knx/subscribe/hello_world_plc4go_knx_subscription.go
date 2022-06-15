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
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	"strings"
	"time"
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

	// Prepare a subscription-request
	if subscriptionRequest, err := connection.SubscriptionRequestBuilder().
		// Intentionally catching all without datatype and the temperature values of the first floor with type
		AddChangeOfStateQuery("all", "*/*/*").
		AddChangeOfStateQuery("firstFlorTemperatures", "2/[1,2,4,6]/10:DPT_Value_Temp").
		AddItemHandler(func(event model.PlcSubscriptionEvent) {
			// Iterate over all fields that were triggered in the current event.
			for _, fieldName := range event.GetFieldNames() {
				if event.GetResponseCode(fieldName) == model.PlcResponseCode_OK {
					address := event.GetAddress(fieldName)
					value := event.GetValue(fieldName)
					// If the plc-value was a raw-plcValue, we will try lazily decode the value
					// In my installation all group addresses ending with "/10" are temperature values
					// and ending on "/0" are light switch actions.
					// So if I find a group address ending on that, decode it with a given type name,
					// If not, simply output it as array of USINT values.
					switch value.(type) {
					case values.RawPlcValue:
						rawValue := value.(values.RawPlcValue)
						datatypeName := "USINT"
						if strings.HasSuffix(address, "/10") {
							datatypeName = "DPT_Value_Temp"
						} else if strings.HasSuffix(address, "/0") {
							datatypeName = "BOOL"
						}
						fmt.Printf("Got raw-value event for address %s: ", address)
						if !rawValue.RawHasMore() {
							fmt.Printf("nothing")
						}
						for rawValue.RawHasMore() {
							value = rawValue.RawDecodeValue(datatypeName)
							fmt.Printf(" '%s'", value.GetString())
						}
						fmt.Printf("\n")
					default:
						fmt.Printf("Got event for address %s: %s\n", address, value.GetString())
					}
				}
			}
		}).Build(); err == nil {
		// Execute a subscription-request
		rrc := subscriptionRequest.Execute()

		// Wait for the response to finish
		rrr := <-rrc
		if rrr.GetErr() != nil {
			fmt.Printf("error executing subscription-request: %s", rrr.GetErr().Error())
			return
		}

		// Do something with the response
		for _, fieldName := range rrr.GetResponse().GetFieldNames() {
			if rrr.GetResponse().GetResponseCode(fieldName) != model.PlcResponseCode_OK {
				fmt.Printf("error an non-ok return code for field %s: %s\n", fieldName, rrr.GetResponse().GetResponseCode(fieldName).GetName())
				continue
			}
		}

		time.Sleep(time.Minute * 5)
	} else {
		fmt.Printf("error preparing subscription-request: %s", connectionResult.GetErr().Error())
		return
	}
}
