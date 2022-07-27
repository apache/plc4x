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
	"github.com/apache/plc4x/plc4go/internal/cbus"
	"github.com/apache/plc4x/plc4go/internal/spi/testutils"
	"github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/transports"
	_ "github.com/apache/plc4x/plc4go/tests/initializetest"
	"testing"
)

func TestManualCBusDriver(t *testing.T) {
	t.Skip()

	connectionString := "c-bus://192.168.178.101?srchk=true"
	driverManager := plc4go.NewPlcDriverManager()
	driverManager.RegisterDriver(cbus.NewDriver())
	transports.RegisterTcpTransport(driverManager)
	test := testutils.NewManualTestSuite(connectionString, driverManager, t)

	test.AddTestCase("status/binary/0x04", true)
	// TODO: apparently a level means that we get a extended status reply but at the moment it is guarded by exstat
	test.AddTestCase("status/level=0x40/0x04", true)
	//test.AddTestCase("cal/0/recall=[INTERFACE_OPTIONS_1, 1]", true)
	//test.AddTestCase("cal/0/identify=[FirmwareVersion]", true)
	//test.AddTestCase("cal/0/gestatus=[0xFF, 1]", true)

	test.Run()
}
