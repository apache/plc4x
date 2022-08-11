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

package tests

import (
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/cbus"
	"github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/config"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/transports"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	_ "github.com/apache/plc4x/plc4go/tests/initializetest"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestManualCBusDriver(t *testing.T) {
	log.Logger = log.
		With().Caller().Logger().
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		Level(zerolog.TraceLevel)
	config.TraceTransactionManagerWorkers = true
	config.TraceTransactionManagerTransactions = true
	config.TraceDefaultMessageCodecWorker = true
	t.Skip()

	connectionString := "c-bus://192.168.178.101"
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(cbus.NewDriver())
	transports.RegisterTcpTransport(driverManager)
	test := testutils.NewManualTestSuite(connectionString, driverManager, t)

	test.AddTestCase("status/binary/0x04", "DOES_NOT_EXIST, OFF, ERROR, ON")
	test.AddTestCase("status/level=0x40/0x04", 255)
	//test.AddTestCase("cal/0/recall=[INTERFACE_OPTIONS_1, 1]", true)
	//test.AddTestCase("cal/0/identify=[FirmwareVersion]", true)
	//test.AddTestCase("cal/0/gestatus=[0xFF, 1]", true)

	plcConnection := test.Run()
	t.Run("Subscription test", func(t *testing.T) {
		gotMonitor := make(chan bool)
		subscriptionRequest, err := plcConnection.SubscriptionRequestBuilder().
			AddEventQuery("something", "monitor/*/*").
			AddItemHandler(func(event model.PlcSubscriptionEvent) {
				fmt.Printf("\n%s", event)
				select {
				case gotMonitor <- true:
				default:
				}
			}).
			Build()
		require.NoError(t, err)
		subscriptionRequest.Execute()
		timeout := time.After(30 * time.Second)
		// We expect couple monitors
		monitorCount := 0
	waitingForMonitors:
		for {
			select {
			case at := <-timeout:
				t.Errorf("timeout at %s", at)
				break waitingForMonitors
			case <-gotMonitor:
				monitorCount++
				println(monitorCount)
				if monitorCount > 3 {
					break waitingForMonitors
				}
			}
		}
		t.Logf("Got %d monitors", monitorCount)
	})
}

func TestManualCBusBrowse(t *testing.T) {
	log.Logger = log.
		With().Caller().Logger().
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		Level(zerolog.InfoLevel)
	config.TraceTransactionManagerWorkers = false
	config.TraceTransactionManagerTransactions = false
	config.TraceDefaultMessageCodecWorker = false
	t.Skip()

	connectionString := "c-bus://192.168.178.101?Monitor=false&MonitoredApplication1=0x00&MonitoredApplication2=0x00"
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(cbus.NewDriver())
	transports.RegisterTcpTransport(driverManager)
	connectionResult := <-driverManager.GetConnection(connectionString)
	if err := connectionResult.GetErr(); err != nil {
		t.Error(err)
		t.FailNow()
	}
	connection := connectionResult.GetConnection()
	defer connection.Close()
	browseRequest, err := connection.BrowseRequestBuilder().
		AddQuery("asd", "info/*/*").
		Build()
	if err != nil {
		panic(err)
	}
	browseRequestResult := <-browseRequest.ExecuteWithInterceptor(func(result model.PlcBrowseEvent) bool {
		fmt.Printf("%s", result)
		return true
	})
	fmt.Printf("%s", browseRequestResult.GetResponse())
}
