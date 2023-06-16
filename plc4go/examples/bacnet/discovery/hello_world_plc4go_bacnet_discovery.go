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
	"os"
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/drivers"
	"github.com/apache/plc4x/plc4go/pkg/api/logging"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"

	"github.com/rs/zerolog/log"
)

func main() {
	logging.InfoLevel()

	driverManager := plc4go.NewPlcDriverManager()
	defer func() {
		if err := driverManager.Close(); err != nil {
			panic(err)
		}
	}()
	drivers.RegisterBacnetDriver(driverManager)

	var connectionStrings []string
	if len(os.Args) < 2 {
		// Try to auto-find bacnet devices via broadcast-message discovery
		if err := driverManager.Discover(func(event apiModel.PlcDiscoveryItem) {
			connStr := event.GetProtocolCode() + "://" + event.GetTransportUrl().Host
			log.Info().Str("connection string", connStr).Stringer("event", event.(fmt.Stringer)).Msg("Found Bacnet Gateway")

			connectionStrings = append(connectionStrings, connStr)
		},
			plc4go.WithDiscoveryOptionProtocolSpecific("who-is-low-limit", 0),
			plc4go.WithDiscoveryOptionProtocolSpecific("who-is-high-limit", "255"),
			plc4go.WithDiscoveryOptionProtocolSpecific("who-has-device-instance-range-low-limit", "0"),
			plc4go.WithDiscoveryOptionProtocolSpecific("who-has-device-instance-range-high-limit", "255"),
			// Either use name or the 2 identifier tuple below
			plc4go.WithDiscoveryOptionProtocolSpecific("who-has-object-name", "john"),
			//plc4go.WithDiscoveryOptionProtocolSpecific("who-has-object-identifier-type", "DEVICE"),
			//plc4go.WithDiscoveryOptionProtocolSpecific("who-has-object-identifier-instance", "133"),
		); err != nil {
			panic(err)
		}
		// Wait for 5 seconds for incoming responses
		time.Sleep(time.Second * 50)
	} else {
		connStr := "bacnet-ip://" + os.Args[1] + ":47808"
		log.Info().Str("connection string", connStr).Msg("Using manually provided bacnet gateway")
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
		connection.BlockingClose()
	}
}
