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
	"strings"
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/drivers"
	"github.com/apache/plc4x/plc4go/pkg/api/logging"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
)

func main() {
	// Set logging to INFO
	logging.InfoLevel()

	driverManager := plc4go.NewPlcDriverManager()
	defer func() {
		if err := driverManager.Close(); err != nil {
			panic(err)
		}
	}()
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
		// Intentionally catching all without datatype and the temperature apiValues of the first floor with type
		AddChangeOfStateTagAddress("all", "*/*/*").
		AddChangeOfStateTagAddress("firstFlorTemperatures", "2/[1,2,4,6]/10:DPT_Value_Temp").
		AddPreRegisteredConsumer("all", func(event apiModel.PlcSubscriptionEvent) {
			// Iterate over all tags that were triggered in the current event.
			for _, tagName := range event.GetTagNames() {
				if event.GetResponseCode(tagName) == apiModel.PlcResponseCode_OK {
					address := event.GetAddress(tagName)
					value := event.GetValue(tagName)
					// If the plc-value was a raw-plcValue, we will try lazily decode the value
					// In my installation all group addresses ending with "/10" are temperature apiValues
					// and ending on "/0" are light switch actions.
					// So if I find a group address ending on that, decode it with a given type name,
					// If not, simply output it as array of USINT apiValues.
					switch value.(type) {
					case apiValues.RawPlcValue:
						rawValue := value.(apiValues.RawPlcValue)
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
		for _, tagName := range rrr.GetResponse().GetTagNames() {
			if rrr.GetResponse().GetResponseCode(tagName) != apiModel.PlcResponseCode_OK {
				fmt.Printf("error an non-ok return code for tag %s: %s\n", tagName, rrr.GetResponse().GetResponseCode(tagName).GetName())
				continue
			}
		}

		time.Sleep(time.Minute * 5)
	} else {
		fmt.Printf("error preparing subscription-request: %s", connectionResult.GetErr().Error())
		return
	}
}
