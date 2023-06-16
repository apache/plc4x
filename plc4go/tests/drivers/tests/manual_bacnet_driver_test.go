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
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/apache/plc4x/plc4go/internal/bacnetip"
	"github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/transports"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

func TestManualBacnetDriver(t *testing.T) {
	t.Skip()

	connectionString := "bacnet-ip://192.168.178.101"
	withCustomLogger := options.WithCustomLogger(testutils.ProduceTestingLogger(t))
	driverManager := plc4go.NewPlcDriverManager(withCustomLogger)
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	driverManager.RegisterDriver(bacnetip.NewDriver(withCustomLogger))
	transports.RegisterUdpTransport(driverManager, withCustomLogger)
	test := testutils.NewManualTestSuite(t, connectionString, driverManager)

	test.AddTestCase("ANALOG_OUTPUT,133/PRESENT_VALUE", true)
	test.AddTestCase("DEVICE,133/LOCATION&DESCRIPTION", true)

	test.Run()
}
