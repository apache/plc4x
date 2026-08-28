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

package options

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedactConnectionString(t *testing.T) {
	for _, c := range []struct {
		name  string
		given string
		want  string
	}{
		{"a password parameter", "opcua://host?username=op&password=hunter2",
			"opcua://host?username=op&password=******"},
		{"a leading password parameter", "opcua://host?password=hunter2&username=op",
			"opcua://host?password=******&username=op"},
		{"a pre-shared key", "s7:tls-psk://host?tls-psk.psk-key=0011deadbeef",
			"s7:tls-psk://host?tls-psk.psk-key=******"},
		{"a keystore password", "opcua://host?keyStorePassword=abc",
			"opcua://host?keyStorePassword=******"},
		{"regardless of case", "opcua://host?PassWord=hunter2",
			"opcua://host?PassWord=******"},
		{"credentials in the authority, which have no parameter name",
			"s7://operator:hunter2@plc:102", "s7://operator:******@plc:102"},
		// A colon is legal inside a password. The user segment stops at the first colon so the
		// rest of the credential is the password; a greedy one would publish "operator:hun".
		{"a password containing a colon", "s7://operator:hun:ter2@plc:102",
			"s7://operator:******@plc:102"},
		{"both at once", "opcua://op:hunter2@host?password=abc&read-timeout-ms=5000",
			"opcua://op:******@host?password=******&read-timeout-ms=5000"},
		{"nothing to redact", "modbus-tcp://host:502?unit-identifier=1",
			"modbus-tcp://host:502?unit-identifier=1"},
		{"no parameters at all", "modbus-tcp://host:502", "modbus-tcp://host:502"},
		{"empty", "", ""},
		// The identity names which key was refused - the one thing an operator needs when a PSK
		// handshake fails. Hiding it costs the diagnosis and protects nothing.
		{"the psk identity is not a secret", "s7:tls-psk://host?tls-psk.psk-identity=plc4x",
			"s7:tls-psk://host?tls-psk.psk-identity=plc4x"},
		// Masking these would cost the diagnosis: a path, a store type, a boolean.
		{"things shaped like keys that are not keys",
			"opcua://host?keyStoreFile=/etc/client.p12&securityPolicy=None&discovery=false",
			"opcua://host?keyStoreFile=/etc/client.p12&securityPolicy=None&discovery=false"},
		// Names are matched, not values: a username whose value happens to contain "secret" is
		// still a username, and masking it would cost an operator the account name.
		{"a name that merely looks similar", "opcua://host?username=secretive-bob",
			"opcua://host?username=secretive-bob"},
	} {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, RedactConnectionString(c.given))
		})
	}
}
