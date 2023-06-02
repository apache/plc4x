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
	"github.com/apache/plc4x/plc4go/spi/options"
	"testing"
	"time"

	"github.com/apache/plc4x/plc4go/internal/cbus"
	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/transports"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/stretchr/testify/require"
)

func TestManualCBusDriverMixed(t *testing.T) {
	t.Skip()

	testutils.SetToTestingLogger(t, readWriteModel.Plc4xModelLog)

	withCustomLogger := options.WithCustomLogger(testutils.ProduceTestingLogger(t))

	connectionString := "c-bus://192.168.178.101"
	driverManager := plc4go.NewPlcDriverManager(withCustomLogger)
	driverManager.RegisterDriver(cbus.NewDriver(withCustomLogger))
	transports.RegisterTcpTransport(driverManager, withCustomLogger)
	test := testutils.NewManualTestSuite(t, connectionString, driverManager)

	// TODO: fix those test cases
	//test.AddTestCase("status/binary/0x04", "PlcStruct{\n  application: \"LIGHTING_38\"\n  blockStart: \"false, false, false, false, false, false, false, false\"\n  values: \"DOES_NOT_EXIST, OFF, ERROR, ON, DOES_NOT_EXIST, OFF, ERROR, ON, DOES_NOT_EXIST, OFF, ERROR, ON, DOES_NOT_EXIST, OFF, ERROR, ON, DOES_NOT_EXIST, OFF, ERROR, ON, DOES_NOT_EXIST, OFF, ERROR, ON, DOES_NOT_EXIST, OFF, ERROR, ON, DOES_NOT_EXIST, OFF, ERROR, ON, DOES_NOT_EXIST, OFF, ERROR, ON, DOES_NOT_EXIST, OFF, ERROR, ON, DOES_NOT_EXIST, OFF, ERROR, ON, DOES_NOT_EXIST, OFF, ERROR, ON, DOES_NOT_EXIST, OFF, ERROR, ON, DOES_NOT_EXIST, OFF, ERROR, ON, DOES_NOT_EXIST, OFF, ERROR, ON, DOES_NOT_EXIST, OFF, ERROR, ON, DOES_NOT_EXIST, OFF, ERROR, ON, DOES_NOT_EXIST, OFF, ERROR, ON, DOES_NOT_EXIST, OFF, ERROR, ON, DOES_NOT_EXIST, OFF, ERROR, ON, DOES_NOT_EXIST, OFF, ERROR, ON, DOES_NOT_EXIST, OFF, ERROR, ON\"\n}")
	//test.AddTestCase("status/level=0x40/0x04", 255)
	//test.AddTestCase("cal/0/recall=[INTERFACE_OPTIONS_1, 1]", true)
	//test.AddTestCase("cal/0/identify=[FirmwareVersion]", true)
	//test.AddTestCase("cal/0/gestatus=[0xFF, 1]", true)

	plcConnection := test.Run()
	t.Run("Subscription test", func(t *testing.T) {
		gotMMI := make(chan bool)
		gotSAL := make(chan bool)
		subscriptionRequest, err := plcConnection.SubscriptionRequestBuilder().
			AddEventTagAddress("mmi", "mmimonitor/*/*").
			AddEventTagAddress("sal", "salmonitor/*/*").
			AddPreRegisteredConsumer("mmi", func(event apiModel.PlcSubscriptionEvent) {
				fmt.Printf("mmi:\n%s", event)
				if _, ok := event.GetValue("mmi").GetStruct()["SALData"]; ok {
					panic("got sal in mmi")
				}
				select {
				case gotMMI <- true:
				default:
				}
			}).
			AddPreRegisteredConsumer("sal", func(event apiModel.PlcSubscriptionEvent) {
				fmt.Printf("sal:\n%s", event)
				select {
				case gotSAL <- true:
				default:
				}
			}).
			Build()
		require.NoError(t, err)
		subscriptionRequest.Execute()
		timeout := time.NewTimer(30 * time.Second)
		defer utils.CleanupTimer(timeout)
		// We expect couple monitors
		mmiCount := 0
		salCount := 0
		gotEnough := func() bool {
			return mmiCount > 3 && salCount > 3
		}
	waitingForMonitors:
		for {
			select {
			case at := <-timeout.C:
				t.Errorf("timeout at %s", at)
				break waitingForMonitors
			case <-gotMMI:
				mmiCount++
				fmt.Printf("mmi count: %d\n", mmiCount)
				if gotEnough() {
					break waitingForMonitors
				}
			case <-gotSAL:
				salCount++
				fmt.Printf("sal count: %d\n", salCount)
				if gotEnough() {
					break waitingForMonitors
				}
			}
		}
		t.Logf("Got %d mmis and %d sal monitors", mmiCount, salCount)
	})
}

func TestManualCBusBrowse(t *testing.T) {
	t.Skip()

	testutils.SetToTestingLogger(t, readWriteModel.Plc4xModelLog)

	withCustomLogger := options.WithCustomLogger(testutils.ProduceTestingLogger(t))

	connectionString := "c-bus://192.168.178.101?Monitor=false&MonitoredApplication1=0x00&MonitoredApplication2=0x00"
	driverManager := plc4go.NewPlcDriverManager(withCustomLogger)
	driverManager.RegisterDriver(cbus.NewDriver(withCustomLogger))
	transports.RegisterTcpTransport(driverManager, withCustomLogger)
	connectionResult := <-driverManager.GetConnection(connectionString)
	if err := connectionResult.GetErr(); err != nil {
		t.Error(err)
		t.FailNow()
	}
	connection := connectionResult.GetConnection()
	defer connection.Close()
	browseRequest, err := connection.BrowseRequestBuilder().
		AddQuery("infoQuery", "info/*/*").
		Build()
	if err != nil {
		panic(err)
	}
	browseRequestResult := <-browseRequest.ExecuteWithInterceptor(func(result apiModel.PlcBrowseItem) bool {
		fmt.Printf("%s\n", result)
		return true
	})
	fmt.Printf("%v\n", browseRequestResult.GetResponse())
}

func TestManualCBusRead(t *testing.T) {
	t.Skip()

	testutils.SetToTestingLogger(t, readWriteModel.Plc4xModelLog)

	withCustomLogger := options.WithCustomLogger(testutils.ProduceTestingLogger(t))

	connectionString := "c-bus://192.168.178.101?Monitor=false&MonitoredApplication1=0x00&MonitoredApplication2=0x00"
	driverManager := plc4go.NewPlcDriverManager(withCustomLogger)
	driverManager.RegisterDriver(cbus.NewDriver(withCustomLogger))
	transports.RegisterTcpTransport(driverManager, withCustomLogger)
	connectionResult := <-driverManager.GetConnection(connectionString)
	if err := connectionResult.GetErr(); err != nil {
		t.Error(err)
		t.FailNow()
	}
	connection := connectionResult.GetConnection()
	defer connection.Close()
	readRequest, err := connection.ReadRequestBuilder().
		AddTagAddress("asd", "cal/3/identify=OutputUnitSummary").
		Build()
	require.NoError(t, err)
	readRequestResult := <-readRequest.Execute()
	fmt.Printf("%s", readRequestResult.GetResponse())
}

func TestManualDiscovery(t *testing.T) {
	t.Skip()

	testutils.SetToTestingLogger(t, readWriteModel.Plc4xModelLog)

	withCustomLogger := options.WithCustomLogger(testutils.ProduceTestingLogger(t))

	driverManager := plc4go.NewPlcDriverManager(withCustomLogger)
	driver := cbus.NewDriver(withCustomLogger)
	driverManager.RegisterDriver(driver)
	transports.RegisterTcpTransport(driverManager, withCustomLogger)
	err := driver.Discover(func(event apiModel.PlcDiscoveryItem) {
		t.Log(event.(fmt.Stringer).String())
	})
	require.NoError(t, err)
}
