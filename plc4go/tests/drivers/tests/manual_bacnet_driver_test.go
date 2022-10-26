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
	"github.com/apache/plc4x/plc4go/internal/bacnetip"
	"github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/transports"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	_ "github.com/apache/plc4x/plc4go/tests/initializetest"
	"testing"
)

func TestManualBacnetDriver(t *testing.T) {
	t.Skip()

	connectionString := "bacnet-ip://192.168.178.101"
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(bacnetip.NewDriver())
	transports.RegisterUdpTransport(driverManager)
	test := testutils.NewManualTestSuite(connectionString, driverManager, t)

	test.AddTestCase("ANALOG_OUTPUT,133/PRESENT_VALUE", true)
	test.AddTestCase("DEVICE,133/LOCATION&DESCRIPTION", true)

	test.Run()
}
