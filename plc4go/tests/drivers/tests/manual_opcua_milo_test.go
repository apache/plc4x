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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/apache/plc4x/plc4go/internal/opcua"
	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/transports"
	"github.com/apache/plc4x/plc4go/spi/options/converter"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

func TestManualMiloOpcuaRead(t *testing.T) {
	t.Skip()

	optionsForTesting := testutils.EnrichOptionsWithOptionsForTesting(t)

	connectionString := "opcua://milo.digitalpetri.com:62541/milo"
	driverManager := plc4go.NewPlcDriverManager(converter.WithOptionToExternal(optionsForTesting...)...)
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	driverManager.RegisterDriver(opcua.NewDriver(optionsForTesting...))
	transports.RegisterTcpTransport(driverManager, converter.WithOptionToExternal(optionsForTesting...)...)
	test := testutils.NewManualTestSuite(t, connectionString, driverManager)

	test.AddTestCase("ns=2;i=10846;BOOL", 0)

	test.Run()
}
