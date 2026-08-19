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

package modbus

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// Sending something that isn't a modbus ADU is a programming error somewhere up the stack, but it
// must surface as an error rather than take the process down with a failed type assertion.
func TestMessageCodec_SendRejectsAForeignMessage(t *testing.T) {
	codec := NewMessageCodec(nil)

	err := codec.Send(testutils.TestContext(t), "test", notAnAdu{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ModbusTcpADU")
}
