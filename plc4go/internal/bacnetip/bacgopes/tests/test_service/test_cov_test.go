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

package test_service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/service"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/time_machine"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

func TestBasic(t *testing.T) {
	t.Skip("needs capability to work") // TODO: fixme
	testingLogger := testutils.ProduceTestingLogger(t)
	ExclusiveGlobalTimeMachine(t)

	// create a network
	anet, err := NewApplicationNetwork(testingLogger)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, anet.Close())
	})

	// add the service capability
	anet.iut.AddCapability(new(ChangeOfValuesServices))

	// all start states are successful
	anet.td.GetStartState().Success("")
	anet.iut.GetStartState().Success("")

	// run the group
	err = anet.Run(0)
	assert.NoError(t, err)
}
